/*
 * @Author: Hugo
 * @Date: 2022-04-29 10:23:50
 * @Last Modified by: Hugo
 * @Last Modified time: 2022-05-16 10:34:03
 */
package main

import (
	"context"
	"fmt"
	"gitlab.com/bns-engineering/common/telemetry"
	"gitlab.com/bns-engineering/common/telemetry/instrumentation/filter"
	commonlog "gitlab.com/bns-engineering/td/common/log"
	"gitlab.com/bns-engineering/td/common/logger"
	"log"

	"net/http"
	"os"
	"os/signal"
	"regexp"
	"runtime"
	"syscall"
	"time"

	"gitlab.com/bns-engineering/td/common/util"
	"gitlab.com/bns-engineering/td/core/engine"
	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
	"gitlab.com/bns-engineering/td/common/config"
	"gitlab.com/bns-engineering/td/router"
)

// Initial configuration for this app
func init() {
	zone := time.FixedZone("CST", 7*3600)
	time.Local = zone
	err := logger.SetUp(config.Setup("./config.json"))
	if err != nil {
		zap.L().Error("logger init error", zap.Error(err))
	}
}

func main() {

	cores := runtime.NumCPU()
	runtime.GOMAXPROCS(cores)
	gin.SetMode(config.TDConf.Server.RunMode)
	routersInit := router.InitRouter()
	endPoint := fmt.Sprintf(":%d", config.TDConf.Server.HttpPort)

	server := &http.Server{
		Addr:    endPoint,
		Handler: routersInit,
	}


	log.SetFlags(log.LstdFlags | log.Lshortfile)
	// telemetry
	ins, close := getTelemetry()
	defer close()

	ins.Filter = new(telemetry.FilterConfig)
	initLogBodyFilter([]string{"password","nik","motherName"}, ins)
	initLogHeaderFilter([]string{"authorization,Authorization,deviceid"}, ins)

	util.SetTelemetryLog(ins)
	util.SetTelemetryDataDog(ins)
	commonlog.NewLogger(ins)
	util.SetTelemetryFilter(ins)
	util.GetTelemetryDataDogOpt()

	go func() {
		zap.L().Info("start http server listening ", zap.String("endPoint", endPoint))
		// service connections

		err := server.ListenAndServe()
		if err != http.ErrServerClosed {
			util.CheckAndExit(err)
		}
	}()

	// graceful shutdown
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	<-sigs
	// stop gin engine
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	util.CheckAndExit(server.Shutdown(ctx))
	engine.Pool.Release()

}

func getTelemetry() (*telemetry.API, func()) {
	// TODO: please Help to configure
	collectorURL := "http://localhost:14268/api/traces"
	serviceName := "timedeposit"
	sourceEnv := "local"
	metricPort := 8181
	agentAddress := "127.0.0.1:8125"
	config := telemetry.APIConfig{
		LoggerConfig: telemetry.LoggerConfig{},
		TraceConfig:  telemetry.TraceConfig{CollectorEndpoint: collectorURL, ServiceName: serviceName, SourceEnv: sourceEnv},
		MetricConfig: telemetry.MetricConfig{
			Port:         int(metricPort),
			AgentAddress: agentAddress,
			SampleRate:   1,
		},
	}
	ins, fn, _ := telemetry.NewInstrumentation(config)

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