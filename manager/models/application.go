package models

import "time"

//应用信息
type Application struct {
	ID         string    `gorm:"primary_key"`
	AppKey     string    `gorm:"type:varchar(100);unique_index"`
	AppName    string    `gorm:"type:varchar(100);"`
	Hosts      string    `gorm:"type:varchar(500);"`
	Remark     string    `gorm:"type:varchar(500);"`
	Owner      string    `gorm:"type:varchar(50);"`
	CreateTime time.Time `orm:"column(create_time);type(datetime);null" description:"创建时间"`
	UpdateTime time.Time `orm:"column(update_time);type(datetime);null" description:"更新时间"`
}
