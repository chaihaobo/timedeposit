/*
 * @Author: Hugo
 * @Date: 2022-04-29 10:24:08
 * @Last Modified by: Hugo
 * @Last Modified time: 2022-05-06 11:32:01
 */
package config

import (
	"log"
	"testing"
)

func TestNewConfig(t *testing.T) {

	t.Run("test config read", func(t *testing.T) {
		var err error
		var conf Config
		conf, err = NewConfig("./../../config.yaml")
		if err != nil {
			log.Println(err)
			return
		}
		log.Println(TDConf.Hugo)
		log.Println(conf)
		log.Println(TDConf.Server.RunMode)

	})

}
