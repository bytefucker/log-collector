package manager

import (
	"github.com/gin-gonic/gin"
	"github.com/urfave/cli"
	"github.com/yihongzhi/logCollect/common/etcd"
	"github.com/yihongzhi/logCollect/common/logger"
	"github.com/yihongzhi/logCollect/manager/db"
	"github.com/yihongzhi/logCollect/manager/routers"
)

var (
	log       = logger.Instance
	etcdAddrs []string
	port      int
	debug     bool
)

func NewManageServer(c *cli.Context) error {
	var err error
	debug = c.Bool("debug")
	if debug {
		gin.SetMode(gin.DebugMode)
	}
	//初始化数据库
	err = db.Open()
	//初始化etcdc
	etcdAddrs = c.StringSlice("etcd-addr")
	etcd.NewClient(etcdAddrs)
	//初始化Web服务
	port = c.Int("port")
	app := gin.Default()
	api := app.Group("/api")
	{
		api.GET("/demo", routers.Demo)
	}
	err = app.Run(":" + string(port))
	return err
}
