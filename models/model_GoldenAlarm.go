package models

import (
	"fmt"
	"math"
	"time"

	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/bkzy-wangjp/goldengo"
)

type GoldenDataAlarm struct {
	GoldenDB        *MicGolden                //庚顿实时数据库
	PointIds        []int                     //标签点ID列表
	PointTypes      []int                     //标签点类型列表
	SnapShots       map[int]goldengo.SnapData //以庚顿id为Key的快照
	LastAlarmResult map[int]CalcTagAlarm      //以庚顿id为Key的最后报警信息
	TagDesc         map[int]string            //以庚顿id为Key的变量描述
	TagShortDesc    map[int]string            //以庚顿id为Key的变量短描述
	Locked          bool                      //更新时锁定
	LoadedTime      time.Time                 //信息载入时间
}

type GoldenPlatEx struct {
	goldengo.PointMicPlatEx
	MinValue     float64
	MaxValue     float64
	TagDesc      string
	TagShortDesc string
	AlarmOption  int
}

var GDDB *GoldenDataAlarm //庚顿数据库

/****************************************************
功能:获取庚顿数据库的基本信息
输入:
	无
输出:[int] 需要进行报警计算的标签点数
	[error] 错误信息
编辑:wang_jp
时间:2020年6月7日
****************************************************/
func (g *GoldenDataAlarm) GetGoldenPoints() (int, error) {
	err := g.GoldenDB.GetHandel("GoldenDataAlarm", "GetGoldenPoints") //获得句柄
	if err != nil {
		return 0, err
	}
	g.GoldenDB.GetTables(true) //读取数据表信息
	g.Locked = true            //加锁
	defer func() {             //压后释放
		g.Locked = false
		g.GoldenDB.ReleaseHandel()
	}()
	var pids, ptypes []int
	g.LastAlarmResult = make(map[int]CalcTagAlarm)
	g.TagDesc = make(map[int]string)
	g.TagShortDesc = make(map[int]string)
	var alarmCnt int
	for id, p := range g.GoldenDB.Points {
		pids = append(pids, id)
		ptypes = append(ptypes, p.Base.DataType)

		alrm := new(CalcTagAlarm)
		tag := new(OreProcessDTaglist)
		tag.Id = p.PlatEx.Id
		alrm.Tag = tag
		if p.Base.DataType > 0 && p.Scan.Locations[4] > 0 { //需要报警判断的才加载相关信息
			alarmCnt++
			err := alrm.GetLastAlartMsg()
			if err != nil {
				tag.GetTagAttributByTagId()
				g.TagDesc[id] = tag.TagDescription
				g.TagShortDesc[id] = tag.TagPracticalDescription
			} else {
				g.TagDesc[id] = alrm.Tag.TagDescription
				g.TagShortDesc[id] = alrm.Tag.TagPracticalDescription
			}
		}
		g.LastAlarmResult[id] = *alrm
	}
	g.PointIds = pids
	g.PointTypes = ptypes
	g.LoadedTime = time.Now()
	return alarmCnt, nil
}

/****************************************************
功能:获取庚顿数据库的快照信息并进行对比
输入:
	无
输出:错误信息
编辑:wang_jp
时间:2020年6月7日
****************************************************/
func (g *GoldenDataAlarm) GetGoldenSnapShotsAndComp() error {
	err := g.GoldenDB.GetHandel("GoldenDataAlarm", "GetGoldenSnapShotsAndComp") //获得句柄
	if err != nil {
		return err
	}
	defer g.GoldenDB.ReleaseHandel() //释放句柄
	snaps, err := g.GoldenDB.GetSnapShotById(g.PointIds, g.PointTypes)
	if err != nil {
		return err
	}
	for id, snap := range g.SnapShots { //遍历老快照
		if g.GoldenDB.Points[id].Base.DataType > 0 && g.GoldenDB.Points[id].Scan.Locations[4] > 0 { //非BOOL类型
			if IsExistItem(id, snaps) { //新快照中存在老快照的ID
				if math.Abs(snap.Rtsd.Value-snaps[id].Rtsd.Value) > g.GoldenDB.Points[id].Base.ExcDev { //快照值之差大于例外偏差
					platex := new(GoldenPlatEx)
					platex.HHv = g.GoldenDB.Points[id].PlatEx.HHv
					platex.Hv = g.GoldenDB.Points[id].PlatEx.Hv
					platex.LLv = g.GoldenDB.Points[id].PlatEx.LLv
					platex.Lv = g.GoldenDB.Points[id].PlatEx.Lv
					platex.MaxValue = g.GoldenDB.Points[id].Base.HighLimit
					platex.MinValue = g.GoldenDB.Points[id].Base.LowLimit
					platex.Id = g.GoldenDB.Points[id].PlatEx.Id
					platex.TagDesc = g.TagDesc[id]
					platex.TagShortDesc = g.TagShortDesc[id]
					platex.AlarmOption = g.GoldenDB.Points[id].Scan.Locations[4]
					g.LastAlarmResult[id] = platex.AnalogDataAlarmEngine(snaps[id], g.LastAlarmResult[id]) //进行报警判断
				}
			}
		}
	}
	g.SnapShots = make(map[int]goldengo.SnapData)
	g.SnapShots = snaps
	return nil
}

