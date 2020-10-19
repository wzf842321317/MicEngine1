package models

const (
	_version               = "V1.1.2038"
	_defaultDBServer       = "rm-2zehm037lnj7rl64wo.mysql.rds.aliyuncs.com" //默认平台数据库地址
	_defaultDBPort         = "3306"                                         //默认平台数据库端口
	_cfgPath               = "conf/config.ini"                              //配置文件路径
	_machineCodePath       = "MachineCode.txt"                              //机器码保存路径
	_cfgPlatSection        = "plat"                                         //配置文件内平台数据库参数的section名称
	_cfgLogSection         = "log"
	_serialIndicatorAuth   = 50                  //默认串行计算指标授权数量
	_parallelIndicatorAuth = 50                  //默认并行计算指标授权数量
	_reportAuth            = 5                   //默认报表授权数量
	_PlatPicPath           = "common/picShow/"   //平台图片路径
	_ReportSheet           = "Sheet1"            //报表默认sheet名
	_ReportPreFix          = "MicForm"           //报表前缀
	_SqlitePath            = "data/micengine.db" //SQLite3数据库路径
	_SqliteAlias           = "sqlite"            //SQLite3数据库的别名

	_TIMEFOMAT               = "2006-01-02 15:04:05"
	_ResultDBAlias           = "result" //KPI结果数据库别名
	_variableTypeInSysDicCat = 7        //在sys_dictionary_catalog表中定义的【变量类型】的ID
	_PermissionTypeOfReport  = 7        //在sys_permission表中permission_type定义的【报表类型】的权限
)

var EngineCfgMsg EngineConfig //授权信息

type EngineConfig struct { //引擎授权信息
	IsAuth                bool   //是否授权
	UserName              string //用户名
	Version               string //版本号
	StartTime             string //启动时间
	DogStatus             bool   //看门狗状态,true:stop,false:runing
	ResultDBAlias         string //结果数据库别名
	CfgMsg                CalcKpiEngineConfig
	Log                   LogCfg
	Sys                   SysCfg
	Err                   bool  //配置信息存在错误
	totalRunTimeWhenStart int64 //启动时记录的运行总时间
}
type SysCfg struct {
	Debug                bool   //调试模式,会输出数据库的查询语句
	Createtable          bool   //启动时自动创建表
	PeriodKpiEnable      bool   //周期KPI计算使能
	FastKpiEnable        bool   //快速KPI计算使能
	ReportEnable         bool   //报表计算使能
	AlarmEnable          bool   //报警计算使能
	WebEnable            bool   //Web功能使能
	Lang                 string //语言配置
	FirstPage            string //登录之后的首页面
	TimeFormat           string //系统中的时间格式
	SleepTimeOnStart     int64  //启动时停顿的时间(秒),以方便观察启动信息
	SaveTimeInComTable   int64  //计算结果在通用表(非月度表)中的保存时间(天),0为永久保存,默认为90天
	MaxHisSliceTime      int64  //读取历史数据时的分段时间长度(小时)
	GoldenCennectPool    int    //庚顿数据库连接池大小,最小1,最大50
	ExcelFormulaCalcDeep int    //报表在线预览时对Excel公式的计算深度(0时不限制,直到全部计算完毕,但对于嵌套单元格很多且很大的报表，加载会很慢)
	ResetDogFilePath     string //用于重置看门狗的bat文件路径
}
type LogCfg struct {
	Path         string //初始日志文件名及路径
	Maxlines     int    //日志文件最大行数，当append=true时有效
	Maxdays      int    //日志文件最大大小，当append=true时有效
	Level        int    //日志文件日志输出等级,0~7
	Consolelevel int    //控制台日志信息输出等级,0~7
	Daily        bool   //跨天后是否创建新日志文件，当append=true时有效
	Debug        bool   //日志的调试模式,会显示日志发出文件名和行号
}

type SrtData struct { //实时数据结构
	Time  string
	Value float64
}

//质检化验样本树结构
type SampleLabTree struct {
	Id           int64
	Pid          int64
	Name         string
	NodeType     int64
	ItemId       int64
	Seq          int64
	TreeLevel    string
	IsLeaf       bool
	IsRegular    int
	Func         int64
	FuncName     string
	SamplingSite string
	BaseTime     string
	ShiftHour    int64
}

//物耗树结构
type GoodsTree struct {
	Id        int64
	Pid       int64
	Name      string
	NodeType  int64
	TreeLevel string
	IsLeaf    bool
	BaseTime  string
	ShiftHour int
}

//巡检设备点树结构
type CheckTree struct {
	Id        int64
	Pid       int64
	SiteId    int64
	Name      string
	NodeType  int64
	TreeLevel string
	IsLeaf    bool
	BaseTime  string
	ShiftHour int
	LineId    int64
	LineName  string
	DeptId    int64
	DeptName  string
}

type Arithmetic struct {
	Id                      int64
	ArithmeticName          string //算法工具包名称
	ArithmeticName2         string //算法工具包名称
	ArithmeticResultType    string //计算结果数据类型
	ArithmeticResultEcharts string //计算结果展示类型
	ArithmeticUrl           string //所需页面编号
	ArithmeticType          string //计算显示名称（例如：求和）
}

type ArithmeticResult struct {
	Id                         int64
	ResultName                 string
	ArithmeticObjectId         string
	ArithmeticObjectType       string
	ArithmeticObjectRemark     string
	ArithmeticAuxiliaryId1     string
	ArithmeticAuxiliaryRemark1 string
	ArithmeticAuxiliaryId2     string
	ArithmeticAuxiliaryRemark2 string
	BeginTime                  string
	EndTime                    string
	FinalResult                string
	FinalResultType            string
	ArithmeticId               string
	FinalResultEacharts        string
	CreTime                    string
	CrePerson                  string
}
