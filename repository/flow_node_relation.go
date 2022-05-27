// Package repository
// @author： Boice
// @createTime：2022/5/26 13:54
package repository

import (
	"gitlab.com/bns-engineering/td/common/db"
	"gitlab.com/bns-engineering/td/model"
)

var flowNodeRepository *FlowNodeRelationRepository

type IFlowNodeRelationRepository interface {
	GetFlowNodeRelationListByFlowName(flowName string) []*model.TFlowNodeRelation
}

type FlowNodeRelationRepository struct{}

func (FlowNodeRelationRepository *FlowNodeRelationRepository) GetFlowNodeRelationListByFlowName(flowName string) []*model.TFlowNodeRelation {
	nodes := make([]*model.TFlowNodeRelation, 0)
	db.GetDB().Where("flow_name", flowName).Find(&nodes)
	return nodes
}

func init() {
	flowNodeRepository = new(FlowNodeRelationRepository)
}

func GetFlowNodeRelationRepository() IFlowNodeRelationRepository {
	return flowNodeRepository
}
