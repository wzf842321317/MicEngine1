package models

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

/*******************************************************************************
功能:读取所有有效的报表配置信息
输入:
	[readType] int,读取类型，0时仅读报表项,1时仅读文件夹项,其他值时读取全部
输出:读取到的信息行数,报表配置列表,错误信息
说明:
编辑:wang_jp
时间:2020年5月5日
*******************************************************************************/
func (rptl *CalcKpiReportList) GetRoportLists(readType int) (int64, []CalcKpiReportList, error) {
	o := orm.NewOrm()  //新建orm对象
	o.Using("default") //根据别名选择数据库
	var reports []CalcKpiReportList
	qt := o.QueryTable("CalcKpiReportList")
	switch readType {
	case 0: //仅读报表
		cond := orm.NewCondition()
		cond1 := cond.And("Folder", 0).And("Status", 1).AndCond(cond.And("DistributedId", 0).Or("DistributedId", EngineCfgMsg.CfgMsg.Id)) //分布式ID=0或者等于本机ID
		qt = qt.SetCond(cond1)
		break
	case 1: //仅读文件夹
		qt = qt.Filter("Folder__gt", 0)
		break
	default: //其他值,不限制
		break
	}
	num, err := qt.OrderBy("Seq").All(&reports)
	return num, reports, err
}

/*************************************************
功能:通过用户Id获取报表结构树
输入:用户Id
输出:报表结构树信息
说明:
编辑:wang_jp
时间:2020年5月8日
*************************************************/
func (u *SysUser) GetReportNodesByUserId(userid ...int64) ([]*CalcKpiReportList, error) {
	o := orm.NewOrm()
	o.Using("default") //选定数据库
	if len(userid) > 0 {
		u.Id = userid[0]
	}
	var nodes []*CalcKpiReportList
	qt := o.QueryTable("CalcKpiReportList")
	cond := orm.NewCondition()
	cond1 := cond.AndCond(cond.And("DistributedId", 0).Or("DistributedId", EngineCfgMsg.CfgMsg.Id)) //分布式ID=0或者等于本机ID
	qt = qt.SetCond(cond1).Filter("Permissions__Permission__Roles__Role__Users__User__Id", u.Id)
	qt = qt.Filter("Permissions__Permission__PermissionType", _PermissionTypeOfReport) //过滤权限类型
	if _, err := qt.Distinct().OrderBy("Pid", "Seq").All(&nodes); err == nil {
		return nodes, nil
	} else {
		return nil, err
	}
}

/*******************************************************************************
功能:读取授权的所有子节点
输入:[LevelCode] 层级码
	[userid] 用户ID
输出:读取到的信息行数,报表配置列表,错误信息
说明:
编辑:wang_jp
时间:2020年5月8日
*******************************************************************************/
func (rpt *CalcKpiReportList) GetChildNodes(userid int64) (int64, []CalcKpiReportList, error) {
	o := orm.NewOrm()  //新建orm对象
	o.Using("default") //根据别名选择数据库
	var reports []CalcKpiReportList
	cond := orm.NewCondition()
	cond1 := cond.AndCond(cond.And("DistributedId", 0).Or("DistributedId", EngineCfgMsg.CfgMsg.Id)) //分布式ID=0或者等于本机ID
	qt := o.QueryTable("CalcKpiReportList").SetCond(cond1).Filter("Permissions__Permission__Roles__Role__Users__User__Id", userid).Filter("LevelCode__istartswith", rpt.LevelCode)
	qt = qt.Filter("Permissions__Permission__PermissionType", _PermissionTypeOfReport) //过滤权限类型
	num, err := qt.Distinct().OrderBy("Id", "Seq").All(&reports)
	return num, reports, err
}

/*******************************************************************************
功能:保存最后计算时间
输入:无
输出:无
说明:
编辑:wang_jp
时间:2020年5月5日
*******************************************************************************/
func (rpt *CalcKpiReportList) SaveLastCalcTime() error {
	o := orm.NewOrm()  //新建orm对象
	o.Using("default") //根据别名选择数据库
	rpt.UpdateTime = TimeFormat(time.Now())
	_, err := o.Update(rpt, "LastCalcTime", "UpdateTime") //执行查询
	if err != nil {                                       //发生错误
		//logs.Alert(`Update CalcKpiReportList LastCalcTime error,the report id is %d:[%s]`, rpt.Id, err.Error())
		err = fmt.Errorf(`更新报表 LastCalcTime 时发生错误,报表ID=[%d],错误信息:[%s]`, rpt.Id, err.Error())
	}
	return err
}

