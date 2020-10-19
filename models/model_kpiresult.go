package models

import (
	"errors"
	"fmt"
	"time"

	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
)

/*******************************************************************************
功能:删除通用结果表中的过期数据
输入:无
输出:无
说明:
编辑:wang_jp
时间:2019年12月14日
*******************************************************************************/
func (r *CalcKpiResult) DeleteOutTimeResultDataInComTable() {
	o := orm.NewOrm()                   //新建orm对象
	o.Using(EngineCfgMsg.ResultDBAlias) //根据别名选择数据库
	tm := time.Now().Add(-24 * time.Duration(EngineCfgMsg.Sys.SaveTimeInComTable) * time.Hour)
	sql := fmt.Sprintf(`DELETE FROM %s WHERE calc_ending_time <= "%s";`, EngineCfgMsg.CfgMsg.ResultdbTbname, tm.Format(EngineCfgMsg.Sys.TimeFormat))
	rer, err := o.Raw(sql).Exec() //执行查询
	if err != nil {               //发生错误
		logs.Error("Delete expired data error[删除过期数据时遇到错误]:[%s]", err.Error())
	}
	num, err := rer.RowsAffected()
	if err != nil { //发生错误
		logs.Error("Get RowsAffected error after delete expired data[删除过期数据后获取删除结果时遇到错误]:[%s]", err.Error())
	}
	logs.Info("Deleted %d rows expired data[删除了%d条过期数据]", num, num)
}

/*******************************************************************************
功能:批量保存同一个结果表中的KPI指标
输入:表名称，kpi指标结果结构体数组
输出:保存的数量
说明:
编辑:wang_jp
时间:2019年12月14日
*******************************************************************************/
func (r *CalcKpiResult) saveKpiResultToDB(tablename string, kpires []CalcKpiResult) int64 {
	o := orm.NewOrm()                   //新建orm对象
	o.Using(EngineCfgMsg.ResultDBAlias) //根据别名选择数据库
	//保存到通用结果表中
	sql := fmt.Sprintf(`INSERT INTO %s ( tag_type,tag_id, tag_name, kpi_config_list_id,kpi_key, kpi_tag,kpi_name, kpi_period, kpi_value, calc_ending_time, inbase_time )
			VALUES(?,?,?,?,?,?,?,?,?,?,?)`, EngineCfgMsg.CfgMsg.ResultdbTbname)

	p, err := o.Raw(sql).Prepare() //预处理
	if err != nil {
		logs.Alert("Prepare insert to [%s] error:[%s]", EngineCfgMsg.CfgMsg.ResultdbTbname, err.Error())
	}
	for _, res := range kpires {
		r, err := p.Exec(res.TagType, res.TagId, res.TagName, res.KpiConfigListId, res.KpiKey, res.KpiTag, res.KpiName, res.KpiPeriod, res.KpiValue, res.CalcEndingTime, res.InbaseTime)
		if err != nil {
			logs.Alert("Insert kpi value [%s] error:[%s]", res.KpiTag, err.Error())
		} else {
			r.RowsAffected()
		}
	}
	p.Close() //关闭预处理
	//保存到月度结果表中
	var cnt int64
	sql = fmt.Sprintf(`INSERT INTO %s ( tag_type,tag_id, tag_name, kpi_config_list_id,kpi_key, kpi_tag,kpi_name, kpi_period, kpi_value, calc_ending_time, inbase_time )
			VALUES(?,?,?,?,?,?,?,?,?,?,?)`, tablename)

	p, err = o.Raw(sql).Prepare() //预处理
	if err != nil {
		logs.Alert("Prepare insert to [%s] error:[%s]", tablename, err.Error())
	}
	for _, res := range kpires {
		r, err := p.Exec(res.TagType, res.TagId, res.TagName, res.KpiConfigListId, res.KpiKey, res.KpiTag, res.KpiName, res.KpiPeriod, res.KpiValue, res.CalcEndingTime, res.InbaseTime)
		if err != nil {
			logs.Alert("Insert kpi value [%s] error:[%s]", res.KpiTag, err.Error())
		} else {
			i, _ := r.RowsAffected()
			cnt += i
		}
	}
	p.Close() //关闭预处理
	return cnt
}

