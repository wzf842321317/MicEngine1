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
		<a class="btn btn-outline-danger"  role="button"   href="/logout/"   ><img src="../static/img/menuico/quit.svg"  style="width:25px; " id="exit"></a>
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
					<button type="button" class="btn btn-outline-primary" id="ExpandTreeNode" title=""></button>
					<button type="button" class="btn btn-outline-info" id="CollapseTreeNode" title=""></button>
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
		<!-- ===================右侧内容窗口========================== -->
		<div class="col-sm-9 col-md-6 col-lg-8 col-xl-10" id="FloatArea">
			<!--====================================================================================================-->
			<div class="container col-12 border rounded" style="padding:5px"><!--变量基本信息和选择按钮-->
				<div class="container col-12 border rounded" style="padding-top:3px"><!--变量基本信息和选择按钮-->
					<div class="col-12"><!--时间选择区-->
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
								<label class="font-weight-bold" id="Interval" for="Interval"></label>
								<select class="form-control" id="Interval" onchange="onIntervalTimeChange(this.options[this.options.selectedIndex].value)" style="width:80px">
									<option value="1">1s</option>
									<option value="5">5s</option>
									<option value="10">10s</option>
									<option value="30" selected="selected">30s</option>
									<option value="60">60s</option>
									<option value="120">120s</option>
									<option value="180">180s</option>
									<option value="240">240s</option>
								</select>
							</div>
						</div>
						<div class="row">
							<div class="col-sm-6 form-inline">
								<label class="font-weight-bold" id="ThelastTimestemp" for="ThelastTimestemp"></label>
								<label for="rd_1">
									<input class="form-control" type="radio" name="timediff" id="rd_1" onclick="onTimediffClick(60);"></label>
								<label for="rd_2">
									<input class="form-control" type="radio" name="timediff" id="rd_2" onclick="onTimediffClick(480);"></label>
								<label for="rd_3">
									<input class="form-control" type="radio" name="timediff" id="rd_3" onclick="onTimediffClick(720);"></label>
								<label for="rd_4">
									<input class="form-control" type="radio" name="timediff" id="rd_4" onclick="onTimediffClick(1440);"></label>
							</div>
							<div class="col-sm-2 form-inline btn-group"></div>
							<div class="col-sm-4 form-inline btn-group">
								<button class="btn btn-outline-primary" id="Last" onclick="onLast();" > <<</button>
								<button class="btn btn-outline-primary" id="Next" onclick="onNext();" > >></button>
							</div>
							<!--div class="col-sm-2 form-inline">
                                <label class="font-weight-bold" class="form-check-label">
                                    <input id="show_hisinterval_data" type="checkbox" class="form-check-input" value="" onclick="onShowHisIntervalData(this.id)">显示聚合数据表
                                </label>
                            </div>
                            <div class="col-sm-2 form-inline">
                                <label class="font-weight-bold" class="form-check-label">
                                    <input id="show_his_data" type="checkbox" class="form-check-input" value="" onclick="onShowHisData(this.id)">显示原始数据表
                                </label>
                            </div -->

						</div>
					</div><!--/时间选择区-->
				</div><!--/变量基本信息和选择按钮-->
			</div>
			<!--====================================================================================================-->
			<div class="container col-12 border rounded" style="margin-top:3px;padding-top:3px;padding-bottom:3px">	<!-----------数据区域------------------->
				<!-- -------------变量基本信息------------------ -->
				<div class="col-12 border rounded" style="margin-top:3px;padding-top:3px">
					<div class="row col-sm-12">
						<div class="col-sm-5 form-inline" >
							<label class="font-weight-bold" class="form-check-label" id="VariableName"></label>
							<span id="TagName"></span>
						</div>
						<div class="col-sm-2 form-inline" >
							<label class="font-weight-bold" class="form-check-label" id="LatestValue"></label>
							<span id="TagValue"></span>
						</div>
						<div class="col-sm-2 form-inline" >
							<label class="font-weight-bold" class="form-check-label" id="unit"></label>
							<span id="TagUnit"></span>
						</div>
						<div class="col-sm-3 form-inline" >
							<label class="font-weight-bold" class="form-check-label" id="UpdateTime"></label>
							<span id="TagTime"></span>
						</div>
					</div>
					<!-- -------------统计数据表------------------ -->
					<div id="HisSumData">
						<table class="table table-striped table-hover table-sm"><thead class="thead-light"><tr><th id="Min"></th><th id="Max"></th><th id="Range"></th><th id="ArithmeticMean"></th><th id="weightedMean"></th><th id="Mode"></th><th id="Median "></th><th id="Sum"></th><th id="Difference"></th><th id="PositiveDifference"></th><th id="Area "></th><th id="Points "></th><th id="SD"></th><th id="Skewness"></th><th id="Kurtosis"></th></tr></thead><tbody></tbody></table>
					</div><!--/统计数据表-->
				</div><!--/变量基本信息-->

				<!-----------所选变量列表------------------->
				<div class="col-12 border rounded" style="margin-top:3px;padding-top:3px">
					<div class="col-12 form-inline border rounded">
						<div class="col-sm-4 form-inline" id="LoadDataMsg"></div>
						<div class="col-sm-4 form-inline btn-group">
							<button class="btn btn-outline-success" id="AddTagToYTable" onclick="onAddSelectTagToYTable();" ></button>
							<button class="btn btn-outline-primary" id="AddTagToXTable" onclick="onAddSelectTagToXTable();" ></button>
						</div>
						<div class="col-sm-4 form-inline btn-group">
							<button class="btn btn-outline-danger" id="RemoveAll" onclick="onRemoveAll();" ></button>
							<button class="btn btn-outline-warning" id="Submit" onclick="onSubmit();" ></button>
						</div>
					</div>

					<div class="col-12 border rounded">
						<div class="col-12">
							<h5 id="SelectedAnalysisVariable"><small id="TagsTableSmallTitle" class="text-muted"></small></h5>
						</div>
						<div class="row">
							<div id="SelectedTagsY" class="col-sm-6"><!--已选Y变量列表-->
								<table class="table table-striped table-hover table-sm"><tr><th colspan="4" id="SelectedTags_y"></th></tr><tr><th id="Name">名称</th><th id="type">类型</th><th id="Removed">移除</th></tr></thead><tbody></tbody></table>
							</div><!--/已选Y变量列表-->
							<div id="SelectedTagsX" class="col-sm-6"><!--已选X变量列表-->
								<table class="table table-striped table-hover table-sm"><tr><th colspan="4" id="SelectedTags_x"></th></tr><tr><th id="SerialNumber">序号</th><th id="Name1">名称</th><th id="type1">类型</th><th id="Removed1">移除</th></tr></thead><tbody></tbody></table>
							</div><!--/已选X变量列表-->
						</div>
					</div>
				</div><!-----------/所选变量列表------------------->

				<div id="Echarts_His" style="height: 300px;display: ;border: 1px solid #cecece;margin-top:3px;"></div><!--原始历史趋势-->

				<div id="Echarts_HisIntervalSerial" style="height: 500px;display: ;border: 1px solid #cecece;margin-top:3px;"></div><!--等间隔历史趋势序列-->
				<div id="HisSerialRemark" class="col-12 form-inline">
					<span class="alert alert-warning"></span>
				</div><!--备注说明-->
				<!------------------回归分析数据--------------------->
				<div class="col-sm-12 form-inline border rounded" id="RegResult" style="margin-top:3px;padding-top:3px">
					<div class="col-12"><a name="ViewRegResult"></a>
						<!--------------------------------------->
						<div class="col-12" id="result">
							<h2></h2><hr/>
						</div>
						<div class="col-12" id="test">
							<h3></h3>
						</div>
						<div class="col-12">
							<table class="table table-striped table-hover table-sm"><thead><tr><th>序号</th><th>回归系数</th><th>T检验值</th><th>偏回归平方和</th><th>是否显著</th></tr></thead><tbody id="RegCoeff"></tbody></table>
						</div>
						<!--------------------------------------->
						<div class="col-12 card">
							<div class="card-header" id="regressionEquation"><h3></h3></div>
							<div class="card-body">
								<h5><span  id="RegEquation"></span></h5>
							</div>
							<div class="card-footer">
								<div class="col-12" id="InTheFormula"></div>
								<div class="col-12 form-inline" id="EquationRemark">
									<!-- div class="col-sm-4">y:变量1</div><div class="col-sm-4">x1:变量2</div><div class="col-sm-4">x2:变量3</div><div class="col-sm-4">x3:变量4</div -->
								</div>
							</div>
						</div>
						<!--------------------------------------->
						<div class="col-12" id="varianceAnalysis">
							<h3></h3>
						</div>
						<div class="col-12">
							<table class="table table-striped table-hover table-sm"><thead><tr><th>来源</th><th>平方和</th><th>自由度</th><th>方差</th><th>方差比</th><th>标准偏差</th></tr></thead><tbody  id="RegVariancef">
								<tr><td>回归</td><td>U</td><td>m</td><td>U/m</td><td rowspan="2">(U/m)/(Q/(n-m-1))</td><td rowspan="3">SD</td></tr>
								<tr><td>剩余</td><td>Q</td><td>n-m-1</td><td>Q/(n-m-1)</td></tr>
								<tr><td>总计</td><td>Tyy</td><td>n-1</td><td>---</td><td>---</td></tr>
								</tbody></table>
						</div>

						<!--------------------------------------->
						<div class="col-12" id="RegressionSignificanceTest">
							<h3></h3>
						</div>
						<div class="col-12">
							<table class="table table-striped table-hover table-sm"><thead><tr><th>检验项目</th><th>统计值</th><th>临界值</th><th>结论</th></tr></thead><tbody  id="RegSignificance">
								<tr><td>复相关系数</td><td>R</td><td>Ra</td><td><span class="badge badge-danger">不显著</span></td></tr>
								<tr><td>F值</td><td>F</td><td>Fa</td><td><span class="badge badge-success">显著</span></td></tr>
								</tbody></table>
						</div>
						<!--------------------------------------->
						<div class="col-12" id="DataGraphicAnalysis">
							<h3></h3>
						</div>
						<div class="col-12">
							<div id="Echarts_Scatter" style="height: 700px"></div>
							<hr/><!--我是分割线-->
							<div id="Echarts_Trend" style="height: 700px"></div>
						</div>
						<hr/><!--我是分割线-->
						<!--------------------------------------->
						<div class="col-12" id="DataAnalysisTable">
							<h3></h3>
						</div>
						<div class="col-12">
							<table class="table table-striped table-hover table-sm"><thead><tr><th>序号</th><th>实际值</th><th>模型值</th><th>残差</th><th>标准残差</th><th>相对偏差</th></tr></thead><tbody  id="RegDatalist"></tbody></table>
						</div>
						<!--------------------------------------->
					</div>
				</div>
				<!------------------/回归分析数据--------------------->
				<div class="alert alert-success" id="HisDataTable">
				</div>
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