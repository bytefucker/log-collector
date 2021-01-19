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
	app.Commands = []cli.Command{
		{
			Name:        "agent",
			Description: "服务器收集模块",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "config",
					Usage: "配置文件路径",
					Value: "conf/config.yaml",
				},
			},
			Action: func(c *cli.Context) {
				path := c.String("config")
				agentConfig := config.NewAppConfig(path).Agent
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
			Action: func(c *cli.Context) {
				path := c.String("config")
				serverConfig := config.NewAppConfig(path).Manager
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
