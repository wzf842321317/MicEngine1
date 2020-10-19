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
	  <a class="btn btn-outline-danger"  role="button"  href="/logout/"><img src="../static/img/menuico/quit.svg"  style="width:25px;" id="exit" ></a>
  </div>
</nav><!-- /导航栏部分 -->
<!-- ===================页面主体=================== -->	
<div class="container-fluid" style="margin-top:0px">
<div class="row">
	<!-- 左侧树形结构 -->
    <div class="col-sm-3 col-md-6 col-lg-4 col-xl-2  border rounded" id="TreeBar">
	  <div id="TreeButton" class="row border rounded">
		<div class="col-12 btn-group btn-group-sm">
			<button type="button" class="btn btn-outline-success" id="UpdateTree" onclick="onUpdateTree()"></button>
			<button type="button" class="btn btn-outline-primary" id="ExpandTreeNode" ></button>
			<button type="button" class="btn btn-outline-info" id="CollapseTreeNode" ></button>
			<button type="button" class="btn btn-outline-warning" id="HideTreeNode" onclick="onTreeNodeView(0)"></button>
		</div>
		<div  class="col-12">
			<input class="col-12" type="text" id="SearchTreeNode" >
		</div>
	  </div>
	  <div class="row">
		<div id="NodeTree" class="ztree"></div><!--树形结构-->
	  </div>
    </div><!-- /左侧树形结构 -->
	<div id="TreeBarSm" style="display:none">
		<button type="button" class="btn btn-outline-primary btn-sm" id="HideTreeNode" onclick="onTreeNodeView(1)">></button>
	</div>
<!-- ===================右侧内容窗口========================== -->
<div class="col-sm-9 col-md-6 col-lg-8 col-xl-10" id="FloatArea">
<!--====================================================================================================-->
<div class="container col-12 border rounded" style="padding:5px"><!--变量基本信息和选择按钮-->
	<div class="container col-12 border rounded" style="padding-top:3px"><!--变量基本信息和选择按钮-->
	<div class="col-12"><!--时间选择区-->	
		<div class="row">
		<div class="col-sm-4 form-inline">
			<label class="font-weight-bold" id="BeginTimes" for="BeginTime"></label>
			<input class="form-control" type="datetime-local" id="BeginTime" value="{{.BeginTime}}" onchange="onTimeChange()" onfocus="onTimeFocus(this.id)">
		</div>
		<div class="col-sm-4 form-inline">
			<label class="font-weight-bold" id="EndTimes" for="EndTime"></label>
			<input class="form-control" type="datetime-local" id="EndTime" value="{{.EndTime}}" onchange="onTimeChange()" onfocus="onTimeFocus(this.id)">
		</div>
		<div class="col-sm-4 form-inline">
			<label class="font-weight-bold" id="Interval" for="Interval"></label>
			<select class="form-control" id="Interval" onchange="onIntervalTimeChange(this.options[this.options.selectedIndex].value)" style="width:80px">
					<option value="1">1s</option>
					<option value="5">5s</option>
					<option value="10">10s</option>
					<option value="30" selected="selected">30s</option>
					<option value="60">60s</option>
					<option value="120">120s</option>
					<option value="180">180s</option>
					<option value="240">240s</option>
				</select>
		</div>
		</div>
		<div class="row">
		<div class="col-sm-6 form-inline">
			<label class="font-weight-bold" id="ThelastTimestemp" for="ThelastTimestemp"></label>
			<label for="rd_1">
			<input class="form-control" type="radio" name="timediff" id="rd_1" onclick="onTimediffClick(60);"></label>
			<label for="rd_2">
			<input class="form-control" type="radio" name="timediff" id="rd_2" onclick="onTimediffClick(480);"></label>
			<label for="rd_3">
			<input class="form-control" type="radio" name="timediff" id="rd_3" onclick="onTimediffClick(720);"></label>
			<label for="rd_4">
			<input class="form-control" type="radio" name="timediff" id="rd_4" onclick="onTimediffClick(1440);"></label>
		</div>
		
		<!--div class="col-sm-2 form-inline">
			<label class="font-weight-bold" class="form-check-label">
				<input id="show_hisinterval_data" type="checkbox" class="form-check-input" value="" onclick="onShowHisIntervalData(this.id)">显示聚合数据表
			</label>
		</div>
		<div class="col-sm-2 form-inline">
			<label class="font-weight-bold" class="form-check-label">
				<input id="show_his_data" type="checkbox" class="form-check-input" value="" onclick="onShowHisData(this.id)">显示原始数据表
			</label>
		</div -->
   
		</div>
	</div><!--/时间选择区-->
	</div><!--/变量基本信息和选择按钮-->
