package task

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"github.com/hpcloud/tail"
	"github.com/yihongzhi/log-collector/common/etcd"
	"github.com/yihongzhi/log-collector/common/logger"
	"github.com/yihongzhi/log-collector/common/utils"
	"strings"
	"sync"
)

//日志任务管理
type TailTaskManger struct {
	TailTasks []*TailTask      //任务列表
	MsgChan   chan *LogTextMsg //消息通道
	Lock      sync.Mutex
}

//Tail 任务
type TailTask struct {
	TailObj  *tail.Tail
	Details  TailTaskDetails
	Status   int
	ExitChan chan int
}

//任务详情
type TailTaskDetails struct {
	AppKey  string `json:"appKey"`  //应用id
	LogPath string `json:"logPath"` //日志路径
}

//日志消息体
type LogTextMsg struct {
	AppKey string //应用ID
	Msg    string //日志消息
}

const (
	StatusNormal = 1
	StatusDelete = 2
)

var (
	hostIp string
	logs   = logger.Instance
)

//初始化收集任务
func NewTailTaskManger(collectKey string, chanSize int, client *etcd.EtcdClient) (*TailTaskManger, error) {
	tailObjMgr := &TailTaskManger{
		MsgChan:   make(chan *LogTextMsg, chanSize),
		TailTasks: []*TailTask{},
	}
	tasks, err := readTask(collectKey, client)
	if err != nil {
		logs.Fatalln("读取任务失败", err)
		return nil, err
	}
	for _, task := range tasks {
		tailTask := createTailTask(task, tailObjMgr.MsgChan)
		tailObjMgr.TailTasks = append(tailObjMgr.TailTasks, tailTask)
	}
	return tailObjMgr, nil
}

// 从chan中获取一行数据
func (tailObjMgr *TailTaskManger) GetOneLine() *LogTextMsg {
	return <-tailObjMgr.MsgChan
}

//读取etcd任务
func readTask(collectKey string, client *etcd.EtcdClient) ([]TailTaskDetails, error) {
	var tasks []TailTaskDetails
	if strings.HasSuffix(collectKey, "/") == false {
		collectKey = fmt.Sprintf("%s/", collectKey)
	}
	for _, ip := range utils.LocalIpArray {
		etcdKey := fmt.Sprintf("%s%s", collectKey, ip)
		resp, err := client.Client.Get(context.Background(), etcdKey)
		if err != nil {
			logs.Warnf("get key: %s from etcd failed, err: %s", etcdKey, err)
			return nil, err
		}
		for _, v := range resp.Kvs {
			if string(v.Key) == etcdKey {
				var task TailTaskDetails
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
	return tasks, nil
}

//创建一个收集任务
func createTailTask(collectTask TailTaskDetails, msgChan chan<- *LogTextMsg) *TailTask {
	tailObj, err := tail.TailFile(collectTask.LogPath, tail.Config{
		ReOpen: true,
		Follow: true,
		// Location:  &tail.SeekInfo{Offset: 0, Whence: 2},
		MustExist: false,
		Poll:      true,
	})
	if err != nil {
		logs.Warnf("task [%v] create failed, %v", collectTask, err)
		return nil
	}
	tailTask := &TailTask{
		TailObj:  tailObj,
		Details:  collectTask,
		ExitChan: make(chan int, 1),
	}
	go readFromTail(tailTask, msgChan)
	return tailTask
}

// 读取监听日志文件的内容
func readFromTail(tailObj *TailTask, msgChan chan<- *LogTextMsg) {
	for true {
		select {
		case lineMsg, ok := <-tailObj.TailObj.Lines:
			if !ok {
				logs.Warnf("read obj:[%v] topic:[%v] filed continue", tailObj, tailObj.Details.AppKey)
				continue
			}
			// 消息为空跳过
			if lineMsg.Text == "" {
				continue
			}
			msgObj := &LogTextMsg{
				Msg:    lineMsg.Text,
				AppKey: tailObj.Details.AppKey,
			}
			msgChan <- msgObj
		// 任务退出
		case <-tailObj.ExitChan:
			logs.Warnf("task [%v] exit ", tailObj.Details)
			return
		}
	}
}

// 更新tailf任务
/*func (tailObjMgr *TailTaskManger) UpdateTailfTask(collectConfig []TailTaskDetails) (err error) {
	tailObjMgr.Lock.Lock()
	defer tailObjMgr.Lock.Unlock()
	for _, newColl := range collectConfig {
		// 判断tailf运行状态，是否存在
		var isRunning = false
		for _, oldTailObj := range tailObjMgr.TailTasks {
			if newColl.LogPath == oldTailObj.Details.LogPath {
				isRunning = true
				break
			}
		}
		// 如果tailf任务不存在，创建新的任务
		if isRunning == false {
			createTailTask(newColl)
		}
	}

	// 更新tailf任务管理列表内容
	var tailObjs []*TailTask
	for _, oldTailObj := range tailObjMgr.TailTasks {
		oldTailObj.Status = StatusDelete
		for _, newColl := range collectConfig {
			if newColl.LogPath == oldTailObj.Details.LogPath {
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
}*/

// 初始化etcd key监控
/*func initEtcdWatch() {
	for _, key := range etcdClient.collectKeys {
		go etcdWatch(key)
	}
}*/

// 	etcd key监控处理
func etcdWatch(key string, client *etcd.EtcdClient) {
	logs.Debug("start watch key: %s", key)
	for true {
		rech := client.Client.Watch(context.Background(), key)
		var colConfig []TailTaskDetails
		var getConfStatus = true
		for wresp := range rech {
			for _, ev := range wresp.Events {
				// key 删除
				if ev.Type == mvccpb.DELETE {
					logs.Warn("key [%s] is deleted", key)
					continue
				}
				// key 更新
				if ev.Type == mvccpb.PUT && string(ev.Kv.Key) == key {
					err := json.Unmarshal(ev.Kv.Value, &colConfig)
					if err != nil {
						logs.Error("key [%s], Unmarshal[%s], err:%s", key, err)
						getConfStatus = false
						continue
					}
				}
				logs.Debug("get etcd config, %s %q : %q", ev.Type, ev.Kv.Key, ev.Kv.Value)
			}
			if getConfStatus {
				break
			}
		}
		logs.Info("Update task config")
		// 更新tailf任务
		/*err := UpdateTailfTask(colConfig)
		if err != nil {
			logs.Error("Update task task failed, connect: %s, err: %s", colConfig, err)
		}*/
	}
}
