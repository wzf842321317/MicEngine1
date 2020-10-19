package models

import (
	"testing"
)

/*
func TestGoldenGetServerTime(t *testing.T) {
	gdt, err := GoldenGetServerTime()
	if err != nil {
		t.Errorf("错误:[%s]", err.Error())
	} else {
		t.Logf("服务器时间:%s", gdt)
	}
}

func TestGoldenGetApiVersion(t *testing.T) {
	gdt, err := GoldenGetApiVersion()
	if err != nil {
		t.Errorf("错误:[%s]", err.Error())
	} else {
		t.Logf("服务器版本:%s", gdt)
	}
}

func TestGoldenGetSnapShotData(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"sf8kt.x1_zjs_sfc_ps8kt_4-1_100-1_pv:1"},
		{"sf8kt.x1_zjs_sfc_ps8kt_4-1_100-1_sum:1"},
		{""},
	}
	for i, tt := range tests {
		val, dtime, dq, err := GoldenGetSnapShotData(tt.name)
		if err != nil {
			t.Logf("第%d行错误:[%s]", i, err.Error())
		} else {
			t.Logf("行数:%d,数据:%f,时间:%s,质量:%d", i, val, dtime, dq)
		}
	}
}

func TestGoldenGetSingleStatisticsData(t *testing.T) {
	tests := []struct {
		name    string
		bgtime  string
		endtime string
		key     string
	}{
		{"sf8kt.x1_zjs_sfc_ps8kt_4-1_100-1_pv:1", "2020-05-15 15:00:00", "2020-05-15 16:00:00", "point"},
		{"sf8kt.x1_zjs_sfc_ps8kt_4-1_100-1_pv:1", "2020-05-15 15:00:00", "2020-05-15 16:00:00", "snapshot"},
		{"sf8kt.x1_zjs_sfc_ps8kt_4-1_100-1_pv:1", "2020-05-15 15:00:00", "2020-05-15 16:00:00", "min"},
		{"sf8kt.x1_zjs_sfc_ps8kt_4-1_100-1_pv:1", "2020-05-15 15:00:00", "2020-05-15 16:00:00", "max"},
		{"sf8kt.x1_zjs_sfc_ps8kt_4-1_100-1_pv:1", "2020-05-15 15:00:00", "2020-05-15 16:00:00", "range"},
		{"sf8kt.x1_zjs_sfc_ps8kt_4-1_100-1_pv:1", "2020-05-15 15:00:00", "2020-05-15 16:00:00", "total"},
		{"sf8kt.x1_zjs_sfc_ps8kt_4-1_100-1_pv:1", "2020-05-15 15:00:00", "2020-05-15 16:00:00", "sum"},
		{"sf8kt.x1_zjs_sfc_ps8kt_4-1_100-1_pv:1", "2020-05-15 15:00:00", "2020-05-15 16:00:00", "mean"},
		{"sf8kt.x1_zjs_sfc_ps8kt_4-1_100-1_pv:1", "2020-05-15 15:00:00", "2020-05-15 16:00:00", "poweravg"},
		{"sf8kt.x1_zjs_sfc_ps8kt_4-1_100-1_pv:1", "2020-05-15 15:00:00", "2020-05-15 16:00:00", "diff"},
		{"sf8kt.x1_zjs_sfc_ps8kt_4-1_100-1_pv:1", "2020-05-15 15:00:00", "2020-05-15 16:00:00", "duration"},
		{"sf8kt.x1_zjs_sfc_ps8kt_4-1_100-1_pv:1", "2020-05-15 15:00:00", "2020-05-15 16:00:00", "pointcnt"},
		{"sf8kt.x1_zjs_sfc_ps8kt_4-1_100-1_pv:1", "2020-05-15 15:00:00", "2020-05-15 16:00:00", "risingcnt"},
		{"sf8kt.x1_zjs_sfc_ps8kt_4-1_100-1_pv:1", "2020-05-15 15:00:00", "2020-05-15 16:00:00", "fallingcnt"},
		{"sf8kt.x1_zjs_sfc_ps8kt_4-1_100-1_pv:1", "2020-05-15 15:00:00", "2020-05-15 16:00:00", "sd"},
		{"sf8kt.x1_zjs_sfc_ps8kt_4-1_100-1_pv:1", "2020-05-15 15:00:00", "2020-05-15 16:00:00", "stddev"},
		{"sf8kt.x1_zjs_sfc_ps8kt_4-1_100-1_pv:1", "2020-05-15 15:00:00", "2020-05-15 16:00:00", "se"},
		{"sf8kt.x1_zjs_sfc_ps8kt_4-1_100-1_pv:1", "2020-05-15 15:00:00", "2020-05-15 16:00:00", "ske"},
		{"sf8kt.x1_zjs_sfc_ps8kt_4-1_100-1_pv:1", "2020-05-15 15:00:00", "2020-05-15 16:00:00", "kur"},
		{"sf8kt.x1_zjs_sfc_ps8kt_4-1_100-1_pv:1", "2020-05-15 15:00:00", "2020-05-15 16:00:00", "mode"},
		{"sf8kt.x1_zjs_sfc_ps8kt_4-1_100-1_pv:1", "2020-05-15 15:00:00", "2020-05-15 16:00:00", "median"},
	}
	for i, tt := range tests {
		val, err := GoldenGetSingleStatisticsData(tt.name, tt.key, tt.bgtime, tt.endtime)
		if err != nil {
			t.Errorf("第%d行错误:[%s]", i, err.Error())
		} else {
			t.Logf("行数:%d,指标:%s,数据:%f", i, tt.key, val)
		}
	}
}

func TestGoldenGetStatisticsData(t *testing.T) {
	tests := []struct {
		name    string
		bgtime  string
		endtime string
		key     string
	}{
		{"sf8kt.x1_zjs_sfc_ps8kt_4-1_100-1_pv:1", "2020-05-15 15:00:00", "2020-05-15 16:00:00", "point"},
	}
	for i, tt := range tests {
		val, err := GoldenGetStatisticsData(tt.name, tt.bgtime, tt.endtime)
		if err != nil {
			t.Errorf("第%d行错误:[%s]", i, err.Error())
		} else {
			t.Logf("行数:%d,数据:%+v", i, val)
		}
	}
}

func TestGoldenGetSingleStatisticsDataBatch(t *testing.T) {
	tests := []struct {
		names   []string
		bgtime  string
		endtime string
		key     string
	}{
		{[]string{"sf8kt.x1_zjs_sfc_ps8kt_4-1_100-1_pv:1", "sf8kt.x1_zjs_sfc_ps8kt_4-1_100-1_sum:1"}, "2020-05-15 15:00:00", "2020-05-15 16:00:00", "point"},
		{[]string{"sf8kt.x1_zjs_sfc_ps8kt_4-1_100-1_pv:1", "sf8kt.x1_zjs_sfc_ps8kt_4-1_100-1_sum:1"}, "2020-05-15 15:00:00", "2020-05-15 16:00:00", "snapshot"},
		{[]string{"sf8kt.x1_zjs_sfc_ps8kt_4-1_100-1_pv:1", "sf8kt.x1_zjs_sfc_ps8kt_4-1_100-1_sum:1"}, "2020-05-15 15:00:00", "2020-05-15 16:00:00", "min"},
		{[]string{"sf8kt.x1_zjs_sfc_ps8kt_4-1_100-1_pv:1", "sf8kt.x1_zjs_sfc_ps8kt_4-1_100-1_sum:1"}, "2020-05-15 15:00:00", "2020-05-15 16:00:00", "max"},
		{[]string{"sf8kt.x1_zjs_sfc_ps8kt_4-1_100-1_pv:1", "sf8kt.x1_zjs_sfc_ps8kt_4-1_100-1_sum:1"}, "2020-05-15 15:00:00", "2020-05-15 16:00:00", "range"},
		{[]string{"sf8kt.x1_zjs_sfc_ps8kt_4-1_100-1_pv:1", "sf8kt.x1_zjs_sfc_ps8kt_4-1_100-1_sum:1"}, "2020-05-15 15:00:00", "2020-05-15 16:00:00", "total"},
		{[]string{"sf8kt.x1_zjs_sfc_ps8kt_4-1_100-1_pv:1", "sf8kt.x1_zjs_sfc_ps8kt_4-1_100-1_sum:1"}, "2020-05-15 15:00:00", "2020-05-15 16:00:00", "sum"},
		{[]string{"sf8kt.x1_zjs_sfc_ps8kt_4-1_100-1_pv:1", "sf8kt.x1_zjs_sfc_ps8kt_4-1_100-1_sum:1"}, "2020-05-15 15:00:00", "2020-05-15 16:00:00", "mean"},
		{[]string{"sf8kt.x1_zjs_sfc_ps8kt_4-1_100-1_pv:1", "sf8kt.x1_zjs_sfc_ps8kt_4-1_100-1_sum:1"}, "2020-05-15 15:00:00", "2020-05-15 16:00:00", "poweravg"},
		{[]string{"sf8kt.x1_zjs_sfc_ps8kt_4-1_100-1_pv:1", "sf8kt.x1_zjs_sfc_ps8kt_4-1_100-1_sum:1"}, "2020-05-15 15:00:00", "2020-05-15 16:00:00", "diff"},
		{[]string{"sf8kt.x1_zjs_sfc_ps8kt_4-1_100-1_pv:1", "sf8kt.x1_zjs_sfc_ps8kt_4-1_100-1_sum:1"}, "2020-05-15 15:00:00", "2020-05-15 16:00:00", "duration"},
		{[]string{"sf8kt.x1_zjs_sfc_ps8kt_4-1_100-1_pv:1", "sf8kt.x1_zjs_sfc_ps8kt_4-1_100-1_sum:1"}, "2020-05-15 15:00:00", "2020-05-15 16:00:00", "pointcnt"},
		{[]string{"sf8kt.x1_zjs_sfc_ps8kt_4-1_100-1_pv:1", "sf8kt.x1_zjs_sfc_ps8kt_4-1_100-1_sum:1"}, "2020-05-15 15:00:00", "2020-05-15 16:00:00", "risingcnt"},
		{[]string{"sf8kt.x1_zjs_sfc_ps8kt_4-1_100-1_pv:1", "sf8kt.x1_zjs_sfc_ps8kt_4-1_100-1_sum:1"}, "2020-05-15 15:00:00", "2020-05-15 16:00:00", "fallingcnt"},
		{[]string{"sf8kt.x1_zjs_sfc_ps8kt_4-1_100-1_pv:1", "sf8kt.x1_zjs_sfc_ps8kt_4-1_100-1_sum:1"}, "2020-05-15 15:00:00", "2020-05-15 16:00:00", "sd"},
		{[]string{"sf8kt.x1_zjs_sfc_ps8kt_4-1_100-1_pv:1", "sf8kt.x1_zjs_sfc_ps8kt_4-1_100-1_sum:1"}, "2020-05-15 15:00:00", "2020-05-15 16:00:00", "stddev"},
		{[]string{"sf8kt.x1_zjs_sfc_ps8kt_4-1_100-1_pv:1", "sf8kt.x1_zjs_sfc_ps8kt_4-1_100-1_sum:1"}, "2020-05-15 15:00:00", "2020-05-15 16:00:00", "se"},
		{[]string{"sf8kt.x1_zjs_sfc_ps8kt_4-1_100-1_pv:1", "sf8kt.x1_zjs_sfc_ps8kt_4-1_100-1_sum:1"}, "2020-05-15 15:00:00", "2020-05-15 16:00:00", "ske"},
		{[]string{"sf8kt.x1_zjs_sfc_ps8kt_4-1_100-1_pv:1", "sf8kt.x1_zjs_sfc_ps8kt_4-1_100-1_sum:1"}, "2020-05-15 15:00:00", "2020-05-15 16:00:00", "kur"},
		{[]string{"sf8kt.x1_zjs_sfc_ps8kt_4-1_100-1_pv:1", "sf8kt.x1_zjs_sfc_ps8kt_4-1_100-1_sum:1"}, "2020-05-15 15:00:00", "2020-05-15 16:00:00", "mode"},
		{[]string{"sf8kt.x1_zjs_sfc_ps8kt_4-1_100-1_pv:1", "sf8kt.x1_zjs_sfc_ps8kt_4-1_100-1_sum:1"}, "2020-05-15 15:00:00", "2020-05-15 16:00:00", "median"},
	}
	for i, tt := range tests {
		val, err := GoldenGetSingleStatisticsDataBatch(tt.key, tt.bgtime, tt.endtime, tt.names...)
		if err != nil {
			t.Errorf("第%d行错误:[%s]", i, err.Error())
		} else {
			for k, v := range val {
				t.Logf("tag:%s,数据:%+v", k, v)
			}
		}
	}
}

func TestGoldenGetAnalogComparedTotal(t *testing.T) {
	tests := []struct {
		id      int64
		name    string
		bgtime  string
		endtime string
		key     string
	}{
		{6019, "sf8kt.x1_zjs_sfc_ps8kt_4-1_100-1_pv:1", "2020-05-15 15:00:00", "2020-05-15 16:00:00", "gt_l_total"},
		{6019, "sf8kt.x1_zjs_sfc_ps8kt_4-1_100-1_pv:1", "2020-05-15 15:00:00", "2020-05-15 16:00:00", "lt_l_total"},
	}
	for i, tt := range tests {
		val, err := GoldenGetAnalogComparedTotal(tt.id, tt.name, tt.key, tt.bgtime, tt.endtime)
		if err != nil {
			t.Errorf("第%d行错误:[%s]", i, err.Error())
		} else {
			t.Logf("行数:%d,数据:%+v", i, val)
		}
	}
}

func TestGoldenGetHistory(t *testing.T) {
	tests := []struct {
		names   []string
		bgtime  string
		endtime string
		key     string
	}{
		{[]string{"sf8kt.x1_zjs_sfc_ps8kt_4-1_100-1_pv:1", "sf8kt.x1_zjs_sfc_ps8kt_4-1_100-1_sum:1"}, "2020-05-15 15:00:00", "2020-05-15 16:00:00", "gt_l_total"},
		{[]string{}, "2020-05-15 15:00:00", "2020-05-15 16:00:00", "lt_l_total"},
	}
	for i, tt := range tests {
		val, err := GoldenGetHistory(tt.bgtime, tt.endtime, tt.names...)
		t.Logf("序号%d==========", i)
		if err != nil {
			t.Errorf("第%d行错误:[%s]", i, err.Error())
		} else {
			for k, v := range val {
				t.Logf("tag:%s,数据:%+v", k, v)
			}
		}
	}
}

func TestGoldenGetHistoryInterval(t *testing.T) {
	bgt := "2020-06-30 10:00:00"
	edt := "2020-07-01 00:00:00"
	tests := []struct {
		interval int64
		names    []string
		bgtime   string
		endtime  string
	}{
		{0, []string{"sf8kt.x1_zjs_sfc_ps8kt_4-1_100-1_pv:1", "sf8kt.x1_zjs_sfc_ps8kt_4-1_100-1_sum:1"}, bgt, edt},
		//{60, []string{"sf8kt.x1_zjs_sfc_ps8kt_4-1_100-1_pv:1", "sf8kt.x1_zjs_sfc_ps8kt_4-1_100-1_sum:1"}, bgt, edt},
	}
	for i, tt := range tests {
		micgd := new(MicGolden)
		_, gdhis, err := micgd.GoldenGetHistoryInterval(tt.bgtime, tt.endtime, tt.interval, tt.names...)
		t.Logf("序号%d==========", i)
		if err != nil {
			t.Errorf("第%d行错误:[%s]", i, err.Error())
		} else {
			for k, v := range gdhis {
				t.Logf("tag:%s,数据量:%d,错误信息:%s,是否继续:%t------------", k, len(v.HisRtsd), v.Err, v.Continue)
				for _, hd := range v.HisRtsd {
					t.Logf("%s %+v", time.Unix(hd.Time/1e3, 0).Format("2006-01-02 15:04:05"), hd)
				}
			}
		}
	}
}

func TestGoldenGetHistorySummary(t *testing.T) {
	tests := []struct {
		names   []string
		bgtime  string
		endtime string
		key     string
	}{
		{[]string{"sf8kt.x1_zjs_sfc_ps8kt_4-1_100-1_pv:1", "sf8kt.x1_zjs_sfc_ps8kt_4-1_100-1_sum:1"}, "2020-05-15 15:00:00", "2020-05-15 16:00:00", "gt_l_total"},
		{[]string{"sf8kt.x1_zjs_sfc_ps8kt_4-1_100-1_pv:1", "sf8kt.x1_zjs_sfc_ps8kt_4-1_100-1_sum:1"}, "2020-05-15 15:00:00", "2020-05-15 16:00:00", "lt_l_total"},
		{[]string{}, "2020-05-15 15:00:00", "2020-05-15 16:00:00", "lt_l_total"},
		{[]string{}, "2020-05-15 15:00:00", "2020-05-15 16:00:00", "lt_l_total"},
	}
	for i, tt := range tests {
		val, err := GoldenGetHistorySummary(tt.bgtime, tt.endtime, tt.names...)
		t.Logf("序号%d==========", i)
		if err != nil {
			t.Errorf("第%d行错误:[%s]", i, err.Error())
		} else {
			for k, v := range val {
				t.Logf("tag:%s,数据:%+v", k, v)
			}
		}
	}
}

func TestGoldenGetTables(t *testing.T) {
	tests := []struct {
		selector int
	}{
		{0},
		{1},
		{2},
		{3},
		{4},
	}
	for i, tt := range tests {
		micgd := new(MicGolden)
		val, err := micgd.GoldenGetTables(tt.selector)
		t.Logf("序号%d==========", i)
		if err != nil {
			t.Errorf("第%d行错误:[%s]", i, err.Error())
		} else {
			t.Logf("数据:%+v", val)
		}
	}
}

func TestGoldenGetTablePropertyByTableName(t *testing.T) {
	tests := []struct {
		names []string
	}{
		{[]string{"sf8kt", "demo", "ddd", "t3cdcs"}},
		{[]string{}},
	}
	for i, tt := range tests {
		val, err := GoldenGetTablePropertyByTableName(tt.names...)
		t.Logf("序号%d==========", i)
		if err != nil {
			t.Errorf("第%d行错误:[%s]", i, err.Error())
		} else {
			for k, v := range val {
				t.Logf("tag:%s,数据:%+v", k, v)
			}
		}
	}
}

func TestGoldenGetTagNameListInTables(t *testing.T) {
	tests := []struct {
		names []string
	}{
		{[]string{"sf8kt", "test", "t3cdcs", "ddd"}},
		{[]string{}},
	}
	for i, tt := range tests {
		val, err := GoldenGetTagNameListInTables(tt.names...)
		t.Logf("序号%d==========", i)
		if err != nil {
			t.Errorf("第%d行错误:[%s]", i, err.Error())
		} else {
			for k, v := range val {
				t.Logf("tag:%s,数据:%+v", k, v)
			}
		}
	}
}

func TestGoldenGetTagListInTables(t *testing.T) {
	tests := []struct {
		names []string
	}{
		{[]string{"sf8kt", "test", "t3cdcs", "ddd"}},
		{[]string{}},
	}
	for i, tt := range tests {
		val, err := GoldenGetTagListInTables(tt.names...)
		t.Logf("序号%d==========", i)
		if err != nil {
			t.Errorf("第%d行错误:[%s]", i, err.Error())
		} else {
			for k, v := range val {
				t.Logf("tag:%s,数据:%+v", k, v)
			}
		}
	}
}

func TestGoldenGetHistorySinglePoint(t *testing.T) {
	tests := []struct {
		mode    int
		names   []string
		bgtime  string
		endtime string
		key     string
	}{
		{0, []string{"sf8kt.x1_zjs_sfc_ps8kt_4-1_100-1_pv:1", "sf8kt.x1_zjs_sfc_ps8kt_4-1_100-1_sum:1"}, "2020-05-15 15:00:00", "2020-05-15 16:00:00", "gt_l_total"},
		{1, []string{"sf8kt.x1_zjs_sfc_ps8kt_4-1_100-1_pv:1", "sf8kt.x1_zjs_sfc_ps8kt_4-1_100-1_sum:1"}, "2020-05-15 15:00:00", "2020-05-15 16:00:00", "lt_l_total"},
		{2, []string{"sf8kt.x1_zjs_sfc_ps8kt_4-1_100-1_pv:1", "sf8kt.x1_zjs_sfc_ps8kt_4-1_100-1_sum:1"}, "2020-05-15 15:00:00", "2020-05-15 16:00:00", "lt_l_total"},
		{3, []string{"sf8kt.x1_zjs_sfc_ps8kt_4-1_100-1_pv:1", "sf8kt.x1_zjs_sfc_ps8kt_4-1_100-1_sum:1"}, "2020-05-15 15:00:00", "2020-05-15 16:00:00", "lt_l_total"},
	}
	for i, tt := range tests {
		val, err := GoldenGetHistorySinglePoint(tt.mode, tt.bgtime, tt.names...)
		t.Logf("序号%d==========", i)
		if err != nil {
			t.Errorf("第%d行错误:[%s]", i, err.Error())
		} else {
			for k, v := range val {
				t.Logf("tag:%s,时间:%s,数据:%+v", k, Millisecond2Time(v.Time), v)
			}
		}
	}
}

func TestGoldenGetTagPointInfoByName(t *testing.T) {
	tests := []struct {
		names []string
	}{
		{[]string{"sf8kt.x1_zjs_sfc_ps8kt_4-1_100-1_pv:1", "sf8kt.x1_zjs_sfc_ps8kt_4-1_100-1_sum:1"}},
		{[]string{"sf8kt.x1_zjs_sfc_ps8kt_4-1_100-1_pv:1", "sf8kt.x1_zjs_sfc_ps8kt_4-1_100-1_sum:1"}},
		{[]string{"sf8kt.x3_zjs_sfc_ps8kt_4-1_35-4_47-49_run:1", "sf8kt.webinsert_point"}},
		{[]string{"t3cdcs.demo_float_1", "t3cdcs.demo_int_1"}},
		{[]string{}},
	}
	for i, tt := range tests {
		val, err := GoldenGetTagPointInfoByName(tt.names...)
		t.Logf("序号%d==========", i)
		if err != nil {
			t.Errorf("第%d行错误:[%s]", i, err.Error())
		} else {
			for k, v := range val {
				t.Logf("tag:%s,数据:%+v", k, v)
			}
		}
	}
}

func TestSetGoldenSnapShots(t *testing.T) {
	rand.Seed(time.Now().UnixNano()) //设定随机数种子
	tests := []struct {
		names  []string
		values []float64
		dqual  []int
		times  []string
	}{
		{
			[]string{"sf8kt.x1_zjs_sfc_ps8kt_4-1_100-1_pv:1", "sf8kt.x1_zjs_sfc_ps8kt_4-1_100-1_sum:1"},
			[]float64{rand.Float64() * 10, rand.Float64() * 20},
			[]int{0, 0},
			[]string{TimeFormat(time.Now()), TimeFormat(time.Now())},
		},
		{
			[]string{"sf8kt.x1_zjs_sfc_ps8kt_4-1_100-1_pv:1", "sf8kt.x1_zjs_sfc_ps8kt_4-1_100-1_pv:1"},
			[]float64{rand.Float64() * 10, rand.Float64() * 20},
			[]int{0, 0},
			[]string{TimeFormat(time.Now()), TimeFormat(time.Now().Add(1 * time.Minute))},
		},
	}
	for i, tt := range tests {
		err := SetGoldenSnapShots(tt.names, tt.values, tt.dqual, tt.times)
		t.Logf("序号%d==========", i)
		if err != nil {
			t.Errorf("第%d行错误:[%s]", i, err.Error())
		} else {
			t.Logf("第%d行写入快照成功", i)
		}
	}
}

func TestGoldenGetSnapShotMap(t *testing.T) {
	tests := []struct {
		names []string
	}{
		{[]string{"sf8kt.x1_zjs_sfc_ps8kt_4-1_100-1_pv:1", "sf8kt.x1_zjs_sfc_ps8kt_4-1_100-1_sum:1"}},
		{[]string{"Micbox1-2.x1_asl_asl-xc1_MF1_MKⅠ3_MY1-021_sp:1", "Micbox1-2.x1_asl_asl-xc1_MF1_MKⅡ3_MY1-045_sp:1"}},
		{[]string{"Micbox1-2.x3_asl_asl-xc1_MF1_TBYX3_F1-003_FD1-003_run:1", "micbox1-2.x3_asl_asl-xc1_MF1_TBYX3_F1-006_FD1-006_run:1"}},
	}
	for i, tt := range tests {
		datas, err := GoldenGetSnapShotMap(tt.names...)
		if err != nil {
			t.Logf("第%d行错误:[%s]", i, err.Error())
		} else {
			for tag, data := range datas {
				t.Logf("tag:%s,数据:%f,时间:%s,质量:%d,错误信息:%s", tag, data.Rtsd.Value, Millisecond2Time(data.Rtsd.Time), data.Rtsd.Quality, data.Err)
			}
		}
	}
}

func TestInsertAlarm(t *testing.T) {
	tests := []struct {
		dt          string
		tag_id      int64
		tag_desc    string
		tag_s_desc  string
		tag_value   float64
		limit_value float64
		alarm_s     int
	}{
		{TimeFormat(time.Now()), 4337, "desc", "shortdesc", 1234.567, 1000.00, 1},
	}
	for _, tt := range tests {
		alrm := new(CalcTagAlarm)
		tag := new(OreProcessDTaglist)
		tag.Id = tt.tag_id
		alrm.Datatime = tt.dt
		alrm.Tag = tag
		alrm.TagDesc = tt.tag_desc
		alrm.TagShotDesc = tt.tag_s_desc
		alrm.TagValue = tt.tag_value
		alrm.LimitValue = tt.limit_value
		alrm.AlarmStatus = tt.alarm_s
		alrm.InsertAlarm()

		err := alrm.GetLastAlartMsg()
		if err != nil {
			t.Error(err)
		} else {
			t.Logf("%+v", alrm)
		}
	}
}

func TestGoldenGetPoints(t *testing.T) {
	err := GDDB.GoldenDB.GetHandel()
	if err != nil {
		t.Error(err)
	}
	defer GDDB.GoldenDB.ReleaseHandel()

	num, err := GDDB.GetGoldenPoints()
	if err != nil {
		t.Error(err)
	} else {
		t.Logf("需要报警的标签点数:%d,庚顿信息%+v", num, GDDB.GoldenDB.Tables)
	}
	err = GDDB.GetGoldenSnapShotsAndComp()
	if err != nil {
		t.Error(err)
	}
}

func TestGetLastAlartMsg(t *testing.T) {
	tests := []struct {
		tag_id int64
	}{
		{4445},
		{4481},
		{4499},
	}
	for _, tt := range tests {
		alrm := new(CalcTagAlarm)
		tag := new(OreProcessDTaglist)
		tag.Id = tt.tag_id
		alrm.Tag = tag
		err := alrm.GetLastAlartMsg()
		if err != nil {
			t.Error(err)
		} else {
			t.Logf("%+v", alrm.Tag)
		}
	}
}

func TestAlarmOptionMask(t *testing.T) {
	tests := []struct {
		mask   int
		option int
		res    int
	}{
		{1, 3, 1},
		{1, 6, 0},
		{8, 62, 8},
		{8, 54, 0},
	}
	for _, tt := range tests {
		res := alarmOption(tt.mask, tt.option)
		t.Log(res)
	}
}


func TestTable(t *testing.T) {
	tests := []struct {
		id   int
		name string
		desc string
	}{
		{0, "demo8", "测试表"},
	}
	for _, tt := range tests {
		err := GoldenTableRemove(tt.name)
		if err != nil {
			t.Error(err)
		}
	}
}
*/
func TestY(t *testing.T) {
}
