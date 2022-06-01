// Package repository
// @author： Boice
// @createTime：2022/5/26 14:20
package repository

import (
	"gitlab.com/bns-engineering/td/common/constant"
	"gitlab.com/bns-engineering/td/common/db"
	db2 "gitlab.com/bns-engineering/td/model/db"
	"time"
)

var flowTaskInfoRepository *FlowTaskInfoRepository

func GetFlowTaskInfoRepository() IFlowTaskInfoRepository {
	return flowTaskInfoRepository
}

type IFlowTaskInfoRepository interface {
	Get(flowId string) *db2.TFlowTaskInfo
	Update(flowTaskInfo *db2.TFlowTaskInfo)
	FailFlowList(pageNo int, pageSize int, accountId string) ([]*db2.TFlowTaskInfo, int64)
	AllFailFlowIdList() []string
}

type FlowTaskInfoRepository struct{}

func (flowTaskInfoRepository *FlowTaskInfoRepository) Get(flowId string) *db2.TFlowTaskInfo {
	flowTaskInfo := new(db2.TFlowTaskInfo)
	db.GetDB().Where("flow_id", flowId).Last(flowTaskInfo)
	if flowTaskInfo.Id > 0 {
		return flowTaskInfo
	} else {
		return nil
	}
}

func (flowTaskInfoRepository *FlowTaskInfoRepository) Update(flowTaskInfo *db2.TFlowTaskInfo) {
	flowTaskInfo.UpdateTime = time.Now()
	db.GetDB().Save(flowTaskInfo)
}

func (flowTaskInfoRepository *FlowTaskInfoRepository) FailFlowList(pageNo int, pageSize int, accountId string) ([]*db2.TFlowTaskInfo, int64) {
	failTaskInfoList := make([]*db2.TFlowTaskInfo, 0)
	query := db.GetDB().Where("cur_status", string(constant.FlowNodeFailed))
	if accountId != "" {
		query = query.Where("account_id", accountId)
	}
	query.Order("id desc").
		Limit(pageSize).Offset((pageNo - 1) * pageSize).Find(&failTaskInfoList)

	var total int64
	db.GetDB().Model(new(db2.TFlowTaskInfo)).Count(&total)
	return failTaskInfoList, total
}

func (flowTaskInfoRepository *FlowTaskInfoRepository) AllFailFlowIdList() []string {
	failFlowIdList := make([]string, 0)
	db.GetDB().Model(new(db2.TFlowTaskInfo)).Where("cur_status", string(constant.FlowNodeFailed)).Order("id desc").Pluck("flow_id", &failFlowIdList)
	return failFlowIdList
}

func init() {
	flowTaskInfoRepository = new(FlowTaskInfoRepository)
}
