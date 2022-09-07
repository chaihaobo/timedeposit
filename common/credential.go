// Package config /*
package common

import (
	"github.com/spf13/viper"
	"log"
	"os"
)

var CredentialConfig Credential

type Credential struct {
	Server struct {
		AuthToken string
	}
	DB struct {
		Username string
		Password string
		Host     string
		Port     int
		Db       string
	}
	Redis struct {
		Addr     string
		Password string
	}
	Mambu struct {
		Host   string
		ApiKey string
	}
}

func newCredential(path string) *Credential {
	configViper := viper.New()
	configViper.SetConfigType("json")

	gsmCredentialPath := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
	parent := os.Getenv("PARENT_SECRET")
	version := os.Getenv("VERSION")
	if "" != gsmCredentialPath && "" != parent && "" != version {
		setupViperGSM(configViper, parent, version)
	} else {
		envConfigPath := os.Getenv("TD_CREDENTIAL_CONFIG_PATH")
		configViper.SetConfigFile(path)
		if "" != envConfigPath {
			configViper.SetConfigFile(envConfigPath)
		}
		var err error
		if err = configViper.ReadInConfig(); err != nil {
			if _, ok := err.(viper.ConfigFileNotFoundError); ok {
				log.Fatalf("Config file not found. Please check the file path...")
			} else {
				log.Fatalf("Config file read error...")
			}
		}
	}
	err := configViper.Unmarshal(&CredentialConfig)
	if err != nil {
		log.Fatalf("config unmarshal error")
	}
	return &CredentialConfig
}
