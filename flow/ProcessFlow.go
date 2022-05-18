/*
 * @Author: Hugo
 * @Date: 2022-05-05 08:59:52
 * @Last Modified by: Hugo
 * @Last Modified time: 2022-05-18 04:16:07
 */
package flow

import (
	"fmt"
	"reflect"
	"sync"

	"github.com/trustmaster/goflow"
	"gitlab.com/bns-engineering/td/common/log"
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
	registerType(new(timeDepositNode.StartNode))
	registerType(new(timeDepositNode.CalcAdditionalProfitNode))
	registerType(new(timeDepositNode.CloseAccNode))
	registerType(new(timeDepositNode.MaturityDateNode))
	registerType(new(timeDepositNode.ProfitApplyNode))
	registerType(new(timeDepositNode.TransferProfitNode))
	registerType(new(timeDepositNode.UpdateAccNode))
	registerType(new(timeDepositNode.WithdrawBalanceNode))
	registerType(new(timeDepositNode.EndNode))
}

func makeInstance(name string) node.NodeRun {
	return typeRegistry[name]
}

// NewProcessFlow defines the app graph
func NewProcessFlow(flowName string) *goflow.Graph {
	flowNodes, flowNodeRelations := dao.GetProcessFlowByName(flowName)
	if len(flowNodes) == 0 || len(flowNodeRelations) == 0 {
		log.Log.Error("Get Flow Nodes Info Error! flowName=%v, len(flowNodes)=%v, len(flowNodeRelations)=%v", flowName, len(flowNodes), len(flowNodeRelations))
		return nil
	}

	n := goflow.NewGraph()
	for _, tmpNode := range flowNodes {
		tmpNodeRun := makeInstance(tmpNode.NodePath)
		fmt.Println("tmpNode.NodeName:", tmpNode.NodeName)
		n.Add(tmpNode.NodeName, tmpNodeRun)
	}

	// Connect them with a channel
	for _, tmpNodeRelation := range flowNodeRelations {
		fmt.Println("tmpNodeRelation.NodeName:", tmpNodeRelation.NodeName)
		n.Connect(tmpNodeRelation.NodeName, "Output", tmpNodeRelation.NextNode, "Input")
	}

	// Our net has 1 inport mapped to greeter.Name
	n.MapInPort("In", "start_node", "Input")
	n.MapInPort("FlowTaskInfo", "start_node", "FlowTaskInfo")
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
