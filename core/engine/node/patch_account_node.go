// Package node
// @author： Boice
// @createTime：2022/5/26 18:08
package node

import (
	"context"
	"errors"
	"fmt"
	"github.com/uniplaces/carbon"
	"gitlab.com/bns-engineering/td/common/log"
	"gitlab.com/bns-engineering/td/core/engine/mambu/accountservice"
)

type PatchAccountNode struct {
	*Node
}

func (node *PatchAccountNode) Run(ctx context.Context) (INodeResult, error) {
	account, err := node.GetMambuAccount(ctx, node.AccountId, false)
	if err != nil {
		return nil, err
	}
	if account.IsCaseB1_1() || account.IsCaseB2() {
		newDate := carbon.NewCarbon(account.MaturityDate).DateString()
		isApplySucceed := accountservice.UpdateMaturifyDateForTDAccount(node.GetContext(ctx), account.ID, newDate)
		if !isApplySucceed {
			err := errors.New("call mambu service failed")
			log.Error(ctx, fmt.Sprintf("Apply profit failed for account: %v", account.ID), err)
			return nil, err
		}
		log.Info(ctx, fmt.Sprintf("Finish apply profit for account: %v", account.ID))
	} else {
		log.Info(ctx, "not match! skip it")
		return ResultSkip, nil
	}
	return ResultSuccess, nil
}
