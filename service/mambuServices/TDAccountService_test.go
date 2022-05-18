/*
 * @Author: Hugo
 * @Date: 2022-05-11 12:21:05
 * @Last Modified by: Hugo
 * @Last Modified time: 2022-05-17 12:10:33
 */
package mambuservices

import (
	"encoding/json"
	"fmt"
	"testing"

	commonConfig "gitlab.com/hugo.hu/time-deposit-eod-engine/common/config"
	"gitlab.com/hugo.hu/time-deposit-eod-engine/common/log"
	mambuEntity "gitlab.com/hugo.hu/time-deposit-eod-engine/service/mambuEntity"
)

func TestGetTDAccountById(t *testing.T) {
	type args struct {
		tdAccountID string
	}
	tests := []struct {
		name string
		args args
		// want    mambuEntity.TDAccount
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "query td acc by id",
			args: args{tdAccountID: "11114361436"},
			// want:    mambuEntity.TDAccount{},
			wantErr: false,
		},
	}
	conf, _ := commonConfig.NewConfig("./../../config.json")
	log.InitLogConfig(conf)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetTDAccountById(tt.args.tdAccountID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetTDAccountById() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			fmt.Println("=====result struct================")
			log.Log.Info("%v", got)
			// if !reflect.DeepEqual(got, tt.want) {
			// 	t.Errorf("GetTDAccountById() = %v, want %v", got, tt.want)
			// }
		})
	}
}

func TestGetTDAccountListById(t *testing.T) {
	type args struct {
		searchParam mambuEntity.SearchParam
	}
	tests := []struct {
		name    string
		args    args
		want    []mambuEntity.TDAccount
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "query td acc by status and accountType",
			args: args{
				mambuEntity.SearchParam{
					Filtercriteria: []mambuEntity.Filtercriteria{
						{
							Field:    "accountState",
							Operator: "IN",
							Values:   []string{"ACTIVE", "MATURED"},
						},
						{
							Field:    "accountType",
							Operator: "EQUALS",
							Value:    "FIXED_DEPOSIT",
						},
					},
					Sortingcriteria: mambuEntity.Sortingcriteria{
						Field: "id",
						Order: "ASC",
					},
				},
			},
			want:    []mambuEntity.TDAccount{},
			wantErr: false,
		},
		{
			name: "query td acc by expire date",
			args: args{
				mambuEntity.SearchParam{
					Filtercriteria: []mambuEntity.Filtercriteria{
						{
							Field:    "accountState",
							Operator: "IN",
							Values:   []string{"ACTIVE", "MATURED"},
						},
						{
							Field:    "accountType",
							Operator: "EQUALS",
							Value:    "FIXED_DEPOSIT",
						},
						{
							Field:       "_rekening.rekeningTanggalJatohTempo",
							Operator:    "BETWEEN",
							Value:       "2022-05-01",
							Secondvalue: "2022-05-31",
						},
					},
					Sortingcriteria: mambuEntity.Sortingcriteria{
						Field: "id",
						Order: "ASC",
					},
				},
			},
			want:    []mambuEntity.TDAccount{},
			wantErr: false,
		},
	}
	conf, _ := commonConfig.NewConfig("./../../config.json")
	log.InitLogConfig(conf)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetTDAccountListByQueryParam(tt.args.searchParam)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetTDAccountListById() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			for index, tmpTDAcc := range got {
				b, err := json.Marshal(tmpTDAcc)
				if err != nil {
					fmt.Errorf("Json Convert Error! srcData:%v", tmpTDAcc)
				}
				log.Log.Info("QueryTDAccInfo: %v, %v", index, string(b))
			}
		})
	}
}
