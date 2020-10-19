package controllers

import (
	"time"

	"github.com/bkzy-wangjp/MicEngine/models"
)

type DataMonitorController struct {
	MyController //beego.Controller
}

/***********************************************
功能:计算
输入:无
输出:无
说明:
编辑:wang_jp
时间:2020年3月18日
************************************************/
func (c *DataMonitorController) PageCalculate() {

	pagename := "calculate"
	c.Data["JsFileName"] = pagename  //javascript脚本文件名
	c.Data["ModalSize"] = "modal-lg" //模式框大小:大=modal-lg,小=modal-sm,中号省略
	//c.Data["StartTime"] = models.EngineCfgMsg.Calculate.StartTime
	//c.Data["BeginTime"] = models.EngineCfgMsg.Calculate.BeginTime
	//c.Data["EndTime"] = models.EngineCfgMsg.Calculate.EndTime
	//c.Data["SiftHors"] = models.EngineCfgMsg.Calculate.SiftHors
	//c.Data["TagId"] = models.EngineCfgMsg.Calculate.TagId
	//c.Data["LifeTime"] = models.EngineCfgMsg.Calculate.LifeTime
	//c.Data["HandlingCapacity"] = models.EngineCfgMsg.Calculate.HandlingCapacity
	//c.Data["TimeRemaining"] = models.EngineCfgMsg.Calculate.TimeRemaining
	//c.Data["SurplusCapacity"] = models.EngineCfgMsg.Calculate.SurplusCapacity
	//c.Data["AlarmFlag"] = models.EngineCfgMsg.Calculate.AlarmFlag
	//c.FormatCalculateTree("TreeNodes", "RootPid") //左侧树状菜单栏
	//c.FormatTreeNode("TreeNodes", "RootPid", false, false) //左侧树状菜单栏

	c.InitPageTemplate(pagename) //载入模板数据

	c.TplName = pagename + ".tpl"
}

/***********************************************
功能:快照数据
输入:无
输出:无
说明:
编辑:wang_jp
时间:2020年3月18日
************************************************/
func (c *DataMonitorController) PageSnapShot() { //快照数据
	pagename := "snapshot"
	c.Data["JsFileName"] = pagename
	c.Data["ModalSize"] = "modal-lg"
	c.FormatTreeNode("TreeNodes", "RootPid", false, false) //左侧树状菜单栏
	c.InitPageTemplate(pagename)                           //载入模板数据

	c.TplName = pagename + ".tpl"
}

/***********************************************
功能:历史数据
输入:无
输出:无
说明:
编辑:wang_jp
时间:2020年3月18日
************************************************/
func (c *DataMonitorController) PageHistory() { //历史数据
	pagename := "history"
	c.Data["JsFileName"] = pagename
	c.InitPageTemplate(pagename) //载入模板数据
	c.Data["EndTime"] = models.TimeFormat(time.Now().Add(-60*time.Second), "2006-01-02T15:04")
	c.Data["BeginTime"] = models.TimeFormat(time.Now().Add(-60*481*time.Second), "2006-01-02T15:04")
	c.FormatTreeNode("TreeNodes", "RootPid", true, false) //左侧树状菜单栏

	c.TplName = pagename + ".tpl"
}

/***********************************************
功能:数据对比
输入:无
输出:无
说明:
编辑:wang_jp
时间:2020年3月18日
************************************************/
func (c *DataMonitorController) PageCompare() { //数据对比
	pagename := "compare"
	c.Data["JsFileName"] = pagename
	c.InitPageTemplate(pagename) //载入模板数据
	c.Data["EndTime"] = models.TimeFormat(time.Now().Add(-60*time.Second), "2006-01-02T15:04")
	c.Data["BeginTime"] = models.TimeFormat(time.Now().Add(-60*481*time.Second), "2006-01-02T15:04")
	c.FormatTreeNode("TreeNodes", "RootPid", true, false) //左侧树状菜单栏

	c.TplName = pagename + ".tpl"
}

