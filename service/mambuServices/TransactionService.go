/*
 * @Author: Hugo
 * @Date: 2022-05-12 10:48:09
 * @Last Modified by: Hugo
 * @Last Modified time: 2022-05-18 03:52:31
 */
package mambuservices

import (
	"encoding/json"
	"fmt"

	"gitlab.com/hugo.hu/time-deposit-eod-engine/common/constant"
	commonLog "gitlab.com/hugo.hu/time-deposit-eod-engine/common/log"
	"gitlab.com/hugo.hu/time-deposit-eod-engine/common/util"
	mambuEntity "gitlab.com/hugo.hu/time-deposit-eod-engine/service/mambuEntity"
)

// // Get Transaction Info from mambu api with key of account
// func GetTransactionByAccountKey(accountKey string) (entity.TransactionInfo, error) {
// 	var tdAccount mambuEntity.TDAccount
// 	getUrl := fmt.Sprintf(constant.GetTransactionUrl, accountKey)
// 	commonLog.Log.Info("get Transaction Url: %v", getUrl)
// 	resp, err := util.HttpPostData(getUrl)
// 	if err != nil {
// 		commonLog.Log.Error("Query td account Info failed! td acc id: %v", tdAccountID)
// 		return tdAccount, err
// 	}
// 	commonLog.Log.Info("Query td account Info result: %v", resp)
// 	err = json.Unmarshal([]byte(resp), &tdAccount)
// 	if err != nil {
// 		commonLog.Log.Error("Convert Json to TDAccount Failed. json: %v", resp)
// 		return tdAccount, err
// 	}
// 	return tdAccount, nil
// }

func WithdrawNetProfit(latestTDAccount, benefitAccount mambuEntity.TDAccount, netProfit float64) {
	tmpTransaction := mambuEntity.Transaction{
		Metadata: mambuEntity.Metadata{
			MessageType:                    "",
			ExternalTransactionID:          "",
			ExternalTransactionDetailID:    "",
			ExternalOriTransactionID:       "",
			ExternalOriTransactionDetailID: "",
			TransactionType:                "",
			TransactionDateTime:            "",
			TerminalType:                   "",
			TerminalID:                     "",
			TerminalLocation:               "",
			TerminalRRN:                    "",
			ProductCode:                    "",
			AcquirerIID:                    "",
			ForwarderIID:                   "",
			IssuerIID:                      "",
			IssuerName:                     "",
			DestinationIID:                 "",
			SourceAccountNo:                "",
			SourceAccountName:              "",
			BeneficiaryAccountNo:           "",
			BeneficiaryAccountName:         "",
			Currency:                       "",
			TranDesc1:                      "",
			TranDesc2:                      "",
			TranDesc3:                      "",
		},
		Amount: string(netProfit),
		TransactionDetails: mambuEntity.TransactionDetails{
			TransactionChannelID: "",
		},
	}
	queryParamByte, err := json.Marshal(tmpTransaction)
	if err != nil {
		commonLog.Log.Error("Convert searchParam to JsonStr Failed. searchParam: %v", searchParam)
		return
	}
	postJsonStr := string(queryParamByte)

	postUrl := fmt.Sprintf(constant.WithdrawAccountUrl, latestTDAccount.ID)
	util.HttpPostData(postJsonStr, postUrl)
	panic("Not implemented")
}

func DepositNetprofit(latestTDAccount, benefitAccount mambuEntity.TDAccount, netProfitTax float64) {
	panic("Not implemented")
}
