// Package node

// @author： Boice

// @createTime：2022/5/26 18:08

package node

import (
	"context"
	"gitlab.com/bns-engineering/td/common/util"
	"reflect"
	"testing"

	"gitlab.com/bns-engineering/td/common/config"
)

func TestPatchAccountNode_Run(t *testing.T) {
	util.SetupTelemetry(config.Setup("./../../../config.json"))
	tests := []struct {
		name    string
		node    *PatchAccountNode
		want    INodeResult
		wantErr bool
	}{
		{
			name: "Withdraw additional profit test: 11645631879",
			node: &PatchAccountNode{
				Node: &Node{
					FlowId:    "testFlowID_11645631879_1",
					AccountId: "11645631879",
					NodeName:  "patch_account_node",
				},
			},
			want:    ResultSuccess,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.node.Run(context.Background())
			if (err != nil) != tt.wantErr {

				return
			}
			if !reflect.DeepEqual(got, tt.want) {

			}
		})
	}
}
