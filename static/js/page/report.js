//=========变量定义区域===========================================================
var NODE_IS_REPORT=false; //所选节点是报表
var HAS_SELECT_REPORT = false;//已经选择了报表
var NODE_ID=0;//所选层级的ID
var REPORT_ID=0;//当前显示的报表的ID
var NODES_KVA=new Array();//层级节点KV值列表,ID为KEY
var LAST_SELECTED_FILE="";//上一个被选中的报表
var CELLSMAP=new Array();//单元格Map,以单元格坐标字符串为key

//==================动作响应区域==================================================
function zTreeOnClick(event, treeId, treeNode) {
	NODE_ID = treeNode.id;
	if (treeNode.isParent > 0){
		NODE_IS_REPORT = false;//所选节点是文件夹NewLevel
	}else{
		NODE_IS_REPORT = true;//所选节点是报表
	}
	$("#FileVew").hide();
	requestChildNodes(treeNode.levelcode);//载入子集
}
//响应鼠标双击
function zTreeOnDbClick(event, treeId, treeNode){
	
}
//页面初始化
function pageInit(){
	initNodesKva(NODESMSG);
	getReportsDatas(NODESMSG);
	$("#selectTdValue").attr("style","height:="+(winH-70)+"px; max-height:"+(winH)+"px;overflow: auto;width: 100%;outline:none;");
}

function initNodesKva(datas){
	NODES_KVA[0]={Name:"",BaseTime:"2006-01-01 00:00:00",ShiftHour:8};
	for (var j=0;j<datas.length;j++){
		var node=datas[j];
		NODES_KVA[node.Id]={
			DistributedId : node.DistributedId,     //分布式计算ID，为0的时候不区分计算引擎，不为0的时候需要id适配的计算引擎进行计算
			Name          : node.Name, //名称
			Pid           : node.Pid,     //父级菜单ID
			WorkShop      : node.Workshop==null?0:node.Workshop.Id,   //所属车间ID
			Level         : node.Level ,//层级深度
			LevelCode     : node.LevelCode ,//层级码
			Folder        : node.Folder, //是否是文件夹，1-是，0-否
			Debug         : node.Debug, //是否是调试模式，1-是，0-否
			Seq           : node.Seq,    //排序号
			Remark        : node.Remark, //备注
			TemplateUrl   : node.TemplateUrl, //模板文件路径
			TemplateFile  : node.TemplateFile ,//模板文件名称
			ResultUrl     : node.ResultUrl,    //结果地址
			StartTime     : node.StartTime,    //统计计算开始起作用的时间
			Period        : node.Period,      //计算周期,详见KPI表
			OffsetMinutes : node.OffsetMinutes, //偏移时间
			LastCalcTime  : node.LastCalcTime,     //最后计算时间
			BaseTime      : node.BaseTime,     //基准时间
			ShiftHour     : node.ShiftHour,    //每班工作时间
			Status        : node.Status //1有效 0无效
		};
		if(REPORT_ID==0 && node.Folder==0)
			REPORT_ID=node.Id;
	}
}

