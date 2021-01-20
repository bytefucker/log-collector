package task

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"github.com/hpcloud/tail"
	"github.com/yihongzhi/log-collector/common/etcd"
	"github.com/yihongzhi/log-collector/common/logger"
	"github.com/yihongzhi/log-collector/common/utils"
	"path"
	"strings"
	"sync"
)

//日志任务管理
type TailTaskManger struct {
	TailTasks []*TailTask      //任务列表
	MsgChan   chan *LogTextMsg //消息通道
	Key       string
	Lock      sync.Mutex
	BindHost  string //绑定的机器IP
}

//Tail 任务
type TailTask struct {
	TailObj  *tail.Tail
	Details  TailTaskDetails
	Key      string
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

var logs = logger.Instance

//初始化收集任务
func NewTailTaskManger(collectKey string, chanSize int, client *etcd.EtcdClient) (*TailTaskManger, error) {
	ipArray := utils.LocalIpArray
	if len(ipArray) == 0 {
		return nil, errors.New("获取机器绑定ip失败")
	}
	tailObjMgr := &TailTaskManger{
		MsgChan:   make(chan *LogTextMsg, chanSize),
		TailTasks: []*TailTask{},
		BindHost:  ipArray[0],
		Key:       collectKey,
	}
	collectKey = path.Join(collectKey, tailObjMgr.BindHost)
	tasks, err := readTask(collectKey, client)
	if err != nil {
		logs.Fatalln("读取任务失败", err)
		return nil, err
	}
	for _, task := range tasks {
		tailTask := createTailTask(task, tailObjMgr.MsgChan)
		tailObjMgr.TailTasks = append(tailObjMgr.TailTasks, tailTask)
	}
	//开启监听
	go watchTask(tailObjMgr, client)
	return tailObjMgr, nil
}

// 从chan中获取一行数据
func (tailObjMgr *TailTaskManger) GetOneLine() *LogTextMsg {
	return <-tailObjMgr.MsgChan
}

//读取etcd任务
func readTask(collectKey string, client *etcd.EtcdClient) ([]TailTaskDetails, error) {
	var tasks []TailTaskDetails
	resp, err := client.Client.Get(context.Background(), collectKey, clientv3.WithPrefix())
	if err != nil {
		logs.Warnf("get key: %s from etcd failed, err: %s", collectKey, err)
		return nil, err
	}
	for _, v := range resp.Kvs {
		if strings.HasPrefix(string(v.Key), collectKey) {
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
	return tasks, nil
}

//监听任务
func watchTask(mgr *TailTaskManger, client *etcd.EtcdClient) {
	logs.Infof("start watch key: %s", mgr.Key)
	for true {
		rech := client.Client.Watch(context.Background(), mgr.Key, clientv3.WithPrefix())
		for wresp := range rech {
			for _, ev := range wresp.Events {
				tailTask := findTailTasks(string(ev.Kv.Key), mgr.TailTasks)
				if tailTask != nil {
					// key 删除
					if ev.Type == mvccpb.DELETE {
						logs.Infof("key delete, %s %q : %q", ev.Type, ev.Kv.Key, ev.Kv.Value)
						tailTask.ExitChan <- 0
						tailTask.Status = StatusDelete
						continue
					}
					// key 更新
					if ev.Type == mvccpb.PUT {
						logs.Infof("key update, %s %q : %q", ev.Type, ev.Kv.Key, ev.Kv.Value)
						continue
					}
				} else {
					// key 新增
					if ev.Type == mvccpb.PUT {
						logs.Infof("key add, %s %q : %q", ev.Type, ev.Kv.Key, ev.Kv.Value)
						details := TailTaskDetails{}
						if err := json.Unmarshal(ev.Kv.Value, &details); err == nil {
							task := createTailTask(details, mgr.MsgChan)
							mgr.TailTasks = append(mgr.TailTasks, task)
						}
					}
				}
			}
		}
	}
}

//筛选任务
func findTailTasks(key string, tasks []*TailTask) *TailTask {
	for _, task := range tasks {
		if task.Key == key {
			return task
		}
	}
	return nil
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
		Status:   StatusNormal,
		ExitChan: make(chan int, 1),
	}
	go readFromTail(tailTask, msgChan)
	return tailTask
}

// 读取日志文件
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
			logs.Infof("task [%v] exit ", tailObj.Details)
			tailObj.Status = StatusDelete
			return
		}
	}
}
