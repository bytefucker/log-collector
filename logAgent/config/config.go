package config

import (
	"errors"
	"logAgent/model"
	"strings"

	"github.com/astaxie/beego/config"
)

// 存储logAgent配置信息
type Config struct {
	LogLevel     string
	LogPath      string
	ChanSize     int
	SendModel    string //发送方式
	KafkaAddress []string
	EtcdAddress  []string
	CollectKey   string
	CollectTasks []model.CollectTask
	Ip           string
}

// 加载配置信息
func LoadConfig(configType, configPath string) (agentConfig *Config, err error) {
	conf, err := config.NewConfig(configType, configPath)
	if err != nil {
		return
	}
	agentConfig = &Config{}
	// 获取基础配置
	err = parseAgentConfig(agentConfig, conf)
	if err != nil {
		return
	}
	return agentConfig, nil
}

//解析配置文件
func parseAgentConfig(agentConfig *Config, conf config.Configer) (err error) {
	// 获取日志级别
	logLevel := conf.String("base::log_level")
	if len(logLevel) == 0 {
		logLevel = "debug"
	}
	agentConfig.LogLevel = logLevel

	// 获取日志路径
	logPath := conf.String("base::log_path")
	if len(logPath) == 0 {
		logPath = "logs/logagent.log"
	}
	agentConfig.LogPath = logPath

	sendModel := conf.String("base::send_model")
	if len(sendModel) == 0 {
		err = errors.New("Agent config sendModel error")
		return
	}
	agentConfig.SendModel = sendModel

	// 日志收集开启chan大小
	chanSize, chanStatus := conf.Int("base::queue_size")
	if chanStatus != nil {
		chanSize = 200
	}
	agentConfig.ChanSize = chanSize

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
