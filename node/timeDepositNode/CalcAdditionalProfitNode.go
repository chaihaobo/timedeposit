/*
 * @Author: Hugo
 * @Date: 2022-05-16 04:15:54
 * @Last Modified by: Hugo
 * @Last Modified time: 2022-05-18 05:07:26
 */
package timeDepositNode

import (
	"gitlab.com/hugo.hu/time-deposit-eod-engine/common/log"
	"gitlab.com/hugo.hu/time-deposit-eod-engine/node"
)

//Calc the Additional Profit for TD Account
type CalcAdditionalProfitNode struct {
	node.Node
}

func (node *CalcAdditionalProfitNode) Process() {
	tmpTDAccount := <-node.Node.Input
	log.Log.Info("CalcAdditionalProfitNode: OutputData: %v", tmpTDAccount)
	//Todo:implement here
	node.Node.Output <- tmpTDAccount
}
