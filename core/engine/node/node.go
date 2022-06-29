// Package node
// @author： Boice
// @createTime：2022/5/26 15:13
package node

import (
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"gitlab.com/bns-engineering/td/common/log"
	"gitlab.com/bns-engineering/td/core/engine/mambu/accountservice"
	"gitlab.com/bns-engineering/td/core/engine/node/constant"
	"gitlab.com/bns-engineering/td/model/mambu"
	"gitlab.com/bns-engineering/td/repository"
)

const (
	ResultSuccess Result = "success"
	ResultSkip    Result = "skip"
)

type INode interface {
	Run(ctx context.Context) (INodeResult, error)
	SetUp(ctx context.Context, flowId string, accountId string, nodeName string)
}

type INodeResult interface {
	GetNodeResult() Result
}

type Result string

func (nodeResult Result) GetNodeResult() Result {
	return nodeResult
}

type Node struct {
	FlowId    string
	AccountId string
	NodeName  string
}

func (node *Node) SetUp(ctx context.Context, flowId string, accountId string, nodeName string) {
	node.FlowId = flowId
	node.AccountId = accountId
	node.NodeName = nodeName
}

func (node *Node) GetContext(ctx context.Context) context.Context {
	if ctxFlowId := ctx.Value(constant.ContextFlowId); ctxFlowId == nil {
		ctx = context.WithValue(ctx, constant.ContextFlowId, node.FlowId)
	}
	if cxtAccountId := ctx.Value(constant.ContextAccountId); cxtAccountId == nil {
		ctx = context.WithValue(ctx, constant.ContextAccountId, node.AccountId)
	}

	if cxtNodeName := ctx.Value(constant.ContextNodeName); cxtNodeName == nil {
		ctx = context.WithValue(ctx, constant.ContextNodeName, node.NodeName)
	}
	idempotencyKey := repository.GetFlowNodeQueryLogRepository().GetLogValueOr(ctx, node.FlowId, node.NodeName, constant.QueryIdempotencyKey, uuid.New().String)
	ctx = context.WithValue(ctx, constant.ContextIdempotencyKey, idempotencyKey)
	return ctx

}

func (node *Node) GetMambuBenefitAccountAccount(ctx context.Context, accountId string, realTime bool) (*mambu.TDAccount, error) {
	if !realTime {
		// from redis
		account := repository.GetRedisRepository().GetBenefitAccount(ctx, node.AccountId)
		if account != nil {
			return account, nil
		}
		// from db
		flowNodeQueryLog := repository.GetFlowNodeQueryLogRepository().GetNewLog(ctx, node.FlowId, constant.QueryBenefitAccount)
		if flowNodeQueryLog != nil {
			saveDBAccount := new(mambu.TDAccount)
			data := flowNodeQueryLog.Data
			err := json.Unmarshal([]byte(data), saveDBAccount)
			if err != nil {
				log.Error(ctx, "account from db can not map to struct", err)
			}
			if saveDBAccount.ID != "" {
				return saveDBAccount, nil
			}
		}
	}
	// real
	id, err := accountservice.GetAccountById(node.GetContext(ctx), accountId)
	if err == nil {
		marshal, _ := json.Marshal(id)
		err = repository.GetRedisRepository().SaveBenefitAccount(ctx, id, node.FlowId)
		repository.GetFlowNodeQueryLogRepository().SaveLog(ctx, node.FlowId, node.NodeName, constant.QueryBenefitAccount, string(marshal))
	}
	return id, err

}

func (node *Node) GetMambuAccount(ctx context.Context, accountId string, realTime bool) (*mambu.TDAccount, error) {
	if !realTime {
		// from redis
		account := repository.GetRedisRepository().GetTDAccount(ctx, node.FlowId)
		if account != nil {
			return account, nil
		}
		// from db
		nodeQueryLog := repository.GetFlowNodeQueryLogRepository().GetNewLog(ctx, node.FlowId, constant.QueryTDAccount)
		if nodeQueryLog != nil {
			saveDBAccount := new(mambu.TDAccount)
			data := nodeQueryLog.Data
			err := json.Unmarshal([]byte(data), saveDBAccount)
			if err != nil {
				log.Error(ctx, "account from db can not map to struct", err)
			}
			if saveDBAccount.ID != "" {
				return saveDBAccount, nil
			}
		}
	}
	// real
	id, err := accountservice.GetAccountById(node.GetContext(ctx), accountId)
	if err == nil {
		marshal, _ := json.Marshal(id)
		err = repository.GetRedisRepository().SaveTDAccount(ctx, id, node.FlowId)
		repository.GetFlowNodeQueryLogRepository().SaveLog(ctx, node.FlowId, node.NodeName, constant.QueryTDAccount, string(marshal))
	}
	return id, err

}
