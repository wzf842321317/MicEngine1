package goldenweb

import (
	"github.com/bkzy-wangjp/MicEngine/MicScript/statistic"
)

const (
	_SnapCmd     = "SnapShot"            //获取快照
	_HisCmd      = "History"             //获取历史时刻值或者历史范围值
	_HisInterCmd = "HistoryInterval"     //获取插补的等间隔历史数据
	_Table       = "Table"               //获取变量表信息
	_Sever       = "Server"              //获取服务器时间
	_Point       = "Point"               //获取标签点信息
	_HisSum      = "HistorySummary"      //历史统计数据
	_SnapWrite   = "pointDataCollection" //写快照数据
)

type GoldenAdd struct { //庚顿数据库地址数据结构
	Host string //主机地址,ip地址或者域名
	Port int    //端口号
}

type SnapshotCmd struct { //读取快照时的参数结构
	cmd     string   //指令,"Snapshot"
	TagName []string //变量名
}
type SnapWriteCmd struct { //读取快照时的参数结构
	cmd   string      //指令,"pointDataCollection"
	Datas []SnapWrite //需要写的变量列表
}
type TableListCmd struct { //读取数据库中数据表列表时的参数结构(只需要建立变量,不需要设定参数)
	cmd      string //指令,"Table"
	TableMsg string //isAllTalbleInfo为all,获取全部表名,为表名的时候，获取指定表的信息
}

type TagListCmd struct { //读取某表变量标签列表时或者某变量信息的参数结构
	cmd       string   //指令,"Point"
	TableName string   //表名称，与变量名二选一,不选的值为空
	TagName   []string //变量名,与表名称二选一，不选的值为空
}

type HistoryCmd struct { //读取历史数据和历史时刻点数据时的参数结构
	cmd       string   //指令,"History"
	TagName   []string //变量名
	TimePoint string   //获取时刻值的时候用，同时有beginTime和endTime时失效
	BeginTime string   //开始时间，获取时间范围值的时候用
	EndTime   string   //结束时间，获取时间范围值的时候用
}

type HisIntervalCmd struct { //读取等间隔历史数据时的参数结构
	cmd       string   //指令,"History"
	TagName   []string //变量名
	BeginTime string   //开始时间，获取时间范围值的时候用
	EndTime   string   //结束时间，获取时间范围值的时候用
	Interval  int64    //间隔单位,秒
	count     int64    //数据点的个数。要获取全部点的时候,count取值要大于(endTime-beginTime)/interval
}
type HisSumCmd struct { //读取历史数据和历史时刻点数据时的参数结构
	cmd       string   //指令,"HistorySummary"
	TagName   []string //变量名
	BeginTime string   //开始时间，获取时间范围值的时候用
	EndTime   string   //结束时间，获取时间范围值的时候用
	DataType  string   //数据类型."min","max","range","total","sum","mean","poweravg","diff","duration","pointcnt","risingcnt","fallingcnt","advangce", "sd", "stddev", "se", "ske", "kur", "mode", "median","distribution","increment"
	Group     int      //高级统计(求众数)时的分组数,如果为0,则默认分为100组
	Interval  int64    //间隔单位,秒。如果为负数，则读取庚顿原生统计数据;如果为0，读取原始历史数据的统计数据,如果大于0，则读取等间隔历史数据的统计数据
}

type SnapData struct { //快照基础数据结构
	id       int    `json:"ID"`       //测点ID
	TagName  string `json:"TagName"`  //测点全名称
	Time     string `json:"Time"`     //时间
	Value    string `json:"Value"`    //数值
	Quality  int    `json:"Quality"`  //质量码(GOOD = 0,NODATA = 1,CREATED = 2,SHUTDOWN = 3,CALCOFF = 4,BAD = 5,DIVBYZERO = 6,REMOVED = 7,OPC = 256,USER = 512)
	error    int    `json:"Error"`    //错误码
	errorMsg string `json:"ErrorMsg"` //错误码描述
}

type SnapWrite struct {
	TagName string `json:"TagName"` //测点全名称
	Time    string `json:"Time"`    //时间
	Value   string `json:"Value"`   //数值
	Quality int    `json:"Quality"` //质量码(GOOD = 0,NODATA = 1,CREATED = 2,SHUTDOWN = 3,CALCOFF = 4,BAD = 5,DIVBYZERO = 6,REMOVED = 7,OPC = 256,USER = 512)
}
type SnapWriteResponse struct { //写快照响应时间
	SuccessCount int        `json:"SuccessCount"` //写成功的数量
	ErrorList    []SnapData //失败的数据
}
type HisData struct { //历史数据基础数据结构
	ms    string `json:"Ms"`    //时间戳(弃用，使用时间值)
	Time  string `json:"Time"`  //时间
	Value string `json:"Value"` //数值
}
type TableInfo struct { //数据表信息
	Id   int64  `json:"Id"`   //表id
	Name string `json:"Name"` //表名
	Desc string `json:"Desc"` //表描述
}
type HisSumData struct { //庚顿原生历史数据统计数据结构
	Min        float64 `json:"Min"`        //最小值
	Max        float64 `json:"Max"`        //最大值
	Total      float64 `json:"Total"`      //表示统计时间段内的累计值，结果的单位为标签点的工程单位(计算的是面积)
	Sum        float64 `json:"CalcTotal"`  //统计时间段内的算术累积值
	Mean       float64 `json:"CalcAvg"`    //统计时间段内的算术平均值
	PowerAvg   float64 `json:"PowerAvg"`   //统计时间段内的加权平均值
	Diff       float64 `json:"Difference"` //统计时间段内的差值
	BeginTime  string  `json:"StartTime"`  //开始时间
	EndTime    string  `json:"EndTime"`    //结束时间
	eqLastYear float64 `json:"EqLastYear"` //同比去年总量 百分比
	availble   float64 `json:"Availble"`   //可利用率
}

