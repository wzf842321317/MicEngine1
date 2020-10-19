//=========变量定义区域===========================================================
var SELECT_NODE={//当前选中的样本节点
	Id:0,//ID
	Pid:0,//父层级ID
	Name:"",//名称
	Desc:"",//描述
	TagType:"",//节点类型
	IsForder:0,//是否文件夹
	BaseTime:"",//基准时间
	ShiftHour:8//每班时间长度
};
var HAS_SELECT_NODE=false;//已经选择了有效节点
var NODES_CFG;//节点配置信息
var DATAS=[];//物耗数据
var DATA_TIME=[];//数据时间数组
var DATA_VALUE=[];//数据值数组
var SELECT_LEAF;//选择了叶子节点
var PIE_LEGEND=[];//饼图图标
var PIE_DATA=[];//饼图数据

//==================动作响应区域==================================================
//响应鼠标单击
function zTreeOnClick(event, treeId, treeNode) {
	$("#OperationMsg").hide();//隐藏操作提示
	$("#CfgMsg").html("");//清空数据
	$("#DataMsg").html("");//清空数据

	HAS_SELECT_NODE = true;//已经选择了有效节点

	SELECT_NODE.Id=treeNode.id;
	SELECT_NODE.Pid=treeNode.pId;
	SELECT_NODE.Name=treeNode.name;
	SELECT_NODE.Desc=treeNode.desc;
	SELECT_NODE.TagType=treeNode.tagtype;
	SELECT_NODE.IsForder=treeNode.isforder;
	SELECT_NODE.ShiftHour=treeNode.shifthour;
	SELECT_NODE.BaseTime=treeNode.basetime;
	
	if (treeNode.isforder == 0){//选择了叶子节点
		SELECT_LEAF=true;
		requestCfg(SELECT_NODE.Id);//请求读取物耗配置信息
		$("#OperationMsg").hide();//隐藏操作提示
	}else{//选择了文件夹
		SELECT_LEAF=false;
		if(treeNode.id==0){//选择的是顶级文件夹节点
			requestTagTypeList();
		}else{//选择的不是顶级文件夹节点
			requestKpiTagTypePeriodInfo(treeNode.tagtype);
		}
	}
	timeRangeSelectorSet();//时间范围设置
};

//响应鼠标双击
function zTreeOnDbClick(event, treeId, treeNode){
	
};
//页面初始化工作
function pageInit(){
	timeRangeCheck();
	timeRangeSelectorSet();
	requestTagTypeList();
};
//请求数据
function requestDatas(){
	if(HAS_SELECT_NODE == true){//已经选择了有效节点
		if(SELECT_LEAF == true){//
			requestGoodsDatas(SELECT_NODE.Id/100000,1);
		}else{
			requestGoodsDatas(SELECT_NODE_LEVEL,0);
		}
	}
}

