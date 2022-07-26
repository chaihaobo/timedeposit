// Package node
// @author： Boice
// @createTime：2022/5/26 17:24
package node

import (
	"context"
	"github.com/pkg/errors"
)

type UndoMaturityNode struct {
	node
}

func (node *UndoMaturityNode) Run(ctx context.Context) (INodeResult, error) {

	account, err := node.GetMambuAccount(ctx, node.accountId, true)
	if err != nil {
		return nil, err
	}
	if account.IsCaseA(node.taskCreateTime) {
		undoMaturityResult := node.service.Mambu.Account.UndoMaturityDate(node.GetContext(ctx), account.ID)
		if !undoMaturityResult {
			return nil, errors.New("undo Maturity Date Failed")
		}
	} else {
		node.common.Logger.Info(ctx, "not match! skip it")
		return ResultSkip, nil
	}
	return ResultSuccess, nil
}
