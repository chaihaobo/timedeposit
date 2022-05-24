/*
 * @Author: Hugo
 * @Date: 2022-05-18 03:38:56
 * @Last Modified by: Hugo
 * @Last Modified time: 2022-05-18 04:22:32
 */
package dao

import (
	"time"

	"gitlab.com/bns-engineering/td/common/db"
	"gitlab.com/bns-engineering/td/model"
)

func CreateFlowTask(flowId, accountId, flowName string) *model.TFlowTaskInfo {
	tFlowTask := model.TFlowTaskInfo{
		FlowId:      flowId,
		AccountId:   accountId,
		FlowName:    flowName,
		FlowStatus:  "start",
		CurNodeName: "start_node",
		CurStatus:   "start",
		EndStatus:   "",
		StartTime:   time.Now(),
		EndTime:     time.Now(),
		CreateTime:  time.Now(),
		UpdateTime:  time.Now(),
	}
	db := db.GetDB()
	db.Save(&tFlowTask)
	return &tFlowTask
}

func UpdateFlowTask(flowTask *model.TFlowTaskInfo) {
	db := db.GetDB()
	flowTask.UpdateTime = time.Now()
	db.Save(flowTask)
}
