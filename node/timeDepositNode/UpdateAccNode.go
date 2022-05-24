/*
 * @Author: Hugo
 * @Date: 2022-05-16 04:16:12
 * @Last Modified by: Hugo
 * @Last Modified time: 2022-05-23 11:12:13
 */
package timeDepositNode

import (
	"errors"

	"gitlab.com/bns-engineering/td/common/constant"
	"gitlab.com/bns-engineering/td/common/log"
	"gitlab.com/bns-engineering/td/common/util"
	"gitlab.com/bns-engineering/td/node"
	"gitlab.com/bns-engineering/td/service/mambuEntity"
	mambuservices "gitlab.com/bns-engineering/td/service/mambuServices"
)

type UpdateAccNode struct {
	node.Node
	// nodeName string
}

func NewUpdateAccNode() *UpdateAccNode {
	tmpNode := new(UpdateAccNode)
	// tmpNode.nodeName = "update_account_node"
	tmpNode.Node.NodeRun = tmpNode
	return tmpNode
}

func (node *UpdateAccNode) Process() {
	node.RunNode("update_account_node")
}

func (node *UpdateAccNode) RunProcess(tmpTDAccount mambuEntity.TDAccount, flowID string, nodeName string) (constant.FlowNodeStatus, error) {
	if tmpTDAccount.IsCaseB1_1() || tmpTDAccount.IsCaseB2() {
		newDate := util.GetDate(tmpTDAccount.MaturityDate)
		isApplySucceed := mambuservices.UpdateMaturifyDateForTDAccount(tmpTDAccount.ID, newDate)
		if !isApplySucceed {
			log.Log.Error("Apply profit failed for account: %v", tmpTDAccount.ID)
			return constant.FlowNodeFailed, errors.New("call mambu service failed")
		} else {
			log.Log.Info("Finish apply profit for account: %v", tmpTDAccount.ID)
			return constant.FlowNodeFinish, nil
		}
	} else {
		log.Log.Info("No need to update maturity info for td account, accNo: %v", tmpTDAccount.ID)
		return constant.FlowNodeSkip, nil
	}
}
