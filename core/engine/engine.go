// Package engine
// @author： Boice
// @createTime：2022/5/26 10:11
package engine

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"gitlab.com/bns-engineering/common/tracer"
	"gitlab.com/bns-engineering/td/common/config"
	"gitlab.com/bns-engineering/td/common/constant"
	"gitlab.com/bns-engineering/td/common/log"
	"gitlab.com/bns-engineering/td/core/engine/flow"
	"gitlab.com/bns-engineering/td/core/engine/node"
	commomConstant "gitlab.com/bns-engineering/td/core/engine/node/constant"
	"gitlab.com/bns-engineering/td/model/po"
	"gitlab.com/bns-engineering/td/repository"
	"go.uber.org/zap"
	"strings"
	"time"
)

const (
	FlowName  = "eod_flow"
	FirstNode = "start_node"
)

func Start(ctx context.Context, accountId string) error {
	// create task info
	flowId := fmt.Sprintf("%v_%v", time.Now().Format("20060102150405"), accountId)
	createFlowTaskInfo(ctx, flowId, accountId)
	log.Info(ctx, "create task info success!", zap.String("flowId", flowId))
	// run flow by task flow id
	return Run(ctx, flowId)
}

func Run(ctx context.Context, flowId string) error {
	flowTaskInfo := repository.GetFlowTaskInfoRepository().Get(ctx, flowId)
	if flowTaskInfo == nil {
		err := errors.New("invalid flow id")
		log.Error(ctx, "could not find task info by flowId", err, zap.String("flowId", flowId))
		return err
	}
	if flowTaskInfo.CurStatus != string(constant.FlowNodeFailed) && flowTaskInfo.CurStatus != string(constant.FlowNodeStart) {
		err := errors.New("flow status invalid")
		log.Error(ctx, "flow is already running or finished", err, zap.String("curStatus", flowTaskInfo.CurStatus))
		return err
	}
	flowName := flowTaskInfo.FlowName
	nodeName := flowTaskInfo.CurNodeName
	flowNodes := repository.GetFlowNodeRepository().GetFlowNodeListByFlowName(ctx, flowName)
	relationList := repository.GetFlowNodeRelationRepository().GetFlowNodeRelationListByFlowName(ctx, flowName)
	log.Info(ctx, "find engine flow", zap.Int("node size", len(flowNodes)), zap.Int("node relation size", len(relationList)))
	for {
		if nodeName == "" {
			break
		}
		currentNode := getNodeInNodeList(flowNodes, nodeName)
		ctx = getContext(ctx, flowId, flowTaskInfo.AccountId, nodeName)

		tr := tracer.StartTrace(ctx, fmt.Sprintf("%s_%s", flowId, nodeName))
		ctx := tr.Context()
		runNode := getINode(currentNode.NodePath)
		runNode.SetUp(ctx, flowId, flowTaskInfo.AccountId, nodeName)
		// update run status to running
		taskRunning(ctx, flowTaskInfo, nodeName)
		runStartTime := time.Now()
		log.Info(ctx, "flow node run start", zap.String("flowId", flowId), zap.String("currentNodeName", nodeName))
		run, err := runNode.Run(ctx)
		tr.Finish()
		// TODO only use to test case
		if strings.EqualFold(gin.Mode(), "debug") {
			time.Sleep(config.TDConf.Flow.NodeSleepTime)
		}
		if err != nil {
			log.Info(ctx, "node run fail,now retry 3 times", zap.String("current node name", nodeName))
			retry(ctx, func() error {
				run, err = runNode.Run(ctx)
				return err
			}, config.TDConf.Flow.NodeFailRetryTimes, flowId, nodeName)
		}
		useRuntime := time.Since(runStartTime)
		saveNodeRunLog(ctx, flowId, flowTaskInfo.AccountId, flowName, nodeName, run, err)
		if err != nil {
			log.Error(ctx, "flow run failed ", err, zap.String("flowId", flowId), zap.String("currentNodeName", nodeName),
				zap.String("error", fmt.Sprintf("%+v", errors.WithStack(err))),
			)
			taskError(ctx, flowTaskInfo)
			return err
		}
		nodeResult := string(run.GetNodeResult())
		log.Info(ctx, "flow node run finish", zap.String("flowId", flowId), zap.String("currentNodeName", nodeName),
			zap.String("result", nodeResult),
			zap.Int64("useTime", useRuntime.Milliseconds()),
		)
		taskNodeFinish(ctx, flowTaskInfo, nodeResult)
		result := nodeResult
		relation := getNextNodeRelation(nodeName, result, relationList)
		if relation == nil {
			log.Info(ctx, "flow run finished")
			flowRunFinish(ctx, flowTaskInfo)
			break
		}
		nodeName = relation.NextNode
	}
	return nil
}

