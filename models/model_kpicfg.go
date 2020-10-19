package models

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/bkzy-wangjp/MicEngine/MicScript/statistic"
	_ "github.com/go-sql-driver/mysql"
)

/*
功能:获取全部指标的配置信息
输入:kpitype:0获取全部,1获取串行指标,2获取并行计算指标
输出:获取到的指标数量，指标信息数组，错误信息
说明:
编辑:wang_jp
时间:2019年12月13日
*/
func (cfg *CalcKpiConfigList) GetKpiConfigInfo(kpitype int64) (int64, []CalcKpiConfigListExi, error) {
	proCnt, proKpi, proerr := cfg.getKpiConfigInfo_TypeProcess(kpitype) //读取过程变量的KPI配置
	if proerr != nil {
		//return proCnt, proKpi, proerr
		logs.Error("获取实时变量指标时发生错误:[%s]", proerr.Error())
	}
	cnt := proCnt
	if proCnt > 0 {
		logs.Info("载入[%d]个process变量指标", proCnt)
	}

	goodsCnt, goodsKpi, goodserr := cfg.getKpiConfigInfo_TypeGoods(kpitype) //读取物耗变量的KPI配置
	if goodserr != nil {
		//return goodsCnt, goodsKpi, goodserr
		logs.Error("获取物耗变量指标时发生错误:[%s]", goodserr.Error())
	}
	cnt += goodsCnt
	kpis := append(proKpi, goodsKpi...)
	if goodsCnt > 0 {
		logs.Info("载入[%d]个goods变量指标", goodsCnt)
	}

	patrolCnt, patrolKpi, patrolerr := cfg.getKpiConfigInfo_TypePatrol(kpitype) //读取巡检变量的KPI配置
	if patrolerr != nil {
		//return patrolCnt, patrolKpi, patrolerr
		logs.Error("获取巡检变量指标时发生错误:[%s]", patrolerr.Error())
	}
	cnt += patrolCnt
	kpis = append(kpis, patrolKpi...)
	if patrolCnt > 0 {
		logs.Info("载入[%d]个check变量指标", patrolCnt)
	}

	labCnt, labKpi, laberr := cfg.getKpiConfigInfo_TypeProductLab(kpitype) //读取紫金山化验变量的KPI配置
	if laberr != nil {
		//return labCnt, labKpi, laberr
		logs.Error("获取生产化验变量指标时发生错误:[%s]", laberr.Error())
	}
	cnt += labCnt
	kpis = append(kpis, labKpi...)
	if labCnt > 0 {
		logs.Info("载入[%d]个lab化验变量指标", labCnt)
	}

	sampleCnt, sampleKpi, samplerr := cfg.getKpiConfigInfo_TypeSample(kpitype) //读取化验变量的KPI配置
	if samplerr != nil {
		//return sampleCnt, sampleKpi, samplerr
		logs.Error("获取取样化验变量指标时发生错误:[%s]", samplerr.Error())
	}
	cnt += sampleCnt
	kpis = append(kpis, sampleKpi...)
	if sampleCnt > 0 {
		logs.Info("载入[%d]个sample变量指标", sampleCnt)
	}

	pro2Cnt, pro2Kpi, pro2err := cfg.getKpiConfigInfo_TypeProcess2(kpitype) //读取过程变量的KPI配置
	if pro2err != nil {
		//return pro2Cnt, pro2Kpi, pro2err
		logs.Error("获取慢实时变量指标时发生错误:[%s]", pro2err.Error())
	}
	cnt += pro2Cnt
	kpis = append(kpis, pro2Kpi...)
	if pro2Cnt > 0 {
		logs.Info("载入[%d]个Process2过程变量指标", pro2Cnt)
	}

	kpiCnt, kpiKpi, kpierr := cfg.getKpiConfigInfo_TypeKpi(kpitype) //读取过程变量的KPI配置
	if kpierr != nil {
		//return kpiCnt, kpiKpi, kpierr
		logs.Error("获取KPI二次计算变量指标时发生错误:[%s]", kpierr.Error())
	}
	cnt += kpiCnt
	kpis = append(kpis, kpiKpi...)
	if kpiCnt > 0 {
		logs.Info("载入[%d]个KPI二次计算指标", kpiCnt)
	}
	return cnt, kpis, nil
}

/*
功能:获取Process(过程变量)指标的配置信息
输入:kpitype:0获取全部,1获取串行指标,2获取并行计算指标
输出:获取到的指标数量，指标信息数组，错误信息
说明:
编辑:wang_jp
时间:2019年12月13日
*/
func (cfg *CalcKpiConfigList) getKpiConfigInfo_TypeProcess(kpitype int64) (int64, []CalcKpiConfigListExi, error) {
	o := orm.NewOrm() //新建orm对象
	o.Using("default")
	var kpi []CalcKpiConfigListExi
	var num int64
	var err error

	sql := `SELECT
			lst.id,
			lst.distributed_id,
			lst.tag_type,
			lst.tag_id,
			lst.kpi_index_dic_id,
			lst.kpi_tag,
			lst.kpi_name,
			lst.script,
			lst.start_time,
			lst.period,
			lst.offset_minutes,
			lst.last_calc_time,
			lst.supplement,
			lst.description,
			lst.seq,
			lst.status,
			lst.create_user_id,
			lst.create_time,
			lst.update_user_id,
			lst.update_time,
            lst.kpi_base_time,
			lst.kpi_shift_hour,
			tb.table_name,
			tag.tage_name AS tag_name,
			shp.base_time,
			shp.shift_hour,
			dic.script AS kpi_script, 
			dic.kpi_varible_name AS kpi_key
		FROM
			calc_kpi_config_list AS lst
			LEFT JOIN ore_process_d_taglist AS tag ON (lst.tag_id = tag.id)
			LEFT JOIN relevance_dcs_to_dbtable AS rel ON ( rel.dcs_id = tag.dcs_id )
			LEFT JOIN datatable_info AS tb ON ( tb.id = rel.datatable_id )
			LEFT JOIN ore_process_d_workstage stg ON ( stg.id = tag.stage_id )
			LEFT JOIN ore_process_d_workshop shp ON (stg.workshop_id = shp.id)
			LEFT JOIN calc_kpi_index_dic dic ON (lst.kpi_index_dic_id = dic.id )
		WHERE rel.status = 1 AND lst.status = 1 AND lst.tag_type = "process" AND tb.table_name IS NOT NULL AND`
	switch kpitype {
	case 0:
		sql = fmt.Sprintf("%s lst.distributed_id = %d ;", sql, EngineCfgMsg.CfgMsg.DistributedId)
	case 1:
		sql = fmt.Sprintf("%s lst.distributed_id = %d AND lst.period < 0 ;", sql, EngineCfgMsg.CfgMsg.DistributedId)
	case 2:
		sql = fmt.Sprintf("%s lst.distributed_id = %d AND lst.period > 0 ;", sql, EngineCfgMsg.CfgMsg.DistributedId)
	}
	num, err = o.Raw(sql).QueryRows(&kpi)

	return num, kpi, err
}

