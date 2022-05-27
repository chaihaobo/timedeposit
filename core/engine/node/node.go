// Package node
// @author： Boice
// @createTime：2022/5/26 15:13
package node

import (
	"gitlab.com/bns-engineering/td/core/engine/mambu/accountservice"
	"gitlab.com/bns-engineering/td/service/mambuEntity"
)

const (
	NodeResultSuccess NodeResult = "success"
)

type INode interface {
	Run() (INodeResult, error)
	SetUp(flowId string, accountId string, nodeName string)
}

type Node struct {
	FlowId    string
	AccountId string
	NodeName  string
}

func (node *Node) SetUp(flowId string, accountId string, nodeName string) {
	node.FlowId = flowId
	node.AccountId = accountId
	node.NodeName = nodeName
}

func (node *Node) GetMambuAccount() (*mambuEntity.TDAccount, error) {
	id, err := accountservice.GetAccountById(node.AccountId)
	return id, err
}

type INodeResult interface {
	GetNodeResult() NodeResult
}

type NodeResult string

func (nodeResult NodeResult) GetNodeResult() NodeResult {
	return nodeResult
}
