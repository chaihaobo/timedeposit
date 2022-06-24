// Package api

// @author： Boice

// @createTime：2022/5/27 10:03

package api

import (
	"context"
	"gitlab.com/bns-engineering/td/common/log"
	"gitlab.com/bns-engineering/td/common/util"
	"testing"

	"gitlab.com/bns-engineering/td/common/config"
	"gitlab.com/bns-engineering/td/core/engine"
)

func init() {
	util.SetupTelemetry(config.Setup("../../config.json"))
}

func TestStartFlow(t *testing.T) {
	if config.TDConf.SkipTests {
		return
	}
	tmpTDAccountList, err := loadAccountList(context.Background())
	if err != nil {
		log.Error(context.Background(), "load mambu account list error", err)
	}

	if len(tmpTDAccountList) > 0 {
		engine.Start(context.Background(), tmpTDAccountList[0].ID)
	}
}
