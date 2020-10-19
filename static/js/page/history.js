//=========全局变量定义区域========================================================
var TAG;//当前选中的变量
var TAG_HAVE_SELECTED=false;//已经选择了变量
var HIS_TIME=[];//历史数据时间数组
var HIS_DATA=[];//历史数据数据数组
var HIS_TABLE;//历史数据表
var HIS_INTERVAL_TIME=[];//等间隔历史数据时间数组
var HIS_INTERVAL_DATA=[];//等间隔历史数据数据数组
var HIS_INTERVAL_TABLE;//等间隔历史数据表
var HIS_SUMMARY;//历史统计数据
var HIS_SUMMARY_TABLE;//历史统计数据数据表
var HIS_SUM_GROUP_KEY=[];//历史统计数据分组KEY
var HIS_SUM_GROUP_VAL=[];//历史统计数据分组数值
var HIS_INCREMENT_DATA=[];//历史数据增量数据数组,其时间维度与原始历史数据时间HIS_TIME相同
var SHOW_HIS_TABLE=0;//显示历史数据表
var SHOW_HIS_INTERVAL_TABLE=0;//显示等间隔历史数据表
var START_TIME;//加载数据开始时间
var CHART_Q1=0.0;//echart上显示的数值
var CHART_Q3=0.0;//echart上显示的数值
var CHART_Mean=0.0;//echart上显示的数值
var CHART_Median=0.0;//echart上显示的数值
//=========动作响应区域===========================================================
function zTreeOnClick(event, treeId, treeNode) {
	SELECT_LEVEL_CODE=treeNode.treelevel;
	SELECT_NAME=treeNode.name;
	SELECT_IS_TAG=treeNode.istag;
	if(treeNode.istag==1){
		var datatype;
		switch(treeNode.seq){
		case 1:
			datatype="BOOL";
			break;
		case 2:
			datatype="INT";
			break;
		case 3:
			datatype="FLOAT";
			break;
		default:
			datatype="";
			break;
		}
		var tag={
			Name:treeNode.name,//变量描述名称
			FullName:treeNode.treelevel,//变量层级码
			Id:treeNode.itemid,//变量id
			DotNum:treeNode.dotnum,//小数点数量
			TagType:datatype,//数据类型
			Unit:treeNode.unit//单位
		};
		TAG=tag;
		var unit=TAG.Unit;
		if(unit.length ==0){//没有定义单位的时候
			if(treeNode.seq == 1){//如果是BOOL变量
				unit="无";
			}else{
				unit="未设定";
			}
		}
		$("#TagName").text(TAG.Name);
		$("#TagUnit").text(unit);
		TAG_HAVE_SELECTED = true;
		requestSnapshot(TAG.FullName);
		START_TIME = new Date;

		clearView();
	}
};
//响应鼠标双击
function zTreeOnDbClick(event, treeId, treeNode){
	
}
//时间输入框的值发生改变
function onTimeChange(){
	timeDiffCheck();
	if(TAG_HAVE_SELECTED==true){
		START_TIME = new Date;
		requestHistory(TAG.FullName);//读取历史数据统计请求
	}
}
//时间输入框获得输入焦点时设置最大值和最小值
function onTimeFocus(id){
	var now = new Date;
	if(id=="EndTime"){
		var begintime=$("#BeginTime").val();
		var bgstemp = new Date(begintime.replace(/T/," "));
		bgstemp.setTime(bgstemp.getTime() + $("#Interval").val()*1000);
		$("#"+id).attr("max",DateFormat("YYYY-mm-ddTHH:MM",now));
		$("#"+id).attr("min",DateFormat("YYYY-mm-ddTHH:MM",bgstemp));
	}else{
		var endtime=$("#EndTime").val();
		var edstemp = new Date(endtime.replace(/T/," "));
		edstemp.setTime(edstemp.getTime() - $("#Interval").val()*1000);
		$("#"+id).attr("max",DateFormat("YYYY-mm-ddTHH:MM",edstemp));
	}
};
//响应点击上一时间段
function onLast(){
	var endtime=$("#EndTime").val();
	var begintime=$("#BeginTime").val();
	var bgstemp = new Date(begintime.replace(/T/," "));
	var edstemp = new Date(endtime.replace(/T/," "));
	var timediff = (edstemp.getTime() - bgstemp.getTime());
	bgstemp.setTime(bgstemp.getTime() - timediff);
	edstemp.setTime(edstemp.getTime() - timediff);
	$("#BeginTime").val(DateFormat("YYYY-mm-ddTHH:MM",bgstemp));
	$("#EndTime").val(DateFormat("YYYY-mm-ddTHH:MM",edstemp));
	if(TAG_HAVE_SELECTED==true){
		START_TIME = new Date;
		requestHistory(TAG.FullName);//读取历史数据统计请求
	}
}
//响应点击下一时间段
function onNext(){
	var now = new Date;
	var endtime=$("#EndTime").val();
	var begintime=$("#BeginTime").val();
	var bgstemp = new Date(begintime.replace(/T/," "));
	var edstemp = new Date(endtime.replace(/T/," "));
	var timediff = (edstemp.getTime() - bgstemp.getTime());
	edstemp.setTime(edstemp.getTime() + timediff);
	bgstemp.setTime(bgstemp.getTime() + timediff);
	if(edstemp.getTime() > (now.getTime()-60*1000)){
		edstemp.setTime(now.getTime()-60*1000);
		bgstemp.setTime(edstemp.getTime() - timediff);
	}
	$("#BeginTime").val(DateFormat("YYYY-mm-ddTHH:MM",bgstemp));
	$("#EndTime").val(DateFormat("YYYY-mm-ddTHH:MM",edstemp));
	if(TAG_HAVE_SELECTED==true){
		START_TIME = new Date;
		requestHistory(TAG.FullName);//读取历史数据统计请求
	}
}
//时间选择按钮被按下
function onTimediffClick(tdiff){
	var endtime=$("#EndTime").val();
	var begintime=$("#BeginTime").val();
	var bgstemp = new Date(begintime.replace(/T/," "));
	var edstemp = new Date(endtime.replace(/T/," "));
	bgstemp.setTime(edstemp.getTime() - tdiff*60*1000);
	$("#BeginTime").val(DateFormat("YYYY-mm-ddTHH:MM",bgstemp));
	switch(tdiff){
	case 60://1小时
		$("#Interval").val(5);
		break;
	case 480://8小时
		$("#Interval").val(30);
		break;
	case 720://12小时
	    $("#Interval").val(60);
		break;
	case 1440://24小时
		$("#Interval").val(120);
		break;
	default:
		break;
	}
	if(TAG_HAVE_SELECTED==true){
		START_TIME = new Date;
		requestHistory(TAG.FullName);//读取历史数据统计请求
	}
}
//响应显示对齐数据选择框
function onShowHisIntervalData(id){
	SHOW_HIS_INTERVAL_TABLE = 1 - SHOW_HIS_INTERVAL_TABLE;//切换状态
	if(SHOW_HIS_INTERVAL_TABLE==1){
		$("#HisIntervalTable").html(HIS_INTERVAL_TABLE);//显示
	}else{
		$("#HisIntervalTable").html("");//不显示
	}
}
//响应显示历史数据选择框
function onShowHisData(id){
	SHOW_HIS_TABLE = 1 - SHOW_HIS_TABLE;//切换状态
	if(SHOW_HIS_TABLE==1){
		$("#HisDataTable").html(HIS_TABLE);//显示
	}else{
		$("#HisDataTable").html("");//不显示
	}
}
//页面初始化工作
function pageInit(){
	timeDiffCheck();
	$('[data-toggle="tooltip"]').tooltip();   
};

