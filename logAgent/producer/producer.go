package producer

import (
	"encoding/json"
	"github.com/astaxie/beego/logs"
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
		producer = KafkaProducer{}
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
type KafkaProducer struct{}

func (KafkaProducer) SendMsg(appKey string, msg model.LogContent) (err error) {
	panic("implement me")
}
