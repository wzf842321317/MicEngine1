package controllers

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/bkzy-wangjp/MicEngine/MicScript/extable"
	"github.com/bkzy-wangjp/MicEngine/models"
)

type ReportController struct {
	MyController //beego.Controller
}



/***********************************************
功能:获取表单数据
输入:
输出:
说明:
编辑:wangzf
时间:2020年10月13日
************************************************/
func (c *ReportController) SaveCalculate() {
	type KpiArithmeticResult struct {
		BeginTime string `form:"BeginTime"`
		EndTime   string `form:"EndTime"`
		TagId     string `form:"tag_id"`
		Fc        string `form:"fc"`
		Data      string `form:"ty"`
		State     string `form:"state"`
	}

	u := KpiArithmeticResult{}
	err := c.ParseForm(&u) //解析请求信息/接受前台变量
	//初始化数据库对象

	if err != nil {
		c.Data["json"] = err.Error()
	}

	c.Ctx.WriteString("")
}

/***********************************************
功能:获取用户ID被授权的报表层级列表
输入:[levelcode]
输出:报表节点集
说明:
编辑:wang_jp
时间:2020年5月9日
************************************************/
func (c *ReportController) ApiGetReportListsByUserId() {
	//userid := c.CheckSession() //检查授权,返回授权的ID
	userid, err := c.GetInt64("userid")
	if err != nil {
		c.Data["json"] = fmt.Errorf("用户ID必须为数字:[%s]", err.Error())
	} else {
		user := new(models.SysUser)
		user.Id = userid
		nodes, err := user.GetReportNodesByUserId()
		if err != nil {
			c.Data["json"] = err.Error()
		} else {
			c.Data["json"] = nodes
		}
	}
	c.ServeJSON()
}

/***********************************************
功能:获取所选报表层级下的所有子节点
输入:[levelcode]
输出:报表节点集
说明:
编辑:wang_jp
时间:2020年5月9日
************************************************/
func (c *ReportController) ApiGetReportChildNodes() {
	userid := c.CheckSession() //检查授权,返回授权的ID
	level := c.GetString("levelcode")
	rpt := new(models.CalcKpiReportList)
	rpt.LevelCode = level
	_, nodes, err := rpt.GetChildNodes(userid)
	if err != nil {
		c.Data["json"] = err.Error()
	} else {
		c.Data["json"] = nodes
	}
	c.ServeJSON()
}

/***********************************************
功能:获取所选报表的模板列表
输入:[id]
输出:报表节点集
说明:
编辑:wang_jp
时间:2020年5月9日
************************************************/
func (c *ReportController) ApiGetReportTplList() {
	id, err := c.GetInt64("id")
	if err != nil {
		c.Data["json"] = fmt.Sprintf("Id必须是整数,传入的是:[%s]", c.GetString("id"))
	} else {
		rpt := new(models.CalcKpiReportList) //新建报表结构
		rpt.Id = id
		tpllist, err := rpt.GetTemplatesListById()
		if err != nil {
			c.Data["json"] = err.Error()
		} else {
			c.Data["json"] = tpllist
		}
	}
	c.ServeJSON()
}

/***********************************************
功能:获取所选报表的结果列表
输入:[id]
输出:报表节点集
说明:
编辑:wang_jp
时间:2020年5月9日
************************************************/
func (c *ReportController) ApiGetReportResultList() {
	id, err := c.GetInt64("id")
	bgtime := c.GetString("begintime")
	endtime := c.GetString("endtime")
	if err != nil {
		c.Data["json"] = fmt.Sprintf("Id必须是整数,传入的是:[%s]", c.GetString("id"))
	} else {
		rpt := new(models.CalcKpiReportList) //新建报表结构
		rpt.Id = id
		tpllist, err := rpt.GetResultListById(bgtime, endtime)
		if err != nil {
			c.Data["json"] = err.Error()
		} else {
			c.Data["json"] = tpllist
		}
	}
	c.ServeJSON()
}

