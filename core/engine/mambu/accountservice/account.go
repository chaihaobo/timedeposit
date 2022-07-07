// Package accountservice
// @author： Boice
// @createTime：2022/5/26 17:26
package accountservice

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/shopspring/decimal"
	"gitlab.com/bns-engineering/common/tracer"
	"gitlab.com/bns-engineering/td/common/config"
	"gitlab.com/bns-engineering/td/common/constant"
	"gitlab.com/bns-engineering/td/common/log"
	"gitlab.com/bns-engineering/td/common/util/mambu_http"
	"gitlab.com/bns-engineering/td/common/util/mambu_http/persistence"
	"gitlab.com/bns-engineering/td/model/mambu"
	"go.uber.org/zap"
	"time"
)

// GetTDAccountListByQueryParam Get TDAccount List from mambu api
func GetTDAccountListByQueryParam(ctx context.Context, searchParam mambu.SearchParam) ([]mambu.TDAccount, error) {
	tr := tracer.StartTrace(ctx, "accountService-GetTDAccountListByQueryParam")
	ctx = tr.Context()
	defer tr.Finish()
	var tdAccountList []mambu.TDAccount
	postUrl := fmt.Sprintf(constant.UrlOf(constant.SearchTDAccountListUrl), 0, config.TDConf.Flow.MaxLimitSearchAccount)

	queryParamByte, err := json.Marshal(searchParam)
	if err != nil {
		log.Error(ctx, fmt.Sprintf("Convert searchParam to JsonStr Failed. searchParam: %v", searchParam), err)
		return tdAccountList, nil
	}
	postJsonStr := string(queryParamByte)

	type RspHeader struct {
		Total int32 `header:"items-total"`
	}
	responseHeader := new(RspHeader)
	err = mambu_http.Post(ctx, postUrl, postJsonStr, &tdAccountList, responseHeader, persistence.DBPersistence(nil, "GetTDAccountListByQueryParam"))
	if err != nil {
		log.Error(ctx, fmt.Sprintf("Search td account Info List failed! queryParam: %v", postJsonStr), err)
		return tdAccountList, err
	}

	if responseHeader.Total > config.TDConf.Flow.MaxLimitSearchAccount {
		log.Info(ctx, "Total has gt configured limit, Load more..")
		// get all accounts
		pages := int(decimal.NewFromInt32(responseHeader.Total).Div(decimal.NewFromInt32(config.TDConf.Flow.MaxLimitSearchAccount)).Ceil().IntPart())
		for i := 1; i <= pages-1; i++ {
			var moreAccountList []mambu.TDAccount
			offset := i * int(config.TDConf.Flow.MaxLimitSearchAccount)
			log.Info(ctx, "get more td account", zap.Int("offset", offset))
			postUrl = fmt.Sprintf(constant.UrlOf(constant.SearchTDAccountListUrl), offset, config.TDConf.Flow.MaxLimitSearchAccount)
			err = mambu_http.Post(ctx, postUrl, postJsonStr, &moreAccountList, responseHeader, persistence.DBPersistence(nil, "GetMoreTDAccountListByQueryParam"))
			if err == nil {
				tdAccountList = append(tdAccountList, moreAccountList...)
			}
		}
	}

	return tdAccountList, nil
}

func GetAccountById(context context.Context, tdAccountID string) (*mambu.TDAccount, error) {
	tr := tracer.StartTrace(context, "accountService-GetAccountById")
	context = tr.Context()
	defer tr.Finish()
	var tdAccount = new(mambu.TDAccount)
	getUrl := fmt.Sprintf(constant.UrlOf(constant.GetTDAccountUrl), tdAccountID)
	err := mambu_http.Get(context, getUrl, tdAccount, persistence.DBPersistence(context, "GetAccountById"))
	if err != nil {
		log.Error(context, "get account Failed", err)
		return nil, err
	}
	return tdAccount, nil
}

