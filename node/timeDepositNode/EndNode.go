/*
 * @Author: Hugo
 * @Date: 2022-04-29 10:24:39
 * @Last Modified by: Hugo
 * @Last Modified time: 2022-05-24 01:19:23
 */
package timeDepositNode

import (
	"gitlab.com/bns-engineering/td/common/constant"
	"gitlab.com/bns-engineering/td/dao"
	"gitlab.com/bns-engineering/td/node"
	"gitlab.com/bns-engineering/td/service/mambuEntity"
)

//End Node
type EndNode struct {
	// Input <-chan string // input port
	// Input <-chan node.NodeData
	node.Node
	// nodeName string
}

func NewEndNode() *EndNode {
	tmpNode := new(EndNode)
	// tmpNode.No.nodeName = "end_node"
	return tmpNode
}

// In start node, will try to get the detail info of this td account.
func (tmpNode *EndNode) Process() {
	nodeDataInfo := <-tmpNode.Input
	tmpTDAccount := nodeDataInfo.TDAccountInfo
	tmpFlowTask := nodeDataInfo.FlowTaskInfo

	tmpFlowTask.CurNodeName = "end_node"
	tmpFlowTask.CurStatus = string(constant.FlowNodeFinish)
	tmpFlowTask.FlowStatus = constant.FlowFinished
	dao.UpdateFlowTask(tmpFlowTask)

	nodeLog := dao.CreateFlowNodeLog(tmpFlowTask.FlowId, tmpTDAccount.ID, tmpFlowTask.FlowName, "end_node")

	nodeLog.NodeResult = string(constant.FlowNodeFinish)
	dao.UpdateFlowNodeLog(nodeLog)
}

// Update maturity date for this account
func (node *EndNode) RunProcess(tmpTDAccount mambuEntity.TDAccount, flowID string, nodeName string) (constant.FlowNodeStatus, error) {
	return constant.FlowNodeFinish, nil
}

// func (node *EndNode) Process() {
// 	CurNodeName := "end_node"
// 	_, tmpFlowTask, nodeLog := node.GetAccAndFlowLog(CurNodeName)

// 	node.UpdateLogWhenNodeFinish(tmpFlowTask, nodeLog)
// 	tmpFlowTask.EndStatus = constant.FlowFinished
// 	dao.UpdateFlowTask(tmpFlowTask)
// }
