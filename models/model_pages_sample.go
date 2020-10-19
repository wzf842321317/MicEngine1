package models

import (
	"fmt"

	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

/*************************************************
功能:通过用户ID获取用户授权的质检节点
输入:userid
输出:[]*MsMiddleCorrelation, error
说明:
编辑:wang_jp
时间:2020年4月8日
*************************************************/
func (u *SysUser) GetSampleLabTreeNodesByUserId(userid ...int64) ([]*SampleLabTree, error) {
	var nodes []*MsMiddleCorrelation
	if len(userid) > 0 {
		u.Id = userid[0]
	}
	for i := 1; i <= 4; i++ { //查询质检样本所关联的树节点
		md := new(MsMiddleCorrelation)
		nds, err := md.GetSampleLabNodesByLevel(i)
		if err != nil {
			return nil, err
		}
		nodes = append(nodes, nds...)
	}

	//查询用户被授权的树节点
	user_nodes, err := u.GetTreeNodesByUserId(false)
	if err != nil {
		return nil, err
	}
	var sample_nodes []*MsMiddleCorrelation
	for _, node := range nodes { //遍历样本节点
		for _, unode := range user_nodes { //遍历所有用户节点
			if node.Id == unode.Id { //寻找与样本节点相同的用户节点
				sample_nodes = append(sample_nodes, node) //找到后保存样本节点
				break
			}
		}
	}

	var tree []*SampleLabTree
	for _, mid := range sample_nodes {
		nd := new(SampleLabTree)
		nd.Id = mid.Id
		nd.Pid = mid.Pid
		nd.Name = mid.LevelName
		nd.NodeType = mid.LevelCategory.Id
		nd.ItemId = mid.ItemIdInTable
		nd.Seq = mid.Seq
		nd.TreeLevel = mid.TreeLevelCode
		nd.BaseTime = mid.ConstrutionTableCode
		nd.ShiftHour = mid.ConstrutionCode
		tree = append(tree, nd)
	}

	o := orm.NewOrm()
	o.Using("default") //选定数据库
	var samples []*SamplingManage
	qt := o.QueryTable("SamplingManage")
	_, err = qt.Filter("Id__gt", 0).RelatedSel("SampleFunction").All(&samples) //读取质检样本管理数据
	if err != nil {
		return nil, err
	}

	for _, sample := range samples { //遍历所有样本
		for _, node := range sample_nodes { //在树节点中查找样本的父节点
			if node.LevelCategory.Id == 4 && int64(sample.Stage.Id) == node.ItemIdInTable { //节点类型为作业,且样本父节点id与作业表ID相等
				nd := new(SampleLabTree)
				nd.Id = sample.Id * 100000
				nd.Pid = node.Id
				nd.Name = sample.SampleName
				nd.NodeType = 9999
				nd.ItemId = sample.SampleType.Id
				nd.Seq = sample.Seq
				nd.TreeLevel = sample.TreeLevelCode
				nd.SamplingSite = sample.SamplingSite
				nd.Func = sample.SampleFunction.Id
				nd.IsLeaf = true
				nd.IsRegular = sample.IsRegular
				nd.FuncName = sample.SampleFunction.Name
				nd.BaseTime = node.ConstrutionTableCode
				nd.ShiftHour = node.ConstrutionCode
				tree = append(tree, nd)
				break
			}
		}
	}
	return tree, nil
}

/*************************************************
功能:获取质检数据节点树中的分类节点信息
输入:level:
	矿山=1;厂=2;车间=3;作业=4
输出:[]MsMiddleCorrelation, error
说明:
编辑:wang_jp
时间:2020年4月8日
*************************************************/
func (md *MsMiddleCorrelation) GetSampleLabNodesByLevel(level int) ([]*MsMiddleCorrelation, error) {
	o := orm.NewOrm()
	o.Using("default") //选定数据库
	var nodes []*MsMiddleCorrelation
	sqlstr := `SELECT
	DISTINCT(middle.id) AS id,
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
	sampling_manage AS samp
	LEFT JOIN lab_sampletype_pool AS lab ON samp.sample_type_id = lab.id
	LEFT JOIN ore_process_d_workstage AS stage ON samp.sample_from_stage_id = stage.id
	LEFT JOIN ore_process_d_workshop AS workshop ON stage.workshop_id = workshop.id
	LEFT JOIN ore_process_concentration AS plant ON workshop.ore_process_concentration_id = plant.id
	LEFT JOIN mine_basic_info AS mine ON plant.mine_basic_info_id = mine.id
	LEFT JOIN ms_middle_correlation AS middle ON ( %s.id = middle.table_id AND middle.level_category = %d )`
	var code string
	switch level {
	case 1:
		code = "mine"
	case 2:
		code = "plant"
	case 3:
		code = "workshop"
	case 4:
		code = "stage"
	}
	sql := fmt.Sprintf(sqlstr, code, level)
	_, err := o.Raw(sql).QueryRows(&nodes)

	return nodes, err
}

/*************************************************
功能:通过样本模板ID获取样本模板的指标信息
输入:样本模板ID
输出:[]*SamplingManageSub,error
说明:
编辑:wang_jp
时间:2020年4月9日
*************************************************/
func (sub *SamplingManageSub) GetSampleSubBySampleId(sampleid int64) ([]*SamplingManageSub, error) {
	o := orm.NewOrm()
	o.Using("default") //选定数据库
	var rel []*SamplingManageSub
	qt := o.QueryTable("SamplingManageSub")
	if _, err := qt.Filter("Sample__Id", sampleid).RelatedSel().All(&rel); err == nil {
		return rel, nil
	} else {
		return nil, err
	}
}

/*************************************************
功能:通过样本模板ID获取样本指标的取送样及化验结果信息
输入:样本模板ID,开始时间,结束时间
输出:[]*LabAnaResultTsd,error
说明:
编辑:wang_jp
时间:2020年4月9日
*************************************************/
func (lb *LabAnaResultTsd) GetSampleLabResultBySampleId(sampleid int64, bgtime, endtime string) ([]*LabAnaResultTsd, error) {
	o := orm.NewOrm()
	o.Using("default") //选定数据库
	var rel []*LabAnaResultTsd
	qt := o.QueryTable("LabAnaResultTsd")
	if _, err := qt.Filter("SampleToLab__Sample__Id", sampleid).Filter("SampleToLab__SamplingTime__gt", bgtime).Filter("SampleToLab__SamplingTime__lte", endtime).RelatedSel().OrderBy("-SampleToLab__SamplingTime").All(&rel); err == nil {
		return rel, nil
	} else {
		return nil, err
	}
}
