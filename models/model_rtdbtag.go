package models

import (
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/bkzy-wangjp/goldengo"
	_ "github.com/go-sql-driver/mysql"
)

/*******************************************************************************
功能:通过tagid获取tag全名(含表名)
输入:tagid
输出:表名.tag_name字符串和错误信息
说明:
编辑:wang_jp
时间:2019年12月9日
*******************************************************************************/
func (this *OreProcessDTaglist) GetTagFullNameByID() (string, error) {
	o := orm.NewOrm() //新建orm对象
	o.Using("default")
	type taginfo struct {
		TageName  string
		TableName string
	}
	var tag taginfo
	err := o.Raw(`SELECT
					lst.tage_name,
					tb.table_name 
				FROM
					ore_process_d_taglist AS lst
					LEFT JOIN relevance_dcs_to_dbtable AS rel ON ( rel.dcs_id = lst.dcs_id )
					LEFT JOIN datatable_info AS tb ON ( tb.id = rel.datatable_id ) 
				WHERE
					lst.id = ?
					AND lst.status > 0`, this.Id).QueryRow(&tag)
	if err != nil {
		return "", err
	}
	fullname := fmt.Sprintf("%s.%s", tag.TableName, tag.TageName)
	this.TagFullName = fullname
	return this.TagFullName, nil
}

/*******************************************************************************
功能:通过tagFullname(含表名)获取tagId
输入:tagFullname
输出:tagId和错误信息
说明:
编辑:wang_jp
时间:2020年1月3日
*******************************************************************************/
func (this *OreProcessDTaglist) GetTagIDByFullName(tagFullName string) (int64, error) {
	o := orm.NewOrm() //新建orm对象
	o.Using("default")
	tagname := strings.Split(tagFullName, ".")
	if len(tagname) < 2 {
		//return 0, fmt.Errorf("TagName [%s] fomart error,must like XXX.YYY[变量名格式错误]", tagFullName)
		return 0, fmt.Errorf("实时变量全名 [%s] 格式错误,变量全名的格式为:tablename.tagname", tagFullName)
	}
	this.TagName = tagname[1]
	if err := o.Read(this, "TagName"); err != nil {
		return 0, err
	} else {
		return this.Id, nil
	}
}

/*******************************************************************************
功能:通过tag ID获取TagList信息
输入:tagId
输出:taglist
说明:
编辑:wang_jp
时间:2020年1月3日
*******************************************************************************/
func (this *OreProcessDTaglist) GetTagById(tagid ...int64) error {
	o := orm.NewOrm() //新建orm对象
	o.Using("default")
	if len(tagid) > 0 {
		this.Id = tagid[0]
	}
	return o.Read(this, "Id")
}

/*************************************************
功能:通过Tag_id获取该Tag的基本信息
输入:可选的TagId
输出:Taglists
说明:
编辑:wang_jp
时间:2020年3月17日
*************************************************/
func (this *OreProcessDTaglist) GetTagAttributByTagId(tagid ...int64) error {
	o := orm.NewOrm()
	o.Using("default") //选定数据库
	if len(tagid) > 0 {
		this.Id = tagid[0]
	}
	qt := o.QueryTable("OreProcessDTaglist")
	if this.Id > 0 {
		if err := qt.Filter("Id", this.Id).RelatedSel().One(this); err == nil {
			if _, err := this.GetTagFullNameByID(); err == nil {
			}
			return nil
		} else {
			return err
		}
	} else {
		return fmt.Errorf("Tag.Id must biger than 0[变量ID必须大于0]")
	}
}

