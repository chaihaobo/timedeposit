// Package node
// @author： Boice
// @createTime：2022/5/26 11:07
package node

import (
	"errors"
	"fmt"
	"gitlab.com/bns-engineering/td/core/engine/mambu/transactionservice"
	"go.uber.org/zap"
)

type DepositBalanceNode struct {
	*Node
}

func (node *DepositBalanceNode) Run() (INodeResult, error) {
	account, err := node.GetMambuAccount(node.AccountId, false)
	if err != nil {
		return nil, err
	}
	totalBalance := account.Balances.TotalBalance
	if (account.IsCaseB3() || account.IsCaseC()) && totalBalance > 0 {
		// Get benefit account info
		benefitAccount, err := node.GetMambuBenefitAccountAccount(account.OtherInformation.BhdNomorRekPencairan, false)
		if err != nil {
			zap.L().Error(fmt.Sprintf("Failed to get benefit acc info of td account: %v, benefit acc id:%v", account.ID, account.OtherInformation.BhdNomorRekPencairan))
			return nil, errors.New("call mambu get benefit acc info failed")
		}
		channelID := fmt.Sprintf("RAKTRAN_DEPMUDC_%vM", account.OtherInformation.Tenor)
		depositTransID := node.FlowId + "-" + node.NodeName + "-" + "Deposit"
		depositResp, err := transactionservice.DepositTransaction(account, benefitAccount, totalBalance, depositTransID, channelID)
		if err != nil {
			zap.L().Error(fmt.Sprintf("Failed to deposit for td account: %v", account.ID))
			// todo: Add reverse withdraw here
			zap.L().Error(fmt.Sprintf("depositResp: %v", depositResp))
			return nil, errors.New("call mambu deposit failed")
		}
		zap.L().Info(fmt.Sprintf("Finish deposit balance for accNo: %v, encodedKey:%v", account.ID, depositResp.EncodedKey))
	} else {
		zap.L().Info("not match! skip it")
	}

	return ResultSuccess, nil
}