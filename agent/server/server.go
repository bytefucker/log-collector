package server

import (
	"github.com/yihongzhi/logCollect/agent/producer"
	"github.com/yihongzhi/logCollect/agent/task"
	"github.com/yihongzhi/logCollect/common/logger"
	"time"
)

var (
	pro  producer.Producer
	logs = logger.Instance
)

// 启动logagent服务
func ServerRun(producer producer.Producer) (err error) {
	logs.Info("Log Agent start is running...")
	pro = producer
	for true {
		// 获取一行日志数据
		msg := task.GetOneLine()
		// 发送一行日志数据到kafka
		err = pro.SendMsg(msg.AppKey, msg.Msg)
		if err != nil {
			logs.Errorf("send msg:[%v] topic:[%v] failed, err:[%v]", msg.Msg, msg.AppKey, err)
			time.Sleep(time.Second)
			continue
		}
	}
	return
}