/*************************************************
功能:通过Tag_Name获取该Tag的信息（基本信息和关联信息）
输入:可选的TagName
输出:Taglists
说明:
编辑:wang_jp
时间:2020年6月2日
*************************************************/
func (this *OreProcessDTaglist) GetTagAttributByTagName(tagname ...string) error {
	o := orm.NewOrm()
	o.Using("default") //选定数据库
	if len(tagname) > 0 {
		this.TagName = tagname[0]
	}
	qt := o.QueryTable("OreProcessDTaglist")
	if err := qt.Filter("TagName", this.TagName).RelatedSel().One(this); err == nil {
		if _, err := this.GetTagFullNameByID(); err == nil {
		}
		return nil
	} else {
		return err
	}
}

/*******************************************************************************
功能:通过TagName或者tagFullname(含表名)获取tag的基本属性信息
输入:[tagFullname] 可以是全名，也可以是不含表名的tag名
输出:tag属性信息
说明:
编辑:wang_jp
时间:2020年6月17日
*******************************************************************************/
func (this *OreProcessDTaglist) GetTagBaseAttributByName(tagFullName ...string) error {
	o := orm.NewOrm() //新建orm对象
	o.Using("default")
	if len(tagFullName) > 0 {
		tagname := strings.Split(tagFullName[0], ".")
		if len(tagname) > 1 { //如果分割成功
			this.TagName = tagname[1]
		} else {
			this.TagName = tagFullName[0] //分割不成功
		}
	}
	err := o.Read(this, "TagName")
	return err
}

/*************************************************
功能:通过Tag_id更新tag的庚顿信息
输入:可选的TagId
输出:错误信息
说明:
编辑:wang_jp
时间:2020年6月6日
*************************************************/
func (this *OreProcessDTaglist) UpdateTagGoldenAttrByTagId(tagid ...int64) error {
	o := orm.NewOrm()
	o.Using("default") //选定数据库
	if len(tagid) > 0 {
		this.Id = tagid[0]
	}

	_, err := o.Update(this,
		"ClassOf",
		"CompDev",
		"CompDevPercent",
		"CompTimeMax",
		"CompTimeMin",
		"Digits",
		"ExcDev",
		"ExcDevPercent",
		"ExcTimeMax",
		"ExcTimeMin",
		"GoldenId",
		"IsArchive",
		"IsCompress",
		"IsShutDown",
		"IsStep",
		"IsSummary",
		"MilliSecond",
		"Mirror",
		"Instrument",
		"IsScan",
		"Source",
		"Location1", "Location2", "Location3", "Location4", "Location5",
		"UserInt1", "UserInt2",
		"UserReal1", "UserReal2",
		"Equation", "Period", "TimeCopy", "Trigger")
	return err
}

/*******************************************************************************
功能:通过tagid获取CategoryName和ItemID
输入:tagid
输出:CategoryName和ItemID和错误信息
说明:如果有错误，CategoryName为空,ItemID=tagid
编辑:wang_jp
时间:2019年12月9日
*******************************************************************************/
func (this *OreProcessDTaglist) GetTagCategoryNameAndItemIDByID() (string, int64, error) {
	o := orm.NewOrm() //新建orm对象
	o.Using("default")
	type taginfo struct {
		TableProgramName string
		Itemid           int64
	}
	var tag taginfo
	err := o.Raw(`SELECT
			tb.table_program_name,tg.itemid
		FROM
			mine_table_list tb,
			( SELECT level_category AS category,resource_id AS itemid FROM ore_process_d_taglist WHERE id = ? ) tg 
		WHERE
			tb.id = tg.category`, this.Id).QueryRow(&tag)
	if err != nil { //错误的时候返回错误信息
		return "", this.Id, err
	}
	return tag.TableProgramName, tag.Itemid, nil
}

