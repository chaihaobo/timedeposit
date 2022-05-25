/*
 * @Author: Hugo
 * @Date: 2022-05-19 07:08:05
 * @Last Modified by: Hugo
 * @Last Modified time: 2022-05-19 08:02:34
 */
package dao

import (
	"time"

	"gitlab.com/bns-engineering/td/common/constant"
	"gitlab.com/bns-engineering/td/common/db"
	"gitlab.com/bns-engineering/td/model"
	"gitlab.com/bns-engineering/td/service/mambuEntity"
)

func CreateSucceedFlowTransaction(transactionResp mambuEntity.TransactionResp) model.TFlowTransactions {
	tFlowTask := model.TFlowTransactions{
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
	return tFlowTask
}

func CreateFailedTransaction(transactionReq *mambuEntity.TransactionReq, transType string, errorMsg string) model.TFlowTransactions {
	tFlowTask := model.TFlowTransactions{
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
	SaveFailTransactionLog(&tFlowTask)
	return tFlowTask
}
