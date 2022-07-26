// Package node
// @author： Boice
// @createTime：2022/5/26 11:07
package node

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"gitlab.com/bns-engineering/td/common/constant"
	"gitlab.com/bns-engineering/td/model/mambu"
	"go.uber.org/zap"
	"strings"
	"time"
)

type WithdrawBalanceNode struct {
	node
}

func (node *WithdrawBalanceNode) Run(ctx context.Context) (INodeResult, error) {
	account, err := node.GetMambuAccount(ctx, node.accountId, false)
	if err != nil {
		return nil, err
	}
	// Get benefit account info
	benefitAccount, err := node.GetMambuBenefitAccountAccount(ctx, account.OtherInformation.BhdNomorRekPencairan)
	if err != nil {
		node.common.Logger.Error(ctx, "Failed to get benefit acc info of td account: %v, benefit acc id:%v", err, zap.String("account", account.ID), zap.String("benefit acc id", account.OtherInformation.BhdNomorRekPencairan))
		return nil, errors.New("call mambu get benefit acc info failed")
	}

	totalBalance := decimal.NewFromFloat(account.Balances.TotalBalance).Round(2).InexactFloat64()
	if (account.IsCaseB3(node.taskCreateTime) || account.IsCaseC()) && totalBalance > 0 {
		if !account.IsValidBenefitAccount(benefitAccount, node.common.Config.TransactionReqMetaData.LocalHolderKey) {
			node.common.Logger.Error(ctx, "is not a valid benefit account!", constant.ErrBenefitAccountInvalid)
			return nil, constant.ErrBenefitAccountInvalid
		}
		rrn := node.repository.FlowNodeQueryLog.GetLogValueOr(ctx, node.flowId, node.nodeName, constant.QueryTerminalRRN, node.service.Mambu.Transaction.GenerationTerminalRRN)
		channelID := fmt.Sprintf("RAKTRAN_DEPMUDC_%vM", account.OtherInformation.Tenor)
		withdrawTransID := node.flowId + "-" + node.nodeName + "-" + "Withdraw"
		withrawResp, err := node.service.Mambu.Transaction.WithdrawTransaction(node.GetContext(ctx), account, benefitAccount, totalBalance,
			node.common.Config.TransactionReqMetaData.TranDesc.WithdrawBalanceTranDesc1,
			node.common.Config.TransactionReqMetaData.TranDesc.WithdrawBalanceTranDesc3,
			withdrawTransID, channelID, func(transactionReq *mambu.TransactionReq) {
				transactionReq.Metadata.TranDesc2 = account.ID
				transactionReq.Metadata.TerminalRRN = rrn
			})
		if err != nil {
			node.common.Logger.Error(ctx, fmt.Sprintf("Failed to withdraw for td account: %v", account.ID), err)
			return nil, errors.Wrap(err, "call mambu withdraw failed")
		}
		node.common.Logger.Info(ctx, fmt.Sprintf("Finish withdraw balance for accNo: %v, encodedKey:%v", account.ID, withrawResp.EncodedKey))
		depositTransID := node.flowId + "-" + node.nodeName + "-" + "Deposit"
		// TODO only use to test case
		if strings.EqualFold(gin.Mode(), "debug") {
			time.Sleep(node.common.Config.Flow.NodeSleepTime)
		}

		// deposit
		depositResp, err := node.service.Mambu.Transaction.DepositTransaction(node.GetContext(ctx), account, benefitAccount, totalBalance,
			node.common.Config.TransactionReqMetaData.TranDesc.DepositBalanceTranDesc1,
			node.common.Config.TransactionReqMetaData.TranDesc.DepositBalanceTranDesc3,
			depositTransID, channelID, func(transactionReq *mambu.TransactionReq) {
				transactionReq.Metadata.TranDesc2 = account.ID
				transactionReq.Metadata.TerminalRRN = rrn
			})
		if err != nil {
			node.common.Logger.Error(ctx, fmt.Sprintf("Failed to deposit for td account: %v", account.ID), err)
			node.common.Logger.Error(ctx, fmt.Sprintf("depositResp: %v", depositResp), err)

			notes := fmt.Sprintf("eod_engine-%s-%s-%s", node.flowId, node.nodeName, depositResp.ID)
			// rollback withdraw transaction
			err := node.service.Mambu.Transaction.AdjustTransaction(node.GetContext(ctx), withrawResp.ID, notes)
			if err != nil {
				err = errors.Wrap(err, "adjust rollback fail!")
			}
			return nil, errors.Wrap(err, "call mambu deposit failed")
		}
		node.common.Logger.Info(ctx, fmt.Sprintf("Finish deposit balance for accNo: %v, encodedKey:%v", account.ID, depositResp.EncodedKey))

	} else {
		node.common.Logger.Info(ctx, "not match! skip it")
		return ResultSkip, nil
	}

	return ResultSuccess, nil
}
