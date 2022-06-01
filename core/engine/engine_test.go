// Package engine
// @author： Boice
// @createTime：2022/5/26 13:59
package engine

import (
	"gitlab.com/bns-engineering/td/common/config"
	"testing"
)

func init() {
	logger.SetUp(config.Setup("../../config.yaml"))
}

func TestEngine(t *testing.T) {
	Start("11645631879")
}
