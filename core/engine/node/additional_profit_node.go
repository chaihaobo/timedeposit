// Package node
// @author： Boice
// @createTime：2022/5/27 09:18
package node

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"gitlab.com/bns-engineering/td/common/config"
	"gitlab.com/bns-engineering/td/common/log"
	"gitlab.com/bns-engineering/td/core/engine/mambu/transactionservice"
	"gitlab.com/bns-engineering/td/core/engine/node/constant"
	"gitlab.com/bns-engineering/td/model/mambu"
	"gitlab.com/bns-engineering/td/repository"
	"go.uber.org/zap"
	"strings"
	"time"
)

type AdditionalProfitNode struct {
	*Node
}

func (node *AdditionalProfitNode) Run(ctx context.Context) (INodeResult, error) {
	account, err := node.GetMambuAccount(ctx, node.AccountId, false)
	if err != nil {
		return nil, err
	}
	// Get benefit account info
	benefitAccount, err := node.GetMambuBenefitAccountAccount(ctx, account.OtherInformation.BhdNomorRekPencairan)
	if err != nil {
		log.Error(ctx, "Failed to get benefit acc info of td account: %v, benefit acc id:%v", err, zap.String("account", account.ID), zap.String("benefit acc id", account.OtherInformation.BhdNomorRekPencairan))
		return nil, errors.New("call mambu get benefit acc info failed")
	}

	if account.IsCaseB1_1_1_1() ||
		account.IsCaseB2_1_1() ||
		(account.IsCaseB3() &&
			account.Balances.TotalBalance > 0 &&
			strings.ToUpper(account.OtherInformation.IsSpecialER) == "TRUE") ||
		(account.IsCaseC() &&
			account.Balances.TotalBalance > 0 &&
			strings.ToUpper(account.OtherInformation.IsSpecialER) == "TRUE") {
		if account.IsSpecialERExpired() {
			log.Info(ctx, "special ER expired! skip it")
			return ResultSkip, nil
		}
		if !account.IsValidBenefitAccount(benefitAccount, config.TDConf.TransactionReqMetaData.LocalHolderKey) {
			log.Error(ctx, "is not a valid benefit account!", constant.ErrBenefitAccountInvalid)
			return nil, constant.ErrBenefitAccountInvalid
		}

		// Get last applied interest info
		transList, err := transactionservice.GetTransactionByQueryParam(node.GetContext(ctx), account.EncodedKey)

		if err != nil || len(transList) <= 0 {
			log.Info(ctx, "No applied profit, skip")
			return ResultSkip, nil
		}
		lastAppliedInterestTrans := transList[0]

		// Calculate additionalProfit & tax of additionalProfit
		additionalProfit, additionalProfitTax := transactionservice.GetAdditionProfitAndTax(account, lastAppliedInterestTrans)
		rrn := repository.GetFlowNodeQueryLogRepository().GetLogValueOr(ctx, node.FlowId, node.NodeName, constant.QueryTerminalRRN, transactionservice.GenerationTerminalRRN)
		// Deposit additional profit
		depositTransID := node.FlowId + "-" + node.NodeName + "-" + "Deposit"
		depositChannelID := "BBN_BONUS_DEPMUDC"
		depositResp, err := transactionservice.DepositTransaction(node.GetContext(ctx), account, benefitAccount, additionalProfit,
			config.TDConf.TransactionReqMetaData.TranDesc.DepositAdditionalProfitTranDesc1,
			config.TDConf.TransactionReqMetaData.TranDesc.DepositAdditionalProfitTranDesc3,
			depositTransID, depositChannelID, func(transactionReq *mambu.TransactionReq) {
				transactionReq.Metadata.TranDesc2 = account.ID
				transactionReq.Metadata.SourceAccountNo = ""
				transactionReq.Metadata.SourceAccountName = ""
				transactionReq.Metadata.TerminalRRN = rrn

			})
		if err != nil {
			log.Error(ctx, "Failed to deposit for td account", err, zap.String("account", account.ID))

			return nil, errors.Wrap(err, "call mambu deposit failed")
		}
		log.Info(ctx, "Finish deposit additional profit tax", zap.String("account", account.ID), zap.String("encodedKey", depositResp.EncodedKey))
		// TODO only use to test case
		if strings.EqualFold(gin.Mode(), "debug") {
			time.Sleep(config.TDConf.Flow.NodeSleepTime)
		}

		// Withdraw additional profit
		withdrawTransID := node.FlowId + "-" + node.NodeName + "-" + "Withdraw"
		channelID := "PPH_PS42_DEPOSITO"
		withrawResp, err := transactionservice.WithdrawTransaction(node.GetContext(ctx), benefitAccount,
			benefitAccount, additionalProfitTax,
			config.TDConf.TransactionReqMetaData.TranDesc.WithdrawAdditionalProfitTranDesc1,
			config.TDConf.TransactionReqMetaData.TranDesc.WithdrawAdditionalProfitTranDesc3,
			withdrawTransID, channelID, func(transactionReq *mambu.TransactionReq) {
				transactionReq.Metadata.TranDesc2 = account.ID
				transactionReq.Metadata.BeneficiaryAccountNo = ""
				transactionReq.Metadata.BeneficiaryAccountName = ""
				transactionReq.Metadata.TerminalRRN = rrn
			})
		if err != nil {
			log.Error(ctx, "Failed to withdraw for td account", err, zap.String("account", account.ID))
			// rollback deposit transaction
			notes := fmt.Sprintf("eod_engine-%s-%s-%s", node.FlowId, node.NodeName, depositResp.ID)
			adjustErr := transactionservice.AdjustTransaction(node.GetContext(ctx), depositResp.ID, notes)
			if adjustErr != nil {
				err = fmt.Errorf("adjust error:%w", err)
			}
			return nil, errors.Wrap(err, "call mambu withdraw failed")
		}
		log.Info(ctx, "Finish withdraw balance", zap.String("account", account.ID), zap.String("encodedKey", withrawResp.EncodedKey))
	} else {
		log.Info(ctx, "not match! skip it")
		return ResultSkip, nil
	}
	return ResultSuccess, nil

}
