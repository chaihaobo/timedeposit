// Package node

// @author： Boice

// @createTime：2022/5/26 11:07

package node

import (
	"gitlab.com/bns-engineering/td/common/config"
	"gitlab.com/bns-engineering/td/common/logger"
	"go.uber.org/zap"
	"reflect"
	"testing"
)

func TestWithdrawBalanceNode_Run(t *testing.T) {
	config.Setup("./../../../config.json")
	err := logger.SetUp(config.TDConf)
	if err != nil {
		zap.L().Error("logger init error", zap.Error(err))
	}

	tests := []struct {
		name    string
		node    *WithdrawBalanceNode
		want    INodeResult
		wantErr bool
	}{
		{
			name: "Withdraw total balance test: 11645631879",
			node: &WithdrawBalanceNode{
				Node: &Node{
					FlowId:    "testFlowID_11645631879_2",
					AccountId: "11645631879",
					NodeName:  "withdraw_balance_node",
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
