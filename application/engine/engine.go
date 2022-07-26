// Package engine
// @author： Boice
// @createTime：2022/7/22 16:22
package engine

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"gitlab.com/bns-engineering/common/tracer"
	"gitlab.com/bns-engineering/td/application/engine/node"
	"gitlab.com/bns-engineering/td/common"
	"gitlab.com/bns-engineering/td/common/constant"
	"gitlab.com/bns-engineering/td/model/po"
	"gitlab.com/bns-engineering/td/repository"
	"gitlab.com/bns-engineering/td/service"
	"go.uber.org/zap"
	"reflect"
	"strings"
	"time"
)

const (
	FlowName  = "eod_flow"
	FirstNode = "start_node"
)

type Engine interface {
	Start(ctx context.Context, accountId string) error
	Run(ctx context.Context, flowId string) error
}

type engine struct {
	common     *common.Common
	repository *repository.Repository
	service    *service.Service
	nodeList   map[string]node.INode
}

func NewEngine(common *common.Common, repository *repository.Repository, service *service.Service) Engine {
	return &engine{
		common:     common,
		repository: repository,
		nodeList: createNodeList(
			new(node.StartNode),
			new(node.EndNode),
			new(node.UndoMaturityNode),
			new(node.StartNewMaturityNode),
			new(node.ApplyProfitNode),
			new(node.WithdrawNetprofitNode),
			new(node.DepositNetprofitNode),
			new(node.WithdrawBalanceNode),
			new(node.PatchAccountNode),
			new(node.CloseAccountNode),
			new(node.AdditionalProfitNode),
		),
		service: service,
	}
}

func createNodeList(list ...node.INode) map[string]node.INode {
	var nodeList = make(map[string]node.INode)
	for _, nodeObj := range list {
		nodeName := reflect.TypeOf(nodeObj).Elem().Name()
		nodeList[nodeName] = nodeObj
	}
	return nodeList
}

func (e *engine) Start(ctx context.Context, accountId string) error {
	// create task info
	flowId := fmt.Sprintf("%v_%v", time.Now().Format("20060102150405"), accountId)
	e.createFlowTaskInfo(ctx, flowId, accountId)
	// save account maturity date
	// saveAccountMaturityDate(ctx, flowId, accountId, maturiryDate)
	e.common.Logger.Info(ctx, "create task info success!", zap.String("flowId", flowId))
	// run flow by task flow id
	return e.Run(ctx, flowId)
}

func (e *engine) Run(ctx context.Context, flowId string) error {
	flowTaskInfo := e.repository.FlowTaskInfo.Get(ctx, flowId)
	if flowTaskInfo == nil {
		err := errors.New("invalid flow id")
		e.common.Logger.Error(ctx, "could not find task info by flowId", err, zap.String("flowId", flowId))
		return err
	}
	if flowTaskInfo.CurStatus != string(constant.FlowNodeFailed) && flowTaskInfo.CurStatus != string(constant.FlowNodeStart) {
		err := errors.New("flow status invalid")
		e.common.Logger.Error(ctx, "flow is already running or finished", err, zap.String("curStatus", flowTaskInfo.CurStatus))
		return err
	}
	flowName := flowTaskInfo.FlowName
	nodeName := flowTaskInfo.CurNodeName
	flowNodes := e.repository.FlowNode.GetFlowNodeListByFlowName(ctx, flowName)
	relationList := e.repository.FlowNodeRelation.GetFlowNodeRelationListByFlowName(ctx, flowName)
	e.common.Logger.Info(ctx, "find engine flow", zap.Int("node size", len(flowNodes)), zap.Int("node relation size", len(relationList)))
	for {
		if nodeName == "" {
			break
		}
		currentNode := getNodeInNodeList(flowNodes, nodeName)
		ctx = e.getContext(ctx, flowId, flowTaskInfo.AccountId, nodeName)

		tr := tracer.StartTrace(ctx, fmt.Sprintf("%s_%s", flowId, nodeName))
		ctx := tr.Context()
		runNode := e.getNode(currentNode.NodePath)
		runNode.SetUp(ctx, *flowTaskInfo, e.common, e.repository, e.service)
		// update run status to running
		e.taskRunning(ctx, flowTaskInfo, nodeName)
		runStartTime := time.Now()
		e.common.Logger.Info(ctx, "flow node run start", zap.String("flowId", flowId), zap.String("currentNodeName", nodeName))
		run, err := runNode.Run(ctx)
		tr.Finish()
		// TODO only use to test case
		if strings.EqualFold(gin.Mode(), "debug") {
			time.Sleep(e.common.Config.Flow.NodeSleepTime)
		}
		if err != nil {
			e.common.Logger.Info(ctx, "node run fail,now retry 3 times", zap.String("current node name", nodeName))
			e.retry(ctx, func() error {
				run, err = runNode.Run(ctx)
				return err
			}, e.common.Config.Flow.NodeFailRetryTimes, flowId, nodeName)
		}
		useRuntime := time.Since(runStartTime)
		e.saveNodeRunLog(ctx, flowId, flowTaskInfo.AccountId, flowName, nodeName, run, err)
		if err != nil {
			e.common.Logger.Error(ctx, "flow run failed ", err, zap.String("flowId", flowId), zap.String("currentNodeName", nodeName),
				zap.String("error", fmt.Sprintf("%+v", errors.WithStack(err))),
			)
			e.taskError(ctx, flowTaskInfo)
			return err
		}
		nodeResult := string(run.GetNodeResult())
		e.common.Logger.Info(ctx, "flow node run finish", zap.String("flowId", flowId), zap.String("currentNodeName", nodeName),
			zap.String("result", nodeResult),
			zap.Int64("useTime", useRuntime.Milliseconds()),
		)
		e.taskNodeFinish(ctx, flowTaskInfo, nodeResult)
		result := nodeResult
		relation := getNextNodeRelation(nodeName, result, relationList)
		if relation == nil {
			e.common.Logger.Info(ctx, "flow run finished")
			e.taskRunFinish(ctx, flowTaskInfo)
			break
		}
		nodeName = relation.NextNode
	}
	return nil
}