</div>
<!--====================================================================================================-->
<div class="container col-12 border rounded" style="margin-top:3px;padding-top:3px;padding-bottom:3px">	<!-----------数据区域------------------->
    <!-- -------------变量基本信息------------------ -->
	<div class="col-12 border rounded" style="margin-top:3px;padding-top:3px">
	<div class="row col-sm-12">
		<div class="col-sm-5 form-inline" >
			<label class="font-weight-bold" id="VariableName" class="form-check-label"></label>
			<span id="TagName"></span>
		</div>
		<div class="col-sm-2 form-inline" >
			<label class="font-weight-bold" id="LatestValue" class="form-check-label"></label>
			<span id="TagValue"></span>
		</div>
		<div class="col-sm-2 form-inline" >
			<label class="font-weight-bold" id="unit" class="form-check-label"></label>
			<span id="TagUnit"></span>
		</div>
		<div class="col-sm-3 form-inline" >
			<label class="font-weight-bold" id="UpdateTime" class="form-check-label"></label>
			<span id="TagTime"></span>
		</div>	
	</div>	
	<!-- -------------统计数据表------------------ -->
	<div id="HisSumData">
	<table class="table table-striped table-hover table-sm"><thead class="thead-light"><tr><th id="Min"></th><th id="Max"></th><th id="Range"></th><th id="ArithmeticMean"></th><th id="weightedMean"></th><th id="Mode"></th><th id="Median "></th><th id="Sum"></th><th id="Difference"></th><th id="PositiveDifference"></th><th id="Area "></th><th id="Points "></th><th id="SD"></th><th id="Skewness"></th><th id="Kurtosis"></th></tr></thead><tbody></tbody></table>
	</div><!--/统计数据表-->
	</div><!--/变量基本信息-->
	
	<!-----------所选变量列表------------------->
	<div class="col-12 border rounded" style="margin-top:3px;padding-top:3px">
	<div class="col-12 form-inline border rounded">
		<div class="col-sm-4 form-inline" id="LoadDataMsg"></div>
		<div class="col-sm-4 text-center ">
			<button class="btn btn-outline-primary" id="AddTagToTable" onclick="onAddSelectTagToTable();"></button>
		</div>
		<div class="col-sm-4 text-right form-inline btn-group">
			<button class="btn btn-outline-primary" id="Last" onclick="onLast();" > <<</button>
			<button class="btn btn-outline-primary" id="Next" onclick="onNext();" > >></button>
		</div>
	</div>
	<div id="SelectedTags"><h5><small id="TagsTableSmallTitle" class="text-muted"></small></h5><table class="table table-striped table-hover table-sm"><tr><th id="SerialNumber"></th><th id="Name"></th><th id="tag"></th><th id="type"></th><th id="unit1"></th><th id="removed"></th></tr></thead><tbody></tbody></table></div><!--已选变量列表-->
	</div><!-----------/所选变量列表------------------->
	
	<div id="Echarts_His" style="height: 300px;display: ;border: 1px solid #cecece;margin-top:3px;"></div><!--原始历史趋势-->
	
	<div id="Echarts_HisIntervalSerial" style="height: 500px;display: ;border: 1px solid #cecece;margin-top:3px;"></div><!--等间隔历史趋势序列-->
	<div id="HisSerialRemark" class="col-12 form-inline">
		<span class="alert alert-warning"></span>
	</div><!--备注说明-->
	
	<div >
      <div class="alert alert-success" id="HisDataTable" >
      </div>
  </div><!--原始历史数据表-->
	<div id="HisIntervalTable"></div><!--等间隔历史数据表-->

  <div id="TestMsg"></div>
</div><!-----------/数据区域------------------->
<!--====================================================================================================-->
</div><!-- /右侧内容窗口 -->
</div>
</div><!-- ===================/页面主体=================== -->
<!-- ===================外部Js引用 =================== -->
{{template "pageInclude/scriptlink.tpl"}}
<script type="text/javascript" src="static/js/page/{{.JsFileName}}.js"></script>  
<script type="text/javascript" src="static/js/page/{{.JsFileName}}_echarts.js"></script> 
<!-- ===================内部JavaScript部分 =================== -->
<!-- ---------------加载结构树数据----------------- -->
{{str2html .TreeNodes}}	
{{template "pageInclude/treeseting.tpl" .RootPid}}
</body>
</html>
