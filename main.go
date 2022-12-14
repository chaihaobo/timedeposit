/*
 * @Author: Hugo
 * @Date: 2022-04-29 10:23:50
 * @Last Modified by: Hugo
 * @Last Modified time: 2022-05-16 10:34:03
 */
package main

import (
	"gitlab.com/bns-engineering/td/application"
	"gitlab.com/bns-engineering/td/common"
	"gitlab.com/bns-engineering/td/repository"
	"gitlab.com/bns-engineering/td/service"
	"gitlab.com/bns-engineering/td/transport"
	"time"

	"os"
	"os/signal"
	"syscall"
)

func main() {

	zone := time.FixedZone("CST", 7*3600)
	time.Local = zone
	com := common.NewCommon("./config.json", "./credential.json")
	repo := repository.NewRepository(com)
	svc := service.NewService(com, repo)
	app := application.NewApplication(com, repo, svc)
	tp := transport.NewTransport(com, app)
	tp.Start()
	// graceful shutdown
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs
	// stop gin engine
	tp.Stop()

}
