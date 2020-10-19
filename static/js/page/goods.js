//=========变量定义区域===========================================================
var SELECT_NODE={//当前选中的样本节点
	Id:0,//样本ID
	Pid:0,//所在作业ID
	Name:"",//样本名称
	TreeLevel:"",//层级码
	BaseTime:"",//所在车间基准时间
	ShiftHour:0//所在车间每班工作时间
};
var HAS_SELECT_NODE=false;//已经选择了有效节点
var NODES_CFG;//节点配置信息
var DATAS=[];//物耗数据
var SELECT_NODE_LEVEL;//所选节点的层级码
var SELECT_LEAF;//选择了叶子节点
var LEVEL_NAMES=new Array();//节点层级名称键值对
//==================动作响应区域==================================================
//响应鼠标单击
function zTreeOnClick(event, treeId, treeNode) {
	$("#OperationMsg").hide();//隐藏操作提示
	$("#CfgMsg").html("");//清空数据
	$("#DataMsg").html("");//清空数据
	SELECT_NODE_LEVEL = treeNode.treelevel;
	HAS_SELECT_NODE = true;//已经选择了有效节点

	SELECT_NODE.Id=treeNode.id;
	SELECT_NODE.Pid=treeNode.pId;
	SELECT_NODE.Name=treeNode.name;
	SELECT_NODE.TreeLevel=treeNode.treelevel;
	SELECT_NODE.BaseTime=treeNode.basetime;
	SELECT_NODE.ShiftHour=treeNode.shifthour;
	
	if (treeNode.nodetype == 9999){
		SELECT_LEAF=true;//选择了叶子节点
		requestGoodsCfg(SELECT_NODE.Id);//请求读取物耗配置信息
		$("#OperationMsg").hide();//隐藏操作提示
	}else{
		SELECT_LEAF = false;//没有选择叶子节点
		requestGoodsDatas(SELECT_NODE_LEVEL,0);
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
			requestGoodsDatas(SELECT_NODE.Id/100000,1);
		}else{
			requestGoodsDatas(SELECT_NODE_LEVEL,0);
		}
	}
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