func (e *engine) createFlowTaskInfo(ctx context.Context, flowId string, accountId string) string {
	taskInfo := new(po.TFlowTaskInfo)
	taskInfo.FlowId = flowId
	taskInfo.FlowStatus = constant.FlowStart
	taskInfo.FlowName = FlowName
	taskInfo.AccountId = accountId
	taskInfo.CurNodeName = FirstNode
	taskInfo.CurStatus = string(constant.FlowNodeStart)
	taskInfo.StartTime = time.Now()
	taskInfo.EndTime = time.Now()
	taskInfo.CreateTime = time.Now()
	taskInfo.UpdateTime = time.Now()
	taskInfo.Enable = true
	e.repository.FlowTaskInfo.Update(ctx, taskInfo)
	return taskInfo.FlowId
}

func (e *engine) taskRunning(ctx context.Context, info *po.TFlowTaskInfo, nodeName string) {
	info.CurNodeName = nodeName
	info.CurStatus = string(constant.FlowNodeRunning)
	info.FlowStatus = constant.FlowRunning
	info.UpdateTime = time.Now()
	e.repository.FlowTaskInfo.Update(ctx, info)
}

func (e *engine) saveNodeRunLog(ctx context.Context, flowId string, accountId string, flowName string, nodeName string, nodeResult node.INodeResult, err error) {
	log := new(po.TFlowNodeLog)
	log.FlowId = flowId
	log.FlowName = flowName
	log.NodeName = nodeName
	log.AccountId = accountId
	log.CreateTime = time.Now()
	log.UpdateTime = time.Now()
	if nodeResult != nil {
		log.NodeResult = string(nodeResult.GetNodeResult())
	}
	if err != nil {
		log.NodeMsg = fmt.Sprintf("%v", errors.WithStack(err))
		log.NodeResult = "exception"
	}
	e.repository.FlowNodeLog.Save(ctx, log)

}

func (e *engine) taskError(ctx context.Context, taskInfo *po.TFlowTaskInfo) {
	taskInfo.CurStatus = string(constant.FlowNodeFailed)
	taskInfo.FlowStatus = constant.FlowFailed
	taskInfo.EndStatus = constant.FlowFailed
	taskInfo.UpdateTime = time.Now()
	e.repository.FlowTaskInfo.Update(ctx, taskInfo)

}

