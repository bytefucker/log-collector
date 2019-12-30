package model

import "time"

//结构化日志
type StructuredLog struct {
	AppKey     string            `json:"appKey"`     //应用id
	Host       string            `json:"host"`       //机器Ip
	DateTime   time.Time         `json:"datetime"`   //日志时间
	Level      string            `json:"level"`      //日志等级
	Catalog    string            `json:"catalog"`    //大类
	SubCatalog string            `json:"subCatalog"` //小类
	Content    string            `json:"content"`    //日志内容
	Exception  string            `json:"exception"`  //异常
	Attributes map[string]string `json:"attributes"` //附件属性
}
