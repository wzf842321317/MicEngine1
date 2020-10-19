//=========变量定义区域===========================================================


//==================动作响应区域==================================================
//页面初始化
function pageInit(){
	
}

function onEditProjectMsg(){
	modaltext = `	
	<form action="/api/editprojectmsg/" method="post">
	<strong>Essential information</strong>
	<div class="col-sm-12 border rounded" style="margin-top:3px;padding-top:3px">
	 
	  <div class="form-group">
		<label for="ProjectName">Entry name:</label>
		<input type="text" class="form-control" name="ProjectName" value="`+$("#_ProjectName").text()+`">
	  </div>
	  <div class="form-group">
		<label for="ProjectLogo">Project logo:</label>
		<input type="text" class="form-control" name="LogoPath" value="`+$("#_ProjectLogo").text()+`">
	  </div>
	  <div class="form-group">
		<label for="ProjectIcon">Project ICON:</label>
		<input type="text" class="form-control" name="IcoPath" value="`+$("#_ProjectIcon").text()+`">
	  </div>
	  <div class="form-group">
		<label for="PlatPath">Platform path:</label>
		<input type="text" class="form-control" name="PlatPath" value="`+$("#_PlatPath").text()+`">
	  </div>
	  </div><strong>Authorization information</strong>
	  <div class="col-sm-12 border rounded" style="margin-top:3px;padding-top:3px">
	  <div class="form-group">
		<label for="AuthCode">Project authorization code:</label>
		<input type="text" class="form-control" name="AuthCode" value="`+$("#_AuthCode").text()+`">
	  </div>
	  <div class="form-group">
		<label for="DistributedId">Distributed ID:</label>
		<input type="number" class="form-control" name="DistributedId" value="`+$("#_DistributedId").text()+`">
	  </div>
	  <div class="form-group">
		<label for="Description">Remark:</label>
		<input type="text" class="form-control" name="Description" value="`+$("#_Description").text()+`">
	  </div>
	  </div><strong>Real time database information</strong>
	  <div class="col-sm-12 border rounded" style="margin-top:3px;padding-top:3px">
	 <div class="form-group">
		<label for="RtdbServer">Real time database address:</label>
		<input type="text" class="form-control" name="RtdbServer" value="`+$("#_RtdbServer").text()+`">
	  </div>
	  <div class="form-group">
	    <label for="RtdbPort">Real time database port:</label>
	    <input type="number" class="form-control" name="RtdbPort" value="`+$("#_RtdbPort").text()+`">
	  </div>
	  <div class="form-group">
	    <label for="RtdbUser">Real time database user name:</label>
	    <input type="text" class="form-control" name="RtdbUser" value="`+$("#_RtdbUser").text()+`">
	  </div>
	  <div class="form-group">
	    <label for="RtdbPsw">Real time database password:</label>
	    <input type="password" class="form-control" name="RtdbPsw" value="`+$("#_RtdbPsw").text()+`">
	  </div>
	  </div> <strong>Calculation result database information</strong>
	  <div class="col-sm-12 border rounded" style="margin-top:3px;padding-top:3px">
	  <div class="form-group">
		<label for="ResultDbServer">Calculation result database address:</label>
		<input type="text" class="form-control" name="ResultdbServer" value="`+$("#_ResultDbServer").text()+`">
	  </div>
	  <div class="form-group">
	    <label for="ResultDbPort">Calculation result database port:</label>
	    <input type="number" class="form-control" name="ResultdbPort" value="`+$("#_ResultDbPort").text()+`">
	  </div>
	  <div class="form-group">
	    <label for="ResultDbName">Calculation result database name:</label>
	    <input type="text" class="form-control" name="ResultdbDbname" value="`+$("#_ResultDbName").text()+`">
	  </div>
	  <div class="form-group">
	    <label for="ResultDbTableName">Calculation result data table name:</label>
	    <input type="text" class="form-control" name="ResultdbTbname" value="`+$("#_ResultDbTableName").text()+`">
	  </div>
	  <div class="form-group">
	    <label for="ResultDbUser">User name of calculation result database:</label>
	    <input type="text" class="form-control" name="ResultdbUser" value="`+$("#_ResultDbUser").text()+`">
	  </div>
	  <div class="form-group">
	    <label for="ResultDbPsw">Calculation result database password:</label>
	    <input type="password" class="form-control" name="ResultdbPsw" value="`+$("#_ResultDbPsw").text()+`">
	  </div></div>
	  <button type="submit" class="btn btn-success">Submit</button></form>

	  
	`;
	$("#modal-btn").text('close');
	ShowModal("Edit project information",modaltext);
}
//=========AJAX函数定义区域=======================================================


$(document).ready(function() {
	//Ztree国际化
	$("#exit").after('Exit');
	$("#edit").html('edit information');
	$("#pm").html('<strong>Project Management</strong>');
	$("#ProjectTable").html('Project basic information table');
	$("#t1").html('Main program version');
	$("#t2").html('Total running time');
	$("#t3").html('This run time');
	$("#t4").html('Project name');
	$("#t5").html('Project Logo');
	$("#t6").html('Project ICON');
	$("#t7").html('Copyright unit');
	$("#t8").html('Platform path');


	$("#r1").html('Project core information table');
	$("#r2").html('Item identification code');
	$("#r3").html('KPI indicator calculates the number of authorizations');
	$("#r4").html('Number of report template authorizations');
	$("#r5").html('Distributed computing ID');
	$("#r6").html('Real-time database information');
	$("#r7").html('IP address');
	$("#r8").html('database name');
	$("#r9").html('User name');
	$("#r10").html('Calculated results database information');
	$("#r11").html('IP address');
	$("#r12").html('database name');
	$("#r13").html('User name');
	$("#r14").html('Project authorization code');
	$("#r15").html('Quickly calculate the number of authorization');
	$("#r16").html('remark');
	$("#r17").html('port number');
	$("#r18").html('Table Name');
	$("#r19").html('password');
	$("#r20").html('port number');
	$("#r21").html('Table Name');
	$("#r22").html('password');
	$('#dog_1').html('Calculate service monitoring status');
	if (dog_status == false){
		$('#dog_state').html('stop');
		$('#dog_state').attr("class","badge badge-pill badge-warning");
	}else{
		$('#dog_state').html('active');
		$('#dog_state').attr("class","badge badge-success");
	}


});