type HisSumDataExi statistic.StatisticData

type basePointInfo struct {
	Tag            string  //Tag 标签点名
	Id             int     //Id 标签点ID
	DataType       string  //DataType 数据类型
	TableId        int     //TableId 标签点所属表ID
	Desc           string  //Desc 有关标签点的描述性文字
	Unit           string  //Unit 工程单位
	Archive        int     //Archive 是否存档
	Digits         int     //Digits 数值
	ShutDown       int     //ShutDown 停机状态字（Shutdown）
	LowLimit       float64 //LowLimit 量程下限
	HighLimit      float64 //HighLimit 量程上限
	Step           int     //Step 是否阶跃
	Typical        float64 //Typical 典型值
	Compress       int     //Compress 是否压缩
	Compdev        float64 //Compdev 压缩偏差
	CompdevPercent float64 //CompdevPercent 压缩偏差百分比
	CompTimeMax    int     //CompTimeMax 最大压缩间隔
	CompTimeMin    int     //CompTimeMin 最短压缩间隔
	Excdev         float64 //Excdev 例外偏差
	ExcdevPercent  float64 //ExcdevPercent 例外偏差百分比
	ExcTimeMax     int     //ExcTimeMax 最大例外间隔
	ExcTimeMin     int     //ExcTimeMin 最短例外间隔
	Classof        int     //Classof 标签点类别 （Classof=0表示基本点；Classof=1表示采集点；Classof=2表示计算点；Classof=3表示采集计算点）
	ChangeDate     int     //ChangeDate 标签点属性最后一次被修改的时间
	Changer        string  //Changer 标签点属性最后一次被修改的用户名GOLDEN_USER_SIZE
	CreateDate     int     //CreateDate 标签点被创建的时间
	Creator        string  //Creator 标签点创建者的用户名 GOLDEN_USER_SIZE
	Mirror         int     //Mirror 允许镜像 0：OFF,1：ON
	MicroSecond    int     //MicroSecond 时间戳精度 0:表示秒，1表示毫秒
	ScanIndex      int     //ScanIndex 采集点扩展属性集存储地址索引
	CalcIndex      int     //CalcIndex 计算点扩展属性集存储地址索引
	AlarmIndex     int     //AlarmIndex 报警点扩展属性集存储地址索引
	TableDotTag    string  //TableDotTag 标签点全名，格式为“表名称.标签点名称” GOLDEN_TAG_SIZE, GOLDEN_TAG_SIZE
	Padding        string  //Padding 基本标签点备用字节
}
type scanPointInfo struct {
	Id         int    //Id 全库唯一标识。0表示无效
	Source     string //Source 数据源
	Scan       int    //Scan 是否采集。
	Instrument string //Instrument 设备标签
	Locations  string //Locations 共包含五个设备位址，缺省值全部为0。GOLDEN_LOCATIONS_SIZE
	UserInts   string //UserInts 共包含两个自定义整数，缺省值全部为0。GOLDEN_USERINT_SIZE
	UserReals  string //UserReals 共包含两个自定义单精度浮点数，缺省值全部为0。GOLDEN_USERREAL_SIZE
	Padding    string //Padding 采集标签点备用字节。GOLDEN_PACK_OF_SCAN
}
type calcPointInfo struct {
	Id       int    //Id 标签点ID
	Equation string //Equation 实时方程式
	Trigger  int    //Trigger 计算触发机制。（Trigger=1表示事件触发，Trigger=2表示周期触发，Trigger=0表示无触发,Trigger=3表示定时触发，定时触发的设置Period定时周期）
	Timecopy int    //Timecopy 算结果时间戳参考
	Period   int    //Period 对于“定时触发”的计算点，设定其计算周期，单位：秒
}
type PointInfoStruct struct { //标签点属性
	BasePointInfo basePointInfo `json:"BasePointInfo"`
	ScanPointInfo scanPointInfo `json:"ScanPointInfo"`
	CalcPointInfo calcPointInfo `json:"ClacPointInfo"`
}
type HisDataMap map[string][]HisData           //读取历史数据和等间隔历史数据时返回的数据结构,string是标签点名称
type HisPointDataMap map[string]HisData        //读取某时间点历史数据时返回的数据的结构,string是标签点名称
type SnapDataMap []SnapData                    //读取快照时返回的数据结构
type HisSumDataMap map[string]*HisSumData      //读取历史统计数据的数据结构(庚顿原生),string是标签点名称
type HisSumDataMapExi map[string]HisSumDataExi //读取历史统计数据的数据结构(增强型),string是标签点名称
type PointInfo []PointInfoStruct
