package server

import (
	"github.com/yihongzhi/logCollect/agent/producer"
	"github.com/yihongzhi/logCollect/agent/task"
	"github.com/yihongzhi/logCollect/common/logger"
	"github.com/yihongzhi/logCollect/common/utils"
	"time"
)

var (
	logs    = logger.Instance
	localIp string
)

// 启动logagent服务
func ServerRun(pro producer.Producer) (err error) {
	localIp = utils.LocalIpArray[0]
	logs.Infof("logagent start is running bind ip %s....", localIp)
	for true {
		msg := task.GetOneLine()
		err = pro.SendMsg(msg.AppKey, producer.LogContent{Ip: localIp, Msg: msg.Msg})
		if err != nil {
			logs.Errorf("send msg:[%v] topic:[%v] failed, err:[%v]", msg.Msg, msg.AppKey, err)
			time.Sleep(time.Second)
			continue
		}
	}
	return
}
