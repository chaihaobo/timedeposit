/*
 * @Author: Hugo
 * @Date: 2022-05-16 04:16:12
 * @Last Modified by: Hugo
 * @Last Modified time: 2022-05-23 11:12:13
 */
package timeDepositNode

import (
	"errors"
	"fmt"
	"go.uber.org/zap"

	"gitlab.com/bns-engineering/td/common/constant"
	"gitlab.com/bns-engineering/td/common/util"
	"gitlab.com/bns-engineering/td/node"
	"gitlab.com/bns-engineering/td/service/mambuEntity"
	mambuservices "gitlab.com/bns-engineering/td/service/mambuServices"
)

type UpdateAccNode struct {
	node.Node
}

func NewUpdateAccNode() *UpdateAccNode {
	tmpNode := new(UpdateAccNode)
	tmpNode.Name = constant.UpdateAccountNode
	tmpNode.Node.NodeRun = tmpNode
	return tmpNode
}

func (node *UpdateAccNode) RunProcess(tmpTDAccount mambuEntity.TDAccount, flowID string, nodeName string) (constant.FlowNodeStatus, error) {
	if tmpTDAccount.IsCaseB1_1() || tmpTDAccount.IsCaseB2() {
		newDate := util.GetDate(tmpTDAccount.MaturityDate)
		isApplySucceed := mambuservices.UpdateMaturifyDateForTDAccount(tmpTDAccount.ID, newDate)
		if !isApplySucceed {
			zap.L().Error(fmt.Sprintf("Apply profit failed for account: %v", tmpTDAccount.ID))
			return constant.FlowNodeFailed, errors.New("call mambu service failed")
		} else {
			zap.L().Info(fmt.Sprintf("Finish apply profit for account: %v", tmpTDAccount.ID))
			return constant.FlowNodeFinish, nil
		}
	} else {
		zap.L().Info(fmt.Sprintf("No need to update maturity info for td account, accNo: %v", tmpTDAccount.ID))
		return constant.FlowNodeSkip, nil
	}
}
