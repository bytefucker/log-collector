package main

import (
	"logTransfer/elasticsearch"
	"logTransfer/kafka"

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
		logs.Error("get config file failed, err: %s", err)
		return
	}
	logs.Debug("get config file success, file: %s", configPath)

	// 加载配置文件
	err = loadConfig(configType, configPath)
	if err != nil {
		logs.Error("load config failed, err:%s", err)
		return
	}
	logs.Debug("load log transfer success")

	// 初始化日志
	err = initTransgerLog()
	if err != nil {
		logs.Error("init log failed, err:%s", err)
		return
	}
	logs.Debug("init transfer log success")

	// 初始化etcd
	err = initEtcd(transferConfig.EtcdAddress)
	if err != nil {
		logs.Error("init etcd failed, err:%s", err)
		return
	}
	logs.Debug("init transfer etcd success")

	// 初始化elasticsearch
	err = elasticsearch.InitElastic(transferConfig.EsAddress, transferConfig.Chansize)
	if err != nil {
		logs.Error("init elasticsearch failed, err: %s", err)
		return
	}
	logs.Debug("int transger elasticsearch success")

	// 初始化kafka
	err = kafka.InitKafka(transferConfig.KafkaAddress, transferConfig.Topics)
	if err != nil {
		logs.Error("init kafka failed, err: %s", err)
		return
	}
	logs.Debug("init transger kafka success")

	// 启动服务
	err = serverRun()
	if err != nil {
		logs.Error("server run failed, err: %s", err)
		return
	}
	logs.Info("Log Transger exited")

}
