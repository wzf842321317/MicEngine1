package models

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/astaxie/beego/orm"
	"github.com/bkzy-wangjp/MicEngine/MicScript/statistic"
	_ "github.com/go-sql-driver/mysql"
)

/*********************************************************************
功能:根据tagid或者tagname以及时间参数读取Sys_real_data表中的数据
输入:tagid,tagname,tagtype,key,beginTime,endTime
输出:tagId和错误信息
说明:如果tagid为0,则通过tagname查询数据
编辑:wang_jp
时间:2020年1月3日
*********************************************************************/
func (srd *SysRealData) GetSysRealTimeDataStatisticByKey(tagid int64, tagname, tagtype, key, beginTime, endTime string) (float64, bool, error) {
	switch strings.ToLower(key) {
	case "point": //读取单点数据
		datas, err := srd.GetSysRealTimeDatas(tagid, tagname, tagtype, beginTime, endTime, -1)
		if err != nil {
			if strings.Contains(err.Error(), "No data found in the set condition") { //没有查询到数据
				res, e := srd.GetSysRealTimeDatas(tagid, tagname, tagtype, beginTime, endTime, 1) //读取指定时间范围之后的数据
				if len(res) > 0 && e == nil {                                                     //后面有新数据，但是当期没有数据
					return 0.0, true, fmt.Errorf("%s,最近的数据时间是:%s", err, res[0].Datatime) //输出特殊错误标志"0",计算继续往下进行
				} else {
					return 0.0, false, err
				}
			} else {
				return 0.0, false, err
			}
		}
		if len(datas) > 0 {
			return datas[0].Value, false, nil
		}
	default: //读取统计结果数据
		sttt, ctnue, err := srd.GetSysRealTimeDataStatistic(tagid, tagname, tagtype, key, beginTime, endTime)
		if err != nil {
			return 0.0, ctnue, err
		}
		return statistic.SelectValueFromStatisticData(sttt, key), false, nil
	}
	return 0.0, false, nil
}

/*********************************************************************
功能:根据tagid或者tagname以及时间参数对读取自Sys_real_data表中的数据进行统计计算
输入:tagid,tagname,tagtype,key,beginTime,endTime
	filloutliers:填充异常值的方法,可选 midean(中位数) 或者 extremum(四分位极值),不选不处理
输出:
	statistic.StatisticData:统计数据结构体
	bool:1.设定时间范围内数据正常,输出为true
		2.设定时间范围内无数据。但endtime之后有数据时输出为true,否则为false
	error:错误信息
说明:1.如果tagid为0,则通过tagname查询数据
编辑:wang_jp
时间:2020年1月12日
*********************************************************************/
func (srd *SysRealData) GetSysRealTimeDataStatistic(tagid int64, tagname, tagtype, key, beginTime, endTime string, filloutliers ...string) (statistic.StatisticData, bool, error) {
	var sttt statistic.StatisticData //返回的统计数据
	var tsds statistic.Tsds          //时间序列数组
	var tsd statistic.TimeSeriesData
	if key == "diff" || key == "plusdiff" { //求差的时候需要带上上一个计算周期的最后一个值
		rtds, err := srd.GetSysRealTimeDatas(tagid, tagname, tagtype, beginTime, endTime, -1) //读取指定时间点的数据
		if err == nil {                                                                       //如果没有读到数据开始时间点之前的数据
			for _, rtd := range rtds { //整理出标准的时间序列数据
				tsd.Time, _ = TimeParse(rtd.Datatime)
				tsd.Value = rtd.Value
				tsds = append(tsds, tsd)
			}
		}
	}
	rtds, err := srd.GetSysRealTimeDatas(tagid, tagname, tagtype, beginTime, endTime, 0) //读取指定时间范围内的数据
	if err != nil {                                                                      //当期查询错误
		if strings.Contains(err.Error(), "No data found in the set condition") { //没有查询到数据
			res, e := srd.GetSysRealTimeDatas(tagid, tagname, tagtype, beginTime, endTime, 1) //读取指定时间范围之后的数据
			if len(res) > 0 && e == nil {                                                     //后面有新数据，但是当期没有数据
				return sttt, true, err //输出特殊错误标志"true",计算继续往下进行
			} else {
				return sttt, false, err
			}
		} else {
			return sttt, false, err
		}
	}
	for _, rtd := range rtds { //整理出标准的时间序列数据
		tsd.Time, _ = TimeParse(rtd.Datatime)
		tsd.Value = rtd.Value
		tsds = append(tsds, tsd)
	}
	if len(filloutliers) > 0 {
		switch strings.ToLower(filloutliers[0]) {
		case "midean": //用中位数填充异常值
			tsds.FillOutliersByMidean()
		case "extremum": //用四分位极值填充异常值
			tsds.FillOutliersByExtremum()
		default:
		}
	}
	switch strings.ToLower(key) {
	case "advangce", "sd", "stddev", "se", "ske", "kur", "mode", "median", "groupdist", "distribution":
		sttt = tsds.Statistics(1) //包含高级历史统计
	case "needraw":
		sttt = tsds.Statistics(2) //包含原始数据的高级历史统计
	default:
		sttt = tsds.Statistics(0) //只有普通统计
	}
	sttt.BeginTime = beginTime
	sttt.EndTime = endTime

	return sttt, true, nil
}

