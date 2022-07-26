// Package node
// @author： Boice
// @createTime：2022/5/26 18:29
package node

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"gitlab.com/bns-engineering/td/common/constant"
	"gitlab.com/bns-engineering/td/model/mambu"
)

const (
	lastNodeName = "withdraw_netprofit_node"
)

type DepositNetprofitNode struct {
	node
}

func (node *DepositNetprofitNode) Run(ctx context.Context) (INodeResult, error) {
	account, err := node.GetMambuAccount(ctx, node.accountId, false)
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
			node.common.Logger.Error(ctx, "is not a valid benefit account!", constant.ErrBenefitAccountInvalid)
			return nil, constant.ErrBenefitAccountInvalid
		}

		netProfit, err := account.GetNetProfit()
		if err != nil {
			return nil, err
		}
		// Deposit netProfit to benefit account
		// get last node rnn
		rrn := node.repository.FlowNodeQueryLog.GetLogValueOr(ctx, node.flowId, lastNodeName, constant.QueryTerminalRRN, node.service.Mambu.Transaction.GenerationTerminalRRN)
		idKey := node.repository.FlowNodeQueryLog.GetLogValueOr(ctx, node.flowId, lastNodeName, constant.QueryIdempotencyKey, uuid.New().String)
		ctx = context.WithValue(ctx, constant.ContextIdempotencyKey, idKey)
		channelID := fmt.Sprintf("RAKTRAN_DEPMUDC_%vM", account.OtherInformation.Tenor)
		depositTransID := node.flowId + "-" + node.nodeName + "-" + "Deposit"
		depositResp, err := node.service.Mambu.Transaction.DepositTransaction(node.GetContext(ctx), account, benefitAccount, netProfit,
			node.common.Config.TransactionReqMetaData.TranDesc.DepositNetprofitTranDesc1,
			node.common.Config.TransactionReqMetaData.TranDesc.DepositNetprofitTranDesc3,
			depositTransID, channelID, func(transactionReq *mambu.TransactionReq) {
				transactionReq.Metadata.TranDesc2 = account.ID
				if lastTransaction := node.repository.FlowTransaction.GetLastByFlowId(ctx, node.flowId); lastTransaction.TerminalRrn != "" {
					transactionReq.Metadata.TerminalRRN = rrn
				}

			})
		if err != nil {
			node.common.Logger.Error(ctx, fmt.Sprintf("Failed to deposit for td account: %v", account.ID), err)
			return nil, errors.New("call mambu deposit failed")
		}
		node.common.Logger.Info(ctx, fmt.Sprintf("Finish deposit balance for accNo: %v, encodedKey:%v", account.ID, depositResp.EncodedKey))
	} else {
		node.common.Logger.Info(ctx, "not match! skip it")
		return ResultSkip, nil
	}

	return ResultSuccess, nil
}
