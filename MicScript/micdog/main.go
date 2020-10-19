package main

import (
	"database/sql"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/bkzy-wangjp/MicEngine/models"

	iniconf "github.com/clod-moon/goconf"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	fmt.Println("监测计算服务启动...")
	run_state, _, _, doc_timer, _ := readConf()
	if run_state == 0 {
		for {
			time.Sleep(1 * time.Second)
			startRun()
			time.Sleep(time.Duration(doc_timer) * time.Second)
		}
	} else if run_state == 1 {
		startRun()
	}
}

/****************************************************
功能:获取所在路径
输入:
输出:
说明:
时间:2020年9月9日
编辑:wang_zf .
****************************************************/
func GetCurrentPath() string {
	dir, err := os.Getwd()
	if err != nil {
		err.Error()
	}
	return strings.Replace(dir, "\\", "/", -1)
}

/****************************************************
功能:执行看门狗监测计算服务，超时重启
输入:
输出:
说明:
时间:2020年9月9日
编辑:wang_zf .
****************************************************/
func startRun() {
	_, doc_state, doc_timeout, doc_timer, _ := readConf()
	//定时器
	if doc_state == 1 {
		if int(timeDifference()) > doc_timeout {
			reBootServer()
			time.Sleep(time.Duration(doc_timer) * time.Second)
			fmt.Println("服务未启动,正在启动...", int(timeDifference()), "秒")

		} else {
			fmt.Println("服务已启动,正在运行...", int(timeDifference()), "秒")

		}
	}
}

/****************************************************
功能:读取配置文件
输入:
输出:执行状态，运行状态，最大超时时间，定时执行时间
说明:
时间:2020年9月9日
编辑:wang_zf .
****************************************************/
func readConf() (int, int, int, int, string) {
	conf := iniconf.InitConfig(GetCurrentPath() + "/conf/micdog.ini")
	run_state, err := strconv.Atoi(conf.GetValue("plat", "run_state"))
	if err != nil {
		fmt.Println("配置文件:run_state格式错误!!!")
	}
	doc_state, err := strconv.Atoi(conf.GetValue("plat", "doc_state"))
	if err != nil {
		fmt.Println("配置文件:doc_state格式错误!!!")
	}
	doc_timeout, err := strconv.Atoi(conf.GetValue("plat", "doc_timeout"))
	if err != nil {
		fmt.Println("配置文件:doc_timeout格式错误!!!")
	}
	doc_timer, err := strconv.Atoi(conf.GetValue("plat", "doc_timer"))
	if err != nil {
		fmt.Println("配置文件:doc_timer格式错误!!!")
	}
	path := conf.GetValue("plat", "path")
	return run_state, doc_state, doc_timeout, doc_timer, path
}

/****************************************************
功能:看门狗心跳写入
输入:sqliter路径
输出:最大时间
说明:
时间:2020年9月2日
编辑:wang_zf .
****************************************************/
func UpdateFlag() {
	db, err := sql.Open("sqlite3", GetCurrentPath()+"/data/micengine.db")
	if err != nil {
		fmt.Println("/data/micengine.db 文件不存在...")
	}
	//写入标记位
	startTime := models.EngineCfgMsg.StartTime
	timeLayout := "2006-01-02 15:04:05"                            //转化所需模板
	loc, _ := time.LoadLocation("Local")                           //重要：获取时区
	theTime, _ := time.ParseInLocation(timeLayout, startTime, loc) //使用模板在对应时区转化为time.time类型
	sr := theTime.Unix()                                           //转化为时间戳 类型是int64

	timeUnix := time.Now().Unix()
	timeDiff := (timeUnix - sr) / 60
	if theTime == 0 {
		timeDiff = 1
	}
	stmt, err := db.Prepare("update main.heart_beat set dog_checked = ?  where id = (select max(id) from main.heart_beat);commit ")
	checkErr(err)
	stmt.Exec(timeDiff)
	checkErr(err)
	stmt.Close()
}

/****************************************************
功能:获取数据库heart_beat最大的时间
输入:sqliter路径
输出:最大时间
说明:
时间:2020年9月2日
编辑:wang_zf .
****************************************************/
func timeDifference() int {
	db, err := sql.Open("sqlite3", GetCurrentPath()+"/data/micengine.db")
	if err != nil {
		fmt.Println("/data/micengine.db 文件不存在...")
	}
	//upstr := "update main.heart_beat set dog_checked = 1  where id = (select max(id) from main.heart_beat)"

	//查询数据库
	rows, err := db.Query("SELECT data_time FROM main.heart_beat ORDER BY ID DESC LIMIT 0,1  ")
	checkErr(err)
	rows.Next()
	var data_time string
	err = rows.Scan(&data_time)
	checkErr(err)
	max_date := strings.ReplaceAll(data_time, "T", " ")
	max_date = strings.ReplaceAll(max_date, "Z", "")
	//待转化为时间戳的字符串 注意 这里的小时和分钟还要秒必须写 因为是跟着模板走的 修改模板的话也可以不写
	toBeCharge := max_date
	//转化所需模板
	timeLayout := "2006-01-02 15:04:05"
	//重要：获取时区
	loc, _ := time.LoadLocation("Local")
	//使用模板在对应时区转化为time.time类型
	theTime, _ := time.ParseInLocation(timeLayout, toBeCharge, loc)
	sr := theTime.Unix()
	//转化为时间戳 类型是int64
	unix := time.Now().Unix()
	max_timetamp := int(unix - sr)
	rows.Close()
	UpdateFlag()
	return max_timetamp

}

/****************************************************
功能:连接sqliter错误日志
输入:错误信息
输出:错误
说明:
时间:2020年9月2日
编辑:wang_zf
****************************************************/
func checkErr(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

/****************************************************
功能:重启服务
输入:路径
输出:
说明:重启服务
时间:2020年9月2日
编辑:wang_zf
****************************************************/
func reBootServer() {
	_, _, _, _, path := readConf()
	c := exec.Command(path)
	if err := c.Run(); err != nil {
		fmt.Println("path not fond: ", err)
	}
}
