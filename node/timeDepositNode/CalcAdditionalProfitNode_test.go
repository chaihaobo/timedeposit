/*
 * @Author: Hugo
 * @Date: 2022-05-16 04:15:54
 * @Last Modified by: Hugo
 * @Last Modified time: 2022-05-23 09:48:32
 */
package timeDepositNode

import (
	"reflect"
	"testing"

	"gitlab.com/bns-engineering/td/common/constant"
	"gitlab.com/bns-engineering/td/node"
	"gitlab.com/bns-engineering/td/service/mambuEntity"
)

func TestCalcAdditionalProfitNode_RunProcess(t *testing.T) {
	type fields struct {
		Node node.Node
	}
	type args struct {
		tmpTDAccount mambuEntity.TDAccount
		flowID       string
		nodeName string
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
			node := &CalcAdditionalProfitNode{
				Node: tt.fields.Node,
			}
			got, err := node.RunProcess(tt.args.tmpTDAccount, tt.args.flowID, tt.args.nodeName)
			if (err != nil) != tt.wantErr {
				t.Errorf("CalcAdditionalProfitNode.RunProcess() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CalcAdditionalProfitNode.RunProcess() = %v, want %v", got, tt.want)
			}
		})
	}
}
