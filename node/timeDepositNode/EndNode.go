/*
 * @Author: Hugo
 * @Date: 2022-04-29 10:24:39
 * @Last Modified by: Hugo
 * @Last Modified time: 2022-05-24 01:19:23
 */
package timeDepositNode

import (
	"gitlab.com/bns-engineering/td/common/constant"
	"gitlab.com/bns-engineering/td/node"
	"gitlab.com/bns-engineering/td/service/mambuEntity"
)

//End Node
type EndNode struct {
	// Input <-chan string // input port
	// Input <-chan node.NodeData
	node.Node
}

func NewEndNode() *EndNode {
	tmpNode := new(EndNode)
	tmpNode.Name = constant.EndNode
	tmpNode.Node.NodeRun = tmpNode
	return tmpNode
}

// Update maturity date for this account
func (node *EndNode) RunProcess(tmpTDAccount mambuEntity.TDAccount, flowID string, nodeName string) (constant.FlowNodeStatus, error) {
	return constant.FlowNodeFinish, nil
}
