package models

import (
	"database/sql"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/astaxie/beego/orm"
	auth "github.com/bkzy-wangjp/Author/EngineAuth"

	"github.com/astaxie/beego/logs"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
	"gopkg.in/ini.v1"
)

func init() {
	//读取配置信息
	logs.Info("读取配置文件……")
	time.Sleep(1 * time.Second)
	lgcfg, dbstr, syscfg := readCfg()

	// 参数1   driverName
	// 参数2   数据库类型
	// 这个用来设置 driverName 对应的数据库类型
	// mysql / sqlite3 / postgres 这三种是默认已经注册过的，所以可以无需设置
	//orm.RegisterDriver("mysql", orm.DRMySQL)

	//注册数据库，并设定为默认连接数据库
	// 参数1        数据库的别名，用来在 ORM 中切换数据库使用
	// 参数2        driverName
	// 参数3        对应的链接字符串
	// 参数4(可选)  设置最大空闲连接
	// 参数5(可选)  设置最大数据库连接 (go >= 1.2)
	logs.Info("软件名称:MicEngine")
	logs.Info("软件版本:%s", _version)
	logs.Info("测试打开平台数据库和检查授权……")
	time.Sleep(1 * time.Second)
	srcdb, err := sql.Open("mysql", dbstr)
	if err != nil {
		logs.Emergency("配置文件中数据库信息填写错误，无法打开数据库，请检查:%s", dbstr)
	}
	err = srcdb.Ping()
	if err == nil {
		if orm.RegisterDataBase("default", "mysql", dbstr, 30) != nil {
			logs.Emergency("配置文件中数据库信息填写错误，无法打开数据库，请检查:%s", dbstr)
			EngineCfgMsg.Err = true
		} else {
			logs.Info("平台数据库打开成功")
		}
	} else {
		logs.Emergency("配置文件中数据库信息填写错误，无法打开数据库，请检查:%s", dbstr)
		EngineCfgMsg.Err = true
	}

	if orm.RegisterDataBase(_SqliteAlias, "sqlite3", _SqlitePath) != nil {
		logs.Warn("打开本地SQLite数据库失败！")
	}

	//注册模型
	orm.RegisterModel(new(SysMenuPermission), new(SysMiddlePermission), new(SysRoleUser), new(SysRolePermission))
	orm.RegisterModel(new(CalcFunctionDic), new(CalcKpiEngineConfig), new(CalcKpiIndexDic), new(CalcKpiConfigList))
	orm.RegisterModel(new(MineBasicInfo), new(MineTableList), new(MsMiddleCorrelation), new(MineDDictionaryOfVaribleset))
	orm.RegisterModel(new(OreProcessDWorkshop), new(OreProcessDWorkstage), new(OreProcessConcentration))
	orm.RegisterModel(new(OreProcessDTaglist), new(MsMiddleCorrelationString))
	orm.RegisterModel(new(GoodsConsumePool), new(GoodsConfigInfo), new(GoodsConsumeInfo))
	orm.RegisterModel(new(SysDictionaryCatalog), new(SysDictionary), new(SysRealData), new(SysUnit))
	orm.RegisterModel(new(SysUser), new(SysRole), new(SysMenu), new(SysPermission), new(SysLog))
	orm.RegisterModel(new(MineDcs), new(DatatableInfo), new(DataAcqStation), new(RelevanceDcsToDbtable))
	orm.RegisterModel(new(RealTimeMonitor), new(RealTimeMonitorPermission), new(SysReportPermission))
	orm.RegisterModel(new(CalcKpiReportList), new(MineDeptDic), new(MineDeptInfo), new(MineConstructionList))
	orm.RegisterModel(new(OreProcessEquipmentDic), new(LabAnaIndexUnitDic), new(SamplingManage), new(LabSampletypePool))
	orm.RegisterModel(new(SamplingManageSub), new(SamplelistToLab), new(LabAnaResultTsd))
	orm.RegisterModel(new(CheckVariableSet), new(CheckItemExe), new(CheckLine), new(CheckPlan))
	orm.RegisterModel(new(CheckPlanExe), new(CheckSite), new(CheckSiteEquipRel), new(CheckSiteExe))
	orm.RegisterModel(new(CheckTagList), new(Period), new(PeriodItem), new(MineCheckImg))
	orm.RegisterModel(new(CalcTagAlarm))
	//SQLite3数据库
	orm.RegisterModel(new(HeartBeat))
	orm.RegisterModel(new(KpiArithmetic), new(KpiArithmeticResult))

	//自动创建表
	if syscfg.Createtable == true {
		orm.RunSyncdb("default", false, true)
	}
	//SQLite3数据库始终自动创建
	orm.RunSyncdb(_SqliteAlias, false, true)

	//开启调试模式?
	if syscfg.Debug == true {
		orm.Debug = true

	} else {
		orm.Debug = false
	}
	EngineCfgMsg.getEngineConfig() //读取计算引擎配置信息
	EngineCfgMsg.Log = lgcfg
	EngineCfgMsg.Sys = syscfg
	EngineCfgMsg.Version = _version

	//计算结果存储数据库配置
	resdbstr := EngineCfgMsg.CfgMsg.getResultDBString()
	if len(resdbstr) < 33 { //结果数据库没有配置地址
		resdbstr = dbstr //使用系统数据库
	}
	if strings.ToLower(dbstr) == strings.ToLower(resdbstr) { //计算结果数据库与系统数据库相同
		EngineCfgMsg.ResultDBAlias = "default"
	} else {
		logs.Info("测试打开计算结果数据库……")
		time.Sleep(1 * time.Second)
		EngineCfgMsg.ResultDBAlias = _ResultDBAlias
		resdb, err := sql.Open("mysql", resdbstr)
		err = resdb.Ping()
		if err == nil {
			if orm.RegisterDataBase(EngineCfgMsg.ResultDBAlias, "mysql", resdbstr, 30) != nil {
				logs.Emergency("无法打开计算结果数据库，请检查计算服务引擎中的相关配置:%s", resdbstr)
				EngineCfgMsg.Err = true
			} else {
				logs.Info("计算结果数据库打开成功")
			}
		} else {
			logs.Emergency("无法打开计算结果数据库，请检查计算服务引擎中的相关配置:%s", resdbstr)
			EngineCfgMsg.Err = true
		}
	}
	if EngineCfgMsg.Err == false {
		setMysqlMaxPreparedStmtCount(EngineCfgMsg.ResultDBAlias, 500000)
		res := new(CalcKpiResult)
		_, err := res.createResultTable(EngineCfgMsg.ResultDBAlias, EngineCfgMsg.CfgMsg.ResultdbTbname) //检查结果表是否存在，不存在就创建
		if err != nil {
			logs.Info("在计算结果数据库中创建计算结果数据表失败")
		}
	}

	logs.Info("测试打开实时数据库……")
	time.Sleep(1 * time.Second)
	micgd := new(MicGolden)
	GDPOOL, err = micgd.NewGoldenPool()
	GDDB = new(GoldenDataAlarm)
	GDDB.GoldenDB = micgd
	if err != nil {
		logs.Emergency("建立庚顿连接池失败:[%s]", err.Error())
		EngineCfgMsg.Err = true
		time.Sleep(10 * time.Second)
	} else {
		if vers, err := micgd.GoldenGetApiVersion(); err != nil { //检测庚顿数据库是否设置成功
			logs.Emergency("无法打开实时数据库，请检查配置:%s", err.Error())
			EngineCfgMsg.Err = true
		} else {
			logs.Info("实时数据库打开成功,数据库API接口版本:", vers)
		}
	}
	if EngineCfgMsg.Err == true { //存在系统错误的时候,禁止开启计算
		EngineCfgMsg.Sys.FastKpiEnable = false
		EngineCfgMsg.Sys.PeriodKpiEnable = false
		EngineCfgMsg.Sys.AlarmEnable = false
	}
	lst := new(CalcKpiConfigList)
	lst.InsertTagKpi2CfgList() //检查是否有KPI指标，没有则自动创建（前提是要有kpi_dic到variableset的多对多关系表）
	logs.Info("检查菜单数据……")
	menu := new(SysMenu)
	menu.CheckWebMenu() //检查Web菜单是否存在，不存在就创建
}

