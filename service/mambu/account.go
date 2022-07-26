// Package service
// @author： Boice
// @createTime：2022/7/26 11:27
package mambu

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/shopspring/decimal"
	"gitlab.com/bns-engineering/common/tracer"
	"gitlab.com/bns-engineering/td/common"
	"gitlab.com/bns-engineering/td/common/constant"
	"gitlab.com/bns-engineering/td/model/mambu"
	"gitlab.com/bns-engineering/td/repository"
	"go.uber.org/zap"
	"time"
)

type AccountService interface {
	GetTDAccountListByQueryParam(ctx context.Context, searchParam mambu.SearchParam) ([]mambu.TDAccount, error)
	GetAccountById(context context.Context, tdAccountID string) (*mambu.TDAccount, error)
	UndoMaturityDate(context context.Context, accountID string) bool
	ChangeMaturityDate(ctx context.Context, accountID, maturityDate, note string) (mambu.TDAccount, error)
	ApplyProfit(context context.Context, accountID, note string, interestApplicationDate time.Time) bool
	UpdateMaturifyDateForTDAccount(context context.Context, accountID, newDate string) bool
	CloseAccount(context context.Context, accID, notes string) bool
}

func newAccountService(common *common.Common, repo *repository.Repository, mambuClient Client) AccountService {
	return &accountService{
		common:      common,
		repository:  repo,
		mambuClient: mambuClient,
	}
}

type accountService struct {
	common      *common.Common
	repository  *repository.Repository
	mambuClient Client
}

func (a *accountService) GetTDAccountListByQueryParam(ctx context.Context, searchParam mambu.SearchParam) ([]mambu.TDAccount, error) {
	tr := tracer.StartTrace(ctx, "accountService-GetTDAccountListByQueryParam")
	ctx = tr.Context()
	defer tr.Finish()
	var tdAccountList []mambu.TDAccount
	postUrl := fmt.Sprintf(constant.UrlOf(constant.SearchTDAccountListUrl), 0, a.common.Config.Flow.MaxLimitSearchAccount)

	queryParamByte, err := json.Marshal(searchParam)
	if err != nil {
		a.common.Logger.Error(ctx, fmt.Sprintf("Convert searchParam to JsonStr Failed. searchParam: %v", searchParam), err)
		return tdAccountList, nil
	}
	postJsonStr := string(queryParamByte)

	type RspHeader struct {
		Total int32 `header:"items-total"`
	}
	responseHeader := new(RspHeader)
	err = a.mambuClient.Post(ctx, postUrl, postJsonStr, &tdAccountList, responseHeader, a.mambuClient.DBPersistence(nil, "GetTDAccountListByQueryParam"))
	if err != nil {
		a.common.Logger.Error(ctx, fmt.Sprintf("Search td account Info List failed! queryParam: %v", postJsonStr), err)
		return tdAccountList, err
	}

	if responseHeader.Total > a.common.Config.Flow.MaxLimitSearchAccount {
		a.common.Logger.Info(ctx, "Total has gt configured limit, Load more..")
		// get all accounts
		pages := int(decimal.NewFromInt32(responseHeader.Total).Div(decimal.NewFromInt32(a.common.Config.Flow.MaxLimitSearchAccount)).Ceil().IntPart())
		for i := 1; i <= pages-1; i++ {
			var moreAccountList []mambu.TDAccount
			offset := i * int(a.common.Config.Flow.MaxLimitSearchAccount)
			a.common.Logger.Info(ctx, "get more td account", zap.Int("offset", offset))
			postUrl = fmt.Sprintf(constant.UrlOf(constant.SearchTDAccountListUrl), offset, a.common.Config.Flow.MaxLimitSearchAccount)
			err = a.mambuClient.Post(ctx, postUrl, postJsonStr, &moreAccountList, responseHeader, a.mambuClient.DBPersistence(nil, "GetMoreTDAccountListByQueryParam"))
			if err == nil {
				tdAccountList = append(tdAccountList, moreAccountList...)
			}
		}
	}
	return tdAccountList, nil
}

func (a *accountService) GetAccountById(context context.Context, tdAccountID string) (*mambu.TDAccount, error) {
	tr := tracer.StartTrace(context, "accountService-GetAccountById")
	context = tr.Context()
	defer tr.Finish()
	var tdAccount = new(mambu.TDAccount)
	getUrl := fmt.Sprintf(constant.UrlOf(constant.GetTDAccountUrl), tdAccountID)
	err := a.mambuClient.Get(context, getUrl, tdAccount, a.mambuClient.DBPersistence(context, "GetAccountById"))
	if err != nil {
		a.common.Logger.Error(context, "get account Failed", err)
		return nil, err
	}
	return tdAccount, nil
}

func (a *accountService) UndoMaturityDate(context context.Context, accountID string) bool {
	tr := tracer.StartTrace(context, "accountService-UndoMaturityDate")
	context = tr.Context()
	defer tr.Finish()

	postUrl := fmt.Sprintf(constant.UrlOf(constant.UndoMaturityDateUrl), accountID)
	err := a.mambuClient.Post(context, postUrl, "", nil, nil, a.mambuClient.DBPersistence(context, "UndoMaturityDate"))
	if err != nil {
		a.common.Logger.Error(context, fmt.Sprintf("Undo MaturityDate for td account failed! td acc id: %v", accountID), err)
		return false
	}
	return true
}

func (a *accountService) ChangeMaturityDate(ctx context.Context, accountID, maturityDate, note string) (mambu.TDAccount, error) {
	tr := tracer.StartTrace(ctx, "accountService-ChangeMaturityDate")
	ctx = tr.Context()
	defer tr.Finish()

	postUrl := fmt.Sprintf(constant.UrlOf(constant.StartMaturityDateUrl), accountID)
	a.common.Logger.Info(ctx, fmt.Sprintf("StartMaturityDateUrl: %v", postUrl))

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
	err := a.mambuClient.Post(ctx, postUrl, postJsonStr, &resultTDAccount, nil, a.mambuClient.DBPersistence(ctx, "ChangeMaturityDate"))
	if err != nil {
		a.common.Logger.Error(ctx, "Create MaturityDate for td account failed! ", err, zap.String("accountId", accountID))
		return resultTDAccount, err
	}
	return resultTDAccount, nil
}

func (a *accountService) ApplyProfit(context context.Context, accountID, note string, interestApplicationDate time.Time) bool {
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
	err := a.mambuClient.Post(context, postUrl, postJsonStr, nil, nil, a.mambuClient.DBPersistence(context, "ApplyProfit"))
	if err != nil {
		a.common.Logger.Error(context, fmt.Sprintf("ApplyProfit for td account failed! td acc id: %v", accountID), err)
		return false
	}
	return true
}

func (a *accountService) UpdateMaturifyDateForTDAccount(context context.Context, accountID, newDate string) bool {
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
	err := a.mambuClient.Patch(context, postUrl, postJsonStr, nil, a.mambuClient.DBPersistence(context, "UpdateMaturifyDateForTDAccount"))
	if err != nil {
		a.common.Logger.Error(context, fmt.Sprintf("Undo MaturityDate for td account failed! td acc id: %v", accountID), err)
		return false
	}
	return true
}

func (a *accountService) CloseAccount(context context.Context, accID, notes string) bool {
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
	err := a.mambuClient.Post(context, postUrl, postJsonStr, nil, nil, a.mambuClient.DBPersistence(context, "CloseAccount"))
	if err != nil {
		a.common.Logger.Error(context, fmt.Sprintf("CloseAccount for td account failed! td acc id: %v", accID), err)
		return false
	}
	return true
}