func (e *engine) taskNodeFinish(ctx context.Context, info *po.TFlowTaskInfo, result string) {
	info.CurStatus = string(constant.FlowNodeFinish)
	info.UpdateTime = time.Now()
	info.EndStatus = result
	e.repository.FlowTaskInfo.Update(ctx, info)

}

func (e *engine) taskRunFinish(ctx context.Context, info *po.TFlowTaskInfo) {
	info.FlowStatus = constant.FlowFinished
	info.UpdateTime = time.Now()
	info.EndTime = time.Now()
	e.repository.FlowTaskInfo.Update(ctx, info)
}

func (e *engine) getContext(ctx context.Context, flowId string, accountId string, nodeName string) context.Context {
	if ctxFlowId := ctx.Value(constant.ContextFlowId); ctxFlowId == nil {
		ctx = context.WithValue(ctx, constant.ContextFlowId, flowId)
	}
	if cxtAccountId := ctx.Value(constant.ContextAccountId); cxtAccountId == nil {
		ctx = context.WithValue(ctx, constant.ContextAccountId, accountId)
	}
	ctxNodeName := ctx.Value(constant.ContextNodeName)
	if ctxNodeName == nil || ctxNodeName.(string) != nodeName {
		ctx = context.WithValue(ctx, constant.ContextNodeName, nodeName)
	}
	idempotencyKey := e.repository.FlowNodeQueryLog.GetLogValueOr(ctx, flowId, nodeName, constant.QueryIdempotencyKey, uuid.New().String)
	ctx = context.WithValue(ctx, constant.ContextIdempotencyKey, idempotencyKey)
	return ctx
}

func (e *engine) getNode(nodeName string) node.INode {
	unKnowNode := e.nodeList[nodeName]
	switch unKnowNode.(type) {
	case *node.StartNode:
		startNode := new(node.StartNode)
		return startNode
	case *node.EndNode:
		endNode := new(node.EndNode)
		return endNode
	case *node.UndoMaturityNode:
		realNode := new(node.UndoMaturityNode)
		return realNode
	case *node.StartNewMaturityNode:
		realNode := new(node.StartNewMaturityNode)
		return realNode
	case *node.ApplyProfitNode:
		realNode := new(node.ApplyProfitNode)
		return realNode
	case *node.WithdrawNetprofitNode:
		realNode := new(node.WithdrawNetprofitNode)
		return realNode
	case *node.DepositNetprofitNode:
		realNode := new(node.DepositNetprofitNode)
		return realNode
	case *node.WithdrawBalanceNode:
		realNode := new(node.WithdrawBalanceNode)
		return realNode
	case *node.PatchAccountNode:
		realNode := new(node.PatchAccountNode)
		return realNode
	case *node.CloseAccountNode:
		realNode := new(node.CloseAccountNode)
		return realNode
	case *node.AdditionalProfitNode:
		realNode := new(node.AdditionalProfitNode)
		return realNode
	default:
		return nil
	}

}

func (e *engine) retry(ctx context.Context, retryFun func() error, times int, flowId string, nodeName string) {
	e.common.Logger.Info(ctx, "now retry start ..............", zap.Int("allTime", times), zap.String("flowId", flowId), zap.String("nodeName", nodeName))
	for count := 1; count <= times; count++ {
		e.common.Logger.Info(ctx, "start retry..........", zap.Int("times", count), zap.String("flowId", flowId), zap.String("nodeName", nodeName))
		err := retryFun()
		if err != nil {
			e.common.Logger.Info(ctx, "retry fail........ ", zap.Int("times", count), zap.Error(err), zap.String("flowId", flowId), zap.String("nodeName", nodeName))
		} else {
			e.common.Logger.Info(ctx, "retry success use count ", zap.Int("times", count), zap.String("flowId", flowId), zap.String("nodeName", nodeName))
			break
		}

	}
}

func getNodeInNodeList(flowNodeList []*po.TFlowNode, nodeName string) *po.TFlowNode {
	for _, flowNode := range flowNodeList {
		if nodeName == flowNode.NodeName {
			return flowNode
		}
	}
	return nil
}

func getNextNodeRelation(currentNodeName string, resultCode string, nodeRelationList []*po.TFlowNodeRelation) *po.TFlowNodeRelation {
	for _, relation := range nodeRelationList {
		if currentNodeName == relation.NodeName && relation.ResultCode == resultCode {
			return relation
		}
	}
	return nil
}
