// Package http
// @author： Boice
// @createTime：2022/7/22 15:21
package http

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"gitlab.com/bns-engineering/td/application"
	"gitlab.com/bns-engineering/td/common"
	"gitlab.com/bns-engineering/td/common/util"
	"gitlab.com/bns-engineering/td/middleware"
	"gitlab.com/bns-engineering/td/transport/http/controller"
	"go.uber.org/zap"
	"net/http"
	"runtime"
	"time"
)

type Http struct {
	common      *common.Common
	app         *application.Application
	server      *http.Server
	controllers []controller.Controller
}

func (h *Http) Serve() {
	cores := runtime.NumCPU()
	runtime.GOMAXPROCS(cores)
	gin.SetMode(h.common.Config.Server.RunMode)
	ginEngine := gin.New()
	ginEngine.Use(common.GinRecovery(h.common.Telemetry.API, true))
	ginEngine.Use(middleware.TraceMiddleware(h.common.Telemetry.API))
	ginEngine.Use(middleware.AuthMiddleware(h.common))
	ginEngine.Use(common.GinLoggerRequest(h.common.Telemetry.API))
	ginEngine.Use(common.GinLoggerResponse(h.common.Telemetry.API))
	for _, ctl := range h.controllers {
		ctl.Apply(ginEngine)
	}
	endPoint := fmt.Sprintf(":%d", h.common.Config.Server.HttpPort)
	h.server = &http.Server{
		Addr:    endPoint,
		Handler: ginEngine,
	}
	go func() {
		defer h.common.Telemetry.Closer()
		h.common.Logger.Info(context.Background(), "start http server listening", zap.String("endPoint", endPoint))
		err := h.server.ListenAndServe()
		if err != http.ErrServerClosed {
			util.CheckAndExit(err)
		}
	}()
}

func (h *Http) Shutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	util.CheckAndExit(h.server.Shutdown(ctx))
}

func NewHttp(common *common.Common, app *application.Application) *Http {
	controllers := make([]controller.Controller, 0)
	controllers = append(controllers,
		controller.NewFlowController(common, app),
	)
	return &Http{
		common:      common,
		app:         app,
		controllers: controllers,
	}
}
