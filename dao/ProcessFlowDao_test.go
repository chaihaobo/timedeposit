/*
 * @Author: Hugo
 * @Date: 2022-05-07 01:48:47
 * @Last Modified by: Hugo
 * @Last Modified time: 2022-05-07 03:50:57
 */
package dao

import (
	"testing"

	"gitlab.com/hugo.hu/time-deposit-eod-engine/model"
)

func TestGetProcessFlowByName(t *testing.T) {
	type args struct {
		flowName string
	}
	tests := []struct {
		name  string
		args  args
		want  []model.TFlowNode
		want1 []model.TFlowNodeRelation
	}{
		// TODO: Add test cases.
		{
			name:  "test TimeDepoistFlow Info",
			args:  args{flowName: "time_deposit_flow"},
			want:  nil,
			want1: nil,
			// want1: []model.TFlowNodeRelation{
			// 	model.TFlowNodeRelation{
			// 		ID:         1,
			// 		FlowName:   "time_deposit_flow",
			// 		NodeName:   "start_node",
			// 		ResultCode: "1",
			// 		NextNode:   "additional_profit_cal_node",
			// 		CreateTime: time.Date(2022, 05, 07, 10, 46, 12, 0, time.UTC),
			// 		UpdateTime: time.Date(2022, 05, 07, 10, 46, 15, 0, time.UTC),
			// 	},
			// 	model.TFlowNodeRelation{
			// 		ID:         2,
			// 		FlowName:   "time_deposit_flow",
			// 		NodeName:   "additional_profit_cal_node",
			// 		ResultCode: "1",
			// 		NextNode:   "end_node",
			// 		CreateTime: time.Date(2022, 05, 07, 10, 46, 41, 0, time.UTC),
			// 		UpdateTime: time.Date(2022, 05, 07, 10, 46, 44, 0, time.UTC),
			// 	},
			// },
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := GetProcessFlowByName(tt.args.flowName)
			if len(got) == 0 {
				t.Errorf("GetProcessFlowByName() got is empty")
			} else {
				t.Log("FlowNode num: ", len(got))
			}
			if len(got1) == 0 {
				t.Errorf("GetProcessFlowByName() got1 is empty")
			} else {
				t.Log("FlowNodeRelation num: ", len(got1))
			}
		})
	}
}