//请求数据
function requestDatas(){
	if(HAS_SELECT_REPORT == true){//已经选择了有效节点
		requestTplLists(REPORT_ID)
	}
}
//下载
function onTplDown(filename){
	var urlstr = "api/download?filepath="+NODES_KVA[REPORT_ID].ResultUrl+"&filename="+filename+"&id="+REPORT_ID;
	//console.log(urlstr);
	window.location.href=urlstr;
}
//预览
function onTplView(filename,thisid){
	$("#"+thisid).attr("class","btn btn-primary btn-sm");
	if(LAST_SELECTED_FILE!=""){
		$("#"+LAST_SELECTED_FILE).attr("class","btn btn-outline-primary btn-sm");
	}
	LAST_SELECTED_FILE=thisid;
	var urlstr = "api/viewexcel?filepath="+NODES_KVA[REPORT_ID].ResultUrl+"&filename="+filename+"&id="+REPORT_ID;
	loadViewFile(urlstr);
}
//=========AJAX请求定义区域=======================================================
//读取所选层级下的所有子集
function requestChildNodes(levelcode){
	HAS_SELECT_REPORT = false;
	var urlstr = "api/getreportchildnodes?levelcode="+levelcode;
	loadChildNodes(urlstr);
}
//读取所选报表的模板列表
function requestTplLists(id){
	REPORT_ID=id;
	HAS_SELECT_REPORT = true;
	$("#FileVew").hide();
	var urlstr = "api/getreportreultlist?id="+id+"&begintime="+InputDateTimeToString("YYYY-mm-dd HH:MM:SS",$("#BeginTime").val())+"&endtime="+InputDateTimeToString("YYYY-mm-dd HH:MM:SS",$("#EndTime").val());
	loadReportLists(urlstr);
}
//=========AJAX函数定义区域=======================================================
function loadChildNodes(urlstr)//读取所选层级下的所有子集
{	
	$("#DataFrame").html('<div class="alert alert-warning">正在加载数据……</div>');
	$("#TplLists").hide();
	//调用公用的loadXMLDoc函数
    loadXMLDoc(urlstr,function()
    {
        if (xmlhttp.readyState==4 && xmlhttp.status==200)//请求处理完成，且状态OK
        {//添加请求完成后的处理功能---------------------------------------
			getReportsDatas(eval("("+xmlhttp.responseText+")"));//解读数据
			//下一步：	
			if (NODE_IS_REPORT){
				requestTplLists(NODE_ID)
			} 	
        }//请求完成后的处理功能结束---------------------------------------
    });
}

function loadReportLists(urlstr)//读取所选报表的模板列表
{	
	$("#TplLists").show();
	$("#TplLists").html('<div class="alert alert-warning">正在加载数据……</div>');
	//调用公用的loadXMLDoc函数
    loadXMLDoc(urlstr,function()
    {
        if (xmlhttp.readyState==4 && xmlhttp.status==200)//请求处理完成，且状态OK
        {//添加请求完成后的处理功能---------------------------------------
			getReportLists(eval("("+xmlhttp.responseText+")"));//解读数据
			//下一步：		
        }//请求完成后的处理功能结束---------------------------------------
    });
}
function loadViewFile(urlstr)//加载文件
{	
	$("#FileVew").show();
	$("#FileVew").html('<div class="alert alert-warning">正在加载文件……</div>');
	//调用公用的loadXMLDoc函数
    loadXMLDoc(urlstr,function()
    {
        if (xmlhttp.readyState==4 && xmlhttp.status==200)//请求处理完成，且状态OK
		{//添加请求完成后的处理功能---------------------------------------
			$("#FileVew").html(xmlhttp.responseText);//显示数据
			$("#FileVew").attr("style","max-height:"+(winH-130)+"px;overflow: auto;");
			$("#EditForm").attr("style","width: 100%;display: flex;padding:1px;");//显示编辑区域
			$("#selectTdValue").val("");
			$("#CellAxis").val("");
			//下一步：		
        }//请求完成后的处理功能结束---------------------------------------
    });
}
//==================AJAX数据解析区域==================================================
//解析报表列表
function getReportsDatas(datas){
	var htmlstr='<table class="table table-striped table-hover table-sm"><thead class="thead-light"><tr><th style="width:50px">序号</th><th style="width:150px">名称</th><th  style="width:150px">父层级</th><th style="width:50px">模板</th><th style="width:50px">调试</th><th style="width:50px">状态</th><th style="width:50px">排序号</th><th style="width:150px">计算截至时间</th><th style="width:150px">更新时间</th><th style="width:150px">操作</tr></thead><tbody>';
	if (datas!=null){
		var j=0
		for (var i=0;i<datas.length;i++){
			var node=datas[i];
			var status = '';
			var isdebug= '';
			var isreport='<span class="badge badge-secondary">否</span>';
			var lasttime='';
			if (node.Folder==0){
				isreport='<span class="badge badge-success">是</span>';
				lasttime = node.LastCalcTime;
				if (node.Status==0){
					status = '<span class="badge badge-secondary">停用</span>';
				}else{
					status = '<span class="badge badge-success">启用</span>';
				}
				if (node.Debug==0){
					isdebug = '<span class="badge badge-success">否</span>';
				}else{
					isdebug = '<span class="badge badge-danger">是</span>';
				}
				if(j==0){
					requestTplLists(node.Id);
				}
				j++;
				var btedit='<button class="btn btn-outline-primary btn-sm" onclick="requestTplLists('+node.Id+');">查看结果</button>';
				//var btdelete='<button class="btn btn-outline-danger btn-sm" onclick="onDeleteLevel('+node.Id+');">删除</button>';
				htmlstr+='<tr><td>'+j+'</td><td>'+node.Name+'</td><td>'+NODES_KVA[node.Pid].Name+'</td><td>'+isreport+'</td><td>'+isdebug+'</td><td>'+status+'</td><td>'+node.Seq+'</td><td>'+lasttime+'</td><td>'+node.UpdateTime+'</td><td><div class="btn-group">'+btedit+'</div></td></tr>';
			}
			
		}
	}
	htmlstr+='</tbody></table>';
	$("#DataFrame").html(htmlstr);
}
//冒泡排序
function bubbleSort(arr) {
    var i = arr.length, j;
    var tempExchangVal;
    while (i > 0) {
        for (j = 0; j < i - 1; j++) {
			var t1=new Date(arr[j].FileTime.replace(/T/," "));
			var t2=new Date(arr[j+1].FileTime.replace(/T/," "));
            if (t1.getTime() < t2.getTime()) {
                tempExchangVal = arr[j];
                arr[j] = arr[j + 1];
                arr[j + 1] = tempExchangVal;
            }
        }
        i--;
    }
    return arr;
}

