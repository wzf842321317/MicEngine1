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
	ShowModal("Edit project information",modaltext);
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
		$("#TongjiLog").text("Turn off statistics");//改变按钮文字
		ECHARTS_SHOW=true;
	}else{
		$("#DataAnalysMsg").hide();//隐藏图形
		$("#TongjiLog").attr("class","btn btn-outline-primary btn-sm");//改变按钮颜色
		$("#TongjiLog").text("Log statistics");//改变按钮文字
		ECHARTS_SHOW=false;
	}
	
}
//=========AJAX请求URL定义区域=======================================================
//请求数据
function requestDatas(){//
	$("#DataAnalysMsg").hide();//隐藏图形
	$("#TongjiLog").attr("class","btn btn-outline-primary btn-sm");//改变按钮颜色
	$("#TongjiLog").text("Log statistics");//改变按钮文字
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
	$("#DataMsg").html('<div class="alert alert-warning">Loading data information……</div>');
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
				if (SOURCE_ANALYSE.hasOwnProperty("Management platform")){//如果存在
					SOURCE_ANALYSE["Management platform"]+=data.Cnts;//加一
				}else{//如果不存在
					SOURCE_ANALYSE["Management platform"]=data.Cnts;//赋值
				}
				if (PLAT_ANALYSE.hasOwnProperty(data.ReqUrl)){//如果存在
					PLAT_ANALYSE[data.ReqUrl]+=data.Cnts;//加一
				}else{//如果不存在
					PLAT_ANALYSE[data.ReqUrl]=data.Cnts;//赋值
				}
				break;
			case 2:
				if (SOURCE_ANALYSE.hasOwnProperty("Mobile app")){//如果存在
					SOURCE_ANALYSE["Mobile app"]+=data.Cnts;//加一
				}else{//如果不存在
					SOURCE_ANALYSE["Mobile app"]=data.Cnts;//赋值
				}
				if (APP_ANALYSE.hasOwnProperty(data.ReqUrl)){//如果存在
					APP_ANALYSE[data.ReqUrl]+=data.Cnts;//加一
				}else{//如果不存在
					APP_ANALYSE[data.ReqUrl]=data.Cnts;//赋值
				}
				break;
			case 3:
				if (SOURCE_ANALYSE.hasOwnProperty("Analysis system")){//如果存在
					SOURCE_ANALYSE["Analysis system"]+=data.Cnts;//加一
				}else{//如果不存在
					SOURCE_ANALYSE["Analysis system"]=data.Cnts;//赋值
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
	var htmlstr=`<div class="col-12 border-bottom"><h5>Log data</h5></div>
	<table class="table table-striped table-hover table-sm">
		<thead class="thead-light">
			<tr><th style="width:50px">Number</th><th style="width:180px">Time</th><th  style="width:75px">User</th><th style="width:80px">Source</th><th style="width:200px">Operation content</th><th style="width:300px">Request URL</th><th>Request parameters</th><th>Request IP</th></tr>
		</thead><tbody>`;
	var logs = nodes.Logs;
	if(logs!=null){
		for(var i=0;i<logs.length;i++){
			var data=logs[i];

			var systype="Undefined";
			switch(data.SysType){
			case 1:
				systype="platform";
				break;
			case 2:
				systype="APP";
				break;
			case 3:
				systype="analysis";
				break;
			default:
				break;
			}
			var oprtype="Undefined";
			switch(data.OprType){
			case 0:
				oprtype="Other";
				break;
			case 1:
				oprtype="Add to";
				break;
			case 2:
				oprtype="Delete";
				break;
			case 3:
				oprtype="Update";
				break;
			case 4:
				oprtype="Check";
				break;
			case 5:
				oprtype="Add / update";
				break;
			case 6,7,8:
				oprtype="Login";
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
		$("#LogUserName").append('<option value="0">Unlimited</option>');
		for(let  k in users){
			var option='<option value="'+k+'">'+users[k]+'</option>';
			$("#LogUserName").append(option); 
		}
	}
	CHANGE_UNAME=false;
}
//设置页面选择按钮
function setPageMenus(logs){
	var htmlstr='<strong> '+logs.MaxPage+'page in total </strong>';
	var btstr='';
	for(var i=1;i<=logs.MaxPage;i++){
		if(logs.PageNo<8){
			if(i<=10){
				if(i==logs.PageNo+1){
					btstr +='<input class="form-control" type="number" id="PageNo" name="PageNo" value="'+( i)+'" onchange="onPageChange(this.value)" style="width:80px" min="1" max="'+logs.MaxPage+'"></input>';
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
	$("#exit").after('Exit');
	$("#BeginTimes").html('Begin Time');
	$("#EndTimes").html('End Time');
	$("#TimeRange").html('Time Range');
	$("#custom").html("custom");
	$("#eight_hours").html('8h');
	$("#twelve_hours").html('12h');
	$("#Twenty-four_hours").html('24h');
	$("#seven_days").html('7d');
	$("#Last").attr("title","The previous period!");
	$("#Next").attr("title",'The next time!');
	$("#keyword").html('&nbsp Key word：');
	$("#DescKeyWord").attr("placeholder","Please input keyword");
	$("#SysType").html("Source:");
	$("#selected").html('Unlimited');
	$("#Management").html('Management platform');
	$("#Mobile").html('Mobile app');
	$("#Analysis").html('Analysis system');
	$("#OprType").html("Operation type:");

	$("#selected1").html('Unlimited');
	$("#1").html('Add');
	$("#2").html('Delete');
	$("#3").html('Update');
	$("#4").html('View');
	$("#5").html("Add/Update");
	$("#6").html('Login');
	$("#0").html('Others');
	$("#LogUserName1").html('User name：');
	$("#selected2").html('Unlimited');
	$("#MsgCntsPerPages").html('Number of pages: ');
	$('#MsgCntsPerPage').attr("title",'When 0, all information is displayed on one page');
	$("#TongjiLog").html("Log statistics");
	$('#TongjiLog').attr("title","Log information statistics");
});