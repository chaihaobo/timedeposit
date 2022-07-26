// Package mambu
// @author： Boice
// @createTime：2022/7/26 13:54
package mambu

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/uniplaces/carbon"
	"gitlab.com/bns-engineering/common/tracer"
	"gitlab.com/bns-engineering/td/common"
	"gitlab.com/bns-engineering/td/common/constant"
	"gitlab.com/bns-engineering/td/common/util"
	"gitlab.com/bns-engineering/td/model/mambu"
	"gitlab.com/bns-engineering/td/repository"
	"go.uber.org/zap"
	"time"
)

type TransactionService interface {
	GetTransactionByQueryParam(context context.Context, enCodeKey string, taskCreateTime time.Time) ([]mambu.TransactionBrief, error)
	AdjustTransaction(ctx context.Context, transactionId string, notes string) error
	WithdrawTransaction(context context.Context, tdAccount, benefitAccount *mambu.TDAccount,
		amount float64, tranDesc1 string, tranDesc3 string,
		transactionID, channelID string,
		transactionReqConfigure func(transactionReq *mambu.TransactionReq),
	) (mambu.TransactionResp, error)

	DepositTransaction(context context.Context, tdAccount, benefitAccount *mambu.TDAccount, amount float64,
		tranDesc1 string, tranDesc3 string,
		transactionID, channelID string,
		transactionReqConfigure func(transactionReq *mambu.TransactionReq),
	) (mambu.TransactionResp, error)
	GenerationTerminalRRN() string
}

func newTransactionService(common *common.Common, repo *repository.Repository, mambuClient Client) TransactionService {
	return &transactionService{
		common:      common,
		repository:  repo,
		mambuClient: mambuClient,
	}
}

type transactionService struct {
	common      *common.Common
	repository  *repository.Repository
	mambuClient Client
}

func (t *transactionService) GetTransactionByQueryParam(context context.Context, enCodeKey string, taskCreateTime time.Time) ([]mambu.TransactionBrief, error) {
	searchParam := generateTransactionSearchParam(enCodeKey, taskCreateTime)
	var tmpTransList []mambu.TransactionBrief
	postUrl := constant.UrlOf(constant.SearchTransactionUrl)

	queryParamByte, err := json.Marshal(searchParam)
	if err != nil {
		t.common.Logger.Error(context, "Convert searchParam to JsonStr Failed.", err, zap.Any("searchParam", searchParam))
		return tmpTransList, nil
	}
	postJsonStr := string(queryParamByte)

	err = t.mambuClient.Post(context, postUrl, postJsonStr, &tmpTransList, nil, t.mambuClient.DBPersistence(context, "GetTransactionByQueryParam"))

	if err != nil {
		t.common.Logger.Error(context, "Search td account Info List failed!", err, zap.String("queryParam", postJsonStr))
		return tmpTransList, err
	}
	return tmpTransList, nil
}

func (t *transactionService) AdjustTransaction(ctx context.Context, transactionId string, notes string) error {
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
	err := t.mambuClient.Post(ctx, adjustUrl, string(marshal), nil, nil, t.mambuClient.DBPersistence(ctx, "AdjustTransaction"))
	return err
}

func (t *transactionService) WithdrawTransaction(context context.Context, tdAccount, benefitAccount *mambu.TDAccount,
	amount float64, tranDesc1 string, tranDesc3 string,
	transactionID, channelID string,
	transactionReqConfigure func(transactionReq *mambu.TransactionReq),
) (mambu.TransactionResp, error) {
	tr := tracer.StartTrace(context, "transactionservice-WithdrawTransaction")
	context = tr.Context()
	defer tr.Finish()
	transaction := t.repository.FlowTransaction.GetTransactionByTransId(context, transactionID)
	if transaction != nil {
		errMsg := "transaction is ready submit"
		err := errors.New(errMsg)
		t.common.Logger.Error(context, errMsg, err, zap.String("transactionID", transactionID))
		return mambu.TransactionResp{}, err
	}
	transactionDetailID := transactionID + "-" + time.Now().Format("20060102150405")
	custMessage := fmt.Sprintf("Withdraw for flowTask: %v", transactionID)
	tmpTransaction := t.BuildTransactionReq(tdAccount, benefitAccount, transactionID, transactionDetailID, custMessage, tranDesc1, tranDesc3, amount, channelID)
	if transactionReqConfigure != nil {
		transactionReqConfigure(tmpTransaction)
	}
	var transactionResp mambu.TransactionResp
	queryParamByte, err := json.Marshal(tmpTransaction)
	if err != nil {
		t.common.Logger.Error(context, fmt.Sprintf("Convert searchParam to JsonStr Failed. searchParam: %v", queryParamByte), err)
		t.repository.FlowTransaction.CreateFailedTransaction(context, tmpTransaction, constant.TransactionWithdraw, err.Error())
		return transactionResp, errors.New("build withdraw parameters failed")
	}
	postJsonStr := string(queryParamByte)

	postUrl := fmt.Sprintf(constant.UrlOf(constant.WithdrawTransactionUrl), tdAccount.ID)
	err = t.mambuClient.Post(context, postUrl, postJsonStr, &transactionResp, nil, t.mambuClient.DBPersistence(context, "WithdrawTransaction"))

	if err != nil {
		t.common.Logger.Error(context, fmt.Sprintf("Withdraw Transaction Error! td acc id: %v", tdAccount.ID), err)
		t.repository.FlowTransaction.CreateFailedTransaction(context, tmpTransaction, constant.TransactionWithdraw, err.Error())
		return transactionResp, err
	}

	t.repository.FlowTransaction.CreateSucceedFlowTransaction(context, &transactionResp)
	return transactionResp, nil
}

