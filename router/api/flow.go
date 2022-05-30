// Package api
// @author： Boice
// @createTime：2022/5/27 10:03
package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/thoas/go-funk"
	"gitlab.com/bns-engineering/td/core/engine"
	"gitlab.com/bns-engineering/td/model"
	"gitlab.com/bns-engineering/td/repository"
	"gitlab.com/bns-engineering/td/router/api/dto"
	mambuservices "gitlab.com/bns-engineering/td/service/mambuServices"
	"go.uber.org/zap"
	"net/http"
)

func StartFlow(c *gin.Context) {
	//Get all td accounts which need to process
	tmpQueryParam := generateSearchTDAccountParam()
	tmpTDAccountList, err := mambuservices.GetTDAccountListByQueryParam(tmpQueryParam)
	if err != nil {
		zap.L().Error(fmt.Sprintf("Query mambu service for TD Account List failed! error: %v", err))
		return
	}
	if len(tmpTDAccountList) == 0 {
		zap.L().Info("Query mambu service for TD Account List get empty! No TD Account need to process")
		return
	}

	for _, tmpTDAcc := range tmpTDAccountList {
		zap.L().Info(fmt.Sprintf("Before Run Flow for Account: %v", tmpTDAcc.ID))
	}
	zap.L().Info("=======================================")

	for _, tmpTDAcc := range tmpTDAccountList {
		err := engine.Pool.Submit(func() {
			engine.Start(tmpTDAcc.ID)
		})
		if err != nil {
			c.String(http.StatusInternalServerError, "%s", err.Error())
		}
	}

}

func FailFlowList(c *gin.Context) {
	page := dto.DefaultPage()
	_ = c.BindJSON(page)
	list, total := repository.GetFlowTaskInfoRepository().FailFlowList(page.PageNo, page.PageSize)
	result := funk.Map(list, func(taskInfo *model.TFlowTaskInfo) *dto.FailFlowModel {
		d := new(dto.FailFlowModel)
		d.Id = taskInfo.Id
		d.FlowId = taskInfo.FlowId
		d.AccountId = taskInfo.AccountId
		d.CurStatus = taskInfo.CurStatus
		d.CurNodeName = taskInfo.CurNodeName
		d.CreateTime = taskInfo.CreateTime
		return d
	})
	c.JSON(http.StatusOK, successData(dto.NewPageResult(total, result)))
}

func Retry(c *gin.Context) {
	m := new(dto.RetryFlowModel)
	_ = c.BindJSON(m)
	list := m.FlowIdList
	for _, flowId := range list {
		_ = engine.Pool.Submit(func() {
			engine.Run(flowId)
		})
	}
	c.JSON(http.StatusOK, success())
}

func RetryAll(c *gin.Context) {
	failFlowIdList := repository.GetFlowTaskInfoRepository().AllFailFlowIdList()

	for _, flowId := range failFlowIdList {
		zap.L().Info("retry flow", zap.String("flowId", flowId))
		_ = engine.Pool.Submit(func() {
			engine.Run(flowId)
		})
	}
	c.JSON(http.StatusOK, success())
}