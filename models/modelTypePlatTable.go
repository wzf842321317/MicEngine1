package models

//计算服务引擎配置表(calc_kpi_engine_config)
type CalcKpiEngineConfig struct {
	Id                          int64  `orm:"auto"` //自动id
	MachineCode                 string //机器码
	AuthCode                    string //授权码
	DistributedId               int64  `orm:"default(0)"`  //分布式计算ID，只计算calc_indicator_config_list表中与该id相匹配的指标
	SerialIndicatorAuth         int64  `orm:"default(50)"` //串行指标授权数量
	ParallelIndicatorAuth       int64  `orm:"default(50)"` //并行计算指标授权数量
	ReportAuth                  int64  `orm:"default(5)"`  //报表授权数量
	RtdbClass                   string //实时数据库类型，可选：golden,influxdb,mysql,mssql,sqlite,odbc
	RtdbServer                  string //实时数据库主机地址
	RtdbPort                    int    `orm:"default(6327)"` //实时数据库端口号
	RtdbUser                    string //实时数据库用户名
	RtdbPsw                     string //实时数据库密码
	RtdbDbname                  string //实时数据库的库名称
	RtdbTbname                  string //实时数据库的表名称
	ResultdbClass               string //计算结果数据库类型,可选：golden,influxdb,mysql,mssql,sqlite,odbc
	ResultdbServer              string //计算结果数据库主机地址
	ResultdbPort                int    `orm:"default(3306)"` //计算结果数据库端口号
	ResultdbUser                string //计算结果数据库用户名
	ResultdbPsw                 string //计算结果数据库密码
	ResultdbDbname              string //计算结果数据库的库名称
	ResultdbTbname              string //计算结果数据库的表名称
	ParallelIndicReloadInterval int64  `orm:"default(600)"` //并行计算指标配置信息重新加载时间间隔,单位:秒
	SerialCalcInterval          int64  `orm:"default(300)"` //串行计算指标扫描时间间隔,单位:秒
	SerialCalcDelaySec          int64  `orm:"default(120)"` //串行计算延迟时间（延迟计算，以便给数据同步留出足够时间），单位:秒
	SerialIndicNumPerThread     int64  `orm:"default(500)"` //每个线程管理的串行指标数量，单位:条
	ReportCalcInterval          int64  `orm:"default(300)"` //报表计算扫描时间间隔,单位:秒
	AlarmScanInterval           int64  `orm:"default(3)"`   //报警计算扫描时间间隔,单位:秒
	AlarmMsgReloadInterval      int64  `orm:"default(600)"` //报警基础信息重载时间间隔,单位:秒
	Gomaxprocs                  int64  //执行程序的CPU核心数量，可以有如下几种数值：小于1，按最大CPU核心数并发执行；等于1，单核心执行； 大于1，按指定数量的CPU核心多核并发执行，但最大到CPU的核心数量。
	Description                 string //备注描述
	Seq                         int64  //排序序号
	Status                      int64  `orm:"default(1)"` //状态值,取值范围:0和1,0的时候本条配置无效,1的时候本条配置有效
	RunMinutes                  int64  //本次启动运行时间(分钟)
	TotalMinutes                int64  //总运行时间(分钟)
	Version                     string //程序版本号
	CreateUserId                int64  //创建者用户名
	CreateTime                  string `orm:"type(datetime)"` //创建时间
	UpdateUserId                int64  //最后编辑者用户名
	UpdateTime                  string `orm:"auto_now;type(datetime)"` //最后编辑时间
	ProjectName                 string //项目名称
	LogoPath                    string //Logo路径
	IcoPath                     string //Ico路径
	HeadBgPicPath               string //header背景图片路径
	PlatPath                    string //平台配置地址
	Copyright                   string //版权单位名称
}

//计算功能定义表（calc_function_dic)
type CalcFunctionDic struct {
	Id               int64  `orm:"auto"` //自动ID
	Function         string //功能
	InputParaDefine  string //输入数据说明
	OutputParaDefine string //输出数据说明
	Description      string //描述
	Seq              int64  //排序,默认等于id
	Status           int64  //状态，1有效，0无效
	CreateUserId     int64  //创建者id
	CreateTime       string `orm:"type(datetime)"` //创建时间
	UpdateUserId     int64  //修改者id
	UpdateTime       string `orm:"type(datetime)"` //修改时间
}

//KPI指标定义表(calc_kpi_index_dic)
type CalcKpiIndexDic struct {
	Id                int64                          `orm:"auto"`                         //自动ID
	KpiTag            string                         `orm:"column(kpi_varible_name)"`     //KPI变量
	KpiNameEn         string                         `orm:"column(kpi_varible_name_eng)"` //英文名称
	KpiNameCn         string                         `orm:"column(kpi_varible_name_cn)"`  //中文名称
	TagType           string                         //变量标签类型,process,goods等
	Script            string                         //计算脚本
	Description       string                         //含义描述
	Seq               int64                          //排序,默认等于id
	Status            int64                          //状态，1有效，0无效
	CreateUserId      int                            //创建者用户名
	CreateTime        string                         `orm:"type(datetime)"` //创建时间
	UpdateUserId      int                            //最后编辑者用户名
	UpdateTime        string                         `orm:"type(datetime)"` //最后编辑时间
	CalcKpiConfigList []*CalcKpiConfigList           `orm:"reverse(many)"`  // 设置一对多的反向关系
	Variblesets       []*MineDDictionaryOfVaribleset `orm:"rel(m2m)"`       // 设置与VariableSet表的多对多关系
}

//计算指标配置表(calc_kpi_config_list)
type CalcKpiConfigList struct {
	Id              int64            `orm:"auto"` //自动ID
	DistributedId   int64            //分布式计算ID，为0的时候不区分计算引擎，不为0的时候需要id适配的计算引擎进行计算
	TagType         string           //变量标签类型,process,goods等
	TagId           int64            //tag变量的id,与tag_type指定的taglist表相关联
	CalcKpiIndexDic *CalcKpiIndexDic `orm:"column(kpi_index_dic_id);rel(fk)"` //kpi指标的id，与calc_kpi_index_dic的id关联。如果设置为0，则可以在script列中自定义计算脚本
	KpiTag          string           //tag变量的ore_process_d_taglist.tage_name 与 period、 offset_minute、 calc_kpi_define.name_en拼接而成的字符串,中间以“.”连接。kpi_index_dic_id=0的时候，允许自行编辑定义该标签。该标签不可重复
	KpiName         string           //kpi指标的名称，默认由oreprocess_d_taglist.tag_description 与calc_kpi_define.description加“”拼接而成，kpi_index_dic_id=0的时候，允许用户自行编辑定义
	Script          string           `orm:"type(text)"` //自定义脚本。kpi_index_dic_id=0的时候有效。 脚本语法规则见[计算服务引擎脚本语法]
	StartTime       string           //统计计算开始起作用的时间。只有当该时间小于当前时间时，计算才进行；该时间大于当前时间时，计算不进行。
	Period          int64            //计算周期,可选：0=停止,-1=每分钟,-2=每小时,-3=每班,-4=每日,-5=每月,-6=每季度,-7=每年。大于零时，为自由周期的时间值，单位为秒。其小于0时，计算的基准时间为[ore_process_d_workshop]中的[base_time（datatime）]字段，按多任务的形式串行执行；其大于0时，计算的基准时间为当前行的start_time，并发执行。
	OffsetMinutes   int64            //计算偏移时间，单位：分钟，取值范围 -120~120。period选择为3(每班)的时候必填，其他选项的时候不需要
	LastCalcTime    string           //最后一次计算的时间，由程序自动填写，严禁手工编辑。程序启动时，如果其为空，自动从start_time开始计算至当前时间，每计算一个周期更新一次该值；如果启动时该时间与当前时间的差超过了一个计算周期，自动循环执行至当前时间，同步更新该值
	Supplement      int64            //补充计算开关。当其为1时，快速从start_time计算至当前时间，计算结果覆盖结果表中的同时间戳的项目，随后将该值置为0
	Description     string           //备注描述
	Seq             int64            //排序序号
	Status          int64            //状态值,取值范围:0和1,0的时候本条配置无效,1的时候本条配置有效
	CreateUserId    int              //创建者用户名
	CreateTime      string           `orm:"type(datetime)"` //创建时间
	UpdateUserId    int              //最后编辑者用户名
	UpdateTime      string           `orm:"type(datetime)"` //最后编辑时间
	KpiBaseTime     string           `orm:"type(datetime)"` //自定义基准时间
	KpiShiftHour    int64            //自定义每班工作时间
	//CfgUnitId       int64            //单位ID
	CfgUnit string //单位
	//SelectStartTime time.Time //为预判型故障加的这个字段，当用预判型故障相关联的KPI指标去kpi_result表查询时作为calc_ending_time的筛选条件，查询本时间以后的计算结果，默认为空，维修后更新该时间
}

