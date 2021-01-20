package agent

import (
	"github.com/yihongzhi/log-collector/agent/task"
	"github.com/yihongzhi/log-collector/common/etcd"
	"github.com/yihongzhi/log-collector/common/kafka"
	"github.com/yihongzhi/log-collector/common/logger"
	"github.com/yihongzhi/log-collector/config"
	"time"
)

var (
	log = logger.Instance
)

//日志收集代理
type logAgent struct {
	etcdClient     *etcd.EtcdClient
	tailTaskManger *task.TailTaskManger
	kafkaClient    *kafka.KafkaClient
}

//开启一个收集代理
func NewAgent(c *config.AgentConfig) (*logAgent, error) {
	var err error
	//1.初始化etcd
	etcdClient, err := etcd.NewClient(c.EtcdAddr)
	if err != nil {
		log.Fatalf("init etcd client %s failed...", c.EtcdAddr)
	}
	//2.初始化producer
	kafkaClient, err := kafka.NewKafkaClient(c.KafKaAddr)
	if err != nil {
		log.Fatalf("init kafka producer %s failed...", c.KafKaAddr)
	}
	//3.初始化任务
	taskMgr, err := task.NewTailTaskManger(c.CollectorKey, c.ChanSize, etcdClient)
	if err != nil {
		log.Fatal("init task failed ...", err)
	}
	agent := &logAgent{
		etcdClient:     etcdClient,
		kafkaClient:    kafkaClient,
		tailTaskManger: taskMgr,
	}
	return agent, err
}

func (agent *logAgent) StartAgent() {
	log.Infof("agent start is running bind ip %s....", agent.tailTaskManger.BindHost)
	for true {
		msg := agent.tailTaskManger.GetOneLine()
		err := agent.kafkaClient.SendMsg(msg.AppKey, kafka.LogContent{Ip: agent.tailTaskManger.BindHost, Msg: msg.Msg})
		if err != nil {
			log.Errorf("send msg:[%v] topic:[%v] failed, err:[%v]", msg.Msg, msg.AppKey, err)
			time.Sleep(time.Second)
			continue
		}
	}
	return
}
