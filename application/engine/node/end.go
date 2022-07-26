// Package node
// @author： Boice
// @createTime：2022/5/26 11:21
package node

import (
	"context"
)

type EndNode struct {
	node
}

func (node *EndNode) Run(ctx context.Context) (INodeResult, error) {
	node.common.Logger.Info(ctx, "starting End node")
	return ResultSuccess, nil
}
