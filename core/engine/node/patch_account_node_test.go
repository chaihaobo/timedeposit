// Package node

// @author： Boice

// @createTime：2022/5/26 18:08

package node

import (
	"reflect"
	"testing"

	"gitlab.com/bns-engineering/td/common/config"
	"gitlab.com/bns-engineering/td/common/logger"
	"go.uber.org/zap"
)

func TestPatchAccountNode_Run(t *testing.T) {
	config.Setup("./../../../config.json")
	err := logger.SetUp(config.TDConf)
	if err != nil {
		zap.L().Error("logger init error", zap.Error(err))
	}
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
			got, err := tt.node.Run()
			if (err != nil) != tt.wantErr {

				return
			}
			if !reflect.DeepEqual(got, tt.want) {

			}
		})
	}
}
