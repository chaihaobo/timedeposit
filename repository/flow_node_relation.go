// Package repository
// @author： Boice
// @createTime：2022/5/26 13:54
package repository

import (
	"gitlab.com/bns-engineering/td/common/db"
	db2 "gitlab.com/bns-engineering/td/model/db"
)

var flowNodeRepository *FlowNodeRelationRepository

type IFlowNodeRelationRepository interface {
	GetFlowNodeRelationListByFlowName(flowName string) []*db2.TFlowNodeRelation
}

type FlowNodeRelationRepository struct{}

func (FlowNodeRelationRepository *FlowNodeRelationRepository) GetFlowNodeRelationListByFlowName(flowName string) []*db2.TFlowNodeRelation {
	nodes := make([]*db2.TFlowNodeRelation, 0)
	db.GetDB().Where("flow_name", flowName).Find(&nodes)
	return nodes
}

func init() {
	flowNodeRepository = new(FlowNodeRelationRepository)
}

func GetFlowNodeRelationRepository() IFlowNodeRelationRepository {
	return flowNodeRepository
}
