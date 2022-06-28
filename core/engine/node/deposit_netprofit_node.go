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

const (
	lastNodeName = "withdraw_netprofit_node"
)

type DepositNetprofitNode struct {
	*Node
}

func (node *DepositNetprofitNode) Run(ctx context.Context) (INodeResult, error) {
	account, err := node.GetMambuAccount(ctx, node.AccountId, false)
	if err != nil {
		return nil, err
	}
	// Get benefit account info
	benefitAccount, err := node.GetMambuBenefitAccountAccount(ctx, account.OtherInformation.BhdNomorRekPencairan, false)
	if err != nil {
		log.Error(ctx, fmt.Sprintf("Failed to get benefit acc info of td account: %v, benefit acc id:%v", account.ID, account.OtherInformation.BhdNomorRekPencairan), err)
		return nil, errors.New("call mambu get benefit acc info failed")
	}

	if account.IsCaseB1_1() {
		if !account.IsValidBenefitAccount(benefitAccount, config.TDConf.TransactionReqMetaData.LocalHolderKey) {
			log.Error(ctx, "is not a valid benefit account!", constant.ErrBenefitAccountInvalid)
			return nil, constant.ErrBenefitAccountInvalid
		}

		netProfit, err := account.GetNetProfit()
		if err != nil {
			return nil, err
		}
		// Deposit netProfit to benefit account
		// get last node rnn
		rrn := repository.GetRedisRepository().GetTerminalRRN(ctx, node.FlowId, lastNodeName, transactionservice.GenerationTerminalRRN)
		channelID := fmt.Sprintf("RAKTRAN_DEPMUDC_%vM", account.OtherInformation.Tenor)
		depositTransID := node.FlowId + "-" + node.NodeName + "-" + "Deposit"
		depositResp, err := transactionservice.DepositTransaction(node.GetContext(ctx), account, benefitAccount, netProfit,
			config.TDConf.TransactionReqMetaData.TranDesc.DepositNetprofitTranDesc1,
			config.TDConf.TransactionReqMetaData.TranDesc.DepositNetprofitTranDesc3,
			depositTransID, channelID, func(transactionReq *mambu.TransactionReq) {
				transactionReq.Metadata.TranDesc2 = account.ID
				if lastTransaction := repository.GetFlowTransactionRepository().GetLastByFlowId(ctx, node.FlowId); lastTransaction.TerminalRrn != "" {
					transactionReq.Metadata.TerminalRRN = rrn
				}

			})
		if err != nil {
			log.Error(ctx, fmt.Sprintf("Failed to deposit for td account: %v", account.ID), err)
			return nil, errors.New("call mambu deposit failed")
		}
		log.Info(ctx, fmt.Sprintf("Finish deposit balance for accNo: %v, encodedKey:%v", account.ID, depositResp.EncodedKey))
	} else {
		log.Info(ctx, "not match! skip it")
		return ResultSkip, nil
	}

	return ResultSuccess, nil
}
