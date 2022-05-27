// Package node
// @author： Boice
// @createTime：2022/5/26 17:59
package node

import (
	"reflect"
	"testing"
)

func TestStartNewMaturityNode_Run(t *testing.T) {
	type fields struct {
		Node *Node
	}
	tests := []struct {
		name    string
		fields  fields
		want    INodeResult
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name : "Undo maturity test: 11747126703",
			fields: fields{
				Node : &Node{
					FlowId    : "testFlowID",
					AccountId : "11747126703",
					NodeName  : "undo_maturity_node",
				},
			},
			want :NodeResultSuccess,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node := &StartNewMaturityNode{
				Node: tt.fields.Node,
			}
			got, err := node.Run()
			if (err != nil) != tt.wantErr {
				t.Errorf("StartNewMaturityNode.Run() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("StartNewMaturityNode.Run() = %v, want %v", got, tt.want)
			}
		})
	}
}
