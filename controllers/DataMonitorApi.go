package controllers

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/bkzy-wangjp/MicEngine/MicScript/statistic"
	"github.com/bkzy-wangjp/MicEngine/models"
)

const (
	//日志操作类型:0:其他,1:"添加",2:"删除",3:"更新",4:"查看",5:添加/更新",6:"登录"
	_LOG_OPR_TYPE_OTHER          = 0
	_LOG_OPR_TYPE_INSERT         = 1
	_LOG_OPR_TYPE_DELETE         = 2
	_LOG_OPR_TYPE_UPDATE         = 3
	_LOG_OPR_TYPE_SELECT         = 4
	_LOG_OPR_TYPE_UPDATEORSELECT = 5
	_LOG_OPR_TYPE_LOGIN          = 6
)

/***********************************************
功能:快照的变量列表
输入:无
输出:无
说明:
编辑:wang_jp
时间:2020年3月18日
************************************************/
func (c *DataMonitorController) ApiGetTaglistByLevelCode() { //快照的变量列表
	//c.CheckSession() //权限验证
	//通过层级码读取taglists
	var lcode string
	lcode = c.GetString("levelcode", "0")
	if models.EngineCfgMsg.Sys.Debug {
		c.SaveUserActionMsg("获取层级下的变量列表", _LOG_OPR_TYPE_SELECT)
	}
	taglists, err := models.GetTaglistByNodeLevelCode(lcode)
	if err != nil {
		c.Data["json"] = err.Error()
	} else {
		c.Data["json"] = taglists
	}
	c.ServeJSON()
}

/***********************************************
功能:获取等间隔历史数据
输入:无
输出:无
说明:
编辑:wang_jp
时间:2020年3月19日
************************************************/
func (c *DataMonitorController) ApiGetHisInterData() { //获取等间隔历史数据
	//c.CheckSession() //权限验证
	//通过层级码读取taglists
	tagname := c.GetString("tagname")
	begintime := c.GetString("begintime")
	endtime := c.GetString("endtime")
	interval, err := c.GetInt64("interval")
	if err != nil {
		interval = 60
	}
	micgd := new(models.MicGolden)
	hiss, _, err := micgd.GoldenGetHistoryInterval(begintime, endtime, interval, tagname)
	if err != nil {
		c.Data["json"] = err.Error()
	} else {
		c.Data["json"] = hiss
	}
	c.ServeJSON()
}

/***********************************************
功能:获取历史数据
输入:无
输出:无
说明:
编辑:wang_jp
时间:2020年3月19日
************************************************/
func (c *DataMonitorController) ApiGetHistoryData() { //获取历史数据
	//c.CheckSession() //权限验证
	//通过层级码读取taglists
	tagname := c.GetString("tagname")
	begintime := c.GetString("begintime")
	endtime := c.GetString("endtime")
	micgd := new(models.MicGolden)
	hiss, err := micgd.GoldenGetHistory(begintime, endtime, tagname)
	if err != nil {
		c.Data["json"] = err.Error()
	} else {
		c.Data["json"] = hiss
	}
	c.ServeJSON()
}

/***********************************************
功能:获取历史统计数据
输入:无
输出:无
说明:
编辑:wang_jp
时间:2020年3月21日
************************************************/
func (c *DataMonitorController) ApiGetHistorySummaryData() { //获取历史统计数据
	//c.CheckSession() //权限验证
	//通过层级码读取taglists
	tagname := c.GetString("tagname")
	begintime := c.GetString("begintime")
	endtime := c.GetString("endtime")
	tags := strings.Split(tagname, ",")
	micgd := new(models.MicGolden)
	hiss, err := micgd.GoldenGetHistorySummary(begintime, endtime, tags...)
	if err != nil {
		c.Data["json"] = err.Error()
	} else {
		c.Data["json"] = hiss
	}
	c.ServeJSON()
}

