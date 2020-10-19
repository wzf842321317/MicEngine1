package models

import (
	"fmt"
	"math"
	"strings"
	"sync"
	"time"

	"github.com/astaxie/beego/logs"
	"github.com/bkzy-wangjp/MicEngine/MicScript/statistic"
	"github.com/bkzy-wangjp/goldengo"
)

type MicGolden struct {
	goldengo.Golden
}

var GDPOOL *goldengo.GoldenPool //庚顿连接资源池

func (m *MicGolden) NewGoldenPool() (*goldengo.GoldenPool, error) {
	if len(EngineCfgMsg.CfgMsg.RtdbServer) <= 5 || EngineCfgMsg.CfgMsg.RtdbPort == 0 { //没有配置庚顿数据库
		return nil, fmt.Errorf("没有配置实时数据库接口")
	}
	return goldengo.NewGoldenPool(EngineCfgMsg.CfgMsg.RtdbServer,
		EngineCfgMsg.CfgMsg.RtdbUser,
		EngineCfgMsg.CfgMsg.RtdbPsw,
		int(EngineCfgMsg.CfgMsg.RtdbPort),
		EngineCfgMsg.Sys.GoldenCennectPool)
}

/****************************************************
功能:建立与庚顿数据库的网络连接
输入:[tenant ...string] 可选的租借句柄的字符串
输出:庚顿数据库连接指针,错误信息
编辑:wang_jp
时间:2020年5月16日
****************************************************/
func (m *MicGolden) GetHandel(tenant ...string) error {
	if EngineCfgMsg.Sys.Debug {
		logs.Debug("当前在用资源数:[%d],当前请求数:[%d]", len(GDPOOL.Worker), len(GDPOOL.Req))
	}
	if m.Handle > 0 { //已经有连接句柄
		return nil
	}
	return m.GetConnect(GDPOOL, tenant...)
}

/****************************************************
功能:断开与庚顿数据库的网络连接
输入:无
输出:无
编辑:wang_jp
时间:2020年5月16日
****************************************************/
func (m *MicGolden) ReleaseHandel() {
	if m.Handle > 0 { //已经有连接
		m.DisConnect(GDPOOL)
		m.Handle = 0
	}
}

/***************************************************
功能:读取庚顿服务器时间
输入:无
输出:时间,err
时间:2020年3月19日
编辑:wang_jp
***************************************************/
func (m *MicGolden) GoldenGetServerTime() (string, error) {
	err := m.GetHandel("GoldenGetServerTime") //建立到庚顿数据库的连接
	if err != nil {                           //判断连接是否有错误
		return "", err
	}
	defer m.ReleaseHandel() //压后断开连接
	t, err := m.HostTime()
	if err != nil {
		return "", err
	}
	return TimeFormat(Millisecond2Time(int64(t) * 1e3)), nil
}

/***************************************************
功能:读取庚顿服务器版本号
输入:无
输出:版本号,err
时间:2020年3月19日
编辑:wang_jp
***************************************************/
func (m *MicGolden) GoldenGetApiVersion() (string, error) {
	err := m.GetHandel("GoldenGetApiVersion") //建立到庚顿数据库的连接
	if err != nil {                           //判断连接是否有错误
		return "", err
	}
	defer m.ReleaseHandel() //压后断开连接
	return m.GetAPIVersion()
}

/***************************************************
功能:读取庚顿服务器时间
输入:无
输出:[int64] Golden服务器的当前UTC时间，表示距离1970年1月1日08:00:00的秒数
	[error] 错误信息
时间:2020年3月19日
编辑:wang_jp
***************************************************/
func (m *MicGolden) GoldenGetHostTime() (int64, error) {
	err := m.GetHandel("GoldenGetHostTime") //建立到庚顿数据库的连接
	if err != nil {                         //判断连接是否有错误
		return 0, err
	}
	defer m.ReleaseHandel() //压后断开连接
	return m.HostTime()
}

/****************************************************
功能:获取庚顿数据库中指定变量的快照
输入:
	tagname:变量名
输出:数值,时间,数据质量,错误信息
编辑:wang_jp
时间:2019年12月8日
****************************************************/
func (m *MicGolden) GoldenGetSnapShotData(tagname string) (float64, string, int, error) {
	err := m.GetHandel("GoldenGetSnapShotData", tagname) //建立到庚顿数据库的连接
	if err != nil {                                      //判断连接是否有错误
		return 0.0, "", 0, err
	}
	defer m.ReleaseHandel() //压后断开连接
	snames := strings.Split(tagname, ".")
	if len(snames) < 2 {
		return 0.0, "", 0, fmt.Errorf("变量[%s]的名称不合法,合法的格式是:[ tablename.tagname ]", tagname)
	}
	snapmap, err := m.GetSnapShotByName(tagname)
	if err != nil { //有错误
		return 0.0, "", 0, fmt.Errorf("读取变量[%s]快照时发生错误:[%s]", tagname, err.Error())
	}
	if len(snapmap[tagname].Err) > 0 {
		err = fmt.Errorf(snapmap[tagname].Err)
	}
	return snapmap[tagname].Rtsd.Value,
		TimeFormat(Millisecond2Time(snapmap[tagname].Rtsd.Time)),
		snapmap[tagname].Rtsd.Quality,
		err
}

/****************************************************
功能:获取庚顿数据库中指定变量的快照结构
输入:
	tagname:变量名
输出:数值,时间,数据质量,错误信息
编辑:wang_jp
时间:2019年12月8日
****************************************************/
func (m *MicGolden) GoldenGetSnapShotMap(tagname ...string) (map[string]goldengo.SnapData, error) {
	snp := make(map[string]goldengo.SnapData)
	err := m.GetHandel("GoldenGetSnapShotMap") //建立到庚顿数据库的连接
	if err != nil {                            //判断连接是否有错误
		return snp, err
	}
	defer m.ReleaseHandel() //压后断开连接

	snapmap, err := m.GetSnapShotByName(tagname...)
	if err != nil { //有错误
		return snp, fmt.Errorf("读取变量[%s]快照时发生错误:[%s]", tagname, err.Error())
	}
	return snapmap, err
}

