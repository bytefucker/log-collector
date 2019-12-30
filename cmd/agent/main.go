package main

import (
	"github.com/MarciaYi/logCollect/common/config"
	"github.com/MarciaYi/logCollect/internal/agent/model"
	"github.com/MarciaYi/logCollect/internal/agent/producer"
	"github.com/MarciaYi/logCollect/internal/agent/server"
	"github.com/MarciaYi/logCollect/internal/agent/task"
	"github.com/astaxie/beego/logs"
	"os"
)

var (
	configType  = "ini"
	configPath  string
	agentConfig *model.Config
)

func main() {
	var err error

	// 获取配置文件路径
	err = getConfigPath()
	if err != nil {
		logs.Error("获取配置文件路径失败，%s", err)
		return
	}

	// 加载配置文件
	agentConfig, err := config.LoadConfig(configType, configPath)
	if err != nil {
		logs.Error("配置文件加载失败, %s", err)
		return
	}

	// 初始化日志
	err = config.InitAgentLog(agentConfig)
	if err != nil {
		logs.Error("初始化日志失败,%s", err)
		return
	}

	// 初始化Etcd
	err = config.InitEtcd(agentConfig)
	if err != nil {
		logs.Error("初始化etcd失败，%s", err)
		return
	}

	//初始化producer
	pr, err := producer.InitProducer(agentConfig)
	if err != nil {
		logs.Error("初始化producer失败", err)
		return
	}

	// 初始化task
	err = task.InitTailfTask(agentConfig)
	if err != nil {
		logs.Error("初始化tailf tasks失败", err)
		return
	}

	// 启动logagent服务
	err = server.ServerRun(pr)
	if err != nil {
		logs.Error("启动logagent服务失败", err)
		return
	}
}

// 通过传参的方式获取配置文件的路径
func getConfigPath() (err error) {
	cmdArgs := os.Args
	if len(cmdArgs) < 2 {
		logs.Warn("配置加载失败，默认加载config/logagent.ini")
		configPath = "conf/logagent.ini"
	} else {
		configPath = cmdArgs[1]
	}
	return
}
