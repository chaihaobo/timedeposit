// Package repository
// @author： Boice
// @createTime：2022/5/31 16:42
package repository

import (
	"gitlab.com/bns-engineering/td/common/db"
	"gitlab.com/bns-engineering/td/model"
)

var flowTransactionRepository = new(FlowTransactionRepository)

type IFlowTransactionRepository interface {
	GetTransactionByTransId(transId string) *model.TFlowTransactions
}

type FlowTransactionRepository int

func (flowTransactionRepository *FlowTransactionRepository) GetTransactionByTransId(transId string) *model.TFlowTransactions {
	flowTransaction := new(model.TFlowTransactions)
	rowsAffected := db.GetDB().Where("trans_id", transId).Where("result", 1).Last(flowTransaction).RowsAffected
	if rowsAffected > 0 {
		return flowTransaction
	}
	return nil
}

func GetFlowTransactionRepository() IFlowTransactionRepository {
	return flowTransactionRepository
}
