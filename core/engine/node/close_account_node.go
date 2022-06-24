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

type CloseAccountNode struct {
	*Node
}

func (node *CloseAccountNode) Run(ctx context.Context) (INodeResult, error) {
	account, err := node.GetMambuAccount(ctx, node.AccountId, false)
	if err != nil {
		return nil, err
	}

	totalBalance := account.Balances.TotalBalance
	if (account.IsCaseB3() && totalBalance > 0) ||
		(account.IsCaseC() && totalBalance > 0) {
		notes := fmt.Sprintf("AccountNo:%v, FlowID:%v", account.ID, node.FlowId)
		isApplySucceed := accountservice.CloseAccount(node.GetContext(ctx), account.ID, notes)
		if !isApplySucceed {
			err := errors.New("call Mambu service failed")
			log.Error(ctx, fmt.Sprintf("close account failed for account: %v", account.ID), err)
			return nil, err
		} else {
			log.Info(ctx, fmt.Sprintf("Finish close account for account: %v", account.ID))
		}
	} else {
		log.Info(ctx, "not match! skip it")
		return ResultSkip, nil
	}
	return ResultSuccess, nil
}
