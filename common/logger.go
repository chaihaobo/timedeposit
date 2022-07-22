// Package common
// @author： Boice
// @createTime：2022/7/22 16:10
package common

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"gitlab.com/bns-engineering/common/telemetry"
	ins "gitlab.com/bns-engineering/common/telemetry/instrumentation"
	"gitlab.com/bns-engineering/common/telemetry/instrumentation/filter"
	"go.uber.org/zap"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"strconv"
	"strings"
	"time"
)

type Logger telemetry.Logger

func NewLogger(telemetryAPI *telemetry.API) Logger {
	return telemetryAPI.Logger()
}

func GinLoggerResponse(telemetry *telemetry.API) gin.HandlerFunc {

	return func(c *gin.Context) {
		c.Next()

		defer func(start time.Time) {
			elapsedTime := time.Since(start).Milliseconds()

			// send metrics to datadog
			method := strings.ToLower(c.Request.Method)
			name := telemetry.ServiceAPI + fmt.Sprintf("%s_%s", method, c.Request.URL.Path)
			if telemetry.ServiceAPI != "" {
				name = telemetry.ServiceAPI
			}
			src_env := telemetry.SourceEnv

			urlPath := c.Request.URL.Path

			endpoint := fmt.Sprintf("%s.%s", urlPath, method)

			// send metrics to datadog
			tags := []string{
				"http_endpoint:" + endpoint,
				"src_env:" + src_env,
				"response_code:" + strconv.FormatUint(uint64(c.Writer.Status()), 10),
			}

			telemetry.Metric().Count(name, 1, tags)
			telemetry.Metric().Histogram(name+".histogram", float64(elapsedTime), tags)
			telemetry.Metric().Distribution(name+".distribution", float64(elapsedTime), tags)
		}(time.Now())

		defer func() {
			var makemapresp interface{}
			var responseBytes []byte // response writer in byte
			filterConfig := telemetry.Filter
			makemapresp = string(responseBytes)

			if string(responseBytes) != "" && filterConfig != nil {
				rules := filterConfig.PayloadFilter(&filter.TargetFilter{
					Method: c.Request.URL.Path,
				})

				makemapresp = filter.BodyFilter(rules, makemapresp)
			}

			telemetry.Logger().Info(c.Request.Context(), "Http Response",
				[]zap.Field{
					zap.String(ins.LabelHTTPService, c.Request.URL.Path),
					zap.Any(ins.LabelHTTPResponse, makemapresp),
					zap.Any(ins.LabelHTTPStatus, c.Writer.Status()),
				}...,
			)

		}()

	}
}

func GinLoggerRequest(telemetry *telemetry.API) gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path
		header := c.Request.Header
		if c.Request.Method == http.MethodPost || c.Request.Method == http.MethodPatch || c.Request.Method == http.MethodPut {
			var requestBody string
			rawBody, err := ioutil.ReadAll(c.Request.Body)
			if err == nil {
				c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(rawBody))
				requestBody = string(rawBody)
			}

			var filteredBody interface{}
			headerKeyFilter := filter.DefaultHeaderFilter
			// set request and response default
			filteredBody = requestBody
			if telemetry.Filter != nil {
				rules := telemetry.Filter.PayloadFilter(&filter.TargetFilter{
					Method: path,
				})

				filteredBody = filter.BodyFilter(rules, requestBody)
				headerKeyFilter = append(headerKeyFilter, telemetry.Filter.HeaderFilter...)
			}

			filteredHeader := filter.HeaderFilter(header, headerKeyFilter)

			telemetry.Logger().Info(c.Request.Context(), "Http Request",
				[]zap.Field{
					zap.String(ins.LabelHTTPService, c.Request.URL.Path),
					zap.Any(ins.LabelHTTPHeader, filteredHeader),
					zap.Any(ins.LabelHTTPRequest, filteredBody),
					zap.Any(ins.LabelHTTPMethod, c.Writer.Status()),
				}...,
			)

		} else {
			headerKeyFilter := filter.DefaultHeaderFilter
			if telemetry.Filter != nil {
				headerKeyFilter = append(headerKeyFilter, telemetry.Filter.HeaderFilter...)
			}

			filteredHeader := filter.HeaderFilter(header, headerKeyFilter)

			telemetry.Logger().Info(c.Request.Context(), "Http Request",
				[]zap.Field{
					zap.String(ins.LabelHTTPService, c.Request.URL.Path),
					zap.Any(ins.LabelHTTPHeader, filteredHeader),
					zap.Any(ins.LabelHTTPMethod, c.Writer.Status()),
				}...,
			)
		}
		c.Next()
	}

}

func GinRecovery(telemetry *telemetry.API, stack bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Check for a broken connection, as it is not really a
				// condition that warrants a panic stack trace.
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}

				httpRequest, _ := httputil.DumpRequest(c.Request, false)
				if brokenPipe {
					telemetry.Logger().Error(c.Request.Context(), c.Request.URL.Path,
						errors.New(fmt.Sprintf("%v", err)),
						zap.String("request", string(httpRequest)),
					)
					// If the connection is dead, we can't write a status to it.
					c.Error(err.(error)) // nolint: errcheck
					c.Abort()
					return
				}

				if stack {
					telemetry.Logger().Error(c.Request.Context(), "[Recovery from panic]",
						errors.New(fmt.Sprintf("%v", err)),
						zap.String("request", string(httpRequest)),
						zap.String("stack", string(debug.Stack())),
					)
				} else {
					telemetry.Logger().Error(c.Request.Context(), "[Recovery from panic]",
						errors.New(fmt.Sprintf("%v", err)),
						zap.String("request", string(httpRequest)),
					)
				}
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		c.Next()
	}
}
