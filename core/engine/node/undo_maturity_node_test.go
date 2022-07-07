// Package node
// @author： Boice
// @createTime：2022/5/26 17:24
package node

import (
	"context"
	"fmt"
	"github.com/uniplaces/carbon"
	"gitlab.com/bns-engineering/td/common/util"
	"reflect"
	"testing"
	"time"

	"gitlab.com/bns-engineering/td/common/config"
)

func TestUndoMaturityNode_Run(t *testing.T) {

	util.SetupTelemetry(config.Setup("./../../../config.json"))

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

func TestName(t *testing.T) {
	zone := time.FixedZone("CST", 7*3600)
	time.Local = zone
	parse, _ := carbon.Parse(carbon.DateFormat, "2022-01-31", "")
	println(parse.AddMonthsNoOverflow(1).DateString())

}