//计算指标配置表(视图)(calc_kpi_config_list)
type CalcKpiConfigListExi struct {
	Id            int64  `orm:"auto"` //自动ID
	DistributedId int64  //分布式计算ID，为0的时候不区分计算引擎，不为0的时候需要id适配的计算引擎进行计算
	TagType       string //变量标签类型,process,goods等
	TagId         int64  //tag变量的id,与tag_type指定的taglist表相关联
	KpiIndexDicId int64  //kpi指标的id，与calc_kpi_index_dic的id关联。如果设置为0，则可以在script列中自定义计算脚本
	KpiTag        string //tag变量的ore_process_d_taglist.tage_name 与 period、 offset_minute、 calc_kpi_define.name_en拼接而成的字符串,中间以“__”连接。kpi_index_dic_id=0的时候，允许自行编辑定义该标签。该标签不可重复
	KpiName       string //kpi指标的名称，默认由oreprocess_d_taglist.tag_description 与calc_kpi_define.description加“__”拼接而成，kpi_index_dic_id=0的时候，允许用户自行编辑定义
	Script        string `orm:"type(text)"` //自定义脚本。kpi_index_dic_id=0的时候有效。 脚本语法规则见[计算服务引擎脚本语法]
	StartTime     string //统计计算开始起作用的时间。只有当该时间小于当前时间时，计算才进行；该时间大于当前时间时，计算不进行。
	Period        int64  //计算周期,可选：0=停止,-1=每分钟,-2=每小时,-3=每班,-4=每日,-5=每月,-6=每季度,-7=每年。大于零时，为自由周期的时间值，单位为秒。其小于0时，计算的基准时间为[ore_process_d_workshop]中的[base_time（datatime）]字段，按多任务的形式串行执行；其大于0时，计算的基准时间为当前行的start_time，并发执行。
	OffsetMinutes int64  //计算偏移时间，单位：分钟，取值范围 -120~120。period选择为3(每班)的时候必填，其他选项的时候不需要
	LastCalcTime  string //最后一次计算的时间，由程序自动填写，严禁手工编辑。程序启动时，如果其为空，自动从start_time开始计算至当前时间，每计算一个周期更新一次该值；如果启动时该时间与当前时间的差超过了一个计算周期，自动循环执行至当前时间，同步更新该值
	Supplement    int64  //补充计算开关。当其为1时，快速从start_time计算至当前时间，计算结果覆盖结果表中的同时间戳的项目，随后将该值置为0
	Description   string //备注描述
	Seq           int64  //排序序号
	Status        int64  //状态值,取值范围:0和1,0的时候本条配置无效,1的时候本条配置有效
	CreateUserId  int    //创建者用户名
	CreateTime    string `orm:"type(datetime)"` //创建时间
	UpdateUserId  int    //最后编辑者用户名
	UpdateTime    string `orm:"type(datetime)"` //最后编辑时间
	KpiBaseTime   string `orm:"type(datetime)"` //自定义基准时间
	KpiShiftHour  int64  //自定义每班工作时间
	//CfgUnitId     int64  //单位ID
	CfgUnit string //单位

	TableName string //实时变量所在数据库表的表名
	TagName   string //实时变量的变量名
	BaseTime  string `orm:"type(datetime)"` //基准时间
	ShiftHour int64  //每班工作时间
	KpiScript string //KPI指标的脚本
	KpiKey    string //KPI的关键词,对应 kpi_varible_name
}

type CalcKpiResult struct {
	Id              int64   `orm:"auto"` //自动id
	TagType         string  //变量标签类型,process,goods等
	TagId           int64   //tag变量的id,与tag_type指定的taglist表相关联
	TagName         string  //原始tag在原始taglist表中的tag_name
	KpiConfigListId int64   //指标id,指标在calc_indicator_config_list中的id
	KpiKey          string  //Kpi关键词,与kpi_index_dic中的 kpi_varible_name 对应
	KpiTag          string  //指标名称
	KpiName         string  //指标中文名称
	KpiPeriod       int64   //指标计算周期：-1=每小时,-2=每班,-3=每日,-4=每月,-5=每季度,-6=每年。大于零时，为自由周期的时间值，单位为秒。
	KpiValue        float64 //指标的值
	CalcEndingTime  string  `orm:"type(datetime)"` //指标对应的时间。该时间是指标统计期末的时间值
	InbaseTime      string  `orm:"type(datetime)"` //当前记录生成时的时间
}

type CalcKpiReportList struct { //报表定义列表
	Id            int64                `orm:"auto"` //ID
	DistributedId int64                //分布式计算ID，为0的时候不区分计算引擎，不为0的时候需要id适配的计算引擎进行计算
	Name          string               //名称
	Pid           int64                //父级菜单ID
	Workshop      *OreProcessDWorkshop `orm:"rel(fk)"` //所属车间ID
	Level         int64                //层级深度
	LevelCode     string               //层级码
	Folder        int64                //是否是文件夹，1-是，0-否
	Debug         int64                //调试模式,1:调试模式,0:费调试模式
	Seq           int64                //排序号
	Remark        string               //备注
	TemplateUrl   string               //模板文件路径
	TemplateFile  string               //模板文件名称
	ResultUrl     string               //结果地址
	StartTime     string               `orm:"type(datetime)"` //统计计算开始起作用的时间
	Period        int64                //计算周期,详见KPI表
	OffsetMinutes int64                //偏移时间
	LastCalcTime  string               `orm:"type(datetime)"` //最后计算时间
	BaseTime      string               `orm:"type(datetime)"` //基准时间
	ShiftHour     int64                //每班工作时间
	CreateUserId  int64                //创建人ID
	CreateTime    string               `orm:"type(datetime)"`          //创建时间
	UpdateTime    string               `orm:"auto_now;type(datetime)"` //更新时间
	Status        int64                `orm:"default(1)"`              //1有效 0无效
	Permissions   []*SysPermission     `orm:"reverse(many)"`           //设置与权限表的反向多对多关系
}

//变量报警信息表
type CalcTagAlarm struct {
	Id          int64               `orm:"auto"`           //自动id
	Datatime    string              `orm:"type(datetime)"` //当前记录生成时的时间
	Tag         *OreProcessDTaglist `orm:"rel(fk)"`        //tag变量的id
	TagDesc     string              //变量描述
	TagShotDesc string              //变量短描述
	TagValue    float64             //指标的值
	LimitValue  float64             //报警限
	AlarmStatus int                 //报警状态:-3=小于低限;-2=小于低低报警值;-1=小于低报警值;0=无报警;1=大于高报警值;2=大于高高报警值;3=大于上限
}

//变量列表
type OreProcessDTaglist struct {
	Id       int64                        `orm:"auto"`    //	自动编号。
	Variable *MineDDictionaryOfVaribleset `orm:"rel(fk)"` //	mine_d_dictionary_of_varibleset 表的 id
	//LevelCategory int64                        //
	ResourceType  int64    //MOTOR(1,"电机"),METER(2,"仪表"),SINGLE_ANALYZER(3,"单流道分析仪器"),DEVICE(4,"设备"),ASSAY(5,"取样点"),METERING(6,"电能计量表"),PLANT(7,"供变配电设备");
	ItemIDinTable int64    `orm:"column(resource_id)"` //关联表的ID，比如存入的记录设备的，name这条table_id存储的就是设备表的ID
	TreeLevelCode string   //	关联ms_middle_correlation.tree_level_code,tree_level_code为树形节点层级id拼接而成
	Dcs           *MineDcs `orm:"rel(fk)"` //	自动化系统 id
	//DcsStationId    int64                        //DCS站ID
	//DcsModuleId     int64                        //DCS模块ID
	//DcsModuleAdd    string                       //DCS模块地址
	FirstTagClass string                //	tag的varible类型
	Stage         *OreProcessDWorkstage `orm:"rel(fk)"` //	作业id
	DeviceId      int64                 //	设备ID
	EquipmentName string                //	设备名称
	Seq           int64                 //	排序
	Status        string                //	基本本条记录的有效性
	TagId         string                //	标签的id，由标签的数据源实体的id组合而成。tag的前缀，即第一位在mine_dictionary_tag_type表中可见
	//CheckStatus     string                //	记录该条记录是否被放入已选列表，1-放入，0-未放入
	//RegistStatus    string                //	登记状态，1-已登记，0-未登记
	//TagCheckinTime  string                `orm:"type(datetime)"` //	标签登记时间，记录标签登记的时间。标签登记后，才能进行数据采集策略配置。
	//TagCheckoutTime string                `orm:"type(datetime)"` //	标签注销时间，记录注销标签的时间。标签注销后，不能对其进行数据采集功能的配置。
	//CreateTime        string `orm:"type(datetime)"` //创建时间
	//CreateUserId      int64  //创建者id
	//UpdateTime        string `orm:"type(datetime)"` //更新时间
	//UpdateUserId      int64  //更新者
	DecimalNum int64 //	小数点位数，控制页面每条记录中所有浮点数字段的显示位
	//MiningClass string //	矿业类别：1-选矿厂，2-采矿厂
	GoldenId   int    `orm:"column(alias_tag_id)"` //点表在庚顿数据库中的ID
	TagPlcName string //变量在PLC中的名称
	TagHmiName string //变量在HMI中的名称
	//IoDataType              int     //	数值类型，整理IO变量时填的数值类型
	Valid                   int64   //	IO变量有效状态，-1无效,0未采集,1有效
	NormalValue             float64 //	正常值
	LimitLl                 float64 //	低低报警值
	LimitL                  float64 //	低报警值
	LimitH                  float64 //	高报警值
	LimitHh                 float64 //	高高报警值
	MinValue                float64 `orm:"column(meas_minvalue)"` //量程下限
	MaxValue                float64 `orm:"column(meas_maxvalue)"` //量程上限
	Scale                   float64 //	比例。比例和偏移量用于对采集到的原始值进行修正,y=x*比例+偏移量。比例默认为1
	Offset                  float64 //	偏移量。比例和偏移量用于对采集到的原始值进行修正,y=x*比例+偏移量。偏移量默认为0
	TagFullName             string  `orm:"column(io_comment)"` //借用io_comment用作变量全名:tablename.tagname
	TagName                 string  `orm:"column(tage_name)"`  //	标签名组合，由标签的数据源实体的编号组合而成。tag的前缀，即第一位在mine_dictionary_tag_type表中可见
	TagType                 string  `orm:"column(value_type)"` //	数值类型
	TagDescription          string  //	标签描述
	TagUnit                 string  `orm:"column(tag_varible_unit)"` //	tag的varible工程单位,对应varibleSet的varible_unit
	TagPracticalDescription string  //变量的实际描述

	Unit       string `orm:"column(unit)"`        //	测量值的工程单位：
	IsArchive  int    `orm:"column(keep)"`        //	存档
	Digits     int    `orm:"column(value_digit)"` //	数值位数
	IsShutDown int    `orm:"column(halt_write)"`  //	停机补写开关;
	//LowLimit       float64 `orm:"column(meas_maxvalue)"`                //	量程上限
	//HighLimit      float64 `orm:"column(meas_minvalue)"`                //	量程下限
	IsStep         int     `orm:"column(phase_step)"`                   //	阶跃开关；
	Typical        float64 `orm:"column(typical_value)"`                //	典型值；
	IsCompress     int     `orm:"column(compress)"`                     //	压缩开关；
	CompDev        float64 `orm:"column(compress_deviation)"`           //	压缩偏差;
	CompDevPercent float64 `orm:"column(compress_percentage)"`          //	压缩百分比;
	CompTimeMax    int     `orm:"column(longest_compress_interval)"`    //	最长压缩间隔；
	CompTimeMin    int     `orm:"column(shortest_compress_interval)"`   //	最短压缩间隔
	ExcDev         float64 `orm:"column(exception_deviation)"`          //	例外偏差
	ExcDevPercent  float64 `orm:"column(exception_percentage)"`         //	例外百分比
	ExcTimeMax     int     `orm:"column(longest_exception_deviation)"`  //	最长例外间隔：单位：秒
	ExcTimeMin     int     `orm:"column(shortest_exception_deviation)"` //	最短例外间隔，单位：秒
	ClassOf        int     `orm:"column(category)"`                     //	类别，基本点—base=0;采集点—base,sacn=1;计算点—base,calc=2;采集计算点—base,scan,calc=3
	//ChangeDate     string  `orm:"column(point_update_time)"`            //	点表中标签点的更新时间
	//Changer        string  `orm:"column(update_user)"`                  //	修改者
	//CreateDate     string  `orm:"column(point_create_time)"`            //	点表中标签点的创建时间
	//Creator        string  `orm:"column(create_user)"`                  //	创建者
	Mirror      int     `orm:"column(mirror)"`              //	镜像模式
	MilliSecond int     `orm:"column(timestamp_accuracy)"`  //	时间戳精度
	IsSummary   int     `orm:"column(speed_up)"`            //	加速统计
	Source      string  `orm:"column(data_origin)"`         //	数据源
	IsScan      int     `orm:"column(collect)"`             //	采集开关
	Instrument  string  `orm:"column(device_tag)"`          //	设备标签
	Location1   int     `orm:"column(location1)"`           //	五个位地址之一
	Location2   int     `orm:"column(location2)"`           //	五个位地址之一
	Location3   int     `orm:"column(location3)"`           //	五个位地址之一
	Location4   int     `orm:"column(location4)"`           //	五个位地址之一
	Location5   int     `orm:"column(location5)"`           //	五个位地址之一
	UserInt1    int     `orm:"column(user_int1)"`           //	供用户使用的两个32位整数之一。
	UserInt2    int     `orm:"column(user_int2)"`           //	供用户使用的两个32位整数之一。
	UserReal1   float64 `orm:"column(user_real1)"`          //	供用户使用的两个32位浮点数之一。
	UserReal2   float64 `orm:"column(user_real2)"`          //	供用户使用的两个32位浮点数之一。
	Equation    string  `orm:"column(equation)"`            //	方程式
	Trigger     int     `orm:"column(trigger_way)"`         //	触发方式，用0-3表示：0—无触发；1-事件触发；2—周期触发；3-定时触发；
	TimeCopy    int     `orm:"column(timestamp_reference)"` //	时间戳参考，0—使用计算发生时间；1—使用最晚标签点时间；2—使用最早标签点时间
	Period      int     `orm:"column(calcu_period)"`        //	计算周期,单位：秒

	TagAlarm []*CalcTagAlarm `orm:"reverse(many)"` //变量报警信息
}

