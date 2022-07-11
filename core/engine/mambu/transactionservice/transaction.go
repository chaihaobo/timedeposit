// Package transactionservice
// @author： Boice
// @createTime：2022/5/26 18:37
package transactionservice

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/shopspring/decimal"
	"github.com/uniplaces/carbon"
	"gitlab.com/bns-engineering/common/tracer"
	"gitlab.com/bns-engineering/td/common/config"
	"gitlab.com/bns-engineering/td/common/constant"
	"gitlab.com/bns-engineering/td/common/log"
	"gitlab.com/bns-engineering/td/common/util/id"
	"gitlab.com/bns-engineering/td/common/util/mambu_http"
	"gitlab.com/bns-engineering/td/common/util/mambu_http/persistence"
	"gitlab.com/bns-engineering/td/model/mambu"
	"gitlab.com/bns-engineering/td/repository"
	"go.uber.org/zap"
	"time"
)

func GetTransactionByQueryParam(context context.Context, enCodeKey string, taskCreateTime time.Time) ([]mambu.TransactionBrief, error) {
	searchParam := generateTransactionSearchParam(enCodeKey, taskCreateTime)
	tmpTransList := []mambu.TransactionBrief{}
	postUrl := constant.UrlOf(constant.SearchTransactionUrl)

	queryParamByte, err := json.Marshal(searchParam)
	if err != nil {
		log.Error(context, "Convert searchParam to JsonStr Failed.", err, zap.Any("searchParam", searchParam))
		return tmpTransList, nil
	}
	postJsonStr := string(queryParamByte)

	err = mambu_http.Post(context, postUrl, postJsonStr, &tmpTransList, nil, persistence.DBPersistence(context, "GetTransactionByQueryParam"))

	if err != nil {
		log.Error(context, "Search td account Info List failed!", err, zap.String("queryParam", postJsonStr))
		return tmpTransList, err
	}
	return tmpTransList, nil
}

func GetAdditionProfitAndTax(tmpTDAccount *mambu.TDAccount, lastAppliedInterestTrans mambu.TransactionBrief) (float64, float64) {
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

func generateTransactionSearchParam(encodedKey string, taskCreateTime time.Time) mambu.SearchParam {
	tmpQueryParam := mambu.SearchParam{
		FilterCriteria: []mambu.FilterCriteria{
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
				Field:       "creationDate",
				Operator:    "BETWEEN",
				Value:       carbon.NewCarbon(taskCreateTime).DateString(),            // today
				SecondValue: carbon.NewCarbon(taskCreateTime).AddDays(1).DateString(), // tomorrow
			},
		},
		SortingCriteria: mambu.SortingCriteria{
			Field: "id",
			Order: "DESC",
		},
	}
	return tmpQueryParam
}

func AdjustTransaction(ctx context.Context, transactionId string, notes string) error {
	tr := tracer.StartTrace(ctx, "transactionservice-AdjustTransaction")
	ctx = tr.Context()
	defer tr.Finish()
	adjustUrl := fmt.Sprintf(constant.UrlOf(constant.AdjustTransactionUrl), transactionId)
	noteBody := struct {
		Notes string `json:"notes"`
	}{
		Notes: notes,
	}
	marshal, _ := json.Marshal(noteBody)
	err := mambu_http.Post(ctx, adjustUrl, string(marshal), nil, nil, persistence.DBPersistence(ctx, "AdjustTransaction"))
	return err
}

