package models

import (
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"strconv"
	"strings"
)

func (hb *OreProcessDTaglist) OreProcessTag() []string {
	o := orm.NewOrm() //新建orm对象
	o.Using("default")
	var hbs []OreProcessDTaglist
	qt := o.QueryTable("OreProcessDTaglist")
	qt.Filter("Id__gt", 0).All(&hbs)
	//读取数据到前台
	res := []string{}
	for _, data := range hbs {
		if strings.Contains(data.TagDescription, "车间") {
			id := strconv.Itoa(int(data.Id)) + "_" + data.TagDescription
			res = append(res, id)
		}
	}
	res = RemoveRepeatedElement(res)
	return res
}

func RemoveRepeatedElement(arr []string) (newArr []string) {
	newArr = make([]string, 0)
	for i := 0; i < len(arr); i++ {
		repeat := false
		for j := i + 1; j < len(arr); j++ {
			if arr[i] == arr[j] {
				repeat = true
				break
			}
		}
		if !repeat {
			newArr = append(newArr, arr[i])
		}
	}
	return
}

func (hb *KpiArithmeticResult) SaveDataLog(res *KpiArithmeticResult) {
	o := orm.NewOrm()
	o.Using("default")

	res.CrePerson = EngineCfgMsg.CfgMsg.Description
	_, err := o.Insert(res)
	if err != nil {
		logs.Warning("Insert Local Engine runtime error[插入失败]")
		logs.Warning(err)
	}
}

func (hb *KpiArithmetic) ReaderArithmeticData(ArithmeticName string) []string {
	o := orm.NewOrm() //新建orm对象
	o.Using("default")
	var hbs []KpiArithmetic
	res := []string{}
	qt := o.QueryTable("KpiArithmetic")
	qt.Filter("ArithmeticName2__iexact", ArithmeticName).All(&hbs)
	for _, data := range hbs {
		res = append(res, strconv.Itoa(int(data.Id)), data.ArithmeticName, data.ArithmeticResultType)
		return res
	}
	return res
}
