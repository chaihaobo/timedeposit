/*
 * @Author: Hugo
 * @Date: 2022-04-29 10:24:08
 * @Last Modified by: Hugo
 * @Last Modified time: 2022-05-06 11:32:01
 */
package config

import (
	"fmt"
	"log"
	"testing"
)

func TestNewConfig(t *testing.T) {
	var err error
	var conf Config
	conf, err = NewConfig("./../../config.json")
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println(conf.GetString("hugo"))
	fmt.Println(conf.GetString("system.mode"))
	fmt.Println(conf.GetString("log.prefix"))
}
