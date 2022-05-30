// Package transactionservice
// @author： Boice
// @createTime：2022/5/26 18:37
package transactionservice

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/uniplaces/carbon"
	"gitlab.com/bns-engineering/td/common/config"
	"gitlab.com/bns-engineering/td/common/constant"
	"gitlab.com/bns-engineering/td/common/util"
	"gitlab.com/bns-engineering/td/common/util/mambu_http"
	"gitlab.com/bns-engineering/td/core/engine/mambu"
	"gitlab.com/bns-engineering/td/dao"
	"gitlab.com/bns-engineering/td/service/mambuEntity"
	"go.uber.org/zap"
	"strconv"
	"time"
)

func GetTransactionByQueryParam(context context.Context, enCodeKey string) ([]mambuEntity.TransactionBrief, error) {
	searchParam := generateTransactionSearchParam(enCodeKey)
	tmpTransList := []mambuEntity.TransactionBrief{}
	postUrl := constant.UrlOf(constant.SearchTransactionUrl)
	zap.L().Debug(fmt.Sprintf("postUrl: %v", postUrl))
	queryParamByte, err := json.Marshal(searchParam)
	if err != nil {
		zap.L().Error("Convert searchParam to JsonStr Failed.", zap.Any("searchParam", searchParam))
		return tmpTransList, nil
	}
	postJsonStr := string(queryParamByte)

	zap.L().Debug("transaction service", zap.String("postUrl", postUrl))
	zap.L().Debug("transaction service", zap.String("postJsonStr", postJsonStr))
	err = mambu_http.Post(postUrl, postJsonStr, &tmpTransList, mambu.SaveMambuRequestLog(context, "GetTransactionByQueryParam"))

	if err != nil {
		zap.L().Error("Search td account Info List failed!", zap.String("queryParam", postJsonStr))
		return tmpTransList, err
	}
	return tmpTransList, nil
}

func GetAdditionProfitAndTax(tmpTDAccount *mambuEntity.TDAccount, lastAppliedInterestTrans mambuEntity.TransactionBrief) (float64, float64) {
	specialER, _ := strconv.ParseFloat(tmpTDAccount.OtherInformation.SpecialER, 64)
	ER := tmpTDAccount.InterestSettings.InterestRateSettings.InterestRate
	appliedInterest := lastAppliedInterestTrans.Amount
	additionalProfit := (specialER/ER)*appliedInterest - appliedInterest
	taxRate, _ := strconv.ParseFloat(tmpTDAccount.OtherInformation.NisbahPajak, 64)
	taxRateReal := taxRate / 100
	additionalProfitTax := additionalProfit * taxRateReal
	return additionalProfit, additionalProfitTax
}

func generateTransactionSearchParam(encodedKey string) mambuEntity.SearchParam {
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
				// todo: Remember to set the value to today!
				Operator:    "BETWEEN",
				Value:       util.GetDate(time.Now().AddDate(0, 0, -20)), // today
				SecondValue: util.GetDate(time.Now().AddDate(0, 0, 1)),   // tomorrow
			},
		},
		SortingCriteria: mambuEntity.SortingCriteria{
			Field: "id",
			Order: "DESC",
		},
	}
	return tmpQueryParam
}

