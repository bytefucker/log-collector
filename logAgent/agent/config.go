package main

import (
	"errors"
	"github.com/astaxie/beego/logs"
	"logAgent/tailf"
	"os"
	"strings"

	"github.com/astaxie/beego/config"
)

// 存储logAgent配置信息
type Config struct {
	LogLevel     string
	LogPath      string
	Chansize     int
	KafkaAddress []string
	EtcdAddress  []string
	CollectKey   string
	Collects     []tailf.Collect
	Ip           string
}

var (
	// 配置信息对象
	agentConfig *Config
)

// 加载配置信息
func loadConfig(configType, configPath string) (err error) {
	conf, err := config.NewConfig(configType, configPath)
	if err != nil {
		return
	}

	agentConfig = &Config{}

	// 获取基础配置
	err = getAgentConfig(conf)
	if err != nil {
		return
	}

	return
}

func getAgentConfig(conf config.Configer) (err error) {
	// 获取日志级别
	logLevel := conf.String("base::log_level")
	if len(logLevel) == 0 {
		logLevel = "debug"
	}
	agentConfig.LogLevel = logLevel

	// 获取日志路径
	logPath := conf.String("base::log_path")
	if len(logPath) == 0 {
		logPath = "/Users/aery/Data/code/Go/go_old/logs/logagent.log"
	}
	agentConfig.LogPath = logPath

	// 日志收集开启chan大小
	chanSize, chanStatus := conf.Int("base::queue_size")
	if chanStatus != nil {
		chanSize = 200
	}
	agentConfig.Chansize = chanSize

	// etcd 地址
	etcdAddress := conf.String("etcd::etcd_address")
	if len(etcdAddress) == 0 {
		err = errors.New("Agent config etcd address error")
		return
	}
	agentConfig.EtcdAddress = strings.Split(etcdAddress, ",")

	// kafka 地址
	kafkaAddress := conf.String("kafka::kafka_address")
	if len(kafkaAddress) == 0 {
		err = errors.New("Agent config kafka address error")
		return
	}
	agentConfig.KafkaAddress = strings.Split(kafkaAddress, ",")

	// 获取日志收集前缀key
	collectKey := conf.String("collect::collectKey")
	if len(collectKey) == 0 {
		err = errors.New("Agent config collectKey error")
		return
	}
	agentConfig.CollectKey = collectKey

	return
}

// 通过传参的方式获取配置文件的路径
func getConfigPath() (err error) {
	cmdArgs := os.Args
	if len(cmdArgs) < 2 {
		logs.Warn("配置加载失败，默认加载config/logagent.ini")
		configPath = "config/logagent.ini"
	} else {
		configPath = cmdArgs[1]
	}
	return
}