/****************************************************
功能:获取庚顿数据库中的变量的统计数据
输入:
	tagname:变量名
	beginTime:开始时间,读取快照时无用
	endTime:结束时间,读取快照时无用
输出:查询到的数据
编辑:wang_jp
时间:2019年12月8日
****************************************************/
func (m *MicGolden) GoldenGetStatisticsData(tagname, beginTime, endTime string) (statistic.StatisticData, error) {
	var datas statistic.StatisticData
	err := m.GetHandel("GoldenGetStatisticsData", tagname) //建立到庚顿数据库的连接
	if err != nil {                                        //判断连接是否有错误
		return datas, err
	}
	defer m.ReleaseHandel() //压后断开连接
	snames := strings.Split(tagname, ".")
	if len(snames) < 2 {
		return datas, fmt.Errorf("变量[%s]的名称不合法,合法的格式是:[ tablename.tagname ]", tagname)
	}

	_, err = TimeParse(beginTime)
	if err != nil {
		return datas, fmt.Errorf("变量[%s]获取历史统计数据时beginTime时间参数错误:[%s]", tagname, err.Error())
	}
	_, err = TimeParse(endTime)
	if err != nil {
		return datas, fmt.Errorf("变量[%s]获取历史统计数据时endTime时间参数错误:[%s]", tagname, err.Error())
	}

	_, vmap, err := m.GoldenGetHistoryInterval(beginTime, endTime, 0, tagname)
	if err != nil {
		return datas, err
	}

	if vmap[tagname].Err != nil {
		return datas, vmap[tagname].Err
	}

	var tsds statistic.Tsds
	//获取变量基本信息
	tag := new(OreProcessDTaglist)
	tag.GetTagBaseAttributByName(tagname)
	min := tag.MinValue
	max := tag.MaxValue
	dtype := strings.ToLower(tag.TagType)
	for _, v := range vmap[tagname].HisRtsd {
		//如果最大最小值有效,且数据超过最大最小值,跳过当前数据
		if (dtype != "bool" && min != max && (v.Value <= min || v.Value >= max)) || v.Quality != 0 {
			continue
		}
		var tsd statistic.TimeSeriesData
		tsd.Time = Millisecond2Time(v.Time)
		tsd.Value = v.Value
		tsds = append(tsds, tsd)
	}
	return tsds.Statistics(1), nil
}

/****************************************************
功能:获取庚顿数据库中的变量的单个统计数据
输入:
	tagname:变量名
	key:获取数据的指令，不区分大小写,指令列表如下
		"point"		-获取指定时间点(endtime)的值,此时begintime无效
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
	beginTime:开始时间,读取快照时无用
	endTime:结束时间,读取快照时无用
输出:[float64] 查询到的数据
	[bool] 如果endTime之后还有数据，则为true
	[error] 错误信息
编辑:wang_jp
时间:2019年12月8日
****************************************************/
func (m *MicGolden) GoldenGetSingleStatisticsData(tagname, key, beginTime, endTime string) (float64, bool, error) {
	err := m.GetHandel("GoldenGetSingleStatisticsData", tagname, key) //建立到庚顿数据库的连接
	if err != nil {                                                   //判断连接是否有错误
		return 0.0, false, err
	}
	defer func() {
		if err := recover(); err != nil {
			logs.Critical("MicGolden.GoldenGetSingleStatisticsData 中遇到错误:Tag:%s,Key:%s,BeginTime:%s,EndTime:%s;[%#v]", tagname, key, beginTime, endTime, err)
		}
		m.ReleaseHandel() //压后断开连接
	}()

	snames := strings.Split(tagname, ".")
	if len(snames) < 2 {
		return 0.0, false, fmt.Errorf("变量[%s]的名称不合法,合法的格式是:[ tablename.tagname ]", tagname)
	}
	var data float64 = 0.0
	_, err = TimeParse(beginTime)
	if err != nil {
		return 0.0, false, fmt.Errorf("变量[%s]获取单点历史数据时beginTime时间参数错误:[%s]", tagname, err.Error())
	}
	endt, err := TimeParse(endTime)
	if err != nil {
		return 0.0, false, fmt.Errorf("变量[%s]获取单点历史数据时endTime时间参数错误:[%s]", tagname, err.Error())
	}
	snap_new_than_endtime := false //快照时间比endTime新
	switch strings.ToLower(key) {
	case "snapshot":
		snaps, err := m.GetSnapShotByName(tagname)
		if err != nil {
			return data, false, err
		}
		data = snaps[tagname].Rtsd.Value
	case "advangce": //advangce在此函数中无效
		data = 0.0
	case "point": //获取指定时间点的值
		vmap, err := m.GetHistorySingleByName(1, endt.UnixNano(), tagname)
		if err != nil {
			return data, false, err
		}
		data = vmap[tagname].Value
	default:
		_, vmap, err := m.GoldenGetHistoryInterval(beginTime, endTime, 0, tagname)
		if err != nil {
			return data, vmap[tagname].Continue, err
		}
		if vmap[tagname].Err != nil {
			return data, vmap[tagname].Continue, vmap[tagname].Err
		}

		snap_new_than_endtime = vmap[tagname].Continue
		var tsds statistic.Tsds

		//获取变量基本信息
		tag := new(OreProcessDTaglist)
		tag.GetTagBaseAttributByName(tagname)
		min := tag.MinValue
		max := tag.MaxValue
		dtype := strings.ToLower(tag.TagType)
		for _, v := range vmap[tagname].HisRtsd {
			//如果最大最小值有效,且数据超过最大最小值,跳过当前数据
			if (dtype != "bool" && min != max && (v.Value <= min || v.Value >= max)) || v.Quality != 0 {
				continue
			}
			var tsd statistic.TimeSeriesData
			tsd.Time = Millisecond2Time(v.Time)
			tsd.Value = v.Value
			tsds = append(tsds, tsd)
		}
		data = statistic.SelectValueFromStatisticData(tsds.Statistics(1), key)
	}
	return data, snap_new_than_endtime, nil
}

