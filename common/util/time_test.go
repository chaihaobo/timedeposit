/*
 * @Author: Hugo
 * @Date: 2022-05-11 10:49:31
 * @Last Modified by: Hugo
 * @Last Modified time: 2022-05-11 11:00:09
 */
package util

import (
	"testing"
	"time"
)

func TestGetDaysBetweenTime(t *testing.T) {
	type args struct {
		startTime time.Time
		endTime   time.Time
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		// TODO: Add test cases.
		{
			name: "test 1",
			args: args{
				startTime: time.Date(2022, 5, 11, 11, 22, 33, 00, time.Local),
				endTime:   time.Date(2022, 5, 10, 10, 22, 33, 00, time.Local),
			},
			want: 1,
		},
		{
			name: "test 2",
			args: args{
				startTime: time.Date(2022, 5, 11, 11, 22, 33, 00, time.Local),
				endTime:   time.Date(2022, 5, 12, 12, 22, 33, 00, time.Local),
			},
			want: 1,
		},
		{
			name: "test 3",
			args: args{
				startTime: time.Date(2022, 5, 10, 11, 22, 33, 00, time.Local),
				endTime:   time.Date(2022, 5, 10, 12, 22, 33, 00, time.Local),
			},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetDaysBetweenTime(tt.args.startTime, tt.args.endTime); got != tt.want {
				t.Errorf("GetDaysBetweenTime() = %v, want %v", got, tt.want)
			}
		})
	}
}
