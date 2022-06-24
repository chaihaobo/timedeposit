// Package node
// @author： Boice
// @createTime：2022/6/9 10:41
package node

import (
	"context"
	"gitlab.com/bns-engineering/td/common/config"
	"gitlab.com/bns-engineering/td/transport"
	"testing"
)

func init() {
	transport.NewTdServer(config.Setup("../../../../config.json")).SetUp()
}

func TestAdditionalProfitNode_Run(t *testing.T) {
	node := &AdditionalProfitNode{
		Node: &Node{
			FlowId:    "test_flow_111695044011",
			AccountId: "11312188579",
			NodeName:  "additional_profit_node",
		},
	}
	node.Run(context.Background())
}
