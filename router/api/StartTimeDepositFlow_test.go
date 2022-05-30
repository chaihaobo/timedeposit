/*
 * @Author: Hugo
 * @Date: 2022-05-16 09:08:29
 * @Last Modified by: Hugo
 * @Last Modified time: 2022-05-23 11:25:01
 */
package api

import (
	"fmt"
	"gitlab.com/bns-engineering/td/core/engine"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	commonConfig "gitlab.com/bns-engineering/td/common/config"
	logger "gitlab.com/bns-engineering/td/common/log"
	"gitlab.com/bns-engineering/td/flow"
	"go.uber.org/zap"
)

func init() {
	conf := commonConfig.Setup("../../config.yaml")
	logger.SetUp(conf)
	zap.L().Info("===============Start Test Whole Flow==============")
	flow.InitWorkflow()
}

func TestStartTDFlow(t *testing.T) {
	type args struct {
		c *gin.Context
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			StartTDFlow(tt.args.c)
		})
	}
}

func TestRunFlow(t *testing.T) {
	for i := 0; i < 100; i++ {
		_ = engine.Pool.Invoke(fmt.Sprintf("%d\n", i))

	}

	time.Sleep(time.Minute * 10)
}
