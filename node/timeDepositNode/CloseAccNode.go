/*
 * @Author: Hugo
 * @Date: 2022-05-16 04:16:26
 * @Last Modified by: Hugo
 * @Last Modified time: 2022-05-18 05:07:11
 */
package timeDepositNode

import (
	"errors"
	"fmt"

	"gitlab.com/bns-engineering/td/common/log"
	"gitlab.com/bns-engineering/td/node"
	mambuservices "gitlab.com/bns-engineering/td/service/mambuServices"
)

//Close this TD Account
type CloseAccNode struct {
	node.Node
}

func (node *CloseAccNode) Process() {
	CurNodeName := "close_account_node"
	tmpTDAccount, tmpFlowTask, nodeLog := node.GetAccAndFlowLog(CurNodeName)
	totalBalance := tmpTDAccount.Balances.TotalBalance
	if !(tmpTDAccount.IsCaseB3() && totalBalance > 0) &&
		!(tmpTDAccount.IsCaseC() && totalBalance > 0) {
		node.UpdateLogWhenSkipNode(tmpFlowTask, CurNodeName, nodeLog)
		log.Log.Info("No need to close account, accNo: %v", tmpFlowTask.FlowId)
	} else {
		notes := fmt.Sprintf("AccountNo:%v, FlowID:%v", tmpTDAccount.ID, tmpFlowTask.FlowId)
		isApplySucceed := mambuservices.CloseAccount(tmpTDAccount.ID, notes)
		if !isApplySucceed {
			log.Log.Error("close account failed for account: %v", tmpTDAccount.ID)
			node.UpdateLogWhenNodeFailed(tmpFlowTask, nodeLog, errors.New("call Mambu service failed"))
		} else {
			log.Log.Info("Finish close account for account: %v", tmpTDAccount.ID)
			node.UpdateLogWhenNodeFinish(tmpFlowTask, nodeLog)
		}
	}
	node.Node.Output <- tmpTDAccount
}
