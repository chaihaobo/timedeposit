// Package repository
// @author： Boice
// @createTime：2022/5/26 13:41
package repository

import (
	"gitlab.com/bns-engineering/td/common/db"
	db2 "gitlab.com/bns-engineering/td/model/db"
)

var repository *FlowNodeRepository

type IFlowNodeRepository interface {
	GetFlowNodeListByFlowName(flowName string) []*db2.TFlowNode
}

type FlowNodeRepository struct{}

func (flowNodeRepository *FlowNodeRepository) GetFlowNodeListByFlowName(flowName string) []*db2.TFlowNode {
	var flowNodes = make([]*db2.TFlowNode, 0)
	db.GetDB().Where("flow_name", flowName).Find(&flowNodes)
	return flowNodes
}

func init() {
	repository = new(FlowNodeRepository)
}

func GetFlowNodeRepository() IFlowNodeRepository {
	return repository
}
