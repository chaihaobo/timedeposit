// Package engine

// @author： Boice

// @createTime：2022/5/26 13:59

package engine

import (
	"gitlab.com/bns-engineering/td/common/config"
	"gitlab.com/bns-engineering/td/common/logger"
	"testing"
)

func init() {
	logger.SetUp(config.Setup("../../config.json"))
}

func TestEngine(t *testing.T) {
	Start("11246851925")
}

func TestStart(t *testing.T) {
	type args struct {
		accountId string
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
		{
			name: "test  Case B1.1.1.1",
			args: args{
				accountId: "11169504404",
			},
		},
		{
			name: "test  Case B1.1.1.1 again",
			args: args{
				accountId: "11246851925",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Start(tt.args.accountId)
		})
	}
}

func TestRun(t *testing.T) {
	type args struct {
		flowId string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Test Retry failed flows",
			args: args{
				flowId: "20220606072530_11563057399",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Run(tt.args.flowId)
		})
	}
}