/*
功能:获取Process2(tag定义在process_taglist,但数据存储在sys_real_data数据表中的)指标的配置信息
输入:kpitype:0获取全部,1获取串行指标,2获取并行计算指标
输出:获取到的指标数量，指标信息数组，错误信息
说明:
编辑:wang_jp
时间:2020年2月18日
*/
func (cfg *CalcKpiConfigList) getKpiConfigInfo_TypeProcess2(kpitype int64) (int64, []CalcKpiConfigListExi, error) {
	o := orm.NewOrm() //新建orm对象
	o.Using("default")
	var kpi []CalcKpiConfigListExi
	var num int64
	var err error

	sql := `SELECT
			lst.id,
			lst.distributed_id,
			lst.tag_type,
			lst.tag_id,
			lst.kpi_index_dic_id,
			lst.kpi_tag,
			lst.kpi_name,
			lst.script,
			lst.start_time,
			lst.period,
			lst.offset_minutes,
			lst.last_calc_time,
			lst.supplement,
			lst.description,
			lst.seq,
			lst.status,
			lst.create_user_id,
			lst.create_time,
			lst.update_user_id,
			lst.update_time,
            lst.kpi_base_time,
			lst.kpi_shift_hour,
			tag.id AS table_name,
			tag.tage_name AS tag_name,
			shp.base_time,
			shp.shift_hour,
			dic.script AS kpi_script, 
			dic.kpi_varible_name AS kpi_key
		FROM
			calc_kpi_config_list AS lst
			LEFT JOIN ore_process_d_taglist AS tag ON (lst.tag_id = tag.id)
			LEFT JOIN ore_process_d_workstage stg ON ( stg.id = tag.stage_id )
			LEFT JOIN ore_process_d_workshop shp ON (stg.workshop_id = shp.id)
			LEFT JOIN calc_kpi_index_dic dic ON (lst.kpi_index_dic_id = dic.id )
		WHERE lst.status = 1 AND lst.tag_type = "process2" AND`
	switch kpitype {
	case 0:
		sql = fmt.Sprintf("%s lst.distributed_id = %d ;", sql, EngineCfgMsg.CfgMsg.DistributedId)
	case 1:
		sql = fmt.Sprintf("%s lst.distributed_id = %d AND lst.period < 0 ;", sql, EngineCfgMsg.CfgMsg.DistributedId)
	case 2:
		sql = fmt.Sprintf("%s lst.distributed_id = %d AND lst.period > 0 ;", sql, EngineCfgMsg.CfgMsg.DistributedId)
	}
	num, err = o.Raw(sql).QueryRows(&kpi)

	return num, kpi, err
}

/*
功能:获取Goods(物耗)指标的配置信息
输入:kpitype:0获取全部,1获取串行指标,2获取并行计算指标
输出:获取到的指标数量，指标信息数组，错误信息
说明:
编辑:wang_jp
时间:2019年12月13日
*/
func (cfg *CalcKpiConfigList) getKpiConfigInfo_TypeGoods(kpitype int64) (int64, []CalcKpiConfigListExi, error) {
	o := orm.NewOrm() //新建orm对象
	o.Using("default")
	var kpi []CalcKpiConfigListExi
	var num int64
	var err error

	sql := `SELECT
			lst.id,
			lst.distributed_id,
			lst.tag_type,
			lst.tag_id,
			lst.kpi_index_dic_id,
			lst.kpi_tag,
			lst.kpi_name,
			lst.script,
			lst.start_time,
			lst.period,
			lst.offset_minutes,
			lst.last_calc_time,
			lst.supplement,
			lst.description,
			lst.seq,
			lst.status,
			lst.create_user_id,
			lst.create_time,
			lst.update_user_id,
			lst.update_time,
			lst.kpi_base_time,
			lst.kpi_shift_hour,
			tag.seq AS table_name,
			tag.goods_tag_name AS tag_name,
			shp.base_time,
			shp.shift_hour,
			dic.script AS kpi_script, 
			dic.kpi_varible_name AS kpi_key
		FROM
			calc_kpi_config_list AS lst
			LEFT JOIN goods_config_info AS tag ON (lst.tag_id = tag.id)
			LEFT JOIN ore_process_d_workshop shp ON (tag.workshop_id = shp.id)
			LEFT JOIN calc_kpi_index_dic dic ON (lst.kpi_index_dic_id = dic.id )
		WHERE lst.status = 1 AND lst.tag_type = "goods"  AND`
	switch kpitype {
	case 0:
		sql = fmt.Sprintf("%s lst.distributed_id = %d ;", sql, EngineCfgMsg.CfgMsg.DistributedId)
	case 1:
		sql = fmt.Sprintf("%s lst.distributed_id = %d AND lst.period < 0 ;", sql, EngineCfgMsg.CfgMsg.DistributedId)
	case 2:
		sql = fmt.Sprintf("%s lst.distributed_id = %d AND lst.period > 0 ;", sql, EngineCfgMsg.CfgMsg.DistributedId)
	}
	num, err = o.Raw(sql).QueryRows(&kpi)

	return num, kpi, err
}