/*********************************************************************
功能:根据tagid或者tagname以及时间参数读取Sys_real_data表中的数据
输入:tagid,tagname,tagtype,beginTime,endTime,readType
输出:结果数组和错误信息
说明:1.如果tagid为0,则通过tagname查询数据
	2.如果readType为负数,则取小于等于beginTime的一个数据点,endTime无效
	3.如果readType为0,则取beginTime和endTime之间的数
	4.如果readType为正,则取大于等于endTime的数据点,beginTime无效
编辑:wang_jp
时间:2020年1月12日
*********************************************************************/
func (srd *SysRealData) GetSysRealTimeDatas(tagid int64, tagname, tagtype, beginTime, endTime string, readType int) ([]SysRealDataExi, error) {
	var srtds []SysRealData
	var srtdsexi []SysRealDataExi
	o := orm.NewOrm() //新建orm对象
	o.Using("default")
	dic := new(SysDictionary)
	sysdicid, err := dic.GetSysDicIdByNameCode(tagtype) //校验/读取tagtype的id
	if err != nil {
		return srtdsexi, err
	}
	if tagid == 0 { //tagid无效,按照tagname读取数据
		srtd := new(SysRealData)
		qs := o.QueryTable(srtd)
		if readType < 0 { //读取时间点上的单点数据
			qs = qs.Filter("SysDictionary", sysdicid).Filter("TagName", tagname).Filter("Datatime__lte", beginTime).OrderBy("-Datatime").Limit(1)
		} else if readType == 0 { //读取时间范围内的数据
			qs = qs.Filter("SysDictionary", sysdicid).Filter("TagName", tagname).Filter("Datatime__gt", beginTime).Filter("Datatime__lte", endTime)
		} else {
			qs = qs.Filter("SysDictionary", sysdicid).Filter("TagName", tagname).Filter("Datatime__gt", endTime).OrderBy("Datatime").Limit(1)
		}
		rows, err := qs.All(&srtds)
		if err != nil {
			return srtdsexi, err
		}
		if rows == 0 { //没有读取到数据
			var sql string
			if readType < 0 { //读取时间点上的单点数据
				sql = fmt.Sprintf("SELECT * FROM sys_real_data WHERE sys_dictionary_id=%d AND tag_name = %q AND datatime <= %q LIMIT 1;", sysdicid, tagname, beginTime)
			} else if readType == 0 { //读取时间范围内的数据
				sql = fmt.Sprintf("SELECT * FROM sys_real_data WHERE sys_dictionary_id=%d AND tag_name = %q AND datatime > %q AND datatime <= %q;", sysdicid, tagname, beginTime, endTime)
			} else { //读取时间点后的单点数据
				sql = fmt.Sprintf("SELECT * FROM sys_real_data WHERE sys_dictionary_id=%d AND tag_name = %q AND datatime > %q LIMIT 1;", sysdicid, tagname, endTime)
			}
			return srtdsexi, fmt.Errorf("type[%s].tagname[%s] No data found in the set condition[在设定的条件下没有查询到数据]; The SQL is [%s]", tagtype, tagname, sql)
		}
	} else { //tagid有效，按照tagid读取数据
		srtd := new(SysRealData)
		qs := o.QueryTable(srtd)
		if readType < 0 { //读取时间点上的单点数据
			qs = qs.Filter("SysDictionary", sysdicid).Filter("TagID", tagid).Filter("Datatime__lte", beginTime).OrderBy("-Datatime").Limit(1)
		} else if readType == 0 { //读取时间范围内的数据
			qs = qs.Filter("SysDictionary", sysdicid).Filter("TagID", tagid).Filter("Datatime__gt", beginTime).Filter("Datatime__lte", endTime)
		} else { //读取大于结束时间的数据
			qs = qs.Filter("SysDictionary", sysdicid).Filter("TagID", tagid).Filter("Datatime__gt", endTime).OrderBy("Datatime").Limit(1)
		}
		rows, err := qs.All(&srtds)
		if err != nil {
			return srtdsexi, err
		}
		if rows == 0 { //没有读取到数据
			var sql string
			if readType < 0 {
				sql = fmt.Sprintf("SELECT * FROM sys_real_data WHERE sys_dictionary_id=%d AND tag_id = %d AND datatime <= %q LIMIT 1;", sysdicid, tagid, beginTime)
			} else if readType == 0 {

				sql = fmt.Sprintf("SELECT * FROM sys_real_data WHERE sys_dictionary_id=%d AND tag_id = %d AND datatime > %q AND datatime <= %q;", sysdicid, tagid, beginTime, endTime)
			} else {
				sql = fmt.Sprintf("SELECT * FROM sys_real_data WHERE sys_dictionary_id=%d AND tag_id = %d AND datatime > %q LIMIT 1;", sysdicid, tagid, endTime)
			}
			return srtdsexi, fmt.Errorf("type[%s].tag[%d] No data found in the set condition[在设定的条件下没有查询到数据]; The SQL is [%s]", tagtype, tagid, sql)
		}
	}
	var sr SysRealDataExi
	for _, v := range srtds { //整理成输出结构数据
		sr.Datatime = v.Datatime
		sr.Id = v.Id
		sr.SysDictionaryId = v.SysDictionary.Id
		sr.TagExeId = v.TagExeId
		sr.TagId = v.TagId
		sr.TagName = v.TagName
		sr.Value = v.Value
		srtdsexi = append(srtdsexi, sr)
	}
	return srtdsexi, nil
}