func getContext(ctx context.Context, flowId string, accountId string, nodeName string) context.Context {
	if ctxFlowId := ctx.Value(commomConstant.ContextFlowId); ctxFlowId == nil {
		ctx = context.WithValue(ctx, commomConstant.ContextFlowId, flowId)
	}
	if cxtAccountId := ctx.Value(commomConstant.ContextAccountId); cxtAccountId == nil {
		ctx = context.WithValue(ctx, commomConstant.ContextAccountId, accountId)
	}
	ctxNodeName := ctx.Value(commomConstant.ContextNodeName)
	if ctxNodeName == nil || ctxNodeName.(string) != nodeName {
		ctx = context.WithValue(ctx, commomConstant.ContextNodeName, nodeName)
	}
	idempotencyKey := repository.GetFlowNodeQueryLogRepository().GetLogValueOr(ctx, flowId, nodeName, commomConstant.QueryIdempotencyKey, uuid.New().String)
	ctx = context.WithValue(ctx, commomConstant.ContextIdempotencyKey, idempotencyKey)
	return ctx
}

func retry(ctx context.Context, retryFun func() error, times int, flowId string, nodeName string) {
	log.Info(ctx, "now retry start ..............", zap.Int("allTime", times), zap.String("flowId", flowId), zap.String("nodeName", nodeName))
	for count := 1; count <= times; count++ {
		log.Info(ctx, "start retry..........", zap.Int("times", count), zap.String("flowId", flowId), zap.String("nodeName", nodeName))
		err := retryFun()
		if err != nil {
			log.Info(ctx, "retry fail........ ", zap.Int("times", count), zap.Error(err), zap.String("flowId", flowId), zap.String("nodeName", nodeName))
		} else {
			log.Info(ctx, "retry success use count ", zap.Int("times", count), zap.String("flowId", flowId), zap.String("nodeName", nodeName))
			break
		}

	}
}

func createFlowTaskInfo(ctx context.Context, flowId string, accountId string) string {
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
	repository.GetFlowTaskInfoRepository().Update(ctx, taskInfo)
	return taskInfo.FlowId
}

func flowRunFinish(ctx context.Context, info *po.TFlowTaskInfo) {
	info.FlowStatus = constant.FlowFinished
	info.UpdateTime = time.Now()
	info.EndTime = time.Now()
	repository.GetFlowTaskInfoRepository().Update(ctx, info)
}

func taskNodeFinish(ctx context.Context, info *po.TFlowTaskInfo, result string) {
	info.CurStatus = string(constant.FlowNodeFinish)
	info.UpdateTime = time.Now()
	info.EndStatus = result
	repository.GetFlowTaskInfoRepository().Update(ctx, info)

}

func taskError(ctx context.Context, taskInfo *po.TFlowTaskInfo) {
	taskInfo.CurStatus = string(constant.FlowNodeFailed)
	taskInfo.FlowStatus = constant.FlowFailed
	taskInfo.EndStatus = constant.FlowFailed
	taskInfo.UpdateTime = time.Now()
	repository.GetFlowTaskInfoRepository().Update(ctx, taskInfo)

}

func saveNodeRunLog(ctx context.Context, flowId string, accountId string, flowName string, nodeName string, nodeResult node.INodeResult, err error) {
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
	repository.GetFlowNodeLogRepository().Save(ctx, log)

}

func taskRunning(ctx context.Context, info *po.TFlowTaskInfo, nodeName string) {
	info.CurNodeName = nodeName
	info.CurStatus = string(constant.FlowNodeRunning)
	info.FlowStatus = constant.FlowRunning
	info.UpdateTime = time.Now()
	repository.GetFlowTaskInfoRepository().Update(ctx, info)
}

func getINode(nodePath string) node.INode {
	return flow.GetNode(nodePath)
}

func getNextNodeRelation(currentNodeName string, resultCode string, nodeRelationList []*po.TFlowNodeRelation) *po.TFlowNodeRelation {
	for _, relation := range nodeRelationList {
		if currentNodeName == relation.NodeName && relation.ResultCode == resultCode {
			return relation
		}
	}
	return nil
}

func getNodeInNodeList(flowNodeList []*po.TFlowNode, nodeName string) *po.TFlowNode {
	for _, flowNode := range flowNodeList {
		if nodeName == flowNode.NodeName {
			return flowNode
		}

	}
	return nil
}

func init() {
	flow.SetUp()
}
