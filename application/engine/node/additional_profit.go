// Package node
// @author： Boice
// @createTime：2022/5/27 09:18
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

type AdditionalProfitNode struct {
	node
}

func (node *AdditionalProfitNode) Run(ctx context.Context) (INodeResult, error) {
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

	if account.IsCaseB1_1_1_1(node.taskCreateTime) ||
		account.IsCaseB2_1_1(node.taskCreateTime) ||
		(account.IsCaseB3(node.taskCreateTime) &&
			account.Balances.TotalBalance > 0 &&
			strings.ToUpper(account.OtherInformation.IsSpecialER) == "TRUE") ||
		(account.IsCaseC() &&
			account.Balances.TotalBalance > 0 &&
			strings.ToUpper(account.OtherInformation.IsSpecialER) == "TRUE") {
		if account.IsSpecialERExpired(node.taskCreateTime) {
			node.common.Logger.Info(ctx, "special ER expired! skip it")
			return ResultSkip, nil
		}
		if !account.IsValidBenefitAccount(benefitAccount, node.common.Config.TransactionReqMetaData.LocalHolderKey) {
			node.common.Logger.Error(ctx, "is not a valid benefit account!", constant.ErrBenefitAccountInvalid)
			return nil, constant.ErrBenefitAccountInvalid
		}

		// Get last applied interest info
		transList, err := node.service.Mambu.Transaction.GetTransactionByQueryParam(node.GetContext(ctx), account.EncodedKey, node.taskCreateTime)

		if err != nil || len(transList) <= 0 {
			node.common.Logger.Info(ctx, "No applied profit, skip")
			return ResultSkip, nil
		}
		// check ER
		specialER, erError := decimal.NewFromString(account.OtherInformation.SpecialER)
		ER := decimal.NewFromFloat(account.InterestSettings.InterestRateSettings.InterestRate)
		if erError != nil || specialER.LessThanOrEqual(ER) {
			return ResultSkip, nil
		}
		lastAppliedInterestTrans := transList[0]

		// Calculate additionalProfit & tax of additionalProfit
		additionalProfit, additionalProfitTax := getAdditionProfitAndTax(account, lastAppliedInterestTrans)
		rrn := node.repository.FlowNodeQueryLog.GetLogValueOr(ctx, node.flowId, node.nodeName, constant.QueryTerminalRRN, node.service.Mambu.Transaction.GenerationTerminalRRN)
		// Deposit additional profit
		depositTransID := node.flowId + "-" + node.nodeName + "-" + "Deposit"
		depositChannelID := "BBN_BONUS_DEPMUDC"
		depositResp, err := node.service.Mambu.Transaction.DepositTransaction(node.GetContext(ctx), account, benefitAccount, additionalProfit,
			node.common.Config.TransactionReqMetaData.TranDesc.DepositAdditionalProfitTranDesc1,
			node.common.Config.TransactionReqMetaData.TranDesc.DepositAdditionalProfitTranDesc3,
			depositTransID, depositChannelID, func(transactionReq *mambu.TransactionReq) {
				transactionReq.Metadata.TranDesc2 = account.ID
				transactionReq.Metadata.SourceAccountNo = ""
				transactionReq.Metadata.SourceAccountName = ""
				transactionReq.Metadata.TerminalRRN = rrn

			})
		if err != nil {
			node.common.Logger.Error(ctx, "Failed to deposit for td account", err, zap.String("account", account.ID))

			return nil, errors.Wrap(err, "call mambu deposit failed")
		}
		node.common.Logger.Info(ctx, "Finish deposit additional profit tax", zap.String("account", account.ID), zap.String("encodedKey", depositResp.EncodedKey))
		// TODO only use to test case
		if strings.EqualFold(gin.Mode(), "debug") {
			time.Sleep(node.common.Config.Flow.NodeSleepTime)
		}

		// Withdraw additional profit
		withdrawTransID := node.flowId + "-" + node.nodeName + "-" + "Withdraw"
		channelID := "PPH_PS42_DEPOSITO"
		withrawResp, err := node.service.Mambu.Transaction.WithdrawTransaction(node.GetContext(ctx), benefitAccount,
			benefitAccount, additionalProfitTax,
			node.common.Config.TransactionReqMetaData.TranDesc.WithdrawAdditionalProfitTranDesc1,
			node.common.Config.TransactionReqMetaData.TranDesc.WithdrawAdditionalProfitTranDesc3,
			withdrawTransID, channelID, func(transactionReq *mambu.TransactionReq) {
				transactionReq.Metadata.TranDesc2 = account.ID
				transactionReq.Metadata.BeneficiaryAccountNo = ""
				transactionReq.Metadata.BeneficiaryAccountName = ""
				transactionReq.Metadata.TerminalRRN = rrn
			})
		if err != nil {
			node.common.Logger.Error(ctx, "Failed to withdraw for td account", err, zap.String("account", account.ID))
			// rollback deposit transaction
			notes := fmt.Sprintf("eod_engine-%s-%s-%s", node.flowId, node.nodeName, depositResp.ID)
			adjustErr := node.service.Mambu.Transaction.AdjustTransaction(node.GetContext(ctx), depositResp.ID, notes)
			if adjustErr != nil {
				err = fmt.Errorf("adjust error:%w", err)
			}
			return nil, errors.Wrap(err, "call mambu withdraw failed")
		}
		node.common.Logger.Info(ctx, "Finish withdraw balance", zap.String("account", account.ID), zap.String("encodedKey", withrawResp.EncodedKey))
	} else {
		node.common.Logger.Info(ctx, "not match! skip it")
		return ResultSkip, nil
	}
	return ResultSuccess, nil

}

func getAdditionProfitAndTax(tmpTDAccount *mambu.TDAccount, lastAppliedInterestTrans mambu.TransactionBrief) (float64, float64) {
	// specialER, _ := strconv.ParseFloat(tmpTDAccount.OtherInformation.SpecialER, 64)
	// ER := tmpTDAccount.InterestSettings.InterestRateSettings.InterestRate
	// appliedInterest := lastAppliedInterestTrans.Amount
	// additionalProfit := (specialER/ER)*appliedInterest - appliedInterest
	// taxRate, _ := strconv.ParseFloat(tmpTDAccount.OtherInformation.NisbahPajak, 64)
	// taxRateReal := taxRate / 100
	// additionalProfitTax := additionalProfit * taxRateReal
	specialER, _ := decimal.NewFromString(tmpTDAccount.OtherInformation.SpecialER)
	ER := decimal.NewFromFloat(tmpTDAccount.InterestSettings.InterestRateSettings.InterestRate)
	appliedInterest := decimal.NewFromFloat(lastAppliedInterestTrans.Amount)
	additionalProfit := specialER.Div(ER).Mul(appliedInterest).Sub(appliedInterest)
	taxRate, _ := decimal.NewFromString(tmpTDAccount.OtherInformation.NisbahPajak)
	taxRateReal := taxRate.Div(decimal.NewFromInt(100))
	additionalProfitTax := additionalProfit.Mul(taxRateReal)
	return additionalProfit.Round(2).InexactFloat64(), additionalProfitTax.Round(2).InexactFloat64()
}
