// Package node
// @author： Boice
// @createTime：2022/5/26 18:08
package node

import (
	"errors"
	"fmt"
	"gitlab.com/bns-engineering/td/core/engine/mambu/accountservice"
	"go.uber.org/zap"
)

type CloseAccountNode struct {
	*Node
}

func (node *CloseAccountNode) Run() (INodeResult, error) {
	account, err := node.GetMambuAccount()
	if err != nil {
		return nil, err
	}

	totalBalance := account.Balances.TotalBalance
	if (account.IsCaseB3() && totalBalance > 0) ||
		(account.IsCaseC() && totalBalance > 0) {
		notes := fmt.Sprintf("AccountNo:%v, FlowID:%v", account.ID, account)
		isApplySucceed := accountservice.CloseAccount(account.ID, notes)
		if !isApplySucceed {
			zap.L().Error(fmt.Sprintf("close account failed for account: %v", account.ID))
			return nil, errors.New("call Mambu service failed")
		} else {
			zap.L().Info(fmt.Sprintf("Finish close account for account: %v", account.ID))
		}
	} else {
		zap.L().Info("not match! skip it")
	}
	return NodeResultSuccess, nil
}
