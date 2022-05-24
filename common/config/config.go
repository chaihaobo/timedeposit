/*
 * @Author: Hugo
 * @Date: 2022-04-29 10:24:08
 * @Last Modified by:   Hugo
 * @Last Modified time: 2022-04-29 10:24:08
 */
package config

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
)

var TDConf = new(TDConfig)

type TDConfig struct {
	Hugo   string
	System *struct {
		Mode string
	}
	Log *struct {
		Filename   string
		Maxsize    int
		MaxBackups int
		MaxAge     int
		Level      string
	}
	Server *struct {
		RunMode      string
		HttpPort     int
		ReadTimeout  int
		WriteTimeout int
	}
	Db *struct {
		Username    string
		Password    string
		Host        string
		Port        int
		Db          string
		MaxOpenConn int
		MaxIdleConn int
	}
}

func Setup(path string) *TDConfig {
	envConfigPath := os.Getenv("TD_CONFIG_PATH")
	configViper := viper.New()
	configViper.SetConfigFile(path)
	configViper.SetConfigType("yaml")
	if "" != envConfigPath {
		configViper.SetConfigFile(envConfigPath)
	}

	var err error
	if err = configViper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			fmt.Println("Config file not found. Please check the file path...")
		} else {
			fmt.Println("Config file read error...")
		}
	}
	err = configViper.Unmarshal(TDConf)
	if err != nil {
		fmt.Println("config unmarshal error")
	}
	return TDConf
}
