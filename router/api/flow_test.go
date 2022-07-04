// Package api

// @author： Boice

// @createTime：2022/5/27 10:03

package api

import (
	"gitlab.com/bns-engineering/td/common/config"
	"gitlab.com/bns-engineering/td/common/util"
	"testing"
)

func init() {
	util.SetupTelemetry(config.Setup("../../config.json"))
}

func TestStartFlow(t *testing.T) {
	// engine.Start(context.Background(), "11979664472")
}
