/*
 * @Author: Hugo
 * @Date: 2022-04-29 10:24:39
 * @Last Modified by: Hugo
 * @Last Modified time: 2022-05-16 10:34:51
 */
package timeDepositNode

import (
	"gitlab.com/hugo.hu/time-deposit-eod-engine/common/log"
	"gitlab.com/hugo.hu/time-deposit-eod-engine/node"
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
