// Package accountservice
// @author： Boice
// @createTime：2022/5/30 10:18
package accountservice

import (
	"gitlab.com/bns-engineering/td/common/config"
	logger "gitlab.com/bns-engineering/td/common/log"
	"testing"
)

func init() {
	logger.SetUp(config.Setup("../../../../config.yaml"))
}

func TestGetAccountById(t *testing.T) {
	_, err := GetAccountById("11249460359")
	if err != nil {
		t.Errorf("test error")
	}
}

func TestUndoMaturityDate(t *testing.T) {
	UndoMaturityDate("11249460359")
}

func TestApplyProfit(t *testing.T) {
	ApplyProfit("11249460359", "ok")
}
