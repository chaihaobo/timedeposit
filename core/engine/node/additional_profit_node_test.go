// Package node
// @author： Boice
// @createTime：2022/6/9 10:41
package node

import (
	"gitlab.com/bns-engineering/td/common/config"
	"gitlab.com/bns-engineering/td/common/logger"
	"testing"
)

func init() {
	logger.SetUp(config.Setup("../../../config.json"))
}

func TestAdditionalProfitNode_Run(t *testing.T) {
	node := &AdditionalProfitNode{
		Node: &Node{
			FlowId:    "test_flow_111695044011",
			AccountId: "11312188579",
			NodeName:  "additional_profit_node",
		},
	}
	node.Run()

}
