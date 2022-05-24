/*
 * @Author: Hugo
 * @Date: 2022-04-29 10:24:31
 * @Last Modified by: Hugo
 * @Last Modified time: 2022-05-24 03:00:01
 */
package node

import (
	"fmt"
	"gitlab.com/bns-engineering/td/common/constant"
	"gitlab.com/bns-engineering/td/dao"
	"gitlab.com/bns-engineering/td/model"
	"gitlab.com/bns-engineering/td/service/mambuEntity"
	"go.uber.org/zap"
)

type NodeData struct {
	FlowTaskInfo  *model.TFlowTaskInfo
	TDAccountInfo mambuEntity.TDAccount // input port
}
type Node struct {
	Input   <-chan NodeData
	Output  chan<- NodeData // output port
	Name    string
	NodeRun NodeRun
}

type NodeRun interface {
	RunProcess(tdAccount mambuEntity.TDAccount, flowID, nodeName string) (constant.FlowNodeStatus, error)
	Process()
}

func (node *Node) Process() {
	nodeDataInfo := <-node.Input
	tmpTDAccount := nodeDataInfo.TDAccountInfo
	tmpFlowTask := nodeDataInfo.FlowTaskInfo

	zap.L().Info(fmt.Sprintf("FlowID: %v, flowCurStatus:%v, flowStatus:%v, CurNodeName:%v", tmpFlowTask.FlowId, tmpFlowTask.CurStatus, tmpFlowTask.FlowStatus, node.Name))

	if tmpFlowTask.CurStatus == string(constant.FlowNodeFailed) {
		node.Output <- nodeDataInfo
		return
	}

	tmpFlowTask.CurNodeName = node.Name
	tmpFlowTask.CurStatus = string(constant.FlowNodeStart)
	tmpFlowTask.FlowStatus = constant.FlowRunning
	dao.UpdateFlowTask(tmpFlowTask)

	nodeLog := dao.CreateFlowNodeLog(tmpFlowTask.FlowId, tmpTDAccount.ID, tmpFlowTask.FlowName, node.Name)

	tmpNodeStatus, err := node.NodeRun.RunProcess(tmpTDAccount, tmpFlowTask.FlowId, node.Name)
	if node.Name != "end_node" {
		switch tmpNodeStatus {
		case constant.FlowNodeSkip:
			node.UpdateLogWhenSkipNode(tmpFlowTask, nodeLog)
		case constant.FlowNodeFailed:
			node.UpdateLogWhenNodeFailed(tmpFlowTask, nodeLog, err)
			nodeDataInfo.FlowTaskInfo.CurStatus = string(constant.FlowNodeFailed)
		case constant.FlowNodeFinish:
			node.UpdateLogWhenNodeFinish(tmpFlowTask, nodeLog)
		}
		node.Output <- nodeDataInfo
	} else {
		node.UpdateLogWhenNodeFinish(tmpFlowTask, nodeLog)
		tmpFlowTask.EndStatus = constant.FlowFinished
		tmpFlowTask.FlowStatus = constant.FlowFinished
		dao.UpdateFlowTask(tmpFlowTask)
	}
}

func (*Node) UpdateLogWhenSkipNode(tmpFlowTask *model.TFlowTaskInfo, nodeLog model.TFlowNodeLog) {
	tmpFlowTask.CurStatus = string(constant.FlowNodeSkip)
	dao.UpdateFlowTask(tmpFlowTask)

	nodeLog.NodeResult = string(constant.FlowNodeSkip)
	dao.UpdateFlowNodeLog(nodeLog)
}

func (node *Node) UpdateLogWhenNodeFailed(tmpFlowTask *model.TFlowTaskInfo, nodeLog model.TFlowNodeLog, err error) {
	tmpFlowTask.CurStatus = string(constant.FlowNodeFailed)
	tmpFlowTask.FlowStatus = constant.FlowFailed
	dao.UpdateFlowTask(tmpFlowTask)

	nodeLog.NodeResult = string(constant.FlowNodeFailed)
	nodeLog.NodeMsg = err.Error()
	dao.UpdateFlowNodeLog(nodeLog)
}

func (node *Node) UpdateLogWhenNodeFinish(tmpFlowTask *model.TFlowTaskInfo, nodeLog model.TFlowNodeLog) {
	tmpFlowTask.CurStatus = string(constant.FlowNodeFinish)
	dao.UpdateFlowTask(tmpFlowTask)

	nodeLog.NodeResult = string(constant.FlowNodeFinish)
	dao.UpdateFlowNodeLog(nodeLog)
}