/*******************************************************************************
功能:更新指定列或者全部更新
输入:可选的列名
输出:无
说明:如果不设置列名,则更新所有列。列名可以设置一个，也可以设置多个。
编辑:wang_jp
时间:2020年5月5日
*******************************************************************************/
func (rpt *CalcKpiReportList) Update(filedname ...string) error {
	o := orm.NewOrm()  //新建orm对象
	o.Using("default") //根据别名选择数据库
	rpt.UpdateTime = TimeFormat(time.Now())
	var err error
	if len(filedname) > 0 {
		for _, fname := range filedname {
			_, err = o.Update(rpt, fname) //更新指定列
		}
	} else {
		_, err = o.Update(rpt) //全部
	}
	if err != nil { //发生错误
		err = fmt.Errorf(`更新时发生错误,报表ID=[%d],错误信息:[%s]`, rpt.Id, err.Error())
	}
	return err
}

/*******************************************************************************
功能:获取本次计算的开始时间点和结束时间点，并判断是否可以开始计算
输入:无
输出:是否可以计算\计算开始时间点\结束时间点\错误信息
说明:
编辑:wang_jp
时间:2020年5月6日
*******************************************************************************/
func (rpt *CalcKpiReportList) PrevTime() (bool, string, string, error) {
	if rpt.ShiftHour == 0 || len(rpt.BaseTime) < 10 { //如果没有定义基准时间和班组时间
		if rpt.Workshop.Id > 0 { //查看是否定义了所属车间,如果定义了所属车间
			w := new(OreProcessDWorkshop)
			err := w.GetAttributeById(rpt.Workshop.Id) //从所属车间读取数据
			if err == nil {                            //如果正确读取到了数据
				rpt.Workshop = w //将数据保存到报表属性中
				rpt.ShiftHour = int64(rpt.Workshop.ShiftHour)
				rpt.BaseTime = rpt.Workshop.BaseTime
			}
		}
	}
	if rpt.ShiftHour == 0 || len(rpt.BaseTime) < 10 {
		return false, "", "", fmt.Errorf("报表没有设定基准时间和每班工作时间,ID=[ %d ],Name=[ %s ]", rpt.Id, rpt.Name)
	}
	if len(rpt.StartTime) < 10 { //如果没有设定起始时间
		rpt.StartTime = TimeFormat(time.Now()) //取当前时间
	}
	if len(rpt.LastCalcTime) < 10 { //如果没有设定最后计算时间
		rpt.LastCalcTime = rpt.StartTime //取开始时间
	}
	return prevTime(rpt.BaseTime, rpt.StartTime, rpt.LastCalcTime, rpt.ShiftHour, rpt.Period, rpt.OffsetMinutes)
}

