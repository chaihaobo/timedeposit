// Package repository
// @author： Boice
// @createTime：2022/5/26 14:20
package repository

import (
	"context"
	"gitlab.com/bns-engineering/common/tracer"
	"gitlab.com/bns-engineering/td/common/db"
	"gitlab.com/bns-engineering/td/model/po"
)

var flowNodeLogRepository *FlowNodeLogRepository

func GetFlowNodeLogRepository() IFlowNodeLogRepository {
	return flowNodeLogRepository
}

type IFlowNodeLogRepository interface {
	Save(ctx context.Context, log *po.TFlowNodeLog)
}

type FlowNodeLogRepository struct{}

func (flowNodeLogRepository *FlowNodeLogRepository) Save(ctx context.Context, log *po.TFlowNodeLog) {
	tr := tracer.StartTrace(ctx, "flow_node_log_repository-Save")
	ctx = tr.Context()
	defer tr.Finish()
	db.GetDB().Save(log)
}

func init() {
	flowTaskInfoRepository = new(FlowTaskInfoRepository)
}
