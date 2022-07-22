// Package common
// @author： Boice
// @createTime：2022/7/22 15:57
package common

import (
	"gitlab.com/bns-engineering/common/telemetry"
	"time"
)

type DataDog struct {
	Metric      telemetry.Metrics
	ServiceName string
	SourceEnv   string
}

type DataDogMetric struct {
	metric     telemetry.Metrics
	tags       []string
	startTime  time.Time
	metricName string
}

func NewDataDog(telemetryAPI *telemetry.API) *DataDog {
	return &DataDog{
		Metric:      telemetryAPI.Metric(),
		ServiceName: telemetryAPI.ServiceAPI,
		SourceEnv:   telemetryAPI.SourceEnv,
	}
}

func (d DataDog) StartMetric(metricName string) DataDogMetric {
	metricName = d.ServiceName + "." + metricName
	tags := []string{"src_env:" + d.SourceEnv, "name:" + metricName}
	startTime := time.Now()
	return DataDogMetric{
		metric:     d.Metric,
		tags:       tags,
		startTime:  startTime,
		metricName: metricName,
	}
}

func (m *DataDogMetric) appendTags(input []string) {
	m.tags = append(m.tags, input...)
}

func (m *DataDogMetric) Tags() []string {
	return m.tags
}
func (m *DataDogMetric) PushMetric(tags []string) {
	m.appendTags(tags)
	m.metric.Count(m.metricName, 1, m.tags)
	elapsedTime := time.Since(m.startTime).Milliseconds()
	m.metric.Histogram(m.metricName+".histogram", float64(elapsedTime), m.tags)
}
