package configs

import (
	"github.com/astaxie/beego/config"
	"github.com/astaxie/beego/logs"
)

//系统配置
type AppConfig struct {
	Port       int
	DataSource string
}

func InitConf() (appConfig *AppConfig, err error) {
	conf, err := config.NewConfig("yaml", "conf/app.yaml")
	if err != nil {
		logs.Error("解析配置文件错误", err)
		return
	}
	appConfig = &AppConfig{}
	appConfig.Port, _ = conf.Int("port")
	return
}
