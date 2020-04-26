package agent

import (
	"github.com/urfave/cli"
	"github.com/yihongzhi/logCollect/agent/producer"
	"github.com/yihongzhi/logCollect/agent/server"
	"github.com/yihongzhi/logCollect/agent/task"
	"github.com/yihongzhi/logCollect/common/etcd"
	"github.com/yihongzhi/logCollect/common/logger"
)

var (
	log        = logger.Instance
	etcdAddrs  []string
	kafkaAddrs []string
	collectKey string
	chanSize   int
)

func initArgs(cli *cli.Context) {
	etcdAddrs = cli.StringSlice("etcd-addr")
	kafkaAddrs = cli.StringSlice("kafka-addr")
	collectKey = cli.String("collect-key")
	chanSize = cli.Int("chan-size")
	logger.EnableDebugLevel(cli.Bool("debug"))
	log.Infof("initArgs etcd-addr=%s,kafka-addr=%s,collect-key=%s,chan-size=%d", etcdAddrs, kafkaAddrs, collectKey, chanSize)
}

//开启一个收集代理
func StartAgent(cli *cli.Context) error {
	initArgs(cli)
	var err error
	//1.初始化etcd
	etcdClient, err := etcd.NewClient(etcdAddrs)
	if err != nil {
		log.Fatalf("init etcdclient %s failed...", etcdAddrs)
	}
	//2.初始化producer
	kafkaProducer, err := producer.InitKafkaProducer(kafkaAddrs)
	if err != nil {
		log.Fatalf("init kafka producer %s failed...", kafkaAddrs)
	}
	//3.初始化任务
	err = task.InitTailfTask(collectKey, chanSize, etcdClient)
	if err != nil {
		log.Fatal("init task failed ...", err)
	}
	//4.初始化收集服务
	err = server.ServerRun(kafkaProducer)
	if err != nil {
		log.Fatal("log agent server start failed ...", err)
	}
	return err
}
