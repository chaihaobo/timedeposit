/*
 * @Author: Hugo
 * @Date: 2022-05-16 04:14:01
 * @Last Modified by: Hugo
 * @Last Modified time: 2022-05-18 05:06:05
 */
package timeDepositNode

import (
	"errors"

	"gitlab.com/bns-engineering/td/common/log"
	"gitlab.com/bns-engineering/td/node"
	mambuservices "gitlab.com/bns-engineering/td/service/mambuServices"
)

//Calc the Additional Profit for TD Account
type ProfitApplyNode struct {
	node.Node
}

func (node *ProfitApplyNode) Process() {
	CurNodeName := "profit_apply_node"
	tmpTDAccount, tmpFlowTask, nodeLog := node.GetAccAndFlowLog(CurNodeName)
	if !tmpTDAccount.IsCaseB() {
		node.UpdateLogWhenSkipNode(tmpFlowTask, CurNodeName, nodeLog)
		log.Log.Info("No need to apply profit, accNo: %v", tmpFlowTask.FlowId)
	} else {
		isApplySucceed := mambuservices.ApplyProfit(tmpTDAccount.ID, tmpFlowTask.FlowId)
		if !isApplySucceed {
			log.Log.Error("Apply profit failed for account: %v", tmpTDAccount.ID)
			node.UpdateLogWhenNodeFailed(tmpFlowTask, nodeLog, errors.New("call Mambu service failed"))
		} else {
			log.Log.Info("Finish apply profit for account: %v", tmpTDAccount.ID)
			node.UpdateLogWhenNodeFinish(tmpFlowTask, nodeLog)
		}
	}
	node.Node.Output <- tmpTDAccount
}
