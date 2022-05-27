// Package node
// @author： Boice
// @createTime：2022/5/26 18:29
package node

import (
	"errors"
	"fmt"
	"gitlab.com/bns-engineering/td/core/engine/mambu/accountservice"
	"gitlab.com/bns-engineering/td/core/engine/mambu/transactionservice"
	"go.uber.org/zap"
	"strconv"
)

type WithdrawNetprofitNode struct {
	*Node
}

func (node *WithdrawNetprofitNode) Run() (INodeResult, error) {

	account, err := node.GetMambuAccount()
	if err != nil {
		return nil, err
	}
	if account.IsCaseB1_1() {
		principal, err := strconv.ParseFloat(account.Rekening.RekeningPrincipalAmount, 64)
		if err != nil {
			errMsg := fmt.Sprintf("Failed to convert Rekening.RekeningPrincipalAmount from string to float64, value:%v", account.Rekening.RekeningPrincipalAmount)
			zap.L().Error(errMsg)
			return nil, errors.New(errMsg)
		}
		//Calculate the profit
		netProfit := account.Balances.TotalBalance - principal
		if netProfit > 0 {
			// Get benefit account info
			benefitAccount, err := accountservice.GetAccountById(account.OtherInformation.BhdNomorRekPencairan)
			if err != nil {
				zap.L().Error(fmt.Sprintf("Failed to get benefit acc info of td account: %v, benefit acc id:%v", account.ID, account.OtherInformation.BhdNomorRekPencairan))
				return nil, errors.New("call mambu get benefit acc info failed")
			}

			//Withdraw netProfit for deposit account
			channelID := fmt.Sprintf("RAKTRAN_DEPMUDC_%vM", account.OtherInformation.Tenor)
			withdrawTransID := node.FlowId + "-" + node.NodeName + "-" + "Withdraw"
			withrawResp, err := transactionservice.WithdrawTransaction(account, benefitAccount, netProfit, withdrawTransID, channelID)
			if err != nil {
				zap.L().Error(fmt.Sprintf("Failed to withdraw for td account: %v", account.ID))
				//Just return error, no need to reverse
				return nil, errors.New("call mambu withdraw failed")
			}
			zap.L().Info(fmt.Sprintf("Finish withdraw balance for accNo: %v, encodedKey:%v", account.ID, withrawResp.EncodedKey))

		}
	} else {
		zap.L().Info("not match! skip it")
	}

	return NodeResultSuccess, nil
}
