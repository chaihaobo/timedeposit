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
	logger.SetUp(config.Setup("../../config.yaml"))
}

func TestStartFlow(t *testing.T) {
	tmpTDAccountList, err := loadAccountList()
	if err != nil {
		zap.L().Error("load mambu account list error")
	}

	for _, tmpTDAcc := range tmpTDAccountList {
		engine.Start(tmpTDAcc.ID)
		// go engine.Start(tmpTDAcc.ID)
		zap.L().Info("commit task success!", zap.String("account", tmpTDAcc.ID))
	}
}