/*
功能:获取Patrol(巡检)指标的配置信息
输入:kpitype:0获取全部,1获取串行指标,2获取并行计算指标
输出:获取到的指标数量，指标信息数组，错误信息
说明:
编辑:wang_jp
时间:2019年12月13日
*/
func (cfg *CalcKpiConfigList) getKpiConfigInfo_TypePatrol(kpitype int64) (int64, []CalcKpiConfigListExi, error) {
	o := orm.NewOrm() //新建orm对象
	o.Using("default")
	var kpi []CalcKpiConfigListExi
	var num int64
	var err error

	sql := `SELECT
			lst.id,
			lst.distributed_id,
			lst.tag_type,
			lst.tag_id,
			lst.kpi_index_dic_id,
			lst.kpi_tag,
			lst.kpi_name,
			lst.script,
			lst.start_time,
			lst.period,
			lst.offset_minutes,
			lst.last_calc_time,
			lst.supplement,
			lst.description,
			lst.seq,
			lst.status,
			lst.create_user_id,
			lst.create_time,
			lst.update_user_id,
			lst.update_time,
            lst.kpi_base_time,
			lst.kpi_shift_hour,
			tag.id AS table_name,
			tag.tag_name AS tag_name,
			shp.base_time,
			shp.shift_hour,
			dic.script AS kpi_script, 
			dic.kpi_varible_name AS kpi_key
		FROM
			calc_kpi_config_list AS lst
			LEFT JOIN check_tag_list AS tag ON (lst.tag_id = tag.id)
			LEFT JOIN ore_process_d_workstage stg ON ( stg.id = tag.stage_id )
			LEFT JOIN ore_process_d_workshop shp ON (stg.workshop_id = shp.id)
			LEFT JOIN calc_kpi_index_dic dic ON (lst.kpi_index_dic_id = dic.id )
		WHERE lst.status = 1 AND (lst.tag_type = "patrol" OR lst.tag_type = "check")  AND`
	switch kpitype {
	case 0:
		sql = fmt.Sprintf("%s lst.distributed_id = %d ;", sql, EngineCfgMsg.CfgMsg.DistributedId)
	case 1:
		sql = fmt.Sprintf("%s lst.distributed_id = %d AND lst.period < 0 ;", sql, EngineCfgMsg.CfgMsg.DistributedId)
	case 2:
		sql = fmt.Sprintf("%s lst.distributed_id = %d AND lst.period > 0 ;", sql, EngineCfgMsg.CfgMsg.DistributedId)
	}
	num, err = o.Raw(sql).QueryRows(&kpi)

	return num, kpi, err
}

/*
功能:获取Lab(化验,仅用于紫金山Mes化验数据)指标的配置信息
输入:kpitype:0获取全部,1获取串行指标,2获取并行计算指标
输出:获取到的指标数量，指标信息数组，错误信息
说明:
编辑:wang_jp
时间:2019年12月13日
*/
func (cfg *CalcKpiConfigList) getKpiConfigInfo_TypeLab(kpitype int64) (int64, []CalcKpiConfigListExi, error) {
	o := orm.NewOrm() //新建orm对象
	o.Using("default")
	var kpi []CalcKpiConfigListExi
	var num int64
	var err error
	if exist, e := isTableExist("default", "two_assay_lab_list"); exist == true && e == nil {
		sql := `SELECT
			lst.id,
			lst.distributed_id,
			lst.tag_type,
			lst.tag_id,
			lst.kpi_index_dic_id,
			lst.kpi_tag,
			lst.kpi_name,
			lst.script,
			lst.start_time,
			lst.period,
			lst.offset_minutes,
			lst.last_calc_time,
			lst.supplement,
			lst.description,
			lst.seq,
			lst.status,
			lst.create_user_id,
			lst.create_time,
			lst.update_user_id,
			lst.update_time,
			lst.id AS table_name,
            lst.kpi_base_time,
			lst.kpi_shift_hour,
			tag.AssayTagName AS tag_name,
			shp.base_time,
			shp.shift_hour,
			dic.script AS kpi_script, 
			dic.kpi_varible_name AS kpi_key
		FROM
			calc_kpi_config_list AS lst
			LEFT JOIN two_assay_lab_list AS tag ON (lst.tag_id = tag.AssayTagID)
			LEFT JOIN ore_process_d_workshop shp ON (tag.Workshopid = shp.id)
			LEFT JOIN calc_kpi_index_dic dic ON (lst.kpi_index_dic_id = dic.id )
		WHERE lst.status = 1 AND lst.tag_type = "lab"  AND`
		switch kpitype {
		case 0:
			sql = fmt.Sprintf("%s lst.distributed_id = %d ;", sql, EngineCfgMsg.CfgMsg.DistributedId)
		case 1:
			sql = fmt.Sprintf("%s lst.distributed_id = %d AND lst.period < 0 ;", sql, EngineCfgMsg.CfgMsg.DistributedId)
		case 2:
			sql = fmt.Sprintf("%s lst.distributed_id = %d AND lst.period > 0 ;", sql, EngineCfgMsg.CfgMsg.DistributedId)
		}
		num, err = o.Raw(sql).QueryRows(&kpi)
	}
	return num, kpi, err
}

/*
功能:获取Lab(化验,仅用于紫金山Mes化验数据)指标的配置信息
输入:kpitype:0获取全部,1获取串行指标,2获取并行计算指标
输出:获取到的指标数量，指标信息数组，错误信息
说明:
编辑:wang_jp
时间:2019年12月13日
*/
func (cfg *CalcKpiConfigList) getKpiConfigInfo_TypeProductLab(kpitype int64) (int64, []CalcKpiConfigListExi, error) {
	o := orm.NewOrm() //新建orm对象
	o.Using("default")
	var kpi []CalcKpiConfigListExi
	var num int64
	var err error
	if exist, e := isTableExist("default", "trans_other_system_list"); exist == true && e == nil {
		sql := `SELECT
			lst.id,
			lst.distributed_id,
			lst.tag_type,
			lst.tag_id,
			lst.kpi_index_dic_id,
			lst.kpi_tag,
			lst.kpi_name,
			lst.script,
			lst.start_time,
			lst.period,
			lst.offset_minutes,
			lst.last_calc_time,
			lst.supplement,
			lst.description,
			lst.seq,
			lst.status,
			lst.create_user_id,
			lst.create_time,
			lst.update_user_id,
			lst.update_time,
			lst.id AS table_name,
            lst.kpi_base_time,
			lst.kpi_shift_hour,
			tag.tag_name AS tag_name,
			shp.base_time,
			shp.shift_hour,
			dic.script AS kpi_script, 
			dic.kpi_varible_name AS kpi_key
		FROM
			calc_kpi_config_list AS lst
			LEFT JOIN trans_other_system_list AS tag ON (lst.tag_id = tag.id)
			LEFT JOIN ore_process_d_workshop shp ON (tag.workshop_id = shp.id)
			LEFT JOIN calc_kpi_index_dic dic ON (lst.kpi_index_dic_id = dic.id )
		WHERE lst.status = 1 AND (lst.tag_type = "lab" OR lst.tag_type = "product_lab") AND`
		switch kpitype {
		case 0:
			sql = fmt.Sprintf("%s lst.distributed_id = %d ;", sql, EngineCfgMsg.CfgMsg.DistributedId)
		case 1:
			sql = fmt.Sprintf("%s lst.distributed_id = %d AND lst.period < 0 ;", sql, EngineCfgMsg.CfgMsg.DistributedId)
		case 2:
			sql = fmt.Sprintf("%s lst.distributed_id = %d AND lst.period > 0 ;", sql, EngineCfgMsg.CfgMsg.DistributedId)
		}
		num, err = o.Raw(sql).QueryRows(&kpi)
	}
	return num, kpi, err
}

