// Package node
// @author： Boice
// @createTime：2022/5/27 15:15
package node

import (
	"encoding/json"
	"gitlab.com/bns-engineering/td/core/engine/node/constant"
	"gitlab.com/bns-engineering/td/repository"
)

type GetBenefitAccountNode struct {
	*Node
}

func (node *GetBenefitAccountNode) Run() (INodeResult, error) {
	account, err := node.GetMambuAccount(node.AccountId, false)
	if err != nil {
		return nil, err
	}
	benefitAccountId := account.OtherInformation.BhdNomorRekPencairan
	benefitAccount, err := node.GetMambuBenefitAccountAccount(benefitAccountId, true)
	if err != nil {
		return nil, err
	}
	marshal, err := json.Marshal(benefitAccount)
	if err != nil {
		return nil, err
	}
	repository.GetFlowNodeQueryLogRepository().SaveLog(node.FlowId, node.NodeName, constant.QueryBenefitAccount, string(marshal))
	return NodeResultSuccess, nil
}
