/*
 * @Author: Hugo
 * @Date: 2022-05-16 04:15:23
 * @Last Modified by: Hugo
 * @Last Modified time: 2022-05-23 11:12:19
 */
package timeDepositNode

import (
	"errors"
	"fmt"
	"go.uber.org/zap"

	"gitlab.com/bns-engineering/td/common/constant"
	"gitlab.com/bns-engineering/td/node"
	"gitlab.com/bns-engineering/td/service/mambuEntity"
	mambuservices "gitlab.com/bns-engineering/td/service/mambuServices"
)

type WithdrawBalanceNode struct {
	node.Node
	// nodeName string
}

func NewWithdrawBalanceNode() *WithdrawBalanceNode {
	tmpNode := new(WithdrawBalanceNode)
	// tmpNode.nodeName = "withdraw_balance_node"
	tmpNode.Node.NodeRun = tmpNode
	return tmpNode
}

func (node *WithdrawBalanceNode) Process() {
	node.RunNode("withdraw_balance_node")
}

func (node *WithdrawBalanceNode) RunProcess(tmpTDAccount mambuEntity.TDAccount, flowID string, nodeName string) (constant.FlowNodeStatus, error) {
	newTDAccount, err := mambuservices.GetTDAccountById(tmpTDAccount.ID)
	if err != nil {
		zap.L().Error(fmt.Sprintf("Failed to get info of td account: %v", tmpTDAccount.ID))
		errMsg := "Failed to get detail info of td account"
		zap.L().Error(errMsg)
		return constant.FlowNodeFailed, errors.New(errMsg)
	}

	totalBalance := newTDAccount.Balances.TotalBalance
	if !(newTDAccount.IsCaseB() && totalBalance > 0) {
		zap.L().Info(fmt.Sprintf("No need to withdraw balance, accNo: %v", flowID))
		return constant.FlowNodeSkip, nil
	} else {
		// Get benefit account info
		benefitAccount, err := mambuservices.GetTDAccountById(newTDAccount.OtherInformation.BhdNomorRekPencairan)
		if err != nil {
			zap.L().Error(fmt.Sprintf("Failed to get benefit acc info of td account: %v, benefit acc id:%v", newTDAccount.ID, newTDAccount.OtherInformation.BhdNomorRekPencairan))
			return constant.FlowNodeSkip, errors.New("call mambu get benefit acc info failed")
		}
		channelID := fmt.Sprintf("RAKTRAN_DEPMUDC_%vM", newTDAccount.OtherInformation.Tenor)
		withdrawTransID := flowID + "-" + nodeName + "-" + "Withdraw"
		withrawResp, err := mambuservices.WithdrawTransaction(newTDAccount, benefitAccount, totalBalance, withdrawTransID, channelID)
		if err != nil {
			zap.L().Error(fmt.Sprintf("Failed to withdraw for td account: %v", newTDAccount.ID))
			return constant.FlowNodeFailed, errors.New("call mambu withdraw failed")
		}
		zap.L().Info(fmt.Sprintf("Finish withdraw balance for accNo: %v, encodedKey:%v", tmpTDAccount.ID, withrawResp.EncodedKey))

		depositTransID := flowID + "-" + nodeName + "-" + "Deposit"
		depositResp, err := mambuservices.DepositTransaction(tmpTDAccount, benefitAccount, totalBalance, depositTransID, channelID)
		if err != nil {
			zap.L().Error(fmt.Sprintf("Failed to deposit for td account: %v", newTDAccount.ID))
			//todo: Add reverse withdraw here
			zap.L().Error(fmt.Sprintf("depositResp: %v", depositResp))

			return constant.FlowNodeFailed, errors.New("call mambu deposit failed")
		}
		zap.L().Info(fmt.Sprintf("Finish deposit balance for accNo: %v, encodedKey:%v", tmpTDAccount.ID, depositResp.EncodedKey))
		return constant.FlowNodeFinish, nil
	}
}
