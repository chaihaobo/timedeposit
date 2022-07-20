// Package api
// @author： Boice
// @createTime：2022/5/27 10:03
package api

import (
	"context"
	"fmt"
	carbonv2 "github.com/golang-module/carbon/v2"
	"github.com/pkg/errors"
	"gitlab.com/bns-engineering/common/tracer"
	"gitlab.com/bns-engineering/td/common/config"
	"gitlab.com/bns-engineering/td/common/constant"
	"gitlab.com/bns-engineering/td/common/log"
	"go.uber.org/zap"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/thoas/go-funk"
	"github.com/uniplaces/carbon"
	timeUtil "gitlab.com/bns-engineering/td/common/util/time"
	"gitlab.com/bns-engineering/td/core/engine"
	"gitlab.com/bns-engineering/td/core/engine/mambu/accountservice"
	"gitlab.com/bns-engineering/td/model/dto"
	"gitlab.com/bns-engineering/td/model/mambu"
	"gitlab.com/bns-engineering/td/model/po"
	"gitlab.com/bns-engineering/td/repository"
)

func StartFlow(c *gin.Context) (interface{}, error) {
	go func() {
		log.Info(c.Request.Context(), "flow start sleep")
		time.Sleep(config.TDConf.Flow.FlowStartSleepTime)
		log.Info(c.Request.Context(), "flow start sleep end")
		tmpTDAccountList, err := loadAccountList(c.Request.Context())
		if err != nil {
			log.Error(c, "[StartFlow] load mambu account list error : ", err)
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
	}()

	return OK, nil
}

func FailFlowList(c *gin.Context) (interface{}, error) {
	tr := tracer.StartTrace(c.Request.Context(), "controller-FailFlowList")
	ctx := tr.Context()
	defer tr.Finish()
	flowSearchModel := dto.DefaultRetryFlowSearchModel()
	err := c.ShouldBindQuery(flowSearchModel)
	if err != nil {
		return nil, err
	}
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
	return dto.NewPageResult(flowSearchModel.Pagination, result), nil
}

func Retry(c *gin.Context) (interface{}, error) {
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
	return &dto.RetryResponseDTO{
		FlowCount:     len(list),
		FlowFailCount: failCount,
		FailInfo:      errInfo,
	}, nil
}

func RetryAll(c *gin.Context) (interface{}, error) {
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
	return &dto.RetryResponseDTO{
		FlowCount:     len(failFlowIdList),
		FlowFailCount: failCount,
		FailInfo:      errInfo,
	}, nil
}

func Remove(c *gin.Context) (interface{}, error) {
	flowId := c.Param("flowId")
	if flowId == "" {
		return nil, errors.New("flow id must not bee null")
	}

	flow := repository.GetFlowTaskInfoRepository().Get(c.Request.Context(), flowId)

	if flow != nil {
		if flow.FlowStatus != constant.FlowFailed {
			return nil, errors.New("Only tasks with failed status can be removed")
		}

		flow.Enable = false
		repository.GetFlowTaskInfoRepository().Update(c.Request.Context(), flow)
	} else {
		return nil, errors.New("flow id not exist!")
	}
	return OK, nil
}

func Metric(c *gin.Context) (interface{}, error) {
	metricRequest := new(dto.FlowMetricQueryModel)
	err := c.ShouldBindJSON(metricRequest)
	if err != nil {
		return nil, err
	}
	startDate := metricRequest.StartDate
	endDate := metricRequest.EndDate

	if startDate == "" || endDate == "" {
		return nil, errors.New("These parameters are empty(start_date or end_date)")
	}

	carbonStart := carbonv2.ParseByFormat(startDate, carbonv2.DateLayout)
	if carbonStart.Error != nil {
		return nil, errors.New("start date format error")
	}

	carbonEnd := carbonv2.ParseByFormat(endDate, carbonv2.DateLayout)
	if carbonEnd.Error != nil {
		return nil, errors.New("end date format error")
	}
	if carbonStart.Gt(carbonEnd) {
		return nil, errors.New("start date can not after end date")
	}

	var metricDays = make([]string, 0)
	for carbonStart.Lte(carbonEnd) {
		metricDays = append(metricDays, carbonStart.ToDateString())
		carbonStart = carbonStart.AddDay()
	}

	return repository.GetFlowTaskInfoRepository().MetricByDay(c.Request.Context(), metricDays), nil
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
