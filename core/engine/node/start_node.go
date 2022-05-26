// Package node
// @author： Boice
// @createTime：2022/5/26 11:07
package node

import (
	"go.uber.org/zap"
)

type StartNode struct {
	*Node
}

func (node *StartNode) Run() (INodeResult, error) {
	zap.L().Info("starting start node")

	return NodeResultSuccess, nil
}
