package controllers

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego"
	"github.com/bkzy-wangjp/MicEngine/models"
)

type GrafanaApiController struct {
	beego.Controller
}

func (g *GrafanaApiController) Search() {
	micgd := new(models.MicGolden)
	res, _ := micgd.GoldenGetTables(2)
	value, _ := res.([]string)
	g.Data["json"] = value
	g.ServeJSON()
}

func (this *GrafanaApiController) GrafanaJsonQuery() {
	ob := new(grafanaRequestDataStruct)
	var err error
	if err = json.Unmarshal(this.Ctx.Input.RequestBody, ob); err == nil {
		responsData, err := grafanaRequestResponse(*ob)
		if err != nil {
			this.Data["json"] = err.Error()
		} else {
			this.Data["json"] = responsData //"{数据解析完成}"
		}
	} else {
		this.Data["json"] = err.Error()
	}
	this.ServeJSON()
}

/**********************************************
功能：Grafana从庚顿数据库中获取数据的响应
输入：请求数据的结构
输出：获取到的数据及错误信息
说明：
时间：2019年12月22日
编辑：wang_jp
**********************************************/
func grafanaRequestResponse(req grafanaRequestDataStruct) (interface{}, error) {
	//fmt.Printf("========请求数据:%+v\n", req)
	//时间转换为标准时间格式
	t_from, _ := strconv.ParseInt(req.ScopedVars.From.Text, 10, 64)
	t := time.Unix(t_from/1000, t_from%1000*1000)
	begineTimeStr := t.Format(models.EngineCfgMsg.Sys.TimeFormat)

	t_to, _ := strconv.ParseInt(req.ScopedVars.To.Text, 10, 64)
	t = time.Unix(t_to/1000, t_to%1000*1000)
	if req.Range.Raw.To == "now" {
		t = t.Add(time.Second * time.Duration(_DelaySecWhenNow))
	}
	endTimeStr := t.Format(models.EngineCfgMsg.Sys.TimeFormat)

	//时间间隔
	intervalMs := req.IntervalMs
	if intervalMs < 1000 {
		intervalMs = 1000
	}
	//开始遍历循环查看目标指令
	var responsData []interface{}
	for _, taget := range req.Targets {
		order := strings.ToLower(taget.Data.Order) //读取指令
		table := strings.ToLower(taget.Target)     //数据表名称
		if len(table) == 0 {
			return nil, fmt.Errorf("No valid table name selected[没有选择有效的数据表名]")
		}
		responsType := strings.ToLower(taget.Type) //需要的反馈类型,timeseries 或者 table
		var l = 1
		if len(taget.Data.TagNames) > 0 {
			l = len(taget.Data.TagNames)
		}
		tagNames := make([]string, l)                  //变量名
		tagDescs := make([]string, l)                  //变量描述
		if strings.Contains(order, "goods") == false { //不是读取物耗，按照庚顿格式整理变量名和标签
			var readTagDesc bool
			if len(taget.Data.TagNames) > 0 { //变量名数组有效
				for i, tag := range taget.Data.TagNames {
					tagNames[i] = fmt.Sprintf("%s.%s", table, tag)
				}
				if len(taget.Data.TagDescs) != len(taget.Data.TagNames) { //如果设定的变量描述与设定的变量名不一致
					readTagDesc = true //需要读取变量描述
				} else { //不需要读取变量名
					for j, desc := range taget.Data.TagDescs {
						tagDescs[j] = desc
					}
				}
			} else if len(taget.Data.TagName) > 0 { //单一变量名有效
				tagNames[0] = fmt.Sprintf("%s.%s", table, taget.Data.TagName)
				if len(taget.Data.TagDesc) == 0 { //如果没有设置变量描述
					readTagDesc = true //需要读取变量描述
				} else {
					tagDescs[0] = taget.Data.TagDesc
				}
			} else {
				return nil, fmt.Errorf("No valid variable name set[没有设置有效的变量名]")
			}
			if readTagDesc == true {
				tagDescs = getTagDescFromDB(tagNames) //获取变量描述
			}
		} else { //读取物耗
			if len(taget.Data.TagNames) > 0 { //变量名数组有效
				for i, tag := range taget.Data.TagNames {
					tagNames[i] = fmt.Sprintf("%s", strings.ToLower(tag))
				}
				if len(taget.Data.TagDescs) != len(taget.Data.TagNames) { //如果设定的变量描述与设定的变量名不一致
					tagDescs = tagNames //描述等于变量名
				} else { //不需要读取变量名
					for j, desc := range taget.Data.TagDescs {
						tagDescs[j] = desc
					}
				}
			} else if len(taget.Data.TagName) > 0 { //单一变量名有效
				tagNames[0] = fmt.Sprintf("%s", strings.ToLower(taget.Data.TagName))
				if len(taget.Data.TagDesc) == 0 { //如果没有设置变量描述
					tagDescs = tagNames //描述等于变量名
				} else {
					tagDescs[0] = taget.Data.TagDesc
				}
			} else {
				return nil, fmt.Errorf("No valid variable name set[没有设置有效的变量名]")
			}
		}

		interval := taget.Data.Interval
		if interval < 0 { //间隔设置小于0
			interval = intervalMs / 1000 //取自动间隔
		}
		switch order {
		case "goods", "物耗":
			timePointStr := ""
			if len(taget.Data.TimePoint) > 10 {
				t, err := models.TimeParse(taget.Data.TimePoint)
				if err != nil {
					return nil, fmt.Errorf("The TimePoint formate error[时间点格式错误];[%s]", err.Error())
				}
				timePointStr = t.Format(models.EngineCfgMsg.Sys.TimeFormat)
			}
			respData, err := grafanaGoodsDataRes(tagNames, tagDescs, responsType, begineTimeStr, endTimeStr, timePointStr)
			if err != nil {
				return nil, err
			} else {
				responsData = append(responsData, respData...)
			}
		case "goodssum", "物耗累积":
			respData, err := grafanaGoodsDataRes(tagNames, tagDescs, responsType, begineTimeStr, endTimeStr, "")
			if err != nil {
				return nil, err
			} else {
				responsData = append(responsData, respData...)
			}
		case "snapshot", "快照":
			respData, err := getSnapShotDataFromGoldenDB(tagNames, tagDescs, responsType)
			if err != nil {
				return nil, err
			} else {
				responsData = append(responsData, respData...)
			}
		case "historyinterval", "历史", "history":
			timePointStr := ""
			if len(taget.Data.TimePoint) > 10 {
				t, err := models.TimeParse(taget.Data.TimePoint)
				if err != nil {
					return nil, fmt.Errorf("The TimePoint formate error[时间点格式错误];[%s]", err.Error())
				}
				timePointStr = t.Format(models.EngineCfgMsg.Sys.TimeFormat)
			}

			respData, err := getHistoryDataFromGoldenDB(tagNames, tagDescs, begineTimeStr, endTimeStr, timePointStr, responsType, interval)
			if err != nil {
				return nil, err
			} else {
				responsData = append(responsData, respData...)
			}
		case "historysummary", "统计", "advangce", "min", "max", "range", "total", "sum", "mean", "poweravg", "diff", "duration", "pointcnt", "risingcnt", "fallingcnt", "sd", "stddev", "se", "ske", "kur", "mode", "median", "distribution", "groupdist", "increment":
			respData, err := getHistorySummaryDataFromGoldenDB(tagNames, tagDescs, order, begineTimeStr, endTimeStr, responsType)
			if err != nil {
				return nil, err
			} else {
				responsData = append(responsData, respData...)
			}
		case "lt_l_total", "lte_l_total", "gt_l_total", "gte_l_total",
			"lt_ll_total", "lte_ll_total", "gt_ll_total", "gte_ll_total",
			"lt_h_total", "lte_h_total", "gt_h_total", "gte_h_total",
			"lt_hh_total", "lte_hh_total", "gt_hh_total", "gte_hh_total",
			"lt_min_total", "lte_min_total", "gt_min_total", "gte_min_total",
			"lt_max_total", "lte_max_total", "gt_max_total", "gte_max_total":
			respData, err := dataCompareTotalForGrafana(tagNames, tagDescs, order, begineTimeStr, endTimeStr, responsType)
			if err != nil {
				return nil, err
			} else {
				responsData = append(responsData, respData...)
			}
		default:
			return nil, fmt.Errorf("The 'Order' is not valid.['Order'命令设置的不正确]")
		}
	}
	return responsData, nil
}

