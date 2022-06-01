// Package repository
// @author： Boice
// @createTime：2022/5/26 14:20
package repository

import (
	"gitlab.com/bns-engineering/td/common/db"
	db2 "gitlab.com/bns-engineering/td/model/db"
)

var flowNodeLogRepository *FlowNodeLogRepository

func GetFlowNodeLogRepository() IFlowNodeLogRepository {
	return flowNodeLogRepository
}

type IFlowNodeLogRepository interface {
	Save(log *db2.TFlowNodeLog)
}

type FlowNodeLogRepository struct{}

func (flowNodeLogRepository *FlowNodeLogRepository) Save(log *db2.TFlowNodeLog) {
	db.GetDB().Save(log)
}

func init() {
	flowTaskInfoRepository = new(FlowTaskInfoRepository)
}
