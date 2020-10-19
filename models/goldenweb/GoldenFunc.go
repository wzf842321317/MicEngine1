package goldenweb

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego/logs"
	"github.com/bkzy-wangjp/MicEngine/MicScript/statistic"
)

/*
功能：根据庚顿数据库的地址和读取指令获取数据
输入：add:Golden API的地址
	itf:读取指令接口
输出：数据接口
时间：2019年11月28日
编辑：wang_jp
*/
func GetGoldenData(add GoldenAdd, itf interface{}) (interface{}, error) {
	switch itf.(type) {
	case SnapshotCmd:
		cmd, _ := itf.(SnapshotCmd)
		return getSnapShotData(add, cmd)
	case HistoryCmd:
		cmd, _ := itf.(HistoryCmd)
		if cmd.TimePoint == "" { //如果时间点为空，则读取一段范围内的历史值
			return getHistoryData(add, cmd)
		} else { //否则读取指定时间点上的历史值(如果该时间点没有值，则取该点之前的第一个有效值)
			return getHisPointData(add, cmd)
		}
	case HisIntervalCmd: //读取等间隔历史值
		cmd, _ := itf.(HisIntervalCmd)
		return getHisIntervalData(add, cmd)
	case TableListCmd: //读取表格列表和变量列表
		return getTableList(add, itf)
	case TagListCmd: //读取变量列表或者指定变量的信息
		cmd, _ := itf.(TagListCmd)
		if len(cmd.TagName) > 0 { //读取指定变量名的基本信息
			return getTagPointInfo(add, itf)
		} else { //读取指定数据表中的变量列表
			return getTagList(add, itf)
		}
	case HisSumCmd: //读取历史统计
		cmd, _ := itf.(HisSumCmd)
		if cmd.Interval < 0 {
			return getHistorySummaryData(add, cmd) //返回庚顿原生统计数据
		} else {
			return getHistorySummaryIndicator(add, cmd) //返回智云科技计算的统计数据
		}
	default: //没有匹配的类型情况下,返回服务器时间
		return getServerTime(add)
	}
}

/*
功能：根据庚顿数据库的地址和读取指令获取历史时刻数据
输入：add:Golden API的地址
	cmd:读取指令接口
输出：数据接口
时间：2019年11月28日
编辑：wang_jp
*/
func getHisPointData(add GoldenAdd, cmd HistoryCmd) (HisPointDataMap, error) {
	var data HisPointDataMap
	urlstr, err := getUrl(add, cmd) //拼接URL
	if err != nil {
		return data, err
	}

	resp, err := http.Get(urlstr) //Get数据
	if err != nil {
		return data, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body) //接收反馈数据
	if err != nil {
		return data, err
	}

	json.Unmarshal([]byte(body), &data)

	return data, nil
}

/*
功能：根据数据库的地址和读取指令获取等间隔历史数据
输入：add:Golden API的地址
	cmd:读取指令接口
输出：数据接口
时间：2019年11月28日
编辑：wang_jp
*/
func getHisIntervalData(add GoldenAdd, cmd HisIntervalCmd) (HisDataMap, error) {
	var data HisDataMap
	urlstr, err := getUrl(add, cmd) //拼接URL
	if err != nil {
		return data, err
	}
	resp, err := http.Get(urlstr) //Get数据
	if err != nil {
		return data, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body) //接收反馈数据
	if err != nil {
		return data, err
	}

	json.Unmarshal([]byte(body), &data)

	return data, nil
}

/*
功能：根据庚顿数据库的地址和读取指令获取历史数据
输入：add:Golden API的地址
	cmd:读取指令接口
输出：数据接口
时间：2019年11月28日
编辑：wang_jp
*/
func getHistoryData(add GoldenAdd, cmd HistoryCmd) (HisDataMap, error) {
	var data HisDataMap
	urlstr, err := getUrl(add, cmd) //拼接URL
	if err != nil {
		return data, err
	}
	resp, err := http.Get(urlstr) //Get数据
	if err != nil {
		return data, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body) //接收反馈数据
	if err != nil {
		return data, err
	}

	json.Unmarshal([]byte(body), &data)

	return data, nil
}