/****************************************************
功能:批量获取庚顿数据库中变量的单个统计数据
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
	beginTime:开始时间,读取快照时无用
	endTime:结束时间,读取快照时无用
输出:查询到的数据
编辑:wang_jp
时间:2019年12月8日
****************************************************/
func (m *MicGolden) GoldenGetSingleStatisticsDataBatch(key, beginTime, endTime string, tagname ...string) (map[string]float64, error) {
	data := make(map[string]float64)
	var tn string
	if len(tagname) > 0 {
		tn = tagname[0]
	}
	err := m.GetHandel("GoldenGetSingleStatisticsDataBatch", tn, key) //建立到庚顿数据库的连接
	if err != nil {                                                   //判断连接是否有错误
		return data, err
	}
	defer m.ReleaseHandel() //压后断开连接

	_, err = TimeParse(beginTime)
	if err != nil {
		return data, fmt.Errorf("变量[%s]获取单点历史数据时beginTime时间参数错误:[%s]", tagname, err.Error())
	}
	endt, err := TimeParse(endTime)
	if err != nil {
		return data, fmt.Errorf("变量[%s]获取单点历史数据时endTime时间参数错误:[%s]", tagname, err.Error())
	}
	switch strings.ToLower(key) {
	case "snapshot":
		snaps, err := m.GetSnapShotByName(tagname...)
		if err != nil {
			return data, err
		}
		for tag, snap := range snaps {
			data[tag] = snap.Rtsd.Value
		}
	case "advangce": //advangce在此函数中无效
		break
	case "point": //获取指定时间点的值
		vmap, err := m.GetHistorySingleByName(1, endt.UnixNano(), tagname...)
		if err != nil {
			return data, err
		}
		for tag, v := range vmap {
			data[tag] = v.Value
		}
	default:
		_, vmap, err := m.GoldenGetHistoryInterval(beginTime, endTime, 0, tagname...)
		if err != nil {
			return data, err
		}
		for tag, hisv := range vmap {
			if hisv.Err != nil {
				return data, hisv.Err
			}
			var tsds statistic.Tsds
			t := new(OreProcessDTaglist)
			t.GetTagBaseAttributByName(tag)
			min := t.MinValue
			max := t.MaxValue
			dtype := strings.ToLower(t.TagType)
			for _, v := range hisv.HisRtsd {
				//如果最大最小值有效,且数据超过最大最小值,跳过当前数据
				if (dtype != "bool" && min != max && (v.Value <= min || v.Value >= max)) || v.Quality != 0 {
					continue
				}
				var tsd statistic.TimeSeriesData
				tsd.Time = Millisecond2Time(v.Time)
				tsd.Value = v.Value
				tsds = append(tsds, tsd)
			}
			data[tag] = statistic.SelectValueFromStatisticData(tsds.Statistics(1), key)
		}
	}
	return data, nil
}

/*****************************************************
功能:根据需要的数据指令从历史统计结构体中返回相应的数据
输入:历史统计结构体,数据类型指令
输出:浮点数
时间:2020年5月17日
编辑:wang_jp
*****************************************************/
func (m *MicGolden) GetSingleDataFromStatistics(v statistic.StatisticData, key string) float64 {
	return statistic.SelectValueFromStatisticData(v, key)
}

/****************************************************
功能:对庚顿历史数据进行与限制值进行比较，输出比较结果的累积值
输入:
	[tagid] 变量ID
	[tagname] 变量名
	[key] 比较指令
	[beginTime] 开始时间
	[endTime] 结束时间
输出:[float64] 比较结果累积值
	[bool] 如果endTime之后还有数据，则为true
	[error] 错误信息
编辑:wang_jp
时间:2020年1月3日
****************************************************/
func (m *MicGolden) GoldenGetAnalogComparedTotal(tagid int64, tagname, key, beginTime, endTime string) (float64, bool, error) {
	err := m.GetHandel("GoldenGetAnalogComparedTotal", tagname, key) //建立到庚顿数据库的连接
	if err != nil {                                                  //判断连接是否有错误
		return 0.0, false, err
	}
	defer m.ReleaseHandel() //压后断开连接

	_, err = TimeParse(beginTime)
	if err != nil {
		return 0.0, false, fmt.Errorf("变量[%s]模拟量比较时beginTime时间参数错误:[%s]", tagname, err.Error())
	}
	_, err = TimeParse(endTime)
	if err != nil {
		return 0.0, false, fmt.Errorf("变量[%s]模拟量比较时endTime时间参数错误:[%s]", tagname, err.Error())
	}
	_, vmap, err := m.GoldenGetHistoryInterval(beginTime, endTime, 0, tagname)
	if err != nil {
		return 0, vmap[tagname].Continue, err
	}
	if vmap[tagname].Err != nil {
		return 0, vmap[tagname].Continue, vmap[tagname].Err
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
	case "lt_min_total", "lte_min_total", "gt_min_total", "gte_min_total":
		limit = tag.MinValue
	case "lt_max_total", "lte_max_total", "gt_max_total", "gte_max_total":
		limit = tag.MaxValue
	}

	return m.analogComparedTotal(vmap[tagname].HisRtsd, key, limit), vmap[tagname].Continue, nil
}

/****************************************************
功能:对庚顿历史数据进行与限制值进行比较，输出比较结果的累积值
输入:
	hisdata:历史数据map,key:比较指令,limit:比较值
输出:比较结果累积值
编辑:wang_jp
时间:2020年1月3日
****************************************************/
func (m *MicGolden) analogComparedTotal(hisdata []goldengo.RealTimeSeriesData, key string, limit float64) float64 {
	var v, v0, total float64
	var t0 time.Time

	for k, dv := range hisdata {
		if dv.Quality == 0 { //只有GOOD值有效
			t := Millisecond2Time(dv.Time)
			if k > 0 { //非第一个数
				if m.compare(key, dv.Value, limit) {
					v = 1.0
				} else {
					v = 0.0
				}
				total += (v0 + v) / 2 * float64(t.Unix()-t0.Unix())
				v0 = v
			} else { //第一个数
				if m.compare(key, dv.Value, limit) {
					v0 = 1.0
				} else {
					v0 = 0.0
				}
			}
			t0 = t
		}
	}

	return total
}

