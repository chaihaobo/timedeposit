package logger

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"strconv"
	"strings"
	"time"

	"github.com/natefinch/lumberjack"
	telemetry "gitlab.com/bns-engineering/common/telemetry"
	"gitlab.com/bns-engineering/common/telemetry/instrumentation/filter"
	"gitlab.com/bns-engineering/td/common/config"
	commonlog "gitlab.com/bns-engineering/td/common/log"

	"github.com/gin-gonic/gin"
	ins "gitlab.com/bns-engineering/common/telemetry/instrumentation"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var lg *zap.Logger
var SERVICENAME = "time_deposit" // please help configure

// InitLogger init logger
func SetUp(cfg *config.TDConfig) (err error) {
	writeSyncer := getLogWriter(cfg.Log.Filename, cfg.Log.Maxsize, cfg.Log.MaxBackups, cfg.Log.MaxAge)
	encoder := getEncoder()
	var l = new(zapcore.Level)
	err = l.UnmarshalText([]byte(cfg.Log.Level))
	if err != nil {
		return
	}
	core := zapcore.NewCore(encoder, writeSyncer, l)

	lg = zap.New(core, zap.AddCaller())
	zap.ReplaceGlobals(lg) // replace package zap global logger instance
	return
}

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.TimeKey = "time"
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.EncodeDuration = zapcore.SecondsDurationEncoder
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	return zapcore.NewJSONEncoder(encoderConfig)
}

func getLogWriter(filename string, maxSize, maxBackup, maxAge int) zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   filename,
		MaxSize:    maxSize,
		MaxBackups: maxBackup,
		MaxAge:     maxAge,
	}
	return zapcore.AddSync(io.MultiWriter(lumberJackLogger, os.Stdout))
}

// GinLogger use gin logger
func GinLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		c.Next()

		cost := time.Since(start)
		lg.Info(path,
			zap.Int("status", c.Writer.Status()),
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.String("query", query),
			zap.String("ip", c.ClientIP()),
			zap.String("user-agent", c.Request.UserAgent()),
			zap.String("errors", c.Errors.ByType(gin.ErrorTypePrivate).String()),
			zap.Duration("cost", cost),
		)
	}
}
func GinLoggerResponse(telemetry *telemetry.API) gin.HandlerFunc {
	return func(c *gin.Context) {

		defer func(start time.Time) {
			elapsedTime := time.Since(start).Milliseconds()

			// send metrics to datadog
			method := strings.ToLower(c.Request.Method)
			name := SERVICENAME + fmt.Sprintf("%s_%s", method, c.Request.URL.Path)
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

			var bytes = []byte{} // response writer in byte
			filterConfig := telemetry.Filter
			makemapresp = string(bytes)
			if string(bytes) != "" && filterConfig != nil {
				rules := filterConfig.PayloadFilter(&filter.TargetFilter{
					Method: c.Request.URL.Path,
				})

				makemapresp = filter.BodyFilter(rules, makemapresp)
			}

			commonlog.Info(context.Background(), "Http Response",
				[]commonlog.Field{
					{Field: zap.String(ins.LabelHTTPService, c.Request.URL.Path)},
					{Field: zap.Any(ins.LabelHTTPResponse, makemapresp)},
					{Field: zap.Any(ins.LabelHTTPStatus, c.Writer.Status())},
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

			commonlog.Info(context.Background(), "Http Request",
				[]commonlog.Field{
					{Field: zap.String(ins.LabelHTTPService, c.Request.URL.Path)},
					{Field: zap.Any(ins.LabelHTTPHeader, filteredHeader)},
					{Field: zap.Any(ins.LabelHTTPRequest, filteredBody)},
					{Field: zap.Any(ins.LabelHTTPMethod, c.Writer.Status())},
				}...,
			)

		} else {
			headerKeyFilter := filter.DefaultHeaderFilter
			if telemetry.Filter != nil {
				headerKeyFilter = append(headerKeyFilter, telemetry.Filter.HeaderFilter...)
			}

			filteredHeader := filter.HeaderFilter(header, headerKeyFilter)

			commonlog.Info(context.Background(), "Http Request",
				[]commonlog.Field{
					{Field: zap.String(ins.LabelHTTPService, c.Request.URL.Path)},
					{Field: zap.Any(ins.LabelHTTPHeader, filteredHeader)},
					{Field: zap.Any(ins.LabelHTTPMethod, c.Writer.Status())},
				}...,
			)
		}
	}
}

// GinRecovery recover panic
func GinRecovery(stack bool) gin.HandlerFunc {
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
					lg.Error(c.Request.URL.Path,
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
					// If the connection is dead, we can't write a status to it.
					c.Error(err.(error)) // nolint: errcheck
					c.Abort()
					return
				}

				if stack {
					lg.Error("[Recovery from panic]",
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
						zap.String("stack", string(debug.Stack())),
					)
				} else {
					lg.Error("[Recovery from panic]",
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
				}
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		c.Next()
	}
}
