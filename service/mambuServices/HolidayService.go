/*
 * @Author: Hugo
 * @Date: 2022-05-23 02:32:25
 * @Last Modified by: Hugo
 * @Last Modified time: 2022-05-23 02:45:05
 */
package mambuservices

import (
	"encoding/json"
	"time"

	"gitlab.com/bns-engineering/td/common/constant"
	"gitlab.com/bns-engineering/td/common/log"
	commonLog "gitlab.com/bns-engineering/td/common/log"
	"gitlab.com/bns-engineering/td/common/util"
)

func GetHolidayList() []time.Time {
	holidayList := []time.Time{}
	commonLog.Log.Info("getUrl: %v", constant.HolidayInfoUrl)
	resp, code, err := util.HttpGetData(constant.HolidayInfoUrl)
	if err != nil || code != constant.HttpStatusCodeSucceed {
		commonLog.Log.Error("Query holiday Info failed! query url: %v", constant.HolidayInfoUrl)
		return holidayList
	}
	commonLog.Log.Info("Query td account Info result: %v", resp)
	var holidayInfo HolidayInfo
	err = json.Unmarshal([]byte(resp), &holidayInfo)
	if err != nil {
		commonLog.Log.Error("Convert Json to HolidayInfo Failed. json: %v, err:%v", resp, err.Error())
		return holidayList
	}
	for _, tmpHoliday := range holidayInfo.Holidays {
		tmpDate, err := time.Parse("2006-01-02", tmpHoliday.Date)
		if err != nil {
			log.Log.Error("Parse holiday error, src holiday info:%v", tmpHoliday.Date)
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
