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
<div class="col-sm-9 col-md-6 col-lg-8 col-xl-10" id="FloatArea">
<!--====================================================================================================-->
<div class="container col-12 border rounded" style="padding-top:3px">	<!-----------数据区域------------------->
	<div class="col-12 form-inline border rounded" style="padding:3px">
		<button class="btn btn-outline-primary btn" id="NewLevel" onclick="onNewLevel();" title="添加一个新的层级!">添加层级</button>
		<button class="btn btn-outline-success btn" id="UploadTpl" onclick="onUpLoadTpl();" title="为所选报表上传一个模板文件!" disabled>上传模板</button>
	</div>
	<div class="col-12 border rounded" id="DataFrame" style="padding:3px"></div>
	<div class="col-12 border rounded" id="TplLists" style="padding:3px"></div>
	<div class="col-12 border rounded" style="padding:3px">
		<div class="alert alert-warning" style="display:none" id="Atention"><strong>注意:</strong>在线预览的格式与实际Excel文件格式并非完全一致，请以下载后的Excel文件为准！<br>在线编辑仅支持单元格中的文本和公式的编辑，不支持格式的编辑，要编辑文本和表格格式请离线编辑！</div>
		<div style="width: 100%;display: none;padding:0px;" id="EditForm">
			<input type="text" style="width: 60px; text-align: center;" id="CellAxis" disabled title="选中单元格的坐标">
			<button class="btn btn-outline-primary btn-sm" id="AddFunc" onclick="onAddFunc();" disabled title="公式编辑器暂未开通,敬请期待！"><img src="../static/img/func.svg" style="width:25px;"></button>
			<input type="text" id="selectTdValue" style="width: 90%;outline:none;" onchange="onCellValueChange()" oninput="onCellValueInput()" title="选中单元格的值或者公式,可编辑,编辑后自动保存">
			<button class="btn btn-outline-success btn-sm" style="width:60px;" id="CheckCellValue" onclick="onCheckCellValue();" >校验</button>
		</div>
		<div class="row" id = "ExcelFile" style="display: inline-block;">
			<div class="col-12" id="FileVew" style="padding:1px;display:none"></div>
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
<script type="text/javascript" src="static/js/page/reportedit_advance.js"></script>  
<!-- ===================内部JavaScript部分 =================== -->
<!-- ---------------加载结构树数据----------------- -->
{{str2html .TreeNodes}}	
{{template "pageInclude/treeseting.tpl" .RootPid}}
</body>
</html>


