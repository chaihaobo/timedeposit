// Package node
// @author： Boice
// @createTime：2022/5/27 09:18
package node

import (
	"errors"
	"gitlab.com/bns-engineering/td/core/engine/mambu/accountservice"
	"gitlab.com/bns-engineering/td/core/engine/mambu/transactionservice"
	"go.uber.org/zap"
	"strings"
)

type DepositAdditionalProfitNode struct {
	*Node
}

func (node *DepositAdditionalProfitNode) Run() (INodeResult, error) {
	account, err := node.GetMambuAccount()
	if err != nil {
		return nil, err
	}

	if account.IsCaseB1_1_1_1() ||
		account.IsCaseB2_1_1() ||
		(account.IsCaseB3() &&
			account.Balances.TotalBalance > 0 &&
			strings.ToUpper(account.OtherInformation.IsSpecialRate) == "TRUE") ||
		(account.IsCaseC() &&
			account.Balances.TotalBalance > 0 &&
			strings.ToUpper(account.OtherInformation.IsSpecialRate) == "TRUE") {
		// Get last applied interest info
		transList, err := transactionservice.GetTransactionByQueryParam(account.EncodedKey)
		if err != nil || len(transList) <= 0 {
			zap.L().Info("No applied profit, skip")
			return nil, errors.New("No applied profit, skip")
		}
		lastAppliedInterestTrans := transList[0]

		// Get benefit account info
		benefitAccount, err := accountservice.GetAccountById(account.OtherInformation.BhdNomorRekPencairan)
		if err != nil {
			zap.L().Error("Failed to get benefit acc info of td account: %v, benefit acc id:%v", zap.String("account", account.ID), zap.String("benefit acc id", account.OtherInformation.BhdNomorRekPencairan))
			return nil, errors.New("call mambu get benefit acc info failed")
		}
		//Calculate additionalProfit & tax of additionalProfit
		_, additionalProfitTax := transactionservice.GetAdditionProfitAndTax(account, lastAppliedInterestTrans)
		//Deposit additional profit
		depositTransID := node.FlowId + "-" + node.NodeName + "-" + "Deposit"
		depositChannelID := "PPH_PS42_DEPOSITO"
		depositResp, err := transactionservice.DepositTransaction(account, benefitAccount, additionalProfitTax, depositTransID, depositChannelID)
		if err != nil {
			zap.L().Error("Failed to deposit for td account", zap.String("account", account.ID))
			//todo: Add reverse withdraw here
			zap.L().Error("depositResp error", zap.Any("depositResp", depositResp), zap.Error(err))

			return nil, errors.New("call mambu deposit failed")
		}
		zap.L().Info("Finish deposit additional profit tax", zap.String("account", account.ID), zap.String("encodedKey", depositResp.EncodedKey))
	} else {
		zap.L().Info("not match! skip it")
	}
	return NodeResultSuccess, nil

}
