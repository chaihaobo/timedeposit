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
	"go.uber.org/zap"
	"testing"

	"gitlab.com/bns-engineering/td/flow"
)

func TestStartTDFlow(t *testing.T) {
	tests := []struct {
		name string
	}{
		// TODO: Add test cases.
		{
			name: "test start flow",
		},
	}
	initConfig()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			StartTDFlow(nil)
		})
	}
}

func initConfig() {
	conf := commonConfig.Setup("./../config.json")
	logger.SetUp(conf)
	zap.L().Info("===============Start Test Whole Flow==============")
	flow.InitWorkflow()
}
