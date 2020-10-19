package models

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego/logs"

	"github.com/astaxie/beego/orm"

	"github.com/bkzy-wangjp/MicEngine/MicScript/Calc"
	"github.com/bkzy-wangjp/MicEngine/MicScript/statistic"
)

const (
	_timeStr  = `\d{4}(\-|\/)\d{1,2}(\-|\/)\d{1,2}(\s|T)\d{1,2}\:\d{1,2}\:\d{1,2}` //纯时间格式
	_timeDev  = `t_add\(\s*` + _timePara + `{1}\s*\,\s*\-?[0-9nuµmsh\.]+\s*\)`     //时间运算函数格式
	_timePara = `((` + _timeStr + `)|(` + _Field + `)|bgoy|bgom|bgod|bgos|now|0)`  //组成函数中时间的参数格式
	_funcTime = `((` + _timeDev + `)|(` + _timePara + `))`

	//解析函数
	//过程实时数据统计算法
	_funcReg = `(this\.)?(tag\(\s*[\w\.\:\-]+\s*\)\.)?fc\(\s*\w+(\,\s*` + _funcTime + `{1}\s*\,\s*` + _funcTime + `{1})?\s*\)`
	//平台内的动态数据
	_srtd = `(tag\(\s*[\w\.\:\-]+\s*\)\.)?srtd\(\s*\w+\,\w+(\,\s*` + _funcTime + `{1}\s*\,\s*` + _funcTime + `{1})?\s*\)`
	_kpi  = `(kpi\(\s*[\w\:\.\-]+\,\w+(\,\s*` + _funcTime + `{1}\s*\,\s*` + _funcTime + `{1})?\s*\))` //KPI计算结果数据
	//数学符号
	_mathSymbol = `(\+|\-|\*|\/|\%|\>|\<|\=|\(\s*|\s*\)){1}|(\>\=|\<\=|\=\=){2}`
	//SQL语句
	//_select = `select\(([a-zA-Z_\,\.]+\([\w\s\.\,\=]+\))*[\w\s\.\*\+\-\/\=\>\<\'\,\:\%\!]*\)\.`
	//_from       = `from\((\([\w\s\.\=]+\))*[\w\s\.\*\+\-\/\=]*\)\.?`
	//_where      = `where\((\([\w\s\.\*\+\-\/\=\>\<\'\,\:\%\!]+\))*[\w\s\.\*\+\-\/\=\>\<\'\,\:\%\!]*\)\.?`
	_select     = `select\([[:print:]]*\)\.`
	_from       = `from\([[:print:]]*\)\.?`
	_timecolumn = `timecolumn\([\w\s\.]+\)\.?`
	_where      = `where\([[:print:]]*\)\.?`
	_timefilter = `timefilter\(\s*` + _funcTime + `\,` + _funcTime + `\s*\)\.?`
	_groupby    = `groupby\([\w\s\,\.]+\)\.?`
	_orderby    = `orderby\([\w\s\,\.]+\)\.?`
	_limit      = `limit\([\s\d]+\)\.?`
	_as         = `as\((string|value|map|json)\)`
	_sql        = `(` + _select + _from + `(` + _where + `)?` + `(` + _timecolumn + `)?` + `(` + _timefilter + `)?` + `(` + _groupby + `)?` + `(` + _orderby + `)?` + `(` + _limit + `)?` + `(` + _as + `)?` + `)`

	_gdsum  = `gdsum\(\s*[\w\-]*\s*\)`                           //物耗求和
	_date   = `date\(` + _funcTime + `\,[tTYMDymdWwhHfFsS]{1}\)` //日期
	_number = `\d+\.*\d*`                                        //解析数字
	//解析数据表单元格
	_Field = `((table\(\s*\w+\s*\)\.row\(\s*\d+\s*\))|(` + _tag + `\.table\.row)|(this\.table\.row)|(` + _tag + `)|(this)){1}\.field\(\s*\w+\s*\)`

	_thisField         = `this\.field\(\s*\w+\s*\)`                                 //当前tag在taglist当中的某列值
	_tagField          = _tag + `\.field\(\s*\w+\s*\)`                              //某个tag在taglist中的某列值
	_thisTableRowField = `this\.table\.row\.field\(\s*\w+\s*\)`                     //当前tag在其所属的 电机/仪表/分析仪/设备 表中某列中的值
	_tagTableRowField  = _tag + `\.table\.row\.field\(\s*\w+\s*\)`                  //某个tag在其所属的 电机/仪表/分析仪/设备 表中某列中的值
	_tableRowField     = `table\(\s*\w+\s*\)\.row\(\s*\d+\s*\)\.field\(\s*\w+\s*\)` //某个表某行某列中的值
	_tag               = `tag\(\s*[\w\.\:\-]+\s*\)`                                 //某个tag
)

//脚本结构体
type Script struct {
	Id              int64  //ID
	BeginTime       string //本周期开始时间
	EndTime         string //本周期结束时间
	BaseTime        string //基准时间
	ShiftHour       int64  //每班工作时间
	ScriptStr       string //脚本程序
	MainTagId       int64  //主变量ID(针对对于实时数据库变量,可选)
	MainTagFullName string //主变量全名(针对对于实时数据库变量,可选)
	OnlySQL         bool   //只有SQL脚本
}

/*******************************************************************************
功能:运行脚本程序
输入:当前tagid,每班工作时间,脚本字符串,统计起始时间,统计结束时间,基准时间
输出:计算结果
编辑:wang_jp
时间:2019年12月8日
*******************************************************************************/
func (s *Script) Run() (interface{}, bool, error) {
	defer func() {
		if err := recover(); err != nil {
			logs.Critical("Script.Run中遇到错误:%d;[%#v]", s.Id, err)
		}
	}()
	errs := s.Compile()
	if len(errs) > 0 { //编译错误
		estr := ""
		for i, err := range errs {
			estr += fmt.Sprintf("错误%d %s; ", i+1, err.Error())
		}
		return nil, false, fmt.Errorf(estr)
	} else {
		if s.OnlySQL {
			return s.extractSql(s.ScriptStr)
		}
		if mat, err := regexp.MatchString(_date, s.ScriptStr); err == nil && mat == true { //匹配时间参数，解析date函数
			data, err := s.extractDateTime()
			return data, false, err
		} else {
			if expression, ctinue, err := s.analysisScript(true); err != nil { //分析脚本
				return 0, ctinue, err
			} else {
				data, err := calc.CalcAndCompare(expression) //执行计算
				return data, false, err
			}
		}
	}
}

