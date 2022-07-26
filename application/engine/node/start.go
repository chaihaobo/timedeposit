// Package node
// @author： Boice
// @createTime：2022/5/26 11:07
package node

import "context"

type StartNode struct {
	node
}

func (node *StartNode) Run(ctx context.Context) (INodeResult, error) {
	return ResultSuccess, nil
}
