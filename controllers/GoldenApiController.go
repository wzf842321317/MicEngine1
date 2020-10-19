package controllers

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/astaxie/beego"
	"github.com/bkzy-wangjp/MicEngine/models"
)

type GoldenApiController struct {
	beego.Controller
}

func (g *GoldenApiController) Get() {
	g.Data["HeaderTitle"] = "数据计算服务引擎"
	g.Data["Website"] = "mining-icloud.com"
	g.Data["Email"] = "xuning@bksmartmining.com"
	g.Data["Version"] = models.EngineCfgMsg.Version
	g.Data["Copyright"] = models.EngineCfgMsg.CfgMsg.Copyright
	g.TplName = "GoldenApi/index.tpl"
}

func (this *GoldenApiController) GetQuery() {
	if cmd, err := discernGetParemater(this.Input()); err != nil {
		this.Data["json"] = err.Error()
	} else {
		responsData, err := checkGetCmd(cmd)
		if err != nil {
			this.Data["json"] = err.Error()
		} else {
			this.Data["json"] = responsData
		}
	}
	this.ServeJSON()
}

/**********************************************
功能：识别Get参数
输入：url Values
输出：GetCmd,err
说明：
时间：2019年12月27日
编辑：wang_jp
**********************************************/
func discernGetParemater(input url.Values) (getCmd, error) {
	var cmd getCmd
	var err error
	for k, v := range input {
		switch strings.ToLower(k) {
		case "order":
			for i, od := range v {
				if i == 0 {
					cmd.Order = strings.ToLower(od)
				} else {
					err = fmt.Errorf("Too many order parameters[Order参数太多了]")
					return cmd, err
				}
			}
			if len(v) == 0 {
				err = fmt.Errorf("No Order parameter[没有 order 参数]")
				return cmd, err
			}
		case "tagnames":
			for _, tags := range v {
				s := strings.Split(tags, ",") //分割
				for _, tag := range s {
					cmd.TagNames = append(cmd.TagNames, tag)
				}
			}
			if len(v) == 0 {
				err = fmt.Errorf("No tagname parameter[没有 tagname 参数]")
				return cmd, err
			}
		case "tablename":
			for i, s := range v {
				if i == 0 {
					cmd.TableName = s
				} else {
					err = fmt.Errorf("Too many tablename parameters[tablename 参数太多了]")
					return cmd, err
				}
			}
			if len(v) == 0 {
				err = fmt.Errorf("No begintime parameter[没有 begintime 参数]")
				return cmd, err
			}
		case "begintime":
			for i, s := range v {
				if i == 0 {
					if _, e := models.TimeParse(s); e != nil {
						return cmd, fmt.Errorf("begintime fomate error[begintime 格式错误];%s", e.Error())
					}
					cmd.BeginTime = s
				} else {
					err = fmt.Errorf("Too many begintime parameters[begintime 参数太多了]")
					return cmd, err
				}
			}
			if len(v) == 0 {
				err = fmt.Errorf("No begintime parameter[没有 begintime 参数]")
				return cmd, err
			}
		case "endtime":
			for i, s := range v {
				if i == 0 {
					if _, e := models.TimeParse(s); e != nil {
						return cmd, fmt.Errorf("endtime fomate error[endtime 格式错误];%s", e.Error())
					}
					cmd.EndTime = s
				} else {
					err = fmt.Errorf("Too many endtime parameters[endtime 参数太多了]")
					return cmd, err
				}
			}
			if len(v) == 0 {
				err = fmt.Errorf("No endtime parameter[没有 endtime 参数]")
				return cmd, err
			}
		case "timepoint":
			for i, s := range v {
				if i == 0 {
					if _, e := models.TimeParse(s); e != nil {
						return cmd, fmt.Errorf("timepoint fomate error[timepoint 格式错误];%s", e.Error())
					}
					cmd.TimePoint = s
				} else {
					err = fmt.Errorf("Too many timepoint parameters[timepoint 参数太多了]")
					return cmd, err
				}
			}
			if len(v) == 0 {
				err = fmt.Errorf("No timepoint parameter[没有 timepoint 参数]")
				return cmd, err
			}
		case "interval":
			for i, s := range v {
				if i == 0 {
					cmd.Interval, err = strconv.ParseInt(s, 10, 64)
					if err != nil {
						return cmd, fmt.Errorf("interval fomate error,must be an integer[interval 格式错误,必须为整数];%s", err.Error())
					}
				} else {
					err = fmt.Errorf("Too many interval parameters[interval 参数太多了]")
					return cmd, err
				}
			}
			if len(v) == 0 {
				err = fmt.Errorf("No interval parameter[没有 interval 参数]")
				return cmd, err
			}
		default:
			err = fmt.Errorf("There are unrecognized parameters[有未识别的参数]")
			return cmd, err
		}
	}
	return cmd, err
}

