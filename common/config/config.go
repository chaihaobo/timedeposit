/*
 * @Author: Hugo 
 * @Date: 2022-04-29 10:24:08 
 * @Last Modified by:   Hugo 
 * @Last Modified time: 2022-04-29 10:24:08 
 */
package config

import (
	"fmt"

	viper "github.com/spf13/viper"
)

type config struct {
	data *viper.Viper
}

type Config interface {
	GetInt(key string) int64
	GetString(key string) string
	GetBool(key string) bool
	GetFloat(key string) float64
	GetMap(key string) map[string]interface{}
}

func NewConfig(path string) (Config, error) {
	configViper := viper.New()
	configViper.SetConfigFile(path)

	var err error
	if err = configViper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			fmt.Println("Config file not found. Please check the file path...")
		} else {
			fmt.Println("Config file read error...")
		}
	}

	var confData config
	confData.data = configViper
	return &confData, err
}

func (c *config) GetInt(key string) int64 {
	return c.data.GetInt64(key)
}

func (c *config) GetString(key string) string {
	return c.data.GetString(key)
}

func (c *config) GetBool(key string) bool {
	return c.data.GetBool(key)
}
func (c *config) GetFloat(key string) float64 {
	return c.data.GetFloat64(key)
}

func (c *config) GetMap(key string) map[string]interface{} {
	return c.data.GetStringMap(key)
}