/*******************************************************************************
功能:把脚本字符串转换为可以计算的词元数组
输入:当前tagid,每班工作时间,thisid的tagfullname,脚本字符串,统计起始时间,统计结束时间,基准时间,是否校验float类型参数
输出:[string] 词元字符串(该词元字符串仅包含数字和运算符)
	[bool] 本周期没有数据,但本周期之后有数据时,为true,否则为false
	[error] 错误信息
编辑:wang_jp
时间:2019年12月8日
*******************************************************************************/
func (s *Script) analysisScript(checkfloat bool) (string, bool, error) {
	defer func() {
		if err := recover(); err != nil {
			logs.Critical("Script.analysisScript 中遇到错误:%d;[%#v]", s.Id, err)
		}
	}()

	token := s.scriptStr2TokenArr()
	for i, tok := range token {
		//检查是否包含fc(xxx)/this.fc(xxx)/tag(y).fc(xxx)
		if mat, err := regexp.Match(_funcReg, []byte(tok)); err != nil {
			return tok, false, err
		} else if mat == true { //包含fc的表达式
			if v, ctinue, err := s.extractFc(tok); err != nil {
				return tok, ctinue, err
			} else {
				token[i] = v
				continue
			}
		}
		//检查是否包含srtd(type,func)/tag(y).srtd(type,func)
		if mat, err := regexp.Match(_srtd, []byte(tok)); err != nil {
			return tok, false, err
		} else if mat == true { //包含srtd的表达式
			if v, ctinue, err := s.extractSrtd(tok); err != nil {
				return tok, ctinue, err
			} else {
				token[i] = v
				continue
			}
		}
		//检查是否包含 kpi(tag,func)
		if mat, err := regexp.Match(_kpi, []byte(tok)); err != nil {
			return tok, false, err
		} else if mat == true { //包含kpi的表达式
			if v, ctinue, err := s.extractKpi(tok); err != nil {
				return tok, ctinue, err
			} else {
				token[i] = v
				continue
			}
		}
		//检查是否包含this.field(yyy)、tag(x).field(yyy)的脚本
		pat := fmt.Sprintf("%s|%s", _thisField, _tagField)
		if mat, err := regexp.Match(pat, []byte(tok)); err != nil {
			return tok, false, err
		} else if mat == true { //检查是否包含的表达式
			if res, err := s.extractMsgInBracket(tok, 2); err != nil { //解析脚本
				return tok, false, err //有错误
			} else { //没有错误
				tagid, _ := strconv.ParseInt(res[0], 0, 64)
				field_name := res[1]
				tag := new(OreProcessDTaglist)
				tag.Id = tagid
				if v, err := tag.GetTagListFieldValueByID(field_name); err != nil { //获取数据替换token
					return tok, false, err //有错误的时候返回存在问题的脚本字符串和错误信息
				} else {
					if checkfloat == true { //检验是否可以转换为float
						_, err := strconv.ParseFloat(v, 64)
						if err != nil {
							return tok, false, err
						}
					}
					token[i] = v //没有错误的时候继续
					continue
				}
			}
		}
		//检查是否包含this.table.row.field(yyy)\tag(x).table.row.field(yyy) 的脚本
		pat = fmt.Sprintf("%s|%s", _thisTableRowField, _tagTableRowField)
		if mat, err := regexp.Match(pat, []byte(tok)); err != nil {
			return tok, false, err
		} else if mat == true { //检查是否包含的表达式
			if res, err := s.extractMsgInBracket(tok, 2); err != nil { //解析脚本
				return tok, false, err //有错误
			} else { //没有错误
				tagid, _ := strconv.ParseInt(res[0], 0, 64)
				field_name := res[1]
				tag := new(OreProcessDTaglist)
				tag.Id = tagid
				if v, err := tag.GetTagTableRowFieldValueByID(field_name); err != nil { //获取数据替换token
					return tok, false, err //有错误的时候返回存在问题的脚本字符串和错误信息
				} else {
					if checkfloat == true { //检验是否可以转换为float
						_, err := strconv.ParseFloat(v, 64)
						if err != nil {
							return tok, false, err
						}
					}
					token[i] = v //没有错误的时候继续
					continue
				}
			}
		}
		//检查是否包含取table(xxxx).row(y).field(zzz)的脚本
		pat = fmt.Sprintf("%s", _tableRowField)
		if mat, err := regexp.Match(pat, []byte(tok)); err != nil {
			return tok, false, err
		} else if mat == true { //检查是否包含的表达式
			if res, err := s.extractMsgInBracket(tok, 3); err != nil { //解析脚本
				return tok, false, err //有错误
			} else { //没有错误
				category := res[0] //表名或者表ID
				itemid, _ := strconv.ParseInt(res[1], 0, 64)
				field_name := res[2]
				if v, err := GetTableRowFieldValueByID(category, itemid, field_name); err != nil { //获取数据替换token
					return tok, false, err //有错误的时候返回存在问题的脚本字符串和错误信息
				} else {
					//测试是否可以转换为浮点数,不可以的话抛出错误信息
					if checkfloat == true { //检验是否可以转换为float
						_, err := strconv.ParseFloat(v, 64)
						if err != nil {
							return tok, false, err
						}
					}
					token[i] = v //没有错误的时候继续
					continue
				}
			}
		}
		//检查是否包含gdsum(xxx)
		if strings.Contains(tok, "gdsum(") { //检查是否包含fc的表达式
			if res, err := s.extractMsgInBracket(tok, -1); err != nil { //解析脚本
				return tok, false, err //有错误
			} else { //没有错误
				goodsid, err := strconv.ParseInt(res[0], 0, 64)
				if err != nil {
					gd := new(GoodsConfigInfo)
					goodsid, err = gd.GetGoodsTagIDByTagName(res[0])
					if err != nil {
						return tok, false, err //有错误的时候返回存在问题的脚本字符串和错误信息
					}
				}
				gd := new(GoodsConsumeInfo)
				if v, err := gd.GetGoodsSumByID(goodsid, s.BeginTime, s.EndTime); err != nil { //获取数据替换token
					return tok, false, err //有错误的时候返回存在问题的脚本字符串和错误信息
				} else {
					token[i] = fmt.Sprint(v) //没有错误的时候继续
					continue
				}
			}
		}
		//检查是否包含select语句
		if strings.Contains(tok, "select(") { //检查是否包含select语句
			if v, ctinue, err := s.extractSql(tok); err != nil { //解析脚本
				return tok, ctinue, err //有错误
			} else { //没有错误
				value, ok := v.(string)
				if ok {
					token[i] = value //没有错误的时候继续
				}
				continue
			}
		}
	}
	var tokstr string
	for _, v := range token {
		tokstr = fmt.Sprintf("%s%s", tokstr, v)
	}
	return tokstr, true, nil
}

