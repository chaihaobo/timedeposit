// Package http
// @author： Boice
// @createTime：2022/5/30 09:35
package http

import (
	"github.com/guonaihong/gout"
	"github.com/pkg/errors"
	"gitlab.com/bns-engineering/td/common/constant"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
)

//
//  getMambuHeader mambu base header
//  @Description: get the base header
//  @return map[string][]string
//
func getMambuHeader() map[string][]string {
	return map[string][]string{
		"Content-Type": {constant.ContentType},
		"Accept":       {constant.Accept},
		"Apikey":       {constant.Apikey},
	}
}

func Post(postJsonStr, postUrl string, resultBind interface{}) (int, error) {
	var code int
	err := gout.POST(postUrl).
		Debug(true).
		SetHeader(getMambuHeader()).
		SetJSON(postJsonStr).
		Code(&code).BindJSON(resultBind).
		Do()
	logError(err)
	return code, err
}

func Get(urlStr string, resultBind interface{}) (int, error) {
	var code int
	err := gout.GET(urlStr).SetHeader(getMambuHeader()).Code(&code).BindJSON(resultBind).RequestUse().Do()
	logError(err)
	return code, err
}

func logError(err error) {
	zap.L().Error("call api error", zap.Error(errors.WithStack(err)))
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
	return nil
}
