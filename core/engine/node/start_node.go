// Package node
// @author： Boice
// @createTime：2022/5/26 11:07
package node

import (
	"encoding/json"
	"gitlab.com/bns-engineering/td/core/engine/node/constant"
	"gitlab.com/bns-engineering/td/repository"
)

type StartNode struct {
	*Node
}

func (node *StartNode) Run() (INodeResult, error) {
	//query account save account to log
	account, err := node.GetMambuAccount(node.AccountId, true)
	if err != nil {
		return nil, err
	}
	marshal, err := json.Marshal(account)
	if err != nil {
		return nil, err
	}
	repository.GetFlowNodeQueryLogRepository().SaveLog(node.FlowId, node.NodeName, constant.QueryTDAccount, string(marshal))
	return NodeResultSuccess, nil
}
