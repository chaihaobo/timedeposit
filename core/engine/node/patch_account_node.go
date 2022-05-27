// Package node
// @author： Boice
// @createTime：2022/5/26 18:08
package node

import (
	"errors"
	"fmt"
	"gitlab.com/bns-engineering/td/common/util"
	"gitlab.com/bns-engineering/td/core/engine/mambu/accountservice"
	"go.uber.org/zap"
)

type PatchAccountNode struct {
	*Node
}

func (node *PatchAccountNode) Run() (INodeResult, error) {
	account, err := node.GetMambuAccount(node.AccountId, false)
	if err != nil {
		return nil, err
	}
	if account.IsCaseB1_1() || account.IsCaseB2() {
		newDate := util.GetDate(account.MaturityDate)
		isApplySucceed := accountservice.UpdateMaturifyDateForTDAccount(account.ID, newDate)
		if !isApplySucceed {
			zap.L().Error(fmt.Sprintf("Apply profit failed for account: %v", account.ID))
			return nil, errors.New("call mambu service failed")
		}
		zap.L().Info(fmt.Sprintf("Finish apply profit for account: %v", account.ID))
	} else {
		zap.L().Info("not match! skip it")
	}
	return NodeResultSuccess, nil
}
