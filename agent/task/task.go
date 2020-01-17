package task

import (
	"github.com/yihongzhi/logCollect/agent/model"
	"github.com/yihongzhi/logCollect/common/logger"

	"github.com/hpcloud/tail"
)

const (
	StatusNormal = 1
	StatusDelete = 2
)

var (
	tailObjMgr *model.TailsTaskMgr
	hostIp     string
	logs       = logger.Instance
)

//初始化收集任务
func InitTailfTask(agentConfig *model.Config) (err error) {
	tailObjMgr = &model.TailsTaskMgr{
		MsgChan: make(chan *model.LogTextMsg, agentConfig.ChanSize),
	}
	if len(agentConfig.CollectTasks) == 0 {
		logs.Warn("没有收集任务")
	}
	hostIp = agentConfig.Ip

	for _, v := range agentConfig.CollectTasks {
		createTask(v)
	}
	return
}

//创建一个收集任务
func createTask(collectTask model.CollectTask) {
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
	tailObj := &model.TailTask{
		TailObj:  obj,
		Collect:  collectTask,
		ExitChan: make(chan int, 1),
	}
	go readFromTail(tailObj, collectTask.AppKey)
	tailObjMgr.TailObjs = append(tailObjMgr.TailObjs, tailObj)
}

// 读取监听日志文件的内容
func readFromTail(tailObj *model.TailTask, appKey string) {
	for true {
		select {
		// 任务正常运行
		case lineMsg, ok := <-tailObj.TailObj.Lines:
			if !ok {
				logs.Warn("read obj:[%v] topic:[%v] filed continue", tailObj, appKey)
				continue
			}
			// 消息为空跳过
			if lineMsg.Text == "" {
				continue
			}
			kfMsg := model.LogContent{
				Msg: lineMsg.Text,
				Ip:  hostIp,
			}
			msgObj := &model.LogTextMsg{
				Msg:    kfMsg,
				AppKey: appKey,
			}
			tailObjMgr.MsgChan <- msgObj
		// 任务退出
		case <-tailObj.ExitChan:
			logs.Warn("收集任务退出[%v]", tailObj.Collect)
			return
		}
	}
}

// 从chan中获取一行数据
func GetOneLine() (msg *model.LogTextMsg) {
	msg = <-tailObjMgr.MsgChan
	return
}

// 更新tailf任务
func UpdateTailfTask(collectConfig []model.CollectTask) (err error) {
	tailObjMgr.Lock.Lock()
	defer tailObjMgr.Lock.Unlock()
	for _, newColl := range collectConfig {
		// 判断tailf运行状态，是否存在
		var isRunning = false
		for _, oldTailObj := range tailObjMgr.TailObjs {
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
	var tailObjs []*model.TailTask
	for _, oldTailObj := range tailObjMgr.TailObjs {
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
	tailObjMgr.TailObjs = tailObjs
	return
}
