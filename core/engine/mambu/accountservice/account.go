// Package accountservice
// @author： Boice
// @createTime：2022/5/26 17:26
package accountservice

import (
	"encoding/json"
	"fmt"
	"gitlab.com/bns-engineering/td/common/constant"
	"gitlab.com/bns-engineering/td/common/util"
	"gitlab.com/bns-engineering/td/service/mambuEntity"
	"go.uber.org/zap"
)

func GetAccountById(tdAccountID string) (*mambuEntity.TDAccount, error) {
	var tdAccount = new(mambuEntity.TDAccount)
	getUrl := fmt.Sprintf(constant.GetTDAccountUrl, tdAccountID)
	resp, code, err := util.HttpGetData(getUrl)
	if err != nil || code != constant.HttpStatusCodeSucceed {
		zap.L().Error(fmt.Sprintf("Query td account Info failed! td acc id: %v", tdAccountID, zap.Error(err)))
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
