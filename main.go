/*
 * @Author: Hugo
 * @Date: 2022-04-29 10:23:50
 * @Last Modified by: Hugo
 * @Last Modified time: 2022-05-16 10:34:03
 */
package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	commonConfig "gitlab.com/bns-engineering/td/common/config"
	commonLog "gitlab.com/bns-engineering/td/common/log"
	"gitlab.com/bns-engineering/td/flow"
	"gitlab.com/bns-engineering/td/router"
)

const (
	filename = "./config.json"
)

var config commonConfig.Config

func initConfig() {
	config, _ = commonConfig.NewConfig("./config.json")
}

// Initial configuration for this app
func init() {
	initConfig()
	commonLog.InitLogConfig(config)
	flow.InitWorkflow()
}

func main() {

	gin.SetMode(config.GetString("server.RunMode"))

	routersInit := router.InitRouter()
	endPoint := fmt.Sprintf(":%d", config.GetInt("ServerSetting.HttpPort"))

	server := &http.Server{
		Addr:    endPoint,
		Handler: routersInit,
	}
	commonLog.Log.Info("[info] start http server listening %s", endPoint)
	server.ListenAndServe()
}

func getConfig() (commonConfig.Config, error) {
	return commonConfig.NewConfig("./config.json")
}
