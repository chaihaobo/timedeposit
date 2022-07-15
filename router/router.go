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
		flowGroup.POST("/start", GinWrapper(api.StartFlow))
		flowGroup.GET("/failFlows", GinWrapper(api.FailFlowList))
		flowGroup.POST("/retry", GinWrapper(api.Retry))
		flowGroup.POST("/retryAll", GinWrapper(api.RetryAll))
		flowGroup.DELETE("/:flowId", GinWrapper(api.Remove))
		flowGroup.POST("/metric", GinWrapper(api.Metric))
	}
	return r
}
