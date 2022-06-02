// Package node
// @author： Boice
// @createTime：2022/5/26 18:29
package node

import (
	"errors"
	"fmt"
	"gitlab.com/bns-engineering/td/common/config"
	"gitlab.com/bns-engineering/td/core/engine/mambu/transactionservice"
	"gitlab.com/bns-engineering/td/core/engine/node/constant"
	"go.uber.org/zap"
)

type DepositNetprofitNode struct {
	*Node
}

func (node *DepositNetprofitNode) Run() (INodeResult, error) {

	account, err := node.GetMambuAccount(node.AccountId, false)
	if err != nil {
		return nil, err
	}
	// Get benefit account info
	benefitAccount, err := node.GetMambuBenefitAccountAccount(account.OtherInformation.BhdNomorRekPencairan, false)
	if err != nil {
		zap.L().Error(fmt.Sprintf("Failed to get benefit acc info of td account: %v, benefit acc id:%v", account.ID, account.OtherInformation.BhdNomorRekPencairan))
		return nil, errors.New("call mambu get benefit acc info failed")
	}
	if !account.IsValidBenefitAccount(benefitAccount, config.TDConf.TransactionReqMetaData.LocalHolderKey) {
		zap.L().Error("is not a valid benefit account!")
		return nil, constant.ErrBenefitAccountInvalid
	}

	if account.IsCaseB1_1() {
		netProfit, err := account.GetNetProfit()
		if err != nil {
			return nil, err
		}
		// Deposit netProfit to benefit account
		channelID := fmt.Sprintf("RAKTRAN_DEPMUDC_%vM", account.OtherInformation.Tenor)
		depositTransID := node.FlowId + "-" + node.NodeName + "-" + "Deposit"
		depositResp, err := transactionservice.DepositTransaction(node.GetContext(), account, benefitAccount, netProfit,
			config.TDConf.TransactionReqMetaData.TranDesc.DepositNetprofitTranDesc1,
			config.TDConf.TransactionReqMetaData.TranDesc.DepositNetprofitTranDesc3,
			depositTransID, channelID, nil)
		if err != nil {
			zap.L().Error(fmt.Sprintf("Failed to deposit for td account: %v", account.ID))
			return nil, errors.New("call mambu deposit failed")
		}
		zap.L().Info(fmt.Sprintf("Finish deposit balance for accNo: %v, encodedKey:%v", account.ID, depositResp.EncodedKey))
	} else {
		zap.L().Info("not match! skip it")
		return ResultSkip, nil
	}

	return ResultSuccess, nil
}
