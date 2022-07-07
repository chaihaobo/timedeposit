// Package accountservice
// @author： Boice
// @createTime：2022/5/30 10:18
package accountservice

import (
	"context"
	"gitlab.com/bns-engineering/td/common/config"
	"gitlab.com/bns-engineering/td/common/util"
	"testing"
)

func init() {
	util.SetupTelemetry(config.Setup("../../../../config.json"))
}

func TestGetAccountById(t *testing.T) {
	_, err := GetAccountById(context.Background(), "11249460359")
	if err != nil {

	}
}

func TestUndoMaturityDate(t *testing.T) {
	UndoMaturityDate(context.Background(), "11249460359")
}
