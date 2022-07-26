// Package node
// @author： Boice
// @createTime：2022/5/26 18:08
package node

import (
	"context"
	"errors"
	"fmt"
	"github.com/uniplaces/carbon"
)

type PatchAccountNode struct {
	node
}

func (node *PatchAccountNode) Run(ctx context.Context) (INodeResult, error) {
	account, err := node.GetMambuAccount(ctx, node.accountId, false)
	if err != nil {
		return nil, err
	}
	if account.IsCaseB1_1(node.taskCreateTime) || account.IsCaseB2(node.taskCreateTime) {
		newDate := carbon.NewCarbon(account.MaturityDate).DateString()
		isApplySucceed := node.service.Mambu.Account.UpdateMaturifyDateForTDAccount(node.GetContext(ctx), account.ID, newDate)
		if !isApplySucceed {
			err := errors.New("call mambu service failed")
			node.common.Logger.Error(ctx, fmt.Sprintf("Apply profit failed for account: %v", account.ID), err)
			return nil, err
		}
		node.common.Logger.Info(ctx, fmt.Sprintf("Finish apply profit for account: %v", account.ID))
	} else {
		node.common.Logger.Info(ctx, "not match! skip it")
		return ResultSkip, nil
	}
	return ResultSuccess, nil
}
