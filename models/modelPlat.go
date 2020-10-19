package models

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"time"

	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/bkzy-wangjp/MicEngine/MicScript/statistic"
	_ "github.com/go-sql-driver/mysql"
)

/*
功能:检查数据库中是否存在指定表
输入:dbalias:数据库别名
	tablename:要检查的数据表名称
输出:bool:如果存在true,不存在false
	error:错误信息
说明:有错误信息的时候查询结果为false
编辑:wang_jp
时间:2019年12月13日
*/
func isTableExist(dbalias, tablename string) (bool, error) {
	o := orm.NewOrm() //新建orm对象
	o.Using(dbalias)  //根据别名选择数据库
	var res []orm.Params
	sql := fmt.Sprintf("SHOW TABLES LIKE '%s'", tablename)
	num, err := o.Raw(sql).Values(&res) //执行查询
	if err != nil {                     //发生错误
		return false, errors.New(fmt.Sprintf(`SHOW TABLES LIKE "%s" Error:[%s]`, tablename, err.Error()))
	}
	return num > 0, nil
}

/*
功能:检查数据库中是否存在指定表,如果不存在，则创建
输入:dbalias:数据库别名
	tablename:要检查的数据表名称
输出:bool:如果创建成功或者已经存在,true,没有创建成功false
	error:错误信息
说明:有错误信息的时候返回结果为false
编辑:wang_jp
时间:2019年12月13日
*/
/*
func createSysLogTable(dbalias, tablename string) (bool, error) {
	if exist, err := isTableExist(dbalias, tablename); err != nil { //检查表是否存在
		return false, err //发生错误
	} else if exist == true {
		return true, nil //表已经存在
	}

	o := orm.NewOrm() //新建orm对象
	o.Using(dbalias)  //根据别名选择数据库
	var res []orm.Params
	sql := fmt.Sprintf(`CREATE TABLE %s  (
  id bigint(20) NOT NULL AUTO_INCREMENT COMMENT '自动ID',
  user_id bigint(20) NOT NULL DEFAULT 0 COMMENT '用户ID',
  user_name varchar(50) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT '' COMMENT '用户名',
  action varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT '' COMMENT '动作',
  message longtext CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL COMMENT '详细信息',
  log_time datetime(0) NULL DEFAULT NULL COMMENT '时间',
  ip varchar(50) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT '' COMMENT '登录IP',
  PRIMARY KEY (id) USING BTREE
) ENGINE = MyISAM AUTO_INCREMENT = 1 CHARACTER SET = utf8 COLLATE = utf8_general_ci ROW_FORMAT = Dynamic;`, tablename)
	_, err := o.Raw(sql).Values(&res) //执行查询
	if err != nil {                   //发生错误
		return false, errors.New(fmt.Sprintf(`Create table "%s" Error:[%s]`, tablename, err.Error()))
	}
	return true, nil //建表成功
}
*/

/*
功能:设定Mysql的最大预处理数量
输入:数据库别名,最大数量
输出:无
说明:
编辑:wang_jp
时间:2019年12月14日
*/
func setMysqlMaxPreparedStmtCount(dbalias string, maxcount int) {
	o := orm.NewOrm() //新建orm对象
	o.Using(dbalias)  //根据别名选择数据库
	sql := fmt.Sprintf(`set global max_prepared_stmt_count=%d;`, maxcount)
	_, err := o.Raw(sql).Exec() //执行查询
	if err != nil {             //发生错误
		logs.Info("setMysqlMaxPreparedStmtCount error:[%s]", err.Error())
	}
}

/*********************************************
功能:时间字符串转换为时间对象
输入:时间字符串,可选的时区
输出:格式化后的时间变量,错误信息
时间:2019年11月28日
编辑:wang_jp
*********************************************/
func TimeParse(s string, loc ...*time.Location) (time.Time, error) {
	if len(loc) > 0 {
		return statistic.TimeParse(s, loc[0])
	} else {
		return statistic.TimeParse(s)
	}
}

/*********************************************
功能：时间对象格式化为时间字符串
输入：时间字符串,可选的时区
输出：格式化后的时间变量,错误信息
时间：2019年11月28日
编辑：wang_jp
*********************************************/
func TimeFormat(t time.Time, layout ...string) string {
	if len(layout) > 0 {
		return t.Format(layout[0])
	} else {
		return t.Format(EngineCfgMsg.Sys.TimeFormat)
	}
}

