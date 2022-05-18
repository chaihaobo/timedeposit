/*
 * @Author: Hugo
 * @Date: 2022-05-16 04:16:12
 * @Last Modified by: Hugo
 * @Last Modified time: 2022-05-18 05:22:13
 */
package timeDepositNode

import (
	"errors"

	"gitlab.com/bns-engineering/td/common/log"
	"gitlab.com/bns-engineering/td/common/util"
	"gitlab.com/bns-engineering/td/node"
	"gitlab.com/bns-engineering/td/service/mambuEntity"
	mambuservices "gitlab.com/bns-engineering/td/service/mambuServices"
)

type UpdateAccNode struct {
	node.Node
}

func (node *UpdateAccNode) Process() {
	CurNodeName := "update_account_node"
	tmpTDAccount, tmpFlowTask, nodeLog := node.GetAccAndFlowLog(CurNodeName)
	if !needToUpdateTDAccInfo(tmpTDAccount) {
		log.Log.Info("No need to update maturity info for td account, accNo: %v", tmpTDAccount.ID)
		//Todo: db log
		node.UpdateLogWhenSkipNode(tmpFlowTask, CurNodeName, nodeLog)
	} else {
		//Todo: update new date for /_rekening/rekeningTanggalJatohTempo
		newDate := util.GetDate(tmpTDAccount.Maturitydate)
		isApplySucceed := mambuservices.UpdateMaturifyDateForTDAccount(tmpTDAccount.ID, newDate)
		if !isApplySucceed {
			log.Log.Error("Apply profit failed for account: %v", tmpTDAccount.ID)
			//Todo: db log
			node.UpdateLogWhenNodeFailed(tmpFlowTask, nodeLog, errors.New("Call mambu service failed"))
			//Failed in this node, skip all the other steps and finish this order
			return
		} else {
			node.UpdateLogWhenNodeFinish(tmpFlowTask, nodeLog)
			log.Log.Info("Finish apply profit for account: %v", tmpTDAccount.ID)
		}
	}

	log.Log.Info("ProfitApplyNode: OutputData: %v", tmpTDAccount)
	node.Node.Output <- tmpTDAccount
}

func needToUpdateTDAccInfo(tmpTDAccount mambuEntity.TDAccount) bool {
	//Todo: not implemented
	panic("unimplemented")
}
