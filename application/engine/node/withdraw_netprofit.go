// Package node
// @author： Boice
// @createTime：2022/5/26 18:29
package node

import (
	"context"
	"errors"
	"fmt"
	"gitlab.com/bns-engineering/td/common/constant"
	"gitlab.com/bns-engineering/td/model/mambu"
)

type WithdrawNetprofitNode struct {
	node
}

func (node *WithdrawNetprofitNode) Run(ctx context.Context) (INodeResult, error) {
	account, err := node.GetMambuAccount(ctx, node.accountId, true)
	if err != nil {
		return nil, err
	}
	// Get benefit account info
	benefitAccount, err := node.GetMambuBenefitAccountAccount(ctx, account.OtherInformation.BhdNomorRekPencairan)
	if err != nil {
		node.common.Logger.Error(ctx, fmt.Sprintf("Failed to get benefit acc info of td account: %v, benefit acc id:%v", account.ID, account.OtherInformation.BhdNomorRekPencairan), err)
		return nil, errors.New("call mambu get benefit acc info failed")
	}
	if account.IsCaseB1_1(node.taskCreateTime) {
		if !account.IsValidBenefitAccount(benefitAccount, node.common.Config.TransactionReqMetaData.LocalHolderKey) {
			node.common.Logger.Error(ctx, "is not a valid benefit account!", err)
			return nil, constant.ErrBenefitAccountInvalid
		}
		netProfit, err := account.GetNetProfit()
		if err != nil {
			return nil, err
		}
		rrn := node.repository.FlowNodeQueryLog.GetLogValueOr(ctx, node.flowId, node.nodeName, constant.QueryTerminalRRN, node.service.Mambu.Transaction.GenerationTerminalRRN)
		// Withdraw netProfit for deposit account
		channelID := fmt.Sprintf("RAKTRAN_DEPMUDC_%vM", account.OtherInformation.Tenor)
		withdrawTransID := node.flowId + "-" + node.nodeName + "-" + "Withdraw"
		withrawResp, err := node.service.Mambu.Transaction.WithdrawTransaction(node.GetContext(ctx), account, benefitAccount, netProfit,
			node.common.Config.TransactionReqMetaData.TranDesc.WithdrawNetprofitTranDesc1,
			node.common.Config.TransactionReqMetaData.TranDesc.WithdrawNetprofitTranDesc3,
			withdrawTransID, channelID, func(transactionReq *mambu.TransactionReq) {
				transactionReq.Metadata.TranDesc2 = account.ID
				transactionReq.Metadata.TerminalRRN = rrn
			})
		if err != nil {
			node.common.Logger.Error(ctx, fmt.Sprintf("Failed to withdraw for td account: %v", account.ID), err)
			// Just return error, no need to reverse
			return nil, errors.New("call mambu withdraw failed")
		}
		node.common.Logger.Info(ctx, fmt.Sprintf("Finish withdraw balance for accNo: %v, encodedKey:%v", account.ID, withrawResp.EncodedKey))

	} else {
		node.common.Logger.Info(ctx, "not match! skip it")
		return ResultSkip, nil
	}

	return ResultSuccess, nil
}
