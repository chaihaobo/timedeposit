/*
 * @Author: Hugo
 * @Date: 2022-05-11 12:21:05
 * @Last Modified by: Hugo
 * @Last Modified time: 2022-05-23 02:46:56
 */
package mambuservices

import (
	"encoding/json"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"time"

	"gitlab.com/bns-engineering/td/common/constant"
	"gitlab.com/bns-engineering/td/common/util"
	mambuEntity "gitlab.com/bns-engineering/td/service/mambuEntity"
)

// Get TDAccount Info from mambu api
func GetTDAccountById(tdAccountID string) (mambuEntity.TDAccount, error) {
	var tdAccount mambuEntity.TDAccount
	getUrl := fmt.Sprintf(constant.GetTDAccountUrl, tdAccountID)
	zap.L().Info(fmt.Sprintf("getUrl: %v", getUrl))
	resp, code, err := util.HttpGetData(getUrl)
	if err != nil || code != constant.HttpStatusCodeSucceed {
		zap.L().Error(fmt.Sprintf("Query td account Info failed! td acc id: %v", tdAccountID))
		return tdAccount, err
	}
	zap.L().Info(fmt.Sprintf("Query td account Info result: %v", resp))
	err = json.Unmarshal([]byte(resp), &tdAccount)
	if err != nil {
		zap.L().Error(fmt.Sprintf("Convert Json to TDAccount Failed. json: %v, err:%v", resp, err.Error()))
		return tdAccount, err
	}
	return tdAccount, nil
}

// Get TDAccount List from mambu api
func GetTDAccountListByQueryParam(searchParam mambuEntity.SearchParam) ([]mambuEntity.TDAccount, error) {
	tdAccountList := []mambuEntity.TDAccount{}
	postUrl := constant.SearchTDAccountListUrl
	zap.L().Debug(fmt.Sprintf("postUrl: %v", postUrl))
	queryParamByte, err := json.Marshal(searchParam)
	if err != nil {
		zap.L().Error(fmt.Sprintf("Convert searchParam to JsonStr Failed. searchParam: %v", searchParam))
		return tdAccountList, nil
	}
	postJsonStr := string(queryParamByte)
	zap.L().Debug(fmt.Sprintf("PostUrl:%v", postUrl))
	zap.L().Debug(fmt.Sprintf("postJsonStr:%v", postJsonStr))
	resp, code, err := util.HttpPostData(postJsonStr, postUrl)
	zap.L().Debug(fmt.Sprintf("responseStr:%v", resp))

	if err != nil || code != constant.HttpStatusCodeSucceed {
		zap.L().Error(fmt.Sprintf("Search td account Info List failed! queryParam: %v", postJsonStr))
		return tdAccountList, err
	}
	zap.L().Debug(fmt.Sprintf("Query td account Info result: %v", resp))
	err = json.Unmarshal([]byte(resp), &tdAccountList)
	if err != nil {
		zap.L().Error(fmt.Sprintf("Convert Json to TDAccount Failed. json: %v", resp))
		return tdAccountList, err
	}
	return tdAccountList, nil
}

// Disable MaturityDate
func UndoMaturityDate(accountID string) bool {
	postUrl := fmt.Sprintf(constant.UndoMaturityDateUrl, accountID)
	zap.L().Info(fmt.Sprintf("getUrl: %v", postUrl))
	_, code, err := util.HttpPostData("", postUrl)
	if err != nil && code != constant.HttpStatusCodeSucceed && code != constant.HttpStatusCodeSucceedNoContent {
		zap.L().Error(fmt.Sprintf("Undo MaturityDate for td account failed! td acc id: %v", accountID))
		return false
	}
	return true
}

//Create New Maturity Date for this TD account
func ChangeMaturityDate(accountID, maturityDate, note string) (mambuEntity.TDAccount, error) {
	postUrl := fmt.Sprintf(constant.StartMaturityDateUrl, accountID)
	zap.L().Info(fmt.Sprintf("StartMaturityDateUrl: %v", postUrl))

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
	if err != nil &&
		code != constant.HttpStatusCodeSucceed &&
		code != constant.HttpStatusCodeSucceedNoContent &&
		code != constant.HttpStatusCodeSucceedCreate {
		zap.L().Error(fmt.Sprintf("Create MaturityDate for td account failed! td acc id: %v, error:%v", accountID, respBody))
		return resultTDAccount, errors.New("mambu process StartMaturityDate Error")
	}

	zap.L().Debug(fmt.Sprintf("Create MaturityDate for td account succeed. Result: %v", respBody))
	err = json.Unmarshal([]byte(respBody), &resultTDAccount)
	if err != nil {
		zap.L().Error(fmt.Sprintf("Convert Json to TDAccount Failed. json: %v", respBody))
		return resultTDAccount, errors.New("mambu process StartMaturityDate Error, the response data error")
	}
	return resultTDAccount, nil
}

// Apply profit
func ApplyProfit(accountID, note string) bool {
	postUrl := fmt.Sprintf(constant.ApplyProfitUrl, accountID)
	zap.L().Debug(fmt.Sprintf("applyProfitUrl: %v", postUrl))

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
	if err != nil && code != constant.HttpStatusCodeSucceed && code != constant.HttpStatusCodeSucceedNoContent {
		zap.L().Error(fmt.Sprintf("Undo MaturityDate for td account failed! td acc id: %v", accountID))
		return false
	}
	return true
}

func UpdateMaturifyDateForTDAccount(accountID, newDate string) bool {
	postUrl := fmt.Sprintf(constant.ApplyProfitUrl, accountID)
	zap.L().Debug(fmt.Sprintf("applyProfitUrl: %v", postUrl))

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
	if err != nil &&
		code != constant.HttpStatusCodeSucceed &&
		code != constant.HttpStatusCodeSucceedNoContent {
		zap.L().Error(fmt.Sprintf("Undo MaturityDate for td account failed! td acc id: %v", accountID))
		return false
	}
	return true
}

func CloseAccount(accID, notes string) bool {
	postUrl := fmt.Sprintf(constant.CloseAccountUrl, accID)
	zap.L().Debug(fmt.Sprintf("CloseAccountUrl: %v", postUrl))

	//Build the update maturity json struct
	postJsonByte, _ := json.Marshal([]struct {
		Action string `json:"action"`
		Notes  string `json:"notes"`
	}{
		{
			Action: "CLOSE",
			Notes:  notes,
		},
	})

	postJsonStr := string(postJsonByte)

	_, code, err := util.HttpPatchData(postJsonStr, postUrl)
	if err != nil &&
		code != constant.HttpStatusCodeSucceed &&
		code != constant.HttpStatusCodeSucceedNoContent {
		zap.L().Error(fmt.Sprintf("Undo MaturityDate for td account failed! td acc id: %v", accID))
		return false
	}
	return true
}