//根据选择的时间设置时间区间选择框
function timeDiffCheck(){
	var endtime=$("#EndTime").val()
	var begintime=$("#BeginTime").val()
	var bgstemp = new Date(begintime.replace(/T/," "));
	var edstemp = new Date(endtime.replace(/T/," "));
	var timediff = (edstemp.getTime() - bgstemp.getTime());
	
	switch(timediff/1000){
	case 60*60:
		$("#rd_1").attr("checked","checked");
		break;
	case 60*60*8:
		$("#rd_2").attr("checked","checked");
		break;
	case 60*60*12:
		$("#rd_3").attr("checked","checked");
		break;
	case 60*60*24:
		$("#rd_4").attr("checked","checked");
		break;
	default:
		$("#rd_1").removeAttr("checked");
		$("#rd_2").removeAttr("checked");
		$("#rd_3").removeAttr("checked");
		$("#rd_4").removeAttr("checked");
		break;
	}
}
function clearView(){
	HIS_SUMMARY_TABLE='<table class="table table-striped table-hover table-sm"><thead class="thead-light"><tr><th>最小</th><th>最大</th><th title="最大值-最小值">极差</th><th title="算术平均值">算术平均</th><th title="加权平均值">加权平均</th><th>众数</th><th title="下四分位数">Q1</th><th title="中位数">Q2</th><th title="上四分位数">Q3</th><th title="四分位差">Qd</th><th>差</th><th>正差</th><th>面积</th><th>点数</th><th title="标准差">SD</th><th>偏度</th><th>峰度</th></tr></thead><tbody>';
	$("#HisSumData").html(HIS_SUMMARY_TABLE);
	$("#HisDataTable").html("");//不显示
	$("#HisIntervalTable").html("");//不显示
	$("#Echarts_His").hide();
	$("#Echarts_HisInterval").hide();
	$("#Echarts_HisGroup").hide();

}
//=========数据接收解析区域========================================================
//接收AJAX反馈的快照数据并解析
function getTagSnapshotData(ajaxdata){
	var dtarr = eval("("+ajaxdata+")"); 
	var d = new Date();
	var t;
	var snap;
	snap=dtarr[TAG.FullName];
	if(snap.Id>0){
		$("#TagValue").text(DataToFixed(snap.Rtsd.Value,TAG.TagType,TAG.DotNum));//更新TagValue
		d.setTime(snap.Rtsd.Time);//将2006-05-06T00:00:00Z格式的时间转换为UTC时间戳
		$("#TagTime").text(d.toLocaleString());//更新时间戳:转换为当地时间格式
	}else{
		$("#TagValue").html('<span class="badge badge-danger">#Error</span>');//更新TagValue
		$("#TagTime").text('');//更新时间戳:转换为当地时间格式
		alert('没有在数据库中找到匹配变量名['+TAG.FullName+']的变量,请检查!');
	}
}
//接收AJAX反馈的历史统计数据并解析
function getTagHistorySummary(ajaxdata){
	var ajax = eval("("+ajaxdata+")"); 
	dtarr = ajax[TAG.FullName];
	var suma={
		Min:DataToFixed(dtarr.Min,TAG.TagType,TAG.DotNum),            //最小值(基本)
		Max:DataToFixed(dtarr.Max,TAG.TagType,TAG.DotNum),            //最大值(基本)
		Range:DataToFixed(dtarr.Range,TAG.TagType,TAG.DotNum),        //数据范围(Max-Min)(基本)
		Total:DataToFixed(dtarr.Total,TAG.TagType,TAG.DotNum),          //表示统计时间段内的累计值，结果的单位为标签点的工程单位(面积,值*时间(s))(基本)
		Sum:DataToFixed(dtarr.Sum,TAG.TagType,TAG.DotNum),            //统计时间段内的算术累积值(值相加)(基本)
		Mean:DataToFixed(dtarr.Mean,TAG.TagType,TAG.DotNum),          //统计时间段内的算术平均值(Mean = Sum/PointCnt)(基本)
		PowerAvg:DataToFixed(dtarr.PowerAvg,TAG.TagType,TAG.DotNum),       //统计时间段内的加权平均值,对BOOL量而言是ON率（Total/Duration）(基本)
		Diff:DataToFixed(dtarr.Diff,TAG.TagType,TAG.DotNum),          //统计时间段内的差值(最后一个值减去第一个值)(基本)
		PlusDiff:DataToFixed(dtarr.PlusDiff,TAG.TagType,TAG.DotNum),       //正差值,用于累计值求差,可以削除清零对值的影响(统计周期内只可以有一次清零动作)
		Duration:dtarr.Duration,     //统计时间段内的秒数(EndTime - BeginTime)(基本)
		PointCnt:dtarr.PointCnt,     //统计时间段内的数据点数(基本)
		RisingCnt:dtarr.RisingCnt,   //统计时间段内数据上升的次数(基本)
		FallingCnt:dtarr.FallingCnt, //统计时间段内数据下降的次数(基本)
		LtzCnt:dtarr.LtzCnt,         //小于0的次数
		GtzCnt:dtarr.GtzCnt,         //大于0的次数
		EzCnt:dtarr.EzCnt,           //等于0的次数
		BeginTime:dtarr.BeginTime,   //开始时间(基本)
		EndTime:dtarr.EndTime,       //结束时间(基本)
		SD:DataToFixed(dtarr.SD,TAG.TagType,TAG.DotNum),             //总体标准差(高级)
		STDDEV:DataToFixed(dtarr.STDDEV,TAG.TagType,TAG.DotNum),     //样本标准差(高级)
		SE:DataToFixed(dtarr.SE,TAG.TagType,TAG.DotNum),             //标准误差(SE = STDDEV / PointCnt)(高级)
		Ske:DataToFixed(dtarr.Ske,TAG.TagType,TAG.DotNum),           //偏度(高级)
		Kur:DataToFixed(dtarr.Kur,TAG.TagType,TAG.DotNum),           //峰度(高级)
		Mode:DataToFixed(dtarr.Mode,TAG.TagType,TAG.DotNum),         //众数(高级)
		Median:DataToFixed(dtarr.Median,TAG.TagType,TAG.DotNum),     //中位数(高级)
		Q1:DataToFixed(dtarr.Q1,TAG.TagType,TAG.DotNum),     //下四分位(高级)
		Q3:DataToFixed(dtarr.Q3,TAG.TagType,TAG.DotNum),     //上四分位(高级)
		Qd:DataToFixed(dtarr.Qd,TAG.TagType,TAG.DotNum),     //四分位差(高级)
		Lower:DataToFixed(dtarr.Lower,TAG.TagType,TAG.DotNum),     //下限
		Upper:DataToFixed(dtarr.Upper,TAG.TagType,TAG.DotNum),     //上限
		OutliersCnt:dtarr.OutliersCnt,     //离群点数
		GroupDist:dtarr.GroupDist       //组距GroupDistance(高级),DataGroup中两组数之间的距离
	};
	HIS_SUMMARY = suma;
	HIS_SUMMARY_TABLE='<table class="table table-striped table-hover table-sm"><thead class="thead-light"><tr><th>最小</th><th>最大</th><th title="最大值-最小值">极差</th><th title="算术平均值">算术平均</th><th title="加权平均值">加权平均</th><th>众数</th><th title="下四分位数">Q1</th><th title="中位数">Q2</th><th title="上四分位数">Q3</th><th title="四分位差">Qd</th><th title="Q1-1.5*Qd">下限</th><th title="Q3+1.5*Qd">上限</th><th>差</th><th>正差</th><th>面积</th><th>点数</th><th>离群点数</th><th title="标准差">SD</th></tr></thead><tbody>';
	HIS_SUMMARY_TABLE +='<tr><td>'+suma.Min+'</td><td>'+suma.Max+'</td><td>'+suma.Range+'</td><td>'+suma.Mean+'</td><td>'+suma.PowerAvg+'</td><td>'+suma.Mode+'</td><td>'+suma.Q1+'</td><td>'+suma.Median+'</td><td>'+suma.Q3+'</td><td>'+suma.Qd+'</td><td>'+suma.Lower+'</td><td>'+suma.Upper+'</td><td>'+suma.Diff+'</td><td>'+suma.PlusDiff+'</td><td>'+suma.Total+'</td><td>'+suma.PointCnt+'</td><td>'+suma.OutliersCnt+'</td><td>'+suma.SD+'</td></tr>';
	HIS_SUMMARY_TABLE +='</tbody></table>';
	$("#HisSumData").html(HIS_SUMMARY_TABLE);
	
	var dgroup = dtarr.DataGroup;
	var increment = dtarr.Increment;
	var point=suma.PointCnt;
	if(point == 0){
		point = 1;
	}
	HIS_SUM_GROUP_KEY.splice(0,HIS_SUM_GROUP_KEY.length);//清空数组
	HIS_SUM_GROUP_VAL.splice(0,HIS_SUM_GROUP_VAL.length);//清空数组
	HIS_INCREMENT_DATA.splice(0,HIS_INCREMENT_DATA.length);//清空数组
	
	var groupkey=[];
	var groupval=[];
	var keyarray=[];
	var k=0;
	for(var key in dgroup){
		var keyd=DataToFixed(parseFloat(key),TAG.TagType,TAG.DotNum);
		keyarray[k]=keyd;
		groupkey[keyd] = keyd;
		groupval[keyd] = DataToFixed(dgroup[key],'float',2);
		k++;
	}
	keyarray=BubbleAsc(keyarray);
	for(var k in keyarray){
		HIS_SUM_GROUP_KEY[k] = keyarray[k];
		HIS_SUM_GROUP_VAL[k] = groupval[keyarray[k]];
	}
	var v1,v2;
	for(var k=0;k<HIS_SUM_GROUP_KEY.length-1;k++){
		v1=HIS_SUM_GROUP_KEY[k];
		v2=HIS_SUM_GROUP_KEY[k+1];
		//console.log(v1,v2);
		if (parseFloat(suma.Q1)>=parseFloat(v1) && parseFloat(suma.Q1)<parseFloat(v2)){
			CHART_Q1 = v1;
		}
		if (parseFloat(suma.Q3)>=parseFloat(v1) && parseFloat(suma.Q3)<parseFloat(v2)){
			CHART_Q3 = v2;
		}
		if (parseFloat(suma.Mean)>=parseFloat(v1) && parseFloat(suma.Mean)<parseFloat(v2)){
			CHART_Mean = v1;
		}
		if (parseFloat(suma.Median)>=parseFloat(v1) && parseFloat(suma.Median)<parseFloat(v2)){
			CHART_Median = v1;
		}
	}
	//console.log(suma.Q1,suma.Mean,suma.Median,suma.Q3);
	for(var i=0;i<HIS_TIME.length;i++){
		HIS_INCREMENT_DATA[i] = DataToFixed(increment[HIS_TIME[i]],TAG.TagType,TAG.DotNum);
	}
	refreshEcharts_his();//刷新Echarts 
	refreshEcharts_hisGroup();//刷新Echarts 
}
//接收AJAX反馈的等间隔历史数据并解析
function getTagHistoryInterval(ajaxdata){//根据Ajax反馈的结果更新Tag的实时数据
	var d = new Date();
	var dtarr = eval("("+ajaxdata+")"); 
	HIS_INTERVAL_TIME.splice(0,HIS_INTERVAL_TIME.length);//清空数组
	HIS_INTERVAL_DATA.splice(0,HIS_INTERVAL_DATA.length);//清空数组
	HIS_INTERVAL_TABLE='<hr/><h3>等间隔聚合历史数据</h1><br/><table class="table table-striped table-hover table-sm"><thead class="thead-light"><tr><th>序号</th><th>时间</th><th>数据</th></tr></thead><tbody>';
	if(dtarr != null){
		var hisdata=dtarr[TAG.FullName];
		for(var i=0;i< hisdata.length;i++){
			HIS_INTERVAL_TIME[i] = hisdata[i].Time;//更新时间戳:转换为当地时间格式
			HIS_INTERVAL_DATA[i] = DataToFixed(hisdata[i].Value,TAG.TagType,TAG.DotNum);//更新TagValue
			HIS_INTERVAL_TABLE += '<tr><td>'+(i+1)+'</td><td>'+HIS_INTERVAL_TIME[i]+'</td><td>'+HIS_INTERVAL_DATA[i]+'</td>  </tr>';
		}
	}
	HIS_INTERVAL_TABLE +='</tbody></table>';
	if(SHOW_HIS_INTERVAL_TABLE==1){
		$("#HisIntervalTable").html(HIS_INTERVAL_TABLE);//显示
	}
	refreshEcharts_hisInterval();//刷新Echarts 

	var timediff = new Date() - START_TIME;//计算完成耗时,ms
	$("#LoadDataMsg").html("数据加载完成,耗时:"+timediff/1000+"秒");
}
//接收AJAX反馈的历史数据并解析
function getTagHistory(ajaxdata){//根据Ajax反馈的结果更新Tag的实时数据
	var d = new Date();
	var dtarr = eval("("+ajaxdata+")"); 
	HIS_TIME.splice(0,HIS_TIME.length);//清空数组
	HIS_DATA.splice(0,HIS_DATA.length);//清空数组
	HIS_TABLE='<hr/><h3>原始历史数据</h1><br/><table class="table table-striped table-hover table-sm"><thead class="thead-light"><tr><th>序号</th><th>时间</th><th>数据</th></tr></thead><tbody>';
	if(dtarr != null){
		var hisdata=dtarr[TAG.FullName];
		for(var i=0;i< hisdata.length;i++){
			HIS_TIME[i] = hisdata[i].Time;//更新时间戳:转换为当地时间格式
			HIS_DATA[i] = DataToFixed(hisdata[i].Value,TAG.TagType,TAG.DotNum);//更新TagValue
			HIS_TABLE += '<tr><td>'+(i+1)+'</td><td>'+HIS_TIME[i]+'</td><td>'+HIS_DATA[i]+'</td>  </tr>';
		}
	}
	HIS_TABLE +='</tbody></table>';
	if(SHOW_HIS_TABLE==1){
		$("#HisDataTable").html(HIS_TABLE);//显示
	}
}
//=========读取请求区域============================================================
//读取变量快照请求
function requestSnapshot(tagname){
	var urlstr = "api/snapshot?tagnames="+tagname;
	loadTagSnapshot(urlstr);
}
////读取历史数据统计请求
function requestHistorySummary(tagname){
	var urlstr = "api/historysummary?tagname="+tagname+"&begintime="+InputDateTimeToString("YYYY-mm-dd HH:MM:SS",$("#BeginTime").val())+"&endtime="+InputDateTimeToString("YYYY-mm-dd HH:MM:SS",$("#EndTime").val());
	loadTagHistorySummary(urlstr);
}
//读取等间隔历史数据请求
function requestHistoryInterval(tagname){
	var urlstr = "api/hisinterval?tagname="+tagname+"&begintime="+InputDateTimeToString("YYYY-mm-dd HH:MM:SS",$("#BeginTime").val())+"&endtime="+InputDateTimeToString("YYYY-mm-dd HH:MM:SS",$("#EndTime").val())+"&interval="+$("#Interval").val();
	loadTagHistoryInterval(urlstr);
}
//读取原始历史数据请求
function requestHistory(tagname){
	var urlstr = "api/history?tagname="+tagname+"&begintime="+InputDateTimeToString("YYYY-mm-dd HH:MM:SS",$("#BeginTime").val())+"&endtime="+InputDateTimeToString("YYYY-mm-dd HH:MM:SS",$("#EndTime").val());
	loadTagHistory(urlstr);
}

