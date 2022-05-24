/*
 * @Author: Hugo
 * @Date: 2022-05-07 01:48:47
 * @Last Modified by: Hugo
 * @Last Modified time: 2022-05-07 03:52:14
 */
package dao

import (
	"go.uber.org/zap"

	"gitlab.com/bns-engineering/td/common/db"
	"gitlab.com/bns-engineering/td/model"
)

func GetProcessFlowByName(flowName string) ([]model.TFlowNode, []model.TFlowNodeRelation) {
	db := db.GetDB()
	var flowNodes []model.TFlowNode
	result := db.Where("flow_name = ?", flowName).Find(&flowNodes)
	zap.L().Info("flow node result", zap.Int64("RowsAffected", result.RowsAffected), zap.Error(result.Error))

	var flowNodeRelation []model.TFlowNodeRelation
	result = db.Where("flow_name = ?", flowName).Find(&flowNodeRelation)
	zap.L().Info("flow node relation", zap.Int64("RowsAffected", result.RowsAffected), zap.Error(result.Error))
	return flowNodes, flowNodeRelation
}
