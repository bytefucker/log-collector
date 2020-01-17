package kafka

import (
	"github.com/Shopify/sarama"
	"github.com/yihongzhi/logCollect/common/logger"
)

var log = logger.Instance

//kafka客户端
type KafkaClient struct {
	Client sarama.SyncProducer
}

// 初始化kafka生产者
func InitKafkaClient(addrs []string) (client KafkaClient, err error) {
	conf := sarama.NewConfig()
	conf.Producer.RequiredAcks = sarama.WaitForAll
	conf.Producer.Partitioner = sarama.NewRandomPartitioner
	conf.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer(addrs, conf)
	if err != nil {
		log.Error("SyncProduce create failed !", err)
		return
	}
	client = KafkaClient{
		Client: producer,
	}
	return
}
