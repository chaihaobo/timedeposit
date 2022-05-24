/*
 * @Author: Hugo
 * @Date: 2022-05-16 04:16:12
 * @Last Modified by: Hugo
 * @Last Modified time: 2022-05-23 09:46:36
 */
package timeDepositNode

import (
	"reflect"
	"testing"

	"gitlab.com/bns-engineering/td/common/constant"
	"gitlab.com/bns-engineering/td/node"
	"gitlab.com/bns-engineering/td/service/mambuEntity"
)

func TestUpdateAccNode_RunProcess(t *testing.T) {
	type fields struct {
		Node node.Node
	}
	type args struct {
		tmpTDAccount mambuEntity.TDAccount
		flowID       string
		nodeName     string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    constant.FlowNodeStatus
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node := &UpdateAccNode{
				Node: tt.fields.Node,
			}
			got, err := node.RunProcess(tt.args.tmpTDAccount, tt.args.flowID, tt.args.nodeName)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateAccNode.RunProcess() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UpdateAccNode.RunProcess() = %v, want %v", got, tt.want)
			}
		})
	}
}