/*******************************************************************************
功能:通过tagid获取tag在taglist中指定列的信息
输入:tagid,field_name
输出:float64结果和error
说明:field_name所在列必须可转换为数字，否则输出0和错误信息
编辑:wang_jp
时间:2019年12月9日
*******************************************************************************/
func (this *OreProcessDTaglist) GetTagListFieldValueByID(field_name string) (string, error) {
	o := orm.NewOrm() //新建orm对象
	o.Using("default")
	var lists []orm.Params
	sql := fmt.Sprintf("SELECT * FROM ore_process_d_taglist WHERE id = %d", this.Id)
	if num, err := o.Raw(sql).Values(&lists); err != nil {
		return "", err
	} else if num == 0 {
		return "", errors.New(fmt.Sprintf("%s no data return", sql))
	} else { //取得了数据
		return fmt.Sprint(lists[0][field_name]), nil
	}
}

/*******************************************************************************
功能:通过tagid获取tag的上级Category中指定列的信息
输入:tagid,field_name
输出:float64结果和error
说明:field_name所在列必须可转换为数字，否则输出0和错误信息
编辑:wang_jp
时间:2019年12月9日
*******************************************************************************/
func (this *OreProcessDTaglist) GetTagTableRowFieldValueByID(field_name string) (string, error) {
	if tablename, item_id, err := this.GetTagCategoryNameAndItemIDByID(); err != nil {
		return "", err
	} else {
		return GetTableFieldValueByName(tablename, item_id, field_name)
	}
}

/*
功能:通过过程变量的变量名获取该名称在taglist中的id
输入:tagname:过程变量名
输出:ID结果和error
说明:
编辑:wang_jp
时间:2019年12月26日
*/
func (this *OreProcessDTaglist) GetTagIDByTagName() (int64, error) {
	o := orm.NewOrm() //新建orm对象
	o.Using("default")
	qt := o.QueryTable("OreProcessDTaglist").Filter("TagName", this.TagName)
	err := qt.One(this, "Id")
	if err != nil {
		return 0, err
	}
	return this.Id, nil
}

/*
功能:通过过程变量的变量名获取该名称在taglist中的描述
输入:tagname:过程变量名
输出:长描述、短描述和错误信息
说明:
编辑:wang_jp
时间:2020年5月20日
*/
func (this *OreProcessDTaglist) GetTagDescByTagName() (string, string, error) {
	o := orm.NewOrm() //新建orm对象
	o.Using("default")
	qt := o.QueryTable("OreProcessDTaglist").Filter("TagName", this.TagName)
	err := qt.One(this, "TagDescription", "TagPracticalDescription")
	if err != nil {
		return "", "", err
	}
	return this.TagDescription, this.TagPracticalDescription, nil
}

/*
功能:通过过程变量的变量名获取该变量的全名
输入:tagname:过程变量名
输出:变量全名和error
说明:
编辑:wang_jp
时间:2020年5月20日
*/
func (this *OreProcessDTaglist) GetTagFullNameByTagName() (string, error) {
	o := orm.NewOrm() //新建orm对象
	o.Using("default")
	qt := o.QueryTable("OreProcessDTaglist").Filter("TagName", this.TagName)
	err := qt.One(this, "Id")
	if err != nil {
		return "", err
	}
	return this.GetTagFullNameByID()
}

/*
功能:通过过程变量的变量名获取该变量的全名,并判断获得的全名与输入的参数是否相同(不区分大小写)
输入:fullname:过程变量名全名
输出:[bool] 输入参数与变量全名相符时为true，否则为false
	[string] 变量全名,正确的变量全名
	[error] 错误信息
说明:如果变量名在数据库中不存在,返回错误信息
编辑:wang_jp
时间:2020年5月20日
*/
func (this *OreProcessDTaglist) CheckFullName(fullname string) (bool, string, error) {
	_, err := this.GetTagFullNameByTagName() //获取变量全名
	if err != nil {                          //
		return false, "", err
	}
	if strings.ToLower(this.TagFullName) == strings.ToLower(fullname) {
		return true, fullname, nil
	} else {
		return false, this.TagFullName, nil
	}
}

