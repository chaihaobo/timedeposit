/*
 * @Author: Hugo
 * @Date: 2022-05-16 04:15:54
 * @Last Modified by: Hugo
 * @Last Modified time: 2022-05-19 11:23:10
 */
package timeDepositNode

import (
	"errors"
	"strconv"
	"time"

	"gitlab.com/bns-engineering/td/common/log"
	"gitlab.com/bns-engineering/td/common/util"
	"gitlab.com/bns-engineering/td/node"
	"gitlab.com/bns-engineering/td/service/mambuEntity"
	mambuservices "gitlab.com/bns-engineering/td/service/mambuServices"
)

//Calc the Additional Profit for TD Account
type CalcAdditionalProfitNode struct {
	node.Node
}

func (node *CalcAdditionalProfitNode) Process() {
	CurNodeName := "transfer_profit_node"
	tmpTDAccount, tmpFlowTask, nodeLog := node.GetAccAndFlowLog(CurNodeName)

	//Get latest info of td account
	newTDAccount, err := mambuservices.GetTDAccountById(tmpTDAccount.ID)
	if err != nil {
		log.Log.Error("Failed to get info of td account: %v", tmpTDAccount.ID)
		node.UpdateLogWhenNodeFailed(tmpFlowTask, nodeLog, errors.New("call mambu get td acc info failed"))
		return
	}
	tmpTDAccount = newTDAccount

	// Get last applied interest info
	transList, err := mambuservices.GetTransactionByQueryParam(generateSearchParam(tmpTDAccount.EncodedKey))
	if err != nil || len(transList) <= 0 {
		log.Log.Info("No applied profit, skip")
		node.UpdateLogWhenSkipNode(tmpFlowTask, CurNodeName, nodeLog)
		return
	}
	lastAppliedInterestTrans := transList[0]

	//Get benefit account
	benefitAccount, err := mambuservices.GetTDAccountById(tmpTDAccount.OtherInformation.BhdNomorRekPencairan)
	if err != nil {
		log.Log.Error("Failed to get benefit acc info of td account: %v, benefit acc id:%v", tmpTDAccount.ID, tmpTDAccount.OtherInformation.BhdNomorRekPencairan)
		node.UpdateLogWhenNodeFailed(tmpFlowTask, nodeLog, errors.New("call mambu get benefit acc info failed"))
	}

	//Calculate additionalProfit & tax of additionalProfit
	additionalProfit, additionalProfitTax := getAdditionProfitAndTax(tmpTDAccount, lastAppliedInterestTrans)

	if tmpTDAccount.IsCaseB1_1_1_1() ||
		tmpTDAccount.IsCaseB2_1_1() ||
		(tmpTDAccount.IsCaseB3() && newTDAccount.Balances.TotalBalance > 0) {
		withrawResp, err := mambuservices.WithdrawTransaction(tmpTDAccount, benefitAccount, nodeLog, additionalProfit, "BBN_BAGHAS_DEPMUDC")
		if err != nil {
			log.Log.Error("Failed to withdraw for td account: %v", tmpTDAccount.ID)
			node.UpdateLogWhenNodeFailed(tmpFlowTask, nodeLog, errors.New("call mambu withdraw failed"))
			//todo: Log failed transaction info here
			return
		}
		log.Log.Info("Finish withdraw balance for accNo: %v, encodedKey:%v", tmpTDAccount.ID, withrawResp.EncodedKey)
		depositResp, err := mambuservices.DepositTransaction(tmpTDAccount, benefitAccount, nodeLog, additionalProfitTax, "PPH_PS42_DEPOSITO")
		if err != nil {
			log.Log.Error("Failed to deposit for td account: %v", tmpTDAccount.ID)
			node.UpdateLogWhenNodeFailed(tmpFlowTask, nodeLog, errors.New("call mambu deposit failed"))
			//todo: Add reverse withdraw here
			//todo: Log failed transaction info here
		}
		log.Log.Info("Finish deposit balance for accNo: %v, encodedKey:%v", tmpTDAccount.ID, depositResp.EncodedKey)
		node.Node.Output <- tmpTDAccount
	} else {
		log.Log.Info("No need to withdraw profit, accNo: %v", tmpTDAccount.ID)
		node.UpdateLogWhenSkipNode(tmpFlowTask, CurNodeName, nodeLog)
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
