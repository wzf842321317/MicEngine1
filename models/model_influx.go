package models

import (
	"github.com/influxdata/influxdb1-client/v2"
)

/***************************************************************
函数名：queryInfluxDB(clnt client.Client, MyDB, cmd string) (res []client.Result, err error)
功  能：InfluxDB查询函数
参  数:clnt client.Client:连接信息
      MyDB string:数据库名
	  cmd string:查询命令
返回值:res []client.Result:查询结果集
      err error:错误信息
创建时间:2019年1月12日
修订信息:
***************************************************************/
func queryInfluxDB(clnt client.Client, MyDB, cmd string) (res []client.Result, err error) {
	q := client.Query{
		Command:  cmd,
		Database: MyDB,
	}
	if response, err := clnt.Query(q); err == nil {
		if response.Error() != nil {
			return res, response.Error()
		}
		res = response.Results
	} else {
		return res, err
	}
	return res, nil
}
