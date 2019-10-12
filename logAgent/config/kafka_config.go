package config

import (
	"fmt"
	"github.com/Shopify/sarama"
	"logAgent/model"
)

type KafkaClient struct {
	Client sarama.SyncProducer
}

// 初始化kafka生产者
func InitKafka(agentConfig *model.Config) (client *KafkaClient, err error) {
	conf := sarama.NewConfig()
	conf.Producer.RequiredAcks = sarama.WaitForAll
	conf.Producer.Partitioner = sarama.NewRandomPartitioner
	conf.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer(agentConfig.KafkaAddress, conf)
	if err != nil {
		fmt.Println("producer create failed,", err)
		return
	}
	client = &KafkaClient{
		Client: producer,
	}
	return
}
