package models

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

/*************************************************
功能:通过用户ID获取用户授权的巡检节点
输入:userid
输出:[]*CheckTree, error
说明:
编辑:wang_jp
时间:2020年4月15日
*************************************************/
func (u *SysUser) GetCheckTreeNodesByUserId(userid ...int64) ([]*CheckTree, error) {
	var nodes map[int64]*MsMiddleCorrelation
	//作业及以上节点全选，作业以下层级的节点有物耗配置信息的才选
	md := new(MsMiddleCorrelation)
	nodes, err := md.GetCheckNodesByLevel()
	if len(userid) > 0 {
		u.Id = userid[0]
	}
	//查询用户被授权的树节点
	user_nodes, err := u.GetTreeNodesByUserId(false)
	if err != nil {
		return nil, err
	}
	//挑选出被授权的节点(level:1~4)
	var sample_nodes []*MsMiddleCorrelation
	for _, unode := range user_nodes { //遍历用户节点
		if node, ok := nodes[unode.Id]; ok == true {
			sample_nodes = append(sample_nodes, node) //找到后保存巡检节点
		}
	}
	//挑选出作业以下节点，并在被授权节点中寻找其父节点
	var sub_nodes []*MsMiddleCorrelation
	for _, node := range nodes { //遍历巡检节点
		if node.LevelCategory.Id > 4 { //挑选出作业以下的节点
			for _, unode := range sample_nodes { //遍历所有用户节点
				if node.Pid == unode.Id { //寻找与巡检节点相同的用户节点
					sub_nodes = append(sub_nodes, node) //找到后保存巡检节点
					break
				}
			}
		}
	}

	//对子节点排序
	flag := true
	vLen := len(sub_nodes)
	for i := 0; i < vLen-1; i++ {
		flag = true
		for j := 0; j < vLen-i-1; j++ {
			if sub_nodes[j].Seq > sub_nodes[j+1].Seq {
				sub_nodes[j], sub_nodes[j+1] = sub_nodes[j+1], sub_nodes[j]
				flag = false
				continue
			}
		}
		if flag {
			break
		}
	}
	sample_nodes = append(sample_nodes, sub_nodes...)

	o := orm.NewOrm()
	o.Using("default") //选定数据库
	var checks []*CheckSiteEquipRel
	qt := o.QueryTable("CheckSiteEquipRel")
	_, err = qt.Filter("Id__gt", 0).RelatedSel().OrderBy("Id").All(&checks) //读取物耗管理数据
	if err != nil {
		return nil, err
	}

	//提出各个物耗数据的上级层级ID
	ckser := new(CheckSiteEquipRel)
	cknodes, err := ckser.GetCheckSiteEquipLevelMap()
	if err != nil {
		return nil, err
	}

	var tree []*CheckTree
	for _, mid := range sample_nodes {
		key := mid.Id
		if _, ok := cknodes[key]; ok { //判断层级ID是否有物耗数据,有的话才添加
			nd := new(CheckTree)
			nd.Id = mid.Id
			nd.Pid = mid.Pid
			nd.Name = mid.LevelName
			nd.NodeType = mid.LevelCategory.Id
			nd.TreeLevel = mid.TreeLevelCode
			nd.BaseTime = mid.ConstrutionTableCode
			nd.ShiftHour = int(mid.ConstrutionCode)
			tree = append(tree, nd)
		}
	}

	for _, gd := range checks { //遍历所有巡检列表
		for _, node := range sample_nodes { //在树节点中查找样本的父节点
			if gd.ResourceType == node.LevelCategory.Id && gd.ResourceId == node.ItemIdInTable { //节点类型为作业,且样本父节点id与作业表ID相等
				key := node.Id
				nd := new(CheckTree)
				nd.Id = gd.Id * 100000
				nd.Pid = node.Id
				nd.Name = gd.CheckSite.Name
				nd.SiteId = gd.CheckSite.Id
				nd.NodeType = 9999
				nd.TreeLevel = node.TreeLevelCode
				nd.IsLeaf = true
				nd.BaseTime = nodes[key].ConstrutionTableCode
				nd.ShiftHour = int(nodes[key].ConstrutionCode)
				nd.LineId = gd.CheckSite.CheckLine.Id
				nd.LineName = gd.CheckSite.CheckLine.LineName
				nd.DeptId = gd.CheckSite.CheckLine.Dept.Id
				nd.DeptName = gd.CheckSite.CheckLine.Dept.DepartmentName
				tree = append(tree, nd)
				break
			}
		}
	}
	return tree, nil
}

