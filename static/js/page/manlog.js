//=========变量定义区域===========================================================
var LOGS;//日志信息
var CHANGE_UNAME=false;//更改了用户选项
var SOURCE_ANALYSE;//数据来源信息
var APP_ANALYSE=new Array();//app数据统计信息
var PLAT_ANALYSE=new Array();//平台数据统计信息
var ENG_ANALYSE=new Array();//计算服务数据统计信息
var ECHARTS_SHOW=false;//图形已经显示
//==================动作响应区域==================================================
//页面初始化
function pageInit(){
	var maxwidth=$("#BeginTime").width();
	$("#DescKeyWord").width(maxwidth);
	$("#MsgCntsPerPage").width(maxwidth);
	requestDatas();//请求数据
}

function onEditProjectMsg(){
	modaltext = ``;
	ShowModal("编辑项目信息",modaltext);
}
//每页显示条数改变
function onPerPageChange(value){
	requestDatas();//请求数据
}
//关键词改变
function onDescKeyWordChange(key){
	requestDatas();//请求数据
}
//信息来源发生改变
function onSysTypeChange(value){
	requestDatas();//请求数据
}
//操作类型发生改变
function onOprTypeChange(value){
	requestDatas();//请求数据
}
//日志用户名发生改变
function onLogUserNameChange(userid){
	CHANGE_UNAME=true;
	requestDatas();//请求数据
}
//页数按钮被按下改变
function onPages(pageno){
	$("#PageNo").val(pageno)
	requestDatas();//请求数据
}
//页数输入框发生改变
function onPageChange(pageno){
	requestDatas();//请求数据
}
//统计按钮被按下
function onTongjiLog(){
	if(ECHARTS_SHOW==false){
		requestDatasAnalys();//请求数据
		$("#TongjiLog").attr("class","btn btn-outline-danger btn-sm");//改变按钮颜色
		$("#TongjiLog").text("关闭统计");//改变按钮文字
		ECHARTS_SHOW=true;
	}else{
		$("#DataAnalysMsg").hide();//隐藏图形
		$("#TongjiLog").attr("class","btn btn-outline-primary btn-sm");//改变按钮颜色
		$("#TongjiLog").text("日志统计");//改变按钮文字
		ECHARTS_SHOW=false;
	}
	
}
//=========AJAX请求URL定义区域=======================================================
//请求数据
function requestDatas(){//
	$("#DataAnalysMsg").hide();//隐藏图形
	$("#TongjiLog").attr("class","btn btn-outline-primary btn-sm");//改变按钮颜色
	$("#TongjiLog").text("日志统计");//改变按钮文字
	ECHARTS_SHOW=false;

	var pageno=$("#PageNo").val();
	if(pageno==null){
		pageno = 0;
	}
	var urlstr="api/userlog?begintime="+InputDateTimeToString("YYYY-mm-dd HH:MM:SS",$("#BeginTime").val())+"&endtime="+InputDateTimeToString("YYYY-mm-dd HH:MM:SS",$("#EndTime").val())+"&userid="+$("#LogUserName").val()+"&systype="+$("#SysType").val()+"&oprtype="+$("#OprType").val()+"&pagesize="+$("#MsgCntsPerPage").val()+"&pageno="+(pageno-1)+"&desc="+$("#DescKeyWord").val();
	loadDatas(urlstr);
}
//请求分析数据
function requestDatasAnalys(){//
	var urlstr="api/userloganalys?begintime="+InputDateTimeToString("YYYY-mm-dd HH:MM:SS",$("#BeginTime").val())+"&endtime="+InputDateTimeToString("YYYY-mm-dd HH:MM:SS",$("#EndTime").val())+"&userid="+$("#LogUserName").val()+"&systype="+$("#SysType").val()+"&oprtype="+$("#OprType").val()+"&desc="+$("#DescKeyWord").val();
	console.log(urlstr);
	loadDatasAnalys(urlstr);
}
//=========AJAX函数定义区域=======================================================
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
        }//请求完成后的处理功能结束---------------------------------------
    });
}
function loadDatasAnalys(urlstr)//读取物耗配置信息
{	
	//$("#DataAnalysMsg").html('<div class="alert alert-warning">正在加载数据信息……</div>');
	//调用公用的loadXMLDoc函数
    loadXMLDoc(urlstr,function()
    {
        if (xmlhttp.readyState==4 && xmlhttp.status==200)//请求处理完成，且状态OK
        {//添加请求完成后的处理功能---------------------------------------
			getDatasAnalys(xmlhttp.responseText);//解读数据
			//下一步：		
        }//请求完成后的处理功能结束---------------------------------------
    });
}
//=========AJAX数据接收解析区域====================================================
//获取分析统计数据
function getDatasAnalys(ajaxdata){
	var datas = eval("("+ajaxdata+")"); 
	console.log(datas);
	if(datas!=null){
		SOURCE_ANALYSE=new Array();
		PLAT_ANALYSE=new Array();
		APP_ANALYSE=new Array();
		ENG_ANALYSE=new Array();

		for(var i=0;i<datas.length;i++){
			data=datas[i];
			switch(data.SysType){//来源类型统计
			case 1:
				if (SOURCE_ANALYSE.hasOwnProperty("管理平台")){//如果存在
					SOURCE_ANALYSE["管理平台"]+=data.Cnts;//加一
				}else{//如果不存在
					SOURCE_ANALYSE["管理平台"]=data.Cnts;//赋值
				}
				if (PLAT_ANALYSE.hasOwnProperty(data.ReqUrl)){//如果存在
					PLAT_ANALYSE[data.ReqUrl]+=data.Cnts;//加一
				}else{//如果不存在
					PLAT_ANALYSE[data.ReqUrl]=data.Cnts;//赋值
				}
				break;
			case 2:
				if (SOURCE_ANALYSE.hasOwnProperty("移动端APP")){//如果存在
					SOURCE_ANALYSE["移动端APP"]+=data.Cnts;//加一
				}else{//如果不存在
					SOURCE_ANALYSE["移动端APP"]=data.Cnts;//赋值
				}
				if (APP_ANALYSE.hasOwnProperty(data.ReqUrl)){//如果存在
					APP_ANALYSE[data.ReqUrl]+=data.Cnts;//加一
				}else{//如果不存在
					APP_ANALYSE[data.ReqUrl]=data.Cnts;//赋值
				}
				break;
			case 3:
				if (SOURCE_ANALYSE.hasOwnProperty("分析系统")){//如果存在
					SOURCE_ANALYSE["分析系统"]+=data.Cnts;//加一
				}else{//如果不存在
					SOURCE_ANALYSE["分析系统"]=data.Cnts;//赋值
				}
				if (ENG_ANALYSE.hasOwnProperty(data.ReqUrl)){//如果存在
					ENG_ANALYSE[data.ReqUrl]+=data.Cnts;//加一
				}else{//如果不存在
					ENG_ANALYSE[data.ReqUrl]=data.Cnts;//赋值
				}
				break;
			}
		}
	}
	refreshPageEcharts(0,0);//刷新图形
}
//解析配置信息
function getDatas(ajaxdata){
	var nodes = eval("("+ajaxdata+")"); 
	LOGS = nodes;
	if(nodes.TotalRows==0){
		$("#TongjiLog").attr("disabled","disabled")
	}else{
		$("#TongjiLog").removeAttr("disabled")
	}
	var htmlstr=`<div class="col-12 border-bottom"><h5>日志数据</h5></div>
	<table class="table table-striped table-hover table-sm">
		<thead class="thead-light">
			<tr><th style="width:50px">序号</th><th style="width:180px">时间</th><th  style="width:75px">用户名</th><th style="width:80px">来源</th><th style="width:200px">操作内容</th><th style="width:300px">请求URL</th><th>请求参数</th><th>请求IP</th></tr>
		</thead><tbody>`;
	var logs = nodes.Logs;
	if(logs!=null){
		for(var i=0;i<logs.length;i++){
			var data=logs[i];

			var systype="未定义";
			switch(data.SysType){
			case 1:
				systype="平台";
				break;
			case 2:
				systype="APP";
				break;
			case 3:
				systype="分析";
				break;
			default:
				break;
			}
			var oprtype="未定义";
			switch(data.OprType){
			case 0:
				oprtype="其他";
				break;
			case 1:
				oprtype="添加";
				break;
			case 2:
				oprtype="删除";
				break;
			case 3:
				oprtype="更新";
				break;
			case 4:
				oprtype="查看";
				break;
			case 5:
				oprtype="添加/更新";
				break;
			case 6,7,8:
				oprtype="登录";
				break;
			default:
				break;
			}

			htmlstr+='<tr><td>'+((nodes.PageNo*nodes.PageSize)+i+1)+'</td><td>'+data.StartTime+'</td><td>'+data.User.Name+'</td><td>'+systype+'</td><td>'+data.Description+'</td><td>'+data.ReqUrl+'</td><td>'+data.ReqParams+'</td><td>'+data.RemoteIp+'</td></tr>';
		}
	}
	htmlstr+='</tbody></table>';
	$("#DataMsg").html(htmlstr);
	setUserNameSelector(nodes.UserNames);//设置用户名选择框
	setPageMenus(nodes);//设置页面选择按钮
}
//设置用户名选择框
function setUserNameSelector(users){
	if(CHANGE_UNAME==false){
		$("#LogUserName").empty();
		$("#LogUserName").append('<option value="0">不限</option>'); 
		for(let  k in users){
			var option='<option value="'+k+'">'+users[k]+'</option>';
			$("#LogUserName").append(option); 
		}
	}
	CHANGE_UNAME=false;
}
//设置页面选择按钮
function setPageMenus(logs){
	var htmlstr='<strong>共'+logs.MaxPage+'页</strong>';
	var btstr='';
	for(var i=1;i<=logs.MaxPage;i++){
		if(logs.PageNo<8){
			if(i<=10){
				if(i==logs.PageNo+1){
					btstr +='<input class="form-control" type="number" id="PageNo" name="PageNo" value="'+(i)+'" onchange="onPageChange(this.value)" style="width:80px" min="1" max="'+logs.MaxPage+'"></input>';
				}else{
					btstr +='<boutton class="btn btn-sm" id="page_'+(i)+'" onclick="onPages('+(i)+');">'+(i)+'</boutton>';
				}
			}else{
				btstr+='......';
				break;
			}
		}else{
			if(i<(logs.PageNo - 5)){
				if(btstr.length<1){
					btstr+='......';
					continue;
				}
			}else{
				if(i<=(logs.PageNo + 5)){
					if(i==logs.PageNo+1){
						btstr +='<input class="form-control" type="number" id="PageNo" name="PageNo" value="'+(i)+'" onchange="onPageChange(this.value)" style="width:80px" min="1" max="'+logs.MaxPage+'"></input>';
					}else{
						btstr +='<boutton class="btn btn-sm" id="page_'+(i)+'" onclick="onPages('+(i)+');">'+(i)+'</boutton>';
					}
				}else{
					btstr+='......';
					break;
				}
			}

		}
	}
	htmlstr+=btstr;
	$("#PageMsg").html(htmlstr);
}
//=========时间控件控制函数=======================================================
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

