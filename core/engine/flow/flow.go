// Package engine
// @author： Boice
// @createTime：2022/5/26 10:14
package flow

import (
	"gitlab.com/bns-engineering/td/core/engine/node"
	"reflect"
)

var NodeList = make(map[string]node.INode)

func registerNode(nodeList ...node.INode) {
	for _, iNode := range nodeList {
		iNodeName := reflect.TypeOf(iNode).Elem().Name()
		NodeList[iNodeName] = iNode
	}

}

func SetUp() {
	registerNode(
		&node.StartNode{Node: &node.Node{}},
		&node.EndNode{Node: &node.Node{}},
		&node.UndoMaturityNode{Node: &node.Node{}},
		&node.StartNewMaturityNode{Node: &node.Node{}},
		&node.ApplyProfitNode{Node: &node.Node{}},
		&node.WithdrawNetprofitNode{Node: &node.Node{}},
		&node.DepositNetprofitNode{Node: &node.Node{}},
		&node.WithdrawBalanceNode{Node: &node.Node{}},
		&node.DepositBalanceNode{Node: &node.Node{}},
		&node.PatchAccountNode{Node: &node.Node{}},
		&node.CloseAccountNode{Node: &node.Node{}},
	)
}
