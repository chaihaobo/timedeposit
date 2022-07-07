// Package node
// @author： Boice
// @createTime：2022/5/26 18:08
package node

import (
	"context"
	"errors"
	"fmt"
	"gitlab.com/bns-engineering/td/common/log"
	"gitlab.com/bns-engineering/td/core/engine/mambu/accountservice"
)

type ApplyProfitNode struct {
	*Node
}

func (node *ApplyProfitNode) Run(ctx context.Context) (INodeResult, error) {
	account, err := node.GetMambuAccount(ctx, node.AccountId, false)
	if err != nil {
		return nil, err
	}
	if account.IsCaseB(node.TaskCreateTime) {
		isApplySucceed := accountservice.ApplyProfit(node.GetContext(ctx), account.ID, fmt.Sprintf("TDE-AUTO-%s", node.FlowId))
		if !isApplySucceed {
			err := errors.New("call Mambu service failed")
			log.Error(ctx, fmt.Sprintf("Apply profit failed for account: %v", account.ID), err)
			return nil, err
		}
	} else {
		log.Info(ctx, "not match! skip it")
		return ResultSkip, nil

	}
	return ResultSuccess, nil
}
