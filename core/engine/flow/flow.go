// Package engine
// @author： Boice
// @createTime：2022/5/26 10:14
package flow

import (
	"gitlab.com/bns-engineering/td/core/engine/node"
	"reflect"
)

var nodeList = make(map[string]interface{})

func GetNode(nodeName string) node.INode {
	unKnowNode := nodeList[nodeName]
	switch unKnowNode.(type) {
	case *node.StartNode:
		startNode := new(node.StartNode)
		startNode.Node = new(node.Node)
		return startNode
	case *node.EndNode:
		endNode := new(node.EndNode)
		endNode.Node = new(node.Node)
		return endNode
	case *node.UndoMaturityNode:
		realNode := new(node.UndoMaturityNode)
		realNode.Node = new(node.Node)
		return realNode
	case *node.StartNewMaturityNode:
		realNode := new(node.StartNewMaturityNode)
		realNode.Node = new(node.Node)
		return realNode
	case *node.ApplyProfitNode:
		realNode := new(node.ApplyProfitNode)
		realNode.Node = new(node.Node)
		return realNode
	case *node.WithdrawNetprofitNode:
		realNode := new(node.WithdrawNetprofitNode)
		realNode.Node = new(node.Node)
		return realNode
	case *node.DepositNetprofitNode:
		realNode := new(node.DepositNetprofitNode)
		realNode.Node = new(node.Node)
		return realNode
	case *node.WithdrawBalanceNode:
		realNode := new(node.WithdrawBalanceNode)
		realNode.Node = new(node.Node)
		return realNode
	case *node.DepositBalanceNode:
		realNode := new(node.DepositBalanceNode)
		realNode.Node = new(node.Node)
		return realNode
	case *node.WithdrawAdditionalProfitNode:
		realNode := new(node.WithdrawAdditionalProfitNode)
		realNode.Node = new(node.Node)
		return realNode
	case *node.DepositAdditionalProfitNode:
		realNode := new(node.DepositAdditionalProfitNode)
		realNode.Node = new(node.Node)
		return realNode
	case *node.PatchAccountNode:
		realNode := new(node.PatchAccountNode)
		realNode.Node = new(node.Node)
		return realNode
	case *node.CloseAccountNode:
		realNode := new(node.CloseAccountNode)
		realNode.Node = new(node.Node)
		return realNode
	case *node.AdditionalProfitNode:
		realNode := new(node.AdditionalProfitNode)
		realNode.Node = new(node.Node)
		return realNode
	default:
		return nil
	}

}

func register(list ...interface{}) {
	for _, nodeObj := range list {
		nodeName := reflect.TypeOf(nodeObj).Elem().Name()
		nodeList[nodeName] = nodeObj
	}
}

func SetUp() {
	register(
		new(node.StartNode),
		new(node.EndNode),
		new(node.UndoMaturityNode),
		new(node.StartNewMaturityNode),
		new(node.ApplyProfitNode),
		new(node.WithdrawNetprofitNode),
		new(node.DepositNetprofitNode),
		new(node.WithdrawBalanceNode),
		new(node.DepositBalanceNode),
		new(node.WithdrawAdditionalProfitNode),
		new(node.DepositAdditionalProfitNode),
		new(node.PatchAccountNode),
		new(node.CloseAccountNode),
		new(node.AdditionalProfitNode),
	)

}
