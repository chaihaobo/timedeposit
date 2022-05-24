/*
 * @Author: Hugo
 * @Date: 2022-04-29 10:23:50
 * @Last Modified by: Hugo
 * @Last Modified time: 2022-05-16 10:34:03
 */
package main

import (
	"fmt"
	"go.uber.org/zap"
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/bns-engineering/td/common/config"
	"gitlab.com/bns-engineering/td/common/log"
	"gitlab.com/bns-engineering/td/flow"
	"gitlab.com/bns-engineering/td/router"
)

// Initial configuration for this app
func init() {
	config.Setup("./config.yaml")
	err := logger.SetUp(config.TDConf)
	if err != nil {
		zap.L().Error("logger init error", zap.Error(err))
	}
	flow.InitWorkflow()
}

func main() {

	gin.SetMode(config.TDConf.Server.RunMode)

	routersInit := router.InitRouter()
	endPoint := fmt.Sprintf(":%d", config.TDConf.Server.HttpPort)

	server := &http.Server{
		Addr:    endPoint,
		Handler: routersInit,
	}
	zap.L().Info("start http server listening ", zap.String("endPoint", endPoint))
	server.ListenAndServe()
}
