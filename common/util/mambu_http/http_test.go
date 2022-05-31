// Package http
// @author： Boice
// @createTime：2022/5/30 10:58
package mambu_http

import (
	"fmt"
	"github.com/pkg/errors"
	"gitlab.com/bns-engineering/td/common/config"
	"gitlab.com/bns-engineering/td/common/constant"
	"gitlab.com/bns-engineering/td/service/mambuEntity"
	"testing"
)

func init() {
	logger.SetUp(config.Setup("../../../config.yaml"))
}

func TestGet(t *testing.T) {
	var tdAccount = new(mambuEntity.TDAccount)
	getUrl := fmt.Sprintf(constant.GetTDAccountUrl, "11249460359")
	_, err := Get(getUrl, tdAccount)
	if err != nil {
		t.Errorf("test get error %v", errors.WithStack(err))
	}
}
