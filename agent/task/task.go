package task

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hpcloud/tail"
	"github.com/yihongzhi/logCollect/common/etcd"
	"github.com/yihongzhi/logCollect/common/logger"
	"github.com/yihongzhi/logCollect/common/utils"
	"strings"
	"sync"
)

//日志任务管理
type TailsTaskMgr struct {
	TailTasks []*TailTask      //任务列表
	MsgChan   chan *LogTextMsg //消息通道
	Lock      sync.Mutex
}

//Tail 任务
type TailTask struct {
	TailObj  *tail.Tail
	Collect  CollectTask
	Status   int
	ExitChan chan int
}

//任务详情
type CollectTask struct {
	AppKey  string `json:"appKey"`  //应用id
	LogPath string `json:"logPath"` //日志路径
}

//日志消息体
type LogTextMsg struct {
	Msg    string //日志消息
	AppKey string //应用ID
}

const (
	StatusNormal = 1
	StatusDelete = 2
)

var (
	tailObjMgr *TailsTaskMgr
	hostIp     string
	logs       = logger.Instance
)

//初始化收集任务
func InitTailfTask(collectKey string, chanSize int, client *etcd.EtcdClient) (err error) {
	var tasks []CollectTask
	if strings.HasSuffix(collectKey, "/") == false {
		collectKey = fmt.Sprintf("%s/", collectKey)
	}
	for _, ip := range utils.LocalIpArray {
		etcdKey := fmt.Sprintf("%s%s", collectKey, ip)
		resp, err := client.Client.Get(context.Background(), etcdKey)
		if err != nil {
			logs.Warn("get key: %s from etcd failed, err: %s", etcdKey, err)
			continue
		}
		for _, v := range resp.Kvs {
			if string(v.Key) == etcdKey {
				var task CollectTask
				err = json.Unmarshal(v.Value, &task)
				if err != nil {
					logs.Warnf("json Unmarshal key: %s failed, err: %s", v.Key, err)
					continue
				}
				logs.Debugf("log agent task is: %v", task)
				tasks = append(tasks, task)
			}
		}
	}
	tailObjMgr = &TailsTaskMgr{
		MsgChan: make(chan *LogTextMsg, chanSize),
	}
	if len(tasks) == 0 {
		logs.Warn("no task running....")
	}
	for _, task := range tasks {
		createTask(task)
	}
	return
}

// 从chan中获取一行数据
func GetOneLine() (msg *LogTextMsg) {
	msg = <-tailObjMgr.MsgChan
	return
}

//创建一个收集任务
func createTask(collectTask CollectTask) {
	tailObj, err := tail.TailFile(collectTask.LogPath, tail.Config{
		ReOpen: true,
		Follow: true,
		// Location:  &tail.SeekInfo{Offset: 0, Whence: 2},
		MustExist: false,
		Poll:      true,
	})
	if err != nil {
		logs.Warnf("task [%v] create failed, %v", collectTask, err)
		return
	}
	tailTask := &TailTask{
		TailObj:  tailObj,
		Collect:  collectTask,
		ExitChan: make(chan int, 1),
	}
	go readFromTail(tailTask, collectTask.AppKey)
	tailObjMgr.TailTasks = append(tailObjMgr.TailTasks, tailTask)
}

// 读取监听日志文件的内容
func readFromTail(tailObj *TailTask, appKey string) {
	for true {
		select {
		case lineMsg, ok := <-tailObj.TailObj.Lines:
			if !ok {
				logs.Warnf("read obj:[%v] topic:[%v] filed continue", tailObj, appKey)
				continue
			}
			// 消息为空跳过
			if lineMsg.Text == "" {
				continue
			}
			msgObj := &LogTextMsg{
				Msg:    lineMsg.Text,
				AppKey: appKey,
			}
			tailObjMgr.MsgChan <- msgObj
		// 任务退出
		case <-tailObj.ExitChan:
			logs.Warnf("task [%v] exit ", tailObj.Collect)
			return
		}
	}
}

// 更新tailf任务
func UpdateTailfTask(collectConfig []CollectTask) (err error) {
	tailObjMgr.Lock.Lock()
	defer tailObjMgr.Lock.Unlock()
	for _, newColl := range collectConfig {
		// 判断tailf运行状态，是否存在
		var isRunning = false
		for _, oldTailObj := range tailObjMgr.TailTasks {
			if newColl.LogPath == oldTailObj.Collect.LogPath {
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
	for _, oldTailObj := range tailObjMgr.TailTasks {
		oldTailObj.Status = StatusDelete
		for _, newColl := range collectConfig {
			if newColl.LogPath == oldTailObj.Collect.LogPath {
				oldTailObj.Status = StatusNormal
				break
			}
		}
		if oldTailObj.Status == StatusDelete {
			oldTailObj.ExitChan <- 1
			continue
		}
		tailObjs = append(tailObjs, oldTailObj)
	}
	tailObjMgr.TailTasks = tailObjs
	return
}
