/*
 * @Author: Hugo
 * @Date: 2022-05-16 09:08:29
 * @Last Modified by: Hugo
 * @Last Modified time: 2022-05-24 01:17:37
 */
package api

import (
	"fmt"
	"go.uber.org/zap"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/trustmaster/goflow"
	"gitlab.com/bns-engineering/td/common/constant"
	"gitlab.com/bns-engineering/td/common/util"
	"gitlab.com/bns-engineering/td/dao"
	"gitlab.com/bns-engineering/td/flow"
	"gitlab.com/bns-engineering/td/node"
	"gitlab.com/bns-engineering/td/service/mambuEntity"
	mambuservices "gitlab.com/bns-engineering/td/service/mambuServices"
)

func StartTDFlow(c *gin.Context) {
	//Get all td accounts which need to process
	tmpQueryParam := generateSearchTDAccountParam()
	tmpTDAccountList, err := mambuservices.GetTDAccountListByQueryParam(tmpQueryParam)
	if err != nil {
		zap.L().Error(fmt.Sprintf("Query mambu service for TD Account List failed! error: %v", err))
		return
	}
	if len(tmpTDAccountList) == 0 {
		zap.L().Info("Query mambu service for TD Account List get empty! No TD Account need to process")
		return
	}

	for _, tmpTDAcc := range tmpTDAccountList {
		zap.L().Info(fmt.Sprintf("Before Run Flow for Account: %v", tmpTDAcc.ID))
	}
	zap.L().Info("=======================================")

	for _, tmpTDAcc := range tmpTDAccountList {
		RunFlow(&tmpTDAcc)
	}

}

//RunFlow run flow by account
func RunFlow(tmpTDAcc *mambuEntity.TDAccount) {
	tmpFlow := flow.GetProcessFlow("time_deposit_flow")
	inputNodeDataChan := make(chan node.NodeData)
	tmpFlow.SetInPort("In", inputNodeDataChan)

	// Run the net
	wait := goflow.Run(tmpFlow)

	// Now we can send some names and see what happens
	zap.L().Info(fmt.Sprintf("Start Run Flow for Account: %v", tmpTDAcc.ID))
	flowID := fmt.Sprintf("%v_%v", time.Now().Format("20060102150405"), tmpTDAcc.ID)
	flowTaskInfo := dao.CreateFlowTask(flowID, tmpTDAcc.ID, "time_deposit_flow")
	newTDAccount, err := mambuservices.GetTDAccountById(tmpTDAcc.ID)
	if err != nil {
		zap.L().Error(fmt.Sprintf("Failed to get info of td account: %v", tmpTDAcc.ID))
		flowTaskInfo.EndStatus = constant.FlowFailed
		flowTaskInfo.CurStatus = string(constant.FlowNodeFailed)
		flowTaskInfo.FlowStatus = constant.FlowFailed
		dao.UpdateFlowTask(flowTaskInfo)
		return
	}

	inputNodeDataChan <- node.NodeData{
		FlowTaskInfo:  flowTaskInfo,
		TDAccountInfo: newTDAccount,
	}

	// Send end of input
	close(inputNodeDataChan)
	// Wait until the net has completed its job
	result := <-wait
	zap.L().Info(fmt.Sprintf("Flow End result: %v", result))

	zap.L().Info(fmt.Sprintf("Flow Run Finishd for Account: %v", tmpTDAcc.ID))
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
