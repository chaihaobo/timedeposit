// Package dao
// @author： Boice
// @createTime：
package dao

import (
	"gitlab.com/bns-engineering/td/common/config"
	logger "gitlab.com/bns-engineering/td/common/log"
	"gitlab.com/bns-engineering/td/model"
	"testing"
	"time"
)

func init() {
	logger.SetUp(config.Setup("../config.yaml"))
}

func TestFailTransactionCreate(t *testing.T) {
	t.Run("test fail transaction create", func(t *testing.T) {
		flowID := "20220523092803_11137796099"
		nodeName := "cal_additional_profit_node"

		transactions := &model.TFlowTransactions{
			Id:                 1,
			TransId:            flowID + "-" + nodeName + "-" + "Deposit",
			TerminalRrn:        config.TDConf.TransactionReqMetaData.TerminalRRN,
			SourceAccountNo:    "123456",
			SourceAccountName:  "test",
			BenefitAccountNo:   "2222222",
			BenefitAccountName: "test benefit",
			Amount:             10.1,
			Channel:            "PPH_PS42_DEPOSITO",
			TransactionType:    config.TDConf.TransactionReqMetaData.TransactionType,
			Result:             0,
			EncodedKey:         "",
			ErrorMsg:           "transfer error",
			CreateTime:         time.Now(),
			UpdateTime:         time.Now(),
		}

		SaveFailTransactionLog(transactions)
		t.Log("test success")

	})
}

func TestFailTransactionGet(t *testing.T) {
	t.Run("test get fail transaction", func(t *testing.T) {
		failTransaction := GetFailTransactionLog("20220523092803_11137796099-cal_additional_profit_node-Deposit")
		if failTransaction == nil {
			t.Error("test get fail transaction fail,because this id is present")
		}
		f := GetFailTransactionLog("s")
		if f != nil {
			t.Error("test get fail transaction fail,because this id is not present")
		}
	})

}
