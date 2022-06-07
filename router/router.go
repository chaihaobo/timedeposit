/*
 * @Author: Hugo
 * @Date: 2022-05-16 08:47:12
 * @Last Modified by: Hugo
 * @Last Modified time: 2022-05-16 10:29:41
 */
package router

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/bns-engineering/td/common/logger"
	"gitlab.com/bns-engineering/td/router/api"
)

// InitRouter initialize routing information
func InitRouter() *gin.Engine {
	r := gin.New()
	r.Use(logger.GinLogger())
	r.Use(logger.GinRecovery(true))

	flowGroup := r.Group("/flow")
	{
		flowGroup.POST("/start", api.StartFlow)
		flowGroup.POST("/failFlows", api.FailFlowList)
		flowGroup.POST("/retry", api.Retry)
		flowGroup.POST("/retryAll", api.RetryAll)
		flowGroup.DELETE("/:flowId", api.Remove)
	}
	return r
}
