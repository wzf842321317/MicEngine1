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
        <!-- ===================右侧内容窗口========================== -->
            <!--====================================================================================================-->
            <div class="container-fluid col-12 border rounded" style="padding:5px"><!--变量基本信息和选择按钮-->
                <div class="container col-12 border rounded" style="padding-top:3px"><!--变量基本信息和选择按钮-->
                    <div class="col-12"><!--时间选择区-->
                        <form action="/api/calculate" method="post"  target="frameName">
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
                                <label class="font-weight-bold" id="tag_id" for="tag_id"></label>
                                 <select class="form-control" id="tag_id1" style="width:205px">
                                </select>
                            </div>
                            <div class="col-sm-4 form-inline">
                                <label class="font-weight-bold" id="fc" for="fc"></label>
                                <select class="form-control" id="fc1" style="width:225px">
                                    <option value="Total" id="Total"></option>
                                    <option value="plusDiff" id="plusDiff"></option>
                                </select>
                            </div>
                            <div class="col-sm-4 form-inline">
                            <label class="font-weight-bold" id="cj1" for="cj"></label>
                            <select class="form-control" id="cj" style="width:227px">
                            </select>
                            </div>
                            <div class="col-sm-4 form-inline">
                            <label class="font-weight-bold" id="sb1" for="sb"></label>
                            <select class="form-control" id="sb" style="width:205px">
                            </select>
                            </div>
                            <div class="col-sm-4 form-inline">
                            <label class="font-weight-bold" id="bl1" for="bl"></label>
                            <select class="form-control" id="bl" style="width:223px">
                            </select>
                            </div>
                            <div class="btn-group col-sm-2">
                                <button  class=" btn  btn-outline-primary col-sm--1" onclick="echartsUtls()" id="search" ></button>
                                <button  class=" btn  btn-outline-primary col-sm--1" onclick="submitAction()" id="save" ></button>
                            </div>
                        </div><!--/时间选择区-->
                            <div class="alert alert-success" id="HisDataTable" ></div>
                            <div class="container col-12 border rounded" style="margin-top:3px;padding-top:3px;padding-bottom:3px">	<!-----------数据区域------------------->
                                <div id="present" style="height: 300px;border: 2px solid #ccc;"></div>
                                <div id="chart_bar" style="height: 300px;border: 2px solid #ccc;"></div>
                            </div><!-----------/数据区域------------------->
                        </form>
                    </div><!--/变量基本信息和选择按钮-->
                </div>
            </div><!-- /右侧内容窗口 -->
        </div>
</div><!-- ===================/页面主体=================== -->
<iframe name="frameName" scrolling="no" frameborder="0" src=" " style="width:100%;height:100%"></iframe>
<!-- ===================外部Js引用 =================== -->
{{template "pageInclude/scriptlink.tpl"}}
<script type="text/javascript" src="static/js/page/{{.JsFileName}}.js"></script>
<!-- ===================内部JavaScript部分 =================== -->
<!-- ---------------加载结构树数据----------------- -->
</body>
</html>
