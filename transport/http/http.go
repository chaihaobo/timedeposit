// Package http
// @author： Boice
// @createTime：2022/7/22 15:21
package http

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"gitlab.com/bns-engineering/common/telemetry"
	"gitlab.com/bns-engineering/td/common/config"
	commonlog "gitlab.com/bns-engineering/td/common/log"
	"gitlab.com/bns-engineering/td/common/util"
	"gitlab.com/bns-engineering/td/router"
	"go.uber.org/zap"
	"net/http"
	"runtime"
)

type Http struct {
	telemetryApi    *telemetry.API
	telemetryCloser func()
}

func (h Http) Serve() {
	cores := runtime.NumCPU()
	runtime.GOMAXPROCS(cores)
	gin.SetMode(config.TDConf.Server.RunMode)
	routersInit := router.InitRouter(h.telemetryApi)
	endPoint := fmt.Sprintf(":%d", config.TDConf.Server.HttpPort)

	server := &http.Server{
		Addr:    endPoint,
		Handler: routersInit,
	}
	err := server.ListenAndServe()
	commonlog.Info(context.Background(), "start http server listening", zap.String("endPoint", endPoint))
	if err != http.ErrServerClosed {
		util.CheckAndExit(err)
	}
}
