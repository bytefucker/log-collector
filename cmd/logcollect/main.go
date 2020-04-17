package main

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"github.com/yihongzhi/logCollect/agent"
)

func main() {
	app := cli.NewApp()
	app.Name = "logcollect"
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
				cli.StringFlag{
					Name:  "collect-key",
					Value: "/logagent",
				},
				cli.IntFlag{
					Name:  "chan-size",
					Value: 1024,
				},
			},
			Action: func(c *cli.Context) {
				err := agent.StartAgent(c)
				logrus.Error("StartAgent Failed...", err)
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
				fmt.Println("manager")
			},
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		logrus.Fatal(err)
	}
}