/****************************************************
功能:根据指令比较输入的比较值和被比较值
输入:
	key:比较指令,value:被比较值,limit:比较值
输出:比较结果
编辑:wang_jp
时间:2020年1月3日
****************************************************/
func (m *MicGolden) compare(key string, value, limit float64) bool {
	switch key {
	case "lt_l_total", "lt_ll_total", "lt_h_total", "lt_hh_total", "lt_min_total", "lt_max_total":
		return value < limit
	case "lte_l_total", "lte_ll_total", "lte_h_total", "lte_hh_total", "lte_min_total", "lte_max_total":
		return value <= limit
	case "gt_l_total", "gt_ll_total", "gt_h_total", "gt_hh_total", "gt_min_total", "gt_max_total":
		return value > limit
	case "gte_l_total", "gte_ll_total", "gte_h_total", "gte_hh_total", "gte_min_total", "gte_max_total":
		return value >= limit
	default:
		return false
	}
}

/****************************************************
功能:获取庚顿历史数据库中历史数据的峰谷值统计
输入:
	[tagid] 变量ID
	[tagname] 变量名
	[key] 数据指令
	[beginTime] 开始时间
	[endTime] 结束时间
输出:[interface{}] 结果值
	[bool] 如果endTime之后还有数据，则为true
	[error] 错误信息
输出:数据接口,err
时间:2020年1月19日
编辑:wang_jp
****************************************************/
func (m *MicGolden) GoldenGetPeakValleyOfHisData(tagid int64, tagname, key, beginTime, endTime string) (interface{}, bool, error) {
	tag := new(OreProcessDTaglist)
	tag.Id = tagid
	if tagid == 0 { //没有设置ID
		id, err := tag.GetTagIDByFullName(tagname)
		if err != nil {
			fmt.Println("获取ID失败", err.Error())
			return 0.0, false, err
		}
		tagid = id
	} else { //设置了ID
		fullname, err := tag.GetTagFullNameByID()
		if err != nil {
			return 0.0, false, err
		}
		tagname = fullname
	}
	err := tag.GetTagById()
	if err != nil {
		return 0.0, false, err
	}

	stdv := tag.UserReal1          //user_real1作为稳态判据
	increment := tag.UserReal2     //user_real2作为拐点增量值
	cp := tag.UserInt1             //user_int1作为最大连续点数
	negativeAsZero := tag.UserInt2 //user_int2如果为0，保留负数;如果为1,将负数作为0处理

	err = m.GetHandel("GoldenGetPeakValleyOfHisData", tagname) //建立到庚顿数据库的连接
	if err != nil {                                            //判断连接是否有错误
		return 0.0, false, err
	}
	defer m.ReleaseHandel() //压后断开连接
	_, err = TimeParse(beginTime)
	if err != nil {
		return 0.0, false, err
	}
	_, err = TimeParse(endTime)
	if err != nil {
		return 0.0, false, err
	}

	_, vmap, err := m.GoldenGetHistoryInterval(beginTime, endTime, 0, tagname)
	if err != nil {
		return 0.0, vmap[tagname].Continue, err
	}
	if vmap[tagname].Err != nil {
		return 0.0, vmap[tagname].Continue, vmap[tagname].Err
	}
	pvs := new(statistic.PeakValleySelector)
	pvs.New(stdv, int(cp))
	var tsdt []statistic.TimeSeriesData
	for _, his := range vmap {
		var tsd statistic.TimeSeriesData
		for _, data := range his.HisRtsd {
			if data.Quality == 0 { //只有GOOD时有效
				tsd.Time = Millisecond2Time(data.Time)
				tsd.Value = data.Value
				tsdt = append(tsdt, tsd)
			}
		}
	}
	pvs.DataFillter(tsdt, increment, int(negativeAsZero))
	switch strings.ToLower(key) {
	case "pvs_peaksum":
		return pvs.PeakSum, vmap[tagname].Continue, nil
	case "pvs_valleysum":
		return pvs.ValleySum, vmap[tagname].Continue, nil
	case "pvs_pvdiffsum":
		return pvs.PVDiffSum, vmap[tagname].Continue, nil
	case "pvs_periodcnt":
		return pvs.PeriodCnt, vmap[tagname].Continue, nil
	case "pvs_peek_valley_datas":
		return pvs.PvDatas, vmap[tagname].Continue, nil
	case "peakvalley":
		return *pvs, vmap[tagname].Continue, nil
	default:
		return 0.0, vmap[tagname].Continue, fmt.Errorf("Invalid instruction;[无效的指令]")
	}

}

/***************************************************
功能:读取庚顿补齐头尾的历史数据
输入:
	tagname:变量全名,beginTime:起始时间,endTime:结束时间
输出:时间序列数据,err
时间:2020年3月19日
编辑:wang_jp
***************************************************/
func (m *MicGolden) GoldenGetHistory(begintime, endtime string, tagname ...string) (map[string][]SrtData, error) {
	srtdmap, _, err := m.GoldenGetHistoryInterval(begintime, endtime, 0, tagname...)
	return srtdmap, err
}

/***************************************************
功能:将庚顿历史数据转换为实时数据
输入:[hismap] 庚顿历史数据结构Map
输出:时间序列数据Map
时间:2020年6月17日
编辑:wang_jp
***************************************************/
func (m *MicGolden) goldenHisDataToSrtData(hismap map[string]goldengo.HisData) map[string][]SrtData {
	srtdmap := make(map[string][]SrtData)
	for tag, his := range hismap {
		var srtds []SrtData
		for _, hd := range his.HisRtsd {
			var srtd SrtData
			srtd.Time = TimeFormat(Millisecond2Time(hd.Time))
			srtd.Value = hd.Value
			srtds = append(srtds, srtd)
		}
		srtdmap[tag] = srtds
	}
	return srtdmap
}

