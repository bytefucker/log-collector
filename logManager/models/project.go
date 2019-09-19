package models

import (
	"fmt"
	"strings"

	"github.com/astaxie/beego/logs"
	_ "github.com/go-sql-driver/mysql"
)

// project数据信息
type ProjectDataInfo struct {
	Id         int    `db:"id" json:"id"`
	Name       string `db:"pname" json:"name"`
	Type       string `db:"type" json:"type"`
	CreateTime string `db:"createTime" json:"createTime"`
	ApplyPath  string `db:"applyPath" json:"applyPath"`
}

// 创建project form
type CreateProjectData struct {
	Name      string `form:"pname"`
	Type      string `form:"type"`
	ApplyPath string `form:"applyPath"`
}

// 获取项目列表数据
func GetProjectData() (listData []ProjectDataInfo, err error) {
	err = mysqlClient.Select(&listData, "SELECT * FROM log_project")
	if err != nil {
		logs.Error("Select all project list failed, err: %s", err)
		return
	}
	return
}

// 创建项目
func CreateProject(data *CreateProjectData) (err error) {
	ts, err := mysqlClient.Begin()
	if err != nil {
		return
	}
	if !strings.HasSuffix(data.ApplyPath, "/") {
		data.ApplyPath = fmt.Sprintf("%s/", data.ApplyPath)
	}
	_, err = ts.Exec("INSERT INTO log_project (pname, type, applyPath, createTime) VALUES (?,?,?,?)", data.Name, data.Type, data.ApplyPath, NowTime())
	if err != nil {
		ts.Rollback()
		return
	}
	ts.Commit()
	return
}
