package main

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "logcollect"
	app.Description = "分布式日志收集组件"
	app.Version = "1.0.0"
	app.Commands = []cli.Command{
		{
			Name:        "agent",
			Description: "服务器收集模块",
			Action: func(c *cli.Context) {
				fmt.Println("agent")
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
