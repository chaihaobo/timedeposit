// Package api
// @author： Boice
// @createTime：2022/5/27 10:03
package api

import (
	"github.com/gin-gonic/gin"
	"github.com/thoas/go-funk"
	timeUtil "gitlab.com/bns-engineering/td/common/util/time"
	"gitlab.com/bns-engineering/td/core/engine"
	"gitlab.com/bns-engineering/td/core/engine/mambu/accountservice"
	"gitlab.com/bns-engineering/td/model/dto"
	"gitlab.com/bns-engineering/td/model/mambu"
	"gitlab.com/bns-engineering/td/model/po"
	"gitlab.com/bns-engineering/td/repository"
	"go.uber.org/zap"
	"net/http"
	"time"
)

func StartFlow(c *gin.Context) {
	tmpTDAccountList, err := loadAccountList()
	if err != nil {
		zap.L().Error("load mambu account list error")
		c.JSON(http.StatusOK, Error("load mambu account list error"))
	}

	for _, tmpTDAcc := range tmpTDAccountList {
		_ = engine.Pool.Invoke(tmpTDAcc.ID)
		// go engine.Start(tmpTDAcc.ID)
		zap.L().Info("commit task success!", zap.String("account", tmpTDAcc.ID))

	}
	c.JSON(http.StatusOK, Success())
}

func FailFlowList(c *gin.Context) {
	retryFlowSearchModel := dto.DefaultRetryFlowSearchModel()
	_ = c.BindJSON(retryFlowSearchModel)
	list, total := repository.GetFlowTaskInfoRepository().FailFlowList(retryFlowSearchModel.Page.PageNo, retryFlowSearchModel.Page.PageSize, retryFlowSearchModel.Search.AccountId)
	result := funk.Map(list, func(taskInfo *po.TFlowTaskInfo) *dto.FailFlowModel {
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
	c.JSON(http.StatusOK, SuccessData(dto.NewPageResult(total, result)))
}

func Retry(c *gin.Context) {
	m := new(dto.RetryFlowReqModel)
	_ = c.BindJSON(m)
	list := m.FlowIdList
	for _, flowId := range list {
		_ = engine.RetryPool.Invoke(flowId)
	}
	c.JSON(http.StatusOK, Success())
}

func RetryAll(c *gin.Context) {
	failFlowIdList := repository.GetFlowTaskInfoRepository().AllFailFlowIdList()
	for _, flowId := range failFlowIdList {
		zap.L().Info("retry flow", zap.String("flowId", flowId))
		_ = engine.RetryPool.Invoke(flowId)
	}
	c.JSON(http.StatusOK, Success())
}

func Remove(c *gin.Context) {
	flowId := c.Param("flowId")
	if flowId == "" {
		c.JSON(http.StatusOK, Error("flow id must not bee null"))
		return
	}

	flow := repository.GetFlowTaskInfoRepository().Get(flowId)
	if flow != nil {
		flow.Enable = false
		repository.GetFlowTaskInfoRepository().Update(flow)
	}
	c.JSON(http.StatusOK, Success())
}

func loadAccountList() ([]mambu.TDAccount, error) {
	// Get all td accounts which need to process
	tmpQueryParam := generateSearchTDAccountParam()
	tmpTDAccountList, err := accountservice.GetTDAccountListByQueryParam(tmpQueryParam)
	return tmpTDAccountList, err
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
				Value:       timeUtil.GetDate(time.Now()),                  // today
				SecondValue: timeUtil.GetDate(time.Now().AddDate(0, 0, 1)), // tomorrow
			},
		},
		SortingCriteria: mambu.SortingCriteria{
			Field: "id",
			Order: "ASC",
		},
	}
	return tmpQueryParam
}