/*************************************************
功能:获取巡检数据节点树中的节点信息
输入:无
输出:[]MsMiddleCorrelation, error
说明:作业及以上节点全选，作业以下层级的节点有物耗配置信息的才选。
编辑:wang_jp
时间:2020年4月15日
*************************************************/
func (md *MsMiddleCorrelation) GetCheckNodesByLevel() (map[int64]*MsMiddleCorrelation, error) {
	o := orm.NewOrm()
	o.Using("default") //选定数据库
	var nodes []*MsMiddleCorrelation
	var nds []*MsMiddleCorrelation
	sqlstr1 := `SELECT
			DISTINCT middle.*
		FROM
			ms_middle_correlation AS middle
		WHERE middle.level_category = %d
		ORDER BY middle.sortNum_in_level asc,middle.id asc`
	sqlstr2 := `SELECT
			DISTINCT middle.id,
			middle.pId,
			middle.level_category,
			middle.level_name,
			middle.level_num,
			middle.status,
			middle.sortNum_in_level,
			middle.table_id,
			middle.tree_level_code,
			workshop.base_time AS constrution_table_code,
			workshop.shift_hour AS constrution_code
		FROM
			ms_middle_correlation AS middle
			%s
		WHERE middle.level_category = %d
		ORDER BY middle.sortNum_in_level asc,middle.id asc`
	sqlstr3 := `SELECT
			DISTINCT middle.id,
			middle.pId,
			middle.level_category,
			middle.level_name,
			middle.level_num,
			middle.status,
			middle.sortNum_in_level,
			middle.table_id,
			middle.tree_level_code,
			workshop.base_time AS constrution_table_code,
			workshop.shift_hour AS constrution_code
		FROM
			check_site_equip_rel AS list
			LEFT JOIN ms_middle_correlation AS middle ON ( list.level_category = middle.level_category AND middle.table_id = list.resource_id )
			%s
		WHERE list.level_category = %d
		ORDER BY middle.sortNum_in_level asc,middle.id asc`
	for i := 1; i <= 9; i++ {
		var sql string
		switch i {
		case 1, 2: //矿、厂
			sql = fmt.Sprintf(sqlstr1, i)
			break
		case 3: //车间
			leftsql := `LEFT JOIN ore_process_d_workshop AS workshop ON middle.table_id = workshop.id`
			sql = fmt.Sprintf(sqlstr2, leftsql, i)
			break
		case 4: //作业
			leftsql := `LEFT JOIN ore_process_d_workstage AS stage ON middle.table_id = stage.id
					LEFT JOIN ore_process_d_workshop AS workshop ON stage.workshop_id = workshop.id`
			sql = fmt.Sprintf(sqlstr2, leftsql, i)
			break
		case 5: //设备
			leftsql := `LEFT JOIN ore_process_d_workstage_device AS device ON device.id = middle.table_id
					LEFT JOIN ore_process_d_workstage AS stage ON device.stage_id = stage.id
					LEFT JOIN ore_process_d_workshop AS workshop ON stage.workshop_id = workshop.id`
			sql = fmt.Sprintf(sqlstr3, leftsql, i)
			break
		case 6: //仪表
			leftsql := `LEFT JOIN ore_process_d_workstage_meter AS meter ON meter.id = middle.table_id
					LEFT JOIN ore_process_d_workstage AS stage ON meter.stage_id = stage.id
					LEFT JOIN ore_process_d_workshop AS workshop ON stage.workshop_id = workshop.id`
			sql = fmt.Sprintf(sqlstr3, leftsql, i)
			break
		case 7: //电机
			leftsql := `LEFT JOIN ore_process_d_workstage_motor AS motor ON motor.id = middle.table_id
					LEFT JOIN ore_process_d_workstage_device AS device ON motor.device_id = device.id
					LEFT JOIN ore_process_d_workstage AS stage ON device.stage_id = stage.id
					LEFT JOIN ore_process_d_workshop AS workshop ON stage.workshop_id = workshop.id`
			sql = fmt.Sprintf(sqlstr3, leftsql, i)
			break
		case 8: //分析仪
			leftsql := `LEFT JOIN ore_process_d_workstage_analyzer AS analyzer ON analyzer.id = middle.table_id
					LEFT JOIN ore_process_d_workstage AS stage ON analyzer.stage_id = stage.id
					LEFT JOIN ore_process_d_workshop AS workshop ON stage.workshop_id = workshop.id`
			sql = fmt.Sprintf(sqlstr3, leftsql, i)
			break
		case 9: //取样点
			leftsql := `LEFT JOIN ore_workstage_assay AS assay ON assay.id = middle.table_id
					LEFT JOIN ore_process_d_workstage AS stage ON assay.stage_id = stage.id
					LEFT JOIN ore_process_d_workshop AS workshop ON stage.workshop_id = workshop.id`
			sql = fmt.Sprintf(sqlstr3, leftsql, i)
			break
		default:
			sql = fmt.Sprintf(sqlstr1, i)
			break
		}
		_, err := o.Raw(sql).QueryRows(&nds)
		if err == nil {
			nodes = append(nodes, nds...)
		}
	}
	nodesmap := make(map[int64]*MsMiddleCorrelation)
	for _, node := range nodes {
		key := node.Id
		nodesmap[key] = node
	}
	return nodesmap, nil
}

