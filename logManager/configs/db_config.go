package configs

import "github.com/astaxie/beego/orm"
import _ "github.com/go-sql-driver/mysql"

func init() {
	orm.RegisterDriver("mysql", orm.DRMySQL)
	orm.RegisterDataBase("default", "mysql", "root:peE5RTd2@tcp(10.231.50.30:13306)/log_collect?charset=utf8")
	//orm.Debug = true
}
