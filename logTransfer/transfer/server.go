package main

import (
	"logTransfer/elasticsearch"

	"github.com/astaxie/beego/logs"
)

// 启动服务
func serverRun() (err error) {
	logs.Info("Log Transfer start is running...")
	for true {
		// 获取一条日志消息
		msg := elasticsearch.GetOneMsg()
		// 发送一条日志消息到elasticsearch
		err = elasticsearch.SendMsgToEs(msg, msg.Topic)
		if err != nil {
			logs.Error("send msg [%s] from topic: [%s] to es failed, err: %s", msg.Msg, msg.Topic, err)
			continue
		}
	}
	return
}