/*
功能:通过过程变量的变量名获取该变量的全名
输入:tagname:过程变量名
输出:变量全名和error
说明:
编辑:wang_jp
时间:2020年5月20日
*/
func (this *OreProcessDTaglist) GetAllTagLists() ([]OreProcessDTaglist, error) {
	o := orm.NewOrm() //新建orm对象
	o.Using("default")
	var tags []OreProcessDTaglist
	qt := o.QueryTable("OreProcessDTaglist").Filter("TagName__isnull", false)
	_, err := qt.All(&tags)
	for i, tag := range tags {
		if _, err := tag.GetTagFullNameByID(); err == nil {
			tags[i].TagFullName = tag.TagFullName
		}
	}
	return tags, err
}

/****************************************************
功能:删除庚顿标签点
输入:无
输出:错误信息
说明:
编辑:wang_jp
时间:2020年6月22日
****************************************************/
func (t *OreProcessDTaglist) RemoveGoldenTagByGoldenId(tagid int64) error {
	var gdp goldengo.GoldenPoint
	micgd := new(MicGolden)
	err := micgd.GetHandel("OreProcessDTaglist", "RemoveGoldenTagByGoldenId")
	if err != nil {
		return err
	}
	defer micgd.ReleaseHandel()
	gdp.Base.Id = int(tagid)
	return gdp.RemovePointById(micgd.Handle)
}

/****************************************************
功能:更新单个平台标签点信息到庚顿数据库
输入:无
输出:错误信息
说明:前提是平台标签点中已经有对应的庚顿数据库ID
编辑:wang_jp
时间:2020年6月8日
****************************************************/
func (t *OreProcessDTaglist) UpdateTagToGolden() (goldengo.GoldenPoint, error) {
	var gdp goldengo.GoldenPoint
	micgd := new(MicGolden)
	err := micgd.GetHandel("OreProcessDTaglist", "UpdateTagToGolden")
	if err != nil {
		return gdp, err
	}
	defer micgd.ReleaseHandel()

	tagname := t.TagFullName
	if len(tagname) < 3 {
		return gdp, fmt.Errorf("变量TagId=[%d]尚未设置关联数据表", t.Id)
	}

	gdmap, err := micgd.GetTagPointInfoByName(tagname)
	if err != nil {
		return gdp, fmt.Errorf("标签点[%s]在庚顿数据库中不存在,请先建立标签点！", tagname)
	}
	gdp = gdmap[tagname]

	gdp = t.parsTaglistToGoldenTag(gdp)
	t.GoldenId = gdp.Base.Id
	err = gdp.UpdatePointById(micgd.Handle)
	if err != nil {
		return gdp, fmt.Errorf("更新标签点[%s]属性时发生错误:[%s]", tagname, err.Error())
	}
	return gdp, t.UpdateTagGoldenAttrByTagId()
}

