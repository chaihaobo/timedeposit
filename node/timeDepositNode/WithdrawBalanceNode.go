/*
 * @Author: Hugo
 * @Date: 2022-05-16 04:15:23
 * @Last Modified by: Hugo
 * @Last Modified time: 2022-05-18 05:07:04
 */
package timeDepositNode

import (
	"gitlab.com/hugo.hu/time-deposit-eod-engine/node"
)

type WithdrawBalanceNode struct {
	node.Node
}

func (node *WithdrawBalanceNode) Process() {
	// tmpTDAccount := <-node.Node.Input
	// latestTDAccount, err := mambuservices.GetTDAccountById(tmpTDAccount.ID)
	// if err != nil {
	// 	//Todo: log
	// 	log.Log.Error("Failed to get info of td account: %v", tmpTDAccount.ID)
	// 	//Finish current flow
	// 	return
	// }

	// //_otherInformation.bhdNomorRekPencairan
	// benefitAccount, err := mambuservices.GetTDAccountById(latestTDAccount.Otherinformation.Bhdnomorrekpencairan)
	// if err != nil {
	// 	//Todo: log
	// 	log.Log.Error("Failed to get benefit acc info of td account: %v, benefit acc id:%v", tmpTDAccount.ID, tmpTDAccount.Otherinformation.Bhdnomorrekpencairan)
	// 	//Finish current flow
	// }

	// if !needToWithdrawProfit(tmpTDAccount) {
	// 	log.Log.Info("No need to withdraw profit, accNo: %v", tmpTDAccount.ID)
	// } else {
	// 	principalAmount := latestTDAccount.Rekening.RekeningPrincipalAmount //Not sure for this
	// 	netProfit := latestTDAccount.Balances.Totalbalance - principalAmount
	// 	mambuservices.WithdrawNetProfit(latestTDAccount, benefitAccount, netProfit)
	// 	log.Log.Info("Finish withdraw profit for account: %v", tmpTDAccount.ID)
	// 	mambuservices.DepositNetprofit(latestTDAccount, benefitAccount, netProfit)
	// }

	// log.Log.Info("WithdrawBalanceNode: OutputData: %v", tmpTDAccount)
	// node.Node.Output <- tmpTDAccount
}