/**********************************************
功能：获取统计数据
输入：变量名数组，变量描述，开始时间点，结束时间点，历史时刻点，反馈类型，间隔时间
输出：获取到的数据及错误信息
说明：
时间：2019年12月22日
编辑：wang_jp
**********************************************/
func getHistorySummaryDataFromGoldenDB(tagnames, tagDescs []string, order, begineTime, endTime, reqType string) ([]interface{}, error) {
	var errCnt int
	micgd := new(models.MicGolden)
ReadGolden:
	dmap, err := micgd.GoldenGetHistorySummary(begineTime, endTime, tagnames...)
	if err != nil { //如果有错误
		errCnt += 1
		if strings.Contains(err.Error(), "snapshot stoped,the last time is") { //判断是否是快照时间未到的错误
			reg := regexp.MustCompile(`\d{4}\-\d{1,2}\-\d{1,2}\s\d{1,2}\:\d{1,2}\:\d{1,2}\.*\d*`) //提取最后时间
			endstr := reg.FindAllString(err.Error(), -1)
			for _, tstr := range endstr {
				endTime = tstr //重置最后时间
			}
			//beego.Debug("查询统计时间修正完成,修正前:", endTime, "修正后:", cmd.EndTime, "快照时间值:", endstr)
			if errCnt <= len(tagnames) {
				goto ReadGolden //重新读取数据
			}
			return nil, err
		} else { //如果不是,返回错误
			return nil, err
		}
	}
	hismap := dmap

	var hisTable []interface{}
	type SummaryTableTypeData struct { //表格形式的统计数据反馈结构
		Columns [23]gfTableColumn `json:"columns"`
		Rows    [][23]interface{} `json:"rows"`
		Type    string            `json:"type"` //固定返回值"table"
	}
	switch order {
	case "historysummary", "统计", "advangce": //返回所有统计值
		var tb SummaryTableTypeData
		for i, _ := range tb.Columns {
			tb.Columns[i].Type = "number"
		}
		tb.Columns[0].Text = "BeginTime"
		tb.Columns[0].Type = "time"
		tb.Columns[1].Text = "EndTime"
		tb.Columns[1].Type = "time"
		tb.Columns[2].Text = "TagName"
		tb.Columns[2].Type = "string"
		tb.Columns[3].Text = "Duration"
		tb.Columns[4].Text = "Min"
		tb.Columns[5].Text = "Max"
		tb.Columns[6].Text = "Range"
		tb.Columns[7].Text = "Total"
		tb.Columns[8].Text = "Sum"
		tb.Columns[9].Text = "Mean"
		tb.Columns[10].Text = "PowerAvg"
		tb.Columns[11].Text = "Mode"
		tb.Columns[12].Text = "Median"
		tb.Columns[13].Text = "Diff"
		tb.Columns[14].Text = "PointCnt"
		tb.Columns[15].Text = "RisingCnt"
		tb.Columns[16].Text = "FallingCnt"
		tb.Columns[17].Text = "SD"
		tb.Columns[18].Text = "STDDEV"
		tb.Columns[19].Text = "SE"
		tb.Columns[20].Text = "Ske"
		tb.Columns[21].Text = "Kur"
		tb.Columns[22].Text = "GroupDist"
		tb.Type = "table"

		for i, tag := range tagnames { //遍历原始变量数组，而不是结果数组，防止有变量而没有结果的情况
			var stb SummaryTableTypeData
			stb = tb
			for k, v := range hismap {
				if strings.ToLower(k) == strings.ToLower(tag) {
					var av [23]interface{}
					t, _ := models.TimeParse(v.BeginTime) //时间转换
					tv := t.Unix() * 1000
					av[0] = tv //时间

					t, _ = models.TimeParse(v.EndTime) //时间转换
					tv = t.Unix() * 1000
					av[1] = tv //时间

					av[2] = tagDescs[i] //变量名
					av[3] = v.Duration
					av[4] = v.Min
					av[5] = v.Max
					av[6] = v.Range
					av[7] = v.Total
					av[8] = v.Sum
					av[9] = v.Mean
					av[10] = v.PowerAvg
					av[11] = v.Mode
					av[12] = v.Median
					av[13] = v.Diff
					av[14] = v.PointCnt
					av[15] = v.RisingCnt
					av[16] = v.FallingCnt
					av[17] = v.SD
					av[18] = v.STDDEV
					av[19] = v.SE
					av[20] = v.Ske
					av[21] = v.Kur
					av[22] = v.GroupDist
					stb.Rows = append(stb.Rows, av)
					break
				}
			}
			hisTable = append(hisTable, stb)
		}
		return hisTable, nil
	case "distribution": //返回数据分布曲线
		var tb gfTableTypeData
		tb.Columns[0].Text = "TagName"
		tb.Columns[0].Type = "string"
		tb.Columns[1].Text = "Value"
		tb.Columns[1].Type = "string"
		tb.Columns[2].Text = "Count"
		tb.Columns[2].Type = "number"
		tb.Type = "table"
		for i, tag := range tagnames { //遍历原始变量数组，而不是结果数组，防止有变量而没有结果的情况
			var stb gfTableTypeData
			stb = tb
			for k, hisv := range hismap {
				if strings.ToLower(k) == strings.ToLower(tag) {
					for j, v := range hisv.DataGroup {
						var av [3]interface{}
						av[0] = tagDescs[i] //变量名
						av[1] = j           //fmt.Sprintf("%.2f", hisv.GroupDist*float64(j+1))
						av[2] = v           //值
						stb.Rows = append(stb.Rows, av)
					}
					break
				}
			}
			hisTable = append(hisTable, stb)
		}
		return hisTable, nil
	case "increment": //返回相邻两点之间的增量
		switch reqType {
		case "timeseries": //返回时间序列数据
			var series []interface{}
			for i, tag := range tagnames { //遍历原始变量数组，而不是结果数组，防止有变量而没有结果的情况
				var sdata gfTimeSeriesTypeData
				sdata.Target = tagDescs[i]
				for k, hisv := range hismap {
					if strings.ToLower(k) == strings.ToLower(tag) {
						for t, v := range hisv.Increment {
							var tv [2]interface{}
							tm, err := models.TimeParse(t) //时间转换
							if err != nil {
								beego.Warn("The tag [%s] history time fomate error;[%s]", tag, err.Error())
							}
							tv[1] = tm.Unix() * 1000
							tv[0] = v
							sdata.Datapoints = append(sdata.Datapoints, tv)
						}
						break
					}
				}
				series = append(series, sdata)
			}
			return series, nil
		default:
			var tb gfTableTypeData
			tb.Columns[0].Text = "Time"
			tb.Columns[0].Type = "time"
			tb.Columns[1].Text = "TagName"
			tb.Columns[1].Type = "string"
			tb.Columns[2].Text = "Value"
			tb.Columns[2].Type = "number"
			tb.Type = "table"
			for i, tag := range tagnames { //遍历原始变量数组，而不是结果数组，防止有变量而没有结果的情况
				var stb gfTableTypeData
				stb = tb
				for k, hisv := range hismap {
					if strings.ToLower(k) == strings.ToLower(tag) {
						for t, v := range hisv.Increment {
							var av [3]interface{}
							tm, err := models.TimeParse(t) //时间转换
							if err != nil {
								beego.Warn("The tag [%s] history time fomate error;[%s]", tag, err.Error())
							}
							av[0] = tm.Unix() * 1000 //时间
							av[1] = tagDescs[i]      //变量名
							av[2] = v                //值
							stb.Rows = append(stb.Rows, av)
						}
						break
					}
				}
				hisTable = append(hisTable, stb)
			}
			return hisTable, nil
		}
	default: //其他单值统计
		switch reqType {
		case "timeseries": //返回时间序列数据
			var series []interface{}
			for i, tag := range tagnames { //遍历原始变量数组，而不是结果数组，防止有变量而没有结果的情况
				var sdata gfTimeSeriesTypeData
				sdata.Target = tagDescs[i]
				for k, v := range hismap {
					if strings.ToLower(k) == strings.ToLower(tag) {
						var tv [2]interface{}
						tm, err := models.TimeParse(v.EndTime) //时间转换
						if err != nil {
							beego.Warn("变量 [%s] 历史数据时间格式错误:[%s]", tag, err.Error())
						}
						tv[1] = tm.Unix() * 1000
						tv[0] = micgd.GetSingleDataFromStatistics(v, order) //Golden.SelectValueFromHisSummryExi(v, order)
						sdata.Datapoints = append(sdata.Datapoints, tv)
						break
					}
				}
				series = append(series, sdata)
			}
			return series, nil
		default:
			var tb gfTableTypeData
			tb.Columns[0].Text = "Time"
			tb.Columns[0].Type = "time"
			tb.Columns[1].Text = "TagName"
			tb.Columns[1].Type = "string"
			tb.Columns[2].Text = "Value"
			tb.Columns[2].Type = "number"
			tb.Type = "table"
			for i, tag := range tagnames { //遍历原始变量数组，而不是结果数组，防止有变量而没有结果的情况
				var stb gfTableTypeData
				stb = tb
				for k, v := range hismap {
					if strings.ToLower(k) == strings.ToLower(tag) {
						var av [3]interface{}
						tm, err := models.TimeParse(v.EndTime) //时间转换
						if err != nil {
							beego.Warn("The tag [%s] history time fomate error;[%s]", tag, err.Error())
						}
						av[0] = tm.Unix() * 1000                            //时间
						av[1] = tagDescs[i]                                 //变量名
						av[2] = micgd.GetSingleDataFromStatistics(v, order) //值
						stb.Rows = append(stb.Rows, av)
						break
					}
				}
				hisTable = append(hisTable, stb)
			}
			return hisTable, nil
		}
	}
	return nil, fmt.Errorf("待完善……")
}