/*******************************************************************************
功能:计算报表
输入:无
输出:
	[int]   进行计算的单元格数
	[string]报表名称
	[error] 错误信息
说明:
编辑:wang_jp
时间:2020年5月6日
*******************************************************************************/
func (rpt *CalcKpiReportList) ReportCalc(beginTime, endTime string) (int, string, error) {
	var calc_cell int = 0 //计算的单元格数

	dt := strings.ReplaceAll(endTime, " ", "")                            //去掉空格
	dt = strings.ReplaceAll(dt, ":", "")                                  //去掉":"
	dt = strings.ReplaceAll(dt, "-", "")                                  //去掉"-"
	reportname := fmt.Sprintf("%s-%d_%s.xlsx", _ReportPreFix, rpt.Id, dt) //生成报表名称

	tplpath := rpt.TemplateUrl
	if tplpath[len(tplpath)-1:] != "/" { //如果最后一个字符不是"/"
		tplpath += "/"
	}
	tplpath += rpt.TemplateFile //模板文件路径
	if len(rpt.TemplateFile) < 5 {
		return calc_cell, reportname, fmt.Errorf("报表[%s]没有上传有效模板", rpt.Name)
	}

	tplf, err := excelize.OpenFile(tplpath) //读取模板
	if err != nil {                         //读取模板错误
		return calc_cell, reportname, fmt.Errorf("读取报表模板信息失败,报表名称:[%s],错误信息:[%s]", reportname, err.Error())
	}
	// 获取 报表Sheet 上所有单元格
	rows, err := tplf.GetRows(_ReportSheet)
	if err != nil { //读取Sheet1错误
		return calc_cell, reportname, fmt.Errorf("读取报表模板[%s]信息失败,报表名称:[%s],错误信息:[%s]", _ReportSheet, reportname, err.Error())
	}

	reg := regexp.MustCompile(`\{\{[^\{\}]+\}\}`) //正则提取大括号内的值，含大括号
	rowCnt := len(rows)                           //总行数
	for r, row := range rows {                    //遍历所有行
		for c, colCell := range row { //遍历行中的所有列
			if reg.MatchString(colCell) { //含有脚本
				cellname, e := excelize.CoordinatesToCellName(c+1, r+1) //计算单元格坐标名
				if e == nil {
					strs := reg.FindAllString(colCell, -1)       //提取脚本
					indxs := reg.FindAllStringIndex(colCell, -1) //提取位置
					ldif := 0
					for i, scr := range strs {
						script := new(Script) //新建脚本结构
						script.Id = rpt.Id
						script.BaseTime = rpt.BaseTime
						script.BeginTime = beginTime
						script.EndTime = endTime
						script.ShiftHour = rpt.ShiftHour
						script.ScriptStr = scr[2 : len(scr)-2] //去掉两头的大括号,提取脚本

						calc_cell += 1              //计算脚本的数量
						res, _, err := script.Run() //运行脚本
						resstr := fmt.Sprint(res)
						if err != nil { //脚本有错误
							estr := fmt.Sprintf("报表[%s]单元格[%s]脚本[%s]错误,错误信息:[%s]", reportname, cellname, script.ScriptStr, err.Error())
							if rpt.Debug > 0 { //调试模式
								rowCnt += 1                                              //最后一行的行号
								emsgcell, e := excelize.CoordinatesToCellName(1, rowCnt) //最后一行第一列的坐标
								if e == nil {
									tplf.SetCellValue(_ReportSheet, emsgcell, estr) //将错误信息保存到最后一行第一列
								}
							} else { //运行模式
								logs.Warn(estr)                                                            //将错误记录到日志中
								resstr = "0"                                                               //空白字符替换公式
								colCell = colCell[:indxs[i][0]-ldif] + resstr + colCell[indxs[i][1]-ldif:] //用计算结果替换掉单元格中的相应脚本
								ldif += len(scr) - len(resstr)
							}
						} else { //无错误
							colCell = colCell[:indxs[i][0]-ldif] + resstr + colCell[indxs[i][1]-ldif:] //用计算结果替换掉单元格中的相应脚本
							ldif += len(scr) - len(resstr)                                             //脚本与计算结果长度的差
						}
					}
					fval, err := strconv.ParseFloat(colCell, 64) //转换成浮点数
					if err == nil {                              //可以转换
						er := tplf.SetCellValue(_ReportSheet, cellname, fval) //以浮点数的形式写入文件
						if er != nil {
							return calc_cell, reportname, fmt.Errorf("设置单元格[%s]的值失败,报表名称:[%s],错误信息:[%s]", cellname, reportname, er.Error())
						}
					} else { //不可以转换
						er := tplf.SetCellValue(_ReportSheet, cellname, colCell) //以文本的形式写入文件
						if er != nil {
							return calc_cell, reportname, fmt.Errorf("设置单元格[%s]的值失败,报表名称:[%s],错误信息:[%s]", cellname, reportname, er.Error())
						}
					}
				} //单元格别名
			} //含有脚本
		} //遍历列
	} //遍历行

	rpt.LastCalcTime = endTime   //最后计算时间
	err = rpt.SaveLastCalcTime() //保存最后计算时间
	if err == nil {              //保存时间成功后将报表另存为
		filepath := rpt.ResultUrl
		if filepath[len(filepath)-1:] != "/" { //如果最后一个字符不是"/"
			filepath += "/"
		}
		filepath += fmt.Sprintf("%d/", rpt.Id)
		direrr := MakeDir(filepath) //检查路径是否存在,不存在就创建
		if direrr != nil {          //创建路径错误
			return calc_cell, reportname, fmt.Errorf("报表[%s]的存储路径[%s]不存在且创建失败,错误信息:[%s]", reportname, filepath, direrr.Error())
		}
		filepath += reportname
		err = tplf.SaveAs(filepath) //报表另存为
		if err != nil {
			return calc_cell, "", fmt.Errorf("保存报表失败,报表名称:[%s],错误信息:[%s]", reportname, err.Error())
		}
	} else {
		return calc_cell, reportname, fmt.Errorf("保存报表最后计算时间失败,报表名称:[%s],错误信息:[%s]", reportname, err.Error())
	}
	return calc_cell, reportname, err
}

