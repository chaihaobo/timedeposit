// Package config /*
package config

import (
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"os"
)

var TDConf = new(TDConfig)

type TDConfig struct {
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
	Flow *struct {
		NodeFailRetryTimes int
	}
	TransactionReqMetaData *struct {
		MessageType                    string
		ExternalOriTransactionID       string
		ExternalOriTransactionDetailID string
		TransactionType                string
		TerminalType                   string
		TerminalID                     string
		TerminalLocation               string
		ProductCode                    string
		AcquirerIID                    string
		ForwarderIID                   string
		IssuerIID                      string
		IssuerIName                    string
		DestinationIID                 string
		Currency                       string
		TranDesc                       *struct {
			WithdrawAdditionalProfitTranDesc1 string
			WithdrawAdditionalProfitTranDesc3 string
			WithdrawBalanceTranDesc1          string
			WithdrawBalanceTranDesc3          string
			WithdrawNetprofitTranDesc1        string
			WithdrawNetprofitTranDesc3        string
			DepositAdditionalProfitTranDesc1  string
			DepositAdditionalProfitTranDesc3  string
			DepositBalanceTranDesc1           string
			DepositBalanceTranDesc3           string
			DepositNetprofitTranDesc1         string
			DepositNetprofitTranDesc3         string
		}
	}
	Redis *struct {
		Addr     string
		Password string
		PoolSize int
		DB       int
	}
	Mambu *struct {
		Host   string
		ApiKey string
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
			zap.L().Error("Config file not found. Please check the file path...")
		} else {
			zap.L().Error("Config file read error...")
		}
	}
	err = configViper.Unmarshal(TDConf)
	if err != nil {
		zap.L().Error("config unmarshal error")
	}
	return TDConf
}
