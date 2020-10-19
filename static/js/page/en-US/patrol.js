//=========变量定义区域===========================================================
var SELECT_NODE={//当前选中的样本节点
	Id:0,//样本ID
	Pid:0,//所在作业ID
	Name:"",//样本名称
	TreeLevel:"",//层级码
	BaseTime:"",//所在车间基准时间
	ShiftHour:0,//所在车间每班工作时间
	LineId:0  ,  //线路ID
	LineName:"",  //线路名
	DeptId:0  ,  //所属部门ID
	DeptName:"",  //所属部门名称
};
var HAS_SELECT_NODE=false;//已经选择了有效节点
var NODES_CFG;//节点配置信息
var DATAS=new Array();;//巡检数据
var SELECT_NODE_LEVEL;//所选节点的层级码
var SELECT_NODE_NAME;//所选节点的名称
var SELECT_LEAF;//选择了叶子节点
var LEVEL_NAMES=new Array();//节点层级名称键值对
var ONLY_SHOW_ABNORMAL=0;//仅显示异常项
var PERSON_GROUP=new Array();//按人统计的检查次数,键为检查人,值为检查次数
var CHECK_DATA=new Array();//巡检数据,键为taglist_id,值为检查值
var SITE_LIST_UNDER_NODE=new Array();//所选层级下的巡检站点信息
var CHECK_LEVEL_RELATION;//巡检层级关系图数据
//==================动作响应区域==================================================
//响应鼠标单击
function zTreeOnClick(event, treeId, treeNode) {
	$("#OperationMsg").hide();//隐藏操作提示
	$("#CfgMsg").html("");//清空数据
	$("#DataMsg").html("");//清空数据
	$("#CheckSiteMsg").html("");//清空数据
	$("#TongJi").html("");
	$("#PageEcharts").hide();

	SELECT_NODE_LEVEL = treeNode.treelevel;
	SELECT_NODE_NAME = treeNode.name;
	HAS_SELECT_NODE = true;//已经选择了有效节点
	PERSON_GROUP=new Array();
	CHECK_DATA=new Array();
	SITE_LIST_UDER_NODE=new Array();
	SELECT_NODE.Id=treeNode.id;
	SELECT_NODE.Pid=treeNode.pId;
	SELECT_NODE.Name=treeNode.name;
	SELECT_NODE.SiteId=treeNode.siteid;
	SELECT_NODE.TreeLevel=treeNode.treelevel;
	SELECT_NODE.BaseTime=treeNode.basetime;
	SELECT_NODE.ShiftHour=treeNode.shifthour;
	SELECT_NODE.LineId = treeNode.lineid;
	SELECT_NODE.LineName = treeNode.linename;
	SELECT_NODE.DeptId = treeNode.deptid;
	SELECT_NODE.DeptName = treeNode.deptname;

	if (treeNode.nodetype == 9999){
		SELECT_NODE.Id=treeNode.id/100000;
		SELECT_LEAF=true;//选择了叶子节点
		showCheckSiteMsg();
		requestPatrolCfg(SELECT_NODE.Id,1);//请求该巡检点中的点检项信息
		$("#OperationMsg").hide();//隐藏操作提示
	}else{
		SELECT_LEAF = false;//没有选择叶子节点
		requestPatrolCfg(SELECT_NODE_LEVEL,0);//请求该层级中的点检项信息
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
	initLevelNames();//初始化键值对
};
//请求数据
function requestDatas(){
	if(HAS_SELECT_NODE == true){//已经选择了有效节点
		if(SELECT_LEAF == true){//
			requestPatrolDatas(SELECT_NODE.Id,1);
		}else{
			requestPatrolCfg(SELECT_NODE_LEVEL,0);
		}
	}
}
//显示巡检站点信息
function showCheckSiteMsg(){
	var htmlstr='<div class="col-12 border-bottom"><h5>Inspection point information</h5></div>';
	htmlstr+='<div class="col form-inline"><strong>Device point ID:</strong>'+SELECT_NODE.Id+'</div>';
	htmlstr+='<div class="col form-inline"><strong>Patrol point ID:</strong>'+SELECT_NODE.SiteId+'</div>';
	htmlstr+='<div class="col form-inline"><strong>Inspection point name:</strong>'+SELECT_NODE.Name+'</div>';
	htmlstr+='<div class="col form-inline"><strong>Line:</strong>'+SELECT_NODE.LineName+'</div>';
	htmlstr+='<div class="col form-inline"><strong>Department:</strong>'+SELECT_NODE.DeptName+'</div>';
	$("#CheckSiteMsg").html(htmlstr);
}

//初始化节点层级名称键值对
function initLevelNames(){
	for(var i=0;i<zNodes.length;i++){
		var node = zNodes[i];
		if(node.nodetype < 9999){
			LEVEL_NAMES[node.treelevel]=node.name;
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
//显示巡检照片
function onShowPic(siteexeid){
	requestPicUrlBySiteExeId(siteexeid);
}
//仅显示异常信息
function onlyshowabnomal(){
	ONLY_SHOW_ABNORMAL=1-ONLY_SHOW_ABNORMAL;
	getLeafsDatas(DATAS);
}
//显示饼图Echarts
function onShowEcharts(){
	ShowModal("Inspection statistics","");
	refreshEcharts(PERSON_GROUP);
}
//显示趋势Echarts
function onShowTrenth(tagname,taglistid){
	ShowModal(tagname+"-Data trends","");
	refreshTrendEcharts(tagname,CHECK_DATA[taglistid]);
}

//显示统计分析结果
function onShowAnalyse(tagname,taglistid){
	//ShowModal(tagname+"-统计数据","");
	requestPatrolStatisticDatas(tagname,taglistid);
}

function onShowRelation(){
	ShowModal("Inspection relationship map","");
	refreshRelitionEcharts(CHECK_LEVEL_RELATION);
}
//=========AJAX请求定义区域=======================================================
//读取巡检配置信息请求
function requestPatrolCfg(patrolid,isid){
	if(isid==1){
		var urlstr = "api/patrollist?leveltag="+patrolid+"&isid="+isid;
		loadPatrolCfg(urlstr,patrolid);
	}else{
		var urlstr = "api/patrolsitelist?levelcode="+patrolid+"&begintime="+InputDateTimeToString("YYYY-mm-dd HH:MM:SS",$("#BeginTime").val())+"&endtime="+InputDateTimeToString("YYYY-mm-dd HH:MM:SS",$("#EndTime").val());
		loadPatrolCfg(urlstr,patrolid);
	}
}
//读取巡检数据信息请求
function requestPatrolDatas(idorlevel,isid){
	var urlstr = "api/patrolresult?leveltag="+idorlevel+"&begintime="+InputDateTimeToString("YYYY-mm-dd HH:MM:SS",$("#BeginTime").val())+"&endtime="+InputDateTimeToString("YYYY-mm-dd HH:MM:SS",$("#EndTime").val())+"&startonly=1&isid="+isid;
	loadPatrolDatas(urlstr);
}

//读取巡检数据统计信息请求
function requestPatrolStatisticDatas(tagname,tagid){
	var urlstr = "api/srtd/statisticauto?tagid="+tagid+"&tagtype=check&begintime="+InputDateTimeToString("YYYY-mm-dd HH:MM:SS",$("#BeginTime").val())+"&endtime="+InputDateTimeToString("YYYY-mm-dd HH:MM:SS",$("#EndTime").val())+"&filloutliers=midean&needraw=0";
	loadPatrolStatisticDatas(tagname,urlstr);
}

//读取巡检图片信息请求
function requestPicUrlBySiteExeId(siteexeid){
	var urlstr = "api/patrolpicurl?siteexeid="+siteexeid;
	loadCheckSiteExePicUrl(urlstr);
}
//=========AJAX加载定义区域=======================================================
function loadPatrolCfg(urlstr,patrolid)//读取巡检配置信息
{	
	$("#CfgMsg").html('<div class="alert alert-warning">Loading configuration information……</div>');
	//调用公用的loadXMLDoc函数
    loadXMLDoc(urlstr,function()
    {
        if (xmlhttp.readyState==4 && xmlhttp.status==200)//请求处理完成，且状态OK
        {//添加请求完成后的处理功能---------------------------------------
			if(SELECT_LEAF){//如果选择的是叶子节点,请求叶子节点的历史数据		
				getPatrolCfg(xmlhttp.responseText);//解读数据
				requestPatrolDatas(patrolid,1)
			}else{
				getCheckSiteMsg(xmlhttp.responseText);
			};
        }//请求完成后的处理功能结束---------------------------------------
    });
}
function loadPatrolDatas(urlstr)//读取巡检配置信息
{	
	$("#DataMsg").html('<div class="alert alert-warning">Loading data information……</div>');
	//调用公用的loadXMLDoc函数
    loadXMLDoc(urlstr,function()
    {
        if (xmlhttp.readyState==4 && xmlhttp.status==200)//请求处理完成，且状态OK
        {//添加请求完成后的处理功能---------------------------------------
			getPatrolDatas(xmlhttp.responseText);//解读数据
			//下一步：		
        }//请求完成后的处理功能结束---------------------------------------
    });
}

function loadPatrolStatisticDatas(tagname,urlstr)//读取巡检统计数据
{	
	//$("#DataMsg").html('<div class="alert alert-warning">正在加载数据信息……</div>');
	//调用公用的loadXMLDoc函数
    loadXMLDoc(urlstr,function()
    {
        if (xmlhttp.readyState==4 && xmlhttp.status==200)//请求处理完成，且状态OK
        {//添加请求完成后的处理功能---------------------------------------
			getStatisticDatas(tagname,xmlhttp.responseText);//解读数据
			//下一步：		
		}else//请求完成后的处理功能结束---------------------------------------
		{
			ShowModal('Error','<div class="alert alert-danger">'+decodeURI(xmlhttp.responseText)+'</div>');
		}
    });
}

function loadCheckSiteExePicUrl(urlstr)//读取巡检照片信息
{	
	//$("#DataMsg").html('<div class="alert alert-warning">正在加载数据信息……</div>');
	//调用公用的loadXMLDoc函数
    loadXMLDoc(urlstr,function()
    {
        if (xmlhttp.readyState==4 && xmlhttp.status==200)//请求处理完成，且状态OK
        {//添加请求完成后的处理功能---------------------------------------
			getCheckSiteExePicUrl(xmlhttp.responseText);//解读数据
			//下一步：		
        }//请求完成后的处理功能结束---------------------------------------
    });
}
//=========AJAX数据接收解析区域====================================================
function getCheckSiteExePicUrl(ajaxdata){
	var imgmsg = eval("("+ajaxdata+")");
	var htmlstr = '<img src="'+imgmsg.FileName+'" class="rounded" alt="巡检照片" style="height:400px;width:768px;display: ;border: 1px solid #cecece;">';
	ShowModal("Inspection photos",htmlstr);
}
//解析配置信息
function getPatrolCfg(ajaxdata){
	var nodes = eval("("+ajaxdata+")"); 
	NODES_CFG = nodes;
	var htmlstr='<div class="col-12 border-bottom"><h5>Check item information</h5></div><table class="table table-striped table-hover table-sm"><thead class="thead-light">';
	htmlstr+='<tr><th>Number</th><th>Inspection items</th><th>Inspection method</th><th>Inspection contents</th><th>Company</th><th>Min</th><th>Max</th><th>Lower limit</th><th>lower limit</th><th>Upper limit</th><th>Uppest limit</th><th>Trend</th><th>Statistics</th></tr></thead><tbody>';
	for(var i=0;i<nodes.length;i++){
		var data=nodes[i];

		//对于有策略上限的检查项(数字类型)，显示趋势按钮，否则禁用
		var trenthbtn=data.MeasMaxvalue==0?('<button type="button" class="btn btn-outline-secondary btn-sm" onclick="onShowTrenth(this.name,this.id)" id="'+data.Id+'" disabled="disabled">Trend</button>'):('<button type="button" class="btn btn-outline-primary btn-sm" onclick="onShowTrenth(this.name,this.id)" id="'+data.Id+'" name="'+data.CheckName+'">Trend</button>');
		//对于有策略上限的检查项(数字类型)，显示分析按钮，否则禁用
		var analysebtn=data.MeasMaxvalue==0?('<button type="button" class="btn btn-outline-secondary btn-sm" onclick="onShowAnalyse(this.name,this.id)" id="'+data.Id+'" disabled="disabled">Statistics</button>'):('<button type="button" class="btn btn-outline-primary btn-sm" onclick="onShowAnalyse(this.name,this.id)" id="'+data.Id+'" name="'+data.CheckName+'">Statistics</button>');

		htmlstr+='<tr><td>'+(i+1)+'</td><td>'+data.CheckName+'</td><td>'+(data.Variable==null ? "":data.Variable.Type.NameCn)+'</td><td>'+(data.Variable==null ? "":data.Variable.CheckContent)+'</td><td>'+(data.Unit==null ? "":data.Unit.UnitName)+'</td><td>'+data.MeasMinvalue+'</td><td>'+data.MeasMaxvalue+'</td><td>'+data.LimitLl+'</td><td>'+data.LimitL+'</td><td>'+data.LimitH+'</td><td>'+data.LimitHh+'</td><td>'+trenthbtn+'</td><td>'+analysebtn+'</td></tr>';
	}

	htmlstr+='</tbody></table>';
	$("#CfgMsg").html(htmlstr);
}

//解析配置信息
function getCheckSiteMsg(ajaxdata){
	var nodes = eval("("+ajaxdata+")"); 
	SITE_LIST_UNDER_NODE = nodes;
	var siteCnt=0;//站点数量
	var lineCnt=0;//线路数量
	var deptCnt=0;//部门数量
	var siteKv=new Array();//站点数量
	var lineKv=new Array();//线路数量
	var deptKv=new Array();//部门数量
	var sites=new Array();
	var doneRates=new Array();

	var htmlstr='<div class="col-12 border-bottom"><h5>Patrol site information</h5></div><table class="table table-striped table-hover table-sm"><thead class="thead-light">';
	htmlstr+='<tr><th>Number</th><th>Checkpoint ID</th><th>Checkpoint name</th><th>Site ID</th><th>Site name</th><th>Inspection line</th><th>Department</th><th>Planned inspection items(项)</th><th>Actual inspection items</th><th>Completion rate(%)</th></tr></thead><tbody>';
	if(nodes!=null){
		for(var i=0;i<nodes.length;i++){
			var data=nodes[i];

			if ((deptKv.hasOwnProperty(data.DeptName+data.DeptId))==false){//如果不存在
				deptKv[data.DeptName+data.DeptId]={name:data.DeptName,children:[]};//赋值
				deptCnt++;
			}
			if (lineKv.hasOwnProperty(data.LineName+data.LineId)==false){//如果不存在
				lineKv[data.LineName+data.LineId]=deptKv[data.DeptName+data.DeptId].children.length;//赋值
				deptKv[data.DeptName+data.DeptId].children[data.LineName+data.LineId]={name:data.LineName,children:[]};
				lineCnt++;
			}else{
				lineKv[data.LineName+data.LineId]++;
			}
			if ((siteKv.hasOwnProperty(data.SiteName+data.SiteId))==false){//如果不存在
				siteKv[data.SiteName+data.SiteId]=deptKv[data.DeptName+data.DeptId].children.length;//赋值
				deptKv[data.DeptName+data.DeptId].children[data.LineName+data.LineId].children[data.SiteName+data.SiteId]={name:data.SiteName,children:[]};
				siteCnt++;
			}else{
				siteKv[data.SiteName+data.SiteId]++;
			}
			var doneRate=data.TotalItemCnt==0?0.00:DataToFixed(data.CheckItemCnt/data.TotalItemCnt*100,"float",2);//完成率
			
			deptKv[data.DeptName+data.DeptId].children[data.LineName+data.LineId].children[data.SiteName+data.SiteId].children.push({name:data.EquipName,value:doneRate});

			htmlstr+='<tr><td>'+(i+1)+'</td><td>'+data.SiteEquipId+'</td><td>'+data.EquipName+'</td><td>'+data.SiteId+'</td><td>'+data.SiteName+'</td><td>'+data.LineName+'</td><td>'+data.DeptName+'</td><td>'+data.TotalItemCnt+'</td><td>'+data.CheckItemCnt+'</td><td>'+doneRate+'</td></tr>';

			sites.push(data.EquipName);
			doneRates.push(doneRate);
		}
	}
	htmlstr+='</tbody></table>';
	$("#CfgMsg").html(htmlstr);

	htmlstr='<div class="col-12 border-bottom"><h5>Checkpoint statistics</h5></div>';
	htmlstr+='<div class="col form-inline"><strong>Number of checkpoints:</strong>'+(nodes==null?0:nodes.length)+'</div>';
	htmlstr+='<div class="col form-inline"><strong>Number of sites:</strong>'+siteCnt+'</div>';
	htmlstr+='<div class="col form-inline"><strong>Number of inspection lines:</strong>'+lineCnt+'</div>';
	htmlstr+='<div class="col form-inline"><strong>Number of management departments:</strong>'+deptCnt+'th</div>';
	htmlstr+='<button type="button" class="btn btn-outline-primary btn-sm" onclick="onShowRelation()">View the atlas</button>'
	$("#CheckSiteMsg").html(htmlstr);

	refreshPageEcharts(sites,doneRates);

	if(deptCnt>1){//部门数大于1
		CHECK_LEVEL_RELATION={name:SELECT_NODE_NAME,children:[]}
		for(var key in deptKv){
			CHECK_LEVEL_RELATION.children.push(keyValueToArray(deptKv[key]));
		}
	}else{
		for(var key in deptKv){
			CHECK_LEVEL_RELATION=keyValueToArray(deptKv[key]);
		}
	}
	//console.log(CHECK_LEVEL_RELATION);
}

function keyValueToArray(kv){
	var rt={name:'',children:[],value:0};
	rt.name = kv.name;
	rt.value = kv.value;
	var i=0;
	for(var key in kv.children){
		rt.children[i]=keyValueToArray(kv.children[key]);
		i++;
	}
	return rt;
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
//解析标签点数据
function getLeafsDatas(datas){
	var htmlstr='<div class="col-12 border-bottom"><h5>Inspection results</h5></div><table class="table table-striped table-hover table-sm"><thead class="thead-light">';
	if(SELECT_LEAF==true){//选择的是巡检标签节点
		var checkcnt=0;//已经检查项数量
		var yichangcnt=0;//异常项数量
		htmlstr+='<tr><th>Number</th><th>Inspection items</th><th>Start time</th><th>End time</th><th>Actual inspection time</th><th>Value</th><th>Company</th><th>Check status</th><th>Value status</th><th>Abnormal</th><th>Remark</th><th>Inspector</th><th>Picture</th></tr></thead><tbody>';
		for(var i=0;i<datas.length;i++){
			data=datas[i];
			checkcnt += data.CheckStatus;//检查项累加
			yichangcnt += data.Abnormal;//异常项累加
			if(data.CheckStatus==0){
				if (PERSON_GROUP.hasOwnProperty("Not checked")){//如果存在
					PERSON_GROUP["Not checked"]+=1;//加一
				}else{//如果不存在
					PERSON_GROUP["Not checked"]=1;//赋值
				}
			}else{//已检
				if (PERSON_GROUP.hasOwnProperty(data.ActualExecutorName)){//如果存在
					PERSON_GROUP[data.ActualExecutorName]+=1;//加一
				}else{//如果不存在
					PERSON_GROUP[data.ActualExecutorName]=1;//赋值
				}
				if(data.MeasMaxvalue!=0){
					if(CHECK_DATA.hasOwnProperty(data.CheckItem.Id)==false){//如果不存在
						CHECK_DATA[data.CheckItem.Id]=[];//创建
					}
					var item={time:data.CheckTime,value:data.CheckResult}
					CHECK_DATA[data.CheckItem.Id].push(item);//插入值
				}
			}

			var checkstatus=data.CheckStatus==1?'<span class="badge badge-success">Inspected</span>':'<span class="badge badge-warning">Not checked</span>';
			var abnormal=data.Abnormal==1?'<span class="badge badge-success">No</span>':'<span class="badge badge-danger">Yes</span>';

			var alarmtype='';
			switch(data.AlarmType){
			case 1:
				alarmtype='<span class="badge badge-danger">Lower than lower limit</span>';
				break;
			case 2:
				alarmtype='<span class="badge badge-warning">Below lower limit</span>';
				break;
			case 3:
				alarmtype='<span class="badge badge-success">Normal</span>';
				break;
			case 4:
				alarmtype='<span class="badge badge-warning">Above the upper limit</span>';
				break;
			case 5:
				alarmtype='<span class="badge badge-danger">Above uppest limit</span>';
				break;
			case 6:
				alarmtype='<span class="badge badge-dark">Over range</span>';
				break;
			default:
				break;
			}
			var picbtn=data.CheckStatus==1?'<button type="button" class="btn btn-outline-primary btn-sm" onclick="onShowPic(this.id)" id="'+data.CheckSiteExe.Id+'">显示</button>':'<button type="button" class="btn btn-outline-secondary btn-sm" onclick="onShowPic(this.id)" id="'+data.CheckSiteExe.Id+'" disabled="disabled">显示</button>';

			if(ONLY_SHOW_ABNORMAL==0 || (ONLY_SHOW_ABNORMAL>0 && data.Abnormal==0)){
				htmlstr+='<tr><td>'+(i+1)+'</td><td>'+data.CheckItemName+'</td><td>'+data.CheckSiteExe.CheckPlanExe.AllStartTime+'</td><td>'+data.CheckSiteExe.CheckPlanExe.AllEndTime+'</td><td>'+data.CheckTime+'</td><td>'+data.CheckResult+'</td><td>'+data.UnitName+'</td><td>'+checkstatus+'</td><td>'+alarmtype+'</td><td>'+abnormal+'</td><td>'+data.Remark+'</td><td>'+data.ActualExecutorName+'</td><td>'+picbtn+'</td></tr>';
			}
		}
	}
	htmlstr+='</tbody></table>';
	$("#DataMsg").html(htmlstr);

	var tjhtml='<div class="col-12 border-bottom"><h5>Statistics of inspection items</h5></div>';
	tjhtml+='<div class="col form-inline"><strong>Inspected items:</strong>'+checkcnt+'</div>';
	tjhtml+='<div class="col form-inline"><strong>Undetected items:</strong>'+(i-checkcnt)+'</div>';
	tjhtml+='<div class="col form-inline"><strong>Inspection rate:</strong>'+DataToFixed((i>0?checkcnt/i*100:0),"float",2)+'%</div>';
	tjhtml+='<div class="col form-inline"><strong>Exception:</strong>'+(i-yichangcnt)+'</div>';
	tjhtml+='<div class="col form-inline"><strong>Abnormal rate:</strong>'+DataToFixed((i>0?(i-yichangcnt)/i*100:0),"float",2)+'%</div>';
	tjhtml+='<div class="col form-inline"><strong>Show exceptions only:</strong><input id="checkbox" type="checkbox" class="form-check-input" '+(ONLY_SHOW_ABNORMAL==1?"checked":"")+' onclick="onlyshowabnomal()"/></div>';
	tjhtml+='<div class="col form-inline"><button type="button" class="btn btn-outline-primary btn-sm" onclick="onShowEcharts()">Statistics by person</button></div>';
	$("#TongJi").show();
	$("#TongJi").html(tjhtml);
}
//解析节点数据
function getNodesDatas(datas){
	var htmlstr='<div class="col-12 border-bottom"><h5>Inspection results</h5></div><table class="table table-striped table-hover table-sm"><thead class="thead-light">';
	if(SELECT_LEAF==true){//选择的是巡检标签节点
		htmlstr+='<tr><th>Number</th><th>Inspection items</th><th>Check item check status</th><th>Start time</th><th>End time</th><th>Actual inspection time</th><th>Value</th><th>Company</th><th>Abnormal</th><th>Remark</th><th>Inspector</th></tr></thead><tbody>';
		for(var i=0;i<datas.length;i++){
			data=datas[i];
			htmlstr+='<tr><td>'+(i+1)+'</td><td>'+data.CheckItemName+'</td><td>'+data.CheckStatus+'</td><td>'+data.StartTime+'</td><td>'+data.EndTime+'</td><td>'+data.CheckTime+'</td><td>'+data.CheckResult+'</td><td>'+data.UnitName+'</td><td>'+data.Abnormal+'</td><td>'+data.Remark+'</td><td>'+data.ActualExecutorName+'</td></tr>';
		}
	}
	htmlstr+='</tbody></table>';
	$("#DataMsg").html(htmlstr);
}

//解析数据信息
function getPatrolDatas(ajaxdata){
	var datas = eval("("+ajaxdata+")"); 
	DATAS = datas;
	if(SELECT_LEAF==true){//选择的是巡检标签节点
		getLeafsDatas(datas);
	}else{//当前选择的是工厂层级节点
		getNodesDatas(datas);
	}
}

//解析数据信息
function getStatisticDatas(tagname,ajaxdata){
	var datas = eval("("+ajaxdata+")"); 
/*
{
    "Min": 1348,
    "Max": 3371,
    "Range": 2023,
    "Total": 1573647238.5,
    "Sum": 45549,
    "Mean": 1897.875,
    "PowerAvg": 1947.575984866393,
    "Diff": -889,
    "PlusDiff": 2569,
    "Duration": 808003,
    "PointCnt": 24,
    "RisingCnt": 9,
    "FallingCnt": 13,
    "LtzCnt": 0,
    "GtzCnt": 24,
    "EzCnt": 0,
    "OutliersCnt": 2,
    "Lower": 472.5,
    "Q1": 1488.75,
    "Q3": 2166.25,
    "Qd": 677.5,
    "Upper": 3182.5,
    "BeginTime": "2020-02-17 07:38:50",
    "EndTime": "2020-02-26 16:05:33",
    "SD": 608.3934659206984,
    "STDDEV": 621.4786933631336,
    "SE": 25.349727746695766,
    "Ske": 1.5170928348103614,
    "Kur": 0.9617627479890722,
    "Mode": 1376.9,
    "Median": 1636,
    "GroupDist": 404.6,
    "DataGroup": {
        "1348": 66.66666666666666,
        "1752.6": 8.333333333333332,
        "2157.2": 4.166666666666666,
        "2561.8": 8.333333333333332,
        "2966.4": 8.333333333333332,
        "3371": 4.166666666666666
    },
    "Increment": {
        "2020-02-17 16:11:08": -1093,
        "2020-02-18 07:20:58": 1153,
        "2020-02-26 16:05:33": -1678
    },
    "RawData": [
        {
            "Time": "2020-02-17 07:38:50",
            "Value": 2569
        },
        {
            "Time": "2020-02-17 16:11:08",
            "Value": 1476
        }
    ]
}
*/
	var items=["Min","Max","Range","=Standard deviation","Avg","Mode","Median","Lower quartile(Q1)","Upper quartile(Q3)","Quartile difference(IQR)","Lower bound of outlier","Upper limit of outlier","Outlier points","Data points"];
	var stat_datas=new Array();
	for(var i=0;i<datas.length;i++){
		dtarr = datas[i];
		var suma={
			Min:DataToFixed(dtarr.Min,'float',2),            //最小值(基本)
			Max:DataToFixed(dtarr.Max,'float',2),            //最大值(基本)
			Range:DataToFixed(dtarr.Range,'float',2),        //数据范围(Max-Min)(基本)
			Total:DataToFixed(dtarr.Total,'float',2),          //表示统计时间段内的累计值，结果的单位为标签点的工程单位(面积,值*时间(s))(基本)
			Sum:DataToFixed(dtarr.Sum,'float',2),            //统计时间段内的算术累积值(值相加)(基本)
			Mean:DataToFixed(dtarr.Mean,'float',2),          //统计时间段内的算术平均值(Mean = Sum/PointCnt)(基本)
			PowerAvg:DataToFixed(dtarr.PowerAvg,'float',2),       //统计时间段内的加权平均值,对BOOL量而言是ON率（Total/Duration）(基本)
			Diff:DataToFixed(dtarr.Diff,'float',2),          //统计时间段内的差值(最后一个值减去第一个值)(基本)
			PlusDiff:DataToFixed(dtarr.PlusDiff,'float',2),       //正差值,用于累计值求差,可以削除清零对值的影响(统计周期内只可以有一次清零动作)
			Duration:dtarr.Duration,     //统计时间段内的秒数(EndTime - BeginTime)(基本)
			PointCnt:dtarr.PointCnt,     //统计时间段内的数据点数(基本)
			RisingCnt:dtarr.RisingCnt,   //统计时间段内数据上升的次数(基本)
			FallingCnt:dtarr.FallingCnt, //统计时间段内数据下降的次数(基本)
			LtzCnt:dtarr.LtzCnt,         //小于0的次数
			GtzCnt:dtarr.GtzCnt,         //大于0的次数
			EzCnt:dtarr.EzCnt,           //等于0的次数
			BeginTime:dtarr.BeginTime,   //开始时间(基本)
			EndTime:dtarr.EndTime,       //结束时间(基本)
			SD:DataToFixed(dtarr.SD,'float',2),             //总体标准差(高级)
			STDDEV:DataToFixed(dtarr.STDDEV,'float',2),     //样本标准差(高级)
			SE:DataToFixed(dtarr.SE,'float',2),             //标准误差(SE = STDDEV / PointCnt)(高级)
			Ske:DataToFixed(dtarr.Ske,'float',2),           //偏度(高级)
			Kur:DataToFixed(dtarr.Kur,'float',2),           //峰度(高级)
			Mode:DataToFixed(dtarr.Mode,'float',2),         //众数(高级)
			Median:DataToFixed(dtarr.Median,'float',2),     //中位数(高级)
			Q1:DataToFixed(dtarr.Q1,'float',2),     //下四分位(高级)
			Q3:DataToFixed(dtarr.Q3,'float',2),     //上四分位(高级)
			Qd:DataToFixed(dtarr.Qd,'float',2),     //四分位差(高级)
			Lower:DataToFixed(dtarr.Lower,'float',2),     //下限
			Upper:DataToFixed(dtarr.Upper,'float',2),     //上限
			OutliersCnt:dtarr.OutliersCnt,     //离群点数
			GroupDist:dtarr.GroupDist       //组距GroupDistance(高级),DataGroup中两组数之间的距离
		};
		stat_datas[i]=[suma.Min,suma.Max,suma.Range,suma.SD,suma.Mean,suma.Mode,suma.Median,suma.Q1,suma.Q3,suma.Qd,suma.Lower,suma.Upper,suma.OutliersCnt,suma.PointCnt];
	}
	var htmlstr='<table class="table table-bordered table-hover table-sm"><thead class="thead-light">';
	if(stat_datas.length>2){
		htmlstr+='<tr><th rowspan="2">Statistical items</th><th colspan="2">Statistics of raw data</th><th colspan="2">Statistical data after outliers are replaced by median</th></tr>'
		htmlstr+='<tr><th>Raw data</th><th>Incremental data</th><th>Original data</th><th>Incremental data</th></tr></thead><tbody>'
	}else{
		htmlstr+='<tr><th>Statistical items</th><th>Raw data</th><th>Incremental data</th></tr></thead><tbody>'
	}
	for(let i in items){
		htmlstr+='<tr>';
		for(var j=-1;j<stat_datas.length;j++){
			if(j<0){
				htmlstr+='<th>'+items[i]+'</th>';
			}else{
				htmlstr+='<td>'+stat_datas[j][i]+'</td>';
			}
		}
		htmlstr+='</tr>';
	}
	htmlstr+='</tbody></table>';
	ShowModal(tagname+"Statistical data",htmlstr);
}

//=========筛选数据====================================================
//按物资名称筛选
function onFillterByName(name){
	var patroltypes=new Array();
	var htmlstr='<div class="col-12 border-bottom"><h5>Inspection results</h5></div><table class="table table-striped table-hover table-sm"><thead class="thead-light">';
	htmlstr+='<tr><th>Number</th><th>Material name</th><th>Specification and model</th><th>Workshop</th><th>Level</th><th>Start time</th><th>End time</th><th>Interval(h)</th><th>Number</th><th>Company</th><th>Remark</th></tr></thead><tbody>';
	var dsum=0;//和
	var unit;//单位
	var j=0;
	for(var i=0;i<DATAS.length;i++){
		data=DATAS[i];
		var patrolname=name;
		if(name=='0'){
			patrolname = data.PatrolConfigInfo.Patrol.PatrolName;
		}
		if(data.PatrolConfigInfo.Patrol.PatrolName == patrolname){
			patroltypes[data.PatrolConfigInfo.Patrol.Specifications]=data.PatrolConsumeAmount;//按名称求和
			htmlstr+='<tr><td>'+(++j)+'</td><td>'+data.PatrolConfigInfo.Patrol.PatrolName+'</td><td>'+data.PatrolConfigInfo.Patrol.Specifications+'</td><td>'+data.PatrolConfigInfo.Workshop.WorkshopName+'</td><td>'+LEVEL_NAMES[data.PatrolConfigInfo.TreeLevelCode]+'</td><td>'+data.UseStartTime+'</td><td>'+data.UseEndTime+'</td><td>'+timeDiff(data.UseStartTime,data.UseEndTime)+'</td><td>'+data.PatrolConsumeAmount+'</td><td>'+data.PatrolConfigInfo.Patrol.PatrolConsumeUnit+'</td><td>'+data.Remark+'</td></tr>';
			dsum+=data.PatrolConsumeAmount;//累加
			unit=data.PatrolConfigInfo.Patrol.PatrolConsumeUnit;//单位
		}
	}
	if(name !='0'){//不是全部的时候才显示累积,更新型号筛选
		htmlstr+='<tr><th colspan="8">Total</th><th>'+DataToFixed(dsum,"float",2)+'</th><th>'+unit+'</th><th></th><tr>';

		$("#TypeFillter").html("");
		var fillter='<label class="font-weight-bold" for="PatrolType">Specification and model:</label><select class="form-control" id="PatrolType" onchange="onFillterByType(this.options[this.options.selectedIndex].value)"><option value="0">All</option>';
		for (let k in patroltypes) {//型号过滤选择框
			fillter+='<option value="'+k+'">'+k+'</option>';
		}
		fillter+='</select>';
		$("#PatrolSum").html(DataToFixed(dsum,"float",2)+unit);//显示累积
		$("#PatrolSumDiv").show();
		$("#TypeFillter").html(fillter);
	}else{
		$("#TypeFillter").html("");
		var fillter='<label class="font-weight-bold" for="PatrolType">Specification and model:</label><select class="form-control" id="PatrolType" onchange="onFillterByType(this.options[this.options.selectedIndex].value)" disabled="disabled"><option value="0">All</option>';
		fillter+='</select>';
		$("#TypeFillter").html(fillter);
		$("#PatrolSumDiv").hide();
	}
	htmlstr+='</tbody></table>';
	$("#DataMsg").html(htmlstr);
}

//按物资名称筛选
function onFillterByType(type){
	var htmlstr='<table class="table table-striped table-hover table-sm"><tr colspan="5"><h5>Material consumption data</h5></tr><thead class="thead-light">';
	htmlstr+='<tr><th>Number</th><th>Material name</th><th>Specification and model</th><th>Workshop</th><th>Level</th><th>Start time</th><th>End time </th><th>Interval(h)</th><th>umber</th><th>Unit</th><th>Remark</th></tr></thead><tbody>';
	var patrolname=$("#PatrolName").val();
	var dsum=0;//和
	var unit;//单位
	var j=0;
	for(var i=0;i<DATAS.length;i++){
		data=DATAS[i];
		var patroltype=type;
		if(type=='0'){
			patroltype = data.PatrolConfigInfo.Patrol.Specifications;
		}
		if(data.PatrolConfigInfo.Patrol.PatrolName == patrolname && data.PatrolConfigInfo.Patrol.Specifications==patroltype){
			htmlstr+='<tr><td>'+(++j)+'</td><td>'+data.PatrolConfigInfo.Patrol.PatrolName+'</td><td>'+data.PatrolConfigInfo.Patrol.Specifications+'</td><td>'+data.PatrolConfigInfo.Workshop.WorkshopName+'</td><td>'+LEVEL_NAMES[data.PatrolConfigInfo.TreeLevelCode]+'</td><td>'+data.UseStartTime+'</td><td>'+data.UseEndTime+'</td><td>'+timeDiff(data.UseStartTime,data.UseEndTime)+'</td><td>'+data.PatrolConsumeAmount+'</td><td>'+data.PatrolConfigInfo.Patrol.PatrolConsumeUnit+'</td><td>'+data.Remark+'</td></tr>';
			dsum+=data.PatrolConsumeAmount;//累加
			unit=data.PatrolConfigInfo.Patrol.PatrolConsumeUnit;//单位
		}
	}
	htmlstr+='<tr><th colspan="8">合计</th><th>'+DataToFixed(dsum,"float",2)+'</th><th>'+unit+'</th><th></th><tr>';
	htmlstr+='</tbody></table>';
	$("#PatrolSum").html(DataToFixed(dsum,"float",2) + unit);
	$("#DataMsg").html(htmlstr);
}
$("#modal-btn").text('close');

$(document).ready(function() {
	$("#exit").after('Exit');
	$("#ExpandTreeNode").html('Open');
	$("#CollapseTreeNode").html('Fold');
	$("#HideTreeNode").html('Hide');
	$("#SearchTreeNode").attr("placeholder",'Search');

	$('#ExpandTreeNode').attr("title","Expand all nodes when no nodes are selected, and expand selected nodes when nodes are selected");
	$('#CollapseTreeNode').attr("title","Collapse all nodes when no nodes are selected, and collapse selected nodes when nodes are selected");
	$("#BeginTimes").text('Begin Time:');
	$("#EndTimes").text('End Time:');
	$("#TimeRanges").text('Time Range:');

	$("#custom").text("custom");
	$("#8h").text("8h");
	$("#12h").text("12h");
	$("#24h").text("24h");
	$("#7d").text("7d");
	$("#op_this_shift").text("class");
	$("#op_this_day").text("day");
	$("#op_this_month").text("month");

	$('#Last').attr("title","Previous period!");
	$('#Next').attr("title","The latter period!");

	$('#OperationMsg').html('<strong>handle：</strong> Please select a sample node in the left structure tree to display the data!');

});

