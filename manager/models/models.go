package models

import "github.com/astaxie/beego/orm"
import _ "github.com/go-sql-driver/mysql"

func init() {
	orm.RegisterDriver("mysql", orm.DRMySQL)
	orm.RegisterDataBase("default", "mysql", "root:peE5RTd2@tcp(10.231.50.30:3306)/log_collect?charset=utf8")
	//orm.Debug = true
}

//分页实体
type Pagination struct {
	Current  int32 `json:"current"`
	PageSize int32 `json:"pageSize"`
	Total    int64 `json:"total"`
}