/*
功能:获取Sample(化验)指标的配置信息
输入:kpitype:0获取全部,1获取串行指标,2获取并行计算指标
输出:获取到的指标数量，指标信息数组，错误信息
说明:
编辑:wang_jp
时间:2020年2月9日
*/
func (cfg *CalcKpiConfigList) getKpiConfigInfo_TypeSample(kpitype int64) (int64, []CalcKpiConfigListExi, error) {
	o := orm.NewOrm() //新建orm对象
	o.Using("default")
	var kpi []CalcKpiConfigListExi
	var num int64
	var err error
	if exist, e := isTableExist("default", "sampling_manage_sub"); exist == true && e == nil {
		sql := `SELECT
			lst.id,
			lst.distributed_id,
			lst.tag_type,
			lst.tag_id,
			lst.kpi_index_dic_id,
			lst.kpi_tag,
			lst.kpi_name,
			lst.script,
			lst.start_time,
			lst.period,
			lst.offset_minutes,
			lst.last_calc_time,
			lst.supplement,
			lst.description,
			lst.seq,
			lst.status,
			lst.create_user_id,
			lst.create_time,
			lst.update_user_id,
			lst.update_time,
			lst.id AS table_name,
            lst.kpi_base_time,
			lst.kpi_shift_hour,
			tag.sample_index_tag AS tag_name,
			shp.base_time,
			shp.shift_hour,
			dic.script AS kpi_script, 
			dic.kpi_varible_name AS kpi_key
		FROM
			calc_kpi_config_list AS lst
			LEFT JOIN sampling_manage_sub AS tag ON (lst.tag_id = tag.id)
			LEFT JOIN ore_process_d_workstage stage ON (tag.sample_from_stage_id = stage.id)
			LEFT JOIN ore_process_d_workshop shp ON (stage.Workshop_id = shp.id)
			LEFT JOIN calc_kpi_index_dic dic ON (lst.kpi_index_dic_id = dic.id )
		WHERE lst.status = 1 AND lst.tag_type = "sample"  AND`
		switch kpitype {
		case 0:
			sql = fmt.Sprintf("%s lst.distributed_id = %d ;", sql, EngineCfgMsg.CfgMsg.DistributedId)
		case 1:
			sql = fmt.Sprintf("%s lst.distributed_id = %d AND lst.period < 0 ;", sql, EngineCfgMsg.CfgMsg.DistributedId)
		case 2:
			sql = fmt.Sprintf("%s lst.distributed_id = %d AND lst.period > 0 ;", sql, EngineCfgMsg.CfgMsg.DistributedId)
		}
		num, err = o.Raw(sql).QueryRows(&kpi)
	}
	return num, kpi, err
}

/*
功能:获取kpi二次计算指标的配置信息
输入:kpitype:0获取全部,1获取串行指标,2获取并行计算指标
输出:获取到的指标数量，指标信息数组，错误信息
说明:
编辑:wang_jp
时间:2020年5月5日
*/
func (cfg *CalcKpiConfigList) getKpiConfigInfo_TypeKpi(kpitype int64) (int64, []CalcKpiConfigListExi, error) {
	o := orm.NewOrm() //新建orm对象
	o.Using("default")
	var kpi []CalcKpiConfigListExi
	var num int64
	var err error
	sql := `SELECT
			lst.id,
			lst.distributed_id,
			lst.tag_type,
			lst.tag_id,
			lst.kpi_index_dic_id,
			lst.kpi_tag,
			lst.kpi_name,
			lst.script,
			lst.start_time,
			lst.period,
			lst.offset_minutes,
			lst.last_calc_time,
			lst.supplement,
			lst.description,
			lst.seq,
			lst.status,
			lst.create_user_id,
			lst.create_time,
			lst.update_user_id,
			lst.update_time,
			lst.id AS table_name,
			lst.kpi_base_time,
			lst.kpi_shift_hour,
			lst.kpi_name AS tag_name,
			lst.kpi_base_time AS base_time,
			lst.kpi_shift_hour AS shift_hour,
			dic.script AS kpi_script, 
			dic.kpi_varible_name AS kpi_key
		FROM
			calc_kpi_config_list AS lst
			LEFT JOIN calc_kpi_index_dic dic ON (lst.kpi_index_dic_id = dic.id )
		WHERE lst.status = 1 AND lst.tag_type = "kpi"  AND`
	switch kpitype {
	case 0:
		sql = fmt.Sprintf("%s lst.distributed_id = %d ;", sql, EngineCfgMsg.CfgMsg.DistributedId)
	case 1:
		sql = fmt.Sprintf("%s lst.distributed_id = %d AND lst.period < 0 ;", sql, EngineCfgMsg.CfgMsg.DistributedId)
	case 2:
		sql = fmt.Sprintf("%s lst.distributed_id = %d AND lst.period > 0 ;", sql, EngineCfgMsg.CfgMsg.DistributedId)
	}
	num, err = o.Raw(sql).QueryRows(&kpi)

	return num, kpi, err
}

