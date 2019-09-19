package main

import (
	"errors"
	"logManager/models"
	_ "logManager/routers"
	"strings"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
)

// 初始化mysq
func initMysql() (err error) {
	mysqlInfo := models.MysqlInfo{
		Address:  beego.AppConfig.String("mysqladdr"),
		Port:     beego.AppConfig.String("mysqlport"),
		DbName:   beego.AppConfig.String("mysqldbname"),
		User:     beego.AppConfig.String("mysqluser"),
		Password: beego.AppConfig.String("mysqlpwd"),
	}
	err = models.InitMySql(mysqlInfo)
	if err != nil {
		return
	}
	return
}

// 初始化etcd
func initEtcd() (err error) {
	address := beego.AppConfig.String("etcdaddr")
	if len(address) == 0 {
		errors.New("etcd addresss config error")
		return
	}
	etcdAddress := strings.Split(address, ",")
	err = models.InitEtcd(etcdAddress)
	if err != nil {
		return
	}
	return
}

func main() {
	var err error

	// 初始化log
	err = initLog()
	if err != nil {
		logs.Error("init log manager failed, err: %s", err)
	}
	// 初始化mysql
	err = initMysql()
	if err != nil {
		logs.Error("init mysql failed, err: %s", err)
	}
	// 初始化etcd
	err = initEtcd()
	if err != nil {
		logs.Error("init etcd failed, err: %s", err)
	}

	beego.Run()
}
