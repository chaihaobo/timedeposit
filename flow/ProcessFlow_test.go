/*
 * @Author: Hugo
 * @Date: 2022-05-05 08:59:52
 * @Last Modified by: Hugo
 * @Last Modified time: 2022-05-23 11:39:00
 */
package flow

import (
	"fmt"
	logger "gitlab.com/bns-engineering/td/common/log"
	"go.uber.org/zap"
	"testing"

	"github.com/trustmaster/goflow"
	commonConfig "gitlab.com/bns-engineering/td/common/config"
)

func TestInitProcessFlow(t *testing.T) {
	conf := commonConfig.Setup("./../config.json")
	logger.SetUp(conf)
	InitWorkflow()
	for key, _ := range typeRegistry {
		zap.L().Info(fmt.Sprintf("%v", key))
		// zap.L().Info("%v", value)
	}
}

func TestNewProcessFlow(t *testing.T) {
	conf := commonConfig.Setup("./../config.json")
	logger.SetUp(conf)

	type args struct {
		flowName string
	}
	tests := []struct {
		name     string
		args     args
		inputStr string
	}{
		// TODO: Add test cases.
		{
			name:     "time deposit flow 0",
			args:     args{flowName: "time_deposit_flow"},
			inputStr: "111111",
		},
		{
			name:     "time deposit flow 1",
			args:     args{flowName: "time_deposit_flow"},
			inputStr: "111111",
		},
		{
			name:     "time deposit flow 2",
			args:     args{flowName: "time_deposit_flow"},
			inputStr: "22222",
		},
		{
			name:     "time deposit flow 3",
			args:     args{flowName: "time_deposit_flow"},
			inputStr: "33333",
		},
		{
			name:     "time deposit flow 4",
			args:     args{flowName: "time_deposit_flow"},
			inputStr: "44444",
		},
		{
			name:     "time deposit flow 5",
			args:     args{flowName: "time_deposit_flow"},
			inputStr: "55555",
		},
	}
	InitWorkflow()
	var got *goflow.Graph
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got == nil {
				got = NewProcessFlow(tt.args.flowName)
			}
			if got == nil {
				t.Errorf("NewProcessFlow() did not generate a flow")
			}
			// We need a channel to talk to it
			in := make(chan string)
			got.SetInPort("In", in)
			// Run the net
			wait := goflow.Run(got)
			// Now we can send some names and see what happens
			in <- tt.inputStr
			// Send end of input
			fmt.Println("Before Close Channel In")
			close(in)
			fmt.Println("End of Close Channel In")
			// Wait until the net has completed its job
			result := <-wait
			fmt.Println(result)
			fmt.Println("End of Working flow")
		})
	}
}
