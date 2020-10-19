//=========变量定义区域===========================================================
var NODE_IS_REPORT=false; //所选节点是报表
var NODE_ID=0;//所选层级的ID
var REPORT_ID=0;//当前查看的报表模板列表的ID
var WORKSHOPS=new Array();//车间列表
var WORKSHOP_KVA=new Array();//车间KV值列表,ID为KEY
var NODES_KVA=new Array();//层级节点KV值列表,ID为KEY
var WORKSHOP_SELECTOR="";//车间选择器
var PARENT_LEVEL_SELECTOR="";//父层级选择器
var PERIOD_UDF=false;//自定义周期
var LAST_SELECTED_FILE="";//上一个被选中的报表

//==================动作响应区域==================================================
function zTreeOnClick(event, treeId, treeNode) {
	NODE_ID = treeNode.id;
	if (treeNode.isParent > 0){
		NODE_IS_REPORT = false;//所选节点是文件夹NewLevel
		$("#NewLevel").removeAttr("disabled");
		$("#UploadTpl").attr("disabled","disabled");
	}else{
		NODE_IS_REPORT = true;//所选节点是报表
		$("#UploadTpl").removeAttr("disabled");
		$("#NewLevel").attr("disabled","disabled");
	}
	$("#FileVew").hide();
	requestChildNodes(treeNode.levelcode);//载入子集
}
//响应鼠标双击
function zTreeOnDbClick(event, treeId, treeNode){
	
}
//页面初始化
function pageInit(){
	getReportsDatas(NODESMSG);
	requestWorkshopLists();//获取车间列表
}

function onDeleteLevel(id){
	var htmlstr=`<div class="alert alert-danger">确定要删除[<strong>`+NODES_KVA[id].Name+`</strong>]吗！</div>`;
	var btsubmit='<button class="btn btn-danger btn" onclick="onDelete('+id+');">确定</button>';
	var btcancel='<button class="btn btn-success btn" data-dismiss="modal">取消</button>';
	htmlstr+=`<div class="row"><div class="col-4"></div><div class="col-4  btn-group">`+btsubmit+btcancel+`</div><div class="col-4"></div></div>`;
	ShowModal("删除层级",htmlstr);
}
function onDelete(id){
	requestDeleteLevel(id);
}
//父层级发生变化
function onParent(id){
	//{BaseTime:node.BaseTime,ShiftHour:node.ShiftHour,WorkShop:node.Workshop.Id,TemplateUrl:node.TemplateUrl,ResultUrl:node.ResultUrl}
	if(id>0){
		if(NODES_KVA[id].BaseTime.length > 18)
			$("#ReportBaseTime").val(InputDateTimeToString("YYYY-mm-ddTHH:MM:SS",NODES_KVA[id].BaseTime));
		if(NODES_KVA[id].ShiftHour>8)
			$("#ReportShifthour").val(NODES_KVA[id].ShiftHour);
		if(NODES_KVA[id].WorkShop > 0)
			$("#ReportWorkShop").val(NODES_KVA[id].WorkShop);
		if(NODES_KVA[id].TemplateUrl.length>5)
			$("#ReportTplUrl").val(NODES_KVA[id].TemplateUrl);
		if(NODES_KVA[id].ResultUrl.length>5)
			$("#ReportResultUrl").val(NODES_KVA[id].ResultUrl);
	}
}
//所属车间发生变化
function onWorkShop(id){
	if(id>0){
		if(WORKSHOP_KVA[id].BaseTime.length > 18)
			$("#ReportBaseTime").val(InputDateTimeToString("YYYY-mm-ddTHH:MM:SS",WORKSHOP_KVA[id].BaseTime));
		if(WORKSHOP_KVA[id].ShiftHour>8)
			$("#ReportShifthour").val(WORKSHOP_KVA[id].ShiftHour);
	}
}
//周期发生变化
function onPeriod(period){
	var udf = PERIOD_UDF;
	var htmlstr=``;
	if (period>=0){
		PERIOD_UDF = true;
		htmlstr+=`<label for="ReportPeriod">计算周期(秒):</label><input type="number" class="form-control" name="ReportPeriod" id="ReportPeriod" oninput ="onPeriod(this.value)" value="3600">`;
	}else{
		PERIOD_UDF = false;
		htmlstr+=`<label for="ReportPeriod">计算周期:</label><select class="form-control" name="ReportPeriod" id="ReportPeriod" oninput ="onPeriod(this.options[this.options.selectedIndex].value)">
		<option value="0">自定义</option>
		<option value="-1">每小时</option>
		<option value="-2">每班</option>
		<option value="-3">每日</option>
		<option value="-4" >每月</option>
		<option value="-5" >每季度</option>
		<option value="-6" >每年</option>
</select>`;
	}
	if(PERIOD_UDF != udf){//状态改变的时候才变更
		$("#SetPeriod").html(htmlstr);
		if(PERIOD_UDF==false){
			if(period<-6){
				period=-3;
			}
			$("#ReportPeriod").val(period);
		}
	}
}