/*************************************************
功能:通过tag ID获取该tag的执行数据
输入:tag ID,查询起始时间,查询结束时间,查询时间仅与开始时间匹配
输出:[]*CheckItemExe,error
说明:
编辑:wang_jp
时间:2020年4月15日
*************************************************/
func (ckexe *CheckItemExe) GetCheckResultsByTagId(tagid int64, bgtime, endtime string, startonly bool) ([]*CheckItemExe, error) {
	o := orm.NewOrm()
	o.Using("default") //选定数据库
	var rel []*CheckItemExe
	qt := o.QueryTable("CheckItemExe").Filter("CheckItem__SiteEquipRel__Id", tagid)

	if startonly == true { //起止时间与开始时间仅与开始时间匹配
		qt = qt.Filter("CheckSiteExe__CheckPlanExe__AllStartTime__gte", bgtime).Filter("CheckSiteExe__CheckPlanExe__AllStartTime__lte", endtime)
	} else { //起止时间与开始时间和结束时间匹配
		qt = qt.Filter("CheckSiteExe__CheckPlanExe__AllStartTime__gte", bgtime).Filter("CheckSiteExe__CheckPlanExe__AllEndTime__lte", endtime)
	}
	if _, err := qt.RelatedSel("CheckItem", "CheckSiteExe__CheckPlanExe").OrderBy("-Id").All(&rel); err == nil {
		return rel, nil
	} else {
		return nil, err
	}
}

/*************************************************
功能:通过层级码获该层级下的物资的消耗数据
输入:层级码,查询起始时间,查询结束时间,查询时间仅与开始时间匹配
输出:[]*CheckItemExe,error
说明:
编辑:wang_jp
时间:2020年4月15日
*************************************************/
func (ckexe *CheckItemExe) GetCheckResultsByLevelCode(levelcode, bgtime, endtime string, startonly bool) ([]*CheckItemExe, error) {
	o := orm.NewOrm()
	o.Using("default") //选定数据库
	sqlstr := `SELECT
			T0.*
		FROM
			check_item_exe AS T0
			LEFT JOIN ms_middle_correlation AS T1 ON ( T0.resource_type = T1.level_category AND T1.table_id = T0.resource_id )
			LEFT JOIN check_site_exe AS T2 ON T2.id = T0.check_site_exe_id
			LEFT JOIN check_plan_exe AS T3 ON T3.id = T2.check_plan_exe_id
		WHERE T1.tree_level_code LIKE "%s" %s`
	var rel []*CheckItemExe

	var filter string
	if startonly == true { //起止时间与开始时间仅与开始时间匹配
		filter = fmt.Sprintf(`AND T3.all_start_time >= "%s" AND T3.all_start_time <="%s" ORDER BY T3.all_start_time DESC`, bgtime, endtime)
	} else { //起止时间与开始时间和结束时间匹配
		filter = fmt.Sprintf(`AND T3.all_start_time >= "%s" AND T3.all_end_time <="%s" ORDER BY T3.all_start_time DESC`, bgtime, endtime)
	}
	sql := fmt.Sprintf(sqlstr, levelcode+"%", filter)
	if _, err := o.Raw(sql).QueryRows(&rel); err == nil {
		return rel, nil
	} else {
		return nil, err
	}
}

