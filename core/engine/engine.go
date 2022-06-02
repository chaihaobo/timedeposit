// Package engine
// @author： Boice
// @createTime：2022/5/26 10:11
package engine

import (
	"fmt"
	"github.com/panjf2000/ants/v2"
	"github.com/pkg/errors"
	"gitlab.com/bns-engineering/td/common/constant"
	"gitlab.com/bns-engineering/td/core/engine/flow"
	"gitlab.com/bns-engineering/td/core/engine/node"
	"gitlab.com/bns-engineering/td/model/po"
	"gitlab.com/bns-engineering/td/repository"
	"go.uber.org/zap"
	"time"
)

const (
	FlowName  = "eod_flow"
	FirstNode = "start_node"
)

var Pool *ants.PoolWithFunc
var RetryPool *ants.PoolWithFunc

func Start(accountId string) {

	// create task info
	flowId := fmt.Sprintf("%v_%v", time.Now().Format("20060102150405"), accountId)
	createFlowTaskInfo(flowId, accountId)
	zap.L().Info("create task info success!", zap.String("flowId", flowId))
	// run flow by task flow id
	Run(flowId)
}

func Run(flowId string) {
	flowTaskInfo := repository.GetFlowTaskInfoRepository().Get(flowId)
	if flowTaskInfo == nil {
		zap.L().Error("could not find task info by flowId", zap.String("flowId", flowId))
		return
	}
	if flowTaskInfo.CurStatus != string(constant.FlowNodeFailed) && flowTaskInfo.CurStatus != string(constant.FlowNodeStart) {
		zap.L().Error("flow is already running or finished", zap.String("curStatus", flowTaskInfo.CurStatus))
		return
	}

	flowName := flowTaskInfo.FlowName
	nodeName := flowTaskInfo.CurNodeName
	flowNodes := repository.GetFlowNodeRepository().GetFlowNodeListByFlowName(flowName)
	relationList := repository.GetFlowNodeRelationRepository().GetFlowNodeRelationListByFlowName(flowName)
	zap.L().Info("find engine flow", zap.Int("node size", len(flowNodes)), zap.Int("node relation size", len(relationList)))
	for {
		if nodeName == "" {
			break
		}
		currentNode := getNodeInNodeList(flowNodes, nodeName)

		runNode := getINode(currentNode.NodePath)

		runNode.SetUp(flowId, flowTaskInfo.AccountId, nodeName)
		// update run status to running
		taskRunning(flowTaskInfo, nodeName)
		runStartTime := time.Now()
		zap.L().Info("flow node run start", zap.String("flowId", flowId), zap.String("currentNodeName", nodeName))
		run, err := runNode.Run()
		if err != nil {
			zap.L().Info("node run fail,now retry 3 times", zap.String("current node name", nodeName))
			retry(func() error {
				run, err = runNode.Run()
				return err
			}, 3)
		}
		useRuntime := time.Now().Sub(runStartTime)
		saveNodeRunLog(flowId, flowTaskInfo.AccountId, flowName, nodeName, run, err)
		if err != nil {
			zap.L().Error("flow run failed ", zap.String("flowId", flowId), zap.String("currentNodeName", nodeName),
				zap.String("error", fmt.Sprintf("%v", errors.WithStack(err))),
			)
			taskError(flowTaskInfo)
			break
		}
		nodeResult := string(run.GetNodeResult())
		zap.L().Info("flow node run finish", zap.String("flowId", flowId), zap.String("currentNodeName", nodeName),
			zap.String("result", nodeResult),
			zap.Int64("useTime", useRuntime.Milliseconds()),
		)
		taskNodeFinish(flowTaskInfo, nodeResult)
		result := nodeResult
		relation := getNextNodeRelation(nodeName, result, relationList)
		if relation == nil {
			zap.L().Info("flow run finished")
			flowRunFinish(flowTaskInfo)
			break
		}
		nodeName = relation.NextNode
	}
}

func retry(retryFun func() error, times int) {
	zap.L().Info("now retry start ..............", zap.Int("allTime", times))
	for count := 0; count <= times; count++ {
		zap.L().Info("start retry..........", zap.Int("times", count))
		err := retryFun()
		if err != nil {
			zap.L().Info("retry fail........ ", zap.Int("times", count), zap.Error(err))
		} else {
			zap.L().Info("retry success use count ", zap.Int("times", count))
		}

	}
}

func createFlowTaskInfo(flowId string, accountId string) string {
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
	repository.GetFlowTaskInfoRepository().Update(taskInfo)
	return taskInfo.FlowId
}

func flowRunFinish(info *po.TFlowTaskInfo) {
	info.FlowStatus = constant.FlowFinished
	info.UpdateTime = time.Now()
	info.EndTime = time.Now()
	repository.GetFlowTaskInfoRepository().Update(info)
}

func taskNodeFinish(info *po.TFlowTaskInfo, result string) {
	info.CurStatus = string(constant.FlowNodeFinish)
	info.UpdateTime = time.Now()
	info.EndStatus = result
	repository.GetFlowTaskInfoRepository().Update(info)

}

func taskError(taskInfo *po.TFlowTaskInfo) {
	taskInfo.CurStatus = string(constant.FlowNodeFailed)
	taskInfo.FlowStatus = constant.FlowFailed
	taskInfo.EndStatus = constant.FlowFailed
	taskInfo.UpdateTime = time.Now()
	repository.GetFlowTaskInfoRepository().Update(taskInfo)

}

func saveNodeRunLog(flowId string, accountId string, flowName string, nodeName string, nodeResult node.INodeResult, err error) {
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
	repository.GetFlowNodeLogRepository().Save(log)

}

func taskRunning(info *po.TFlowTaskInfo, nodeName string) {
	info.CurNodeName = nodeName
	info.CurStatus = string(constant.FlowNodeRunning)
	info.FlowStatus = constant.FlowRunning
	info.UpdateTime = time.Now()
	repository.GetFlowTaskInfoRepository().Update(info)
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
	Pool, _ = ants.NewPoolWithFunc(100, func(accountId interface{}) {
		Start(accountId.(string))
	})

	RetryPool, _ = ants.NewPoolWithFunc(100, func(flowId interface{}) {
		Run(flowId.(string))
	})
}
