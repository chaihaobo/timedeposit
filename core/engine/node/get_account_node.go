// Package node
// @author： Boice
// @createTime：2022/5/27 15:15
package node

import (
	"encoding/json"
	"gitlab.com/bns-engineering/td/core/engine/node/constant"
	"gitlab.com/bns-engineering/td/repository"
)

type GetAccountNode struct {
	*Node
}

func (node *GetAccountNode) Run() (INodeResult, error) {
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
