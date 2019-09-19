package elasticsearch

import (
	"encoding/json"

	"github.com/astaxie/beego/logs"
	elastic "gopkg.in/olivere/elastic.v2"
)

// elasticsearch 客户端
type EsClient struct {
	client *elastic.Client
}

// 发送到es消息存储结构
type SesMsg struct {
	Msg   string `json:"msg"`
	Ip    string `json:"ip"`
	Topic string `json:"topic"`
}

// kafka获取消息存储结构
type EsMsg struct {
	Msg   []byte
	Topic string
}

var (
	//初始化elasticsearch 客户端
	esClient *EsClient
	// es消息队列
	EsChan chan EsMsg
)

// 初始化elasticsearch
func InitElastic(esAddress []string, chanSize int) (err error) {
	client, err := elastic.NewClient(elastic.SetSniff(false), elastic.SetURL(esAddress...))
	if err != nil {
		return
	}
	esClient = &EsClient{
		client: client,
	}
	EsChan = make(chan EsMsg, chanSize)
	return
}

// 从eschan中获取一条数据
func GetOneMsg() (msg EsMsg) {
	msg = <-EsChan
	return
}

// 发送数据到elasticsearch
func SendMsgToEs(msg EsMsg, topic string) (err error) {
	logs.Debug("send msg: %s to topic: %s", msg, topic)
	kafkaMsg := SesMsg{}
	err = json.Unmarshal(msg.Msg, &kafkaMsg)
	if err != nil {
		logs.Error("send msg to es, ummarshal failed, msg: [%s], topic:[%s], err:[%s]",
			string(msg.Msg), topic, err)
		return
	}
	kafkaMsg.Topic = topic
	_, err = esClient.client.Index().
		Index(topic).
		Type(topic).
		//Id(interface{}...).
		BodyJson(kafkaMsg).
		Do()
	if err != nil {
		return
	}
	logs.Debug("send msg: %s to topic: %s success", msg, topic)
	return
}