/***********************************************
功能:获取快照数据
输入:无
输出:无
说明:
编辑:wang_jp
时间:2020年3月19日
************************************************/
func (c *DataMonitorController) ApiGetSnapshotData() { //获取快照数据
	//c.CheckSession() //权限验证
	//通过层级码读取taglists
	tagnames := c.GetString("tagnames")
	tags := strings.Split(tagnames, ",")
	micgd := new(models.MicGolden)
	snap, err := micgd.GoldenGetSnapShotMap(tags...)
	if err != nil {
		c.Data["json"] = err.Error()
	} else {
		c.Data["json"] = snap
	}
	c.ServeJSON()
}

/***********************************************
功能:写数据进入快照
输入:无
输出:无
说明:
编辑:wang_jp
时间:2020年3月19日
************************************************/
func (c *DataMonitorController) ApiWriteSnapshot() {
	//c.CheckSession()                                 //权限验证
	c.SaveUserActionMsg("写快照", _LOG_OPR_TYPE_INSERT) //记录用动作作信息
	type SnapWrite struct {                          //写庚顿快照数据结构
		TagName string  `json:"tagname"` //测点全名称
		Value   float64 `json:"value"`   //数值
		Time    string  `json:"time"`
		Quality int     `json:"quality"` //质量码(GOOD = 0,NODATA = 1,CREATED = 2,SHUTDOWN = 3,CALCOFF = 4,BAD = 5,DIVBYZERO = 6,REMOVED = 7,OPC = 256,USER = 512)
	}
	var writers []SnapWrite
	datas := c.GetString("datas")
	err := json.Unmarshal([]byte(datas), &writers)
	if err != nil {
		c.Data["json"] = err.Error()
	} else {
		var tags []string
		var values []float64
		var dtimes []string
		var dqs []int
		for _, snap := range writers {
			tags = append(tags, snap.TagName)
			values = append(values, snap.Value)
			dtimes = append(dtimes, snap.Time)
			dqs = append(dqs, snap.Quality)
		}
		micgd := new(models.MicGolden)
		err := micgd.GoldenSetSnapShots(tags, values, dqs, dtimes)
		if err != nil {
			c.Data["json"] = err.Error()
		} else {
			c.Data["json"] = "ok"
		}
	}
	c.ServeJSON()
}

/***********************************************
功能:获取变量的属性信息
输入:无
输出:无
说明:
编辑:wang_jp
时间:2020年3月19日
************************************************/
func (c *DataMonitorController) ApiGetTagAttribut() { //获取变量的属性信息
	//c.CheckSession() //权限验证
	//通过tagid读取tag属性
	tagid, err := c.GetInt64("tagid")
	if err != nil {
		c.Data["json"] = err.Error()
	} else {
		tag := new(models.OreProcessDTaglist)
		err := tag.GetTagAttributByTagId(tagid)
		if err != nil {
			c.Data["json"] = err.Error()
		} else {
			c.Data["json"] = tag
		}
	}
	c.ServeJSON()
}

/***********************************************
功能:最小二乘回归分析
输入:无
输出:无
说明:
编辑:wang_jp
时间:2020年3月19日
************************************************/
func (c *DataMonitorController) ApiRegression() { //最小二乘回归分析
	//c.CheckSession() //权限验证
	c.SaveUserActionMsg("回归分析", _LOG_OPR_TYPE_OTHER)
	//通过tagid读取tag属性
	tagy := c.GetString("tagy")
	tagxs := c.GetString("tagxs")
	begintime := c.GetString("begintime")
	endtime := c.GetString("endtime")
	interval, err := c.GetInt64("interval")

	if err != nil {
		c.Data["json"] = err.Error()
	} else {
		micgd := new(models.MicGolden)
		reg, err := micgd.Regression(tagy, tagxs, begintime, endtime, interval)
		if err != nil {
			c.Data["json"] = err.Error()
		} else {
			c.Data["json"] = reg
		}
	}
	c.ServeJSON()
}

