//=========全局变量定义区域========================================================
var TAG;//当前选中的变量
var TAGS=[];//当前选中的变量组,key为变量ID,值为变量信息结构
var TAGS_SERIAL=[];//变量组序列数组,key为序号,值为Tag
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
var TIME_CHANGE;//时间范围发生了改变
//=========动作响应区域===========================================================
function zTreeOnClick(event, treeId, treeNode) {
	SELECT_LEVEL_CODE=treeNode.treelevel;
	SELECT_NAME=treeNode.name;
	SELECT_IS_TAG=treeNode.istag;
	if(treeNode.istag==1){
		TAG=getTagInfo(treeNode);

		$("#TagName").text(TAG.Name);
		$("#TagUnit").text(TAG.Unit);
		TAG_HAVE_SELECTED = true;
		requestSnapshot(TAG.FullName);
		$("#AddTagToTable").removeAttr("disabled");
		$("#AddTagToTable").attr("class","btn btn-primary");
	}else{
		$("#AddTagToTable").attr("class","btn btn-outline-primary");
		$("#AddTagToTable").attr("disabled","disabled");
	}
};
//响应鼠标双击,添加所选变量入列表
function zTreeOnDbClick(event, treeId, treeNode){
	if(treeNode.istag==1){//如果所选变量是tag		
		if (!!!TAGS[treeNode.itemid]) { //值不存在
			TAGS_SERIAL[TAGS_SERIAL.length] = getTagInfo(treeNode);//保存序列ID
		}	
		TAGS[treeNode.itemid]=getTagInfo(treeNode);//移入变量列表
		showTagsTable();
	}
};
//相应按钮,添加所选变量入列表
function onAddSelectTagToTable(){
	if(SELECT_IS_TAG==1){//如果所选变量是tag
		if (!!!TAGS[TAG.Id]) { //值不存在
			TAGS_SERIAL[TAGS_SERIAL.length] = TAG;//保存序列ID
		}
		TAGS[TAG.Id]=TAG;//移入变量列表
		showTagsTable();
	}
}
//显示已选变量
function showTagsTable(){
	var tbstr='<h5>Selected comparison variable<small id="TagsTableSmallTitle" class="text-muted">You can double-click the variable label in the tree structure on the left or select the label and click it 👆 To add variables to the list</small></h5><table class="table table-striped table-hover table-sm"><thead class="thead-light"><tr><th>Number</th><th>name</th><th>label</th><th>type</th><th>unit</th><th>remove</th></tr></thead><tbody>';
	for(var i=0;i<TAGS_SERIAL.length;i++){
		var key = TAGS_SERIAL[i].Id;
		tbstr +='<tr id="TagsRow_'+key+'"><td>'+(i+1)+'</td><td>'+TAGS[key].Name+'</td><td>'+TAGS[key].FullName+'</td><td>'+TAGS[key].TagType+'</td><td>'+TAGS[key].Unit+'</td><td><div><button type="button" class="btn btn-outline-danger btn-sm" onclick="onRemoveTag('+key+')" title="Remove the variable from the list">remove</button></div></td></tr>';
	}
	tbstr +='</tbody></table>';
	$("#SelectedTags").html(tbstr);
	HIS_INTERVAL_TIME.splice(0,HIS_INTERVAL_TIME.length);//清空数组
	HIS_INTERVAL_DATA.splice(0,HIS_INTERVAL_DATA.length);//清空数组
	if(TAGS_SERIAL.length > 0){
		requestHistoryInterval(TAGS_SERIAL[0].FullName,0);
		$("#Echarts_HisIntervalSerial").show();
	}else{
		$("#Echarts_HisIntervalSerial").hide();
		$("#HisSerialRemark").hide();
	}
}
//从已选列表中移除变量
function onRemoveTag(key){
	if (!!TAGS[key]) { 
		delete (TAGS[key]);//在列表中删除
		for(var i=0;i<TAGS_SERIAL.length;i++){
			if(TAGS_SERIAL[i].Id==key){
				TAGS_SERIAL.splice(i,1);//删除序列
				break;
			}
		}
		showTagsTable();
	} 
}
function getTagInfo(treeNode){
	var tag;
	if(treeNode.istag==1){
		var datatype;//数据类型
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
		var unit=treeNode.unit//单位
		if(unit.length ==0){//没有定义单位的时候
			if(treeNode.seq == 1){//如果是BOOL变量
				unit="Nil";
			}else{
				unit="Not set";
			}
		}
		tag={
			Name:treeNode.name,//变量描述名称
			FullName:treeNode.treelevel,//变量层级码
			Id:treeNode.itemid,//变量id
			DotNum:treeNode.dotnum,//小数点数量
			TagType:datatype,//数据类型
			Unit:unit//单位
		};
	}
	return tag;
}
//时间输入框的值发生改变
function onTimeChange(){
	timeDiffCheck();
	TIME_CHANGE=true;//时间范围繁盛了改变
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
	TIME_CHANGE=true;//时间范围繁盛了改变
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
	TIME_CHANGE=true;//时间范围繁盛了改变
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
	TIME_CHANGE=true;//时间范围繁盛了改变
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
		alert('No variable matching variable name ['+TAG.FullName+'] was found in the database, please check!');
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
		SD:DataToFixed(dtarr.SD,TAG.TagType,TAG.DotNum),             //总体标准偏差(高级)
		STDDEV:DataToFixed(dtarr.STDDEV,TAG.TagType,TAG.DotNum),     //样本标准偏差(高级)
		SE:DataToFixed(dtarr.SE,TAG.TagType,TAG.DotNum),             //标准误差(SE = STDDEV / PointCnt)(高级)
		Ske:DataToFixed(dtarr.Ske,TAG.TagType,TAG.DotNum),           //偏度(高级)
		Kur:DataToFixed(dtarr.Kur,TAG.TagType,TAG.DotNum),           //峰度(高级)
		Mode:DataToFixed(dtarr.Mode,TAG.TagType,TAG.DotNum),         //众数(高级)
		Median:DataToFixed(dtarr.Median,TAG.TagType,TAG.DotNum),     //中位数(高级)
		GroupDist:dtarr.GroupDist       //组距GroupDistance(高级),DataGroup中两组数之间的距离
	};
	HIS_SUMMARY = suma;
	HIS_SUMMARY_TABLE='<table class="table table-striped table-hover table-sm"><thead class="thead-light"><tr><th>Min</th><th>Max</th><th>Range</th><th>ArithmeticMean</th><th>weightedMean</th><th>Mode</th><th>Median</th><th>Sum</th><th>Difference</th><th>PositiveDifference</th><th>Area</th><th>Points</th><th>SD</th><th>Skewness</th><th>Kurtosis</th></tr></thead><tbody>';
	HIS_SUMMARY_TABLE +='<tr><td>'+suma.Min+'</td><td>'+suma.Max+'</td><td>'+suma.Range+'</td><td>'+suma.Mean+'</td><td>'+suma.PowerAvg+'</td><td>'+suma.Mode+'</td><td>'+suma.Median+'</td><td>'+suma.Sum+'</td><td>'+suma.Diff+'</td><td>'+suma.PlusDiff+'</td><td>'+suma.Total+'</td><td>'+suma.PointCnt+'</td><td>'+suma.SD+'</td><td>'+suma.Ske+'</td><td>'+suma.Kur+'</td></tr>';
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
	for(var key in dgroup){
		HIS_SUM_GROUP_KEY[key] = DataToFixed(parseFloat(key)*parseFloat(dtarr.GroupDist)+parseFloat(suma.Min),TAG.TagType,TAG.DotNum);
		HIS_SUM_GROUP_VAL[key] = DataToFixed(dgroup[key]/point * 100,'int',2);
	}
	for(var i=0;i<HIS_TIME.length;i++){
		HIS_INCREMENT_DATA[i] = DataToFixed(increment[HIS_TIME[i]],TAG.TagType,TAG.DotNum);
		HIS_TIME[i] = HIS_TIME[i].split(".",1);//去掉毫秒
	}
	refreshEcharts_his();//刷新Echarts 
	//refreshEcharts_hisGroup();//刷新Echarts 
}
//接收AJAX反馈的等间隔历史数据并解析
function getTagHistoryInterval(ajaxdata,groupnum){//根据Ajax反馈的结果更新Tag的实时数据
	var dtarr = eval("("+ajaxdata+")"); 
	var tagtype=TAGS_SERIAL[groupnum].TagType;
	var dotnum=TAGS_SERIAL[groupnum].DotNum;
	var tagfullname=TAGS_SERIAL[groupnum].FullName;
	var histime=[];
	var hisdata=[];
	if (groupnum==0){
		HIS_INTERVAL_DATA.splice(0,HIS_INTERVAL_DATA.length);//清空数组
		HIS_INTERVAL_TIME.splice(0,HIS_INTERVAL_TIME.length);//清空数组
	}
	if(dtarr != null){
		var hisdata=dtarr[tagfullname];
		for(var i=0;i< hisdata.length;i++){
			histime[i] = hisdata[i].Time;//更新时间戳:转换为当地时间格式
			hisdata[i] = DataToFixed(hisdata[i].Value,tagtype,dotnum);//更新TagValue
		}

		/*for(var i=0;i< dtarr.length;i++){
			histime[i] = dtarr[i].Datatime.split(".",1);//更新时间戳:去掉毫秒
			hisdata[i] = DataToFixed(dtarr[i].Value,tagtype,dotnum);//更新TagValue
		}*/
		if(histime.length > HIS_INTERVAL_TIME.length){
			HIS_INTERVAL_TIME = histime;
		}
		HIS_INTERVAL_DATA[groupnum] = hisdata;
	}else{
		$("#TagsRow_"+TAGS_SERIAL[groupnum].Id).attr("class","table-danger");
		$("#TagsTableSmallTitle").html('<span class="text-danger">Red background line loading data failed！</span>');
	}

	refreshEcharts_hisInterval();//刷新Echarts 

	groupnum++;
	if(groupnum < TAGS_SERIAL.length){
		requestHistoryInterval(TAGS_SERIAL[groupnum].FullName,groupnum);
	}
}

