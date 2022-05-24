/*
 * @Author: Hugo
 * @Date: 2022-04-29 10:24:43
 * @Last Modified by: Hugo
 * @Last Modified time: 2022-05-23 09:47:57
 */

package timeDepositNode

import (
	"reflect"
	"testing"
	"time"

	"gitlab.com/bns-engineering/td/common/constant"
	"gitlab.com/bns-engineering/td/service/mambuEntity"
	mambuservices "gitlab.com/bns-engineering/td/service/mambuServices"
)

func TestMaturityDateNode_RunProcess(t *testing.T) {

	testAccountID := "11114361436"
	testTDAccount, err := mambuservices.GetTDAccountById(testAccountID)
	if err != nil {
		t.Errorf("get test td account error! accountID: %v, err: %v", testAccountID, err.Error())
	}

	type args struct {
		tmpTDAccount mambuEntity.TDAccount
		flowID       string
		nodeName     string
	}
	tests := []struct {
		name    string
		args    args
		want    constant.FlowNodeStatus
		wantErr bool
	}{
		{
			name: "test Update Maturity succeed",
			args: args{
				tmpTDAccount: testTDAccount,
				flowID:       "flowID in testCase: " + time.Now().Format("20060102150405"),
			},
			want:    constant.FlowNodeSkip,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node := new(MaturityDateNode)
			got, err := node.RunProcess(tt.args.tmpTDAccount, tt.args.flowID, tt.args.nodeName)
			if (err != nil) != tt.wantErr {
				t.Errorf("MaturityDateNode.RunProcess() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MaturityDateNode.RunProcess() = %v, want %v", got, tt.want)
			}
		})
	}
}
