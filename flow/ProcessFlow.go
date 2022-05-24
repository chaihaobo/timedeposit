/*
 * @Author: Hugo
 * @Date: 2022-05-05 08:59:52
 * @Last Modified by: Hugo
 * @Last Modified time: 2022-05-24 01:18:51
 */
package flow

import (
	"fmt"
	"go.uber.org/zap"
	"reflect"
	"sync"

	"github.com/trustmaster/goflow"
	"gitlab.com/bns-engineering/td/dao"
	"gitlab.com/bns-engineering/td/node"
	"gitlab.com/bns-engineering/td/node/timeDepositNode"
)

var typeRegistry = make(map[string]node.NodeRun)

func registerType(typeStruct node.NodeRun) {
	t := reflect.TypeOf(typeStruct).Elem()
	typeRegistry[t.PkgPath()+"."+t.Name()] = typeStruct
}

func InitWorkflow() {
	_workFlowDic = make(map[string]*goflow.Graph)
	registerType(timeDepositNode.NewStartNode())
	registerType(timeDepositNode.NewMaturityDateNode())
	registerType(timeDepositNode.NewCalcAdditionalProfitNode())
	registerType(timeDepositNode.NewCloseAccNode())
	registerType(timeDepositNode.NewProfitApplyNode())
	registerType(timeDepositNode.NewTransferProfitNode())
	registerType(timeDepositNode.NewUpdateAccNode())
	registerType(timeDepositNode.NewWithdrawBalanceNode())
	registerType(timeDepositNode.NewEndNode())
}

func makeInstance(name string) node.NodeRun {
	return typeRegistry[name]
}

// NewProcessFlow defines the app graph
func NewProcessFlow(flowName string) *goflow.Graph {
	flowNodes, flowNodeRelations := dao.GetProcessFlowByName(flowName)
	if len(flowNodes) == 0 || len(flowNodeRelations) == 0 {
		zap.L().Error(fmt.Sprintf("Get Flow Nodes Info Error! flowName=%v, len(flowNodes)=%v, len(flowNodeRelations)=%v", flowName, len(flowNodes), len(flowNodeRelations)))
		return nil
	}

	n := goflow.NewGraph()
	for _, tmpNode := range flowNodes {
		tmpNodeRun := makeInstance(tmpNode.NodePath)
		zap.L().Info("tmpNode.NodeName", zap.String("NodeName", tmpNode.NodeName))
		n.Add(tmpNode.NodeName, tmpNodeRun)
	}

	// Connect them with a channel
	for _, tmpNodeRelation := range flowNodeRelations {
		fmt.Println("tmpNodeRelation.NodeName:", tmpNodeRelation.NodeName)
		n.Connect(tmpNodeRelation.NodeName, "Output", tmpNodeRelation.NextNode, "Input")
	}

	// Our net has 1 inport mapped to greeter.Name
	n.MapInPort("In", "start_node", "Input")
	return n
}

// define work flow map, all the working flow will use the graph here
var _workFlowDic map[string]*goflow.Graph
var once sync.Once

// Get the flow by name
func GetProcessFlow(flowName string) *goflow.Graph {
	if tmpWorkFlow, ok := _workFlowDic[flowName]; ok {
		return tmpWorkFlow
	} else {
		once.Do(func() {
			tmpWorkFlow = NewProcessFlow(flowName)
			_workFlowDic[flowName] = tmpWorkFlow
		})
		return tmpWorkFlow
	}

}
