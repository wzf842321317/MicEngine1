//=========变量定义区域===========================================================

var SAMPLE_NODE={//当前选中的样本节点
	Id:0,//样本ID
	Pid:0,//所在作业ID
	Name:"",//样本名称
	PoolId:0,//所属样本池ID
	TreeLevel:"",//层级码
	SamplingSite:"",//取样地点
	FuncType:0,//样本类型码
	FuncName:"",//样本类型名
	IsRegular:0,//是否常规样/临时样
	BaseTime:"",//所在车间基准时间
	ShiftHour:0//所在车间每班工作时间
};
var HAS_SELECT_SAMPLE=false;//已经选择了样本节点
var SAMPLE_LAB_TAGS=[];//化验标签信息
var SAMPLE_LAB_RESULT=[];//化验结果信息

//==================动作响应区域==================================================
//响应鼠标单击
function zTreeOnClick(event, treeId, treeNode) {
	if (treeNode.nodetype == 9999){
		SAMPLE_NODE.Id = treeNode.id;
		SAMPLE_NODE.Pid = treeNode.pId;
		SAMPLE_NODE.Name = treeNode.name;
		SAMPLE_NODE.poolId = treeNode.itemid;
		SAMPLE_NODE.TreeLevel = treeNode.treelevel;
		SAMPLE_NODE.SamplingSite = treeNode.sampsite;
		SAMPLE_NODE.FuncType = treeNode.functype;
		SAMPLE_NODE.FuncName = treeNode.funcname;
		SAMPLE_NODE.IsRegular = treeNode.isregular;
		SAMPLE_NODE.BaseTime = treeNode.basetime;
		SAMPLE_NODE.ShiftHour = treeNode.shifthour;
		
		HAS_SELECT_SAMPLE = true;//设定已选样本模板标记
		$("#SampleLabResult").html("");//清空结果
		$("#SampleLabTags").html("");//清空结果
		$("#OperationMsg").hide();//隐藏操作提示
		requestSampleLabTags(treeNode.id);//请求样本化验标签信息
		showSampleMsg(SAMPLE_NODE);//显示样本模板基本信息
	}
	
	timeRangeSelectorSet();//时间范围设置
}

