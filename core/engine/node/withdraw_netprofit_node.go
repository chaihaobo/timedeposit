// Package node
// @author： Boice
// @createTime：2022/5/26 18:29
package node

import (
	"context"
	"errors"
	"fmt"
	"gitlab.com/bns-engineering/td/common/config"
	"gitlab.com/bns-engineering/td/common/log"
	"gitlab.com/bns-engineering/td/core/engine/mambu/transactionservice"
	"gitlab.com/bns-engineering/td/core/engine/node/constant"
	"gitlab.com/bns-engineering/td/model/mambu"
	"gitlab.com/bns-engineering/td/repository"
)

type WithdrawNetprofitNode struct {
	*Node
}

func (node *WithdrawNetprofitNode) Run(ctx context.Context) (INodeResult, error) {

	account, err := node.GetMambuAccount(ctx, node.AccountId, true)
	if err != nil {
		return nil, err
	}
	// Get benefit account info
	benefitAccount, err := node.GetMambuBenefitAccountAccount(ctx, account.OtherInformation.BhdNomorRekPencairan, true)
	if err != nil {
		log.Error(ctx, fmt.Sprintf("Failed to get benefit acc info of td account: %v, benefit acc id:%v", account.ID, account.OtherInformation.BhdNomorRekPencairan), err)
		return nil, errors.New("call mambu get benefit acc info failed")
	}
	if account.IsCaseB1_1() {
		if !account.IsValidBenefitAccount(benefitAccount, config.TDConf.TransactionReqMetaData.LocalHolderKey) {
			log.Error(ctx, "is not a valid benefit account!", err)
			return nil, constant.ErrBenefitAccountInvalid
		}
		netProfit, err := account.GetNetProfit()
		if err != nil {
			return nil, err
		}
		rrn := repository.GetFlowNodeQueryLogRepository().GetLogValueOr(ctx, node.FlowId, node.NodeName, constant.QueryTerminalRRN, transactionservice.GenerationTerminalRRN)
		// Withdraw netProfit for deposit account
		channelID := fmt.Sprintf("RAKTRAN_DEPMUDC_%vM", account.OtherInformation.Tenor)
		withdrawTransID := node.FlowId + "-" + node.NodeName + "-" + "Withdraw"
		withrawResp, err := transactionservice.WithdrawTransaction(node.GetContext(ctx), account, benefitAccount, netProfit,
			config.TDConf.TransactionReqMetaData.TranDesc.WithdrawNetprofitTranDesc1,
			config.TDConf.TransactionReqMetaData.TranDesc.WithdrawNetprofitTranDesc3,
			withdrawTransID, channelID, func(transactionReq *mambu.TransactionReq) {
				transactionReq.Metadata.TranDesc2 = account.ID
				transactionReq.Metadata.TerminalRRN = rrn
			})
		if err != nil {
			log.Error(ctx, fmt.Sprintf("Failed to withdraw for td account: %v", account.ID), err)
			// Just return error, no need to reverse
			return nil, errors.New("call mambu withdraw failed")
		}
		log.Info(ctx, fmt.Sprintf("Finish withdraw balance for accNo: %v, encodedKey:%v", account.ID, withrawResp.EncodedKey))

	} else {
		log.Info(ctx, "not match! skip it")
		return ResultSkip, nil
	}

	return ResultSuccess, nil
}
