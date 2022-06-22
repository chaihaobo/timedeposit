package util

import (

	"gitlab.com/bns-engineering/common/telemetry"
)

var telemetryLog telemetry.Logger

var filter *telemetry.FilterConfig

func SetTelemetryLog(telemetryAPI *telemetry.API) {
	telemetryLog = telemetryAPI.Logger()
}

// SetTelemetryMockLog for testing purposes
func SetTelemetryMockLog(mock telemetry.Logger) {
	telemetryLog = mock
}

func GetTelemetryLog() telemetry.Logger {
	return telemetryLog
}

func SetTelemetryFilter(telemetryAPI *telemetry.API) {
	filter = telemetryAPI.Filter
}

func GetTelemetryFilter() *telemetry.FilterConfig {
	return filter
}
