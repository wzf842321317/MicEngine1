<!DOCTYPE html>
<html lang="zh">
<html>
<head>
	<!-- =============CSS等引用 =================== -->
	{{template "pageInclude/headlink.tpl"}}

	<title>{{.WebTitle}}</title>
	<link rel="Bookmark" href="{{.IcoPic}}" />
	<link rel="shortcut icon" href="{{.IcoPic}}" />
</head>
<body>
<!-- =============预定义的模态框 =================== -->
{{template "pageInclude/modal.tpl" .ModalSize}}

<!-- =============导航栏部分============= -->
<nav class="navbar navbar-expand-sm bg-light navbar-light">
	<!--Logo区域-->
	<a class="navbar-brand" href="/"><img src="{{.LogoPic}}" alt= "Logo" style="width:40px;"></a>
	<!-- 小屏幕时折叠导航栏属性 -->
	<button class="navbar-toggler navbar-toggler-right" type="button" data-toggle="collapse" data-target="#collapsibleNavbar">
		<span class="navbar-toggler-icon"></span>
	</button>

	<!-- Links -->
	<div class="collapse navbar-collapse" id="collapsibleNavbar">
		<ul class="nav nav-tabs">
			{{str2html .Navs}}
		</ul>
	</div>
	<script type="text/javascript" src="static/js/page/lang.js"></script>
	<div class="btn-group btn-group-x1" role="group" aria-label="...">
		<a class="btn btn-outline-primary" role="button"  onclick="cookieSwitch()" id = "language" ><img src="../static/img/menuico/lang.svg" style="width:25px;" ></a>
		<a class="btn btn-outline-success" role="button" href="/usercenter"><img src="../static/img/menuico/man.svg" style="width:25px;">{{.UserName}}</a>
		<a class="btn btn-outline-danger" role="button"  href="/logout/"><img src="../static/img/menuico/quit.svg" style="width:25px;" id="exit"></a>
	</div>
