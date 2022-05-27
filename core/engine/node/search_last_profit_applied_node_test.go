// Package node
// @author： Boice
// @createTime：2022/5/27 11:18
package node

import (
	"testing"
)

func TestSearchLastProfitAppliedNode_Run(t *testing.T) {
	t.Run("run search_last_profit_applied_node.go", func(t *testing.T) {
		node := &SearchLastProfitAppliedNode{Node: &Node{}}
		node.SetUp("test", "11249460359", "test")
		node.Run()
	})

}
