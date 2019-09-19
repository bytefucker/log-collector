package main

import (
	"github.com/astaxie/beego/logs"
)

var (
	configType = "ini"
	configPath string
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
	err = loadConfig(configType, configPath)
	if err != nil {
		logs.Error("Start logAgent [init loadConfig] failed, err: %s", err)
		return
	}
	logs.Debug("load Agent [config] success")

	// 初始化日志
	err = initAgentLog()
	if err != nil {
		logs.Error("Start logAgent [init agentLog] failed, err: %s", err)
		return
	}
	logs.Debug("Init Agent [log] success")

	// 初始化Etcd
	err = initEtcd(agentConfig.EtcdAddress, agentConfig.CollectKey)
	if err != nil {
		logs.Error("Start logAgent [init etcd] failed, err:", err)
		return
	}
	logs.Debug("Init Agent [etcd] success")

	// 初始化tailf
	err = initTailf(agentConfig.Collects, agentConfig.Chansize, agentConfig.Ip)
	if err != nil {
		logs.Error("Start logAgent [init tailf] failed, err:", err)
		return
	}
	logs.Debug("Init Agent [tailf] success")

	// 初始化kafka
	err = initKafka(agentConfig.KafkaAddress)
	if err != nil {
		logs.Error("Start logAgent [init kafka] failed, err:", err)
		return
	}
	logs.Debug("Init Agent [kafka] success")

	// 启动logagent服务
	err = serverRun()
	if err != nil {
		logs.Error("Start logAgent [init serverRun] failed, err:", err)
		return
	}
	logs.Info("Log Agent exit")
}
