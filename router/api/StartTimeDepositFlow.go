/*
 * @Author: Hugo
 * @Date: 2022-05-16 09:08:29
 * @Last Modified by: Hugo
 * @Last Modified time: 2022-05-19 11:01:26
 */
package api

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/trustmaster/goflow"
	"gitlab.com/bns-engineering/td/common/log"
	"gitlab.com/bns-engineering/td/common/util"
	"gitlab.com/bns-engineering/td/dao"
	"gitlab.com/bns-engineering/td/flow"
	"gitlab.com/bns-engineering/td/service/mambuEntity"
	mambuservices "gitlab.com/bns-engineering/td/service/mambuServices"
)

func StartTDFlow(c *gin.Context) {
	//Get all td accounts which need to process
	tmpQueryParam := generateSearchTDAccountParam()
	tmpTDAccountList, err := mambuservices.GetTDAccountListByQueryParam(tmpQueryParam)
	if err != nil {
		log.Log.Error("Query mambu service for TD Account List failed! error: %v", err)
		return
	}
	if len(tmpTDAccountList) == 0 {
		log.Log.Info("Query mambu service for TD Account List get empty! No TD Account need to process")
		return
	}

	for _, tmpTDAcc := range tmpTDAccountList {
		tmpFlow := flow.GetProcessFlow("time_deposit_flow")
		in := make(chan mambuEntity.TDAccount)
		tmpFlow.SetInPort("In", in)
		// Run the net
		wait := goflow.Run(tmpFlow)
		// Now we can send some names and see what happens
		log.Log.Info("Start Run Flow for Account: %v", tmpTDAcc.ID)
		flowID := fmt.Sprintf("%v_%v", time.Now().Format("20060102150405"), tmpTDAcc.ID)
		flowTaskInfo := dao.CreateFlowTask(flowID, tmpTDAcc.ID, "time_deposit_flow")
		tmpFlow.SetInPort("FlowTaskInfo", flowTaskInfo)
		in <- tmpTDAcc

		// Send end of input
		close(in)
		// Wait until the net has completed its job
		result := <-wait
		log.Log.Info("Flow End for Account: %v, result: %v", tmpTDAcc.ID, result)
	}

}

func generateSearchTDAccountParam() mambuEntity.SearchParam {
	tmpQueryParam := mambuEntity.SearchParam{
		FilterCriteria: []mambuEntity.FilterCriteria{
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
				Field: "_rekening.rekeningTanggalJatohTempo",
				//todo: Remember to set the value to today!
				Operator:    "BETWEEN",
				Value:       util.GetDate(time.Now().AddDate(0, 0, -20)), //today
				SecondValue: util.GetDate(time.Now().AddDate(0, 0, 1)),   //tomorrow
			},
		},
		SortingCriteria: mambuEntity.SortingCriteria{
			Field: "id",
			Order: "ASC",
		},
	}
	return tmpQueryParam
}