func WithdrawTransaction(context context.Context, tdAccount, benefitAccount *mambuEntity.TDAccount, amount float64, transactionID, channelID string) (mambuEntity.TransactionResp, error) {

	transactionDetailID := transactionID + "-" + time.Now().Format("20060102150405")
	custMessage := fmt.Sprintf("Withdraw for flowTask: %v", transactionID)
	tmpTransaction := BuildTransactionReq(tdAccount, transactionID, transactionDetailID, custMessage, amount, channelID)
	var transactionResp mambuEntity.TransactionResp
	queryParamByte, err := json.Marshal(tmpTransaction)
	if err != nil {
		zap.L().Error(fmt.Sprintf("Convert searchParam to JsonStr Failed. searchParam: %v", queryParamByte))
		dao.CreateFailedTransaction(tmpTransaction, constant.TransactionWithdraw, err.Error())
		return transactionResp, errors.New("build withdraw parameters failed")
	}
	postJsonStr := string(queryParamByte)

	postUrl := fmt.Sprintf(constant.UrlOf(constant.WithdrawTransactiontUrl), tdAccount.ID)
	err = mambu_http.Post(postUrl, postJsonStr, &transactionResp, mambu.SaveMambuRequestLog(context, "WithdrawTransaction"))

	if err != nil {
		zap.L().Error(fmt.Sprintf("Withdraw Transaction Error! td acc id: %v", tdAccount.ID))
		dao.CreateFailedTransaction(tmpTransaction, constant.TransactionWithdraw, err.Error())
		return transactionResp, err
	}

	zap.L().Debug(fmt.Sprintf("Withdraw Transaction for td account succeed. Result: %v", transactionResp))
	dao.CreateSucceedFlowTransaction(transactionResp)
	return transactionResp, nil
}
func DepositTransaction(context context.Context, tdAccount, benefitAccount *mambuEntity.TDAccount, amount float64, transactionID, channelID string) (mambuEntity.TransactionResp, error) {
	transactionDetailID := transactionID + "-" + time.Now().Format("20060102150405")
	custMessage := fmt.Sprintf("Deposit for flowTask: %v", transactionID)
	tmpTransaction := BuildTransactionReq(tdAccount, transactionID, transactionDetailID, custMessage, amount, channelID)
	var transactionResp mambuEntity.TransactionResp
	queryParamByte, err := json.Marshal(tmpTransaction)
	if err != nil {
		zap.L().Error(fmt.Sprintf("Convert searchParam to JsonStr Failed. searchParam: %v", queryParamByte))
		dao.CreateFailedTransaction(tmpTransaction, constant.TransactionWithdraw, err.Error())
		return transactionResp, errors.New("build withdraw parameters failed")
	}
	postJsonStr := string(queryParamByte)

	postUrl := fmt.Sprintf(constant.UrlOf(constant.DepositTransactiontUrl), benefitAccount.ID)
	err = mambu_http.Post(postUrl, postJsonStr, &transactionResp, mambu.SaveMambuRequestLog(context, "DepositTransaction"))

	if err != nil {
		zap.L().Error(fmt.Sprintf("Deposit Transaction Error! td acc id: %v", tdAccount.ID))
		dao.CreateFailedTransaction(tmpTransaction, constant.TransactionWithdraw, err.Error())
		return transactionResp, err
	}
	dao.CreateSucceedFlowTransaction(transactionResp)
	return transactionResp, nil
}

func BuildTransactionReq(tdAccount *mambuEntity.TDAccount, transactionID string, transactionDetailID string, custMessage string, amount float64, channelID string) *mambuEntity.TransactionReq {
	tmpTransaction := &mambuEntity.TransactionReq{
		Metadata: mambuEntity.TransactionReqMetadata{
			MessageType:                    config.TDConf.TransactionReqMetaData.MessageType,
			ExternalTransactionID:          transactionID,
			ExternalTransactionDetailID:    transactionDetailID,
			ExternalOriTransactionID:       transactionID,
			ExternalOriTransactionDetailID: transactionDetailID,
			TransactionType:                config.TDConf.TransactionReqMetaData.TransactionType,
			TransactionDateTime:            carbon.Now().String(),
			TerminalType:                   config.TDConf.TransactionReqMetaData.TerminalType,
			TerminalID:                     config.TDConf.TransactionReqMetaData.TerminalID,
			TerminalLocation:               config.TDConf.TransactionReqMetaData.TerminalLocation,
			TerminalRRN:                    generationTerminalRRN(),
			ProductCode:                    config.TDConf.TransactionReqMetaData.ProductCode,
			AcquirerIID:                    config.TDConf.TransactionReqMetaData.AcquirerIID,
			ForwarderIID:                   config.TDConf.TransactionReqMetaData.ForwarderIID,
			IssuerIID:                      config.TDConf.TransactionReqMetaData.IssuerIID,
			IssuerIName:                    config.TDConf.TransactionReqMetaData.IssuerIName,
			DestinationIID:                 config.TDConf.TransactionReqMetaData.DestinationIID,
			SourceAccountNo:                tdAccount.ID,
			SourceAccountName:              tdAccount.Name,
			BeneficiaryAccountNo:           tdAccount.OtherInformation.BhdNomorRekPencairan,
			BeneficiaryAccountName:         tdAccount.OtherInformation.BhdNamaRekPencairan,
			Currency:                       config.TDConf.TransactionReqMetaData.Currency,
			TranDesc1:                      config.TDConf.TransactionReqMetaData.TranDesc1,
			TranDesc2:                      custMessage,
			TranDesc3:                      config.TDConf.TransactionReqMetaData.TranDesc3,
		},
		Amount: amount,
		TransactionDetails: mambuEntity.TransactionReqDetails{
			TransactionChannelID: channelID,
		},
	}
	return tmpTransaction
}

func generationTerminalRRN() string {
	return "TDE-" + util.RandomSnowFlakeId()

}
