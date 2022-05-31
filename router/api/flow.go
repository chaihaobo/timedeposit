// Package api
// @author： Boice
// @createTime：2022/5/27 10:03
package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/thoas/go-funk"
	"gitlab.com/bns-engineering/td/common/util"
	"gitlab.com/bns-engineering/td/core/engine"
	"gitlab.com/bns-engineering/td/core/engine/mambu/accountservice"
	"gitlab.com/bns-engineering/td/model"
	"gitlab.com/bns-engineering/td/repository"
	"gitlab.com/bns-engineering/td/router/api/dto"
	"gitlab.com/bns-engineering/td/service/mambuEntity"
	"go.uber.org/zap"
	"net/http"
	"time"
)

func StartFlow(c *gin.Context) {
	// Get all td accounts which need to process
	tmpQueryParam := generateSearchTDAccountParam()
	tmpTDAccountList, err := accountservice.GetTDAccountListByQueryParam(tmpQueryParam)
	if err != nil {
		zap.L().Error(fmt.Sprintf("Query mambu service for TD Account List failed! error: %v", err))
		return
	}
	if len(tmpTDAccountList) == 0 {
		zap.L().Info("Query mambu service for TD Account List get empty! No TD Account need to process")
		return
	}

	for _, tmpTDAcc := range tmpTDAccountList {
		_ = engine.Pool.Invoke(tmpTDAcc.ID)
		// go engine.Start(tmpTDAcc.ID)
		zap.L().Info("commit task success!", zap.String("account", tmpTDAcc.ID))

	}
	c.JSON(http.StatusOK, success())

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
		d.FlowName = taskInfo.FlowName
		d.FlowStatus = taskInfo.FlowStatus
		d.FailedOperation = taskInfo.CurStatus
		d.CreateTime = taskInfo.CreateTime
		d.UpdateTime = taskInfo.UpdateTime
		return d
	})
	c.JSON(http.StatusOK, successData(dto.NewPageResult(total, result)))
}

func Retry(c *gin.Context) {
	m := new(dto.RetryFlowModel)
	_ = c.BindJSON(m)
	list := m.FlowIdList
	for _, flowId := range list {
		_ = engine.RetryPool.Invoke(flowId)
	}
	c.JSON(http.StatusOK, success())
}

func RetryAll(c *gin.Context) {
	failFlowIdList := repository.GetFlowTaskInfoRepository().AllFailFlowIdList()
	for _, flowId := range failFlowIdList {
		zap.L().Info("retry flow", zap.String("flowId", flowId))
		_ = engine.RetryPool.Invoke(flowId)
	}
	c.JSON(http.StatusOK, success())
}

func generateSearchTDAccountParam() mambuEntity.SearchParam {
	tmpQueryParam := mambuEntity.SearchParam{
		FilterCriteria: []mambuEntity.FilterCriteria{
			{
				Field:    "accountState",
				Operator: "IN",
				Values:   []string{"ACTIVE", "MATURED"},
			},
			{
				Field:    "accountType",
				Operator: "EQUALS",
				Value:    "FIXED_DEPOSIT",
			},
			{
				Field: "_rekening.rekeningTanggalJatohTempo",
				// todo: Remember to set the value to today!
				Operator:    "BETWEEN",
				Value:       util.GetDate(time.Now().AddDate(0, 0, -20)), // today
				SecondValue: util.GetDate(time.Now().AddDate(0, 0, 1)),   // tomorrow
			},
		},
		SortingCriteria: mambuEntity.SortingCriteria{
			Field: "id",
			Order: "ASC",
		},
	}
	return tmpQueryParam
}