/***********************************************
功能:批量自动追加kpi指标
输入:beginid:起始taglist的id值
	period:周期,,-1:小时;-2:班;-3:日;-4:月;-5:季度;-6:年；如果大于0则单位为秒。
	begintime:kpi指标起始计算时间
输出:追加的kpi指标数量
说明:
编辑:wang_jp
时间:2020年3月19日
************************************************/
func (c *DataMonitorController) ApiAppendTagKpi2CfgList() { //批量自动追加kpi指标
	//c.CheckSession() //权限验证
	//通过tagid读取tag属性
	c.SaveUserActionMsg("追加KPI指标", _LOG_OPR_TYPE_INSERT)
	beginid, err1 := c.GetInt("beginid")
	begintime := c.GetString("begintime")
	period, err2 := c.GetInt("period")
	if err1 != nil {
		c.Data["json"] = err1.Error()
	} else if err2 != nil {
		c.Data["json"] = err2.Error()
	} else {
		lst := new(models.CalcKpiConfigList)
		num, err := lst.InsertTagKpi2CfgListByAppend(beginid, period, begintime)
		if err != nil {
			c.Data["json"] = err.Error()
		} else {
			c.Data["json"] = fmt.Sprintf("共成功追加了%d条KPI配置指标", num)
		}
	}
	c.ServeJSON()
}

/***********************************************
功能:获取样本的化验变量标签信息
输入:样本ID
输出:样本的化验标签信息
说明:
编辑:wang_jp
时间:2020年4月9日
************************************************/
func (c *DataMonitorController) ApiSampleLabTag() { //获取样本的化验变量标签信息
	//c.CheckSession() //权限验证
	//通过tagid读取tag属性
	sampleid, err := c.GetInt64("sampleid")
	if err != nil {
		c.Data["json"] = err.Error()
	} else {
		sub := new(models.SamplingManageSub)
		tag, err := sub.GetSampleSubBySampleId(sampleid)
		if err != nil {
			c.Data["json"] = err.Error()
		} else {
			c.Data["json"] = tag
		}
	}
	c.ServeJSON()
}

/***********************************************
功能:获取样本的化验结果信息
输入:样本ID,开始时间,结束时间
输出:样本的化验标签信息
说明:
编辑:wang_jp
时间:2020年4月9日
************************************************/
func (c *DataMonitorController) ApiSampleLabResult() { //获取样本的化验结果信息
	//c.CheckSession() //权限验证
	//通过tagid读取tag属性
	sampleid, err := c.GetInt64("sampleid")
	if err != nil {
		c.Data["json"] = err.Error()
	} else {
		bgtime := c.GetString("begintime")
		endtime := c.GetString("endtime")

		lab := new(models.LabAnaResultTsd)
		tag, err := lab.GetSampleLabResultBySampleId(sampleid, bgtime, endtime)
		if err != nil {
			c.Data["json"] = err.Error()
		} else {
			c.Data["json"] = tag
		}
	}
	c.ServeJSON()
}

/***********************************************
功能:获取物耗标签信息
输入:样本ID
输出:样本的化验标签信息
说明:
编辑:wang_jp
时间:2020年4月9日
************************************************/
func (c *DataMonitorController) ApiGoodsCfg() { //获取样本的化验变量标签信息
	//c.CheckSession() //权限验证
	//通过tagid读取tag属性
	goodsid, err := c.GetInt64("goodsid")
	if err != nil {
		c.Data["json"] = err.Error()
	} else {
		gd := new(models.GoodsConfigInfo)
		tag, err := gd.GetGoodsInfoByGoodsId(goodsid)
		if err != nil {
			c.Data["json"] = err.Error()
		} else {
			c.Data["json"] = tag
		}
	}
	c.ServeJSON()
}

