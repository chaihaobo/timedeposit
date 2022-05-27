// Package node
// @author： Boice
// @createTime：2022/5/26 18:08
package node

import (
	"reflect"
	"testing"
	"gitlab.com/bns-engineering/td/common/config"
	"gitlab.com/bns-engineering/td/common/log"
	"go.uber.org/zap"
)

func TestApplyProfitNode_Run(t *testing.T) {
	
	config.Setup("./../../../config.yaml")
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
			name : "apply profit test: 11714744288",
			fields: fields{
				Node : &Node{
					FlowId    : "testFlowID",
					AccountId : "11714744288",
					NodeName  : "apply_profit_node",
				},
			},
			want :NodeResultSuccess,
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