//=========AJAX函数定义区域=======================================================
function loadTagSnapshot(urlstr)//从数据库中读取单一变量的最新值
{
	$("#LoadDataMsg").html("正在加载快照数据……");	
	//调用公用的loadXMLDoc函数
    loadXMLDoc(urlstr,function()
    {
        if (xmlhttp.readyState==4 && xmlhttp.status==200)//请求处理完成，且状态OK
        {//添加请求完成后的处理功能---------------------------------------
			getTagSnapshotData(xmlhttp.responseText);//解读数据
			//下一步：加载历史数据
			requestHistory(TAG.FullName);
			
        }//请求完成后的处理功能结束---------------------------------------
    });
}
function loadTagHistorySummary(urlstr)//从数据库读取统计值
{
	$("#LoadDataMsg").html("正在加载统计数据……");	
	//调用公用的loadXMLDoc函数
    loadXMLDoc(urlstr,function()
    {
        if (xmlhttp.readyState==4 && xmlhttp.status==200)//请求处理完成，且状态OK
        {//添加请求完成后的处理功能---------------------------------------
			getTagHistorySummary(xmlhttp.responseText);
			//下一步：加载等间隔历史数据
			requestHistoryInterval(TAG.FullName);
        }//请求完成后的处理功能结束---------------------------------------
    });
}
function loadTagHistory(urlstr)//从数据库中读取单一变量指定时间段的原始历史数据
{
	$("#LoadDataMsg").html("正在加载历史数据……");	
	//调用公用的loadXMLDoc函数
    loadXMLDoc(urlstr,function()
    {
        if (xmlhttp.readyState==4 && xmlhttp.status==200)//请求处理完成，且状态OK
        {//添加请求完成后的处理功能---------------------------------------
			getTagHistory(xmlhttp.responseText);
			//下一步：加载历史统计数据
			requestHistorySummary(TAG.FullName);			
        }//请求完成后的处理功能结束---------------------------------------
    });
}

