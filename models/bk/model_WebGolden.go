package models

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/bkzy-wangjp/MicEngine/MicScript/statistic"
	"github.com/bkzy-wangjp/MicEngine/models/Golden"
)

/****************************************************
功能:获取庚顿数据库中的单个数据
输入:
	tagname:变量名
	key:获取数据的指令，不区分大小写,指令列表如下
		"point"		-获取指定时间点的值,此时beginetime和endtime无效
		"snapshot"	-获取当前快照值,此时beginetime和endtime无效
		"min" 		-获取指定时间范围内的最小值
		"max"		-获取指定时间范围内的最大值
		"range"		-获取指定范围内的数据范围(最大值-最小值)
		"total"		-获取指定范围内的数据轴与时间轴(秒)围成的面积
		"sum"		-获取指定范围内的算术累积值
		"mean"		-获取指定范围内的算术平均值
		"poweravg"	-获取指定范围内的加权平均值(total/duration)
		"diff"		-获取指定范围内的最新值与最老值的差值
		"duration","dur","times"	-获取指定的时间持续的秒数
		"pointcnt", "count","pcnt"	-获取指定范围内的历史数据点数
		"risingcnt", "rising", "rising_counts","rcnt"	-获取指定范围内的数据上升次数
		"fallingcnt", "falling", "falling_counts","fcnt"-获取指定范围内的数据下降次数
		"sd"		-获取指定范围内的总体标准偏差
		"stddev"	-获取指定范围内的样本标准变差
		"se"		-获取指定范围内的标准误差
		"ske"		-获取指定范围内的数据偏度
		"kur"		-获取指定范围内的数据峰度
		"mode"		-获取指定范围内的数据的众数
		"median"	-获取指定范围内的数据的中位数
		"advangce"	-获取上述除快照外的所有数据,此函数中不可用
	begineTime:开始时间,读取快照时无用
	endTime:结束时间,读取快照时无用
输出：查询到的数据
编辑：wang_jp
时间：2019年12月8日
****************************************************/
func GetGoldenSingleData(tagname, key, begineTime, endTime string) (float64, error) {
	add := new(Golden.GoldenAdd)
	add.Host = EngineCfgMsg.CfgMsg.RtdbServer
	add.Port = int(EngineCfgMsg.CfgMsg.RtdbPort)
	var data float64 = 0.0

	switch strings.ToLower(key) {
	case "snapshot":
		tagvalue, _, _, err := GetGoldenSnapShotData(tagname)
		if err != nil {
			return data, err
		}
		data = tagvalue
	case "advangce": //advangce在此函数中无效
		data = 0.0
	case "point": //获取指定时间点的值
		cmd := new(Golden.HistoryCmd)              //新建历史统计数据读取指令
		cmd.TagName = append(cmd.TagName, tagname) //变量名
		cmd.TimePoint = endTime                    //结束时间
		vmap, err := Golden.GetGoldenData(*add, *cmd)
		if err != nil {
			return data, err
		}
		value := vmap.(Golden.HisPointDataMap)
		if Golden.IsExistItem(tagname, value) {
			data, err = strconv.ParseFloat(value[tagname].Value, 64)
			if err != nil {
				return data, err
			}
		}
	default:
		cmd := new(Golden.HisSumCmd)               //新建历史统计数据读取指令
		cmd.TagName = append(cmd.TagName, tagname) //变量名
		cmd.BeginTime = begineTime                 //起始时间
		cmd.EndTime = endTime                      //结束时间
		cmd.DataType = key                         //关键词
		cmd.Interval = 0                           //间隔单位,秒。如果为负数，则读取庚顿原生统计数据;如果为0，读取原始历史数据的统计数据,如果大于0，则读取等间隔历史数据的统计数据
		cmd.Group = 0

		vmap, err := Golden.GetGoldenData(*add, *cmd)
		if err != nil {
			return data, err
		}
		value := vmap.(map[string]interface{})
		if Golden.IsExistItem(tagname, value) {
			data = value[tagname].(float64)
		}
	}
	return data, nil
}

