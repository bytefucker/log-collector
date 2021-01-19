package config

import (
	"github.com/spf13/viper"
	"github.com/yihongzhi/log-collector/common/logger"
)

const CollectorKey = "log-collector"

var (
	log = logger.Instance
)

type AppConfig struct {
	Debug    bool            `mapstructure:debug`
	Agent    *AgentConfig    `mapstructure:"agent"`
	Analyzer *AnalyzerConfig `mapstructure:"analyzer"`
	Manager  *ManagerConfig  `mapstructure:"manager"`
}

type AgentConfig struct {
	EtcdAddr     []string `mapstructure:"etcd-addr"`
	KafKaAddr    []string `mapstructure:"kafka-addr"`
	ChanSize     int      `mapstructure:"chan-size"`
	CollectorKey string   `mapstructure:"collector-key"`
}

type AnalyzerConfig struct {
}

type ManagerConfig struct {
	Port         int      `mapstructure:"port"`
	EtcdAdrr     []string `mapstructure:"etcd-addr"`
	DBConnectStr string   `mapstructure:"db-connect-str"`
	CollectorKey string   `mapstructure:"collector-key"`
}

//读取yaml配置
func NewAppConfig(path string) *AppConfig {
	config := viper.New()
	config.SetConfigFile(path)
	//尝试进行配置读取
	if err := config.ReadInConfig(); err != nil {
		log.Fatalln("读取配置失败", err)
	}
	appConfig := &AppConfig{
		Debug:    false,
		Agent:    &AgentConfig{},
		Analyzer: &AnalyzerConfig{},
		Manager:  &ManagerConfig{},
	}
	if err := config.Unmarshal(appConfig); err != nil {
		log.Fatalln("转换配置失败", err)
	}
	return appConfig
}
