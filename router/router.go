/*
 * @Author: Hugo
 * @Date: 2022-05-16 08:47:12
 * @Last Modified by: Hugo
 * @Last Modified time: 2022-05-16 10:29:41
 */
package router

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"gitlab.com/bns-engineering/td/common/log"
	"gitlab.com/bns-engineering/td/common/logger"
	"gitlab.com/bns-engineering/td/router/api"
	"net/http"
)

// InitRouter initialize routing information
func InitRouter() *gin.Engine {
	r := gin.New()
	r.Use(logger.GinLogger())
	r.Use(logger.GinRecovery(true))
	//r.Use(middleware.AuthMiddleware())

	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "check success!")
	})

	r.GET("/ping", func(c *gin.Context) {

		ctx:=context.Background()
		err := errors.New("Internal server Error")
		log.Error(ctx, "[UseCaseName] some_error_message : ", err)


		c.JSON(http.StatusOK, api.Success())


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
