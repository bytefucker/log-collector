package models

import (
	"time"
)

type ApplicationConfig struct {
	Id         int       `orm:"column(id);pk" description:"主键"`
	AppKey     string    `orm:"column(app_key);size(255)" description:"唯一应用标识"`
	AppName    string    `orm:"column(app_name);size(255)" description:"应用名称"`
	Hosts      string    `orm:"column(hosts);size(500);null" description:"部署服务器"`
	Remark     string    `orm:"column(remark);size(500);null" description:"应用描述"`
	Owner      string    `orm:"column(owner);size(50);null" description:"负责人"`
	CreateTime time.Time `orm:"column(create_time);type(datetime);null" description:"创建时间"`
	UpdateTime time.Time `orm:"column(update_time);type(datetime);null" description:"更新时间"`
}
