// Package node

// @author： Boice

// @createTime：2022/5/27 09:18

package node

import (
	"reflect"
	"testing"
	
	"gitlab.com/bns-engineering/td/common/config"
	"gitlab.com/bns-engineering/td/common/log"
	"go.uber.org/zap"
)

func TestWithdrawAdditionalProfitNode_Run(t *testing.T) {
	config.Setup("./../../../config.yaml")
	err := logger.SetUp(config.TDConf)
	if err != nil {
		zap.L().Error("logger init error", zap.Error(err))
	}

	tests := []struct {
		name    string
		node    *WithdrawAdditionalProfitNode
		want    INodeResult
		wantErr bool
	}{
		{
			name: "Withdraw additional profit test: 11645631879",
			node: &WithdrawAdditionalProfitNode{
				Node:&Node{
					FlowId:    "testFlowID_11645631879_1",
					AccountId: "11645631879",
					NodeName:  "withdraw_additional_profit_node",
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
				t.Errorf("WithdrawAdditionalProfitNode.Run() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithdrawAdditionalProfitNode.Run() = %v, want %v", got, tt.want)
			}
		})
	}
}
