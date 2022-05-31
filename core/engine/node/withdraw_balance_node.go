// Package node
// @author： Boice
// @createTime：2022/5/26 11:07
package node

import (
	"errors"
	"fmt"
	"github.com/shopspring/decimal"
	"gitlab.com/bns-engineering/td/common/config"
	"gitlab.com/bns-engineering/td/core/engine/mambu/transactionservice"
	"go.uber.org/zap"
)

type WithdrawBalanceNode struct {
	*Node
}

func (node *WithdrawBalanceNode) Run() (INodeResult, error) {
	account, err := node.GetMambuAccount(node.AccountId, false)
	if err != nil {
		return nil, err
	}

	totalBalance := decimal.NewFromFloat(account.Balances.TotalBalance).RoundFloor(2).InexactFloat64()
	if (account.IsCaseB3() || account.IsCaseC()) && totalBalance > 0 {
		// Get benefit account info
		benefitAccount, err := node.GetMambuBenefitAccountAccount(account.OtherInformation.BhdNomorRekPencairan, false)
		if err != nil {
			zap.L().Error(fmt.Sprintf("Failed to get benefit acc info of td account: %v, benefit acc id:%v", account.ID, account.OtherInformation.BhdNomorRekPencairan))
			return nil, errors.New("call mambu get benefit acc info failed")
		}
		channelID := fmt.Sprintf("RAKTRAN_DEPMUDC_%vM", account.OtherInformation.Tenor)
		withdrawTransID := node.FlowId + "-" + node.NodeName + "-" + "Withdraw"
		withrawResp, err := transactionservice.WithdrawTransaction(node.GetContext(), account, benefitAccount, totalBalance,
			config.TDConf.TransactionReqMetaData.TranDesc.WithdrawBalanceTranDesc1,
			config.TDConf.TransactionReqMetaData.TranDesc.WithdrawBalanceTranDesc3,
			withdrawTransID, channelID)
		if err != nil {
			zap.L().Error(fmt.Sprintf("Failed to withdraw for td account: %v", account.ID))
			return nil, errors.New("call mambu withdraw failed")
		}
		zap.L().Info(fmt.Sprintf("Finish withdraw balance for accNo: %v, encodedKey:%v", account.ID, withrawResp.EncodedKey))
	} else {
		zap.L().Info("not match! skip it")
	}

	return ResultSuccess, nil
}