//工序级别信息表
type OreProcessDWorkstage struct {
	Id                  int                   `orm:"auto"`    //工艺自动编号。
	Workshop            *OreProcessDWorkshop  `orm:"rel(fk)"` //工艺所属车间的id编号，和ore_process_d_workshop.id关联。
	OrganizationOrgCode int                   //sys_organization表org_code，用于数据权限
	ConstructionCode    int                   //工序作业类别，该字段的赋值由constructionlist的选矿车间以下所有第4级实体作为赋值选项。
	ConstructionName    string                //工序作业类别名称，该字段的赋值由constructionlist的选矿车间以下所有第4级实体作为赋值选项。
	StageName           string                //工序作业名称
	StageNameCode       string                //作业编号
	StageStatus         int                   //工序作业的状态，1：生产；2：停产；3：搁置。
	StageTechpara       string                //生产作业工序的专有技术参数，该字段存储相应工艺技术参数的数据表名。
	Seq                 int                   //序号
	CreateUserId        int                   //记录创建人的用户id，关联表待确认。
	CreateTime          string                //记录创建时间。
	UpdateUserId        int                   //记录更新人的用户id，关联表待确认。
	UpdateTime          string                //记录更新时间。
	Status              int                   //本条记录的有效性 1有效 0无效
	UserNameCode        string                //用户编号
	GraphServiceUrl     string                //外部图形页面URL
	TableServiceUrl     string                //外部报表页面URL
	TreeLevelCode       string                //中间表树形结构ID拼接
	Taglists            []*OreProcessDTaglist `orm:"reverse(many)"` // 设置一对多的反向关系
	CheckTagLists       []*CheckTagList       `orm:"reverse(many)"` //巡检标签定义表设置一对多的反向关系
	Samples             []*SamplingManage     `orm:"reverse(many)"` //
}

//车间级信息表
type OreProcessDWorkshop struct {
	Id                      int                      `orm:"auto"`    //表ID
	OreProcessConcentration *OreProcessConcentration `orm:"rel(fk)"` //车间所属工厂厂的id号，与ore_process_concentration.id关联。
	OrganizationOrgCode     int                      //关联 sys_organization.org_code 用于数据权限
	ConstructionCode        int                      //关联mine_constructionlist.level_num_code字段，
	WorkshopName            string                   //生产车间的名称，以厂矿管理系统的命名为准。
	WorkshopNameCode        string                   //车间编号，手工设定
	DayProcess              string                   //车间的日处理量，单位：t/d。
	WorkerNum               int                      //车间人员数量，单位：人，可用于计算人均劳动生产率等指标
	ShiftHour               int                      //车间每班工作额定时长，单位：小时。
	ShifttimesDay           int                      //车间一个工作日的班次数量，单位：次。
	Seq                     int                      //排序序号
	Status                  int                      //本条记录的有效性，1有效 0无效
	Remarkes                string                   //备注
	CreateUserId            int                      //记录创建人的用户id，关联表待确认。
	CreateTime              string                   `orm:"type(datetime)"` //记录创建时间。
	UpdateUserId            int                      //记录更新人的用户id，关联表待确认。
	UpdateTime              string                   `orm:"type(datetime)"` //记录更新时间。
	UserNameCode            string                   //用户编码
	ShiftHourBias           int                      //时差，单位：h
	BaseTime                string                   `orm:"type(datetime)"` //定义统计计算的基准时间
	GraphServiceUrl         string                   //外部图形页面URL
	TableServiceUrl         string                   //外部报表图形页面URL
	TreeLevelCode           string                   //中间表树形结构ID拼接
	Workstages              []*OreProcessDWorkstage  `orm:"reverse(many)"` // 设置一对多的反向关系
	GoodsLists              []*GoodsConfigInfo       `orm:"reverse(many)"` // 设置一对多的反向关系
	Reports                 []*CalcKpiReportList     `orm:"reverse(many)"` // 设置一对多的反向关系
}

//厂级信息表
type OreProcessConcentration struct {
	Id                       int                    `orm:"auto"`    //编号。
	MineBasicInfo            *MineBasicInfo         `orm:"rel(fk)"` //关联mine_basic_info.id字段，用于所属矿山的相关信息。
	ConstructionCode         int                    //关联mine_construction_list.level_num_code。
	PlantName                string                 //工厂的名称，如泗州选矿厂，按照矿山企业管理命名执行。
	PlantNameCode            string                 //工厂名称的编号，只能由平台的后台工作人员来设定。此字段只能由平台的后台管理员来设置。
	PlantSiteAtitude         int                    //工厂的海拔高度，单位：m；该指标对选矿厂的能源监控非常有用。字段设置为4位，即最大海拔高度是9999m。
	OreProcessPerday         int                    //工厂的矿石日处理量，单位：t/d。字段长度设置为6位，即日处理量最大999999t/d。
	WorkdayPeryear           int                    //工厂的年度额定工作时长，单位：天，最多365天。字段长度设定为3位。
	WorkerTotalNum           int                    //工厂的生产定员总数，单位：人，可用于计算人均劳动生产率等指标。字段长度设置为4位，即最多是9999人。
	PowerSupplyInlineVoltage int                    //工厂的总进线电压，单位：V。字段长度设置为6，即最大进线电压999.999KV。
	PowerSupplyTotalload     int                    //工厂供电的总进线功率，指额定功率，单位：KW。字段长度设置成9位，即999999999KW。
	ProductionUnitNum        int                    //工厂管理的生产车间数量，单位：个；字段长度设定成2，即加工厂最多有99个车间。
	PlantIsIndoor            int                    //工厂是室内还是室外的？1：室内；0：室外。
	ConcentrationNum         int                    //工厂的精矿产品数量，单位：个。字段长度设置为1，即精矿产品数量最多是9个。
	ConcentrationClass       string                 //精矿产品的种类。
	MetalElementNum          int                    //生产关注的金属主元素数量。字段长度设置为1，即金属元素数量最多是9个。
	MetalElementClass        string                 //生产关注的金属元素种类。
	Seq                      int                    //排序序号
	Status                   int                    //本条记录有效性，1有效，0无效
	CreateUserId             int                    //记录创建人的用户id，关联表待确认。
	CreateTime               string                 `orm:"type(datetime)"` //记录创建时间。
	UpdateUserId             int                    //记录更新人的用户id，关联表待确认。
	UpdateTime               string                 `orm:"type(datetime)"` //记录更新时间。
	UserNameCode             string                 //用户编码
	OreOutputPerday          int                    //每日可利用矿石产量，单位:t/d
	SpoilProductionPerday    int                    //每日废弃的矿石量，单位:t/d
	PlantClass               int                    //工厂类别 1：选矿厂 2：湿法厂 3：地下采矿厂 4：露天采矿厂
	GraphServiceUrl          string                 //外部图形页面URL
	TableServiceUrl          string                 //外部报表页面URL
	TreeLevelCode            string                 //中间表树形结构ID拼接
	Workshops                []*OreProcessDWorkshop `orm:"reverse(many)"` // 设置一对多的反向关系
}