/*********************************************
功能:获取指定的时间所在月的起始时间
输入:指定的时间
输出:指定时间所在月的起始时间
时间:2020年7月15日
编辑:wang_jp
*********************************************/
func TimeMonthStart(t time.Time) time.Time {
	str := t.Format("2006-01") + "-01 00:00:00"
	t2, _ := TimeParse(str)
	return t2
}

/*********************************************
功能:根据基准时间合成当前输入时间所在的起止时间点
输入:基准时间\最后时间标签\每班的小时数\结果类型
	结果类型：
	-- 0:返回值是lastTime所在工作年的开始时间
    -- 1:返回值是lastTime所在工作季度的开始时间
    -- 2:返回值是lastTime所在工作月的开始时间
    -- 3:返回值是lastTime中日期部分所在工作日的开始时间
    -- 4:返回值是lastTime中时间部分所在工作日的开始时间
    -- 5或其他:返回值是lastTime所在工作班的开始时间
输出:时间点\错误信息
编辑:wang_jp
时间:2020年3月1日
*********************************************/
func ComposeTime(baseTime, lastTime string, shifthour, resultType int64) (string, error) {
	base, err := TimeParse(baseTime) //基准时间转换
	if err != nil {
		//return "", fmt.Errorf("Time formate error when [baseTime] parse in ComposeTime(),the value is [%s],the wangt is %s;[%s]", baseTime, EngineCfgMsg.Sys.TimeFormat, err.Error())
		return "", fmt.Errorf("[baseTime] 时间格式错误,无法在 ComposeTime() 函数中完成转换,当前值是: [%s],期望的格式是:[%s];错误信息:[%s]", baseTime, EngineCfgMsg.Sys.TimeFormat, err.Error())
	}
	last, err := TimeParse(lastTime) //最后计算时间标签转换
	if err != nil {                  //如果转换失败或者last更靠前
		//return "", fmt.Errorf("Time formate error when [lastTime] parse in ComposeTime(),the value is [%s],the wangt is %s;[%s]", lastTime, EngineCfgMsg.Sys.TimeFormat, err.Error())
		return "", fmt.Errorf("[lastTime] 时间格式错误,无法在 ComposeTime(),函数中完成转换,当前值是: [%s],期望的格式是:[%s];错误信息:[%s]", lastTime, EngineCfgMsg.Sys.TimeFormat, err.Error())
	}

	_, bm, bd := base.Date()        //获取基准月、日
	ly, lm, ld := last.Date()       //获取最后标记的年
	blocal := time.Now().Location() //获取基准时区
	var restime time.Time
	switch resultType { //根据不同的周期设定,计算开始计算时间和结束计算时间
	case 0: //返回值是lastTime所在工作年的开始时间
		restime = time.Date(ly, bm, bd, base.Hour(), base.Minute(), base.Second(), 0, blocal)
		if restime.After(last) { //第一个边界不能大于最后标记
			restime = restime.AddDate(0, -12, 0)
		}
	case 1: //返回值是lastTime所在工作季度的开始时间
		var sesion []time.Time //建立季度边界数组
		sesion = append(sesion, time.Date(ly, bm, bd, base.Hour(), base.Minute(), base.Second(), 0, blocal))
		if sesion[0].After(last) { //第一个边界不能大于最后标记
			sesion[0] = sesion[0].AddDate(0, -12, 0)
		}
		sesion = append(sesion, sesion[0].AddDate(0, 3, 0))
		sesion = append(sesion, sesion[1].AddDate(0, 3, 0))
		sesion = append(sesion, sesion[2].AddDate(0, 3, 0))
		sesion = append(sesion, sesion[3].AddDate(0, 3, 0))

		for i, s := range sesion { //遍历季度边界
			if i < len(sesion)-1 { //只对前四个边界进行比较
				if last == s { //如果最后标记与季度边界相等
					restime = s //开始时间为本边界时间
					break       //已经找到结果，结束遍历季节
				} else { //如果不与边界相等
					if last.Before(sesion[i+1]) { //判断是否处于当前边界与下一个边界之间
						restime = s //如果是,则开始时间是本季节的边界
						break       //已经找到结果，结束遍历季节
					}
				}
			}
		}
	case 2: //返回值是lastTime所在工作月的开始时间
		restime = time.Date(ly, lm, bd, base.Hour(), base.Minute(), base.Second(), 0, blocal)
		if restime.After(last) { //第一个边界不能大于最后标记
			restime = restime.AddDate(0, -1, 0)
		}
	case 3: //返回值是lastTime中日期部分所在工作日的开始时间
		restime = time.Date(ly, lm, ld, base.Hour(), base.Minute(), base.Second(), 0, blocal)
	case 4: //返回值是lastTime中时间部分所在工作日的开始时间
		period := time.Hour * 24
		pcnt := last.Add(period).Sub(base).Nanoseconds() / period.Nanoseconds() //计算从基准时间到最近时间+1个周期的总周期数
		end := base.Add(time.Duration(pcnt * period.Nanoseconds()))
		restime = end.Add(period * -1)
	default: //返回值是lastTime所在工作班的开始时间
		period := time.Hour * time.Duration(shifthour)
		pcnt := last.Add(period).Sub(base).Nanoseconds() / period.Nanoseconds() //计算从基准时间到最近时间+1个周期的总周期数
		end := base.Add(time.Duration(pcnt * period.Nanoseconds()))
		restime = end.Add(period * -1)
	}

	return restime.Format("2006-01-02 15:04:05"), nil
}

