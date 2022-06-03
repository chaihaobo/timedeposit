//
// @author： Boice
// @createTime：
package logger

import (
	"gitlab.com/bns-engineering/td/common/config"
	"go.uber.org/zap"
	"testing"
)

func TestLogger(t *testing.T) {
	t.Run("test logger", func(t *testing.T) {
		config.Setup("../../config.json")
		SetUp(config.TDConf)
		zap.L().Info("ok", zap.String("123", "123"))

	})
}
