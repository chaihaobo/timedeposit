/*
 * @Author: Hugo
 * @Date: 2022-05-07 01:48:47
 * @Last Modified by: Hugo
 * @Last Modified time: 2022-05-07 03:52:14
 */
package dao

import (
	"fmt"

	"gitlab.com/hugo.hu/time-deposit-eod-engine/common/db"
	"gitlab.com/hugo.hu/time-deposit-eod-engine/model"
)

func GetProcessFlowByName(flowName string) ([]model.TFlowNode, []model.TFlowNodeRelation) {
	db := db.GetDB()
	var flowNodes []model.TFlowNode
	result := db.Where("flow_name = ?", flowName).Find(&flowNodes)
	fmt.Println(result.Error)
	fmt.Println(result.RowsAffected)

	var flowNodeRelation []model.TFlowNodeRelation
	result = db.Where("flow_name = ?", flowName).Find(&flowNodeRelation)
	fmt.Println(result.Error)
	fmt.Println(result.RowsAffected)
	return flowNodes, flowNodeRelation
}