/*
功能:获取本次计算的开始时间点和结束时间点，并判断是否可以开始计算
输入:基准时间\设定的计算起始时间\最后时间标签\每班的小时数\周期设定\时间偏移量
输出:是否可以计算\计算开始时间点\结束时间点\错误信息
说明:
编辑:wang_jp
时间:2019年12月14日
*/
func prevTime(baseTime, startTime, lastTime string, shifthour, periodset, offsetminute int64) (bool, string, string, error) {
	base, err := time.ParseInLocation(EngineCfgMsg.Sys.TimeFormat, baseTime, time.Local) //基准时间转换
	if err != nil {
		//return false, "", "", fmt.Errorf("Time format error when [baseTime] parse,the value is [%s],the wangt is %s;[%s]", baseTime, EngineCfgMsg.Sys.TimeFormat, err.Error())
		return false, "", "", fmt.Errorf("基准时间格式错误,设置的值是 [%s],期望的值是 [%s];错误信息:[%s]", baseTime, EngineCfgMsg.Sys.TimeFormat, err.Error())
	}
	start, err := time.ParseInLocation(EngineCfgMsg.Sys.TimeFormat, startTime, time.Local) //开始时间转换
	if err != nil {
		return false, "", "", fmt.Errorf("计算开始时间格式错误,设置的值是 [%s],期望的值是 [%s];错误信息:[%s]", startTime, EngineCfgMsg.Sys.TimeFormat, err.Error())
	}
	if start.After(time.Now()) { //还没到开始时间
		return false, "", "", nil //返回
	}
	last, err := time.ParseInLocation(EngineCfgMsg.Sys.TimeFormat, lastTime, time.Local) //最后计算时间标签转换
	if err != nil || last.Before(start) {                                                //如果转换失败或者last更靠前
		last = start //设定为开始时间
	}

	delay := time.Duration(EngineCfgMsg.CfgMsg.SerialCalcDelaySec) * time.Second //延迟计算时间转换
	offset := time.Duration(offsetminute) * time.Minute                          //时间偏移转换
	if offsetminute < 0 {                                                        //小于零时,beginTime和endTime要做迁移
		base = base.Add(offset) //设置时间偏移
	}
	if offsetminute > 0 { //大于零时,beginTime和endTime不做迁移,但是启动计算的时间要迁移
		delay += offset
	}

	var period time.Duration //周期
	var begine, end time.Time
	var timeDone bool = false //用于返回的开始计算时间和结束计算时间
	switch periodset {        //根据不同的周期设定,计算开始计算时间和结束计算时间
	case -1: //小时
		period = time.Hour
		begine, end = prevTimeInDay(base, last, period)
		if end.Add(delay).Before(time.Now()) { //结束时间（+延迟）必须在当前时间之前才可以计算
			timeDone = true
		}
	case -2: //班
		period = time.Hour * time.Duration(shifthour)
		begine, end = prevTimeInDay(base, last, period)
		if end.Add(delay).Before(time.Now()) { //结束时间（+延迟）必须在当前时间之前才可以计算
			timeDone = true
		}
	case -3: //日
		period = time.Hour * 24
		begine, end = prevTimeInDay(base, last, period)
		if end.Add(delay).Before(time.Now()) { //结束时间（+延迟）必须在当前时间之前才可以计算
			timeDone = true
		}
	case -4, -5, -6: //月
		begine, end = prevTimeOutDay(base, last, periodset)
		if end.Add(delay).Before(time.Now()) { //结束时间（+延迟）必须在当前时间之前才可以计算
			timeDone = true
		}
	case 0:
		//return false, "", "", fmt.Errorf("Period error,the value is [%d]", period)
		return false, "", "", fmt.Errorf("周期设置错误,设定的周期值是: [%d]", period)
	default:
		period = time.Second * time.Duration(periodset)
		begine, end = prevTimeInDay(base, last, period)
		if end.Before(time.Now()) { //结束时间（+延迟）必须在当前时间之前才可以计算
			timeDone = true
		}
	}

	return timeDone, begine.Format(EngineCfgMsg.Sys.TimeFormat), end.Format(EngineCfgMsg.Sys.TimeFormat), nil
}

