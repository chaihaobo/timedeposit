// Package node
// @author： Boice
// @createTime：2022/5/26 17:24
package node

import (
	"github.com/pkg/errors"
	"gitlab.com/bns-engineering/td/core/engine/mambu/accountservice"
	"go.uber.org/zap"
)

type UndoMaturityNode struct {
	*Node
}

func (node *UndoMaturityNode) Run() (INodeResult, error) {

	account, err := node.GetMambuAccount()
	if err != nil {
		return nil, err
	}
	if account.IsCaseA() {
		undoMaturityResult := accountservice.UndoMaturityDate(account.ID)
		if !undoMaturityResult {
			return nil, errors.New("undo Maturity Date Failed")
		}
	} else {
		zap.L().Info("not match! skip it")
	}
	return NodeResultSuccess, nil
}
