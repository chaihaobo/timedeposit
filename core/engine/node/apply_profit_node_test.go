// Package node
// @author： Boice
// @createTime：2022/5/26 18:08
package node

import (
	"context"
	"gitlab.com/bns-engineering/td/common/config"
	"gitlab.com/bns-engineering/td/common/util"
	"reflect"
	"testing"
)

func TestApplyProfitNode_Run(t *testing.T) {

	util.SetupTelemetry(config.Setup("./../../../config.json"))
	type fields struct {
		Node *Node
	}
	tests := []struct {
		name    string
		fields  fields
		want    INodeResult
		wantErr bool
	}{
		{
			name: "apply profit test: 11645631879",
			fields: fields{
				Node: &Node{
					FlowId:    "testFlowID_11645631879_12",
					AccountId: "11312188579",
					NodeName:  "apply_profit_node",
				},
			},
			want:    ResultSuccess,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node := &ApplyProfitNode{
				Node: tt.fields.Node,
			}
			got, err := node.Run(context.Background())
			if (err != nil) != tt.wantErr {

				return
			}
			if !reflect.DeepEqual(got, tt.want) {

			}
		})
	}
}