/*
功能:计算周期长度在一天之内的周期的开始时间和结束时间
输入:基准时间\最后时间标签\周期
输出:计算开始时间点和结束时间点
说明:
编辑:wang_jp
时间:2019年12月14日
*/
func prevTimeInDay(base, last time.Time, period time.Duration) (time.Time, time.Time) {
	var begine, end time.Time
	var pcnt int64
	pcnt = last.Add(period).Sub(base).Nanoseconds() / period.Nanoseconds() //计算从基准时间到最近时间+1个周期的总周期数
	end = base.Add(time.Duration(pcnt * period.Nanoseconds()))
	begine = end.Add(period * -1)
	if begine.UnixNano() > end.UnixNano() {
		begine = begine.Add(period * -1)
	}
	return begine, end
}

/*
功能:计算周期为月、季度、年的周期的开始时间和结束时间
输入:基准时间\最后时间标签\周期
输出:计算开始时间点和结束时间点
说明:
编辑:wang_jp
时间:2019年12月14日
*/
func prevTimeOutDay(base, last time.Time, period int64) (time.Time, time.Time) {
	var begine, end time.Time
	switch period {
	case -4: //月
		_, _, bd := base.Date()
		ly, lm, _ := last.Date()
		blocal := time.Now().Location()
		begine = time.Date(ly, lm, bd, base.Hour(), base.Minute(), base.Second(), 0, blocal)
		if begine.After(last) {
			begine = begine.AddDate(0, -1, 0)
		}
		end = begine.AddDate(0, 1, 0)
	case -5: //季度
		_, bm, bd := base.Date()        //获取基准月、日
		blocal := time.Now().Location() //获取基准时区
		ly, _, _ := last.Date()         //获取最后标记的年
		var sesion []time.Time          //建立季度边界数组
		sesion = append(sesion, time.Date(ly, bm, bd, base.Hour(), base.Minute(), base.Second(), 0, blocal))
		if sesion[0].After(last) { //第一个边界不能大于最后标记
			sesion[0] = sesion[0].AddDate(0, -12, 0)
		}
		sesion = append(sesion, sesion[0].AddDate(0, 3, 0))
		sesion = append(sesion, sesion[1].AddDate(0, 3, 0))
		sesion = append(sesion, sesion[2].AddDate(0, 3, 0))
		sesion = append(sesion, sesion[3].AddDate(0, 3, 0))

		for i, s := range sesion { //遍历季度边界
			if i < len(sesion)-1 { //只对前四个边界进行比较
				if last == s { //如果最后标记与季度边界相等
					begine = s               //开始时间为本边界时间
					end = s.AddDate(0, 3, 0) //结束时间为下一个边界
					break                    //已经找到结果，结束遍历季节
				} else { //如果不与边界相等
					if last.Before(sesion[i+1]) { //判断是否处于当前边界与下一个边界之间
						begine = s               //如果是,则开始时间是本季节的边界
						end = s.AddDate(0, 3, 0) //结束时间是下一个节边界
						break                    //已经找到结果，结束遍历季节
					}
				}
			}
		}
	case -6: //年
		_, bm, bd := base.Date()        //获取基准月、日
		blocal := time.Now().Location() //获取基准时区
		ly, _, _ := last.Date()         //获取最后标记的年
		begine = time.Date(ly, bm, bd, base.Hour(), base.Minute(), base.Second(), 0, blocal)
		if begine.After(last) { //第一个边界不能大于最后标记
			begine = begine.AddDate(0, -12, 0)
		}
		end = begine.AddDate(0, 12, 0)
	}
	return begine, end
}