/***********************************************
功能：将获取的物耗数据转换为grafana接收的格式
输入：标签名数组，时间点,开始时间，结束时间
输出：获取到的数据及错误信息
说明：
时间：2019年12月22日
编辑：wang_jp
***********************************************/
func grafanaGoodsDataRes(tagnames, tagDescs []string, reqType string, beginTime, endTime, timePoint string) ([]interface{}, error) {
	goodsData, err := getGoodsDataFromDB(tagnames, beginTime, endTime, timePoint)
	if err != nil {
		return nil, err
	}
	var resdata []interface{}
	if len(timePoint) > 0 { //读取的时间点数据
		goodsv := goodsData.(map[string]goodsPointData)
		switch reqType {
		case "timeseries": //返回时间序列数据
			for i, tag := range tagnames { //遍历原始变量数组，而不是结果数组，防止有变量而没有结果的情况
				var sns gfTimeSeriesTypeData
				sns.Target = tagDescs[i]
				for k, v := range goodsv {
					if strings.ToLower(k) == strings.ToLower(tag) {
						var tv [2]interface{}
						t, err := models.TimeParse(v.Time) //时间转换
						if err != nil {
							beego.Warn("The goods tag [%s] time fomate error;[%s]", tag, err.Error())
						}
						tv[1] = t.Unix() * 1000
						tv[0] = v.Value
						sns.Datapoints = append(sns.Datapoints, tv)
						break
					}
				}
				resdata = append(resdata, sns)
			}
			return resdata, nil
		case "table": //返回表格形式数据
			var tb gfTableTypeData
			tb.Columns[0].Text = "Time"
			tb.Columns[0].Type = "time"
			tb.Columns[1].Text = "TagName"
			tb.Columns[1].Type = "string"
			tb.Columns[2].Text = "Value"
			tb.Columns[2].Type = "number"
			tb.Type = "table"
			tb.Rows = make([][3]interface{}, 1)

			for i, tag := range tagnames { //遍历原始变量数组，而不是结果数组，防止有变量而没有结果的情况
				var stb gfTableTypeData
				stb = tb
				for k, v := range goodsv {
					if strings.ToLower(k) == strings.ToLower(tag) {
						var tv, vl interface{}
						t, err := models.TimeParse(v.Time) //时间转换
						if err != nil {
							beego.Warn("The tag [%s] snapshot time fomate error;[%s]", tag, err.Error())
						}
						tv = t.Unix() * 1000
						stb.Rows[0][0] = tv          //时间
						stb.Rows[0][1] = tagDescs[i] //变量名
						vl = v.Value
						stb.Rows[0][2] = vl //值
						break
					}
				}
				resdata = append(resdata, stb)
			}
			return resdata, nil
		default:
			return nil, fmt.Errorf("Output foamet select error,must be 'timeseries' or 'table'")
		}
	} else { //获取的是时间范围内定累积数据
		goodsv := goodsData.(map[string]goodsSumData)
		switch reqType {
		case "timeseries": //返回时间序列数据
			for i, tag := range tagnames { //遍历原始变量数组，而不是结果数组，防止有变量而没有结果的情况
				var sns gfTimeSeriesTypeData
				sns.Target = tagDescs[i]
				for k, v := range goodsv {
					if strings.ToLower(k) == strings.ToLower(tag) {
						var tv [2]interface{}
						t, err := models.TimeParse(v.EndTime) //时间转换
						if err != nil {
							beego.Warn("The goods tag [%s] time fomate error;[%s]", tag, err.Error())
						}
						tv[1] = t.Unix() * 1000
						tv[0] = v.Sum
						sns.Datapoints = append(sns.Datapoints, tv)
						break
					}
				}
				resdata = append(resdata, sns)
			}
			return resdata, nil
		case "table": //返回表格形式数据
			var tb gfTableTypeData
			tb.Columns[0].Text = "Time"
			tb.Columns[0].Type = "time"
			tb.Columns[1].Text = "TagName"
			tb.Columns[1].Type = "string"
			tb.Columns[2].Text = "Value"
			tb.Columns[2].Type = "number"
			tb.Type = "table"
			tb.Rows = make([][3]interface{}, 1)

			for i, tag := range tagnames { //遍历原始变量数组，而不是结果数组，防止有变量而没有结果的情况
				var stb gfTableTypeData
				stb = tb
				for k, v := range goodsv {
					if strings.ToLower(k) == strings.ToLower(tag) {
						var tv, vl interface{}
						t, err := models.TimeParse(v.EndTime) //时间转换
						if err != nil {
							beego.Warn("The tag [%s] snapshot time fomate error;[%s]", tag, err.Error())
						}
						tv = t.Unix() * 1000
						stb.Rows[0][0] = tv          //时间
						stb.Rows[0][1] = tagDescs[i] //变量名
						vl = v.Sum
						stb.Rows[0][2] = vl //值
						break
					}
				}
				resdata = append(resdata, stb)
			}
			return resdata, nil
		default:
			return nil, fmt.Errorf("Output foamet select error,must be 'timeseries' or 'table'")
		}
	}
	return nil, nil
}