/***********************************************
功能:添加报表层级
输入:
输出:报表节点集
说明:
编辑:wang_jp
时间:2020年5月9日
************************************************/
func (c *ReportController) ApiAddReportLevel() {
	userid := c.CheckSession() //检查授权,返回授权的ID
	type nodemsg struct {      //请求参数结构
		DistributedId int64  `form:"DistributedId"` //分布式计算ID，为0的时候不区分计算引擎，不为0的时候需要id适配的计算引擎进行计算
		Name          string `form:"Name"`          //名称
		Pid           int64  `form:"Pid"`           //父级菜单ID
		WorkshopId    int    `form:"WorkshopId"`    //所属车间ID
		Level         int64  `form:"Level"`         //层级深度
		LevelCode     string `form:"LevelCode"`     //层级码
		Folder        int64  `form:"Folder"`        //是否是文件夹，1-是，0-否
		Debug         int64  `form:"Debug"`         //是否是调试模式，1-是，0-否
		Seq           int64  `form:"Seq"`           //排序号
		Remark        string `form:"Remark"`        //备注
		TemplateUrl   string `form:"TemplateUrl"`   //模板文件路径
		TemplateFile  string `form:"TemplateFile"`  //模板文件名称
		ResultUrl     string `form:"ResultUrl"`     //结果地址
		StartTime     string `form:"StartTime"`     //统计计算开始起作用的时间
		Period        int64  `form:"Period"`        //计算周期,详见KPI表
		OffsetMinutes int64  `form:"OffsetMinutes"` //偏移时间
		LastCalcTime  string `form:"LastCalcTime"`  //最后计算时间
		BaseTime      string `form:"BaseTime"`      //基准时间
		ShiftHour     int64  `form:"ShiftHour"`     //每班工作时间
		Status        int64  `form:"Status"`        //1有效 0无效
	}

	res := new(nodemsg)     //请求信息
	err := c.ParseForm(res) //解析请求信息
	if err != nil {         //解析请求错误
		c.Data["json"] = err.Error()
	} else { //
		w := new(models.OreProcessDWorkshop)
		w.Id = res.WorkshopId
		//新建节点
		rpt := new(models.CalcKpiReportList)
		rpt.DistributedId = res.DistributedId //分布式计算ID，为0的时候不区分计算引擎，不为0的时候需要id适配的计算引擎进行计算
		rpt.Name = res.Name                   //名称
		rpt.Pid = res.Pid                     //父级菜单ID
		rpt.Workshop = w
		rpt.Folder = res.Folder
		rpt.Debug = res.Debug
		rpt.Seq = res.Seq
		rpt.Remark = res.Remark
		rpt.TemplateUrl = res.TemplateUrl     //模板文件路径
		rpt.ResultUrl = res.ResultUrl         //结果地址
		rpt.StartTime = res.StartTime         //统计计算开始起作用的时间
		rpt.Period = res.Period               //计算周期,详见KPI表
		rpt.OffsetMinutes = res.OffsetMinutes //偏移时间
		rpt.LastCalcTime = res.LastCalcTime   //最后计算时间
		rpt.BaseTime = res.BaseTime           //基准时间
		rpt.ShiftHour = res.ShiftHour         //每班工作时间
		rpt.Status = 1                        //1有效 0无效
		rpt.CreateUserId = userid
		err := rpt.InsertNode()
		if err != nil { //解析请求错误
			c.Data["json"] = err.Error()
		} else {
			c.SaveUserActionMsg("添加报表层级节点", _LOG_OPR_TYPE_INSERT, fmt.Sprintf("%+v", *rpt))
			c.Data["json"] = "层级节点[" + rpt.Name + "]添加成功！"
		}
	}
	c.ServeJSON()
}