/*******************************************************************************
功能:合成正则字符串
输入:无
输出:正则字符串
编辑:wang_jp
时间:2019年12月8日
*******************************************************************************/
func (s *Script) getRegexpString() string {
	//合成正则表达式
	return fmt.Sprintf("%s|%s|%s|%s|%s|%s|%s|%s",
		_sql,
		_funcReg,
		_srtd,
		_kpi,
		_mathSymbol,
		_Field,
		_gdsum,
		_number)
}

/*******************************************************************************
功能:把脚本字符串转换为词元数组
输入:脚本字符串
输出:词元数组(词元包含脚本函数,扔需要进一步处理)
编辑:wang_jp
时间:2019年12月8日
*******************************************************************************/
func (s *Script) scriptStr2TokenArr() []string {
	//合成正则表达式
	pattern := s.getRegexpString()
	reg := regexp.MustCompile(pattern)
	return reg.FindAllString(s.ScriptStr, -1)
}

/*******************************************************************************
功能:从词元中解析词元括号中的内容
输入:当前tagid,待解析的词元,预期结果数组的类型(-1:gdsum(id)中的id,1:tagid和函数名,2=tagid和列名,3=tableid\itemid\列名)
输出:解析出的信息数组、错误信息
编辑:wang_jp
时间:2019年12月8日
*******************************************************************************/
func (s *Script) extractMsgInBracket(token string, restype int) ([]string, error) {
	reg := regexp.MustCompile(`\(\s*[\w\.\:\-]*\s*\)`) //提取括号内的值，含括号
	str := reg.FindAllString(token, -1)
	//提取括号内的值
	reg = regexp.MustCompile(`[\w\.\:\-]*`)
	for i, s := range str {
		c := reg.FindAllString(s, -1)
		s = ""
		for _, m := range c {
			s = fmt.Sprintf("%s%s", s, m)
		}
		str[i] = s
	}
	if len(str) > 0 { //提取到的数据长度必须大于0
		if restype < 3 { //期望结果小于3的情况下
			if len(str) < 3 { //只有长度小于3才合法
				if len(str) == 2 { //长度等于2的情况
					if tagid, err := strconv.ParseInt(str[0], 0, 64); err != nil { //测试脚本中提取出的tagid是否数字
						tagnames := strings.Split(str[0], ".")
						var tagname string
						if len(tagnames) > 1 {
							tagname = tagnames[1]
						} else {
							tagname = str[0]
						}
						var rtag OreProcessDTaglist
						rtag.TagName = tagname
						tid, err := rtag.GetTagIDByTagName()
						if err != nil {
							return str, fmt.Errorf("错误:TagId=%d,脚本[%s]中的Tag(%s)中参数错误;错误信息:%s", s.MainTagId, token, str[0], err.Error())
						}
						str[0] = fmt.Sprint(tid)
					} else {
						if tagid <= 0 { //tagid不可以小于或者等于0
							str[0] = fmt.Sprint(s.MainTagId)
						}
					}
				} else if len(str) == 1 { //如果只提取到了一个值,
					if restype == -1 { //是gdsum(id)或者gdsum()的情况
						if len(str[0]) == 0 {
							str[0] = fmt.Sprint(s.MainTagId)
						}
						return str, nil
					} else { //不是gdsum(id)的情况，扩充成两个,第一个是tagid,第二个是函数名/列名
						str = append(str, str[0])
						str[0] = fmt.Sprint(s.MainTagId)
					}
				}
				if len(str[1]) < 2 { //列名和函数名不可能小于两个字符
					t := "函数名"
					if restype == 2 {
						t = "列名"
					}
					return str, fmt.Errorf("错误:TagId=%d,脚本[%s]中没有包含%s", s.MainTagId, token, t)
				} else { //返回正确信息
					return str, nil
				}
			} else { //结果长度大于3的报错信息
				return str, fmt.Errorf("错误:TagID=%d,解析出的数据过多,期望为2个,解析出了%d个;脚本为[%s]", s.MainTagId, len(str), token)
			}
		} else { //期望结果不小于3的情况
			if len(str) == 3 { //如果提取到了三个值,分别是category,itemid,列名
				if category, err := strconv.ParseInt(str[0], 0, 64); err == nil && category <= 0 {
					return str, fmt.Errorf("错误:TagId=%d,错误脚本[%s]中的category()值不应小于或等于0;", s.MainTagId, token)
				}
				if itemid, err := strconv.ParseInt(str[1], 0, 64); err != nil {
					return str, fmt.Errorf("错误:TagId=%d,错误脚本[%s]中的item()中包含非数字字符;错误信息:%s", s.MainTagId, token, err.Error())
				} else if itemid <= 0 {
					return str, fmt.Errorf("错误:TagId=%d,错误脚本[%s]中的item()值不应小于或等于0;", s.MainTagId, token)
				}
				return str, nil //没有错误后返回正确值
			} else {
				return str, fmt.Errorf("错误:TagID=%d,解析出的数据数量不对,期望为3个,解析出了%d个;脚本为[%s]", s.MainTagId, len(str), token)
			}
		}
	} else if restype == -1 { //是gdsum()的情况
		str = append(str, fmt.Sprint(s.MainTagId))
		return str, nil
	}
	return str, fmt.Errorf("错误:TagID=%d,没有解析出和合法数据;脚本[%s]", s.MainTagId, token)
}