//解析模板列表
function getReportLists(datas){
	/*
		FileName string //文件名
		FileTime string //文件上传时间
		ModTime  string //文件编辑时间
		Size     int64  //单位:字节
	*/
	var htmlstr=`<div class="col-12"><strong>`+NODES_KVA[REPORT_ID].Name+`</strong></div>
	<table class="table table-striped table-hover table-sm"><thead class="thead-light"><tr><th style="width:50px">序号</th><th style="width:200px">报表时间</th><th style="width:200px">生成/编辑时间</th><th>文件大小</th><th>操作</th></tr></thead><tbody>`;
	if (datas!=null){
		datas = bubbleSort(datas);
		for (var i=0;i<datas.length;i++){
			var node=datas[i];
			var fname = node.FileName.split(".",1);
			var btview=`<button class="btn btn-outline-primary btn-sm" id="`+fname+`" onclick="onTplView('`+node.FileName+`',this.id);" title="在线预览文件">预览</button>`;
			var btdown=`<button class="btn btn-outline-success btn-sm" onclick="onTplDown('`+node.FileName+`');" title="下载报表文件到本地">下载</button>`;

			if(i==0){
				onTplView(node.FileName);
				LAST_SELECTED_FILE = fname;
				btview=`<button class="btn btn-primary btn-sm" id="`+fname+`" onclick="onTplView('`+node.FileName+`',this.id);" title="在线预览文件">预览</button>`;
			}
			

			htmlstr+='<tr><td>'+(i+1)+'</td><td>'+node.FileTime+'</td><td>'+node.ModTime+'</td><td>'+DataToFixed(node.Size/1024,"float",2)+'KB</td><td><div class="btn-group">'+btview+btdown+'</div></td></tr>';
		}
	}
	htmlstr+='</tbody></table>';
	$("#TplLists").html(htmlstr);
	
}

