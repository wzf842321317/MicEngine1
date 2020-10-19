package models

import (
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

/*
功能:通过物耗配置id和起止时间求物耗的和
输入:goodscfgid:物耗配置id,begineTime, endTime
输出:float64结果和error
说明:
编辑:wang_jp
时间:2019年12月26日
*/
func (gd *GoodsConsumeInfo) GetGoodsSumByID(goodscfgid int64, begineTime, endTime string) (float64, error) {
	o := orm.NewOrm() //新建orm对象
	o.Using("default")
	goods := new(GoodsConsumeInfo)
	var maps []orm.Params
	_, err := o.QueryTable(goods).Filter("GoodsConfigInfo__Id", goodscfgid).Filter("UseStartTime__gte", begineTime).Filter("UseStartTime__lte", endTime).Values(&maps)
	if err != nil {
		return 0.0, err
	}
	var sum float64
	for _, m := range maps {
		v, _ := m["GoodsConsumeAmount"].(float64)
		sum += v
	}
	return sum, nil
}

/*
功能:通过物耗配置id和时间点查询该时间点(或之前第一个)的物耗值
输入:goodscfgid:物耗配置id,timePoint
输出:float64结果和error
说明:
编辑:wang_jp
时间:2019年12月26日
*/
func (gd *GoodsConsumeInfo) GetGoodsValueByIDandTime(goodscfgid int64, timePoint string) (string, float64, error) {
	o := orm.NewOrm() //新建orm对象
	o.Using("default")
	goods := new(GoodsConsumeInfo)
	var maps []orm.Params
	_, err := o.QueryTable(goods).Filter("GoodsConfigInfo__Id", goodscfgid).Filter("UseStartTime__lte", timePoint).Limit(1).Values(&maps)
	if err != nil {
		return "", 0.0, err
	}
	var value float64
	var tstr string
	for _, m := range maps {
		value, _ = m["GoodsConsumeAmount"].(float64)
		tstr, _ = m["UseStartTime"].(string)
	}
	return tstr, value, nil
}

/*
功能:通过物耗名称获取物耗ID
输入:tagname:物耗名称
输出:id结果和error
说明:
编辑:wang_jp
时间:2019年12月26日
*/
func (gd *GoodsConfigInfo) GetGoodsTagIDByTagName(tagname string) (int64, error) {
	o := orm.NewOrm() //新建orm对象
	o.Using("default")
	goods := GoodsConfigInfo{GoodsTagName: tagname}
	err := o.Read(&goods, "GoodsTagName")
	if err != nil {
		return 0, err
	}
	return goods.Id, nil
}
