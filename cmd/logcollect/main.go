package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"os"
)

func main() {
	app := cli.NewApp()
	app.Name = "分布式日志收集"
	app.Version = "1.0.0"
	app.Flags = []cli.Flag{
		cli.StringFlag{Name: "ip"},
	}
	app.Commands = []cli.Command{
		cli.Command{
			Name: "agent",
			Action: func(c *cli.Context) {
				fmt.Println("agent")
			},
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		logrus.Fatal(err)
	}
}