/***************************************************
功能:读取庚顿等间隔历史数据
输入:
	tagname:变量全名,beginTime:起始时间,endTime:结束时间,interval:时间间隔,以秒为单位
输出:时间序列数据,err
说明:interval为0时,读取原始历史数据
	时间长度如果超过8小时，分为多个线程并发读取
时间:2020年3月19日
编辑:wang_jp
***************************************************/
func (m *MicGolden) GoldenGetHistoryInterval(begintime, endtime string, interval int64, tagname ...string) (map[string][]SrtData, map[string]goldengo.HisData, error) {
	datamap := make(map[string][]SrtData)
	hisdatamap := make(map[string]goldengo.HisData)
	bgtime, err := TimeParse(begintime)
	if err != nil {
		return datamap, hisdatamap, err
	}
	edtime, err := TimeParse(endtime)
	if err != nil {
		return datamap, hisdatamap, err
	}

	maxslice := EngineCfgMsg.Sys.MaxHisSliceTime
	if EngineCfgMsg.Sys.MaxHisSliceTime == 0 { //不可为0
		EngineCfgMsg.Sys.MaxHisSliceTime = 8
		maxslice = 8
	}
	var maxgocnt int64 = 10
	gocnt := int64(math.Ceil(float64(edtime.Unix()-bgtime.Unix()) / 3600.0 / float64(maxslice)))
	if gocnt > maxgocnt {
		gocnt = maxgocnt
		maxslice = int64(math.Ceil(float64(edtime.Unix()-bgtime.Unix()) / float64(gocnt) / 3600.0))
	}
	if gocnt <= 1 { //时间不超过1个最大切片周期
		hmp, err := m.getGoldenHistoryInterval(begintime, endtime, interval, tagname...)
		return m.goldenHisDataToSrtData(hmp), hmp, err
	}

	hisdata := make(chan map[string]goldengo.HisData, gocnt)
	defer func() {
		if err := recover(); err != nil {
			logs.Critical("MicGolden.GoldenGetHistoryInterval 中遇到错误:Tags:%s,BetinTime:%s,EndTime:%s;[%#v]", tagname, begintime, endtime, err)
		}
		close(hisdata) //压后关闭chan道
	}()
	var hiswait sync.WaitGroup
	var i int64
	for i = 0; i < gocnt; i++ { //遍历，生成并发的Go程
		hiswait.Add(1)
		st := bgtime.Unix() + i*3600*maxslice
		ed := st + 3600*maxslice
		if ed > edtime.Unix() {
			ed = edtime.Unix()
		}
		go func(res chan map[string]goldengo.HisData) {
			defer func() {
				if err := recover(); err != nil {
					logs.Critical("%#v", err)
				}
				hiswait.Done()
			}()
			micgd := new(MicGolden)
			his, err := micgd.getGoldenHistoryInterval(TimeFormat(time.Unix(st, 0)), TimeFormat(time.Unix(ed, 0)), interval, tagname...)
			if err != nil {
				logs.Warning("读取庚顿历史数据时发生错误,Tag=%+v,BeginTime=[%s],EndTime=[%s]。错误信息:[%s]",
					tagname, TimeFormat(time.Unix(st, 0)), TimeFormat(time.Unix(ed, 0)), err.Error())
			}
			res <- his
		}(hisdata)
	}

	hismaparr := make([]map[string]goldengo.HisData, gocnt)
	for i = 0; i < gocnt; i++ {
		his := <-hisdata //接收计算结果
		for _, h := range his {
			if len(h.HisRtsd) > 0 { //结果中有历史数据
				ht := time.Unix(h.HisRtsd[0].Time/1e3, h.HisRtsd[0].Time%1e3*1e6)
				if ht.UnixNano() < bgtime.UnixNano() { //读取到的历史数据的第一个时间小于起始时间
					ht = bgtime
				}
				s := (ht.Unix() - bgtime.Unix()) / 3600 / EngineCfgMsg.Sys.MaxHisSliceTime //计算序列号
				if s < gocnt && s >= 0 {
					hismaparr[s] = his
					break
				}
			}
		}
	}
	hiswait.Wait()
	tagerrmap := make(map[string][]error)
	for i, hismap := range hismaparr {
		for tag, datas := range hismap {
			var errs []error
			if _, ok := tagerrmap[tag]; ok {
				errs = tagerrmap[tag] //存在
			}
			if datas.Err != nil {
				errs = append(errs, datas.Err)
			}
			tagerrmap[tag] = errs

			var hisdata goldengo.HisData
			hisdata = datas

			if i == 0 {
				hisdatamap[tag] = datas
			} else {
				dl := len(hisdatamap[tag].HisRtsd)    //已有数据长度
				if len(datas.HisRtsd) > 0 && dl > 1 { //有数据
					//已有数据最后一个与现数据第一个重合
					if hisdatamap[tag].HisRtsd[dl-1].Time == datas.HisRtsd[0].Time || hisdatamap[tag].HisRtsd[dl-1].Value == datas.HisRtsd[0].Value {
						hisdata.HisRtsd = append(hisdatamap[tag].HisRtsd, datas.HisRtsd[1:]...)
					} else {
						hisdata.HisRtsd = append(hisdatamap[tag].HisRtsd, datas.HisRtsd...)
					}
				} else {
					hisdata.HisRtsd = append(hisdatamap[tag].HisRtsd, datas.HisRtsd...)
				}
			}
			hisdatamap[tag] = hisdata
		}
	}
	for tag, hisdata := range hisdatamap {
		if len(tagerrmap[tag]) == len(hismaparr) {
			hisdata.Err = fmt.Errorf("变量[%s]没有读取到从[%s]到[%s]的历史数据",
				tag, begintime, endtime)
			hisdatamap[tag] = hisdata
		}
	}
	return m.goldenHisDataToSrtData(hisdatamap), hisdatamap, nil
}

/***************************************************
功能:读取庚顿等间隔历史数据
输入:
	tagname:变量全名,beginTime:起始时间,endTime:结束时间,interval:时间间隔,以秒为单位
输出:庚顿历史数据,err
说明:interval为0时,读取原始历史数据
时间:2020年3月19日
编辑:wang_jp
***************************************************/
func (m *MicGolden) getGoldenHistoryInterval(begintime, endtime string, interval int64, tagname ...string) (map[string]goldengo.HisData, error) {
	datamap := make(map[string]goldengo.HisData)

	var tenantmsg []string
	tenantmsg = append(tenantmsg, "getGoldenHistoryInterval")
	tenantmsg = append(tenantmsg, tagname...)
	tenantmsg = append(tenantmsg, begintime)
	tenantmsg = append(tenantmsg, endtime)
	err := m.GetHandel(tenantmsg...) //建立到庚顿数据库的连接
	if err != nil {                  //判断连接是否有错误
		return datamap, err
	}

	defer func() {
		if err := recover(); err != nil {
			logs.Critical(err)
		}
		m.ReleaseHandel() //压后断开连接
	}()
	bgtime, err := TimeParse(begintime)
	if err != nil {
		return datamap, err
	}
	edtime, err := TimeParse(endtime)
	if err != nil {
		return datamap, err
	}
	return m.GetHistoryDataAlignHeadAndTail(bgtime.UnixNano(), edtime.UnixNano(), int(interval), tagname...)
}

