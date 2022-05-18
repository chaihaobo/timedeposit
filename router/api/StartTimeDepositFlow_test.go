/*
 * @Author: Hugo
 * @Date: 2022-05-16 09:08:29
 * @Last Modified by: Hugo
 * @Last Modified time: 2022-05-16 10:38:40
 */
package api

import (
	"testing"

	commonConfig "gitlab.com/hugo.hu/time-deposit-eod-engine/common/config"
	commonLog "gitlab.com/hugo.hu/time-deposit-eod-engine/common/log"
	"gitlab.com/hugo.hu/time-deposit-eod-engine/flow"
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
	config, _ := commonConfig.NewConfig("./../../config.json")
	commonLog.InitLogConfig(config)
	flow.InitWorkflow()
}