/****************************************************
功能:插入或者更新单个平台标签点信息到庚顿数据库（插入或者更新）
输入:无
输出:[goldengo.GoldenPoint] 庚顿标签点信息
	[bool] 更新为false,插入为true
	[error] 错误信息
说明:需要先获取变量的全部信息
编辑:wang_jp
时间:2020年6月8日
****************************************************/
func (t *OreProcessDTaglist) InsertOrUpdateTagToGolden() (goldengo.GoldenPoint, bool, error) {
	var gdp goldengo.GoldenPoint
	isnew := false
	micgd := new(MicGolden)
	err := micgd.GetHandel("OreProcessDTaglist", "InsertOrUpdateTagToGolden")
	if err != nil {
		return gdp, isnew, err
	}
	defer micgd.ReleaseHandel()

	tagname := t.TagFullName
	if len(tagname) < 3 {
		return gdp, isnew, fmt.Errorf("变量TagId=[%d]尚未设置关联数据表", t.Id)
	}

	gdmap, err := micgd.GetTagPointInfoByName(tagname)
	if err != nil { //没有读取到标签,则新建标签点
		isnew = true
		gdtb := new(goldengo.GoldenTable)
		names := strings.Split(tagname, ".")
		if len(names) < 2 {
			return gdp, isnew, fmt.Errorf("平台变量标签全名[%s]不符合要求,标签全名应该是tablename.tagname", tagname)
		}
		gdtb.Name = names[0]                             //获取表名
		err := gdtb.GetTablePropertyByName(micgd.Handle) //获取表属性
		if err != nil {                                  //尚未建立表
			tb := new(DatatableInfo)
			tb.GetTableByName(names[0]) //获取描述信息
			gdtb.Desc = tb.TableDescription
			gdtb.AppendTable(micgd.Handle) //新建庚顿表
		}
		var base goldengo.GoldenBasePoint
		base.TableId = gdtb.Id
		gdp.Base = base
		gdp = t.parsTaglistToGoldenTag(gdp)

		err = gdp.InsertPoint(micgd.Handle)
		if err != nil {
			return gdp, isnew, fmt.Errorf("插入新变量标签[%s]失败:[%s]", gdp.Base.Tag, err.Error())
		}
	} else { //读取到了标签点,则更新标签点
		gdp = gdmap[tagname]

		gdp = t.parsTaglistToGoldenTag(gdp)
		err := gdp.UpdatePointById(micgd.Handle)
		if err != nil {
			return gdp, isnew, fmt.Errorf("更新变量标签失败:[%s]", err.Error())
		}
	}
	t.GoldenId = gdp.Base.Id
	return gdp, isnew, t.UpdateTagGoldenAttrByTagId()
}

/****************************************************
功能:将平台所有标签一次性更新到庚顿数据库
输入:无
输出:错误信息
说明:
编辑:wang_jp
时间:2020年6月8日
****************************************************/
func (t *OreProcessDTaglist) TaglistsToGolden() (int, int, int, float64, error) {
	st := time.Now()
	tags, err := t.GetAllTagLists() //获取所有Tag
	if err != nil {
		return 0, 0, 0, 0.0, fmt.Errorf("获取平台标签点时出错:[%s]", err.Error())
	}
	total := len(tags)
	var insertCnt, updateCnt int
	gocnt := 10
	tagcnts := make(chan [2]int, gocnt) //新插入数量
	defer close(tagcnts)

	var wait sync.WaitGroup
	cpg := (total + gocnt) / gocnt //每个go程处理的标签数

	var i int
	//fmt.Printf("共[%d]点Tag\n", total) //===============================================
	for i = 0; i < gocnt; i++ { //遍历，生成并发的Go程
		wait.Add(1)
		st := i * cpg
		ed := st + cpg
		if i == gocnt-1 {
			ed = total
		}
		tgs := tags[st:ed]
		//fmt.Printf("第[%d]个Go程,从[%d~%d]\n", i, st, ed) //===================================
		go func(taglists []OreProcessDTaglist, cnt chan [2]int) {
			defer wait.Done()
			var ncnt, ucnt int
			for _, tag := range taglists {
				_, isnew, err := tag.InsertOrUpdateTagToGolden()
				if err == nil {
					if isnew {
						ncnt++
					} else {
						ucnt++
					}
				}
			}
			var cnts [2]int
			cnts[0] = ncnt
			cnts[1] = ucnt
			cnt <- cnts
			//fmt.Printf("完成一个Go程,有[%d]点,新建[%d]个标签,更新[%d]个标签\n", len(taglists), ncnt, ucnt) //======================
		}(tgs, tagcnts)
	}
	for i = 0; i < gocnt; i++ {
		ct := <-tagcnts //接收计算结果
		insertCnt += ct[0]
		updateCnt += ct[1] //接收计算结果
	}
	wait.Wait()
	return total, insertCnt, updateCnt, time.Since(st).Seconds(), nil
}