/**********************************************
功能：数据比较后的结果求累积面积
输入：请求数据的结构
输出：获取到的数据及错误信息
说明：
时间：2020年1月3日
编辑：wang_jp
**********************************************/
func dataCompareTotal(input getCmd) (map[string]float64, error) {
	res := make(map[string]float64, len(input.TagNames))
	if len(input.TagNames) == 0 {
		return nil, fmt.Errorf("Tagnames is a required parameter[Tagnames 是必须的参数]")
	}
	if len(input.BeginTime) == 0 {
		return nil, fmt.Errorf("BeginTime is a required parameter[BeginTime 是必须的参数]")
	}
	if len(input.EndTime) == 0 {
		return nil, fmt.Errorf("EndTime is a required parameter[EndTime 是必须的参数]")
	}
	for _, tagname := range input.TagNames {
		tag := new(models.OreProcessDTaglist)
		tagid, err := tag.GetTagIDByFullName(tagname)
		if err != nil {
			return nil, err
		}
		micgd := new(models.MicGolden)
		v, _, err := micgd.GoldenGetAnalogComparedTotal(tagid, tagname, input.Order, input.BeginTime, input.EndTime)
		if err != nil {
			return nil, err
		}
		res[tagname] = v
	}
	return res, nil
}

/**********************************************
功能：数据比较后的结果求累积面积,为Grafana响应
输入：请求数据的结构
输出：获取到的数据及错误信息
说明：
时间：2020年1月3日
编辑：wang_jp
**********************************************/
func dataCompareTotalForGrafana(tagNames, tagDescs []string, order, begineTimeStr, endTimeStr, reqType string) ([]interface{}, error) {
	res := make(map[string]float64, len(tagNames))
	micgd := new(models.MicGolden)
	for _, tagname := range tagNames {
		tag := new(models.OreProcessDTaglist)
		tagid, err := tag.GetTagIDByFullName(tagname)
		if err != nil {
			return nil, err
		}
		errCnt := 0
	ReadGolden:
		v, _, err := micgd.GoldenGetAnalogComparedTotal(tagid, tagname, order, begineTimeStr, endTimeStr)
		if err != nil { //如果有错误
			errCnt += 1
			if strings.Contains(err.Error(), "snapshot stoped,the last time is") { //判断是否是快照时间未到的错误
				reg := regexp.MustCompile(`\d{4}\-\d{1,2}\-\d{1,2}\s\d{1,2}\:\d{1,2}\:\d{1,2}\.*\d*`) //提取最后时间
				endstr := reg.FindAllString(err.Error(), -1)
				for _, tstr := range endstr {
					endTimeStr = tstr //重置最后时间
				}
				if errCnt <= 1 {
					goto ReadGolden //重新读取数据
				}
				return nil, err
			} else { //如果不是,返回错误
				return nil, err
			}
		}
		res[tagname] = v
	}
	var resdata []interface{}
	switch reqType {
	case "timeseries": //返回时间序列数据
		for i, tag := range tagNames { //遍历原始变量数组，而不是结果数组，防止有变量而没有结果的情况
			var sns gfTimeSeriesTypeData
			sns.Target = tagDescs[i]
			for k, v := range res {
				if strings.ToLower(k) == strings.ToLower(tag) {
					var tv [2]interface{}
					t, err := models.TimeParse(endTimeStr) //时间转换
					if err != nil {
						beego.Warn("The tag [%s] time fomate error;[%s]", tag, err.Error())
					}
					tv[1] = t.Unix() * 1000
					tv[0] = v
					sns.Datapoints = append(sns.Datapoints, tv)
					break
				}
			}
			resdata = append(resdata, sns)
		}
		return resdata, nil
	case "table": //返回表格形式数据
		var tb gfTableTypeData
		tb.Columns[0].Text = "Time"
		tb.Columns[0].Type = "time"
		tb.Columns[1].Text = "TagName"
		tb.Columns[1].Type = "string"
		tb.Columns[2].Text = "Value"
		tb.Columns[2].Type = "number"
		tb.Type = "table"
		tb.Rows = make([][3]interface{}, 1)

		for i, tag := range tagNames { //遍历原始变量数组，而不是结果数组，防止有变量而没有结果的情况
			var stb gfTableTypeData
			stb = tb
			for k, v := range res {
				if strings.ToLower(k) == strings.ToLower(tag) {
					var tv, vl interface{}
					t, err := models.TimeParse(endTimeStr) //时间转换
					if err != nil {
						beego.Warn("The tag [%s] time fomate error;[%s]", tag, err.Error())
					}
					tv = t.Unix() * 1000
					stb.Rows[0][0] = tv          //时间
					stb.Rows[0][1] = tagDescs[i] //变量名
					vl = v
					stb.Rows[0][2] = vl //值
					break
				}
			}
			resdata = append(resdata, stb)
		}
		return resdata, nil
	}
	return resdata, nil
}