/*
功能：根据庚顿数据库的地址和读取指令获取庚顿原生历史统计数据
输入：add:Golden API的地址
	cmd:读取指令接口
输出：数据接口
时间：2019年11月28日
编辑：wang_jp
*/
func getHistorySummaryData(add GoldenAdd, cmd HisSumCmd) (HisSumDataMap, error) {
	var data HisSumDataMap
	urlstr, err := getUrl(add, cmd) //拼接URL
	if err != nil {
		return data, err
	}
	resp, err := http.Get(urlstr) //Get数据
	if err != nil {
		return data, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body) //接收反馈数据
	if err != nil {
		return data, err
	}

	json.Unmarshal([]byte(body), &data)

	return data, nil
}

/*
功能：根据庚顿数据库的地址和读取指令获取历史统计数据
输入：add:Golden API的地址
	cmd:读取指令接口
输出：数据接口
时间：2019年12月7日
编辑：wang_jp
*/
func getHistorySummaryIndicator(add GoldenAdd, cmd HisSumCmd) (interface{}, error) {
	his, err := getHistorySummaryDataExi(add, cmd) //获取统计数据
	if err != nil {
		return "", err
	}
	//定义可能的返回map

	switch strings.ToLower(cmd.DataType) { //查询并返回结果
	case "distribution", "datagroup":
		fdis := make(map[string]interface{}, len(cmd.TagName))
		for k, v := range his {
			fdis[k] = v.DataGroup
		}
		return fdis, nil
	case "increment":
		finc := make(map[string]interface{}, len(cmd.TagName))
		for k, v := range his {
			finc[k] = v.Increment
		}
		return finc, nil
	case "advangce", "统计":
		return his, nil
	default:
		fv := make(map[string]interface{}, len(cmd.TagName))
		for k, v := range his {
			fv[k] = SelectValueFromHisSummryExi(v, cmd.DataType)
		}
		return fv, nil
	}
}

/*
功能：根据需要的数据指令从历史统计结构体中返回相应的数据
输入：历史统计结构体,数据类型指令
输出：浮点数
时间：2019年12月23日
编辑：wang_jp
*/
func SelectValueFromHisSummryExi(v HisSumDataExi, dataType string) float64 {
	return statistic.SelectValueFromStatisticData(statistic.StatisticData(v), dataType)
}