/***********************************************
功能:获取物耗录入结果信息
输入:样本ID,开始时间,结束时间
输出:物耗录入结果信息
说明:
编辑:wang_jp
时间:2020年4月9日
************************************************/
func (c *DataMonitorController) ApiGoodsDatas() {
	//c.CheckSession() //权限验证
	bgtime := c.GetString("begintime")
	endtime := c.GetString("endtime")
	isid := c.GetString("isid", "0")
	startonly, _ := c.GetBool("startonly", true)
	if isid == "0" {
		level := c.GetString("leveltag")
		gd := new(models.GoodsConsumeInfo)
		tag, err := gd.GetGoodsResultsByLevelCode(level, bgtime, endtime, startonly)
		if err != nil {
			c.Data["json"] = err.Error()
		} else {
			c.Data["json"] = tag
		}
	} else {
		goodsid, err := c.GetInt64("leveltag")
		if err != nil {
			c.Data["json"] = err.Error()
		} else {
			gd := new(models.GoodsConsumeInfo)
			tag, err := gd.GetGoodsResultsByGoodsId(goodsid, bgtime, endtime, startonly)
			if err != nil {
				c.Data["json"] = err.Error()
			} else {
				c.Data["json"] = tag
			}
		}
	}
	c.ServeJSON()
}

/***********************************************
功能:获取巡检标签信息
输入:样本ID
输出:巡检标签信息
说明:
编辑:wang_jp
时间:2020年4月9日
************************************************/
func (c *DataMonitorController) ApiPatrolTagList() {
	//c.CheckSession() //权限验证
	isid := c.GetString("isid", "1")
	id, err := c.GetInt64("leveltag")
	if err != nil {
		isid = "0"
	}
	if isid == "0" { //不是ID
		level := c.GetString("leveltag")
		cklst := new(models.CheckTagList)
		tag, err := cklst.GetCheckTagListByLevelCode(level)
		if err != nil {
			c.Data["json"] = err.Error()
		} else {
			c.Data["json"] = tag
		}
	} else { //是ID
		cklst := new(models.CheckTagList)
		tag, err := cklst.GetCheckTagListsBySiteId(id)
		if err != nil {
			c.Data["json"] = err.Error()
		} else {
			c.Data["json"] = tag
		}
	}
	c.ServeJSON()
}

/***********************************************
功能:获取巡检结果结果信息
输入:样本ID,开始时间,结束时间
输出:巡检结果信息
说明:
编辑:wang_jp
时间:2020年4月9日
************************************************/
func (c *DataMonitorController) ApiPatrolResult() {
	//c.CheckSession() //权限验证
	bgtime := c.GetString("begintime")
	endtime := c.GetString("endtime")
	isid := c.GetString("isid", "0")
	startonly, _ := c.GetBool("startonly", true)
	if isid == "0" {
		level := c.GetString("leveltag")
		ckexe := new(models.CheckItemExe)
		tag, err := ckexe.GetCheckResultsByLevelCode(level, bgtime, endtime, startonly)
		if err != nil {
			c.Data["json"] = err.Error()
		} else {
			c.Data["json"] = tag
		}
	} else {
		goodsid, err := c.GetInt64("leveltag")
		if err != nil {
			c.Data["json"] = err.Error()
		} else {
			ckexe := new(models.CheckItemExe)
			tag, err := ckexe.GetCheckResultsByTagId(goodsid, bgtime, endtime, startonly)
			if err != nil {
				c.Data["json"] = err.Error()
			} else {
				c.Data["json"] = tag
			}
		}
	}
	c.ServeJSON()
}

/***********************************************
功能:通过CheckSiteExeID获取巡检点的照片
输入:CheckSiteExeID
输出:样本的化验标签信息
说明:
编辑:wang_jp
时间:2020年4月9日
************************************************/
func (c *DataMonitorController) ApiGetPatrolPicUrl() { //获取样本的化验变量标签信息
	//c.CheckSession() //权限验证
	siteexeid, err := c.GetInt64("siteexeid")
	if err != nil {
		c.Data["json"] = err.Error()
	} else {
		img := new(models.MineCheckImg)
		tag, err := img.GetMineCheckImg(48, siteexeid)
		if err != nil {
			c.Data["json"] = err.Error()
		} else {
			tag.FileName = models.EngineCfgMsg.CfgMsg.PlatPath + "common/imgShow/" + tag.FileName
			c.Data["json"] = tag
		}
	}
	c.ServeJSON()
}