/**********************************************
功能：获取历史时刻点数据
输入：变量名数组，变量描述，历史时刻点，反馈类型
输出：获取到的数据及错误信息
说明：
时间：2020年5月17日
编辑：wang_jp
**********************************************/
func getHistoryPointDataFromGoldenDB(tagnames, tagDescs []string, timepoint, reqType string) ([]interface{}, error) {
	micgd := new(models.MicGolden)
	hispmp, err := micgd.GoldenGetHistorySinglePoint(1, timepoint, tagnames...)
	if err != nil {
		return nil, err
	}

	switch reqType {
	case "timeseries": //返回时间序列数据
		var series []interface{}
		for i, tag := range tagnames { //遍历原始变量数组，而不是结果数组，防止有变量而没有结果的情况
			var sdata gfTimeSeriesTypeData
			sdata.Target = tagDescs[i]
			for k, v := range hispmp {
				if strings.ToLower(k) == strings.ToLower(tag) {
					var tv [2]interface{}
					tv[1] = v.Time
					tv[0] = v.Value
					sdata.Datapoints = append(sdata.Datapoints, tv)
					break
				}
			}
			series = append(series, sdata)
		}
		return series, nil
	case "table": //返回表格形式数据
		var hisTable []interface{}
		var tb gfTableTypeData
		tb.Columns[0].Text = "Time"
		tb.Columns[0].Type = "time"
		tb.Columns[1].Text = "TagName"
		tb.Columns[1].Type = "string"
		tb.Columns[2].Text = "Value"
		tb.Columns[2].Type = "number"
		tb.Type = "table"

		for i, tag := range tagnames { //遍历原始变量数组，而不是结果数组，防止有变量而没有结果的情况
			var stb gfTableTypeData
			stb = tb
			for k, v := range hispmp {
				if strings.ToLower(k) == strings.ToLower(tag) {
					var av [3]interface{}
					av[0] = v.Time      //时间
					av[1] = tagDescs[i] //变量名
					av[2] = v.Value     //值
					stb.Rows = append(stb.Rows, av)
					break
				}
			}
			hisTable = append(hisTable, stb)
		}
		return hisTable, nil
	default:
		return nil, fmt.Errorf("输出格式选择错误,必须是'timeseries' 或者 'table'")
	}
	return nil, fmt.Errorf("待完善……")
}

