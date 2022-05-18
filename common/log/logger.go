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
	commonConfig "gitlab.com/hugo.hu/time-deposit-eod-engine/common/config"
)

var (
	// Log 日志记录句柄
	Log *logs.BeeLogger //日志记录句柄
)

// InitLogConfig 初始化Log设置
func InitLogConfig(configData commonConfig.Config) {
	Log = logs.NewLogger(100000)                        // 创建一个日志记录器，参数为缓冲区的大小
	logFileName := configData.GetString("log.filename") // 日志文件名
	maxlines := configData.GetInt("log.maxlines")       // 日志最大行数
	maxsize := configData.GetInt("log.maxsize")         // 日志文件大小的最大值
	maxdays := configData.GetInt("log.maxdays")         // 日志文件保留最长周期
	logLevel := configData.GetString("log.logLevel")    // 日志级别

	logConifgMap := make(map[string]interface{})
	logConifgMap["filename"] = logFileName
	logConifgMap["maxlines"] = maxlines
	logConifgMap["maxsize"] = maxsize
	logConifgMap["maxdays"] = maxdays

	// 设置配置文
	jsonConfig, _ := json.Marshal(logConifgMap)

	Log.SetLogger("file", string(jsonConfig)) // 设置日志记录方式：本地文件记录
	Log.SetLevel(getLogLevel(logLevel))       // 设置日志写入缓冲区的等级
	Log.EnableFuncCallDepth(true)             // 输出log时能显示输出文件名和行号（非必须）
}

// 获取日志级别
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
