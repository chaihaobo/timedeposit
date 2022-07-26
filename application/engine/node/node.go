// Package node
// @author： Boice
// @createTime：2022/5/26 15:13
package node

import (
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"gitlab.com/bns-engineering/td/common"
	constant "gitlab.com/bns-engineering/td/common/constant"
	"gitlab.com/bns-engineering/td/model/mambu"
	"gitlab.com/bns-engineering/td/model/po"
	"gitlab.com/bns-engineering/td/repository"
	"gitlab.com/bns-engineering/td/service"
	"time"
)

const (
	ResultSuccess Result = "success"
	ResultSkip    Result = "skip"
)

type INode interface {
	Run(ctx context.Context) (INodeResult, error)
	SetUp(ctx context.Context, taskInfo po.TFlowTaskInfo, common *common.Common, repo *repository.Repository, service *service.Service)
}

type INodeResult interface {
	GetNodeResult() Result
}

type Result string

func (nodeResult Result) GetNodeResult() Result {
	return nodeResult
}

type node struct {
	ctx            context.Context
	flowId         string
	accountId      string
	nodeName       string
	taskCreateTime time.Time
	common         *common.Common
	repository     *repository.Repository
	service        *service.Service
}

func (node *node) SetUp(ctx context.Context, taskInfo po.TFlowTaskInfo, common *common.Common, repo *repository.Repository, service *service.Service) {
	node.ctx = ctx
	node.flowId = taskInfo.FlowId
	node.accountId = taskInfo.AccountId
	node.nodeName = taskInfo.CurNodeName
	node.taskCreateTime = taskInfo.CreateTime
	node.common = common
	node.repository = repo
	node.service = service
}

func (node *node) GetContext(ctx context.Context) context.Context {
	if ctxFlowId := ctx.Value(constant.ContextFlowId); ctxFlowId == nil {
		ctx = context.WithValue(ctx, constant.ContextFlowId, node.flowId)
	}
	if cxtAccountId := ctx.Value(constant.ContextAccountId); cxtAccountId == nil {
		ctx = context.WithValue(ctx, constant.ContextAccountId, node.accountId)
	}

	if cxtNodeName := ctx.Value(constant.ContextNodeName); cxtNodeName == nil {
		ctx = context.WithValue(ctx, constant.ContextNodeName, node.nodeName)
	}
	idempotencyKey := node.repository.FlowNodeQueryLog.GetLogValueOr(ctx, node.flowId, node.nodeName, constant.QueryIdempotencyKey, uuid.New().String)
	ctx = context.WithValue(ctx, constant.ContextIdempotencyKey, idempotencyKey)
	return ctx

}

func (node *node) GetMambuBenefitAccountAccount(ctx context.Context, accountId string) (*mambu.TDAccount, error) {
	id, err := node.service.Mambu.Account.GetAccountById(node.GetContext(ctx), accountId)
	return id, err

}

func (node *node) GetMambuAccount(ctx context.Context, accountId string, realTime bool) (*mambu.TDAccount, error) {
	// when retry get new from mambu
	if !realTime {
		// from redis
		account := node.repository.Redis.GetTDAccount(ctx, node.flowId)
		if account != nil {
			return node.loadRealTimeFields(ctx, account)
		}
		// from db
		nodeQueryLog := node.repository.FlowNodeQueryLog.GetNewLog(ctx, node.flowId, constant.QueryTDAccount)
		if nodeQueryLog != nil {
			saveDBAccount := new(mambu.TDAccount)
			data := nodeQueryLog.Data
			err := json.Unmarshal([]byte(data), saveDBAccount)
			if err != nil {
				node.common.Logger.Error(ctx, "account from db can not map to struct", err)
			}
			if saveDBAccount.ID != "" {
				return node.loadRealTimeFields(ctx, saveDBAccount)
			}
		}
	}
	// real
	id, err := node.service.Mambu.Account.GetAccountById(node.GetContext(ctx), accountId)
	if err == nil {
		marshal, _ := json.Marshal(id)
		err = node.repository.Redis.SaveTDAccount(ctx, id, node.flowId)
		node.repository.FlowNodeQueryLog.SaveLog(ctx, node.flowId, node.nodeName, constant.QueryTDAccount, string(marshal))
	}
	return id, err

}

func (node *node) loadRealTimeFields(ctx context.Context, account *mambu.TDAccount) (*mambu.TDAccount, error) {
	realTimeAccount, err := node.service.Mambu.Account.GetAccountById(ctx, account.ID)
	if err != nil {
		node.common.Logger.Error(ctx, "get real time account failed", err)
		return nil, err
	}
	account.OtherInformation.BhdNomorRekPencairan = realTimeAccount.OtherInformation.BhdNomorRekPencairan
	return account, nil
}