//时间范围选择框设置
function timeRangeSelectorSet(){
	if(NODES_KVA[REPORT_ID].BaseTime.length > 10){//已经选择了样本模板,取消禁用选择本班、今日、本月、上月
		if (NODES_KVA[REPORT_ID].ShiftHour > 0){
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
	requestDatas();//请求数据
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
		timediff = NODES_KVA[REPORT_ID].ShiftHour * 3600 * 1000;
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
		getBeginTimeOfMonth(NODES_KVA[REPORT_ID].BaseTime,begintime);
		break;
	default:
		bgstemp.setTime(bgstemp.getTime() - timediff);
		edstemp.setTime(edstemp.getTime() - timediff);
		$("#BeginTime").val(DateFormat("YYYY-mm-ddTHH:MM",bgstemp));
		$("#EndTime").val(DateFormat("YYYY-mm-ddTHH:MM",edstemp));
		break;
	}
	
	requestDatas();//请求数据
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
		timediff = NODES_KVA[REPORT_ID].ShiftHour * 3600 * 1000;
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

	requestDatas();//请求数据
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
		getBeginTimeInDay(NODES_KVA[REPORT_ID].BaseTime,DateFormat("YYYY-mm-ddTHH:MM",nowTime),NODES_KVA[REPORT_ID].ShiftHour);
		break;
	case '2'://今日
		getBeginTimeInDay(NODES_KVA[REPORT_ID].BaseTime,DateFormat("YYYY-mm-ddTHH:MM",nowTime),24);
		break;
	case '3'://本月
		getBeginTimeOfMonth(NODES_KVA[REPORT_ID].BaseTime,DateFormat("YYYY-mm-ddTHH:MM",nowTime))
		break;
	default:
		bgstemp.setTime(edstemp.getTime() - tdiff*60*1000);
		$("#BeginTime").val(DateFormat("YYYY-mm-ddTHH:MM",bgstemp));
		break;
	}	

	requestDatas();//请求数据
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

var LAST_SELECTED_CELL='A1';//上次选择的单元格

function onClickCell(axis){
    var msg = $("#"+axis).text();
	$("#"+axis).attr("style","border:blue solid 2px;");
	if(LAST_SELECTED_CELL!=axis){//如果与上次单击的不是同一个单元格
		$("#"+LAST_SELECTED_CELL).removeAttr("style");
	}
	LAST_SELECTED_CELL = axis;
	$("#coordinate").html(axis);
	$("#CellAxis").val(axis);
	$("#selectTdValue").val(msg);
	if(CELLSMAP[axis].Formula.length>0){
		$("#selectTdValue").val("="+CELLSMAP[axis].Formula);
	}
}
$(document).ready(function() {
	//Ztree国际化
	$("#exit").after('退出');
	$("#ExpandTreeNode").html('展开');
	$("#CollapseTreeNode").html('折叠');
	$("#HideTreeNode").html('隐藏');
	$("#SearchTreeNode").attr("placeholder",'搜索');

	$('#ExpandTreeNode').attr("title","未选中节点时展开所有节点,选中节点时展开选中节点");
	$('#CollapseTreeNode').attr("title","未选中节点时折叠所有节点,选中节点时折叠选中节点");
	$("#BeginTimes").text('起始时间：');
	$("#EndTimes").text('结束时间：');
	$("#TimeRanges").text('时间范围:');

	$("#custom").text("自定义");
	$("#8h").text("8小时");
	$("#12h").text("12小时");
	$("#24h").text("24小时");
	$("#7d").text("7天");
	$("#op_this_shift").text("班");
	$("#op_this_day").text("日");
	$("#op_this_month").text("月");

	$('#Last').attr("title","前一时间段!");
	$('#Next').attr("title","后一时间段!");

	$('#OperationMsg').html('<strong id = "works">注意：</strong>  在线预览的格式与实际Excel文件格式并非完全一致，请以下载后的Excel文件为准！在线预览目前只具备简单的Excel公式的计算能力，使用Excel公式的运算结果仅供参考，请以下载后的Excel文件中的结果为准！');
	$('#CellAxis').attr("title","选择单元格的坐标");
	$('#selectTdValue').attr("title","选择单元格的值或公式");
	$("#CheckCellValue").html('检验');


});