package controllers

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/bkzy-wangjp/MicEngine/models"
)

const (
	_SERVER_ERR    = 500 //服务器反馈错误信息
	_PAROMATER_ERR = 400 //输入参数错误
)

type EngineController struct {
	MyController
}

type engineScriptCmd struct {
	Script     string //脚本
	BeginTime  string //起始时间
	EndTime    string //结束时间
	BaseTime   string //基准时间
	ShiftHours int64  //每班工作时间
}

/***********************************************
功能:API检查
输入:无
输出:
说明:POST/GET	方法
编辑:wang_jp
时间:2020年6月1日
************************************************/
func (this *EngineController) ApiTest() {
	this.Data["json"] = 1
	this.ServeJSON()
}

/***********************************************
功能:运行脚本
输入:
	[script] 脚本程序
	[begintime] 开始时间
	[endtime] 结束时间
	[basetime] 基准时间
	[shifthours] 每班工作时间
输出:
说明:POST/GET	方法
编辑:wang_jp
时间:2020年5月25日
************************************************/
func (this *EngineController) ScriptRun() {
	if cmd, err := engineParemater(this.Input()); err != nil {
		this.Data["json"] = err.Error()
	} else {
		script := new(models.Script)
		script.MainTagId = 0
		script.MainTagFullName = ""
		script.ShiftHour = cmd.ShiftHours
		script.ScriptStr = cmd.Script
		script.BeginTime = cmd.BeginTime
		script.BaseTime = cmd.BaseTime
		script.EndTime = cmd.EndTime
		responsData, _, err := script.Run()
		if err != nil {
			this.Data["json"] = err.Error()
		} else {
			this.Data["json"] = responsData
		}
	}
	this.ServeJSON()
}

/***********************************************
功能:运行SQL脚本
输入:
	[micsql] sql脚本程序
输出:
说明:POST/GET	方法
编辑:wang_jp
时间:2020年07月07日
************************************************/
func (this *EngineController) ScriptMicSql() {
	scriptstr := this.GetString("micsql")
	if len(scriptstr) < 15 {
		this.Ctx.ResponseWriter.WriteHeader(_PAROMATER_ERR)
		this.Data["json"] = "没有输入合法的SQL脚本"
	} else {
		script := new(models.Script)
		script.ScriptStr = strings.ReplaceAll(scriptstr, "as(json)", "as(map)")

		responsData, _, err := script.Run()
		if err != nil {
			this.Ctx.ResponseWriter.WriteHeader(_SERVER_ERR)
			this.Data["json"] = err.Error()
		} else {
			this.Data["json"] = responsData
		}
	}
	this.ServeJSON()
}

/***********************************************
功能:编译校验脚本
输入:
	[script] 脚本程序
输出:
说明:POST/GET	方法
编辑:wang_jp
时间:2020年5月25日
************************************************/
func (this *EngineController) ScriptCompile() {
	scriptstr := this.GetString("script")
	scrpts := strings.Split(scriptstr, ";")
	var errstr string //错误信息
	if len(scriptstr) > 0 {
		for _, scr := range scrpts {
			script := new(models.Script)
			script.ScriptStr = scr
			errs := script.Compile()
			if len(errs) > 0 {
				errstr += fmt.Sprintf("脚本[%s]中的错误:</br>", scr)
				for _, e := range errs {
					errstr += fmt.Sprintf("---%s;</br>", e.Error())
				}
			}
		}
		if len(errstr) > 0 {
			this.Ctx.WriteString(errstr)
		} else {
			this.Ctx.WriteString("ok")
		}
	} else {
		this.Ctx.WriteString("没有输入有效的脚本程序")
	}
}