//时间范围选择框设置
function timeRangeSelectorSet(){
	if(SELECT_NODE.BaseTime.length > 10){//已经选择了样本模板,取消禁用选择本班、今日、本月、上月
		if (SELECT_NODE.ShiftHour > 0){
			$("#op_this_shift").removeAttr("disabled");
		}else{
			$("#op_this_shift").attr("disabled","disabled");
		}
		$("#op_this_day").removeAttr("disabled");
		$("#op_this_month").removeAttr("disabled");
	}else{//没有选择样本模板,禁用选择本班、今日、本月、上月
		$("#op_this_shift").attr("disabled","disabled");
		$("#op_this_day").attr("disabled","disabled");
		$("#op_this_month").attr("disabled","disabled");
	}
}
//时间输入框的值发生改变
function onTimeChange(){
	timeRangeCheck();
	if(SELECT_LEAF==true){
		requestCfg(SELECT_NODE.Id);//请求读取物耗配置信息
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
	var range = $("#TimeRange").val();
	switch(range){
	case '1':
		timediff = SELECT_NODE.ShiftHour * 3600 * 1000;
		bgstemp.setTime(bgstemp.getTime() - timediff);
		edstemp.setTime(bgstemp.getTime() + timediff);
		$("#BeginTime").val(DateFormat("YYYY-mm-ddTHH:MM",bgstemp));
		$("#EndTime").val(DateFormat("YYYY-mm-ddTHH:MM",edstemp));
		break;
	case '2':
		timediff = 24 * 3600 * 1000;
		bgstemp.setTime(bgstemp.getTime() - timediff);
		edstemp.setTime(bgstemp.getTime() + timediff);
		$("#BeginTime").val(DateFormat("YYYY-mm-ddTHH:MM",bgstemp));
		$("#EndTime").val(DateFormat("YYYY-mm-ddTHH:MM",edstemp));
		break;
	case '3':
		getBeginTimeOfMonth(SELECT_NODE.BaseTime,begintime);
		break;
	default:
		bgstemp.setTime(bgstemp.getTime() - timediff);
		edstemp.setTime(edstemp.getTime() - timediff);
		$("#BeginTime").val(DateFormat("YYYY-mm-ddTHH:MM",bgstemp));
		$("#EndTime").val(DateFormat("YYYY-mm-ddTHH:MM",edstemp));
		break;
	}
	
	if(SELECT_LEAF==true){
		requestCfg(SELECT_NODE.Id);//请求读取物耗配置信息
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
	var range = $("#TimeRange").val();
	switch(range){
	case '1':
		timediff = SELECT_NODE.ShiftHour * 3600 * 1000;
		edstemp.setTime(edstemp.getTime() + timediff);
		bgstemp.setTime(bgstemp.getTime() + timediff);
		if(edstemp.getTime() > (now.getTime()-60*1000)){
			edstemp.setTime(now.getTime()-60*1000);
		}
		if(bgstemp.getTime() > edstemp.getTime()){
			bgstemp.setTime(bgstemp.getTime() - timediff);
		}
		$("#BeginTime").val(DateFormat("YYYY-mm-ddTHH:MM",bgstemp));
		$("#EndTime").val(DateFormat("YYYY-mm-ddTHH:MM",edstemp));
		break;
	case '2':
		timediff = 24 * 3600 * 1000;
		edstemp.setTime(edstemp.getTime() + timediff);
		bgstemp.setTime(bgstemp.getTime() + timediff);
		if(edstemp.getTime() > (now.getTime()-60*1000)){
			edstemp.setTime(now.getTime()-60*1000);
		}
		if(bgstemp.getTime() > edstemp.getTime()){
			bgstemp.setTime(bgstemp.getTime() - timediff);
		}
		$("#BeginTime").val(DateFormat("YYYY-mm-ddTHH:MM",bgstemp));
		$("#EndTime").val(DateFormat("YYYY-mm-ddTHH:MM",edstemp));
		break;
	case '3':
		if(edstemp.getMonth()==11){//12月份
			edstemp.setFullYear(edstemp.getFullYear()+1,0);//年加1,月归零
		}else{
			edstemp.setMonth(edstemp.getMonth()+1);//月加1
		}
		if(bgstemp.getMonth()==11){//12月份
			bgstemp.setFullYear(bgstemp.getFullYear()+1,0);//年加1,月归零
		}else{
			bgstemp.setMonth(bgstemp.getMonth()+1);//月加1
		}
		if(edstemp.getTime() > (now.getTime()-60*1000)){//结束时间大于当前时间
			edstemp.setTime(now.getTime()-60*1000);//等于当前时间
		}
		if(bgstemp.getTime() >= edstemp.getTime()){//开始时间大于结束时间
			if(bgstemp.getMonth()==0){//12月份
				bgstemp.setFullYear(bgstemp.getFullYear()-1,11);//年减1,月为满
			}else{
				bgstemp.setMonth(bgstemp.getMonth()-1);//月减1
			}
		}
		$("#BeginTime").val(DateFormat("YYYY-mm-ddTHH:MM",bgstemp));
		$("#EndTime").val(DateFormat("YYYY-mm-ddTHH:MM",edstemp));
		break;
	default:
		edstemp.setTime(edstemp.getTime() + timediff);
		bgstemp.setTime(bgstemp.getTime() + timediff);
		if(edstemp.getTime() > (now.getTime()-60*1000)){
			edstemp.setTime(now.getTime()-60*1000);
			bgstemp.setTime(edstemp.getTime() - timediff);
		}
		$("#BeginTime").val(DateFormat("YYYY-mm-ddTHH:MM",bgstemp));
		$("#EndTime").val(DateFormat("YYYY-mm-ddTHH:MM",edstemp));
		break;
	}
	if(SELECT_LEAF==true){
		requestCfg(SELECT_NODE.Id);//请求读取物耗配置信息
	}
	//requestDatas();//请求数据
}
//时间范围发生改变
function onTimeRangeChange(tdiff){
	var endtime=$("#EndTime").val();
	var begintime=$("#BeginTime").val();
	var bgstemp = new Date(begintime.replace(/T/," "));
	var edstemp = new Date(endtime.replace(/T/," "));
	var nowTime = new Date;
	switch(tdiff){
	case '0'://主动选择自定义时无效
		break;
	case '1'://本班
		getBeginTimeInDay(SELECT_NODE.BaseTime,DateFormat("YYYY-mm-ddTHH:MM",nowTime),SELECT_NODE.ShiftHour);
		break;
	case '2'://今日
		getBeginTimeInDay(SELECT_NODE.BaseTime,DateFormat("YYYY-mm-ddTHH:MM",nowTime),24);
		break;
	case '3'://本月
		getBeginTimeOfMonth(SELECT_NODE.BaseTime,DateFormat("YYYY-mm-ddTHH:MM",nowTime))
		break;
	default:
		bgstemp.setTime(edstemp.getTime() - tdiff*60*1000);
		$("#BeginTime").val(DateFormat("YYYY-mm-ddTHH:MM",bgstemp));
		break;
	}	

	if(SELECT_LEAF==true){
		requestCfg(SELECT_NODE.Id);//请求读取物耗配置信息
	}
}
//根据选择的时间设置时间区间选择框
function timeRangeCheck(){
	var endtime=$("#EndTime").val()
	var begintime=$("#BeginTime").val()
	var bgstemp = new Date(begintime.replace(/T/," "));
	var edstemp = new Date(endtime.replace(/T/," "));
	var timediff = (edstemp.getTime() - bgstemp.getTime());
	
	switch(timediff/1000){
	case 60*60*8://8小时
		$("#TimeRange").val(480);
		break;
	case 60*60*12://12小时
		$("#TimeRange").val(720);
		break;
	case 60*60*24://24小时
		$("#TimeRange").val(1440);
		break;
	case 60*60*24*7://7天
		$("#TimeRange").val(10080);
		break;
	default:
		$("#TimeRange").val(0);
		break;
	}
}

//本班或者今日的开始时间
function getBeginTimeInDay(basetime,lasttime,period_h){
	var baseT= new Date(basetime.replace(/T/," "));//基准时间
	var nowT= new Date(lasttime.replace(/T/," "));//当前时间
	var nowTime= new Date(lasttime.replace(/T/," "));//当前时间
	var period = period_h * 3600*1000;//周期
	var endT=nowT;
	var bgT=nowT;
	endT.setTime(Math.round(((nowT.getTime() + period) - baseT.getTime())/period) * period + baseT.getTime());
	bgT.setTime(endT.getTime() - period);
	if(bgT.getTime() > nowTime.getTime()){
		bgT.setTime(endT.getTime() - period);
	}
	beginTime = DateFormat("YYYY-mm-ddTHH:MM",bgT);
	$("#BeginTime").val(beginTime);
	
	endTime = DateFormat("YYYY-mm-ddTHH:MM",nowTime);
	$("#EndTime").val(endTime);
}
//本月开始时间
function getBeginTimeOfMonth(basetime,lasttime){
	var baseT= new Date(basetime.replace(/T/," "));//基准时间
    var baseM = baseT.getMonth();    // 月
    var baseD = baseT.getDate();     // 日
    var baseH = baseT.getHours();    // 时
    var baseMi = baseT.getMinutes(); // 分
	var baseS = baseT.getSeconds();  // 秒

	var bgT= new Date(lasttime.replace(/T/," "));//当前时间
	var bgY = bgT.getFullYear();     // 年
    var bgM = bgT.getMonth();        // 月
    var bgD = bgT.getDate();         // 日
	var bgH = bgT.getHours();        // 时
    var bgMi = bgT.getMinutes();     // 分
	bgT.setHours(baseH,baseMi,baseS,0);
	if (bgD > baseD && bgH > baseH && bgMi > baseMi){//当前日大于等于基准日
		bgT.setDate(baseD);//设置基准日		
	}else{
		if(bgM==0){
			bgT.setFullYear(bgY - 1,11,baseD);//设置基准日
		}else{
			bgT.setMonth(bgM - 1,baseD);//设置基准日
		}
	}
	beginTime = DateFormat("YYYY-mm-ddTHH:MM",bgT);
	$("#BeginTime").val(beginTime);
	var now = new Date;
	var edstemp= new Date(lasttime.replace(/T/," "));//当前时间
	if(edstemp.getTime() > (now.getTime()-60*1000)){
		edstemp.setTime(now.getTime()-60*1000);
	}
	endTime = DateFormat("YYYY-mm-ddTHH:MM",edstemp);
	$("#EndTime").val(endTime);
}

//=========AJAX请求定义区域=======================================================
//读取物耗配置信息请求
function requestCfg(id){
	var urlstr = "api/kpi/getconfig?kpicfgid="+id;
	loadCfg(urlstr,id);
}
//读取物耗数据信息请求
function requestDatas(id){
	var urlstr = "api/kpi/getresult?kpicfgid="+id+"&readtype=0"+"&begintime="+InputDateTimeToString("YYYY-mm-dd HH:MM:SS",$("#BeginTime").val())+"&endtime="+InputDateTimeToString("YYYY-mm-dd HH:MM:SS",$("#EndTime").val());
	loadDatas(urlstr);
}

function requestTagTypeList(){
	var urlstr = "api/script/sql?micsql=select(lst.tag_type,COUNT(lst.tag_type) AS cnt,dic.name).from(calc_kpi_config_list lst join sys_dictionary dic ON (dic.dictionary_name_code=lst.tag_type)).where(lst.status=1).groupby(lst.tag_type).orderby(dic.id).as(map)";
	loadTagTypeList(urlstr);
}

function requestKpiTagTypePeriodInfo(tagtype){
	var urlstr = "api/script/sql?micsql=select(period,COUNT(period) AS cnt).from(calc_kpi_config_list).where(status=1 and tag_type='"+tagtype+"').groupby(period).as(map)";
	loadKpiTagTypePeriodInfo(urlstr);
}
//=========AJAX加载定义区域=======================================================
function loadCfg(urlstr,id)//读取物耗配置信息
{	
	$("#CfgMsg").html('<div class="alert alert-warning">正在加载配置信息……</div>');
	//调用公用的loadXMLDoc函数
    loadXMLDoc(urlstr,function()
    {
        if (xmlhttp.readyState==4 && xmlhttp.status==200)//请求处理完成，且状态OK
        {//添加请求完成后的处理功能---------------------------------------
			getCfg(xmlhttp.responseText);//解读数据
			//下一步：
			requestDatas(id);			
        }//请求完成后的处理功能结束---------------------------------------
    });
}
function loadTagTypeList(urlstr)//读取物耗配置信息
{	
	$("#CfgMsg").html('<div class="alert alert-warning">正在加载配置信息……</div>');
	//调用公用的loadXMLDoc函数
    loadXMLDoc(urlstr,function()
    {
        if (xmlhttp.readyState==4 && xmlhttp.status==200)//请求处理完成，且状态OK
        {//添加请求完成后的处理功能---------------------------------------
			getTagTypeList(xmlhttp.responseText);//解读数据
			//下一步：		
        }//请求完成后的处理功能结束---------------------------------------
    });
}
function loadKpiTagTypePeriodInfo(urlstr)//读取物耗配置信息
{	
	$("#CfgMsg").html('<div class="alert alert-warning">正在加载配置信息……</div>');
	//调用公用的loadXMLDoc函数
    loadXMLDoc(urlstr,function()
    {
        if (xmlhttp.readyState==4 && xmlhttp.status==200)//请求处理完成，且状态OK
        {//添加请求完成后的处理功能---------------------------------------
			getKpiTagTypePeriodInfo(xmlhttp.responseText);//解读数据
			//下一步：		
        }//请求完成后的处理功能结束---------------------------------------
    });
}
function loadDatas(urlstr)//读取物耗配置信息
{	
	$("#DataMsg").html('<div class="alert alert-warning">正在加载数据信息……</div>');
	//调用公用的loadXMLDoc函数
    loadXMLDoc(urlstr,function()
    {
        if (xmlhttp.readyState==4 && xmlhttp.status==200)//请求处理完成，且状态OK
        {//添加请求完成后的处理功能---------------------------------------
			getDatas(xmlhttp.responseText);//解读数据
			//下一步：		
		}else{//请求完成后的处理功能结束---------------------------------------
			$("#PageEcharts").hide();
			$("#DataMsg").html('<div class="alert alert-danger">'+xmlhttp.responseText+'</div>');
		}
    });
}
//=========AJAX数据接收解析区域====================================================
//解析配置信息
function getCfg(ajaxdata){
	/*
	{
		"Id": 1,
		"DistributedId": 0,
		"TagType": "process",
		"TagId": 4337,
		"CalcKpiIndexDic": {
			"Id": 8,
			"KpiTag": "PowerAvg",
			"KpiNameEn": "PowerAvg",
			"KpiNameCn": "加权平均值",
			"TagType": "process",
			"Script": "fc(PowerAvg)",
			"Description": "统计时间段内的加权平均值,对BOOL量而言是ON率（Total/Duration）(基本)",
			"Seq": 8,
			"Status": 1,
			"CreateUserId": 115,
			"CreateTime": "2019-12-01 17:58:08",
			"UpdateUserId": 0,
			"UpdateTime": "2019-12-26 13:09:26",
			"CalcKpiConfigList": null,
			"Variblesets": null
		},
		"KpiTag": "x2_asl_asl-xc1_SK1_GK1_S1-001_SY1-001_meas-value:1__hour_PowerAvg",
		"KpiName": "破碎一车间_新原矿仓料位计_测量值__加权平均值",
		"Script": "",
		"StartTime": "2006-01-02 15:04:00",
		"Period": -1,
		"OffsetMinutes": 0,
		"LastCalcTime": "2020-02-27 15:30:00",
		"Supplement": 0,
		"Description": "",
		"Seq": 0,
		"Status": 0,
		"CreateUserId": 0,
		"CreateTime": "0000-00-00 00:00:00",
		"UpdateUserId": 0,
		"UpdateTime": "2020-02-27 17:19:57",
		"KpiBaseTime": "",
		"KpiShiftHour": 0,
		"CfgUnit": ""
	}
	*/
	var node = eval("("+ajaxdata+")"); 
	var htmlstr='<div class="col-12 border-bottom"><h5>KPI配置信息</h5></div>';
	NODES_CFG = node;
		
	htmlstr+='<div class="col-sm-3 form-inline"><strong>ID:</strong>'+node.Id+'</div>';
	htmlstr+='<div class="col-sm-3 form-inline"><strong>分布式ID:</strong>'+node.DistributedId+'</div>';
	htmlstr+='<div class="col-sm-3 form-inline"><strong>主标签类型:</strong>'+node.TagType+'</div>';
	htmlstr+='<div class="col-sm-3 form-inline"><strong>主标签ID:</strong>'+node.TagId+'</div>';
	htmlstr+='<div class="col-sm-12 form-inline"><strong>KPI标签:</strong>'+node.KpiTag+'</div>';
	htmlstr+='<div class="col-sm-12 form-inline"><strong>KPI名称:</strong>'+node.KpiName+'</div>';
	if(node.CalcKpiIndexDic.Id>0){
		htmlstr+='<div class="col-sm-6 form-inline"><strong>计算脚本:</strong><div style="width:60%"><form><textarea class="form-control" style="width:100%" rows="1" disabled>'+node.CalcKpiIndexDic.Script+'</textarea></form></div></div>';
		htmlstr+='<div class="col-sm-6 form-inline"><strong>脚本说明:</strong>'+node.CalcKpiIndexDic.KpiNameCn+'</div>';
	}else{
		htmlstr+='<div class="col-sm-12 form-inline"><strong>计算脚本:</strong><div style="width:100%"><form><textarea class="form-control" style="width:100%" disabled>'+node.Script+'</textarea></form></div></div>';
	}
	htmlstr+='<div class="col-sm-3 form-inline"><strong>计量单位:</strong>'+node.CfgUnit+'</div>';
	htmlstr+='<div class="col-sm-3 form-inline"><strong>备注:</strong>'+node.Description+'</div>';
	htmlstr+='<div class="col-sm-3 form-inline"><strong>计算周期:</strong>'+getPeriodStr(node.Period)+'</div>';
	htmlstr+='<div class="col-sm-3 form-inline"><strong>计算偏移时间(分钟):</strong>'+node.OffsetMinutes+'</div>';
	htmlstr+='<div class="col-sm-3 form-inline"><strong>开始时间:</strong>'+node.StartTime+'</div>';
	htmlstr+='<div class="col-sm-3 form-inline"><strong>最新数据时间:</strong>'+node.LastCalcTime+'</div>';
	htmlstr+='<div class="col-sm-3 form-inline"><strong>最后计算时间:</strong>'+node.UpdateTime+'</div>';
	htmlstr+='<div class="col-sm-3 form-inline"><strong>创建时间:</strong>'+node.CreateTime+'</div>';
	$("#CfgMsg").html(htmlstr);
}

function getTagTypeList(ajaxdata){
	/*
	[
		{
			"cnt": "635",
			"name": "巡检变量",
			"tag_type": "check"
		},
		{
			"cnt": "919",
			"name": "物耗变量",
			"tag_type": "goods"
		},
		{
			"cnt": "1017",
			"name": "过程变量",
			"tag_type": "process"
		},
		{
			"cnt": "216",
			"name": "非实时过程变量",
			"tag_type": "process2"
		},
		{
			"cnt": "34",
			"name": "样本化验变量",
			"tag_type": "sample"
		}
	]
	*/
	var nodes = eval("("+ajaxdata+")"); 
	PIE_DATA.splice(0,PIE_DATA.length);//清空数组
	PIE_LEGEND.splice(0,PIE_LEGEND.length);//清空数组
	var htmlstr='<div class="col-12 border-bottom"><h5>KPI统计信息</h5></div>';
	htmlstr+='<div class="col-12"><table class="table table-striped table-hover table-sm"><tbody><tr><th>序号</th><th>主变量类型</th><th>主变量类型名称</th><th>指标数量</th></tr>';
	var kpicnt=0;
	for(var i=0;i<nodes.length;i++){
		var node=nodes[i];
		htmlstr+='<tr>';
		htmlstr+='<td>'+(i+1)+'</td>';
		htmlstr+='<td>'+node.tag_type+'</td>';
		htmlstr+='<td>'+node.name+'</td>';
		htmlstr+='<td>'+node.cnt+'</td>';
		htmlstr+='</tr>';
		kpicnt+=parseInt(node.cnt);
		PIE_LEGEND[i]=node.name;
		PIE_DATA[i]={value:node.cnt,name:node.name};
	}
	htmlstr+='<tr>';
	htmlstr+='<th colspan="3">合计</th>';
	htmlstr+='<th>'+kpicnt+'</th>';
	htmlstr+='</tr>';
	htmlstr+='</tbody></table></div>';
	$("#CfgMsg").html(htmlstr);
	refreshEcharts_Pie('KPI统计信息');
}

function getKpiTagTypePeriodInfo(ajaxdata){
	/*
	[
		{
			"cnt": "70",
			"period": "-4"
		},
		{
			"cnt": "320",
			"period": "-3"
		},
		{
			"cnt": "529",
			"period": "-2"
		}
	]
	*/
	var nodes = eval("("+ajaxdata+")"); 
	PIE_DATA.splice(0,PIE_DATA.length);//清空数组
	PIE_LEGEND.splice(0,PIE_LEGEND.length);//清空数组
	var htmlstr='<div class="col-12 border-bottom"><h5>KPI统计信息</h5></div>';
	htmlstr+='<div class="col-12"><table class="table table-striped table-hover table-sm"><tbody><tr><th>序号</th><th>周期类型</th><th>指标数量</th></tr>';
	var kpicnt=0;
	for(var i=0;i<nodes.length;i++){
		var node=nodes[i];
		htmlstr+='<tr>';
		htmlstr+='<td>'+(i+1)+'</td>';
		htmlstr+='<td>'+getPeriodStr(node.period)+'</td>';
		htmlstr+='<td>'+node.cnt+'</td>';
		htmlstr+='</tr>';
		kpicnt+=parseInt(node.cnt);
		PIE_LEGEND[i]=getPeriodStr(node.period);
		PIE_DATA[i]={value:node.cnt,name:getPeriodStr(node.period)};
	}
	htmlstr+='<tr>';
	htmlstr+='<th colspan="2">合计</th>';
	htmlstr+='<th>'+kpicnt+'</th>';
	htmlstr+='</tr>';
	htmlstr+='</tbody></table></div>';
	$("#CfgMsg").html(htmlstr);
	refreshEcharts_Pie('KPI统计信息');
}
function getPeriodStr(period){
	var periodstr = '每小时';
	switch(parseInt(period)){
		case -1:
			periodstr='每小时';
			break;
		case -2:
			periodstr='每班';
			break;
		case -3:
			periodstr='每日';
			break;
		case -4:
			periodstr='每月';
			break;
		case -5:
			periodstr='每季度';
			break;
		case -6:
			periodstr='每年';
			break;
		case 0:
			periodstr='停止';
			break;
		default:
			periodstr='每'+period+'秒';
			break;
	}
	return periodstr;
}

function timeDiff(bgtime,endtime){
	if(bgtime.length > 10 && endtime.length > 10){
		var bgstemp= new Date(bgtime.replace(/T/," "));//开始时间
		var edstemp= new Date(endtime.replace(/T/," "));//结束时间
		var diff=edstemp.getTime() - bgstemp.getTime();
		diff /=(1000*3600);
		return DataToFixed(diff,"float",2);
	}
	return ""
}

//解析数据信息
function getDatas(ajaxdata){
	/*
	[
		{
			"Id": 12087739,
			"TagType": "process",
			"TagId": 4615,
			"TagName": "x3_asl_asl-xc1_SK1_ZXS1_S1-013_SD1-010_run:1",
			"KpiConfigListId": 102,
			"KpiKey": "on_time",
			"KpiTag": "x3_asl_asl-xc1_SK1_ZXS1_S1-013_SD1-010_run:1__-2_0_on_time",
			"KpiName": "破碎一车间_1号破碎机电机_运行__累计时间",
			"KpiPeriod": -2,
			"KpiValue": 28800,
			"CalcEndingTime": "2020-07-10 01:00:00",
			"InbaseTime": "2020-07-10 01:02:54"
		}
	]
	*/
	var datas = eval("("+ajaxdata+")"); 
	DATAS.splice(0,DATAS.length);//清空数组
	DATA_TIME.splice(0,DATA_TIME.length);//清空数组
	DATA_VALUE.splice(0,DATA_VALUE.length);//清空数组
	DATAS = datas;
	var htmlstr='<table class="table table-striped table-hover table-sm"><tr colspan="5"><h5>KPI计算结果数据</h5></tr><thead class="thead-light">';
	htmlstr+='<tr><th>序号</th><th>数据时间</th><th>计算时间</th><th>数值</th><th>单位</th></tr></thead><tbody>';
	var kpisum=0.0,max=0,min=0,avg=0,dot=2;
	for(var i=0;i<datas.length;i++){
		var data=datas[i];
		if (i==0){
			min=data.KpiValue;
			max=data.KpiValue;
		}
		if(data.KpiValue>max){
			max=data.KpiValue;
		}
		if(data.KpiValue<min){
			min=data.KpiValue;
		}
		kpisum+=data.KpiValue;
	
		if(data.KpiValue<1 && data.KpiValue!=0){
			dot=3;
			if(data.KpiValue<0.1 && data.KpiValue!=0){
				dot=4;
			}
		}
		DATA_VALUE[i]=DataToFixed(data.KpiValue,"float",dot)
		DATA_TIME[i] =data.CalcEndingTime

		htmlstr+='<tr>';
		htmlstr+='<td>'+(i+1)+'</td>';
		htmlstr+='<td>'+data.CalcEndingTime+'</td>';
		htmlstr+='<td>'+data.InbaseTime+'</td>';
		htmlstr+='<td>'+DataToFixed(data.KpiValue,"float",dot)+'</td>';
		htmlstr+='<td>'+NODES_CFG.CfgUnit+'</td>';
		htmlstr+='</tr>';		
	}
	htmlstr+='</tbody></table>';
	if (i>0){
		avg=DataToFixed(kpisum/i,"float",dot);
	}
	htmlstr+='<strong>最小值:</strong>'+DataToFixed(min,"float",dot)+'&nbsp;';
	htmlstr+='<strong>最大值:</strong>'+DataToFixed(max,"float",dot)+'&nbsp;';
	htmlstr+='<strong>平均值:</strong>'+avg+'&nbsp;';
	htmlstr+='<strong>和:</strong>'+DataToFixed(kpisum,"float",dot);
	$("#DataMsg").html(htmlstr);
	refreshEcharts_Bar();
}
$(document).ready(function (){
	$("#exit").after('退出');
})
