package producer

import (
	"encoding/json"
	"github.com/Shopify/sarama"
	"github.com/astaxie/beego/logs"
	"logAgent/config"
	"logAgent/model"
)

//消费者接口
type Producer interface {
	//发送消息
	SendMsg(appKey string, msg model.LogContent) error
}

//初始化producer
func InitProducer(agentConfig *model.Config) (producer Producer, err error) {
	switch agentConfig.SendModel {
	case "http":
		producer = HttpProducer{}
	case "kafka":
		kafkaClient, err := config.InitKafka(agentConfig)
		if err != nil {
			logs.Error("初始化kafka失败", err)
		}
		producer = KafkaProducer{
			KafkaClient: kafkaClient,
		}
	}
	return
}

//Http消费者
type HttpProducer struct{}

func (HttpProducer) SendMsg(topic string, msg model.LogContent) (err error) {
	json, err := json.Marshal(&msg)
	if err != nil {
		logs.Error("send to http marshal failed --> msg: [%v], appKey:[%s], error: %s", msg, topic, err)
		return
	}
	logs.Debug("send to http -->appKey:[%v],msg:[%v]", topic, string(json))
	return
}

//Kafka消费者
type KafkaProducer struct {
	config.KafkaClient
}

func (producer KafkaProducer) SendMsg(appKey string, msg model.LogContent) (err error) {
	text, err := json.Marshal(&msg)
	if err != nil {
		logs.Error("序列化kafka消息失败", err)
		return
	}
	_, _, err = producer.Client.SendMessage(&sarama.ProducerMessage{Topic: appKey, Value: sarama.StringEncoder(text)})
	if err != nil {
		logs.Error("kafka消费消息失败", err)
		return
	}
	logs.Debug("send to http -->appKey:[%v],msg:[%v]", appKey, string(text))
	return
}