/****************************************************
功能:模拟量报警计算
输入:
	[snap] 快照数据
	[lastalarm] 上次报警信息
输出:[CalcTagAlarm] 计算之后的报警信息
编辑:wang_jp
时间:2020年6月8日
****************************************************/
func (e *GoldenPlatEx) AnalogDataAlarmEngine(snap goldengo.SnapData, lastalrm CalcTagAlarm) CalcTagAlarm {
	//没有设置报警值
	if (e.LLv == 0 && e.Lv == 0 && e.HHv == 0 && e.Hv == 0 && e.MaxValue == 0 && e.MinValue == 0) || snap.Rtsd.Quality != 0 {
		return lastalrm
	}
	//logs.Debug("进入报警比较:[%+v]", e)
	val := snap.Rtsd.Value        //取快照值
	limitv := lastalrm.LimitValue //取报警限
	var alv int = 0

	if e.alarmOption(1) && val < e.MinValue {
		alv = -3
		limitv = e.MinValue
	}
	if e.alarmOption(2) && val < e.LLv {
		alv = -2
		limitv = e.LLv
	}
	if e.alarmOption(4) && val < e.Lv {
		alv = -1
		limitv = e.Lv
	}
	if e.alarmOption(8) && val > e.Hv {
		alv = 1
		limitv = e.Hv
	}
	if e.alarmOption(16) && val > e.HHv {
		alv = 2
		limitv = e.HHv
	}
	if e.alarmOption(32) && val > e.MaxValue {
		alv = 3
		limitv = e.MaxValue
	}

	if alv != lastalrm.AlarmStatus || lastalrm.Id == 0 { //报警状态发生改变或者没有初始状态
		if lastalrm.Id == 0 {
			lastalrm.Id = 1
		}
		var alrm CalcTagAlarm
		lastalrm.AlarmStatus = alv
		lastalrm.Datatime = TimeFormat(time.Unix(snap.Rtsd.Time/1e3, snap.Rtsd.Time%1e3))
		lastalrm.LimitValue = limitv
		tag := new(OreProcessDTaglist)
		tag.Id = e.Id
		lastalrm.Tag = tag
		lastalrm.TagDesc = e.TagDesc
		lastalrm.TagShotDesc = e.TagShortDesc
		lastalrm.TagValue = val
		alrm = lastalrm
		alrm.Id = 0
		go alrm.InsertAlarm()
	}
	return lastalrm
}

/****************************************************
功能:根据报警配置和掩码，求是否设置了报警开关
输入:[mask] 掩码;Min=1,LL=2,L=4,H=8,HH=16,Max=32
	[option] 报警设置,在Location5中设置
输出:如果设置了报警，返回true，否则返回false
编辑:wang_jp
时间:2020年6月8日
****************************************************/
func (e *GoldenPlatEx) alarmOption(mask int) bool {
	return mask&e.AlarmOption > 0
}

/****************************************************
功能:模拟量报警信息存储
输入:无
输出:无
编辑:wang_jp
时间:2020年6月8日
****************************************************/
func (a *CalcTagAlarm) InsertAlarm() {
	o := orm.NewOrm()  //新建orm对象
	o.Using("default") //选定数据库
	_, err := o.Insert(a)
	if err != nil {
		logs.Warn("插入报警信息错误,待插入信息:[%+v];错误信息:[%s]", a, err.Error())
	}
}

/****************************************************
功能:获取最后的报警信息
输入:无
输出:无
编辑:wang_jp
时间:2020年6月8日
****************************************************/
func (a *CalcTagAlarm) GetLastAlartMsg() error {
	o := orm.NewOrm()  //新建orm对象
	o.Using("default") //选定数据库
	tagid := a.Tag.Id
	err := o.QueryTable("CalcTagAlarm").Filter("Tag__Id", tagid).RelatedSel("Tag").OrderBy("-Datatime").One(a)
	return err
}

/****************************************************
功能:获取庚顿报警标签点的当前信息
输入:[pointid] 庚顿标签点的ID
输出:标签点信息,错误信息
编辑:wang_jp
时间:2020年6月9日
****************************************************/
func (g *GoldenDataAlarm) SelectPoint(pointid int) (interface{}, error) {
	if g.Locked == false {
		if IsExistItem(pointid, g.GoldenDB.Points) {
			type gd struct {
				Gdp   goldengo.GoldenPoint
				Snap  goldengo.SnapData
				Alarm CalcTagAlarm
			}
			point := new(gd)
			point.Alarm = g.LastAlarmResult[pointid]
			point.Gdp = g.GoldenDB.Points[pointid]
			point.Snap = g.SnapShots[pointid]
			return point, nil
		} else {
			return nil, fmt.Errorf("标签点[%d]不存在", pointid)
		}
	} else {
		return nil, fmt.Errorf("标签点Map正在更新,请稍后访问")
	}
}
