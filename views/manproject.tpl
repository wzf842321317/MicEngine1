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
    <script type="text/javascript">var dog_status={{.DogStatus}}</script>
    <script type="text/javascript" src="static/js/page/lang.js"></script>
    <div class="btn-group btn-group-x1" role="group" aria-label="...">
        <a class="btn btn-outline-primary" role="button"  onclick="cookieSwitch()" id = "language" ><img src="../static/img/menuico/lang.svg" style="width:25px;" ></a>
        <a class="btn btn-outline-success" role="button" href="/usercenter"><img src="../static/img/menuico/man.svg" style="width:25px;">{{.UserName}}</a>
        <a class="btn btn-outline-danger" role="button"  href="/logout/"><img src="../static/img/menuico/quit.svg" style="width:25px;" id = "exit"></a>
    </div>

  <!-- Links -->

</nav><!-- /导航栏部分 -->
<!-- ===================页面主体=================== -->
<div class="container-fluid" style="margin-top:0px">
    <!-- ===================主窗口========================== -->
    <div class="row">
        <!-- 左侧树形结构 -->
        <div class="col-sm-3 col-md-6 col-lg-4 col-xl-2">
            <ul class="list-group">
                <li class="list-group-item"><div class="col-12"><span class="text-center" id = "pm"></span></div></li>
                <li class="list-group-item"><button type="button" class="btn btn-primary" style="width:180px" onclick="onEditProjectMsg()"  id = "edit"></button></li>
            </ul>
        </div>

        <!-- /左侧树形结构 -->
        <!-- ===================右侧内容窗口========================== -->
        <div class="col-sm-9 col-md-6 col-lg-8 col-xl-10" id="FloatArea">
            <!--====================================================================================================-->
            <div class="col-sm-12 form-inline border rounded" id="MineMsg" style="margin-top:3px;padding-top:3px">
                <table class="table table-striped table-hover">
                    <thead class="thead-light"><tr><th colspan="3" id ="ProjectTable">项目基本信息表</th></tr></thead><tbody>
                    <tr><td id = "t1">主程序版本</td><td id="_Version">{{.Version}}</td><td></td></tr>
                    <tr><td id = "t2">总运行时间</td><td id = "minute">{{.TotalRunTime}} min</td><td></td></tr>
                    <tr><td id = "t3">本次运行时间</td><td id = "minute2">{{.ThisRunTime}} min</td><td></td></tr>
                    <tr><td id = "t4">项目名称</td><td id="_ProjectName">{{.ProjectName}}</td><td></td></tr>
                    <tr><td id = "dog_1">计算服务监控</td><td ><div class="" id = "dog_state" ></div></td><td></td></tr>
                    <tr><td id = "t5">项目Logo</td><td id="_ProjectLogo">{{.ProjectLogo}}</td><td><img src="{{.ProjectLogo}}" alt= "Logo" style="width:40px;"></td></tr>
                    <tr><td id = "t6">项目ICON</td><td id="_ProjectIcon">{{.ProjectIcon}}</td><td><img src="{{.ProjectIcon}}" alt= "Logo" style="width:40px;"></td></tr>
                    <tr><td id = "t7">版权单位</td><td id="_Copyright">{{.Copyright}}</td><td></td></tr>
                    <tr><td id = "t8">平台路径</td><td id="_PlatPath"><a href="{{.PlatPath}}" target="_blank">{{.PlatPath}}</a></td><td></td></tr>
                    </tbody></table>
            </div>
            <div class="col-sm-12 form-inline border rounded" style="margin-top:3px;padding-top:3px">
                <table class="table table-striped table-hover">
                    <thead class="thead-light"><tr><th colspan="4" id = "r1">项目核心信息表</th></tr></thead><tbody>
                    <tr><td id = "r2">项目识别码</td><td>{{.MachineCode}}</td><td id = "r14">项目授权码</td><td id="_AuthCode">{{.AuthCode}}</td></tr>
                    <tr><td id = "r3">KPI指标计算授权数</td><td>{{.KpiAuth}}</td><td id = "r15">快速计算授权数</td><td>{{.FastAuth}}</td></tr>
                    <tr><td id = "r4">报表模板授权数</td><td>{{.ReportAuth}}</td><td></td><td></td></tr>
                    <tr><td id = "r5">分布式计算ID</td><td id="_DistributedId">{{.DistributedId}}</td><td id = "r16">备注</td><td id="_Description">{{.Description}}</td></tr>
                    <tr><th colspan="4" id = "r6">实时数据库信息:</th></tr>
                    <tr><td id = "r7">IP地址</td><td id="_RtdbServer">{{.RtdbServer}}</td><td id = "r17">端口号</td><td id="_RtdbPort">{{.RtdbPort}}</td></tr>
                    <tr><td id = "r8">数据库名</td><td id="_RtdbDbName">{{.RtdbDbName}}</td><td id = "r18">表名</td><td id="_RtdbTableName">{{.RtdbTableName}}</td></tr>
                    <tr><td id = "r9">用户名</td><td id="_RtdbUser">{{.RtdbUser}}</td><td id = "r19">密码</td><td id="_RtdbPsw">{{.RtdbPsw}}</td></tr>
                    <tr><th colspan="4" id = "r10">计算结果数据库信息:</th></tr>
                    <tr><td id = "r11">IP地址</td><td id="_ResultDbServer">{{.ResultDbServer}}</td><td id = "r20">端口号</td><td id="_ResultDbPort">{{.ResultDbPort}}</td></tr>
                    <tr><td id = "r12">数据库名</td><td id="_ResultDbName">{{.ResultDbName}}</td><td id = "r21">表名</td><td id="_ResultDbTableName">{{.ResultDbTableName}}</td></tr>
                    <tr><td id = "r13">用户名</td><td id="_ResultDbUser">{{.ResultDbUser}}</td><td id = "r22">密码</td><td id="_ResultDbPsw">{{.ResultDbPsw}}</td></tr>
                    </tbody></table>
            </div>
            <!--====================================================================================================-->
        </div>
    </div>
    <!-- ===================/主窗口========================== -->
<!-- ===================主窗口========================== -->

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
