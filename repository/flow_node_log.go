// Package repository
// @author： Boice
// @createTime：2022/5/26 14:20
package repository

import (
	"context"
	"gitlab.com/bns-engineering/common/tracer"
	"gitlab.com/bns-engineering/td/common"
	"gitlab.com/bns-engineering/td/model/po"
)

type IFlowNodeLogRepository interface {
	Save(ctx context.Context, log *po.TFlowNodeLog)
}

type flowNodeLogRepository struct {
	common *common.Common
}

func (flowNodeLogRepository *flowNodeLogRepository) Save(ctx context.Context, log *po.TFlowNodeLog) {
	tr := tracer.StartTrace(ctx, "flow_node_log_repository-Save")
	ctx = tr.Context()
	defer tr.Finish()
	flowNodeLogRepository.common.DB.Save(log)
}

func newFlowNodeLogRepository(common *common.Common) IFlowNodeLogRepository {
	return &flowNodeLogRepository{}
}
