// Package mambu
// @author： Boice
// @createTime：2022/5/30 14:10
package persistence

import (
	"context"
	"gitlab.com/bns-engineering/td/common/db"
	"gitlab.com/bns-engineering/td/common/util/mambu_http"
	"gitlab.com/bns-engineering/td/model"
	"time"
)

func DBPersistence(context context.Context, requestType string) mambu_http.RequestCallbackFun {
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
		requestLog := model.NewTMambuRequestLogsBuilder().FlowId(flowId).
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
		db.GetDB().Save(requestLog.Build())

	}

}
