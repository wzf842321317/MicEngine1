package models

import (
	"os/exec"
	"time"

	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
)

/***************************************************
功能:每分钟保存心跳信号
输入:无
输出:无
说明:计时器由调用者提供
编辑:wang_jp
时间:2020年8月31日
****************************************************/
func (hb *HeartBeat) SaveHeartBeatPerMinute() {
	o := orm.NewOrm() //新建orm对象
	o.Using(_SqliteAlias)

	hb.RunMinutes = EngineCfgMsg.CfgMsg.RunMinutes
	hb.TotalMinutes = EngineCfgMsg.CfgMsg.TotalMinutes
	hb.DataTime = time.Now().Format(_TIMEFOMAT)
	_, err := o.Insert(hb)
	if err != nil {
		logs.Warning("Update Local Engine runtime error[更新计算引擎本地运行时间失败]")
		logs.Warning(err)
	}
}

/***************************************************
功能:获取看门狗状态
输入:无
输出:无
说明:读取最近3分钟的心跳信息数据,如果看门狗状态大于0，说明看门狗在运行
编辑:wang_jp
时间:2020年9月22日
****************************************************/
func (hb *HeartBeat) GetDogStatus() {
	o := orm.NewOrm() //新建orm对象
	o.Using(_SqliteAlias)

	var hbs []HeartBeat
	nowt := time.Now().Add(-3 * time.Minute)
	qt := o.QueryTable("HeartBeat")
	qt.Filter("DataTime__gt", nowt.Format(_TIMEFOMAT)).All(&hbs)

	EngineCfgMsg.DogStatus = false
	for _, heart := range hbs {
		if heart.DogChecked > 0 {
			EngineCfgMsg.DogStatus = true
			break
		}
	}
	//如果看门狗没有启动，重启之
	if EngineCfgMsg.DogStatus == false {
		hb.ResetDog()
	}
}

/***************************************************
功能:重启看门狗
输入:无
输出:无
说明:
编辑:wang_jp
时间:2020年10月13日
****************************************************/
func (hb *HeartBeat) ResetDog() {
	if len(EngineCfgMsg.Sys.ResetDogFilePath) > 0 { //配置了重启看门狗bat文件路径
		//看门狗未启动,且计算服务已经启动超过3分钟
		if EngineCfgMsg.DogStatus == false && EngineCfgMsg.CfgMsg.RunMinutes >= 3 {
			logs.Info("Restart watchdog...[重启看门狗]")
			//重启看门狗
			c := exec.Command(EngineCfgMsg.Sys.ResetDogFilePath)
			if err := c.Run(); err != nil {
				logs.Info("Failed to restart watchdog[重启看门狗失败]")
			}
		}
	}
}

/***************************************************
功能:删除过期数据
输入:无
输出:无
说明:删除过期数据
编辑:wang_jp
时间:2020年8月31日
****************************************************/
func (hb *HeartBeat) DeleteOutTimeData() {
	o := orm.NewOrm()     //新建orm对象
	o.Using(_SqliteAlias) //根据别名选择数据库
	if EngineCfgMsg.Sys.SaveTimeInComTable > 0 {
		tm := time.Now().Add(-24 * time.Duration(EngineCfgMsg.Sys.SaveTimeInComTable) * time.Hour)
		o.QueryTable(hb).Filter("DataTime__lt", tm.Format(EngineCfgMsg.Sys.TimeFormat)).Delete()
	}
}
