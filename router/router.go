/*
 * @Author: Hugo
 * @Date: 2022-05-16 08:47:12
 * @Last Modified by: Hugo
 * @Last Modified time: 2022-05-16 10:29:41
 */
package router

import (
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	logger "gitlab.com/bns-engineering/td/common/log"
	"gitlab.com/bns-engineering/td/router/api"
)

// InitRouter initialize routing information
func InitRouter() *gin.Engine {
	r := gin.New()
	r.Use(logger.GinLogger())
	r.Use(logger.GinRecovery(true))

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.POST("/timeDeposite/start", api.StartTDFlow)
	r.POST("/timeDeposite/start2", api.StartFlow)

	transactionGroup := r.Group("/transaction")
	transactionGroup.POST("/retry/:transactionId", api.TransactionRetry)

	return r
}
