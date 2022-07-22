// Package transport
// @author： Boice
// @createTime：2022/6/24 09:49
package transport

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
	"time"
)

type Server interface {
	Start()
	Shutdown()
}

type tdServer struct {
	telemetryApi    *telemetry.API
	telemetryCloser func()
	server          *http.Server
}

func (s *tdServer) Start() {
	cores := runtime.NumCPU()
	runtime.GOMAXPROCS(cores)
	gin.SetMode(config.TDConf.Server.RunMode)
	routersInit := router.InitRouter(s.telemetryApi)
	endPoint := fmt.Sprintf(":%d", config.TDConf.Server.HttpPort)

	s.server = &http.Server{
		Addr:    endPoint,
		Handler: routersInit,
	}
	go func() {
		defer s.telemetryCloser()
		commonlog.Info(context.Background(), "start http server listening", zap.String("endPoint", endPoint))
		err := s.server.ListenAndServe()
		if err != http.ErrServerClosed {
			util.CheckAndExit(err)
		}
	}()
}

func (s *tdServer) Shutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	util.CheckAndExit(s.server.Shutdown(ctx))
}

func NewTdServer(config *config.Config) Server {
	ins, closer := util.SetupTelemetry(config)
	return &tdServer{
		telemetryApi:    ins,
		telemetryCloser: closer,
	}

}
