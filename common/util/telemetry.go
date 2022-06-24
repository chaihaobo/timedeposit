package util

import (
	"gitlab.com/bns-engineering/common/telemetry"
)

var telemetryApi *telemetry.API

func SetTelemetry(telemetryAPI *telemetry.API) {
	telemetryApi = telemetryAPI
}

func GetTelemetry() *telemetry.API {
	return telemetryApi
}
