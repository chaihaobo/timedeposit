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
	GetLogValueOr(ctx context.Context, flowId string, nodeName string, queryType string, valueGenerator func() string) string
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

func (f *FlowNodeQueryLogRepository) GetLogValueOr(ctx context.Context, flowId string, nodeName string, queryType string, valueGenerator func() string) string {
	tr := tracer.StartTrace(ctx, "flow_node_query_log_repository-GetLogValueOr")
	ctx = tr.Context()
	defer tr.Finish()
	log := new(po.TFlowNodeQueryLog)
	genValue := valueGenerator()
	result := db.GetDB().Where("flow_id = ? and node_name=? and query_type = ?", flowId, nodeName, queryType).Order("id desc").First(log)
	if result.RowsAffected <= 0 {
		db.GetDB().Save(&po.TFlowNodeQueryLog{
			FLowId:     flowId,
			NodeName:   nodeName,
			QueryType:  queryType,
			Data:       genValue,
			CreateTime: time.Now(),
			UpdateTime: time.Now(),
		})
		return genValue
	}
	return log.Data
}

func GetFlowNodeQueryLogRepository() IFlowNodeQueryLogRepository {
	return flowNodeQueryLogRepository
}

func init() {
	flowNodeQueryLogRepository = new(FlowNodeQueryLogRepository)
}
