// Package api

// @author： Boice

// @createTime：2022/5/27 10:03

package api

import (
	"context"
	"github.com/uniplaces/carbon"
	"gitlab.com/bns-engineering/td/common/config"
	"gitlab.com/bns-engineering/td/common/util"
	"gitlab.com/bns-engineering/td/core/engine"
	"testing"
)

func init() {
	util.SetupTelemetry(config.Setup("../../config.json"))
}

func TestStartFlow(t *testing.T) {
	at, _ := carbon.Parse(carbon.DateFormat, "2022-08-31", "")
	engine.Start(context.Background(), "11979664472", at.Time)
}
