// Package repository
// @author： Boice
// @createTime：2022/5/26 14:20
package repository

import (
	"gitlab.com/bns-engineering/td/common/db"
	"gitlab.com/bns-engineering/td/model"
	"time"
)

var flowTaskInfoRepository *FlowTaskInfoRepository

func GetFlowTaskInfoRepository() IFlowTaskInfoRepository {
	return flowTaskInfoRepository
}

type IFlowTaskInfoRepository interface {
	Get(flowId string) *model.TFlowTaskInfo
	Update(flowTaskInfo *model.TFlowTaskInfo)
}

type FlowTaskInfoRepository struct{}

func (flowTaskInfoRepository *FlowTaskInfoRepository) Get(flowId string) *model.TFlowTaskInfo {
	flowTaskInfo := new(model.TFlowTaskInfo)
	db.GetDB().Where("flow_id", flowId).First(flowTaskInfo)
	if flowTaskInfo.Id > 0 {
		return flowTaskInfo
	} else {
		return nil
	}
}

func (flowTaskInfoRepository *FlowTaskInfoRepository) Update(flowTaskInfo *model.TFlowTaskInfo) {
	flowTaskInfo.UpdateTime = time.Now()
	db.GetDB().Save(flowTaskInfo)
}

func init() {
	flowTaskInfoRepository = new(FlowTaskInfoRepository)
}
