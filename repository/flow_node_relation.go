// Package repository
// @author： Boice
// @createTime：2022/5/26 13:54
package repository

import (
	"context"
	"gitlab.com/bns-engineering/common/tracer"
	"gitlab.com/bns-engineering/td/common"
	"gitlab.com/bns-engineering/td/model/po"
)

type IFlowNodeRelationRepository interface {
	GetFlowNodeRelationListByFlowName(ctx context.Context, flowName string) []*po.TFlowNodeRelation
}

type flowNodeRelationRepository struct {
	common *common.Common
}

func (f *flowNodeRelationRepository) GetFlowNodeRelationListByFlowName(ctx context.Context, flowName string) []*po.TFlowNodeRelation {
	tr := tracer.StartTrace(ctx, "flow_node_relation_repository-GetFlowNodeRelationListByFlowName")
	ctx = tr.Context()
	defer tr.Finish()
	nodes := make([]*po.TFlowNodeRelation, 0)
	f.common.DB.Where("flow_name", flowName).Find(&nodes)
	return nodes
}

func newFlowNodeRelationRepository(common *common.Common) IFlowNodeRelationRepository {
	return &flowNodeRelationRepository{
		common: common,
	}
}