/****************************************************
功能:把TagList中的参数传送到Golden Tag
输入:[gdp] 庚顿标签点指针
输出:无
说明:前提是平台标签点中已经有对应的庚顿数据库ID
编辑:wang_jp
时间:2020年6月8日
****************************************************/
func (t *OreProcessDTaglist) parsTaglistToGoldenTag(gdp goldengo.GoldenPoint) goldengo.GoldenPoint {
	exp := gdp.PlatEx
	exp.Id = t.Id
	exp.HHv = t.LimitHh
	exp.Hv = t.LimitH
	exp.LLv = t.LimitLl
	exp.Lv = t.LimitL

	gdp.PlatEx = exp

	bp := gdp.Base
	bp.Tag = t.TagName
	bp.TableDotTag = t.TagFullName
	switch strings.ToLower(t.TagType) {
	case "bool":
		bp.DataType = 0
	case "uint8":
		bp.DataType = 1
	case "int8":
		bp.DataType = 2
	case "char":
		bp.DataType = 3
	case "uint16":
		bp.DataType = 4
	case "int16":
		bp.DataType = 5
	case "uint32":
		bp.DataType = 6
	case "int32":
		bp.DataType = 7
	case "int64":
		bp.DataType = 8
	case "float16":
		bp.DataType = 9
	case "float32":
		bp.DataType = 10
	case "float64":
		bp.DataType = 11
	}

	bp.Unit = t.Unit
	bp.Desc = t.TagDescription
	bp.LowLimit = t.MinValue
	bp.HighLimit = t.MaxValue
	bp.Typical = t.Typical

	bp.ClassOf = uint(t.ClassOf)
	bp.CompDev = t.CompDev
	bp.CompDevPercent = t.CompDevPercent
	bp.CompTimeMax = t.CompTimeMax
	bp.CompTimeMin = t.CompTimeMin
	bp.Digits = t.Digits
	bp.ExcDev = t.ExcDev
	bp.ExcDevPercent = t.ExcDevPercent
	bp.ExcTimeMax = t.ExcTimeMax
	bp.ExcTimeMin = t.ExcTimeMin
	bp.IsArchive = t.IsArchive > 0
	bp.IsCompress = t.IsCompress > 0
	bp.IsShutDown = t.IsShutDown > 0
	bp.IsStep = t.IsStep > 0
	bp.IsSummary = t.IsSummary > 0
	bp.MilliSecond = t.MilliSecond
	bp.Mirror = t.Mirror

	gdp.Base = bp

	scan := gdp.Scan
	scan.Instrument = t.Instrument
	scan.IsScan = t.IsScan > 0
	scan.Source = t.Source
	scan.Locations[0] = t.Location1
	scan.Locations[1] = t.Location2
	scan.Locations[2] = t.Location3
	scan.Locations[3] = t.Location4
	scan.Locations[4] = t.Location5
	scan.UserInts[0] = t.UserInt1
	scan.UserInts[1] = t.UserInt2
	scan.UserReals[0] = t.UserReal1
	scan.UserReals[1] = t.UserReal2

	gdp.Scan = scan

	calc := gdp.Calc
	calc.Equation = t.Equation
	calc.Period = t.Period
	calc.TimeCopy = t.TimeCopy
	calc.Trigger = t.Trigger

	gdp.Calc = calc
	return gdp
}

