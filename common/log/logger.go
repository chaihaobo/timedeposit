/*
 * @Author: Hugo
 * @Date: 2022-04-29 10:24:24
 * @Last Modified by:   Hugo
 * @Last Modified time: 2022-04-29 10:24:24
 */
package log

import (
	"encoding/json"

	"github.com/astaxie/beego/logs"
	commonConfig "gitlab.com/bns-engineering/td/common/config"
)

var (	
	Log *logs.BeeLogger // Log handler
)

// InitLogConfig initial the log configuration
func InitLogConfig(configData commonConfig.Config) {
	Log = logs.NewLogger(100000)                        // create log hanlder, set the buffer size
	logFileName := configData.GetString("log.filename") // log file name
	maxlines := configData.GetInt("log.maxlines")       // max lines for log file
	maxsize := configData.GetInt("log.maxsize")         // max size of log file
	maxdays := configData.GetInt("log.maxdays")         // log file expire time
	logLevel := configData.GetString("log.logLevel")    // log level

	logConifgMap := make(map[string]interface{})
	logConifgMap["filename"] = logFileName
	logConifgMap["maxlines"] = maxlines
	logConifgMap["maxsize"] = maxsize
	logConifgMap["maxdays"] = maxdays

	// load the config map
	jsonConfig, _ := json.Marshal(logConifgMap)

	Log.SetLogger("file", string(jsonConfig)) // Set the log type：log into file
	Log.SetLevel(getLogLevel(logLevel))       // set the log level
	Log.EnableFuncCallDepth(true)             // whether to show the line number?
}

// Get log level
func getLogLevel(logLevel string) int {
	switch logLevel {
	case "INFO":
		return logs.LevelInfo
	case "WARN":
		return logs.LevelWarn
	case "DEBUG":
		return logs.LevelDebug
	case "ERROR":
		return logs.LevelError
	default:
		return logs.LevelInformational
	}
}