/***********************************************
功能:编辑报表层级
输入:
输出:报表节点集
说明:
编辑:wang_jp
时间:2020年5月11日
************************************************/
func (c *ReportController) ApiEditReportLevel() {
	c.CheckSession()      //检查授权,返回授权的ID
	type nodemsg struct { //请求参数结构
		Id            int64  `form:"Id"`
		DistributedId int64  `form:"DistributedId"` //分布式计算ID，为0的时候不区分计算引擎，不为0的时候需要id适配的计算引擎进行计算
		Name          string `form:"Name"`          //名称
		Pid           int64  `form:"Pid"`           //父级菜单ID
		WorkshopId    int    `form:"WorkshopId"`    //所属车间ID
		Level         int64  `form:"Level"`         //层级深度
		LevelCode     string `form:"LevelCode"`     //层级码
		Folder        int64  `form:"Folder"`        //是否是文件夹，1-是，0-否
		Debug         int64  `form:"Debug"`         //是否是调试模式，1-是，0-否
		Seq           int64  `form:"Seq"`           //排序号
		Remark        string `form:"Remark"`        //备注
		TemplateUrl   string `form:"TemplateUrl"`   //模板文件路径
		TemplateFile  string `form:"TemplateFile"`  //模板文件名称
		ResultUrl     string `form:"ResultUrl"`     //结果地址
		StartTime     string `form:"StartTime"`     //统计计算开始起作用的时间
		Period        int64  `form:"Period"`        //计算周期,详见KPI表
		OffsetMinutes int64  `form:"OffsetMinutes"` //偏移时间
		LastCalcTime  string `form:"LastCalcTime"`  //最后计算时间
		BaseTime      string `form:"BaseTime"`      //基准时间
		ShiftHour     int64  `form:"ShiftHour"`     //每班工作时间
		Status        int64  `form:"Status"`        //1有效 0无效
	}

	res := new(nodemsg)     //请求信息
	err := c.ParseForm(res) //解析请求信息
	if err != nil {         //解析请求错误
		c.Data["json"] = err.Error()
	} else { //
		w := new(models.OreProcessDWorkshop)
		w.Id = res.WorkshopId
		//新建节点
		rpt := new(models.CalcKpiReportList)
		rpt.Id = res.Id
		rpt.DistributedId = res.DistributedId //分布式计算ID，为0的时候不区分计算引擎，不为0的时候需要id适配的计算引擎进行计算
		rpt.Name = res.Name                   //名称
		rpt.Pid = res.Pid                     //父级菜单ID
		rpt.Workshop = w
		rpt.Level = res.Level
		rpt.LevelCode = res.LevelCode
		rpt.Folder = res.Folder
		rpt.Debug = res.Debug
		rpt.Seq = res.Seq
		rpt.Remark = res.Remark
		rpt.TemplateUrl = res.TemplateUrl                                  //模板文件路径
		rpt.ResultUrl = res.ResultUrl                                      //结果地址
		rpt.StartTime = strings.Replace(res.StartTime, "T", " ", -1)       //统计计算开始起作用的时间
		rpt.Period = res.Period                                            //计算周期,详见KPI表
		rpt.OffsetMinutes = res.OffsetMinutes                              //偏移时间
		rpt.LastCalcTime = strings.Replace(res.LastCalcTime, "T", " ", -1) //最后计算时间
		rpt.BaseTime = strings.Replace(res.BaseTime, "T", " ", -1)         //基准时间
		rpt.ShiftHour = res.ShiftHour                                      //每班工作时间
		rpt.Status = res.Status                                            //1有效 0无效

		//c.SaveUserActionMsg("临时记录", _LOG_OPR_TYPE_OTHER, fmt.Sprintf("%+v", *rpt))
		orpt := new(models.CalcKpiReportList)
		orpt.Id = rpt.Id
		orpt.GetAttributeById()

		var fieldnames []string
		updateMsg := fmt.Sprintf("修改层级节点ID=%d,Name=%s[", rpt.Id, rpt.Name)
		if rpt.BaseTime != orpt.BaseTime {
			fieldnames = append(fieldnames, "BaseTime")
			updateMsg += fmt.Sprintf("BaseTime:From %s To %s;", orpt.BaseTime, rpt.BaseTime)
		}
		if rpt.DistributedId != orpt.DistributedId {
			fieldnames = append(fieldnames, "DistributedId")
			updateMsg += fmt.Sprintf("DistributedId:From %d To %d;", orpt.DistributedId, rpt.DistributedId)
		}

		if rpt.Folder != orpt.Folder {
			fieldnames = append(fieldnames, "Folder")
			updateMsg += fmt.Sprintf("Folder:From %d To %d;", orpt.Folder, rpt.Folder)
		}
		if rpt.Debug != orpt.Debug {
			fieldnames = append(fieldnames, "Debug")
			updateMsg += fmt.Sprintf("Debug:From %d To %d;", orpt.Debug, rpt.Debug)
		}
		if rpt.LastCalcTime != orpt.LastCalcTime {
			fieldnames = append(fieldnames, "LastCalcTime")
			updateMsg += fmt.Sprintf("LastCalcTime:From %s To %s;", orpt.LastCalcTime, rpt.LastCalcTime)
		}
		if rpt.Name != orpt.Name {
			fieldnames = append(fieldnames, "Name")
			updateMsg += fmt.Sprintf("Name:From %s To %s;", orpt.Name, rpt.Name)
		}
		if rpt.OffsetMinutes != orpt.OffsetMinutes {
			fieldnames = append(fieldnames, "OffsetMinutes")
			updateMsg += fmt.Sprintf("OffsetMinutes:From %d To %d;", orpt.OffsetMinutes, rpt.OffsetMinutes)
		}
		if rpt.Period != orpt.Period {
			fieldnames = append(fieldnames, "Period")
			updateMsg += fmt.Sprintf("Period:From %d To %d;", orpt.Period, rpt.Period)
		}

		if rpt.Pid != orpt.Pid && rpt.Folder == 0 { //不是文件夹的情况下父层级才可修改
			fieldnames = append(fieldnames, "Pid")
			updateMsg += fmt.Sprintf("Pid:From %d To %d;", orpt.Pid, rpt.Pid)

			prpt := new(models.CalcKpiReportList)
			prpt.Id = rpt.Pid
			e := prpt.GetAttributeById()
			if e == nil { //父层级修改时,要同时修改Level和LevelCode
				rpt.Level = prpt.Level + 1
				rpt.LevelCode = fmt.Sprintf("%s-%d", prpt.LevelCode, rpt.Id)
				fieldnames = append(fieldnames, "Level")
				fieldnames = append(fieldnames, "LevelCode")
				updateMsg += fmt.Sprintf("Level:From %d To %d;", orpt.Level, rpt.Level)
				updateMsg += fmt.Sprintf("LevelCode:From %s To %s;", orpt.LevelCode, rpt.LevelCode)
			}
		}
		if rpt.Remark != orpt.Remark {
			fieldnames = append(fieldnames, "Remark")
			updateMsg += fmt.Sprintf("Remark:From %s To %s;", orpt.Remark, rpt.Remark)
		}
		if rpt.ResultUrl != orpt.ResultUrl {
			fieldnames = append(fieldnames, "ResultUrl")
			updateMsg += fmt.Sprintf("ResultUrl:From %s To %s;", orpt.ResultUrl, rpt.ResultUrl)
		}
		if rpt.Seq != orpt.Seq {
			fieldnames = append(fieldnames, "Seq")
			updateMsg += fmt.Sprintf("Seq:From %d To %d;", orpt.Seq, rpt.Seq)
		}
		if rpt.ShiftHour != orpt.ShiftHour {
			fieldnames = append(fieldnames, "ShiftHour")
			updateMsg += fmt.Sprintf("ShiftHour:From %d To %d;", orpt.ShiftHour, rpt.ShiftHour)
		}
		if rpt.StartTime != orpt.StartTime {
			fieldnames = append(fieldnames, "StartTime")
			updateMsg += fmt.Sprintf("StartTime:From %s To %s;", orpt.StartTime, rpt.StartTime)
		}
		if rpt.Status != orpt.Status {
			fieldnames = append(fieldnames, "Status")
			updateMsg += fmt.Sprintf("Status:From %d To %d;", orpt.Status, rpt.Status)
		}
		if rpt.TemplateUrl != orpt.TemplateUrl {
			fieldnames = append(fieldnames, "TemplateUrl")
			updateMsg += fmt.Sprintf("TemplateUrl:From %s To %s;", orpt.TemplateUrl, rpt.TemplateUrl)
		}
		if rpt.Workshop != nil && orpt.Workshop != nil {
			if rpt.Workshop.Id != orpt.Workshop.Id {
				fieldnames = append(fieldnames, "Workshop")
				updateMsg += fmt.Sprintf("Workshop:From %d To %d;", orpt.Workshop.Id, rpt.Workshop.Id)
			}
		}
		updateMsg += "]"

		err := rpt.Update(fieldnames...) //执行更新

		if err != nil { //解析请求错误
			c.Data["json"] = err.Error()
		} else {
			c.SaveUserActionMsg("修改报表层级节点", _LOG_OPR_TYPE_UPDATE, updateMsg)
			c.Data["json"] = "层级节点[" + rpt.Name + "]修改成功！"
		}
	}
	c.ServeJSON()
}