/****************************************************
功能:同步庚顿数据库的标签点与平台数据库taglist标签点
输入:
	无
输出:[int] 进行同步的数据条数
	[int64] 耗时,毫秒数
	[error] 错误信息
编辑:wang_jp
时间:2020年6月2日
****************************************************/
func (t *OreProcessDTaglist) SynchGoldenPointAndPlatTaglist() (int, int, int, int64, error) {
	bg := time.Now()
	var total, g2p, p2g int
	micgd := new(MicGolden)
	err := micgd.GetHandel("OreProcessDTaglist", "SynchGoldenPointAndPlatTaglist") //建立到庚顿数据库的连接
	if err != nil {                                                                //判断连接是否有错误
		return total, g2p, p2g, 0, err
	}
	defer micgd.ReleaseHandel() //压后断开连接

	err = micgd.GetTables(true) //获取庚顿数据库详细信息
	if err != nil {             //判断连接是否有错误
		return total, g2p, p2g, 0, fmt.Errorf("获取庚顿数据库详细信息时错误:[%s]", err.Error())
	}

	for id, point := range micgd.Points { //遍历所有标签点
		tag := new(OreProcessDTaglist)
		tag.TagName = point.Base.Tag
		if tag.GetTagAttributByTagName() == nil { //读取到了tag变量
			total++
			//if point.PlatEx.Id == 0 {
			exp := point.PlatEx
			exp.Id = tag.Id
			exp.HHv = tag.LimitHh
			exp.Hv = tag.LimitH
			exp.LLv = tag.LimitLl
			exp.Lv = tag.LimitL

			point.PlatEx = exp

			bp := point.Base
			bp.Unit = tag.Unit
			bp.Desc = tag.TagDescription
			bp.LowLimit = tag.MinValue
			bp.HighLimit = tag.MaxValue
			bp.Typical = tag.Typical

			point.Base = bp

			tag.ClassOf = int(bp.ClassOf)
			tag.CompDev = bp.CompDev
			tag.CompDevPercent = bp.CompDevPercent
			tag.CompTimeMax = bp.CompTimeMax
			tag.CompTimeMin = bp.CompTimeMin
			tag.Digits = bp.Digits
			tag.ExcDev = bp.ExcDev
			tag.ExcDevPercent = bp.ExcDevPercent
			tag.ExcTimeMax = bp.ExcTimeMax
			tag.ExcTimeMin = bp.ExcTimeMin
			tag.GoldenId = bp.Id
			if bp.IsArchive {
				tag.IsArchive = 1
			}
			if bp.IsCompress {
				tag.IsCompress = 1
			}
			if bp.IsShutDown {
				tag.IsShutDown = 1
			}
			if bp.IsStep {
				tag.IsStep = 1
			}
			if bp.IsSummary {
				tag.IsSummary = 1
			}
			tag.MilliSecond = bp.MilliSecond
			tag.Mirror = bp.Mirror

			scan := point.Scan
			tag.Instrument = scan.Instrument
			if scan.IsScan {
				tag.IsScan = 1
			}
			tag.Source = scan.Source
			tag.Location1 = scan.Locations[0]
			tag.Location2 = scan.Locations[1]
			tag.Location3 = scan.Locations[2]
			tag.Location4 = scan.Locations[3]
			scan.Locations[4] = tag.Location5
			//tag.UserInt1 = scan.UserInts[0]
			//tag.UserInt2 = scan.UserInts[1]
			//tag.UserReal1 = scan.UserReals[0]
			//tag.UserReal2 = scan.UserReals[1]

			point.Scan = scan
			calc := point.Calc
			tag.Equation = calc.Equation
			tag.Period = calc.Period
			tag.TimeCopy = calc.TimeCopy
			tag.Trigger = calc.Trigger

			micgd.Points[id] = point
			//logs.Debug("更新平台标签点[%d]", tag.Id)
			err := tag.UpdateTagGoldenAttrByTagId()
			if err != nil {
				logs.Warn("同步庚顿标签点信息到平台Taglist时错误,庚顿id=%d,平台id=%d:[%s]", point.Base.Id, tag.Id, err.Error())
			} else {
				g2p++
			}
			logs.Debug("更新庚顿标签点[%d]", point.Base.Id)
			err = point.UpdatePointById(micgd.Handle)
			if err != nil {
				logs.Warn("同步平台Taglist信息到庚顿标签点时错误,庚顿id=%d,平台id=%d:[%s]", point.Base.Id, tag.Id, err.Error())
			} else {
				p2g++
			}
			//}
		}
	}
	return total, g2p, p2g, time.Since(bg).Milliseconds(), nil
}

