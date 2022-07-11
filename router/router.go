/*
 * @Author: Hugo
 * @Date: 2022-05-16 08:47:12
 * @Last Modified by: Hugo
 * @Last Modified time: 2022-05-16 10:29:41
 */
package router

import (
	"gitlab.com/bns-engineering/td/common/log"
	"gitlab.com/bns-engineering/td/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/bns-engineering/common/telemetry"
	"gitlab.com/bns-engineering/td/router/api"
)

// InitRouter initialize routing information
func InitRouter(telemetry *telemetry.API) *gin.Engine {
	r := gin.New()
	r.Use(log.GinRecovery(true))
	r.Use(middleware.TraceMiddleware(telemetry))
	r.Use(middleware.AuthMiddleware())
	r.Use(log.GinLoggerRequest(telemetry))
	r.Use(log.GinLoggerResponse(telemetry))

	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "check success!")
	})

	flowGroup := r.Group("/flow")
	{
		flowGroup.POST("/start", api.StartFlow)
		flowGroup.GET("/failFlows", api.FailFlowList)
		flowGroup.POST("/retry", api.Retry)
		flowGroup.POST("/retryAll", api.RetryAll)
		flowGroup.DELETE("/:flowId", api.Remove)
	}
	return r
}
