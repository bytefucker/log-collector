package model

import (
	"github.com/hpcloud/tail"
	"sync"
)

// 存储logAgent配置信息
type Config struct {
	LogLevel     string
	LogPath      string
	ChanSize     int
	SendModel    string //发送方式
	KafkaAddress []string
	EtcdAddress  []string
	CollectKey   string
	CollectTasks []CollectTask
	Ip           string
}


type CollectTask struct {
	Topic   string `json:"topic"`
	LogPath string `json:"logPath"`
}

type LogContent struct {
	Msg string `json:"msg"`
	Ip  string `json:"ip"`
}

type LogTextMsg struct {
	Msg   LogContent
	Topic string
}

type TailTask struct {
	TailObj  *tail.Tail
	Collect  CollectTask
	Status   int
	ExitChan chan int
}

type TailsTaskMgr struct {
	TailObjs []*TailTask
	MsgChan  chan *LogTextMsg
	Lock     sync.Mutex
}
