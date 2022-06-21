// Package node
// @author： Boice
// @createTime：2022/5/26 18:29
package node

import (
	"errors"
	"fmt"
	"gitlab.com/bns-engineering/td/common/config"
	"gitlab.com/bns-engineering/td/core/engine/mambu/transactionservice"
	"gitlab.com/bns-engineering/td/core/engine/node/constant"
	"gitlab.com/bns-engineering/td/model/mambu"
	"go.uber.org/zap"
)

type WithdrawNetprofitNode struct {
	*Node
}

func (node *WithdrawNetprofitNode) Run() (INodeResult, error) {

	account, err := node.GetMambuAccount(node.AccountId, true)
	if err != nil {
		return nil, err
	}
	// Get benefit account info
	benefitAccount, err := node.GetMambuBenefitAccountAccount(account.OtherInformation.BhdNomorRekPencairan, true)
	if err != nil {
		zap.L().Error(fmt.Sprintf("Failed to get benefit acc info of td account: %v, benefit acc id:%v", account.ID, account.OtherInformation.BhdNomorRekPencairan))
		return nil, errors.New("call mambu get benefit acc info failed")
	}
	if account.IsCaseB1_1() {
		if !account.IsValidBenefitAccount(benefitAccount, config.TDConf.TransactionReqMetaData.LocalHolderKey) {
			zap.L().Error("is not a valid benefit account!")
			return nil, constant.ErrBenefitAccountInvalid
		}
		netProfit, err := account.GetNetProfit()
		if err != nil {
			return nil, err
		}
		// Withdraw netProfit for deposit account
		channelID := fmt.Sprintf("RAKTRAN_DEPMUDC_%vM", account.OtherInformation.Tenor)
		withdrawTransID := node.FlowId + "-" + node.NodeName + "-" + "Withdraw"
		withrawResp, err := transactionservice.WithdrawTransaction(node.GetContext(), account, benefitAccount, netProfit,
			config.TDConf.TransactionReqMetaData.TranDesc.WithdrawNetprofitTranDesc1,

			config.TDConf.TransactionReqMetaData.TranDesc.WithdrawNetprofitTranDesc3,
			withdrawTransID, channelID, func(transactionReq *mambu.TransactionReq) {
				transactionReq.Metadata.TranDesc2 = account.ID
			})
		if err != nil {
			zap.L().Error(fmt.Sprintf("Failed to withdraw for td account: %v", account.ID))
			// Just return error, no need to reverse
			return nil, errors.New("call mambu withdraw failed")
		}
		zap.L().Info(fmt.Sprintf("Finish withdraw balance for accNo: %v, encodedKey:%v", account.ID, withrawResp.EncodedKey))

	} else {
		zap.L().Info("not match! skip it")
		return ResultSkip, nil
	}

	return ResultSuccess, nil
}