func WithdrawTransaction(context context.Context, tdAccount, benefitAccount *mambu.TDAccount,
	amount float64, tranDesc1 string, tranDesc3 string,
	transactionID, channelID string,
	transactionReqConfigure func(transactionReq *mambu.TransactionReq),
) (mambu.TransactionResp, error) {
	tr := tracer.StartTrace(context, "transactionservice-WithdrawTransaction")
	context = tr.Context()
	defer tr.Finish()
	transaction := repository.GetFlowTransactionRepository().GetTransactionByTransId(context, transactionID)
	if transaction != nil {
		errMsg := "transaction is ready submit"
		err := errors.New(errMsg)
		log.Error(context, errMsg, err, zap.String("transactionID", transactionID))
		return mambu.TransactionResp{}, err
	}
	transactionDetailID := transactionID + "-" + time.Now().Format("20060102150405")
	custMessage := fmt.Sprintf("Withdraw for flowTask: %v", transactionID)
	tmpTransaction := BuildTransactionReq(tdAccount, benefitAccount, transactionID, transactionDetailID, custMessage, tranDesc1, tranDesc3, amount, channelID)
	if transactionReqConfigure != nil {
		transactionReqConfigure(tmpTransaction)
	}
	var transactionResp mambu.TransactionResp
	queryParamByte, err := json.Marshal(tmpTransaction)
	if err != nil {
		log.Error(context, fmt.Sprintf("Convert searchParam to JsonStr Failed. searchParam: %v", queryParamByte), err)
		repository.GetFlowTransactionRepository().CreateFailedTransaction(context, tmpTransaction, constant.TransactionWithdraw, err.Error())
		return transactionResp, errors.New("build withdraw parameters failed")
	}
	postJsonStr := string(queryParamByte)

	postUrl := fmt.Sprintf(constant.UrlOf(constant.WithdrawTransactionUrl), tdAccount.ID)
	err = mambu_http.Post(context, postUrl, postJsonStr, &transactionResp, nil, persistence.DBPersistence(context, "WithdrawTransaction"))

	if err != nil {
		log.Error(context, fmt.Sprintf("Withdraw Transaction Error! td acc id: %v", tdAccount.ID), err)
		repository.GetFlowTransactionRepository().CreateFailedTransaction(context, tmpTransaction, constant.TransactionWithdraw, err.Error())
		return transactionResp, err
	}

	repository.GetFlowTransactionRepository().CreateSucceedFlowTransaction(context, &transactionResp)
	return transactionResp, nil
}
func DepositTransaction(context context.Context, tdAccount, benefitAccount *mambu.TDAccount, amount float64,
	tranDesc1 string, tranDesc3 string,
	transactionID, channelID string,
	transactionReqConfigure func(transactionReq *mambu.TransactionReq),
) (mambu.TransactionResp, error) {
	tr := tracer.StartTrace(context, "transactionservice-DepositTransaction")
	context = tr.Context()
	defer tr.Finish()
	transaction := repository.GetFlowTransactionRepository().GetTransactionByTransId(context, transactionID)
	if transaction != nil {
		errMsg := "transaction is ready submit"
		err := errors.New(errMsg)
		log.Error(context, errMsg, err, zap.String("transactionID", transactionID))
		return mambu.TransactionResp{}, nil
	}
	transactionDetailID := transactionID + "-" + time.Now().Format("20060102150405")
	custMessage := fmt.Sprintf("Deposit for flowTask: %v", transactionID)
	tmpTransaction := BuildTransactionReq(tdAccount, benefitAccount, transactionID, transactionDetailID, custMessage, tranDesc1, tranDesc3, amount, channelID)
	if transactionReqConfigure != nil {
		transactionReqConfigure(tmpTransaction)
	}
	var transactionResp mambu.TransactionResp
	queryParamByte, err := json.Marshal(tmpTransaction)
	if err != nil {
		log.Error(context, fmt.Sprintf("Convert searchParam to JsonStr Failed. searchParam: %v", queryParamByte), err)
		repository.GetFlowTransactionRepository().CreateFailedTransaction(context, tmpTransaction, constant.TransactionDeposit, err.Error())
		return transactionResp, errors.New("build withdraw parameters failed")
	}
	postJsonStr := string(queryParamByte)

	postUrl := fmt.Sprintf(constant.UrlOf(constant.DepositTransactionUrl), benefitAccount.ID)
	err = mambu_http.Post(context, postUrl, postJsonStr, &transactionResp, nil, persistence.DBPersistence(context, "DepositTransaction"))

	if err != nil {
		log.Error(context, fmt.Sprintf("Deposit Transaction Error! td acc id: %v", tdAccount.ID), err)
		repository.GetFlowTransactionRepository().CreateFailedTransaction(context, tmpTransaction, constant.TransactionDeposit, err.Error())
		return transactionResp, err
	}
	repository.GetFlowTransactionRepository().CreateSucceedFlowTransaction(context, &transactionResp)
	return transactionResp, nil
}

func BuildTransactionReq(tdAccount *mambu.TDAccount,
	benefitAccount *mambu.TDAccount,
	transactionID string,
	transactionDetailID string,
	custMessage string,
	tranDesc1 string,
	tranDesc3 string,
	amount float64,
	channelID string) *mambu.TransactionReq {
	tmpTransaction := &mambu.TransactionReq{
		Metadata: mambu.TransactionReqMetadata{
			MessageType:                 config.TDConf.TransactionReqMetaData.MessageType,
			ExternalTransactionID:       transactionID,
			ExternalTransactionDetailID: transactionDetailID,
			TransactionType:             config.TDConf.TransactionReqMetaData.TransactionType,
			TransactionDateTime:         carbon.Now().String(),
			TerminalType:                config.TDConf.TransactionReqMetaData.TerminalType,
			TerminalID:                  config.TDConf.TransactionReqMetaData.TerminalID,
			TerminalLocation:            config.TDConf.TransactionReqMetaData.TerminalLocation,
			TerminalRRN:                 GenerationTerminalRRN(),
			ProductCode:                 config.TDConf.TransactionReqMetaData.ProductCode,
			AcquirerIID:                 config.TDConf.TransactionReqMetaData.AcquirerIID,
			ForwarderIID:                config.TDConf.TransactionReqMetaData.ForwarderIID,
			IssuerIID:                   config.TDConf.TransactionReqMetaData.IssuerIID,
			IssuerIName:                 config.TDConf.TransactionReqMetaData.IssuerIName,
			DestinationIID:              config.TDConf.TransactionReqMetaData.DestinationIID,
			SourceAccountNo:             tdAccount.ID,
			SourceAccountName:           tdAccount.Name,
			BeneficiaryAccountNo:        benefitAccount.ID,
			BeneficiaryAccountName:      benefitAccount.Name,
			Currency:                    config.TDConf.TransactionReqMetaData.Currency,
			TranDesc1:                   tranDesc1,
			TranDesc2:                   custMessage,
			TranDesc3:                   tranDesc3,
		},
		Amount: amount,
		TransactionDetails: mambu.TransactionReqDetails{
			TransactionChannelID: channelID,
		},
	}
	return tmpTransaction
}

func GenerationTerminalRRN() string {
	return "TDE-" + id.RandomSnowFlakeId()
}