function loadTagHistoryInterval(urlstr)//从数据库中读取单一变量指定时间段的等间隔历史数据
{
	$("#LoadDataMsg").html("正在加载等间隔历史数据……");	
	//调用公用的loadXMLDoc函数
    loadXMLDoc(urlstr,function()
    {
        if (xmlhttp.readyState==4 && xmlhttp.status==200)//请求处理完成，且状态OK
        {//添加请求完成后的处理功能---------------------------------------
			getTagHistoryInterval(xmlhttp.responseText);
        }//请求完成后的处理功能结束---------------------------------------
    });
}

function loadUpdateTree(urlstr)//更新结构树
{
	//调用公用的loadXMLDoc函数
    loadXMLDoc(urlstr,function()
    {
        if (xmlhttp.readyState==4 && xmlhttp.status==200)//请求处理完成，且状态OK
		{//添加请求完成后的处理功能---------------------------------------
			ShowModal("更新提示",'<div class="alert alert-success">'+xmlhttp.responseText+'</div>');
        }//请求完成后的处理功能结束---------------------------------------
    });
}

$(document).ready(function() {
	$("#exit").after('退出');
	$("#UpdateTree").html("更新")
	$('#ExpandTreeNode').attr("title","未选中节点时展开所有节点,选中节点时展开选中节点")
	$("#ExpandTreeNode").html('展开');
	$('#CollapseTreeNode').attr("title","未选中节点时折叠所有节点,选中节点时折叠选中节点")
	$("#CollapseTreeNode").html('折叠');
	$("#HideTreeNode").html('隐藏：');
	$('#SearchTreeNode').attr("placeholder","搜索")

	$("#BeginTimes").text('起始时间：');
	$("#EndTimes").before('结束时间：');
	$("#Interval").before("聚合时间：")
	$("#ThelastTimestemp").text('时间范围：');
	$("#rd_1").before("&nbsp; 1小时")
	$("#rd_2").before("&nbsp; 8小时")
	$("#rd_3").before("&nbsp; 12小时")
	$("#rd_4").before("&nbsp; 24小时")

	$("#show_hisinterval_data").before("显示聚合数据表")
	$("#show_his_data").before("显示原始数据表")

	$("#TagName").before("无")
	$("#TagValue").before("无")
	$("#TagUnit").before("无")
	$("#TagTime").before("无")

	$("#Last").attr("title",'前一时间段!');
	$("#Next").attr("title",'后一时间段!');

	$("#VariableName").html("变量名:")
	$("#LatestValue").html('最新值:');
	$("#unit").html('单位:');
	$("#UpdateTime").text('更新时间:');
	$("#HisDataTable").html('<strong>操作：</strong> 请在左侧结构树中选择一个节点以便显示其下的变量快照数据！');
});

function onUpdateTree(){
	var url="api/updatetagnodetree?withtag=1";
	loadUpdateTree(url);
}