/*************************************************
功能:通过层级码获该层级下的巡检点信息
输入:层级码,查询起始时间,查询结束时间,查询时间仅与开始时间匹配
输出:[]*CheckItemExe,error
说明:
编辑:wang_jp
时间:2020年4月15日
*************************************************/
func (ck *CheckTagList) GetCheckTagListByLevelCode(levelcode string) ([]*CheckTagList, error) {
	o := orm.NewOrm()
	o.Using("default") //选定数据库
	sqlstr := `SELECT
			T0.*
		FROM
			check_tag_list T0
			LEFT JOIN ms_middle_correlation T1 ON ( T0.resource_type = T1.level_category AND T1.table_id = T0.resource_id )
			LEFT JOIN sys_unit T2 ON T2.id = T0.unit_id
			LEFT JOIN check_variable_set T3 ON T3.id = T0.variable_id 
			LEFT JOIN ore_process_equipment_dic T4 ON T4.id = T3.type
			LEFT JOIN check_type_conf T5 ON T5.check_type_id = T4.id
			LEFT JOIN check_site_equip_rel T6 ON T6.id = T0.check_site_equip_exe_id
			LEFT JOIN check_site T7 ON T7.id = T6.check_site_id
			LEFT JOIN check_line T8 ON T8.id = T7.check_line_id  
			LEFT JOIN mine_dept_info T9 ON T9.id = T8.dept_id    
		WHERE T1.tree_level_code LIKE "%s"`
	var rel []*CheckTagList

	sql := fmt.Sprintf(sqlstr, levelcode+"%")
	if _, err := o.Raw(sql).QueryRows(&rel); err == nil {
		return rel, nil
	} else {
		return nil, err
	}
}

/*************************************************
功能:通过层级码获该层级下的巡检点信息
输入:层级码,查询起始时间,查询结束时间,查询时间仅与开始时间匹配
输出:[]*CheckItemExe,error
说明:
编辑:wang_jp
时间:2020年4月15日
*************************************************/
func (ck *CheckSite) GetCheckSiteByLevelCode(levelcode, bgtime, endtime string) (interface{}, error) {
	o := orm.NewOrm()
	type sitemsg struct {
		SiteEquipId  int64
		EquipName    string
		SiteId       int64
		SiteName     string
		LineId       int64
		LineName     string
		DeptId       int64
		DeptName     string
		CheckItemCnt int64
		TotalItemCnt int64
	}
	o.Using("default") //选定数据库
	/*	sqlstr := `SELECT
			DISTINCT T2.id AS site_equip_id
			,T1.level_name AS equip_name
			,T3.id AS site_id
			,T3.name AS site_name
			,T4.id AS line_id
			,T4.line_name
			,T5.id AS dept_id
			,T5.department_name AS dept_name
		FROM
			check_tag_list T0
			LEFT JOIN ms_middle_correlation T1 ON ( T0.resource_type = T1.level_category AND T1.table_id = T0.resource_id )
			LEFT JOIN check_site_equip_rel T2 ON T2.id = T0.check_site_equip_exe_id
			LEFT JOIN check_site T3 ON T3.id = T2.check_site_id
			LEFT JOIN check_line T4 ON T4.id = T3.check_line_id
			LEFT JOIN mine_dept_info T5 ON T5.id = T4.dept_id
		WHERE T1.tree_level_code LIKE "%s" AND T0.check_site_equip_exe_id is not null ORDER BY check_name`
	*/
	sqlstr := `SELECT
			T2.id AS site_equip_id
			,T1.level_name AS equip_name
			,T3.id AS site_id
			,T3.name AS site_name
			,T4.id AS line_id
			,T4.line_name
			,T5.id AS dept_id
			,T5.department_name AS dept_name
			,SUM(T6.check_status) AS check_item_cnt
			,COUNT(T6.id) AS total_item_cnt
		FROM
			check_tag_list T0
			LEFT JOIN ms_middle_correlation T1 ON ( T0.resource_type = T1.level_category AND T1.table_id = T0.resource_id )
			LEFT JOIN check_site_equip_rel T2 ON T2.id = T0.check_site_equip_exe_id
			LEFT JOIN check_site T3 ON T3.id = T2.check_site_id
			LEFT JOIN check_line T4 ON T4.id = T3.check_line_id  
			LEFT JOIN mine_dept_info T5 ON T5.id = T4.dept_id  
			INNER JOIN check_item_exe T6 ON T6.check_item_id = T0.id
			LEFT JOIN check_plan_exe T7 ON T7.id = T6.check_plan_exe_id
		WHERE T1.tree_level_code LIKE "%s" AND T7.all_start_time >= "%s" AND T7.all_start_time < "%s" AND T0.check_site_equip_exe_id is not null  GROUP BY T2.id  ORDER BY T3.id`
	var rel []sitemsg

	sql := fmt.Sprintf(sqlstr, levelcode+"%", bgtime, endtime)
	if _, err := o.Raw(sql).QueryRows(&rel); err == nil {
		return rel, nil
	} else {
		return nil, err
	}
}