/***********************************************
功能:设置文件作为报表模板
输入:
	[id] 报表ID
	[reportname]  报表名称
	[oldfilename] 新文件名
	[newfilename] 新文件名
输出:
说明:
编辑:wang_jp
时间:2020年5月9日
************************************************/
func (c *ReportController) ApiSetFileAsReportTpl() {
	c.CheckSession() //检查授权,返回授权的ID
	newfilename := c.GetString("newfilename")
	oldfilename := c.GetString("oldfilename")
	reportname := c.GetString("reportname")
	id, err := c.GetInt64("id")
	c.SaveUserActionMsg("更新模板文件名", _LOG_OPR_TYPE_UPDATE, fmt.Sprintf("将报表[Id:%d,Name:%s]的模板文件由[%s]修改为[%s]", id, reportname, oldfilename, newfilename)) //记录信息
	if err != nil {
		c.Data["json"] = fmt.Sprintf("Id必须是整数,传入的是:[%s]", c.GetString("id"))
	} else {
		rpt := new(models.CalcKpiReportList) //新建报表结构
		rpt.Id = id
		rpt.TemplateFile = newfilename
		err := rpt.Update("TemplateFile") //更新模板文件名
		if err == nil {
			c.SaveUserActionMsg("更新模板文件名", _LOG_OPR_TYPE_UPDATE, fmt.Sprintf("将报表[Id:%d,Name:%s]的模板文件由[%s]修改为[%s]", id, reportname, oldfilename, newfilename)) //记录信息
			c.Data["json"] = fmt.Sprintf("将报表[%s]的模板文件由[%s]修改为[%s]", reportname, oldfilename, newfilename)
		} else {
			c.Data["json"] = err.Error()
		}
	}
	c.ServeJSON()
}

