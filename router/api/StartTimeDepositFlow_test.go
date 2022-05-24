/*
 * @Author: Hugo
 * @Date: 2022-05-16 09:08:29
 * @Last Modified by: Hugo
 * @Last Modified time: 2022-05-23 11:25:01
 */
package api

import (
	commonConfig "gitlab.com/bns-engineering/td/common/config"
	logger "gitlab.com/bns-engineering/td/common/log"
	"gitlab.com/bns-engineering/td/flow"
	"gitlab.com/bns-engineering/td/service/mambuEntity"
	"go.uber.org/zap"
	"testing"
)

func init() {
	conf := commonConfig.Setup("../../config.yaml")
	logger.SetUp(conf)
	zap.L().Info("===============Start Test Whole Flow==============")
	flow.InitWorkflow()
}

func TestRunFlow(t *testing.T) {
	t.Run("run flow by account", func(t *testing.T) {
		RunFlow(&mambuEntity.TDAccount{
			ID: "11747126703",
		})

		t.Log("run flow by account success")
	})

}
