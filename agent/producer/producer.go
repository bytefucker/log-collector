package producer

import (
	"encoding/json"
	"github.com/Shopify/sarama"
	"github.com/yihongzhi/logCollect/agent/model"
	"github.com/yihongzhi/logCollect/common/kafka"
	"github.com/yihongzhi/logCollect/common/logger"
)

var logs = logger.Instance

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
		kafkaClient, err := kafka.InitKafkaClient(agentConfig.KafkaAddress)
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
		logs.Errorf("send to http marshal failed --> msg: [%v], appKey:[%s], error: %s", msg, topic, err)
		return
	}
	logs.Debugf("send to http -->appKey:[%v],msg:[%v]", topic, string(json))
	return
}

//Kafka消费者
type KafkaProducer struct {
	kafka.KafkaClient
}

func (producer KafkaProducer) SendMsg(appKey string, msg model.LogContent) (err error) {
	text, err := json.Marshal(&msg)
	if err != nil {
		logs.Debug("序列化kafka消息失败", err)
		return
	}
	_, _, err = producer.Client.SendMessage(&sarama.ProducerMessage{Topic: appKey, Value: sarama.StringEncoder(text)})
	if err != nil {
		logs.Debug("kafka消费消息失败", err)
		return
	}
	logs.Debugf("send to http -->appKey:[%v],msg:[%v]", appKey, string(text))
	return
}