/*
功能：根据庚顿数据库的地址和读取指令获取补全了开头和结尾时刻数据的历史数据
输入：add:Golden API的地址
	cmd:读取指令接口
输出：数据接口
时间：2020年1月3日
编辑：wang_jp
*/
func GetHistoryDataAlignHeadAndTail(add GoldenAdd, cmd HisSumCmd) (HisDataMap, error) {
	var his, his_st, his_ed HistoryCmd
	var hisinter HisIntervalCmd
	var snap SnapshotCmd
	hismap := make(HisDataMap, len(cmd.TagName))
	hisd := make(HisDataMap, len(cmd.TagName))

	begintime, err := timeFormat(cmd.BeginTime)
	if err != nil {
		return hismap, err
	}
	endtime, err := timeFormat(cmd.EndTime)
	if err != nil {
		return hismap, err
	}
	//快照
	snap.TagName = cmd.TagName
	//开始时间点
	his_st.TagName = cmd.TagName
	his_st.TimePoint = cmd.BeginTime
	//结束时间点
	his_ed.TagName = cmd.TagName
	his_ed.TimePoint = cmd.EndTime
	//历史区间
	his.TagName = cmd.TagName
	his.BeginTime = cmd.BeginTime
	his.EndTime = cmd.EndTime
	//等间隔历史
	hisinter.TagName = cmd.TagName
	hisinter.BeginTime = cmd.BeginTime
	hisinter.EndTime = cmd.EndTime
	hisinter.Interval = cmd.Interval

	snapmp, err := getSnapShotData(add, snap) //获取快照数据
	if err != nil {
		return hismap, fmt.Errorf("Error when getSnapShorData. Error message is [%s]", err.Error())
	}
	if len(snapmp) < len(cmd.TagName) {
		return hismap, fmt.Errorf("The tag %v no snapshot data", cmd.TagName)
	}
	st, err := getHisPointData(add, his_st) //读取开始时间点数据
	if err != nil {
		return hismap, fmt.Errorf("Error when getHisPointData. Error message is [%s]", err.Error())
	}
	if cmd.Interval == 0 {
		hisd, err = getHistoryData(add, his) //读取区间历史数据
		if err != nil {
			return hismap, fmt.Errorf("Error when getHistoryData. Error message is [%s]", err.Error())
		}
	} else {
		hisd, err = getHisIntervalData(add, hisinter) //读取等间隔区间历史数据
		if err != nil {
			return hismap, fmt.Errorf("Error when getHisIntervalData. Error message is [%s]", err.Error())
		}
	}

	for i, tagname := range cmd.TagName {
		snapTime, err := timeFormat(snapmp[i].Time) //快照数据时间转换
		if err != nil {
			//logs.Alert("Time format error:%s", err.Error())
			return hismap, err
		}
		if snapTime.Before(endtime) {
			return hismap, fmt.Errorf("Tag %s snapshot stoped,the last time is:%s", tagname, snapmp[i].Time)
		}
		var hd []HisData
		if len(st) > 0 { //起始时间点之前有数据
			if IsExistItem(tagname, st) == false {
				logs.Warn("Tag [%s] no data in database before [%s]", tagname, begintime)
				return hismap, fmt.Errorf("0")
			}
			if IsExistItem(tagname, hisd) == false {
				return hismap, fmt.Errorf("Tag [%s] not exist in database", tagname)
			}
			if len(hisd[tagname]) > 0 { //有读取到历史数据
				//字符串转换为时间
				stdatatime, err := timeFormat(st[tagname].Time) //开始时间点的数据点的时间
				if err != nil {
					//logs.Alert("Time format error:%s", err.Error())
					return hismap, err
				}
				sthisdatatime, err := timeFormat(hisd[tagname][0].Time) //第一个历史数据点的时间
				if err != nil {
					//logs.Alert("Time format error:%s", err.Error())
					return hismap, err
				}
				//历史数据时间点小于开始时间点,或者开始时间点小于第一个历史数据的时间点
				if stdatatime.Unix() < begintime.Unix() || begintime.Unix() < sthisdatatime.Unix() {
					//要区分数据断更和数据有更新但值不变的情况*********
					var h HisData
					h.Time = cmd.BeginTime
					h.Value = st[tagname].Value
					hd = append(hd, h) //复制起始点数据
				}

				for _, v := range hisd[tagname] {
					hd = append(hd, v) //复制区间数据点
				}
				edv := hd[len(hd)-1].Value                        //最后一个历史数据点的值
				eddatatime, err := timeFormat(hd[len(hd)-1].Time) //最后一个历史数据点的时间
				if err != nil {
					return hismap, err
				}

				if eddatatime.UnixNano() < endtime.UnixNano() { //如果读取到的历史数据中的最后一个点的时间小于指定的结束时间
					var h HisData
					h.Time = cmd.EndTime //将指定的结束时间点
					h.Value = edv        //和最后一个历史数据点的值作为一个新数据点
					hd = append(hd, h)   //插入到数据数组中
				}

			} else { //没有读到历史数据
				var h HisData
				h.Time = cmd.BeginTime
				h.Value = st[tagname].Value
				hd = append(hd, h) //复制起始点数据
				h.Time = cmd.EndTime
				hd = append(hd, h) //复制结束点数据
			}
		} else { //起始时间点之前没有数据
			if len(hisd[tagname]) > 0 { //有读取到历史数据
				for _, v := range hisd[tagname] {
					hd = append(hd, v) //复制区间数据点
				}
				edv := hd[len(hd)-1].Value                        //最后一个历史数据点的值
				eddatatime, err := timeFormat(hd[len(hd)-1].Time) //最后一个历史数据点的时间
				if err != nil {
					return hismap, err
				}
				if eddatatime.UnixNano() < endtime.UnixNano() { //如果读取到的历史数据中的最后一个点的时间小于指定的结束时间
					var h HisData
					h.Time = cmd.EndTime //将指定的结束时间点
					h.Value = edv        //和最后一个历史数据点的值作为一个新数据点
					hd = append(hd, h)   //插入到数据数组中
				}

			} else { //没有读到历史数据
				logs.Warn("Tag %s no data in database before %s,the last data time is:%s", tagname, cmd.EndTime, snapmp[i].Time)
				return hismap, errors.New("0")
			}
		}
		hismap[tagname] = hd
	}

	return hismap, nil
}

