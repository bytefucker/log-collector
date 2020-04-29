package manager

import (
	"github.com/gin-gonic/gin"
	"github.com/urfave/cli"
	"github.com/yihongzhi/logCollect/common/etcd"
	"github.com/yihongzhi/logCollect/common/logger"
	"github.com/yihongzhi/logCollect/manager/database"
	"github.com/yihongzhi/logCollect/manager/routers"
)

var (
	log       = logger.Instance
	etcdAddrs []string
	port      int
	debug     bool
)

func StartManageServer(c *cli.Context) error {
	var err error
	//初始化数据库
	err = database.Open("root:y@4216160@/manager?charset=utf8&parseTime=True&loc=Local")
	//初始化etcdc
	etcdAddrs = c.StringSlice("etcd-addr")
	etcd.NewClient(etcdAddrs)
	//初始化Web服务
	port = c.Int("port")
	gin.SetMode(gin.ReleaseMode)
	app := gin.Default()
	api := app.Group("/api")
	{
		api.GET("/application/list", routers.ApplicationList)
	}
	err = app.Run(":8080")
	return err
}