//接收AJAX反馈的历史数据并解析
function getTagHistory(ajaxdata){//根据Ajax反馈的结果更新Tag的实时数据
	var dtarr = eval("("+ajaxdata+")"); 
	HIS_TIME.splice(0,HIS_TIME.length);//清空数组
	HIS_DATA.splice(0,HIS_DATA.length);//清空数组
	HIS_TABLE='<hr/><h3>Original historical data</h1><br/><table class="table table-striped table-hover table-sm"><thead class="thead-light"><tr><th>SerialNumber</th><th>Time </th><th>Data </th></tr></thead><tbody>';
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
function requestHistoryInterval(tagname,groupnum){
	var urlstr = "api/hisinterval?tagname="+tagname+"&begintime="+InputDateTimeToString("YYYY-mm-dd HH:MM:SS",$("#BeginTime").val())+"&endtime="+InputDateTimeToString("YYYY-mm-dd HH:MM:SS",$("#EndTime").val())+"&interval="+$("#Interval").val();
	loadTagHistoryInterval(urlstr,groupnum);
}
//读取原始历史数据请求
function requestHistory(tagname){
	var urlstr = "api/history?tagname="+tagname+"&begintime="+InputDateTimeToString("YYYY-mm-dd HH:MM:SS",$("#BeginTime").val())+"&endtime="+InputDateTimeToString("YYYY-mm-dd HH:MM:SS",$("#EndTime").val());
	loadTagHistory(urlstr);
}

//=========AJAX函数定义区域=======================================================
function loadTagSnapshot(urlstr)//从数据库中读取单一变量的最新值
{
	$("#LoadDataMsg").html('<div class="text-info">Loading snapshot data for the currently selected variable……</div>');
	//调用公用的loadXMLDoc函数
    loadXMLDoc(urlstr,function()
    {
        if (xmlhttp.readyState==4 && xmlhttp.status==200)//请求处理完成，且状态OK
        {//添加请求完成后的处理功能---------------------------------------
			getTagSnapshotData(xmlhttp.responseText);//解读数据
			$("#LoadDataMsg").html('<div class="text-success">Snapshot data loading complete</div>');
			//下一步：加载历史数据
			requestHistory(TAG.FullName);
        }//请求完成后的处理功能结束---------------------------------------
    });
}
function loadTagHistory(urlstr)//从数据库中读取单一变量指定时间段的原始历史数据
{
	$("#LoadDataMsg").html('<div class="text-info">Loading historical data for the currently selected variable……</div>');
	//调用公用的loadXMLDoc函数
    loadXMLDoc(urlstr,function()
    {
        if (xmlhttp.readyState==4 && xmlhttp.status==200)//请求处理完成，且状态OK
        {//添加请求完成后的处理功能---------------------------------------
			getTagHistory(xmlhttp.responseText);
			$("#LoadDataMsg").html('<div class="text-success">Historical data loading complete</div>');
			//下一步：加载历史统计数据
			requestHistorySummary(TAG.FullName);			
        }//请求完成后的处理功能结束---------------------------------------
    });
}
function loadTagHistorySummary(urlstr)//从数据库读取统计值
{
	$("#LoadDataMsg").html('<div class="text-info">Loading statistics for the currently selected variable……</div>');
	//调用公用的loadXMLDoc函数
    loadXMLDoc(urlstr,function()
    {
        if (xmlhttp.readyState==4 && xmlhttp.status==200)//请求处理完成，且状态OK
        {//添加请求完成后的处理功能---------------------------------------
			getTagHistorySummary(xmlhttp.responseText);
			//下一步：加载等间隔历史数据
			//requestHistoryInterval(TAG.FullName);
			$("#LoadDataMsg").html('<div class="text-success">Statistics loading complete</div>');
			if (TIME_CHANGE==true){//时间范围发生了改变
				if(TAGS_SERIAL.length > 0){
					requestHistoryInterval(TAGS_SERIAL[0].FullName,0);
				}
			}
        }//请求完成后的处理功能结束---------------------------------------
    });
}

function loadTagHistoryInterval(urlstr,groupnum)//从数据库中读取单一变量指定时间段的等间隔历史数据
{
	$("#LoadDataMsg").html('<div class="text-info">Loading data for '+(groupnum+1)+'th comparison variable……</div>');
	//调用公用的loadXMLDoc函数
    loadXMLDoc(urlstr,function()
    {
        if (xmlhttp.readyState==4 && xmlhttp.status==200)//请求处理完成，且状态OK
        {//添加请求完成后的处理功能---------------------------------------
			getTagHistoryInterval(xmlhttp.responseText,groupnum);
			$("#LoadDataMsg").html('<div class="text-success">Data loading of the ['+(groupnum+1)+']th comparison variable is completed</div>');
			TIME_CHANGE = false;
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
			ShowModal("Update tips",'<div class="alert alert-success">'+xmlhttp.responseText+'</div>');
        }//请求完成后的处理功能结束---------------------------------------
    });
}
$(document).ready(function() {
	//Ztree国际化
	$("#exit").after('Exit');
	$("#UpdateTree").html('Update')
	$('#ExpandTreeNode').attr("title","Expand all nodes when no nodes are selected, and expand selected nodes when nodes are selected")
	$("#ExpandTreeNode").html('Open');
	$('#CollapseTreeNode').attr("title","Collapse all nodes when no nodes are selected, and collapse selected nodes when nodes are selected")
	$("#CollapseTreeNode").html('Fold');
	$("#HideTreeNode").html('Hide');
	$('#SearchTreeNode').attr("placeholder","Search")
	$("#SearchTreeNode").before("Search")

	$("#BeginTimes").text('Begin Time:');
	$("#EndTimes").before('End Time:');
	$("#Interval").before("Aggregation time:")
	$("#ThelastTimestemp").text('last Timestemp:');
	$("#rd_1").before("&nbsp; 1h ")
	$("#rd_2").before("&nbsp; 8h ")
	$("#rd_3").before("&nbsp; 12h ")
	$("#rd_4").before("&nbsp; 24h ")

	$("#VariableName").html("Variable Name:")
	$("#LatestValue").html('Latest Value:');
	$("#unit").html('Unit:');
	$("#UpdateTime").text('Update Time:');
	$("#TagName").html("Nil")
	$("#TagValue").html("Nil")
	$("#TagUnit").html("Nil")
	$("#TagTime").html("Nil")

	//数据表HisSumData表头国际化
	$("#Min").html("Min")
	$("#Max").html("Max")
	$("#Range").html("Range")
	$("#ArithmeticMean").html("ArithmeticMean")
	$("#weightedMean").html("weightedMean")
	$("#Mode").html("Mode")
	$("#Median").html("Median")
	$("#Sum").html("Sum")
	$("#Difference").html("Difference")
	$("#PositiveDifference").html("PositiveDifference")
	$("#Area").html("Area")
	$("#Points").html("Points")
	$("#SD").html("SD")
	$("#Skewness").html("Skewness")
	$("#Kurtosis").html("Kurtosis")

	//数据表TagsTableSmallTitle表头国际化
	$("#SerialNumber").html("Number")
	$("#Name").html("Name")
	$("#tag").html("tag")
	$("#type").html("type")
	$("#unit1").html("unit")
	$("#removed").html("removed")

	$('#AddTagToTable').text('Add the current variable to enter the comparison variable table');
	$("#Last").html('Previous period!');
	$("#Next").html('Later period!');
	$("#AddTagToTable").attr("title","This button can only be used after the valid variable is selected");
	//$("#HisSerialRemark").html('<strong>Note:</strong> the trend chart data of the currently selected variable is the original historical data, and the data in the comparative trend chart is the data after the equal interval aggregation');
	$("#HisDataTable").html('<strong>Operation: </strong> Please select a variable node in the left structure tree to display data!');


});

function onUpdateTree(){
	var url="api/updatetagnodetree?withtag=1";
	loadUpdateTree(url);
}
