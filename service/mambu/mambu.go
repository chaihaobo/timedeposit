// Package mambu
// @author： Boice
// @createTime：2022/7/26 11:44
package mambu

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/guonaihong/gout"
	"github.com/pkg/errors"
	"gitlab.com/bns-engineering/td/common"
	"gitlab.com/bns-engineering/td/common/constant"
	"gitlab.com/bns-engineering/td/model/po"
	"gitlab.com/bns-engineering/td/repository"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
	"time"
)

type RequestCallbackFun func(url string, code int, requestBody string, responseBody string, err error)

type Mambu struct {
	Account     AccountService
	Holiday     HolidayService
	Transaction TransactionService
}

func NewMambuService(common *common.Common, repo *repository.Repository) *Mambu {
	mambuClient := newClient(common, repo)

	return &Mambu{
		Account:     newAccountService(common, repo, mambuClient),
		Holiday:     newHolidayService(common, repo, mambuClient),
		Transaction: newTransactionService(common, repo, mambuClient),
	}
}

type Client interface {
	Patch(ctx context.Context, url, body string, resultBind interface{}, callback RequestCallbackFun) error
	Post(ctx context.Context, url, body string, resultBind, headerBind interface{}, callback RequestCallbackFun) error
	Get(ctx context.Context, url string, resultBind interface{}, callback RequestCallbackFun) error
	DBPersistence(context context.Context, requestType string) RequestCallbackFun
}

func newClient(common *common.Common, repo *repository.Repository) Client {
	return &client{
		common: common,
		repo:   repo,
	}
}

type client struct {
	common *common.Common
	repo   *repository.Repository
}

func (c *client) Patch(ctx context.Context, url, body string, resultBind interface{}, callback RequestCallbackFun) error {
	var code int
	var response string
	call := gout.PATCH(url).SetHeader(c.getMambuHeader(ctx)).
		RequestUse(&logRequestMiddler{ctx}).
		ResponseUse(&logResponseMiddler{ctx}).
		BindBody(&response)
	if body != "" {
		call.SetJSON(body)
	}
	err := call.Code(&code).Do()
	if resultBind != nil && response != "" {
		_ = json.Unmarshal([]byte(response), resultBind)
	}
	c.logError(ctx, err)
	if err != nil {
		return err
	}
	if err == nil && code != http.StatusOK && code != http.StatusNoContent && code != http.StatusCreated {
		c.common.Logger.Error(ctx, "http response status code is not success", errors.New("status code error"), zap.Int("status", code))
		err = errors.WithStack(errors.New("http response status not success"))
	}
	if callback != nil {
		callback(url, code, body, response, err)
	}
	return err
}

func (c *client) DBPersistence(context context.Context, requestType string) RequestCallbackFun {
	return func(url string, code int, requestBody string, responseBody string, err error) {
		var flowId = ""
		var nodeName = ""
		var accountId = ""
		if context != nil {
			contextFlowId := context.Value("flowId")
			contextNodeName := context.Value("nodeName")
			contextAccountId := context.Value("accountId")
			if contextFlowId != nil {
				flowId = contextFlowId.(string)
			}
			if contextNodeName != nil {
				nodeName = contextNodeName.(string)
			}
			if contextAccountId != nil {
				accountId = contextAccountId.(string)
			}

		}
		requestLog := po.NewTMambuRequestLogsBuilder().FlowId(flowId).
			NodeName(nodeName).
			Type(requestType).
			RequestUrl(url).
			AccountId(accountId).
			RequestBody(requestBody).
			ResponseCode(code).
			ResponseBody(responseBody).
			CreateTime(time.Now()).
			UpdateTime(time.Now())
		if err != nil {
			requestLog.Error(err.Error())
		}
		c.common.DB.Save(requestLog.Build())
	}

}

func (c *client) logError(ctx context.Context, err error) {
	if err != nil {
		c.common.Logger.Error(ctx, "call api error", err, zap.Error(errors.WithStack(err)))
	}
}

func (c *client) getMambuHeader(ctx context.Context) map[string][]string {
	value := ctx.Value(constant.ContextIdempotencyKey)
	key := ""
	if value != nil {
		key = value.(string)
	}
	return map[string][]string{
		"Accept":                       {constant.Accept},
		"Apikey":                       {c.common.Config.Mambu.ApiKey},
		constant.ContextIdempotencyKey: {key},
	}

}

func (c *client) Post(ctx context.Context, url, body string, resultBind, headerBind interface{}, callback RequestCallbackFun) error {
	var code int
	var response string
	call := gout.POST(url).SetHeader(c.getMambuHeader(ctx)).
		RequestUse(&logRequestMiddler{ctx}).
		ResponseUse(&logResponseMiddler{ctx}).
		BindBody(&response)
	if body != "" {
		call.SetJSON(body)
	}
	if headerBind != nil {
		call.BindHeader(headerBind)
	}
	err := call.Code(&code).Do()
	if resultBind != nil && response != "" {
		_ = json.Unmarshal([]byte(response), resultBind)
	}
	c.logError(ctx, err)

	if err == nil && code != http.StatusOK && code != http.StatusNoContent && code != http.StatusCreated {
		c.common.Logger.Error(ctx, "http response status code is not success", errors.New("status error"), zap.Int("status", code))
		err = errors.WithStack(errors.New("http response status not success"))
	}
	if callback != nil {
		callback(url, code, body, response, err)
	}
	return err
}

func (c *client) Get(ctx context.Context, url string, resultBind interface{}, callback RequestCallbackFun) error {
	var code int
	var response string
	call := gout.GET(url).SetHeader(c.getMambuHeader(ctx)).
		RequestUse(&logRequestMiddler{ctx}).
		ResponseUse(&logResponseMiddler{ctx}).
		BindBody(&response).Code(&code)
	err := call.Do()
	if resultBind != nil {
		_ = json.Unmarshal([]byte(response), resultBind)
	}
	c.logError(ctx, err)
	if err == nil && code != http.StatusOK && code != http.StatusNoContent && code != http.StatusCreated {
		c.common.Logger.Error(ctx, "http response status code is not success", errors.New("status error"), zap.Int("status", code))
		err = errors.WithStack(errors.New("http response status not success"))
	}
	if callback != nil {
		callback(url, code, "", response, err)
	}
	return err
}

type logResponseMiddler struct {
	context.Context
}

func (d *logResponseMiddler) ModifyResponse(response *http.Response) error {
	all, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}
	code := response.StatusCode
	common.L.Info(d.Context, "http response", zap.String("url", response.Request.URL.String()),
		zap.Int("response code", code),
		zap.String("response body", string(all)),
	)

	response.Body = ioutil.NopCloser(bytes.NewReader(all))
	return nil
}

type logRequestMiddler struct {
	context.Context
}

func (d *logRequestMiddler) ModifyRequest(request *http.Request) error {
	all, err := ioutil.ReadAll(request.Body)
	if err != nil {
		return err
	}
	common.L.Info(d.Context, "http request", zap.String("url", request.URL.String()),
		zap.Any("request headers", request.Header),
		zap.String("request body", string(all)),
	)
	request.Body = ioutil.NopCloser(bytes.NewReader(all))
	return nil
}