/*************************************************
功能:通过SiteEquipID获取该TagList的基本数据
输入:SiteEquipID
输出:[]*CheckTagList,error
说明:
编辑:wang_jp
时间:2020年4月15日
*************************************************/
func (ck *CheckTagList) GetCheckTagListsBySiteId(tagid int64) ([]*CheckTagList, error) {
	o := orm.NewOrm()
	o.Using("default") //选定数据库
	var rel []*CheckTagList
	qt := o.QueryTable("CheckTagList").Filter("SiteEquipRel__Id", tagid)

	if _, err := qt.RelatedSel().All(&rel); err == nil {
		return rel, nil
	} else {
		return nil, err
	}
}

/*************************************************
功能:通过与taglist相关联的中间表中的levelcode获取与taglist相关的所有层级map
输入:无
输出:map[int64]int64,error
说明:
编辑:wang_jp
时间:2020年4月15日
*************************************************/
func (ck *CheckSiteEquipRel) GetCheckSiteEquipLevelMap() (map[int64]int64, error) {
	o := orm.NewOrm()
	o.Using("default") //选定数据库
	sql := `SELECT
			middle.tree_level_code
		FROM
			check_site_equip_rel AS list
			LEFT JOIN ms_middle_correlation AS middle ON ( list.level_category = middle.level_category AND middle.table_id = list.resource_id )
		WHERE list.id > 0`
	var rel []string
	if _, err := o.Raw(sql).QueryRows(&rel); err == nil {
		nodesidmap := make(map[int64]int64)
		for _, lc := range rel { //遍历所有物耗列表
			idstrs := strings.Split(lc, "-")
			for _, idstr := range idstrs {
				id, e := strconv.ParseInt(idstr, 10, 64)
				if e == nil {
					nodesidmap[id] = id
				}
			}
		}
		return nodesidmap, nil
	} else {
		return nil, err
	}
}

/*************************************************
功能:通过表名ID和表项目ID在img表中查找照片地址
输入:tableid, itemid
输出:MineCheckImg,error
说明:
编辑:wang_jp
时间:2020年4月17日
*************************************************/
func (img *MineCheckImg) GetMineCheckImg(tableid, itemid int64) (MineCheckImg, error) {
	o := orm.NewOrm()
	o.Using("default") //选定数据库
	rel := MineCheckImg{TableNameId: tableid, TableId: itemid}
	err := o.Read(&rel, "TableNameId", "TableId")
	return rel, err
}
