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
	"sync"
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
	err := c.ShouldBindQuery(flowSearchModel)
	if err != nil {
		c.JSON(http.StatusOK, Error(err.Error()))
		return
	}
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
	var lock sync.Mutex
	var wait sync.WaitGroup
	wait.Add(len(list))
	errInfo := make([]dto.RetryFailInfoResponse, 0)
	failCount := 0
	for _, flowId := range list {
		id := flowId
		go func() {
			defer wait.Done()
			if err := engine.Run(c.Request.Context(), id); err != nil {
				lock.Lock()
				defer lock.Unlock()
				failCount++
				errInfo = append(errInfo, dto.RetryFailInfoResponse{
					FlowId:  id,
					Message: err.Error(),
				})
			}
		}()
	}
	wait.Wait()
	c.JSON(http.StatusOK, SuccessData(&dto.RetryResponseDTO{
		FlowCount:     len(list),
		FlowFailCount: failCount,
		FailInfo:      errInfo,
	}))
}

func RetryAll(c *gin.Context) {
	failFlowIdList := repository.GetFlowTaskInfoRepository().AllFailFlowIdList(c.Request.Context())
	var lock sync.Mutex
	var wait sync.WaitGroup
	wait.Add(len(failFlowIdList))
	errInfo := make([]dto.RetryFailInfoResponse, 0)
	failCount := 0
	for i := range failFlowIdList {
		id := failFlowIdList[i]
		go func() {
			lock.Lock()
			if err := engine.Run(c.Request.Context(), id); err != nil {
				defer lock.Unlock()
				defer wait.Done()
				failCount++
				errInfo = append(errInfo, dto.RetryFailInfoResponse{
					FlowId:  id,
					Message: err.Error(),
				})
			}
		}()
	}
	wait.Wait()
	c.JSON(http.StatusOK, SuccessData(&dto.RetryResponseDTO{
		FlowCount:     len(failFlowIdList),
		FlowFailCount: failCount,
		FailInfo:      errInfo,
	}))
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
	} else {
		c.JSON(http.StatusOK, Error("flow id not exist!"))
		return
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
