package db

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var DB *gorm.DB

func Open() error {
	var err error
	DB, err = gorm.Open("mysql", "root:password@10.231.50.26/manager?charset=utf8&parseTime=True")
	return err
}