/*********************************************************************
功能:根据输入的字符串(DictionaryNameCode)查询在SysDictionary中的ID值
输入:namecode
输出:Id和错误信息
说明:如果namecode可转换为数字,则验证该数字是否是SysDictionary中的Id,验证成功直接返回该数据
编辑:wang_jp
时间:2020年1月12日
*********************************************************************/
func (dic *SysDictionary) GetSysDicIdByNameCode(namecode string) (int64, error) {
	o := orm.NewOrm() //新建orm对象
	o.Using("default")
	//dicc := SysDictionaryCatalog{Id: _variableTypeInSysDicCat}

	id, err := strconv.ParseInt(namecode, 10, 64)
	if err == nil { //能转换为整数
		sysdic := SysDictionary{Id: id, DicCatalogId: _variableTypeInSysDicCat}
		err := o.Read(&sysdic, "Id", "DicCatalogId")
		if err == orm.ErrNoRows {
			return 0, fmt.Errorf("[%s] is a invalid variable type value[数值非法,不是有效的变量类型ID]", namecode)
		}
		return id, nil //直接返回
	} //不能转换为整数则继续查询

	sysdic := SysDictionary{DictionaryNameCode: namecode, DicCatalogId: _variableTypeInSysDicCat}
	err = o.Read(&sysdic, "DictionaryNameCode", "DicCatalogId")
	if err == orm.ErrNoRows {
		return 0, fmt.Errorf("[%s] is a invalid variable type string[类型非法,不是有效的变量类型]", namecode)
	} else if err == orm.ErrMissPK {
		return 0, fmt.Errorf("Miss primary key")
	}
	if sysdic.Id == 0 {
		return sysdic.Id, fmt.Errorf("没有在sys_dictionary中找到dictionary_name_code为[%s]的值", namecode)
	}
	return sysdic.Id, nil
}

