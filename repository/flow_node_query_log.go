// Package repository
// @author： Boice
// @createTime：2022/5/27 10:58
package repository

import (
	"gitlab.com/bns-engineering/td/common/db"
	db2 "gitlab.com/bns-engineering/td/model/db"
	"time"
)

var flowNodeQueryLogRepository *FlowNodeQueryLogRepository

type IFlowNodeQueryLogRepository interface {
	SaveLog(flowId string, nodeName string, queryType string, data string)
	GetLog(flowId string, nodeName string, queryType string) *db2.TFlowNodeQueryLog
	GetNewLog(flowId string, queryType string) *db2.TFlowNodeQueryLog
}

type FlowNodeQueryLogRepository struct {
}

func (f *FlowNodeQueryLogRepository) SaveLog(flowId string, nodeName string, queryType string, data string) {
	log := new(db2.TFlowNodeQueryLog)
	log.FLowId = flowId
	log.NodeName = nodeName
	log.QueryType = queryType
	log.Data = data
	log.CreateTime = time.Now()
	log.UpdateTime = time.Now()
	db.GetDB().Save(log)
}

func (f *FlowNodeQueryLogRepository) GetLog(flowId string, nodeName string, queryType string) *db2.TFlowNodeQueryLog {
	log := new(db2.TFlowNodeQueryLog)
	db.GetDB().Where("flow_id", flowId).Where("node_name", nodeName).Where("query_type", queryType).First(log)
	if log.ID > 0 {
		return log
	}
	return nil
}

func (f *FlowNodeQueryLogRepository) GetNewLog(flowId string, queryType string) *db2.TFlowNodeQueryLog {
	log := new(db2.TFlowNodeQueryLog)
	db.GetDB().Where("flow_id", flowId).Where("query_type", queryType).Last(log)
	if log.ID > 0 {
		return log
	}
	return nil
}

func GetFlowNodeQueryLogRepository() IFlowNodeQueryLogRepository {
	return flowNodeQueryLogRepository
}

func init() {
	flowNodeQueryLogRepository = new(FlowNodeQueryLogRepository)
}
