package main

import (
	"sync"
	"time"

	"github.com/bkzy-wangjp/MicEngine/models"

	"github.com/astaxie/beego/logs"
)

//定义任务队列
var _WaitSerial sync.WaitGroup

func SerialRun() {
	logs.Info("周期计算功能已开启……")
	time.Sleep(time.Second * time.Duration(models.EngineCfgMsg.Sys.SleepTimeOnStart)) //延迟执行
	for {
		serialCalc()
		time.Sleep(time.Second * time.Duration(models.EngineCfgMsg.CfgMsg.SerialCalcInterval)) //执行间隔
	}
}

func serialCalc() {

	start := time.Now()
	logs.Info("准备开启新一轮的长周期KPI指标计算……")
	var kpicnt int
	var calcSerialCnt int64 = 0
	cfg := new(models.CalcKpiConfigList)
	if serialCnt, kpis, err := cfg.GetKpiConfigInfo(1); err != nil { //获取串行指标数量
		logs.Error("获取串行KPI计算指标信息时发生错误,错误信息为:", err)
	} else {
		calcSerialCnt = serialCnt                                       //获取本次计算的总数量
		if serialCnt > models.EngineCfgMsg.CfgMsg.SerialIndicatorAuth { //如果总的指标数量大于授权数量
			calcSerialCnt = models.EngineCfgMsg.CfgMsg.SerialIndicatorAuth //已经授权数量为本次计算的总数量
		}
		if models.EngineCfgMsg.CfgMsg.SerialIndicNumPerThread == 0 { //每Go程计算的指标数量不能为0
			models.EngineCfgMsg.CfgMsg.SerialIndicNumPerThread = 600 //入股为0设定为默认值
		}
		kpiNumPerGo := models.EngineCfgMsg.CfgMsg.SerialIndicNumPerThread //记录每Go程指标数量
		serialGoCnt := calcSerialCnt / kpiNumPerGo                        //计算需要的Go程数
		var i int64
		logs.Info("共查询到[%d]个长周期KPI指标,授权计算[%d]个,实际计算[%d]个,每个线程[%d]个指标,共划分了[%d]个线程",
			serialCnt, models.EngineCfgMsg.CfgMsg.SerialIndicatorAuth, calcSerialCnt, kpiNumPerGo, serialGoCnt+1)

		resKpi := make(chan int, serialGoCnt+1) //结果集chan
		defer func() {
			if err := recover(); err != nil {
				logs.Critical("serialCalc时遇到错误:[%#v]", err)
			}
			defer close(resKpi) //压后关闭chan道
		}()

		for i = 0; i <= serialGoCnt; i++ { //遍历，生成并发的Go程
			_WaitSerial.Add(1)
			st := kpiNumPerGo * i
			ed := st + kpiNumPerGo
			if ed > serialCnt {
				ed = serialCnt
			}
			kpi := kpis[st:ed]
			go serialKpiCalc(kpi, resKpi) //开启计算Go程
		}

		for i = 0; i <= serialGoCnt; i++ {
			k := <-resKpi //接收计算结果
			kpicnt += k
		}
		_WaitSerial.Wait()
	}
	end := time.Now()
	//输出执行时间，单位为秒。
	logs.Info("本轮KPI计算开始于[%s],执行耗时[%f]秒.共扫描[%d]条数据,生成[%d]条KPI数据",
		start.Format(models.EngineCfgMsg.Sys.TimeFormat), end.Sub(start).Seconds(), calcSerialCnt, kpicnt)
}

/*
功能:串行KPI指标计算GO程程序
输入:
	kpis:指标配置信息数组
	resKpi:KPI计算结果程道数组
输出:无
说明:
编辑:wang_jp
时间:2019年12月13日
*/
func serialKpiCalc(kpis []models.CalcKpiConfigListExi, resKpi chan int) {
	var kpi_result []models.CalcKpiResult
	defer func() {
		if err := recover(); err != nil {
			logs.Critical("serialKpiCalc时遇到错误:%d~%d;[%#v]", kpis[0].Id, kpis[len(kpis)-1].Id, err)
		}
		resKpi <- len(kpi_result)
		_WaitSerial.Done()
	}()

	for _, kpi := range kpis {
		if iscalc, bgTime, endTime, err := kpi.PrevTime(); err != nil {
			logs.Error(err.Error())
		} else {
			if iscalc == true { //到了计算时间点，允许计算
				if models.EngineCfgMsg.Sys.Debug {
					logs.Debug("执行计算:%d", kpi.Id) //==========================================
				}
				res, ctinue, err := kpi.KpiCalc(bgTime, endTime) //执行计算
				if err != nil {                                  //计算结果有错误
					if ctinue {
						logs.Warn(err.Error())
					} else {
						logs.Error(err.Error())
					}
				} else { //没有错误
					if res.KpiConfigListId != 0 { //返回了计算结果
						kpi_result = append(kpi_result, res)
					}
				}
			}
		}
		if models.EngineCfgMsg.Sys.Debug {
			logs.Debug("计算完成:%d", kpi.Id) //==========================================
		}
	}

	go func(result []models.CalcKpiResult) {
		defer func() {
			if err := recover(); err != nil {
				logs.Critical("批量保存周期计算结果集时遇到错误:[%#v]", err)
			}
		}()
		r := new(models.CalcKpiResult)
		r.SaveBatchKpiResultToDB(result) //批量保存结果集
	}(kpi_result)

	go func(result []models.CalcKpiResult) {
		defer func() {
			if err := recover(); err != nil {
				logs.Critical("批量保存KPI最后计算时间时遇到错误:[%#v]", err)
			}
		}()
		r := new(models.CalcKpiConfigList)
		r.SetKpiLastCalcTime(result)
	}(kpi_result)
}