/*********************************************************************
功能:获取KPI TagType 列表
输入:无
输出:[[]SysDictionary] KPI TagType 列表
	[error] 错误信息
说明:
编辑:wang_jp
时间:2020年7月14日
*********************************************************************/
func (dic *SysDictionary) GetKpiTagTypeList() ([]SysDictionary, error) {
	o := orm.NewOrm() //新建orm对象
	o.Using("default")
	var sysdic []SysDictionary
	qt := o.QueryTable("SysDictionary").Filter("DicCatalogId", 7).Filter("DictionaryCode__gt", 0).Filter("Status", 1).OrderBy("Seq")
	_, err := qt.All(&sysdic)

	return sysdic, err
}

/*********************************************************************
功能:获取Grafana的访问地址
输入:[iswan ...bool],省略或者False为获取局域网地址,True为获取广域网地址
输出:[string] Grafana访问地址,如: http://127.0.0.1:3000
	[error] 错误信息
说明:
编辑:wang_jp
时间:2020年8月14日
*********************************************************************/
func (dic *SysDictionary) GetGrafanaHost(iswan ...bool) (string, error) {
	o := orm.NewOrm() //新建orm对象
	o.Using("default")
	dic_code := 1 //获取内网地址
	if len(iswan) > 0 {
		if iswan[0] {
			dic_code = 2 //获取外网地址
		}
	}
	qt := o.QueryTable("SysDictionary").Filter("DicCatalogId", 26).Filter("DictionaryCode", dic_code).Filter("Status", 1)
	err := qt.One(dic)

	host := dic.Name
	if err == nil {
		if len(host) > 1 {
			if host[len(host)-1] == '/' {
				host = host[:len(host)-1]
			}
		}
	}

	return host, err
}

/*********************************************************************
功能:获取MicEngine的访问地址
输入:[iswan ...bool],省略或者False为获取局域网地址,True为获取广域网地址
输出:[string] MicEngine访问地址,如: http://127.0.0.1:8080
	[error] 错误信息
说明:
编辑:wang_jp
时间:2020年8月14日
*********************************************************************/
func (dic *SysDictionary) GetMicEngineHost(iswan ...bool) (string, error) {
	o := orm.NewOrm() //新建orm对象
	o.Using("default")
	dic_code := 1 //获取内网地址
	if len(iswan) > 0 {
		if iswan[0] {
			dic_code = 2 //获取外网地址
		}
	}
	qt := o.QueryTable("SysDictionary").Filter("DicCatalogId", 13).Filter("DictionaryCode", dic_code).Filter("Status", 1)
	err := qt.One(dic)

	host := dic.Name
	if err == nil {
		if len(host) > 1 {
			if host[len(host)-1] == '/' {
				host = host[:len(host)-1]
			}
		}
	}

	return host, err
}

/*********************************************************************
功能:获取平台的访问地址
输入:[iswan ...bool],省略或者False为获取局域网地址,True为获取广域网地址
输出:[string] 平台访问地址,如: http://127.0.0.1:8080
	[error] 错误信息
说明:
编辑:wang_jp
时间:2020年8月14日
*********************************************************************/
func (dic *SysDictionary) GetPlatHost(iswan ...bool) (string, error) {
	o := orm.NewOrm() //新建orm对象
	o.Using("default")
	dic_code := 1 //获取内网地址
	if len(iswan) > 0 {
		if iswan[0] {
			dic_code = 2 //获取外网地址
		}
	}
	qt := o.QueryTable("SysDictionary").Filter("DicCatalogId", 27).Filter("DictionaryCode", dic_code).Filter("Status", 1)
	err := qt.One(dic)

	host := dic.Name
	if err == nil {
		if len(host) > 1 {
			if host[len(host)-1] == '/' {
				host = host[:len(host)-1]
			}
		}
	}

	return host, err
}
