/*
 * @Author: Hugo
 * @Date: 2022-04-29 10:24:31
 * @Last Modified by: Hugo
 * @Last Modified time: 2022-05-18 05:03:52
 */
package node

import (
	"gitlab.com/hugo.hu/time-deposit-eod-engine/common/constant"
	"gitlab.com/hugo.hu/time-deposit-eod-engine/dao"
	"gitlab.com/hugo.hu/time-deposit-eod-engine/model"
	"gitlab.com/hugo.hu/time-deposit-eod-engine/service/mambuEntity"
)

type Node struct {
	FlowTaskInfo model.TFlowTaskInfo
	Input        <-chan mambuEntity.TDAccount // input port
	Output       chan<- mambuEntity.TDAccount // output port
}

type NodeRun interface {
	Process()
}

func (*Node) UpdateLogWhenSkipNode(tmpFlowTask model.TFlowTaskInfo, CurNodeName string, nodeLog model.TFlowNodeLog) {
	tmpFlowTask.CurNodeName = CurNodeName
	tmpFlowTask.CurStatus = constant.FlowNodeSkip
	dao.UpdateFlowTask(tmpFlowTask)

	nodeLog.NodeResult = constant.FlowNodeSkip
	dao.UpdateFlowNodeLog(nodeLog)
}

// Get the input data and registed flow log info for this node
func (node *Node) GetAccAndFlowLog(CurNodeName string) (mambuEntity.TDAccount, model.TFlowTaskInfo, model.TFlowNodeLog) {
	tmpTDAccount := <-node.Input
	tmpFlowTask := node.FlowTaskInfo
	tmpFlowTask.CurNodeName = CurNodeName
	tmpFlowTask.CurStatus = constant.FlowNodeStart
	tmpFlowTask.FlowStatus = constant.FlowRunning
	dao.UpdateFlowTask(tmpFlowTask)

	nodeLog := dao.CreateFlowNodeLog(tmpFlowTask.FlowId, tmpTDAccount.ID, tmpFlowTask.FlowName, CurNodeName)
	return tmpTDAccount, tmpFlowTask, nodeLog
}

func (node *Node) UpdateLogWhenNodeFailed(tmpFlowTask model.TFlowTaskInfo, nodeLog model.TFlowNodeLog, err error) {
	tmpFlowTask.CurStatus = constant.FlowNodeFailed
	tmpFlowTask.FlowStatus = constant.FlowFailed
	dao.UpdateFlowTask(tmpFlowTask)

	nodeLog.NodeResult = constant.FlowNodeFailed
	nodeLog.NodeMsg = err.Error()
	dao.UpdateFlowNodeLog(nodeLog)
}

func (node *Node) UpdateLogWhenNodeFinish(tmpFlowTask model.TFlowTaskInfo, nodeLog model.TFlowNodeLog) {
	tmpFlowTask.CurStatus = constant.FlowNodeFinish
	dao.UpdateFlowTask(tmpFlowTask)

	nodeLog.NodeResult = constant.FlowNodeFinish
	dao.UpdateFlowNodeLog(nodeLog)
}