/***********************************************
功能:删除报表层级
输入:[id]
输出:报表节点集
说明:
编辑:wang_jp
时间:2020年5月9日
************************************************/
func (c *ReportController) ApiDeleteReportLevel() {
	c.CheckSession() //检查授权,返回授权的ID
	id, err := c.GetInt64("id")
	if err != nil {
		c.Data["json"] = fmt.Sprintf("Id必须是整数,传入的是:[%s]", c.GetString("id"))
	} else {
		rpt := new(models.CalcKpiReportList) //新建报表结构
		rpt.Id = id
		err := rpt.GetAttributeById() //读取节点信息
		if err == nil {
			c.SaveUserActionMsg("删除报表层级节点", _LOG_OPR_TYPE_DELETE, fmt.Sprintf("%+v", *rpt)) //记录信息
			err := rpt.DeleteNode()
			if err != nil {
				c.Data["json"] = err.Error()
			} else {
				c.Data["json"] = "ok"
			}
		} else {
			c.Data["json"] = err.Error()
		}
	}
	c.ServeJSON()
}

/***********************************************
功能:下载文件
输入:
	[filepath] 文件路径
	[filename] 文件名
输出:
说明:
编辑:wang_jp
时间:2020年5月11日
************************************************/
func (c *ReportController) ApiDownLoadFile() {
	//c.CheckSession() //检查授权,返回授权的ID
	filepath := c.GetString("filepath")
	filename := c.GetString("filename")
	id := c.GetString("id")
	if filepath[len(filepath)-1:] != "/" { //如果最后一个字符不是"/"
		filepath += "/"
	}
	if id != "0" {
		filepath += id + "/"
	}
	filepath += filename
	c.SaveUserActionMsg("下载文件", _LOG_OPR_TYPE_OTHER, filepath) //记录信息
	//第一个参数是文件的地址，第二个参数是下载显示的文件的名称
	c.Ctx.Output.ContentType("application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Ctx.Output.Download(filepath, filename)
}

/***********************************************
功能:上传文件
输入:
	[filepath] 文件路径
	[filename] 文件名
输出:
说明:
编辑:wang_jp
时间:2020年5月11日
************************************************/
func (c *ReportController) ApiUpLoadFile() {
	c.CheckSession() //检查授权,返回授权的ID
	filepath := c.GetString("filepath")
	reportname := c.GetString("reportname")
	id := c.GetString("id")
	defaulttpl := c.GetString("defaulttpl")
	c.SaveUserActionMsg("上传模板", _LOG_OPR_TYPE_OTHER, filepath, reportname, id, defaulttpl) //记录信息
	if filepath[len(filepath)-1:] != "/" {                                                 //如果最后一个字符不是"/"
		filepath += "/"
	}

	file, head, err := c.GetFile("file")
	if err != nil {
		c.Ctx.WriteString("获取文件失败")
		return
	}
	defer file.Close()

	rawname := head.Filename
	filename := fmt.Sprintf("tpl-%s_%s.xlsx", id, time.Now().Format("20060102150405"))
	err = c.SaveToFile("file", filepath+filename)
	if err != nil {
		c.Ctx.WriteString(fmt.Sprintf("上传失败:[%s]", err.Error()))
	} else {
		if defaulttpl == "on" {
			rpt := new(models.CalcKpiReportList)
			rpt.Id, _ = strconv.ParseInt(id, 10, 64)
			rpt.TemplateFile = filename
			rpt.Update("TemplateFile")
		}
		c.Ctx.WriteString("上传成功")
		c.SaveUserActionMsg("上传模板", _LOG_OPR_TYPE_OTHER, fmt.Sprintf("为报表[Id:%s,Name:%s]上传模板文件,文件原名[%s],系统另存为[%s]", id, reportname, rawname, filename)) //记录信息
	}
}

/***********************************************
功能:在线预览文件
输入:
	[filepath] 文件路径
	[filename] 文件名
输出:
说明:
编辑:wang_jp
时间:2020年5月11日
************************************************/
func (c *ReportController) ApiViewExcel() {
	//c.CheckSession() //检查授权,返回授权的ID
	filepath := c.GetString("filepath")
	filename := c.GetString("filename")
	id := c.GetString("id")
	if filepath[len(filepath)-1:] != "/" { //如果最后一个字符不是"/"
		filepath += "/"
	}
	calcformula := false //是否计算Excel公式
	if id != "0" {       //带有ID,是报表
		filepath += id + "/"
		calcformula = true
	}

	filepath += filename
	c.SaveUserActionMsg("查看报表", _LOG_OPR_TYPE_OTHER, filepath) //记录信息
	table := new(extable.Table)
	err := table.OpenFile(filepath, calcformula, models.EngineCfgMsg.Sys.ExcelFormulaCalcDeep)
	if err != nil {
		c.Ctx.WriteString(err.Error())
	} else {
		c.Ctx.WriteString(table.FomatToHtml())
	}
}

/***********************************************
功能:获取车间列表
输入:无
输出:车间列表
说明:
编辑:wang_jp
时间:2020年5月9日
************************************************/
func (c *ReportController) ApiGetWorkShops() {
	wshp := new(models.OreProcessDWorkshop)
	workshops, err := wshp.GetWorkshopLists()
	if err != nil {
		c.Data["json"] = err.Error()
	} else {
		c.Data["json"] = workshops
	}
	c.ServeJSON()
}

/***********************************************
功能:保存单元格的值
输入:
	[filepath] 文件路径
	[filename] 文件名
	[cellaxis] 单元格坐标字符串
	[cellvalue] 单元格的值
输出:
说明:POST方法
编辑:wang_jp
时间:2020年5月11日
************************************************/
func (c *ReportController) ApiSetExcelCellValue() {
	c.CheckSession() //检查授权,返回授权的ID
	filepath := c.GetString("filepath")
	filename := c.GetString("filename")
	cellaxis := c.GetString("cellaxis")
	cellvalue := c.GetString("cellvalue")

	if filepath[len(filepath)-1:] != "/" { //如果最后一个字符不是"/"
		filepath += "/"
	}
	filepath += filename
	cell := new(extable.TableCell)
	cell.Axis = cellaxis
	cell.Value = cellvalue
	err := cell.SetCellValue(filepath)
	if err != nil {
		c.Ctx.WriteString(err.Error())
	} else {
		c.SaveUserActionMsg("编辑报表模板单元格的内容",
			_LOG_OPR_TYPE_UPDATE,
			fmt.Sprintf("模板文件:[%s]单元格[%s]中的内容变更为[%s]", filename, cellaxis, cellvalue))
		c.Ctx.WriteString("ok")
	}
}

/***********************************************
功能:保存单元格的公式
输入:
	[filepath] 文件路径
	[filename] 文件名
	[cellaxis] 单元格坐标字符串
	[cellformula] 单元格的值
输出:
说明:POST方法
编辑:wang_jp
时间:2020年5月11日
************************************************/
func (c *ReportController) ApiSetExcelCellFormula() {
	c.CheckSession() //检查授权,返回授权的ID
	filepath := c.GetString("filepath")
	filename := c.GetString("filename")
	cellaxis := c.GetString("cellaxis")
	cellformula := c.GetString("cellformula")

	if filepath[len(filepath)-1:] != "/" { //如果最后一个字符不是"/"
		filepath += "/"
	}
	filepath += filename
	cell := new(extable.TableCell)
	cell.Axis = cellaxis
	cell.Formula = cellformula
	err := cell.SetCellFormula(filepath)
	if err != nil {
		c.Ctx.WriteString(err.Error())
	} else {
		c.SaveUserActionMsg("编辑报表模板单元格的Excel公式",
			_LOG_OPR_TYPE_UPDATE,
			fmt.Sprintf("模板文件:[%s]单元格[%s]中的公式变更为[%s]", filename, cellaxis, cellformula))
		c.Ctx.WriteString("ok")
	}
}
/***********************************************
功能:获取数据
输入:[levelcode]
输出:报表节点集
说明:
编辑:wangzf
时间:2020年10月14日
************************************************/
func (c *ReportController) ApiGetCalculate() {
	data := c.GetString("data")
	begin := c.GetString("begin")
	end := c.GetString("end")
	fc := c.GetString("fc")
	tag_id := c.GetString("tag_id")
	objtype := c.GetString("objtype")
	fz := c.GetString("fz")

	rpt := new(models.KpiArithmeticResult)
	ari := new(models.KpiArithmetic)
	ariType := ari.ReaderArithmeticData(fc)
	rpt.ResultName = ariType[1]
	rpt.ArithmeticId = ariType[0]
	rpt.FinalResultType = ariType[2]

	rpt.ArithmeticObjectType = objtype
	rpt.ArithmeticAuxiliaryId1 = tag_id
	rpt.ArithmeticAuxiliaryRemark1 = fz

	rpt.ArithmeticObjectId = tag_id
	rpt.CreTime = time.Now().Format("2006-01-02 15:04:05")
	rpt.FinalResult = data
	rpt.BeginTime = begin
	rpt.EndTime = end
	rpt.SaveDataLog(rpt)

	c.Ctx.WriteString("")
}