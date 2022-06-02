/*
 * @Author: Hugo
 * @Date: 2022-04-29 10:23:50
 * @Last Modified by: Hugo
 * @Last Modified time: 2022-05-16 10:34:03
 */
package main

import (
	"context"
	"fmt"
	"gitlab.com/bns-engineering/td/common/logger"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"gitlab.com/bns-engineering/td/common/util"
	"gitlab.com/bns-engineering/td/core/engine"
	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
	"gitlab.com/bns-engineering/td/common/config"
	"gitlab.com/bns-engineering/td/router"
)

// Initial configuration for this app
func init() {
	err := logger.SetUp(config.Setup("./config.yaml"))
	if err != nil {
		zap.L().Error("logger init error", zap.Error(err))
	}
}

func main() {
	cores := runtime.NumCPU()
	runtime.GOMAXPROCS(cores)
	gin.SetMode(config.TDConf.Server.RunMode)
	routersInit := router.InitRouter()
	endPoint := fmt.Sprintf(":%d", config.TDConf.Server.HttpPort)

	server := &http.Server{
		Addr:    endPoint,
		Handler: routersInit,
	}

	go func() {
		zap.L().Info("start http server listening ", zap.String("endPoint", endPoint))
		// service connections
		err := server.ListenAndServe()
		if err != http.ErrServerClosed {
			util.CheckAndExit(err)
		}
	}()

	// graceful shutdown
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	<-sigs
	// stop gin engine
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	util.CheckAndExit(server.Shutdown(ctx))
	engine.Pool.Release()

}
