package producer

import (
	"encoding/json"
	"github.com/Shopify/sarama"
	"github.com/yihongzhi/logCollect/common/kafka"
	"github.com/yihongzhi/logCollect/common/logger"
)

type LogContent struct {
	Msg string `json:"msg"` //日志内容
	Ip  string `json:"ip"`  //机器IP
}

var logs = logger.Instance

//消费者接口
type Producer interface {
	//发送消息
	SendMsg(appKey string, msg LogContent) error
}

//Kafka消费者
type KafkaProducer struct {
	kafka.KafkaClient
}

func (producer KafkaProducer) SendMsg(appKey string, msg LogContent) (err error) {
	text, err := json.Marshal(&msg)
	if err != nil {
		logs.Debug("serilazid msg failed ", err)
		return
	}
	_, _, err = producer.Client.SendMessage(&sarama.ProducerMessage{Topic: appKey, Key: sarama.StringEncoder(appKey), Value: sarama.StringEncoder(text)})
	if err != nil {
		logs.Error("kafka producer msg failed", err)
		return
	}
	logs.Debugf("send to kafka -->appKey:[%v],msg:[%v]", appKey, string(text))
	return
}

//初始化producer
func InitKafkaProducer(addrs []string) (producer Producer, err error) {
	kafkaClient, err := kafka.InitKafkaClient(addrs)
	if err != nil {
		logs.Error("int kafka failed ", err)
	}
	producer = KafkaProducer{
		KafkaClient: kafkaClient,
	}
	return
}
