// Package accountservice
// @author： Boice
// @createTime：2022/5/26 17:26
package accountservice

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/shopspring/decimal"
	"github.com/uniplaces/carbon"
	"gitlab.com/bns-engineering/td/common/config"
	"gitlab.com/bns-engineering/td/common/constant"
	"gitlab.com/bns-engineering/td/common/util/mambu_http"
	"gitlab.com/bns-engineering/td/common/util/mambu_http/persistence"
	"gitlab.com/bns-engineering/td/model/mambu"
	"go.uber.org/zap"
	"time"
)

// GetTDAccountListByQueryParam Get TDAccount List from mambu api
func GetTDAccountListByQueryParam(searchParam mambu.SearchParam) ([]mambu.TDAccount, error) {
	var tdAccountList []mambu.TDAccount
	postUrl := fmt.Sprintf(constant.UrlOf(constant.SearchTDAccountListUrl), 0, config.TDConf.Flow.MaxLimitSearchAccount)
	zap.L().Debug(fmt.Sprintf("postUrl: %v", postUrl))
	queryParamByte, err := json.Marshal(searchParam)
	if err != nil {
		zap.L().Error(fmt.Sprintf("Convert searchParam to JsonStr Failed. searchParam: %v", searchParam))
		return tdAccountList, nil
	}
	postJsonStr := string(queryParamByte)
	zap.L().Debug(fmt.Sprintf("PostUrl:%v", postUrl))
	zap.L().Debug(fmt.Sprintf("postJsonStr:%v", postJsonStr))
	type RspHeader struct {
		Total int32 `header:"items-total"`
	}
	responseHeader := new(RspHeader)
	err = mambu_http.Post(postUrl, postJsonStr, &tdAccountList, responseHeader, persistence.DBPersistence(nil, "GetTDAccountListByQueryParam"))
	if err != nil {
		zap.L().Error(fmt.Sprintf("Search td account Info List failed! queryParam: %v", postJsonStr))
		return tdAccountList, err
	}

	if responseHeader.Total > config.TDConf.Flow.MaxLimitSearchAccount {
		zap.L().Info("Total has gt configured limit, Load more..")
		// get all accounts
		pages := int(decimal.NewFromInt32(responseHeader.Total).Div(decimal.NewFromInt32(config.TDConf.Flow.MaxLimitSearchAccount)).Ceil().IntPart())
		for i := 1; i <= pages-1; i++ {
			var moreAccountList []mambu.TDAccount
			offset := i * int(config.TDConf.Flow.MaxLimitSearchAccount)
			zap.L().Info("get more td account", zap.Int("offset", offset))
			postUrl = fmt.Sprintf(constant.UrlOf(constant.SearchTDAccountListUrl), offset, config.TDConf.Flow.MaxLimitSearchAccount)
			err = mambu_http.Post(postUrl, postJsonStr, &moreAccountList, responseHeader, persistence.DBPersistence(nil, "GetMoreTDAccountListByQueryParam"))
			if err == nil {
				tdAccountList = append(tdAccountList, moreAccountList...)
			}
		}
	}

	return tdAccountList, nil
}

func GetAccountById(context context.Context, tdAccountID string) (*mambu.TDAccount, error) {
	var tdAccount = new(mambu.TDAccount)
	getUrl := fmt.Sprintf(constant.UrlOf(constant.GetTDAccountUrl), tdAccountID)
	err := mambu_http.Get(getUrl, tdAccount, persistence.DBPersistence(context, "GetAccountById"))
	if err != nil {
		zap.L().Error("get account Failed", zap.Error(err))
		return nil, err
	}
	return tdAccount, nil
}

func UndoMaturityDate(context context.Context, accountID string) bool {
	postUrl := fmt.Sprintf(constant.UrlOf(constant.UndoMaturityDateUrl), accountID)
	err := mambu_http.Post(postUrl, "", nil, nil, persistence.DBPersistence(context, "UndoMaturityDate"))
	if err != nil {
		zap.L().Error(fmt.Sprintf("Undo MaturityDate for td account failed! td acc id: %v", accountID))
		return false
	}
	return true
}

// ChangeMaturityDate Create New Maturity Date for this TD account
func ChangeMaturityDate(context context.Context, accountID, maturityDate, note string) (mambu.TDAccount, error) {
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

	var resultTDAccount mambu.TDAccount
	err := mambu_http.Post(postUrl, postJsonStr, &resultTDAccount, nil, persistence.DBPersistence(context, "ChangeMaturityDate"))
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
	err := mambu_http.Post(postUrl, postJsonStr, nil, nil, persistence.DBPersistence(context, "ApplyProfit"))
	if err != nil {
		zap.L().Error(fmt.Sprintf("ApplyProfit for td account failed! td acc id: %v", accountID))
		return false
	}
	return true
}

func UpdateMaturifyDateForTDAccount(context context.Context, accountID, newDate string) bool {
	postUrl := fmt.Sprintf(constant.UrlOf(constant.UpdateTDAccountUrl), accountID)
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
	err := mambu_http.Patch(postUrl, postJsonStr, nil, persistence.DBPersistence(context, "UpdateMaturifyDateForTDAccount"))
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
	postJsonByte, _ := json.Marshal(struct {
		Action string `json:"action"`
		Notes  string `json:"notes"`
	}{
		Action: "CLOSE",
		Notes:  notes,
	})

	postJsonStr := string(postJsonByte)
	err := mambu_http.Post(postUrl, postJsonStr, nil, nil, persistence.DBPersistence(context, "CloseAccount"))
	if err != nil {
		zap.L().Error(fmt.Sprintf("CloseAccount for td account failed! td acc id: %v", accID))
		return false
	}
	return true
}
