package config

import "github.com/urfave/cli"

const CollectorKey = "log-collector"

type systemConfig struct {
	Debug bool
}

type AgentConfig struct {
	EtcdAdrr     []string
	KafKaAddr    []string
	ChanSize     int
	CollectorKey string
}

type ManagerServerConfig struct {
	Port         int
	EtcdAdrr     []string
	CollectorKey string
}

//初始化配置
func InitAgentConfig(c *cli.Context) *AgentConfig {
	return &AgentConfig{
		EtcdAdrr:     c.StringSlice("etcd-addr"),
		KafKaAddr:    c.StringSlice("kafka-addr"),
		ChanSize:     c.Int("chan-size"),
		CollectorKey: CollectorKey,
	}
}

func InitManagerServerConfig(c *cli.Context) *ManagerServerConfig {
	return &ManagerServerConfig{
		Port:     c.Int("port"),
		EtcdAdrr: c.StringSlice("etcd-addr"),
	}
}
