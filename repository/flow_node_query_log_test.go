// Package repository
// @author： Boice
// @createTime：2022/5/31 11:28
package repository

import (
	"fmt"
	"gitlab.com/bns-engineering/td/common/config"
	"gitlab.com/bns-engineering/td/common/logger"
	"testing"
)

func init() {
	logger.SetUp(config.Setup("../config.yaml"))
}

func TestFlowNodeQueryLogRepository_GetNewLog(t *testing.T) {
	log := GetFlowNodeQueryLogRepository().GetNewLog("20220531110451_11747126703", "QueryTDAccount")
	fmt.Printf("%v", log)
}
