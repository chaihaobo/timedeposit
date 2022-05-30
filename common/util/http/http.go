// Package http
// @author： Boice
// @createTime：2022/5/30 09:35
package http

import (
	"encoding/json"
	"github.com/guonaihong/gout"
	"github.com/pkg/errors"
	"gitlab.com/bns-engineering/td/common/constant"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
	"strings"
)

type RequestCallbackFun func(url string, code int, requestBody string, responseBody string, err error)

//
//  getMambuHeader mambu base header
//  @Description: get the base header
//  @return map[string][]string
//
func getMambuHeader() map[string][]string {
	return map[string][]string{
		"Accept": {constant.Accept},
		"Apikey": {constant.Apikey},
	}
}

func Post(url, body string, resultBind interface{}, callback RequestCallbackFun) error {
	var code int
	data := strings.NewReader(body)
	req, _ := http.NewRequest("POST", url, data)
	req.Header = getMambuHeader()
	call := gout.POST(url).SetHeader(getMambuHeader()).
		Debug(true).SetJSON(body)
	if resultBind != nil {
		call.BindJSON(resultBind)
	}
	err := call.Code(&code).Do()

	logError(err)
	if callback != nil {
		marshal, _ := json.Marshal(resultBind)
		callback(url, code, body, string(marshal), err)
	}
	return err
}

func Get(urlStr string, resultBind interface{}, callback RequestCallbackFun) error {
	var code int
	err := gout.GET(urlStr).SetHeader(getMambuHeader()).Code(&code).BindJSON(resultBind).
		RequestUse(new(logRequestMiddler)).
		ResponseUse(new(logResponseMiddler)).Do()
	logError(err)
	if callback != nil {
		marshal, _ := json.Marshal(resultBind)
		callback(urlStr, code, "", string(marshal), err)
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
	zap.L().Info("http response==============", zap.String("url", response.Request.URL.String()),
		zap.Int("response code", code),
		zap.String("response body", string(all)),
	)
	if code != http.StatusOK {
		return errors.WithStack(errors.New("http response status not success"))
	}
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
	return nil
}
