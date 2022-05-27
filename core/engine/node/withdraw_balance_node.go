// Package node
// @author： Boice
// @createTime：2022/5/26 11:07
package node

import (
	"errors"
	"fmt"
	"gitlab.com/bns-engineering/td/core/engine/mambu/accountservice"
	"gitlab.com/bns-engineering/td/core/engine/mambu/transactionservice"
	"go.uber.org/zap"
)

type WithdrawBalanceNode struct {
	*Node
}

func (node *WithdrawBalanceNode) Run() (INodeResult, error) {
	account, err := node.GetMambuAccount()
	if err != nil {
		return nil, err
	}

	totalBalance := account.Balances.TotalBalance
	if (account.IsCaseB3() || account.IsCaseC()) && totalBalance > 0 {
		// Get benefit account info
		benefitAccount, err := accountservice.GetAccountById(account.OtherInformation.BhdNomorRekPencairan)
		if err != nil {
			zap.L().Error(fmt.Sprintf("Failed to get benefit acc info of td account: %v, benefit acc id:%v", account.ID, account.OtherInformation.BhdNomorRekPencairan))
			return nil, errors.New("call mambu get benefit acc info failed")
		}
		channelID := fmt.Sprintf("RAKTRAN_DEPMUDC_%vM", account.OtherInformation.Tenor)
		withdrawTransID := node.FlowId + "-" + node.NodeName + "-" + "Withdraw"
		withrawResp, err := transactionservice.WithdrawTransaction(account, benefitAccount, totalBalance, withdrawTransID, channelID)
		if err != nil {
			zap.L().Error(fmt.Sprintf("Failed to withdraw for td account: %v", account.ID))
			return nil, errors.New("call mambu withdraw failed")
		}
		zap.L().Info(fmt.Sprintf("Finish withdraw balance for accNo: %v, encodedKey:%v", account.ID, withrawResp.EncodedKey))
	} else {
		zap.L().Info("not match! skip it")
	}

	return NodeResultSuccess, nil
}
