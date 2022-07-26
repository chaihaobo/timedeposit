// Package node
// @author： Boice
// @createTime：2022/5/26 18:08
package node

import (
	"context"
	"errors"
	"fmt"
)

type ApplyProfitNode struct {
	node
}

func (node *ApplyProfitNode) Run(ctx context.Context) (INodeResult, error) {
	account, err := node.GetMambuAccount(ctx, node.accountId, false)
	if err != nil {
		return nil, err
	}
	if account.IsCaseB(node.taskCreateTime) {
		isApplySucceed := node.service.Mambu.Account.ApplyProfit(node.GetContext(ctx), account.ID, fmt.Sprintf("TDE-AUTO-%s", node.flowId), node.taskCreateTime)
		if !isApplySucceed {
			err := errors.New("call Mambu service failed")
			node.common.Logger.Error(ctx, fmt.Sprintf("Apply profit failed for account: %v", account.ID), err)
			return nil, err
		}
	} else {
		node.common.Logger.Info(ctx, "not match! skip it")
		return ResultSkip, nil
	}
	return ResultSuccess, nil
}
