//=========变量定义区域===========================================================
var IsShowTagName = false;//是否显示标签名
var Data_update_interval = "30000";
var Data_update_enable = true;
var Select_node_levelcode='';
var Select_node_desc='';
var Taglists=[];//当前选中节点下的标签列表
var FULLNAME_ID_MAP=[];//名称到ID的键值对
var Firstupdatetime;//开始更新快照的时间点
var Thelasttimestemp=0;
var Id_of_settimeout;//定时器的ID
var His_info_data_Ta=[];//历史数据
var His_info_time_Ta=[];//历史数据时间
var HIS_TAG_TYPE;//查看的历史数据的类型

//==================动作响应区域==================================================
function zTreeOnClick(event, treeId, treeNode) {
	Select_node_levelcode=treeNode.treelevel;
	Select_node_desc=treeNode.name;
	Thelasttimestemp = 0;//最新时间戳清零
	$("#ThelastTimestemp").html("");//最新时间戳清零
	$("#SelectNodeName").html(Select_node_desc);
	$("#SelectNodeLevelCode").html(Select_node_levelcode);
	var url='api/taglist?levelcode='
	url+=Select_node_levelcode
	loadRealTagInfo(url)
}
//响应鼠标双击
function zTreeOnDbClick(event, treeId, treeNode){
	
}
//页面初始化工作
function pageInit(){
	
};
function OnShowTagName(){//显示变量名标签
	IsShowTagName = 1 - IsShowTagName;
	getTagInfo(null,true);
}

function OnDataUpdateSwitch(id){//自动更新开关按钮
	if(Data_update_enable==true){
		Data_update_enable=false;
		clearTimeout(Id_of_settimeout);//停止定时器
		$("#"+id).text("Start updating");
		$("#"+id).attr("class","btn btn-outline-danger btn-sm");
	}else{
		Data_update_enable=true;
		$("#"+id).text("Stop updating");
		$("#"+id).attr("class","btn btn-outline-success btn-sm");
		if(Taglists.length > 0)
			GetTagRealTimeData(0);//获取第一个标签的实时值指令
	}
};

function OnUpDateIntervalTimeChange(v){//更新间隔时间改变
	Data_update_interval = v;
};

function OnViewHistory(id){
	var end = new Date();
	var begin = new Date();
	var url;
	HIS_TAG_TYPE = Taglists[id].TagType;
	end.setTime(end.getTime());
	begin.setTime(begin.getTime()-60*60*8*1000);
	url="api/history?tagname="+Taglists[id].TagFullName+"&begintime="+DateFormat("YYYY-mm-dd HH:MM:SS",begin)+"&endtime="+DateFormat("YYYY-mm-dd HH:MM:SS",end);
	loadTagHistoryDataFromDB(url,id);
}

