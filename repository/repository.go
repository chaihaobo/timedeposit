// Package repository
// @author： Boice
// @createTime：2022/7/26 10:08
package repository

import "gitlab.com/bns-engineering/td/common"

type Repository struct {
	FlowNode         IFlowNodeRepository
	FlowNodeLog      IFlowNodeLogRepository
	FlowNodeQueryLog IFlowNodeQueryLogRepository
	FlowNodeRelation IFlowNodeRelationRepository
	FlowTaskInfo     IFlowTaskInfoRepository
	FlowTransaction  IFlowTransactionRepository
	Redis            IRedisRepository
}

func NewRepository(common *common.Common) *Repository {
	return &Repository{
		FlowNode:         newFlowNodeRepository(common),
		FlowNodeLog:      newFlowNodeLogRepository(common),
		FlowNodeQueryLog: newFlowNodeQueryLogRepository(common),
		FlowNodeRelation: newFlowNodeRelationRepository(common),
		FlowTaskInfo:     newFlowTaskInfoRepository(common),
		FlowTransaction:  newFlowTransactionRepository(common),
		Redis:            newRedisRepository(common),
	}
}