/***********************************************
功能:通过层级码获取该层级下的巡检点
输入:层级码
输出:巡检点信息
说明:
编辑:wang_jp
时间:2020年4月18日
************************************************/
func (c *DataMonitorController) ApiGetPatrolSiteListByLevelCode() {
	level := c.GetString("levelcode")
	bgtime := c.GetString("begintime")
	endtime := c.GetString("endtime")
	sit := new(models.CheckSite)
	tag, err := sit.GetCheckSiteByLevelCode(level, bgtime, endtime)
	if err != nil {
		c.Data["json"] = err.Error()
	} else {
		c.Data["json"] = tag
	}
	c.ServeJSON()
}

/***********************************************
功能:获取用户日志信息
输入:begintime:开始时间,格式:2006-01-02 15:04:05
	endtime:结束时间,格式:2006-01-02 15:04:05
	userid:用户ID,小于1时取所有用户ID
	systype:系统类型,1：PC端 2：移动端 3:计算服务,小于1时取前述所有
	oprtype:操作类型,0:其他,1:"添加",2:"删除",3:"更新",4:"查看",5:添加/更新",6:"登录"
			小于0时取上述所有。
	pagesize:每页显示的信息数,小于等于0时不限制
	pageno:偏移页数，最小为0
	desc:描述,可选。
输出:巡检点信息
说明:
编辑:wang_jp
时间:2020年4月18日
************************************************/
func (c *DataMonitorController) ApiGetUserLogs() {
	type logReqMsg struct { //请求参数结构
		Begintime string `form:"begintime"`
		Endtime   string `form:"endtime"`
		Userid    int64  `form:"userid"`
		Systype   int64  `form:"systype"`
		Oprtype   int64  `form:"oprtype"`
		Pagesize  int64  `form:"pagesize"`
		Pageno    int64  `form:"pageno"`
		Desc      string `form:"desc"`
	}
	type logMsg struct { //结果数据结构
		TotalRows int64            //总行数
		PageSize  int64            //页面显示行数
		PageNo    int64            //页面号码
		MaxPage   int64            //最大页面数
		FirstPage bool             //是否第一页
		LastPage  bool             //是否最后一页
		UserNames map[int64]string //用户名列表
		Logs      []*models.SysLog //日志信息
	}
	//c.SaveUserActionMsg("读取系统日志", _LOG_OPR_TYPE_SELECT)
	req := new(logReqMsg)   //请求信息
	err := c.ParseForm(req) //解析请求信息
	rst := new(logMsg)      //结果结构

	if err != nil { //解析请求错误
		c.Data["json"] = err.Error()
	} else { //
		//读取日志
		if req.Pageno < 0 {
			req.Pageno = 0
		}
		log := new(models.SysLog)
		rst.TotalRows, rst.MaxPage, rst.Logs, err = log.GetSysLog(req.Begintime, req.Endtime, req.Userid, req.Systype, req.Oprtype, req.Pagesize, req.Pageno, req.Desc)
		if req.Pageno >= rst.MaxPage { //当前页不能大于最大页
			rst.PageNo = rst.MaxPage - 1
		} else {
			rst.PageNo = req.Pageno
		}
		if rst.PageNo < 0 {
			rst.PageNo = 0
		}
		rst.FirstPage = rst.PageNo == 0          //是否首页
		rst.LastPage = rst.MaxPage == rst.PageNo //是否尾页
		names := make(map[int64]string)
		for _, log := range rst.Logs {
			names[log.User.Id] = log.User.Name
		}
		rst.UserNames = names
		rst.PageSize = req.Pagesize
		c.Data["json"] = rst
	}
	c.ServeJSON()
}

