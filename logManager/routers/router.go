package routers

import (
	"logManager/controllers"

	"github.com/astaxie/beego"
)

func init() {
	// 路由
	// index
	beego.Router("/", &controllers.IndexControll{}, "get,post:Index")
	beego.Router("/home", &controllers.HomeControll{}, "get,post:Home")

	// project
	beego.Router("/project", &controllers.ProjectListControll{})
	beego.Router("/project/getList", &controllers.ProjectGetListControll{}, "get,post:GetData")
	beego.Router("/project/create", &controllers.ProjectCreateControll{}, "post:CreateProject;get:Get")

	// log
	beego.Router("/log", &controllers.LogListControll{})
	beego.Router("/log/getList", &controllers.LogGetListControll{}, "get,post:GetData")
	beego.Router("/log/create", &controllers.LogCreateControll{}, "post:CreateLog;get:Get")
	beego.Router("/log/delete", &controllers.LogDeleteControll{}, "post:DeleteLog")

}
