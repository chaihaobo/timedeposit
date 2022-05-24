/*
 * @Author: Hugo
 * @Date: 2022-05-16 04:16:26
 * @Last Modified by: Hugo
 * @Last Modified time: 2022-05-23 11:11:03
 */
package timeDepositNode

import (
	"errors"
	"fmt"

	"gitlab.com/bns-engineering/td/common/constant"
	"gitlab.com/bns-engineering/td/common/log"
	"gitlab.com/bns-engineering/td/node"
	"gitlab.com/bns-engineering/td/service/mambuEntity"
	mambuservices "gitlab.com/bns-engineering/td/service/mambuServices"
)

//Close this TD Account
type CloseAccNode struct {
	node.Node
	// nodeName string
}

func NewCloseAccNode() *CloseAccNode {
	tmpNode := new(CloseAccNode)
	// tmpNode.nodeName = "close_account_node"
	tmpNode.Node.NodeRun = tmpNode
	return tmpNode
}

func (node *CloseAccNode) Process() {
	node.RunNode("close_account_node")
}

func (node *CloseAccNode) RunProcess(tmpTDAccount mambuEntity.TDAccount, flowID string, nodeName string) (constant.FlowNodeStatus, error) {
	totalBalance := tmpTDAccount.Balances.TotalBalance
	if (tmpTDAccount.IsCaseB3() && totalBalance > 0) ||
		(tmpTDAccount.IsCaseC() && totalBalance > 0) {
		notes := fmt.Sprintf("AccountNo:%v, FlowID:%v", tmpTDAccount.ID, flowID)
		isApplySucceed := mambuservices.CloseAccount(tmpTDAccount.ID, notes)
		if !isApplySucceed {
			log.Log.Error("close account failed for account: %v", tmpTDAccount.ID)
			return constant.FlowNodeFailed, errors.New("call Mambu service failed")
		} else {
			log.Log.Info("Finish close account for account: %v", tmpTDAccount.ID)
			return constant.FlowNodeFinish, nil
		}
	} else {
		log.Log.Info("No need to close account, accNo: %v", flowID)
		return constant.FlowNodeSkip, nil
	}
}