/*
功能：根据庚顿数据库的地址和读取指令获取历史统计数据
输入：add:Golden API的地址
	cmd:读取指令接口
输出：数据接口
时间：2019年12月7日
编辑：wang_jp
*/
func getHistorySummaryDataExi(add GoldenAdd, cmd HisSumCmd) (HisSumDataMapExi, error) {
	hismap := make(HisSumDataMapExi, len(cmd.TagName))
	hisdata, err := GetHistoryDataAlignHeadAndTail(add, cmd)
	if err != nil {
		return hismap, err
	}
	for tagname, hd := range hisdata {
		var hissum HisSumDataExi
		switch strings.ToLower(cmd.DataType) { //包含高级历史统计
		case "advangce", "sd", "stddev", "se", "ske", "kur", "mode", "median", "groupdist", "distribution":
			hissum = calcBaseStatistics(hd, 1, cmd.Group)
		default:
			hissum = calcBaseStatistics(hd, 0, cmd.Group)
		}
		hissum.BeginTime = cmd.BeginTime
		hissum.EndTime = cmd.EndTime
		hismap[tagname] = hissum
	}

	return hismap, nil
}

/*
功能：根据庚顿数据库的地址和读取指令获取快照数据
输入：add:Golden API的地址
	cmd:读取指令接口
输出：数据接口
时间：2019年11月28日
编辑：wang_jp
*/
func getSnapShotData(add GoldenAdd, cmd SnapshotCmd) (SnapDataMap, error) {
	data := SnapDataMap{}
	urlstr, err := getUrl(add, cmd) //拼接URL
	if err != nil {
		return data, err
	}
	resp, err := http.Get(urlstr) //Get数据
	if err != nil {
		return data, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body) //接收反馈数据
	if err != nil {
		return data, err
	}

	json.Unmarshal([]byte(body), &data)

	return data, nil
}

/*
功能：根据庚顿数据库的地址和指令往庚顿数据库快照中写数据
输入：add:Golden API的地址
	cmd:写数据指令接口
输出：写成功的数量,失败的数据
时间：2020年3月30日
编辑：wang_jp
说明:未通过测试,反馈为“未知方法”
*/
func WriteSnapShotData(add GoldenAdd, cmd SnapWriteCmd) (interface{}, error) {
	urlstr, _ := getUrl(add, cmd) //拼接URL
	resp, err := Post(urlstr, cmd.Datas, "charset=UTF-8")
	return resp, err
}

//发送POST请求
//url:请求地址，data:POST请求提交的数据,contentType:请求体格式，如：application/json
//content:请求放回的内容
func Post(url string, data interface{}, contentType string) (string, error) {
	jsonStr, _ := json.Marshal(data)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Add("content-type", contentType)
	req.Header.Add("Content-Length", "0")
	req.Header.Add("Host", "<calculated when request is sent>")
	req.Header.Add("Connection", "keep-alive")
	if err != nil {
		return "", err
	}
	defer req.Body.Close()

	client := &http.Client{Timeout: 1 * time.Second}
	resp, error := client.Do(req)
	if error != nil {
		return "", err
	}
	defer resp.Body.Close()

	result, _ := ioutil.ReadAll(resp.Body)
	content := string(result)
	return content, nil
}

/*
功能：变量信息
输入：add:Golden API的地址
	cmd:读取指令接口
输出：字符串数组
时间：2019年11月28日
编辑：wang_jp
*/
func getTagList(add GoldenAdd, cmd interface{}) ([]string, error) {
	var data []string
	urlstr, err := getUrl(add, cmd) //拼接URL
	if err != nil {
		return data, err
	}
	resp, err := http.Get(urlstr) //Get数据
	if err != nil {
		return data, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body) //接收反馈数据
	if err != nil {
		return data, err
	}

	json.Unmarshal([]byte(body), &data)

	return data, nil
}

