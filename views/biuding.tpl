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
  <a  class="btn btn-outline-success"  src="static/js/page/lang.js"  role="button"  onclick="cookieSwitch()" id = "language"><img src="../static/img/menuico/lang.svg" style="width:25px;" ></a>
	  <a class="btn btn-outline-success" role="button" href="/usercenter"><img src="../static/img/menuico/man.svg" style="width:25px;">{{.UserName}}</a>
	  <a class="btn btn-outline-danger" role="button"  href="/logout/"><img src="../static/img/menuico/quit.svg" style="width:25px;">退出</a>
  </div>
</nav><!-- /导航栏部分 -->
<!-- ===================页面主体=================== -->	
<div class="container-fluid" style="margin-top:0px">
	<img src="../static/img/biuding.jfif" style="width:100%;" alt="页面建设中……">
</div><!-- ===================/页面主体=================== -->
<div id="PlatPicPath" style="display:none;">{{.PlatPicPath}}</div>
<!-- ===================外部Js引用 =================== -->
{{template "pageInclude/scriptlink.tpl"}}
<script type="text/javascript" src="static/js/page/{{.JsFileName}}.js"></script>  
<!-- ===================内部JavaScript部分 =================== -->
<!-- ---------------加载结构树数据----------------- -->
{{str2html .TreeNodes}}	
{{template "pageInclude/treeseting.tpl" .RootPid}}
</body>
</html>
