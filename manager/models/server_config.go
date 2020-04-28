package models

import (
	"time"
)

type ServerConfig struct {
	Id         string    `orm:"column(id);pk" description:"主键"`
	Hostname   string    `orm:"column(hostname);size(255)" description:"hostname"`
	Ip         string    `orm:"column(ip);size(50)" description:"ip"`
	Remark     string    `orm:"column(remark);size(500)" description:"用途描述"`
	Owner      string    `orm:"column(owner);size(50);null" description:"拥有者"`
	CreateTime time.Time `orm:"column(create_time);type(datetime);null" description:"创建时间"`
	UpdateTime time.Time `orm:"column(update_time);type(datetime);null" description:"更新时间"`
}
