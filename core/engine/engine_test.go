// Package engine
// @author： Boice
// @createTime：2022/5/26 13:59
package engine

import (
	"encoding/json"
	"gitlab.com/bns-engineering/td/common/config"
	logger "gitlab.com/bns-engineering/td/common/log"
	"gitlab.com/bns-engineering/td/core/engine/node/constant"
	"gitlab.com/bns-engineering/td/repository"
	"gitlab.com/bns-engineering/td/service/mambuEntity"
	"go.uber.org/zap"
	"testing"
)

func init() {
	logger.SetUp(config.Setup("../../config.yaml"))
}

func TestEngine(t *testing.T) {

	t.Run("test engine run", func(t *testing.T) {
		Start("11249460359")

	})
}

func TestGetTDAccount(t *testing.T) {

	t.Run("test engine run", func(t *testing.T) {
		log := repository.GetFlowNodeQueryLogRepository().GetNewLog("20220527151700_11249460359", constant.QueryTDAccount)
		if log != nil {
			saveDBAccount := new(mambuEntity.TDAccount)
			data := log.Data
			err := json.Unmarshal([]byte(data), saveDBAccount)
			if err != nil {
				zap.L().Error("account from db can not map to struct")
			}
		}

	})
}
