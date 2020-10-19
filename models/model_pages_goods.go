package models

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

/*************************************************
功能:通过用户ID获取用户授权的物耗节点
输入:userid
输出:[]*GoodsTree, error
说明:
编辑:wang_jp
时间:2020年4月11日
*************************************************/
func (u *SysUser) GetGoodsTreeNodesByUserId(userid ...int64) ([]*GoodsTree, error) {
	var nodes map[int64]*MsMiddleCorrelation
	//作业及以上节点全选，作业以下层级的节点有物耗配置信息的才选
	md := new(MsMiddleCorrelation)
	nodes, err := md.GetGoodsNodesByLevel()
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
	for _, unode := range user_nodes { //遍历物耗节点
		for _, node := range nodes { //遍历所有用户节点
			if node.Id == unode.Id { //寻找与物耗节点相同的用户节点
				sample_nodes = append(sample_nodes, node) //找到后保存物耗节点
				break
			}
		}
	}
	//挑选出作业以下节点，并在被授权节点中寻找其父节点
	var sub_nodes []*MsMiddleCorrelation
	for _, node := range nodes { //遍历物耗节点
		if node.LevelCategory.Id > 4 { //挑选出作业以下的节点
			for _, unode := range sample_nodes { //遍历所有用户节点
				if node.Pid == unode.Id { //寻找与物耗节点相同的用户节点
					sub_nodes = append(sub_nodes, node) //找到后保存物耗节点
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
	var goods []*GoodsConfigInfo
	qt := o.QueryTable("GoodsConfigInfo")
	_, err = qt.Filter("Id__gt", 0).RelatedSel("Goods", "Workshop").All(&goods) //读取物耗管理数据
	if err != nil {
		return nil, err
	}
	//提出各个物耗数据的上级层级ID
	gdnodes := make(map[int64]int64)
	for _, gd := range goods { //遍历所有物耗列表
		idstrs := strings.Split(gd.TreeLevelCode, "-")
		for _, idstr := range idstrs {
			id, e := strconv.ParseInt(idstr, 10, 64)
			if e == nil {
				gdnodes[id] = id
			}
		}
	}

	var tree []*GoodsTree
	for _, mid := range sample_nodes {
		if _, ok := gdnodes[mid.Id]; ok { //判断层级ID是否有物耗数据,有的话才添加
			nd := new(GoodsTree)
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

	for _, gd := range goods { //遍历所有物耗列表
		for _, node := range sample_nodes { //在树节点中查找样本的父节点
			if gd.ResourceType == node.LevelCategory.Id && gd.ResourceId == node.ItemIdInTable { //节点类型为作业,且样本父节点id与作业表ID相等
				nd := new(GoodsTree)
				nd.Id = gd.Id * 100000
				nd.Pid = node.Id
				nd.Name = gd.Goods.GoodsName
				nd.NodeType = 9999
				nd.TreeLevel = gd.TreeLevelCode
				nd.IsLeaf = true
				nd.BaseTime = gd.Workshop.BaseTime
				nd.ShiftHour = gd.Workshop.ShiftHour
				tree = append(tree, nd)
				break
			}
		}
	}
	return tree, nil
}

/*************************************************
功能:获取物耗数据节点树中的节点信息
输入:无
输出:[]MsMiddleCorrelation, error
说明:作业及以上节点全选，作业以下层级的节点有物耗配置信息的才选。
编辑:wang_jp
时间:2020年4月11日
*************************************************/
func (md *MsMiddleCorrelation) GetGoodsNodesByLevel() (map[int64]*MsMiddleCorrelation, error) {
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
			goods_config_info AS gds
			LEFT JOIN ms_middle_correlation AS middle ON ( gds.resource_type = middle.level_category AND middle.table_id = gds.resource_id )
			%s
		WHERE gds.resource_type = %d
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
		nodesmap[node.Id] = node
	}
	return nodesmap, nil
}

/*************************************************
功能:通过物耗ID获取该种物资的消耗数据
输入:物耗ID,查询起始时间,查询结束时间,查询时间仅与开始时间匹配
输出:[]*GoodsConsumeInfo,error
说明:
编辑:wang_jp
时间:2020年4月11日
*************************************************/
func (gd *GoodsConsumeInfo) GetGoodsResultsByGoodsId(goodscfgid int64, bgtime, endtime string, startonly bool) ([]*GoodsConsumeInfo, error) {
	o := orm.NewOrm()
	o.Using("default") //选定数据库
	var rel []*GoodsConsumeInfo
	qt := o.QueryTable("GoodsConsumeInfo").Filter("GoodsConfigInfo__Id", goodscfgid)

	if startonly == false { //起止时间与开始时间和结束时间匹配
		qt = qt.Filter("UseStartTime__gte", bgtime).Filter("UseEndTime__lte", endtime)
	} else { //起止时间与开始时间仅与开始时间匹配
		qt = qt.Filter("UseStartTime__gte", bgtime).Filter("UseStartTime__lte", endtime)
	}
	if _, err := qt.RelatedSel().OrderBy("-UseStartTime").All(&rel); err == nil {
		return rel, nil
	} else {
		return nil, err
	}
}

/*************************************************
功能:通过层级码获该层级下的物资的消耗数据
输入:层级码,查询起始时间,查询结束时间,查询时间仅与开始时间匹配
输出:[]*GoodsConsumeInfo,error
说明:
编辑:wang_jp
时间:2020年4月11日
*************************************************/
func (gd *GoodsConsumeInfo) GetGoodsResultsByLevelCode(levelcode, bgtime, endtime string, startonly bool) ([]*GoodsConsumeInfo, error) {
	o := orm.NewOrm()
	o.Using("default") //选定数据库
	var rel []*GoodsConsumeInfo
	qt := o.QueryTable("GoodsConsumeInfo").Filter("GoodsConfigInfo__TreeLevelCode__istartswith", levelcode)

	if startonly == false { //起止时间与开始时间和结束时间匹配
		qt = qt.Filter("UseStartTime__gte", bgtime).Filter("UseEndTime__lte", endtime)
	} else { //起止时间与开始时间仅与开始时间匹配
		qt = qt.Filter("UseStartTime__gte", bgtime).Filter("UseStartTime__lte", endtime)
	}
	if _, err := qt.RelatedSel().OrderBy("-UseStartTime").All(&rel); err == nil {
		return rel, nil
	} else {
		return nil, err
	}
}

/*************************************************
功能:通过物耗ID获取该种物资的基本数据
输入:物耗ID
输出:[]*GoodsConfigInfo,error
说明:
编辑:wang_jp
时间:2020年4月11日
*************************************************/
func (gd *GoodsConfigInfo) GetGoodsInfoByGoodsId(goodscfgid int64) ([]*GoodsConfigInfo, error) {
	o := orm.NewOrm()
	o.Using("default") //选定数据库
	var rel []*GoodsConfigInfo
	qt := o.QueryTable("GoodsConfigInfo").Filter("Id", goodscfgid)

	if _, err := qt.RelatedSel().All(&rel); err == nil {
		return rel, nil
	} else {
		return nil, err
	}

}
