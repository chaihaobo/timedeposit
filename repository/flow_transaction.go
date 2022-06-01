// Package repository
// @author： Boice
// @createTime：2022/5/31 16:42
package repository

import (
	"gitlab.com/bns-engineering/td/common/constant"
	"gitlab.com/bns-engineering/td/common/db"
	db2 "gitlab.com/bns-engineering/td/model/db"
	"gitlab.com/bns-engineering/td/model/mambu"
	"time"
)

var flowTransactionRepository = new(FlowTransactionRepository)

type IFlowTransactionRepository interface {
	GetTransactionByTransId(transId string) *db2.TFlowTransactions
	CreateSucceedFlowTransaction(transactionResp *mambu.TransactionResp) *db2.TFlowTransactions
	CreateFailedTransaction(transactionReq *mambu.TransactionReq, transType string, errorMsg string) *db2.TFlowTransactions
}

type FlowTransactionRepository int

func (flowTransactionRepository *FlowTransactionRepository) GetTransactionByTransId(transId string) *db2.TFlowTransactions {
	flowTransaction := new(db2.TFlowTransactions)
	rowsAffected := db.GetDB().Where("trans_id", transId).Where("result", 1).Last(flowTransaction).RowsAffected
	if rowsAffected > 0 {
		return flowTransaction
	}
	return nil
}

func (flowTransactionRepository *FlowTransactionRepository) CreateSucceedFlowTransaction(transactionResp *mambu.TransactionResp) *db2.TFlowTransactions {
	tFlowTask := db2.TFlowTransactions{
		TransId:            transactionResp.Metadata.ExternalTransactionID,
		TerminalRrn:        transactionResp.Metadata.TerminalRRN,
		SourceAccountNo:    transactionResp.Metadata.SourceAccountNo,
		SourceAccountName:  transactionResp.Metadata.SourceAccountName,
		BenefitAccountNo:   transactionResp.Metadata.BeneficiaryAccountNo,
		BenefitAccountName: transactionResp.Metadata.BeneficiaryAccountName,
		Amount:             transactionResp.Amount,
		Channel:            transactionResp.TransactionDetails.TransactionChannelID,
		TransactionType:    transactionResp.Type,
		Result:             constant.TransactionSucceed,
		EncodedKey:         transactionResp.EncodedKey,
		CreateTime:         time.Now(),
		UpdateTime:         time.Now(),
		ErrorMsg:           "",
	}
	db := db.GetDB()
	db.Save(&tFlowTask)
	return &tFlowTask
}

func (flowTransactionRepository *FlowTransactionRepository) CreateFailedTransaction(transactionReq *mambu.TransactionReq, transType string, errorMsg string) *db2.TFlowTransactions {
	tFlowTask := db2.TFlowTransactions{
		TransId:            transactionReq.Metadata.ExternalTransactionID,
		TerminalRrn:        transactionReq.Metadata.TerminalRRN,
		SourceAccountNo:    transactionReq.Metadata.SourceAccountNo,
		SourceAccountName:  transactionReq.Metadata.SourceAccountName,
		BenefitAccountNo:   transactionReq.Metadata.BeneficiaryAccountNo,
		BenefitAccountName: transactionReq.Metadata.BeneficiaryAccountName,
		Amount:             transactionReq.Amount,
		Channel:            transactionReq.TransactionDetails.TransactionChannelID,
		TransactionType:    transType,
		Result:             constant.TransactionFailed,
		EncodedKey:         "",
		CreateTime:         time.Now(),
		UpdateTime:         time.Now(),
		ErrorMsg:           errorMsg,
	}
	db := db.GetDB()
	db.Save(&tFlowTask)
	return &tFlowTask
}

func GetFlowTransactionRepository() IFlowTransactionRepository {
	return flowTransactionRepository
}
