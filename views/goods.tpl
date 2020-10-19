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
		<a class="btn btn-outline-primary" role="button"  onclick="cookieSwitch()" id = "language" ><img src="../static/img/menuico/lang.svg" style="width:25px;"></a>
	    <a class="btn btn-outline-success" role="button" href="/usercenter"><img src="../static/img/menuico/man.svg" style="width:25px;">{{.UserName}}</a>
	    <a class="btn btn-outline-danger"  role="button"  href="/logout/"> <img src="../static/img/menuico/quit.svg"  style="width:25px;"  id = "exit"></a>
  </div>
</nav><!-- /导航栏部分 -->
<!-- ===================页面主体=================== -->	
<div class="container-fluid" style="margin-top:0px">
<div class="row">
	<!-- 左侧树形结构 -->
    <div class="col-sm-3 col-md-6 col-lg-4 col-xl-2  border rounded" id="TreeBar">
	  <div id="TreeButton" class="row border rounded">
		<div class="col-12 btn-group btn-group-sm">
			<button type="button" class="btn btn-outline-primary" id="ExpandTreeNode"></button>
			<button type="button" class="btn btn-outline-info" id="CollapseTreeNode"></button>
			<button type="button" class="btn btn-outline-warning" id="HideTreeNode" onclick="onTreeNodeView(0)"  ></button>
		</div>
		<div  class="col-12">
			<input class="col-12" type="text" id="SearchTreeNode" placeholder="搜索">
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
			<label class="font-weight-bold" for="BeginTime" id = "BeginTimes" ></label>
			<input class="form-control" type="datetime-local" id="BeginTime" value="{{.BeginTime}}" onchange="onTimeChange()" onfocus="onTimeFocus(this.id)">
		</div>
		<div class="col-sm-4 form-inline">
			<label class="font-weight-bold" for="EndTime" id = "EndTimes"></label>
			<input class="form-control" type="datetime-local" id="EndTime" value="{{.EndTime}}" onchange="onTimeChange()" onfocus="onTimeFocus(this.id)">
		</div>
		<div class="col-sm-3 form-inline">
			<label class="font-weight-bold" for="TimeRange" id = "TimeRanges"></label>
			<select class="form-control" id="TimeRange" onchange="onTimeRangeChange(this.options[this.options.selectedIndex].value)" >
					<option value="0" disabled="disabled" id = "custom"></option>
					<option value="480" id = "8h"></option>
					<option value="720" id = "12h"></option>
					<option value="1440" id = "24h"></option>
					<option value="10080" selected="selected" id = "7d"></option>
					<option id="op_this_shift" value="1"></option>
					<option id="op_this_day" value="2"></option>
					<option id="op_this_month" value="3"></option>
			</select>
		</div>
		<div class="col-sm-1 text-right form-inline btn-group">
			<button class="btn btn-outline-primary btn-sm" id="Last" onclick="onLast();"> <</button>
			<button class="btn btn-outline-primary btn-sm" id="Next" onclick="onNext();"> ></button>
		</div>
		</div>
	</div><!--/时间选择区-->
	</div><!--/变量基本信息和选择按钮-->
</div>
<!--====================================================================================================-->
<div class="container col-12 border rounded">	<!-----------数据区域------------------->
    <!--------变量信息-------->
	<div class="container col-12 border rounded">
		<div class="row border rounded" id="CfgMsg" style="margin:3px"></div><!-->配置信息<-->
		<div class="row border rounded" id="DataMsg" style="margin:3px;padding:5px"></div><!-->数据区域<-->
		<div class="row border rounded" id="PageEcharts" style="height: 300px;display: ;border: 1px solid #cecece;"></div><!-->样本化验指标结果图<-->
	</div> 
	<!--------/变量信息-------->
    <div class="alert alert-success" id="OperationMsg">
    </div>

	<div id="TestMsg"></div>
	<div id="PlatPicPath" style="display:none;">{{.PlatPicPath}}</div>
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
