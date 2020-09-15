package agent

import (
	"github.com/yihongzhi/logCollect/agent/task"
	"github.com/yihongzhi/logCollect/common/etcd"
	"github.com/yihongzhi/logCollect/common/kafka"
	"github.com/yihongzhi/logCollect/common/logger"
	"github.com/yihongzhi/logCollect/common/utils"
	"github.com/yihongzhi/logCollect/config"
	"time"
)

var (
	log = logger.Instance
)

//日志收集代理
type logAgent struct {
	etcdClient  *etcd.EtcdClient
	taskMgr     *task.TaskManger
	kafkaClient *kafka.KafkaClient
}

//开启一个收集代理
func NewAgent(config *config.AgentConfig) (*logAgent, error) {
	var err error
	//1.初始化etcd
	etcdClient, err := etcd.NewClient(config.EtcdAdrr)
	if err != nil {
		log.Fatalf("init etcdclient %s failed...", config.EtcdAdrr)
	}
	//2.初始化producer
	kafkaClient, err := kafka.NewKafkaClient(config.KafKaAddr)
	if err != nil {
		log.Fatalf("init kafka producer %s failed...", config.KafKaAddr)
	}
	//3.初始化任务
	taskMgr, err := task.NewTaskManger(config.CollectorKey, config.ChanSize, etcdClient)
	if err != nil {
		log.Fatal("init task failed ...", err)
	}
	agent := &logAgent{
		etcdClient:  etcdClient,
		kafkaClient: kafkaClient,
		taskMgr:     taskMgr,
	}
	return agent, err
}

func (agent *logAgent) StartAgent() {
	localIp := utils.LocalIpArray[0]
	log.Infof("agent start is running bind ip %s....", localIp)
	for true {
		msg := task.GetOneLine()
		err := agent.kafkaClient.SendMsg(msg.AppKey, kafka.LogContent{Ip: localIp, Msg: msg.Msg})
		if err != nil {
			log.Errorf("send msg:[%v] topic:[%v] failed, err:[%v]", msg.Msg, msg.AppKey, err)
			time.Sleep(time.Second)
			continue
		}
	}
	return
}
