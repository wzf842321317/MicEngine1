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
	  <a class="btn btn-outline-danger" role="button"  href="/logout/"><img src="../static/img/menuico/quit.svg" style="width:25px;">退出</a>
  </div>
</nav><!-- /导航栏部分 -->
<!-- ===================页面主体=================== -->	
<div class="container-fluid" style="margin-top:0px">
<div class="row">
	<!-- 左侧树形结构 -->
    <div class="col-sm-3 col-md-6 col-lg-4 col-xl-2  border rounded" id="TreeBar">
	  <div id="TreeButton" class="row border rounded">
		<div class="col-12 btn-group btn-group-sm">
			<button type="button" class="btn btn-outline-primary" id="ExpandTreeNode" title="未选中节点时展开所有节点,选中节点时展开选中节点">展开</button>
			<button type="button" class="btn btn-outline-info" id="CollapseTreeNode" title="未选中节点时折叠所有节点,选中节点时折叠选中节点">折叠</button>
			<button type="button" class="btn btn-outline-warning" id="HideTreeNode" onclick="onTreeNodeView(0)">隐藏</button>
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
		<button type="button" class="btn btn-outline-primary btn-sm" id="ShowTreeNode" onclick="onTreeNodeView(1)">></button>
	</div>
<!-- ===================右侧内容窗口========================== -->
<div class="col-sm-9 col-md-6 col-lg-8 col-xl-10" id="FloatArea">
<!--====================================================================================================-->
<div class="container col-12 border rounded" style="padding:5px"><!--变量基本信息和选择按钮-->
	<div class="container col-12 border rounded" style="padding-top:3px"><!--变量基本信息和选择按钮-->
	<div class="col-12"><!--时间选择区-->	
		<div class="row">
		<div class="col-sm-4 form-inline">
			<label class="font-weight-bold" for="BeginTime">起始时间：</label>
			<input class="form-control" type="datetime-local" id="BeginTime" value="{{.BeginTime}}" onchange="onTimeChange()" onfocus="onTimeFocus(this.id)">
		</div>
		<div class="col-sm-4 form-inline">
			<label class="font-weight-bold" for="EndTime">结束时间：</label>
			<input class="form-control" type="datetime-local" id="EndTime" value="{{.EndTime}}" onchange="onTimeChange()" onfocus="onTimeFocus(this.id)">
		</div>
		<div class="col-sm-3 form-inline">
			<label class="font-weight-bold" for="TimeRange">时间范围：</label>
			<select class="form-control" id="TimeRange" onchange="onTimeRangeChange(this.options[this.options.selectedIndex].value)" >
					<option value="0" disabled="disabled">自定义</option>
					<option value="480">8小时</option>
					<option value="720">12小时</option>
					<option value="1440">24小时</option>
					<option value="10080" selected="selected">7天</option>
					<option id="op_this_shift" value="1">班</option>
					<option id="op_this_day" value="2">日</option>
					<option id="op_this_month" value="3">月</option>
			</select>
		</div>
		<div class="col-sm-1 text-right form-inline btn-group">
			<button class="btn btn-outline-primary btn-sm" id="Last" onclick="onLast();" title="前一时间段!"> <</button>
			<button class="btn btn-outline-primary btn-sm" id="Next" onclick="onNext();" title="后一时间段!"> ></button>
		</div>
		</div>
	</div><!--/时间选择区-->
	</div><!--/变量基本信息和选择按钮-->
</div>
<!--====================================================================================================-->
<div class="container col-12 border rounded">	<!-----------数据区域------------------->
    <!--------变量信息-------->
	<div class="container col-12 border rounded">
		<div class="row border rounded" id="SampleMsg" style="margin:3px"></div><!-->样本模板信息<-->
		<div class="row border rounded" id="SampleLabTags" style="margin:3px;padding:5px"></div><!-->样本化验指标信息<-->
		<div class="row border rounded" id="SampleLabResult" style="margin:3px;padding:5px"></div><!-->样本化验指标结果信息<-->
		<div class="row border rounded" id="SampleEcharts" style="height: 300px;display: ;border: 1px solid #cecece;"></div><!-->样本化验指标结果图<-->
	</div> 
	<!--------/变量信息-------->
    <div class="alert alert-success" id="OperationMsg">
       <strong>操作：</strong> 请在左侧结构树中选择一个样本节点以便显示数据！
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
