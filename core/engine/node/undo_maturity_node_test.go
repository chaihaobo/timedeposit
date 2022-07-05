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
	carbonMaturityDate, _ := carbon.Parse(carbon.RFC3339Format, "2022-04-16T00:00:00+07:00", "")
	carbonActivationDate, _ := carbon.Parse(carbon.RFC3339Format, "2022-03-15T07:00:00+07:00", "")
	carbonActivationDate = carbonActivationDate.StartOfDay()
	diffInMonths := carbonMaturityDate.DiffInMonths(carbonActivationDate, true)
	if carbonActivationDate.Day() == 31 && (carbonMaturityDate.Day() == 30 || (carbonMaturityDate.Month() == 2 && carbonMaturityDate.LastDayOfMonth().Day() == carbonMaturityDate.Day())) {
		diffInMonths++
	}
	// 正常情况下 == 1 如果

	resultDate := carbonActivationDate.AddMonthsNoOverflow(int(diffInMonths) + 1)
	println(resultDate.DateString())
}
