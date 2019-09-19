package controllers

import (
	"logManager/models"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
)

// 项目列表页面controll
type ProjectListControll struct {
	beego.Controller
}

// 项目列表数据获取 controll
type ProjectGetListControll struct {
	beego.Controller
}

//创建项目 coltroll
type ProjectCreateControll struct {
	beego.Controller
}

// 项目列表数据结构体
type PorjectList struct {
	Code  int                      `json:"code"`
	Msg   string                   `json:"msg"`
	Count int                      `json:"count"`
	Data  []models.ProjectDataInfo `json:"data"`
}

//  project 列表模版
func (this *ProjectListControll) Get() {
	this.Layout = "layout/layout.html"
	this.TplName = "project/list.html"
}

// 获取project列表数据
func (this *ProjectGetListControll) GetData() {
	var err error
	listData := PorjectList{}
	listData.Data, err = models.GetProjectData()
	listData.Count = len(listData.Data)
	if err != nil {
		logs.Error("get project list failed, err: %s", err)
	}
	this.Data["json"] = listData
	this.ServeJSON()
}

func (this *ProjectCreateControll) Get() {
	this.Layout = "layout/layout.html"
	this.TplName = "project/add.html"
}

func (this *ProjectCreateControll) CreateProject() {
	var err error
	respon := make(map[string]interface{})
	reqData := models.CreateProjectData{}
	err = this.ParseForm(&reqData)
	if err != nil {
		logs.Error("get create project data failed, err: %s", err)
		respon["code"] = 1
	} else {
		err = models.CreateProject(&reqData)
		if err != nil {
			logs.Error("create project failed, err: %s", err)
			respon["code"] = 1
		} else {
			respon["code"] = 0
		}
	}

	this.Data["json"] = respon
	this.ServeJSON()
}