//===============数据接收解析区域==================================================
/*
[
  {
    "Id": 1792,
    "Variable": {
      "Id": 123,
      "VaribleType": 0,
      "EquipDicId": 0,
      "EquipCode": 0,
      "VaribleKeyClass": "",
      "VaribleName": "",
      "VaribleDescription": "",
      "VaribleValueDesc": "",
      "Seq": 0,
      "Status": 0,
      "ValueType": "",
      "PointValueType": "",
      "VaribleUnit": "",
      "FirstTagClass": "",
      "CalcKpiId": "",
      "KpiIndexDics": null,
      "Taglists": null
    },
    "ResourceType": 2,
    "ItemIDinTable": 710,
    "TreeLevelCode": "1-2-4-16393-16394-1792",
    "Dcs": {
      "Id": 64,
      "Plant": null,
      "DcsName": "",
      "Seq": 0,
      "Status": 0,
      "Taglists": null,
      "Reltable": null
    },
    "FirstTagClass": "value-meas",
    "Stage": {
      "Id": 77,
      "Workshop": null,
      "OrganizationOrgCode": 0,
      "ConstructionCode": 0,
      "ConstructionName": "",
      "StageName": "",
      "StageNameCode": "",
      "StageStatus": 0,
      "StageTechpara": "",
      "Seq": 0,
      "CreateUserId": 0,
      "CreateTime": "",
      "UpdateUserId": 0,
      "UpdateTime": "",
      "Status": 0,
      "UserNameCode": "",
      "GraphServiceUrl": "",
      "TableServiceUrl": "",
      "TreeLevelCode": "",
      "Taglists": null,
      "CheckTagLists": null,
      "Samples": null
    },
    "DeviceId": 0,
    "EquipmentName": "铜精扫II尾泵B 电磁阀",
    "Seq": 0,
    "Status": "1",
    "TagId": "x1_1_2_102_77_710_CV:1",
    "DecimalNum": 0,
    "GoldenId": 0,
    "TagPlcName": "Q12.3",
    "TagHmiName": "DB11.DBX3.5",
    "Valid": 1,
    "NormalValue": 0,
    "LimitLl": 0,
    "LimitL": 0,
    "LimitH": 0,
    "LimitHh": 0,
    "MinValue": 0,
    "MaxValue": 0,
    "Scale": 0,
    "Offset": 0,
    "TagFullName": ".x1_CMS_CT-XK_mf_fon77_OV710_CV:1",
    "TagName": "x1_CMS_CT-XK_mf_fon77_OV710_CV:1",
    "TagType": "FLOAT32",
    "TagDescription": "磨浮系统_铜精扫II尾泵B 电磁阀_阀门控制值",
    "TagUnit": "",
    "TagPracticalDescription": "精扫二尾泵B电磁阀",
    "Unit": "",
    "IsArchive": 0,
    "Digits": 0,
    "IsShutDown": 0,
    "IsStep": 0,
    "Typical": 0,
    "IsCompress": 0,
    "CompDev": 0,
    "CompDevPercent": 0,
    "CompTimeMax": 0,
    "CompTimeMin": 0,
    "ExcDev": 0,
    "ExcDevPercent": 0,
    "ExcTimeMax": 0,
    "ExcTimeMin": 0,
    "ClassOf": 0,
    "Mirror": 0,
    "MilliSecond": 0,
    "IsSummary": 0,
    "Source": "",
    "IsScan": 0,
    "Instrument": "",
    "Location1": 0,
    "Location2": 0,
    "Location3": 0,
    "Location4": 0,
    "Location5": 0,
    "UserInt1": 0,
    "UserInt2": 0,
    "UserReal1": 0,
    "UserReal2": 0,
    "Equation": "",
    "Trigger": 0,
    "TimeCopy": 0,
    "Period": 0,
    "TagAlarm": null
  }
]
*/
function getTagAlarmStatus(tag,snapv,tagid){
	var alarmmsg='';
	if(tag.TagType=='BOOL' || tag.TagType=='bool'){
		alarmmsg= "";
	}else{
		if(tag.Location5==0){
			alarmmsg= `<div><button type="button" class="btn btn-sm" onclick="onShowSnapAlarm(0,0,`+tagid+','+snapv+`);"><span class="badge badge-secondary" title="The alarm function is not enabled">Nil</span></button></div>`;
		}else{
			alarmmsg= `<div><button type="button" class="btn btn-sm" onclick="onShowSnapAlarm(0,0,`+tagid+','+snapv+`);"><span class="badge badge-success" title="The state of the variable is normal and there is no alarm information">Normal</span></button></div>`;
			if(alarmOption(4,tag.Location5) && snapv<tag.LimitL){
				alarmmsg= `<div><button type="button" class="btn btn-sm" onclick="onShowSnapAlarm(4,`+tag.LimitL+','+tagid+','+snapv+`);"><span class="badge badge-warning" title="The variable value is less than the set low limit">Low</span></button></div>`;
			}
			if(alarmOption(8,tag.Location5) && snapv>tag.LimitH){
				alarmmsg= `<div><button type="button" class="btn btn-sm" onclick="onShowSnapAlarm(8,`+tag.LimitH+','+tagid+','+snapv+`);"><span class="badge badge-warning" title="The variable value is greater than the set upper limit">High</span></button></div>`;
			}
			if(alarmOption(2,tag.Location5) && snapv<tag.LimitLl){
				alarmmsg= `<div><button type="button" class="btn btn-sm" onclick="onShowSnapAlarm(2,`+tag.LimitLl+','+tagid+','+snapv+`);"><span class="badge badge-danger" title="The variable value is less than the set lower limit">Lower</span></button></div>`;
			}
			if(alarmOption(16,tag.Location5) && snapv>tag.LimitHh){
				alarmmsg= `<div><button type="button" class="btn btn-sm" onclick="onShowSnapAlarm(16,`+tag.LimitHh+','+tagid+','+snapv+`);"><span class="badge badge-danger" title="The variable value is greater than the set higher limit">Higher</span></button></div>`;
			}
			if(alarmOption(1,tag.Location5) && snapv<tag.MinValue){
				alarmmsg= `<div><button type="button" class="btn btn-sm" onclick="onShowSnapAlarm(1,`+tag.MinValue+','+tagid+','+snapv+`);"><span class="badge badge-dark" title="The variable value is less than the set minimum value">Underflow</span></button></div>`;
			}
			if(alarmOption(32,tag.Location5) && snapv>tag.MaxValue){
				alarmmsg= `<div><button type="button" class="btn btn-sm" onclick="onShowSnapAlarm(32,`+tag.MaxValue+','+tagid+','+snapv+`);"><span class="badge badge-dark" title="The variable value is greater than the set maximum value">Overflow</span></button></div>`;
			}
		}
	}
	return alarmmsg;
}