/***********************************************
功能:获取用户日志信息
输入:begintime:开始时间,格式:2006-01-02 15:04:05
	endtime:结束时间,格式:2006-01-02 15:04:05
	userid:用户ID,小于1时取所有用户ID
	systype:系统类型,1：PC端 2：移动端 3:计算服务,小于1时取前述所有
	oprtype:操作类型,0:其他,1:"添加",2:"删除",3:"更新",4:"查看",5:添加/更新",6:"登录"
			小于0时取上述所有。
	pagesize:每页显示的信息数,小于等于0时不限制
	pageno:偏移页数，最小为0
	desc:描述,可选。
输出:巡检点信息
说明:
编辑:wang_jp
时间:2020年4月18日
************************************************/
func (c *DataMonitorController) ApiGetUserLogsAnalys() {
	type logReqMsg struct { //请求参数结构
		Begintime string `form:"begintime"`
		Endtime   string `form:"endtime"`
		Userid    int64  `form:"userid"`
		Systype   int64  `form:"systype"`
		Oprtype   int64  `form:"oprtype"`
		Desc      string `form:"desc"`
	}

	c.SaveUserActionMsg("系统日志分析", _LOG_OPR_TYPE_SELECT)
	req := new(logReqMsg)   //请求信息
	err := c.ParseForm(req) //解析请求信息
	if err != nil {         //解析请求错误
		c.Data["json"] = err.Error()
	} else { //
		//读取日志统计结果
		log := new(models.SysLog)
		logs, err := log.GetSysLogAnalyse(req.Begintime, req.Endtime, req.Userid, req.Systype, req.Oprtype, req.Desc)
		if err != nil { //解析请求错误
			c.Data["json"] = err.Error()
		} else {
			c.Data["json"] = logs
		}
	}
	c.ServeJSON()
}

func (c *DataMonitorController) ApiPortListen() {
	c.SaveUserActionMsg("端口监听", _LOG_OPR_TYPE_SELECT)
	c.Data["json"] = c.Ctx.Request.RequestURI
	c.ServeJSON()
}

/***********************************************
功能:更新用户节点树
输入:withtag:是否带节点,是为1,否为0
输出:错误信息或者正确结果
说明:
编辑:wang_jp
时间:2020年5月20日
************************************************/
func (c *DataMonitorController) ApiUpdateNodeTree() {
	withtag, err := c.GetBool("withtag")
	if err != nil {
		c.Data["json"] = err.Error()
	}
	c.FormatTreeNode("", "", withtag, true)
	c.Data["json"] = "更新完成,请刷新页面"
	c.ServeJSON()
}

/***********************************************
功能:获取KPI计算结果
输入:kpicfgid,begintime,endtime,readtype
输出:错误信息或者正确结果
说明:
编辑:wang_jp
时间:2020年7月16日
************************************************/
func (c *DataMonitorController) ApiGetKpiConfig() {
	kpicfgid, err := c.GetInt64("kpicfgid")
	if err != nil {
		c.Ctx.ResponseWriter.WriteHeader(_PAROMATER_ERR)
		c.Data["json"] = "没有设定有效的kpicfgid"
	} else {
		cfg := new(models.CalcKpiConfigList)
		cfg.Id = kpicfgid
		err := cfg.GetKpiConfig()
		if err != nil {
			c.Ctx.ResponseWriter.WriteHeader(_SERVER_ERR)
			c.Data["json"] = err.Error()
		} else {
			c.Data["json"] = cfg
		}
	}
	c.ServeJSON()
}

/***********************************************
功能:获取KPI计算结果
输入:kpicfgid,begintime,endtime,readtype
输出:错误信息或者正确结果
说明:
编辑:wang_jp
时间:2020年7月16日
************************************************/
func (c *DataMonitorController) ApiGetKpiResult() {
	kpicfgid, err := c.GetInt64("kpicfgid")
	if err != nil {
		c.Ctx.ResponseWriter.WriteHeader(_PAROMATER_ERR)
		c.Data["json"] = "没有设定有效的kpicfgid"
	} else {
		readtype, err := c.GetInt("readtype")
		if err != nil {
			c.Ctx.ResponseWriter.WriteHeader(_PAROMATER_ERR)
			c.Data["json"] = "没有设定有效的readtype"
		} else {
			begintime := c.GetString("begintime")
			endtime := c.GetString("endtime")
			cfg := new(models.CalcKpiConfigList)
			cfg.Id = kpicfgid
			kpis, err := cfg.GetKpiDatas(begintime, endtime, readtype)
			if err != nil {
				c.Ctx.ResponseWriter.WriteHeader(_SERVER_ERR)
				c.Data["json"] = err.Error()
			} else {
				c.Data["json"] = kpis
			}
		}
	}
	c.ServeJSON()
}

