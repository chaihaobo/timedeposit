// Package node
// @author： Boice
// @createTime：2022/5/26 11:21
package node

import (
	"context"
	"gitlab.com/bns-engineering/td/common/log"
)

type EndNode struct {
	*Node
}

func (node *EndNode) Run(ctx context.Context) (INodeResult, error) {
	log.Info(ctx, "starting End node")

	return ResultSuccess, nil
}
