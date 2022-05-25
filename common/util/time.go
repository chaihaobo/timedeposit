/*
 * @Author: Hugo
 * @Date: 2022-05-11 10:49:31
 * @Last Modified by: Hugo
 * @Last Modified time: 2022-05-16 09:40:22
 */
package util

import (
	"github.com/uniplaces/carbon"
	"time"
)

func GetDaysBetweenTime(startTime, endTime time.Time) int {
	if startTime.Before(endTime) {
		return int(endTime.Sub(startTime).Hours() / 24)
	} else {
		return int(startTime.Sub(endTime).Hours() / 24)
	}
}

// Check whether 2 days in a same day
func InSameDay(t1, t2 time.Time) bool {
	return carbon.NewCarbon(t1).IsSameDay(carbon.NewCarbon(t2))
}

// Get the date string of current time
func GetDate(tmpTime time.Time) string {
	return time.Unix(tmpTime.Local().Unix(), 0).Format("2006-01-02") //打印结果：2017-04-11 13:30:39
}
