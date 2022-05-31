// Package holidayservice
// @author： Boice
// @createTime：2022/5/31 17:32
package holidayservice

import (
	"context"
	"fmt"
	"gitlab.com/bns-engineering/td/common/constant"
	"gitlab.com/bns-engineering/td/common/util/mambu_http"
	"gitlab.com/bns-engineering/td/core/engine/mambu"
	"go.uber.org/zap"
	"time"
)

func GetHolidayList(ctx context.Context) []time.Time {
	var holidayList []time.Time
	zap.L().Info(fmt.Sprintf("getUrl: %v", constant.HolidayInfoUrl))
	var holidayInfo HolidayInfo
	err := mambu_http.Get(constant.UrlOf(constant.HolidayInfoUrl), &holidayInfo, mambu.SaveMambuRequestLog(ctx, "GetHolidayList"))
	if err != nil {
		zap.L().Error(fmt.Sprintf("Query holiday Info failed! query url: %v", constant.HolidayInfoUrl))
		return holidayList
	}
	zap.L().Info(fmt.Sprintf("Query td account Info result: %v", holidayList))

	for _, tmpHoliday := range holidayInfo.Holidays {
		tmpDate, err := time.Parse("2006-01-02", tmpHoliday.Date)
		if err != nil {
			zap.L().Error(fmt.Sprintf("Parse holiday error, src holiday info:%v", tmpHoliday.Date))
			continue
		}
		holidayList = append(holidayList, tmpDate)
	}
	return holidayList
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
