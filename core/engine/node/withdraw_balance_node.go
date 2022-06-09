// Package node
// @author： Boice
// @createTime：2022/5/26 11:07
package node

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"gitlab.com/bns-engineering/td/common/config"
	"gitlab.com/bns-engineering/td/core/engine/mambu/transactionservice"
	"gitlab.com/bns-engineering/td/core/engine/node/constant"
	"gitlab.com/bns-engineering/td/model/mambu"
	"go.uber.org/zap"
)

type WithdrawBalanceNode struct {
	*Node
}

func (node *WithdrawBalanceNode) Run() (INodeResult, error) {
	account, err := node.GetMambuAccount(node.AccountId, false)
	if err != nil {
		return nil, err
	}
	// Get benefit account info
	benefitAccount, err := node.GetMambuBenefitAccountAccount(account.OtherInformation.BhdNomorRekPencairan, false)
	if err != nil {
		zap.L().Error("Failed to get benefit acc info of td account: %v, benefit acc id:%v", zap.String("account", account.ID), zap.String("benefit acc id", account.OtherInformation.BhdNomorRekPencairan))
		return nil, errors.New("call mambu get benefit acc info failed")
	}

	totalBalance := decimal.NewFromFloat(account.Balances.TotalBalance).RoundFloor(2).InexactFloat64()
	if (account.IsCaseB3() || account.IsCaseC()) && totalBalance > 0 {
		if !account.IsValidBenefitAccount(benefitAccount, config.TDConf.TransactionReqMetaData.LocalHolderKey) {
			zap.L().Error("is not a valid benefit account!")
			return nil, constant.ErrBenefitAccountInvalid
		}
		channelID := fmt.Sprintf("RAKTRAN_DEPMUDC_%vM", account.OtherInformation.Tenor)
		withdrawTransID := node.FlowId + "-" + node.NodeName + "-" + "Withdraw"
		withrawResp, err := transactionservice.WithdrawTransaction(node.GetContext(), account, benefitAccount, totalBalance,
			config.TDConf.TransactionReqMetaData.TranDesc.WithdrawBalanceTranDesc1,
			config.TDConf.TransactionReqMetaData.TranDesc.WithdrawBalanceTranDesc3,
			withdrawTransID, channelID, func(transactionReq *mambu.TransactionReq) {
				transactionReq.Metadata.TranDesc2 = account.ID
			})
		if err != nil {
			zap.L().Error(fmt.Sprintf("Failed to withdraw for td account: %v", account.ID))
			return nil, errors.Wrap(err, "call mambu withdraw failed")
		}
		zap.L().Info(fmt.Sprintf("Finish withdraw balance for accNo: %v, encodedKey:%v", account.ID, withrawResp.EncodedKey))
		depositTransID := node.FlowId + "-" + node.NodeName + "-" + "Deposit"
		// deposit
		depositResp, err := transactionservice.DepositTransaction(node.GetContext(), account, benefitAccount, totalBalance,
			config.TDConf.TransactionReqMetaData.TranDesc.DepositBalanceTranDesc1,
			config.TDConf.TransactionReqMetaData.TranDesc.DepositBalanceTranDesc3,
			depositTransID, channelID, func(transactionReq *mambu.TransactionReq) {
				transactionReq.Metadata.TranDesc2 = account.ID
			})
		if err != nil {
			zap.L().Error(fmt.Sprintf("Failed to deposit for td account: %v", account.ID))
			zap.L().Error(fmt.Sprintf("depositResp: %v", depositResp))

			notes := fmt.Sprintf("eod_engine-%s-%s-%s", node.FlowId, node.NodeName, depositResp.ID)
			// rollback withdraw transaction
			err := transactionservice.AdjustTransaction(node.GetContext(), withrawResp.ID, notes)
			if err != nil {
				err = errors.Wrap(err, "adjust rollback fail!")
			}
			return nil, errors.Wrap(err, "call mambu deposit failed")
		}
		zap.L().Info(fmt.Sprintf("Finish deposit balance for accNo: %v, encodedKey:%v", account.ID, depositResp.EncodedKey))

	} else {
		zap.L().Info("not match! skip it")
		return ResultSkip, nil
	}

	return ResultSuccess, nil
}
