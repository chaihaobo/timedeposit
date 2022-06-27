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
	"gitlab.com/bns-engineering/td/common/config"
	"gitlab.com/bns-engineering/td/common/log"
	"gitlab.com/bns-engineering/td/core/engine/mambu/transactionservice"
	"gitlab.com/bns-engineering/td/core/engine/node/constant"
	"gitlab.com/bns-engineering/td/model/mambu"
	"go.uber.org/zap"
	"strings"
	"time"
)

type WithdrawBalanceNode struct {
	*Node
}

func (node *WithdrawBalanceNode) Run(ctx context.Context) (INodeResult, error) {
	account, err := node.GetMambuAccount(ctx, node.AccountId, false)
	if err != nil {
		return nil, err
	}
	// Get benefit account info
	benefitAccount, err := node.GetMambuBenefitAccountAccount(ctx, account.OtherInformation.BhdNomorRekPencairan, false)
	if err != nil {
		log.Error(ctx, "Failed to get benefit acc info of td account: %v, benefit acc id:%v", err, zap.String("account", account.ID), zap.String("benefit acc id", account.OtherInformation.BhdNomorRekPencairan))
		return nil, errors.New("call mambu get benefit acc info failed")
	}

	totalBalance := decimal.NewFromFloat(account.Balances.TotalBalance).Round(2).InexactFloat64()
	if (account.IsCaseB3() || account.IsCaseC()) && totalBalance > 0 {
		if !account.IsValidBenefitAccount(benefitAccount, config.TDConf.TransactionReqMetaData.LocalHolderKey) {
			log.Error(ctx, "is not a valid benefit account!", constant.ErrBenefitAccountInvalid)
			return nil, constant.ErrBenefitAccountInvalid
		}
		channelID := fmt.Sprintf("RAKTRAN_DEPMUDC_%vM", account.OtherInformation.Tenor)
		withdrawTransID := node.FlowId + "-" + node.NodeName + "-" + "Withdraw"
		rrn := transactionservice.GenerationTerminalRRN()
		withrawResp, err := transactionservice.WithdrawTransaction(node.GetContext(ctx), account, benefitAccount, totalBalance,
			config.TDConf.TransactionReqMetaData.TranDesc.WithdrawBalanceTranDesc1,
			config.TDConf.TransactionReqMetaData.TranDesc.WithdrawBalanceTranDesc3,
			withdrawTransID, channelID, func(transactionReq *mambu.TransactionReq) {
				transactionReq.Metadata.TranDesc2 = account.ID
				transactionReq.Metadata.TerminalRRN = rrn
			})
		if err != nil {
			log.Error(ctx, fmt.Sprintf("Failed to withdraw for td account: %v", account.ID), err)
			return nil, errors.Wrap(err, "call mambu withdraw failed")
		}
		log.Info(ctx, fmt.Sprintf("Finish withdraw balance for accNo: %v, encodedKey:%v", account.ID, withrawResp.EncodedKey))
		depositTransID := node.FlowId + "-" + node.NodeName + "-" + "Deposit"
		// TODO only use to test case
		if strings.EqualFold(gin.Mode(), "debug") {
			time.Sleep(config.TDConf.Flow.NodeSleepTime)
		}

		// deposit
		depositResp, err := transactionservice.DepositTransaction(node.GetContext(ctx), account, benefitAccount, totalBalance,
			config.TDConf.TransactionReqMetaData.TranDesc.DepositBalanceTranDesc1,
			config.TDConf.TransactionReqMetaData.TranDesc.DepositBalanceTranDesc3,
			depositTransID, channelID, func(transactionReq *mambu.TransactionReq) {
				transactionReq.Metadata.TranDesc2 = account.ID
				transactionReq.Metadata.TerminalRRN = rrn
			})
		if err != nil {
			log.Error(ctx, fmt.Sprintf("Failed to deposit for td account: %v", account.ID), err)
			log.Error(ctx, fmt.Sprintf("depositResp: %v", depositResp), err)

			notes := fmt.Sprintf("eod_engine-%s-%s-%s", node.FlowId, node.NodeName, depositResp.ID)
			// rollback withdraw transaction
			err := transactionservice.AdjustTransaction(node.GetContext(ctx), withrawResp.ID, notes)
			if err != nil {
				err = errors.Wrap(err, "adjust rollback fail!")
			}
			return nil, errors.Wrap(err, "call mambu deposit failed")
		}
		log.Info(ctx, fmt.Sprintf("Finish deposit balance for accNo: %v, encodedKey:%v", account.ID, depositResp.EncodedKey))

	} else {
		log.Info(ctx, "not match! skip it")
		return ResultSkip, nil
	}

	return ResultSuccess, nil
}