//读取配置文件
func readCfg() (LogCfg, string, SysCfg) {
	cfg, err := ini.Load(_cfgPath)
	if err != nil {
		logs.Emergency("Fail to read Config file[读取配置文件失败]: %v", err)
		fmt.Printf("The program will closed after 20 second[程序将在20秒后关闭]:")
		for i := 20; i > 0; i-- {
			time.Sleep(time.Second)
			fmt.Printf(".")
		}
		os.Exit(1)
	}
	lg := new(LogCfg)
	sys := new(SysCfg)
	if cfg.Section(_cfgPlatSection).Key("server").String() == "" { //默认平台数据库地址
		cfg.Section(_cfgPlatSection).Key("server").SetValue(_defaultDBServer)
	}
	if cfg.Section(_cfgPlatSection).Key("port").String() == "" { //默认平台数据库端口
		cfg.Section(_cfgPlatSection).Key("port").SetValue(_defaultDBPort)
	}
	dbstr := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8", //mysql地址
		cfg.Section(_cfgPlatSection).Key("user_name").String(),
		cfg.Section(_cfgPlatSection).Key("password").String(),
		cfg.Section(_cfgPlatSection).Key("server").String(),
		cfg.Section(_cfgPlatSection).Key("port").String(),
		cfg.Section(_cfgPlatSection).Key("db_name").String())

	lg.Daily = true
	lg.Level, err = cfg.Section(_cfgLogSection).Key("level").Int()
	if err != nil {
		logs.Alert("Config Error,the [%s].[level] must be int,the config value is %s", _cfgLogSection, cfg.Section(_cfgLogSection).Key("level"))
		lg.Level = 7
	}
	lg.Maxdays, err = cfg.Section(_cfgLogSection).Key("maxdays").Int()
	if err != nil {
		logs.Alert("Config Error,the [%s].[maxdays] must be int,the config value is %s", _cfgLogSection, cfg.Section(_cfgLogSection).Key("maxdays"))
		lg.Maxdays = 30
	}
	lg.Maxlines, err = cfg.Section(_cfgLogSection).Key("maxlines").Int()
	if err != nil {
		logs.Alert("Config Error,the [%s].[maxlines] must be int,the config value is %s", _cfgLogSection, cfg.Section(_cfgLogSection).Key("maxlines"))
		lg.Maxlines = 100000
	}
	sys.TimeFormat = cfg.Section(_cfgPlatSection).Key("timeFormat").String()
	lg.Path = cfg.Section(_cfgLogSection).Key("path").String()
	lg.Consolelevel, err = cfg.Section(_cfgLogSection).Key("consolelevel").Int()
	if err != nil {
		logs.Alert("Config Error,the [%s].[consolelevel] must be int,the config value is %s", _cfgLogSection, cfg.Section(_cfgLogSection).Key("consolelevel"))
		lg.Consolelevel = 7
	}
	lg.Debug, err = cfg.Section(_cfgLogSection).Key("debug").Bool()
	if err != nil {
		logs.Alert("Config Error,the [%s].[debug] must be bool,the config value is %s", _cfgLogSection, cfg.Section(_cfgLogSection).Key("debug"))
		lg.Debug = true
	}
	sys.Debug, err = cfg.Section(_cfgPlatSection).Key("debug").Bool()
	if err != nil {
		logs.Alert("Config Error,the [%s].[debug] must be bool,the config value is %s", _cfgPlatSection, cfg.Section(_cfgLogSection).Key("debug"))
		sys.Debug = false
	}
	sys.Createtable, err = cfg.Section(_cfgPlatSection).Key("createtable").Bool()
	if err != nil {
		logs.Alert("Config Error,the [%s].[createtable] must be bool,the config value is %s", _cfgPlatSection, cfg.Section(_cfgPlatSection).Key("createtable"))
		sys.Createtable = false
	}
	sys.PeriodKpiEnable, err = cfg.Section(_cfgPlatSection).Key("PeriodKpiEnable").Bool()
	if err != nil {
		logs.Alert("Config Error,the [%s].[PeriodKpiEnable] must be bool,the config value is %s", _cfgPlatSection, cfg.Section(_cfgPlatSection).Key("PeriodKpiEnable"))
		sys.PeriodKpiEnable = true
	}
	sys.FastKpiEnable, err = cfg.Section(_cfgPlatSection).Key("FastKpiEnable").Bool()
	if err != nil {
		logs.Alert("Config Error,the [%s].[FastKpiEnable] must be bool,the config value is %s", _cfgPlatSection, cfg.Section(_cfgPlatSection).Key("FastKpiEnable"))
		sys.FastKpiEnable = true
	}
	sys.ReportEnable, err = cfg.Section(_cfgPlatSection).Key("ReportEnable").Bool()
	if err != nil {
		logs.Alert("Config Error,the [%s].[ReportEnable] must be bool,the config value is %s", _cfgPlatSection, cfg.Section(_cfgPlatSection).Key("ReportEnable"))
		sys.ReportEnable = true
	}
	sys.AlarmEnable, err = cfg.Section(_cfgPlatSection).Key("AlarmEnable").Bool()
	if err != nil {
		logs.Alert("Config Error,the [%s].[AlarmEnable] must be bool,the config value is %s", _cfgPlatSection, cfg.Section(_cfgPlatSection).Key("AlarmEnable"))
		sys.AlarmEnable = true
	}
	sys.WebEnable, err = cfg.Section(_cfgPlatSection).Key("WebEnable").Bool()
	if err != nil {
		logs.Alert("Config Error,the [%s].[WebEnable] must be bool,the config value is %s", _cfgPlatSection, cfg.Section(_cfgPlatSection).Key("WebEnable"))
		sys.WebEnable = true
	}
	sys.SleepTimeOnStart, err = cfg.Section(_cfgPlatSection).Key("SleepTimeOnStart").Int64()
	if err != nil {
		logs.Alert("Config Error,the [%s].[SleepTimeOnStart] must be int,the config value is %s", _cfgPlatSection, cfg.Section(_cfgPlatSection).Key("SleepTimeOnStart"))
		sys.SleepTimeOnStart = 10
	}
	sys.SaveTimeInComTable, err = cfg.Section(_cfgPlatSection).Key("SaveTimeInComTable").Int64()
	if err != nil {
		logs.Alert("Config Error,the [%s].[SaveTimeInComTable] must be int,the config value is %s", _cfgPlatSection, cfg.Section(_cfgPlatSection).Key("SaveTimeInComTable"))
		sys.SaveTimeInComTable = 90
	}
	sys.MaxHisSliceTime, err = cfg.Section(_cfgPlatSection).Key("MaxHisSliceTime").Int64()
	if err != nil {
		logs.Alert("Config Error,the [%s].[MaxHisSliceTime] must be int,the config value is %s", _cfgPlatSection, cfg.Section(_cfgPlatSection).Key("MaxHisSliceTime"))
		sys.MaxHisSliceTime = 8
	}
	sys.GoldenCennectPool, err = cfg.Section(_cfgPlatSection).Key("GoldenCennectPool").Int()
	if err != nil {
		logs.Alert("Config Error,the [%s].[GoldenCennectPool] must be int,the config value is %s", _cfgPlatSection, cfg.Section(_cfgPlatSection).Key("GoldenCennectPool"))
		sys.GoldenCennectPool = 50
	}
	//ExcelFormulaCalcDeep
	sys.ExcelFormulaCalcDeep, err = cfg.Section(_cfgPlatSection).Key("ExcelFormulaCalcDeep").Int()
	if err != nil {
		logs.Alert("Config Error,the [%s].[ExcelFormulaCalcDeep] must be int,the config value is %s", _cfgPlatSection, cfg.Section(_cfgPlatSection).Key("ExcelFormulaCalcDeep"))
		sys.ExcelFormulaCalcDeep = 5
	}
	sys.FirstPage = cfg.Section(_cfgPlatSection).Key("FirstPage").String()
	sys.ResetDogFilePath = cfg.Section(_cfgPlatSection).Key("ResetWatchdogFilePath").String()
	return *lg, dbstr, *sys
}

