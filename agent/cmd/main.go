package main

import (
	"flag"
	"github.com/yihongzhi/logCollect/common/kafka"
	"github.com/yihongzhi/logCollect/common/logger"
	"strings"
)

var logs = logger.Instance

var (
	help       bool
	etcdAddrs  string
	kafkaAddrs string
	logLevel   string
)

func init() {
	flag.BoolVar(&help, "h", false, "帮助")
	flag.StringVar(&etcdAddrs, "etcd-address", "", "Etcd服务地址")
	flag.StringVar(&kafkaAddrs, "kafka-address", "", "Kafka服务地址")
	flag.StringVar(&logLevel, "log-level", "info", "Log级别")
}

func main() {
	flag.Parse()
	if help {
		flag.PrintDefaults()
	}
	logs.Infof("begin init etcd %s", etcdAddrs)
	//todo
	logs.Infof("begin init kafka %s", kafkaAddrs)
	_, err := kafka.InitKafkaClient(strings.Split(kafkaAddrs, ","))
	if err != nil {
		logs.Fatalf("init kafka %s failed!", kafkaAddrs)
	}
}
