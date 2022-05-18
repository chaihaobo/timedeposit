/*
 * @Author: Hugo
 * @Date: 2022-05-11 12:21:05
 * @Last Modified by: Hugo
 * @Last Modified time: 2022-05-18 05:19:39
 */
package mambuservices

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"gitlab.com/hugo.hu/time-deposit-eod-engine/common/constant"
	commonLog "gitlab.com/hugo.hu/time-deposit-eod-engine/common/log"
	"gitlab.com/hugo.hu/time-deposit-eod-engine/common/util"
	mambuEntity "gitlab.com/hugo.hu/time-deposit-eod-engine/service/mambuEntity"
)

// Get TDAccount Info from mambu api
func GetTDAccountById(tdAccountID string) (mambuEntity.TDAccount, error) {
	var tdAccount mambuEntity.TDAccount
	getUrl := fmt.Sprintf(constant.GetTDAccountUrl, tdAccountID)
	commonLog.Log.Info("getUrl: %v", getUrl)
	resp, code, err := util.HttpGetData(getUrl)
	if err != nil || code != constant.HttpStatusCodeSucceed {
		commonLog.Log.Error("Query td account Info failed! td acc id: %v", tdAccountID)
		return tdAccount, err
	}
	commonLog.Log.Info("Query td account Info result: %v", resp)
	err = json.Unmarshal([]byte(resp), &tdAccount)
	if err != nil {
		commonLog.Log.Error("Convert Json to TDAccount Failed. json: %v", resp)
		return tdAccount, err
	}
	return tdAccount, nil
}

// Get TDAccount List from mambu api
func GetTDAccountListByQueryParam(searchParam mambuEntity.SearchParam) ([]mambuEntity.TDAccount, error) {
	tdAccountList := []mambuEntity.TDAccount{}
	postUrl := constant.SearchTDAccountListUrl
	commonLog.Log.Debug("postUrl: %v", postUrl)
	queryParamByte, err := json.Marshal(searchParam)
	if err != nil {
		commonLog.Log.Error("Convert searchParam to JsonStr Failed. searchParam: %v", searchParam)
		return tdAccountList, nil
	}
	postJsonStr := string(queryParamByte)
	commonLog.Log.Debug("PostUrl:%v", postUrl)
	commonLog.Log.Debug("postJsonStr:%v", postJsonStr)
	resp, code, err := util.HttpPostData(postJsonStr, postUrl)
	commonLog.Log.Debug("responseStr:%v", resp)

	if err != nil || code != constant.HttpStatusCodeSucceed {
		commonLog.Log.Error("Search td account Info List failed! queryParam: %v", postJsonStr)
		return tdAccountList, err
	}
	commonLog.Log.Debug("Query td account Info result: %v", resp)
	err = json.Unmarshal([]byte(resp), &tdAccountList)
	if err != nil {
		commonLog.Log.Error("Convert Json to TDAccount Failed. json: %v", resp)
		return tdAccountList, err
	}
	return tdAccountList, nil
}

// Disable MaturityDate
func UndoMaturityDate(accountID string) bool {
	postUrl := fmt.Sprintf(constant.UndoMaturityDateUrl, accountID)
	commonLog.Log.Info("getUrl: %v", postUrl)
	_, code, err := util.HttpPostData("", postUrl)
	if err != nil || code != constant.HttpStatusCodeSucceed || code != constant.HttpStatusCodeSucceedNoContent {
		commonLog.Log.Error("Undo MaturityDate for td account failed! td acc id: %v", accountID)
		return false
	}
	return true
}

//Create New Maturity Date for this TD account
func ChangeMaturityDate(accountID, maturityDate, note string) (mambuEntity.TDAccount, error) {
	postUrl := fmt.Sprintf(constant.StartMaturityDateUrl, accountID)
	commonLog.Log.Info("StartMaturityDateUrl: %v", postUrl)

	//Build the update maturity json struct
	postJsonByte, _ := json.Marshal(struct {
		MaturityDate string `json:"maturityDate"`
		Notes        string `json:"notes"`
	}{
		MaturityDate: maturityDate,
		Notes:        note,
	})
	postJsonStr := string(postJsonByte)

	var resultTDAccount mambuEntity.TDAccount

	respBody, code, err := util.HttpPostData(postJsonStr, postUrl)
	if err != nil || code != constant.HttpStatusCodeSucceed || code != constant.HttpStatusCodeSucceedNoContent {
		commonLog.Log.Error("Create MaturityDate for td account failed! td acc id: %v, error:%v", accountID, respBody)
		return resultTDAccount, errors.New("mambu process StartMaturityDate Error")
	}

	commonLog.Log.Debug("Create MaturityDate for td account succeed. Result: %v", respBody)
	err = json.Unmarshal([]byte(respBody), &resultTDAccount)
	if err != nil {
		commonLog.Log.Error("Convert Json to TDAccount Failed. json: %v", respBody)
		return resultTDAccount, errors.New("mambu process StartMaturityDate Error, the response data error")
	}
	return resultTDAccount, nil
}

// Apply profit
func ApplyProfit(accountID, note string) bool {
	postUrl := fmt.Sprintf(constant.ApplyProfitUrl, accountID)
	commonLog.Log.Debug("applyProfitUrl: %v", postUrl)

	//Build the update maturity json struct
	postJsonByte, _ := json.Marshal(struct {
		InterestApplicationDate time.Time `json:"interestApplicationDate"`
		Notes                   string    `json:"notes"`
	}{
		InterestApplicationDate: time.Now(),
		Notes:                   note,
	})
	postJsonStr := string(postJsonByte)

	_, code, err := util.HttpPostData(postJsonStr, postUrl)
	if err != nil || code != constant.HttpStatusCodeSucceed || code != constant.HttpStatusCodeSucceedNoContent {
		commonLog.Log.Error("Undo MaturityDate for td account failed! td acc id: %v", accountID)
		return false
	}
	return true
}

func UpdateMaturifyDateForTDAccount(accountID, newDate string) bool {
	postUrl := fmt.Sprintf(constant.ApplyProfitUrl, accountID)
	commonLog.Log.Debug("applyProfitUrl: %v", postUrl)

	//Build the update maturity json struct
	postJsonByte, _ := json.Marshal([]struct {
		Op    string `json:"op"`
		Path  string `json:"path"`
		Value string `json:"value"`
	}{
		{
			Op:    "REPLACE",
			Path:  "/_rekening/rekeningTanggalJatohTempo",
			Value: newDate,
		},
	})

	postJsonStr := string(postJsonByte)

	_, code, err := util.HttpPatchData(postJsonStr, postUrl)
	if err != nil || code != constant.HttpStatusCodeSucceed || code != constant.HttpStatusCodeSucceedNoContent {
		commonLog.Log.Error("Undo MaturityDate for td account failed! td acc id: %v", accountID)
		return false
	}
	return true
}