//获取计算引擎配置信息
func (cfg *EngineConfig) getEngineConfig() {
	o := orm.NewOrm()                                         //新建orm对象
	o.Using("default")                                        //选定数据库
	cfg.CfgMsg.MachineCode = auth.MachineCodeEncrypt()        //获取机器码
	cfg.CfgMsg.ParallelIndicatorAuth = _parallelIndicatorAuth //默认授权数量
	cfg.CfgMsg.SerialIndicatorAuth = _serialIndicatorAuth     //默认授权数量
	cfg.CfgMsg.ReportAuth = _reportAuth                       //默认授权数量
	cfg.StartTime = time.Now().Format(_TIMEFOMAT)             //启动时间
	cfg.CfgMsg.CreateTime = cfg.StartTime                     //创建时间
	cfg.CfgMsg.UpdateTime = cfg.StartTime                     //更新时间
	cfg.CfgMsg.Status = 1
	if created, _, err := o.ReadOrCreate(&cfg.CfgMsg, "MachineCode"); err == nil { //读取机器码对应的配置信息,如果没有就创建一个
		if created {
			mc, _ := ini.Load([]byte(fmt.Sprintf("MachineCode = %s", cfg.CfgMsg.MachineCode)), []byte(fmt.Sprintf("Time = %s", time.Now().Format("2006-01-02 15:04:05"))))
			mc.SaveTo(_machineCodePath)
			cfg.CfgMsg.CreateTime = cfg.StartTime //创建时间
			o.Update(&cfg.CfgMsg, "CreateTime", "Status")
			infostr := fmt.Sprintf("系统初次运行，需要进行授权，请提供本机机器码给供应商\n")
			logs.Info(infostr)
			infostr = fmt.Sprintf("本机机器码为:%s，已保存在程序安装目录下的[%s]文件中。\n", cfg.CfgMsg.MachineCode, _machineCodePath)
			logs.Info(infostr)
		} else {
			var pcnt, scnt, rcnt int
			pcnt, scnt, rcnt, cfg.UserName, cfg.IsAuth = auth.AuthorizationCheck(cfg.CfgMsg.AuthCode)
			if cfg.IsAuth == false { //尚未被授权
				mc, _ := ini.Load([]byte(fmt.Sprintf("MachineCode = %s", cfg.CfgMsg.MachineCode)), []byte(fmt.Sprintf("Time = %s", time.Now().Format("2006-01-02 15:04:05"))))
				mc.SaveTo(_machineCodePath)
				infostr := fmt.Sprintf("本系统尚未被授权，允许临时使用的授权数量为:并行授权[%d]个,串行授权[%d]个,报表授权[%d]个\n", cfg.CfgMsg.ParallelIndicatorAuth, cfg.CfgMsg.SerialIndicatorAuth, cfg.CfgMsg.ReportAuth)
				logs.Info(infostr)
				infostr = fmt.Sprintf("本机机器码为:%s，已保存在程序安装目录下的[%s]文件中,请提供给供应商以获取授权。\n", cfg.CfgMsg.MachineCode, _machineCodePath)
				logs.Info(infostr)
			} else { //已经授权
				cfg.CfgMsg.ParallelIndicatorAuth = int64(pcnt) //授权数量
				cfg.CfgMsg.SerialIndicatorAuth = int64(scnt)   //授权数量
				cfg.CfgMsg.ReportAuth = int64(rcnt)            //授权数量
				infostr := fmt.Sprintf("本系统已获得授权,授权给%s,并行计算授权[%d]个,串行计算授权[%d]个,报表授权[%d]", cfg.UserName, cfg.CfgMsg.ParallelIndicatorAuth, cfg.CfgMsg.SerialIndicatorAuth, cfg.CfgMsg.ReportAuth)
				logs.Info(infostr)
			}
		}
		cfg.CfgMsg.Version = _version
		o.Update(&cfg.CfgMsg, "ParallelIndicatorAuth", "SerialIndicatorAuth", "ReportAuth", "Version") //保存授权数量
		if cfg.CfgMsg.SerialCalcDelaySec == 0 || cfg.CfgMsg.SerialIndicNumPerThread == 0 || cfg.CfgMsg.ParallelIndicReloadInterval == 0 || cfg.CfgMsg.ReportCalcInterval == 0 {
			if cfg.CfgMsg.SerialCalcDelaySec == 0 {
				cfg.CfgMsg.SerialCalcDelaySec = 120
			}
			if cfg.CfgMsg.SerialIndicNumPerThread == 0 {
				cfg.CfgMsg.SerialIndicNumPerThread = 500
			}
			if cfg.CfgMsg.ParallelIndicReloadInterval == 0 {
				cfg.CfgMsg.ParallelIndicReloadInterval = 600
			}
			if cfg.CfgMsg.ReportCalcInterval == 0 {
				cfg.CfgMsg.ReportCalcInterval = 300
			}
			o.Update(&cfg.CfgMsg, "SerialCalcDelaySec", "SerialIndicNumPerThread", "ParallelIndicReloadInterval", "ReportCalcInterval")
		}
	}
	cfg.totalRunTimeWhenStart = cfg.CfgMsg.TotalMinutes
}

//合成结果数据库的字符串
func (cfg *CalcKpiEngineConfig) getResultDBString() string {
	dbstr := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8", //mysql地址
		cfg.ResultdbUser,
		cfg.ResultdbPsw,
		cfg.ResultdbServer,
		cfg.ResultdbPort,
		cfg.ResultdbDbname)
	return dbstr
}
