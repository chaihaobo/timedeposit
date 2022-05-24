/*
 * @Author: Hugo
 * @Date: 2022-05-16 04:15:07
 * @Last Modified by: Hugo
 * @Last Modified time: 2022-05-23 11:12:05
 */
package timeDepositNode

import (
	"errors"
	"fmt"
	"strconv"

	"gitlab.com/bns-engineering/td/common/constant"
	"gitlab.com/bns-engineering/td/common/log"
	"gitlab.com/bns-engineering/td/node"
	"gitlab.com/bns-engineering/td/service/mambuEntity"
	mambuservices "gitlab.com/bns-engineering/td/service/mambuServices"
)

type TransferProfitNode struct {
	node.Node
	// nodeName string
}

func NewTransferProfitNode() *TransferProfitNode {
	tmpNode := new(TransferProfitNode)
	// tmpNode.nodeName = "transfer_profit_node"
	tmpNode.Node.NodeRun = tmpNode
	return tmpNode
}

func (node *TransferProfitNode) Process() {
	node.RunNode("transfer_profit_node")
}

func (node *TransferProfitNode) RunProcess(tmpTDAccount mambuEntity.TDAccount, flowID string, nodeName string) (constant.FlowNodeStatus, error) {
	// Get the latest info of TD Account
	newTDAccount, err := mambuservices.GetTDAccountById(tmpTDAccount.ID)
	if err != nil {
		log.Log.Error("Failed to get info of td account: %v", tmpTDAccount.ID)
		errMsg := "Failed to get detail info of td account"
		log.Log.Error(errMsg)
		return constant.FlowNodeFailed, errors.New(errMsg)
	}

	// Get the principal amount of td account
	principal, err := strconv.ParseFloat(newTDAccount.Rekening.RekeningPrincipalAmount, 64)
	if err != nil {
		errMsg := fmt.Sprintf("Failed to convert Rekening.RekeningPrincipalAmount from string to float64, value:%v", newTDAccount.Rekening.RekeningPrincipalAmount)
		log.Log.Error(errMsg)
		return constant.FlowNodeFailed, errors.New(errMsg)
	}
	//Calculate the profit
	netProfit := newTDAccount.Balances.TotalBalance - principal

	//Need to transfer profit or not.
	if !newTDAccount.IsCaseB1_1() || netProfit <= 0 {
		log.Log.Info("No need to withdraw profit, accNo: %v", newTDAccount.ID)
		return constant.FlowNodeSkip, errors.New("No need to withdraw profit")
	}

	// Get benefit account info
	benefitAccount, err := mambuservices.GetTDAccountById(newTDAccount.OtherInformation.BhdNomorRekPencairan)
	if err != nil {
		log.Log.Error("Failed to get benefit acc info of td account: %v, benefit acc id:%v", newTDAccount.ID, newTDAccount.OtherInformation.BhdNomorRekPencairan)
		return constant.FlowNodeSkip, errors.New("call mambu get benefit acc info failed")
	}

	//Withdraw netProfit for deposit account
	channelID := fmt.Sprintf("RAKTRAN_DEPMUDC_%vM", newTDAccount.OtherInformation.Tenor)
	withdrawTransID := flowID + "-" + nodeName + "-" + "Withdraw"
	withrawResp, err := mambuservices.WithdrawTransaction(newTDAccount, benefitAccount, netProfit, withdrawTransID, channelID)
	if err != nil {
		log.Log.Error("Failed to withdraw for td account: %v", newTDAccount.ID)
		//Just return error, no need to reverse
		return constant.FlowNodeFailed, errors.New("call mambu withdraw failed")
	}
	log.Log.Info("Finish withdraw balance for accNo: %v, encodedKey:%v", newTDAccount.ID, withrawResp.EncodedKey)

	//Deposit netProfit to benefit account
	depositTransID := flowID + "-" + nodeName + "-" + "Deposit"
	depositResp, err := mambuservices.DepositTransaction(newTDAccount, benefitAccount, netProfit, depositTransID, channelID)
	if err != nil {
		log.Log.Error("Failed to deposit for td account: %v", newTDAccount.ID)
		//todo: Add reverse withdraw here

		return constant.FlowNodeFailed, errors.New("call mambu deposit failed")
	}
	log.Log.Info("Finish deposit balance for accNo: %v, encodedKey:%v", newTDAccount.ID, depositResp.EncodedKey)

	return constant.FlowNodeFinish, nil
}