//矿山信息表
type MineBasicInfo struct {
	Id                   int                        `orm:"auto"` //矿山企业的自动编号。
	CorparationInfoId    int                        //关联corparation_info.id字段，用于查询矿山所属集团公司的相关信息。
	OrganizationOrgCode  int                        //sys_organization表org_code，用于数据权限。
	ConstructionCode     int                        //关联mine_constructionlist.level_num_code
	FormalName           string                     //矿山官方名称，以矿山门户网站名称为准，如江铜集团城门山铜矿
	FormalNameShort      string                     //矿山名称简写编码，手工设定，编码原则：中文名称以名称的拼音首字母代替，如城门山铜矿，编码为CMS。
	Address              string                     //矿山的所在地区和地址，如城门山铜矿位于江西生产九江市柴桑区城门山镇。
	IsUndergroundMining  int                        //矿山是否包括地下采矿，有，则值为1；无，则值为0。
	UndergroundMiningNum int                        //矿山拥有的地下采矿场的数量，如果undergroud_mining==1,则该字段的数值至少为1；否则为0。
	IsOpenpitMining      int                        //矿山是否包括露天采矿，有，则值为1；无，则值为0。
	OpenpitMiningNum     int                        //矿山拥有的露天采矿场的数量，如果openpit_mining==1,则该字段的数值至少为1；否则为0。
	IsConcentrator       int                        //矿山是否包括选矿厂，有，则值为1；无，则值为0。
	ConcentratorNum      int                        //矿山拥有的选矿厂的数量，如果mineral_concentrator==1,则该字段的数值至少为1;否则为0。
	IsHydrometallurgy    int                        //矿山是否包括湿法冶炼厂，有，则值为1；无，则值为0。
	HydrometallurgyNum   int                        //矿山拥有的湿法冶炼厂的数量，如果mineral_hydrometallurg==1,则该字段的数值至少为1；否则为0。
	IsTailingReservior   int                        //矿山是否有尾矿库，有，则值为1；无，则值为0
	TailingReserviorNum  int                        //矿山拥有的尾矿库的数量，如果tailing_reservior==1,则该字段的数值至少为1；否则为0。
	Seq                  int                        //排序序号
	OperatorId           int                        //对于矿业集团公司，平台制定一名信息维护人员（operator）对该集团的所属信息进行维护。此字段关联mine_group_info.operator_id；矿山必须和集团共用一个信息维护员。
	OperatingTime        string                     `orm:"type(datetime)"` //信息维护人员（operator）进行有效操作的时间，此操作时间以信息维护人员登出的时间为准。
	OperationClass       string                     //信息维护人员操作类别：1：创建和修改相关账号信息；2：冻结集团相关账号；3：删除集团相关账号；4：未进行任何操作。
	CreateUserId         int                        //记录创建人的用户id，关联表待确认。默认为北矿智云科技的后台维护人员。
	CreateTime           string                     `orm:"type(datetime)"` //记录创建时间。
	UpdateUserId         int                        //记录更新人的用户id，关联表待确认。默认为北矿智云科技的后台维护人员。
	UpdateTime           string                     `orm:"type(datetime)"` //记录更新时间。
	Status               int                        //本条记录是否有效，只有具有删除权限的人可以改变该记录的状态，状态为1：有效，状态为0：无效。查询时必须附带查询条件status=1/status=0。
	UserNameCode         string                     //用户编码
	GraphServiceUrl      string                     //外部图形页面URL
	TableServiceUrl      string                     //外部报表页面URL
	TreeLevelCode        string                     //中间表树形结构ID拼接
	Plants               []*OreProcessConcentration `orm:"reverse(many)"` // 设置一对多的反向关系
	MineDeptDics         []*MineDeptDic             `orm:"reverse(many)"` // 设置一对多的反向关系
}

//层级中间表ms_middle_correlation
type MsMiddleCorrelation struct {
	Id                   int64            `orm:"auto"`                           //本表ID
	Pid                  int64            `orm:"column(pId)"`                    //父ID
	LevelCategory        *MineTableList   `orm:"rel(fk);column(level_category)"` //属性所在表的表ID
	ItemIdInTable        int64            `orm:"column(table_id)"`               //属性在其所在表中的id值
	ConstrutionCode      int64            //(未使用)关联constructionList.level_num_code字段
	ConstrutionTableCode string           //(未使用)关联constructionList.level_mine_table_id字段
	LevelName            string           //结构树名称，即本层级的实际名称，与id关联的本表的名称属性值一致
	Seq                  int64            `orm:"column(sortNum_in_level)"` //本层级排序号
	LevelNum             int64            //结构树的层级序号
	Status               int64            //状态 0无效。 1有效
	TreeLevelCode        string           //树层级节点ID拼接
	Permissions          []*SysPermission `orm:"reverse(many)"` //设置与权限表的多对多关系
}

//中间表合成树结构缓存字符串
type MsMiddleCorrelationString struct {
	Id        int64  `orm:"auto;column(ID)"`                  //自动ID
	IdCode    string `orm:"column(IdCode)"`                   //标识码
	MenuValue string `orm:"column(menuValue);type(text)"`     //字符串
	EditTime  string `orm:"column(Edit_Time);type(datetime)"` //编辑时间
	RootPid   int64  //根节点的父ID
}

//平台数据表列表(mine_table_list)
type MineTableList struct {
	Id               int64                  //id
	TableProgramName string                 //表名
	LevelCategory    int64                  //表的代号，一旦定义好，禁止修改。
	Middles          []*MsMiddleCorrelation `orm:"reverse(many)"` // 设置一对多的反向关系
}

//变量类型表(mine_d_dictionary_of_varibleset)
type MineDDictionaryOfVaribleset struct {
	Id                 int64                 `orm:"auto"` //
	VaribleType        int64                 //0是全部可以看，1代表电机，2代表仪表及单流道分析仪器 ，3 代表多流道分析仪器，设备全部显示0-3///////设备电机都可以看是4。供电线路5，供配电设备是6
	EquipDicId         int64                 //关联equipmentdic.id，值不为0时代表这个小类型专属变量
	EquipCode          int64                 //关联equipmentdic.equipcode
	VaribleKeyClass    string                //类型名称
	VaribleName        string                //变量名称
	VaribleDescription string                //变量物理含义的描述。
	VaribleValueDesc   string                //变量值描述
	Seq                int64                 `orm:"column(varible_seq)"`    //顺序
	Status             int64                 `orm:"column(varible_status)"` //是否有效,0无效，1有效
	ValueType          string                //变量数据类型
	PointValueType     string                //庚顿点表的数值类型
	VaribleUnit        string                //变量的工程单位。
	FirstTagClass      string                //变量类型
	CalcKpiId          string                //默认统计计算方式
	KpiIndexDics       []*CalcKpiIndexDic    `orm:"reverse(many)"` // 设置一对多的反向关系
	Taglists           []*OreProcessDTaglist `orm:"reverse(many)"` // 设置一对多的反向关系
}

//物耗池
type GoodsConsumePool struct {
	Id                 int64                   `orm:"auto"` //主键自增ID
	GoodsName          string                  //物品名称
	DeptManage         *MineDeptInfo           `orm:"rel(fk)"`                          //管理部门
	GoodsType          *OreProcessEquipmentDic `orm:"rel(fk);column(equipment_dic_id)"` //样本二级分类ID(物耗子类型)
	GoodsNameCode      string                  //平台编码
	UserNameCode       string                  //用户编码
	Specifications     string                  //规格型号
	GoodsConsumeUnit   string                  //消耗计量单位
	GoodsForm          int64                   //物品形态(1气体，2液体，3粉末，4固态，5胶体)
	PackingMode        string                  //包装方式
	PackingUnit        string                  //包装单位
	PackingModel       string                  //包装规格
	PreserveCondition  string                  //存放条件
	TransportCondition string                  //运输条件
	TransportMode      string                  //运输方式
	GoodsIsDanger      int64                   //是否是属于危害品(true:是，false:否)
	FactorySerialNum   string                  //出厂编号
	ProductionBatchNum string                  //生产批次号
	ManufacturerName   string                  //生产厂家
	SupplierName       string                  //供应商
	ConsumeUnitPrice   string                  //参考单价
	ScalingFactor      string                  //换算比例
	Remark             string                  //备注
	Seq                int64                   //排序号
	CreateUserId       int64                   //创建人
	CreateTime         string                  `orm:"type(datetime)"` //创建时间
	UpdateUserId       int64                   //更新人
	UpdateTime         string                  `orm:"type(datetime)"` //更新时间
	Status             int64                   //有效性
	GoodsLists         []*GoodsConfigInfo      `orm:"reverse(many)"` // 设置一对多的反向关系
}

//物资-工厂层级关系配置表
type GoodsConfigInfo struct { //物耗配置表
	Id                int64                `orm:"auto"`    //ID
	Goods             *GoodsConsumePool    `orm:"rel(fk)"` //物品ID
	ResourceType      int64                //层级类型
	ResourceId        int64                //层级主体ID
	TreeLevelCode     string               //树形结构节点编号
	Workshop          *OreProcessDWorkshop `orm:"rel(fk)"` //厂级ID
	PlanDosage        string               //计划用量
	ConsumptionMode   string               //消耗方式
	GoodsTagName      string               //按照各个层级生成的tag
	Remark            string               //备注
	Status            int64                //有效性
	Seq               int64                //排序号
	CreateUserId      int64                //创建人
	CreateTime        string               `orm:"type(datetime)"` //创建时间
	UpdateUserId      int64                //更新人
	UpdateTime        string               `orm:"type(datetime)"` //更新时间
	GoodsConsumeInfos []*GoodsConsumeInfo  `orm:"reverse(many)"`  // 设置一对多的反向关系
}

//物资消耗记录表
type GoodsConsumeInfo struct {
	Id                 int64            //主键自增ID
	GoodsConfigInfo    *GoodsConfigInfo `orm:"rel(fk)"` //物耗配置表
	GoodsConsumeAmount float64          //物资消耗量(用户填写)
	UseStartTime       string           `orm:"type(datetime)"` //物资开始消耗时间（用户填写）
	UseEndTime         string           `orm:"type(datetime)"` //物资结束消耗时间（用户填写）
	Remark             string           //备注
	CreateUser         *SysUser         `orm:"rel(fk)"`        //创建人
	CreateTime         string           `orm:"type(datetime)"` //创建时间
	UpdateUser         *SysUser         `orm:"rel(fk)"`        //更新人
	UpdateTime         string           `orm:"type(datetime)"` //更新时间
}