function onShowSnapAlarm(mask,limitv,tagid,snapv){
	var msgstr="";
	switch(mask){
		case 1:
			msgstr=`<div class="alert alert-dark">The variable value is less than the set minimum value:`+limitv+`</div>`;
			break;
		case 2:
			msgstr=`<div class="alert alert-danger">The variable value is less than the set lower limit:`+limitv+`</div>`;
			break;
		case 4:
			msgstr=`<div class="alert alert-warning">The variable value is less than the set low limit:`+limitv+`</div>`;
			break;	
		case 8:
			msgstr=`<div class="alert alert-warning">The variable value is greater than the set upper limit:`+limitv+`</div>`;
			break;
		case 16:
			msgstr=`<div class="alert alert-danger">The variable value is greater than the set high limit:`+limitv+`</div>`;
			break;
		case 32:
			msgstr=`<div class="alert alert-dark">The variable value is greater than the set maximum value:`+limitv+`</div>`;
			break;
		default:
			msgstr=`<div class="alert alert-success">The state of the variable is normal and there is no alarm information。</div>`;
			break;
	}
	var tag=Taglists[tagid];
	msgstr+='<strong>Current value</strong>:<h3 style="text-align:center">'+snapv+'<small>'+tag.Unit+'</small></h3>';
	msgstr+='<strong>Alarm setting</strong></br><table class="table table-striped table-hover table-sm"><thead class="thead-light"><tr><th>Number</th><th>Name of alarm</th><th>Setting value</th><th>Alarm enable</th></tr></thead><tbody>';
	msgstr+='<tr class="table-dark text-dark"><td>1</td><td>Max</td><td>'+tag.MaxValue+'</td><td>'+ischecked(32,tag.Location5) +'</td></tr>';
	msgstr+='<tr class="table-danger"><td>2</td><td>High limit</td><td>'+tag.LimitHh+'</td><td>'+ischecked(16,tag.Location5) +'</td></tr>';
	msgstr+='<tr class="table-warning"><td>3</td><td>Upper limit</td><td>'+tag.LimitH+'</td><td>'+ischecked(8,tag.Location5) +'</td></tr>';
	msgstr+='<tr class="table-warning"><td>4</td><td>Lower limit</td><td>'+tag.LimitL+'</td><td>'+ischecked(4,tag.Location5) +'</td></tr>';
	msgstr+='<tr class="table-danger"><td>5</td><td></td>Lower limit<td>'+tag.LimitLl+'</td><td>'+ischecked(2,tag.Location5) +'</td></tr>';
	msgstr+='<tr class="table-dark text-dark"><td>6</td><td>Min</td><td>'+tag.MinValue+'</td><td>'+ischecked(1,tag.Location5) +'</td></tr>';
	msgstr+='</tbody></table>';

	ShowModal("Alarm status information",msgstr);
}
function ischecked(mask,option){
	if(alarmOption(mask,option)){
		return '<span class="badge badge-success">Yes</span>';
	}else{
		return '<span class="badge badge-secondary">No</span>';
	}
}
function alarmOption(mask,option){
	return (mask&option)>0
}