$(document).ready(function(){
	pageInit();
});

$(document).ready(function() {
	//Ztree国际化
	$("#exit").after('退出');
	$("#BeginTimes").html('起始时间：');
	$("#EndTimes").html('结束时间：');
	$("#TimeRange").html('时间范围：');
	$("#custom").html("自定义");
	$("#eight_hours").html('8小时');
	$("#twelve_hours").html('12小时');
	$("#Twenty-four_hours").html('24小时');
	$("#seven_days").html('7天');
	$("#Last").attr("title","前一时间段!");
	$("#Next").attr("title",'后一时间段!');
	$("#keyword").html('&nbsp关 键 词：');
	$("#DescKeyWord").attr("placeholder","请输入操作关键词");
	$("#SysType").html("信息来源：");
	$("#selected").html('不限');
	$("#Management").html('管理平台');
	$("#Mobile").html('移动端APP');
	$("#Analysis").html('分析系统');
	$("#OprType").html("操作类型：");

	$("#selected1").html('不限');
	$("#1").html('添加');
	$("#2").html("删除");
	$("#3").html("更新");
	$("#4").html('查看');
	$("#5").html('添加/更新');
	$("#6").html('登录');
	$("#0").html('其他');
	$("#LogUserName1").html("用户名：");
	$("#selected2").html('不限');
	$("#MsgCntsPerPages").html("每页条数: ");
	$('#MsgCntsPerPage').attr("title","0时所有信息在一页显示");
	$("#TongjiLog").html('日志统计');
	$('#TongjiLog').attr("title","日志信息统计");
});