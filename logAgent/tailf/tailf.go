package tailf

import (
	"sync"

	"github.com/astaxie/beego/logs"

	"github.com/hpcloud/tail"
)

const (
	StatusNormal = 1
	StatusDelete = 2
)

// 日志收集任务定义
type Collect struct {
	Topic   string `json:"topic"`
	LogPath string `json:"logPath"`
}

// 发送到kafka的消息
type KafkaMsg struct {
	// Msg *tail.Line `json:"msg"`
	Msg string `json:"msg"`
	Ip  string `json:"ip"`
}

// 发送消息结构体
type TextMsg struct {
	Msg   KafkaMsg
	Topic string
}

// tailf任务对象
type TailObj struct {
	tailobj  *tail.Tail
	collect  Collect
	status   int
	exitChan chan int
}

// tailf任务对象管理
type TailsObjMgr struct {
	tailObjs []*TailObj
	msgChan  chan *TextMsg
	lock     sync.Mutex
}

var (
	//初始化tailf任务对象管理
	tailObjMgr *TailsObjMgr
	hostIp     string
)

// 初始化tailf
func InitTailf(collects []Collect, chanSize int, ip string) (err error) {
	tailObjMgr = &TailsObjMgr{
		msgChan: make(chan *TextMsg, chanSize),
	}

	if len(collects) == 0 {
		// err = errors.New("collect task is nill")
		logs.Warn("collect task is nill")
	}

	hostIp = ip

	// 创建tailf task
	for _, v := range collects {
		createTask(v)
	}
	return
}

// 创建tailf task
func createTask(collect Collect) {
	obj, err := tail.TailFile(collect.LogPath, tail.Config{
		ReOpen: true,
		Follow: true,
		// Location:  &tail.SeekInfo{Offset: 0, Whence: 2},
		MustExist: false,
		Poll:      true,
	})
	if err != nil {
		logs.Warn("tailf create [%v] failed, %v", collect.LogPath, err)
	}
	tailObj := &TailObj{
		tailobj:  obj,
		collect:  collect,
		exitChan: make(chan int, 1),
	}
	// 开启goroute去读取监听日志的内容
	go readFromTail(tailObj, collect.Topic)
	tailObjMgr.tailObjs = append(tailObjMgr.tailObjs, tailObj)
}

// 读取监听日志文件的内容
func readFromTail(tailObj *TailObj, topic string) {
	for true {
		select {
		// 任务正常运行
		case lineMsg, ok := <-tailObj.tailobj.Lines:
			if !ok {
				logs.Warn("read obj:[%v] topic:[%v] filed continue", tailObj, topic)
				continue
			}
			// 消息为空跳过
			if lineMsg.Text == "" {
				continue
			}
			kfMsg := KafkaMsg{
				Msg: lineMsg.Text,
				Ip:  hostIp,
			}
			msgObj := &TextMsg{
				Msg:   kfMsg,
				Topic: topic,
			}
			tailObjMgr.msgChan <- msgObj
		// 任务退出
		case <-tailObj.exitChan:
			logs.Warn("tail obj will exited, conf:%v", tailObj.collect)
			return
		}
	}
}

// 从chan中获取一行数据
func GetOneLine() (msg *TextMsg) {
	msg = <-tailObjMgr.msgChan
	return
}

// 更新tailf任务
func UpdateTailfTask(collectConfig []Collect) (err error) {
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
	var tailObjs []*TailObj
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
