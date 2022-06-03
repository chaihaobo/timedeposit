// Package node

// @author： Boice

// @createTime：2022/5/26 18:08

package node

import (
	"gitlab.com/bns-engineering/td/common/config"
	"gitlab.com/bns-engineering/td/common/logger"
	"go.uber.org/zap"
	"reflect"
	"testing"
)

func TestCloseAccountNode_Run(t *testing.T) {
	config.Setup("./../../../config.json")
	err := logger.SetUp(config.TDConf)
	if err != nil {
		zap.L().Error("logger init error", zap.Error(err))
	}
	tests := []struct {
		name    string
		node    *CloseAccountNode
		want    INodeResult
		wantErr bool
	}{
		{
			name: "Close account test: 11645631879",
			node: &CloseAccountNode{
				Node: &Node{
					FlowId:    "testFlowID_11645631879_2",
					AccountId: "11645631879",
					NodeName:  "close_account_node",
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
				t.Errorf("CloseAccountNode.Run() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CloseAccountNode.Run() = %v, want %v", got, tt.want)
			}
		})
	}
}
