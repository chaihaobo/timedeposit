// Package repository
// @author： Boice
// @createTime：2022/5/27 10:58
package repository

import (
	"context"
	"gitlab.com/bns-engineering/common/tracer"
	"gitlab.com/bns-engineering/td/common/db"
	"gitlab.com/bns-engineering/td/model/po"
	"time"
)

var flowNodeQueryLogRepository *FlowNodeQueryLogRepository

type IFlowNodeQueryLogRepository interface {
	SaveLog(ctx context.Context, flowId string, nodeName string, queryType string, data string)
	GetLog(ctx context.Context, flowId string, nodeName string, queryType string) *po.TFlowNodeQueryLog
	GetNewLog(ctx context.Context, flowId string, queryType string) *po.TFlowNodeQueryLog
}

type FlowNodeQueryLogRepository struct {
}

func (f *FlowNodeQueryLogRepository) SaveLog(ctx context.Context, flowId string, nodeName string, queryType string, data string) {
	log := new(po.TFlowNodeQueryLog)
	log.FLowId = flowId
	log.NodeName = nodeName
	log.QueryType = queryType
	log.Data = data
	log.CreateTime = time.Now()
	log.UpdateTime = time.Now()
	db.GetDB().Save(log)
}

func (f *FlowNodeQueryLogRepository) GetLog(ctx context.Context, flowId string, nodeName string, queryType string) *po.TFlowNodeQueryLog {
	tr := tracer.StartTrace(ctx, "flow_node_query_log_repository-GetLog")
	ctx = tr.Context()
	defer tr.Finish()
	log := new(po.TFlowNodeQueryLog)
	db.GetDB().Where("flow_id", flowId).Where("node_name", nodeName).Where("query_type", queryType).First(log)
	if log.ID > 0 {
		return log
	}
	return nil
}

func (f *FlowNodeQueryLogRepository) GetNewLog(ctx context.Context, flowId string, queryType string) *po.TFlowNodeQueryLog {
	tr := tracer.StartTrace(ctx, "flow_node_query_log_repository-GetNewLog")
	ctx = tr.Context()
	defer tr.Finish()
	log := new(po.TFlowNodeQueryLog)
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
