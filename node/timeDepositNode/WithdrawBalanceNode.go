/*
 * @Author: Hugo
 * @Date: 2022-05-16 04:15:23
 * @Last Modified by: Hugo
 * @Last Modified time: 2022-05-19 09:00:00
 */
package timeDepositNode

import (
	"errors"
	"fmt"

	"gitlab.com/bns-engineering/td/common/log"
	"gitlab.com/bns-engineering/td/node"
	mambuservices "gitlab.com/bns-engineering/td/service/mambuServices"
)

type WithdrawBalanceNode struct {
	node.Node
}

func (node *WithdrawBalanceNode) Process() {
	CurNodeName := "withdraw_balance_node"
	tmpTDAccount, tmpFlowTask, nodeLog := node.GetAccAndFlowLog(CurNodeName)
	totalBalance := tmpTDAccount.Balances.TotalBalance
	if !(tmpTDAccount.IsCaseB() && totalBalance > 0) {
		node.UpdateLogWhenSkipNode(tmpFlowTask, CurNodeName, nodeLog)
		log.Log.Info("No need to withdraw balance, accNo: %v", tmpFlowTask.FlowId)
	} else {

		//_otherInformation.bhdNomorRekPencairan
		benefitAccount, err := mambuservices.GetTDAccountById(tmpTDAccount.OtherInformation.BhdNomorRekPencairan)
		if err != nil {
			log.Log.Error("Failed to get benefit acc info of td account: %v, benefit acc id:%v", tmpTDAccount.ID, tmpTDAccount.OtherInformation.BhdNomorRekPencairan)
			node.UpdateLogWhenNodeFailed(tmpFlowTask, nodeLog, errors.New("call mambu get benefit acc info failed"))
		}

		channelID := fmt.Sprintf("RAKTRAN_DEPMUDC_%vM", tmpTDAccount.OtherInformation.Tenor)
		withrawResp, err := mambuservices.WithdrawTransaction(tmpTDAccount, benefitAccount, nodeLog, totalBalance, channelID)
		if err != nil {
			log.Log.Error("Failed to withdraw for td account: %v", tmpTDAccount.ID)
			node.UpdateLogWhenNodeFailed(tmpFlowTask, nodeLog, errors.New("call mambu withdraw failed"))
			//todo: Log failed transaction info here
			return
		}
		log.Log.Info("Finish withdraw balance for accNo: %v, encodedKey:%v", tmpTDAccount.ID, withrawResp.EncodedKey)
		depositResp, err := mambuservices.DepositTransaction(tmpTDAccount, benefitAccount, nodeLog, totalBalance, channelID)
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
