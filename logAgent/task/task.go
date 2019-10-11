package task

import (
	"sync"

	"github.com/astaxie/beego/logs"

	"github.com/hpcloud/tail"
)

const (
	StatusNormal = 1
	StatusDelete = 2
)

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
	tailObj  *tail.Tail
	collect  CollectTask
	status   int
	exitChan chan int
}

type TailsTaskMgr struct {
	tailObjs []*TailTask
	msgChan  chan *LogTextMsg
	lock     sync.Mutex
}

var (
	tailObjMgr *TailsTaskMgr
	hostIp     string
)

//初始化收集任务
func InitTailfTask(collectTasks []CollectTask, chanSize int, ip string) (err error) {
	tailObjMgr = &TailsTaskMgr{
		msgChan: make(chan *LogTextMsg, chanSize),
	}
	if len(collectTasks) == 0 {
		logs.Warn("没有收集任务")
	}
	hostIp = ip

	for _, v := range collectTasks {
		createTask(v)
	}
	return
}

//创建一个收集任务
func createTask(collectTask CollectTask) {
	obj, err := tail.TailFile(collectTask.LogPath, tail.Config{
		ReOpen: true,
		Follow: true,
		// Location:  &tail.SeekInfo{Offset: 0, Whence: 2},
		MustExist: false,
		Poll:      true,
	})
	if err != nil {
		logs.Warn("收集任务[%v]创建失败, %v", collectTask.LogPath, err)
	}
	tailObj := &TailTask{
		tailObj:  obj,
		collect:  collectTask,
		exitChan: make(chan int, 1),
	}
	go readFromTail(tailObj, collectTask.Topic)
	tailObjMgr.tailObjs = append(tailObjMgr.tailObjs, tailObj)
}

// 读取监听日志文件的内容
func readFromTail(tailObj *TailTask, topic string) {
	for true {
		select {
		// 任务正常运行
		case lineMsg, ok := <-tailObj.tailObj.Lines:
			if !ok {
				logs.Warn("read obj:[%v] topic:[%v] filed continue", tailObj, topic)
				continue
			}
			// 消息为空跳过
			if lineMsg.Text == "" {
				continue
			}
			kfMsg := LogContent{
				Msg: lineMsg.Text,
				Ip:  hostIp,
			}
			msgObj := &LogTextMsg{
				Msg:   kfMsg,
				Topic: topic,
			}
			tailObjMgr.msgChan <- msgObj
		// 任务退出
		case <-tailObj.exitChan:
			logs.Warn("收集任务退出[%v]", tailObj.collect)
			return
		}
	}
}

// 从chan中获取一行数据
func GetOneLine() (msg *LogTextMsg) {
	msg = <-tailObjMgr.msgChan
	return
}

// 更新tailf任务
func UpdateTailfTask(collectConfig []CollectTask) (err error) {
	tailObjMgr.lock.Lock()
	defer tailObjMgr.lock.Unlock()
	for _, newColl := range collectConfig {
		// 判断tailf运行状态，是否存在
		var isRunning = false
		for _, oldTailObj := range tailObjMgr.tailObjs {
			if newColl.LogPath == oldTailObj.collect.LogPath {
				isRunning = true
				break
			}
		}
		// 如果tailf任务不存在，创建新的任务
		if isRunning == false {
			createTask(newColl)
		}
	}

	// 更新tailf任务管理列表内容
	var tailObjs []*TailTask
	for _, oldTailObj := range tailObjMgr.tailObjs {
		oldTailObj.status = StatusDelete
		for _, newColl := range collectConfig {
			if newColl.LogPath == oldTailObj.collect.LogPath {
				oldTailObj.status = StatusNormal
				break
			}
		}
		if oldTailObj.status == StatusDelete {
			oldTailObj.exitChan <- 1
			continue
		}
		tailObjs = append(tailObjs, oldTailObj)
	}
	tailObjMgr.tailObjs = tailObjs
	return
}
