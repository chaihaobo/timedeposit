// Package holidayservice
// @author： Boice
// @createTime：2022/5/31 18:17
package holidayservice

import (
	"context"
	"gitlab.com/bns-engineering/td/common/config"
	"gitlab.com/bns-engineering/td/common/logger"
	"testing"
)

func init() {
	logger.SetUp(config.Setup("../../../../config.json"))
}
func TestGetHolidayList(t *testing.T) {

	GetHolidayList(context.Background())
}