/******************************************************************************
功能:根据tagid或者tagname以及时间参数读取 kpi_result 表中的数据
输入:tagid,tagname,key,beginTime,endTime
输出:value和错误信息
说明:如果tagid为0,则通过tagname查询数据
编辑:wang_jp
时间:2020年05月04日
******************************************************************************/
func (cfg *CalcKpiConfigList) GetKpiDataStatisticByKey(key, beginTime, endTime string) (float64, bool, error) {
	_, err := cfg.GetKpiConfigListInfo()
	if err != nil {
		return 0.0, false, err
	}

	switch strings.ToLower(key) {
	case "bgpoint": //读取单点数据(起始时间地点)
		value, err := cfg.GetKpiSingleData(beginTime, "")
		if err != nil {
			if strings.Contains(err.Error(), "no data found") { //没有查询到数据
				res, e := cfg.GetKpiDatas(beginTime, endTime, 1) //读取指定时间范围之后的数据
				if len(res) > 0 && e == nil {                    //后面有新数据，但是当期没有数据
					return 0.0, true, err //输出特殊错误标志"0",计算继续往下进行
				} else {
					return 0.0, false, err
				}
			} else {
				return 0.0, false, err
			}
		}
		return value, false, nil
	case "endpoint": //读取单点数据(结束时间点)
		value, err := cfg.GetKpiSingleData("", endTime)
		if err != nil {
			if strings.Contains(err.Error(), "no data found") { //没有查询到数据
				res, e := cfg.GetKpiDatas(beginTime, endTime, 1) //读取指定时间范围之后的数据
				if len(res) > 0 && e == nil {                    //后面有新数据，但是当期没有数据
					return 0.0, true, err //输出特殊错误标志"0",计算继续往下进行
				} else {
					return 0.0, false, err
				}
			} else {
				return 0.0, false, err
			}
		}
		return value, false, nil
	default: //读取统计结果数据
		sttt, ctnue, err := cfg.GetKpiStatistic(key, beginTime, endTime)
		if err != nil {
			return 0.0, ctnue, err
		}
		return statistic.SelectValueFromStatisticData(sttt, key), false, nil
	}
	return 0.0, false, nil
}

/******************************************************************************
功能:通过kpi_id查询单个KPI值
输入:beginTime,endTime
输出:kpivalue,err
说明:当endTime有值时,返回endTime时刻的kpi值;当beginTime有值时,返回beginTime时刻的kpi值
	endTime值优先。
编辑:wang_jp
时间:2020年05月04日
*******************************************************************************/
func (cfg *CalcKpiConfigList) GetKpiSingleData(beginTime, endTime string) (float64, error) {
	var value float64
	if len(endTime) > 0 {
		kpis, err := cfg.GetKpiDatas(beginTime, endTime, 2)
		if err != nil {
			return value, err
		}
		for _, kpi := range kpis {
			value = kpi.KpiValue
		}
	} else if len(beginTime) > 0 {
		kpis, err := cfg.GetKpiDatas(beginTime, endTime, -2)
		if err != nil {
			return value, err
		}
		for _, kpi := range kpis {
			value = kpi.KpiValue
		}
	} else {
		return value, fmt.Errorf("KPI id=%d invalid begintime and endtime in GetKpiSingleData[获取单点KPI数据时没有设定有效的起止时间]", cfg.Id)
	}

	return value, nil
}

