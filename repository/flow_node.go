// Package repository
// @author： Boice
// @createTime：2022/5/26 13:41
package repository

import (
	"context"
	"gitlab.com/bns-engineering/common/tracer"
	"gitlab.com/bns-engineering/td/common"
	"gitlab.com/bns-engineering/td/model/po"
)

type IFlowNodeRepository interface {
	GetFlowNodeListByFlowName(ctx context.Context, flowName string) []*po.TFlowNode
}

type flowNodeRepository struct {
	common *common.Common
}

func (flowNodeRepository *flowNodeRepository) GetFlowNodeListByFlowName(ctx context.Context, flowName string) []*po.TFlowNode {
	tr := tracer.StartTrace(ctx, "flow_node_repository-GetFlowNodeListByFlowName")
	ctx = tr.Context()
	defer tr.Finish()
	var flowNodes = make([]*po.TFlowNode, 0)
	flowNodeRepository.common.DB.Where("flow_name", flowName).Find(&flowNodes)
	return flowNodes
}

func newFlowNodeRepository(common *common.Common) IFlowNodeRepository {
	return &flowNodeRepository{
		common: common,
	}
}