/*
功能:每分钟保存心跳信号
输入:无
输出:无
说明:计时器由调用者提供
编辑:wang_jp
时间:2020年2月13日
*/

func (cfg *CalcKpiEngineConfig) SaveHeartBeatPerMinute() {
	o := orm.NewOrm() //新建orm对象
	o.Using("default")
	st, err := time.ParseInLocation(_TIMEFOMAT, EngineCfgMsg.StartTime, time.Local) //启动时间
	if err != nil {
		st = time.Now()
	}
	EngineCfgMsg.CfgMsg.RunMinutes = int64(time.Since(st).Minutes())
	EngineCfgMsg.CfgMsg.TotalMinutes = EngineCfgMsg.totalRunTimeWhenStart + EngineCfgMsg.CfgMsg.RunMinutes
	EngineCfgMsg.CfgMsg.UpdateTime = time.Now().Format(_TIMEFOMAT)
	if EngineCfgMsg.Sys.PeriodKpiEnable || EngineCfgMsg.Sys.ReportEnable {
		_, err = o.Update(&EngineCfgMsg.CfgMsg, "RunMinutes", "TotalMinutes", "Updatetime")
		if err != nil {
			logs.Warning("Update Engine runtime error[更新计算引擎运行时间失败]" )
			logs.Warning(err)
		}
	}
}

/*************************************************
功能:获取部门信息
输入:部门ID或者省略
输出:[]MineDeptInfo, error
说明:
编辑:wang_jp
时间:2020年4月15日
*************************************************/
func (dep *MineDeptInfo) GetMineDeptLists(ids ...int64) ([]*MineDeptInfo, error) {
	o := orm.NewOrm()
	o.Using("default") //选定数据库
	var depts []*MineDeptInfo
	qt := o.QueryTable("MineDeptInfo")
	if len(ids) == 0 {
		_, err := qt.Filter("Id__gt", 0).RelatedSel().OrderBy("Id").All(&depts) //
		return depts, err
	} else {
		var dpts []*MineDeptInfo
		for _, id := range ids {
			_, err := qt.Filter("Id", id).RelatedSel().All(&dpts) //
			if err == nil {
				depts = append(depts, dpts...)
			}
		}
		return depts, nil
	}
}

/**************************************************
功能:通过表ID获取表名
输入:表名
输出:表名和错误信息
说明:如果有错误，返回值为0
编辑:wang_jp
时间:2019年12月9日
**************************************************/
func (tb *MineTableList) GetTableNameByTableID(table_id int64) (string, error) {
	o := orm.NewOrm() //新建orm对象
	o.Using("default")
	tableInfo := MineTableList{Id: table_id}
	if err := o.Read(&tableInfo); err != nil {
		return "", err
	} else {
		return tableInfo.TableProgramName, nil
	}
}

/**************************************************
功能:通过表名、ID和列名查询值
输入:表名、ID和列名
输出:值和错误信息
说明:如果有错误，返回值为空
编辑:wang_jp
时间:2019年12月9日
**************************************************/
func GetTableFieldValueByName(tablename string, itemid int64, field_name string) (string, error) {
	o := orm.NewOrm() //新建orm对象
	o.Using("default")
	var lists []orm.Params
	sql := fmt.Sprintf("SELECT id,%s FROM %s WHERE id = %d", field_name, tablename, itemid)
	num, err := o.Raw(sql).Values(&lists)
	if err != nil { //错误的时候返回错误信息
		return "", err
	} else if num == 0 {
		return "", errors.New(fmt.Sprintf("%s no data return", sql))
	} else {
		return fmt.Sprint(lists[0][field_name]), nil

	}
}

/**************************************************
功能:通过表名或者表ID, 行号,列名 获取指定表指定行指定列的值
输入:
	[tablenameorid] 表名或者表ID
	[item_id]		行号
	[field_name]	列名
输出:float64结果和error
说明:field_name所在列必须可转换为数字，否则输出0和错误信息
编辑:wang_jp
时间:2019年12月9日
**************************************************/
func GetTableRowFieldValueByID(tablenameorid string, item_id int64, field_name string) (string, error) {
	tablename := tablenameorid //如果不可以转换为数字,默认category是tablename
	table_id, err := strconv.ParseInt(tablenameorid, 0, 64)
	if err == nil { //如果可以转换为数字,那么category是table_id,需要通过id获取tablename
		tb := new(MineTableList)
		if tablename, err = tb.GetTableNameByTableID(table_id); err != nil {
			//return "", fmt.Errorf("Can't get table name by table id %d,the error massage is:%s", table_id, err.Error())
			return "", fmt.Errorf("无法通过指定的表ID[%d]获取到表名称,错误信息:[%s]", table_id, err.Error())
		}
	}
	return GetTableFieldValueByName(tablename, item_id, field_name)
}