/*******************************************************************************
功能:通过ID获取报表的属性信息
输入:[ID] 报表ID
输出:读取到的信息行数,报表配置列表,错误信息
说明:
编辑:wang_jp
时间:2020年5月8日
*******************************************************************************/
func (rpt *CalcKpiReportList) GetAttributeById() error {
	o := orm.NewOrm()  //新建orm对象
	o.Using("default") //根据别名选择数据库
	qt := o.QueryTable("CalcKpiReportList").Filter("Id", rpt.Id)
	_, err := qt.All(rpt)
	if rpt.Workshop != nil {
		if rpt.Workshop.Id > 0 {
			w := new(OreProcessDWorkshop)
			e := w.GetAttributeById(rpt.Workshop.Id)
			if e == nil {
				rpt.Workshop = w
			}
		}
	}
	return err
}

/*******************************************************************************
功能:读取该报表的所有模板
输入:ID
输出:读取到的信息行数,报表配置列表,错误信息
说明:
编辑:wang_jp
时间:2020年5月8日
*******************************************************************************/
func (rpt *CalcKpiReportList) GetTemplatesListById() (interface{}, error) {
	err := rpt.GetAttributeById() //通过ID读取属性
	if err != nil {
		return nil, err
	}

	type tplstruct struct { //模板列表结构
		FileName  string //文件名
		FileTime  string //文件上传时间
		ModTime   string //文件编辑时间
		Size      int64  //单位:字节
		IsCurrent bool   //是否当前所选模板
		Path      string //路径
	}
	var tpls []tplstruct

	if rpt.Folder == 0 {
		files, err := ioutil.ReadDir(rpt.TemplateUrl) //读取模板路径下的文件
		if err != nil {
			return nil, err
		}

		for _, file := range files {
			if file.IsDir() { //如果是文件夹,跳过
				continue
			} else { //如果是文件
				names := strings.Split(file.Name(), "_") //用下划线分割文件名
				if len(names) > 1 {
					if names[0] != fmt.Sprintf("tpl-%d", rpt.Id) { //如果文件名第一部分不等于层级码,跳过
						continue
					} else { //等于层级码
						tms := strings.Split(names[1], ".") //用点分割时间和后缀名部分
						if len(tms) > 0 {
							tpl := new(tplstruct)
							t, e := TimeParse(tms[0])
							if e == nil {
								tpl.FileTime = TimeFormat(t) //文件上传时间
							}
							tpl.FileName = file.Name() //文件名
							tpl.ModTime = TimeFormat(file.ModTime())
							tpl.Size = file.Size()
							tpl.IsCurrent = strings.Contains(rpt.TemplateFile, file.Name())
							tpl.Path = rpt.TemplateUrl
							tpls = append(tpls, *tpl)
						}
					}
				}
			}
		}
	}
	return tpls, nil
}

/*******************************************************************************
功能:读取该报表的所有结果
输入:
	[id] 报表定义的ID
	[timerange] 可选的时间串,0个、1个或2个
		0个时,不限制时间范围
		1个时,限制开始时间，不限制结束时间
		2个时,第一个为开始时间，第二个为结束时间
输出:报表结果列表,错误信息
说明:
编辑:wang_jp
时间:2020年5月11日
*******************************************************************************/
func (rpt *CalcKpiReportList) GetResultListById(timerange ...string) (interface{}, error) {
	err := rpt.GetAttributeById() //通过ID读取属性
	if err != nil {
		return nil, err
	}

	type tplstruct struct { //模板列表结构
		FileName string //文件名
		FileTime string //文件上传时间
		ModTime  string //文件编辑时间
		Size     int64  //单位:字节
	}
	var tpls []tplstruct
	haverange := len(timerange) //设置了时间范围
	var bgt, edt time.Time
	if haverange > 0 {
		bgt, err = TimeParse(timerange[0]) //起始时间
		if err != nil {
			haverange = 0
		}
		if haverange > 1 {
			edt, err = TimeParse(timerange[1]) //结束时间
			if err != nil {
				haverange = 0
			}
		}
	}

	if rpt.Folder == 0 {
		filepath := rpt.ResultUrl
		if filepath[len(filepath)-1:] != "/" { //如果最后一个字符不是"/"
			filepath += "/"
		}
		filepath += fmt.Sprintf("%d/", rpt.Id)
		direrr := MakeDir(filepath) //检查路径是否存在,不存在就创建
		if direrr == nil {          //创建路径没有错误
		}
		files, err := ioutil.ReadDir(filepath) //读取模板路径下的文件
		if err != nil {
			return nil, err
		}

		for _, file := range files {
			if file.IsDir() { //如果是文件夹,跳过
				continue
			} else { //如果是文件
				names := strings.Split(file.Name(), "_") //用下划线分割文件名
				if len(names) > 1 {
					if names[0] != fmt.Sprintf("%s-%d", _ReportPreFix, rpt.Id) { //如果文件名第一部分不等于文件名,跳过
						continue
					} else { //等于
						tms := strings.Split(names[1], ".") //用点分割时间和后缀名部分
						if len(tms) > 0 {
							tpl := new(tplstruct)
							t, e := TimeParse(tms[0])
							if e == nil {
								tpl.FileTime = TimeFormat(t) //文件上传时间
							}
							tpl.FileName = file.Name() //文件名
							tpl.ModTime = TimeFormat(file.ModTime())
							tpl.Size = file.Size()
							if haverange > 0 { //设置了时间范围
								if haverange > 1 {
									if t.Unix() > bgt.Unix() && t.Unix() <= edt.Unix() { //设置了起始时间和结束时间
										tpls = append(tpls, *tpl)
									}
								} else { //只设置了起始时间
									if t.Unix() > bgt.Unix() {
										tpls = append(tpls, *tpl)
									}
								}
							} else { //没有设置时间范围
								tpls = append(tpls, *tpl)
							}
						}
					}
				}
			}
		}
	}
	return tpls, nil
}

