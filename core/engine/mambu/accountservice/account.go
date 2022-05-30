// Package accountservice
// @author： Boice
// @createTime：2022/5/26 17:26
package accountservice

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/uniplaces/carbon"
	"gitlab.com/bns-engineering/td/common/constant"
	"gitlab.com/bns-engineering/td/common/util/mambu_http"
	"gitlab.com/bns-engineering/td/core/engine/mambu"
	"gitlab.com/bns-engineering/td/service/mambuEntity"
	"go.uber.org/zap"
	"time"
)

func GetAccountById(context context.Context, tdAccountID string) (*mambuEntity.TDAccount, error) {
	var tdAccount = new(mambuEntity.TDAccount)
	getUrl := fmt.Sprintf(constant.UrlOf(constant.GetTDAccountUrl), tdAccountID)
	err := mambu_http.Get(getUrl, tdAccount, mambu.SaveMambuRequestLog(context, "GetAccountById"))
	if err != nil {
		zap.L().Error("get account Failed", zap.Error(err))
		return nil, err
	}
	return tdAccount, nil
}

func UndoMaturityDate(context context.Context, accountID string) bool {
	postUrl := fmt.Sprintf(constant.UrlOf(constant.UndoMaturityDateUrl), accountID)
	zap.L().Info(fmt.Sprintf("getUrl: %v", postUrl))
	err := mambu_http.Post(postUrl, "", nil, mambu.SaveMambuRequestLog(context, "UndoMaturityDate"))
	if err != nil {
		zap.L().Error(fmt.Sprintf("Undo MaturityDate for td account failed! td acc id: %v", accountID))
		return false
	}
	return true
}

// Create New Maturity Date for this TD account
func ChangeMaturityDate(context context.Context, accountID, maturityDate, note string) (mambuEntity.TDAccount, error) {
	postUrl := fmt.Sprintf(constant.UrlOf(constant.StartMaturityDateUrl), accountID)
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
	err := mambu_http.Post(postUrl, postJsonStr, &resultTDAccount, mambu.SaveMambuRequestLog(context, "ChangeMaturityDate"))
	if err != nil {
		zap.L().Error("Create MaturityDate for td account failed! ", zap.String("accountId", accountID))
		return resultTDAccount, err
	}
	return resultTDAccount, nil
}

func ApplyProfit(context context.Context, accountID, note string) bool {
	postUrl := fmt.Sprintf(constant.UrlOf(constant.ApplyProfitUrl), accountID)
	zap.L().Debug(fmt.Sprintf("applyProfitUrl: %v", postUrl))
	// Build the update maturity json struct

	postJsonByte, _ := json.Marshal(struct {
		InterestApplicationDate time.Time `json:"interestApplicationDate"`
		Notes                   string    `json:"notes"`
	}{
		InterestApplicationDate: carbon.NewCarbon(time.Now().In(time.FixedZone("CST", 7*3600))).StartOfDay().Time,
		Notes:                   note,
	})
	postJsonStr := string(postJsonByte)
	err := mambu_http.Post(postUrl, postJsonStr, nil, mambu.SaveMambuRequestLog(context, "ApplyProfit"))
	if err != nil {
		zap.L().Error(fmt.Sprintf("ApplyProfit for td account failed! td acc id: %v", accountID))
		return false
	}
	return true
}

func UpdateMaturifyDateForTDAccount(context context.Context, accountID, newDate string) bool {
	postUrl := fmt.Sprintf(constant.UrlOf(constant.ApplyProfitUrl), accountID)
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
	err := mambu_http.Patch(postUrl, postJsonStr, nil, mambu.SaveMambuRequestLog(context, "UpdateMaturifyDateForTDAccount"))
	if err != nil {
		zap.L().Error(fmt.Sprintf("Undo MaturityDate for td account failed! td acc id: %v", accountID))
		return false
	}
	return true
}

func CloseAccount(context context.Context, accID, notes string) bool {
	postUrl := fmt.Sprintf(constant.UrlOf(constant.CloseAccountUrl), accID)
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
	err := mambu_http.Patch(postUrl, postJsonStr, nil, mambu.SaveMambuRequestLog(context, "CloseAccount"))
	if err != nil {
		zap.L().Error(fmt.Sprintf("Undo MaturityDate for td account failed! td acc id: %v", accID))
		return false
	}
	return true
}
