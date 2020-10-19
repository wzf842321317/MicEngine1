package main

import (
	"sync"
	"time"

	"github.com/bkzy-wangjp/MicEngine/models"

	"github.com/astaxie/beego/logs"
)

//定义任务队列
var _WaitReport sync.WaitGroup

/*******************************************************************************
功能:报表计算程序
输入:无
输出:无
说明:
编辑:wang_jp
时间:2020年5月6日
*******************************************************************************/
func ReportRun() {
	logs.Info("报表计算功能已开启……")
	time.Sleep(time.Second * time.Duration(models.EngineCfgMsg.Sys.SleepTimeOnStart)) //延迟执行
	for {
		reportCalc()
		time.Sleep(time.Second * time.Duration(models.EngineCfgMsg.CfgMsg.ReportCalcInterval)) //执行间隔
	}
}

/*******************************************************************************
功能:报表计算分配程序
输入:无
输出:无
说明:
编辑:wang_jp
时间:2020年5月6日
*******************************************************************************/
func reportCalc() {
	start := time.Now()
	logs.Info("准备开启新一轮的报表计算……")
	var reportcnts int
	var calcReportCnt int64 = 0
	rptl := new(models.CalcKpiReportList)                              //允许计算的报表总数量
	if reportCnt, reports, err := rptl.GetRoportLists(0); err != nil { //获取报表列表
		logs.Error("获取报表列表信息时发生错误,错误信息为:", err)
	} else {
		calcReportCnt = reportCnt                              //获取本次允许计算的总数量
		if reportCnt > models.EngineCfgMsg.CfgMsg.ReportAuth { //如果总的指标数量大于授权数量
			calcReportCnt = models.EngineCfgMsg.CfgMsg.ReportAuth //已经授权数量为本次允许计算的总数量
		}

		reportsNumPerGo := models.EngineCfgMsg.CfgMsg.SerialIndicNumPerThread / 10 //每Go程报表数量,每GO程串行指标数量的十分之一
		if reportsNumPerGo < 10 {                                                  //不小于10
			reportsNumPerGo = 10
		}
		reportGoCnt := calcReportCnt / reportsNumPerGo //计算需要的Go程数
		var i int64
		logs.Info("共查询到[%d]个报表模板,授权计算[%d]个,实际计算[%d]个,每个线程[%d]个报表,共划分了[%d]个线程",
			reportCnt, models.EngineCfgMsg.CfgMsg.ReportAuth, calcReportCnt, reportsNumPerGo, reportGoCnt+1)

		resReport := make(chan int, reportGoCnt+1) //结果集chan
		defer close(resReport)                     //压后关闭chan道
		for i = 0; i <= reportGoCnt; i++ {         //遍历，生成并发的Go程
			_WaitReport.Add(1)
			st := reportsNumPerGo * i
			ed := st + reportsNumPerGo
			if ed > reportCnt {
				ed = reportCnt
			}
			report := reports[st:ed]
			go reportKpiCalc(report, resReport) //开启计算Go程
		}

		for i = 0; i <= reportGoCnt; i++ {
			k := <-resReport //接收计算结果
			reportcnts += k
		}
		_WaitReport.Wait()
	}
	end := time.Now()
	//输出执行时间，单位为秒。
	logs.Info("本轮报表计算开始于[%s],执行耗时[%f]秒.共扫描[%d]个报表模板,生成[%d]个报表",
		start.Format(models.EngineCfgMsg.Sys.TimeFormat), end.Sub(start).Seconds(), calcReportCnt, reportcnts)
}

/*******************************************************************************
功能:报表计算GO程程序
输入:
	reports:报表配置信息数组
	resReport:实际计算的报表数量
输出:无
说明:
编辑:wang_jp
时间:2020年5月6日
*******************************************************************************/
func reportKpiCalc(reports []models.CalcKpiReportList, resReport chan int) {
	defer _WaitReport.Done()
	var calcCnt int = 0 //实际计算的报表数量
	var firstid, lastid int64

	for i, rpt := range reports {
		if i == 0 {
			logs.Debug("报表开始计算，起始ID：%d", rpt.Id)
			firstid = rpt.Id
		}
		if i == len(reports)-1 {
			lastid = rpt.Id
		}
		if iscalc, bgTime, endTime, err := rpt.PrevTime(); err != nil {
			logs.Error("报表时间格式设置错误,报表模板ID是[%d];[%s]", rpt.Id, err.Error())
		} else {
			if iscalc == true { //到了计算时间点，允许计算
				num, fname, err := rpt.ReportCalc(bgTime, endTime)
				if err != nil {
					logs.Error(err.Error())
				} else {
					calcCnt += 1
					logs.Info("生成报表[%s],表中计算了[%d]个KPI指标", fname, num)
				}
			}
		}
	}
	logs.Debug("报表模板 [%d]~[%d]扫描完成", firstid, lastid)
	resReport <- calcCnt
}