/*
功能：获取表信息
输入：add:Golden API的地址
	cmd:读取指令接口
输出：字符串数组
时间：2019年11月28日
编辑：wang_jp
*/
func getTableList(add GoldenAdd, cmd interface{}) (interface{}, error) {
	var tablelist []string
	var idlist []int
	var tbinfo TableInfo
	ord, _ := cmd.(TableListCmd)    //断言接口数据格式
	urlstr, err := getUrl(add, cmd) //拼接URL
	if err != nil {
		return nil, err
	}
	resp, err := http.Get(urlstr) //Get数据
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body) //接收反馈数据
	if err != nil {
		return nil, err
	}
	switch strings.ToLower(ord.TableMsg) {
	case "":
		json.Unmarshal([]byte(body), &idlist)
		return idlist, nil
	case "all":
		json.Unmarshal([]byte(body), &tablelist)
		return tablelist, nil
	default:
		json.Unmarshal([]byte(body), &tbinfo)
		return tbinfo, nil
	}
	return nil, nil
}

/*
功能：根据庚顿数据库的变量标签名获取标签的基本信息
输入：add:Golden API的地址
	cmd:读取指令接口
输出：标签点信息
时间：2019年11月28日
编辑：wang_jp
*/
func getTagPointInfo(add GoldenAdd, cmd interface{}) (PointInfo, error) {
	var data PointInfo
	urlstr, err := getUrl(add, cmd) //拼接URL
	if err != nil {
		return data, err
	}
	resp, err := http.Get(urlstr) //Get数据
	if err != nil {
		return data, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body) //接收反馈数据
	if err != nil {
		return data, err
	}

	json.Unmarshal([]byte(body), &data)

	return data, nil
}

/*
功能：根据庚顿数据库服务器时间
输入：add:Golden API的地址
输出：字符串
时间：2019年11月28日
编辑：wang_jp
*/
func getServerTime(add GoldenAdd) (string, error) {
	urlstr, err := getUrl(add, "") //拼接URL
	if err != nil {
		return "", err
	}
	resp, err := http.Get(urlstr) //Get数据
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body) //接收反馈数据
	if err != nil {
		return "", err
	}
	return string(body), nil
}

/*
功能：根据庚顿数据库的地址和读取指令拼接URL
输入：add:Golden API的地址
	cmd:读取指令接口
输出：URL字符串
时间：2019年11月28日
编辑：wang_jp
*/
func getUrl(add GoldenAdd, cmd interface{}) (string, error) {
	url := fmt.Sprintf("http://%s:%d/api/", add.Host, add.Port)

	switch cmd.(type) {
	case SnapshotCmd:
		url = fmt.Sprintf("%s%s", url, snapshot_cmd(cmd))
		break
	case SnapWriteCmd:
		wdatas, err := snapwrite_cmd(cmd)
		if err != nil {
			return "", err
		}
		url = fmt.Sprintf("%s%s", url, wdatas)
		break
	case HistoryCmd:
		url = fmt.Sprintf("%s%s", url, history_cmd(cmd))
		break
	case HisIntervalCmd:
		hiscmd, err := hisInterval_cmd(cmd)
		if err != nil {
			return "", err
		}
		url = fmt.Sprintf("%s%s", url, hiscmd)
		break
	case TableListCmd:
		url = fmt.Sprintf("%s%s", url, table_cmd(cmd))
		break
	case TagListCmd:
		url = fmt.Sprintf("%s%s", url, taglist_cmd(cmd))
		break
	case HisSumCmd:
		url = fmt.Sprintf("%s%s", url, hissummary_cmd(cmd))
		break
	default: //没有匹配的类型情况下返回服务器时间
		url = fmt.Sprintf("%s%s", url, _Sever) //读取服务器时间
		break
	}
	return url, nil
}

/*
功能：获取读取快照指令参数
输入：读取快照指令结构体
输出：读取快照指令字符串
时间：2019年11月28日
编辑：wang_jp
*/
func snapshot_cmd(itf interface{}) string {
	cmd, _ := itf.(SnapshotCmd)                 //断言接口数据格式
	str := fmt.Sprintf("%s?tagName=", _SnapCmd) //读取快照指令
	for i := 0; i < len(cmd.TagName); i++ {     //逐条填写变量名
		if i == 0 {
			str = fmt.Sprintf("%s%s", str, url.QueryEscape(cmd.TagName[i]))
		} else {
			str = fmt.Sprintf("%s,%s", str, url.QueryEscape(cmd.TagName[i]))
		}
	}
	return str
}

