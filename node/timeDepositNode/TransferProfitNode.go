/*
 * @Author: Hugo
 * @Date: 2022-05-16 04:15:07
 * @Last Modified by: Hugo
 * @Last Modified time: 2022-05-19 09:05:46
 */
package timeDepositNode

import (
	"errors"
	"fmt"

	"gitlab.com/bns-engineering/td/common/log"
	"gitlab.com/bns-engineering/td/node"
	mambuservices "gitlab.com/bns-engineering/td/service/mambuServices"
)

type TransferProfitNode struct {
	node.Node
}

func (node *TransferProfitNode) Process() {
	CurNodeName := "transfer_profit_node"
	tmpTDAccount, tmpFlowTask, nodeLog := node.GetAccAndFlowLog(CurNodeName)
	netProfit := tmpTDAccount.Balances.TotalBalance - tmpTDAccount.Rekening.RekeningPrincipalAmount

	newTDAccount, err := mambuservices.GetTDAccountById(tmpTDAccount.ID)
	if err != nil {
		log.Log.Error("Failed to get info of td account: %v", tmpTDAccount.ID)
		node.UpdateLogWhenNodeFailed(tmpFlowTask, nodeLog, errors.New("call mambu get td acc info failed"))
		return
	}
	tmpTDAccount = newTDAccount

	if !tmpTDAccount.IsCaseB1_1() {
		log.Log.Info("No need to withdraw profit, accNo: %v", tmpTDAccount.ID)
	} else {
		//_otherInformation.bhdNomorRekPencairan
		benefitAccount, err := mambuservices.GetTDAccountById(tmpTDAccount.OtherInformation.BhdNomorRekPencairan)
		if err != nil {
			log.Log.Error("Failed to get benefit acc info of td account: %v, benefit acc id:%v", tmpTDAccount.ID, tmpTDAccount.OtherInformation.BhdNomorRekPencairan)
			node.UpdateLogWhenNodeFailed(tmpFlowTask, nodeLog, errors.New("call mambu get benefit acc info failed"))
		}

		channelID := fmt.Sprintf("RAKTRAN_DEPMUDC_%vM", tmpTDAccount.OtherInformation.Tenor)
		withrawResp, err := mambuservices.WithdrawTransaction(tmpTDAccount, benefitAccount, nodeLog, netProfit, channelID)
		if err != nil {
			log.Log.Error("Failed to withdraw for td account: %v", tmpTDAccount.ID)
			node.UpdateLogWhenNodeFailed(tmpFlowTask, nodeLog, errors.New("call mambu withdraw failed"))
			//todo: Log failed transaction info here
			return
		}
		log.Log.Info("Finish withdraw balance for accNo: %v, encodedKey:%v", tmpTDAccount.ID, withrawResp.EncodedKey)
		depositResp, err := mambuservices.DepositTransaction(tmpTDAccount, benefitAccount, nodeLog, netProfit, channelID)
		if err != nil {
			log.Log.Error("Failed to deposit for td account: %v", tmpTDAccount.ID)
			node.UpdateLogWhenNodeFailed(tmpFlowTask, nodeLog, errors.New("call mambu deposit failed"))
			//todo: Add reverse withdraw here
			//todo: Log failed transaction info here
		}
		log.Log.Info("Finish deposit balance for accNo: %v, encodedKey:%v", tmpTDAccount.ID, depositResp.EncodedKey)
		node.Node.Output <- tmpTDAccount
	}

}
