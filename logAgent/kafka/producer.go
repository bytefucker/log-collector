package kafka

import (
	"fmt"

	"github.com/Shopify/sarama"
	"github.com/astaxie/beego/logs"
)

var (
	// kafka生产者对象
	kafkaServerClient sarama.SyncProducer
)

// 初始化kafka生产者
func InitKafka(address []string) (err error) {
	kafkaServerConf := sarama.NewConfig()
	kafkaServerConf.Producer.RequiredAcks = sarama.WaitForAll
	kafkaServerConf.Producer.Partitioner = sarama.NewRandomPartitioner
	kafkaServerConf.Producer.Return.Successes = true

	kafkaServerClient, err = sarama.NewSyncProducer(address, kafkaServerConf)
	if err != nil {
		fmt.Println("producer create failed,", err)
		return
	}

	return
}

// 发送消息到kafka
func SendMsgToKafka(msg, topic string) (err error) {
	msgobj := &sarama.ProducerMessage{}
	msgobj.Topic = topic
	msgobj.Value = sarama.StringEncoder(msg)

	pid, offset, err := kafkaServerClient.SendMessage(msgobj)

	if err != nil {
		logs.Error("send msg to kafka topic:[%v] msg:[%v] faile, %v", msg, topic, err)
		return
	}

	logs.Debug("topic: [%v] pid: [%v], offset: [%v]", topic, pid, offset)
	return
}
