package main

import (
	"encoding/json"

	"github.com/astaxie/beego"

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
func initLog() (err error) {
	logLevel := beego.AppConfig.String("loglevel")
	logFile := beego.AppConfig.String("logfile")
	logConfig := make(map[string]interface{})
	logConfig["filename"] = logFile
	logConfig["level"] = caseLogLevel(logLevel)
	logConfig["color"] = true

	logConfigString, err := json.Marshal(logConfig)
	if err != nil {
		return
	}

	beego.BConfig.Log.AccessLogs = true
	beego.BConfig.Log.FileLineNum = true

	// beego.BConfig.Log.Outputs["file"] = string(logConfigString)
	//beego.BConfig.Log.Outputs["console"] = string(logConfigString)

	logs.SetLogger(logs.AdapterConsole, string(logConfigString))
	logs.SetLogger(logs.AdapterFile, string(logConfigString))
	///logs.EnableFuncCallDepth(true)

	return
}
