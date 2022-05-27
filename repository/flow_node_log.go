// Package repository
// @author： Boice
// @createTime：2022/5/26 14:20
package repository

import (
	"gitlab.com/bns-engineering/td/common/db"
	"gitlab.com/bns-engineering/td/model"
)

var flowNodeLogRepository *FlowNodeLogRepository

func GetFlowNodeLogRepository() IFlowNodeLogRepository {
	return flowNodeLogRepository
}

type IFlowNodeLogRepository interface {
	Save(log *model.TFlowNodeLog)
}

type FlowNodeLogRepository struct{}

func (flowNodeLogRepository *FlowNodeLogRepository) Save(log *model.TFlowNodeLog) {
	db.GetDB().Save(log)
}

func init() {
	flowTaskInfoRepository = new(FlowTaskInfoRepository)
}