type SysRealData struct { //平台中的慢动态数据
	Id            int64          `orm:"auto"` //主键自增ID
	TagName       string         //标签名
	SysDictionary *SysDictionary `orm:"rel(fk)"` //标签类型字典ID，对应 sys_dictionary 表的 id
	TagId         int64          //标签所在表的ID,例如*_tag_list 表的 ID
	TagExeId      string         //标签的执行记录所在表的ID，例如 *_item_exe 表的ID
	Value         float64        //值
	Datatime      string         `orm:"type(datetime)"` //数据所代表的标签的时间
}

type SysRealDataExi struct { //平台中的慢动态数据（非关联,保存结果用）
	Id              int64   //主键自增ID
	TagName         string  //标签名
	SysDictionaryId int64   //标签类型字典ID，对应 sys_dictionary 表的 id
	TagId           int64   //标签所在表的ID,例如*_tag_list 表的 ID
	TagExeId        string  //标签的执行记录所在表的ID，例如 *_item_exe 表的ID
	Value           float64 //值
	Datatime        string  //数据所代表的标签的时间
}

type SysDictionary struct { //系统字典表
	Id                 int64  `orm:"auto"` //
	DicCatalogId       int64  //*SysDictionaryCatalog `orm:"rel(fk)"` //sys_dictionary_catalog表的id
	DictionaryCode     string //字典数据编号，如果不能使用ID，则可以使用。
	DictionaryNameCode string //字典数据英文名称。每个类型里唯一。
	Name               string //名称
	Seq                int64  //排序号
	Status             int64  //0 无效 1 有效
	Remarks            string //
}

type SysDictionaryCatalog struct { //系统字典目录表
	Id      int64  `orm:"auto"` //
	Name    string //名称
	Remarke string //说明
	Status  int64  //0 无效 1 有效
	//SysDictionarys []*SysDictionary `orm:"reverse(many)"` // 设置一对多的反向关系
}

//DCS和数据表关系表(relevance_dcs_to_dbtable)
type RelevanceDcsToDbtable struct {
	Id           int64           `orm:"auto"`
	Dcs          *MineDcs        `orm:"rel(one)"` //	自动化系统 id
	AcqStation   *DataAcqStation `orm:"rel(fk)"`  //采集机ID
	Datatable    *DatatableInfo  `orm:"rel(fk)"`  //数据表ID
	Seq          int64           //序号
	Remarks      string          //备注
	CreateUserId int64           //创建人id
	CreateTime   string          `orm:"type(datetime)"` //创建时间
	UpdateUserId int64           //修改人id
	UpdateTime   string          `orm:"type(datetime)"` //更新时间
	Status       int64
}

type MineDcs struct { //DCS系统表
	Id       int64                    `orm:"auto"`                     //
	Plant    *OreProcessConcentration `orm:"rel(fk);column(plant_id)"` //所在厂的ID
	DcsName  string                   //名称
	Seq      int64                    //排序号
	Status   int64                    //有效状态,1有效,0无效
	Taglists []*OreProcessDTaglist    `orm:"reverse(many)"`
	Reltable *RelevanceDcsToDbtable   `orm:"reverse(one)"`
}

type DatatableInfo struct { //实时数据库表信息表
	Id               int64  `orm:"auto"` //ID
	TableName        string //表名
	TableDescription string //表描述
	GoldenTableId    int    `orm:"column(technical_doc)"` //庚顿表ID(利用技术文件自动存储)
	//Status           int64             //有效状态,1有效,0无效
	AcqStations []*DataAcqStation `orm:"rel(m2m);rel_through(github.com/bkzy-wangjp/MicEngine/models.RelevanceDcsToDbtable)"`
}

type DataAcqStation struct { //数据采集机信息表
	Id                  int64            `orm:"auto"` //id自动生成
	AcqSysId            int64            //工作站所属数据采集系统id，外键
	AcqServerId         int64            //工作站连接服务器id
	AcqStationName      string           //数据采集工作站名称
	StationSite         string           //数据采集工作站安装位置
	ComputerName        string           //工作站的计算机名称
	LogAccount          string           //工作站登录账号
	LogPassword         string           //工作站登录密码
	ConstructionCode    int64            //工作站的层级编码，外键
	UserNameCode        string           //用户资产编码
	StationNameCode     string           //平台编码
	OrganizationOrgCode int64            //备用
	HardwareManufacture string           //工作站硬件供应商
	HardwareModel       string           //工作站硬件规格型号
	Mac1Add             string           //1#网卡MAC1地址
	Mac1Type            int64            //1#网卡类型 1有线网卡 2无线网卡
	Mac1InUse           int64            //1#网卡启用状态  true:启用 false:未启用
	Mac1LineNo          string           //1#网卡连接线路编号
	Mac1Ip1             string           //1#网卡IP1地址
	Mac1Ip1Mask1        string           `orm:"column(mac1_ip1_mask_1)"`    //1#网卡IP1掩码地址
	Mac1Ip1Gateway1     string           `orm:"column(mac1_ip1_gateway_1)"` //1#网卡IP1网关地址
	Mac1Ip2             string           //1#网卡IP2地址
	Mac1Ip2Mask1        string           `orm:"column(mac1_ip2_mask_1)"`    //1#网卡IP2掩码地址
	Mac1Ip2Gateway1     string           `orm:"column(mac1_ip2_gateway_1)"` //1#网卡IP2网关地址
	Mac1Ip3             string           //1#网卡IP3地址
	Mac1Ip3Mask1        string           `orm:"column(mac1_ip3_mask_1)"`    //1#网卡IP3掩码地址
	Mac1Ip3Gateway1     string           `orm:"column(mac1_ip3_gateway_1)"` //1#网卡IP3网关地址
	Mac2Add             string           //2#网卡MAC2地址
	Mac2Type            int64            //2#网卡类型
	Mac2InUse           int64            //2#网卡启用状态
	Mac2LineNo          string           //2#网卡连接线路编号
	Mac2Ip1             string           //2#网卡IP1地址
	Mac2Ip1Mask1        string           `orm:"column(mac2_ip1_mask_1)"`    //2#网卡IP1掩码地址
	Mac2Ip1Gateway1     string           `orm:"column(mac2_ip1_gateway_1)"` //2#网卡IP1网关地址
	Mac2Ip2             string           //2#网卡IP2地址
	Mac2Ip2Mask1        string           `orm:"column(mac2_ip2_mask_1)"`    //2#网卡IP2掩码地址
	Mac2Ip2Gateway1     string           `orm:"column(mac2_ip2_gateway_1)"` //2#网卡IP2网关地址
	Mac2Ip3             string           //2#网卡IP3地址
	Mac2Ip3Mask1        string           `orm:"column(mac2_ip3_mask_1)"`    //2#网卡IP3掩码地址
	Mac2Ip3Gateway1     string           `orm:"column(mac2_ip3_gateway_1)"` //2#网卡IP3网关地址
	Mac3Add             string           //3#网卡MAC3地址
	Mac3Type            int64            //3#网卡类型
	Mac3InUse           int64            //3#网卡启用状态
	Mac3LineNo          string           //3#网卡连接线路编号
	Mac3Ip1             string           //3#网卡IP1地址
	Mac3Ip1Mask1        string           `orm:"column(mac3_ip1_mask_1)"`    //3#网卡IP1掩码地址
	Mac3Ip1Gateway1     string           `orm:"column(mac3_ip1_gateway_1)"` //3#网卡IP1网关地址
	Mac3Ip2             string           //3#网卡IP2地址
	Mac3Ip2Mask1        string           `orm:"column(mac3_ip2_mask_1)"`    //3#网卡IP2掩码地址
	Mac3Ip2Gateway1     string           `orm:"column(mac3_ip2_gateway_1)"` //3#网卡IP2网关地址
	Mac3Ip3             string           //3#网卡IP3地址
	Mac3Ip3Mask1        string           `orm:"column(mac3_ip3_mask_1)"`    //3#网卡IP3掩码地址
	Mac3Ip3Gateway1     string           `orm:"column(mac3_ip3_gateway_1)"` //3#网卡IP3网关地址
	SoftwareManufacture string           //工作站软件供应商
	SoftwareName        string           //工作站软件名称
	SoftwareVersion     string           //工作站软件版本号
	OprSysName          string           //工作站操作系统
	OprSysVersion       string           //工作站操作系统版本
	OpcserverIsValid    int64            //工作站的OPC server是否可用 true可用，false未启用
	OpcserverName       string           //工作站的OPC server名称
	Seq                 int64            //排序
	PurchaseTime        string           `orm:"type(datetime)"` //采购时间
	StartUseTime        string           `orm:"type(datetime)"` //投入使用时间
	OriginalValue       string           //原值
	ScrapValue          string           //残值
	TechnicalDoc        int64            //技术文件
	QualificationDoc    int64            //产品合格证
	AssetStatus         int64            //资产状态
	Status              int64            //有效状态
	CreateUserId        int64            //创建者ID
	CreateTime          string           `orm:"type(datetime)"` //创建时间
	UpdateUserId        int64            //更新者ID
	UpdateTime          string           `orm:"type(datetime)"` //更新时间
	Datatables          []*DatatableInfo `orm:"reverse(many)"`  //设置一对多反向关系
}

type RealTimeMonitor struct { //数据监控配置表
	Id           int64            `orm:"auto"`              //ID
	Name         string           `orm:"column(title)"`     //层级名称
	Pid          int64            `orm:"column(parent_id)"` //父级菜单ID
	Level        int64            //层级深度
	Seq          int64            //排序号
	Remark       string           //备注
	Url          string           //grafana URL
	Folder       int64            //是否是文件夹，1-是，0-否
	Status       int64            //1有效 0无效
	CreateUserId int64            //创建人ID
	CreateTime   string           `orm:"type(datetime)"` //创建时间
	UpdateTime   string           `orm:"type(datetime)"` //更新时间
	Permissions  []*SysPermission `orm:"reverse(many)"`  //设置与权限表的反向多对多关系
}

