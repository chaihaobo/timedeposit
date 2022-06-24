// Package node
// @author： Boice
// @createTime：2022/5/26 17:24
package node

import (
	"context"
	"fmt"
	"gitlab.com/bns-engineering/td/transport"
	"reflect"
	"testing"

	"gitlab.com/bns-engineering/td/common/config"
)

func TestUndoMaturityNode_Run(t *testing.T) {

	transport.NewTdServer(config.Setup("./../../../config.json")).SetUp()

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
			got, err := node.Run(context.Background())
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
