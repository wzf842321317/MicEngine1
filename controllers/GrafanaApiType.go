package controllers

const (
	_DelaySecWhenNow = -15 //当选择的最后时间为当前时间时,往前推迟的时间（防止时间报警）
)

type getCmd struct {
	Order     string   //指令
	TagNames  []string //变量名
	TableName string   //表名。用于查询表中的变量列表
	BeginTime string   //起始时间
	EndTime   string   //结束时间
	TimePoint string   //时间点。如果时间点不为空，优先读取时间点数据，起始时间、结束时间无效
	Interval  int64    //历史数据间隔时间,单位秒。为0时，读取原始历史数据或者时间点数据,不为零时,读取等间隔历史值.
}

type gfTableColumn struct { //表格形式的数据反馈表的列定义
	Text string `json:"text"` //列名
	Type string `json:"type"` //列类型,可选:time,string,number
}

type gfTableTypeData struct { //表格形式的数据反馈结构
	Columns [3]gfTableColumn `json:"columns"`
	Rows    [][3]interface{} `json:"rows"`
	Type    string           `json:"type"` //固定返回值"table"
}

type gfTimeSeriesTypeData struct { //时间序列形式的数据反馈结构
	Target     string           `json:"target"`
	Datapoints [][2]interface{} `json:"datapoints"`
}
type goodsPointData struct { //物耗时间点数据结构
	Time  string
	Value float64
}
type goodsSumData struct { //物耗累积数据结构
	BeginTime string
	EndTime   string
	Sum       float64
}
type gfScoped struct {
	Text  string
	Value string
}
type gfScoped2 struct {
	Text  string
	Value int64
}
type gfScopedVar struct {
	From       gfScoped  `json:"__from"`
	To         gfScoped  `json:"__to"`
	Interval   gfScoped  `json:"__interval"`
	IntervalMs gfScoped2 `json:"__interval_ms"`
}
type gftimeRangeRaw struct {
	From string
	To   string
}
type gftimeRange struct {
	From string
	To   string
	Raw  gftimeRangeRaw
}

type gfReadDataCmd struct {
	Order        string   //"快照" 或者 "Snapshot","历史" 或者 "HistoryInterval","统计" 或者 "HistorySummary"
	TagNames     []string `json:"tagNames"`  //变量名数组,不需要带表名
	TagName      string   `json:"tagName"`   //变量名,不需要带表名
	TimePoint    string   `json:"timePoint"` //仅在Order为"历史"\"HistoryInterval"且Interval为0的时候有效,用于读取历史时刻值
	Interval     int64    //间隔,仅对“历史”命令有效。可以为负值、0和正值.为负值时,根据所选的时间范围自适应数据点的间隔,为0时返回的是压缩后没有进行时间对齐的数据,正值时是相邻两个数据点之间的对齐间隔,时间为秒.如果没有设置，则取0
	DecimalPoint int64    `json:"decimalPoint"` //定义小数点后的数据位数,不填写或者填写为0表示不限小数位数(废弃)
	TagDescs     []string `json:"tagDescs"`     //变量描述数组.可选项，如果填写则以填写的值优先,不填写则从庚顿数据库中的标签描述处取,如果标签描述没填,则用变量名代替
	TagDesc      string   `json:"tagDesc"`      //变量描述.可选项，如果填写则以填写的值优先,不填写则从庚顿数据库中的标签描述处取,如果标签描述没填,则用变量名代替
}

type gfDataRequest struct {
	Data   gfReadDataCmd
	Target string //数据表名
	RefId  string `json:"refId"`
	Hide   bool
	Type   string //timeseries或者table
}

type grafanaRequestDataStruct struct {
	RequestId     string `json:"requestId"`
	Timezone      string
	PanelId       int64 `json:"panelId"`
	DashboardId   int64 `json:"dashboardId"`
	Range         gftimeRange
	Interval      string
	IntervalMs    int64 `json:"intervalMs"`
	Targets       []gfDataRequest
	MaxDataPoints int64 `json:"maxDataPoints"`
	ScopedVars    gfScopedVar
	CacheTimeout  string         `json:"cacheTimeout"`
	StartTime     int64          `json:"startTime"`
	EndTime       int64          `json:"endTime"`
	RangeRaw      gftimeRangeRaw `json:"rangeRaw"`
	AdhocFilters  []string       `json:"adhocFilters"`
}

type grafanaRequestStruct struct {
	Method string
	Url    string
	Data   grafanaRequestDataStruct
}

type grafanaRequest struct {
	XhrStatus string `json:"xhrStatus"`
	Request   grafanaRequestStruct
}