/****************************************************
功能:对庚顿历史数据进行与限制值进行比较，输出比较结果的累积值
输入:
	hisdata:历史数据map,key:比较指令,limit:比较值
输出：比较结果累积值
编辑：wang_jp
时间：2020年1月3日
****************************************************/
func GetGoldenAnalogComparedTotal(tagid int64, tagname, key, begineTime, endTime string) (float64, error) {
	add := new(Golden.GoldenAdd)
	add.Host = EngineCfgMsg.CfgMsg.RtdbServer
	add.Port = int(EngineCfgMsg.CfgMsg.RtdbPort)

	cmd := new(Golden.HisSumCmd)               //新建历史统计数据读取指令
	cmd.TagName = append(cmd.TagName, tagname) //变量名
	cmd.BeginTime = begineTime                 //起始时间
	cmd.EndTime = endTime                      //结束时间
	vmap, err := Golden.GetHistoryDataAlignHeadAndTail(*add, *cmd)
	if err != nil {
		return 0, err
	}
	tag := new(OreProcessDTaglist)
	tag.Id = tagid
	tag.GetTagById()
	var limit float64
	switch key {
	case "lt_l_total", "lte_l_total", "gt_l_total", "gte_l_total":
		limit = tag.LimitL
	case "lt_ll_total", "lte_ll_total", "gt_ll_total", "gte_ll_total":
		limit = tag.LimitLl
	case "lt_h_total", "lte_h_total", "gt_h_total", "gte_h_total":
		limit = tag.LimitH
	case "lt_hh_total", "lte_hh_total", "gt_hh_total", "gte_hh_total":
		limit = tag.LimitHh
	}

	return analogComparedTotal(vmap, key, limit)
}

/****************************************************
功能:对庚顿历史数据进行与限制值进行比较，输出比较结果的累积值
输入:
	hisdata:历史数据map,key:比较指令,limit:比较值
输出：比较结果累积值
编辑：wang_jp
时间：2020年1月3日
****************************************************/
func analogComparedTotal(hisdata Golden.HisDataMap, key string, limit float64) (float64, error) {
	var v, v0, total float64
	var t0 time.Time
	for _, hisd := range hisdata {
		for k, dv := range hisd {
			t, err := TimeParse(dv.Time)
			if err != nil {
				return 0.0, err
			}
			if k > 0 { //非第一个数
				f, err := strconv.ParseFloat(dv.Value, 64)
				if err != nil {
					return 0.0, err
				}
				if compare(key, f, limit) {
					v = 1.0
				} else {
					v = 0.0
				}
				total += (v0 + v) / 2 * float64(t.Unix()-t0.Unix())
				v0 = v
			} else { //第一个数
				f, err := strconv.ParseFloat(dv.Value, 64)
				if err != nil {
					return 0.0, err
				}
				if compare(key, f, limit) {
					v0 = 1.0
				} else {
					v0 = 0.0
				}
			}
			t0 = t
		}
	}
	return total, nil
}

/****************************************************
功能:根据指令比较输入的比较值和被比较值
输入:
	key:比较指令,value:被比较值,limit:比较值
输出：比较结果
编辑：wang_jp
时间：2020年1月3日
****************************************************/
func compare(key string, value, limit float64) bool {
	switch key {
	case "lt_l_total", "lt_ll_total", "lt_h_total", "lt_hh_total":
		return value < limit
	case "lte_l_total", "lte_ll_total", "lte_h_total", "lte_hh_total":
		return value <= limit
	case "gt_l_total", "gt_ll_total", "gt_h_total", "gt_hh_total":
		return value > limit
	case "gte_l_total", "gte_ll_total", "gte_h_total", "gte_hh_total":
		return value > limit
	default:
		return false
	}
}

/****************************************************
功能:获取庚顿数据库中指定变量的快照
输入:
	tagname:变量名
输出：数值,时间,数据质量,错误信息
编辑：wang_jp
时间：2019年12月8日
****************************************************/
func GetGoldenSnapShotData(tagname string) (float64, string, int, error) {
	add := new(Golden.GoldenAdd) //新建庚顿地址
	add.Host = EngineCfgMsg.CfgMsg.RtdbServer
	add.Port = EngineCfgMsg.CfgMsg.RtdbPort
	var datatime string
	var quality int
	var datavalue float64
	var err error

	cmd := new(Golden.SnapshotCmd) //新建读取快照指令
	cmd.TagName = append(cmd.TagName, tagname)
	vmap, err1 := Golden.GetGoldenData(*add, *cmd) //读取快照
	if err1 != nil {
		return datavalue, datatime, quality, err
	}
	value := vmap.(Golden.SnapDataMap)
	if len(value) > 0 { //读取快照成功
		datatime = value[0].Time
		datavalue, err = strconv.ParseFloat(value[0].Value, 64)
		quality = value[0].Quality
	} else { //读取快照失败
		err = errors.New(fmt.Sprintf("The snapshot value of tag %s was not found.", tagname))
	}
	return datavalue, datatime, quality, err
}

