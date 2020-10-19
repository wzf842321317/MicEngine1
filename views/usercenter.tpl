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
	<!-- 左侧树形结构 -->
    <div class="col-sm-3 col-md-6 col-lg-4 col-xl-2">
      <ul class="list-group">
		<li class="list-group-item"><button type="button" class="btn btn-primary" style="width:180px" onclick="onNewPsw()" id="editPd"></button></li>
		<li class="list-group-item"><button type="button" class="btn btn-primary" style="width:180px" onclick="onEditUserMsg()" id="editMs"></button></li>
	  </ul>
    </div>
	<!-- /左侧树形结构 -->
	<!-- ===================右侧内容窗口========================== -->
	<div class="col-sm-9 col-md-6 col-lg-8 col-xl-10">
	<!--====================================================================================================-->
		<div class="col-sm-12 form-inline border rounded" id="MineMsg" style="margin-top:3px;padding-top:3px">	
		<table class="table table-striped table-hover" style="width:60%"><tbody>
			<tr><th colspan="2" id="userInformation"></th></tr>
			<tr><th id="loginName"></th><td>{{.Username}}</td></tr>
			<tr><th id="name"></th><td>{{.Name}}</td></tr>
			<tr><th id="IDcard"></th><td>{{.CardNo}}</td></tr>
			<tr><th id="phoneNumber"></th><td>{{.Telephone}}</td></tr>
			<tr><th id="email"></th><td>{{.Email}}</td></tr>
			<tr><th id="sex"></th><td>{{.Sex}}</td></tr>
			<tr><th id="qq"></th><td>{{.Qq}}</td></tr>
			<tr><th id="wechat"></th><td>{{.Weixin}}</td></tr>
		</tbody></table>
		</div>
	<!--====================================================================================================-->
	</div>
</div>
<!-- ===================/主窗口========================== -->
</div>
<!-- ===================/页面主体=================== -->
<!-- ===================外部Js引用 =================== -->
{{template "pageInclude/scriptlink.tpl"}}
<script type="text/javascript" src="static/js/page/{{.JsFileName}}.js"></script>  
<!-- ===================内部JavaScript部分 =================== -->
<!-- ---------------加载结构树数据----------------- -->
</body>
</html>
