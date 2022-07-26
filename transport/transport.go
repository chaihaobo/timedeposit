// Package transport
// @author： Boice
// @createTime：2022/7/22 15:20
package transport

import (
	"gitlab.com/bns-engineering/td/application"
	"gitlab.com/bns-engineering/td/common"
	"gitlab.com/bns-engineering/td/transport/http"
)

type Transport struct {
	http *http.Http
}

func NewTransport(common *common.Common, app *application.Application) *Transport {
	return &Transport{http: http.NewHttp(common, app)}
}

func (t *Transport) Start() {
	t.http.Serve()
}

func (t *Transport) Stop() {
	t.http.Shutdown()
}
