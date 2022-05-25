// Package api
// @author： Boice
// @createTime：
package api

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/bns-engineering/td/dao"
	"go.uber.org/zap"
	"net/http"
)

func TransactionRetry(c *gin.Context) {

	transactionId := c.Param("transactionId")
	zap.L().Info("transaction retry", zap.String("transactionId", transactionId))
	failTransactionLog := dao.GetFailTransactionLog(transactionId)
	if failTransactionLog == nil {

	}
	c.JSON(http.StatusOK, success())
}
