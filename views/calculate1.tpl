<!DOCTYPE html>
<html lang="zh">
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
{{/*{{template "pageInclude/modal.tpl" .ModalSize}}*/}}

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
    </div>
</div>
<form action="/api/calculate" method="post"  target="frameName">
<div class="container-fluid" style="margin-top:0px">
    <div class="row">
        <div class="col-sm-9 col-md-6 col-lg-8 col-xl-10" id="FloatArea">
            <div class="container col-12 border rounded" style="padding:5px"><!--变量基本信息和选择按钮-->
                <div class="container col-12 border rounded" style="padding-top:3px"><!--变量基本信息和选择按钮-->
                    <div class="col-12"><!--时间选择区-->
                        <div class="row" form-inline>
                            <div class="col-sm-4 form-inline">
                                <label class="font-weight-bold" id="BeginTimes" for="BeginTime"></label>
                                <input class="form-control" type="datetime-local" id="BeginTime" value="{{.BeginTime}}" >
                            </div>
                            <div class="col-sm-4 form-inline">
                                <label class="font-weight-bold" id="EndTimes" for="EndTime"></label>
                                <input class="form-control" type="datetime-local" id="EndTime" value="{{.EndTime}}" >
                            </div>
                            <div class="col-sm-4 form-inline">
                                <label class="font-weight-bold" id="tag_id" for="tag_id"></label>
                                 <input class="form-control" type="text" id="tag_id1" name="tag_id">
                            </div>
{{/*                            <div class="col-sm-5 form-inline">*/}}
{{/*                                <label class="font-weight-bold" id="ThelastTimestemp" for="ThelastTimestemp"></label>*/}}
{{/*                                <label for="rd_1">*/}}
{{/*                                    <input class="form-control" type="radio" name="timediff" id="rd_1" ></label>*/}}
{{/*                                <label for="rd_2">*/}}
{{/*                                    <input class="form-control" type="radio" name="timediff" id="rd_2" ></label>*/}}
{{/*                                <label for="rd_3">*/}}
{{/*                                    <input class="form-control" type="radio" name="timediff" id="rd_3" ></label>*/}}
{{/*                                <label for="rd_4">*/}}
{{/*                                    <input class="form-control" type="radio" name="timediff" id="rd_4" ></label>*/}}
{{/*                            </div>*/}}
                            <div class="col-sm-5 form-inline">
                                <label class="font-weight-bold" id="fc" for="fc"></label>
                                <select class="form-control" id="fc1" style="width:223px">
                                    <option value="Total" id="Total"></option>
                                    <option value="plusDiff" id="plusDiff"></option>
                                </select>
                            </div>
                            <div class="btn-group col-sm-4">
                                <button  class=" btn  btn-outline-primary btn-sm-2" onclick="echartsUtls()" id="search"></button>
                                <button  class=" btn  btn-outline-primary btn-sm-2" onclick="submitAction()" id="save"></button>
                            </div>
                            <div class="alert alert-success" id="HisDataTable" ></div>
                    </div><!--/时间选择区-->
                </div><!--/变量基本信息和选择按钮-->
            </div>
        </div><!-- /右侧内容窗口 -->
        </div>
    </div>
    </div>
</form>
<div class="container-fluid" style="margin-top:0px">
    <div id="present" style="height: 300px;border: 2px solid #ccc;"></div><br>
    <div id="chart_bar" style="height: 300px;border: 2px solid #ccc;"></div>
</div>
<iframe name="frameName" scrolling="no" frameborder="0" src=" " style="width:100%;height:100%"></iframe>
<!-- ===================外部Js引用 =================== -->
{{template "pageInclude/scriptlink.tpl"}}
<script type="text/javascript" src="static/js/page/{{.JsFileName}}.js"></script>
<!-- ===================内部JavaScript部分 =================== -->
<!-- ---------------加载结构树数据----------------- -->
{{str2html .TreeNodes}}
{{template "pageInclude/treeseting.tpl" .RootPid}}
</body>
</html>