//响应鼠标双击
function zTreeOnDbClick(event, treeId, treeNode){
	
}
//页面初始化工作
function pageInit(){
	timeRangeCheck();
	timeRangeSelectorSet();
};
//显示所选样本模板的基本信息
function showSampleMsg(samp){
	var htmlstr='<div class="col-12 border-bottom"><h5>样本基本信息</h5></div>';
	htmlstr+='<div class="col-sm-3 form-inline"><strong>样本名称:</strong>'+samp.Name+'</div>';
	htmlstr+='<div class="col-sm-3 form-inline"><strong>取样地点:</strong>'+samp.SamplingSite+'</div>';
	htmlstr+='<div class="col-sm-3 form-inline"><strong>样本类型:</strong>'+samp.FuncName+'</div>';
	if(samp.IsRegular==1){
		reg='<span class="badge badge-success">是</span>';
	}else{
		reg='<span class="badge badge-danger">否</span>';
	}
	htmlstr+='<div class="col-sm-3 form-inline"><strong>规律样:</strong>'+reg+'</div>';
	$("#SampleMsg").html(htmlstr);
}
//时间范围选择框设置
function timeRangeSelectorSet(){
	if(HAS_SELECT_SAMPLE == true && SAMPLE_NODE.BaseTime.length > 10){//已经选择了样本模板,取消禁用选择本班、今日、本月、上月
		if (SAMPLE_NODE.ShiftHour > 0){
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
	if(HAS_SELECT_SAMPLE==true){//如果已经选择了样本
		requestSampleLabResult(SAMPLE_NODE.Id);//读取化验结果请求
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
		timediff = SAMPLE_NODE.ShiftHour * 3600 * 1000;
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
		getBeginTimeOfMonth(SAMPLE_NODE.BaseTime,begintime);
		break;
	default:
		bgstemp.setTime(bgstemp.getTime() - timediff);
		edstemp.setTime(edstemp.getTime() - timediff);
		$("#BeginTime").val(DateFormat("YYYY-mm-ddTHH:MM",bgstemp));
		$("#EndTime").val(DateFormat("YYYY-mm-ddTHH:MM",edstemp));
		break;
	}
	
	if(HAS_SELECT_SAMPLE==true){//如果已经选择了样本
		requestSampleLabResult(SAMPLE_NODE.Id);//读取化验结果请求
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
		timediff = SAMPLE_NODE.ShiftHour * 3600 * 1000;
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

	if(HAS_SELECT_SAMPLE==true){//如果已经选择了样本
		requestSampleLabResult(SAMPLE_NODE.Id);//读取化验结果请求
	}
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
		getBeginTimeInDay(SAMPLE_NODE.BaseTime,DateFormat("YYYY-mm-ddTHH:MM",nowTime),SAMPLE_NODE.ShiftHour);
		break;
	case '2'://今日
		getBeginTimeInDay(SAMPLE_NODE.BaseTime,DateFormat("YYYY-mm-ddTHH:MM",nowTime),24);
		break;
	case '3'://本月
		getBeginTimeOfMonth(SAMPLE_NODE.BaseTime,DateFormat("YYYY-mm-ddTHH:MM",nowTime))
		break;
	default:
		bgstemp.setTime(edstemp.getTime() - tdiff*60*1000);
		$("#BeginTime").val(DateFormat("YYYY-mm-ddTHH:MM",bgstemp));
		break;
	}	

	if(HAS_SELECT_SAMPLE==true){//如果已经选择了样本
		requestSampleLabResult(SAMPLE_NODE.Id);//读取化验结果请求
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

//显示样本照片
function onShowPic(picurl){
	var htmlstr = '<img src="'+picurl+'" class="rounded" alt="样本照片" style="height:400px;width:768px;display: ;border: 1px solid #cecece;">';
	ShowModal("样本照片",htmlstr);
}

//本班或者今日的开始时间
function getBeginTimeInDay(basetime,lasttime,period_h){
	var baseT= new Date(basetime.replace(/T/," "));//基准时间
	var nowT= new Date(lasttime.replace(/T/," "));//当前时间
	var period = period_h * 3600*1000;//周期
	var endT=nowT;
	var bgT=nowT;
	endT.setTime(Math.round(((nowT.getTime() + period) - baseT.getTime())/period) * period + baseT.getTime());
	bgT.setTime(endT.getTime() - period);
	beginTime = DateFormat("YYYY-mm-ddTHH:MM",bgT);
	$("#BeginTime").val(beginTime);
	var nowTime= new Date(lasttime.replace(/T/," "));//当前时间
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
//=========AJAX请求区域=======================================================
//读取样本化验变量信息请求
function requestSampleLabTags(sampleid){
	var urlstr = "api/samplelabtag?sampleid="+sampleid/100000;
	loadSampleLabTags(urlstr,sampleid);
}
//读取样本化验结果信息请求
function requestSampleLabResult(sampleid){
	var urlstr = "api/samplelabresult?sampleid="+sampleid/100000+"&begintime="+InputDateTimeToString("YYYY-mm-dd HH:MM:SS",$("#BeginTime").val())+"&endtime="+InputDateTimeToString("YYYY-mm-dd HH:MM:SS",$("#EndTime").val());
	loadSampleLabResult(urlstr);
}
//=========AJAX函数定义区域=======================================================
function loadSampleLabTags(urlstr,sampleid)//读取样本化验变量信息请求
{	
	$("#SampleLabTags").html('<div class="alert alert-warning">正在加载配置指标……</div>');
	//调用公用的loadXMLDoc函数
    loadXMLDoc(urlstr,function()
    {
        if (xmlhttp.readyState==4 && xmlhttp.status==200)//请求处理完成，且状态OK
        {//添加请求完成后的处理功能---------------------------------------
			getSampleLabTags(xmlhttp.responseText);//解读数据
			//下一步：加载历史数据
			requestSampleLabResult(sampleid);			
        }//请求完成后的处理功能结束---------------------------------------
    });
}

function loadSampleLabResult(urlstr)//读取样本化验结果信息请求
{	
	$("#SampleLabResult").html('<div class="alert alert-warning">正在加载化验结果……</div>');
	//调用公用的loadXMLDoc函数
    loadXMLDoc(urlstr,function()
    {
        if (xmlhttp.readyState==4 && xmlhttp.status==200)//请求处理完成，且状态OK
        {//添加请求完成后的处理功能---------------------------------------
			getSampleLabResult(xmlhttp.responseText);//解读数据
			//下一步：		
        }//请求完成后的处理功能结束---------------------------------------
    });
}

//=========AJAX反馈数据解析=======================================================
//解析样本化验变量信息
function getSampleLabTags(ajaxdata){
	var dtarr = eval("("+ajaxdata+")"); 
	SAMPLE_LAB_TAGS.splice(0,SAMPLE_LAB_TAGS.length);//清空数组
	SAMPLE_LAB_TAGS=dtarr
	var htmlstr='<table class="table table-striped table-hover table-sm"><tr colspan="5"><h5>化验指标信息</h5></tr><thead class="thead-light"><tr><th>序号</th><th>名称</th><th>关键词</th><th>单位</th><th>标签</th></tr></thead><tbody>';
	for(var i=0;i<SAMPLE_LAB_TAGS.length;i++){
		var samp=SAMPLE_LAB_TAGS[i];
		htmlstr+='<tr><td>'+(i+1)+'</td><td>'+samp.LabIndex.NameCn+'</td><td>'+samp.LabIndex.KeyWord+'</td><td>'+samp.Unit.IndexUnit+'</td><td>'+samp.SampleIndexTag+'</td></tr>';
	}
	htmlstr+='</tbody></table>';
	$("#SampleLabTags").html(htmlstr);
}
//解析样本化验结果信息
function getSampleLabResult(ajaxdata){
	var dtarr = eval("("+ajaxdata+")"); 
	SAMPLE_LAB_RESULT.splice(0,SAMPLE_LAB_RESULT.length);//清空数组
	SAMPLE_LAB_RESULT=dtarr;
	
	sampPicPath=$("#PlatPicPath").text();//图片主路径
	
	var htmlstr='<table class="table table-striped table-hover table-sm"><tr colspan="5"><h5>化验结果信息</h5></tr><thead class="thead-light"><tr><th>序号</th><th>取样时间</th><th>送样时间</th><th>样本编号</th><th>样本状态</th><th>指标</th><th>关键词</th><th>结果</th><th>单位</th><th>样本照片</th></tr></thead><tbody>';
	for(var i=0;i<SAMPLE_LAB_RESULT.length;i++){
		var samp=SAMPLE_LAB_RESULT[i];
		
		var labtype;
		switch(samp.SampleToLab.LabType){
		case 0:
			labtype='<span class="badge badge-secondary">未接收</span>';
			break;
		case 1:
			labtype='<span class="badge badge-info">未化验</span>';
			break;
		case 2:
			labtype='<span class="badge badge-warning">未审核</span>';
			break;
		case 3:
			labtype='<span class="badge badge-success">已审核</span>';
			break;
		default:
			break;
		}
		var picbtn='<button type="button" class="btn btn-outline-primary btn-sm" onclick="onShowPic(this.id)" id="'+sampPicPath+samp.SampleToLab.SampleSitePhoto+'">显示</button>';
		
		htmlstr+='<tr><td>'+(i+1)+'</td><td>'+samp.SampleToLab.SamplingTime+'</td><td>'+samp.SampleToLab.SampleToLabTime+'</td><td>'+samp.SampleToLab.SampleCodeNum+'</td><td>'+labtype+'</td><td>'+samp.LabTag.LabIndex.NameCn+'</td><td>'+samp.LabTag.LabIndex.KeyWord+'</td><td>'+samp.SampleIndexTagValue+'</td><td>'+samp.LabTag.Unit.IndexUnit+'</td><td>'+picbtn+'</td></tr>';
	}
	htmlstr+='</tbody></table>';
	$("#SampleLabResult").html(htmlstr);
}
