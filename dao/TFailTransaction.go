package dao

import (
	"gitlab.com/bns-engineering/td/common/constant"
	"gitlab.com/bns-engineering/td/common/db"
	"gitlab.com/bns-engineering/td/model"
	"go.uber.org/zap"
	"time"
)

const (
	NoRetryStatus      = 0
	RetrySuccessStatus = 1
	RetryFailStatus    = 2
)

func SaveFailTransactionLog(flowTransaction *model.TFlowTransactions) {
	result := flowTransaction.Result
	if constant.TransactionFailed != result {
		zap.L().Error("transaction is not failed,not need to save")
		return
	}
	flowTransactionId := flowTransaction.Id
	if flowTransactionId <= 0 {
		zap.L().Error("flow transaction id not set,are you sure save flow transaction?")
		return
	}
	failTransaction := &model.TFailTransactions{
		FlowTransactionId:  int64(flowTransactionId),
		TransId:            flowTransaction.TransId,
		TerminalRrn:        flowTransaction.TerminalRrn,
		SourceAccountNo:    flowTransaction.SourceAccountNo,
		SourceAccountName:  flowTransaction.SourceAccountName,
		BenefitAccountNo:   flowTransaction.BenefitAccountNo,
		BenefitAccountName: flowTransaction.BenefitAccountName,
		Amount:             flowTransaction.Amount,
		Channel:            flowTransaction.Channel,
		TransactionType:    flowTransaction.TransactionType,
		RetryStatus:        NoRetryStatus,
		RetryTimes:         0,
		CreateTime:         time.Now(),
		UpdateTime:         time.Now(),
	}
	db.GetDB().Save(failTransaction)
}

func GetFailTransactionLog(transactionId string) *model.TFailTransactions {
	failTransaction := new(model.TFailTransactions)
	result := db.GetDB().Where("trans_id", transactionId).Last(failTransaction)
	if result.RowsAffected > 0 {
		return failTransaction
	} else {
		return nil
	}
}

func FailTransactionLogRetryFail(failTransaction *model.TFailTransactions) {
	failTransaction.RetryTimes++
	failTransaction.RetryStatus = RetryFailStatus
	failTransaction.UpdateTime = time.Now()
	db.GetDB().Save(failTransaction)
}

func FailTransactionLogRetrySuccess(failTransaction *model.TFailTransactions) {
	failTransaction.RetryTimes++
	failTransaction.RetryStatus = RetrySuccessStatus
	failTransaction.UpdateTime = time.Now()
	db.GetDB().Save(failTransaction)
}
