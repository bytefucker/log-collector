package controllers

import (
	"logManager/models"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
)

// log list 页面 controll
type LogListControll struct {
	beego.Controller
}

// log list 数据controll
type LogGetListControll struct {
	beego.Controller
}

// log list 数据结构体
type LogList struct {
	Code  int                  `json:"code"`
	Msg   string               `json:"msg"`
	Count int                  `json:"count"`
	Data  []models.LogDataInfo `json:"data"`
}

// 创建log task controll
type LogCreateControll struct {
	beego.Controller
}

// 删除log task controll
type LogDeleteControll struct {
	beego.Controller
}

func (this *LogListControll) Get() {
	this.Layout = "layout/layout.html"
	this.TplName = "log/list.html"
}

func (this *LogGetListControll) GetData() {
	var err error
	listData := LogList{}
	listData.Data, err = models.GetLogData()
	listData.Count = len(listData.Data)
	if err != nil {
		logs.Error("get log list failed, err: %s", err)
	}
	this.Data["json"] = &listData
	this.ServeJSON()
}

func (this *LogCreateControll) Get() {
	this.Layout = "layout/layout.html"
	this.TplName = "log/add.html"
}

func (this *LogCreateControll) CreateLog() {
	var err error
	respon := make(map[string]interface{})
	reqData := models.CreateLogData{}
	err = this.ParseForm(&reqData)
	if err != nil {
		logs.Error("get create log data failed, err: %s", err)
		respon["code"] = 1
	} else {
		err = models.CreateLog(&reqData)
		if err != nil {
			logs.Error("create log failed, err: %s", err)
			respon["code"] = 1
		} else {
			respon["code"] = 0
		}
	}

	this.Data["json"] = respon
	this.ServeJSON()
}

func (this *LogDeleteControll) DeleteLog() {
	var err error
	respon := make(map[string]interface{})
	reqData := models.LogDataInfo{}
	err = this.ParseForm(&reqData)
	if err != nil {
		logs.Error("get delete log data failed, err: %s", err)
		respon["code"] = 1
	} else {
		err = models.DeleteLog(&reqData)
		if err != nil {
			logs.Error("delete log failed, err: %s", err)
			respon["code"] = 1
		} else {
			respon["code"] = 0
		}
	}

	this.Data["json"] = respon
	this.ServeJSON()
}