/*******************************************************************************
功能:批量保存KPI指标
输入:kpi指标结果结构体数组
输出:error:错误信息
说明:
编辑:wang_jp
时间:2019年12月17日
*******************************************************************************/
func (r *CalcKpiResult) SaveBatchKpiResultToDB(kpis []CalcKpiResult) int64 {
	var cnt int64
	kpigroups := r.groupByKpiByResultTableName(kpis) //对结果集进行分组
	for tbn, kpig := range kpigroups {               //按照组遍历、创建表并保存结果
		if _, err := r.createResultTable(EngineCfgMsg.ResultDBAlias, tbn); err != nil { //创建表
			logs.Alert(err)
		} else {
			cnt += r.saveKpiResultToDB(tbn, kpig) //保存结果集
		}
	}
	return cnt
}

/*******************************************************************************
功能:对KPI指标按照表名进行分组
输入:kpi指标结果结构体数组
输出:分组后的kpi结果集
说明:
编辑:wang_jp
时间:2019年12月17日
*******************************************************************************/
func (c *CalcKpiResult) groupByKpiByResultTableName(kpis []CalcKpiResult) map[string][]CalcKpiResult {
	var endtimegroup map[string][]CalcKpiResult
	endtimegroup = make(map[string][]CalcKpiResult) //初始化
	for _, k := range kpis {                        //遍历结果集,对结果集按照结束时间分组
		if len(k.KpiTag) > 0 {
			endtimegroup[k.CalcEndingTime] = append(endtimegroup[k.CalcEndingTime], k)
		}
	}

	var kgs map[string][]CalcKpiResult
	kgs = make(map[string][]CalcKpiResult) //初始化
	for k, et := range endtimegroup {      //遍历分组后的结果，再次按照表名前缀与时间后缀组成的表名分组
		if t, err := TimeParse(k); err != nil {
			logs.Error("Kpi result time format error:[%s]", err.Error())
		} else {
			tname := fmt.Sprintf("%s_%s", EngineCfgMsg.CfgMsg.ResultdbTbname, t.Format("2006_01"))
			for _, g := range et {
				kgs[tname] = append(kgs[tname], g)
			}
		}
	}
	return kgs
}

/*******************************************************************************
功能:检查数据库中是否存在指定表,如果不存在，则创建
输入:dbalias:数据库别名
	tablename:要检查的数据表名称
输出:bool:如果创建成功或者已经存在,true,没有创建成功false
	error:错误信息
说明:有错误信息的时候返回结果为false
编辑:wang_jp
时间:2019年12月13日
*******************************************************************************/
func (r *CalcKpiResult) createResultTable(dbalias, tablename string) (bool, error) {
	if exist, err := isTableExist(dbalias, tablename); err != nil { //检查表是否存在
		return false, err //发生错误
	} else if exist == true {
		return true, nil //表已经存在
	}

	o := orm.NewOrm() //新建orm对象
	o.Using(dbalias)  //根据别名选择数据库
	var res []orm.Params
	sql := fmt.Sprintf(`CREATE TABLE %s  (
	  id bigint(20) NOT NULL AUTO_INCREMENT COMMENT '自动id',
	  tag_type varchar(50) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '变量标签类型,process,goods等',
	  tag_id int(11) NULL DEFAULT NULL COMMENT '原始tag在taglist表中的id',
	  tag_name varchar(80) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '原始tag在taglist表中的tag_name',
	  kpi_config_list_id int(11) NULL DEFAULT NULL COMMENT '指标id,指标在calc_indicator_config_list中的id',
	  kpi_key varchar(50) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT 'Kpi关键词,与kpi_index_dic中的 kpi_varible_name 对应',
	  kpi_tag varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '指标名称',
      kpi_name varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '指标中文名称',
	  kpi_period int(11) NULL DEFAULT NULL COMMENT '指标计算周期：-1=每小时,-2=每班,-3=每日,-4=每月,-5=每季度,-6=每年。大于零时，为自由周期的时间值，单位为秒。',
	  kpi_value float NULL DEFAULT NULL COMMENT '指标的值',
	  calc_ending_time datetime(0) NULL DEFAULT NULL COMMENT '指标对应的时间。该时间是指标统计期末的时间值',
	  inbase_time datetime(0) NULL DEFAULT NULL COMMENT '当前记录生成时的时间',
	  PRIMARY KEY (id) USING BTREE
	) ENGINE = MyISAM AUTO_INCREMENT = 1 CHARACTER SET = utf8 COLLATE = utf8_general_ci ROW_FORMAT = Dynamic;`, tablename)
	_, err := o.Raw(sql).Values(&res) //执行查询
	if err != nil {                   //发生错误
		return false, errors.New(fmt.Sprintf(`Create table "%s" Error:[%s]`, tablename, err.Error()))
	}
	return true, nil //建表成功
}