/***********************************************
功能:系统时序数据统计
输入:tagid,tagname,tagtype,key,beginTime,endTime
	needraw:是否需要返回原始数据,0或者省略不需要,大于0时需要
	filloutliers:填充异常值的方法,可选 midean(中位数) 或者 extremum(四分位极值),不选不处理
输出:
说明:如果tagid为0,则通过tagname查询数据
编辑:wang_jp
时间:2020年8月26日
************************************************/
func (c *DataMonitorController) ApiSrtdStatistic() {
	tagid, _ := c.GetInt64("tagid", 0)
	needraw, _ := c.GetInt64("needraw", 0)
	tagname := c.GetString("tagname")
	tagtype := c.GetString("tagtype")
	beginTime := c.GetString("begintime")
	endTime := c.GetString("endtime")
	filloutliers := c.GetString("filloutliers")

	key := "advangce"
	if needraw > 0 {
		key = "needraw"
	}

	_, errb := models.TimeParse(beginTime)
	_, erre := models.TimeParse(endTime)
	if errb != nil || erre != nil {
		c.Ctx.ResponseWriter.WriteHeader(_PAROMATER_ERR)
		if errb != nil {
			c.Data["json"] = fmt.Sprintf("begintime format error[开始时间格式错误]:%s", errb.Error())
		}
		if erre != nil {
			c.Data["json"] = fmt.Sprintf("endtime format error[结束时间格式错误]:%s", erre.Error())
		}
	} else {
		srd := new(models.SysRealData)
		stat, _, err := srd.GetSysRealTimeDataStatistic(tagid, tagname, tagtype, key, beginTime, endTime, filloutliers)
		if err != nil {
			c.Ctx.ResponseWriter.WriteHeader(_SERVER_ERR)
			c.Data["json"] = err.Error()
		} else {
			c.Data["json"] = stat
		}
	}
	c.ServeJSON()
}

/***********************************************
功能:系统时序数据增量统计
输入:tagid,tagname,tagtype,key,beginTime,endTime
	needraw:是否需要返回原始数据,0或者省略不需要,大于0时需要
	filloutliers:填充异常值的方法,可选 midean(中位数) 或者 extremum(四分位极值),不选不处理
输出:
说明:如果tagid为0,则通过tagname查询数据
编辑:wang_jp
时间:2020年8月26日
************************************************/
func (c *DataMonitorController) ApiSrtdIncrStatistic() {
	tagid, _ := c.GetInt64("tagid", 0)
	needraw, _ := c.GetInt("needraw", 0)
	tagname := c.GetString("tagname")
	tagtype := c.GetString("tagtype")
	beginTime := c.GetString("begintime")
	endTime := c.GetString("endtime")
	filloutliers := c.GetString("filloutliers")

	key := "advangce"
	if needraw > 0 {
		key = "needraw"
	}
	_, errb := models.TimeParse(beginTime)
	_, erre := models.TimeParse(endTime)
	if errb != nil || erre != nil {
		c.Ctx.ResponseWriter.WriteHeader(_PAROMATER_ERR)
		if errb != nil {
			c.Data["json"] = fmt.Sprintf("begintime format error[开始时间格式错误]:%s", errb.Error())
		}
		if erre != nil {
			c.Data["json"] = fmt.Sprintf("endtime format error[结束时间格式错误]:%s", erre.Error())
		}
	} else {
		srd := new(models.SysRealData)
		stat_r, _, err := srd.GetSysRealTimeDataStatistic(tagid, tagname, tagtype, key, beginTime, endTime)
		if err != nil {
			c.Ctx.ResponseWriter.WriteHeader(_SERVER_ERR)
			c.Data["json"] = err.Error()
		} else {
			var tsds statistic.Tsds              //时间序列数组
			for t, v := range stat_r.Increment { //整理出标准的时间序列数据
				var tsd statistic.TimeSeriesData
				tsd.Time, _ = models.TimeParse(t)
				tsd.Value = v
				tsds = append(tsds, tsd)
			}
			tsds.SortByTime()
			switch strings.ToLower(filloutliers) {
			case "midean": //用中位数填充异常值
				tsds.FillOutliersByMidean()
			case "extremum": //用四分位极值填充异常值
				tsds.FillOutliersByExtremum()
			default:
			}
			c.Data["json"] = tsds.Statistics(1 + needraw) //包含高级历史统计
		}
	}
	c.ServeJSON()
}