//矿山部门词典
type MineDeptDic struct {
	Id              int64          `orm:"auto"`    //主键自增id
	Mine            *MineBasicInfo `orm:"rel(fk)"` //关联mine_basic_info.id字段，用于所属矿山的相关信息。
	OrgPCode        string         //部门人员编码
	OrgCode         string         //组织编码
	DepartmentClass string         //部门类型
	Seq             int64          //排序
	Status          int64          //有效状态 1：有效；0：无效
	//CreateTime      string          `orm:"type(datetime)"` //创建时间
	//CreateUserId    int64           //创建者id
	//UpdateTime      string          `orm:"type(datetime)"` //更新时间
	//UpdateUserId    int64           //更新者id
	DeptInfos []*MineDeptInfo `orm:"reverse(many)"` //设置一对多反向关系
}

//矿山部门列表
type MineDeptInfo struct {
	Id             int64        `orm:"auto"` //主键自增id
	DepartmentName string       //部门名称
	DeptDic        *MineDeptDic `orm:"rel(fk);column(department_class)"` //部门类型，关联到MineDeptDic
	PId            int64        //父级部门id
	Status         int64        //本条记录的有效性，1-有效，0-无效
	//CreateTime     string               `orm:"type(datetime)"` //创建时间
	//CreateUserId   int64                //创建者id
	//UpdateTime     string               `orm:"type(datetime)"` //更新时间
	//UpdateUserId   int64                //更新者
	Seq           int64                //排序
	TreeLevelCode string               //层级编码
	SamplePools   []*LabSampletypePool `orm:"reverse(many)"` //设置一对多反向关系
	GoodsPools    []*GoodsConsumePool  `orm:"reverse(many)"` //设置一对多反向关系
	CheckLines    []*CheckLine         `orm:"reverse(many)"` //设置一对多反向关系
}

//矿山构成列表
type MineConstructionList struct {
	Id          int64  `orm:"auto"`                   //主键自增id
	NameCh      string `orm:"column(level_name_ch)"`  //层级中文名字，在树形结构菜单右键显示的中文名字
	NameEng     string `orm:"column(level_name_eng)"` //层级英文名字
	SubMenulist string //存储的右键菜单，对应的是level_num_code
	NumCode     int64  `orm:"column(level_num_code)"`      //每种分类唯一标识，关联其他表会关联这个字段。
	MainTableId int64  `orm:"column(level_main_table_id)"` //关联tableList.id
	Category    int64  `orm:"column(level_category)"`      //严禁修改，分类每一种表对应一种分类 ，这个数值已经加入项目常量中，严禁修改。
	ExTable     int64  `orm:"column(level_ex_table)"`      //1有固有技术参数 0没有固有技术参数
	TermCode    string //标签简称，注意不能出现下划线。
	LevelNum    int64  //层级结构的第几级。
	Seq         int64  `orm:"column(sortnum_in_level)"` //本层级的排序号
	Remarkes    string //
	Status      int64  //
	//CreateUserId  int64                     //
	//CreateTime    string                    `orm:"type(datetime)"` //
	//UpdataUserId  int64                     //
	//UpdateTime    string                    `orm:"type(datetime)"` //
	MiningClass   int64                     //矿业类别，1-选矿厂，2-采矿场
	EquipmentDics []*OreProcessEquipmentDic `orm:"reverse(many)"` //设置一对多反向关系
}

//矿山设备词典
type OreProcessEquipmentDic struct {
	Id           int64                 `orm:"auto"`                       //主键自增id
	Construction *MineConstructionList `orm:"rel(fk);column(equip_code)"` //constructionList类型的种类表，例如破碎设备种类：旋回破碎等。 关联constructionList.level_num_code
	NameCn       string                `orm:"column(term_name_ch)"`       //种类的中文名称
	NameEn       string                `orm:"column(term_eng)"`           //英文名称
	KeyWord      string                `orm:"column(term_key_word)"`      //关键词
	Code         string                `orm:"column(term_code)"`          //
	CodeShort    string                `orm:"column(term_short)"`         //
	Remarkes     string                //备注
	//CreateUserId   int64                 //
	//CreateTime     string                `orm:"type(datetime)"` //
	//UpdataUserId   int64                 //
	//UpdateTime     string                `orm:"type(datetime)"`
	Status         int64
	SamplePools    []*LabSampletypePool `orm:"reverse(many)"` //设置一对多反向关系
	LabType        []*SamplingManageSub `orm:"reverse(many)"` //设置一对多反向关系
	GoodsPools     []*GoodsConsumePool  `orm:"reverse(many)"` //设置一对多反向关系
	CheckVariables []*CheckVariableSet  `orm:"reverse(many)"` //设置一对多反向关系
}

//质检取样样本池
type LabSampletypePool struct {
	Id                      int64                   `orm:"auto"` //主键自增id
	SampleTypeName          string                  //样本类型名称
	DeptManage              *MineDeptInfo           `orm:"rel(fk)"`                                 //管理部门
	SampleType              *OreProcessEquipmentDic `orm:"rel(fk);column(sample_equipment_dic_id)"` //样本二级分类ID(样本子类型:)
	SampleNameCode          string                  //平台编码
	UserNameCode            string                  //用户编码
	DeptSamplingId          int64                   //采样部门ID
	SampleTimespan          int64                   //样本时间跨度
	SampleTimespanUnit      int64                   //样本时间跨度单位
	SampleSingleAmount      int64                   //单个样本数量
	SampleSingleAmountUnit  string                  //单个样本计量单位
	SampleSingleAmountForm  int64                   //单个样本形态  0.气体，1.液体，2.粉末，3.固态，4.胶体，5.固液两相，6.气液两相
	SampleSingleLabMethod   string                  //样本化验方法
	SampleSingleLabTime     int64                   //样本化验周期
	SampleSingleLabTimeUnit int64                   //样本化验周期单位
	SampleIsDanger          int64                   //样本是否属于危险品
	Remark                  string                  //备注
	Seq                     int64                   //序号
	//CreateUserId            int64                   //创建者ID
	//CreateTime              string                  `orm:"type(datetime)"` //创建时间
	//UpdateUserId            int64                   //更新者ID
	//UpdateTime              string                  `orm:"type(datetime)"` //更新时间
	Status          int64             //记录的有效状态 1：有效；0：无效
	SamplingManages []*SamplingManage `orm:"reverse(many)"` //设置一对多反向关系
}

//质检化验样本单位定义表
type LabAnaIndexUnitDic struct {
	Id           int64  `orm:"column(index_type)"` //主键自增ID
	IndexUnit    string //化验指标单位
	TermEnglish  string //化验指标单位英文名称
	TermsChinese string //化验指标单位中文名称
	Remark       string //备注
	Seq          int64  //排序号
	//CreateUserId int64  //创建人
	//CreateTime   string `orm:"type(datetime)"` //创建时间
	//UpdateUserId int64  //更新人
	//UpdateTime   string `orm:"type(datetime)"` //更新时间
	Status  int64
	LabType []*SamplingManageSub `orm:"reverse(many)"` //设置一对多反向关系
}

//质检样本管理表
type SamplingManage struct {
	Id             int64                 `orm:"auto"` //主键自增ID
	SampleName     string                //样本名称
	SampleType     *LabSampletypePool    `orm:"rel(fk)"`                              //样本类型ID,sample_pool.id
	SampleFunction *SysDictionary        `orm:"rel(fk);column(sample_function)"`      //样本化验属性,标定样=29,流程样=9
	Stage          *OreProcessDWorkstage `orm:"rel(fk);column(sample_from_stage_id)"` //样本来源作业ID
	IsRegular      int                   `orm:"column(is_regtem_sample)"`             //是否常规样,1=常规样,0=临时样
	SamplingSite   string                //取样点位置
	TreeLevelCode  string                //树形结构节点编码
	Remark         string                //备注
	Seq            int64                 //序号
	//CreateUserId   int64                 //创建人
	//CreateTime     string                `orm:"type(datetime)"` //创建时间
	//UpdateUserId   int64                 //更新人
	//UpdateTime     string                `orm:"type(datetime)"` //更新时间
	Status int64                //有效性
	Subs   []*SamplingManageSub `orm:"reverse(many)"` //设置一对多反向关系
	ToLabs []*SamplelistToLab   `orm:"reverse(many)"` //设置一对多反向关系
}

//质检样本管理表子表（样本化验指标表）
type SamplingManageSub struct {
	Id             int64                   `orm:"auto"`                                   //主键自增ID
	Sample         *SamplingManage         `orm:"rel(fk);column(sample_define_id)"`       //sampling_manage.id
	LabIndex       *OreProcessEquipmentDic `orm:"rel(fk);column(index_equipment_dic_id)"` //化验指标二级分类ID（化验指标类型）
	Unit           *LabAnaIndexUnitDic     `orm:"rel(fk);column(index_analysis_unit)"`    //化验单位
	SampleIndexTag string                  //样本指标标签
	TreeLevelCode  string                  //树形结构节点编码
	Remark         string                  //备注
	Seq            int64                   //序号
	//CreateUserId   int64                   //创建人
	//CreateTime     string                  `orm:"type(datetime)"` //创建时间
	//UpdateUserId   int64                   //更新人
	//UpdateTime     string                  `orm:"type(datetime)"` //更新时间
	Status int64 //有效性
}

//取样送样表
type SamplelistToLab struct {
	Id              int64              `orm:"auto"`           //主键自增ID
	Sample          *SamplingManage    `orm:"rel(fk)"`        //样本ID，sampling_manage.id
	SamplingTime    string             `orm:"type(datetime)"` //取样时间
	SampleToLabTime string             `orm:"type(datetime)"` //送样时间
	SampleSitePhoto string             //取样点照片
	SampleCodeNum   string             //样本编号
	SampleQrCode    string             //样本二维码
	ProcessIsMobile int64              //是否手机端作业
	TreeLevelCode   string             //树形结构节点编码
	Remark          string             //备注
	Seq             int64              //序号
	CreateUser      *SysUser           `orm:"rel(fk)"`        //创建人
	CreateTime      string             `orm:"type(datetime)"` //创建时间
	UpdateUserId    int64              //更新人
	UpdateTime      string             `orm:"type(datetime)"` //更新时间
	Status          int64              //有效性
	LabType         int64              //化验类型0:未接收， 1：接收未化验，2：化验未审核， 3：已审核
	ListingNum      string             //清单编号
	Results         []*LabAnaResultTsd `orm:"reverse(many)"` //设置一对多反向关系
}