/***********************************************
功能:回归分析
输入:无
输出:无
说明:
编辑:wang_jp
时间:2020年3月18日
************************************************/
func (c *DataMonitorController) PageRegression() { //回归分析
	pagename := "regression"
	c.Data["JsFileName"] = pagename
	c.InitPageTemplate(pagename) //载入模板数据
	c.Data["EndTime"] = models.TimeFormat(time.Now().Add(-60*time.Second), "2006-01-02T15:04")
	c.Data["BeginTime"] = models.TimeFormat(time.Now().Add(-60*481*time.Second), "2006-01-02T15:04")
	c.FormatTreeNode("TreeNodes", "RootPid", true, false) //左侧树状菜单栏

	c.TplName = pagename + ".tpl"
}

/***********************************************
功能:数据监控
输入:无
输出:无
说明:
编辑:wang_jp
时间:2020年3月18日
************************************************/
func (c *DataMonitorController) PageMonitor() { //数据监控
	pagename := "monitor"
	c.Data["JsFileName"] = pagename
	c.InitPageTemplate(pagename)                            //载入模板数据
	c.FormatMonitorTree("TreeNodes", "RootPid", "FrameUrl") //左侧树状菜单栏

	c.TplName = pagename + ".tpl"
}

/***********************************************
功能:质检化验
输入:无
输出:无
说明:
编辑:wang_jp
时间:2020年4月9日
************************************************/
func (c *DataMonitorController) PageSampleLab() { //质检化验
	platPicPath := "common/picShow/" //平台图片路径
	pagename := "samplelab"
	c.InitPageTemplate()                                                                               //载入模板数据
	c.Data["JsFileName"] = pagename                                                                    //javascript脚本文件名
	c.Data["ModalSize"] = "modal-lg"                                                                   //模式框大小:大=modal-lg,小=modal-sm,中号省略
	c.Data["EndTime"] = models.TimeFormat(time.Now().Add(-60*time.Second), "2006-01-02T15:04")         //结束时间
	c.Data["BeginTime"] = models.TimeFormat(time.Now().Add(-60*10081*time.Second), "2006-01-02T15:04") //开始时间
	c.Data["PlatPicPath"] = c.GetPlatHost() + platPicPath

	c.FormatSampleTree("TreeNodes", "RootPid") //左侧树状菜单栏

	c.TplName = pagename + ".tpl"
}

/***********************************************
功能:物耗模块
输入:无
输出:无
说明:
编辑:wang_jp
时间:2020年4月12日
************************************************/
func (c *DataMonitorController) PageGoods() {
	platPicPath := "common/picShow/" //平台图片路径
	pagename := "goods"
	c.Data["JsFileName"] = pagename                                                                    //javascript脚本文件名
	c.Data["ModalSize"] = "modal-lg"                                                                   //模式框大小:大=modal-lg,小=modal-sm,中号省略
	c.InitPageTemplate(pagename)                                                                       //载入模板数据
	c.Data["EndTime"] = models.TimeFormat(time.Now().Add(-60*time.Second), "2006-01-02T15:04")         //结束时间
	c.Data["BeginTime"] = models.TimeFormat(time.Now().Add(-60*10081*time.Second), "2006-01-02T15:04") //开始时间
	c.Data["PlatPicPath"] = c.GetPlatHost() + platPicPath

	c.FormatGoodsTree("TreeNodes", "RootPid") //左侧树状菜单栏

	c.TplName = pagename + ".tpl"
}

/***********************************************
功能:巡检模块
输入:无
输出:无
说明:
编辑:wang_jp
时间:2020年4月12日
************************************************/
func (c *DataMonitorController) PagePatrol() {
	platPicPath := "common/picShow/" //平台图片路径
	pagename := "patrol"
	c.Data["JsFileName"] = pagename  //javascript脚本文件名
	c.Data["ModalSize"] = "modal-lg" //模式框大小:大=modal-lg,小=modal-sm,中号省略
	c.InitPageTemplate(pagename)     //载入模板数据

	c.Data["EndTime"] = models.TimeFormat(time.Now().Add(-60*time.Second), "2006-01-02T15:04")         //结束时间
	c.Data["BeginTime"] = models.TimeFormat(time.Now().Add(-60*10081*time.Second), "2006-01-02T15:04") //开始时间
	c.Data["PlatPicPath"] = c.GetPlatHost() + platPicPath

	c.FormatPatrolTree("TreeNodes", "RootPid") //左侧树状菜单栏

	c.TplName = pagename + ".tpl"
}