func UndoMaturityDate(context context.Context, accountID string) bool {
	tr := tracer.StartTrace(context, "accountService-UndoMaturityDate")
	context = tr.Context()
	defer tr.Finish()

	postUrl := fmt.Sprintf(constant.UrlOf(constant.UndoMaturityDateUrl), accountID)
	err := mambu_http.Post(context, postUrl, "", nil, nil, persistence.DBPersistence(context, "UndoMaturityDate"))
	if err != nil {
		log.Error(context, fmt.Sprintf("Undo MaturityDate for td account failed! td acc id: %v", accountID), err)
		return false
	}
	return true
}

// ChangeMaturityDate Create New Maturity Date for this TD account
func ChangeMaturityDate(ctx context.Context, accountID, maturityDate, note string) (mambu.TDAccount, error) {
	tr := tracer.StartTrace(ctx, "accountService-ChangeMaturityDate")
	ctx = tr.Context()
	defer tr.Finish()

	postUrl := fmt.Sprintf(constant.UrlOf(constant.StartMaturityDateUrl), accountID)
	log.Info(ctx, fmt.Sprintf("StartMaturityDateUrl: %v", postUrl))

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
	err := mambu_http.Post(ctx, postUrl, postJsonStr, &resultTDAccount, nil, persistence.DBPersistence(ctx, "ChangeMaturityDate"))
	if err != nil {
		log.Error(ctx, "Create MaturityDate for td account failed! ", err, zap.String("accountId", accountID))
		return resultTDAccount, err
	}
	return resultTDAccount, nil
}

func ApplyProfit(context context.Context, accountID, note string, interestApplicationDate time.Time) bool {
	tr := tracer.StartTrace(context, "accountService-ApplyProfit")
	context = tr.Context()
	defer tr.Finish()

	postUrl := fmt.Sprintf(constant.UrlOf(constant.ApplyProfitUrl), accountID)

	// Build the update maturity json struct

	postJsonByte, _ := json.Marshal(struct {
		InterestApplicationDate time.Time `json:"interestApplicationDate"`
		Notes                   string    `json:"notes"`
	}{
		InterestApplicationDate: interestApplicationDate,
		Notes:                   note,
	})
	postJsonStr := string(postJsonByte)
	err := mambu_http.Post(context, postUrl, postJsonStr, nil, nil, persistence.DBPersistence(context, "ApplyProfit"))
	if err != nil {
		log.Error(context, fmt.Sprintf("ApplyProfit for td account failed! td acc id: %v", accountID), err)
		return false
	}
	return true
}

func UpdateMaturifyDateForTDAccount(context context.Context, accountID, newDate string) bool {
	tr := tracer.StartTrace(context, "accountService-UpdateMaturifyDateForTDAccount")
	context = tr.Context()
	defer tr.Finish()

	postUrl := fmt.Sprintf(constant.UrlOf(constant.UpdateTDAccountUrl), accountID)

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
	err := mambu_http.Patch(context, postUrl, postJsonStr, nil, persistence.DBPersistence(context, "UpdateMaturifyDateForTDAccount"))
	if err != nil {
		log.Error(context, fmt.Sprintf("Undo MaturityDate for td account failed! td acc id: %v", accountID), err)
		return false
	}
	return true
}

func CloseAccount(context context.Context, accID, notes string) bool {
	tr := tracer.StartTrace(context, "accountService-CloseAccount")
	context = tr.Context()
	defer tr.Finish()
	postUrl := fmt.Sprintf(constant.UrlOf(constant.CloseAccountUrl), accID)

	// Build the update maturity json struct
	postJsonByte, _ := json.Marshal(struct {
		Action string `json:"action"`
		Notes  string `json:"notes"`
	}{
		Action: "CLOSE",
		Notes:  notes,
	})

	postJsonStr := string(postJsonByte)
	err := mambu_http.Post(context, postUrl, postJsonStr, nil, nil, persistence.DBPersistence(context, "CloseAccount"))
	if err != nil {
		log.Error(context, fmt.Sprintf("CloseAccount for td account failed! td acc id: %v", accID), err)
		return false
	}
	return true
}
