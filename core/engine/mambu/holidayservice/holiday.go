// Package holidayservice
// @author： Boice
// @createTime：2022/5/31 17:32
package holidayservice

import (
	"context"
	"fmt"
	"gitlab.com/bns-engineering/common/tracer"
	"gitlab.com/bns-engineering/td/common/constant"
	"gitlab.com/bns-engineering/td/common/log"
	"gitlab.com/bns-engineering/td/common/util/mambu_http"
	"gitlab.com/bns-engineering/td/common/util/mambu_http/persistence"
	"time"
)

func GetHolidayList(ctx context.Context) HolidayInfo {
	tr := tracer.StartTrace(ctx, "holidayservice-GetHolidayList")
	ctx = tr.Context()
	defer tr.Finish()
	log.Info(ctx, fmt.Sprintf("getUrl: %v", constant.HolidayInfoUrl))
	var holidayInfo HolidayInfo
	err := mambu_http.Get(ctx, constant.UrlOf(constant.HolidayInfoUrl), &holidayInfo, persistence.DBPersistence(ctx, "GetHolidayList"))
	if err != nil {
		log.Error(ctx, fmt.Sprintf("Query holiday Info failed! query url: %v", constant.HolidayInfoUrl), err)
		return holidayInfo
	}
	return holidayInfo
}

type HolidayInfo struct {
	NonWorkingDays []string        `json:"nonWorkingDays"`
	Holidays       []ModelHolidays `json:"holidays"`
}

type ModelHolidays struct {
	EncodedKey          string    `json:"encodedKey"`
	Name                string    `json:"name"`
	ID                  int       `json:"id"`
	Date                string    `json:"date"`
	IsAnnuallyRecurring bool      `json:"isAnnuallyRecurring"`
	CreationDate        time.Time `json:"creationDate"`
}
