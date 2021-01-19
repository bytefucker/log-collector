package manager

import (
	"github.com/gin-gonic/gin"
	"github.com/yihongzhi/log-collector/common/etcd"
	"github.com/yihongzhi/log-collector/common/logger"
	"github.com/yihongzhi/log-collector/config"
)

var (
	log       = logger.Instance
	etcdAddrs []string
	port      int
	debug     bool
)

type manageServer struct {
	etcd *etcd.EtcdClient
	gin  *gin.Engine
}

func NewManageServer(c *config.ManagerConfig) (*manageServer, error) {
	var err error
	etcdClient, err := etcd.NewClient(c.EtcdAdrr)
	if err != nil {
		log.Fatalf("init etcd client %s failed...", c.EtcdAdrr)
	}
	gin.SetMode(gin.ReleaseMode)
	engine := gin.Default()
	server := &manageServer{
		etcd: etcdClient,
		gin:  engine,
	}
	return server, err
}

func (server *manageServer) StartManageServer() error {
	return server.gin.Run("")
}
