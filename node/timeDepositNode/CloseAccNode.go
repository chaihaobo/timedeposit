/*
 * @Author: Hugo
 * @Date: 2022-05-16 04:16:26
 * @Last Modified by: Hugo
 * @Last Modified time: 2022-05-23 11:11:03
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

//Close this TD Account
type CloseAccNode struct {
	node.Node
}

func NewCloseAccNode() *CloseAccNode {
	tmpNode := new(CloseAccNode)
	tmpNode.Name = constant.CloseAccountNode
	tmpNode.Node.NodeRun = tmpNode
	return tmpNode
}

func (node *CloseAccNode) RunProcess(tmpTDAccount mambuEntity.TDAccount, flowID string, nodeName string) (constant.FlowNodeStatus, error) {
	totalBalance := tmpTDAccount.Balances.TotalBalance
	if (tmpTDAccount.IsCaseB3() && totalBalance > 0) ||
		(tmpTDAccount.IsCaseC() && totalBalance > 0) {
		notes := fmt.Sprintf("AccountNo:%v, FlowID:%v", tmpTDAccount.ID, flowID)
		isApplySucceed := mambuservices.CloseAccount(tmpTDAccount.ID, notes)
		if !isApplySucceed {
			zap.L().Error(fmt.Sprintf("close account failed for account: %v", tmpTDAccount.ID))
			return constant.FlowNodeFailed, errors.New("call Mambu service failed")
		} else {
			zap.L().Info(fmt.Sprintf("Finish close account for account: %v", tmpTDAccount.ID))
			return constant.FlowNodeFinish, nil
		}
	} else {
		zap.L().Info(fmt.Sprintf("No need to close account, accNo: %v", flowID))
		return constant.FlowNodeSkip, nil
	}
}
