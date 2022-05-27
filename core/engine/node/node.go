// Package node
// @author： Boice
// @createTime：2022/5/26 15:13
package node

import (
	"encoding/json"
	"gitlab.com/bns-engineering/td/core/engine/mambu/accountservice"
	"gitlab.com/bns-engineering/td/repository"
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

func (node *Node) GetMambuAccount(accountId string, realTime ...bool) (*mambuEntity.TDAccount, error) {
	isReal := false
	for _, rel := range realTime {
		if rel {
			isReal = true
			break
		}
	}
	if isReal {
		account, err := accountservice.GetAccountById(accountId)
		if err == nil {
			//save account to redis
			marshal, _ := json.Marshal(account)
			repository.GetRedisRepository().Set(accountId, string(marshal))
		}
		return account, err
	} else {
		//read from redis
		accountString := repository.GetRedisRepository().Get(accountId)
		account := new(mambuEntity.TDAccount)
		err := json.Unmarshal([]byte(accountString), account)
		return account, err
	}

}

type INodeResult interface {
	GetNodeResult() NodeResult
}

type NodeResult string

func (nodeResult NodeResult) GetNodeResult() NodeResult {
	return nodeResult
}
