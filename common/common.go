// Package common
// @author： Boice
// @createTime：2022/7/22 15:27
package common

type Common struct {
	Config     *Config
	Cache      *Cache
	DB         *Database
	Telemetry  *Telemetry
	DataDog    *DataDog
	Credential *Credential
	Logger     Logger
}

func NewCommon(configPath string, credentialPath string) *Common {
	config := newConfig(configPath)
	credential := newCredential(credentialPath)
	cache := newCache(config, credential)
	database := newDatabase(config, credential)
	telemetry := newTelemetry(config, credential)
	dataDog := newDataDog(telemetry.API)
	logger := newLogger(telemetry.API)
	return &Common{
		Config:     config,
		Credential: credential,
		Cache:      cache,
		DB:         database,
		Telemetry:  telemetry,
		DataDog:    dataDog,
		Logger:     logger,
	}

}