</nav><!-- /导航栏部分 -->
<!-- ===================页面主体=================== -->
<div class="container-fluid" style="margin-top:0px">
	<!-- ===================主窗口========================== -->
	<div class="row">
		<!-- ===================右侧内容窗口========================== -->
		<div class="col-12" id="FloatArea">
			<!--====================================================================================================-->
			<div class="col-sm-12 border rounded" id="MineMsg" style="margin-top:3px;padding:3px">
				<div class="col-12"><!--时间选择区-->
					<div class="row">
						<div class="col-sm-3 form-inline">
							<label class="font-weight-bold" for="BeginTime" id="BeginTimes">起始时间：</label>
							<input class="form-control" type="datetime-local" id="BeginTime" value="{{.BeginTime}}" onchange="onTimeChange()" onfocus="onTimeFocus(this.id)">
						</div>
						<div class="col-sm-3 form-inline">
							<label class="font-weight-bold" for="EndTime" id="EndTimes">结束时间：</label>
							<input class="form-control" type="datetime-local" id="EndTime" value="{{.EndTime}}" onchange="onTimeChange()" onfocus="onTimeFocus(this.id)">
						</div>
						<div class="col-sm-3 form-inline">
							<label class="font-weight-bold" for="TimeRange" id="TimeRange">时间范围：</label>
							<select class="form-control" id="TimeRange" onchange="onTimeRangeChange(this.options[this.options.selectedIndex].value)" >
								<option value="0" disabled="disabled" id="custom">自定义</option>
								<option value="480" id="eight_hours">8小时</option>
								<option value="720" id="twelve_hours">12小时</option>
								<option value="1440" selected="selected" id="Twenty-four_hours">24小时</option>
								<option value="10080" id="seven_days">7天</option>
							</select>
						</div>
						<div class="col-sm-3 text-right form-inline btn-group">
							<button class="btn btn-outline-primary btn-sm" id="Last" onclick="onLast();" title="前一时间段!"> <</button>
							<button class="btn btn-outline-primary btn-sm" id="Next" onclick="onNext();" title="后一时间段!"> ></button>
						</div>
					</div>
				</div><!--/时间选择区-->

				<div class="col-12 form-inline" style="padding:3px"><!--筛选条件区-->
					<div class="col-sm-3 form-inline">
						<label class="font-weight-bold" for="DescKeyWord" id="keyword">&nbsp关 键 词：</label>
						<input class="form-control" type="text" id="DescKeyWord" name="DescKeyWord" placeholder="请输入操作关键词" oninput ="onDescKeyWordChange(this.value);"></input>
					</div>
					<div class="col-sm-3 form-inline">
						<label class="font-weight-bold" for="SysType" id="SysType">信息来源：</label>
						<select class="form-control" id="SysType" onchange="onSysTypeChange(this.options[this.options.selectedIndex].value)" >
							<option value="0" selected="selected" id="selected">不限</option>
							<option value="1" id="Management">管理平台</option>
							<option value="2" id="Mobile">移动端APP</option>
							<option value="3" id="Analysis">分析系统</option>
						</select>
					</div>
					<div class="col-sm-3 form-inline">
						<label class="font-weight-bold" for="OprType" id="OprType">操作类型：</label>
						<select class="form-control" id="OprType" onchange="onOprTypeChange(this.options[this.options.selectedIndex].value)" >
							<option value="-1" selected="selected" id="selected1">不限</option>
							<option value="1" id="1">添加</option>
							<option value="2" id="2">删除</option>
							<option value="3" id="3">更新</option>
							<option value="4" id="4">查看</option>
							<option value="5" id="5">添加/更新</option>
							<option value="6" id="6">登录</option>
							<option value="0" id="0">其他</option>
						</select>
					</div>
					<div class="col-sm-3 form-inline" id="LogUserList">
						<label class="font-weight-bold" for="LogUserName" id="LogUserName1">用户名：</label>
						<select class="form-control" id="LogUserName" onchange="onLogUserNameChange(this.options[this.options.selectedIndex].value)" >
							<option value="0" selected="selected" id="selected2">不限</option>
						</select>
					</div>
				</div><!--/筛选条件区-->
				<div class="col-12 form-inline" style="margin-left:px;padding-left:0px;">
					<div class="col-8" style="margin-left:px;padding-left:0px;">
						<form class="form-inline" role="form">
							<div class="col-sm-6 col-lg-6">
								<div class="form-group">
									<label class="font-weight-bold" for="MsgCntsPerPage" id="MsgCntsPerPages">每页条数：</label>
									<div class="form-group">
										<input class="form-control" type="number" id="MsgCntsPerPage" name="MsgCntsPerPage" min="0" value="50" onblur="onPerPageChange(this.value);" title="0时所有信息在一页显示"></input>
									</div>
								</div>
							</div>
							<div class="col-sm-6 col-lg-6"><div class="form-group">
									<div id="PageMsg">
										<input class="form-control" type="number" id="PageNo" name="PageNo" value="1" onchange="onPageChange(this.value-1)" style="width:50px"></input>
									</div>

								</div></div>
						</form>
					</div>
					<div class="col-4">
						<button class="btn btn-outline-primary btn-sm" id="TongjiLog" onclick="onTongjiLog();" title="日志信息统计!" id="TongjiLog">日志统计</button>
					</div>
				</div>
			</div>
			<div class="col-sm-12 border rounded" id="DataAnalysMsg" style="height:520px;width:100%;margin-top:3px;padding:3px;display:noen">	</div>
			<div class="col-sm-12 border rounded" id="DataMsg" style="margin-top:3px;padding:3px">	</div>
			<!--====================================================================================================-->
		</div>
	</div>
	<!-- ===================/主窗口========================== -->
</div>
<!-- ===================/页面主体=================== -->
<!-- ===================外部Js引用 =================== -->
{{template "pageInclude/scriptlink.tpl"}}
<script type="text/javascript" src="static/js/page/{{.JsFileName}}.js"></script>
<script type="text/javascript" src="static/js/page/{{.JsFileName}}_echarts.js"></script>
<!-- ===================内部JavaScript部分 =================== -->
<!-- ---------------加载结构树数据----------------- -->
</body>
</html>
