package models

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/astaxie/beego/logs"
)

const (
	TopicsName = "topics"
)

// 日志收集数据结构
type LogDataInfo struct {
	Id         int    `json:"id" db:"id" form:"id"`
	Host       string `json:"host" db:"host"`
	Topic      string `json:"topic" db:"topic" form:"topic"`
	Pname      string `json:"pname" db:"pname"`
	LogPath    string `json:"logPath" db:"logPath" form:"logPath"`
	ApplyPath  string `json:"applyPath" db:"applyPath" form:"applyPath"`
	CreateTime string `json:"createTime" db:"createTime"`
}

// 创建日志任务form数据
type CreateLogData struct {
	Hosts     string `form:"hosts"`
	Topic     string `form:"topic"`
	Pname     string `form:"pname"`
	LogPath   string `form:"logPath"`
	ApplyPath string `form: "applyPath"`
}

// 日志收集任务结构体
type LogCollect struct {
	Topic   string `json:"topic"`
	LogPath string `json:"logPath"`
}

// 获取日志列表数据
func GetLogData() (listData []LogDataInfo, err error) {
	err = mysqlClient.Select(&listData, "SELECT a.id, a.host, a.topic, b.pname, a.applyPath, a.createTime, a.logPath FROM log_collects a, log_project b WHERE a.pid=b.id")
	if err != nil {
		logs.Error("Select all log list failed, err: %s", err)
		return
	}
	return
}

// 创建日志任务
func CreateLog(data *CreateLogData) (err error) {
	ts, err := mysqlClient.Begin()
	if err != nil {
		return
	}
	// 获取pname id and applyPath
	var pinfo []ProjectDataInfo
	err = mysqlClient.Select(&pinfo, "SELECT * FROM log_project WHERE pname=?", data.Pname)
	if err != nil {
		logs.Error("select project name id and applyPath failed, err:%s", err)
	} else {
		data.ApplyPath = pinfo[0].ApplyPath
		for _, host := range strings.Split(data.Hosts, ",") {
			applyPath := fmt.Sprintf("%s%s", data.ApplyPath, host)
			// 插入数据库数据
			_, err = ts.Exec("INSERT INTO log_collects (host, topic, logPath, pid, applyPath, createTime) VALUES (?,?,?,?,?,?)",
				host, data.Topic, data.LogPath, pinfo[0].Id, applyPath, NowTime())
			if err != nil {
				logs.Error("Insert into to log_collects failed, err:%s", err)
				break
			}
			// 更新数据库插入的key到etcd
			err = setEtcdLogCollect(applyPath, data.LogPath, data.Topic)
			if err != nil {
				logs.Error("set etcd info faile failed. applyPath:[%s], logPath:[%s], topic:[%s], err:%s",
					applyPath, data.LogPath, data.Topic, err)
				break
			}
		}
	}
	if err != nil {
		ts.Rollback()
		return
	}
	// 向etcd设置topics列表
	err = setEtcdTopics(data.Topic)
	if err != nil {
		logs.Error("set etcd topics failed, err: %s", err)
		return
	}
	ts.Commit()
	return
}

// 删除日志任务
func DeleteLog(data *LogDataInfo) (err error) {
	ts, err := mysqlClient.Begin()
	if err != nil {
		return
	}
	// 数据库删除日志任务
	_, err = ts.Exec("DELETE FROM log_collects WHERE id=?", data.Id)
	if err != nil {
		logs.Error("delete log failed for mysql, key:[%s], topic:[%s], logPath:[%s], err: %s",
			data.ApplyPath, data.Topic, data.LogPath, err)
		ts.Rollback()
		return
	}
	// 删除etcd中日志任务
	err = deleteEtcdLogCollect(data.ApplyPath, data.LogPath, data.Topic)
	if err != nil {
		ts.Rollback()
		logs.Error("delete log failed for etcd, key:[], topic:[%s], logPath:[%s], err: %s",
			data.ApplyPath, data.LogPath, data.Topic, err)
		return
	}
	ts.Commit()
	return
}

// 设置etcd 日志任务
func setEtcdLogCollect(key, logPath, topic string) (err error) {
	var logCollects []*LogCollect
	// 获取etcd中存在的数据
	oldData, err := GetEtcdKeyData(key, logCollects)
	if err != nil {
		logs.Error(err)
		return
	}
	err = json.Unmarshal(oldData, &logCollects)
	if err != nil {
		logs.Error("json unmarshal data failed, err: %s", err)
		return
	}
	rowData := &LogCollect{
		LogPath: logPath,
		Topic:   topic,
	}
	logCollects = append(logCollects, rowData)
	newData, err := json.Marshal(logCollects)
	if err != nil {
		logs.Error("json marshal data failed, err: %s", err)
		return
	}
	// 向etcd中设置新的数据
	err = SetEtcdKeyData(key, string(newData))
	if err != nil {
		logs.Error(err)
		return
	}
	// debug
	tmpData, _ := GetEtcdKeyData(key, nil)
	logs.Debug("key: [%s], value: [%s] form etcd", key, tmpData)
	return
}

// 设置topic到etcd topics
func setEtcdTopics(topic string) (err error) {
	var topics []string
	// 获取旧的数据
	oldData, err := GetEtcdKeyData(TopicsName, topic)
	if err != nil {
		logs.Error(err)
		return
	}
	if len(oldData) == 0 {
		err = json.Unmarshal(oldData, &topics)
		if err != nil {
			logs.Error("json unmarshal data failed -> a, err: %s", err)
			return
		}
	}
	var existsStatus = false
	for _, tv := range oldData {
		if topic == string(tv) {
			existsStatus = true
			break
		}
	}
	if existsStatus {
		return
	}
	topics = append(topics, topic)
	newData, err := json.Marshal(topics)
	if err != nil {
		logs.Error("json marshal data failed -> b, err: %s", err)
		return
	}
	// 设置新的数据
	err = SetEtcdKeyData(TopicsName, string(newData))
	if err != nil {
		logs.Error(err)
		return
	}
	// debug
	tmpData, _ := GetEtcdKeyData(TopicsName, nil)
	logs.Debug("key: [%s], value: [%s] form etcd", TopicsName, tmpData)
	return
}

// 删除etcd log 指定任务
func deleteEtcdLogCollect(key, logPath, topic string) (err error) {
	var logCollects []*LogCollect
	oldData, err := GetEtcdKeyData(key, logCollects)
	if err != nil {
		logs.Error(err)
		return
	}
	err = json.Unmarshal(oldData, &logCollects)
	if err != nil || len(logCollects) == 0 {
		logs.Error("json unmarshal data failed -> c, logCollects: %s. err: %s", logCollects, err)
		return
	}
	var newLogCollects []*LogCollect

	for _, lv := range newLogCollects {
		if lv.LogPath == logPath && lv.Topic == topic {
			continue
		}
		newLogCollects = append(newLogCollects, lv)
	}

	newData, err := json.Marshal(newLogCollects)
	if err != nil {
		logs.Error("json marshal data failed -> d, err: %s", err)
		return
	}
	err = SetEtcdKeyData(key, string(newData))
	if err != nil {
		logs.Error(err)
		return
	}
	tmpData, _ := GetEtcdKeyData(key, nil)
	logs.Debug("key: [%s], value: [%s] form etcd", key, tmpData)
	return
}