/******************************************************************************
功能:通过kpi_id查询KPI值
输入:beginTime,endTime
输出:kpivalue,err
说明:
	0.如果readType为-2,则取等于beginTime的一个数据点,endTime无效
	1.如果readType为-1,则取小于等于beginTime的一个数据点,endTime无效
	2.如果readType为 0,则取大于beginTime和小于等于endTime之间的数
	3.如果readType为 1,则取大于等于endTime的数据点,beginTime无效
	4.如果readType为 2,则取等于endTime的数据点,beginTime无效
编辑:wang_jp
时间:2020年05月04日
*******************************************************************************/
func (cfg *CalcKpiConfigList) GetKpiDatas(beginTime, endTime string, readType int) ([]CalcKpiResult, error) {
	if cfg.Id <= 0 {
		_, err := cfg.GetKpiConfigListInfo()
		if err != nil {
			return nil, err
		}
	}
	o := orm.NewOrm()                   //新建orm对象
	o.Using(EngineCfgMsg.ResultDBAlias) //根据别名选择数据库

	bgtunset := false //未设定有效的开始时间
	oldesttime := time.Now().Add(-24 * time.Duration(EngineCfgMsg.Sys.SaveTimeInComTable) * time.Hour)
	bgtime, err := TimeParse(beginTime)
	if err != nil {
		if readType <= 0 { //需要使用开始时间的时候而没有设定有效的开始时间
			return nil, fmt.Errorf("KPI id=%d begintime format error[查询开始时间格式错误]:[%s]", cfg.Id, err.Error())
		} else {
			bgtunset = true
		}
	}
	edtime, err := TimeParse(endTime)
	if err != nil {
		if readType >= 0 { //需要使用结束时间的时候而没有设定有效的结束时间
			return nil, fmt.Errorf("KPI id=%d endtime format error[查询结束时间格式错误]:[%s]", cfg.Id, err.Error())
		} else {
			edtime = bgtime.Add(1 * time.Hour)
		}
	}
	if bgtunset { //如果未设定起始时间
		bgtime = edtime.Add(-1 * time.Hour)
	}

	type timerange struct {
		Begin string
		End   string
		Table string
	}
	var trs []timerange
	var tr timerange
	if oldesttime.Before(bgtime) { //开始时间在界限范围内
		tr.Begin = beginTime
		tr.End = endTime
		tr.Table = EngineCfgMsg.CfgMsg.ResultdbTbname
		trs = append(trs, tr)
	} else { //开始时间超出了时间界限
		isbreak := false
		tr.End = endTime
		if edtime.Before(oldesttime) { //结束时间也不在界限范围内
			t0, _ := TimeParse(tr.End)
			tms := TimeMonthStart(t0.Add(-1 * time.Second)) //结束时间所在月的开始时间
			if tms.Before(bgtime) {                         //结束时间所在月的开始时间早于起始时间
				tr.Begin = beginTime
				isbreak = true
			} else { //结束时间所在月的开始时间晚于起始时间
				tr.Begin = TimeFormat(tms)
			}
			tr.Table = fmt.Sprintf("%s_%s", EngineCfgMsg.CfgMsg.ResultdbTbname, TimeFormat(tms, "2006_01"))
		} else { //结束时间在界限范围内
			tr.Begin = TimeFormat(oldesttime)
			tr.Table = EngineCfgMsg.CfgMsg.ResultdbTbname
		}
		trs = append(trs, tr)

		for isbreak == false {
			tr.End = tr.Begin
			t0, _ := TimeParse(tr.End)
			tms := TimeMonthStart(t0.Add(-1 * time.Second)) //结束时间所在月的开始时间
			if tms.Before(bgtime) {                         //结束时间所在月的开始时间早于起始时间
				tr.Begin = beginTime
				isbreak = true
			} else { //结束时间所在月的开始时间晚于起始时间
				tr.Begin = TimeFormat(tms)
			}
			tr.Table = fmt.Sprintf("%s_%s", EngineCfgMsg.CfgMsg.ResultdbTbname, TimeFormat(tms, "2006_01"))
			trs = append(trs, tr)
			if isbreak {
				break
			}
		}
	}
	var hiswait sync.WaitGroup
	kpics := make(chan []CalcKpiResult, len(trs))
	gocnt := 0
	defer close(kpics)
	for i, tr := range trs {
		var sqlstr, sql string
		isbreak := false
		iscontinu := false
		switch readType {
		case -1:
			if i == len(trs)-1 {
				sqlstr = `SELECT * FROM %s WHERE kpi_config_list_id = %d AND calc_ending_time <= "%s" ORDER BY calc_ending_time DESC LIMIT 1`
				sql = fmt.Sprintf(sqlstr, tr.Table, cfg.Id, tr.Begin)
			} else {
				iscontinu = true
				continue
			}
		case -2:
			if i == len(trs)-1 {
				sqlstr = `SELECT * FROM %s WHERE kpi_config_list_id = %d AND calc_ending_time = "%s"`
				sql = fmt.Sprintf(sqlstr, tr.Table, cfg.Id, tr.Begin)
			} else {
				iscontinu = true
				continue
			}
		case 0:
			sqlstr = `SELECT * FROM %s WHERE kpi_config_list_id = %d AND calc_ending_time > "%s" AND calc_ending_time <= "%s"`
			sql = fmt.Sprintf(sqlstr, tr.Table, cfg.Id, tr.Begin, tr.End)
		case 1:
			sqlstr = `SELECT * FROM %s WHERE kpi_config_list_id = %d AND calc_ending_time >= "%s" ORDER BY calc_ending_time ASC LIMIT 1`
			sql = fmt.Sprintf(sqlstr, tr.Table, cfg.Id, tr.End)
			isbreak = true
		case 2:
			sqlstr = `SELECT * FROM %s WHERE kpi_config_list_id = %d AND calc_ending_time = "%s"`
			sql = fmt.Sprintf(sqlstr, tr.Table, cfg.Id, tr.End)
			isbreak = true
		}
		if iscontinu == false {
			gocnt++
			hiswait.Add(1)
			go func(sql_str string, ks chan []CalcKpiResult) {
				defer hiswait.Done()
				var kpi []CalcKpiResult
				_, err := o.Raw(sql_str).QueryRows(&kpi)
				if err != nil {
				}
				ks <- kpi
			}(sql, kpics)
		}
		if isbreak {
			break
		}
	}
	var kpis []CalcKpiResult
	for i := 0; i < gocnt; i++ {
		ks := <-kpics
		kpis = append(kpis, ks...)
	}
	hiswait.Wait()
	if len(kpis) == 0 {
		switch readType {
		case -2:
			err = fmt.Errorf("KPI id=%d no data found at the point of time %s[在指定的时间点没有查询到数据]", cfg.Id, beginTime)
		case -1:
			err = fmt.Errorf("KPI id=%d no data found before the point of time %s[在指定的时间点前没有查询到数据]", cfg.Id, beginTime)
		case 0:
			err = fmt.Errorf("KPI id=%d no data found between %s and %s[在指定条件下没有查询到数据]", cfg.Id, beginTime, endTime)
		case 1:
			err = fmt.Errorf("KPI id=%d no data found after the point of time %s[在指定的时间点后没有查询到数据]", cfg.Id, endTime)
		case 2:
			err = fmt.Errorf("KPI id=%d no data found at the point of time %s[在指定的时间点没有查询到数据]", cfg.Id, endTime)
		default:
			err = fmt.Errorf("KPI id=%d no data found between %s and %s[在指定的时间范围内没有查询到数据]", cfg.Id, beginTime, endTime)
		}
	}
	return kpis, err
}

/******************************************************************************
功能:通过kpi_tag查询kpi的id
输入:kpi_tag
输出:kpi_id,err
说明:
编辑:wang_jp
时间:2020年05月04日
*******************************************************************************/
func (cfg *CalcKpiConfigList) GetKpiConfigListInfo() (int64, error) {
	o := orm.NewOrm()  //新建orm对象
	o.Using("default") //根据别名选择数据库
	var e error
	var tag string
	if cfg.Id <= 0 {
		err := o.Read(cfg, "KpiTag")
		e = err
		tag = cfg.KpiTag
	} else {
		err := o.Read(cfg)
		e = err
		tag = fmt.Sprintf("%d", cfg.Id)
	}

	if e == orm.ErrNoRows {
		e = fmt.Errorf("查询不到 KPI TAG [%s]", tag)
	} else if e == orm.ErrMissPK {
		e = fmt.Errorf("通过 [%s] 找不到主键", tag)
	}
	return cfg.Id, e
}

