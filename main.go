/*
 * @Author: Hugo
 * @Date: 2022-04-29 10:23:50
 * @Last Modified by: Hugo
 * @Last Modified time: 2022-05-16 10:34:03
 */
package main

import (
	"github.com/guonaihong/gout/dataflow"
	"gitlab.com/bns-engineering/td/common/config"
	"gitlab.com/bns-engineering/td/transport"
	"time"

	"os"
	"os/signal"
	"syscall"
)

func main() {

	zone := time.FixedZone("CST", 7*3600)
	time.Local = zone
	// set
	dataflow.GlobalSetting.SetTimeout(config.TDConf.Mambu.Timeout)
	server := transport.NewTdServer(config.Setup("./config.json"))
	server.Start()
	// graceful shutdown
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	<-sigs
	// stop gin engine
	server.Shutdown()

}
