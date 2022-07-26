// Package mambu
// @author： Boice
// @createTime：2022/7/26 12:07
package mambu

import (
	"context"
	"fmt"
	"gitlab.com/bns-engineering/common/tracer"
	"gitlab.com/bns-engineering/td/common"
	"gitlab.com/bns-engineering/td/common/constant"
	"gitlab.com/bns-engineering/td/model/mambu"
	"gitlab.com/bns-engineering/td/repository"
)

type HolidayService interface {
	GetHolidayList(ctx context.Context) mambu.HolidayInfo
}

type holidayService struct {
	common      *common.Common
	repository  *repository.Repository
	mambuClient Client
}

func newHolidayService(common *common.Common, repo *repository.Repository, mambuClient Client) HolidayService {
	return &holidayService{
		common:      common,
		repository:  repo,
		mambuClient: mambuClient,
	}
}

func (h *holidayService) GetHolidayList(ctx context.Context) mambu.HolidayInfo {
	tr := tracer.StartTrace(ctx, "holidayservice-GetHolidayList")
	ctx = tr.Context()
	defer tr.Finish()
	h.common.Logger.Info(ctx, fmt.Sprintf("getUrl: %v", constant.HolidayInfoUrl))
	var holidayInfo mambu.HolidayInfo
	err := h.mambuClient.Get(ctx, constant.UrlOf(constant.HolidayInfoUrl), &holidayInfo, h.mambuClient.DBPersistence(ctx, "GetHolidayList"))
	if err != nil {
		h.common.Logger.Error(ctx, fmt.Sprintf("Query holiday Info failed! query url: %v", constant.HolidayInfoUrl), err)
		return holidayInfo
	}
	return holidayInfo
}
