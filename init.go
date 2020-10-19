package main

import (
	"fmt"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/bkzy-wangjp/MicEngine/models"
)

func init() {
	//设置日志信息
	lgcfg := models.EngineCfgMsg.Log
	logs.Async(1e3) //异步模式
	logset := fmt.Sprintf(`{"level":%d,"color":true}`, lgcfg.Consolelevel)
	logs.SetLogger(logs.AdapterConsole, logset) //屏幕输出设置
	logset = fmt.Sprintf(`{"filename":"%s","level":%d,"maxlines":%d,"maxsize":0,"daily":true,"maxdays":%d,"separate":["emergency", "alert", "critical", "error", "warning", "notice", "info", "debug"]}`,
		lgcfg.Path, lgcfg.Level, lgcfg.Maxlines, lgcfg.Maxdays)
	logs.SetLogger(logs.AdapterMultiFile, logset) //文件输出设置
	logs.EnableFuncCallDepth(lgcfg.Debug)

	beego.BConfig.WebConfig.Session.SessionOn = true //开启session,配置文件对应的参数名：sessionon
	//beego.BConfig.WebConfig.Session.SessionName = "MicEngingSeessionId" //设置 cookies 的名字，Session 默认是保存在用户的浏览器 cookies 里面的,配置文件对应的参数名是：sessionname。
	beego.BConfig.WebConfig.Session.SessionGCMaxLifetime = 3600 //设置 Session 过期的时间，默认值是 3600 秒，配置文件对应的参数：sessiongcmaxlifetime
	//beego.BConfig.WebConfig.Session.SessionCookieLifeTime = 0   //Cookie过期的时间，0为浏览器周期
	models.EngineCfgMsg.Sys.Lang = beego.AppConfig.String("lang")
	if len(models.EngineCfgMsg.Sys.Lang) == 0 {
		models.EngineCfgMsg.Sys.Lang = "zh-CN"
	}
}
