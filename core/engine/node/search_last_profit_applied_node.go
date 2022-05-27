// Package node
// @author： Boice
// @createTime：2022/5/27 09:18
package node

import (
	"encoding/json"
	"errors"
	"gitlab.com/bns-engineering/td/core/engine/mambu/transactionservice"
	"gitlab.com/bns-engineering/td/core/engine/node/constant"
	"gitlab.com/bns-engineering/td/repository"
	"go.uber.org/zap"
	"strings"
)

type SearchLastProfitAppliedNode struct {
	*Node
}

func (node *SearchLastProfitAppliedNode) Run() (INodeResult, error) {
	account, err := node.GetMambuAccount(node.AccountId, false)
	if err != nil {
		return nil, err
	}
	if account.IsCaseB1_1_1_1() ||
		account.IsCaseB2_1_1() ||
		(account.IsCaseB3() &&
			account.Balances.TotalBalance > 0 &&
			strings.ToUpper(account.OtherInformation.IsSpecialRate) == "TRUE") ||
		(account.IsCaseC() &&
			account.Balances.TotalBalance > 0 &&
			strings.ToUpper(account.OtherInformation.IsSpecialRate) == "TRUE") {
		// Get last applied interest info
		transList, err := transactionservice.GetTransactionByQueryParam(account.EncodedKey)
		if err != nil || len(transList) <= 0 {
			zap.L().Info("No applied profit, skip")
			return nil, errors.New("No applied profit find")
		}
		lastAppliedInterestTrans := transList[0]
		bytes, _ := json.Marshal(&lastAppliedInterestTrans)
		repository.GetFlowNodeQueryLogRepository().SaveLog(node.FlowId, node.NodeName, constant.QueryLastProfitApplied, string(bytes))

	} else {
		zap.L().Info("not match! skip it")
	}
	return NodeResultSuccess, nil

}
