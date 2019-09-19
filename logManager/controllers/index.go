package controllers

import (
	"github.com/astaxie/beego"
)

// 首页controll
type IndexControll struct {
	beego.Controller
}

// home页controll
type HomeControll struct {
	beego.Controller
}

func (this *IndexControll) Index() {
	this.Layout = "layout/layout.html"
	this.TplName = "app/index.html"
}

func (this *HomeControll) Home() {
	this.Layout = "layout/layout.html"
	this.TplName = "app/home.html"
}
