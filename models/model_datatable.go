package models

import (
	"github.com/astaxie/beego/orm"
)

/*******************************************************************************
功能:获取所有数据表列表
输入:无
输出:[[]DatatableInfo] 数据表列表
	[error] 错误信息
说明:
编辑:wang_jp
时间:2020年6月18日
*******************************************************************************/
func (tb *DatatableInfo) GetTableList() ([]DatatableInfo, error) {
	o := orm.NewOrm()
	o.Using("default") //选定数据库
	var tbs []DatatableInfo
	_, err := o.QueryTable("DatatableInfo").Filter("Id__gt", 0).All(&tbs)
	return tbs, err
}

/*******************************************************************************
功能:更新表关联的庚顿ID
输入:无
输出:[error] 错误信息
说明:需要预先设定表的Id和GoldenTableId
编辑:wang_jp
时间:2020年6月18日
*******************************************************************************/
func (tb *DatatableInfo) UpdateGoldenId() error {
	o := orm.NewOrm()
	o.Using("default") //选定数据库
	_, err := o.Update(tb, "GoldenTableId")
	return err
}

/*******************************************************************************
功能:根据数据表名获取数据表属性
输入:无
输出:[DatatableInfo] 数据表列表
	[error] 错误信息
说明:
编辑:wang_jp
时间:2020年6月18日
*******************************************************************************/
func (tb *DatatableInfo) GetTableByName(name ...string) error {
	o := orm.NewOrm()
	o.Using("default") //选定数据库
	if len(name) > 0 {
		tb.TableName = name[0]
	}
	err := o.QueryTable("DatatableInfo").Filter("TableName", tb.TableName).One(tb)
	return err
}

/*******************************************************************************
功能:通过数据表名更新表属性中的庚顿表ID
输入:[goldentableid] 庚顿数据表ID
	[name] 数据表名
输出:[error] 错误信息
说明:
编辑:wang_jp
时间:2020年6月18日
*******************************************************************************/
func (tb *DatatableInfo) UpdateGoldenIdByName(goldentableid int, name ...string) error {
	o := orm.NewOrm()
	o.Using("default") //选定数据库
	if len(name) > 0 {
		tb.TableName = name[0]
	}
	err := o.QueryTable("DatatableInfo").Filter("TableName", tb.TableName).One(tb)
	if err == nil {
		tb.GoldenTableId = goldentableid
		_, err = o.Update(tb, "GoldenTableId")
	}
	return err
}

/*******************************************************************************
功能:同步庚顿数据表和平台数据表
输入:无
输出:[int] 平台中的总数据表数量
	[int] 在庚顿中新建的数据表数量
	[error] 错误信息
说明:以平台数据表为基准同步庚顿数据表,如果庚顿中已经存在同名表,则将平台中表的描述信息同步到庚顿
	表中. 如果庚顿中没有存在同名表,则在庚顿数据库中新建同名表.
	无论新建还是更新,都会将庚顿表ID取回存入平台表对应字段
编辑:wang_jp
时间:2020年6月18日
*******************************************************************************/
func (tb *DatatableInfo) SynchGoldenTableAndPlatTable() (int, int, error) {
	tbs, err := tb.GetTableList() //获取平台数据表
	if err != nil {
		return 0, 0, err
	}
	newtbs := 0
	for _, t := range tbs { //遍历每个数据表
		micgd := new(MicGolden)
		isnew, id, err := micgd.GoldenTableInsertOrUpdate(t.TableName, t.TableDescription)
		if err == nil {
			if isnew { //新建计数
				newtbs += 1
			}
			t.GoldenTableId = id
			t.UpdateGoldenId() //保存到平台数据库
		}
	}
	return len(tbs), newtbs, nil
}