function getTagInfo(datastr,onshowtagname){//获取并显示标签的静态信息
	var dtarr = eval("("+datastr+")"); 
	var strtmp = '';
	if(onshowtagname==false){
		Taglists.splice(0,Taglists.length);//清空数组
		FULLNAME_ID_MAP.splice(0,FULLNAME_ID_MAP.length);//清空数组
		Taglists = dtarr;
	}
	if(IsShowTagName==true){
		strtmp = '<table class="table table-striped table-hover table-sm"><thead class="thead-light"><tr><th>Number</th><th>Name</th><th>Tag</th><th>Snapshot</th><th>State</th><th>History</th><th>Unit</th><th>Time stamp</th><th>Quality</th><th>Type</th></tr></thead><tbody>';
		for(i=0;i<Taglists.length;i++){
			tag=Taglists[i];
			if(tag.Unit== null){
				tag.Unit = '';
			}
			var id=i
			FULLNAME_ID_MAP[tag.TagFullName]=id;
			strtmp +='<tr id="row_'+id+'"><td>'+ (i+1) +'</td><td id="tablename_'+id+'">'+tag.TagDescription+'</td><td id="tagfullname_'+id+'">'+tag.TagFullName+'</td><td><div id = "TagValue_'+id+'"></div></td><td><div id = "AlarmStatus_'+id+'"></div></td><td><div><button type="button" class="btn btn-success btn-sm" id="History_'+id+'" onclick="OnViewHistory('+id+')" title="Click to view the last 8 hours of historical data">H</button></div></td><td>'+ tag.Unit +'</td><td><div id = "TagValueTime_'+id+'"></div></td><td><div id = "Quality_'+id+'"></div></td><td>'+tag.TagType+'</td></tr>';
		}
	}else{
		strtmp = '<table class="table table-striped table-hover table-sm"><thead class="thead-light"><tr><th>Number</th><th>Name</th><th>Snapshot</th><th>State</th><th>History</th><th>Unit</th><th>Time stamp</th><th>Quality</th><th>Type</th></tr></thead><tbody>';
		for(i=0;i<Taglists.length;i++){
			tag=Taglists[i];
			if(tag.TagUnit== null){
				tag.TagUnit = '';
			}
			var id=i
			FULLNAME_ID_MAP[tag.TagFullName]=id;
			strtmp +='<tr id="row_'+id+'"><td>'+ (i+1) +'</td><td id="tablename_'+id+'">'+tag.TagDescription+'</td><td><div id = "TagValue_'+id+'"></div></td><td><div id = "AlarmStatus_'+id+'"></div></td><td><div><button type="button" class="btn btn-success btn-sm" id="History_'+id+'" onclick="OnViewHistory('+id+')" title="Click to view the last 8 hours of historical data">H</button></div></td><td>'+ tag.Unit +'</td><td><div id = "TagValueTime_'+id+'"></div></td><td><div id = "Quality_'+id+'"></div></td><td>'+tag.TagType+'</td></tr>';
		}
	}
	strtmp +='</tbody></table>'
	$("#taglistinfo").html(strtmp);//更新TagList表
	$("#SelectNodeTagCnt").html(Taglists.length);//更新节点总数
	if(Taglists.length > 0 && Data_update_enable==true)
		GetTagRealTimeData(0);//获取第一个标签的实时值
}

