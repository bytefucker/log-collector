package main

import (
	"flag"
	"github.com/yihongzhi/logCollect/agent/producer"
	"github.com/yihongzhi/logCollect/agent/server"
	"github.com/yihongzhi/logCollect/agent/task"
	"github.com/yihongzhi/logCollect/common/etcd"
	"github.com/yihongzhi/logCollect/common/logger"
	"strings"
)

var logs = logger.Instance

var (
	help       bool
	etcdAddrs  string
	kafkaAddrs string
	logLevel   string
	collectKey string
	chanSize   int
)

func init() {
	flag.BoolVar(&help, "h", false, "帮助")
	flag.StringVar(&etcdAddrs, "etcd-address", "", "Etcd服务地址")
	flag.StringVar(&kafkaAddrs, "kafka-address", "", "Kafka服务地址")
	flag.StringVar(&logLevel, "log-level", "info", "Log级别")
	flag.StringVar(&collectKey, "collect-key", "/logagent", "Log级别")
	flag.IntVar(&chanSize, "chan-size", 1021, "日志通道大小")
}

func main() {
	flagParse()
	var err error
	//1.初始化etcd
	etcdClient, err := etcd.InitEtcdClient(strings.Split(etcdAddrs, ","))
	if err != nil {
		logs.Fatalf("init etcdclient %s failed...", etcdAddrs)
	}
	//2.初始化producer
	kafkaProducer, err := producer.InitKafkaProducer(strings.Split(kafkaAddrs, ","))
	if err != nil {
		logs.Fatalf("init kafka producer %s failed...", kafkaAddrs)
	}
	//3.初始化任务
	err = task.InitTailfTask(collectKey, chanSize, etcdClient)
	if err != nil {
		logs.Fatal("init task failed ...", err)
	}
	//4.初始化收集服务
	err = server.ServerRun(kafkaProducer)
	if err != nil {
		logs.Fatal("log agent server start failed ...", err)
	}
}

func flagParse() {
	flag.Parse()
	logger.SetLevel(logLevel)
	if help {
		flag.PrintDefaults()
	}
}
