// Package repository
// @author： Boice
// @createTime：2022/5/31 16:42
package repository

import (
	"context"
	"gitlab.com/bns-engineering/td/common/constant"
	"gitlab.com/bns-engineering/td/common/db"
	"gitlab.com/bns-engineering/td/model/mambu"
	"gitlab.com/bns-engineering/td/model/po"
	"time"
)

var flowTransactionRepository = new(FlowTransactionRepository)

type IFlowTransactionRepository interface {
	GetTransactionByTransId(ctx context.Context, transId string) *po.TFlowTransactions
	ListErrorTransactionByFlowId(ctx context.Context, flowId string) []po.TFlowTransactions
	CreateSucceedFlowTransaction(ctx context.Context, transactionResp *mambu.TransactionResp) *po.TFlowTransactions
	CreateFailedTransaction(ctx context.Context, transactionReq *mambu.TransactionReq, transType string, errorMsg string) *po.TFlowTransactions
}

type FlowTransactionRepository int

func (flowTransactionRepository *FlowTransactionRepository) GetTransactionByTransId(ctx context.Context, transId string) *po.TFlowTransactions {
	flowTransaction := new(po.TFlowTransactions)
	rowsAffected := db.GetDB().Where("trans_id", transId).Where("result", 1).Last(flowTransaction).RowsAffected
	if rowsAffected > 0 {
		return flowTransaction
	}
	return nil
}

func (flowTransactionRepository *FlowTransactionRepository) ListErrorTransactionByFlowId(ctx context.Context, flowId string) []po.TFlowTransactions {
	failedTransactions := make([]po.TFlowTransactions, 1)
	db.GetDB().Model(new(po.TFlowTransactions)).Where("flow_id = ? and result=0", flowId).Order("id desc").Find(&failedTransactions)
	return failedTransactions
}

func (flowTransactionRepository *FlowTransactionRepository) CreateSucceedFlowTransaction(ctx context.Context, transactionResp *mambu.TransactionResp) *po.TFlowTransactions {
	tFlowTask := po.TFlowTransactions{
		TransId:            transactionResp.Metadata.ExternalTransactionID,
		MambuTransId:       transactionResp.ID,
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
	if ctx != nil && ctx.Value("flowId") != nil {
		tFlowTask.FlowId = ctx.Value("flowId").(string)
	}

	db.GetDB().Save(&tFlowTask)
	return &tFlowTask
}

func (flowTransactionRepository *FlowTransactionRepository) CreateFailedTransaction(ctx context.Context, transactionReq *mambu.TransactionReq, transType string, errorMsg string) *po.TFlowTransactions {
	tFlowTask := po.TFlowTransactions{
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
	if ctx != nil && ctx.Value("flowId") != nil {
		tFlowTask.FlowId = ctx.Value("flowId").(string)
	}
	db.GetDB().Save(&tFlowTask)
	return &tFlowTask
}

func GetFlowTransactionRepository() IFlowTransactionRepository {
	return flowTransactionRepository
}
