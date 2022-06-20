// Package http
// @author： Boice
// @createTime：2022/5/30 09:35
package mambu_http

import (
	"bytes"
	"encoding/json"
	"github.com/guonaihong/gout"
	"github.com/pkg/errors"
	"gitlab.com/bns-engineering/td/common/config"
	"gitlab.com/bns-engineering/td/common/constant"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
)

type RequestCallbackFun func(url string, code int, requestBody string, responseBody string, err error)

// getMambuHeader get the base header
func getMambuHeader() map[string][]string {
	return map[string][]string{
		"Accept": {constant.Accept},
		"Apikey": {config.TDConf.Mambu.ApiKey},
	}
}

// Patch send http patch request
func Patch(url, body string, resultBind interface{}, callback RequestCallbackFun) error {
	var code int
	var response string
	call := gout.PATCH(url).SetHeader(getMambuHeader()).
		RequestUse(new(logRequestMiddler)).
		ResponseUse(new(logResponseMiddler)).
		BindBody(&response)
	if body != "" {
		call.SetJSON(body)
	}
	err := call.Code(&code).Do()
	if resultBind != nil && response != "" {
		_ = json.Unmarshal([]byte(response), resultBind)
	}
	logError(err)
	if err != nil {
		return err
	}
	if err == nil && code != http.StatusOK && code != http.StatusNoContent && code != http.StatusCreated {
		zap.L().Error("http response status code is not success", zap.Int("status", code))
		err = errors.WithStack(errors.New("http response status not success"))
	}
	if callback != nil {
		callback(url, code, body, response, err)
	}
	return err
}

func Post(url, body string, resultBind, headerBind interface{}, callback RequestCallbackFun) error {
	var code int
	var response string
	call := gout.POST(url).SetHeader(getMambuHeader()).
		RequestUse(new(logRequestMiddler)).
		ResponseUse(new(logResponseMiddler)).
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
	logError(err)

	if err == nil && code != http.StatusOK && code != http.StatusNoContent && code != http.StatusCreated {
		zap.L().Error("http response status code is not success", zap.Int("status", code))
		err = errors.WithStack(errors.New("http response status not success"))
	}
	if callback != nil {
		callback(url, code, body, response, err)
	}
	return err
}

func Get(url string, resultBind interface{}, callback RequestCallbackFun) error {
	var code int
	var response string
	call := gout.GET(url).SetHeader(getMambuHeader()).
		RequestUse(new(logRequestMiddler)).
		ResponseUse(new(logResponseMiddler)).
		BindBody(&response).Code(&code)
	err := call.Do()
	if resultBind != nil {
		_ = json.Unmarshal([]byte(response), resultBind)
	}
	logError(err)
	if err == nil && code != http.StatusOK && code != http.StatusNoContent && code != http.StatusCreated {
		zap.L().Error("http response status code is not success", zap.Int("status", code))
		err = errors.WithStack(errors.New("http response status not success"))
	}
	if callback != nil {
		callback(url, code, "", response, err)
	}
	return err
}

func logError(err error) {
	if err != nil {
		zap.L().Error("call api error", zap.Error(errors.WithStack(err)))
	}
}

type logResponseMiddler struct {
}

func (d *logResponseMiddler) ModifyResponse(response *http.Response) error {
	all, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}
	code := response.StatusCode
	zap.L().Info("http response", zap.String("url", response.Request.URL.String()),
		zap.Int("response code", code),
		zap.String("response body", string(all)),
	)

	response.Body = ioutil.NopCloser(bytes.NewReader(all))
	return nil
}

type logRequestMiddler struct {
}

func (d *logRequestMiddler) ModifyRequest(request *http.Request) error {
	all, err := ioutil.ReadAll(request.Body)
	if err != nil {
		return err
	}
	zap.L().Info("http request", zap.String("url", request.URL.String()),
		zap.Any("request headers", request.Header),
		zap.String("request body", string(all)),
	)
	request.Body = ioutil.NopCloser(bytes.NewReader(all))
	return nil
}