/****************************************************
功能:获取庚顿数据库中指定变量的快照结构
输入:
	tagname:变量名
输出：数值,时间,数据质量,错误信息
编辑：wang_jp
时间：2019年12月8日
****************************************************/
func GetGoldenSnapShotMap(tagname string) (Golden.SnapDataMap, error) {
	add := new(Golden.GoldenAdd) //新建庚顿地址
	add.Host = EngineCfgMsg.CfgMsg.RtdbServer
	add.Port = EngineCfgMsg.CfgMsg.RtdbPort

	cmd := new(Golden.SnapshotCmd) //新建读取快照指令
	cmd.TagName = strings.Split(tagname, ",")
	snap, err := Golden.GetGoldenData(*add, *cmd) //读取快照
	return snap.(Golden.SnapDataMap), err
}

/****************************************************
功能:获取庚顿数据库中的数据接口
输入:
	cmd:接口命令,具体参见Golden模块
输出：数据接口
时间：2019年11月28日
编辑：wang_jp
****************************************************/
func GetGoldenData(cmd interface{}) (interface{}, error) {
	add := new(Golden.GoldenAdd) //新建庚顿地址
	add.Host = EngineCfgMsg.CfgMsg.RtdbServer
	add.Port = EngineCfgMsg.CfgMsg.RtdbPort
	return Golden.GetGoldenData(*add, cmd)
}

/****************************************************
功能:获取庚顿历史数据库中历史数据的峰谷值统计
输入:
	tagid,tagname,key,beginTime,endTime
输出：数据接口,err
时间：2020年1月19日
编辑：wang_jp
****************************************************/
func GetPeakValleyOfGoldenHisData(tagid int64, tagname, key, beginTime, endTime string) (interface{}, error) {
	fmt.Println("GetPeakValleyOfGoldenHisData的参数:", tagid, tagname, key, beginTime, endTime)
	tag := new(OreProcessDTaglist)
	tag.Id = tagid
	if tagid == 0 { //没有设置ID
		id, err := tag.GetTagIDByFullName(tagname)
		if err != nil {
			fmt.Println("获取ID失败", err.Error())
			return 0.0, err
		}
		tagid = id
	} else { //设置了ID
		fullname, err := tag.GetTagFullNameByID()
		if err != nil {
			return 0.0, err
		}
		tagname = fullname
	}
	err := tag.GetTagById()
	if err != nil {
		return 0.0, err
	}

	stdv := tag.UserReal1          //user_real1作为稳态判据
	increment := tag.UserReal2     //user_real2作为拐点增量值
	cp := tag.UserInt1             //user_int1作为最大连续点数
	negativeAsZero := tag.UserInt2 //user_int2如果为0，保留负数;如果为1,将负数作为0处理

	cmd := new(Golden.HistoryCmd)
	cmd.BeginTime = beginTime
	cmd.EndTime = endTime
	cmd.TagName = append(cmd.TagName, tagname)
	goldendata, err := GetGoldenData(*cmd)
	if err != nil {
		return 0.0, err
	}
	//fmt.Println(goldendata)
	hisdata := goldendata.(Golden.HisDataMap)
	pvs := new(statistic.PeakValleySelector)
	pvs.New(stdv, int(cp))
	var tsds []statistic.TimeSeriesData
	for _, his := range hisdata {
		var tsd statistic.TimeSeriesData
		for _, data := range his {
			tsd.Time = data.Time
			v, err := strconv.ParseFloat(data.Value, 64)
			if err != nil {
				return 0.0, err
			}
			tsd.Value = v
			tsds = append(tsds, tsd)
		}
	}
	pvs.DataFillter(tsds, increment, int(negativeAsZero))
	switch strings.ToLower(key) {
	case "pvs_peaksum":
		return pvs.PeakSum, nil
	case "pvs_valleysum":
		return pvs.ValleySum, nil
	case "pvs_pvdiffsum":
		return pvs.PVDiffSum, nil
	case "pvs_periodcnt":
		return pvs.PeriodCnt, nil
	case "pvs_peek_valley_datas":
		return pvs.PvDatas, nil
	case "peakvalley":
		return *pvs, nil
	default:
		return 0.0, fmt.Errorf("Invalid instruction;[无效的指令]")
	}

}