/*
功能：获取写快照指令参数
输入：写快照指令结构体
输出：写快照指令字符串
时间：2020年3月30日
编辑：wang_jp
*/
func snapwrite_cmd(itf interface{}) (string, error) {
	//cmd, _ := itf.(SnapWriteCmd)                       //断言接口数据格式
	str := fmt.Sprintf("%s?%s=", _SnapCmd, _SnapWrite) //写快照指令
	/*datastr, err := json.Marshal(cmd.Datas)
	if err != nil {
		return str, err
	}
	str += string(datastr)
	*/
	return str, nil
}

/*
功能：获取读取历史统计指令参数
输入：读取历史统计结构体
输出：读取历史统计字符串
时间：2019年11月28日
编辑：wang_jp
*/
func hissummary_cmd(itf interface{}) string {
	cmd, _ := itf.(HisSumCmd)                  //断言接口数据格式
	str := fmt.Sprintf("%s?tagName=", _HisSum) //读取历史统计指令
	for i := 0; i < len(cmd.TagName); i++ {    //逐条填写变量名
		if i == 0 {
			str = fmt.Sprintf("%s%s", str, url.QueryEscape(cmd.TagName[i]))
		} else {
			str = fmt.Sprintf("%s,%s", str, url.QueryEscape(cmd.TagName[i]))
		}
	}
	str = fmt.Sprintf("%s&beginTime=%s&endTime=%s&dataType=summary", str, url.QueryEscape(cmd.BeginTime), url.QueryEscape(cmd.EndTime))
	return str
}

/*
功能：获取历史数据指令参数
输入：读取历史数据指令结构体
输出：读取历史数据指令字符串
时间：2019年11月28日
编辑：wang_jp
*/
func history_cmd(itf interface{}) string {
	cmd, _ := itf.(HistoryCmd)                 //断言接口数据格式
	str := fmt.Sprintf("%s?tagName=", _HisCmd) //读取历史数据指令
	for i := 0; i < len(cmd.TagName); i++ {    //逐条填写变量名
		if i == 0 {
			str = fmt.Sprintf("%s%s", str, url.QueryEscape(cmd.TagName[i]))
		} else {
			str = fmt.Sprintf("%s,%s", str, url.QueryEscape(cmd.TagName[i]))
		}
	}
	if cmd.BeginTime != "" { //如果开始时间不为空，则读取一段范围内的历史值
		str = fmt.Sprintf("%s&beginTime=%s&endTime=%s", str, url.QueryEscape(cmd.BeginTime), url.QueryEscape(cmd.EndTime))
	} else { //否则读取指定时间点上的历史值(如果该时间点没有值，则取该点之前的第一个有效值)
		str = fmt.Sprintf("%s&Time=%s", str, url.QueryEscape(cmd.TimePoint))
	}
	return str
}

/*
功能：获取等间隔差值历史数据指令参数
输入：读取等间隔差值历史数据指令结构体
输出：读取等间隔差值历史数据指令字符串
时间：2019年11月28日
编辑：wang_jp
*/
func hisInterval_cmd(itf interface{}) (string, error) {
	cmd, _ := itf.(HisIntervalCmd)                  //断言接口数据格式
	str := fmt.Sprintf("%s?tagName=", _HisInterCmd) //读取等间隔历史数据指令
	for i := 0; i < len(cmd.TagName); i++ {         //逐条填写变量名
		if i == 0 {
			str = fmt.Sprintf("%s%s", str, url.QueryEscape(cmd.TagName[i]))
		} else {
			str = fmt.Sprintf("%s,%s", str, url.QueryEscape(cmd.TagName[i]))
		}
	}
	st, err := timeFormat(cmd.BeginTime) //字符串时间转换为标准时间
	if err != nil {
		return "", err
	}
	et, err := timeFormat(cmd.EndTime) //字符串时间转换为标准时间
	if err != nil {
		return "", err
	}
	if cmd.Interval == 0 { //时间间隔如果为0
		cmd.Interval = 1 //则设置为间隔1秒
	}
	cmd.count = (et.Unix()-st.Unix())/cmd.Interval + 1 //计算数据点数量
	str = fmt.Sprintf("%s&beginTime=%s&endTime=%s&interval=%d&count=%d", str, url.QueryEscape(cmd.BeginTime), url.QueryEscape(cmd.EndTime), cmd.Interval, cmd.count)

	return str, nil
}

