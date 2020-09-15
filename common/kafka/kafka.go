package kafka

import (
	"encoding/json"
	"github.com/Shopify/sarama"
	"github.com/yihongzhi/logCollect/common/logger"
	"time"
)

var log = logger.Instance

type LogContent struct {
	Msg string `json:"msg"` //日志内容
	Ip  string `json:"ip"`  //机器IP
}

//kafka客户端
type KafkaClient struct {
	Client sarama.SyncProducer
}

// 初始化kafka生产者
func NewKafkaClient(addrs []string) (client *KafkaClient, err error) {
	conf := sarama.NewConfig()
	conf.Producer.RequiredAcks = sarama.WaitForAll
	conf.Producer.Partitioner = sarama.NewRandomPartitioner
	conf.Producer.Return.Successes = true
	conf.Producer.Timeout = 5 * time.Second

	producer, err := sarama.NewSyncProducer(addrs, conf)
	if err != nil {
		log.Error("SyncProduce create failed !", err)
		return
	}
	client = &KafkaClient{
		Client: producer,
	}
	return
}

func (kafka *KafkaClient) SendMsg(appKey string, msg LogContent) (err error) {
	text, err := json.Marshal(&msg)
	if err != nil {
		log.Debug("serilazid msg failed ", err)
		return
	}
	_, _, err = kafka.Client.SendMessage(&sarama.ProducerMessage{Topic: appKey, Key: sarama.StringEncoder(appKey), Value: sarama.StringEncoder(text)})
	if err != nil {
		log.Error("kafka kafka msg failed", err)
		return
	}
	log.Debugf("send to kafka -->appKey:[%v],msg:[%v]", appKey, string(text))
	return
}
