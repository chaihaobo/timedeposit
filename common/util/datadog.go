package util

import (
	"sync"
	"time"

	"gitlab.com/bns-engineering/common/telemetry"
)

var ddog telemetry.Metrics
var ddogOpt *DataDogOptions
var mu sync.Mutex

type DataDogOptions struct {
	ServiceName string
	SourceEnv   string
}

func SetTelemetryDataDog(telemetryAPI *telemetry.API) *DataDogOptions {
	mu.Lock()
	defer mu.Unlock()
	ddog = telemetryAPI.Metric()
	ddogOpt = &DataDogOptions{
		ServiceName: telemetryAPI.ServiceAPI,
		SourceEnv:   telemetryAPI.SourceEnv,
	}
	return ddogOpt
}

func GetTelemetryDataDog() telemetry.Metrics {
	return ddog
}
func GetTelemetryDataDogOpt() *DataDogOptions {
	return ddogOpt
}

// for unit test initiation
func InitDDogTest() {
	config := telemetry.APIConfig{
		TraceConfig: telemetry.TraceConfig{CollectorEndpoint: "127.0.0.1/api/"},
		MetricConfig: telemetry.MetricConfig{
			AgentAddress: "127.0.0.1:8125",
		},
	}
	ins, _, _ := telemetry.NewInstrumentation(config)
	SetTelemetryDataDog(ins)

}

type datadogMetric struct {
	tags       []string
	startTime  time.Time
	metricName string
}

func StartMetric(metricName string) datadogMetric {
	metricName = GetTelemetryDataDogOpt().ServiceName + "." + metricName
	tags := []string{"src_env:" + GetTelemetryDataDogOpt().SourceEnv, "name:" + metricName}
	startTime := time.Now()

	return datadogMetric{
		tags:       tags,
		startTime:  startTime,
		metricName: metricName,
	}
}

func (m *datadogMetric) appendTags(input []string) {
	m.tags = append(m.tags, input...)
}

func (m *datadogMetric) Tags() []string {
	return m.tags
}
func (m *datadogMetric) PushMetric(tags []string) {
	m.appendTags(tags)
	GetTelemetryDataDog().Count(m.metricName, 1, m.tags)
	elapsedTime := time.Since(m.startTime).Milliseconds()
	GetTelemetryDataDog().Histogram(m.metricName+".histogram", float64(elapsedTime), m.tags)
}
