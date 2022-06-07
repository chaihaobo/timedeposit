// Package repository
// @author： Boice
// @createTime：2022/5/26 14:20
package repository

import (
	"gitlab.com/bns-engineering/td/common/constant"
	"gitlab.com/bns-engineering/td/common/db"
	"gitlab.com/bns-engineering/td/model/po"
	"time"
)

var flowTaskInfoRepository *FlowTaskInfoRepository

func GetFlowTaskInfoRepository() IFlowTaskInfoRepository {
	return flowTaskInfoRepository
}

type IFlowTaskInfoRepository interface {
	Get(flowId string) *po.TFlowTaskInfo
	Update(flowTaskInfo *po.TFlowTaskInfo)
	FailFlowList(pageNo int, pageSize int, accountId string) ([]*po.TFlowTaskInfo, int64)
	AllFailFlowIdList() []string
}

type FlowTaskInfoRepository struct{}

func (flowTaskInfoRepository *FlowTaskInfoRepository) Get(flowId string) *po.TFlowTaskInfo {
	flowTaskInfo := new(po.TFlowTaskInfo)
	db.GetDB().Where("flow_id =? and enable = ?", flowId, 1).Last(flowTaskInfo)
	if flowTaskInfo.Id > 0 {
		return flowTaskInfo
	} else {
		return nil
	}
}

func (flowTaskInfoRepository *FlowTaskInfoRepository) Update(flowTaskInfo *po.TFlowTaskInfo) {
	flowTaskInfo.UpdateTime = time.Now()
	db.GetDB().Save(flowTaskInfo)
}

func (flowTaskInfoRepository *FlowTaskInfoRepository) FailFlowList(pageNo int, pageSize int, accountId string) ([]*po.TFlowTaskInfo, int64) {
	failTaskInfoList := make([]*po.TFlowTaskInfo, 0)
	var total int64
	query := db.GetDB().Model(new(po.TFlowTaskInfo)).Where("cur_status = ? and enable = ?", string(constant.FlowNodeFailed), 1).Order("id desc")
	if accountId != "" {
		query = query.Where("account_id", accountId)
	}
	db.FindPage(query, pageNo, pageSize, &failTaskInfoList, &total)
	return failTaskInfoList, total
}

func (flowTaskInfoRepository *FlowTaskInfoRepository) AllFailFlowIdList() []string {
	failFlowIdList := make([]string, 0)
	db.GetDB().Model(new(po.TFlowTaskInfo)).Where("cur_status = ? and enable = ?", string(constant.FlowNodeFailed), 1).Order("id desc").Pluck("flow_id", &failFlowIdList)
	return failFlowIdList
}

func init() {
	flowTaskInfoRepository = new(FlowTaskInfoRepository)
}