function UpdateTagRealTimeData(ajaxdata,id){//根据Ajax反馈的结果更新Tag的实时数据
	/*
	"micbox1-2.x3_asl_asl-xc1_mf1_tbyx3_f1-002_fd1-002_run-current:1": {
		"Id": 185,
		"Rtsd": {
		"Time": 1589846274000,
		"Value": 2310430464,
		"Quality": 0
		},
		"Err": ""
  	}
	*/

	var dtarr = eval("("+ajaxdata+")"); 
	var d = new Date();
	var t;
	var qualitystr;
	//for(i=0;i<dtarr.length;i++){
	for(let tag in dtarr){
		var tagd = dtarr[tag];//数据
		//console.log(tagd);
		if (tagd.Id>0){
			var value=DataToFixed(tagd.Rtsd.Value,Taglists[FULLNAME_ID_MAP[tag]].TagType,Taglists[FULLNAME_ID_MAP[tag]].DecimalNum);
			$("#TagValue_"+FULLNAME_ID_MAP[tag]).html(value);//更新TagValue
			d.setTime(tagd.Rtsd.Time);//将2006-05-06T00:00:00Z格式的时间转换为UTC时间戳
			$("#TagValueTime_"+FULLNAME_ID_MAP[tag]).html(d.toLocaleString());//更新时间戳:转换为当地时间格式
			$("#AlarmStatus_"+FULLNAME_ID_MAP[tag]).html(getTagAlarmStatus(Taglists[FULLNAME_ID_MAP[tag]],value,FULLNAME_ID_MAP[tag]));//更新报警状态信息
			switch(tagd.Rtsd.Quality){
				case 0:
					qualitystr = "Good";
					break;
				case 1:
					qualitystr = "NODATA";
					break;
				case 2:
					qualitystr = "CREATED";
					break;
				case 3:
					qualitystr = "SHUTDOWN";
					break;
				case 4:
					qualitystr = "CALCOFF";
					break;
				case 5:
					qualitystr = "BAD";
					break;
				case 6:
					qualitystr = "DIVBYZERO";
					break;
				case 7:
					qualitystr = "REMOVED";
					break;
				case 256,511:
					qualitystr = "OPC";
					break;
				case 512,32767:
					qualitystr = "USER";
					break;
				default:
					qualitystr = "Undefined";
					break;
			}
			$("#Quality_"+FULLNAME_ID_MAP[tag]).html(qualitystr);
			
			if(tagd.Rtsd.Time>Thelasttimestemp){//更新最新时间戳
				Thelasttimestemp = tagd.Rtsd.Time;
				$("#ThelastTimestemp").html(d.toLocaleString());
			}
		}else{
			var dataerr=`<div><button type="button" class="btn btn-sm" onclick="onShowSnapErr('`+tagd.Err+`');"><span class="badge badge-danger">#Error</span></button></div>`
			$("#TagValue_"+FULLNAME_ID_MAP[tag]).html(dataerr);//更新TagValue
		}
	}
	if(id < Taglists.length && Data_update_enable==true){//如果id不大于tag总数,继续获取下一个标签的实时值
		GetTagRealTimeData(id);
	}else{
		Id_of_settimeout = setTimeout(function(){if(Taglists.length > 1 && Data_update_enable==true){GetTagRealTimeData(0);}},Data_update_interval);
	}
}



function onShowSnapErr(errmsg){
	ShowModal("错误信息",`<div class="alert alert-danger">`+errmsg+`</div>`);
}

function GetTagRealTimeData(id){//获取Tag的实时数据指令
	var step=50;//读取数据的步长
	var len=Taglists.length;
	if(id==0)
		Firstupdatetime = new Date();//记录启动更新的时间戳
	var urlstr = "api/snapshot?tagnames=";
	var i=0;
	for(i=0;i<step;i++){
		if((id)<len){
			if(i>0){
				urlstr+=',';
			}
			urlstr+=Taglists[id].TagFullName;
		}else{
			break;
		}
		id+=1;
	}
	loadTagRealTimeData(urlstr,id)//启动Ajax
	if(len > 0)//计算更新完成百分比
		var p = id/(len)*100;
	if(p<100)//完成度小于100%
		if(step==1){
			$("#InUpdateID").html(p.toFixed(2)+"%，Updating:"+id);//数据更新进度显示
		}else{
			$("#InUpdateID").html(p.toFixed(2)+"%，Updating:"+(id-i+1)+"~"+id);//数据更新进度显示
		}
	else{//完成度等于100%
		var timediff = new Date() - Firstupdatetime;//计算完成耗时,ms
		$("#InUpdateID").html(p.toFixed(2)+"%,Time:"+timediff/1000+" S");
	}
};

function getTagHistoryData(ajaxdata,id){//获取从DB读到的历史数据
	//console.log(ajaxdata);
	var dtarr = eval("("+ajaxdata+")"); 
	var title;
	His_info_time_Ta.splice(0,His_info_time_Ta.length);//清空数组
	His_info_data_Ta.splice(0,His_info_data_Ta.length);//清空数组
	if(dtarr != null){
		for(let tag in dtarr){
			var tagd = dtarr[tag];
			if (tagd.length>0){
				for(var i=0;i< tagd.length;i++){
					His_info_time_Ta[i] = tagd[i].Time;//d.toLocaleString();//更新时间戳:转换为当地时间格式
					His_info_data_Ta[i] = DataToFixed(tagd[i].Value,Taglists[id].TagType,Taglists[id].DecimalNum);//更新TagValue
				}
				title=Taglists[id].TagDescription;
			}else{
				title='<div class="alert alert-danger">Note: the last 8 hours of historical data are not read！</div>';
			}
		}
	}else{
		title='<div class="alert alert-danger">Note: the last 8 hours of historical data are not read!</div>';
	}
	//showModal("最近8小时历史数据");
	$("#echart_title").html(title);//显示Echarts
	$("#echarts_his_Ta").show();//显示Echarts
	ShowModal("Last 8 hours of historical data","");
	refreshEcharts_his_Ta();//刷新Echarts 
};

