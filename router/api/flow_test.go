// Package api

// @author： Boice

// @createTime：2022/5/27 10:03

package api

import (
	"context"
	"gitlab.com/bns-engineering/td/common/config"
	"gitlab.com/bns-engineering/td/common/util"
	"gitlab.com/bns-engineering/td/core/engine"
	"testing"
	"time"
)

func init() {
	util.SetupTelemetry(config.Setup("../../config.json"))
	zone := time.FixedZone("CST", 7*3600)
	time.Local = zone
}

func TestStartFlow(t *testing.T) {
	engine.Start(context.Background(), "11961471326")
}
