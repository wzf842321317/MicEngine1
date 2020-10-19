package models

import (
	"testing"
)

func TestScript(t *testing.T) {}

/*
func TestReportCalc(t *testing.T) {
	tests := []struct {
		basetime  string
		bgtime    string
		endtime   string
		name      string
		tplname   string
		shifthour int64
		period    int64
	}{
		{"2006-01-01 00:00:00",
			"2020-05-05 08:00:00",
			"2020-05-05 16:00:00",
			"报表测试", "micform.xlsx", 8, -2},
		{"2006-01-01 01:00:00",
			"2020-05-01 01:00:00",
			"2020-05-02 01:00:00",
			"一选厂电耗表", "一选厂电耗表.xlsx", 8, -3},
	}
	var rpt CalcKpiReportList
	rpt.Id = 1
	rpt.OffsetMinutes = 10
	rpt.ResultUrl = "../data/report/form/"
	rpt.Status = 1
	for i, tt := range tests {
		rpt.BaseTime = tt.basetime
		rpt.ShiftHour = tt.shifthour
		rpt.Name = tt.name
		rpt.TemplateUrl = "../data/report/template/" + tt.tplname
		rpt.Period = tt.period
		num, name, err := rpt.ReportCalc(tt.bgtime, tt.endtime)
		if err == nil {
			t.Logf("第%d行 执行结果:%d,%s", i, num, name)
		} else {
			t.Log(err)
		}
	}
}

func TestScriptCompile(t *testing.T) {
	tests := []struct {
		scriptStr string
		res       []string
	}{
		{"fc(total)", []string{"fc(total)"}},
		{"this.fc(sd)", []string{"this.fc(sd)"}},
		{"tag(1234).fc(max)", []string{"tag(1234).fc(max)"}},
		{"fc(total,2019-10-24 03:04:05,2020-10-24 03:04:05)", []string{"fc(total,2019-10-24 03:04:05,2020-10-24 03:04:05)"}},
		{"fc(total,0,2020-10-24 03:04:05)", []string{"fc(total,0,2020-10-24 03:04:05)"}},
		{"fc(total,0,now)", []string{"fc(total,0,now)"}},
		{"fc(total,table(ore_process_d_workstage_meter).row(14).field(field_name),now)", []string{"fc(total,table(ore_process_d_workstage_meter).row(14).field(field_name),now)"}},
		{"fc(total,table(ore_process_d_workstage_meter).row(14).field(field_name),now) + fc(total)", []string{"fc(total,table(ore_process_d_workstage_meter).row(14).field(field_name),now)", "+", "fc(total)"}},
		{"fc(total,table(ore_process_d_workstage_meter).row(14).field(field_name),now) + fc(total) * table(3).row(245).field(field_name)", []string{"fc(total,table(ore_process_d_workstage_meter).row(14).field(field_name),now)", "+", "fc(total)", "*", "table(3).row(245).field(field_name)"}},
		{"tag(123).fc(total,table(ore_process_d_workstage_meter).row(14).field(field_name),now) + fc(total) * table(3).row(245).field(field_name)", []string{"tag(123).fc(total,table(ore_process_d_workstage_meter).row(14).field(field_name),now)", "+", "fc(total)", "*", "table(3).row(245).field(field_name)"}},
		{"tag(123).fc(total,table(equipment_sub_part).row(2).field(start_use_time),now)>=table(equipment_sub_part).row(2).field(work_timespan)", []string{"tag(123).fc(total,table(equipment_sub_part).row(2).field(start_use_time),now)", ">=", "table(equipment_sub_part).row(2).field(work_timespan)"}},
		{"this.field(field_name)", []string{"this.field(field_name)"}},
		{"tag(1234).field(field_name)", []string{"tag(1234).field(field_name)"}},
		{"this.table.row.field(field_name)", []string{"this.table.row.field(field_name)"}},
		{"tag(1234).table.row.field(field_name)", []string{"tag(1234).table.row.field(field_name)"}},
		{"tag(x3_asl_asl-xc1_MF1_FXYJZB3_Y1-010_YD1-008_run:1).table.row.field(field_name)", []string{"tag(x3_asl_asl-xc1_MF1_FXYJZB3_Y1-010_YD1-008_run:1).table.row.field(field_name)"}},
		{"table(3).row(245).field(field_name)", []string{"table(3).row(245).field(field_name)"}},
		{"table(ore_process_d_workstage_meter).row(14).field(field_name)", []string{"table(ore_process_d_workstage_meter).row(14).field(field_name)"}}, //报错信息输出
		{"gdsum(123)+fc(total)", []string{"gdsum(123)", "+", "fc(total)"}},
		{"gdsum()+fc(total)", []string{"gdsum()", "+", "fc(total)"}},
		{"gdsum(x12_asl_asl-ckc_JDCJ_oth196_oil1)+fc(total)", []string{"gdsum(x12_asl_asl-ckc_JDCJ_oth196_oil1)", "+", "fc(total)"}},
		{"srtd(goods,sum)+fc(total)", []string{"srtd(goods,sum)", "+", "fc(total)"}},
		{"tag(355).srtd(goods,sum)+fc(total)", []string{"tag(355).srtd(goods,sum)", "+", "fc(total)"}},
		{"srtd(goods,sum,2020-02-01 12:00:00,2020-02-28 01:02:03)+fc(total)", []string{"srtd(goods,sum,2020-02-01 12:00:00,2020-02-28 01:02:03)", "+", "fc(total)"}},
		{"tag(355).srtd(goods,sum,2020-02-01 12:00:00,2020-02-28 01:02:03)+fc(total)", []string{"tag(355).srtd(goods,sum,2020-02-01 12:00:00,2020-02-28 01:02:03)", "+", "fc(total)"}},
		{"tag(355).srtd(goods,sum,2020-02-01 12:00:00,table(ore_process_d_workstage_meter).row(14).field(field_name))+fc(total)", []string{"tag(355).srtd(goods,sum,2020-02-01 12:00:00,table(ore_process_d_workstage_meter).row(14).field(field_name))", "+", "fc(total)"}},
		{"( (fc(total) + fc(max))*(fc(mean) + fc(min)) )/(srtd(goods,diff)+ this.field(field_name))", []string{"( ", "(", "fc(total)", "+", "fc(max)", ")", "*", "(", "fc(mean)", "+", "fc(min)", ")", " )", "/", "(", "srtd(goods,diff)", "+", "this.field(field_name)", ")"}},
		{"((tag(310).srtd(check,mean) - tag(322).srtd(check,mean)) * (tag(319).srtd(check,mean) - tag(323).srtd(check,mean)) - (tag(318).srtd(check,mean) - tag(322).srtd(check,mean)) * (tag(311).srtd(check,mean) - tag(323).srtd(check,mean))) / ((tag(314).srtd(check,mean) - tag(322).srtd(check,mean)) * (tag(319).srtd(check,mean) - tag(323).srtd(check,mean)) - (tag(318).srtd(check,mean) - tag(322).srtd(check,mean)) * (tag(315).srtd(check,mean) - tag(323).srtd(check,mean)))",
			[]string{"(", "(", "tag(310).srtd(check,mean)", "-", "tag(322).srtd(check,mean)", ")", "*", "(", "tag(319).srtd(check,mean)", "-", "tag(323).srtd(check,mean)", ")", "-", "(", "tag(318).srtd(check,mean)", "-", "tag(322).srtd(check,mean)", ")", "*", "(", "tag(311).srtd(check,mean)", "-", "tag(323).srtd(check,mean)", ")", ")", "/", "(", "(", "tag(314).srtd(check,mean)", "-", "tag(322).srtd(check,mean)", ")", "*", "(", "tag(319).srtd(check,mean)", "-", "tag(323).srtd(check,mean)", ")", "-", "(", "tag(318).srtd(check,mean)", "-", "tag(322).srtd(check,mean)", ")", "*", "(", "tag(315).srtd(check,mean)", "-", "tag(323).srtd(check,mean)", ")", ")"}},
		{"srtd(check,mean)*(tag(10734).fc(diff)*(100 - tag(346).srtd(check,mean))/100)/100",
			[]string{"srtd(check,mean)", "*", "(", "tag(10734).fc(diff)", "*", "(", "100", "-", "tag(346).srtd(check,mean)", ")", "/", "100", ")", "/", "100"}},
		{"srtd(check,sum,t_add(0,-8h),t_add(0,5h))", []string{"srtd(check,sum,t_add(0,-8h),t_add(0,5h))"}},
		{"kpi(x12_asl_asl-ckc_JDCJ_oth196_oil1,sum,t_add(0,-8h),t_add(0,5h))", []string{"kpi(x12_asl_asl-ckc_JDCJ_oth196_oil1,sum,t_add(0,-8h),t_add(0,5h))"}},
		{"kpi(10237,sum,t_add(0,-8h),t_add(0,5h))", []string{"kpi(10237,sum,t_add(0,-8h),t_add(0,5h))"}},
		{"kpi(10237,sum,t_add(2020-02-01 12:00:00,-8h),t_add(0,5h))", []string{"kpi(10237,sum,t_add(2020-02-01 12:00:00,-8h),t_add(0,5h))"}},
		{"kpi(10237,sum,2020-02-01 12:00:00,2020-02-05 12:00:00)", []string{"kpi(10237,sum,2020-02-01 12:00:00,2020-02-05 12:00:00)"}},
		{"kpi(10237,sum,2020-02-01 12:00:00,2020-02-05 12:00:00)+kpi(10237,sum,t_add(0,-8h),t_add(0,5h))", []string{"kpi(10237,sum,2020-02-01 12:00:00,2020-02-05 12:00:00)", "+", "kpi(10237,sum,t_add(0,-8h),t_add(0,5h))"}},
		{"date(0,Y)", []string{"", ""}},
		{"data(0,Y)", []string{"", ""}},
		{"date(0,Y)date(0,w)", []string{"", ""}},
	}
	script := new(Script)
	for i, tt := range tests {
		script.ScriptStr = tt.scriptStr
		res := script.Compile()
		t.Logf("-------第%d行-------", i)
		if len(res) == 0 {
			t.Log("ok")
		} else {
			for _, err := range res {
				t.Log(err.Error())
			}
		}
	}
}

func TestScriptRun(t *testing.T) {
	st := "2020/05/02 01:30:00"
	et := "2020/05/03 01:30:00"
	basetime := "2006-01-01 01:30:00"
	name1 := "x1_asl_asl-xc1_MF1_MKⅠ3_MY1-002_meas-value:1"
	name2 := "ck1_asl_asl-xc1_XWXT_QY5_sap863_sr-ch_record:1__team_PowerAvg"
	tests := []struct {
		thisid    int64
		shifthour int64
		tagname   string
		scriptStr string
		bgt       string
		edt       string
		basetime  string
		result    float64
	}{
		{0, 8, name1, "kpi(10734,sum)", st, et, basetime, 2355.105006},
		{0, 8, name1, "kpi(10734,endpoint,0,t_add(0,-16h))", st, et, basetime, 84786.979611},
		{0, 8, name2, "kpi(10734,endpoint,0,t_add(0,-8h))", st, et, basetime, 8478697.961099},
		{0, 8, name2, "kpi(10734,endpoint)", st, et, basetime, 119},
		{0, 8, name2, "kpi(10734,endpoint)+kpi(10734,endpoint,0,t_add(0,-8h))", st, et, basetime, 119},
		{4867, 8, "", "fc(snapshot)", st, et, basetime, 119},
		{0, 8, "", "tag(4867).fc(snapshot)", st, et, basetime, 119},
		{0, 8, "", "tag(micbox9-liupo.x2_asl_asl-ckc_YK_LP-01_YK-13_W1164_PV:1).fc(snapshot)", st, et, basetime, 119},
		{14545, 8, name1, "tag(micbox9-liupo.x2_asl_asl-ckc_YK_LP-01_YK-13_W1164_PV:1).fc(snapshot)", st, et, basetime, 119},
		{0, 8, "", "date(0,t)", st, et, basetime, 0},
		{0, 8, "", "date(0,Y)", st, et, basetime, 0},
		{0, 8, "", "date(0,M)", st, et, basetime, 0},
		{0, 8, "", "date(0,D)", st, et, basetime, 0},
		{0, 8, "", "date(0,y)", st, et, basetime, 0},
		{0, 8, "", "date(0,m)", st, et, basetime, 0},
		{0, 8, "", "date(0,d)", st, et, basetime, 0},
		{0, 8, "", "date(0,w)", st, et, basetime, 0},
		{0, 8, "", "date(bgos,t)", st, et, basetime, 0},
		{0, 8, "", "date(bgos,Y)", st, et, basetime, 0},
		{0, 8, "", "date(bgos,M)", st, et, basetime, 0},
		{0, 8, "", "date(bgos,D)", st, et, basetime, 0},
		{0, 8, "", "date(bgos,y)", st, et, basetime, 0},
		{0, 8, "", "date(bgos,m)", st, et, basetime, 0},
		{0, 8, "", "date(bgos,d)", st, et, basetime, 0},
		{0, 8, "", "select(count(id)).from(kpi_result).where(id > 11373636 and id <= 11373640)", st, et, basetime, 0},
		{0, 8, "", "select(id).from(kpi_result).where(id > 11373636 and id <= 11373640)", st, et, basetime, 0},           //结果是多行则去第一行
		{0, 8, "", "select(*).from(kpi_result).where(id > 11373636 and id <= 11373640).as(json)", st, et, basetime, 0},   //
		{0, 8, "", "select(*).from(kpi_result).where(id > 11373636 and id <= 11373640).as(map)", st, et, basetime, 0},    //
		{0, 8, "", "select(*).from(kpi_result).where(id > 11373636 and id <= 11373640).as(string)", st, et, basetime, 0}, //
		{0, 8, "", "select(*).from(kpi_result).where(id > 11373636 and id <= 11373640).as(value)", st, et, basetime, 0},  //
		{0, 8, "", "select(*).from(kpi_result).where(id > 11373636 and id <= 11373640)", st, et, basetime, 0},            //
		{0, 8, "", "select(id,kpi_value).from(kpi_result).where(id > 11373636 and id <= 11373640)", st, et, basetime, 0}, //结果是多列会取最后一列第一个数
		{0, 8, "", "select(id,tag_name).from(kpi_result).where(id > 11373636 and id <= 11373640)", st, et, basetime, 0},  //结果不可转换为浮点数则输出0
		{0, 8, "", "select(count(id)).from(kpi_result).where(tag_name='x7_asl_ps-manage5_ps-line51_pmer13_YM-PAP:2').timecolumn(calc_ending_time).timefilter(0,0)", st, et, basetime, 0},              //
		{0, 8, "", "select(count(id)).from(kpi_result).where(tag_name='x7_asl_ps-manage5_ps-line51_pmer13_YM-PAP:2').timecolumn(calc_ending_time)", st, et, basetime, 0},                              //
		{0, 8, "", "select(id).from(kpi_result).where(tag_name='x7_asl_ps-manage5_ps-line51_pmer13_YM-PAP:2').timecolumn(calc_ending_time).limit(5)", st, et, basetime, 0},                            //
		{0, 8, "", "select(id).from(kpi_result).where(tag_name='x7_asl_ps-manage5_ps-line51_pmer13_YM-PAP:2').timecolumn(calc_ending_time).orderby(id desc).limit(5).as(value)", st, et, basetime, 0}, //
	}
	script := new(Script)
	for i, tt := range tests {
		script.Id = 0
		script.BaseTime = tt.basetime
		script.BeginTime = tt.bgt
		script.EndTime = tt.edt
		script.ShiftHour = tt.shifthour
		script.ScriptStr = tt.scriptStr
		script.MainTagId = tt.thisid
		script.MainTagFullName = tt.tagname

		res, _, err := script.Run()
		if err != nil {
			t.Log("No.", i, " ", err)
		} else {
			t.Logf("第%d行:%v", i, res)
		}
	}
}

func TestSQLScriptRun(t *testing.T) {
	st := "2020/05/02 01:30:00"
	et := "2020/05/03 01:30:00"
	basetime := "2006-01-01 01:30:00"
	tests := []struct {
		thisid    int64
		shifthour int64
		tagname   string
		scriptStr string
		bgt       string
		edt       string
		basetime  string
		result    float64
	}{
		{0, 8, "", "select(count(id)).from(sys_real_data).where(id > 1000 and id <= 2000)", st, et, basetime, 0},
		{0, 8, "", "select(id).from(sys_real_data).where(id > 1000 and id <= 2000)", st, et, basetime, 0},                                                                                        //结果是多行则去第一行
		{0, 8, "", "select(*).from(sys_real_data).where(id > 1000 and id <= 2000).as(json)", st, et, basetime, 0},                                                                                //
		{0, 8, "", "select(*).from(sys_real_data).where(id > 1000 and id <= 2000).as(map)", st, et, basetime, 0},                                                                                 //
		{0, 8, "", "select(*).from(sys_real_data).where(id > 1000 and id <= 2000).as(string)", st, et, basetime, 0},                                                                              //
		{0, 8, "", "select(*).from(sys_real_data).where(id > 1000 and id <= 2000).as(value)", st, et, basetime, 0},                                                                               //
		{0, 8, "", "select(*).from(sys_real_data).where(id > 1000 and id <= 2000)", st, et, basetime, 0},                                                                                         //
		{0, 8, "", "select(id,value).from(sys_real_data).where(id > 1000 and id <= 2000)", st, et, basetime, 0},                                                                                  //结果是多列会取最后一列第一个数
		{0, 8, "", "select(id,tag_name).from(sys_real_data).where(id > 1000 and id <= 2000)", st, et, basetime, 0},                                                                               //结果不可转换为浮点数则输出0
		{0, 8, "", "select(count(id)).from(sys_real_data).where(tag_name='x7_asl_ps-manage5_ps-line51_pmer13_YM-PAP:2').timecolumn(datatime).timefilter(0,0)", st, et, basetime, 0},              //
		{0, 8, "", "select(count(id)).from(sys_real_data).where(tag_name='x7_asl_ps-manage5_ps-line51_pmer13_YM-PAP:2').timecolumn(datatime)", st, et, basetime, 0},                              //
		{0, 8, "", "select(id).from(sys_real_data).where(tag_name='x7_asl_ps-manage5_ps-line51_pmer13_YM-PAP:2').timecolumn(datatime).limit(5)", st, et, basetime, 0},                            //
		{0, 8, "", "select(id).from(sys_real_data).where(tag_name='x7_asl_ps-manage5_ps-line51_pmer13_YM-PAP:2').timecolumn(datatime).orderby(id desc).limit(5).as(value)", st, et, basetime, 0}, //
		{0, 8, "", "select(AVG(value)).from(sys_real_data).where(tag_id=29 AND datatime > '2019-12-16 01:00:00' AND datatime <= '2019-12-23 16:00:00').as(json)", st, et, basetime, 0},
		{0, 8, "", "select(AVG(value)).from(sys_real_data).where(tag_id=29 AND datatime > '2019-12-16 01:00:00' AND datatime <= '2019-12-23 16:00:00').as(sql)", st, et, basetime, 0},
		{0, 8, "", "select(AVG(value)).from(sys_real_data).where(tag_id=29 AND datatime > '2019-12-16 01:00:00' AND datatime <= '2019-12-23 16:00:00')", st, et, basetime, 0},
	}
	script := new(Script)
	for i, tt := range tests {
		script.Id = 0
		script.BaseTime = tt.basetime
		script.BeginTime = tt.bgt
		script.EndTime = tt.edt
		script.ShiftHour = tt.shifthour
		script.ScriptStr = tt.scriptStr
		script.MainTagId = tt.thisid
		script.MainTagFullName = tt.tagname

		res, _, err := script.Run()
		if err != nil {
			t.Log("No.", i, " ", err)
		} else {
			t.Logf("第%d行:%v", i, res)
		}
	}
}


func TestScriptStr2TokenArr(t *testing.T) {
	tests := []struct {
		scriptStr string
		res       []string
	}{
		{"fc(total)", []string{"fc(total)"}},
		{"this.fc(sd)", []string{"this.fc(sd)"}},
		{"tag(1234).fc(max)", []string{"tag(1234).fc(max)"}},
		{"fc(total,2019-10-24 03:04:05,2020-10-24 03:04:05)", []string{"fc(total,2019-10-24 03:04:05,2020-10-24 03:04:05)"}},
		{"fc(total,0,2020-10-24 03:04:05)", []string{"fc(total,0,2020-10-24 03:04:05)"}},
		{"fc(total,0,now)", []string{"fc(total,0,now)"}},
		{"fc(total,table(ore_process_d_workstage_meter).row(14).field(field_name),now)", []string{"fc(total,table(ore_process_d_workstage_meter).row(14).field(field_name),now)"}},
		{"fc(total,table(ore_process_d_workstage_meter).row(14).field(field_name),now) + fc(total)", []string{"fc(total,table(ore_process_d_workstage_meter).row(14).field(field_name),now)", "+", "fc(total)"}},
		{"fc(total,table(ore_process_d_workstage_meter).row(14).field(field_name),now) + fc(total) * table(3).row(245).field(field_name)", []string{"fc(total,table(ore_process_d_workstage_meter).row(14).field(field_name),now)", "+", "fc(total)", "*", "table(3).row(245).field(field_name)"}},
		{"this.fc(total,table(ore_process_d_workstage_meter).row(14).field(field_name),now) + fc(total) * table(3).row(245).field(field_name)", []string{"this.fc(total,table(ore_process_d_workstage_meter).row(14).field(field_name),now)", "+", "fc(total)", "*", "table(3).row(245).field(field_name)"}},
		{"this.field(field_name)", []string{"this.field(field_name)"}},
		{"tag(1234).field(field_name)", []string{"tag(1234).field(field_name)"}},
		{"this.table.row.field(field_name)", []string{"this.table.row.field(field_name)"}},
		{"tag(1234).table.row.field(field_name)", []string{"tag(1234).table.row.field(field_name)"}},
		{"tag(x3_asl_asl-xc1_MF1_FXYJZB3_Y1-010_YD1-008_run:1).table.row.field(field_name)", []string{"tag(x3_asl_asl-xc1_MF1_FXYJZB3_Y1-010_YD1-008_run:1).table.row.field(field_name)"}},
		{"table(3).row(245).field(field_name)", []string{"table(3).row(245).field(field_name)"}},
		{"table(ore_process_d_workstage_meter).row(14).field(field_name)", []string{"table(ore_process_d_workstage_meter).row(14).field(field_name)"}}, //报错信息输出
		{"gdsum(123)+fc(total)", []string{"gdsum(123)", "+", "fc(total)"}},
		{"gdsum()+fc(total)", []string{"gdsum()", "+", "fc(total)"}},
		{"gdsum(x12_asl_asl-ckc_JDCJ_oth196_oil1)+fc(total)", []string{"gdsum(x12_asl_asl-ckc_JDCJ_oth196_oil1)", "+", "fc(total)"}},
		{"srtd(goods,sum)+fc(total)", []string{"srtd(goods,sum)", "+", "fc(total)"}},
		{"tag(355).srtd(goods,sum)+fc(total)", []string{"tag(355).srtd(goods,sum)", "+", "fc(total)"}},
		{"srtd(goods,sum,2020-02-01 12:00:00,2020-02-28 01:02:03)+fc(total)", []string{"srtd(goods,sum,2020-02-01 12:00:00,2020-02-28 01:02:03)", "+", "fc(total)"}},
		{"tag(355).srtd(goods,sum,2020-02-01 12:00:00,2020-02-28 01:02:03)+fc(total)", []string{"tag(355).srtd(goods,sum,2020-02-01 12:00:00,2020-02-28 01:02:03)", "+", "fc(total)"}},
		{"tag(355).srtd(goods,sum,2020-02-01 12:00:00,table(ore_process_d_workstage_meter).row(14).field(field_name))+fc(total)", []string{"tag(355).srtd(goods,sum,2020-02-01 12:00:00,table(ore_process_d_workstage_meter).row(14).field(field_name))", "+", "fc(total)"}},
		{"( (fc(total) + fc(max))*(fc(mean) + fc(min)) )/(srtd(goods,diff)+ this.field(field_name))", []string{"( ", "(", "fc(total)", "+", "fc(max)", ")", "*", "(", "fc(mean)", "+", "fc(min)", ")", " )", "/", "(", "srtd(goods,diff)", "+", "this.field(field_name)", ")"}},
		{"((tag(310).srtd(check,mean) - tag(322).srtd(check,mean)) * (tag(319).srtd(check,mean) - tag(323).srtd(check,mean)) - (tag(318).srtd(check,mean) - tag(322).srtd(check,mean)) * (tag(311).srtd(check,mean) - tag(323).srtd(check,mean))) / ((tag(314).srtd(check,mean) - tag(322).srtd(check,mean)) * (tag(319).srtd(check,mean) - tag(323).srtd(check,mean)) - (tag(318).srtd(check,mean) - tag(322).srtd(check,mean)) * (tag(315).srtd(check,mean) - tag(323).srtd(check,mean)))",
			[]string{"(", "(", "tag(310).srtd(check,mean)", "-", "tag(322).srtd(check,mean)", ")", "*", "(", "tag(319).srtd(check,mean)", "-", "tag(323).srtd(check,mean)", ")", "-", "(", "tag(318).srtd(check,mean)", "-", "tag(322).srtd(check,mean)", ")", "*", "(", "tag(311).srtd(check,mean)", "-", "tag(323).srtd(check,mean)", ")", ")", "/", "(", "(", "tag(314).srtd(check,mean)", "-", "tag(322).srtd(check,mean)", ")", "*", "(", "tag(319).srtd(check,mean)", "-", "tag(323).srtd(check,mean)", ")", "-", "(", "tag(318).srtd(check,mean)", "-", "tag(322).srtd(check,mean)", ")", "*", "(", "tag(315).srtd(check,mean)", "-", "tag(323).srtd(check,mean)", ")", ")"}},
		{"srtd(check,mean)*(tag(10734).fc(diff)*(100 - tag(346).srtd(check,mean))/100)/100",
			[]string{"srtd(check,mean)", "*", "(", "tag(10734).fc(diff)", "*", "(", "100", "-", "tag(346).srtd(check,mean)", ")", "/", "100", ")", "/", "100"}},
		{"srtd(check,sum,t_add(0,-8h),t_add(0,5h))", []string{"srtd(check,sum,t_add(0,-8h),t_add(0,5h))"}},
		{"kpi(x12_asl_asl-ckc_JDCJ_oth196_oil1,sum,t_add(0,-8h),t_add(0,5h))", []string{"kpi(x12_asl_asl-ckc_JDCJ_oth196_oil1,sum,t_add(0,-8h),t_add(0,5h))"}},
		{"kpi(10237,sum,t_add(0,-8h),t_add(0,5h))", []string{"kpi(10237,sum,t_add(0,-8h),t_add(0,5h))"}},
		{"kpi(10237,sum,t_add(2020-02-01 12:00:00,-8h),t_add(0,5h))", []string{"kpi(10237,sum,t_add(2020-02-01 12:00:00,-8h),t_add(0,5h))"}},
		{"kpi(10237,sum,2020-02-01 12:00:00,2020-02-05 12:00:00)", []string{"kpi(10237,sum,2020-02-01 12:00:00,2020-02-05 12:00:00)"}},
		{"kpi(10237,sum,2020-02-01 12:00:00,2020-02-05 12:00:00) + kpi(10237,sum,t_add(0,-8h),t_add(0,5h))", []string{"kpi(10237,sum,2020-02-01 12:00:00,2020-02-05 12:00:00)", "+", "kpi(10237,sum,t_add(0,-8h),t_add(0,5h))"}},
		{"select(DISTINCT( middle.id ) AS Id).from(sampling_manage)", []string{"select(DISTINCT( middle.id ) AS Id).from(sampling_manage)"}},
		{"select(*).from(sampling_manage).where(tag_name = 'ck13_asl_asl-xc1_MF1_azr192_MF1-120_sr-ch_record:17')", []string{"select(*).from(sampling_manage).where(tag_name = 'ck13_asl_asl-xc1_MF1_azr192_MF1-120_sr-ch_record:17')"}},
		{"select(kpi_value).from(kpi_result).where(id = 123).timecolumn(calc_ending_time)", []string{"select(kpi_value).from(kpi_result).where(id = 123).timecolumn(calc_ending_time)"}},
		{"select(kpi_value).from(kpi_result).where(id = 123).timecolumn(calc_ending_time).timefilter(0,0)", []string{"select(kpi_value).from(kpi_result).where(id = 123).timecolumn(calc_ending_time).timefilter(0,0)"}},
		{"select(kpi_value).from(kpi_result).where(id = 123).timecolumn(calc_ending_time).timefilter(now,0).groupby(id)", []string{"select(kpi_value).from(kpi_result).where(id = 123).timecolumn(calc_ending_time).timefilter(now,0).groupby(id)"}},
		{"select(kpi_value).from(kpi_result).where(id = 123).timecolumn(calc_ending_time).timefilter(now,0).orderby(id).as(string)", []string{"select(kpi_value).from(kpi_result).where(id = 123).timecolumn(calc_ending_time).timefilter(now,0).orderby(id).as(string)"}},
		{"select(kpi_value+id).from(kpi_result).where(id = 123).timecolumn(calc_ending_time).timefilter(now,0).groupby(id).orderby(id).limit(123)+fc(min)", []string{"select(kpi_value+id).from(kpi_result).where(id = 123).timecolumn(calc_ending_time).timefilter(now,0).groupby(id).orderby(id).limit(123)", "+", "fc(min)"}},
	}
	script := new(Script)
	for _, tt := range tests {
		script.ScriptStr = tt.scriptStr
		res := script.scriptStr2TokenArr()
		for i, r := range res {
			if r != tt.res[i] {
				t.Errorf("Wrong answer, \n got =%v, \n want=%v", res, tt.res)
				break
			}
		}
	}
}

func TestAnalysisScript(t *testing.T) {
	tests := []struct {
		scr string
		res string
	}{
		{"tag(14621).fc(total,table(equipment_sub_part).row(2).field(start_use_time),now)", "6"},
	}
	script := new(Script)
	script.BaseTime = "2020-01-01 00:00:00"
	script.BeginTime = "2020-06-04 08:00:00"
	script.EndTime = "2020-06-04 16:00:00"
	script.ShiftHour = 12

	for _, tt := range tests {
		script.ScriptStr = tt.scr
		rst, _, err := script.Run()
		if err != nil {
			t.Error(err.Error())
		} else {
			t.Logf("%s=%s", script.ScriptStr, rst)
		}
	}
}

func TestExtractDateTime(t *testing.T) {
	tests := []struct {
		scr string
		res string
	}{
		{"date(0,m)", "6"},
		{"date(t_add(0,-12h),T)", "6"},
		{"date(t_add(table(equipment_sub_part).row(2).field(start_use_time),-12h),T)", "6"},
	}
	script := new(Script)
	script.BaseTime = "2020-01-01 00:00:00"
	script.BeginTime = "2020-06-04 08:00:00"
	script.EndTime = "2020-06-04 16:00:00"
	script.ShiftHour = 12

	for _, tt := range tests {
		script.ScriptStr = tt.scr
		rst, err := script.extractDateTime()
		if err != nil {
			t.Error(err.Error())
		} else {
			t.Logf("%s=%s", script.ScriptStr, rst)
		}
	}
}
*/
