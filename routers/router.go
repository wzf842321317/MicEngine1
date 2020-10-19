package routers

import (
	"github.com/astaxie/beego"
	"github.com/bkzy-wangjp/MicEngine/controllers"
)

func init() {
	beego.Router("/", &controllers.MainController{})       //主页面
	beego.Router("/index/", &controllers.MainController{}) //主页面

	beego.Router("/goldenapi/", &controllers.GoldenApiController{})                                 //WebAPI接口index面
	beego.Router("/goldenapi/search/", &controllers.GrafanaApiController{}, "*:Search")             //WebAPI测试
	beego.Router("/goldenapi/query/", &controllers.GrafanaApiController{}, "POST:GrafanaJsonQuery") //Grafana数据源
	beego.Router("/goldenapi/query/", &controllers.GoldenApiController{}, "GET:GetQuery")           //WebAPI

	beego.Router("/api/", &controllers.GoldenApiController{})                                 //WebAPI接口index面
	beego.Router("/api/search/", &controllers.GrafanaApiController{}, "*:Search")             //WebAPI测试
	beego.Router("/api/query/", &controllers.GrafanaApiController{}, "POST:GrafanaJsonQuery") //Grafana数据源
	beego.Router("/api/query/", &controllers.GoldenApiController{}, "GET:GetQuery")           //WebAPI
	beego.Router("/api/config/", &controllers.EngineController{}, "GET:GetConfig")            //读取配置信息（需要先登录）

	beego.Router("/api/script/", &controllers.EngineController{}, "*:ScriptRun")             //校验引擎脚本,直接计算脚本
	beego.Router("/api/script/compile/", &controllers.EngineController{}, "*:ScriptCompile") //校验引擎脚本
	beego.Router("/api/script/sql/", &controllers.EngineController{}, "*:ScriptMicSql")      //MicSQL脚本执行
	beego.Router("/api/test/", &controllers.EngineController{}, "*:ApiTest")                 //API检查接口

	beego.Router("/api/golden/table/insert/", &controllers.EngineController{}, "*:GoldenTableInsert")          //新建表
	beego.Router("/api/golden/table/inorup/", &controllers.EngineController{}, "*:GoldenTableInsertOrUpdate")  //新建表或者更新表
	beego.Router("/api/golden/table/delete/", &controllers.EngineController{}, "*:GoldenTableDelete")          //删除表
	beego.Router("/api/golden/table/update/", &controllers.EngineController{}, "*:GoldenTableUpdate")          //更新表
	beego.Router("/api/golden/table/select/", &controllers.EngineController{}, "*:GoldenTableSelect")          //查询表
	beego.Router("/api/golden/table/plattogolden/", &controllers.EngineController{}, "*:GoldenTablesFromPlat") //以平台数据表为基准创建庚顿表

	beego.Router("/api/golden/point/insert/", &controllers.EngineController{}, "*:GoldenPointInsert")              //新增或者更新标签点
	beego.Router("/api/golden/point/delete/", &controllers.EngineController{}, "*:GoldenPointDelete")              //删除点
	beego.Router("/api/golden/point/update/", &controllers.EngineController{}, "*:GoldenPointUpdate")              //更新点
	beego.Router("/api/golden/point/select/", &controllers.EngineController{}, "*:GoldenPointSelect")              //查询点
	beego.Router("/api/golden/point/plattogolden/", &controllers.EngineController{}, "*:GoldenPointsFromPlat")     //以平台数据Tag为基准创建庚顿标签点
	beego.Router("/api/golden/point/alarm/", &controllers.EngineController{}, "*:GoldenPointAlarm")                //查询点
	beego.Router("/api/golden/point/setsnapshot/", &controllers.EngineController{}, "POST:GoldenPointSetSnapshot") //写快照
	beego.Router("/api/golden/point/sethistory/", &controllers.EngineController{}, "POST:GoldenPointSetHistory")   //写历史
	beego.Router("/api/golden/point/all/", &controllers.EngineController{}, "*:GoldenPointsAll")                   //查看所有庚顿点的配置
	beego.Router("/api/golden/pool/", &controllers.EngineController{}, "*:GetGoldenPoolMsg")                       //获取庚顿连接池的信息
	beego.Router("/api/golden/hosttime/", &controllers.EngineController{}, "*:GetGoldenHostTime")                  //获取庚顿服务器时间

	beego.Router("/api/golden/synch/", &controllers.EngineController{}, "*:SynchGoldenPointAndPlatTaglist")                 //同步庚顿数据库的标签点与平台数据库taglist标签点
	beego.Router("/api/golden/synchgoldenfromplate/", &controllers.EngineController{}, "*:SynchGoldenPointFromPlatTaglist") //从平台数据库taglist标签点同步所有庚顿数据库的标签点

	beego.Router("/api/taglist/", &controllers.DataMonitorController{}, "*:ApiGetTaglistByLevelCode")               //获取变量列表
	beego.Router("/api/tag/", &controllers.DataMonitorController{}, "*:ApiGetTagAttribut")                          //获取Tag属性
	beego.Router("/api/hisinterval/", &controllers.DataMonitorController{}, "*:ApiGetHisInterData")                 //获取等间隔历史数据
	beego.Router("/api/history/", &controllers.DataMonitorController{}, "*:ApiGetHistoryData")                      //获取历史数据
	beego.Router("/api/historysummary/", &controllers.DataMonitorController{}, "*:ApiGetHistorySummaryData")        //获取历史统计数据
	beego.Router("/api/snapshot/", &controllers.DataMonitorController{}, "*:ApiGetSnapshotData")                    //获取快照数据
	beego.Router("/api/writesnapshot/", &controllers.DataMonitorController{}, "*:ApiWriteSnapshot")                 //写快照数据
	beego.Router("/api/regression/", &controllers.DataMonitorController{}, "*:ApiRegression")                       //最小二乘回归分析
	beego.Router("/api/appendkpi/", &controllers.DataMonitorController{}, "*:ApiAppendTagKpi2CfgList")              //添加KPI
	beego.Router("/api/samplelabtag/", &controllers.DataMonitorController{}, "*:ApiSampleLabTag")                   //获取样本的化验变量标签信息
	beego.Router("/api/samplelabresult/", &controllers.DataMonitorController{}, "*:ApiSampleLabResult")             //获取样本的化验结果信息
	beego.Router("/api/goodscfg/", &controllers.DataMonitorController{}, "*:ApiGoodsCfg")                           //获取物耗标签的基本信息
	beego.Router("/api/goodsdatas/", &controllers.DataMonitorController{}, "*:ApiGoodsDatas")                       //获取物耗历史数据
	beego.Router("/api/patrollist/", &controllers.DataMonitorController{}, "*:ApiPatrolTagList")                    //获取巡检标签的基本信息
	beego.Router("/api/patrolresult/", &controllers.DataMonitorController{}, "*:ApiPatrolResult")                   //获取巡检历史数据
	beego.Router("/api/patrolpicurl/", &controllers.DataMonitorController{}, "*:ApiGetPatrolPicUrl")                //获取巡检照片URL
	beego.Router("/api/patrolsitelist/", &controllers.DataMonitorController{}, "*:ApiGetPatrolSiteListByLevelCode") //通过层级码获取该层级下的巡检点
	beego.Router("/api/userlog/", &controllers.DataMonitorController{}, "*:ApiGetUserLogs")                         //获取用户日志信息
	beego.Router("/api/userloganalys/", &controllers.DataMonitorController{}, "*:ApiGetUserLogsAnalys")             //获取用户日志统计数据
	beego.Router("/api/listen/", &controllers.DataMonitorController{}, "POST:ApiPortListen")                        //端口监听
	beego.Router("/api/updatetagnodetree/", &controllers.DataMonitorController{}, "*:ApiUpdateNodeTree")            //更新用户的标签节点树
	beego.Router("/api/kpi/getconfig", &controllers.DataMonitorController{}, "*:ApiGetKpiConfig")                   //获取KPI配置信息
	beego.Router("/api/kpi/getresult", &controllers.DataMonitorController{}, "*:ApiGetKpiResult")                   //获取KPI计算结果

	beego.Router("/api/srtd/statistic", &controllers.DataMonitorController{}, "*:ApiSrtdStatistic")               //系统时序数据统计
	beego.Router("/api/srtd/statisticauto", &controllers.DataMonitorController{}, "*:ApiSrtdStatisticAuto")       //系统时序数据统计(如果有异常数据,自动统计排除异常后的结果)
	beego.Router("/api/srtd/increment/statistic", &controllers.DataMonitorController{}, "*:ApiSrtdIncrStatistic") //系统时序数据增量统计

	beego.Router("/api/getreportlistsbyuserid/", &controllers.ReportController{}, "*:ApiGetReportListsByUserId") //获取用户被授权的报表节点
	beego.Router("/api/getreportchildnodes/", &controllers.ReportController{}, "*:ApiGetReportChildNodes")       //获取所选报表层级下的所有子节点
	beego.Router("/api/getreporttpllist/", &controllers.ReportController{}, "*:ApiGetReportTplList")             //获取所选报表的模板列表
	beego.Router("/api/getreportreultlist/", &controllers.ReportController{}, "*:ApiGetReportResultList")        //获取所选报表的结果列表
	beego.Router("/api/addreportlevel/", &controllers.ReportController{}, "POST:ApiAddReportLevel")              //添加报表层级
	beego.Router("/api/editreportlevel/", &controllers.ReportController{}, "POST:ApiEditReportLevel")            //编辑报表层级
	beego.Router("/api/deletereportlevel/", &controllers.ReportController{}, "*:ApiDeleteReportLevel")           //删除报表层级
	beego.Router("/api/setfileasreporttpl/", &controllers.ReportController{}, "*:ApiSetFileAsReportTpl")         //设置文件作为报表模板
	beego.Router("/api/uploadrpttpl/", &controllers.ReportController{}, "*:ApiUpLoadFile")                       //上传报表模板文件
	beego.Router("/api/download/", &controllers.ReportController{}, "*:ApiDownLoadFile")                         //下载文件
	beego.Router("/api/viewexcel/", &controllers.ReportController{}, "*:ApiViewExcel")                           //在线预览文件
	beego.Router("/api/setcellvalue/", &controllers.ReportController{}, "POST:ApiSetExcelCellValue")             //设置Excel单元格的值
	beego.Router("/api/setcellformula/", &controllers.ReportController{}, "POST:ApiSetExcelCellFormula")         //设置Excel单元格的公式

	beego.Router("/api/getworkshops/", &controllers.ReportController{}, "*:ApiGetWorkShops") //获取车间列表

	beego.Router("/api/updatepswd/", &controllers.LoginController{}, "*:ApiUpdatePswd")          //更新密码
	beego.Router("/api/resetpswd/", &controllers.LoginController{}, "*:ApiResetPswd")            //重置密码
	beego.Router("/api/editprojectmsg/", &controllers.LoginController{}, "*:ApiUpdateEngineCfg") //编辑项目信息

	beego.Router("/login/", &controllers.LoginController{})                      //用户登录
	beego.Router("/logout/", &controllers.LoginController{}, "*:LogOut")         //用户退出
	beego.Router("/usercenter/", &controllers.LoginController{}, "*:UserCenter") //用户中心

	beego.Router("/snapshot/", &controllers.DataMonitorController{}, "*:PageSnapShot")     //数据快照
	beego.Router("/history/", &controllers.DataMonitorController{}, "*:PageHistory")       //数据历史
	beego.Router("/regression/", &controllers.DataMonitorController{}, "*:PageRegression") //回归分析
	beego.Router("/compare/", &controllers.DataMonitorController{}, "*:PageCompare")       //数据比较
	beego.Router("/monitor/", &controllers.DataMonitorController{}, "*:PageMonitor")       //数据监控
	beego.Router("/samplelab/", &controllers.DataMonitorController{}, "*:PageSampleLab")   //质检化验
	beego.Router("/goods/", &controllers.DataMonitorController{}, "*:PageGoods")           //物耗
	beego.Router("/patrol/", &controllers.DataMonitorController{}, "*:PagePatrol")         //巡检
	beego.Router("/kpi/", &controllers.DataMonitorController{}, "*:PageKpi")               //kpi
	beego.Router("/report/", &controllers.DataMonitorController{}, "*:PageReport")         //报表
	beego.Router("/reportedit/", &controllers.DataMonitorController{}, "*:PageReportEdit") //报表管理

	beego.Router("/managerproject/", &controllers.DataMonitorController{}, "*:PageManagerProject")       //项目管理
	beego.Router("/managerlog/", &controllers.DataMonitorController{}, "*:PageManagerLog")               //日志管理
	beego.Router("/managerusers/", &controllers.DataMonitorController{}, "*:PageManagerUers")            //用户管理
	beego.Router("/managerpermission/", &controllers.DataMonitorController{}, "*:PageManagerPermission") //权限管理

	beego.Router("/calculate/", &controllers.DataMonitorController{}, "*:PageCalculate")             //计算
	beego.Router("/api/calculate/", &controllers.ReportController{}, "*:SaveCalculate")              //计算
	beego.Router("/api/calculate/data", &controllers.ReportController{}, "*:ApiGetCalculate")        //计算
	beego.Router("/api/calculate/json", &controllers.DataMonitorController{}, "*:ApiOreProcessJson") //计算
}