type LabAnaResultTsd struct {
	Id                  int64              `orm:"auto"`                   //主键自增ID
	SampleToLab         *SamplelistToLab   `orm:"rel(fk)"`                //送样样本id
	LabTag              *SamplingManageSub `orm:"rel(fk);column(tag_id)"` //sampling_manage_sub.id
	SampleIndexTag      string             //样本指标标签
	SampleIndexTagValue string             //样本化验值
	SampleIsOnlineMeas  int64              //是否是有实时测量？
	SampleMeasValue     string             //样本测量值
	Remark              string             //备注
	Seq                 int64              //序号
	Status              int64              //有效性
	LabType             int64              //化验类型 0：接收未化验，1：化验未审核
	CreateUser          *SysUser           `orm:"rel(fk)"`        //创建人
	CreateTime          string             `orm:"type(datetime)"` //创建时间
	UpdateUserId        int64              //更新人
	UpdateTime          string             `orm:"type(datetime)"` //更新时间
}

//单位表
type SysUnit struct {
	Id         int64  //自增主键
	UnitTypeId int64  //单位类型id
	UnitName   string //单位符号
	UnitNameCh string //中文名称
	Remarke    string //说明
	Status     int64  //0 无效 1 有效
	//CreateTime   string `orm:"type(datetime)"` //创建时间
	//CreateUserId int64  //创建者id
	//UpdateTime   string `orm:"type(datetime)"` //更新时间
	//UpdateUserId int64  //更新者id
	CheckTagLists []*CheckTagList `orm:"reverse(many)"` //设置一对多反向关系
}

//巡检变量定义表
type CheckVariableSet struct {
	Id           int64                   `orm:"auto"`                 //主键自增ID
	Type         *OreProcessEquipmentDic `orm:"rel(fk);column(type)"` //检查类型，对应 equipment_dic.id
	CheckContent string                  //检查内容
	VaribleName  string                  //变量名称
	Status       int64                   //本条记录的有效性，1-有效，0-无效
	Seq          int64                   //排序
	Remark       string                  //备注
	//CreateTime   string                  `orm:"type(datetime)"` //创建时间
	//CreateUserId int64                   //创建者
	//UpdateTime   string                  `orm:"type(datetime)"` //更新时间
	//UpdateUserId int64                   //更新者
	CheckTaglists []*CheckTagList `orm:"reverse(many)"` //设置一对多反向关系
}

type CheckTagList struct { //巡检tag定义表
	Id           int64                 `orm:"auto"` //
	CheckName    string                //检查项名称
	TagName      string                //标签名
	Stage        *OreProcessDWorkstage `orm:"rel(fk)"` //作业id
	ResourceType int64                 //资源（表）类型
	ResourceId   int64                 //资源（表）id
	Variable     *CheckVariableSet     `orm:"rel(fk)"`                                 //检查内容，对应 check_variable_set.id
	Unit         *SysUnit              `orm:"rel(fk);null"`                            //单位 id
	SiteEquipRel *CheckSiteEquipRel    `orm:"rel(fk);column(check_site_equip_exe_id)"` //站点设备表ID
	StopStatus   int64                 //设备状态：0-停车，1-不停车
	AlarmType    int64                 //报警类型
	NormalVal    string                //正常值
	MeasMaxvalue float64               //量程上限
	MeasMinvalue float64               //量程下限
	LimitLl      float64               //下下限报警值
	LimitL       float64               //下限报警值
	LimitH       float64               //上限报警值
	LimitHh      float64               //上上限报警值
	Seq          int64                 //排序
	Status       int64                 //本条记录的有效性：1-有效，0-无效
	//ResourceName  string                //实体名称（设备、电机、仪表、分析仪器的名称）(重复,与resource_type\resource_id)
	//CheckType     int64                 //(重复)检查类型，对应 equipment_dic.id(与CheckVariableSet中的type重复)
	//CheckSiteId   int64                 //站点ID(不用了)
	//CheckLineId   int64                 //线路ID(不用了)
	//CheckPeriodId int64                 //周期ID(无用)
	//CreateUser          *SysUser              `orm:"rel(fk)"`        //创建人
	//CreateTime          string                `orm:"type(datetime)"` //创建时间
	//UpdateUser          *SysUser              `orm:"rel(fk)"`        //更新人
	//UpdateTime          string                `orm:"type(datetime)"` //更新时间
	CheckItemExes []*CheckItemExe `orm:"reverse(many)"` //设置一对多反向关系
}

//站点设备关系表
type CheckSiteEquipRel struct {
	Id           int64      `orm:"auto"`                   //自增主键
	CheckSite    *CheckSite `orm:"rel(fk)"`                //巡检站点ID
	ResourceType int64      `orm:"column(level_category)"` //顶级设备类型
	ResourceId   int64      //顶级设备ID，顶级设备包括设备、仪表、分析仪器
	ResourceName string     //顶级设备名称
	Seq          int64      //排序号
	//CreateTime    string `orm:"type(datetime)"` //创建时间
	//CreateUserId  int64  //创建者ID
	CheckTaglists []*CheckTagList `orm:"reverse(many)"` //设置一对多反向关系
	CheckItemExes []*CheckItemExe `orm:"reverse(many)"` //设置一对多反向关系
}

//站点表
type CheckSite struct {
	Id        int64      `orm:"auto"` //自增主键
	Name      string     //站点名称
	NfcCardId int64      //射频卡编码
	CheckLine *CheckLine `orm:"rel(fk)"` //点检路线id
	NoPhotos  int64      //巡检模式 0：直接进入，1：拍照，2：扫描NFC
	Status    int64      //
	Seq       int64      //排序
	Remark    string     //备注
	//CreateTime   string `orm:"type(datetime)"` //创建时间
	//CreateUserId int64  //创建者
	//UpdateTime   string `orm:"type(datetime)"` //更新时间
	//UpdateUserId int64  //更新者

	SiteEquipRels []*CheckSiteEquipRel `orm:"reverse(many)"` //设置一对多反向关系
	CheckSiteExes []*CheckSiteExe      `orm:"reverse(many)"` //设置一对多反向关系
}

//线路表
type CheckLine struct {
	Id       int64         `orm:"auto"` //自增主键
	LineName string        //路线名称
	Dept     *MineDeptInfo `orm:"rel(fk)"` //部门id
	Status   int64         //本条记录的有效性，1-有效，0-无效
	Remark   string        //备注
	//CreateTime   string `orm:"type(datetime)"` //创建时间
	//CreateUserId int64  //创建者id
	//UpdateTime   string `orm:"type(datetime)"` //更新时间
	//UpdateUserId int64  //更新者
	CheckSites []*CheckSite `orm:"reverse(many)"` //设置一对多反向关系
	CheckPlans []*CheckPlan `orm:"reverse(many)"` //设置一对多反向关系
}

//巡检计划表
type CheckPlan struct {
	Id                int64      `orm:"auto"` //自增主键
	PlanName          string     //计划名称
	Line              *CheckLine `orm:"rel(fk)"`                         //点检路线id
	User              *SysUser   `orm:"rel(fk)"`                         //负责人用户id
	Period            *Period    `orm:"rel(fk);column(check_period_id)"` //点检周期id，period.id
	StartTime         string     `orm:"type(datetime)"`                  //开始时间
	EndTime           string     `orm:"type(datetime)"`                  //结束时间
	Delay             int64      //延迟分钟数
	Status            int64      //本条记录的有效性，1-有效，0-无效
	Remark            string     //备注
	DeptTreeLevelCode string     //部门的 tree_level_code
	Seq               int64      //排序
	DeptId            int64      //部门id（重复）
	//CreateTime        string `orm:"type(datetime)"` //创建时间
	//CreateUserId      int64  //创建者id
	//UpdateTime        string `orm:"type(datetime)"` //更新时间
	//UpdateUserId      int64  //更新者
	CheckPlanExes []*CheckPlanExe `orm:"reverse(many)"` //设置一对多反向关系
}

//周期定义表
type Period struct {
	Id           int64  `orm:"auto"` //自增主键
	Name         string //周期名称
	PeriodType   int64  //周期类型：1.巡检周期  2.保养周期;3可视化周期
	TimeUnitType int64  //时间单位类型
	TimeSpan     int64  //时间跨度
	Time         int64  //次数
	EffectTime   string `orm:"type(datetime)"` //生效时间
	CanDelay     int64  //可否延期，1-可以，0-不可以
	SetTemplate  int64  //是否作为模板，1-是，0-不是
	Status       int64  //本条记录的有效性，1-有效，0-无效
	Seq          int64  //排序
	Remark       string //备注
	//CreateTime   string `orm:"type(datetime)"` //创建时间
	//CreateUserId int64  //创建者
	//UpdateTime   string `orm:"type(datetime)"` //更新时间
	//UpdateUserId int64  //更新者
	//LineId       int64  //路线id
	//PlanId       int64  //计划id
	CheckPlans  []*CheckPlan  `orm:"reverse(many)"` //设置一对多反向关系
	PeriodItems []*PeriodItem `orm:"reverse(many)"` //设置一对多反向关系
}

//周期项目定义表
type PeriodItem struct {
	Id        int64   `orm:"auto"`    //自增主键
	Period    *Period `orm:"rel(fk)"` //周期id,period.id
	Name      string  //子项名称
	RefTime1  int64   //基准时间
	StartTime string  //开始时间
	EndTime   string  //结束时间
	EndType   int64   //结束时间类型，0-今日，1-次日
	Day       int64   //几号，如果time_unit_type 为月类型的话，该字段有值
	WeekType  int64   //星期类型，1-星期一，2：星期2 .. 如果time_unit_type 为周类型的话，该字段有值
	Status    int64   //本条记录的有效性，1-有效，0-无效
	Seq       int64   //排序
	Remark    string  //备注
	//CreateTime   string `orm:"type(datetime)"` //创建时间
	//CreateUserId int64  //创建者
	//UpdateTime   string `orm:"type(datetime)"` //更新时间
	//UpdateUserId int64  //更新者
}