/**********************************************
功能：获取历史时间范围数据
输入：变量名数组，变量描述，开始时间点，结束时间点，历史时刻点，反馈类型，间隔时间
输出：获取到的数据及错误信息
说明：
时间：2019年12月22日
编辑：wang_jp
**********************************************/
func getHistoryRangeDataFromGoldenDB(tagnames, tagDescs []string, beginTime, endTime, reqType string, interval int64) ([]interface{}, error) {
	micgd := new(models.MicGolden)
	hismap, _, err := micgd.GoldenGetHistoryInterval(beginTime, endTime, interval, tagnames...)
	if err != nil {
		return nil, err
	}

	switch reqType {
	case "timeseries": //返回时间序列数据
		var series []interface{}
		for i, tag := range tagnames { //遍历原始变量数组，而不是结果数组，防止有变量而没有结果的情况
			var sdata gfTimeSeriesTypeData
			sdata.Target = tagDescs[i]
			for k, hisv := range hismap {
				if strings.ToLower(k) == strings.ToLower(tag) {
					for _, v := range hisv {
						var tv [2]interface{}
						t, err := models.TimeParse(v.Time) //时间转换
						if err != nil {
							beego.Warn("变量 [%s] 历史数据时间格式错误:[%s]", tag, err.Error())
						}
						tv[1] = t.UnixNano() / 1e6
						tv[0] = v.Value
						sdata.Datapoints = append(sdata.Datapoints, tv)
					}
					break
				}
			}
			series = append(series, sdata)
		}
		return series, nil
	case "table": //返回表格形式数据
		var hisTable []interface{}
		var tb gfTableTypeData
		tb.Columns[0].Text = "Time"
		tb.Columns[0].Type = "time"
		tb.Columns[1].Text = "TagName"
		tb.Columns[1].Type = "string"
		tb.Columns[2].Text = "Value"
		tb.Columns[2].Type = "number"
		tb.Type = "table"
		for i, tag := range tagnames { //遍历原始变量数组，而不是结果数组，防止有变量而没有结果的情况
			var stb gfTableTypeData
			stb = tb
			for k, hisv := range hismap {
				if strings.ToLower(k) == strings.ToLower(tag) {
					for _, v := range hisv {
						var av [3]interface{}
						t, err := models.TimeParse(v.Time) //时间转换
						if err != nil {
							beego.Warn("变量 [%s] 历史数据时间格式错误:[%s]", tag, err.Error())
						}
						av[0] = t.UnixNano() / 1e6 //时间
						av[1] = tagDescs[i]        //变量名
						av[2] = v.Value            //值
						stb.Rows = append(stb.Rows, av)
					}
					break
				}
			}
			hisTable = append(hisTable, stb)
		}
		return hisTable, nil
	default:
		return nil, fmt.Errorf("输出格式选择错误,必须是'timeseries' 或者 'table'")
	}
	return nil, fmt.Errorf("待完善……")
}

