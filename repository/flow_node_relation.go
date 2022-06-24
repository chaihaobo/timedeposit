// Package repository
// @author： Boice
// @createTime：2022/5/26 13:54
package repository

import (
	"context"
	"gitlab.com/bns-engineering/common/tracer"
	"gitlab.com/bns-engineering/td/common/db"
	"gitlab.com/bns-engineering/td/model/po"
)

var flowNodeRepository *FlowNodeRelationRepository

type IFlowNodeRelationRepository interface {
	GetFlowNodeRelationListByFlowName(ctx context.Context, flowName string) []*po.TFlowNodeRelation
}

type FlowNodeRelationRepository struct{}

func (FlowNodeRelationRepository *FlowNodeRelationRepository) GetFlowNodeRelationListByFlowName(ctx context.Context, flowName string) []*po.TFlowNodeRelation {
	tr := tracer.StartTrace(ctx, "flow_node_relation_repository-GetFlowNodeRelationListByFlowName")
	ctx = tr.Context()
	defer tr.Finish()
	nodes := make([]*po.TFlowNodeRelation, 0)
	db.GetDB().Where("flow_name", flowName).Find(&nodes)
	return nodes
}

func init() {
	flowNodeRepository = new(FlowNodeRelationRepository)
}

func GetFlowNodeRelationRepository() IFlowNodeRelationRepository {
	return flowNodeRepository
}
