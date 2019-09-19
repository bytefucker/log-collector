package main

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/astaxie/beego/config"
)

// 配置文件结构体
type Config struct {
	LogLevel     string
	LogPath      string
	Chansize     int
	KafkaAddress []string
	EtcdAddress  []string
	EsAddress    []string
	Topics       []string
}

var (
	// 配置文件对象
	transferConfig *Config
)

// 加载配置文件配置
func loadConfig(configType, configPath string) (err error) {

	conf, err := config.NewConfig(configType, configPath)
	if err != nil {
		return
	}

	transferConfig = &Config{}

	// 获取基础配置
	err = getTransgerConfig(conf)
	if err != nil {
		return
	}
	return
}

func getTransgerConfig(conf config.Configer) (err error) {
	// 获取日志级别
	logLevel := conf.String("base::log_level")
	if len(logLevel) == 0 {
		logLevel = "debug"
	}
	transferConfig.LogLevel = logLevel

	// 获取日志路径
	logPath := conf.String("base::log_path")
	if len(logPath) == 0 {
		logPath = "/Users/aery/Data/code/Go/go_old/logs/logagent.log"
	}
	transferConfig.LogPath = logPath

	// 日志收集开启chan大小
	chanSize, chanStatus := conf.Int("base::queue_size")
	if chanStatus != nil {
		chanSize = 200
	}
	transferConfig.Chansize = chanSize

	// etcd 地址
	etcdAddress := conf.String("etcd::etcd_address")
	if len(etcdAddress) == 0 {
		err = errors.New("Transger config etcd address error")
		return
	}
	transferConfig.EtcdAddress = strings.Split(etcdAddress, ",")

	// kafka 地址
	kafkaAddress := conf.String("kafka::kafka_address")
	if len(kafkaAddress) == 0 {
		err = errors.New("Transger config kafka address error")
		return
	}
	transferConfig.KafkaAddress = strings.Split(kafkaAddress, ",")

	// elasticsearch 地址
	esAddress := conf.String("elasticsearch::es_address")
	if len(kafkaAddress) == 0 {
		err = errors.New("Transger config elasticsearch address error")
		return
	}
	transferConfig.EsAddress = strings.Split(esAddress, ",")

	return
}

// 根据传参的方式获取配置文件路径
func getConfigPath() (err error) {
	cmdArgs := os.Args
	if len(cmdArgs) < 2 {
		err = fmt.Errorf("USAGE: %v  <agent config file>, go to start the log agent.", cmdArgs[0])
		return
	}
	configPath = cmdArgs[1]
	return
}
