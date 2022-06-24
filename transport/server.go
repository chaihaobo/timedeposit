// Package transport
// @author： Boice
// @createTime：2022/6/24 09:49
package transport

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"gitlab.com/bns-engineering/common/telemetry"
	"gitlab.com/bns-engineering/common/telemetry/instrumentation/filter"
	"gitlab.com/bns-engineering/td/common/config"
	commonlog "gitlab.com/bns-engineering/td/common/log"
	"gitlab.com/bns-engineering/td/common/util"
	"gitlab.com/bns-engineering/td/router"
	"go.uber.org/zap"
	"net/http"
	"regexp"
	"runtime"
	"time"
)

type Server interface {
	Start()
	SetUp()
	Shutdown()
}

type tdServer struct {
	telemetryApi    *telemetry.API
	telemetryCloser func()
	conf            *config.TDConfig
	server          *http.Server
}

func (s *tdServer) SetUp() {
	s.telemetryApi.Filter = new(telemetry.FilterConfig)
	initLogBodyFilter([]string{"password", "nik", "motherName"}, s.telemetryApi)
	initLogHeaderFilter([]string{"authorization,Authorization,deviceid"}, s.telemetryApi)
	util.SetTelemetryDataDog(s.telemetryApi)
	util.SetTelemetry(s.telemetryApi)
	commonlog.NewLogger(s.telemetryApi)

}

func (s *tdServer) Start() {
	cores := runtime.NumCPU()
	runtime.GOMAXPROCS(cores)
	gin.SetMode(config.TDConf.Server.RunMode)
	routersInit := router.InitRouter(s.telemetryApi)
	endPoint := fmt.Sprintf(":%d", config.TDConf.Server.HttpPort)

	s.server = &http.Server{
		Addr:    endPoint,
		Handler: routersInit,
	}
	go func() {
		defer s.telemetryCloser()
		commonlog.Info(context.Background(), "start http server listening", zap.String("endPoint", endPoint))
		err := s.server.ListenAndServe()
		if err != http.ErrServerClosed {
			util.CheckAndExit(err)
		}
	}()
}

func (s *tdServer) Shutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	util.CheckAndExit(s.server.Shutdown(ctx))
}

func NewTdServer(config *config.TDConfig) Server {
	ins, closer := getTelemetry()
	return &tdServer{
		telemetryApi:    ins,
		telemetryCloser: closer,
		conf:            config,
	}

}

func getTelemetry() (*telemetry.API, func()) {
	telemetryConfig := telemetry.APIConfig{
		LoggerConfig: telemetry.LoggerConfig{},
		TraceConfig:  telemetry.TraceConfig{CollectorEndpoint: config.TDConf.Trace.CollectorURL, ServiceName: config.TDConf.Trace.ServiceName, SourceEnv: config.TDConf.Trace.SourceEnv},
		MetricConfig: telemetry.MetricConfig{
			Port:         config.TDConf.Metric.Port,
			AgentAddress: config.TDConf.Metric.AgentAddress,
			SampleRate:   1,
		},
	}
	ins, fn, _ := telemetry.NewInstrumentation(telemetryConfig)
	return ins, fn
}

func initLogBodyFilter(configString []string, client *telemetry.API) {
	// overide filter value
	client.Filter.PayloadFilter = func(item *filter.TargetFilter) []*regexp.Regexp {
		var rules []*regexp.Regexp
		for _, v := range configString {
			pattern := fmt.Sprintf(`(%s|\"%s\"\s*):\s?"([\w\s#-@]+)`, v, v)
			regex := regexp.MustCompile(pattern)
			rules = append(rules, regex)
		}
		return rules
	}
}

func initLogHeaderFilter(arrayStringConfig []string, client *telemetry.API) {
	// overide filter value
	client.Filter.HeaderFilter = arrayStringConfig
}
