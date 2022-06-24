// Package node
// @author： Boice
// @createTime：2022/5/26 11:07
package node

import "context"

type StartNode struct {
	*Node
}

func (node *StartNode) Run(ctx context.Context) (INodeResult, error) {
	// query account save account to log
	// account, err := node.GetMambuAccount(ctx,node.AccountId, true)
	// if err != nil {
	//	return nil, err
	// }
	// marshal, err := json.Marshal(account)
	// if err != nil {
	//	return nil, err
	// }
	// repository.GetFlowNodeQueryLogRepository().SaveLog(node.FlowId, node.NodeName, constant.QueryTDAccount, string(marshal))
	return ResultSuccess, nil
}
