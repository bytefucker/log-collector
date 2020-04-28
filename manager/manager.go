package manager

import (
	"github.com/gin-gonic/gin"
	"github.com/urfave/cli"
	"github.com/yihongzhi/logCollect/common/logger"
	"github.com/yihongzhi/logCollect/manager/routers"
)

var (
	log = logger.Instance
)

func NewManageServer(c *cli.Context) {
	port := c.String("port")
	app := gin.Default()
	api := app.Group("/api")
	{
		api.GET("/demo", routers.Demo)
	}
	err := app.Run(":" + port)
	log.Fatalf("start server fail... ", err)
}