//=========AJAX函数定义区域=======================================================
function loadRealTagInfo(urlstr)//加载Taglist信息
{
	//调用公用的loadXMLDoc函数
	$("#taglistinfo").html('<div class="alert alert-warning">Loading data……</div>');//更新TagList表
    loadXMLDoc(urlstr,function()
    {
        if (xmlhttp.readyState==4 && xmlhttp.status==200)//请求处理完成，且状态OK
        {//添加请求完成后的处理功能---------------------------------------
			getTagInfo(xmlhttp.responseText,false);
			//alert(xmlhttp.responseText);
        }//请求完成后的处理功能结束---------------------------------------
    });
}
function loadTagRealTimeData(urlstr,id)//加载Tag的实时数据信息
{
	//调用公用的loadXMLDoc函数
    loadXMLDoc(urlstr,function()
    {
        if (xmlhttp.readyState==4 && xmlhttp.status==200)//请求处理完成，且状态OK
        {//添加请求完成后的处理功能---------------------------------------
			UpdateTagRealTimeData(xmlhttp.responseText,id);
        }//请求完成后的处理功能结束---------------------------------------
    });
}

function loadTagHistoryDataFromDB(urlstr,id)//从数据库中读取单一变量指定时间段的历史数据
{
	//调用公用的loadXMLDoc函数
    loadXMLDoc(urlstr,function()
    {
        if (xmlhttp.readyState==4 && xmlhttp.status==200)//请求处理完成，且状态OK
        {//添加请求完成后的处理功能---------------------------------------
			//document.getElementById("viewinfo2").innerHTML=xmlhttp.responseText;
			getTagHistoryData(xmlhttp.responseText,id);
			//showModal("历史数据",xmlhttp.responseText);
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
			ShowModal("Updating tips",'<div class="alert alert-success">'+xmlhttp.responseText+'</div>');
        }//请求完成后的处理功能结束---------------------------------------
    });
}
function onUpdateTree(){
	var url="api/updatetagnodetree?withtag=0";
	loadUpdateTree(url);
}
$("#modal-btn").text('close');
$(document).ready(function() {
	$("#exit").after('Exit');
	$("#UpdateTree").html('Update')
	$("#ExpandTreeNode").html('Open');
	$('#ExpandTreeNode').attr("title","Expand all nodes when no nodes are selected, and expand selected nodes when nodes are selected");
	$('#CollapseTreeNode').attr("title","Collapse all nodes when no nodes are selected, and collapse selected nodes when nodes are selected");
	$("#CollapseTreeNode").html('Fold');
	$("#HideTreeNode").html('Hide');
	$('#SearchTreeNode').attr("placeholder","Search")
	//$("#webAPI").html('Fold');

	$("#SelectNodeName").before('Node Name: &nbsp;');
	$("#SelectNodeLevelCode").before('Node evel: &nbsp;');
	$("#SelectNodeTagCnt").before('Total number of variables: &nbsp;');
	$("#ThelastTimestemp").before('latest timestamp: &nbsp;');
	$("#UpdateProgress").before('Update progress: &nbsp;');
	$("#InUpdateID").before('');
	$("#UpdateInterval").before('Update interval: &nbsp;');
	$("#show_tagname").before('Display tag name: &nbsp;');

	$("#dataUpdateSwitch").html("Stop updating");
	$("#taglistinfo").html('<strong>Operation: </strong> Please select a node in the left structure tree to display the variable snapshot data under it!');
});

$("#MyModal-Text").html("<div class='col-sm-12 font-weight-bold' id='echart_title'></div><div class='col-sm-12' id='echarts_his_Ta' style='height: 400px;width:768px;display: ;border: 1px solid #cecece;'></div>");//初始化设置模态框内容

