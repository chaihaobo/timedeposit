// Package repository
// @author： Boice
// @createTime：2022/5/26 13:41
package repository

import (
	"gitlab.com/bns-engineering/td/common/db"
	"gitlab.com/bns-engineering/td/model"
)

var repository *FlowNodeRepository

type IFlowNodeRepository interface {
	GetFlowNodeListByFlowName(flowName string) []*model.TFlowNode
}

type FlowNodeRepository struct{}

func (flowNodeRepository *FlowNodeRepository) GetFlowNodeListByFlowName(flowName string) []*model.TFlowNode {
	var flowNodes = make([]*model.TFlowNode, 0)
	db.GetDB().Where("flow_name", flowName).Find(&flowNodes)
	return flowNodes
}

func init() {
	repository = new(FlowNodeRepository)
}

func GetFlowNodeRepository() IFlowNodeRepository {
	return repository
}
