// Package api

// @author： Boice

// @createTime：2022/5/27 10:03

package api

import (
	"testing"

	"gitlab.com/bns-engineering/td/common/config"
	"gitlab.com/bns-engineering/td/common/logger"
	"gitlab.com/bns-engineering/td/core/engine"
	"go.uber.org/zap"
)

func init() {
	logger.SetUp(config.Setup("../../config.json"))
}

func TestStartFlow(t *testing.T) {
	if config.TDConf.SkipTests {
		return
	}
	tmpTDAccountList, err := loadAccountList()
	if err != nil {
		zap.L().Error("load mambu account list error")
	}

	if len(tmpTDAccountList) > 0 {
		engine.Start(tmpTDAccountList[0].ID)
	}
}