/*******************************************************************************
功能:解析FC函数
输入:脚本字符串
输出:[string]词元字符串(该词元字符串仅包含数字)
	[bool] 最新的数据时间大于endTime时为ture
	[error]错误信息
编辑:wang_jp
时间:2020年01月8日
*******************************************************************************/
func (s *Script) extractFc(tok string) (string, bool, error) {
	defer func() {
		if err := recover(); err != nil {
			logs.Critical("Script.extractFc 中遇到错误:%d;[%#v]", s.Id, err)
		}
	}()

	if strings.Contains(tok, ",") == false { //不包含逗号,即fc(xxx)/this.fc(xxx)/tag(y).fc(xxx)模式
		if res, err := s.extractMsgInBracket(tok, 1); err != nil { //解析脚本
			return tok, false, err //有错误
		} else { //没有错误
			var tgname string
			tagid, err := strconv.ParseInt(res[0], 0, 64)
			if err != nil {
				tgname = res[0]
				var rtag OreProcessDTaglist
				rtag.TagName = tgname
				tid, err := rtag.GetTagIDByTagName()
				if err != nil {
					//return tok, fmt.Errorf("The tagname [%s] not exist in database of [%s]", tgname, tok)
					return tok, false, fmt.Errorf("脚本程序 [%s] 中的变量名 [%s] 在数据库中不存在", tok, tgname)
				}
				tagid = tid
			} else if tagid == s.MainTagId { //如果tagid等于当前id
				tgname = s.MainTagFullName //变量名等于传入的变量名
			}
			key := res[1]
			micgd := new(MicGolden)
			value, ctinue, err := micgd.GetTagHisSumValueFromGoldenByID(tagid, tgname, key, s.BeginTime, s.EndTime)
			return fmt.Sprint(value), ctinue, err
		}
	} else { //包含逗号,即fc(xxx,bgt,edt)/this.fc(xxx,bgt,edt)/tag(y).fc(xxx,bgt,edt)模式
		tagid := s.MainTagId
		tgname := ""
		//如果包含tag,则提取tagid
		reg := regexp.MustCompile(_tag)
		if reg.Match([]byte(tok)) { //包含tag
			tag := reg.FindAllString(tok, -1) //提取出tag(x)
			r := regexp.MustCompile(`[\w\:\.\-]+`)
			for i, s := range tag { //遍历tag(x)的提取结果
				if i == 0 && len(s) > 3 { //只对第一个进行处理
					p := s[3:]                    //去除"tag"
					ids := r.FindAllString(p, -1) //提取括号中的参数
					if len(ids) > 0 {
						id, err := strconv.ParseInt(ids[0], 0, 64) //字符串转换成数字
						if err != nil {                            //转换成数字失败
							tgname = ids[0] //参数是tagname
							var rtag OreProcessDTaglist
							rtag.TagName = tgname
							tid, err := rtag.GetTagIDByTagName()
							if err != nil {
								//return tok, fmt.Errorf("The tagname [%s] not exist in database of [%s]", tgname, tok)
								return tok, false, fmt.Errorf("脚本程序 [%s] 中的变量名 [%s] 在数据库中不存在", tok, tgname)
							}
							tagid = tid
						} else {
							tagid = id
						}
					} else { //没有提取到数字
						//return tok, fmt.Errorf("The parameter of tag(*) must be int or string in [%s]", tok)
						return tok, false, fmt.Errorf("脚本程序[%s]中的 tag(*) 的括号中必须是整数或者字符串", tok)
					}
					break
				}
			}
		}
		if tagid == s.MainTagId { //如果tagid等于当前id
			tgname = s.MainTagFullName //变量名等于传入的变量名
		}
		key, bgtime_token, edtime_token, err := s.getFcParameter(tok)
		if err != nil { //提取参数错误
			return tok, false, fmt.Errorf("提取[fc]函数参数时发生错误:[%s]", err.Error())
		}
		bgtime, err := s.extractFuncTimeParamater(bgtime_token, s.BeginTime)
		if err != nil { //解析时间参数错误
			return tok, false, fmt.Errorf("解析时间参数时发生错误,extractFuncTimeParamater(%q,%q) 错误信息:[%s]", bgtime_token, s.BeginTime, err.Error())
		}
		edtime, err := s.extractFuncTimeParamater(edtime_token, s.EndTime)
		if err != nil { //解析时间参数错误
			return tok, false, fmt.Errorf("解析时间参数时发生错误,extractFuncTimeParamater(%q,%q) 错误信息:[%s]", bgtime_token, s.EndTime, err.Error())
		}
		micgd := new(MicGolden)
		value, ctinue, err := micgd.GetTagHisSumValueFromGoldenByID(tagid, tgname, key, bgtime, edtime)
		return fmt.Sprint(value), ctinue, err
	}
}

/*******************************************************************************
功能:解析srtd函数
输入:脚本字符串
输出:[string] 词元字符串(该词元字符串仅包含数字和运算符)
	[bool] 本周期没有数据,但本周期之后有数据时,为true,否则为false
	[error] 错误信息
编辑:wang_jp
时间:2020年01月8日
*******************************************************************************/
func (s *Script) extractSrtd(tok string) (string, bool, error) {
	tagid := s.MainTagId
	tgname := ""
	//如果包含tag,则提取tagid/tagname
	reg := regexp.MustCompile(_tag)
	if reg.Match([]byte(tok)) { //包含tag
		tag := reg.FindAllString(tok, -1) //提取出tag(x)
		r := regexp.MustCompile(`[\w\.\:\-]+`)
		for i, s := range tag { //遍历tag(x)的提取结果
			if i == 0 && len(s) > 3 { //只对第一个进行处理
				p := s[3:]                    //去除"tag"
				ids := r.FindAllString(p, -1) //提取括号中的参数
				if len(ids) > 0 {
					id, err := strconv.ParseInt(ids[0], 0, 64) //字符串转换成数字
					if err != nil {                            //转换成数字失败
						tgname = ids[0] //参数是tagname
						tagid = 0
					} else {
						tagid = id
					}
				} else { //没有提取到数字
					//return tok, fmt.Errorf("The parameter of tag(*) must be int or string in [%s]", tok)
					return tok, false, fmt.Errorf("脚本程序[%s]中的 tag(*) 的括号中必须是整数或者字符串", tok)
				}
				break
			}
		}
	}

	tagtype, key, bgtime_token, edtime_token, err := s.getSrtdParameter(tok) //提取函数的参数
	if err != nil {                                                          //提取参数错误
		//return tok, fmt.Errorf("getSrtdParameter error:[%s]", err.Error())
		return tok, false, fmt.Errorf("提取[srtd]函数参数时发生错误:[%s]", err.Error())
	}
	bgtime, err := s.extractFuncTimeParamater(bgtime_token, s.BeginTime)
	if err != nil { //解析时间参数错误
		return tok, false, fmt.Errorf("解析时间参数时发生错误,extractFuncTimeParamater(%q,%q) 错误信息:[%s]", bgtime_token, s.BeginTime, err.Error())
	}
	edtime, err := s.extractFuncTimeParamater(edtime_token, s.EndTime)
	if err != nil { //解析时间参数错误
		return tok, false, fmt.Errorf("解析时间参数时发生错误,extractFuncTimeParamater(%q,%q) 错误信息:[%s]", edtime_token, s.EndTime, err.Error())
	}
	//fmt.Printf("tagid:%d, tgname:%s, tagtype:%s, key:%s, bgtime:%s, edtime:%s\n", tagid, tgname, tagtype, key, bgtime, edtime)
	srd := new(SysRealData)
	value, ctnue, err := srd.GetSysRealTimeDataStatisticByKey(tagid, tgname, tagtype, key, bgtime, edtime)
	return fmt.Sprint(value), ctnue, err
}

