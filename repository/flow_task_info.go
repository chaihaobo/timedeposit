// Package repository
// @author： Boice
// @createTime：2022/5/26 14:20
package repository

import (
	"gitlab.com/bns-engineering/td/common/constant"
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
	FailFlowList(pageNo int, pageSize int) ([]*model.TFlowTaskInfo, int64)
	AllFailFlowIdList() []string
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

func (flowTaskInfoRepository *FlowTaskInfoRepository) FailFlowList(pageNo int, pageSize int) ([]*model.TFlowTaskInfo, int64) {
	failTaskInfoList := make([]*model.TFlowTaskInfo, 0)
	db.GetDB().Where("cur_status", string(constant.FlowNodeFailed)).Order("id desc").
		Limit(pageSize).Offset((pageNo - 1) * pageSize).
		Find(&failTaskInfoList)
	var total int64
	db.GetDB().Model(new(model.TFlowTaskInfo)).Count(&total)
	return failTaskInfoList, total
}

func (flowTaskInfoRepository *FlowTaskInfoRepository) AllFailFlowIdList() []string {
	failFlowIdList := make([]string, 0)
	db.GetDB().Model(new(model.TFlowTaskInfo)).Where("cur_status", string(constant.FlowNodeFailed)).Order("id desc").Pluck("flow_id", &failFlowIdList)
	return failFlowIdList
}

func init() {
	flowTaskInfoRepository = new(FlowTaskInfoRepository)
}
