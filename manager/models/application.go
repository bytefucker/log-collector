package models

type Application struct {
	ID      int64  `gorm:"primary_key"`
	AppKey  string `gorm:"type:varchar(100);unique_index"`
	AppName string `gorm:"type:varchar(100);"`
	Hosts   string `gorm:"type:varchar(500);"`
	Remark  string `gorm:"type:varchar(500);"`
	Owner   string `gorm:"type:varchar(50);"`
}
