/*
 * @Author: Hugo
 * @Date: 2022-05-16 04:15:54
 * @Last Modified by: Hugo
 * @Last Modified time: 2022-05-23 11:10:53
 */
package timeDepositNode

import (
	"errors"
	"go.uber.org/zap"
	"strconv"
	"strings"
	"time"

	"gitlab.com/bns-engineering/td/common/constant"
	"gitlab.com/bns-engineering/td/common/util"
	"gitlab.com/bns-engineering/td/node"
	"gitlab.com/bns-engineering/td/service/mambuEntity"
	mambuservices "gitlab.com/bns-engineering/td/service/mambuServices"
)

//Calc the Additional Profit for TD Account
type CalcAdditionalProfitNode struct {
	node.Node
}

func NewCalcAdditionalProfitNode() *CalcAdditionalProfitNode {
	tmpNode := new(CalcAdditionalProfitNode)
	tmpNode.Name = constant.CalAdditionalProfitNode
	tmpNode.Node.NodeRun = tmpNode
	return tmpNode
}

func (node *CalcAdditionalProfitNode) RunProcess(tmpTDAccount mambuEntity.TDAccount, flowID string, nodeName string) (constant.FlowNodeStatus, error) {
	// Get the latest info of TD Account
	newTDAccount, err := mambuservices.GetTDAccountById(tmpTDAccount.ID)
	if err != nil {
		zap.L().Error("Failed to get info of td account", zap.String("account", tmpTDAccount.ID))
		errMsg := "Failed to get detail info of td account"
		zap.L().Error(errMsg)
		return constant.FlowNodeFailed, errors.New(errMsg)
	}

	// Get last applied interest info
	transList, err := mambuservices.GetTransactionByQueryParam(generateSearchParam(newTDAccount.EncodedKey))
	if err != nil || len(transList) <= 0 {
		zap.L().Info("No applied profit, skip")
		return constant.FlowNodeSkip, errors.New("No applied profit, skip")
	}
	lastAppliedInterestTrans := transList[0]

	// Get benefit account info
	benefitAccount, err := mambuservices.GetTDAccountById(newTDAccount.OtherInformation.BhdNomorRekPencairan)
	if err != nil {
		zap.L().Error("Failed to get benefit acc info of td account: %v, benefit acc id:%v", zap.String("account", newTDAccount.ID), zap.String("benefit acc id", newTDAccount.OtherInformation.BhdNomorRekPencairan))
		return constant.FlowNodeSkip, errors.New("call mambu get benefit acc info failed")
	}

	//Calculate additionalProfit & tax of additionalProfit
	additionalProfit, additionalProfitTax := getAdditionProfitAndTax(newTDAccount, lastAppliedInterestTrans)

	if newTDAccount.IsCaseB1_1_1_1() ||
		newTDAccount.IsCaseB2_1_1() ||
		(newTDAccount.IsCaseB3() &&
			newTDAccount.Balances.TotalBalance > 0 &&
			strings.ToUpper(newTDAccount.OtherInformation.IsSpecialRate) == "TRUE") ||
		(newTDAccount.IsCaseC() &&
			newTDAccount.Balances.TotalBalance > 0 &&
			strings.ToUpper(newTDAccount.OtherInformation.IsSpecialRate) == "TRUE") {
		//Withdraw additional profit
		withdrawTransID := flowID + "-" + nodeName + "-" + "Withdraw"
		channelID := "BBN_BAGHAS_DEPMUDC"
		withrawResp, err := mambuservices.WithdrawTransaction(newTDAccount, benefitAccount, additionalProfit, withdrawTransID, channelID)
		if err != nil {
			zap.L().Error("Failed to withdraw for td account", zap.String("account", newTDAccount.ID))
			return constant.FlowNodeFailed, errors.New("call mambu withdraw failed")
		}
		zap.L().Info("Finish withdraw balance", zap.String("account", newTDAccount.ID), zap.String("encodedKey", withrawResp.EncodedKey))
		//Deposit additional profit
		depositTransID := flowID + "-" + nodeName + "-" + "Deposit"
		depositChannelID := "PPH_PS42_DEPOSITO"
		depositResp, err := mambuservices.DepositTransaction(newTDAccount, benefitAccount, additionalProfitTax, depositTransID, depositChannelID)
		if err != nil {
			zap.L().Error("Failed to deposit for td account", zap.String("account", newTDAccount.ID))
			//todo: Add reverse withdraw here
			zap.L().Error("depositResp error", zap.Any("depositResp", depositResp), zap.Error(err))

			return constant.FlowNodeFailed, errors.New("call mambu deposit failed")
		}
		zap.L().Info("Finish deposit additional profit tax", zap.String("account", tmpTDAccount.ID), zap.String("encodedKey", depositResp.EncodedKey))
		return constant.FlowNodeFinish, nil
	} else {
		zap.L().Info("No need to withdraw profit", zap.String("account", newTDAccount.ID))
		return constant.FlowNodeSkip, nil
	}
}

func getAdditionProfitAndTax(tmpTDAccount mambuEntity.TDAccount, lastAppliedInterestTrans mambuEntity.TransactionBrief) (float64, float64) {
	specialER, _ := strconv.ParseFloat(tmpTDAccount.OtherInformation.SpecialER, 64)
	ER := tmpTDAccount.InterestSettings.InterestRateSettings.InterestRate
	appliedInterest := lastAppliedInterestTrans.Amount
	additionalProfit := (specialER/ER)*appliedInterest - appliedInterest
	taxRate, _ := strconv.ParseFloat(tmpTDAccount.OtherInformation.NisbahPajak, 64)
	taxRateReal := taxRate / 100
	additionalProfitTax := additionalProfit * taxRateReal
	return additionalProfit, additionalProfitTax
}

func generateSearchParam(encodedKey string) mambuEntity.SearchParam {
	tmpQueryParam := mambuEntity.SearchParam{
		FilterCriteria: []mambuEntity.FilterCriteria{
			{
				Field:    "parentAccountKey",
				Operator: "EQUALS",
				Value:    encodedKey,
			},
			{
				Field:    "type",
				Operator: "EQUALS",
				Value:    "INTEREST_APPLIED",
			},
			{
				Field: "creationDate",
				//todo: Remember to set the value to today!
				Operator:    "BETWEEN",
				Value:       util.GetDate(time.Now().AddDate(0, 0, -20)), //today
				SecondValue: util.GetDate(time.Now().AddDate(0, 0, 1)),   //tomorrow
			},
		},
		SortingCriteria: mambuEntity.SortingCriteria{
			Field: "id",
			Order: "DESC",
		},
	}
	return tmpQueryParam
}
