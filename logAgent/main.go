package main

import (
	"github.com/astaxie/beego/logs"
	"logAgent/config"
	"logAgent/producer"
	"logAgent/server"
	"logAgent/task"
	"os"
)

var (
	configType  = "ini"
	configPath  string
	agentConfig *config.Config
)

func main() {
	var err error

	// 获取配置文件路径
	err = getConfigPath()
	if err != nil {
		logs.Error("%s", err)
		return
	}
	logs.Debug("get config file success, file: %s", configPath)

	// 加载配置文件
	err = config.LoadConfig(configType, configPath)
	if err != nil {
		logs.Error("Start logAgent [init loadConfig] failed, err: %s", err)
		return
	}
	logs.Debug("load Agent [config] success")

	// 初始化日志
	err = config.InitAgentLog()
	if err != nil {
		logs.Error("Start logAgent [init agentLog] failed, err: %s", err)
		return
	}
	logs.Debug("Init Agent [log] success")

	// 初始化Etcd
	err = config.InitEtcd(agentConfig.EtcdAddress, agentConfig.CollectKey)
	if err != nil {
		logs.Error("Start logAgent [init etcd] failed, err:", err)
		return
	}
	logs.Debug("Init Agent [etcd] success")

	// 初始化tailf
	err = task.InitTailfTask(agentConfig.Collects, agentConfig.Chansize, agentConfig.Ip)
	if err != nil {
		logs.Error("Start logAgent [init task] failed, err:", err)
		return
	}
	logs.Debug("Init Agent [task] success")

	// 初始化kafka
	err = producer.InitKafka(agentConfig.KafkaAddress)
	if err != nil {
		logs.Error("Start logAgent [init kafka] failed, err:", err)
		return
	}
	logs.Debug("Init Agent [kafka] success")

	// 启动logagent服务
	err = server.ServerRun()
	if err != nil {
		logs.Error("Start logAgent [init serverRun] failed, err:", err)
		return
	}
	logs.Info("Log Agent exit")
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