/***********************************************
功能:获取系统配置信息
输入:无
输出:
说明:POST/GET	方法
编辑:wang_jp
时间:2020年5月25日
************************************************/
func (this *EngineController) GetConfig() {
	this.Data["json"] = models.EngineCfgMsg
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
func engineParemater(input url.Values) (engineScriptCmd, error) {
	var cmd engineScriptCmd
	var err error
	for k, v := range input {
		switch strings.ToLower(k) {
		case "script":
			for i, od := range v {
				if i == 0 {
					cmd.Script = strings.ToLower(od)
				} else {
					err = fmt.Errorf("Too many script parameters[脚本参数太多了]")
					return cmd, err
				}
			}
			if len(v) == 0 {
				err = fmt.Errorf("No script parameter[没有 script 参数]")
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
		case "basetime":
			for i, s := range v {
				if i == 0 {
					if _, e := models.TimeParse(s); e != nil {
						return cmd, fmt.Errorf("basetime fomate error[basetime 格式错误];%s", e.Error())
					}
					cmd.BaseTime = s
				} else {
					err = fmt.Errorf("Too many basetime parameters[basetime 参数太多了]")
					return cmd, err
				}
			}
			if len(v) == 0 {
				err = fmt.Errorf("No basetime parameter[没有 basetime 参数]")
				return cmd, err
			}
		case "shifthours":
			for i, s := range v {
				if i == 0 {
					cmd.ShiftHours, err = strconv.ParseInt(s, 10, 64)
					if err != nil {
						return cmd, fmt.Errorf("shifthours fomate error,need int[shifthours 格式错误,应为整数];%s", err.Error())
					}
				} else {
					err = fmt.Errorf("Too many shifthours parameters[shifthours 参数太多了]")
					return cmd, err
				}
			}
			if len(v) == 0 {
				err = fmt.Errorf("No shifthours parameter[没有 shifthours 参数]")
				return cmd, err
			}
		default:
			err = fmt.Errorf("There are unrecognized parameters[有未识别的参数]")
			return cmd, err
		}
	}
	return cmd, err
}

/***********************************************
功能:新建表
输入:[tablename] 表名称
	[tabledesc] 表描述
输出:
说明:POST/GET	方法
编辑:wang_jp
时间:2020年6月1日
************************************************/
func (this *EngineController) GoldenTableInsert() {
	tablename := this.GetString("tablename")
	tabledesc := this.GetString("tabledesc", "")
	if len(tablename) == 0 {
		this.Ctx.ResponseWriter.WriteHeader(_PAROMATER_ERR)
		this.Data["json"] = "没有输入合法的表名"
	} else {
		micgd := new(models.MicGolden)
		id, err := micgd.GoldenTableInsert(tablename, tabledesc)
		if err != nil {
			this.Ctx.ResponseWriter.WriteHeader(_SERVER_ERR)
			this.Data["json"] = err.Error()
		} else {
			type tb struct {
				Id   int
				Name string
				Desc string
			}
			table := tb{Id: id, Name: tablename, Desc: tabledesc}

			this.Data["json"] = table
			this.SaveUserActionMsg("新建庚顿数据表", _LOG_OPR_TYPE_INSERT, fmt.Sprintf("TableName=%s,TableDesc=%s,TableId=%d", tablename, tabledesc, id))
		}
	}
	this.ServeJSON()
}

/***********************************************
功能:新建表或者更新表
输入:[tablename] 表名称
	[tabledesc] 表描述
输出:
说明:POST/GET	方法
编辑:wang_jp
时间:2020年6月18日
************************************************/
func (this *EngineController) GoldenTableInsertOrUpdate() {
	tablename := this.GetString("tablename")
	tabledesc := this.GetString("tabledesc", "")
	if len(tablename) == 0 {
		this.Ctx.ResponseWriter.WriteHeader(_PAROMATER_ERR)
		this.Data["json"] = "没有输入合法的表名"
	} else {
		micgd := new(models.MicGolden)
		isinsert, id, err := micgd.GoldenTableInsertOrUpdate(tablename, tabledesc)
		if err != nil {
			this.Ctx.ResponseWriter.WriteHeader(_SERVER_ERR)
			this.Data["json"] = err.Error()
		} else {
			type tb struct {
				Id    int
				Name  string
				Desc  string
				IsNew bool
			}
			table := tb{Id: id, Name: tablename, Desc: tabledesc, IsNew: isinsert}

			this.Data["json"] = table
			if isinsert {
				this.SaveUserActionMsg("新建庚顿数据表", _LOG_OPR_TYPE_INSERT, fmt.Sprintf("TableName=%s,TableDesc=%s,TableId=%d", tablename, tabledesc, id))
			} else {
				this.SaveUserActionMsg("更新庚顿数据表描述", _LOG_OPR_TYPE_UPDATE, fmt.Sprintf("TableName=%s,TableDesc=%s,TableId=%d", tablename, tabledesc, id))
			}
		}
	}
	this.ServeJSON()
}

/***********************************************
功能:以平台数据表为基准创建庚顿表
输入:无
输出:
说明:POST/GET	方法
编辑:wang_jp
时间:2020年6月1日
************************************************/
func (this *EngineController) GoldenTablesFromPlat() {
	tb := new(models.DatatableInfo)
	total, newtb, err := tb.SynchGoldenTableAndPlatTable()
	if err != nil {
		this.Ctx.ResponseWriter.WriteHeader(_SERVER_ERR)
		this.Data["json"] = err.Error()
	} else {
		mp := make(map[string]interface{})
		mp["TotalTalbes"] = total
		mp["NewTables"] = newtb
		mp["UpdateTables"] = total - newtb
		this.Data["json"] = mp
		this.SaveUserActionMsg("同步平台数据表到庚顿", _LOG_OPR_TYPE_INSERT, fmt.Sprintf("共[%d]张表,在庚顿中新建[%d]张,更新[%d]张", total, newtb, total-newtb))
	}
	this.ServeJSON()
}

/***********************************************
功能:删除表
输入:无
输出:
说明:POST/GET	方法
编辑:wang_jp
时间:2020年6月1日
************************************************/
func (this *EngineController) GoldenTableDelete() {
	tablename := this.GetString("tablename")
	if len(tablename) == 0 {
		this.Ctx.ResponseWriter.WriteHeader(_PAROMATER_ERR)
		this.Data["json"] = "没有输入合法的表名"
	} else {
		micgd := new(models.MicGolden)
		err := micgd.GoldenTableRemove(tablename)
		if err != nil {
			this.Ctx.ResponseWriter.WriteHeader(_SERVER_ERR)
			this.Data["json"] = err.Error()
		} else {
			this.Data["json"] = fmt.Sprintf("删除表[%s]成功!", tablename)
			this.SaveUserActionMsg("删除庚顿数据表", _LOG_OPR_TYPE_DELETE, fmt.Sprintf("TableName=%s", tablename))
		}
	}
	this.ServeJSON()
}

/***********************************************
功能:更新表
输入:无
输出:
说明:POST/GET	方法
编辑:wang_jp
时间:2020年6月1日
************************************************/
func (this *EngineController) GoldenTableUpdate() {
	tablename := this.GetString("tablename")
	newname := this.GetString("newname")
	if len(tablename) == 0 || len(newname) == 0 {
		this.Ctx.ResponseWriter.WriteHeader(_PAROMATER_ERR)
		this.Data["json"] = "没有输入合法的表名"
	} else {
		micgd := new(models.MicGolden)
		err := micgd.GoldenTableReName(tablename, newname)
		if err != nil {
			this.Ctx.ResponseWriter.WriteHeader(_SERVER_ERR)
			this.Data["json"] = err.Error()
		} else {
			this.Data["json"] = fmt.Sprintf("将庚顿数据表[%s]更名为[%s]成功!", tablename, newname)
			this.SaveUserActionMsg("重命名庚顿数据表", _LOG_OPR_TYPE_UPDATE, fmt.Sprintf("oldTableName=%s,newTableName=%s", tablename, newname))
		}
	}
	this.ServeJSON()
}

/***********************************************
功能:查询表
输入:无
输出:
说明:POST/GET	方法
编辑:wang_jp
时间:2020年6月1日
************************************************/
func (this *EngineController) GoldenTableSelect() {
	tablenames := this.GetString("tablenames", "")
	var names []string
	if tablenames != "" {
		names = strings.Split(tablenames, ",")
	}
	micgd := new(models.MicGolden)
	tables, err := micgd.GoldenGetTablePropertyByTableName(names...)
	if err != nil {
		this.Data["json"] = err.Error()
	} else {
		this.Data["json"] = tables
	}
	this.ServeJSON()
}

/***********************************************
功能:以平台标签Taglist为基准创建庚顿标签
输入:无
输出:
说明:POST/GET	方法
编辑:wang_jp
时间:2020年6月22日
************************************************/
func (this *EngineController) GoldenPointsFromPlat() {
	tag := new(models.OreProcessDTaglist)
	total, newcnt, upcnt, tsec, err := tag.TaglistsToGolden()
	if err != nil {
		this.Ctx.ResponseWriter.WriteHeader(_SERVER_ERR)
		this.Data["json"] = fmt.Sprintf("同步平台Taglist到庚顿数据库时遇到错误:[%s]", err.Error())
	} else {
		mp := make(map[string]interface{})
		mp["TotalPoints"] = total
		mp["InsertPoints"] = newcnt
		mp["UpdatePoints"] = upcnt
		mp["UseTime"] = tsec
		this.Data["json"] = mp
		this.SaveUserActionMsg("同步平台TagList到庚顿", _LOG_OPR_TYPE_INSERT, fmt.Sprintf("共[%d]标签,在庚顿中新建[%d]个,更新[%d]个,耗时[%f]秒", total, newcnt, upcnt, tsec))
	}
	this.ServeJSON()
}

/***********************************************
功能:新建点或者更新点
输入:标签点ID
输出:
说明:POST/GET	方法
编辑:wang_jp
时间:2020年6月1日
************************************************/
func (this *EngineController) GoldenPointInsert() {
	tagid, err := this.GetInt64("tagid")
	if err != nil {
		this.Ctx.ResponseWriter.WriteHeader(_PAROMATER_ERR)
		this.Data["json"] = fmt.Sprintf("参数tagid必须为整数类型,输入的是:[%s]", this.GetString("tagid"))
	} else {
		tag := new(models.OreProcessDTaglist)
		tag.Id = tagid
		err := tag.GetTagAttributByTagId()
		if err != nil {
			this.Ctx.ResponseWriter.WriteHeader(_SERVER_ERR)
			this.Data["json"] = fmt.Sprintf("获取平台标签属性失败:[%s]", err.Error())
		} else {
			gdp, isnew, err := tag.InsertOrUpdateTagToGolden()
			if err != nil {
				this.Ctx.ResponseWriter.WriteHeader(_SERVER_ERR)
				this.Data["json"] = fmt.Sprintf("插入或者更新庚顿标签失败:[%s]", err.Error())
			} else {
				type dt struct {
					ActionType  string
					GoldenPoint interface{}
				}
				data := new(dt)
				data.GoldenPoint = gdp
				data.ActionType = "Update"
				act := "更新庚顿标签点"
				var oprtype int64 = _LOG_OPR_TYPE_UPDATE
				if isnew {
					data.ActionType = "Insert"
					act = "新建庚顿标签点"
					oprtype = _LOG_OPR_TYPE_INSERT
				}
				this.Data["json"] = data
				this.SaveUserActionMsg(act, oprtype, fmt.Sprintf("PlatTagId=%d,GoldenTagId=%d,TagName=%s", tagid, gdp.Base.Id, gdp.Base.TableDotTag))
			}
		}
	}
	this.ServeJSON()
}

/***********************************************
功能:删除点
输入:无
输出:
说明:POST/GET	方法
编辑:wang_jp
时间:2020年6月1日
************************************************/
func (this *EngineController) GoldenPointDelete() {
	point := this.GetString("point")
	id, err := strconv.ParseInt(point, 10, 64)
	if err != nil {
		micgd := new(models.MicGolden)
		gdp, err := micgd.GoldenGetTagPointInfoByName(point)
		if err != nil {
			this.Ctx.ResponseWriter.WriteHeader(_SERVER_ERR)
			this.Data["json"] = fmt.Sprintf("获取庚顿标签点属性失败:[%s]", err.Error())
		} else {
			id = int64(gdp[point].Base.Id)
		}
	}
	tag := new(models.OreProcessDTaglist)
	err = tag.RemoveGoldenTagByGoldenId(id)
	if err != nil {
		this.Ctx.ResponseWriter.WriteHeader(_SERVER_ERR)
		this.Data["json"] = fmt.Sprintf("删除庚顿标签点失败:[%s]", err.Error())
	} else {
		this.Data["json"] = "ok"
		this.SaveUserActionMsg("删除庚顿标签点", _LOG_OPR_TYPE_DELETE, fmt.Sprintf("GoldenPoint=%s", point))
	}

	this.ServeJSON()
}

/***********************************************
功能:更新点
输入:无
输出:
说明:POST/GET	方法
编辑:wang_jp
时间:2020年6月1日
************************************************/
func (this *EngineController) GoldenPointUpdate() {
	tagid, err := this.GetInt64("tagid")
	if err != nil {
		this.Ctx.ResponseWriter.WriteHeader(_PAROMATER_ERR)
		this.Data["json"] = err.Error()
	} else {
		tag := new(models.OreProcessDTaglist)
		tag.Id = tagid
		err := tag.GetTagAttributByTagId()
		if err != nil {
			this.Ctx.ResponseWriter.WriteHeader(_SERVER_ERR)
			this.Data["json"] = err.Error()
		} else {
			gdp, err := tag.UpdateTagToGolden()
			if err != nil {
				this.Ctx.ResponseWriter.WriteHeader(_SERVER_ERR)
				this.Data["json"] = err.Error()
			} else {
				this.Data["json"] = gdp
				this.SaveUserActionMsg("更新庚顿标签点", _LOG_OPR_TYPE_UPDATE, fmt.Sprintf("PlatTagId=%d,GoldenTagId=%d,TagName=%s", tagid, gdp.Base.Id, gdp.Base.TableDotTag))
			}
		}
	}
	this.ServeJSON()
}

/***********************************************
功能:查询点
输入:无
输出:
说明:POST/GET	方法
编辑:wang_jp
时间:2020年6月1日
************************************************/
func (this *EngineController) GoldenPointSelect() {
	point := this.GetString("point")
	id, err := strconv.ParseInt(point, 10, 64)
	if err != nil {
		micgd := new(models.MicGolden)
		gdp, err := micgd.GoldenGetTagPointInfoByName(point)
		if err != nil {
			this.Ctx.ResponseWriter.WriteHeader(_SERVER_ERR)
			this.Data["json"] = err.Error()
		} else {
			this.Data["json"] = gdp[point]
		}
	} else {
		micgd := new(models.MicGolden)
		gdp, err := micgd.GoldenGetTagPointInfoById(int(id))
		if err != nil {
			this.Ctx.ResponseWriter.WriteHeader(_SERVER_ERR)
			this.Data["json"] = err.Error()
		} else {
			this.Data["json"] = gdp
		}
	}
	this.ServeJSON()
}

/***********************************************
功能:查询含有报警信息和快照信息的单个标签点信息
输入:[pointid] 标签点ID
输出:
说明:POST/GET	方法
编辑:wang_jp
时间:2020年6月1日
************************************************/
func (this *EngineController) GoldenPointAlarm() {
	pointid, err := this.GetInt("pointid")
	if err != nil {
		this.Ctx.ResponseWriter.WriteHeader(_PAROMATER_ERR)
		this.Data["json"] = err.Error()
	} else {
		gdp, err := models.GDDB.SelectPoint(pointid)
		if err != nil {
			this.Ctx.ResponseWriter.WriteHeader(_SERVER_ERR)
			this.Data["json"] = err.Error()
		} else {
			this.Data["json"] = gdp
		}
	}
	this.ServeJSON()
}

/***********************************************
功能:查看所有庚顿点的配置
输入:无
输出:
说明:POST/GET	方法
编辑:wang_jp
时间:2020年6月1日
************************************************/
func (this *EngineController) GoldenPointsAll() {
	this.Data["json"] = models.GDDB
	this.ServeJSON()
}

/***********************************************
功能:同步庚顿数据库的标签点与平台数据库taglist标签点
输入:无
输出:
说明:POST/GET	方法
编辑:wang_jp
时间:2020年6月6日
************************************************/
func (this *EngineController) SynchGoldenPointAndPlatTaglist() {
	tag := new(models.OreProcessDTaglist)
	total, g2p, p2g, millsec, err := tag.SynchGoldenPointAndPlatTaglist()
	if err != nil {
		this.Data["json"] = err.Error()
	} else {
		mp := make(map[string]interface{})
		mp["TotalPoints"] = total
		mp["GoldenToPlatPoints"] = g2p
		mp["PlatToGoldenPoints"] = p2g
		mp["UseTime"] = float64(millsec) / 1000.0
		this.Data["json"] = mp
		this.SaveUserActionMsg("同步平台TagList和庚顿", _LOG_OPR_TYPE_INSERT, fmt.Sprintf("共[%d]标签,庚顿同步到平台[%d]个,平台同步到庚顿[%d]个,耗时[%f]秒", total, g2p, p2g, float64(millsec)/1000.0))
	}
	this.ServeJSON()
}

/***********************************************
功能:从平台数据库taglist标签点同步信息到庚顿数据库的标签点
输入:无
输出:
说明:POST/GET	方法
编辑:wang_jp
时间:2020年6月6日
************************************************/
func (this *EngineController) SynchGoldenPointFromPlatTaglist() {
	tag := new(models.OreProcessDTaglist)
	total, p2g, millsec, err := tag.SynchGoldenPointFromePlatTaglist()
	if err != nil {
		this.Data["json"] = err.Error()
	} else {
		mp := make(map[string]interface{})
		mp["TotalPoints"] = total
		mp["PlatToGoldenPoints"] = p2g
		mp["UseTime"] = float64(millsec) / 1000.0
		this.Data["json"] = mp
		this.SaveUserActionMsg("从平台同步庚顿中已存在的标签点信息", _LOG_OPR_TYPE_INSERT, fmt.Sprintf("共[%d]标签,平台同步到庚顿[%d]个,耗时[%f]秒", total, p2g, float64(millsec)/1000.0))
	}
	this.ServeJSON()
}

/***********************************************
功能:同步单个平台数据库taglist标签点信息到庚顿数据库
输入:[tagid] 标签点ID
输出:
说明:POST/GET	方法
编辑:wang_jp
时间:2020年6月6日
************************************************/
// func (this *EngineController) SynchSiglePlatTagToGolden() {
// 	tagid, err := this.GetInt64("tagid")
// 	if err != nil {
// 		this.Data["json"] = err.Error()
// 	} else {
// 		tag := new(models.OreProcessDTaglist)
// 		tag.Id = tagid
// 		err := tag.GetTagAttributByTagId()
// 		if err != nil {
// 			this.Data["json"] = err.Error()
// 		} else {
// 			err := tag.SynchParameterToGolden()
// 			if err != nil {
// 				this.Data["json"] = err.Error()
// 			} else {
// 				gdp, _ := models.GoldenGetTagPointInfoById(tag.GoldenId)
// 				this.Data["json"] = gdp
// 			}
// 		}
// 	}
// 	this.ServeJSON()
// }

/***********************************************
功能:获取庚顿数据库连接池的信息
输入:无
输出:
说明:POST/GET	方法
编辑:wang_jp
时间:2020年6月6日
************************************************/
func (this *EngineController) GetGoldenPoolMsg() {
	mp := make(map[string]interface{})
	mp["MaxConnect"] = models.EngineCfgMsg.Sys.GoldenCennectPool
	mp["Worker"] = len(models.GDPOOL.Worker)
	mp["Handel"] = len(models.GDPOOL.Handel)
	mp["Req"] = len(models.GDPOOL.Req)
	mp["Tenant"] = models.GDPOOL.ShowTenant()
	this.Data["json"] = mp
	this.ServeJSON()
}

/***********************************************
功能:获取庚顿数据库服务器时间
输入:无
输出:
说明:POST/GET	方法
编辑:wang_jp
时间:2020年6月6日
************************************************/
func (this *EngineController) GetGoldenHostTime() {
	mp := make(map[string]interface{})
	micgd := new(models.MicGolden)
	ht, err := micgd.GoldenGetHostTime()
	if err != nil {
		this.Data["json"] = err.Error()
	} else {
		mp["HostTime"] = time.Unix(ht, 0).Format(models.EngineCfgMsg.Sys.TimeFormat)
		this.Data["json"] = mp
	}
	this.ServeJSON()
}

/***********************************************
功能:写庚顿快照
输入:[tagnames]   标签点全名.同一个标签点标识可以出现多次，但它们的时间戳必需是递增的
	[values] 数值
	[qualities]  质量码
	[times]  时间,UnixNano
输出:
说明:POST/GET	方法
编辑:wang_jp
时间:2020年6月6日
************************************************/
func (c *EngineController) GoldenPointSetSnapshot() {
	c.SaveUserActionMsg("写快照", _LOG_OPR_TYPE_INSERT) //记录用动作作信息
	type SnapWrite struct {                          //写庚顿快照数据结构
		TagName string  `json:"tagname"` //测点全名称
		Value   float64 `json:"value"`   //数值
		Time    string  `json:"time"`    //字符串形式的时间,例如:"2020-8-17 10:11:36",可以为空"",此时取当前时间
		Quality int     `json:"quality"` //质量码(GOOD = 0,NODATA = 1,CREATED = 2,SHUTDOWN = 3,CALCOFF = 4,BAD = 5,DIVBYZERO = 6,REMOVED = 7,OPC = 256,USER = 512)
	}
	var writers []SnapWrite
	datas := c.GetString("datas")
	err := json.Unmarshal([]byte(datas), &writers)
	if err != nil {
		c.Ctx.ResponseWriter.WriteHeader(_PAROMATER_ERR)
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
			c.Ctx.ResponseWriter.WriteHeader(_SERVER_ERR)
			c.Data["json"] = err.Error()
		} else {
			c.Data["json"] = "ok"
		}
	}
	c.ServeJSON()
}

/***********************************************
功能:写庚顿历史数据
输入:无
输出:
说明:POST/GET	方法
编辑:wang_jp
时间:2020年8月17日
************************************************/
func (c *EngineController) GoldenPointSetHistory() {
	c.SaveUserActionMsg("写历史", _LOG_OPR_TYPE_INSERT) //记录用动作作信息
	type Write struct {                              //写庚顿快照数据结构
		TagName string  `json:"tagname"` //测点全名称
		Value   float64 `json:"value"`   //数值
		Time    string  `json:"time"`    //字符串形式的时间,例如:"2020-8-17 10:11:36"
		Quality int     `json:"quality"` //质量码(GOOD = 0,NODATA = 1,CREATED = 2,SHUTDOWN = 3,CALCOFF = 4,BAD = 5,DIVBYZERO = 6,REMOVED = 7,OPC = 256,USER = 512)
	}
	var writers []Write
	datas := c.GetString("datas")
	err := json.Unmarshal([]byte(datas), &writers)
	if err != nil {
		c.Ctx.ResponseWriter.WriteHeader(_PAROMATER_ERR)
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
		err := micgd.GoldenSetArchivedValues(tags, values, dqs, dtimes)
		if err != nil {
			c.Ctx.ResponseWriter.WriteHeader(_SERVER_ERR)
			c.Data["json"] = err.Error()
		} else {
			c.Data["json"] = "ok"
		}
	}
	c.ServeJSON()
}
