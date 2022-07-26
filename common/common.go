// Package common
// @author： Boice
// @createTime：2022/7/22 15:27
package common

type Common struct {
	Config    *Config
	Cache     *Cache
	DB        *Database
	Telemetry *Telemetry
	DataDog   *DataDog
	Logger    Logger
}

func NewCommon(configPath string) *Common {
	config := newConfig(configPath)
	cache := newCache(config)
	database := newDatabase(config)
	telemetry := newTelemetry(config)
	dataDog := newDataDog(telemetry.API)
	logger := newLogger(telemetry.API)
	return &Common{
		Config:    config,
		Cache:     cache,
		DB:        database,
		Telemetry: telemetry,
		DataDog:   dataDog,
		Logger:    logger,
	}

}
