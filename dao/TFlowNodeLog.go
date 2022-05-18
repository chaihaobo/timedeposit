/*
 * @Author: Hugo
 * @Date: 2022-05-18 04:34:21
 * @Last Modified by: Hugo
 * @Last Modified time: 2022-05-18 04:39:54
 */
package dao

import (
	"time"

	"gitlab.com/bns-engineering/td/common/constant"
	"gitlab.com/bns-engineering/td/common/db"
	"gitlab.com/bns-engineering/td/model"
)

func CreateFlowNodeLog(flowId, accountId, flowName, nodeName string) model.TFlowNodeLog {
	tFlowTask := model.TFlowNodeLog{
		AccountId:  accountId,
		FlowId:     flowId,
		FlowName:   flowName,
		NodeName:   nodeName,
		NodeResult: constant.FlowNodeFailed,
		NodeMsg:    "",
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
	}
	db := db.GetDB()
	db.Save(&tFlowTask)
	return tFlowTask
}

func UpdateFlowNodeLog(nodeLog model.TFlowNodeLog) {
	db := db.GetDB()
	nodeLog.UpdateTime = time.Now()
	db.Save(nodeLog)
}