/*
Order=server,查询服务器时间
Order=table;获取全部表的id
Order=table&tableName=all:获取全部表的名称
Order=table&tableName=XXX:获取表属性

Order=point&tableName=XXX:获取指定表下的全部标签信息
Order=point&tagNames=XXX,YYY:获取指定标签点的属性信息

Order=snapshot&tagNames=XXX,YYY:获取指定标签点的快照信息

Order=history&tagNames=XXX,YYY&pointTime=yyyy-MM-dd hh:mm:ss&Interval=0;获取指定标签点在指定时间的历史数据
Order=history&tagNames=XXX,YYY&beginTime=yyyy-MM-dd hh:mm:ss&endTime=yyyy-MM-dd hh:mm:ss&Interval=0;获取指定标签点的原始历史数据
Order=history&tagNames=XXX,YYY&beginTime=yyyy-MM-dd hh:mm:ss&endTime=yyyy-MM-dd hh:mm:ss&Interval=n;获取指定标签点的等间隔历史数据(n>0)
Order=HistorySummary&tagNames=XXX,YYY&beginTime=yyyy-MM-dd hh:mm:ss&endTime=yyyy-MM-dd hh:mm:ss;获取指定标签点在指定时间范围内的全部统计数据
Order=MMM&tagNames=XXX,YYY&beginTime=yyyy-MM-dd hh:mm:ss&endTime=yyyy-MM-dd hh:mm:ss;获取指定标签点的指定统计数据
其中MMM可以是：min,max,range,total,sum,mean,poweravg,diff,duration,pointcnt,risingcnt,fallingcnt,sd, stddev, se, ske, kur,mode, median,distribution,increment,groupdist
Order=goods&tagNames=XXX,YYY&beginTime=yyyy-MM-dd hh:mm:ss&endTime=yyyy-MM-dd hh:mm:ss;获取指定物耗标签点在指定时间范围内的累积量
Order=goods&tagNames=XXX,YYY&pointTime=yyyy-MM-dd hh:mm:ss;获取指定物耗标签点在指定时间点的值(指定的没有值的则取之前的第一个值)
*/
func checkGetCmd(input getCmd) (interface{}, error) {
	micgd := new(models.MicGolden)
	switch input.Order {
	case "goods":
		return getGoodsDataFromDB(input.TagNames, input.BeginTime, input.EndTime, input.TimePoint)
	case "server": //读取服务器时间
		return micgd.GoldenGetServerTime()
	case "table":
		switch strings.ToLower(input.TableName) {
		case "all": //读取全部表列表指令
			return micgd.GoldenGetTables(2)
		case "": //读取全部表ID
			return micgd.GoldenGetTables(1)
		default: //输入为表名,读取表的详细信息列表
			tblist, err := micgd.GoldenGetTablePropertyByTableName(input.TableName)
			return tblist, err
		}
	case "point":
		if len(input.TagNames) > 0 {
			return micgd.GoldenGetTagPointInfoByName(input.TagNames...)
		} else if len(input.TableName) > 0 {
			return micgd.GoldenGetTagNameListInTables(input.TableName)
		} else {
			return nil, fmt.Errorf("没有输入合法的参数")
		}
	case "snapshot":
		if len(input.TagNames) == 0 {
			return nil, fmt.Errorf("Tagnames is a required parameter[Tagnames 是必须的参数]")
		}
		snaps, err := micgd.GoldenGetSnapShotMap(input.TagNames...)
		type snap struct {
			Time    string
			Value   float64
			Quality int
		}
		snatmap := make(map[string]snap)
		for tag, val := range snaps {
			var md snap
			md.Quality = val.Rtsd.Quality
			md.Time = models.TimeFormat(models.Millisecond2Time(val.Rtsd.Time))
			md.Value = val.Rtsd.Value
			snatmap[tag] = md
		}
		return snatmap, err
	case "history":
		type his struct {
			Time  string
			Value float64
		}
		if len(input.TagNames) == 0 {
			return nil, fmt.Errorf("Tagnames is a required parameter[Tagnames 是必须的参数]")
		}
		if len(input.TimePoint) == 0 {
			if len(input.BeginTime) == 0 {
				return nil, fmt.Errorf("BeginTime is a required parameter[BeginTime 是必须的参数]")
			}
			if len(input.EndTime) == 0 {
				return nil, fmt.Errorf("EndTime is a required parameter[EndTime 是必须的参数]")
			}
			hiss, _, err := micgd.GoldenGetHistoryInterval(input.BeginTime, input.EndTime, input.Interval, input.TagNames...)
			hismap := make(map[string][]his)
			for tag, val := range hiss {
				var mds []his
				for _, v := range val {
					var md his
					md.Time = v.Time
					md.Value = v.Value
					mds = append(mds, md)
				}
				hismap[tag] = mds
			}
			return hismap, err
		} else {
			hispoint, err := micgd.GoldenGetHistorySinglePoint(1, input.TimePoint, input.TagNames...)
			hismap := make(map[string]his)
			for tag, val := range hispoint {
				var md his
				md.Time = models.TimeFormat(models.Millisecond2Time(val.Time))
				md.Value = val.Value
				hismap[tag] = md
			}
			return hismap, err
		}
	case "historysummary", "advangce":
		if len(input.TagNames) == 0 {
			return nil, fmt.Errorf("Tagnames is a required parameter[Tagnames 是必须的参数]")
		}
		if len(input.BeginTime) == 0 {
			return nil, fmt.Errorf("BeginTime is a required parameter[BeginTime 是必须的参数]")
		}
		if len(input.EndTime) == 0 {
			return nil, fmt.Errorf("EndTime is a required parameter[EndTime 是必须的参数]")
		}
		return micgd.GoldenGetHistorySummary(input.BeginTime, input.EndTime, input.TagNames...)
	case "min", "max", "range", "total", "sum", "mean", "poweravg", "diff", "duration", "pointcnt", "risingcnt", "fallingcnt", "sd", "stddev", "se", "ske", "kur", "mode", "median", "distribution", "datagroup", "increment", "groupdist":
		if len(input.TagNames) == 0 {
			return nil, fmt.Errorf("Tagnames is a required parameter[Tagnames 是必须的参数]")
		}
		if len(input.BeginTime) == 0 {
			return nil, fmt.Errorf("BeginTime is a required parameter[BeginTime 是必须的参数]")
		}
		if len(input.EndTime) == 0 {
			return nil, fmt.Errorf("EndTime is a required parameter[EndTime 是必须的参数]")
		}
		return micgd.GoldenGetSingleStatisticsDataBatch(input.Order, input.BeginTime, input.EndTime, input.TagNames...)
	case "lt_l_total", "lte_l_total", "gt_l_total", "gte_l_total",
		"lt_ll_total", "lte_ll_total", "gt_ll_total", "gte_ll_total",
		"lt_h_total", "lte_h_total", "gt_h_total", "gte_h_total",
		"lt_hh_total", "lte_hh_total", "gt_hh_total", "gte_hh_total",
		"lt_min_total", "lte_min_total", "gt_min_total", "gte_min_total",
		"lt_max_total", "lte_max_total", "gt_max_total", "gte_max_total":
		return dataCompareTotal(input)
	case "pvs_peaksum", "pvs_valleysum", "pvs_pvdiffsum", "pvs_periodcnt", "pvs_peek_valley_datas", "peakvalley":
		return getPeakValley(input)
	default:
		return nil, fmt.Errorf("Invalid order[无效的指令]")
	}
	return nil, nil
}

