// Package config /*
package config

import (
	"bytes"
	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"context"
	"github.com/spf13/viper"
	secretmanagerpb "google.golang.org/genproto/googleapis/cloud/secretmanager/v1"
	"log"
	"os"
	"time"
)

var TDConf = new(TDConfig)

type TDConfig struct {
	Trace *struct {
		CollectorURL string
		ServiceName  string
		SourceEnv    string
	}
	Metric *struct {
		Port         int
		AgentAddress string
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
		AuthToken    string
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
		NodeFailRetryTimes    int
		MaxLimitSearchAccount int32
		NodeSleepTime         time.Duration
	}
	TransactionReqMetaData *struct {
		MessageType                    string
		LocalHolderKey                 string
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
		Host    string
		ApiKey  string
		Timeout time.Duration
	}
	SkipTests bool
}

func setupViperGSM(viper *viper.Viper, parent string, version string) {
	// Create the client.
	ctx := context.Background()
	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		log.Fatalf("failed to setup client: %v", err)
	}
	defer client.Close()

	// Build the request.
	accessRequest := &secretmanagerpb.AccessSecretVersionRequest{
		Name: parent + "/" + version,
	}
	// Call the API.
	result, err := client.AccessSecretVersion(ctx, accessRequest)
	if err != nil {
		log.Fatalf("failed to access secret version: %v", err)
	}
	err = viper.MergeConfig(bytes.NewReader(result.Payload.Data))
	if err != nil {
		log.Fatalf("failed to merge gsm config to viper %v", err)
	}
}

func Setup(path string) *TDConfig {
	configViper := viper.New()
	configViper.SetConfigType("json")

	gsmCredentialPath := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
	parent := os.Getenv("PARENT")
	version := os.Getenv("VERSION")
	if "" != gsmCredentialPath && "" != parent && "" != version {
		setupViperGSM(configViper, parent, version)

	} else {
		envConfigPath := os.Getenv("TD_CONFIG_PATH")
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
	err := configViper.Unmarshal(TDConf)
	if err != nil {
		log.Fatalf("config unmarshal error")
	}
	return TDConf
}
