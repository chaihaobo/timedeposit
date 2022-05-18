/*
 * @Author: Hugo
 * @Date: 2022-04-29 10:24:35
 * @Last Modified by: Hugo
 * @Last Modified time: 2022-05-18 04:59:36
 */
package timeDepositNode

import (
	"gitlab.com/hugo.hu/time-deposit-eod-engine/common/constant"
	"gitlab.com/hugo.hu/time-deposit-eod-engine/common/log"
	"gitlab.com/hugo.hu/time-deposit-eod-engine/dao"
	"gitlab.com/hugo.hu/time-deposit-eod-engine/node"
	mambuservices "gitlab.com/hugo.hu/time-deposit-eod-engine/service/mambuServices"
)

//AA Time Deposit Engine IWT start Time Deposit
type StartNode struct {
	node.Node
}

// In start node, will try to get the detail info of this td account.
func (node *StartNode) Process() {
	CurNodeName := "start_node"

	tmpTDAccount, tmpFlowTask, nodeLog := node.GetAccAndFlowLog(CurNodeName)
	tdAccInfo, err := mambuservices.GetTDAccountById(tmpTDAccount.ID)
	if err != nil {
		log.Log.Info("Get TDAcc Error!")
		tmpFlowTask.CurNodeName = CurNodeName
		tmpFlowTask.CurStatus = constant.FlowNodeFailed
		tmpFlowTask.FlowStatus = constant.FlowFailed
		dao.UpdateFlowTask(tmpFlowTask)

		nodeLog.NodeResult = constant.FlowNodeFailed
		dao.UpdateFlowNodeLog(nodeLog)
		return
	}
	log.Log.Info("StartNode get the full info of this account, OutputData: %v", tdAccInfo)
	//Update Node Status
	nodeLog.NodeResult = constant.FlowNodeFinish
	dao.UpdateFlowNodeLog(nodeLog)

	tmpFlowTask.CurNodeName = CurNodeName
	tmpFlowTask.CurStatus = constant.FlowNodeFinish
	tmpFlowTask.FlowStatus = constant.FlowRunning
	dao.UpdateFlowTask(tmpFlowTask)
	node.Node.Output <- tdAccInfo
}
