// Package node
// @author： Boice
// @createTime：2022/5/26 17:24
package node

import (
	"context"
	"github.com/pkg/errors"
	"gitlab.com/bns-engineering/td/common/log"
	"gitlab.com/bns-engineering/td/core/engine/mambu/accountservice"
)

type UndoMaturityNode struct {
	*Node
}

func (node *UndoMaturityNode) Run(ctx context.Context) (INodeResult, error) {

	account, err := node.GetMambuAccount(ctx, node.AccountId, true)
	if err != nil {
		return nil, err
	}
	if account.IsCaseA(node.TaskCreateTime) {
		undoMaturityResult := accountservice.UndoMaturityDate(node.GetContext(ctx), account.ID)
		if !undoMaturityResult {
			return nil, errors.New("undo Maturity Date Failed")
		}
	} else {
		log.Info(ctx, "not match! skip it")
		return ResultSkip, nil
	}
	return ResultSuccess, nil
}
