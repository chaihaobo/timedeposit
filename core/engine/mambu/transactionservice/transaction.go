// Package transactionservice
// @author： Boice
// @createTime：2022/5/26 18:37
package transactionservice

import (
	"encoding/json"
	"errors"
	"fmt"
	"gitlab.com/bns-engineering/td/common/config"
	"gitlab.com/bns-engineering/td/common/constant"
	"gitlab.com/bns-engineering/td/common/util"
	"gitlab.com/bns-engineering/td/dao"
	"gitlab.com/bns-engineering/td/service/mambuEntity"
	"go.uber.org/zap"
	"time"
)

func WithdrawTransaction(tdAccount, benefitAccount *mambuEntity.TDAccount,
	amount float64,
	transactionID string,
	channelID string) (mambuEntity.TransactionResp, error) {

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

	postUrl := fmt.Sprintf(constant.WithdrawTransactiontUrl, tdAccount.ID)
	respBody, code, err := util.HttpPostData(postJsonStr, postUrl)
	if err != nil &&
		code != constant.HttpStatusCodeSucceed &&
		code != constant.HttpStatusCodeSucceedNoContent &&
		code != constant.HttpStatusCodeSucceedCreate {
		zap.L().Error(fmt.Sprintf("Withdraw Transaction Error! td acc id: %v, error:%v", tdAccount.ID, respBody))
		dao.CreateFailedTransaction(tmpTransaction, constant.TransactionWithdraw, err.Error())
		return transactionResp, errors.New(respBody)
	}

	zap.L().Debug(fmt.Sprintf("Withdraw Transaction for td account succeed. Result: %v", respBody))
	err = json.Unmarshal([]byte(respBody), &transactionResp)
	if err != nil {
		zap.L().Error(fmt.Sprintf("Convert Json to TransactionResp Failed. json: %v", respBody))
		dao.CreateFailedTransaction(tmpTransaction, constant.TransactionWithdraw, err.Error())
		return transactionResp, errors.New("mambu process Withdraw Transaction Error, the response data error")
	}
	dao.CreateSucceedFlowTransaction(transactionResp)
	return transactionResp, nil
}
func DepositTransaction(tdAccount, benefitAccount *mambuEntity.TDAccount,
	amount float64,
	transactionID string,
	channelID string) (mambuEntity.TransactionResp, error) {
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

	postUrl := fmt.Sprintf(constant.DepositTransactiontUrl, benefitAccount.ID)
	respBody, code, err := util.HttpPostData(postJsonStr, postUrl)
	if err != nil &&
		code != constant.HttpStatusCodeSucceed &&
		code != constant.HttpStatusCodeSucceedNoContent &&
		code != constant.HttpStatusCodeSucceedCreate {
		zap.L().Error(fmt.Sprintf("Deposit Transaction Error! td acc id: %v, error:%v", tdAccount.ID, respBody))
		dao.CreateFailedTransaction(tmpTransaction, constant.TransactionWithdraw, err.Error())
		return transactionResp, errors.New(respBody)
	}

	zap.L().Debug(fmt.Sprintf("Deposit Transaction for td account succeed. Result: %v", respBody))
	err = json.Unmarshal([]byte(respBody), &transactionResp)
	if err != nil {
		zap.L().Error(fmt.Sprintf("Convert Json to TransactionResp Failed. json: %v", respBody))
		dao.CreateFailedTransaction(tmpTransaction, constant.TransactionWithdraw, err.Error())
		return transactionResp, errors.New("mambu process Deposit Transaction Error, the response data error")
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
			ExternalOriTransactionID:       config.TDConf.TransactionReqMetaData.ExternalOriTransactionID,
			ExternalOriTransactionDetailID: config.TDConf.TransactionReqMetaData.ExternalOriTransactionDetailID,
			TransactionType:                config.TDConf.TransactionReqMetaData.TransactionType,
			TransactionDateTime:            time.Now().Format("2006-01-02 15:04:05"),
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
