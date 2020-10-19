//=========变量定义区域===========================================================


//==================动作响应区域==================================================
//页面初始化
function pageInit(){
	
}

function onEditProjectMsg(){
	modaltext = `	
	<form action="/api/editprojectmsg/" method="post">
	<strong>基本信息</strong>
	<div class="col-sm-12 border rounded" style="margin-top:3px;padding-top:3px">
	 
	  <div class="form-group">
		<label for="ProjectName">项目名称:</label>
		<input type="text" class="form-control" name="ProjectName" value="`+$("#_ProjectName").text()+`">
	  </div>
	  <div class="form-group">
		<label for="ProjectLogo">项目Logo:</label>
		<input type="text" class="form-control" name="LogoPath" value="`+$("#_ProjectLogo").text()+`">
	  </div>
	  <div class="form-group">
		<label for="ProjectIcon">项目ICON:</label>
		<input type="text" class="form-control" name="IcoPath" value="`+$("#_ProjectIcon").text()+`">
	  </div>
	  <div class="form-group">
		<label for="PlatPath">平台路径:</label>
		<input type="text" class="form-control" name="PlatPath" value="`+$("#_PlatPath").text()+`">
	  </div>
	  </div><strong>授权信息</strong>
	  <div class="col-sm-12 border rounded" style="margin-top:3px;padding-top:3px">
	  <div class="form-group">
		<label for="AuthCode">项目授权码:</label>
		<input type="text" class="form-control" name="AuthCode" value="`+$("#_AuthCode").text()+`">
	  </div>
	  <div class="form-group">
		<label for="DistributedId">分布式ID:</label>
		<input type="number" class="form-control" name="DistributedId" value="`+$("#_DistributedId").text()+`">
	  </div>
	  <div class="form-group">
		<label for="Description">备注:</label>
		<input type="text" class="form-control" name="Description" value="`+$("#_Description").text()+`">
	  </div>
	  </div><strong>实时数据库信息</strong>
	  <div class="col-sm-12 border rounded" style="margin-top:3px;padding-top:3px">
	 <div class="form-group">
		<label for="RtdbServer">实时数据库地址:</label>
		<input type="text" class="form-control" name="RtdbServer" value="`+$("#_RtdbServer").text()+`">
	  </div>
	  <div class="form-group">
	    <label for="RtdbPort">实时数据库端口:</label>
	    <input type="number" class="form-control" name="RtdbPort" value="`+$("#_RtdbPort").text()+`">
	  </div>
	  <div class="form-group">
	    <label for="RtdbUser">实时数据库用户名:</label>
	    <input type="text" class="form-control" name="RtdbUser" value="`+$("#_RtdbUser").text()+`">
	  </div>
	  <div class="form-group">
	    <label for="RtdbPsw">实时数据库密码:</label>
	    <input type="password" class="form-control" name="RtdbPsw" value="`+$("#_RtdbPsw").text()+`">
	  </div>
	  </div> <strong>计算结果数据库信息</strong>
	  <div class="col-sm-12 border rounded" style="margin-top:3px;padding-top:3px">
	  <div class="form-group">
		<label for="ResultDbServer">计算结果数据库地址:</label>
		<input type="text" class="form-control" name="ResultdbServer" value="`+$("#_ResultDbServer").text()+`">
	  </div>
	  <div class="form-group">
	    <label for="ResultDbPort">计算结果数据库端口:</label>
	    <input type="number" class="form-control" name="ResultdbPort" value="`+$("#_ResultDbPort").text()+`">
	  </div>
	  <div class="form-group">
	    <label for="ResultDbName">计算结果数据库名:</label>
	    <input type="text" class="form-control" name="ResultdbDbname" value="`+$("#_ResultDbName").text()+`">
	  </div>
	  <div class="form-group">
	    <label for="ResultDbTableName">计算结果数据表名:</label>
	    <input type="text" class="form-control" name="ResultdbTbname" value="`+$("#_ResultDbTableName").text()+`">
	  </div>
	  <div class="form-group">
	    <label for="ResultDbUser">计算结果数据库用户名:</label>
	    <input type="text" class="form-control" name="ResultdbUser" value="`+$("#_ResultDbUser").text()+`">
	  </div>
	  <div class="form-group">
	    <label for="ResultDbPsw">计算结果数据库密码:</label>
	    <input type="password" class="form-control" name="ResultdbPsw" value="`+$("#_ResultDbPsw").text()+`">
	  </div></div>
	  <button type="submit" class="btn btn-success">提交</button></form>

	  
	
	`;
	$("#modal-btn").text('退出');
	ShowModal("编辑项目信息",modaltext);
}
//=========AJAX函数定义区域=======================================================


$(document).ready(function() {
	//Ztree国际化
	$("#exit").after('退出');
	$("#edit").html('编辑信息');
	$("#pm").html('<strong>项目管理</strong>');
	$("#ProjectTable").html('项目基本信息表');
	$("#t1").html('主程序版本');
	$("#t2").html('总运行时间');
	$("#t3").html('本次运行时间');
	$("#t4").html('项目名称');
	$("#t5").html('项目Logo');
	$("#t6").html('项目ICON');
	$("#t7").html('版权单位');
	$("#t8").html('平台路径');


	$("#r1").html('项目核心信息表');
	$("#r2").html('项目识别码');
	$("#r3").html('KPI指标计算授权数');
	$("#r4").html('报表模板授权数');
	$("#r5").html('分布式计算ID');
	$("#r6").html('实时数据库信息:');
	$("#r7").html('IP地址');
	$("#r8").html('数据库名');
	$("#r9").html('用户名');
	$("#r10").html('计算结果数据库信息:');
	$("#r11").html('IP地址');
	$("#r12").html('数据库名');
	$("#r13").html('用户名');
	$("#r14").html('项目授权码');
	$("#r15").html('快速计算授权数');
	$("#r16").html('备注');
	$("#r17").html('端口号');
	$("#r18").html('表名');
	$("#r19").html('密码');
	$("#r20").html('端口号');
	$("#r21").html('表名');
	$("#r22").html('密码');
	$('#dog_1').html('计算服务监控状态');

	if (dog_status == false){
		$('#dog_state').html( '停止');
		$('#dog_state').attr("class","badge badge-pill badge-warning");
    }else{
		$('#dog_state').html(' 运行');
		$('#dog_state').attr("class","badge badge-success");
	}
});