/*******************************************************************************
功能:新建节点
输入:无
输出:无
说明:
编辑:wang_jp
时间:2020年5月9日
*******************************************************************************/
func (rpt *CalcKpiReportList) InsertNode() error {
	o := orm.NewOrm()  //新建orm对象
	o.Using("default") //根据别名选择数据库
	rpt.UpdateTime = TimeFormat(time.Now())
	rpt.CreateTime = rpt.UpdateTime

	//获取父层级信息
	prpt := new(CalcKpiReportList)
	if rpt.Pid > 0 { //有父层级的时候才执行查询层级深度
		prpt.Id = rpt.Pid
		err1 := prpt.GetAttributeById()
		if err1 == nil {
			rpt.Level = prpt.Level + 1 //层级深度
		} else {
			rpt.Level = 1
		}
	} else { //没有父层级,深度为1
		rpt.Level = 1
	}

	id, err := o.Insert(rpt) //执行查询
	if err != nil {          //发生错误
		return fmt.Errorf(`新建报表层级失败,拟建报表层级信息:[%+v],错误信息:[%s]`, rpt, err.Error())
	}
	rpt.Id = id
	if rpt.Pid > 0 && prpt.Id > 0 { //设置了父层级,且父层级存在
		rpt.LevelCode = fmt.Sprintf("%s-%d", prpt.LevelCode, rpt.Id) //合成层级码
	} else {
		rpt.LevelCode = fmt.Sprint(rpt.Id)
	}
	err = rpt.Update("LevelCode") //更新层级码
	if err != nil {               //发生错误
		return fmt.Errorf(`更新层级码失败:[%+v],错误信息:[%s]`, rpt, err.Error())
	}
	_, pms, err := prpt.LoadPermissions(rpt.CreateUserId) //获取父节点与建立新节点的用户之间的权限列表
	if err == nil {
		m2m := o.QueryM2M(rpt, "Permissions") //添加权限关系
		_, e := m2m.Add(pms)
		if e == nil {
			return e
		}
	}
	return nil
}

/*******************************************************************************
功能:通过ID和用户ID获取报表与用户关联的权限列表
输入:[userid] 用户ID
输出:数据行数,权限列表,错误信息
说明:
编辑:wang_jp
时间:2020年5月10日
*******************************************************************************/
func (rpt *CalcKpiReportList) LoadPermissions(userid int64) (int64, []SysPermission, error) {
	o := orm.NewOrm()  //新建orm对象
	o.Using("default") //根据别名选择数据库
	var pms []SysPermission
	qt := o.QueryTable("SysPermission").Filter("Roles__Role__Users__User__Id", userid)
	qt = qt.Filter("PermissionType", _PermissionTypeOfReport) //过滤权限类型
	if rpt.Id > 0 {
		qt = qt.Filter("Reports__Report__Id", rpt.Id)
	}
	num, err := qt.Distinct().All(&pms)
	return num, pms, err
}

/*******************************************************************************
功能:删除‘；7
输入:[userid] 用户ID
输出:错误信息
说明:
编辑:wang_jp
时间:2020年5月10日
*******************************************************************************/
func (rpt *CalcKpiReportList) DeleteNode() error {
	o := orm.NewOrm()  //新建orm对象
	o.Using("default") //根据别名选择数据库
	m2m := o.QueryM2M(rpt, "Permissions")
	_, err := m2m.Clear()
	if err != nil {
		return err
	}
	_, err = o.Delete(rpt)
	return err
}