//=========AJAX请求定义区域=======================================================
//读取物耗配置信息请求
function requestGoodsCfg(goodsid){
	var urlstr = "api/goodscfg?goodsid="+goodsid/100000;
	loadGoodsCfg(urlstr,goodsid);
}
//读取物耗数据信息请求
function requestGoodsDatas(idorlevel,isid){
	var urlstr = "api/goodsdatas?leveltag="+idorlevel+"&begintime="+InputDateTimeToString("YYYY-mm-dd HH:MM:SS",$("#BeginTime").val())+"&endtime="+InputDateTimeToString("YYYY-mm-dd HH:MM:SS",$("#EndTime").val())+"&startonly=1&isid="+isid;
	loadGoodsDatas(urlstr);
}
//=========AJAX加载定义区域=======================================================
function loadGoodsCfg(urlstr,goodsid)//读取物耗配置信息
{	
	$("#CfgMsg").html('<div class="alert alert-warning">正在加载配置信息……</div>');
	//调用公用的loadXMLDoc函数
    loadXMLDoc(urlstr,function()
    {
        if (xmlhttp.readyState==4 && xmlhttp.status==200)//请求处理完成，且状态OK
        {//添加请求完成后的处理功能---------------------------------------
			getGoodsCfg(xmlhttp.responseText);//解读数据
			//下一步：
			requestGoodsDatas(goodsid/100000,1);			
        }//请求完成后的处理功能结束---------------------------------------
    });
}
function loadGoodsDatas(urlstr)//读取物耗配置信息
{	
	$("#DataMsg").html('<div class="alert alert-warning">正在加载数据信息……</div>');
	//调用公用的loadXMLDoc函数
    loadXMLDoc(urlstr,function()
    {
        if (xmlhttp.readyState==4 && xmlhttp.status==200)//请求处理完成，且状态OK
        {//添加请求完成后的处理功能---------------------------------------
			getGoodsDatas(xmlhttp.responseText);//解读数据
			//下一步：		
        }//请求完成后的处理功能结束---------------------------------------
    });
}
//=========AJAX数据接收解析区域====================================================
//解析配置信息
function getGoodsCfg(ajaxdata){
	var nodes = eval("("+ajaxdata+")"); 
	var htmlstr='<div class="col-12 border-bottom"><h5>物资基本信息</h5></div>';
	for(var i=0;i<nodes.length;i++){
		var node=nodes[i];
		NODES_CFG = node;
		var danger = '<span class="badge badge-success">否</span>';;
		if(node.Goods.GoodsIsDanger==1){
			danger = '<span class="badge badge-danger">是</span>';
		}
		htmlstr+='<div class="col-sm-4 form-inline"><strong>物品名称:</strong>'+node.Goods.GoodsName+'</div>';
		htmlstr+='<div class="col-sm-4 form-inline"><strong>所属车间:</strong>'+node.Workshop.WorkshopName+'</div>';
		htmlstr+='<div class="col-sm-4 form-inline"><strong>所属层级:</strong>'+LEVEL_NAMES[node.TreeLevelCode]+'</div>';
		htmlstr+='<div class="col-sm-4 form-inline"><strong>管理部门:</strong>'+node.Goods.DeptManage.DepartmentName+'</div>';
		htmlstr+='<div class="col-sm-4 form-inline"><strong>用户编码:</strong>'+node.Goods.UserNameCode+'</div>';
		htmlstr+='<div class="col-sm-4 form-inline"><strong>规格型号:</strong>'+node.Goods.Specifications+'</div>';
		htmlstr+='<div class="col-sm-4 form-inline"><strong>计量单位:</strong>'+node.Goods.GoodsConsumeUnit+'</div>';
		htmlstr+='<div class="col-sm-4 form-inline"><strong>包装方式:</strong>'+node.Goods.PackingMode+'</div>';
		htmlstr+='<div class="col-sm-4 form-inline"><strong>包装单位:</strong>'+node.Goods.PackingUnit+'</div>';
		htmlstr+='<div class="col-sm-4 form-inline"><strong>包装规格:</strong>'+node.Goods.PackingModel+'</div>';
		htmlstr+='<div class="col-sm-4 form-inline"><strong>存放条件:</strong>'+node.Goods.PreserveCondition+'</div>';
		htmlstr+='<div class="col-sm-4 form-inline"><strong>运输条件:</strong>'+node.Goods.TransportCondition+'</div>';
		htmlstr+='<div class="col-sm-4 form-inline"><strong>运输方式:</strong>'+node.Goods.TransportMode+'</div>';
		htmlstr+='<div class="col-sm-4 form-inline"><strong>危 险 品:</strong>'+danger+'</div>';
		htmlstr+='<div class="col-sm-4 form-inline"><strong>出厂编号:</strong>'+node.Goods.FactorySerialNum+'</div>';
		htmlstr+='<div class="col-sm-4 form-inline"><strong>生产厂家:</strong>'+node.Goods.ManufacturerName+'</div>';
		htmlstr+='<div class="col-sm-4 form-inline"><strong>生产批号:</strong>'+node.Goods.ProductionBatchNum+'</div>';
		htmlstr+='<div class="col-sm-4 form-inline"><strong>参考单价:</strong>'+node.Goods.ConsumeUnitPrice+'</div>';
		htmlstr+='<div class="col-sm-4 form-inline"><strong>换算比例:</strong>'+node.Goods.ScalingFactor+'</div>';
		htmlstr+='<div class="col-sm-4 form-inline"><strong>所选时间范围合计用量:</strong><span id="SumData"></span></div>';
	}
	$("#CfgMsg").html(htmlstr);
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
function getGoodsDatas(ajaxdata){
	var datas = eval("("+ajaxdata+")"); 
	DATAS.splice(0,DATAS.length);//清空数组
	DATAS = datas;
	var htmlstr='<table class="table table-striped table-hover table-sm"><tr colspan="5"><h5>物资消耗数据</h5></tr><thead class="thead-light">';
	if(SELECT_LEAF==true){//选择的是物耗标签节点
		htmlstr+='<tr><th>序号</th><th>开始时间</th><th>结束时间</th><th>间隔(h)</th><th>数量</th><th>单位</th><th>备注</th><th>创建人</th><th>创建时间</th><th>编辑人</th><th>编辑时间</th></tr></thead><tbody>';
		var dsum=0;
		var unit=NODES_CFG.Goods.GoodsConsumeUnit;
		for(var i=0;i<datas.length;i++){
			data=datas[i];
			htmlstr+='<tr><td>'+(i+1)+'</td><td>'+data.UseStartTime+'</td><td>'+data.UseEndTime+'</td><td>'+timeDiff(data.UseStartTime,data.UseEndTime)+'</td><td>'+data.GoodsConsumeAmount+'</td><td>'+data.GoodsConfigInfo.Goods.GoodsConsumeUnit+'</td><td>'+data.Remark+'</td><td>'+data.CreateUser.Name+'</td><td>'+data.CreateTime+'</td><td>'+data.UpdateUser.Name+'</td><td>'+data.UpdateTime+'</td></tr>';
			dsum+=data.GoodsConsumeAmount;
		}
		htmlstr+='<tr><th colspan="4">合计</th><th>'+DataToFixed(dsum,"float",2)+'</th><th colspan="6">'+unit+'</th><tr>';
		$("#SumData").text(DataToFixed(dsum,"float",2)+unit);
	}else{//当前选择的是工厂层级节点
		var goodsnames=new Array();
		htmlstr+='<tr><th>序号</th><th>物资名称</th><th>规格型号</th><th>所属车间</th><th>所属层级</th><th>开始时间</th><th>结束时间</th><th>间隔(h)</th><th>数量</th><th>单位</th><th>备注</th></tr></thead><tbody>';
		for(var i=0;i<datas.length;i++){
			data=datas[i];
			goodsnames[data.GoodsConfigInfo.Goods.GoodsName]=data.GoodsConsumeAmount;//按名称求和
			htmlstr+='<tr><td>'+(i+1)+'</td><td>'+data.GoodsConfigInfo.Goods.GoodsName+'</td><td>'+data.GoodsConfigInfo.Goods.Specifications+'</td><td>'+data.GoodsConfigInfo.Workshop.WorkshopName+'</td><td>'+LEVEL_NAMES[data.GoodsConfigInfo.TreeLevelCode]+'</td><td>'+data.UseStartTime+'</td><td>'+data.UseEndTime+'</td><td>'+timeDiff(data.UseStartTime,data.UseEndTime)+'</td><td>'+data.GoodsConsumeAmount+'</td><td>'+data.GoodsConfigInfo.Goods.GoodsConsumeUnit+'</td><td>'+data.Remark+'</td></tr>';
		}
		//按名称、型号过滤
		var fillter='<div class="col-12 border-bottom"><h5>物资筛选</h5></div>';
		fillter+='<div class="col form-inline"><label class="font-weight-bold" for="GoodsName">物资名称:</label><select class="form-control" id="GoodsName" onchange="onFillterByName(this.options[this.options.selectedIndex].value)" ><option value="0">全部</option>';
		for (let k in goodsnames) {//名称过滤选择框
			fillter+='<option value="'+k+'">'+k+'</option>';
		}
		fillter+='</select></div>';
		fillter+='<div class="col form-inline" id="TypeFillter"><label class="font-weight-bold" for="GoodsType">规格型号:</label><select class="form-control" id="GoodsType" onchange="onFillterByType(this.options[this.options.selectedIndex].value)" disabled="disabled"><option value="0">全部</option></select></div>';
		fillter+='<div class="col form-inline" id="GoodsSumDiv" style="display:none"><label class="font-weight-bold" for="GoodsSum">合计:</label><span id="GoodsSum"></span></div>';//显示累积
		$("#CfgMsg").html(fillter);
	}
	htmlstr+='</tbody></table>';
	$("#DataMsg").html(htmlstr);
}

//=========筛选数据====================================================
//按物资名称筛选
function onFillterByName(name){
	var goodstypes=new Array();
	var htmlstr='<table class="table table-striped table-hover table-sm"><tr colspan="5"><h5>物资消耗数据</h5></tr><thead class="thead-light">';
	htmlstr+='<tr><th>序号</th><th>物资名称</th><th>规格型号</th><th>所属车间</th><th>所属层级</th><th>开始时间</th><th>结束时间</th><th>间隔(h)</th><th>数量</th><th>单位</th><th>备注</th></tr></thead><tbody>';
	var dsum=0;//和
	var unit;//单位
	var j=0;
	for(var i=0;i<DATAS.length;i++){
		data=DATAS[i];
		var goodsname=name;
		if(name=='0'){
			goodsname = data.GoodsConfigInfo.Goods.GoodsName;
		}
		if(data.GoodsConfigInfo.Goods.GoodsName == goodsname){
			goodstypes[data.GoodsConfigInfo.Goods.Specifications]=data.GoodsConsumeAmount;//按名称求和
			htmlstr+='<tr><td>'+(++j)+'</td><td>'+data.GoodsConfigInfo.Goods.GoodsName+'</td><td>'+data.GoodsConfigInfo.Goods.Specifications+'</td><td>'+data.GoodsConfigInfo.Workshop.WorkshopName+'</td><td>'+LEVEL_NAMES[data.GoodsConfigInfo.TreeLevelCode]+'</td><td>'+data.UseStartTime+'</td><td>'+data.UseEndTime+'</td><td>'+timeDiff(data.UseStartTime,data.UseEndTime)+'</td><td>'+data.GoodsConsumeAmount+'</td><td>'+data.GoodsConfigInfo.Goods.GoodsConsumeUnit+'</td><td>'+data.Remark+'</td></tr>';
			dsum+=data.GoodsConsumeAmount;//累加
			unit=data.GoodsConfigInfo.Goods.GoodsConsumeUnit;//单位
		}
	}
	if(name !='0'){//不是全部的时候才显示累积,更新型号筛选
		htmlstr+='<tr><th colspan="8">合计</th><th>'+DataToFixed(dsum,"float",2)+'</th><th>'+unit+'</th><th></th><tr>';

		$("#TypeFillter").html("");
		var fillter='<label class="font-weight-bold" for="GoodsType">规格型号:</label><select class="form-control" id="GoodsType" onchange="onFillterByType(this.options[this.options.selectedIndex].value)"><option value="0">全部</option>';
		for (let k in goodstypes) {//型号过滤选择框
			fillter+='<option value="'+k+'">'+k+'</option>';
		}
		fillter+='</select>';
		$("#GoodsSum").html(DataToFixed(dsum,"float",2)+unit);//显示累积
		$("#GoodsSumDiv").show();
		$("#TypeFillter").html(fillter);
	}else{
		$("#TypeFillter").html("");
		var fillter='<label class="font-weight-bold" for="GoodsType">规格型号:</label><select class="form-control" id="GoodsType" onchange="onFillterByType(this.options[this.options.selectedIndex].value)" disabled="disabled"><option value="0">全部</option>';
		fillter+='</select>';
		$("#TypeFillter").html(fillter);
		$("#GoodsSumDiv").hide();
	}
	htmlstr+='</tbody></table>';
	$("#DataMsg").html(htmlstr);
}

//按物资名称筛选
function onFillterByType(type){
	var htmlstr='<table class="table table-striped table-hover table-sm"><tr colspan="5"><h5>物资消耗数据</h5></tr><thead class="thead-light">';
	htmlstr+='<tr><th>序号</th><th>物资名称</th><th>规格型号</th><th>所属车间</th><th>所属层级</th><th>开始时间</th><th>结束时间</th><th>间隔(h)</th><th>数量</th><th>单位</th><th>备注</th></tr></thead><tbody>';
	var goodsname=$("#GoodsName").val();
	var dsum=0;//和
	var unit;//单位
	var j=0;
	for(var i=0;i<DATAS.length;i++){
		data=DATAS[i];
		var goodstype=type;
		if(type=='0'){
			goodstype = data.GoodsConfigInfo.Goods.Specifications;
		}
		if(data.GoodsConfigInfo.Goods.GoodsName == goodsname && data.GoodsConfigInfo.Goods.Specifications==goodstype){
			htmlstr+='<tr><td>'+(++j)+'</td><td>'+data.GoodsConfigInfo.Goods.GoodsName+'</td><td>'+data.GoodsConfigInfo.Goods.Specifications+'</td><td>'+data.GoodsConfigInfo.Workshop.WorkshopName+'</td><td>'+LEVEL_NAMES[data.GoodsConfigInfo.TreeLevelCode]+'</td><td>'+data.UseStartTime+'</td><td>'+data.UseEndTime+'</td><td>'+timeDiff(data.UseStartTime,data.UseEndTime)+'</td><td>'+data.GoodsConsumeAmount+'</td><td>'+data.GoodsConfigInfo.Goods.GoodsConsumeUnit+'</td><td>'+data.Remark+'</td></tr>';
			dsum+=data.GoodsConsumeAmount;//累加
			unit=data.GoodsConfigInfo.Goods.GoodsConsumeUnit;//单位
		}
	}
	htmlstr+='<tr><th colspan="8">合计</th><th>'+DataToFixed(dsum,"float",2)+'</th><th>'+unit+'</th><th></th><tr>';
	htmlstr+='</tbody></table>';
	$("#GoodsSum").html(DataToFixed(dsum,"float",2) + unit);
	$("#DataMsg").html(htmlstr);
}
$(document).ready(function() {
	//Ztree国际化
	$("#exit").after('退出');
	$("#ExpandTreeNode").html('展开');
	$("#CollapseTreeNode").html('折叠');
	$("#HideTreeNode").html('隐藏');
	$("#SearchTreeNode").attr("placeholder",'搜索');

	$('#ExpandTreeNode').attr("title","未选中节点时展开所有节点,选中节点时展开选中节点")
	$('#CollapseTreeNode').attr("title","未选中节点时折叠所有节点,选中节点时折叠选中节点")
	$("#BeginTimes").text('起始时间：');
	$("#EndTimes").before('结束时间：');
	$("#TimeRanges").text('时间范围:');

	$("#custom").before("&nbsp; 自定义");
	$("#8h").before("8小时");
	$("#12h").text("12小时");
	$("#24h").text("24小时");
	$("#7d").text("7天");
	$("#op_this_shift").text("班");
	$("#op_this_day").text("日");
	$("#op_this_month").text("月");

	$('#Last').attr("title","前一时间段!");
	$('#Next').attr("title","后一时间段!");

	$('#OperationMsg').html('<strong id = "works">操作：</strong>  请在左侧结构树中选择一个样本节点以便显示数据！');


});