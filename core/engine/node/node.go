// Package node
// @author： Boice
// @createTime：2022/5/26 15:13
package node

type INode interface {
	Run() (INodeResult, error)
	SetUp(flowId string, accountId string)
}

type Node struct {
	FlowId    string
	AccountId string
}

func (node *Node) SetUp(flowId string, accountId string) {
	node.FlowId = flowId
	node.AccountId = accountId
}

type INodeResult interface {
	GetNodeResult() string
}

type NodeResult struct {
	result string
}

func (nodeResult *NodeResult) GetNodeResult() string {
	return nodeResult.result
}

func NewNodeResult(result string) *NodeResult {
	return &NodeResult{result}
}