/*******************************************************************************
功能:解析Kpi函数
输入:脚本字符串
输出:[string] 词元字符串(该词元字符串仅包含数字和运算符)
	[bool] 本周期没有数据,但本周期之后有数据时,为true,否则为false
	[error] 错误信息
编辑:wang_jp
时间:2020年05月03日
*******************************************************************************/
func (s *Script) extractKpi(tok string) (string, bool, error) {
	tagid := s.MainTagId
	tgname := ""
	//如果包含tag,则提取tagid/tagname

	tag, key, bgtime_token, edtime_token, err := s.getKpiParameter(tok) //提取函数的参数
	if err != nil {                                                     //提取参数错误
		return tok, false, fmt.Errorf("提取函数 [kpi] 参数时发生错误:[%s]", err.Error())
	}
	id, err := strconv.ParseInt(tag, 0, 64) //字符串转换成数字
	if err != nil {                         //转换成数字失败
		tgname = tag //参数是tagname
		tagid = 0
	} else { //参数是tagid
		tagid = id
	}
	bgtime, err := s.extractFuncTimeParamater(bgtime_token, s.BeginTime)
	if err != nil { //解析时间参数错误
		return tok, false, fmt.Errorf("解析时间参数时发生错误,extractFuncTimeParamater(%q,%q) 错误信息:[%s]", bgtime_token, s.BeginTime, err.Error())
	}
	edtime, err := s.extractFuncTimeParamater(edtime_token, s.EndTime)
	if err != nil { //解析时间参数错误
		return tok, false, fmt.Errorf("解析时间参数时发生错误,extractFuncTimeParamater(%q,%q) 错误信息:[%s]", edtime_token, s.EndTime, err.Error())
	}

	var kpicfg CalcKpiConfigList
	kpicfg.Id = tagid
	kpicfg.KpiTag = tgname
	value, ctnue, err := kpicfg.GetKpiDataStatisticByKey(key, bgtime, edtime)
	return fmt.Sprint(value), ctnue, err
}

