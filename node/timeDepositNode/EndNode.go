/*
 * @Author: Hugo
 * @Date: 2022-04-29 10:24:39
 * @Last Modified by: Hugo
 * @Last Modified time: 2022-05-16 10:34:51
 */
package timeDepositNode

import (
	"gitlab.com/bns-engineering/td/common/log"
	"gitlab.com/bns-engineering/td/node"
)

//End Node
type EndNode struct {
	// Input <-chan string // input port
	node.Node
}

func (node *EndNode) Process() {
	tmpTDAccount := <-node.Node.Input
	log.Log.Info("EndNode: InputData: %v", tmpTDAccount)
}
