package main

import (
	"sync"
	"time"

	"github.com/bkzy-wangjp/MicEngine/models"

	"github.com/astaxie/beego/logs"
	drum "github.com/bkzy-wangjp/drumstick"
)

type pkpiresult struct {
	res    models.CalcKpiResult
	number int
}

type kpiInfo struct {
	cfg    models.CalcKpiConfigListExi
	res    chan<- pkpiresult
	number int
}

func ParallelRun() {
	logs.Info("快速计算功能已开启……")
	time.Sleep(time.Second * time.Duration(models.EngineCfgMsg.Sys.SleepTimeOnStart)) //延迟执行
	for {
		cfg := new(models.CalcKpiConfigList)
		if cnt, kpis, err := cfg.GetKpiConfigInfo(2); err != nil { //获取并行指标数量
			logs.Alert("获取并行计算指标信息时发生错误,错误信息为:", err)
			break
		} else {
			pCnt := cnt
			if pCnt > models.EngineCfgMsg.CfgMsg.ParallelIndicatorAuth { //不可大于授权数量
				pCnt = models.EngineCfgMsg.CfgMsg.ParallelIndicatorAuth
			}
			logs.Info("共查询到%d个快速计算KPI指标,授权计算%d个,实际计算%d个", cnt, models.EngineCfgMsg.CfgMsg.ParallelIndicatorAuth, pCnt)
			task := make([]*drum.Task, pCnt)                         //定义循环任务
			kpichan := make(chan kpiInfo, pCnt)                      //传送KPI参数的通道
			resultchan := make(chan pkpiresult, pCnt)                //传送KPI结果的通道
			wait := sync.WaitGroup{}                                 //同步等待组
			done := make(chan struct{})                              //重载通道
			resBuff := make(map[int64][2]models.CalcKpiResult, pCnt) //计算结果缓冲池
			go InitTask(kpichan, resultchan, kpis, pCnt)             //初始化分配任务
			for i := 0; i < int(pCnt); i++ {                         //等待初始化配置结果,逐个设置并行任务
				k := <-kpichan
				k.number = i
				wait.Add(1) //启动一个同步等待
				task[i], err = drum.NewTask(time.Now(), time.Duration(k.cfg.Period)*time.Second, parallelKpiCalc, k, resultchan)
				task[i].Start()
			}
			go ReloadParallelKpiCfgInfoTimer(done) //启动重载计时器
		Listen:
			for { //循环监听计算结果或者计时器完成标志
				select {
				case v := <-resultchan: //监听到有计算完成的结果
					resBuff[v.res.KpiConfigListId] = kpiResultBuffCheck(resBuff[v.res.KpiConfigListId], v) //检查是否需要存储(每个重载周期必须存储一次)
				case <-done: //监听重载计时器
					break Listen //跳出监听
				}
			}
			for i := 0; i < int(pCnt); i++ { //等待开启的计算全部完成
				v := <-resultchan
				resBuff[v.res.KpiConfigListId] = kpiResultBuffCheck(resBuff[v.res.KpiConfigListId], v)
				if task[v.number] != nil {
					task[v.number].Stop() //停止计算任务
				}
				wait.Done() //关闭一个同步等待
			}
			wait.Wait()       //等待所有等待状态完成
			close(resultchan) //关闭结果通道
		}
	}
}

/*
功能:重载并行KPI配置信息计时器
输入:
	done:计时完成通道
输出:无
说明:
编辑:wang_jp
时间:2019年12月19日
*/
func ReloadParallelKpiCfgInfoTimer(done chan<- struct{}) {
	time.Sleep(time.Second * time.Duration(models.EngineCfgMsg.CfgMsg.ParallelIndicReloadInterval))
	logs.Info("准备重新载入快速计算KPI指标……")
	close(done)
}

/*
功能:初始化分配并行计算任务
输入:
	kpichan:指标配置信息通道
	r:结果通道
	kpis:原始配置信息
	pCnt:并行指标数量（不大于授权数量）
输出:无
说明:
编辑:wang_jp
时间:2019年12月19日
*/
func InitTask(kpichan chan<- kpiInfo, r chan pkpiresult, kpis []models.CalcKpiConfigListExi, pCnt int64) {
	defer close(kpichan)
	for i, k := range kpis {
		if int64(i) < pCnt { //实际计算数量不大于授权数量
			kpii := kpiInfo{
				cfg: k,
				res: r,
			}
			kpichan <- kpii
		} else { //超出授权计算数量的不予计算
			break
		}
	}
}

/*
功能:并行KPI指标计算GO程程序
输入:
	kif:指标配置信息
	result:结果输出通道
输出:无
说明:
编辑:wang_jp
时间:2019年12月19日
*/
func parallelKpiCalc(kif kpiInfo, result chan pkpiresult) {
	kpi := kif.cfg
	lasttime := time.Now().Add(time.Duration(models.EngineCfgMsg.CfgMsg.SerialCalcDelaySec) * -1 * time.Second)
	endTime := lasttime.Format(models.EngineCfgMsg.Sys.TimeFormat)
	bgTime := lasttime.Add(time.Duration(kpi.Period) * time.Second * -1).Format(models.EngineCfgMsg.Sys.TimeFormat)

	res, ctinue, err := kpi.KpiCalc(bgTime, endTime) //执行计算
	if err != nil {                                  //计算结果有错误
		if ctinue {
			logs.Warn(err.Error())
		} else {
			logs.Error(err.Error())
		}
	} else { //没有错误
		if res.KpiConfigListId != 0 { //返回了计算结果
			kpi.LastCalcTime = endTime
			go kpi.SetKpiLastCalcTimeSingle(kpi.Id, endTime)
			kpir := pkpiresult{
				res:    res,
				number: kif.number,
			}
			kif.res <- kpir
		}
	}
}

func kpiResultBuffCheck(resb [2]models.CalcKpiResult, resst pkpiresult) [2]models.CalcKpiResult {
	resb[1] = resst.res
	if len(resb[0].CalcEndingTime) == 0 || resb[0].KpiValue != resb[1].KpiValue { //第一指标为空或者两个指标的值不相等时触发存储
		resb[0] = resb[1]
		r := make([]models.CalcKpiResult, 1)
		r = append(r, resb[0])
		go func() {
			result := new(models.CalcKpiResult)
			result.SaveBatchKpiResultToDB(r) //触发存储GO程
		}()
	}
	return resb
}
