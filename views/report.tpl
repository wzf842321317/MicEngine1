<!DOCTYPE html>
<html lang="zh">
<html>
<head>
<!-- =============CSS等引用 =================== -->
{{template "pageInclude/headlink.tpl"}}

  <title>{{.WebTitle}}</title>
  <link rel="Bookmark" href="{{.IcoPic}}" /> 
  <link rel="shortcut icon" href="{{.IcoPic}}" /> 
  <link rel="stylesheet" href="static/css/excel.css" type="text/css"/>
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
	  <a class="btn btn-outline-danger" role="button"  href="/logout/"><img src="../static/img/menuico/quit.svg" style="width:25px;" id = "exit"></a>
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
		<button type="button" class="btn btn-outline-primary btn-sm" id="HideTreeNode" onclick="onTreeNodeView(1)">></button>
	</div>
<!-- ===================右侧内容窗口========================== -->
<div class="col-sm-9 col-md-6 col-lg-8 col-xl-10"  id="FloatArea">
<!--====================================================================================================-->
<div class="container col-12 border rounded" style="padding-top:3px;height=100%">	<!-----------数据区域------------------->
	<div class="col-12 form-inline border rounded" style="padding:3px">
		<div class="col-12"><!--时间选择区-->	
			<div class="row">
			<div class="col-sm-4 form-inline">
				<label class="font-weight-bold" for="BeginTime" id = "BeginTimes">起始时间：</label>
				<input class="form-control" type="datetime-local" id="BeginTime" value="{{.BeginTime}}" onchange="onTimeChange()" onfocus="onTimeFocus(this.id)">
			</div>
			<div class="col-sm-4 form-inline">
				<label class="font-weight-bold" for="EndTime" id = "EndTimes"></label>
				<input class="form-control" type="datetime-local" id="EndTime" value="{{.EndTime}}" onchange="onTimeChange()" onfocus="onTimeFocus(this.id)">
			</div>
			<div class="col-sm-3 form-inline">
				<label class="font-weight-bold" for="TimeRange" id = "TimeRanges">时间范围：</label>
				<select class="form-control" id="TimeRange" onchange="onTimeRangeChange(this.options[this.options.selectedIndex].value)" >
					<option value="0" disabled="disabled" id = "custom"></option>
					<option value="480" id = "8h">8小时</option>
					<option value="720" id = "12h">12小时</option>
					<option value="1440" id = "24h">24小时</option>
					<option value="10080" selected="selected" id = "7d">7天</option>
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
	</div>
	<div class="col-12 border rounded" id="DataFrame" style="padding:3px"></div>
	<div class="col-12 border rounded" id="TplLists" style="padding:3px"></div>
	<div class="col-12 border rounded" style="padding:3px">
		<div class="alert alert-warning" id = "OperationMsg"><strong>注意:</strong>在线预览的格式与实际Excel文件格式并非完全一致，请以下载后的Excel文件为准！在线预览目前只具备简单的Excel公式的计算能力，使用Excel公式的运算结果仅供参考，请以下载后的Excel文件中的结果为准！</div>
		<div style="width: 100%;display: flex;" id="EditForm">
			<input type="text" style="width: 60px; text-align: center;" id="CellAxis" disabled title="选中单元格的坐标">
			<button class="btn btn-outline-primary btn-sm" id="AddFunc" onclick="onAddFunc();" disabled><img src="../static/img/func.svg" style="width:25px;"></button>
			<input type="text" id="selectTdValue" style="width: 90%;outline:none;" disabled title="选中单元格的值或者公式">
			<button class="btn btn-outline-success btn-sm" style="width:60px;" id="CheckCellValue" onclick="onCheckCellValue();" disabled>校验</button>
		</div>
		<div class="row" id = "ExcelFile" style="display: inline-block;">
			<div class="col-12" id="FileVew" style="padding:3px;display:none"></div>
		</div>
	</div>
</div><!-----------/数据区域------------------->
<!--====================================================================================================-->
</div><!-- /右侧内容窗口 -->
</div>
</div><!-- ===================/页面主体=================== -->
<!-- ===================外部Js引用 =================== -->
{{template "pageInclude/scriptlink.tpl"}}
<script type="text/javascript" src="static/js/page/{{.JsFileName}}.js"></script>  
<!-- ===================内部JavaScript部分 =================== -->
<!-- ---------------加载结构树数据----------------- -->
{{str2html .TreeNodes}}	
{{template "pageInclude/treeseting.tpl" .RootPid}}
</body>
</html>
