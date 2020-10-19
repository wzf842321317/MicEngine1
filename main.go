package main

import (
	"runtime"
	"time"

	"github.com/astaxie/beego/plugins/cors"

	"github.com/astaxie/beego"
	"github.com/bkzy-wangjp/MicEngine/models"
	_ "github.com/bkzy-wangjp/MicEngine/routers"
)

func main() {
	defer func() {
		if err := recover(); err != nil {
			beego.Critical("main 捕获到错误:[%#v]", err)
		} else {
			beego.Info("MicEngine主程序退出运行")
		}
		models.GDDB.GoldenDB.DisConnect(models.GDPOOL)
		models.GDPOOL.RemovePool()
	}()
	if models.EngineCfgMsg.Sys.PeriodKpiEnable == true {
		go SerialRun()
	} else {
		beego.Info("周期性KPI指标计算功能已关闭，如要开启请设置配置文件中的 PeriodKpiEnable=true，并重启MicEngine")
	}
	if models.EngineCfgMsg.Sys.FastKpiEnable == true {
		go ParallelRun()
	} else {
		beego.Info("快速KPI指标计算功能已关闭，如要开启请设置配置文件中的 FastKpiEnable=true，并重启MicEngine")
	}
	if models.EngineCfgMsg.Sys.ReportEnable == true {
		go ReportRun()
	} else {
		beego.Info("报表计算功能已关闭，如要开启请设置配置文件中的 ReportEnable=true，并重启MicEngine")
	}
	if models.EngineCfgMsg.Sys.AlarmEnable == true {
		go AlarmRun()
	} else {
		beego.Info("报警计算功能已关闭，如要开启请设置配置文件中的 AlarmEnable=true，并重启MicEngine")
	}

	if models.EngineCfgMsg.Sys.WebEnable == true { //开启了Web功能
		go hertbeat()
		beego.Info("Web服务功能已开启……")
		//允许跨域请求
		beego.InsertFilter("*", beego.BeforeRouter, cors.Allow(&cors.Options{
			AllowOrigins:     []string{"*"},           //允许访问所有源，也可以设置url过滤
			AllowMethods:     []string{"GET", "POST"}, //可选参数"GET", "POST", "PUT", "DELETE", "OPTIONS" (*为所有)
			AllowHeaders:     []string{"*"},           //指的是允许的Header的种类
			ExposeHeaders:    []string{"*"},           //公开的HTTP标头列表
			AllowCredentials: true}),                  //如果设置，则允许共享身份验证凭据，例如cookie
		)
		beego.Run() //启动beego
	} else {
		beego.Info("Web功能已关闭，如要开启请设置配置文件中的 WebEnable=true，并重启MicEngine")
		hertbeat()
	}
}

//心跳信息
func hertbeat() {
	for {
		beat := new(models.HeartBeat)
		beat.GetDogStatus()
		models.EngineCfgMsg.CfgMsg.SaveHeartBeatPerMinute()
		beat.SaveHeartBeatPerMinute()

		if models.EngineCfgMsg.Sys.Debug && models.GDPOOL != nil {
			beego.Debug("当前庚顿连接数:[", len(models.GDPOOL.Worker), "] 连接请求数:[", len(models.GDPOOL.Req), "]")
		}
		beego.Debug("程序已运行[", models.EngineCfgMsg.CfgMsg.RunMinutes, "]分钟,当前共有[", runtime.NumGoroutine(), "]个Go程在运行")
		if models.EngineCfgMsg.CfgMsg.RunMinutes%60 == 0 {
			result := new(models.CalcKpiResult)
			result.DeleteOutTimeResultDataInComTable() //每小时清除一次过期数据
			beat.DeleteOutTimeData()
		}
		time.Sleep(time.Second * 60)
	}
}
