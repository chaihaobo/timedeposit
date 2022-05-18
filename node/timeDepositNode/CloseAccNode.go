/*
 * @Author: Hugo
 * @Date: 2022-05-16 04:16:26
 * @Last Modified by: Hugo
 * @Last Modified time: 2022-05-18 05:07:11
 */
package timeDepositNode

import (
	"gitlab.com/bns-engineering/td/common/log"
	"gitlab.com/bns-engineering/td/node"
)

//Close this TD Account
type CloseAccNode struct {
	node.Node
}

func (node *CloseAccNode) Process() {
	tmpTDAccount := <-node.Node.Input
	log.Log.Info("CloseAccNode: OutputData: %v", tmpTDAccount)
	//Todo:implement here
	node.Node.Output <- tmpTDAccount
}
