/*
 * @Author: Hugo
 * @Date: 2022-05-16 04:14:01
 * @Last Modified by: Hugo
 * @Last Modified time: 2022-05-23 11:11:54
 */
package timeDepositNode

import (
	"errors"
	"fmt"
	"go.uber.org/zap"

	"gitlab.com/bns-engineering/td/common/constant"
	"gitlab.com/bns-engineering/td/node"
	"gitlab.com/bns-engineering/td/service/mambuEntity"
	mambuservices "gitlab.com/bns-engineering/td/service/mambuServices"
)

//Calc the Additional Profit for TD Account
type ProfitApplyNode struct {
	node.Node
	// nodeName string
}

func NewProfitApplyNode() *ProfitApplyNode {
	tmpNode := new(ProfitApplyNode)
	// tmpNode.nodeName = "profit_apply_node"
	tmpNode.Node.NodeRun = tmpNode
	return tmpNode
}

func (node *ProfitApplyNode) Process() {
	node.RunNode("profit_apply_node")
}

// Update maturity date for this account
func (node *ProfitApplyNode) RunProcess(tmpTDAccount mambuEntity.TDAccount, flowID string, nodeName string) (constant.FlowNodeStatus, error) {
	// Skip updating maturity date or not?
	if !tmpTDAccount.IsCaseB() {
		zap.L().Info(fmt.Sprintf("No need to apply profit, accNo: %v", tmpTDAccount.ID))
		return constant.FlowNodeSkip, nil
	}

	isApplySucceed := mambuservices.ApplyProfit(tmpTDAccount.ID, flowID)
	if !isApplySucceed {
		zap.L().Error(fmt.Sprintf("Apply profit failed for account: %v", tmpTDAccount.ID))
		return constant.FlowNodeFailed, errors.New("call Mambu service failed")
	}
	return constant.FlowNodeFinish, nil
}
