package main

import (
	"encoding/json"

	"github.com/astaxie/beego/logs"
)

// 获取日志级别
func caseLogLevel(level string) (logLevel int) {
	switch level {
	case "debug":
		logLevel = logs.LevelDebug
	case "info":
		logLevel = logs.LevelInfo
	case "warn":
		logLevel = logs.LevelWarn
	case "error":
		logLevel = logs.LevelError
	default:
		logLevel = logs.LevelInfo
	}
	return logLevel
}

// 初始化日志
func initTransgerLog() (err error) {
	logConfig := make(map[string]interface{})
	logConfig["filename"] = transferConfig.LogPath
	logConfig["level"] = caseLogLevel(transferConfig.LogLevel)
	logConfig["color"] = true

	logConfigString, err := json.Marshal(logConfig)
	if err != nil {
		return
	}

	logs.SetLogger(logs.AdapterConsole, string(logConfigString))
	logs.SetLogger(logs.AdapterFile, string(logConfigString))
	//logs.EnableFuncCallDepth(true)

	return
}