func (t *transactionService) DepositTransaction(context context.Context, tdAccount, benefitAccount *mambu.TDAccount, amount float64,
	tranDesc1 string, tranDesc3 string,
	transactionID, channelID string,
	transactionReqConfigure func(transactionReq *mambu.TransactionReq),
) (mambu.TransactionResp, error) {
	tr := tracer.StartTrace(context, "transactionservice-DepositTransaction")
	context = tr.Context()
	defer tr.Finish()
	transaction := t.repository.FlowTransaction.GetTransactionByTransId(context, transactionID)
	if transaction != nil {
		errMsg := "transaction is ready submit"
		err := errors.New(errMsg)
		t.common.Logger.Error(context, errMsg, err, zap.String("transactionID", transactionID))
		return mambu.TransactionResp{}, nil
	}
	transactionDetailID := transactionID + "-" + time.Now().Format("20060102150405")
	custMessage := fmt.Sprintf("Deposit for flowTask: %v", transactionID)
	tmpTransaction := t.BuildTransactionReq(tdAccount, benefitAccount, transactionID, transactionDetailID, custMessage, tranDesc1, tranDesc3, amount, channelID)
	if transactionReqConfigure != nil {
		transactionReqConfigure(tmpTransaction)
	}
	var transactionResp mambu.TransactionResp
	queryParamByte, err := json.Marshal(tmpTransaction)
	if err != nil {
		t.common.Logger.Error(context, fmt.Sprintf("Convert searchParam to JsonStr Failed. searchParam: %v", queryParamByte), err)
		t.repository.FlowTransaction.CreateFailedTransaction(context, tmpTransaction, constant.TransactionDeposit, err.Error())
		return transactionResp, errors.New("build withdraw parameters failed")
	}
	postJsonStr := string(queryParamByte)

	postUrl := fmt.Sprintf(constant.UrlOf(constant.DepositTransactionUrl), benefitAccount.ID)
	err = t.mambuClient.Post(context, postUrl, postJsonStr, &transactionResp, nil, t.mambuClient.DBPersistence(context, "DepositTransaction"))

	if err != nil {
		t.common.Logger.Error(context, fmt.Sprintf("Deposit Transaction Error! td acc id: %v", tdAccount.ID), err)
		t.repository.FlowTransaction.CreateFailedTransaction(context, tmpTransaction, constant.TransactionDeposit, err.Error())
		return transactionResp, err
	}
	t.repository.FlowTransaction.CreateSucceedFlowTransaction(context, &transactionResp)
	return transactionResp, nil
}

func (t *transactionService) BuildTransactionReq(tdAccount *mambu.TDAccount,
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
			MessageType:                 t.common.Config.TransactionReqMetaData.MessageType,
			ExternalTransactionID:       transactionID,
			ExternalTransactionDetailID: transactionDetailID,
			TransactionType:             t.common.Config.TransactionReqMetaData.TransactionType,
			TransactionDateTime:         carbon.Now().String(),
			TerminalType:                t.common.Config.TransactionReqMetaData.TerminalType,
			TerminalID:                  t.common.Config.TransactionReqMetaData.TerminalID,
			TerminalLocation:            t.common.Config.TransactionReqMetaData.TerminalLocation,
			TerminalRRN:                 t.GenerationTerminalRRN(),
			ProductCode:                 t.common.Config.TransactionReqMetaData.ProductCode,
			AcquirerIID:                 t.common.Config.TransactionReqMetaData.AcquirerIID,
			ForwarderIID:                t.common.Config.TransactionReqMetaData.ForwarderIID,
			IssuerIID:                   t.common.Config.TransactionReqMetaData.IssuerIID,
			IssuerIName:                 t.common.Config.TransactionReqMetaData.IssuerIName,
			DestinationIID:              t.common.Config.TransactionReqMetaData.DestinationIID,
			SourceAccountNo:             tdAccount.ID,
			SourceAccountName:           tdAccount.Name,
			BeneficiaryAccountNo:        benefitAccount.ID,
			BeneficiaryAccountName:      benefitAccount.Name,
			Currency:                    t.common.Config.TransactionReqMetaData.Currency,
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

func (t *transactionService) GenerationTerminalRRN() string {
	return "TDE-" + util.RandomSnowFlakeId()
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