/*******************************************************************************
功能:解析SQL语句
输入:脚本字符串
输出:[string] 词元字符串(该词元字符串仅包含数字)
	[bool] 本周期没有数据的时候但本周期之后有数据时,为true,否则false
	[error] 错误信息
编辑:wang_jp
时间:2020年05月03日
*******************************************************************************/
func (s *Script) extractSql(token string) (interface{}, bool, error) {
	toks := strings.Split(token, ").") //分割语句
	sql_select := "SELECT "
	sql_from := "FROM "
	sql_where := "WHERE "
	sql_groupby := "GROUP BY "
	sql_orderby := "ORDER BY "
	sql_limit := "LIMIT "
	astype := "value"
	var timecolumn, timefilter, sql_where2 string

	for i, tok := range toks {
		if i == len(toks)-1 { //如果是最后一个,去掉最后的字符
			tok = tok[:len(tok)-1]
		}
		if strings.Contains(strings.ToLower(tok), "select(") {
			sql_select += tok[len("select("):]
		}
		if strings.Contains(strings.ToLower(tok), "from(") {
			sql_from += tok[len("from("):]
		}
		if strings.Contains(strings.ToLower(tok), "where(") {
			sql_where += tok[len("where("):]
		}
		if strings.Contains(strings.ToLower(tok), "groupby(") {
			sql_groupby += tok[len("groupby("):]
		}
		if strings.Contains(strings.ToLower(tok), "orderby(") {
			sql_orderby += tok[len("orderby("):]
		}
		if strings.Contains(strings.ToLower(tok), "limit(") {
			sql_limit += tok[len("limit("):]
		}
		if strings.Contains(strings.ToLower(tok), "timecolumn(") {
			timecolumn += tok[len("timecolumn("):]
		}
		if strings.Contains(strings.ToLower(tok), "timefilter(") {
			timefilter += tok[len("timefilter("):]
		}
		if strings.Contains(strings.ToLower(tok), "as(") {
			astype = tok[len("as("):]
		}
	}

	if len(timefilter) > 0 && len(timecolumn) == 0 { //设置了过滤条件但没有设置时间列
		return token, false, fmt.Errorf("自定义SQL脚本[%s]中配置了时间过滤器但没有设置时间字段!", token)
	}
	if len(timefilter) > 0 { //设置了时间列和过滤条件
		reg := regexp.MustCompile(_funcTime)
		times := reg.FindAllString(timefilter, -1)
		if len(times) < 2 {
			return token, false, fmt.Errorf("自定义SQL脚本[%s]中时间过滤器参数设置错误:[%s]!", token, timefilter)
		}
		bgtime, err := s.extractFuncTimeParamater(times[0], s.BeginTime)
		if err != nil { //解析时间参数错误
			return token, false, fmt.Errorf("自定义SQL脚本[%s]中计算开始时间参数错误,错误信息:[%s]", token, err.Error())
		}
		edtime, err := s.extractFuncTimeParamater(times[1], s.EndTime)
		if err != nil { //解析时间参数错误
			return token, false, fmt.Errorf("自定义SQL脚本[%s]中计算结束时间参数错误,错误信息:[%s]", token, err.Error())
		}
		if len(sql_where) > len("WHERE ") { //已经有了其他WHERE条件
			sql_where += " AND"
		}
		sql_where2 = sql_where + fmt.Sprintf(" %s > '%s'", timecolumn, edtime)
		sql_where += fmt.Sprintf(" %s > '%s' AND %s <= '%s'", timecolumn, bgtime, timecolumn, edtime) //时间过滤
	}
	if len(timecolumn) > 0 && len(timefilter) == 0 { //设置了时间列,但没有设置过滤条件
		if len(sql_where) > len("WHERE ") { //已经有了其他WHERE条件
			sql_where += " AND"
		}
		sql_where2 = sql_where + fmt.Sprintf(" %s > '%s'", timecolumn, s.EndTime)
		sql_where += fmt.Sprintf(" %s > '%s' AND %s <= '%s'", timecolumn, s.BeginTime, timecolumn, s.EndTime) //时间过滤
	}

	sql := sql_select + " " + sql_from
	if len(sql_where) > len("WHERE ") {
		sql += (" " + sql_where)
	}
	if len(sql_groupby) > len("GROUP BY ") {
		sql += (" " + sql_groupby)
	}
	if len(sql_orderby) > len("ORDER BY ") {
		sql += (" " + sql_orderby)
	}
	if len(sql_limit) > len("LIMIT ") {
		sql += (" " + sql_limit)
	}

	cntinu := false
	o := orm.NewOrm() //新建orm对象
	o.Using("default")
	var val []orm.Params
	num, err := o.Raw(sql).Values(&val)
	if err == nil && num == 0 { //现有查询时间范围下没有查询到数据
		err = fmt.Errorf("SQL [%s] 没有查询到数据!", sql)
		if len(sql_where2) > 0 { //有条件2
			sql := sql_select + " " + sql_from
			if len(sql_where2) > len("WHERE ") {
				sql += (" " + sql_where2)
			}
			sql += " LIMIT 1"
			var vp []orm.Params
			num, e := o.Raw(sql).Values(&vp)
			if e == nil && num > 0 { //现有查询时间范围之后有数据
				cntinu = true
			}
		}
	}
	if err != nil {
		return nil, cntinu, fmt.Errorf("SQL Script:[%s], SQL:[%s], Error message:[%s]",
			token, sql, err.Error())
	}

	if s.OnlySQL == false {
		astype = "value"
	}
	switch strings.ToLower(astype) {
	case "map":
		return val, cntinu, err
	case "string":
		str := "0.0"
		if num >= 1 {
			va := val[0]
			for _, v := range va { //随机找到第一个能转换为浮点数的
				str = fmt.Sprint(v)
				break
			}
		}
		return str, cntinu, err
	case "json":
		b, err := json.Marshal(val)
		return string(b), cntinu, err
	case "sql":
		return sql, cntinu, err
	default: //value
		str := "0.0"
		if num >= 1 {
			va := val[0]
			for _, v := range va { //随机找到第一个能转换为浮点数的
				value := fmt.Sprint(v)
				_, e := strconv.ParseFloat(value, 64) //能转换为浮点数
				if e == nil {
					str = value
					break
				}
			}
		}
		return str, cntinu, err
	}
}

/*******************************************************************************
功能:解析DateTime函数
输入:脚本字符串
输出:[interface{}] 词元字符串(该词元字符串仅包含数字或时间)
	[error] 错误信息
编辑:wang_jp
时间:2020年05月07日
*******************************************************************************/
func (s *Script) extractDateTime() (interface{}, error) {
	tok := s.ScriptStr
	reg := regexp.MustCompile(`\([[:print:]]+\)`) //提取括号中的参数(含括号)
	strs := reg.FindAllString(tok, -1)
	if len(strs) > 0 {
		//dtoks := strings.Split(strs[0][1:len(strs[0])-1], ",") //用逗号分隔
		paras := strs[0][1 : len(strs[0])-1]       //不带括号的参数
		paras = strings.ReplaceAll(paras, " ", "") //删除空格

		if len(paras) < 3 {
			return tok, fmt.Errorf("脚本[%s]中没有提取到有效的参数", tok)
		}
		parat := paras[0 : len(paras)-2]     //时间参数
		paraSelector := paras[len(paras)-1:] //时间选择器

		dt, err := s.extractFuncTimeParamater(parat, s.BeginTime)
		if err != nil { //解析时间参数错误
			return tok, fmt.Errorf("解析时间参数时发生错误,extractFuncTimeParamater(%q,%q) 错误信息:[%s]", parat, s.EndTime, err.Error())
		}
		tm, err := TimeParse(dt)
		if err == nil {
			switch paraSelector {
			case "T":
				return tm.Format("2006/01/02 15:04:05"), nil
				break
			case "t":
				return tm.Format("15:04:05"), nil
				break
			case "Y", "y":
				return tm.Format("2006"), nil
				break
			case "M":
				return tm.Format("2006/01"), nil
				break
			case "D":
				return tm.Format("2006/01/02"), nil
				break
			case "m":
				return tm.Format("1"), nil
				break
			case "d":
				return tm.Format("2"), nil
				break
			case "h":
				return tm.Format("3"), nil
				break
			case "H":
				return tm.Format("15"), nil
				break
			case "f", "F":
				return tm.Format("4"), nil
				break
			case "s", "S":
				return tm.Format("5"), nil
				break
			case "w":
				return tm.Weekday(), nil
				break
			case "W":
				return CnWeekday(tm), nil
				break
			default:
				return dt, nil
				break
			}
		} else {
			return tok, err
		}

	} else {
		return tok, fmt.Errorf("脚本[%s]中没有提取到括号及括号中的内容", tok)
	}
	return tok, fmt.Errorf("脚本[%s]中没有提取到括号及括号中的内容", tok)
}

