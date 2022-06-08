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

func TestApplyProfitNode_Run(t *testing.T) {

	config.Setup("./../../../config.json")
	err := logger.SetUp(config.TDConf)
	if err != nil {
		zap.L().Error("logger init error", zap.Error(err))
	}

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
					FlowId:    "testFlowID_11645631879_1",
					AccountId: "11645631879",
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
			got, err := node.Run()
			if (err != nil) != tt.wantErr {
				t.Errorf("ApplyProfitNode.Run() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ApplyProfitNode.Run() = %v, want %v", got, tt.want)
			}
		})
	}
}
