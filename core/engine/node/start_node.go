// Package node
// @author： Boice
// @createTime：2022/5/26 11:07
package node

type StartNode struct {
	*Node
}

func (node *StartNode) Run() (INodeResult, error) {

	return NodeResultSuccess, nil
}