//巡检计划执行表
type CheckPlanExe struct {
	Id            int64      `orm:"auto"` //自增主键
	PlanName      string     //计划名称
	Plan          *CheckPlan `orm:"rel(fk)"` //计划id
	ExecuteDate   string     //计划执行日期
	TotalCheckNum int64      //总检查项数
	UnCheckNum    int64      //未检数
	CheckNum      int64      //已检数
	PunctualNum   int64      //准点数
	AheadNum      int64      //提前数
	BehindNum     int64      //晚点数
	PunctualRate  string     //准点率
	Status        int64      //本条记录的有效性，1-有效，0-无效
	Seq           int64      //排序
	CreateTime    string     `orm:"type(datetime)"` //创建时间
	CreateUserId  int64      //创建者
	UpdateTime    string     `orm:"type(datetime)"` //更新时间
	UpdateUserId  int64      //更新者
	Remark        string     //备注
	StartTime     string     //计划执行开始时间
	EndTime       string     //计划执行结束时间
	EndType       int64      //结束时间类型，0-今日，1-次日
	EndTypeName   string     //结束时间类型名称
	AllStartTime  string     `orm:"type(datetime)"` //开始时间
	AllEndTime    string     `orm:"type(datetime)"` //结束时间
	Delay         int64      //延迟分钟数
	SubmitTime    string     `orm:"type(datetime)"` //提交任务时间
	//ExecutorId    int64      //执行人id(user_id)(未用)
	//ExecutorCode  int64      //执行人编号(未用)
	//ExecutorName  string     //执行人姓名(未用)
	//LineName      string     //路线名称
	//LineId        int64      //路线id
	CheckSiteExes []*CheckSiteExe `orm:"reverse(many)"` //设置一对多反向关系
}

//巡检点执行表
type CheckSiteExe struct {
	Id            int64         `orm:"auto"`    //自增主键
	CheckPlanExe  *CheckPlanExe `orm:"rel(fk)"` //点检每日执行数据计划id
	CheckSite     *CheckSite    `orm:"rel(fk)"` //点检站点id
	CheckSiteName string        //站点名称
	CheckNum      int64         //已检数
	UnCheckNum    int64         //未检数
	CardCode      string        //NFC卡编号
	CardName      string        //NFC卡名称
	NoPhotos      int64         //允许不拍照（0：不允许，1：允许）
	Status        int64         //本条记录的有效性，1-有效，0-无效
	Seq           int64         //排序
	CreateTime    string        `orm:"type(datetime)"` //创建时间
	Remark        string        //备注
	//CheckPlanId    int64  	//点检计划id
	//CreateUserId  int64       //创建者
	//UpdateTime    string      `orm:"type(datetime)"` //更新时间
	//UpdateUserId  int64       //更新者
	CheckItemExes []*CheckItemExe `orm:"reverse(many)"` //设置一对多反向关系
}

//巡检执行结果表
type CheckItemExe struct {
	Id           int64         `orm:"auto"`    //自增主键
	CheckSiteExe *CheckSiteExe `orm:"rel(fk)"` //点检执行数据区域id
	//巡检站点设备关系记录表id
	SiteEquipRel *CheckSiteEquipRel `orm:"rel(fk);column(check_site_equip_rel_exe_id)"`

	CheckItem          *CheckTagList `orm:"rel(fk)"` //检查项id,对应check_tag_list.id
	CheckItemName      string        //检查项名称(重复:Taglist中已定义)
	GdTagName          string        //点检规范里的标签(重复:Taglist中已定义)
	CheckStatus        int64         //检查状态，1-已检，0-未检
	StartTime          string        //开始时间
	EndTime            string        //结束时间
	EndType            int64         //结束时间类型，1-今日，2-次日
	EndTypeName        string        //结束时间类型
	CheckTime          string        `orm:"type(datetime)"` //实际检查时间
	ExecutorName       string        //执行人名称
	ActualExecutorName string        //实际执行者名称
	AlarmType          int64         //报警类型，1-超下下限，2-超下限，3-正常，4-超上限，5-超上上限，6-量程外
	AlarmTypeName      string        //报警类型名称
	ResourceType       int64         //资源（表）类型(重复:Taglist中已定义)
	ResourceId         int64         //资源（表）id(重复:Taglist中已定义)
	ResourceName       string        //实体名称（设备、电机、仪表、分析仪器的名称）(重复:Taglist中已定义)
	LevelsName         string        //仪表、机械构成、部件名称的拼接(重复:Taglist中已定义)
	CheckType          int64         //检查类型，对应 equipment_dic.id(重复:Taglist中已定义)
	CheckTypeName      string        //检查类型名称
	VariableId         int64         //检查内容类型，对应 check_variable_set.id
	VariableName       string        //检查内容类型名称
	StopStatus         int64         //设备状态：0-停车，1-不停车(重复:Taglist中已定义)
	StopStatusName     string        //设备状态名称
	UnitId             int64         //单位 id(重复:Taglist中已定义)
	UnitName           string        //单位名称
	NormalVal          string        //标准值(重复:Taglist中已定义)
	MeasMaxvalue       string        //量程上限(重复:Taglist中已定义)
	MeasMinvalue       string        //量程下限(重复:Taglist中已定义)
	LimitLl            string        //下下限报警值(重复:Taglist中已定义)
	LimitL             string        //下限报警值(重复:Taglist中已定义)
	LimitH             string        //上限报警值(重复:Taglist中已定义)
	LimitHh            string        //上上限报警值(重复:Taglist中已定义)
	RefValue           string        //参考值
	RefValueType       int64         //参考值类型
	CheckResult        string        //检查结果
	InputType          string        //输入类型
	ValueType          string        //值类型
	Abnormal           int64         //是否异常，0-异常，1-正常
	Status             int64         //本条记录的有效性，1-有效，0-无效
	Seq                int64         //排序
	TourRemark         string        //巡视异常备注
	TourAbnormal       int64         //巡视异常（1.该设备处于备用状态  2.该设备已不存在 3.没有可测量工具 4.其他）
	Remark             string        //备注
	CreateTime         string        `orm:"type(datetime)"` //创建时间
	UpdateTime         string        `orm:"type(datetime)"` //更新时间
	EndTime2           string        `orm:"type(datetime)"` //结束时间
	//ExecutorId         int64              //执行人id
	//ActualExecotorId   int64              //实际执行者id
	//CreateUserId       int64  //创建者
	//UpdateUserId       int64  //更新者
	//CheckPlanId            int64  //点检计划id
	//CheckPlanExeId         int64  //点检执行数据计划id
}

//系统日志
type SysLog struct {
	Id          int64  `orm:"auto"` //主键自增ID
	Description string //描述
	SysType     int64  //系统类型 1：PC端 2：移动端 3:计算服务
	OprType     int64  //操作类型 0:其他,1:"添加",2:"删除",3:"更新",4:"查看",5:添加/更新",6:"登录"
	MethodName  string //请求方法名
	ClassName   string //请求类名
	RemoteIp    string //请求IP
	ReqUrl      string //URI
	ReqMethod   string //请求方式
	ReqParams   string //提交参数
	Exception   int64  //异常
	StartTime   string `orm:"type(datetime)"` //开始时间
	//EndTime     string `orm:"type(datetime)"` //结束时间(未用)
	User *SysUser `orm:"rel(fk)"` //用户ID
}

//照片存储表
type MineCheckImg struct {
	Id           int64  `orm:"auto"` //主键自增ID
	TableNameId  int64  //关联表名ID。如tablelist.id
	TableId      int64  //点检结果、站点、维修故障保修等表的id
	TrueFileName string //文件真实名字
	FileName     string //文件保存后的名字
	FileUrl      string //自动生成的唯一文件名
	Status       int64  //记录本条记录的有效性，1-有效，0-无效
	SuffixName   string //文件后缀名
	SaveDate     string `orm:"type(datetime)"` //保存日期
	FileSite     int64  //关联位置,如：技术文件，就存储为1
	LevelNumCode int64  //关联各个层级表的construction_code,如mine_basic_info.construction_code。
	Remark       string //备注
}

//心跳信息存储表(SQLite3)
type HeartBeat struct {
	Id           int64  `orm:"auto"`           //主键自增ID
	DataTime     string `orm:"type(datetime)"` //时间
	RunMinutes   int64  //本次运行时间
	TotalMinutes int64  //总运行时间
	DogChecked   int64  //看门狗检查标记,默认0,看门狗检查完毕后置1
}

type KpiArithmetic struct {
	Id                      int64  `orm:"auto"` //主键自增ID
	ArithmeticName          string //算法工具包名称
	ArithmeticName2         string //算法工具包名称
	ArithmeticResultType    string //计算结果数据类型
	ArithmeticResultEcharts string //计算结果展示类型
	ArithmeticUrl           string //所需页面编号
	ArithmeticType          string //计算显示名称（例如：求和）
}

type KpiArithmeticResult struct {
	Id                         int64  `orm:"auto"` //主键自增ID
	ResultName                 string //结果名称(例如：求半自磨机处理量)
	ArithmeticObjectId         string //计算主体变量ID
	ArithmeticObjectType       string //计算主体变量描述
	ArithmeticObjectRemark     string //计算主体变量类型
	ArithmeticAuxiliaryId1     string //计算辅助变量1ID
	ArithmeticAuxiliaryRemark1 string //计算辅助变量1描述
	ArithmeticAuxiliaryId2     string //计算辅助变量2ID
	ArithmeticAuxiliaryRemark2 string //计算辅助变量2描述
	BeginTime                  string `orm:"type(datetime)"` //统计数据开始时间
	EndTime                    string `orm:"type(datetime)"` //统计数据结束时间
	FinalResult                string //结果数据
	FinalResultType            string //结果数据类型
	ArithmeticId               string //算法工具ID
	FinalResultEacharts        string //结果数据展示类型
	CreTime                    string `orm:"type(datetime)"` //创建时间
	CrePerson                  string //创建人
}
