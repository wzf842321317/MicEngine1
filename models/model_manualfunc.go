package models

import (
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

/************************************************************
功能:根据表MineDDictionaryOfVaribleset中的calc_kpi_id设置,
	添加MineDDictionaryOfVaribleset与CalcKpiIndexDic的多对多关系
	使用的时候手动调用
说明:使用此函数需要先在MineDDictionaryOfVaribleset的calc_kpi_id列中设置对应的kpi_id
	不建议使用该功能，而应在页面上添加多对多关系
时间:2019年12月11日
************************************************************/
func (varset *MineDDictionaryOfVaribleset) InsertKpiDicVariableSetM2M() {
	o := orm.NewOrm() //新建orm对象
	o.Using("default")

	//varset := new(MineDDictionaryOfVaribleset)
	var vset []*MineDDictionaryOfVaribleset
	qs := o.QueryTable(varset)
	//在MineDDictionaryOfVaribleset表中查询calc_kpi_id列的设置,该列设置的格式为(1,3,6)等
	if _, err := qs.Filter("calc_kpi_id__isnull", false).All(&vset, "Id", "calc_kpi_id"); err == nil {
		reg, _ := regexp.Compile(`\d+`) //设置取数字的正则表达式
		for _, v := range vset {        //遍历每个Varibleset
			m2m := o.QueryM2M(v, "KpiIndexDics")      //设置与kpiIndexDics表的多对多关系
			ids := reg.FindAllString(v.CalcKpiId, -1) //提取calc_kpi_id列中设置的每个kpi_dic_id
			for _, s := range ids {                   //遍历每个id，添加多对多关系
				if i, err := strconv.ParseInt(s, 0, 64); err == nil {
					kpi := &CalcKpiIndexDic{Id: i}
					if m2m.Exist(kpi) == false {
						m2m.Add(kpi)
					}
				}
			}
		}
	}
}

/*************************************************
功能:在CalcKpiConfigList中自动添加需要计算的KPI指标
	根据表MineDDictionaryOfVaribleset与CalcKpiIndexDic的多对多关系
	以及tag的varibleSet属性在CalcKpiConfigList中添加tag的对应kpi指标
输入:startid:ore_process_d_taglist中的id.程序将以该id为起始id，为该id之后的tags
	添加kpi指标
	kpi_period:需要添加的指标的计算周期,-1:小时;-2:班;-3:日;-4:月;-5:季度;-6:年；
	如果大于0则单位为秒。
	starttime:kpi指标开始计算的时间
输出:添加的kpi的数量,错误信息
说明:该程序在系统每次运行的时候调用检查是否有KPI指标，没有的时候则自动创建KPI指标
编辑:wang_jp
时间:2019年12月11日
*************************************************/
func (lst *CalcKpiConfigList) InsertTagKpi2CfgList() int64 {
	o := orm.NewOrm() //新建orm对象
	o.Using("default")
	qs := o.QueryTable("calc_kpi_config_list")
	rows, err := qs.Filter("id__gt", 0).Count()
	if err != nil {
		logs.Alert("启动时查询表[calc_kpi_config_list]失败")
	}
	if rows == 0 { //如果表中没有数据
		sql := `SELECT
			lst.id as tag_id,
			lst.tag_description as tag_desc,
			lst.tage_name as tag_name,
			kpi.id as kpi_id,
			kpi.kpi_varible_name as kpi_name,
			kpi.kpi_varible_name_cn as kpi_cname
		FROM
			ore_process_d_taglist lst
			JOIN mine_d_dictionary_of_varibleset vst ON ( lst.variable_id = vst.id )
			JOIN calc_kpi_index_dic_mine_d_dictionary_of_variblesets rel ON vst.id = rel.mine_d_dictionary_of_varibleset_id
			RIGHT JOIN calc_kpi_index_dic kpi ON rel.calc_kpi_index_dic_id = kpi.id 
		WHERE lst.id is not null`
		type msg struct {
			TagId    int64
			TagDesc  string
			TagName  string
			KpiId    int64
			KpiName  string
			KpiCname string
		}
		var res []msg
		o.Raw(sql).QueryRows(&res) //查询

		period := -1
		var periodname string
		switch period {
		case -1:
			periodname = "hour"
		case -2:
			periodname = "shift"
		case -3:
			periodname = "day"
		case -4:
			periodname = "month"
		case -5:
			periodname = "season"
		case -6:
			periodname = "year"
		default:
			periodname = fmt.Sprintf("p%dsec", period)
		}

		var kpis []*CalcKpiConfigList
		for i, r := range res {
			kpidic := new(CalcKpiIndexDic)
			kpi := new(CalcKpiConfigList) //新建kpi指标
			kpidic.Id = r.KpiId
			kpi.Id = int64(i + 1)
			kpi.DistributedId = 0
			kpi.TagType = "process"
			kpi.TagId = r.TagId
			kpi.CalcKpiIndexDic = kpidic
			kpi.KpiTag = fmt.Sprintf("%s__%s_%s", r.TagName, periodname, r.KpiName)
			kpi.KpiName = fmt.Sprintf("%s__%s", r.TagDesc, r.KpiCname)
			kpi.StartTime = _TIMEFOMAT
			kpi.Period = int64(period)
			kpi.OffsetMinutes = 0
			kpi.LastCalcTime = TimeFormat(time.Now())
			kpi.Status = 1
			kpi.CreateTime = TimeFormat(time.Now())
			kpi.UpdateTime = TimeFormat(time.Now())

			kpis = append(kpis, kpi)
			//logger.Debug(r.TagId, r.KpiId, tag.Id, kpidic.Id)
		}

		qs.Filter("id__gt", 0).Delete()

		cnt, err := o.InsertMulti(1, kpis)
		if err != nil {
			logs.Error("创建KPI记录失败:[%s]", err.Error())
		}
		logs.Info("共创建了%d条KPI记录", cnt)
		return cnt
	}
	return 0
}

/*************************************************
功能:以追加的方式在CalcKpiConfigList中自动添加需要计算的KPI指标
	根据表MineDDictionaryOfVaribleset与CalcKpiIndexDic的多对多关系
	以及tag的varibleSet属性在CalcKpiConfigList中添加tag的对应kpi指标
输入:startid:ore_process_d_taglist中的id.程序将以该id为起始id，为该id之后的tags
	添加kpi指标
	kpi_period:需要添加的指标的计算周期,-1:小时;-2:班;-3:日;-4:月;-5:季度;-6:年；
	如果大于0则单位为秒。
	starttime:kpi指标开始计算的时间
输出:添加的kpi的数量,错误信息
说明:该程序手动执行
编辑:wang_jp
时间:2019年12月11日
*************************************************/
func (lst *CalcKpiConfigList) InsertTagKpi2CfgListByAppend(startid, kpi_period int, starttime ...string) (int64, error) {
	o := orm.NewOrm() //新建orm对象
	o.Using("default")
	var lastkpi []*CalcKpiConfigList
	qs := o.QueryTable("calc_kpi_config_list")
	_, err := qs.Filter("id__gt", 0).OrderBy("-id").Limit(1).All(&lastkpi)
	if err != nil {
		return 0, err
	}
	var baseid int64
	if len(lastkpi) > 0 {
		baseid = lastkpi[0].Id + 1
	} else {
		baseid = 1
	}
	//fmt.Printf("当前kpi_config_list中的最大ID是[%d]\n", baseid-1)

	sql := fmt.Sprintf(`SELECT
			lst.id as tag_id,
			lst.tag_description as tag_desc,
			lst.tage_name as tag_name,
			kpi.id as kpi_id,
			kpi.kpi_varible_name as kpi_name,
			kpi.kpi_varible_name_cn as kpi_cname
		FROM
			ore_process_d_taglist lst
			JOIN mine_d_dictionary_of_varibleset vst ON ( lst.variable_id = vst.id )
			JOIN calc_kpi_index_dic_mine_d_dictionary_of_variblesets rel ON vst.id = rel.mine_d_dictionary_of_varibleset_id
			RIGHT JOIN calc_kpi_index_dic kpi ON rel.calc_kpi_index_dic_id = kpi.id 
		WHERE lst.id > %d AND lst.stage_id is not null`, startid)
	type msg struct {
		TagId    int64
		TagDesc  string
		TagName  string
		KpiId    int64
		KpiName  string
		KpiCname string
	}
	var res []msg
	rows, err := o.Raw(sql).QueryRows(&res) //查询
	if err != nil {
		return rows, err
	}
	fmt.Printf("查询到了[%d]行待插入的数据\n", rows)
	if rows > 0 {
		period := kpi_period
		var periodname string
		switch period {
		case -1:
			periodname = "hour"
		case -2:
			periodname = "shift"
		case -3:
			periodname = "day"
		case -4:
			periodname = "month"
		case -5:
			periodname = "season"
		case -6:
			periodname = "year"
		default:
			periodname = fmt.Sprintf("p%dsec", period)
		}
		//开始时间默认为当前时间之前的1天
		sttime := TimeFormat(time.Now().Add(-3600*24*time.Second), _TIMEFOMAT)
		if len(starttime) > 0 {
			sttime = starttime[0]
		}
		var kpis []*CalcKpiConfigList
		for i, r := range res {
			kpidic := new(CalcKpiIndexDic)
			kpi := new(CalcKpiConfigList) //新建kpi指标
			kpidic.Id = r.KpiId
			kpi.Id = int64(i) + baseid
			kpi.DistributedId = 0
			kpi.TagType = "process"
			kpi.TagId = r.TagId
			kpi.CalcKpiIndexDic = kpidic
			kpi.KpiTag = fmt.Sprintf("%s__%s_%s", r.TagName, periodname, r.KpiName)
			kpi.KpiName = fmt.Sprintf("%s__%s", r.TagDesc, r.KpiCname)
			kpi.StartTime = _TIMEFOMAT
			kpi.Period = int64(period)
			kpi.OffsetMinutes = 0
			kpi.LastCalcTime = sttime
			kpi.Status = 1
			kpi.CreateTime = TimeFormat(time.Now())
			kpi.UpdateTime = TimeFormat(time.Now())

			kpis = append(kpis, kpi)
		}

		cnt, err := o.InsertMulti(1, kpis)
		if err != nil {
			return 0, err
		}
		return cnt, nil
	} else {
		return 0, fmt.Errorf("没有查询到符合条件的数据,查询语句是:%s", sql)
	}
}

/*************************************************
功能:检查Web菜单是否存在，不存在就创建
输入:无
输出:无
说明:程序启动时执行
编辑:wang_jp
时间:2020年3月27日
*************************************************/
func (menu *SysMenu) CheckWebMenu() {
	o := orm.NewOrm() //新建orm对象
	o.Using("default")
	thistime := TimeFormat(time.Now(), _TIMEFOMAT)
	menus := []SysMenu{ //PID要填写父节点的顺序(从1开始)
		{Remark: "MicEngineMenu_1", Name: "本地监控Web", NameEng: "LocalMonitor", Pid: 0, Level: 0, Seq: 1, Url: "/no", HasChild: 1, MenuType: 2, IconUrl: "", CreateTime: thistime, State: 1},
		{Remark: "MicEngineMenu_2", Name: "监控", NameEng: "Monitor", Pid: 1, Level: 1, Seq: 1, Url: "/monitor", HasChild: 0, MenuType: 2, IconUrl: "../static/img/menuico/trend.svg", CreateTime: thistime, State: 1},
		{Remark: "MicEngineMenu_3", Name: "数据", NameEng: "Data", Pid: 1, Level: 1, Seq: 2, Url: "/no", HasChild: 1, MenuType: 2, IconUrl: "../static/img/menuico/data.svg", CreateTime: thistime, State: 1},
		{Remark: "MicEngineMenu_4", Name: "快照", NameEng: "SnapShot", Pid: 3, Level: 2, Seq: 1, Url: "/snapshot", HasChild: 0, MenuType: 2, IconUrl: "../static/img/menuico/photo.svg", CreateTime: thistime, State: 1},
		{Remark: "MicEngineMenu_5", Name: "历史", NameEng: "History", Pid: 3, Level: 2, Seq: 2, Url: "/history", HasChild: 0, MenuType: 2, IconUrl: "../static/img/menuico/history.svg", CreateTime: thistime, State: 1},
		{Remark: "MicEngineMenu_6", Name: "分析", NameEng: "Analyse", Pid: 1, Level: 1, Seq: 3, Url: "", HasChild: 1, MenuType: 2, IconUrl: "../static/img/menuico/analys.svg", CreateTime: thistime, State: 1},
		{Remark: "MicEngineMenu_7", Name: "对比", NameEng: "Contrast", Pid: 6, Level: 2, Seq: 1, Url: "/compare", HasChild: 0, MenuType: 2, IconUrl: "../static/img/menuico/vs.svg", CreateTime: thistime, State: 1},
		{Remark: "MicEngineMenu_8", Name: "回归", NameEng: "Regression", Pid: 6, Level: 2, Seq: 2, Url: "/regression", HasChild: 0, MenuType: 2, IconUrl: "../static/img/menuico/huigui1.svg", CreateTime: thistime, State: 1},
		{Remark: "MicEngineMenu_9", Name: "管理", NameEng: "Manager", Pid: 1, Level: 1, Seq: 20, Url: "/no", HasChild: 1, MenuType: 2, IconUrl: "../static/img/menuico/manager.svg", CreateTime: thistime, State: 1},
		{Remark: "MicEngineMenu_10", Name: "项目管理", NameEng: "Project", Pid: 9, Level: 2, Seq: 1, Url: "/managerproject", HasChild: 0, MenuType: 2, IconUrl: "", CreateTime: thistime, State: 1},
		{Remark: "MicEngineMenu_11", Name: "权限管理", NameEng: "Permission", Pid: 9, Level: 2, Seq: 2, Url: "/managerpermission", HasChild: 0, MenuType: 2, IconUrl: "", CreateTime: thistime, State: 1},
		{Remark: "MicEngineMenu_12", Name: "项目日志", NameEng: "Logs", Pid: 9, Level: 2, Seq: 3, Url: "/managerlog", HasChild: 0, MenuType: 2, IconUrl: "", CreateTime: thistime, State: 1},
		{Remark: "MicEngineMenu_13", Name: "用户管理", NameEng: "User", Pid: 9, Level: 2, Seq: 4, Url: "/managerusers", HasChild: 0, MenuType: 2, IconUrl: "", CreateTime: thistime, State: 1},
		{Remark: "MicEngineMenu_14", Name: "检化", NameEng: "Laboratorial", Pid: 1, Level: 1, Seq: 4, Url: "/samplelab", HasChild: 0, MenuType: 2, IconUrl: "../static/img/menuico/lab.svg", CreateTime: thistime, State: 1},
		{Remark: "MicEngineMenu_15", Name: "物耗", NameEng: "Goods", Pid: 1, Level: 1, Seq: 5, Url: "/goods", HasChild: 0, MenuType: 2, IconUrl: "../static/img/menuico/goods.svg", CreateTime: thistime, State: 1},
		{Remark: "MicEngineMenu_16", Name: "巡检", NameEng: "Patrol", Pid: 1, Level: 1, Seq: 6, Url: "/patrol", HasChild: 0, MenuType: 2, IconUrl: "../static/img/menuico/check.svg", CreateTime: thistime, State: 1},
		{Remark: "MicEngineMenu_17", Name: "KPI", NameEng: "KPI", Pid: 1, Level: 1, Seq: 7, Url: "/kpi", HasChild: 0, MenuType: 2, IconUrl: "../static/img/menuico/kpi.svg", CreateTime: thistime, State: 1},
		{Remark: "MicEngineMenu_18", Name: "报表", NameEng: "Report", Pid: 1, Level: 1, Seq: 8, Url: "/report", HasChild: 1, MenuType: 2, IconUrl: "../static/img/menuico/report.svg", CreateTime: thistime, State: 1},
		{Remark: "MicEngineMenu_19", Name: "查看报表", NameEng: "View", Pid: 18, Level: 2, Seq: 1, Url: "/report", HasChild: 0, MenuType: 2, IconUrl: "../static/img/menuico/report.svg", CreateTime: thistime, State: 1},
		{Remark: "MicEngineMenu_20", Name: "管理报表", NameEng: "Edit", Pid: 18, Level: 2, Seq: 2, Url: "/reportedit", HasChild: 0, MenuType: 2, IconUrl: "../static/img/menuico/manager.svg", CreateTime: thistime, State: 1},
	}
	ids := make([]int64, len(menus)) //存放id切片
	for i, menu := range menus {
		if i > 0 { //不是第一条根节点
			menu.Pid = ids[menu.Pid-1] //获取父ID在数据库中的值
		}
		_, id, err := o.ReadOrCreate(&menu, "Remark", "MenuType") //读取或者创建
		if err == nil {
			ids[i] = id //获取ID
		}
	}
}
