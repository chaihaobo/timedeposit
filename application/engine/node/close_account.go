// Package node
// @author： Boice
// @createTime：2022/5/26 18:08
package node

import (
	"context"
	"errors"
	"fmt"
)

type CloseAccountNode struct {
	node
}

func (node *CloseAccountNode) Run(ctx context.Context) (INodeResult, error) {
	account, err := node.GetMambuAccount(ctx, node.accountId, false)
	if err != nil {
		return nil, err
	}

	totalBalance := account.Balances.TotalBalance
	if (account.IsCaseB3(node.taskCreateTime) && totalBalance > 0) ||
		(account.IsCaseC() && totalBalance > 0) {
		notes := fmt.Sprintf("AccountNo:%v, FlowID:%v", account.ID, node.flowId)
		isApplySucceed := node.service.Mambu.Account.CloseAccount(node.GetContext(ctx), account.ID, notes)
		if !isApplySucceed {
			err := errors.New("call Mambu service failed")
			node.common.Logger.Error(ctx, fmt.Sprintf("close account failed for account: %v", account.ID), err)
			return nil, err
		} else {
			node.common.Logger.Info(ctx, fmt.Sprintf("Finish close account for account: %v", account.ID))
		}
	} else {
		node.common.Logger.Info(ctx, "not match! skip it")
		return ResultSkip, nil
	}
	return ResultSuccess, nil
}