/*
功能：获取读取数据表列表指令参数
输入：读取数据表列表指令结构体
输出：读取数据表列表指令字符串
时间：2019年11月28日
编辑：wang_jp
*/
func table_cmd(itf interface{}) string {
	cmd, _ := itf.(TableListCmd) //断言接口数据格式
	var str string
	switch strings.ToLower(cmd.TableMsg) {
	case "all": //读取全部表列表指令
		str = fmt.Sprintf("%s?isAllTalbleInfo=all", _Table)
	case "": //读取全部表ID
		str = fmt.Sprintf("%s", _Table) //读取全部表列表指令
	default:
		str = fmt.Sprintf("%s?tablename=%s", _Table, cmd.TableMsg)
	}
	return str
}

/*
功能：获取读取指定表中变量列表指令参数
输入：读取变量列表指令结构体
输出：读取变量列表指令字符串
时间：2019年11月28日
编辑：wang_jp
*/
func taglist_cmd(itf interface{}) string {
	var str string
	cmd, _ := itf.(TagListCmd) //断言接口数据格式
	if len(cmd.TagName) > 0 {  //填写了变量名称，则指定变量名的变量点信息
		str = fmt.Sprintf("%s?tagName=", _Point) //读取指令
		for i := 0; i < len(cmd.TagName); i++ {  //逐条填写变量名
			if i == 0 {
				str = fmt.Sprintf("%s%s", str, url.QueryEscape(cmd.TagName[i]))
			} else {
				str = fmt.Sprintf("%s,%s", str, url.QueryEscape(cmd.TagName[i]))
			}
		}

	} else { //否则读取变量表中的变量列表
		str = fmt.Sprintf("%s?tableName=%s", _Point, url.QueryEscape(cmd.TableName)) //读取指令
	}
	return str
}

/*
功能：时间参数格式化
输入：时间字符串
输出：格式化后的时间变量,错误信息
时间：2019年11月28日
编辑：wang_jp
*/
func timeFormat(s string) (time.Time, error) {
	t, err := time.ParseInLocation("2006-01-02 15:04:05", s, time.Local)
	if err == nil {
		return t, nil
	}
	t, err = time.ParseInLocation("2006/01/02 15:04:05", s, time.Local)
	if err == nil {
		return t, nil
	}
	t, err = time.ParseInLocation("2006/1/2 15:04:05", s, time.Local)
	if err == nil {
		return t, nil
	}
	return t, err
}

/*
功能：判断元素在数组、Map中是否存在
输入：元素、数组或者Map、Slice
输出：存在输出true，不存在输出false
说明：对于数组、Slice，判断的是值是否存在，对于Map，判断的是Key是否存在
时间：2019年12月15日
编辑：wang_jp
*/
func IsExistItem(obj interface{}, target interface{}) bool {
	targetValue := reflect.ValueOf(target)
	switch reflect.TypeOf(target).Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < targetValue.Len(); i++ {
			if targetValue.Index(i).Interface() == obj {
				return true
			}
		}
	case reflect.Map:
		if targetValue.MapIndex(reflect.ValueOf(obj)).IsValid() {
			return true
		}
	}
	return false
}

/*
功能:基础统计计算
输入：data []hisData:历史数据数组
	isNeedAdvangce int:是否需要进行高级统计计算,1=需要,0=不需要
	group int:高级计算时的分组数量,默认0=100
返回:统计计算结果
编辑:wang_jp
时间:2019年12月7日
*/
func calcBaseStatistics(data []HisData, isNeedAdvangce, group int) HisSumDataExi {
	var tsds []statistic.TimeSeriesData
	var tsd statistic.TimeSeriesData
	for _, v := range data {
		tsd.Time, _ = statistic.TimeParse(v.Time)
		vl, err := strconv.ParseFloat(v.Value, 64)
		if err != nil {
			vl = 0
		}
		tsd.Value = vl
		tsds = append(tsds, tsd)
	}
	return HisSumDataExi(statistic.BaseStatistics(tsds, isNeedAdvangce, group))
}
