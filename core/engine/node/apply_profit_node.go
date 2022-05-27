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

type ApplyProfitNode struct {
	*Node
}

func (node *ApplyProfitNode) Run() (INodeResult, error) {
	account, err := node.GetMambuAccount(node.AccountId)
	if err != nil {
		return nil, err
	}
	if account.IsCaseB() {
		isApplySucceed := accountservice.ApplyProfit(account.ID, node.FlowId)
		if !isApplySucceed {
			zap.L().Error(fmt.Sprintf("Apply profit failed for account: %v", account.ID))
			return nil, errors.New("call Mambu service failed")
		}
	} else {
		zap.L().Info("not match! skip it")
	}
	return NodeResultSuccess, nil
}
