/*
 * @Author: Hugo
 * @Date: 2022-05-11 11:17:19
 * @Last Modified by: Hugo
 * @Last Modified time: 2022-05-17 09:13:19
 */
package util

import (
	"fmt"
	"testing"
)

func TestPostData(t *testing.T) {
	type args struct {
		postJsonStr string
		postUrl     string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "query TD Acc List",
			args: args{
				postJsonStr: "{\"sortingCriteria\":{\"field\":\"id\",\"order\":\"ASC\"},\"filterCriteria\":[{\"field\":\"accountType\",\"operator\":\"EQUALS\",\"value\":\"FIXED_DEPOSIT\"},{\"field\":\"accountState\",\"operator\":\"IN\",\"values\":[\"ACTIVE\",\"MATURED\"]},{\"field\":\"_rekening.rekeningTanggalJatohTempo\",\"operator\":\"BETWEEN\",\"value\":\"2022-05-01T00:00:00+07:00\",\"secondValue\":\"2022-05-09T23:59:59+07:00\"}]}",
				postUrl:     "https://cbs-dev1.aladinbank.id/api/deposits:search?offset=0&limit=100&detailsLevel=FULL",
			},
			want:    "",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _, err := HttpPostData(tt.args.postJsonStr, tt.args.postUrl)
			if (err != nil) != tt.wantErr {
				t.Errorf("PostData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			fmt.Println(got)
			if got != tt.want {
				t.Errorf("PostData() = %v, want %v", got, tt.want)
			}
		})
	}
}
