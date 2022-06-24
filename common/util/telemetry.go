package util

import (
	"fmt"
	"gitlab.com/bns-engineering/common/telemetry"
	"gitlab.com/bns-engineering/common/telemetry/instrumentation/filter"
	"gitlab.com/bns-engineering/td/common/config"
	commonlog "gitlab.com/bns-engineering/td/common/log"
	"regexp"
)

var telemetryApi *telemetry.API
var telemetryApiCloser func()

func GetTelemetry() *telemetry.API {
	return telemetryApi
}

func GetTelemetryCloser() func() {
	return telemetryApiCloser
}

func SetupTelemetry(config *config.TDConfig) (*telemetry.API, func()) {
	telemetryConfig := telemetry.APIConfig{
		LoggerConfig: telemetry.LoggerConfig{},
		TraceConfig:  telemetry.TraceConfig{CollectorEndpoint: config.Trace.CollectorURL, ServiceName: config.Trace.ServiceName, SourceEnv: config.Trace.SourceEnv},
		MetricConfig: telemetry.MetricConfig{
			Port:         config.Metric.Port,
			AgentAddress: config.Metric.AgentAddress,
			SampleRate:   1,
		},
	}
	ins, fn, _ := telemetry.NewInstrumentation(telemetryConfig)
	telemetryApi = ins
	telemetryApiCloser = fn

	ins.Filter = new(telemetry.FilterConfig)
	initLogBodyFilter([]string{"password", "nik", "motherName"}, ins)
	initLogHeaderFilter([]string{"authorization,Authorization,deviceid"}, ins)
	SetTelemetryDataDog(ins)
	commonlog.NewLogger(ins)

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