/***************************************************
功能:读取庚顿历史统计数据
输入:
	tagname:变量全名,beginTime:起始时间,endTime:结束时间
输出:时间序列数据,err
时间:2020年3月19日
编辑:wang_jp
***************************************************/
func (m *MicGolden) GoldenGetHistorySummary(beginTime, endTime string, tagname ...string) (map[string]statistic.StatisticData, error) {
	type static struct {
		Datas statistic.StatisticData
		Err   string
	}
	datas := make(map[string]statistic.StatisticData)

	var tenantmsg []string
	tenantmsg = append(tenantmsg, "GoldenGetHistorySummary")
	tenantmsg = append(tenantmsg, tagname...)
	tenantmsg = append(tenantmsg, beginTime)
	tenantmsg = append(tenantmsg, endTime)
	err := m.GetHandel(tenantmsg...) //建立到庚顿数据库的连接
	if err != nil {                  //判断连接是否有错误
		return datas, err
	}
	defer m.ReleaseHandel() //压后断开连接

	_, err = TimeParse(beginTime)
	if err != nil {
		return datas, fmt.Errorf("变量[%s]获取统计数据时beginTime时间参数错误:[%s]", tagname, err.Error())
	}
	_, err = TimeParse(endTime)
	if err != nil {
		return datas, fmt.Errorf("变量[%s]获取统计数据时endTime时间参数错误:[%s]", tagname, err.Error())
	}

	_, vmap, err := m.GoldenGetHistoryInterval(beginTime, endTime, 0, tagname...)
	if err != nil {
		return datas, err
	}

	for tag, tagvalues := range vmap {
		var tsds statistic.Tsds
		t := new(OreProcessDTaglist)
		t.GetTagBaseAttributByName(tag)
		min := t.MinValue
		max := t.MaxValue
		dtype := strings.ToLower(t.TagType)
		for _, v := range tagvalues.HisRtsd {
			//如果最大最小值有效,且数据超过最大最小值,或者质量码不为GOOD,跳过当前数据
			if (dtype != "bool" && min != max && (v.Value <= min || v.Value >= max)) || v.Quality != 0 {
				continue
			}
			var tsd statistic.TimeSeriesData
			tsd.Time = Millisecond2Time(v.Time)
			tsd.Value = v.Value
			tsds = append(tsds, tsd)
		}
		stat := tsds.Statistics(1)
		datas[tag] = stat
	}

	return datas, nil
}

/***************************************************
功能:读取庚顿服务器数据标签表
输入:[selector] 不输入或者0:读取数据表列表,反馈为数据表列表的map,以表id为key
		1: 读取数据表列表,反馈为数据表ID数组
		2: 读取数据表列表,反馈为数据表名字符串数组
		其他:同0
输出:数据接口,err
时间:2020年3月19日
编辑:wang_jp
***************************************************/
func (m *MicGolden) GoldenGetTables(selector ...int) (interface{}, error) {
	err := m.GetHandel("GoldenGetTables") //建立到庚顿数据库的连接
	if err != nil {                       //判断连接是否有错误
		return nil, err
	}
	defer m.ReleaseHandel() //压后断开连接
	err = m.GetTables()
	if err != nil {
		return nil, err
	}
	if len(selector) > 0 {
		switch selector[0] {
		case 1:
			return m.TableIds, err //返回ID数组
		case 2:
			var names []string
			for _, id := range m.TableIds {
				names = append(names, m.Tables[id].Name)
			}
			return names, nil //返回表名数组
		default:
			return m.Tables, nil //返回表信息Map
		}
	}
	return m.Tables, nil
}

/***************************************************
功能:根据表名读取庚顿服务器数据标签表信息
输入:[tablenames] 零个一个或者多个数据表名,如果不输入,则返回所有表的信息
输出:map[string]goldengo.GoldenTable,以表名为Key
	err,错误信息
时间:2020年3月19日
编辑:wang_jp
***************************************************/
func (m *MicGolden) GoldenGetTablePropertyByTableName(tablenames ...string) (map[string]goldengo.GoldenTable, error) {
	datas := make(map[string]goldengo.GoldenTable)
	tabp, err := m.GoldenGetTables(0)
	if err != nil {
		return datas, err
	}
	tbinfs, _ := tabp.(map[int]goldengo.GoldenTable)
	if len(tablenames) == 0 {
		for _, tb := range tbinfs {
			datas[tb.Name] = tb
		}
	} else {
		var gdtb goldengo.GoldenTable
		for _, name := range tablenames {
			for _, tb := range tbinfs {
				datas[name] = gdtb
				if strings.ToLower(tb.Name) == strings.ToLower(name) {
					datas[name] = tb
					break
				}
			}
		}
	}
	return datas, nil
}

/***************************************************
功能:读取庚顿服务器数据标签表下的标签名列表
输入:[tablenames] 表名切片
	如果不输入表名,则读取数据库中所有表下的标签点
	如果输入表名,则读取输入的每个表下的标签点
	如果输入的表不存在,对应的map下的字符串切片为空
输出:时间,err
时间:2020年5月17日
编辑:wang_jp
***************************************************/
func (m *MicGolden) GoldenGetTagNameListInTables(tablenames ...string) (map[string][]string, error) {
	var tenantmsg []string
	tenantmsg = append(tenantmsg, "GoldenGetTagNameListInTables")
	tenantmsg = append(tenantmsg, tablenames...)
	err := m.GetHandel(tenantmsg...) //建立到庚顿数据库的连接
	if err != nil {                  //判断连接是否有错误
		return nil, err
	}
	defer m.ReleaseHandel() //压后断开连接
	return m.GetTagNameListInTables(tablenames...)
}

