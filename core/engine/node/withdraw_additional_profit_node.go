// Package node
// @author： Boice
// @createTime：2022/5/27 09:18
package node

import (
	"errors"
	"gitlab.com/bns-engineering/td/common/config"
	"gitlab.com/bns-engineering/td/core/engine/mambu/transactionservice"
	"gitlab.com/bns-engineering/td/core/engine/node/constant"
	"gitlab.com/bns-engineering/td/model/mambu"
	"go.uber.org/zap"
	"strings"
)

type WithdrawAdditionalProfitNode struct {
	*Node
}

func (node *WithdrawAdditionalProfitNode) Run() (INodeResult, error) {
	account, err := node.GetMambuAccount(node.AccountId, false)
	if err != nil {
		return nil, err
	}
	// Get benefit account info
	benefitAccount, err := node.GetMambuBenefitAccountAccount(account.OtherInformation.BhdNomorRekPencairan, false)
	if err != nil {
		zap.L().Error("Failed to get benefit acc info of td account: %v, benefit acc id:%v", zap.String("account", account.ID), zap.String("benefit acc id", account.OtherInformation.BhdNomorRekPencairan))
		return nil, errors.New("call mambu get benefit acc info failed")
	}

	if account.IsCaseB1_1_1_1() ||
		account.IsCaseB2_1_1() ||
		(account.IsCaseB3() &&
			account.Balances.TotalBalance > 0 &&
			strings.ToUpper(account.OtherInformation.IsSpecialRate) == "TRUE") ||
		(account.IsCaseC() &&
			account.Balances.TotalBalance > 0 &&
			strings.ToUpper(account.OtherInformation.IsSpecialRate) == "TRUE") {
		if !account.IsValidBenefitAccount(benefitAccount, config.TDConf.TransactionReqMetaData.LocalHolderKey) {
			zap.L().Error("is not a valid benefit account!")
			return nil, constant.ErrBenefitAccountInvalid
		}
		// Get last applied interest info
		transList, err := transactionservice.GetTransactionByQueryParam(node.GetContext(), account.EncodedKey)
		if err != nil || len(transList) <= 0 {
			zap.L().Info("No applied profit, skip")
			return nil, errors.New("no applied profit, skip")
		}
		lastAppliedInterestTrans := transList[0]

		// Calculate additionalProfit & tax of additionalProfit
		additionalProfit, _ := transactionservice.GetAdditionProfitAndTax(account, lastAppliedInterestTrans)

		// Withdraw additional profit
		withdrawTransID := node.FlowId + "-" + node.NodeName + "-" + "Withdraw"
		channelID := "BBN_BAGHAS_DEPMUDC"
		withrawResp, err := transactionservice.WithdrawTransaction(node.GetContext(), benefitAccount,
			benefitAccount, additionalProfit,
			config.TDConf.TransactionReqMetaData.TranDesc.WithdrawAdditionalProfitTranDesc1,
			config.TDConf.TransactionReqMetaData.TranDesc.WithdrawAdditionalProfitTranDesc3,
			withdrawTransID, channelID, func(transactionReq *mambu.TransactionReq) {
				transactionReq.Metadata.TranDesc2 = account.ID
				transactionReq.Metadata.BeneficiaryAccountNo = ""
				transactionReq.Metadata.BeneficiaryAccountName = ""
			})
		if err != nil {
			zap.L().Error("Failed to withdraw for td account", zap.String("account", account.ID))
			return nil, errors.New("call mambu withdraw failed")
		}
		zap.L().Info("Finish withdraw balance", zap.String("account", account.ID), zap.String("encodedKey", withrawResp.EncodedKey))
	} else {
		zap.L().Info("not match! skip it")
		return ResultSkip, nil
	}
	return ResultSuccess, nil

}
