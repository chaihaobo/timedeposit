// Package api
// @author： Boice
// @createTime：2022/5/27 10:03
package api

import (
	"github.com/gin-gonic/gin"
	"github.com/thoas/go-funk"
	time2 "gitlab.com/bns-engineering/td/common/util/time"
	"gitlab.com/bns-engineering/td/core/engine"
	"gitlab.com/bns-engineering/td/core/engine/mambu/accountservice"
	dto2 "gitlab.com/bns-engineering/td/model/dto"
	"gitlab.com/bns-engineering/td/model/mambu"
	"gitlab.com/bns-engineering/td/model/po"
	"gitlab.com/bns-engineering/td/repository"
	"go.uber.org/zap"
	"net/http"
	"time"
)

func StartFlow(c *gin.Context) {
	tmpTDAccountList := loadAccountList()
	if tmpTDAccountList == nil || len(tmpTDAccountList) == 0 {
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

func loadAccountList() []mambu.TDAccount {
	// Get all td accounts which need to process
	tmpQueryParam := generateSearchTDAccountParam()
	tmpTDAccountList, err := accountservice.GetTDAccountListByQueryParam(tmpQueryParam)
	if err != nil {
		return nil
	}
	return tmpTDAccountList
}

func FailFlowList(c *gin.Context) {
	retryFlowSearchModel := dto2.DefaultRetryFlowSearchModel()
	_ = c.BindJSON(retryFlowSearchModel)
	list, total := repository.GetFlowTaskInfoRepository().FailFlowList(retryFlowSearchModel.Page.PageNo, retryFlowSearchModel.Page.PageSize, retryFlowSearchModel.Search.AccountId)
	result := funk.Map(list, func(taskInfo *po.TFlowTaskInfo) *dto2.FailFlowModel {
		d := new(dto2.FailFlowModel)
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
	c.JSON(http.StatusOK, successData(dto2.NewPageResult(total, result)))
}

func Retry(c *gin.Context) {
	m := new(dto2.RetryFlowReqModel)
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

func generateSearchTDAccountParam() mambu.SearchParam {
	tmpQueryParam := mambu.SearchParam{
		FilterCriteria: []mambu.FilterCriteria{
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
				Field:       "_rekening.rekeningTanggalJatohTempo",
				Operator:    "BETWEEN",
				Value:       time2.GetDate(time.Now()),                  // today
				SecondValue: time2.GetDate(time.Now().AddDate(0, 0, 1)), // tomorrow
			},
		},
		SortingCriteria: mambu.SortingCriteria{
			Field: "id",
			Order: "ASC",
		},
	}
	return tmpQueryParam
}
