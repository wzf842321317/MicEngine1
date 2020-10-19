package models

import (
	"testing"
	//"time"
)

//func TestSaveBatchKpiResultToDB(t *testing.T) {
//	res := CalcKpiResult{Id: 0,
//		TagType:         "process",
//		TagId:           5671,
//		TagName:         "x3_asl_asl-xc1_MF1_XLFLFX3_F1-091_FD1-086_run:1",
//		KpiConfigListId: 13268,
//		KpiKey:          "on_time",
//		KpiTag:          "x3_asl_asl-xc1_MF1_XLFLFX3_F1-091_FD1-086_run:1__-2_0_on_time11",
//		KpiName:         "磨浮一车间_锌硫分离精选Ⅳ1号浮选机电机_运行__累计时间",
//		KpiPeriod:       -1,
//		KpiValue:        753.5,
//		CalcEndingTime:  "2020-05-18 12:00:00",
//		InbaseTime:      "2020-05-18 18:58:14",
//	}
//	var ress []CalcKpiResult
//	ress = append(ress, res)
//	num := SaveBatchKpiResultToDB(ress)
//	t.Log(num)
//}

/*
func TestGetRoportLists(t *testing.T) {
	tests := []struct {
		readtype int
	}{
		{0},
		{1},
		{2},
	}
	for i, tt := range tests {
		num, reports, err := GetRoportLists(tt.readtype)
		if err != nil {
			t.Errorf("第%d行错误:[%s]", i, err.Error())
		} else {
			t.Logf("第%d行提取%d条结果", i, num)
			t.Logf("%+v", reports)
		}
	}
}

func TestGetReportNodesByUserId(t *testing.T) {
	tests := []struct {
		id int64
	}{
		{115},
		{117},
	}
	for i, tt := range tests {
		w := new(SysUser)
		rpts, err := w.GetReportNodesByUserId(tt.id)
		if err != nil {
			t.Errorf("第%d行错误:[%s]", i, err.Error())
		} else {
			t.Logf("第%d行Id=%d,共%d条数据", i, tt.id, len(rpts))
			for _, rpt := range rpts {
				t.Log(rpt)
			}
		}
	}
}

func TestGetChildNodes(t *testing.T) {
	tests := []struct {
		id    int64
		level string
	}{
		{115, "1"},
		{115, "1-2"},
		{115, "1-2-3"},
	}
	for i, tt := range tests {
		w := new(CalcKpiReportList)
		w.LevelCode = tt.level
		num, rpts, err := w.GetChildNodes(tt.id)
		if err != nil {
			t.Errorf("第%d行错误:[%s]", i, err.Error())
		} else {
			t.Logf("第%d行Id=%d,共%d条数据", i, tt.id, num)
			for _, rpt := range rpts {
				t.Log(rpt)
			}
		}
	}
}

func TestGetAttributeById(t *testing.T) {
	tests := []struct {
		id    int64
		level string
	}{
		{1, "1"},
		{2, "1-2"},
		{3, "1-2-3"},
		{23, "1-2-3-4"},
	}
	for i, tt := range tests {
		w := new(CalcKpiReportList)
		w.Id = tt.id
		err := w.GetAttributeById()
		if err != nil {
			t.Errorf("第%d行错误:[%s]", i, err.Error())
		} else {
			t.Logf("第%d行Id=%d,数据内容:%+v,车间:%+v", i, tt.id, w, w.Workshop)
		}
	}
}

func TestLoadPermissions(t *testing.T) {
	tests := []struct {
		id  int64
		uid int64
	}{
		{1, 115},
		{2, 115},
		{3, 115},
		{4, 115},
	}
	for i, tt := range tests {
		w := new(CalcKpiReportList)
		w.Id = tt.id
		num, pms, err := w.LoadPermissions(tt.uid)
		if err != nil {
			t.Errorf("第%d行错误:[%s]", i, err.Error())
		} else {
			t.Logf("第%d行Id=%d,权限数:%d", i, tt.id, num)
			for _, p := range pms {
				t.Log(p)
			}
		}
	}
}
func TestUpdate(t *testing.T) {
	tests := []struct {
		id  int64
		uid int64
	}{
		{6, 115},
		{7, 115},
		{8, 115},
	}
	for i, tt := range tests {
		w := new(CalcKpiReportList)
		w.Id = tt.id
		w.GetAttributeById()
		err := w.Update("TemplateUrl", "ResultUrl", "StartTime")
		if err != nil {
			t.Errorf("第%d行错误:[%s]", i, err.Error())
		}
	}
}


func TestGetTemplatesList(t *testing.T) {
	tests := []struct {
		id    int64
		level string
	}{
		{1, "1"},
		{2, "1-2"},
		{3, "1-2-3"},
		{4, "1-2-3-4"},
		{5, "1-2-3-5"},
	}
	for i, tt := range tests {
		w := new(CalcKpiReportList)
		w.Id = tt.id
		tpls, err := w.GetTemplatesListById()
		if err != nil {
			t.Errorf("第%d行错误:[%s]", i, err.Error())
		} else {
			t.Logf("第%d行Id=%d,数据内容:%+v", i, tt.id, tpls)
		}
	}
}

func TestGetWorkshopLists(t *testing.T) {
	res, err := GetWorkshopLists()
	if err != nil {
		t.Error(err.Error())
	} else {
		for _, w := range res {
			t.Log(w)
		}
	}
}


func TestGetWorkshopInfoById(t *testing.T) {
	tests := []struct {
		id int
	}{
		{0},
		{1},
		{2},
	}
	for i, tt := range tests {
		w := new(OreProcessDWorkshop)
		err := w.GetInfoById(tt.id)
		if err != nil {
			t.Errorf("第%d行错误:[%s]", i, err.Error())
		} else {
			t.Logf("第%d行Id=%d:%+v", i, tt.id, w)
		}
	}
}


func TestGetKpiConfigListIdByKpiTag(t *testing.T) {
	tests := []struct {
		sid int64
		tag string
		id  int64
	}{
		{0, "x2_asl_asl-xc1_SK1_ZXS1_S1-017_SY1-009_HL-warning:1__hour_RisingCnt", 10},
		{0, "x3_asl_asl-xc1_SK1_GK1_S1-002_SD1-001_fault:1__hour_on_time", 36},
		{64, "x3_asl_asl-xc1_SK1_GK1_S1-006_SD1-004_fault:1__hour_RisingCnt", 64},
		{1101, "x3_asl_asl-xc1_SK2_ZXS2_S2-013_SD2-008_start:1__hour_RisingCnt", 1101},
		{0, "x3_asl_asl-xc1_SK2_ZXS2_S2-013_SD2-008_start:1__hour_RisingCntxxx", 0},
	}
	var cfg CalcKpiConfigList
	for i, tt := range tests {
		cfg.Id = tt.sid
		cfg.KpiTag = tt.tag
		id, err := cfg.GetKpiConfigListInfo()
		if err == nil {
			if cfg.Id != tt.id {
				t.Errorf("第%d行,期望值是%d,实际值是%d,返回值是%d", i, tt.id, cfg.Id, id)
			} else {
				t.Logf("%+v", cfg)
			}
		} else {
			t.Log(err)
		}
	}
}

func TestGetKpiSingleData(t *testing.T) {
	tests := []struct {
		id      int64
		tag     string
		bgtime  string
		endtime string
	}{
		{0, "x1_asl_asl-xc1_MF1_MKⅠ3_MY1-002_sum:1__team_Diff", "2020-05-04 09:30:00", "2020-05-04 17:30:00"},
		{0, "x1_asl_asl-xc1_MF1_MKⅠ3_MY1-002_sum:1__team_Diff", "", "2020-05-04 17:30:00"},
		{0, "x1_asl_asl-xc1_MF1_MKⅠ3_MY1-002_sum:1__team_Diff", "2020-05-04 09:30:00", ""},
		{10734, "", "2020-05-04 09:30:00", "2020-05-04 17:30:00"},
		{10734, "", "", "2020-05-04 17:30:00"},
		{10734, "", "2020-05-04 09:30:00", ""},
		{10734, "", "2020-05-03 09:30:00", "2020-05-03 17:30:00"},
		{10734, "", "", "2020-05-03 17:30:00"},
		{10734, "", "2020-05-03 09:30:00", ""},
		{10734, "", "", "2020-05-03 09:30:00"},
	}
	var cfg CalcKpiConfigList
	for i, tt := range tests {
		cfg.Id = tt.id
		cfg.KpiTag = tt.tag
		value, err := cfg.GetKpiSingleData(tt.bgtime, tt.endtime)
		if err == nil {
			t.Logf("第%d行:%f", i, value)
		} else {
			t.Log(err)
		}
	}
}

func TestGetKpiDatas(t *testing.T) {
	tests := []struct {
		id       int64
		tag      string
		bgtime   string
		endtime  string
		readType int
	}{
		{0, "x1_asl_asl-xc1_MF1_MKⅠ3_MY1-002_sum:1__team_Diff", "2020-05-01 17:30:00", "2020-05-03 17:30:00", -1},
		{0, "x1_asl_asl-xc1_MF1_MKⅠ3_MY1-002_sum:1__team_Diff", "2020-05-01 17:30:00", "2020-05-03 17:30:00", 0},
		{0, "x1_asl_asl-xc1_MF1_MKⅠ3_MY1-002_sum:1__team_Diff", "2020-05-01 17:30:00", "2020-05-03 17:30:00", 1},
		{0, "x1_asl_asl-xc1_MF1_MKⅠ3_MY1-002_sum:1__team_DiffXXX", "2020-05-01 17:30:00", "2020-05-03 17:30:00", 1},
		{10734, "", "2020-05-01 17:30:00", "2020-05-03 17:30:00", -1},
		{10734, "", "2020-05-01 17:30:00", "2020-05-03 17:30:00", 0},
		{10734, "", "2020-05-01 17:30:00", "2020-05-03 17:30:00", 1},
		{10734, "", "2020-05-05 17:30:00", "2020-05-06 17:30:00", 0},
	}
	var cfg CalcKpiConfigList
	for i, tt := range tests {
		cfg.Id = tt.id
		cfg.KpiTag = tt.tag
		datas, err := cfg.GetKpiDatas(tt.bgtime, tt.endtime, tt.readType)
		t.Logf("第%d行:-----------", i)
		if err == nil {
			for _, data := range datas {
				t.Logf("%+v", data)
			}
		} else {
			t.Log(err)
		}
	}
}

func TestGetKpiStatistic(t *testing.T) {
	tests := []struct {
		id      int64
		tag     string
		bgtime  string
		endtime string
		key     string
	}{
		{0, "x1_asl_asl-xc1_MF1_MKⅠ3_MY1-002_sum:1__team_Diff", "2020-05-01 17:30:00", "2020-05-03 17:30:00", "sd"},
		{0, "x1_asl_asl-xc1_MF1_MKⅠ3_MY1-002_sum:1__team_Diff", "2020-05-01 17:30:00", "2020-05-03 17:30:00", "mode"},
		{0, "x1_asl_asl-xc1_MF1_MKⅠ3_MY1-002_sum:1__team_Diff", "2020-05-01 17:30:00", "2020-05-03 17:30:00", "min"},
		{0, "x1_asl_asl-xc1_MF1_MKⅠ3_MY1-002_sum:1__team_DiffXXX", "2020-05-01 17:30:00", "2020-05-03 17:30:00", "max"},
		{10734, "", "2020-05-01 17:30:00", "2020-05-03 17:30:00", "diff"},
		{10734, "", "2020-05-01 17:30:00", "2020-05-03 17:30:00", "plusdiff"},
		{10734, "", "2020-05-01 17:30:00", "2020-05-03 17:30:00", "gtzcnt"},
		{10734, "", TimeFormat(time.Now()), TimeFormat(time.Now().Add(8 * time.Hour)), "ltzcnt"}, //错误的时间
	}
	var cfg CalcKpiConfigList
	for i, tt := range tests {
		cfg.Id = tt.id
		cfg.KpiTag = tt.tag
		datas, err := cfg.GetKpiStatistic(tt.key, tt.bgtime, tt.endtime)
		t.Logf("第%d行:-----------", i)
		if err == nil {
			t.Logf("%+v", datas)
		} else {
			t.Log(err)
		}
	}
}

func TestGetKpiDataStatisticByKey(t *testing.T) {
	tests := []struct {
		id      int64
		tag     string
		bgtime  string
		endtime string
		key     string
	}{
		{0, "x1_asl_asl-xc1_MF1_MKⅠ3_MY1-002_sum:1__team_Diff", "2020-05-01 17:30:00", "2020-05-03 17:30:00", "sd"},
		{0, "x1_asl_asl-xc1_MF1_MKⅠ3_MY1-002_sum:1__team_Diff", "2020-05-01 17:30:00", "2020-05-03 17:30:00", "mode"},
		{0, "x1_asl_asl-xc1_MF1_MKⅠ3_MY1-002_sum:1__team_Diff", "2020-05-01 17:30:00", "2020-05-03 17:30:00", "min"},
		{0, "x1_asl_asl-xc1_MF1_MKⅠ3_MY1-002_sum:1__team_DiffXXX", "2020-05-01 17:30:00", "2020-05-03 17:30:00", "max"}, //错误的tag
		{10734, "", "2020-05-01 17:30:00", "2020-05-03 17:30:00", "diff"},
		{10734, "", "2020-05-01 17:30:00", "2020-05-03 17:30:00", "plusdiff"},
		{10734, "", "2020-05-01 17:30:00", "2020-05-03 17:30:00", "gtzcnt"},
		{10734, "", TimeFormat(time.Now()), TimeFormat(time.Now().Add(8 * time.Hour)), "ltzcnt"}, //错误的时间
		{100734, "", "2020-05-05 17:30:00", "2020-05-06 17:30:00", "ltzcnt"},                     //错误的ID
	}
	var cfg CalcKpiConfigList
	for i, tt := range tests {
		cfg.Id = tt.id
		cfg.KpiTag = tt.tag
		datas, err := cfg.GetKpiDataStatisticByKey(tt.key, tt.bgtime, tt.endtime)
		t.Logf("第%d行:-----------", i)
		if err == nil {
			t.Logf("%+v", datas)
		} else {
			t.Log(err)
		}
	}
}



func TestGetSysLog(t *testing.T) {
	tests := []struct {
		bgtime  string
		endtime string
		userid  int64
		systype int64
		oprtype int64
		limit   int64
		offset  int64
		desc    string
	}{
		{"2020-04-06 00:00:00", "2020-04-07 00:00:00", -1, -1, -1, 0, 0, ""},
		{"2020-04-06 00:00:00", "2020-04-07 00:00:00", -1, -1, -1, 100, 0, ""},
		{"2020-04-06 00:00:00", "2020-04-07 00:00:00", -1, -1, -1, 100, 1, ""},
	}
	for _, tt := range tests {
		rows, _, logs, err := GetSysLog(tt.bgtime, tt.endtime, tt.userid, tt.systype, tt.oprtype, tt.limit, tt.offset, tt.desc)
		if err != nil {
			t.Error(err.Error())
		} else {
			t.Logf("--------%d--------", rows)
			for k, log := range logs {
				t.Logf("%d:%+v%+v", k, log, log.User)
			}
		}
	}
}
func TestGetSysLogAnalyse(t *testing.T) {
	tests := []struct {
		bgtime  string
		endtime string
		userid  int64
		systype int64
		oprtype int64
		desc    string
	}{
		{"2020-04-06 00:00:00", "2020-04-07 00:00:00", -1, -1, -1, ""},
		{"2020-04-06 00:00:00", "2020-04-07 00:00:00", -1, -1, -1, ""},
		{"2020-04-06 00:00:00", "2020-04-07 00:00:00", -1, -1, -1, ""},
	}
	for _, tt := range tests {
		logs, err := GetSysLogAnalyse(tt.bgtime, tt.endtime, tt.userid, tt.systype, tt.oprtype, tt.desc)
		if err != nil {
			t.Error(err.Error())
		} else {
			t.Log("----------------")
			t.Logf("%+v", logs)
		}
	}
}


func TestGetCheckNodesByLevel(t *testing.T) {
	nodes, err := GetCheckNodesByLevel()
	if err == nil {
		for k, node := range nodes {
			t.Logf("%d:%+v", k, *node)
		}
	}
}

func TestGetCheckTreeNodesByUserId(t *testing.T) {
	tests := []struct {
		uid int64
	}{
		{115},
	}
	for _, tt := range tests {
		nodes, err := GetCheckTreeNodesByUserId(tt.uid)
		if err == nil {
			t.Logf("----%d-----", tt.uid)
			for _, node := range nodes {
				t.Logf("%+v", *node)
			}
		}
	}
}
func TestGetCheckResultsByLevelCode(t *testing.T) {
	tests := []struct {
		level   string
		bgtime  string
		endtime string
		stonly  bool
	}{
		{"1-2-7-35", "2020-04-06 00:00:00", "2020-04-07 00:00:00", true},
	}
	for _, tt := range tests {
		nodes, err := GetCheckResultsByLevelCode(tt.level, tt.bgtime, tt.endtime, tt.stonly)
		if err == nil {
			t.Logf("----%s-----", tt.level)
			for _, node := range nodes {
				t.Logf("%+v", *node)
			}
		}
	}
}

func TestGetCheckResultsByTagId(t *testing.T) {
	tests := []struct {
		id      int64
		bgtime  string
		endtime string
		stonly  bool
	}{
		{351, "2020-04-06 00:00:00", "2020-04-07 00:00:00", true},
	}
	for _, tt := range tests {
		nodes, err := GetCheckResultsByTagId(tt.id, tt.bgtime, tt.endtime, tt.stonly)
		if err == nil {
			t.Logf("----%d-----", tt.id)
			for _, node := range nodes {
				t.Logf("%+v,%+v", *node, *node.CheckItem)
			}
		}
	}
}

func TestGetMineDeptLists(t *testing.T) {
	nodes, err := GetMineDeptLists(1, 2, 3, 4, 5)
	if err == nil {
		for _, node := range nodes {
			t.Logf("%+v", *node)
		}
	}
}

func TestGetMineCheckImg(t *testing.T) {
	tests := []struct {
		tableid int64
		itemid  int64
	}{
		{48, 101157},
		{76, 0},
		{143, 135},
	}
	for _, tt := range tests {
		nodes, err := GetMineCheckImg(tt.tableid, tt.itemid)
		if err == nil {
			t.Logf("%+v", nodes)
		}
	}
}

func TestGetCheckSiteByLevelCode(t *testing.T) {
	tests := []struct {
		level   string
		bgtime  string
		endtime string
	}{
		{"1-2-10-13679", "2020-04-06 00:00:00", "2020-04-07 00:00:00"},
		{"1-2-10", "2020-04-06 00:00:00", "2020-04-07 00:00:00"},
	}
	for _, tt := range tests {
		nodes, err := GetCheckSiteByLevelCode(tt.level, tt.bgtime, tt.endtime)
		if err == nil {
			//for _, node := range nodes {
			t.Logf("%+v", nodes)
			//}
		}
	}
}
func TestGetCheckTagListsLevelMap(t *testing.T) {
	nodes, err := GetCheckTagListsLevelMap()
	if err == nil {
		for _, node := range nodes {
			t.Logf("%d", node)
		}
	}
}


func TestGetGoodsNodesByLevel(t *testing.T) {
	nodes, err := GetGoodsNodesByLevel()
	if err == nil {
		for _, node := range nodes {
			t.Logf("%+v", *node)
		}
	}
}

func TestGetGoodsTreeNodesByUserId(t *testing.T) {
	tests := []struct {
		uid int64
	}{
		{115},
	}
	for _, tt := range tests {
		nodes, err := GetGoodsTreeNodesByUserId(tt.uid)
		if err == nil {
			t.Logf("----%d-----", tt.uid)
			for _, node := range nodes {
				t.Logf("%+v", *node)
			}
		}
	}
}
func TestGetGoodsResultsByLevelCode(t *testing.T) {
	tests := []struct {
		level   string
		bgtime  string
		endtime string
		stonly  bool
	}{
		{"1-2-7-35", "2020-04-06 00:00:00", "2020-04-07 00:00:00", true},
	}
	for _, tt := range tests {
		nodes, err := GetGoodsResultsByLevelCode(tt.level, tt.bgtime, tt.endtime, tt.stonly)
		if err == nil {
			t.Logf("----%s-----", tt.level)
			for _, node := range nodes {
				t.Logf("%+v,%+v", *node, *node.GoodsConfigInfo)
			}
		}
	}
}

func TestGetGoodsInfoByGoodsId(t *testing.T) {
	tests := []struct {
		id int64
	}{
		{351},
	}
	for _, tt := range tests {
		nodes, err := GetGoodsInfoByGoodsId(tt.id)
		if err == nil {
			t.Logf("----%d-----", tt.id)
			for _, node := range nodes {
				t.Logf("%+v,%+v", *node, *node.Goods)
			}
		}
	}
}

func TestGetSampleSubBySampleId(t *testing.T) {
	tests := []struct {
		uid int64
	}{
		{20},
	}
	for _, tt := range tests {
		nodes, err := GetSampleSubBySampleId(tt.uid)
		if err == nil {
			t.Logf("----%d-----", tt.uid)
			for _, node := range nodes {
				t.Logf("%+v", *node)
			}
		}
	}
}

func TestGetSampleLabResultBySampleId(t *testing.T) {
	tests := []struct {
		uid     int64
		bgtime  string
		endtime string
	}{
		{11, "2020-01-10 00:00:00", "2020-01-12 00:00:00"},
	}
	for _, tt := range tests {
		nodes, err := GetSampleLabResultBySampleId(tt.uid, tt.bgtime, tt.endtime)
		if err == nil {
			t.Logf("----%d-----", tt.uid)
			for _, node := range nodes {
				t.Logf("%+v", *node)
			}
		}
	}
}

func TestGetSampleLabTreeNodesByUserId(t *testing.T) {
	tests := []struct {
		uid int64
	}{
		{115},
	}
	for _, tt := range tests {
		nodes, err := GetSampleLabTreeNodesByUserId(tt.uid)
		if err == nil {
			t.Logf("----%d-----", tt.uid)
			for _, node := range nodes {
				t.Logf("%+v", *node)
			}
		}
	}
}


func TestGetSampleLabTreeNodesMine(t *testing.T) {
	tests := []struct {
		level int
	}{
		{1},
		{2},
		{3},
		{4},
	}
	for _, tt := range tests {
		nodes, err := GetSampleLabNodesByLevel(tt.level)
		if err == nil {
			t.Logf("----%d-----", tt.level)
			for _, node := range nodes {
				t.Logf("%+v", *node)
			}
		}
	}
}


func TestGetUserList(t *testing.T) {
	users, err := GetUserList("张景宁", "刘")
	if err == nil {
		for _, user := range users {
			t.Logf("%+v", &user)
		}
	}
}


func TestCheckWebMenu(t *testing.T) {
	CheckWebMenu()
}


func TestMd5str(t *testing.T) {
	tests := []struct {
		str string
		res string
	}{
		{"MD5testing", "f7bb96d1dcd6cfe0e5ce1f03e35f84bf"},
		{"123456", "e10adc3949ba59abbe56e057f20f883e"},
	}
	for _, tt := range tests {
		res := Md5str(tt.str)
		if res != tt.res {
			t.Errorf("错误,得到的结果:%s,期望的结果:%s", res, tt.res)
		}
	}
}

func TestUserLogIn(t *testing.T) {
	tests := []struct {
		name string
		psw  string
		uid  int64
	}{
		{"bkzy1", "123456", 117},
		{"zjs_xxk", "123456", 2},
		{"黄中省", "123456", 118},
	}
	for _, tt := range tests {
		res, err := UserLogIn(tt.name, tt.psw)
		if err != nil {
			t.Error(err)
		} else {
			if tt.uid != res.Id {
				t.Errorf("错误,查询到的id是:%d,期望的id是%d", res.Id, tt.uid)
			} else {
				t.Log(res, res.Roles)
			}
		}
	}
}

func TestGetRolesByUserId(t *testing.T) {
	tests := []struct {
		id int64
	}{
		{117},
		{2},
		{118},
	}
	for _, tt := range tests {
		res, err := GetRolesByUserId(tt.id)
		if err != nil {
			t.Error(err)
		} else {
			t.Logf("----%d-----", tt.id)
			for _, v := range res {
				t.Log(v)
			}
		}
	}
}

func TestGetMenusByUserId(t *testing.T) {
	tests := []struct {
		id int64
	}{
		{1},   //bkzy
		{2},   //zjs_xxk
		{115}, //bkzy1
		{118}, //黄中省
	}
	for _, tt := range tests {
		res, err := GetMenusByUserId(tt.id)
		if err != nil {
			t.Error(err)
		} else {
			t.Logf("----%d-----", tt.id)
			for _, v := range res {
				t.Log(v)
			}
		}
	}
}

func TestGetTreeNodesByUserId(t *testing.T) {
	tests := []struct {
		id int64
	}{
		{1},   //bkzy
		{2},   //zjs_xxk
		{117}, //bkzy1
		{118}, //黄中省
	}
	for _, tt := range tests {
		res, err := GetTreeNodesByUserId(tt.id)
		if err != nil {
			t.Error(err)
		} else {
			t.Logf("----%d-----", tt.id)
			for _, v := range res {
				t.Logf("%v,%v", v, v.LevelCategory)
			}
		}
	}
}

func TestGetSubNodesByStageLevelCode(t *testing.T) {
	tests := []struct {
		levelcode string
	}{
		{"1-2-5-20"},
	}
	for _, tt := range tests {
		res, err := GetSubNodesByStageLevelCode(tt.levelcode)
		if err != nil {
			t.Error(err)
		} else {
			t.Logf("----%s-----", tt.levelcode)
			for _, v := range res {
				t.Logf("%v,%v", v, v.LevelCategory)
			}
		}
	}
}

func TestGetTaglistByNodeLevelCode(t *testing.T) {
	tests := []struct {
		levelcode string
	}{
		{"1-2-5-20"},
	}
	for _, tt := range tests {
		res, err := GetTaglistByNodeLevelCode(tt.levelcode)
		if err != nil {
			t.Error(err)
		} else {
			t.Logf("----%s-----", tt.levelcode)
			for _, v := range res {
				t.Logf("%+v", v)
			}
		}
	}
}

func TestGetTagsDataTableInfoByDcsId(t *testing.T) {
	tests := []struct {
		dcs_id int64
	}{
		{46},
	}
	for _, tt := range tests {
		res, err := GetTagsDataTableInfoByDcsId(tt.dcs_id)
		if err != nil {
			t.Error(err)
		} else {
			t.Logf("----%d-----", tt.dcs_id)
			for _, v := range res {
				t.Logf("%+v,%+v", v, v.Datatable)
			}
		}
	}
}

func TestGetTagAttributByTagId(t *testing.T) {
	tests := []struct {
		id int64
	}{
		{2660},
	}
	for _, tt := range tests {
		res, err := GetTagAttributByTagId(tt.id)
		if err != nil {
			t.Error(err)
		} else {
			t.Logf("%+v", res)
			t.Logf("%+v", res.Variable)
			t.Logf("%+v", res.Dcs)
			t.Logf("%+v", res.Stage)
		}
	}
}

func TestGetMonitorNodesByUserId(t *testing.T) {
	tests := []struct {
		id int64
	}{
		{1},   //bkzy
		{2},   //zjs_xxk
		{117}, //bkzy1
		{118}, //黄中省
	}
	for _, tt := range tests {
		res, err := GetMonitorNodesByUserId(tt.id)
		if err != nil {
			t.Error(err)
		} else {
			t.Logf("----%d-----", tt.id)
			for _, v := range res {
				t.Logf("%+v", v)
			}
		}
	}
}

func TestRegression(t *testing.T) {
	tests := []struct {
		tagYname  string
		tagXnames string
		begintime string
		endtime   string
		interval  int64
	}{
		{"sf8kt.x1_zjs_sfc_ps8kt_4-1_100-1_pv:1", "sf8kt.x1_zjs_sfc_ps8kt_4-1_100-1_pv:1,sf8kt.x1_zjs_sfc_ps8kt_4-1_100-1_sum:1", "2020-03-25 14:17:00", "2020-03-25 15:17:00", 60},
		{"sf8kt.x1_zjs_sfc_ps8kt_4-1_100-1_pv:1", "sf8kt.x1_zjs_sfc_ps8kt_4-1_100-1_pv:1", "2020-03-25 14:17:00", "2020-03-25 15:17:00", 60},
	}
	for _, tt := range tests {
		res, err := Regression(tt.tagYname, tt.tagXnames, tt.begintime, tt.endtime, tt.interval)
		if err != nil {
			t.Error(err)
		} else {
			t.Logf("%+v", res.Coeff)
			t.Logf("%+v", res.Ts)
			t.Logf("%+v", res.Vs)
		}
	}
}


func TestGetPeakValleyOfGoldenHisData(t *testing.T) {
	bt := "2020-02-01 08:00:00"
	et := "2020-02-01 09:00:00"
	tests := []struct {
		tagid     int64
		tagname   string
		key       string
		beginTime string
		endTime   string
	}{
		{0, "SF8KT.x1_zjs_sfc_ps8kt_4-3_100-14_pv:1", "pvs_peaksum", bt, et},
		{0, "SF8KT.x1_zjs_sfc_ps8kt_4-3_100-14_pv:1", "peakvalley", bt, et},
	}
	for _, tt := range tests {
		res, err := GetPeakValleyOfGoldenHisData(tt.tagid, tt.tagname, tt.key, tt.beginTime, tt.endTime)
		if err != nil {
			t.Log(err)
		}
		t.Log(res)
	}
}

func TestGetProcessTagIDByTagName(t *testing.T) {
	sum, err := GetProcessTagIDByTagName("x3_zjs_sfc_hb_25-86_31-446_47-493_freqSP:1")
	if err != nil {
		t.Error(err)
	}
	t.Log(sum)
}

func TestInsertTagKpi2CfgListByAppend(t *testing.T) {
	//res, err := InsertTagKpi2CfgListByAppend(15095, -2, "2020-05-20 00:00:00")
	//if err != nil {
	//	t.Log(err)
	//} else {
	//	t.Logf("成功追加了[%d]条KPI指标", res)
	//}
}

func TestGetSysRealTimeDatas(t *testing.T) {
	bt := "2019-12-23 10:29:35"
	et := "2019-12-25 00:16:00"
	tests := []struct {
		tagid     int64
		tagname   string
		tagtype   string
		beginTime string
		endTime   string
		readType  int
	}{
		{0, "ck1_asl_asl-xc1_SK1_ZXS1_S1-017_sr-ch_record:2", "41", bt, et, 1},
		{0, "ck1_asl_asl-xc1_SK1_ZXS1_S1-017_sr-ch_record:2", "41", bt, et, 0},
		{0, "ck1_asl_asl-xc1_SK1_ZXS1_S1-017_sr-ch_record:2", "check", bt, et, 0},
		{22, "", "41", bt, et, 0},
		{22, "", "check", bt, et, 0},
		{22, "ck1_asl_asl-xc1_SK1_ZXS1_S1-017_sr-ch_record:2", "41", bt, et, 0},
		{22, "ck1_asl_asl-xc1_SK1_ZXS1_S1-017_sr-ch_record:2", "check", bt, et, 0},
		{0, "ck1_asl_asl-xc1_SK1_ZXS1_S1-017_sr-ch_record:2", "41", bt, et, -1},
		{0, "ck1_asl_asl-xc1_SK1_ZXS1_S1-017_sr-ch_record:2", "check", bt, et, -1},
		{22, "", "41", bt, et, -1},
		{22, "", "check", bt, et, -1},
		{22, "ck1_asl_asl-xc1_SK1_ZXS1_S1-017_sr-ch_record:2", "41", bt, et, -1},
		{22, "ck1_asl_asl-xc1_SK1_ZXS1_S1-017_sr-ch_record:2", "check", bt, et, -1},
		{14952, "x7_asl_ps-manage9_ps-line167_pmer64_YM-PAP:1", "47", "2020-03-16 09:00:00", "2020-03-16 17:00:00", -1},
		{14952, "x7_asl_ps-manage9_ps-line167_pmer64_YM-PAP:1", "47", "2020-03-16 09:00:00", "2020-03-16 17:00:00", 0},
	}
	for i, tt := range tests {
		res, err := GetSysRealTimeDatas(tt.tagid, tt.tagname, tt.tagtype, tt.beginTime, tt.endTime, tt.readType)
		if err != nil {
			t.Log(err)
		}
		for n, v := range res {
			t.Logf("No.%d.%d %+v \n", i, n, v)
		}
	}
}

func TestGetSysRealTimeDataStatistic(t *testing.T) {
	bt := "2019-12-23 10:29:35"
	et := "2019-12-25 00:16:00"
	tests := []struct {
		tagid     int64
		tagname   string
		tagtype   string
		keyword   string
		beginTime string
		endTime   string
	}{
		{22, "", "41", "max", bt, et},
		{22, "", "check", "min", bt, et},
		{22, "ck1_asl_asl-xc1_SK1_ZXS1_S1-017_sr-ch_record:2", "41", "sd", bt, et},
		{14952, "x7_asl_ps-manage9_ps-line167_pmer64_YM-PAP:1", "47", "diff", "2020-03-16 09:00:00", "2020-03-16 17:00:00"},
	}
	for i, tt := range tests {
		res, err := GetSysRealTimeDataStatistic(tt.tagid, tt.tagname, tt.tagtype, tt.keyword, tt.beginTime, tt.endTime)
		if err != nil {
			t.Log(err)
		}
		t.Logf("No.%d %+v \n", i, res)
	}
}

func TestGetSysRealTimeDataStatisticByKey(t *testing.T) {
	bt := "2019-12-23 10:29:35"
	et := "2019-12-25 00:16:00"
	tests := []struct {
		tagid     int64
		tagname   string
		tagtype   string
		keyword   string
		beginTime string
		endTime   string
	}{
		{22, "", "41", "max", bt, et},
		{22, "", "check", "min", bt, et},
		{22, "ck1_asl_asl-xc1_SK1_ZXS1_S1-017_sr-ch_record:2", "41", "sd", bt, et},
		{14952, "x7_asl_ps-manage9_ps-line167_pmer64_YM-PAP:1", "47", "diff", "2020-03-16 09:00:00", "2020-03-16 17:00:00"},
	}
	for i, tt := range tests {
		res, err := GetSysRealTimeDataStatisticByKey(tt.tagid, tt.tagname, tt.tagtype, tt.keyword, tt.beginTime, tt.endTime)
		if err != nil {
			t.Log(err)
		}
		t.Logf("No.%d %+v \n", i, res)
	}
}
*/
/*
func TestGetSysDicIdByNameCode(t *testing.T) {
	tests := []struct {
		name string
		res  int64
	}{
		{"process", 32},
		{"goods", 33},
		{"patrol", 41},
		{"32", 32},
		{"33", 33},
		{"41", 41},
		{"25", 0},  //返回错误信息
		{"xxx", 0}, //返回错误信息
		{"235", 0}, //返回错误信息
	}
	for _, tt := range tests {
		res, err := GetSysDicIdByNameCode(tt.name)
		if err != nil {
			t.Log(err)
		}
		if res != tt.res {
			t.Errorf("错误,得到的结果:%d,期望的结果:%d", res, tt.res)
		}
	}
}

func TestGetGoodsValueByIDandTime(t *testing.T) {
	tstr, sum, err := GetGoodsValueByIDandTime(121, "2019-12-26 00:00:00")
	if err != nil {
		t.Error(err)
	}
	t.Log(tstr, sum)
}

func TestGetGoodsTagIDByTagName(t *testing.T) {
	sum, err := GetGoodsTagIDByTagName("x12_asl_asl-ckc_JDCJ_oth196_oil1")
	if err != nil {
		t.Error(err)
	}
	t.Log(sum)
}

func TestGetGoodsSumByID(t *testing.T) {
	tests := []struct {
		id int64
		st string
		ed string
	}{
		{121, "2019-12-18 00:00:00", "2019-12-26 00:00:00"},
		{501, "2020-06-03 09:00:00", "2020-06-03 17:00:00"},
	}
	for _, tt := range tests {
		sum, err := GetGoodsSumByID(tt.id, tt.st, tt.ed)
		if err != nil {
			t.Error(err)
		}
		t.Log(sum)
	}
}

func TestSaveKpiResultToDB(t *testing.T) {
	setMysqlMaxPreparedStmtCount(EngineCfgMsg.Sys.ResultDBAlias, 100000)
}

func TestSaveKpiResultToDB(t *testing.T) {
	kpi := new(CalcKpiResult)
	kpi.CalcEndingTime = "2019-12-14 08:00:00"
	kpi.InbaseTime = time.Now().Format(EngineCfgMsg.TimeFormat)
	kpi.KpiConfigListId = 245
	kpi.KpiPeriod = -1
	kpi.KpiTag = "tag tag tag"
	kpi.KpiValue = 123.456
	kpi.ProcessTaglistId = 432
	kpi.ProcessTaglistTagName = "name name name"

	tests := []struct {
		kpi *CalcKpiResult
		res bool
	}{
		{kpi, true},
	}

	for _, tt := range tests {
		err := SaveKpiResultToDB(tt.kpi)
		if err != nil {
			t.Log(err)
		}
	}
}

func TestSetKpiLastCalcTime(t *testing.T) {
	kpi := CalcKpiConfigList{Id: 1, LastCalcTime: "2019-01-01 00:16:00"}
	tests := []struct {
		kpi CalcKpiConfigList
		res bool
	}{
		{kpi, true},
	}

	for _, tt := range tests {
		res, err := SetKpiLastCalcTime(tt.kpi)
		if err != nil {
			t.Log(err)
		}
		if res != tt.res {
			t.Errorf("错误,得到的结果:%t,期望的结果:%t", res, tt.res)
		}
	}
}

func TestPrevTime(t *testing.T) {
	bt := "2001-01-15 00:00:00"
	st := "2019-01-01 00:16:00"
	tests := []struct {
		bstime  string
		sttime  string
		lasttie string
		shift   int64
		period  int64
		offset  int64
		done    bool
		bgtime  string
		edtime  string
	}{
		{bt, st, "2018-01-01 03:25:00", 8, -1, 0, true, "2019-01-01 00:00:00", "2019-01-01 01:00:00"},   //0
		{bt, st, "2019-07-21 03:25:00", 8, -1, 0, true, "2019-07-21 03:00:00", "2019-07-21 04:00:00"},   //1
		{bt, st, "2019-07-21 03:25:00", 8, -2, 0, true, "2019-07-21 00:00:00", "2019-07-21 08:00:00"},   //2
		{bt, st, "2019-07-21 08:00:00", 8, -2, 0, true, "2019-07-21 08:00:00", "2019-07-21 16:00:00"},   //3
		{bt, st, "2019-07-21 03:25:00", 8, -3, 0, true, "2019-07-21 00:00:00", "2019-07-22 00:00:00"},   //4
		{bt, st, "2019-07-22 00:00:00", 8, -3, 0, true, "2019-07-22 00:00:00", "2019-07-23 00:00:00"},   //5
		{bt, st, "2019-07-21 03:25:00", 8, -4, 0, true, "2019-07-15 00:00:00", "2019-08-15 00:00:00"},   //6
		{bt, st, "2019-08-15 00:00:00", 8, -4, 0, true, "2019-08-15 00:00:00", "2019-09-15 00:00:00"},   //7
		{bt, st, "2019-07-21 03:25:00", 8, -5, 0, true, "2019-07-15 00:00:00", "2019-10-15 00:00:00"},   //8
		{bt, st, "2019-10-15 00:00:00", 8, -5, 0, true, "2019-10-15 00:00:00", "2020-01-15 00:00:00"},   //9
		{bt, st, "2019-07-21 03:25:00", 8, -6, 0, true, "2019-01-15 00:00:00", "2020-01-15 00:00:00"},   //10
		{bt, st, "2020-01-15 00:00:00", 8, -6, 0, false, "2020-01-15 00:00:00", "2021-01-15 00:00:00"},  //11
		{bt, st, "2020-01-03 07:00:00", 8, -1, 0, true, "2020-01-03 07:00:00", "2020-01-03 08:00:00"},   //12
		{bt, st, "2020-01-03 07:00:00", 8, -2, 0, true, "2020-01-03 00:00:00", "2020-01-03 08:00:00"},   //13
		{bt, st, "2020-01-03 08:00:00", 8, -2, 0, true, "2020-01-03 08:00:00", "2020-01-03 16:00:00"},   //14,
		{bt, st, "2020-01-03 00:00:00", 8, -3, 0, true, "2020-01-03 00:00:00", "2020-01-04 00:00:00"},   //15,
		{bt, st, "2020-01-15 08:00:00", 8, -2, -60, true, "2020-01-15 07:00:00", "2020-01-15 15:00:00"}, //16,
		{bt, st, "2020-01-15 08:00:00", 8, -2, 166, true, "2020-01-15 08:00:00", "2020-01-15 16:00:00"}, //16,
	}

	for i, tt := range tests {
		ok, bgt, edt, err := PrevTime(tt.bstime, tt.sttime, tt.lasttie, tt.shift, tt.period, tt.offset)
		if err != nil {
			t.Log(err)
		}
		if ok != tt.done || bgt != tt.bgtime || edt != tt.edtime {
			t.Errorf("%d-错误,得到的结果:%t,%s,%s;期望的结果:%t,%s,%s", i, ok, bgt, edt, tt.done, tt.bgtime, tt.edtime)
		}
	}
}


func TestSaveBatchKpiResultToDB(t *testing.T) {
	kpi := new(CalcKpiResult)
	kpis := make([]CalcKpiResult, 10)
	for i := 0; i < 10; i++ {
		kpi.CalcEndingTime = time.Now().Format(_TIMEFOMAT)
		kpi.InbaseTime = time.Now().Format(_TIMEFOMAT)
		kpi.KpiConfigListId = 123
		kpi.KpiPeriod = -1
		kpi.KpiTag = fmt.Sprint(time.Now().UnixNano())
		kpi.KpiKey = "the key"
		kpi.KpiValue = 33.456
		kpi.ProcessTaglistId = 223
		kpi.ProcessTaglistTagName = "test test test tagname"
		kpis = append(kpis, *kpi)
	}

	tests := []struct {
		tbname string
		res    []CalcKpiResult
		cnt    int64
	}{
		{"testtable", kpis, 10},
	}

	for _, tt := range tests {
		exist, err := isTableExist(EngineCfgMsg.Sys.ResultDBAlias, tt.tbname)
		if err != nil {
			t.Error(err)
		}
		t.Logf("表存[%s]在吗,%v", tt.tbname, exist)
		if exist == false {
			ok, err := createResultTable(EngineCfgMsg.Sys.ResultDBAlias, tt.tbname)
			if err != nil {
				t.Error(err)
			}
			t.Logf("创建表[%s]成功了吗,%v", tt.tbname, ok)
		}

		cnt := saveKpiResultToDB(tt.tbname, tt.res)
		if cnt != tt.cnt {
			t.Errorf("错误,希望插入%v行,实际插入%v行", tt.cnt, cnt)
		}
	}
}

func TestCreateResultTable(t *testing.T) {
	tests := []struct {
		alias  string
		tbname string
		res    bool
	}{
		{"result", "testtable", true},
	}

	for _, tt := range tests {
		res, err := createResultTable(tt.alias, tt.tbname)
		if err != nil {
			t.Log(err)
		}
		if res != tt.res {
			t.Errorf("错误,期望值是%v,得到的值是%v", tt.res, res)
		}
	}
}

func TestIstableExist(t *testing.T) {
	tests := []struct {
		alias  string
		tbname string
		res    bool
	}{
		{"default", "act_ru_job", true},
		{"default", "act_ru_jobxx", false},
		{"result", "act_ru_job", false},
		{"result", "testtable", true},
	}

	for _, tt := range tests {
		res, err := isTableExist(tt.alias, tt.tbname)
		if err != nil {
			t.Log(err)
		}
		if res != tt.res {
			t.Errorf("错误,得到的结果:%t,期望的结果:%t", res, tt.res)
		}
	}
}


func TestGetKpiConfigInfo(t *testing.T) {
	tests := []struct {
		res int64
	}{
		{10000},
	}

	for _, tt := range tests {
		res, _, err := getKpiConfigInfo(1)
		if err != nil {
			t.Log(err)
		}
		t.Log(res)
		if res < tt.res {
			t.Errorf("错误,得到的结果:%d,期望的结果:%d", res, tt.res)
		}
	}
}

func TestGetWorkshopBaseTimeByTagID(t *testing.T) {
	tests := []struct {
		tagid     int
		bstime    string
		shifttime int
	}{
		{5599, "2001-01-01 00:00:00", 8},
	}

	for _, tt := range tests {
		bt, sht, err := GetWorkshopBaseTimeByTagID(tt.tagid)
		if err != nil {
			t.Log(err)
		}
		if bt != tt.bstime || sht != tt.shifttime {
			t.Errorf("错误,得到的结果:%s,%d;期望的结果:%s,%d", tt.bstime, tt.shifttime, bt, sht)
		}
	}
}


func TestGetTagHisSumValueFromGoldenByID(t *testing.T) {
	st := "2020/01/03 00:00:00"
	et := "2020/01/03 08:00:00"
	tests := []struct {
		tagid    int64
		tagname  string
		key      string
		sttime   string
		endtime  string
		tagvalue float64
	}{
		{5599, "SF8KT.x1_zjs_sfc_ps8kt_4-3_100-14_pv:1", "snapshot", "", "", 0.0},
		{5599, "SF8KT.x1_zjs_sfc_ps8kt_4-3_100-14_pv:1", "min", st, et, 0.0},
		{5599, "SF8KT.x1_zjs_sfc_ps8kt_4-3_100-14_pv:1", "max", st, et, 0.0},
		{5599, "SF8KT.x1_zjs_sfc_ps8kt_4-3_100-14_pv:1", "range", st, et, 0.0},
		{5599, "SF8KT.x1_zjs_sfc_ps8kt_4-3_100-14_pv:1", "total", st, et, 0.0},
		{5599, "SF8KT.x1_zjs_sfc_ps8kt_4-3_100-14_pv:1", "sum", st, et, 0.0},
		{5599, "SF8KT.x1_zjs_sfc_ps8kt_4-3_100-14_pv:1", "mean", st, et, 0.0},
		{5599, "SF8KT.x1_zjs_sfc_ps8kt_4-3_100-14_pv:1", "gt_l_total", st, et, 0.0},
		{5599, "SF8KT.x1_zjs_sfc_ps8kt_4-3_100-14_pv:1", "lte_l_total", st, et, 0.0},
	}

	for _, tt := range tests {
		res, _ := GetTagHisSumValueFromGoldenByID(tt.tagid, tt.tagname, tt.key, tt.sttime, tt.endtime)
		t.Log(tt.key, res)
	}
}

func TestGetTagFullNameByID(t *testing.T) {
	tests := []struct {
		tagid   int64
		tagname string
	}{
		{2660, "SFYKCQ.x1_zjs_sfc_ykcq_20-61_105-353_fault-open:1"},
		{2661, "SFYKCQ.x1_zjs_sfc_ykcq_20-61_105-353_flt-close:1"},
	}

	for _, tt := range tests {
		res, _ := GetTagFullNameByID(tt.tagid)
		if res != tt.tagname {
			t.Errorf("Wrong answer, got=%s, want=%s", res, tt.tagname)
		}
	}
}

func TestGetTagIDByFullName(t *testing.T) {
	tests := []struct {
		tagid   int64
		tagname string
	}{
		{2660, "SFYKCQ.x1_zjs_sfc_ykcq_20-61_105-353_fault-open:1"},
		{2661, "SFYKCQ.x1_zjs_sfc_ykcq_20-61_105-353_flt-close:1"},
		{10329, "sfgzl.x6_zjs_sfc_xy_24-60_29-970_current:1"},
	}

	for _, tt := range tests {
		res, _ := GetTagIDByFullName(tt.tagname)
		if res != tt.tagid {
			t.Errorf("Wrong answer, got=%d, want=%d", res, tt.tagid)
		}
	}
}

func TestGetTagListFieldValueByID(t *testing.T) {
	tests := []struct {
		tagid      int64
		field_name string
		value      string
	}{
		{5599, "dcs_id", "50"},
		{5599, "id", "5599"},
		{5599, "seq", "<nil>"},                 //会输出值为空信息
		{5599, "tag_id", "x1_1_1_1_3_14_pv:1"}, //
	}

	for _, tt := range tests {
		res, err := GetTagListFieldValueByID(tt.tagid, tt.field_name)
		if err != nil {
			t.Log(err)
		}
		if res != tt.value {
			t.Errorf("Wrong answer, got=%s, want=%s", res, tt.value)
		}
	}
}

func TestGetTagCategoryNameAndItemIDByID(t *testing.T) {
	tests := []struct {
		tagid        int64
		categoryneme string
		itemid       int64
	}{
		{5599, "ore_process_d_workstage_meter", 14},
	}

	for _, tt := range tests {
		name, id, err := GetTagCategoryNameAndItemIDByID(tt.tagid)
		if err != nil {
			t.Log(err)
		}
		if name != tt.categoryneme || id != tt.itemid {
			t.Errorf("Wrong answer, got=%s,%d, want=%s,%d", name, id, tt.categoryneme, tt.itemid)
		}
	}
}

func TestGetTableFieldValueByName(t *testing.T) {
	tests := []struct {
		tablename  string
		field_name string
		id         int64
		value      string
	}{
		{"ore_process_d_workstage_meter", "meter_class", 123, "106"},
		{"ore_process_d_workstage_meter", "meter_name_code", 123, "106-123"},
	}

	for _, tt := range tests {
		res, err := GetTableFieldValueByName(tt.tablename, tt.id, tt.field_name)
		if err != nil {
			t.Log(err)
		}
		if res != tt.value {
			t.Errorf("Wrong answer, got=%s, want=%s", res, tt.value)
		} else {
			t.Log(res, tt.value)
		}

	}
}

func TestGetTagCategoryItemFieldValueByID(t *testing.T) {
	tests := []struct {
		tagid      int64
		field_name string
		value      string
	}{
		{5599, "meter_class", "100"},
	}

	for _, tt := range tests {
		res, err := GetTagCategoryItemFieldValueByID(tt.tagid, tt.field_name)
		if err != nil {
			t.Log(err)
		}
		if res != tt.value {
			t.Errorf("Wrong answer, got=%s, want=%s", res, tt.value)
		}
	}
}

func TestGetTableNameByID(t *testing.T) {
	tests := []struct {
		table_id   int64
		table_name string
	}{
		{6, "ore_process_d_workstage_meter"},
		{10, "mine_dcs_station"},
		{19, "ore_process_d_taglist"},
	}

	for _, tt := range tests {
		res, err := GetTableNameByID(tt.table_id)
		if err != nil {
			t.Log(err)
		}
		if res != tt.table_name {
			t.Errorf("Wrong answer, got=%s, want=%s", res, tt.table_name)
		}
	}
}

func TestGetCategoryItemFieldValueByID(t *testing.T) {
	tests := []struct {
		table_id   int64
		row_id     int64
		field_name string
		value      string
	}{
		{19, 5600, "resource_id", "14"},
		{7, 70, "rated_power", "30"},
	}

	for _, tt := range tests {
		res, err := GetCategoryItemFieldValueByID(tt.table_id, tt.row_id, tt.field_name)
		if err != nil {
			t.Log(err)
		}
		if res != tt.value {
			t.Errorf("Wrong answer, got=%s, want=%s", res, tt.value)
		}
	}
}

func TestGetTagAttributByTagName(t *testing.T) {
	tests := []struct {
		name string
		id   int64
	}{
		{"x1_asl_asl-xc1_SK1_SF1_SY1-021_sum:1", 4554},
		{"x1_asl_asl-xc1_SK1_SF1_SY1-021_sum:x", 0},
	}

	for _, tt := range tests {
		tag := new(OreProcessDTaglist)
		tag.TagName = tt.name
		err := tag.GetTagAttributByTagName()
		if err != nil {
			t.Error(err)
		} else {
			t.Logf("%+v", tag)
			t.Logf("%+v", tag.Variable)
			t.Logf("%+v", tag.Dcs)
			t.Logf("%+v", tag.Stage)
		}
	}
}


func TestGetTableList(t *testing.T) {
	tb := new(DatatableInfo)
	tbs, err := tb.GetTableList()
	if err != nil {
		t.Error(err.Error())
	} else {
		t.Logf("%+v", tbs)
	}
}

func TestGetTableByName(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"Micbox1-xcdn"},
		{"Micbox3-m2mk"},
		{"XXX"},
	}
	for _, tt := range tests {
		tb := new(DatatableInfo)
		tb.TableName = tt.name
		err := tb.GetTableByName()
		if err != nil {
			t.Log(err.Error())
		} else {
			t.Logf("%+v", tb)
		}
	}
}

func TestUpdateGoldenId(t *testing.T) {
	tests := []struct {
		id   int64
		gdid int
	}{
		{1, 2},
	}
	for _, tt := range tests {
		tb := new(DatatableInfo)
		tb.Id = tt.id
		tb.GoldenTableId = tt.gdid
		err := tb.UpdateGoldenId()
		if err != nil {
			t.Log(err.Error())
		} else {
			t.Logf("%+v", tb)
		}
	}
}

// func TestSynchGoldenTableAndPlatTable(t *testing.T) {
// 	tb := new(DatatableInfo)
// 	total, newtb, err := tb.SynchGoldenTableAndPlatTable()
// 	if err != nil {
// 		t.Error(err.Error())
// 	} else {
// 		t.Logf("总表数:%d,新建数:%d", total, newtb)
// 	}
// }



func TestGetKpiTagTypeInfo(t *testing.T) {
	kpi := new(CalcKpiConfigList)
	tagmp, _ := kpi.GetTagTypeLists()
	for k, v := range tagmp {
		mp, err := kpi.GetKpiTagTypePeriodInfo(k)
		if err != nil {
			t.Error(err)
		} else {
			t.Log(k, v, mp)
		}
	}
}

func TestGetKpiDatas(t *testing.T) {
	tests := []struct {
		id       int64
		sttime   string
		edtime   string
		readtype int
	}{
		{10587, "2007-08-19 00:00:00", "2007-09-30 12:00:00", 0},
		{10587, "2007-08-19 00:00:00", "2007-09-30 12:00:00", -2},
		{10587, "2007-08-19 00:00:00", "2007-09-30 12:00:00", -1},
		{10587, "2007-08-19 00:00:00", "2007-09-30 12:00:00", 1},
		{10587, "2007-08-19 00:00:00", "2007-09-30 12:00:00", 2},
	}
	for _, tt := range tests {
		t.Logf("++++++++++++++++++%d++++++++++++++++++", tt.readtype)
		kpi := new(CalcKpiConfigList)
		kpi.Id = tt.id
		kpires, err := kpi.GetKpiDatas(tt.sttime, tt.edtime, tt.readtype)
		if err != nil {
			t.Log(err.Error())
		} else {
			for _, k := range kpires {
				t.Logf("%+v", k)
			}
		}
	}
}

func TestGetKpiConfig(t *testing.T) {
	tests := []struct {
		id int64
	}{
		{1},
		{10587},
		{100587},
	}
	for _, tt := range tests {
		kpi := new(CalcKpiConfigList)
		kpi.Id = tt.id
		err := kpi.GetKpiConfig()
		if err != nil {
			t.Error(err)
		} else {
			t.Log(kpi)
		}
	}
}


func TestInsertOrUpdateTreeNodesInCache(t *testing.T) {
	tests := []struct {
		idcode  string
		treestr string
		rootpid int64
	}{
		{"test", time.Now().Format("2006-01-02 15:03:04.000"), 1},
	}
	for _, tt := range tests {
		str := new(MsMiddleCorrelationString)
		rcnt, err := str.InsertOrUpdateTreeNodesInCache(tt.idcode, tt.treestr, tt.rootpid)
		if err != nil {
			t.Error(err)
		} else {
			t.Log(rcnt)
		}
	}
}

func TestGetGrafanaIP(t *testing.T) {
	tests := []struct {
		iswan bool
		name  string
	}{
		{false, "http://localhost:3000"},
		{true, "http://127.0.0.1:3000"},
	}
	for _, tt := range tests {
		dic := new(SysDictionary)
		name, err := dic.GetGrafanaHost(tt.iswan)
		if err != nil {
			t.Error(err)
		} else {
			if name != tt.name {
				t.Errorf("获取错误,期望值是:%s,获取到的值是:%s", tt.name, name)
			}
		}
	}
}
*/
func TestX(t *testing.T) {

}