/***************************************************
功能:读取庚顿服务器数据标签表下的标签详细信息列表
输入:[tablenames] 表名切片
输出:[map[string][]goldengo.GoldenPoint],key为标签点全名称
	[err] 错误信息
说明:如果不输入表名,则读取数据库中所有表下的标签点
	如果输入表名,则读取输入的每个表下的标签点
	如果输入的表不存在,对应的map下的字符串切片为空
时间:2020年5月17日
编辑:wang_jp
***************************************************/
func (m *MicGolden) GoldenGetTagListInTables(tablenames ...string) (map[string][]goldengo.GoldenPoint, error) {
	var tenantmsg []string
	tenantmsg = append(tenantmsg, "GoldenGetTagListInTables")
	tenantmsg = append(tenantmsg, tablenames...)
	err := m.GetHandel(tenantmsg...) //建立到庚顿数据库的连接
	if err != nil {                  //判断连接是否有错误
		return nil, err
	}
	defer m.ReleaseHandel() //压后断开连接
	return m.GetTagListInTables(tablenames...)
}

/***************************************************
功能:读取特定时间点的庚顿历史数据
输入:[mode]   整型，输入，取值 GOLDEN_NEXT(0)、GOLDEN_PREVIOUS(1)、GOLDEN_EXACT(2)、
    			GOLDEN_INTER(3) 之一:
    				 GOLDEN_NEXT(0) 寻找下一个最近的数据；
   					 GOLDEN_PREVIOUS(1) 寻找上一个最近的数据；
    				 GOLDEN_EXACT(2) 取指定时间的数据，如果没有则返回错误 GoE_DATA_NOT_FOUND；
    				 GOLDEN_INTER(3) 取指定时间的内插值数据。
    [datatime]  字符串时间
	[tagnames]  字符串切片，输入. 变量全名格式:tablename.tagname
输出:庚顿时间序列数据Map,变量名为key
	err,错误信息
时间:2020年5月17日
编辑:wang_jp
***************************************************/
func (m *MicGolden) GoldenGetHistorySinglePoint(mode int, datatime string, tagnames ...string) (map[string]goldengo.RealTimeSeriesData, error) {
	datas := make(map[string]goldengo.RealTimeSeriesData)
	var tenantmsg []string
	tenantmsg = append(tenantmsg, "GoldenGetHistorySinglePoint")
	tenantmsg = append(tenantmsg, tagnames...)
	err := m.GetHandel(tenantmsg...) //建立到庚顿数据库的连接
	if err != nil {                  //判断连接是否有错误
		return datas, err
	}
	defer m.ReleaseHandel() //压后断开连接
	dt, err := TimeParse(datatime)
	if err != nil {
		return datas, fmt.Errorf("变量[%s]获取单点历史数据时时间点参数[%s]错误:[%s]", tagnames, datatime, err.Error())
	}
	return m.GetHistorySingleByName(mode, dt.UnixNano(), tagnames...)
}

/***************************************************
功能:通过标签点名读取庚顿服务器上的标签详细信息
输入:[tablenames] 表名切片
输出:[map[string]goldengo.GoldenPoint],key为标签点全名称
	[err] 错误信息
时间:2020年5月17日
编辑:wang_jp
***************************************************/
func (m *MicGolden) GoldenGetTagPointInfoByName(tagnames ...string) (map[string]goldengo.GoldenPoint, error) {
	var datas map[string]goldengo.GoldenPoint
	err := m.GetHandel("GoldenGetTagPointInfoByName") //建立到庚顿数据库的连接
	if err != nil {                                   //判断连接是否有错误
		return datas, err
	}
	defer m.ReleaseHandel() //压后断开连接
	return m.GetTagPointInfoByName(tagnames...)
}

/***************************************************
功能:通过标签点名读取庚顿服务器上的标签详细信息
输入:[tablenames] 表名切片
输出:[map[string]goldengo.GoldenPoint],key为标签点全名称
	[err] 错误信息
时间:2020年5月17日
编辑:wang_jp
***************************************************/
func (m *MicGolden) GoldenGetTagPointInfoById(id int) (goldengo.GoldenPoint, error) {
	var datas goldengo.GoldenPoint
	err := m.GetHandel("GoldenGetTagPointInfoById") //建立到庚顿数据库的连接
	if err != nil {                                 //判断连接是否有错误
		return datas, err
	}
	defer m.ReleaseHandel() //压后断开连接
	return m.GetSinglePointPropterty(id)
}

/***************************************************
功能:添加表
输入:[tableName]  表名称
	[tableDesc]  表描述
输出:[int] 表ID
	[error] 错误信息
说明:
时间:2020年6月18日
编辑:wang_jp
***************************************************/
func (m *MicGolden) GoldenTableInsert(tableName, tableDesc string) (int, error) {
	err := m.GetHandel("GoldenTableInsert", tableName) //建立到庚顿数据库的连接
	if err != nil {                                    //判断连接是否有错误
		return 0, err
	}
	defer m.ReleaseHandel() //压后断开连接
	tb := new(goldengo.GoldenTable)
	tb.Name = tableName
	tb.Desc = tableDesc
	err = tb.AppendTable(m.Handle)

	return tb.Id, err
}

/***************************************************
功能:添加表或者更新表
输入:[tableName]  表名称
	[tableDesc]  表描述
输出:[bool] 如果新建表，则为true;更新表为false
	[int] 表ID
	[error] 错误信息
说明:如果表名称不存在，则添加表，如果存在，则更新表的描述信息
时间:2020年6月18日
编辑:wang_jp
***************************************************/
func (m *MicGolden) GoldenTableInsertOrUpdate(tableName, tableDesc string) (bool, int, error) {
	err := m.GetHandel("GoldenTableInsertOrUpdate", tableName) //建立到庚顿数据库的连接
	if err != nil {                                            //判断连接是否有错误
		return false, 0, err
	}
	defer m.ReleaseHandel() //压后断开连接
	isinsert := false
	tb := new(goldengo.GoldenTable)
	tb.Name = tableName
	tb.Desc = tableDesc
	e := tb.GetTablePropertyByName(m.Handle)
	if e != nil || tb.Id == 0 { //如果读取时有错误或者没有读取到ID
		err = tb.AppendTable(m.Handle) //插入新数据
		isinsert = true
	} else {
		err = tb.UpdateTableDescByName(m.Handle) //更新描述信息
	}
	return isinsert, tb.Id, err
}

