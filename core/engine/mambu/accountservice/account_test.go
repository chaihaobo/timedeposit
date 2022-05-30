// Package accountservice
// @author： Boice
// @createTime：2022/5/30 10:18
package accountservice

import (
	"context"
	"gitlab.com/bns-engineering/td/common/config"
	logger "gitlab.com/bns-engineering/td/common/log"
	"testing"
)

func init() {
	logger.SetUp(config.Setup("../../../../config.yaml"))
}

func TestGetAccountById(t *testing.T) {
	_, err := GetAccountById(nil, "11249460359")
	if err != nil {
		t.Errorf("test error")
	}
}

func TestUndoMaturityDate(t *testing.T) {
	UndoMaturityDate(nil, "11249460359")
}

func TestApplyProfit(t *testing.T) {
	ApplyProfit(context.WithValue(context.WithValue(context.Background(), "flowId", "test"), "nodeName", "test Node"), "11714744288", "ok")
}
