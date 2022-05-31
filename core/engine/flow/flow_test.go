// Package node
// @author： Boice
// @createTime：2022/5/26 12:15
package flow

import (
	"fmt"
	"gitlab.com/bns-engineering/td/core/engine/node"
	"reflect"
	"testing"
)

func TestFlowSetup(t *testing.T) {
	register(&node.StartNode{})
	i := nodeList["StartNode"]
	reflect.TypeOf(i).Elem()
	switch i.(type) {
	case *node.StartNode:
		startNode := i.(*node.StartNode)
		startNode.Node = &node.Node{}
		fmt.Printf("StartNode")

	default:
		fmt.Printf("not match")
	}
}