/*******************************************************************************
功能:根据tagid以及时间参数对读取自kpi_result表中的数据,并进行统计计算
输入:key,beginTime,endTime
输出:结果结构和错误信息
说明:
编辑:wang_jp
时间:2020年05月04日
*******************************************************************************/
func (cfg *CalcKpiConfigList) GetKpiStatistic(key, beginTime, endTime string) (statistic.StatisticData, bool, error) {
	var sttt statistic.StatisticData //返回的统计数据
	var tsds statistic.Tsds          //时间序列数组
	var tsd statistic.TimeSeriesData
	if key == "diff" || key == "plusdiff" { //求差的时候需要带上上一个计算周期的最后一个值
		rtds, err := cfg.GetKpiDatas(beginTime, endTime, -1) //读取指定时间点的数据
		if err == nil {                                      //如果有读到数据开始时间点之前的数据
			for _, rtd := range rtds { //整理出标准的时间序列数据
				tsd.Time, err = TimeParse(rtd.CalcEndingTime)
				if err == nil {
					tsd.Value = rtd.KpiValue
					tsds = append(tsds, tsd)
				}
			}
		}
	}
	rtds, err := cfg.GetKpiDatas(beginTime, endTime, 0) //读取指定时间范围内的数据
	if err != nil {                                     //当期查询错误
		if strings.Contains(err.Error(), "no data found in the set condition") || strings.Contains(err.Error(), "在指定条件下没有查询到数据") { //没有查询到数据
			res, e := cfg.GetKpiDatas(beginTime, endTime, 1) //读取指定时间范围之后的数据
			if len(res) > 0 && e == nil {                    //后面有新数据，但是当期没有数据
				return sttt, true, fmt.Errorf("%s,最近的数据时间是:%s", err.Error(), res[0].CalcEndingTime) //输出特殊错误标志"0",计算继续往下进行
			} else {
				return sttt, false, err
			}
		} else {
			return sttt, false, err
		}
	}
	for _, rtd := range rtds { //整理出标准的时间序列数据
		tsd.Time, err = TimeParse(rtd.CalcEndingTime)
		if err == nil {
			tsd.Value = rtd.KpiValue
			tsds = append(tsds, tsd)
		}
	}
	switch strings.ToLower(key) {
	case "advangce", "sd", "stddev", "se", "ske", "kur", "mode", "median", "groupdist", "distribution":
		sttt = tsds.Statistics(1) //包含高级历史统计
	default:
		sttt = tsds.Statistics(0) //只有普通统计
	}
	sttt.BeginTime = beginTime
	sttt.EndTime = endTime

	return sttt, false, nil
}

/*******************************************************************************
功能:获取本次计算的开始时间点和结束时间点，并判断是否可以开始计算
输入:无
输出:是否可以计算\计算开始时间点\结束时间点\错误信息
说明:
编辑:wang_jp
时间:2020年5月7日
*******************************************************************************/
func (cfg *CalcKpiConfigListExi) PrevTime() (bool, string, string, error) {
	if cfg.KpiShiftHour == 0 || len(cfg.KpiBaseTime) < 10 {
		cfg.KpiShiftHour = cfg.ShiftHour
		cfg.KpiBaseTime = cfg.BaseTime
	}
	if cfg.KpiShiftHour == 0 || len(cfg.KpiBaseTime) < 10 {
		return false, "", "", fmt.Errorf("KPI指标没有设定基准时间和每班工作时间,KPI_ID:[ %d ],Name:[ %s ]", cfg.Id, cfg.KpiName)
	}
	calcEn, st, ed, err := prevTime(cfg.KpiBaseTime, cfg.StartTime, cfg.LastCalcTime, cfg.KpiShiftHour, cfg.Period, cfg.OffsetMinutes)
	if err != nil {
		err = fmt.Errorf("KPI_ID:[%d][%s]", cfg.Id, err.Error())
	}
	return calcEn, st, ed, err
}

/*******************************************************************************
功能:KPI计算
输入:无
输出:计算结果\错误信息
说明:
编辑:wang_jp
时间:2020年5月7日
*******************************************************************************/
func (kpi *CalcKpiConfigListExi) KpiCalc(beginTime, endTime string) (CalcKpiResult, bool, error) {
	var res CalcKpiResult //新建结果,这个结果集用于保存计算结果
	defer func() {
		if err := recover(); err != nil {
			logs.Critical("KpiCalc中遇到错误:%d;[%#v]", kpi.Id, err)
		}
	}()
	if len(kpi.KpiScript) > 0 || len(kpi.Script) > 0 { //必须有有效的脚本程序
		script_str := kpi.KpiScript
		if len(kpi.KpiScript) == 0 {
			script_str = kpi.Script
		}
		tagfullname := fmt.Sprintf("%s.%s", kpi.TableName, kpi.TagName)

		script := new(Script)
		script.Id = kpi.Id
		script.BaseTime = kpi.KpiBaseTime
		script.BeginTime = beginTime
		script.EndTime = endTime
		script.ShiftHour = kpi.KpiShiftHour
		script.ScriptStr = script_str
		script.MainTagFullName = tagfullname
		script.MainTagId = kpi.TagId

		if value, ctinue, err := script.Run(); err != nil {
			if ctinue == true { //库中没有数据
				kpi.LastCalcTime = endTime
				go kpi.SetKpiLastCalcTimeSingle(kpi.Id, endTime)
				return res, ctinue, fmt.Errorf("KPI计算脚本运行错误,KPI_ID:[%d],错误信息:[%s]", kpi.Id, err.Error())
			} else {
				return res, ctinue, fmt.Errorf("KPI计算脚本运行错误,KPI_ID:[%d],错误信息:[%s]", kpi.Id, err.Error())
			}
		} else { //正确完成计算
			kpi.LastCalcTime = endTime

			res.CalcEndingTime = endTime
			res.InbaseTime = time.Now().Format(EngineCfgMsg.Sys.TimeFormat)
			res.KpiConfigListId = kpi.Id
			res.KpiPeriod = kpi.Period
			res.KpiKey = kpi.KpiKey
			res.KpiTag = kpi.KpiTag
			res.KpiName = kpi.KpiName
			res.KpiValue = value.(float64)
			res.TagType = kpi.TagType
			res.TagId = kpi.TagId
			res.TagName = kpi.TagName
		}
	} else {
		return res, false, fmt.Errorf("KPI指标中没有定义有效的脚本,KPI定义信息:[%+v]", kpi)
	}
	return res, true, nil
}