/***********************************************
功能:Kpi模块
输入:无
输出:无
说明:
编辑:wang_jp
时间:2020年4月12日
************************************************/
func (c *DataMonitorController) PageKpi() {
	platPicPath := "common/picShow/" //平台图片路径
	pagename := "kpi"
	c.Data["JsFileName"] = pagename  //javascript脚本文件名
	c.Data["ModalSize"] = "modal-lg" //模式框大小:大=modal-lg,小=modal-sm,中号省略
	c.InitPageTemplate(pagename)     //载入模板数据

	c.Data["EndTime"] = models.TimeFormat(time.Now().Add(-60*time.Second), "2006-01-02T15:04")         //结束时间
	c.Data["BeginTime"] = models.TimeFormat(time.Now().Add(-60*10081*time.Second), "2006-01-02T15:04") //开始时间
	c.Data["PlatPicPath"] = c.GetPlatHost() + platPicPath

	c.FormatKpiTree("TreeNodes", "RootPid") //左侧树状菜单栏

	c.TplName = pagename + ".tpl"
}

/***********************************************
功能:报表模块
输入:无
输出:无
说明:
编辑:wang_jp
时间:2020年4月12日
************************************************/
func (c *DataMonitorController) PageReport() {
	platPicPath := "common/picShow/" //平台图片路径
	pagename := "report"
	c.Data["JsFileName"] = pagename  //javascript脚本文件名
	c.Data["ModalSize"] = "modal-lg" //模式框大小:大=modal-lg,小=modal-sm,中号省略
	c.InitPageTemplate(pagename)     //载入模板数据

	c.Data["PlatPicPath"] = c.GetPlatHost() + platPicPath
	c.Data["EndTime"] = models.TimeFormat(time.Now().Add(-60*time.Second), "2006-01-02T15:04")         //结束时间
	c.Data["BeginTime"] = models.TimeFormat(time.Now().Add(-60*10081*time.Second), "2006-01-02T15:04") //开始时间

	c.FormatReportTree("TreeNodes", "RootPid") //左侧树状菜单栏

	c.TplName = pagename + ".tpl"
}

/***********************************************
功能:报表管理模块
输入:无
输出:无
说明:
编辑:wang_jp
时间:2020年5月8日
************************************************/
func (c *DataMonitorController) PageReportEdit() {
	platPicPath := "common/picShow/" //平台图片路径
	pagename := "reportedit"
	c.Data["JsFileName"] = pagename  //javascript脚本文件名
	c.Data["ModalSize"] = "modal-lg" //模式框大小:大=modal-lg,小=modal-sm,中号省略
	c.InitPageTemplate(pagename)     //载入模板数据

	c.Data["PlatPicPath"] = c.GetPlatHost() + platPicPath

	c.FormatReportTree("TreeNodes", "RootPid") //左侧树状菜单栏

	c.TplName = pagename + ".tpl"
}