/*******************************************************************************
功能:解析函数内的时间参数
输入:脚本字符串,tok字符串默认的时间
输出:时间字符串和错误信息
编辑:wang_jp
时间:2020年01月8日
*******************************************************************************/
func (s *Script) extractFuncTimeParamater(tok, defaulttime string) (string, error) {
	var tstr string
	switch tok {
	case "0": //自动时间
		tstr = defaulttime
	case "now": //当前时间(需要减去延迟时间)
		tstr = time.Now().Add(time.Duration(EngineCfgMsg.CfgMsg.SerialCalcDelaySec) * -1 * time.Second).Format(EngineCfgMsg.Sys.TimeFormat)
	case "bgoy": //beginning of the year
		return ComposeTime(s.BaseTime, defaulttime, s.ShiftHour, 0)
	case "bgos": //beginning of the sesion
		return ComposeTime(s.BaseTime, defaulttime, s.ShiftHour, 1)
	case "bgom": //beginning of the month
		return ComposeTime(s.BaseTime, defaulttime, s.ShiftHour, 2)
	case "bgod": //beginning of the day
		return ComposeTime(s.BaseTime, defaulttime, s.ShiftHour, 4)
	default: //
		rg_tfc := regexp.MustCompile(_timeDev) //匹配时间函数参数
		if rg_tfc.MatchString(tok) {           //匹配函数
			return s.funcParaTimeAdd(tok, defaulttime)
		} else {
			rg := regexp.MustCompile(_timeStr) //单纯时间串
			if rg.MatchString(tok) {           //匹配单词时间串
				tstr = tok
			} else { //脚本引用数据
				scr := s
				scr.ScriptStr = tok
				t, _, err := scr.analysisScript(false)
				if err != nil {
					return tok, err
				}
				_, err = statistic.TimeParse(t) //校验时间格式
				if err != nil {
					return tok, err
				}
				tstr = t
			}
		}
	}
	return tstr, nil
}

/*******************************************************************************
功能:解析函数内的时间参数函数
输入:脚本字符串,tok字符串默认的时间
输出:时间字符串和错误信息
编辑:wang_jp
时间:2020年2月23日
*******************************************************************************/
func (s *Script) funcParaTimeAdd(token, defaulttime string) (string, error) {
	l := len(token)
	if l >= 10 {
		paramter_str := token[6 : l-1] //去掉头[ t_add( ]和尾[ ) ]
		paramters := strings.Split(paramter_str, ",")
		if len(paramters) == 2 {
			//解析时间参数
			basetimestr, err := s.extractFuncTimeParamater(paramters[0], defaulttime)
			if err != nil {
				return "", err
			}
			//基准时间格式化
			basetime, err := TimeParse(basetimestr)
			if err != nil {
				//return "", fmt.Errorf("Time format error in [%s],error message is [%s]", token, err.Error())
				return "", fmt.Errorf("时间格式错误 [%s],错误信息: [%s]", token, err.Error())
			}
			//时间偏移量格式化
			timeoffset, err := time.ParseDuration(paramters[1])
			if err != nil {
				//return "", fmt.Errorf("Time duration error in [%s],error message is [%s]", token, err.Error())
				return "", fmt.Errorf("时间范围错误 [%s],错误信息: [%s]", token, err.Error())
			}
			timeresult := basetime.Add(timeoffset)
			return timeresult.Format(EngineCfgMsg.Sys.TimeFormat), nil
		} else {
			//return "", fmt.Errorf("Need 2 parameters in call to t_add() in [%s]", token)
			return "", fmt.Errorf("在脚本程序 [%s] 中调用 t_add()函数需要两个参数", token)
		}
	} else {
		//return "", fmt.Errorf("Script error in [%s]", token)
		return "", fmt.Errorf("脚本程序错误:[%s]", token)
	}
}

/*******************************************************************************
功能:获取fc函数的参数
输入:包含fc()函数及其参数(多参数)的字符串
输出:key,beginTime_token,endTime_token
编辑:wang_jp
时间:2020年01月8日
*******************************************************************************/
func (s *Script) getFcParameter(token string) (string, string, string, error) {
	var key, bgtime, edtime string
	reg := regexp.MustCompile(`fc\([[:print:]]+\)`)
	strs := reg.FindAllString(token, -1)
	if len(strs) > 0 {
		t := strs[0]
		l := len(t) //长度
		if l > 4 {
			paramter_str := t[3 : l-1]
			reg = regexp.MustCompile(fmt.Sprintf(`%s|%s`, _funcTime, `([\w\.\:\-]+)`))
			paramters := reg.FindAllString(paramter_str, -1)
			//paramters := strings.Split(paramter_str, ",")
			if len(paramters) == 3 {
				key = paramters[0]
				bgtime = paramters[1]
				edtime = paramters[2]
			} else {
				if len(paramters) > 3 {
					//return "", "", "", fmt.Errorf("too many parameters in call to fc() in [%s]", token)
					return "", "", "", fmt.Errorf("调用函数 fc() 时输入的参数太多: [%s]", token)
				} else {
					//return "", "", "", fmt.Errorf("not enough parameters in call to fc() in [%s]", token)
					return "", "", "", fmt.Errorf("调用函数 fc() 时输入的参数太少: [%s]", token)
				}
			}
		} else {
			//return "", "", "", fmt.Errorf("fc() paremater error in [%s]", token)
			return "", "", "", fmt.Errorf("函数 fc() 参数错误:[%s]", token)
		}
	} else {
		//return "", "", "", fmt.Errorf("Script error in [%s]", token)
		return "", "", "", fmt.Errorf("脚本程序错误:[%s]", token)
	}
	return key, bgtime, edtime, nil
}

