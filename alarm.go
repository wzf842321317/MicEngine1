package main

import (
	"time"

	"github.com/astaxie/beego/logs"
	"github.com/bkzy-wangjp/MicEngine/models"
)

/*******************************************************************************
功能:报警计算程序
输入:无
输出:无
说明:
编辑:wang_jp
时间:2020年6月9日
*******************************************************************************/
func AlarmRun() {
	logs.Info("报警计算功能已开启……")
	bg := time.Now()
	alarmCnt, err := models.GDDB.GetGoldenPoints()
	if err != nil {
		logs.Emergency("报警模块获取庚顿标签信息失败:[%s]", err.Error())
	}
	logs.Info("加载报警基本信息完成，共[%d]个标签点需要进行报警计算，加载耗时:%f秒", alarmCnt, time.Since(bg).Seconds())
	for {
		bg := time.Now()
		err := models.GDDB.GetGoldenSnapShotsAndComp()
		if err != nil {
			logs.Emergency("报警模块获取快照数据时错误:[%s]", err.Error())
		} else {
			if models.EngineCfgMsg.Sys.Debug {
				logs.Debug("完成一次报警扫描,耗时%f秒", time.Since(bg).Seconds())
			}
		}
		if time.Since(models.GDDB.LoadedTime).Seconds() > float64(models.EngineCfgMsg.CfgMsg.AlarmMsgReloadInterval) {
			bg := time.Now()
			alarmCnt, err := models.GDDB.GetGoldenPoints()
			if err != nil {
				logs.Emergency("报警模块获取庚顿标签信息失败:[%s]", err.Error())
			}
			logs.Debug("重载报警基本信息完成，共[%d]个标签点需要进行报警计算，加载耗时:[%f]秒", alarmCnt, time.Since(bg).Seconds())
		} else {
			time.Sleep(time.Second * time.Duration(models.EngineCfgMsg.CfgMsg.AlarmScanInterval)) //执行间隔
		}

	}
}
