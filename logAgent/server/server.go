package server

import (
	"logAgent/producer"
	"logAgent/task"
	"time"

	"github.com/astaxie/beego/logs"
)

var client producer.Producer

// 启动logagent服务
func ServerRun(producer producer.Producer) (err error) {
	logs.Info("Log Agent start is running...")
	client = producer
	for true {
		// 获取一行日志数据
		msg := task.GetOneLine()
		// 发送一行日志数据到kafka
		err = client.SendMsg(msg.Topic, msg.Msg)
		if err != nil {
			logs.Error("send msg:[%v] topic:[%v] failed, err:[%v]", msg.Msg, msg.Topic, err)
			time.Sleep(time.Second)
			continue
		}
	}
	return
}

// 发送数据到kafka
/*func sendToKafka(msg model.LogContent, topic string) (err error) {
	smsg, err := json.Marshal(&msg)
	if err != nil {
		logs.Error("send to kafka marshal failed --> msg: [%v], topic:[%s], error: %s", msg, topic, err)
		return
	}
	logs.Debug("send to kafka --> msg:[%v], topic:[%v]", string(smsg), topic)
	//err = producer.SendMsg(string(smsg), topic)
	return
}*/
