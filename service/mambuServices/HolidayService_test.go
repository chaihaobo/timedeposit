/*
 * @Author: Hugo
 * @Date: 2022-05-23 02:32:25
 * @Last Modified by: Hugo
 * @Last Modified time: 2022-05-23 02:50:09
 */
package mambuservices

import (
	"testing"

	commonConfig "gitlab.com/bns-engineering/td/common/config"
)

func TestGetHolidayList(t *testing.T) {
	tests := []struct {
		name string
		// want []time.Time
	}{
		// TODO: Add test cases.
		{name: "Test get holiday info"},
	}
	conf, _ := commonConfig.NewConfig("./../../config.json")
	log.InitLogConfig(conf)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetHolidayList(); len(got) > 0 {
				for _, tmpHoliday := range got {
					zap.L().Info("holiday:%v", tmpHoliday.String())
				}
			}
		})
	}
}
