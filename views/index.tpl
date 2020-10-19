<!DOCTYPE html>
<html lang="zh">
<html>
<head>
  <!-- =============CSS等引用 =================== -->
  {{template "pageInclude/headlink.tpl"}}

  <title>{{.WebTitle}}</title>
  <link rel="Bookmark" href="{{.IcoPic}}" />
  <link rel="shortcut icon" href="{{.IcoPic}}" />

  <style type="text/css">
    *,body {
      margin: 0px;
      padding: 0px;
    }

    body {
      margin: 0px;
      font-family: "Helvetica Neue", Helvetica, Arial, sans-serif;
      font-size: 14px;
      line-height: 20px;
      background-color: #fff;
    }

    header,
    footer {
      width: 960px;
      margin-left: auto;
      margin-right: auto;
    }

    .logo {
      background-image: url("{{.LogoImg}}");
      background-repeat: no-repeat;
      -webkit-background-size: 100px 100px;
      background-size: 100px 100px;
      background-position: center center;
      text-align: center;
      font-size: 42px;
      padding: 250px 0 70px;
      font-weight: normal;
      text-shadow: 0px 1px 2px #ddd;
    }

    header {
      padding: 100px 0;
    }

    footer {
      line-height: 1.8;
      text-align: center;
      padding: 50px 0;
      color: #999;
    }

    .description {
      text-align: center;
      font-size: 16px;
    }
    .version {
      text-align: center;
      padding: 0px 0 30px;
    }

    a {
      color: #444;
      text-decoration: none;
    }

    .backdrop {
      position: absolute;
      width: 100%;
      height: 100%;
      box-shadow: inset 0px 0px 100px #ddd;
      z-index: -1;
      top: 0px;
      left: 0px;
    }
  </style>

</head>
<body>
<!-- =============预定义的模态框 =================== -->
{{template "pageInclude/modal.tpl"}}

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
    <a class="btn btn-outline-danger" role="button"  href="/logout/"><img src="../static/img/menuico/quit.svg" style="width:25px; "id="exit"></a>
  </div>
</nav><!-- /导航栏部分 -->
<!-- ===================页面主体=================== -->
<header>
  <h1 class="logo" >{{.HeaderTitle}}</h1>
  <h2  class="version"  >version:<a href="http://mining-icloud.com/showdoc/web/#/21"> {{.Version}} </a></h2>
  <div class="description">
    {{.Copyright}}
  </div>
</header>
<footer>
  <div class="author">
    Official website:
    <a href="http://{{.Website}}">{{.Website}}</a> /
    Contact me:
    <a class="email" href="mailto:{{.Email}}">{{.Email}}</a>
  </div>
  <div class="api">
    <a id="webApi" href="http://mining-icloud.com/showdoc/web/#/9?page_id=50" ></a> /
    <a id="granafa" href="http://mining-icloud.com/showdoc/web/#/1" ></a>
  </div>
</footer>
<div class="backdrop"></div>
<!-- ===================外部Js引用 =================== -->
{{template "pageInclude/scriptlink.tpl"}}
<script type="text/javascript" src="static/js/page/{{.JsFileName}}.js"></script>
<!-- ===================内部JavaScript部分 =================== -->
<!-- ---------------加载结构树数据----------------- -->

</body>
</html>