/*******************************************************************************
功能:获取KPI指标的基本配置列表
输入:无
输出:保存的数据数量
说明:
编辑:wang_jp
时间:2020年7月14日
*******************************************************************************/
func (cfg *CalcKpiConfigList) GetKpiConfigBaseMsg() (ids, tagids, typeids, shifthour []int, tagtypes, names, descs, basetime []string, err error) {
	o := orm.NewOrm()  //新建orm对象
	o.Using("default") //根据别名选择数据库
	sql := `SELECT
		lst.id,
		lst.tag_type,
		lst.tag_id,
		lst.kpi_name,
		lst.description,
		lst.kpi_base_time as base_time,
		lst.kpi_shift_hour as shift_hour,
		dic.id as type_id 
	FROM
		calc_kpi_config_list AS lst
		JOIN sys_dictionary AS dic ON ( dic.dictionary_name_code = lst.tag_type ) 
	WHERE
		dic.dic_catalog_id = 7 
		AND dic.dictionary_code > 0 
		AND lst.status = 1;`
	type list struct {
		Id          int
		Discription string
		TagType     string
		TagId       int
		KpiName     string
		TypeId      int
		BaseTime    string
		ShiftHour   int
	}
	var lsts []list
	_, err = o.Raw(sql).QueryRows(&lsts)
	for _, kpi := range lsts {
		ids = append(ids, kpi.Id)
		tagids = append(tagids, kpi.TagId)
		tagtypes = append(tagtypes, kpi.TagType)
		names = append(names, kpi.KpiName)
		typeids = append(typeids, kpi.TypeId)
		descs = append(descs, kpi.Discription)
		shifthour = append(shifthour, kpi.ShiftHour)
		basetime = append(basetime, kpi.BaseTime)
	}
	return ids, tagids, typeids, shifthour, tagtypes, names, descs, basetime, err
}

/*******************************************************************************
功能:获取指定TagType类型的KPI的周期统计信息
输入:无
输出:[map[int]int] key为周期,value为该周期下的KPI指标数量
	[error] 错误信息
说明:
编辑:wang_jp
时间:2020年7月14日
*******************************************************************************/
func (cfg *CalcKpiConfigList) GetKpiTagTypePeriodInfo(tagtype string) (map[int]int, error) {
	o := orm.NewOrm()  //新建orm对象
	o.Using("default") //根据别名选择数据库
	sql := `SELECT
		period,
		COUNT( period ) AS cnt 
	FROM
		calc_kpi_config_list 
	WHERE
		tag_type = ? 
		AND status = 1 
	GROUP BY
		period;`
	type list struct {
		Period int
		Cnt    int
	}
	var lsts []list
	pmp := make(map[int]int)
	_, err := o.Raw(sql, tagtype).QueryRows(&lsts)
	for _, kpi := range lsts {
		pmp[kpi.Period] = kpi.Cnt
	}
	return pmp, err
}

/*******************************************************************************
功能:获取KPI指标配置表中的Tagtype类型列表
输入:无
输出:[map[string]int] key为TagType,value为该Type下的KPI指标数量
	[error] 错误信息
说明:
编辑:wang_jp
时间:2020年7月14日
*******************************************************************************/
func (cfg *CalcKpiConfigList) GetTagTypeLists() (map[string]int, error) {
	o := orm.NewOrm()  //新建orm对象
	o.Using("default") //根据别名选择数据库
	sql := `SELECT
		tag_type,
		COUNT( tag_type ) AS cnt 
	FROM
		calc_kpi_config_list 
	WHERE
		status = 1 
	GROUP BY
		tag_type;`
	type list struct {
		TagType string
		Cnt     int
	}
	var lsts []list
	pmp := make(map[string]int)
	_, err := o.Raw(sql).QueryRows(&lsts)
	for _, kpi := range lsts {
		pmp[kpi.TagType] = kpi.Cnt
	}
	return pmp, err
}

/*******************************************************************************
功能:获取KPI指标的配置信息
输入:无
输出:[error] 错误信息
说明:
编辑:wang_jp
时间:2020年7月14日
*******************************************************************************/
func (cfg *CalcKpiConfigList) GetKpiConfig() error {
	o := orm.NewOrm()  //新建orm对象
	o.Using("default") //根据别名选择数据库
	err := o.QueryTable("CalcKpiConfigList").Filter("Id", cfg.Id).RelatedSel().One(cfg)
	if err != nil {
		if strings.Contains(err.Error(), "no row found") {
			err = o.QueryTable("CalcKpiConfigList").Filter("Id", cfg.Id).One(cfg)
		}
	}
	return err
}

/*******************************************************************************
功能:保存KPI指标最后计算时间
输入:kpi指标结构体
输出:保存的数据数量
说明:
编辑:wang_jp
时间:2019年12月14日
*******************************************************************************/
func (kpicfg *CalcKpiConfigList) SetKpiLastCalcTime(kpis []CalcKpiResult) int64 {
	o := orm.NewOrm()  //新建orm对象
	o.Using("default") //根据别名选择数据库
	var k CalcKpiConfigList
	var num int64
	for _, kpi := range kpis {
		k.Id = kpi.KpiConfigListId
		k.LastCalcTime = kpi.CalcEndingTime
		i, err := o.Update(&k, "LastCalcTime") //执行查询
		if err != nil {                        //发生错误
			logs.Alert(`更新KPI指标 LastCalcTime 时发生错误,KPI指标ID是[%d];错误信息:[%s]`, k.Id, err.Error())
		}
		num += i
	}
	return num
}

/*******************************************************************************
功能:保存KPI的最后计算时间
输入:[id int64] KPI配置的ID
	[lasttime string] 最后计算时间
输出:无
说明:
编辑:wang_jp
时间:2020年7月14日
*******************************************************************************/
func (kpi *CalcKpiConfigListExi) SetKpiLastCalcTimeSingle(id int64, lasttime string) {
	o := orm.NewOrm()  //新建orm对象
	o.Using("default") //根据别名选择数据库
	var k CalcKpiConfigList
	k.Id = id
	k.LastCalcTime = lasttime

	_, err := o.Update(&k, "LastCalcTime") //执行查询
	if err != nil {                        //发生错误
		logs.Alert(`更新KPI指标 LastCalcTime 时发生错误,KPI指标ID是[%d];错误信息:[%s]`, k.Id, err.Error())
	}
}
