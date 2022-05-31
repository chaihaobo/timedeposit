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
		startNode := unKnowNode.(*node.StartNode)
		startNode.Node = new(node.Node)
		return startNode
	case *node.EndNode:
		endNode := unKnowNode.(*node.EndNode)
		endNode.Node = new(node.Node)
		return endNode
	case *node.UndoMaturityNode:
		realNode := unKnowNode.(*node.UndoMaturityNode)
		realNode.Node = new(node.Node)
		return realNode
	case *node.StartNewMaturityNode:
		realNode := unKnowNode.(*node.StartNewMaturityNode)
		realNode.Node = new(node.Node)
		return realNode
	case *node.ApplyProfitNode:
		realNode := unKnowNode.(*node.ApplyProfitNode)
		realNode.Node = new(node.Node)
		return realNode
	case *node.WithdrawNetprofitNode:
		realNode := unKnowNode.(*node.WithdrawNetprofitNode)
		realNode.Node = new(node.Node)
		return realNode
	case *node.DepositNetprofitNode:
		realNode := unKnowNode.(*node.DepositNetprofitNode)
		realNode.Node = new(node.Node)
		return realNode
	case *node.WithdrawBalanceNode:
		realNode := unKnowNode.(*node.WithdrawBalanceNode)
		realNode.Node = new(node.Node)
		return realNode
	case *node.DepositBalanceNode:
		realNode := unKnowNode.(*node.DepositBalanceNode)
		realNode.Node = new(node.Node)
		return realNode
	case *node.WithdrawAdditionalProfitNode:
		realNode := unKnowNode.(*node.WithdrawAdditionalProfitNode)
		realNode.Node = new(node.Node)
		return realNode
	case *node.DepositAdditionalProfitNode:
		realNode := unKnowNode.(*node.DepositAdditionalProfitNode)
		realNode.Node = new(node.Node)
		return realNode
	case *node.PatchAccountNode:
		realNode := unKnowNode.(*node.PatchAccountNode)
		realNode.Node = new(node.Node)
		return realNode
	case *node.CloseAccountNode:
		realNode := unKnowNode.(*node.CloseAccountNode)
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
	)

}