/**********************************************
功能：获取峰谷值数据
输入：请求数据的结构
输出：获取到的数据及错误信息
说明：
时间：2020年2月5日
编辑：wang_jp
**********************************************/
func getPeakValley(input getCmd) (interface{}, error) {
	fmt.Println("GetCmd指令:", input)
	rslt := make(map[string]interface{})
	for _, tagname := range input.TagNames {
		micgd := new(models.MicGolden)
		res, _, err := micgd.GoldenGetPeakValleyOfHisData(0, tagname, input.Order, input.BeginTime, input.EndTime)
		if err != nil {
			return nil, err
		}
		rslt[tagname] = res
	}
	return rslt, nil
}

/***********************************************
功能：从平台数据库中获取物耗数据
输入：标签名数组，时间点,开始时间，结束时间
输出：获取到的数据及错误信息
说明：
时间：2019年12月22日
编辑：wang_jp
***********************************************/
func getGoodsDataFromDB(tagnames []string, beginTime, endTime, timePoint string) (interface{}, error) {
	if len(tagnames) == 0 {
		return nil, fmt.Errorf("Tagnames is a required parameter[Tagnames 是必须的参数]")
	}
	ids := make([]int64, len(tagnames))
	goodsv := make(map[string]goodsPointData, len(tagnames))
	goodssum := make(map[string]goodsSumData, len(tagnames))
	for i, tag := range tagnames {
		gd := new(models.GoodsConfigInfo)
		id, err := gd.GetGoodsTagIDByTagName(tag) //根据物耗变量名查询其id
		if err != nil {
			return nil, fmt.Errorf("Tagname [%s] does not exist[变量名不存在]", tag)
		}
		ids[i] = id
	}
	if len(timePoint) == 0 { //求指定时间范围内的和
		if len(beginTime) == 0 {
			return nil, fmt.Errorf("BeginTime is a required parameter[BeginTime 是必须的参数]")
		}
		if len(endTime) == 0 {
			return nil, fmt.Errorf("EndTime is a required parameter[EndTime 是必须的参数]")
		}

		var gvalue goodsSumData
		gvalue.BeginTime = beginTime
		gvalue.EndTime = endTime
		for i, id := range ids {
			gd := new(models.GoodsConsumeInfo)
			gv, err := gd.GetGoodsSumByID(id, beginTime, endTime) //查询物耗累积
			if err != nil {
				return nil, err
			}
			gvalue.Sum = gv
			goodssum[tagnames[i]] = gvalue
		}
		return goodssum, nil
	} else { //查询时间点的值
		var gvalue goodsPointData
		for i, id := range ids {
			gd := new(models.GoodsConsumeInfo)
			t, gv, err := gd.GetGoodsValueByIDandTime(id, timePoint) //查询指定时间点的物耗值
			if err != nil {
				return nil, err
			}
			gvalue.Value = gv
			gvalue.Time = t
			goodsv[tagnames[i]] = gvalue
		}
		return goodsv, nil
	}
}