/***************************************************
功能:删除表
输入:[tableName]  表名称
输出:[error] 错误信息
说明:
时间:2020年6月18日
编辑:wang_jp
***************************************************/
func (m *MicGolden) GoldenTableRemove(tableName string) error {
	err := m.GetHandel("GoldenTableRemove", tableName) //建立到庚顿数据库的连接
	if err != nil {                                    //判断连接是否有错误
		return err
	}
	defer m.ReleaseHandel() //压后断开连接
	tb := new(goldengo.GoldenTable)
	tb.Name = tableName
	return tb.RemoveTableByName(m.Handle)
}

/***************************************************
功能:重命名表
输入:[oldtableName]  旧表名称
	[newtableName]  新表名称
输出:[error] 错误信息
说明:
时间:2020年6月18日
编辑:wang_jp
***************************************************/
func (m *MicGolden) GoldenTableReName(oldtableName, newtableName string) error {
	err := m.GetHandel("GoldenTableReName", oldtableName, newtableName) //建立到庚顿数据库的连接
	if err != nil {                                                     //判断连接是否有错误
		return err
	}
	defer m.ReleaseHandel() //压后断开连接
	tb := new(goldengo.GoldenTable)
	tb.Name = oldtableName
	return tb.UpdateTableNameByOldName(m.Handle, newtableName)
}

/***************************************************
功能:通过表名获取表ID
输入:[tableName]  表名称
输出:[int] 表ID
	[error] 错误信息
说明:
时间:2020年6月18日
编辑:wang_jp
***************************************************/
func (m *MicGolden) GoldenTableGetIdByName(tableName string) (int, error) {
	err := m.GetHandel("GoldenTableGetIdByName", tableName) //建立到庚顿数据库的连接
	if err != nil {                                         //判断连接是否有错误
		return 0, err
	}
	defer m.ReleaseHandel() //压后断开连接
	tb := new(goldengo.GoldenTable)
	tb.Name = tableName
	err = tb.GetTablePropertyByName(m.Handle)
	return tb.Id, err
}

/*******************************************************************************
功能:通过tagid,key,begineTime,endTime从庚顿数据库获取历史统计数值或者快照数值
输入:tagid,key, begineTime, endTime
输出:[interface{}]查询结果
	[bool] 快照值是否大于endTime
	[error]错误信息
说明:有错误信息的时候查询结果为0
编辑:wang_jp
时间:2019年12月9日
*******************************************************************************/
func (micgd *MicGolden) GetTagHisSumValueFromGoldenByID(tagid int64, tagname, key, begineTime, endTime string) (interface{}, bool, error) {
	tag := new(OreProcessDTaglist)
	tag.Id = tagid
	if len(tagname) == 0 {
		tname, err := tag.GetTagFullNameByID()
		if err != nil {
			return 0.0, false, err
		}
		tagname = tname
	}
	switch strings.ToLower(key) {
	case "lt_l_total", "lte_l_total", "gt_l_total", "gte_l_total",
		"lt_ll_total", "lte_ll_total", "gt_ll_total", "gte_ll_total",
		"lt_h_total", "lte_h_total", "gt_h_total", "gte_h_total",
		"lt_hh_total", "lte_hh_total", "gt_hh_total", "gte_hh_total",
		"lt_min_total", "lte_min_total", "gt_min_total", "gte_min_total",
		"lt_max_total", "lte_max_total", "gt_max_total", "gte_max_total":
		return micgd.GoldenGetAnalogComparedTotal(tagid, tagname, key, begineTime, endTime)
	case "pvs_peaksum", "pvs_valleysum", "pvs_pvdiffsum", "pvs_periodcnt", "pvs_peek_valley_datas":
		return micgd.GoldenGetPeakValleyOfHisData(tagid, tagname, key, begineTime, endTime)
	default:
		return micgd.GoldenGetSingleStatisticsData(tagname, key, begineTime, endTime)
	}
}

/***************************************************
功能:批量写快照
输入:[tagnames]   标签点全名.同一个标签点标识可以出现多次，但它们的时间戳必需是递增的
	[datavalues] 数值
	[qualities]  质量码
	[datatimes]  时间,UnixNano
输出:
说明:
时间:2020年5月18日
编辑:wang_jp
***************************************************/
func (m *MicGolden) GoldenSetSnapShots(tagnames []string, datavalues []float64, qualities []int, datatimes []string) error {
	err := m.GetHandel("GoldenSetSnapShots") //建立到庚顿数据库的连接
	if err != nil {                          //判断连接是否有错误
		return err
	}
	defer m.ReleaseHandel() //压后断开连接
	var dtimes []int64
	for _, dt := range datatimes { //转换时间
		var t time.Time
		if len(dt) == 0 {
			t = time.Now()
		} else {
			var e error
			t, e = TimeParse(dt)
			if e != nil {
				return fmt.Errorf("批量写快照时时间戳格式错误:[%s]", e.Error())
			}
		}
		dtimes = append(dtimes, t.UnixNano())
	}
	return m.SetSnapShotBatch(tagnames, datavalues, qualities, dtimes)
}

/***************************************************
功能:根据标签点名批量写历史值
输入:[tagnames]   标签点全名.同一个标签点标识可以出现多次，但它们的时间戳必需是递增的
	[datavalues] 数值
	[qualities]  质量码
	[datatimes]  时间,UnixNano
输出:
说明:
时间:2020年8月17日
编辑:wang_jp
***************************************************/
func (m *MicGolden) GoldenSetArchivedValues(tagnames []string, datavalues []float64, qualities []int, datatimes []string) error {
	err := m.GetHandel("GoldenSetArchivedValues") //建立到庚顿数据库的连接
	if err != nil {                               //判断连接是否有错误
		return err
	}
	defer m.ReleaseHandel() //压后断开连接
	var dtimes []int64
	for _, dt := range datatimes { //转换时间
		t, e := TimeParse(dt)
		if e != nil {
			return fmt.Errorf("批量写快照时时间戳格式错误:[%s]", e.Error())
		}
		dtimes = append(dtimes, t.UnixNano())
	}
	return m.SetArchivedValuesBatch(tagnames, datavalues, qualities, dtimes)
}