/*******************************************************************************
功能:判断文件、文件夹是否存在
输入:文件路径
输出:
	[bool]   存在就返回为true，不存在就返回为false
说明:
编辑:wang_jp
时间:2020年5月8日
*******************************************************************************/
func PathExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil { //已存在
		return true
	}
	if os.IsNotExist(err) { //不存在
		return false
	} else { //其他错误
		return false
	}
}

/*******************************************************************************
功能:判断文件夹是否存在,不存在就创建文件夹,存在就什么也不做直接返回
输入:文件夹路径
输出:
	[error]   创建错误返回错误信息，文件夹已存在或者正确创建返回nil
说明:
编辑:wang_jp
时间:2020年5月8日
*******************************************************************************/
func MakeDir(path string) error {
	exist := PathExists(path) //判断是否存在
	if exist == false {       //不存在就创建
		return os.Mkdir(path, os.ModePerm)
	}
	return nil //存在的话不做任何操作
}

/*******************************************************************************
功能:中文表示的星期几
输入:日期时间
输出:
	[string]   星期日~星期六
说明:
编辑:wang_jp
时间:2020年5月8日
*******************************************************************************/
func CnWeekday(t time.Time) (cw string) {
	w := t.Weekday()
	switch int(w) {
	case 0:
		cw = "星期日"
	case 1:
		cw = "星期一"
	case 2:
		cw = "星期二"
	case 3:
		cw = "星期三"
	case 4:
		cw = "星期四"
	case 5:
		cw = "星期五"
	case 6:
		cw = "星期六"
	default:
		cw = "错误"
	}
	return cw
}

/*******************************************************************************
功能:根据ID读取WorkShop信息
输入:id
输出:
	[error]   错误信息
说明:
编辑:wang_jp
时间:2020年5月8日
*******************************************************************************/
func (w *OreProcessDWorkshop) GetAttributeById(id ...int) error {
	o := orm.NewOrm()  //新建orm对象
	o.Using("default") //根据别名选择数据库
	if len(id) > 0 {
		w.Id = id[0]
	}
	return o.Read(w, "Id")
}

/*******************************************************************************
功能:获取车间列表
输入:无
输出:
	[error]   错误信息
说明:
编辑:wang_jp
时间:2020年5月8日
*******************************************************************************/
func (wshp *OreProcessDWorkshop) GetWorkshopLists() ([]OreProcessDWorkshop, error) {
	o := orm.NewOrm()  //新建orm对象
	o.Using("default") //根据别名选择数据库
	var w []OreProcessDWorkshop
	qt := o.QueryTable("OreProcessDWorkshop").Filter("Status", 1)
	_, err := qt.RelatedSel().OrderBy("Seq").All(&w)
	return w, err
}

/****************************************************
功能:毫秒Unix值转换为时间
输入:
	毫秒Unix时间值
输出：时间
编辑：wang_jp
时间：2020年5月16日
****************************************************/
func Millisecond2Time(mill_sec int64) time.Time {
	return time.Unix(mill_sec/1e3, mill_sec%1e3*1e6)
}

/****************************************************
功能：判断元素在数组、Map中是否存在
输入：元素、数组或者Map、Slice
输出：存在输出true，不存在输出false
说明：对于数组、Slice，判断的是值是否存在，对于Map，判断的是Key是否存在
时间：2019年12月15日
编辑：wang_jp
****************************************************/
func IsExistItem(obj interface{}, target interface{}) bool {
	targetValue := reflect.ValueOf(target)
	switch reflect.TypeOf(target).Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < targetValue.Len(); i++ {
			if targetValue.Index(i).Interface() == obj {
				return true
			}
		}
	case reflect.Map:
		if targetValue.MapIndex(reflect.ValueOf(obj)).IsValid() {
			return true
		}
	}
	return false
}
