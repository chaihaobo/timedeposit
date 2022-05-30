// Package accountservice
// @author： Boice
// @createTime：2022/5/26 17:26
package accountservice

import (
	"encoding/json"
	"errors"
	"fmt"
	"gitlab.com/bns-engineering/td/common/constant"
	"gitlab.com/bns-engineering/td/common/util"
	"gitlab.com/bns-engineering/td/common/util/http"
	"gitlab.com/bns-engineering/td/service/mambuEntity"
	"go.uber.org/zap"
	"time"
)

func GetAccountById(tdAccountID string) (*mambuEntity.TDAccount, error) {
	var tdAccount = new(mambuEntity.TDAccount)
	getUrl := fmt.Sprintf(constant.GetTDAccountUrl, tdAccountID)
	resp, code, err := util.HttpGetData(getUrl)
	if err != nil || code != constant.HttpStatusCodeSucceed {
		zap.L().Error(fmt.Sprintf("Query td account Info failed! td acc id: %v", tdAccountID), zap.String("body", resp), zap.Error(err))
		if err == nil {
			err = errors.New("query account status code is not succeed")
		}
		return nil, err
	}
	zap.L().Info(fmt.Sprintf("Query td account Info result: %v", resp))
	err = json.Unmarshal([]byte(resp), &tdAccount)
	if err != nil {
		zap.L().Error(fmt.Sprintf("Convert Json to TDAccount Failed. json: %v, err:%v", resp, err.Error()))
		return nil, err
	}
	return tdAccount, nil
}

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

// Create New Maturity Date for this TD account
func ChangeMaturityDate(accountID, maturityDate, note string) (mambuEntity.TDAccount, error) {
	postUrl := fmt.Sprintf(constant.StartMaturityDateUrl, accountID)
	zap.L().Info(fmt.Sprintf("StartMaturityDateUrl: %v", postUrl))

	// Build the update maturity json struct
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

func ApplyProfit(accountID, note string) bool {
	postUrl := fmt.Sprintf(constant.ApplyProfitUrl, accountID)
	zap.L().Debug(fmt.Sprintf("applyProfitUrl: %v", postUrl))

	// Build the update maturity json struct
	postJsonByte, _ := json.Marshal(struct {
		InterestApplicationDate time.Time `json:"interestApplicationDate"`
		Notes                   string    `json:"notes"`
	}{
		InterestApplicationDate: time.Now().In(time.FixedZone("CST", 7*3600)),
		Notes:                   note,
	})
	postJsonStr := string(postJsonByte)
	err := http.Post(postUrl, postJsonStr, nil, nil)
	if err != nil {
		zap.L().Error(fmt.Sprintf("ApplyProfit for td account failed! td acc id: %v", accountID))
		return false
	}
	return true
}

func UpdateMaturifyDateForTDAccount(accountID, newDate string) bool {
	postUrl := fmt.Sprintf(constant.ApplyProfitUrl, accountID)
	zap.L().Debug(fmt.Sprintf("applyProfitUrl: %v", postUrl))

	// Build the update maturity json struct
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

	// Build the update maturity json struct
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
