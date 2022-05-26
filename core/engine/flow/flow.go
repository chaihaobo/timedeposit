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
	startNode := &node.StartNode{Node: &node.Node{}}
	endNode := &node.EndNode{Node: &node.Node{}}
	registerNode(startNode, endNode)
}
