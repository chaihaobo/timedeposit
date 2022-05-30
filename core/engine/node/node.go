// Package node
// @author： Boice
// @createTime：2022/5/26 15:13
package node

import (
	"encoding/json"
	"gitlab.com/bns-engineering/td/core/engine/mambu/accountservice"
	"gitlab.com/bns-engineering/td/core/engine/node/constant"
	"gitlab.com/bns-engineering/td/repository"
	"gitlab.com/bns-engineering/td/service/mambuEntity"
	"go.uber.org/zap"
)

const (
	ResultSuccess NodeResult = "success"
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

func (node *Node) GetMambuBenefitAccountAccount(accountId string, realTime bool) (*mambuEntity.TDAccount, error) {
	if !realTime {
		// from redis
		account := repository.GetRedisRepository().GetBenefitAccount(accountId)
		if account != nil {
			return account, nil
		}
		// from db
		log := repository.GetFlowNodeQueryLogRepository().GetNewLog(node.FlowId, constant.QueryBenefitAccount)
		if log != nil {
			saveDBAccount := new(mambuEntity.TDAccount)
			data := log.Data
			err := json.Unmarshal([]byte(data), saveDBAccount)
			if err != nil {
				zap.L().Error("account from db can not map to struct")
			}
			if saveDBAccount.ID != "" {
				return saveDBAccount, nil
			}
		}
	}
	// real
	id, err := accountservice.GetAccountById(accountId)
	if err == nil {
		marshal, _ := json.Marshal(id)
		err = repository.GetRedisRepository().SaveBenefitAccount(id)
		repository.GetFlowNodeQueryLogRepository().SaveLog(node.FlowId, node.NodeName, constant.QueryBenefitAccount, string(marshal))
	}
	return id, err

}

func (node *Node) GetMambuAccount(accountId string, realTime bool) (*mambuEntity.TDAccount, error) {
	if !realTime {
		// from redis
		account := repository.GetRedisRepository().GetTDAccount(accountId)
		if account != nil {
			return account, nil
		}
		// from db
		log := repository.GetFlowNodeQueryLogRepository().GetNewLog(node.FlowId, constant.QueryTDAccount)
		if log != nil {
			saveDBAccount := new(mambuEntity.TDAccount)
			data := log.Data
			err := json.Unmarshal([]byte(data), saveDBAccount)
			if err != nil {
				zap.L().Error("account from db can not map to struct")
			}
			if saveDBAccount.ID != "" {
				return saveDBAccount, nil
			}
		}
	}
	// real
	id, err := accountservice.GetAccountById(accountId)
	if err == nil {
		marshal, _ := json.Marshal(id)
		err = repository.GetRedisRepository().SaveTDAccount(id)
		repository.GetFlowNodeQueryLogRepository().SaveLog(node.FlowId, node.NodeName, constant.QueryTDAccount, string(marshal))
	}
	return id, err

}

type INodeResult interface {
	GetNodeResult() NodeResult
}

type NodeResult string

func (nodeResult NodeResult) GetNodeResult() NodeResult {
	return nodeResult
}
