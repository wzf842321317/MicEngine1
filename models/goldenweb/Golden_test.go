package goldenweb

import (
	"testing"
	"time"
)

func TestTableList(t *testing.T) {
	add := GoldenAdd{Host: "zjs-t3.vicp.net", Port: 56732}

	tests := []TableListCmd{
		{"table", ""},
		{"table", "all"},
		{"table", "sf8kt"},
	}
	for _, tt := range tests {
		url, err := getUrl(add, tt)
		if err != nil {
			t.Error(err)
		} else {
			t.Log(url)
		}
		list, err := GetGoldenData(add, tt)
		if err != nil {
			t.Error(err)
		} else {
			t.Log(list)
		}
	}
}

func TestTimeFormat(t *testing.T) {
	tm, _ := time.ParseInLocation("2006-01-02 15:04:05", "2019-01-12 13:04:45", time.Local)
	tests := []struct {
		tstr string
		res  time.Time
	}{
		{"2019-01-12 13:04:45", tm},
		{"2019-01-12 13:04:45.000", tm},
		{"2019/01/12 13:04:45", tm},
		{"2019/01/12 13:04:45.000", tm},
		{"2019/1/12 13:04:45.000", tm},
	}

	for _, tt := range tests {
		res, err := timeFormat(tt.tstr)
		if err != nil {
			t.Log(err)
		}
		if res != tt.res {
			t.Errorf("转换错误，得到的值是[%v],期望值是[%v]", res, tt.res)
		}
	}
}

func TestGetHistoryData(t *testing.T) {
	add := GoldenAdd{Host: "zjs-t3.vicp.net", Port: 56732}
	var cmd HistoryCmd
	cmd.TagName = append(cmd.TagName, "SF2WT.x1_zjs_sfc_ps20kt_4-n_100-1239_pv:1")
	cmd.TagName = append(cmd.TagName, "SF2WT.x1_zjs_sfc_ps20kt_4-5_100-1239_sum:1")
	cmd.BeginTime = "2019-01-02 12:00:00"
	cmd.EndTime = "2019-01-02 13:00:00"
	var resmap HisDataMap
	tests := []struct {
		add GoldenAdd
		cmd HistoryCmd
		res HisDataMap
	}{
		{add, cmd, resmap},
	}

	for _, tt := range tests {
		res, err := getHistoryData(tt.add, tt.cmd)
		if err != nil {
			t.Log(err)
		}
		t.Log(len(res), res)
	}
}

func TestGetHisPointData(t *testing.T) {
	add := GoldenAdd{Host: "zjs-t3.vicp.net", Port: 56732}
	var cmd HistoryCmd
	cmd.TagName = append(cmd.TagName, "SF2WT.x1_zjs_sfc_ps20kt_4-5_100-1239_pv:1")
	cmd.TagName = append(cmd.TagName, "SF2WT.x1_zjs_sfc_ps20kt_4-5_100-1239_sum:1")
	cmd.TimePoint = "2020-01-02 12:00:00"
	var resmap HisDataMap
	tests := []struct {
		add GoldenAdd
		cmd HistoryCmd
		res HisDataMap
	}{
		{add, cmd, resmap},
	}

	for _, tt := range tests {
		res, err := getHisPointData(tt.add, tt.cmd)
		if err != nil {
			t.Log(err)
		}
		t.Logf("是否存在：%t", IsExistItem(tt.cmd.TagName[0], res))
		t.Log(len(res), res)
	}
}

func TestGetSnapShorData(t *testing.T) {
	add := GoldenAdd{Host: "zjs-t3.vicp.net", Port: 56732}
	var cmd SnapshotCmd
	cmd.TagName = append(cmd.TagName, "SF2WT.x1_zjs_sfc_ps20kt_4-5_100-1239_pv:1")
	cmd.TagName = append(cmd.TagName, "SF2WT.x1_zjs_sfc_ps20kt_4-5_100-1239_sum:1")
	cmd.TagName = append(cmd.TagName, "T21DCS.x2_zjs_tec_ycj_7-142_41-757_106-868_enable:1")
	var resmap SnapDataMap
	tests := []struct {
		add GoldenAdd
		cmd SnapshotCmd
		res SnapDataMap
	}{
		{add, cmd, resmap},
	}

	for _, tt := range tests {
		res, err := getSnapShotData(tt.add, tt.cmd)
		if err != nil {
			t.Log(err)
		}
		if len(res) < len(tt.cmd.TagName) {
			t.Errorf("数据长度不匹配,得到的是%d,期望的是%d", len(res), len(tt.cmd.TagName))
		}
	}
}

func TestGetHistorySummaryDataExi(t *testing.T) {
	add := GoldenAdd{Host: "zjs-t3.vicp.net", Port: 56732}
	var cmd HisSumCmd
	cmd.TagName = append(cmd.TagName, "SF2WT.x1_zjs_sfc_ps20kt_4-5_100-1239_pv:1")
	cmd.TagName = append(cmd.TagName, "SF2WT.x1_zjs_sfc_ps20kt_4-5_100-1239_sum:1")
	cmd.BeginTime = "2020-01-01 12:00:00"
	cmd.EndTime = "2020-01-01 13:00:00"
	cmd.DataType = "total"
	var resmap HisSumDataMapExi
	tests := []struct {
		add GoldenAdd
		cmd HisSumCmd
		res HisSumDataMapExi
	}{
		{add, cmd, resmap},
	}

	for _, tt := range tests {
		res, err := getHistorySummaryDataExi(tt.add, tt.cmd)
		if err != nil {
			t.Log(err)
		}
		for k, r := range res {
			t.Log(k, r)
		}

		if len(res) < len(tt.cmd.TagName) {
			t.Errorf("数据长度不匹配,得到的是%d,期望的是%d", len(res), len(tt.cmd.TagName))
		}
	}
}
