/*
 * @Author: Hugo
 * @Date: 2022-05-16 04:14:01
 * @Last Modified by: Hugo
 * @Last Modified time: 2022-05-18 05:06:05
 */
package timeDepositNode

import (
	"errors"
	"time"

	"gitlab.com/bns-engineering/td/common/log"
	"gitlab.com/bns-engineering/td/common/util"
	"gitlab.com/bns-engineering/td/node"
	mambuEntity "gitlab.com/bns-engineering/td/service/mambuEntity"
	mambuservices "gitlab.com/bns-engineering/td/service/mambuServices"
)

//Calc the Additional Profit for TD Account
type ProfitApplyNode struct {
	node.Node
}

func (node *ProfitApplyNode) Process() {
	CurNodeName := "profit_apply_node"
	tmpTDAccount, tmpFlowTask, nodeLog := node.GetAccAndFlowLog(CurNodeName)
	if !needToApplyProfit(tmpTDAccount) {
		node.UpdateLogWhenSkipNode(tmpFlowTask, CurNodeName, nodeLog)
		log.Log.Info("No need to apply profit, accNo: %v", tmpFlowTask.FlowId)
	} else {
		isApplySucceed := mambuservices.ApplyProfit(tmpTDAccount.ID, tmpFlowTask.FlowId)
		if !isApplySucceed {
			log.Log.Error("Apply profit failed for account: %v", tmpTDAccount.ID)
			node.UpdateLogWhenNodeFailed(tmpFlowTask, nodeLog, errors.New("call Mambu service failed"))
		} else {
			log.Log.Info("Finish apply profit for account: %v", tmpTDAccount.ID)
			node.UpdateLogWhenNodeFinish(tmpFlowTask, nodeLog)
		}
	}
	node.Node.Output <- tmpTDAccount
}

// Checking whether this account need to change Maturity date
func needToApplyProfit(tdAccInfo mambuEntity.TDAccount) bool {
	isARO := tdAccInfo.Otherinformation.AroNonAro == "ARO"
	activeState := tdAccInfo.Accountstate == "ACTIVE"
	rekeningTanggalJatohTempoDate, error := time.Parse("2006-01-02", tdAccInfo.Rekening.Rekeningtanggaljatohtempo)
	if error != nil {
		log.Log.Error("Error in parsing timeFormat for rekeningTanggalJatohTempoDate, accNo: %v, rekeningTanggalJatohTempo:%v", tdAccInfo.ID, tdAccInfo.Rekening.Rekeningtanggaljatohtempo)
		return false
	}

	//Note: We should check whether maturityDate is null
	return isARO &&
		activeState &&
		util.InSameDay(rekeningTanggalJatohTempoDate, time.Now()) &&
		rekeningTanggalJatohTempoDate.Before(tdAccInfo.Maturitydate)
}