/***********************************************
功能:系统时序数据增量统计
输入:tagid,tagname,tagtype,key,beginTime,endTime
	needraw:是否需要返回原始数据,0或者省略不需要,大于0时需要
	filloutliers:填充异常值的方法,可选 midean(中位数) 或者 extremum(四分位极值),不选不处理
输出:
说明:如果tagid为0,则通过tagname查询数据
编辑:wang_jp
时间:2020年8月26日
************************************************/
func (c *DataMonitorController) ApiSrtdStatisticAuto() {
	tagid, _ := c.GetInt64("tagid", 0)
	needraw, _ := c.GetInt("needraw", 0)
	tagname := c.GetString("tagname")
	tagtype := c.GetString("tagtype")
	beginTime := c.GetString("begintime")
	endTime := c.GetString("endtime")
	filloutliers := c.GetString("filloutliers")

	_, errb := models.TimeParse(beginTime)
	_, erre := models.TimeParse(endTime)
	if errb != nil || erre != nil {
		c.Ctx.ResponseWriter.WriteHeader(_PAROMATER_ERR)
		if errb != nil {
			c.Data["json"] = fmt.Sprintf("begintime format error[开始时间格式错误]:%s", errb.Error())
		}
		if erre != nil {
			c.Data["json"] = fmt.Sprintf("endtime format error[结束时间格式错误]:%s", erre.Error())
		}
	} else {
		srd := new(models.SysRealData)
		srtds, err := srd.GetSysRealTimeDatas(tagid, tagname, tagtype, beginTime, endTime, 0) //读取指定时间范围内的数据
		if err != nil {
			c.Ctx.ResponseWriter.WriteHeader(_SERVER_ERR)
			c.Data["json"] = err.Error()
		} else {
			var result []statistic.StatisticData
			var tsds statistic.Tsds   //时间序列数组
			for _, v := range srtds { //整理出标准的时间序列数据
				var tsd statistic.TimeSeriesData
				tsd.Time, _ = models.TimeParse(v.Datatime)
				tsd.Value = v.Value
				tsds = append(tsds, tsd)
			}
			tsds.SortByTime()
			stt_raw := tsds.Statistics(1 + needraw) //原始数据的统计
			result = append(result, stt_raw)

			tsds_incr, _ := statistic.ParseFromTimeValueMap(stt_raw.Increment)
			tsds_incr.SortByTime()
			stt_incr := tsds_incr.Statistics(1 + needraw) //增量统计
			result = append(result, stt_incr)

			if stt_raw.OutliersCnt > 0 || stt_incr.OutliersCnt > 0 { //有离群点数据
				switch strings.ToLower(filloutliers) {
				case "midean": //用中位数填充异常值
					tsds.FillOutliersByMidean()
					tsds_incr.FillOutliersByMidean()
				case "extremum": //用四分位极值填充异常值
					tsds.FillOutliersByExtremum()
					tsds_incr.FillOutliersByMidean()
				default:
				}
				result = append(result, tsds.Statistics(1+needraw))
				result = append(result, tsds_incr.Statistics(1+needraw))
			}
			c.Data["json"] = result
		}
	}
	c.ServeJSON()
}

func (c *DataMonitorController) ApiOreProcessJson() {
	tag := new(models.OreProcessDTaglist)
	data := tag.OreProcessTag()

	c.Data["json"] = data
	c.ServeJSON()
}
