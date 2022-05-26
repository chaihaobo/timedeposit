// Package node
// @author： Boice
// @createTime：2022/5/26 11:21
package node

import (
	"go.uber.org/zap"
)

type EndNode struct {
	*Node
}

func (node *EndNode) Run() (INodeResult, error) {
	zap.L().Info("starting End node")

	return NodeResultSuccess, nil
}
