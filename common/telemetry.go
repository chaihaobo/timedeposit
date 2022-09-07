// Package common
// @author： Boice
// @createTime：2022/7/22 15:53
package common

import (
	"fmt"
	"gitlab.com/bns-engineering/common/telemetry"
	"gitlab.com/bns-engineering/common/telemetry/instrumentation/filter"
	"regexp"
)

type Telemetry struct {
	API    *telemetry.API
	Closer func()
}

func newTelemetry(config *Config, credential *Credential) *Telemetry {
	telemetryConfig := telemetry.APIConfig{
		LoggerConfig: telemetry.LoggerConfig{
			FileName: config.Log.Filename,
			MaxSize:  config.Log.Maxsize,
			MaxAge:   config.Log.MaxAge,
		},
		TraceConfig: telemetry.TraceConfig{CollectorEndpoint: credential.Trace.CollectorURL, ServiceName: credential.Trace.ServiceName, SourceEnv: credential.Trace.SourceEnv},
		MetricConfig: telemetry.MetricConfig{
			Port:         credential.Metric.Port,
			AgentAddress: credential.Metric.AgentAddress,
			SampleRate:   1,
		},
	}
	ins, fn, _ := telemetry.NewInstrumentation(telemetryConfig)

	ins.Filter = new(telemetry.FilterConfig)
	initLogBodyFilter([]string{"password", "nik", "motherName"}, ins)
	initLogHeaderFilter([]string{"authorization,Authorization,deviceid"}, ins)
	return &Telemetry{
		ins,
		fn,
	}
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