/*******************************************************************************
功能:获取Srtd函数的参数
输入:包含srtd()函数及其参数的字符串,默认的开始时间、结束时间
输出:type,key,beginTime_token,endTime_token
编辑:wang_jp
时间:2020年01月8日
*******************************************************************************/
func (s *Script) getSrtdParameter(token string) (string, string, string, string, error) {
	var tagtype, key, bgtime, edtime string
	reg := regexp.MustCompile(`srtd\([[:print:]]+\)`)
	strs := reg.FindAllString(token, -1)
	if len(strs) > 0 {
		t := strs[0]
		l := len(t) //长度
		if l > 7 {  //srtd(,)
			paramter_str := t[5 : l-1] //取括号中的参数
			reg = regexp.MustCompile(fmt.Sprintf(`%s|%s`, _funcTime, `([\w\.\:\-]+)`))
			paramters := reg.FindAllString(paramter_str, -1)
			switch len(paramters) {
			case 2: //有两个参数的情况
				{
					tagtype = paramters[0]
					key = paramters[1]
					bgtime = s.BeginTime //使用默认参数
					edtime = s.EndTime   //使用默认参数
				}
			case 4: //有四个参数的情况
				{
					tagtype = paramters[0]
					key = paramters[1]
					bgtime = paramters[2]
					edtime = paramters[3]
				}
			default:
				{
					//return "", "", "", "", fmt.Errorf("Need 2 or 4 parameters in call to srtd() in [%s]", token)
					return "", "", "", "", fmt.Errorf("调用函数 srtd() 时需要2个或者4个参数:[%s]", token)
				}
			}
		} else {
			//return "", "", "", "", fmt.Errorf("srtd() paremater error in [%s]", token)
			return "", "", "", "", fmt.Errorf("调用函数 srtd() 时参数错误:[%s]", token)
		}
	} else {
		//return "", "", "", "", fmt.Errorf("Script error in [%s]", token)
		return "", "", "", "", fmt.Errorf("脚本程序错误:[%s]", token)
	}
	return tagtype, key, bgtime, edtime, nil
}

/*******************************************************************************
功能:获取Kpi函数的参数
输入:包含kpi()函数及其参数的字符串,默认的开始时间、结束时间
输出:tag(tag id或者tag name),key,beginTime_token,endTime_token
编辑:wang_jp
时间:2020年05月4日
*******************************************************************************/
func (s *Script) getKpiParameter(token string) (string, string, string, string, error) {
	var tag, key, bgtime, edtime string
	reg := regexp.MustCompile(`kpi\([[:print:]]+\)`)
	strs := reg.FindAllString(token, -1)
	if len(strs) > 0 {
		t := strs[0]
		l := len(t) //长度
		if l > 6 {  //kpi(,)
			paramter_str := t[4 : l-1] //取括号中的参数
			reg = regexp.MustCompile(fmt.Sprintf(`%s|%s`, _funcTime, `([\w\.\:\-]+)`))
			paramters := reg.FindAllString(paramter_str, -1)
			switch len(paramters) {
			case 2: //有两个参数的情况
				{
					tag = paramters[0]
					key = paramters[1]
					bgtime = s.BeginTime //使用默认参数
					edtime = s.EndTime   //使用默认参数
				}
			case 4: //有四个参数的情况
				{
					tag = paramters[0]
					key = paramters[1]
					bgtime = paramters[2]
					edtime = paramters[3]
				}
			default:
				{
					//return "", "", "", "", fmt.Errorf("Need 2 or 4 parameters in call to kpi() in [%s]", token)
					return "", "", "", "", fmt.Errorf("调用函数 kpi() 时需要2个或者4个参数,实际提取到[%d]个:[%s]", len(paramters), paramters)
				}
			}
		} else {
			//return "", "", "", "", fmt.Errorf("kpi() paremater error in [%s]", token)
			return "", "", "", "", fmt.Errorf("调用函数 kpi() 时参数错误:[%s]", token)
		}
	} else {
		//return "", "", "", "", fmt.Errorf("Script error in [%s]", token)
		return "", "", "", "", fmt.Errorf("脚本程序错误:[%s]", token)
	}
	return tag, key, bgtime, edtime, nil
}

func (s *Script) Compile() []error {
	defer func() {
		if err := recover(); err != nil {
			logs.Critical("Script.Compile 中遇到错误:%s;[%#v]", s.ScriptStr, err)
		}
	}()
	scrStr := s.ScriptStr
	pattern := _date + "|" + s.getRegexpString() //全部脚本词元组
	reg := regexp.MustCompile(pattern)
	toks := reg.FindAllString(scrStr, -1)
	index := reg.FindAllStringIndex(scrStr, -1)

	s.OnlySQL = false
	if len(toks) == 1 { //只有一个脚本词元
		tk := strings.ToLower(toks[0])
		if strings.Contains(tk, "select(") { //包含SELECT
			rg := regexp.MustCompile(`as\((string|map|json|sql)\)`)
			if rg.MatchString(tk) {
				s.OnlySQL = true
			}
		}
	}

	strl := len(scrStr) //总长度
	scrArray := []byte(scrStr)
	var errs []error
	var lastistok, thisistok bool //是否函数词元组
	pattern = fmt.Sprintf("%s|%s|%s|%s|%s|%s|%s|%s",
		_sql,
		_funcReg,
		_srtd,
		_kpi,
		_date,
		_Field,
		_gdsum,
		_number)
	reg = regexp.MustCompile(pattern) //除了数学符号之外的词元组

	for i, id := range index {
		thisistok = reg.MatchString(toks[i])
		if i > 0 {
			if lastistok == true && thisistok == true {
				errs = append(errs, fmt.Errorf("[%d:%d]:两个操作数不可直接相连:[%s%s]", index[i-1][1], id[0]+1, toks[i-1], toks[i]))
			}
		}
		lastistok = thisistok

		switch i {
		case 0: //第一个匹配
			if id[0] != 0 {
				b := scrArray[0:id[0]]
				if isSpace(b) == false {
					errs = append(errs, fmt.Errorf("[1:%d]:无效字符:[%s]", id[0], string(b)))
				}
			}
		case len(index) - 1: //最后一个匹配
			if id[0] > index[i-1][1]+1 {
				b := scrArray[index[i-1][1]:id[0]]
				if isSpace(b) == false {
					errs = append(errs, fmt.Errorf("[%d:%d]:无效字符:[%s]", index[i-1][1]+1, id[0], string(b)))
				}
			}
			if id[1] < strl {
				b := scrArray[id[1]:strl]
				if isSpace(b) == false {
					errs = append(errs, fmt.Errorf("[%d:%d]:无效字符:[%s]", id[1], strl, string(b)))
				}
			}
		default: //其他匹配
			if id[0] > index[i-1][1]+1 {
				b := scrArray[index[i-1][1]:id[0]]
				if isSpace(b) == false {
					errs = append(errs, fmt.Errorf("[%d:%d]:无效字符:[%s]", index[i-1][1]+1, id[0], string(b)))
				}
			}
		}
	}
	return errs
}

func isSpace(b []byte) bool {
	reg := regexp.MustCompile(`[^\s]+`)
	return !reg.Match(b)
}