/***********************************************
功能:项目管理
输入:无
输出:无
说明:
编辑:wang_jp
时间:2020年4月19日
************************************************/
func (c *DataMonitorController) PageManagerProject() {
	pagename := "manproject"

	c.Data["JsFileName"] = pagename  //javascript脚本文件名
	c.Data["ModalSize"] = "modal-lg" //模式框大小:大=modal-lg,小=modal-sm,中号省略

	c.Data["ProjectLogo"] = models.EngineCfgMsg.CfgMsg.LogoPath
	c.Data["ProjectIcon"] = models.EngineCfgMsg.CfgMsg.IcoPath
	c.Data["Copyright"] = models.EngineCfgMsg.CfgMsg.Copyright
	c.Data["PlatPath"] = c.GetPlatHost()

	c.Data["Version"] = models.EngineCfgMsg.CfgMsg.Version
	c.Data["TotalRunTime"] = models.EngineCfgMsg.CfgMsg.TotalMinutes
	c.Data["ThisRunTime"] = models.EngineCfgMsg.CfgMsg.RunMinutes

	c.Data["MachineCode"] = models.EngineCfgMsg.CfgMsg.MachineCode
	c.Data["AuthCode"] = models.EngineCfgMsg.CfgMsg.AuthCode
	c.Data["KpiAuth"] = models.EngineCfgMsg.CfgMsg.SerialIndicatorAuth
	c.Data["FastAuth"] = models.EngineCfgMsg.CfgMsg.ParallelIndicatorAuth
	c.Data["ReportAuth"] = models.EngineCfgMsg.CfgMsg.ReportAuth
	c.Data["DistributedId"] = models.EngineCfgMsg.CfgMsg.DistributedId
	c.Data["Description"] = models.EngineCfgMsg.CfgMsg.Description

	c.Data["RtdbServer"] = models.EngineCfgMsg.CfgMsg.RtdbServer
	c.Data["RtdbPort"] = models.EngineCfgMsg.CfgMsg.RtdbPort
	c.Data["RtdbName"] = models.EngineCfgMsg.CfgMsg.RtdbDbname
	c.Data["RtdbTableName"] = models.EngineCfgMsg.CfgMsg.RtdbTbname
	c.Data["RtdbUser"] = models.EngineCfgMsg.CfgMsg.RtdbUser
	c.Data["RtdbPsw"] = models.EngineCfgMsg.CfgMsg.RtdbPsw

	c.Data["ResultDbServer"] = models.EngineCfgMsg.CfgMsg.ResultdbServer
	c.Data["ResultDbPort"] = models.EngineCfgMsg.CfgMsg.ResultdbPort
	c.Data["ResultDbName"] = models.EngineCfgMsg.CfgMsg.ResultdbDbname
	c.Data["ResultDbTableName"] = models.EngineCfgMsg.CfgMsg.ResultdbTbname
	c.Data["ResultDbUser"] = models.EngineCfgMsg.CfgMsg.ResultdbUser
	c.Data["ResultDbPsw"] = models.EngineCfgMsg.CfgMsg.ResultdbPsw
	c.Data["DogStatus"] = models.EngineCfgMsg.DogStatus
	c.InitPageTemplate(pagename) //载入模板数据
	c.TplName = pagename + ".tpl"
}

/***********************************************
功能:日志管理
输入:无
输出:无
说明:
编辑:wang_jp
时间:2020年4月19日
************************************************/
func (c *DataMonitorController) PageManagerLog() {
	pagename := "manlog"
	c.Data["JsFileName"] = pagename  //javascript脚本文件名
	c.InitPageTemplate(pagename)     //载入模板数据
	c.Data["ModalSize"] = "modal-lg" //模式框大小:大=modal-lg,小=modal-sm,中号省略

	c.Data["EndTime"] = models.TimeFormat(time.Now(), "2006-01-02T15:04")                             //结束时间
	c.Data["BeginTime"] = models.TimeFormat(time.Now().Add(-60*1440*time.Second), "2006-01-02T15:04") //开始时间

	c.TplName = pagename + ".tpl"

}

/***********************************************
功能:用户管理
输入:无
输出:无
说明:
编辑:wang_jp
时间:2020年4月19日
************************************************/
func (c *DataMonitorController) PageManagerUers() {
	pagename := "biuding"
	c.InitPageTemplate()             //载入模板数据
	c.Data["JsFileName"] = pagename  //javascript脚本文件名
	c.Data["ModalSize"] = "modal-lg" //模式框大小:大=modal-lg,小=modal-sm,中号省略

	c.FormatSampleTree("TreeNodes", "RootPid") //左侧树状菜单栏

	c.TplName = pagename + ".tpl"
}

/***********************************************
功能:权限管理
输入:无
输出:无
说明:
编辑:wang_jp
时间:2020年4月19日
************************************************/
func (c *DataMonitorController) PageManagerPermission() {
	pagename := "biuding"
	c.InitPageTemplate()             //载入模板数据
	c.Data["JsFileName"] = pagename  //javascript脚本文件名
	c.Data["ModalSize"] = "modal-lg" //模式框大小:大=modal-lg,小=modal-sm,中号省略

	c.FormatSampleTree("TreeNodes", "RootPid") //左侧树状菜单栏

	c.TplName = pagename + ".tpl"
}
