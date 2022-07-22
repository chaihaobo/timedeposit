// Package common
// @author： Boice
// @createTime：2022/7/22 15:27
package common

import (
	"gitlab.com/bns-engineering/td/common/cache"
	"gitlab.com/bns-engineering/td/common/config"
	"gitlab.com/bns-engineering/td/common/database"
)

type Common struct {
	Config    *config.Config
	Cache     *cache.Cache
	DB        *database.Database
	Telemetry *Telemetry
	DataDog   *DataDog
	Logger    Logger
}

func (c Common) NewCommon(configPath string) *Common {
	config := config.NewConfig(configPath)
	cache := cache.NewCache(config)
	database := database.NewDatabase(config)
	telemetry := NewTelemetry(config)
	dataDog := NewDataDog(telemetry.API)
	logger := NewLogger(telemetry.API)
	return &Common{
		Config:    config,
		Cache:     cache,
		DB:        database,
		Telemetry: telemetry,
		DataDog:   dataDog,
		Logger:    logger,
	}

}
