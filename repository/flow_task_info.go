// Package repository
// @author： Boice
// @createTime：2022/5/26 14:20
package repository

import (
	"gitlab.com/bns-engineering/td/common/constant"
	"gitlab.com/bns-engineering/td/common/db"
	"gitlab.com/bns-engineering/td/model/dto"
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
	GetLastByAccountId(accountId string) *po.TFlowTaskInfo
	FailFlowList(pagination *dto.Pagination, accountId string) []*po.TFlowTaskInfo
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

func (flowTaskInfoRepository *FlowTaskInfoRepository) FailFlowList(pagination *dto.Pagination, accountId string) []*po.TFlowTaskInfo {
	failTaskInfoList := make([]*po.TFlowTaskInfo, 0)
	query := db.GetDB().Model(new(po.TFlowTaskInfo)).Where("cur_status = ? and enable = ?", string(constant.FlowNodeFailed), 1).Order("id desc")
	if accountId != "" {
		query = query.Where("account_id", accountId)
	}
	db.FindPage(query, pagination, &failTaskInfoList)
	return failTaskInfoList
}

func (flowTaskInfoRepository *FlowTaskInfoRepository) AllFailFlowIdList() []string {
	failFlowIdList := make([]string, 0)
	db.GetDB().Model(new(po.TFlowTaskInfo)).Where("cur_status = ? and enable = ?", string(constant.FlowNodeFailed), 1).Order("id desc").Pluck("flow_id", &failFlowIdList)
	return failFlowIdList
}

func (flowTaskInfoRepository *FlowTaskInfoRepository) GetLastByAccountId(accountId string) *po.TFlowTaskInfo {
	taskInfo := new(po.TFlowTaskInfo)
	db.GetDB().Where("account_id = ? and enable = ?", accountId, 1).Order("id desc").First(taskInfo)
	if taskInfo.Id > 0 {
		return taskInfo
	}
	return nil
}

func init() {
	flowTaskInfoRepository = new(FlowTaskInfoRepository)
}
