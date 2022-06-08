// Package node
// @author： Boice
// @createTime：2022/5/26 17:24
package node

import (
	"fmt"
	"gitlab.com/bns-engineering/td/common/logger"
	"reflect"
	"testing"

	"gitlab.com/bns-engineering/td/common/config"
	"go.uber.org/zap"
)

func TestUndoMaturityNode_Run(t *testing.T) {

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
			name: "Undo maturity test: 11747126703",
			fields: fields{
				Node: &Node{
					FlowId:    "testFlowID",
					AccountId: "11747126703",
					NodeName:  "undo_maturity_node",
				},
			},
			want:    ResultSuccess,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node := &UndoMaturityNode{
				Node: tt.fields.Node,
			}
			got, err := node.Run()
			if (err != nil) != tt.wantErr {

				return
			}
			if !reflect.DeepEqual(got, tt.want) {

				return
			}
			fmt.Println("Run Undo Maturity Finished! result:", got)
		})
	}
}
