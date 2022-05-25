// Package api
// @author： Boice
// @createTime：
package api

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/bns-engineering/td/common/constant"
	"gitlab.com/bns-engineering/td/dao"
	mambuservices "gitlab.com/bns-engineering/td/service/mambuServices"
	"go.uber.org/zap"
	"net/http"
	"strings"
)

func TransactionRetry(c *gin.Context) {
	transactionId := c.Param("transactionId")
	zap.L().Info("transaction retry", zap.String("transactionId", transactionId))
	failTransactionLog := dao.GetFailTransactionLog(transactionId)
	if failTransactionLog == nil {
		c.JSON(http.StatusOK, error("this transaction id not exist"))
		return
	}
	transactionSplit := strings.Split(transactionId, "-")
	flowId, nodeName, transactionType := transactionSplit[0], transactionSplit[1], transactionSplit[2]
	zap.L().Info("retry transaction info", zap.String("flowId", flowId), zap.String("nodeName", nodeName), zap.String("transactionType", transactionType))
	account, err := mambuservices.GetTDAccountById(failTransactionLog.SourceAccountNo)
	if err != nil {
		zap.L().Error("retry transaction can not find source account", zap.String("sourceAccount", failTransactionLog.SourceAccountNo))
		c.JSON(http.StatusOK, error(err.Error()))
		return
	}

	// Get benefit account info
	benefitAccount, err := mambuservices.GetTDAccountById(failTransactionLog.BenefitAccountNo)
	if err != nil {
		zap.L().Error("retry transaction can not find source account", zap.String("sourceAccount", failTransactionLog.SourceAccountNo))
		c.JSON(http.StatusOK, error(err.Error()))
		return
	}

	if constant.WithdrawTransactionType == transactionType {
		_, err := mambuservices.RetryWithDrawTransaction(account, failTransactionLog.Amount, failTransactionLog.TransId, failTransactionLog.Channel)
		if err != nil {
			dao.FailTransactionLogRetryFail(failTransactionLog)
			c.JSON(http.StatusOK, error(err.Error()))
			return
		}
		dao.FailTransactionLogRetrySuccess(failTransactionLog)
	}
	if constant.DepositTransactionType == transactionType {
		_, err := mambuservices.RetryDepositTransaction(account, benefitAccount, failTransactionLog.Amount, failTransactionLog.TransId, failTransactionLog.Channel)
		if err != nil {
			dao.FailTransactionLogRetryFail(failTransactionLog)
			c.JSON(http.StatusOK, error(err.Error()))
			return
		}
		dao.FailTransactionLogRetrySuccess(failTransactionLog)
	}

	c.JSON(http.StatusOK, success())
}
