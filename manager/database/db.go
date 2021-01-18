package database

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/yihongzhi/log-collector/manager/models"
	"log"
)

var DB *gorm.DB

func Open(connectStr string) error {
	var err error
	DB, err = gorm.Open("mysql", connectStr)
	if err != nil {
		log.Panicln("open mysql error: ", err)
	}
	DB.LogMode(true)
	DB.SingularTable(true)
	DB.AutoMigrate(&models.Application{}, &models.ServerInfo{})
	return err
}
