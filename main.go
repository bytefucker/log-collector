package main

import (
	"fmt"
	"github.com/yihongzhi/log-collector/common/logger"
	"github.com/yihongzhi/log-collector/config"
	"github.com/yihongzhi/log-collector/manager"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"github.com/yihongzhi/log-collector/agent"
)

var log = logger.Instance

func main() {
	app := cli.NewApp()
	app.Name = "log-collector"
	app.Description = "分布式日志收集组件"
	app.Version = "1.0.0"
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name: "debug",
		},
	}
	app.Commands = []cli.Command{
		{
			Name:        "agent",
			Description: "服务器收集模块",
			Flags: []cli.Flag{
				cli.StringSliceFlag{
					Name: "etcd-addr",
				},
				cli.StringSliceFlag{
					Name: "kafka-addr",
				},
				cli.IntFlag{
					Name:  "chan-size",
					Value: 1024,
				},
			},
			Action: func(c *cli.Context) {
				agentConfig := config.InitAgentConfig(c)
				logAgent, err := agent.NewAgent(agentConfig)
				if err != nil {
					log.Error("Init Agent Failed...", err)
				}
				logAgent.StartAgent()
			},
		},
		{
			Name:        "analysis",
			Description: "日志解析模块",
			Action: func(c *cli.Context) {
				fmt.Println("analysis")
			},
		},
		{
			Name:        "manager",
			Description: "日志查询管理模块",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "port",
					Value: "8080",
					Usage: "API监听地址",
				},
				cli.StringSliceFlag{
					Name:     "etcd-addr",
					Required: true,
					Usage:    "Etcd服务地址",
				},
				cli.StringSliceFlag{
					Name:     "db-connect-string",
					Required: true,
					Usage:    "数据库连接字符串",
				},
			},
			Action: func(c *cli.Context) {
				serverConfig := config.InitManagerServerConfig(c)
				server, err := manager.NewManageServer(serverConfig)
				if err != nil {
					log.Error("Init ManageServer Failed...", err)
				}
				server.StartManageServer()
			},
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		logrus.Fatal(err)
	}
}