/***************************************************
功能:读取庚顿等间隔历史数据
输入:
	tagname:变量全名,beginTime:起始时间,endTime:结束时间
输出：时间序列数据,err
时间：2020年3月19日
编辑：wang_jp
***************************************************/
func GetGoldenHistory(tagname, begintime, endtime string) ([]SrtData, error) {
	var cmd Golden.HistoryCmd //等间隔历史数据命令
	//格式化命令
	cmd.TagName = append(cmd.TagName, tagname)
	cmd.BeginTime = begintime
	cmd.EndTime = endtime
	res, err := GetGoldenData(cmd) //读取数据
	if err != nil {
		return nil, err
	}
	hiss := res.(Golden.HisDataMap) //断言数据类型
	//解析数据
	var srtds []SrtData
	for _, his := range hiss {
		for _, v := range his {
			var srtd SrtData
			srtd.Datatime = v.Time
			f, e := strconv.ParseFloat(v.Value, 64)
			if e == nil {
				srtd.Value = f
			}
			srtds = append(srtds, srtd)
		}
	}
	return srtds, nil
}

/***************************************************
功能:读取庚顿等间隔历史数据
输入:
	tagname:变量全名,beginTime:起始时间,endTime:结束时间,interval:时间间隔,以秒为单位
输出：时间序列数据,err
时间：2020年3月19日
编辑：wang_jp
***************************************************/
func GetGoldenHistoryInterval(tagname, begintime, endtime string, interval int64) ([]SrtData, error) {
	var cmd Golden.HisIntervalCmd //等间隔历史数据命令
	//格式化命令
	cmd.TagName = append(cmd.TagName, tagname)
	cmd.BeginTime = begintime
	cmd.EndTime = endtime
	cmd.Interval = interval
	res, err := GetGoldenData(cmd) //读取数据
	if err != nil {
		return nil, err
	}
	hiss := res.(Golden.HisDataMap) //断言数据类型
	//解析数据
	var srtds []SrtData
	for _, his := range hiss {
		for _, v := range his {
			var srtd SrtData
			srtd.Datatime = v.Time
			f, e := strconv.ParseFloat(v.Value, 64)
			if e == nil {
				srtd.Value = f
			}
			srtds = append(srtds, srtd)
		}
	}
	return srtds, nil
}

/***************************************************
功能:读取庚顿历史统计数据
输入:
	tagname:变量全名,beginTime:起始时间,endTime:结束时间
输出：时间序列数据,err
时间：2020年3月19日
编辑：wang_jp
***************************************************/
func GetGoldenHistorySummary(tagname, begintime, endtime string) (Golden.HisSumDataExi, error) {
	var cmd Golden.HisSumCmd //等间隔历史数据命令
	var hs Golden.HisSumDataExi
	//格式化命令
	cmd.TagName = append(cmd.TagName, tagname)
	cmd.BeginTime = begintime
	cmd.EndTime = endtime
	cmd.Interval = 0
	cmd.Group = 0
	cmd.DataType = "advangce"
	res, err := GetGoldenData(cmd) //读取数据
	if err != nil {
		return hs, err
	}
	//解析数据
	for _, his := range res.(Golden.HisSumDataMapExi) {
		return his, nil
	}
	return hs, nil
}

/***************************************************
功能:写庚顿快照数据
输入:
	datas:写快照的数据结构
输出：时间序列数据,err
时间：2020年3月30日
编辑：wang_jp
***************************************************/
func WriteGoldenSnapShot(datas []GoldenSnapWrite) (interface{}, error) {
	add := new(Golden.GoldenAdd) //数据库地址
	add.Host = EngineCfgMsg.CfgMsg.RtdbServer
	add.Port = EngineCfgMsg.CfgMsg.RtdbPort
	//格式化命令
	var cmd Golden.SnapWriteCmd
	var data Golden.SnapWrite
	data.Time = TimeFormat(time.Now(), _TIMEFOMAT) //当前时间
	for _, v := range datas {
		data.TagName = v.TagName
		data.Value = fmt.Sprintf("%f", v.Value)
		data.Quality = v.Quality
		cmd.Datas = append(cmd.Datas, data)
	}
	return Golden.WriteSnapShotData(*add, cmd)
}