/****************************************************
功能:从平台数据库taglist标签点同步所有庚顿数据库的标签点
输入:
	无
输出:[int] 进行同步的数据条数
	[int64] 耗时,毫秒数
	[error] 错误信息
编辑:wang_jp
时间:2020年6月2日
****************************************************/
func (t *OreProcessDTaglist) SynchGoldenPointFromePlatTaglist() (int, int, int64, error) {
	bg := time.Now()
	var total, p2g int
	micgd := new(MicGolden)
	err := micgd.GetHandel("OreProcessDTaglist", "SynchGoldenPointFromePlatTaglist") //建立到庚顿数据库的连接
	if err != nil {                                                                  //判断连接是否有错误
		return total, p2g, 0, err
	}
	defer micgd.ReleaseHandel() //压后断开连接

	err = micgd.GetTables(true) //获取庚顿数据库详细信息
	if err != nil {             //判断连接是否有错误
		return total, p2g, 0, err
	}

	for id, point := range micgd.Points { //遍历所有标签点
		tag := new(OreProcessDTaglist)
		tag.TagName = point.Base.Tag
		if tag.GetTagAttributByTagName() == nil { //读取到了tag变量
			total++
			exp := point.PlatEx
			exp.Id = tag.Id
			exp.HHv = tag.LimitHh
			exp.Hv = tag.LimitH
			exp.LLv = tag.LimitLl
			exp.Lv = tag.LimitL

			point.PlatEx = exp

			bp := point.Base
			bp.Tag = tag.TagName
			bp.Unit = tag.Unit
			bp.Desc = tag.TagDescription
			bp.LowLimit = tag.MinValue
			bp.HighLimit = tag.MaxValue
			bp.Typical = tag.Typical

			bp.ClassOf = uint(tag.ClassOf)
			bp.CompDev = tag.CompDev
			bp.CompDevPercent = tag.CompDevPercent
			bp.CompTimeMax = tag.CompTimeMax
			bp.CompTimeMin = tag.CompTimeMin
			bp.Digits = tag.Digits
			bp.ExcDev = tag.ExcDev
			bp.ExcDevPercent = tag.ExcDevPercent
			bp.ExcTimeMax = tag.ExcTimeMax
			bp.ExcTimeMin = tag.ExcTimeMin
			bp.IsArchive = tag.IsArchive > 0
			bp.IsCompress = tag.IsCompress > 0
			bp.IsShutDown = tag.IsShutDown > 0
			bp.IsStep = tag.IsStep > 0
			bp.IsSummary = tag.IsSummary > 0
			bp.Mirror = tag.Mirror

			point.Base = bp

			scan := point.Scan
			scan.Instrument = tag.Instrument
			scan.IsScan = tag.IsScan > 0
			scan.Source = tag.Source
			scan.Locations[0] = tag.Location1
			scan.Locations[1] = tag.Location2
			scan.Locations[2] = tag.Location3
			scan.Locations[3] = tag.Location4
			scan.Locations[4] = tag.Location5
			scan.UserInts[0] = tag.UserInt1
			scan.UserInts[1] = tag.UserInt2
			scan.UserReals[0] = tag.UserReal1
			scan.UserReals[1] = tag.UserReal2

			point.Scan = scan

			calc := point.Calc
			calc.Equation = tag.Equation
			calc.Period = tag.Period
			calc.TimeCopy = tag.TimeCopy
			calc.Trigger = tag.Trigger

			point.Calc = calc

			micgd.Points[id] = point

			err = point.UpdatePointById(micgd.Handle)
			if err != nil {
				logs.Warn("同步平台Taglist信息到庚顿标签点时错误,庚顿id=%d,平台id=%d:[%s]", point.Base.Id, tag.Id, err.Error())
			} else {
				p2g++
			}
		}
	}
	return total, p2g, time.Since(bg).Milliseconds(), nil
}
