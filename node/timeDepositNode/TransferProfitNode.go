/*
 * @Author: Hugo
 * @Date: 2022-05-16 04:15:07
 * @Last Modified by: Hugo
 * @Last Modified time: 2022-05-18 05:06:34
 */
package timeDepositNode

import (
	"errors"
	"time"

	"gitlab.com/hugo.hu/time-deposit-eod-engine/common/log"
	"gitlab.com/hugo.hu/time-deposit-eod-engine/common/util"
	"gitlab.com/hugo.hu/time-deposit-eod-engine/node"
	"gitlab.com/hugo.hu/time-deposit-eod-engine/service/mambuEntity"
	mambuservices "gitlab.com/hugo.hu/time-deposit-eod-engine/service/mambuServices"
)

type TransferProfitNode struct {
	node.Node
}

func (node *TransferProfitNode) Process() {
	CurNodeName := "transfer_profit_node"
	tmpTDAccount, tmpFlowTask, nodeLog := node.GetAccAndFlowLog(CurNodeName)

	newTDAccount, err := mambuservices.GetTDAccountById(tmpTDAccount.ID)
	if err != nil {
		log.Log.Error("Failed to get info of td account: %v", tmpTDAccount.ID)
		node.UpdateLogWhenNodeFailed(tmpFlowTask, nodeLog, errors.New("call mambu get td acc info failed"))
		return
	}

	principalAmount := newTDAccount.Rekening.RekeningPrincipalAmount //Not sure for this
	netProfit := newTDAccount.Balances.Totalbalance - principalAmount

	//_otherInformation.bhdNomorRekPencairan
	benefitAccount, err := mambuservices.GetTDAccountById(newTDAccount.Otherinformation.Bhdnomorrekpencairan)
	if err != nil {
		log.Log.Error("Failed to get benefit acc info of td account: %v, benefit acc id:%v", tmpTDAccount.ID, tmpTDAccount.Otherinformation.Bhdnomorrekpencairan)

		node.UpdateLogWhenNodeFailed(tmpFlowTask, nodeLog, errors.New("call mambu get benefit acc info failed"))
	}

	if !needToTransferProfit(tmpTDAccount, netProfit) {
		log.Log.Info("No need to withdraw profit, accNo: %v", tmpTDAccount.ID)
	} else {

		mambuservices.WithdrawNetProfit(newTDAccount, benefitAccount, netProfit)
		log.Log.Info("Finish withdraw profit for account: %v", tmpTDAccount.ID)
		mambuservices.DepositNetprofit(newTDAccount, benefitAccount, netProfit)
	}

	log.Log.Info("WithdrawBalanceNode: OutputData: %v", tmpTDAccount)
	node.Node.Output <- tmpTDAccount
}

func needToTransferProfit(tdAccInfo mambuEntity.TDAccount, netProfit float64) bool {
	isARO := tdAccInfo.Otherinformation.Arononaro == "ARO"
	activeState := tdAccInfo.Accountstate == "ACTIVE"
	isStopARO := tdAccInfo.Otherinformation.s

	aroType := tdAccInfo.Otherinformation.

	rekeningTanggalJatohTempoDate, error := time.Parse("2006-01-02", tdAccInfo.Rekening.Rekeningtanggaljatohtempo)
	if error != nil {
		log.Log.Error("Error in parsing timeFormat for rekeningTanggalJatohTempoDate, accNo: %v, rekeningTanggalJatohTempo:%v", tdAccInfo.ID, tdAccInfo.Rekening.Rekeningtanggaljatohtempo)
		return false
	}

	tomorrow := time.Now().AddDate(0, 0, 1)
	return isARO &&
		activeState &&
		util.InSameDay(rekeningTanggalJatohTempoDate, tomorrow) &&
		util.InSameDay(rekeningTanggalJatohTempoDate, tdAccInfo.Maturitydate)
}
