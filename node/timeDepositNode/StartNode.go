/*
 * @Author: Hugo
 * @Date: 2022-04-29 10:24:35
 * @Last Modified by: Hugo
 * @Last Modified time: 2022-05-23 11:10:14
 */
package timeDepositNode

import (
	"gitlab.com/bns-engineering/td/common/constant"
	"gitlab.com/bns-engineering/td/node"
	"gitlab.com/bns-engineering/td/service/mambuEntity"
)

//AA Time Deposit Engine IWT start Time Deposit
type StartNode struct {
	node.Node
}

// Update maturity date for this account
func (node *StartNode) RunProcess(tmpTDAccount mambuEntity.TDAccount, flowID string, nodeName string) (constant.FlowNodeStatus, error) {
	return constant.FlowNodeFinish, nil
}

func NewStartNode() *StartNode {
	tmpNode := new(StartNode)
	tmpNode.Name = constant.StartNode
	tmpNode.Node.NodeRun = tmpNode
	return tmpNode
}
