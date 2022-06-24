// Package api
// @author： Boice
// @createTime：2022/5/27 10:03
package api

import (
	"context"
	"fmt"
	"gitlab.com/bns-engineering/common/tracer"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/thoas/go-funk"
	"github.com/uniplaces/carbon"
	"gitlab.com/bns-engineering/td/common/log"
	timeUtil "gitlab.com/bns-engineering/td/common/util/time"
	"gitlab.com/bns-engineering/td/core/engine"
	"gitlab.com/bns-engineering/td/core/engine/mambu/accountservice"
	"gitlab.com/bns-engineering/td/model/dto"
	"gitlab.com/bns-engineering/td/model/mambu"
	"gitlab.com/bns-engineering/td/model/po"
	"gitlab.com/bns-engineering/td/repository"
	"go.uber.org/zap"
)

func StartFlow(c *gin.Context) {
	tmpTDAccountList, err := loadAccountList(c.Request.Context())
	if err != nil {
		log.Error(c, "[StartFlow] load mambu account list error : ", err)
		c.JSON(http.StatusOK, Error("load mambu account list error"))
		return
	}

	for _, tmpTDAcc := range tmpTDAccountList {
		accountLastTask := repository.GetFlowTaskInfoRepository().GetLastByAccountId(c.Request.Context(), tmpTDAcc.ID)
		if accountLastTask != nil && carbon.NewCarbon(accountLastTask.CreateTime).IsSameDay(carbon.Now()) {
			log.Info(c, "account today is already has task,skip it!")
			continue
		}

		go engine.Start(c.Request.Context(), tmpTDAcc.ID)
		log.Info(c, "commit task success!", zap.String("account", tmpTDAcc.ID))

	}
	c.JSON(http.StatusOK, Success())
}

func FailFlowList(c *gin.Context) {
	tr := tracer.StartTrace(c.Request.Context(), "controller-FailFlowList")
	ctx := tr.Context()
	defer tr.Finish()
	flowSearchModel := dto.DefaultRetryFlowSearchModel()
	_ = c.BindQuery(flowSearchModel)
	// retryFlowSearchModel := dto.DefaultRetryFlowSearchModel()
	// _ = c.BindJSON(retryFlowSearchModel)
	list := repository.GetFlowTaskInfoRepository().FailFlowList(ctx, flowSearchModel.Pagination, flowSearchModel.AccountId)
	result := funk.Map(list, func(taskInfo *po.TFlowTaskInfo) *dto.FailFlowModel {
		failedTransactions := repository.GetFlowTransactionRepository().ListErrorTransactionByFlowId(c, taskInfo.FlowId)
		currentFaildTransactions := funk.Filter(failedTransactions, func(failTransaction po.TFlowTransactions) bool {
			if strings.Contains(failTransaction.TransId, fmt.Sprintf("%s-%s", taskInfo.FlowId, taskInfo.CurNodeName)) {
				return true
			}
			return false
		}).([]po.TFlowTransactions)
		d := new(dto.FailFlowModel)
		d.Id = taskInfo.Id
		d.FlowId = taskInfo.FlowId
		d.AccountId = taskInfo.AccountId
		d.FlowName = taskInfo.FlowName
		d.FlowStatus = taskInfo.FlowStatus
		d.FailedOperation = taskInfo.CurNodeName
		d.CreateTime = taskInfo.CreateTime
		d.UpdateTime = taskInfo.UpdateTime
		d.AmountToMove = "0"
		if len(currentFaildTransactions) > 0 {
			d.AmountToMove = fmt.Sprintf("%f", currentFaildTransactions[0].Amount)
		}
		return d
	})
	c.JSON(http.StatusOK, SuccessData(dto.NewPageResult(flowSearchModel.Pagination, result)))
}

func Retry(c *gin.Context) {
	m := new(dto.RetryFlowReqModel)
	_ = c.BindJSON(m)
	list := m.FlowIdList
	for _, flowId := range list {
		go engine.Run(c.Request.Context(), flowId)
	}
	c.JSON(http.StatusOK, Success())
}

func RetryAll(c *gin.Context) {
	failFlowIdList := repository.GetFlowTaskInfoRepository().AllFailFlowIdList(c.Request.Context())
	for _, flowId := range failFlowIdList {
		log.Info(c, "retry flow", zap.String("flowId", flowId))
		go engine.Run(c.Request.Context(), flowId)
	}
	c.JSON(http.StatusOK, Success())
}

func Remove(c *gin.Context) {
	flowId := c.Param("flowId")
	if flowId == "" {
		c.JSON(http.StatusOK, Error("flow id must not bee null"))
		return
	}

	flow := repository.GetFlowTaskInfoRepository().Get(c.Request.Context(), flowId)
	if flow != nil {
		flow.Enable = false
		repository.GetFlowTaskInfoRepository().Update(c.Request.Context(), flow)
	}
	c.JSON(http.StatusOK, Success())
}

func loadAccountList(ctx context.Context) ([]mambu.TDAccount, error) {
	// Get all td accounts which need to process
	tmpQueryParam := generateSearchTDAccountParam()
	tmpTDAccountList, err := accountservice.GetTDAccountListByQueryParam(ctx, tmpQueryParam)
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
				Value:       carbon.Now().DateString(),                     // today
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