/***********************************************
功能：从庚顿数据库中获取快照数据
输入：标签名数组，需要的反馈类型
输出：获取到的数据及错误信息
说明：
时间：2019年12月22日
编辑：wang_jp
***********************************************/
func getSnapShotDataFromGoldenDB(tagnames, tagDescs []string, reqType string) ([]interface{}, error) {
	micgd := new(models.MicGolden)
	snap, err := micgd.GoldenGetSnapShotMap(tagnames...)
	if err != nil {
		return nil, err
	}

	switch reqType {
	case "timeseries": //返回时间序列数据
		var snapSeries []interface{}
		for i, tag := range tagnames { //遍历原始变量数组，而不是结果数组，防止有变量而没有结果的情况
			var sns gfTimeSeriesTypeData
			sns.Target = tagDescs[i]
			for k, v := range snap {
				if strings.ToLower(k) == strings.ToLower(tag) {
					var tv [2]interface{}
					tv[1] = v.Rtsd.Time
					tv[0] = v.Rtsd.Value
					sns.Datapoints = append(sns.Datapoints, tv)
					break
				}
			}
			snapSeries = append(snapSeries, sns)
		}
		return snapSeries, nil
	case "table": //返回表格形式数据
		var snaTable []interface{}
		var tb gfTableTypeData
		tb.Columns[0].Text = "Time"
		tb.Columns[0].Type = "time"
		tb.Columns[1].Text = "TagName"
		tb.Columns[1].Type = "string"
		tb.Columns[2].Text = "Value"
		tb.Columns[2].Type = "number"
		tb.Type = "table"
		tb.Rows = make([][3]interface{}, 1)

		for i, tag := range tagnames { //遍历原始变量数组，而不是结果数组，防止有变量而没有结果的情况
			var stb gfTableTypeData
			stb = tb
			for k, v := range snap {
				if strings.ToLower(k) == strings.ToLower(tag) {
					stb.Rows[0][0] = v.Rtsd.Time  //时间
					stb.Rows[0][1] = tagDescs[i]  //变量名
					stb.Rows[0][2] = v.Rtsd.Value //值
					break
				}
			}
			snaTable = append(snaTable, stb)
		}
		return snaTable, nil
	default:
		return nil, fmt.Errorf("输出格式选择错误,必须是'timeseries' 或者 'table'")
	}
	return nil, nil
}