//选择模板作为当前模板
function onTplSelect(filename){
	var formData = new FormData();

	formData.append("id",NODE_ID);
	formData.append("oldfilename",NODES_KVA[NODE_ID].TemplateFile);
	formData.append("newfilename",filename);
	formData.append("reportname",NODES_KVA[NODE_ID].Name);

	var xhr=new XMLHttpRequest();
    xhr.open("post","api/setfileasreporttpl");
    xhr.send(formData);
    xhr.onload=function(){
        if(xhr.status==200){
			//alert(xhr.responseText);
			requestTplLists(NODE_ID); 
        }
    }
}
//下载模板
function onTplDown(filename){
	var urlstr = "api/download?filepath="+NODES_KVA[REPORT_ID].TemplateUrl+"&filename="+filename+"&id=0";
	//console.log(urlstr);
	window.location.href=urlstr;
}
//编辑模板
function onTplEdit(filename){

}
//预览模板
function onTplView(filename,thisid){
	$("#"+thisid).attr("class","btn btn-primary btn-sm");
	if(LAST_SELECTED_FILE!=""){
		$("#"+LAST_SELECTED_FILE).attr("class","btn btn-outline-primary btn-sm");
	}
	LAST_SELECTED_FILE=thisid;
	var urlstr = "api/viewexcel?filepath="+NODES_KVA[REPORT_ID].TemplateUrl+"&filename="+filename+"&id=0";
	//console.log(NODE_ID,urlstr);
	loadViewFile(urlstr);
}
//==================AJAX数据解析区域==================================================
//解析报表列表
function getReportsDatas(datas){
	NODES_KVA[0]={Name:""};
	PARENT_LEVEL_SELECTOR=`<select class="form-control" name="ReportParent" id="ReportParent" oninput="onParent(this.options[this.options.selectedIndex].value)">
					<option value="0" selected="selected">无</option>`;
	
	var htmlstr='<table class="table table-striped table-hover table-sm"><thead class="thead-light"><tr><th style="width:50px">序号</th><th style="width:150px">名称</th><th style="width:150px">父层级</th><th style="width:50px">模板</th><th style="width:50px">调试</th><th style="width:50px">状态</th><th style="width:50px">排序</th><th style="width:150px">计算截至时间</th><th style="width:150px">更新时间</th><th style="width:200px">操作</tr></thead><tbody>';
	var firstRport=0;
	if (datas!=null){
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
			if(REPORT_ID==0 && node.Folder==0){
				REPORT_ID = node.Id;
			}
		}
		if(NODE_ID>0){
			PARENT_LEVEL_SELECTOR+=`<option value="`+NODES_KVA[NODE_ID].Pid+`">`+NODES_KVA[NODES_KVA[NODE_ID].Pid].Name+`</option>`;
		}
		for (var i=0;i<datas.length;i++){
			var node=datas[i];
			var status = '';
			var isdebug= '';
			var isreport='<span class="badge badge-secondary">否</span>';
			var lasttime='';

			var btviewtpl='<button class="btn btn-outline-success btn-sm" onclick="requestTplLists('+node.Id+');">模板列表</button>';
			var btedit='<button class="btn btn-outline-primary btn-sm" onclick="onEdit('+node.Id+');">编辑</button>';
			var btdelete='<button class="btn btn-outline-danger btn-sm" onclick="onDeleteLevel('+node.Id+');">删除</button>';
			var btngroup = btedit+btdelete+btviewtpl;
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
				if(firstRport==0){
					firstRport= node.Id;
				}
			}else{
				btngroup = btedit+btdelete;
				PARENT_LEVEL_SELECTOR+=`<option value="`+node.Id+`">`+node.Name+`</option>`;
			}
			
			htmlstr+='<tr><td>'+(i+1)+'</td><td>'+node.Name+'</td><td>'+NODES_KVA[node.Pid].Name+'</td><td>'+isreport+'</td><td>'+isdebug+'</td><td>'+status+'</td><td>'+node.Seq+'</td><td>'+lasttime+'</td><td>'+node.UpdateTime+'</td><td><div class="btn-group">'+btngroup+'</div></td></tr>';
		}
	}
	htmlstr+='</tbody></table>';
	PARENT_LEVEL_SELECTOR+=`</select>`;
	$("#DataFrame").html(htmlstr);
	if(firstRport>0 && NODE_ID>0)
		requestTplLists(firstRport);
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
function getTplLists(datas){
	/*
		FileName string //文件名
		FileTime string //文件上传时间
		ModTime  string //文件编辑时间
		Size     int64  //单位:字节
		IsCurrent bool   //是否当前所选模板
	*/
	var htmlstr=`<div class="col-12"><strong>`+NODES_KVA[REPORT_ID].Name+`</strong></div>
	<table class="table table-striped table-hover table-sm"><thead class="thead-light"><tr><th style="width:50px">序号</th><th style="width:100px">上传时间</th><th style="width:100px">编辑时间</th><th style="width:80px">文件大小</th><th style="width:80px">当前模板</th><th style="width:200px">操作</th></tr></thead><tbody>`;
	if (datas!=null){
		datas = bubbleSort(datas);
		for (var i=0;i<datas.length;i++){
			var node=datas[i];
			var fname = node.FileName.split(".",1);
			var iscurrent=`<span class="badge badge-light">否</span>`;
			var btview=`<button class="btn btn-outline-primary btn-sm" id="`+fname+`" onclick="onTplView('`+node.FileName+`',this.id);" title="在线预览并编辑报表模板">预览&编辑</button>`;
			var btselect=`<button class="btn btn-outline-success btn-sm" onclick="onTplSelect('`+node.FileName+`');" title="选为当前模板">选为当前</button>`;
			var btdown=`<button class="btn btn-outline-warning btn-sm" onclick="onTplDown('`+node.FileName+`');" title="下载模板文件到本地">下载</button>`;
			if (node.IsCurrent==true){
				onTplView(node.FileName);
				LAST_SELECTED_FILE = fname;
				iscurrent=`<span class="badge badge-success">是</span>`;
				btview=`<button class="btn btn-primary btn-sm" id="`+fname+`" onclick="onTplView('`+node.FileName+`',this.id);" title="在线预览并编辑报表模板">预览&编辑</button>`;
				btselect=`<button class="btn btn-success btn-sm" onclick="onTplSelect('`+node.FileName+`');" disabled title="已是当前模板">当前模板</button>`;
			}

			htmlstr+='<tr><td>'+(i+1)+'</td><td>'+node.FileTime+'</td><td>'+node.ModTime+'</td><td>'+DataToFixed(node.Size/1024,"float",2)+'KB</td><td>'+iscurrent+'</td><td><div class="btn-group">'+btview+btselect+btdown+'</div></td></tr>';
		}
	}
	htmlstr+='</tbody></table>';
	$("#TplLists").html(htmlstr);
}
//解析车间列表
function getWorkshopLists(datas){
	WORKSHOPS=datas;
	var htmlstr=`<select class="form-control" name="ReportWorkShop" id="ReportWorkShop" oninput="onWorkShop(this.options[this.options.selectedIndex].value)">
					<option value="0" selected="selected">无</option>`;
	for(var i=0;i<datas.length;i++){
		data = datas[i];
		WORKSHOP_KVA[data.Id] = {Name:data.WorkshopName,BaseTime:data.BaseTime,ShiftHour:data.ShiftHour};
		htmlstr+=`<option value="`+data.Id+`">`+data.WorkshopName+`</option>`;
	}
	htmlstr+=`</select>`;
	WORKSHOP_SELECTOR = htmlstr;
}
//=========AJAX请求定义区域=======================================================
//读取所选层级下的所有子集
function requestChildNodes(levelcode){
	var urlstr = "api/getreportchildnodes?levelcode="+levelcode;
	loadChildNodes(urlstr);
}
//读取所选报表的模板列表
function requestTplLists(id){
	REPORT_ID = id;
	var urlstr = "api/getreporttpllist?id="+id;
	console.log(urlstr);
	loadTplLists(urlstr);
}
//读取车间列表
function requestWorkshopLists(){
	var urlstr = "api/getworkshops";
	loadWorkshopLists(urlstr);
}
//删除层级
function requestDeleteLevel(id){
	var urlstr = "api/deletereportlevel?id="+id;
	loadDeleteLevel(urlstr);
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

function loadTplLists(urlstr)//读取所选报表的模板列表
{	
	$("#FileVew").hide();
	$("#TplLists").show();
	$("#TplLists").html('<div class="alert alert-warning">正在加载数据……</div>');
	//调用公用的loadXMLDoc函数
    loadXMLDoc(urlstr,function()
    {
        if (xmlhttp.readyState==4 && xmlhttp.status==200)//请求处理完成，且状态OK
        {//添加请求完成后的处理功能---------------------------------------
			getTplLists(eval("("+xmlhttp.responseText+")"));//解读数据
			//下一步：		
		}//请求完成后的处理功能结束---------------------------------------

    });
}

function loadWorkshopLists(urlstr)//读取车间列表
{	
	//调用公用的loadXMLDoc函数
    loadXMLDoc(urlstr,function()
    {
        if (xmlhttp.readyState==4 && xmlhttp.status==200)//请求处理完成，且状态OK
        {//添加请求完成后的处理功能---------------------------------------
			getWorkshopLists(eval("("+xmlhttp.responseText+")"));//解读数据
			//下一步：		
        }//请求完成后的处理功能结束---------------------------------------
    });
}

function loadDeleteLevel(urlstr)//删除报表层级
{	
	//调用公用的loadXMLDoc函数
    loadXMLDoc(urlstr,function()
    {
        if (xmlhttp.readyState==4 && xmlhttp.status==200)//请求处理完成，且状态OK
		{//添加请求完成后的处理功能---------------------------------------
			if(xmlhttp.responseText.indexOf("ok") != -1){
				location.reload();
				ShowModal("删除层级",'<div class="alert alert-success">删除成功！</div>');//解读数据
			}else{
				ShowModal("删除层级",'<div class="alert alert-warning">删除失败:'+xmlhttp.responseText+'</div>');//解读数据
			}
			//下一步：		
        }//请求完成后的处理功能结束---------------------------------------
    });
}
function loadViewFile(urlstr)//加载文件
{	
	$("#FileVew").show();
	$("#FileVew").html('<div class="alert alert-warning">正在加载文件……</div>');
	CELLSMAP=new Array();
	//调用公用的loadXMLDoc函数
    loadXMLDoc(urlstr,function()
    {
        if (xmlhttp.readyState==4 && xmlhttp.status==200)//请求处理完成，且状态OK
		{//添加请求完成后的处理功能---------------------------------------
			$("#FileVew").html(xmlhttp.responseText);//显示数据
			$("#FileVew").attr("style","max-height:"+(winH-130)+"px;overflow: auto;");
			$("#EditForm").show();//显示数据
			$("#Atention").show();//显示注意
			$("#EditForm").attr("style","width: 100%;display: flex;padding:1px;");//显示编辑区域
			$("#selectTdValue").val("");
			$("#CellAxis").val("");
			//下一步：		
        }//请求完成后的处理功能结束---------------------------------------
    });
}

function postMsg(urlstr,msg)//POST
{	
	//调用公用的loadXMLDoc函数
    $.post(urlstr,msg,function(data,status){
		//添加请求完成后的处理功能---------------------------------------
		console.log(data);
		alert(data);
		location.reload();	
        //请求完成后的处理功能结束---------------------------------------
    });
}
function onEdit(id){
	modaltext = `	
<form action="" role="form">
<div class="container col-12 border rounded"><div class="row">
	<div class="col-6">
	<strong>基本信息</strong>
	<div class="col-sm-12 border rounded" style="margin-top:3px;padding-top:3px">
	  <div class="form-group">
		<label for="ReportLeveltName">名称:</label>
		<input type="text" class="form-control" name="ReportLeveltName" id="ReportLeveltName" value="`+NODES_KVA[id].Name+`">
	  </div>
	  <div class="form-group">
		<label for="ReportParent">父层级:</label>`+PARENT_LEVEL_SELECTOR+`
	  </div>
	  <div class="form-group">
		<label for="ReportWorkShop">所属车间:</label>`+WORKSHOP_SELECTOR+`
	  </div>
	  <div class="form-group">
		<label for="ReportRemark">备注:</label>
		<input type="text" class="form-control" name="ReportRemark" id="ReportRemark" value="`+NODES_KVA[id].Remark+`">
	  </div>
	  <div class="form-group" style="display:none">
		<label for="ReportDistId">分布式ID号:</label>
		<input type="number" class="form-control" name="ReportDistId" id="ReportDistId" value="`+NODES_KVA[id].DistributedId+`" disabled>
	  </div>
	  <div class="form-group">
		<label for="ReportSeq">排序号:</label>
		<input type="number" class="form-control" name="ReportSeq" id="ReportSeq" value="`+NODES_KVA[id].Seq+`">
	  </div>
	  <div class="row col-12 container"><div class="form-check-inline">
		<div class="form-group col border rounded" style="margin-right:5px">
			<label for="IsForder">模板:</label><br/>	
			<label class="radio-inline"><input type="radio" name="IsForder" id="Forder_N" onclick="onIsForder(0);" value="0" checked>是</label>
			<label class="radio-inline"><input type="radio" name="IsForder" id="Forder_Y" onclick="onIsForder(1);" value="1" >否</label>
		</div>
		<div class="form-group col border rounded" style="margin-right:5px">
			<label for="IsDebug">调试:</label><br/>	
			<label class="radio-inline"><input type="radio" name="IsDebug" id="Debug_Y" onclick="onIsDebug(1);" value="1" >是</label>
			<label class="radio-inline"><input type="radio" name="IsDebug" id="Debug_N" onclick="onIsDebug(0);" value="0" checked>否</label>
		</div>
		<div class="form-group col border rounded">
		<label for="ReportStatus">状态:</label><br/>	
			<label class="radio-inline"><input type="radio" name="ReportStatus" id="Status_Y" value="1">启用</label>
			<label class="radio-inline"><input type="radio" name="ReportStatus" id="Status_N" value="0">停用</label>
	  	</div>
	  </div></div>
	  </div>
	</div>
	<div class="col-6">
	  <strong>时间信息</strong>
	  <div class="col-sm-12 border rounded" style="margin-top:3px;padding-top:3px">
	  <div class="form-group">
		<label for="ReportBaseTime">基准时间:</label>
		<input type="datetime-local" class="form-control" name="ReportBaseTime" id="ReportBaseTime" value="`+NODES_KVA[id].BaseTime.replace(" ","T")+`">
	  </div>
	  <div class="form-group">
		<label for="ReportShifthour">每班工作时间(小时):</label>
		<input type="number" class="form-control" name="ReportShifthour" id="ReportShifthour" value="`+NODES_KVA[id].ShiftHour+`">
	  </div>
	  <div class="form-group" id="SetPeriod">
		<label for="ReportPeriod">计算周期:</label>
			<select class="form-control" name="ReportPeriod" id="ReportPeriod" oninput ="onPeriod(this.options[this.options.selectedIndex].value)" value="`+NODES_KVA[id].Period+`">
					<option value="0">自定义</option>
					<option value="-1">每小时</option>
					<option value="-2">每班</option>
					<option value="-3">每日</option>
					<option value="-4" >每月</option>
					<option value="-5" >每季度</option>
					<option value="-6" >每年</option>
			</select>
	  </div>
	  <div class="form-group">
		<label for="ReportStartTime">起始时间:</label>
		<input type="datetime-local" class="form-control" name="ReportStartTime" id="ReportStartTime" value="`+NODES_KVA[id].StartTime.replace(" ","T")+`">
	  </div>
	  <div class="form-group">
		<label for="ReportLastTime">最后计算截至时间:</label>
		<input type="datetime-local" class="form-control" name="ReportLastTime" id="ReportLastTime" value="`+NODES_KVA[id].LastCalcTime.replace(" ","T")+`">
	  </div>
	  <div class="form-group">
		<label for="ReportOffsetMinute">偏移时间(分钟):</label>
		<input type="number" class="form-control" name="ReportOffsetMinute" id="ReportOffsetMinute" value="`+NODES_KVA[id].OffsetMinutes+`">
	  </div>
	  </div>
	</div>
  </div>

	<div class="row" style="display:none"><div class="col">
	<strong>路径信息</strong>
	<div class="col-sm-12 border rounded" style="margin-top:3px;padding-top:3px">
	<div class="form-group">
	  <label for="ReportTplUrl">模板路径(默认路径,不建议更改):</label>
	  <input type="text" class="form-control" name="ReportTplUrl" id="ReportTplUrl" value="`+NODES_KVA[id].TemplateUrl+`">
	</div>
	<div class="form-group">
	  <label for="ReportResultUrl">结果存储路径(默认路径,不建议更改):</label>
	  <input type="text" class="form-control" name="ReportResultUrl" id="ReportResultUrl" value="`+NODES_KVA[id].ResultUrl+`">
	</div>
	</div>
	</div></div>
</div></form>
<div class="row" style="margin-top:3px"><div class="col-4"></div><div class="col-4  btn-group"> <button class="btn btn-success" onclick="onSubmitEdit(`+id+`)">提交</button></div><div class="col-4"></div></div>
	`;
	ShowModal("编辑层级",modaltext);
	$("#ReportParent").val(NODES_KVA[id].Pid);
	if(NODES_KVA[id].Folder==1){
		$("#ReportParent").attr("disabled","disabled");//文件夹的父层级不可更改
		$("#Status_Y").attr("disabled","disabled");//文件夹的状态不可更改
		$("#Status_N").attr("disabled","disabled");//文件夹的状态不可更改
		$("#Debug_Y").attr("disabled","disabled");//是否调试模式不可更改
		$("#Debug_N").attr("disabled","disabled");//是否调试模式不可更改
	}
	$("#ReportWorkShop").val(NODES_KVA[id].WorkShop);
	$("#ReportPeriod").val(NODES_KVA[id].Period);
	$("input[name='IsForder'][value='"+NODES_KVA[id].Folder+"']").attr("checked",true);//
	$("input[name='IsDebug'][value='"+NODES_KVA[id].Debug+"']").attr("checked",true);//
	$("input[name='ReportStatus'][value='"+NODES_KVA[id].Status+"']").attr("checked",true);
	
}
function onSubmitEdit(id){
	var msg={
		Id			  : id,
		DistributedId : $("#ReportDistId").val(),     //分布式计算ID，为0的时候不区分计算引擎，不为0的时候需要id适配的计算引擎进行计算
		Name          : $("#ReportLeveltName").val(), //名称
		Pid           : $("#ReportParent").val(),     //父级菜单ID
		WorkshopId    : $("#ReportWorkShop").val(),   //所属车间ID
		Level         : NODES_KVA[id].Level ,//层级深度
		LevelCode     : NODES_KVA[id].LevelCode ,//层级码
		Folder        : $("input[name='IsForder']:checked").val(),//$("#IsForder").val(),     //是否是文件夹，1-是，0-否
		Debug		  : $("input[name='IsDebug']:checked").val(),//$("#IsDebug").val(),     //是否是文件夹，1-是，0-否
		Seq           : $("#ReportSeq").val(),    //排序号
		Remark        : $("#ReportRemark").val(), //备注
		TemplateUrl   : $("#ReportTplUrl").val(), //模板文件路径
		TemplateFile  : NODES_KVA[id].TemplateFile ,//模板文件名称
		ResultUrl     : $("#ReportResultUrl").val(),    //结果地址
		StartTime     : $("#ReportStartTime").val(),    //统计计算开始起作用的时间
		Period        : $("#ReportPeriod").val(),      //计算周期,详见KPI表
		OffsetMinutes : $("#ReportOffsetMinute").val(), //偏移时间
		LastCalcTime  : $("#ReportLastTime").val(),     //最后计算时间
		BaseTime      : $("#ReportBaseTime").val(),     //基准时间
		ShiftHour     : $("#ReportShifthour").val(),    //每班工作时间
		Status        : $("input[name='ReportStatus']:checked").val(), //1有效 0无效
	};
	NODES_KVA[id].DistributedId = $("#ReportDistId").val(),     //分布式计算ID，为0的时候不区分计算引擎，不为0的时候需要id适配的计算引擎进行计算
	NODES_KVA[id].Name          = $("#ReportLeveltName").val(), //名称
	NODES_KVA[id].Pid           = $("#ReportParent").val(),     //父级菜单ID
	NODES_KVA[id].WorkshopId    = $("#ReportWorkShop").val(),   //所属车间ID
	NODES_KVA[id].Folder        = $("input[name='IsForder']:checked").val(),//$("#IsForder").val(),     //是否是文件夹，1-是，0-否
	NODES_KVA[id].Debug		  = $("input[name='IsDebug']:checked").val(),//$("#IsDebug").val(),     //是否是文件夹，1-是，0-否
	NODES_KVA[id].Seq           = $("#ReportSeq").val(),    //排序号
	NODES_KVA[id].Remark        = $("#ReportRemark").val(), //备注
	NODES_KVA[id].TemplateUrl   = $("#ReportTplUrl").val(), //模板文件路径
	NODES_KVA[id].ResultUrl     = $("#ReportResultUrl").val(),    //结果地址
	NODES_KVA[id].StartTime     = $("#ReportStartTime").val(),    //统计计算开始起作用的时间
	NODES_KVA[id].Period        = $("#ReportPeriod").val(),      //计算周期,详见KPI表
	NODES_KVA[id].OffsetMinutes = $("#ReportOffsetMinute").val(), //偏移时间
	NODES_KVA[id].LastCalcTime  = $("#ReportLastTime").val(),     //最后计算时间
	NODES_KVA[id].BaseTime      = $("#ReportBaseTime").val(),     //基准时间
	NODES_KVA[id].ShiftHour     = $("#ReportShifthour").val(),    //每班工作时间
	NODES_KVA[id].Status        = $("input[name='ReportStatus']:checked").val(), //1有效 0无效
		
	console.log(msg);
	HideModal();
	postMsg("/api/editreportlevel",msg);
}

function onNewLevel(){
	var now = new Date;
	var nowstring = DateFormat("YYYY-mm-ddTHH:MM",now);
	modaltext = `	
	<form action="" role="form">

<div class="container col-12 border rounded"><div class="row">
	<div class="col-6">
	<strong>基本信息</strong>
	<div class="col-sm-12 border rounded" style="margin-top:3px;padding-top:3px">
	  <div class="form-group">
		<label for="ReportLeveltName">名称:</label>
		<input type="text" class="form-control" name="ReportLeveltName" id="ReportLeveltName" value="">
	  </div>
	  <div class="form-group">
		<label for="ReportParent">父层级:</label>`+PARENT_LEVEL_SELECTOR+`
	  </div>
	  <div class="form-group">
		<label for="ReportWorkShop">所属车间:</label>`+WORKSHOP_SELECTOR+`
	  </div>
	  <div class="form-group">
		<label for="ReportRemark">备注:</label>
		<input type="text" class="form-control" name="ReportRemark" id="ReportRemark" value="">
	  </div>
	  <div class="form-group" style="display:none;">
		<label for="ReportDistId">分布式ID号:</label>
		<input type="number" class="form-control" name="ReportDistId" id="ReportDistId" value="`+MICENGINEID+`" disabled>
	  </div>
	  <div class="form-group">
		<label for="ReportSeq">排序号:</label>
		<input type="number" class="form-control" name="ReportSeq" id="ReportSeq" value="0">
	  </div>
	  <div class="row col-12 container"><div class="form-check-inline">
		<div class="form-group col border rounded" style="margin-right:5px">
			<label for="IsForder">是否模板:</label><br/>	
			<label class="radio-inline"><input type="radio" name="IsForder" id="Forder_N" onclick="onIsForder(0);" value="0" checked>是</label>
			<label class="radio-inline"><input type="radio" name="IsForder" id="Forder_Y" onclick="onIsForder(1);" value="1" >否</label>
		</div>
		<div class="form-group col border rounded">
			<label for="IsDebug">调试模式:</label><br/>	
			<label class="radio-inline"><input type="radio" name="IsDebug" id="Debug_Y" value="1" >是</label>
			<label class="radio-inline"><input type="radio" name="IsDebug" id="Debug_N" value="0" checked>否</label>
		</div>
	  </div></div>
	  </div>
	</div>
	<div class="col-6">
	  <strong>时间信息</strong>
	  <div class="col-sm-12 border rounded" style="margin-top:3px;padding-top:3px">
	  <div class="form-group">
		<label for="ReportBaseTime">基准时间:</label>
		<input type="datetime-local" class="form-control" name="ReportBaseTime" id="ReportBaseTime" value="">
	  </div>
	  <div class="form-group">
		<label for="ReportShifthour">每班工作时间(小时):</label>
		<input type="number" class="form-control" name="ReportShifthour" id="ReportShifthour" value="8">
	  </div>
	  <div class="form-group" id="SetPeriod">
		<label for="ReportPeriod">计算周期:</label>
			<select class="form-control" name="ReportPeriod" id="ReportPeriod" oninput ="onPeriod(this.options[this.options.selectedIndex].value)">
					<option value="0">自定义</option>
					<option value="-1">每小时</option>
					<option value="-2">每班</option>
					<option value="-3" selected="selected">每日</option>
					<option value="-4" >每月</option>
					<option value="-5" >每季度</option>
					<option value="-6" >每年</option>
			</select>
	  </div>
	  <div class="form-group">
		<label for="ReportStartTime">起始时间:</label>
		<input type="datetime-local" class="form-control" name="ReportStartTime" id="ReportStartTime" value="`+nowstring+`">
	  </div>
	  <div class="form-group">
		<label for="ReportLastTime">最后计算截至时间:</label>
		<input type="datetime-local" class="form-control" name="ReportLastTime" id="ReportLastTime" value="`+ nowstring +`">
	  </div>
	  <div class="form-group">
		<label for="ReportOffsetMinute">偏移时间(分钟):</label>
		<input type="number" class="form-control" name="ReportOffsetMinute" id="ReportOffsetMinute" value="0">
	  </div>
	  </div>
	</div>
  </div>

	<div class="row" style="display:none"><div class="col">
	<strong>路径信息</strong>
	<div class="col-sm-12 border rounded" style="margin-top:3px;padding-top:3px">
	<div class="form-group">
	  <label for="ReportTplUrl">模板路径(默认路径,不建议更改):</label>
	  <input type="text" class="form-control" name="ReportTplUrl" id="ReportTplUrl" value="data/report/template/">
	</div>
	<div class="form-group">
	  <label for="ReportResultUrl">结果存储路径(默认路径,不建议更改):</label>
	  <input type="text" class="form-control" name="ReportResultUrl" id="ReportResultUrl" value="data/report/form/">
	</div>
	</div>
	</div></div>
</div></form>
<div class="row" style="margin-top:3px"><div class="col-4"></div><div class="col-4  btn-group"> <button class="btn btn-success" onclick="onSubmitNew()">提交</button></div><div class="col-4"></div></div>
	`;
	ShowModal("添加层级",modaltext);
}
//是否模板被选择时
function onIsForder(value){
	if(value==1){
		$("#ReportDistId").val(0);
	}else{
		$("#ReportDistId").val(MICENGINEID);
	}
}

function onSubmitNew(){
	var msg={
		DistributedId : $("#ReportDistId").val(),     //分布式计算ID，为0的时候不区分计算引擎，不为0的时候需要id适配的计算引擎进行计算
		Name          : $("#ReportLeveltName").val(), //名称
		Pid           : $("#ReportParent").val(),     //父级菜单ID
		WorkshopId    : $("#ReportWorkShop").val(),   //所属车间ID
		Level         : 0 ,//层级深度
		LevelCode     : "" ,//层级码
		Folder        : $("input[name='IsForder']:checked").val(),//$("#IsForder").val(),     //是否是文件夹，1-是，0-否
		Debug		  : $("input[name='IsDebug']:checked").val(),//$("#IsForder").val(),     //是否是文件夹，1-是，0-否
		Seq           : $("#ReportSeq").val(),    //排序号
		Remark        : $("#ReportRemark").val(), //备注
		TemplateUrl   : $("#ReportTplUrl").val(), //模板文件路径
		TemplateFile  : "" ,//模板文件名称
		ResultUrl     : $("#ReportResultUrl").val(),    //结果地址
		StartTime     : $("#ReportStartTime").val(),    //统计计算开始起作用的时间
		Period        : $("#ReportPeriod").val(),      //计算周期,详见KPI表
		OffsetMinutes : $("#ReportOffsetMinute").val(), //偏移时间
		LastCalcTime  : $("#ReportLastTime").val(),     //最后计算时间
		BaseTime      : $("#ReportBaseTime").val(),     //基准时间
		ShiftHour     : $("#ReportShifthour").val(),    //每班工作时间
		Status        : 1, //1有效 0无效
	};
	//console.log(msg);
	HideModal();
	postMsg("/api/addreportlevel",msg);
}

//上传模板
function onUpLoadTpl(){
	var htmlstr=`
	<div class="row"><div class="col-12">
	<form action="" role="form" id="Uploadfile">
	<div class="custom-file">
	<input type="file" id="file" name ="file" accept="application/vnd.openxmlformats-officedocument.spreadsheetml.sheet" onchange="onFileSelected();">
	</div>
	<div class="custom-control custom-checkbox">
    <input type="checkbox" class="custom-control-input" id="defaulttpl" name="defaulttpl" checked>
    <label class="custom-control-label" for="defaulttpl" >作为默认模板</label>
  	</div>
	</form></div></div>
	<div class="alert alert-warning">注意：仅支持“.xlsx”类型的Excel文件!</div>
	<div class="row"><div class="col-4"></div>
	<div class="col-4  btn-group"> <button class="btn btn-success" onclick="onSubmitUpload()" id="Uploadbtn" disabled>上传</button></div>
	<div class="col-4"></div></div>`;
	ShowModal("上传模板文件",htmlstr);
}
function onFileSelected(){
	$("#Uploadbtn").removeAttr("disabled");
}
function onSubmitUpload(){
	HideModal();
	var formData = new FormData(document.querySelector("#Uploadfile"));

	formData.append("id",NODE_ID);
	formData.append("filepath",NODES_KVA[NODE_ID].TemplateUrl);
	formData.append("reportname",NODES_KVA[NODE_ID].Name);

	var xhr=new XMLHttpRequest();
    xhr.open("post","api/uploadrpttpl");
    xhr.send(formData);
    xhr.onload=function(){
        if(xhr.status==200){
			alert(xhr.responseText);
			requestTplLists(NODE_ID); 
        }
    }
}
$("#modal-btn").text('退出');
$(document).ready(function() {
	//Ztree国际化
	$("#exit").after('退出');
	$("#ExpandTreeNode").html('展开');
	$("#CollapseTreeNode").html('折叠');
	$("#HideTreeNode").html('隐藏');
	$("#SearchTreeNode").attr("placeholder",'搜索');

	$('#ExpandTreeNode').attr("title","未选中节点时展开所有节点,选中节点时展开选中节点");
	$('#CollapseTreeNode').attr("title","未选中节点时折叠所有节点,选中节点时折叠选中节点");

	$("#NewLevel").html('添加层级');
	$("#UploadTpl").html('上传模板');
	$('#NewLevel').attr("title","添加一个新的层级!");
	$('#UploadTpl').attr("title","为所选报表上传一个模板文件!");
	$('#Atention').html('<strong>注意:</strong>在线预览的格式与实际Excel文件格式并非完全一致，请以下载后的Excel文件为准！<br>在线编辑仅支持单元格中的文本和公式的编辑，不支持格式的编辑，要编辑文本和表格格式请离线编辑！');
	$('#CellAxis').attr("title","选中单元格的坐标");
	$('#AddFunc').attr("title","公式编辑器暂未开通,敬请期待！");
	$('#selectTdValue').attr("title","选中单元格的值或者公式,可编辑,编辑后自动保存！");
	$("#CheckCellValue").html('检验');


});
