/*
 * @Author: Hugo
 * @Date: 2022-05-18 03:38:56
 * @Last Modified by: Hugo
 * @Last Modified time: 2022-05-18 04:05:59
 */
package dao

import (
	"fmt"
	"testing"
)

func TestCreateFlowTask(t *testing.T) {
	type args struct {
		flowId    string
		accountId string
		flowName  string
	}
	tests := []struct {
		name string
		args args
		// want model.TFlowTaskInfo
	}{
		// TODO: Add test cases.
		{
			name: "test create",
			args: args{
				flowId:    "123456789",
				accountId: "11114361436",
				flowName:  "testFlowname",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CreateFlowTask(tt.args.flowId, tt.args.accountId, tt.args.flowName)
			fmt.Println(got)
		})
	}
}