/***********************************************
功能：从庚顿数据库中获取标签的描述信息
输入：标签名数组
输出：标签描述信息
说明：1.如果设置的标签描述数量与变量组的数量没有差异，不需要再从数据库中读取
	2.如果数据库中没有设置变量的描述信息，取变量标签名为描述信息(不含表名)
	3.如果在数据库中没有找到正确的变量，则返回的变量描述信息为带表名的标签名，需要检查这个标签哪里错了
时间：2019年12月22日
编辑：wang_jp
***********************************************/
func getTagDescFromDB(tagfullnames []string) []string {
	var tagDescs []string
	var tagdesc string
	for _, tagfullname := range tagfullnames { //遍历所有
		tagdesc = tagfullname                       //默认描述等于变量全名
		tagnames := strings.Split(tagfullname, ".") //分割，去除表名
		if len(tagnames) > 1 {
			tagdesc = tagnames[1]                 //默认描述等于变量名
			tag := new(models.OreProcessDTaglist) //新建tag对象
			tag.TagName = tagnames[1]
			ldesc, sdesc, err := tag.GetTagDescByTagName() //获取tag描述
			if err == nil {                                //无错误
				if len(sdesc) > 1 { //设置了短描述
					tagdesc = sdesc //采用短描述
				} else if len(ldesc) > 1 { //没有设置短描述,而设置了长描述
					tagdesc = ldesc //采用长描述
				}
			}
		}
		tagDescs = append(tagDescs, tagdesc)
	}
	return tagDescs
}

/**********************************************
功能：获取历史数据
输入：变量名数组，变量描述，开始时间点，结束时间点，历史时刻点，反馈类型，间隔时间
输出：获取到的数据及错误信息
说明：
时间：2019年12月22日
编辑：wang_jp
**********************************************/
func getHistoryDataFromGoldenDB(tagnames, tagDescs []string, beginTime, endTime, timepoint, reqType string, interval int64) ([]interface{}, error) {
	if len(timepoint) > 0 { //如果读取的是时间点数据
		return getHistoryPointDataFromGoldenDB(tagnames, tagDescs, timepoint, reqType)
	} else { //否则
		return getHistoryRangeDataFromGoldenDB(tagnames, tagDescs, beginTime, endTime, reqType, interval)
	}
}
