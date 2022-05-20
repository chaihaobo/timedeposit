/*
 * @Author: Hugo
 * @Date: 2022-05-16 04:16:12
 * @Last Modified by: Hugo
 * @Last Modified time: 2022-05-18 05:22:13
 */
package timeDepositNode

import (
	"errors"
	"strings"
	"time"

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
	if !tmpTDAccount.IsCaseB1_1() && !tmpTDAccount.IsCaseB2() {
		log.Log.Info("No need to update maturity info for td account, accNo: %v", tmpTDAccount.ID)
		node.UpdateLogWhenSkipNode(tmpFlowTask, CurNodeName, nodeLog)
	} else {
		newDate := util.GetDate(tmpTDAccount.Maturitydate)
		isApplySucceed := mambuservices.UpdateMaturifyDateForTDAccount(tmpTDAccount.ID, newDate)
		if !isApplySucceed {
			log.Log.Error("Apply profit failed for account: %v", tmpTDAccount.ID)
			node.UpdateLogWhenNodeFailed(tmpFlowTask, nodeLog, errors.New("call mambu service failed"))
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

func needToUpdateTDAccInfo(tdAccInfo mambuEntity.TDAccount) bool {
	isARO := strings.ToUpper(tdAccInfo.Otherinformation.AroNonAro) == "ARO"
	activeState := tdAccInfo.Accountstate == "ACTIVE"
	rekeningTanggalJatohTempoDate, error := time.Parse("2006-01-02", tdAccInfo.Rekening.Rekeningtanggaljatohtempo)
	if error != nil {
		log.Log.Error("Error in parsing timeFormat for rekeningTanggalJatohTempoDate, accNo: %v, rekeningTanggalJatohTempo:%v", tdAccInfo.ID, tdAccInfo.Rekening.Rekeningtanggaljatohtempo)
		return false
	}

	isStopARO := tdAccInfo.Otherinformation.StopAro != "FALSE" 
	aroType := tdAccInfo.Otherinformation.AroType
	
	netProfit := tdAccInfo.Balances.Totalbalance - tdAccInfo.Rekening.RekeningPrincipalAmount

	return isARO &&						//B
		activeState &&//B
		util.InSameDay(rekeningTanggalJatohTempoDate, time.Now()) && //B
		rekeningTanggalJatohTempoDate.Before(tdAccInfo.Maturitydate) && //B
		((!isStopARO && //B1
			aroType == "Principal Only" && //B1
			netProfit > 0 ) || //B1.1
		(!isStopARO &&//B2
			aroType == "Full")) //B2
}


