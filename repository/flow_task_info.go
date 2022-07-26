// Package repository
// @author： Boice
// @createTime：2022/5/26 14:20
package repository

import (
	"context"
	"fmt"
	"github.com/thoas/go-funk"
	"gitlab.com/bns-engineering/common/tracer"
	"gitlab.com/bns-engineering/td/common"
	"gitlab.com/bns-engineering/td/common/constant"
	"gitlab.com/bns-engineering/td/model/dto"
	"gitlab.com/bns-engineering/td/model/po"
	"strings"
	"time"
)

type IFlowTaskInfoRepository interface {
	Get(ctx context.Context, flowId string) *po.TFlowTaskInfo
	Update(ctx context.Context, flowTaskInfo *po.TFlowTaskInfo)
	GetLastByAccountId(ctx context.Context, accountId string) *po.TFlowTaskInfo
	FailFlowList(ctx context.Context, pagination *dto.Pagination, accountId string) []*po.TFlowTaskInfo
	AllFailFlowIdList(ctx context.Context) []string
	MetricByDay(ctx context.Context, dayList []string) []dto.FlowMetricResultModel
}

type flowTaskInfoRepository struct {
	common *common.Common
}

func (f *flowTaskInfoRepository) Get(ctx context.Context, flowId string) *po.TFlowTaskInfo {
	tr := tracer.StartTrace(ctx, "flow_task_info_repository-GetFlowNodeRelatiGetonListByFlowName")
	ctx = tr.Context()
	defer tr.Finish()
	flowTaskInfo := new(po.TFlowTaskInfo)
	f.common.DB.Where("flow_id =? and enable = ?", flowId, 1).Last(flowTaskInfo)
	if flowTaskInfo.Id > 0 {
		return flowTaskInfo
	}
	return nil
}

func (f *flowTaskInfoRepository) Update(ctx context.Context, flowTaskInfo *po.TFlowTaskInfo) {
	tr := tracer.StartTrace(ctx, "flow_task_info_repository-Update")
	ctx = tr.Context()
	defer tr.Finish()
	flowTaskInfo.UpdateTime = time.Now()
	f.common.DB.Save(flowTaskInfo)
}

func (f *flowTaskInfoRepository) FailFlowList(ctx context.Context, pagination *dto.Pagination, accountId string) []*po.TFlowTaskInfo {
	tr := tracer.StartTrace(ctx, "flow_task_info_repository-FailFlowList")
	ctx = tr.Context()
	defer tr.Finish()
	failTaskInfoList := make([]*po.TFlowTaskInfo, 0)
	query := f.common.DB.Model(new(po.TFlowTaskInfo)).Where("cur_status = ? and enable = ?", string(constant.FlowNodeFailed), 1).Order("create_time")
	if accountId != "" {
		query = query.Where("account_id", accountId)
	}
	common.FindPage(ctx, query, pagination, &failTaskInfoList)
	return failTaskInfoList
}

func (f *flowTaskInfoRepository) AllFailFlowIdList(ctx context.Context) []string {
	tr := tracer.StartTrace(ctx, "flow_task_info_repository-AllFailFlowIdList")
	ctx = tr.Context()
	defer tr.Finish()
	failFlowIdList := make([]string, 0)
	f.common.DB.Model(new(po.TFlowTaskInfo)).Where("cur_status = ? and enable = ?", string(constant.FlowNodeFailed), 1).Order("id desc").Pluck("flow_id", &failFlowIdList)
	return failFlowIdList
}

func (f *flowTaskInfoRepository) GetLastByAccountId(ctx context.Context, accountId string) *po.TFlowTaskInfo {
	tr := tracer.StartTrace(ctx, "flow_task_info_repository-GetLastByAccountId")
	ctx = tr.Context()
	defer tr.Finish()
	taskInfo := new(po.TFlowTaskInfo)
	f.common.DB.Where("account_id = ? and enable = ?", accountId, 1).Order("id desc").First(taskInfo)
	if taskInfo.Id > 0 {
		return taskInfo
	}
	return nil
}

func (f *flowTaskInfoRepository) MetricByDay(ctx context.Context, dayList []string) []dto.FlowMetricResultModel {
	tr := tracer.StartTrace(ctx, "flow_task_info_repository-MetricByDay")
	ctx = tr.Context()
	defer tr.Finish()
	dateSql := strings.Join(funk.Map(dayList, func(day string) string {
		return fmt.Sprintf("select '%s' create_time", day)
	}).([]string), " union ")

	sql := fmt.Sprintf(`
	select t1.create_time create_date, count(t2.id) task_cnt,sum(if(t2.flow_status='flow_finish',1,0)) success_cnt,
		sum(if(t2.flow_status='flow_running',1,0)) running_cnt,
		sum(if(t2.flow_status='flow_failed',1,0)) fail_cnt,
		(exists(select id from t_mambu_request_logs t3 where date(t3.create_time)=t1.create_time)) is_call
		from (%s) t1
			left join t_flow_task_infos t2 on t1.create_time = date(t2.create_time)
			group by 1
`, dateSql)

	result := make([]dto.FlowMetricResultModel, 0)
	f.common.DB.Raw(sql).Scan(&result)
	return result
}

func newFlowTaskInfoRepository(common *common.Common) IFlowTaskInfoRepository {
	return &flowTaskInfoRepository{
		common: common,
	}
}
