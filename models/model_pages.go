package models

import (
	"crypto/md5"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/astaxie/beego/orm"
	"github.com/bkzy-wangjp/MicEngine/MicScript/numgo"
	"github.com/bkzy-wangjp/MicEngine/MicScript/regression"
	_ "github.com/go-sql-driver/mysql"
)

/*************************************************
功能:对字符串进行MD5加密
输入:待加密字符串
输出:加密后的字符串
说明:
编辑:wang_jp
时间:2020年3月12日
*************************************************/
func Md5str(str string) string {
	data := []byte(str)
	has := md5.Sum(data)
	md5str := fmt.Sprintf("%x", has)
	return md5str
}

/*************************************************
功能:用户登录
输入:用户名和密码
输出:用户信息和错误信息
说明:
编辑:wang_jp
时间:2020年3月12日
*************************************************/
func (v *SysUser) UserLogIn(psw string, username ...string) (*SysUser, error) {
	o := orm.NewOrm()
	o.Using("default") //选定数据库
	if len(username) > 0 {
		v.Username = username[0]
	}
	uname := v.Username
	ps := Md5str(psw)   //明文密码先加密
	if len(psw) == 32 { //如果密码长度为32位(已经MD5的)
		ps = psw //使用原密码
	}
	if n, err := o.QueryTable("SysUser").Filter("Username", v.Username).Filter("Password", ps).All(v); err == nil {
		if n == 0 {
			return nil, fmt.Errorf("The user [%s] does not exist or the password is wrong.[用户'%s'不存在或者密码错误]", uname, uname)
		} else {
			return v, nil
		}
	} else {
		return nil, err
	}
}

/*************************************************
功能:用户更新密码
输入:用户名、旧密码和新密码
输出:执行结果是否正确和错误信息
说明:
编辑:wang_jp
时间:2020年3月28日
*************************************************/
func (v *SysUser) UserUpdatePswd(oldpswd, newpswd string, username ...string) (bool, error) {
	o := orm.NewOrm()
	o.Using("default") //选定数据库
	if len(username) > 0 {
		v.Username = username[0]
	}
	//用旧密码读取用户信息
	if n, err := o.QueryTable("SysUser").Filter("Username", v.Username).Filter("Password", oldpswd).All(v); err == nil {
		if n == 0 { //没有读取到用户信息，返回错误
			return false, fmt.Errorf("The old password is wrong.[旧密码错误]")
		} else {
			v.Password = newpswd              //设置密码为新密码
			r, err := o.Update(v, "Password") //更新数据库
			if err != nil {                   //如果返回错误信息
				return false, err //返回错误信息
			}
			if r > 0 { //如果受影响的行数大于0
				return true, nil //返回修改成功
			} else { //否则
				return false, nil //返回修改失败
			}
		}
	} else { //读取信息过程错误，范围错误
		return false, err
	}
}

/*************************************************
功能:重置用户密码
输入:用户名、新密码
输出:执行结果是否正确和错误信息
说明:只有有管理员权限的人员才可以重置用户密码，在调用此函数前
	要先验证权限，本函数不进行验证
编辑:wang_jp
时间:2020年4月1日
*************************************************/
func (v *SysUser) UserResetPswd(newpswd string, username ...string) (bool, error) {
	o := orm.NewOrm()
	o.Using("default") //选定数据库
	if len(username) > 0 {
		v.Username = username[0]
	}

	//用旧密码读取用户信息
	if n, err := o.QueryTable("SysUser").Filter("Username", v.Username).All(v); err == nil {
		if n == 0 { //没有读取到用户信息，返回错误
			return false, fmt.Errorf("The user [%s] does not exist.[用户[%s]不存在]", v.Username, v.Username)
		} else {
			v.Password = newpswd              //设置密码为新密码
			r, err := o.Update(v, "Password") //更新数据库
			if err != nil {                   //如果返回错误信息
				return false, err //返回错误信息
			}
			if r > 0 { //如果受影响的行数大于0
				return true, nil //返回修改成功
			} else { //否则
				return false, nil //返回修改失败
			}
		}
	} else { //读取信息过程错误，范围错误
		return false, err
	}
}

/*************************************************
功能:通过用户Id获取菜单
输入:用户Id
输出:菜单信息
说明:
编辑:wang_jp
时间:2020年3月12日
*************************************************/
func (u *SysUser) GetRolesByUserId(userid ...int64) ([]*SysRole, error) {
	o := orm.NewOrm()
	o.Using("default") //选定数据库
	if len(userid) > 0 {
		u.Id = userid[0]
	}
	var v []*SysRole
	if _, err := o.QueryTable("SysRole").Filter("Users__User__Id", u.Id).All(&v); err == nil {
		return v, nil
	} else {
		return nil, err
	}
}

/*************************************************
功能:获取用户列表
输入:可选的用户名,可以没有,也可以多个
输出:用户列表
说明:如果输入了用户名,输入字符串模糊查询,如果没有输入，查询列出所有用户
	有用户名输入时,以UserName Like usernames%的方式模糊查询
编辑:wang_jp
时间:2020年4月1日
*************************************************/
func (su *SysUser) GetUserList(usernames ...string) ([]SysUser, error) {
	o := orm.NewOrm()
	o.Using("default") //选定数据库
	var v []SysUser
	if len(usernames) > 0 { //如果输入了用户名
		for _, user := range usernames { //遍历输入的用户名
			var u []SysUser
			if _, err := o.QueryTable("SysUser").Filter("UserName__startswith", user).All(&u); err == nil { //like username% 模糊查询
				v = append(v, u...)
			}
		}
		return v, nil
	} else {
		if _, err := o.QueryTable("SysUser").Filter("Id__gt", 0).All(&v); err == nil {
			return v, nil
		} else {
			return nil, err
		}
	}
}

/*************************************************
功能:通过用户Id获取菜单
输入:用户Id
输出:菜单信息
说明:
编辑:wang_jp
时间:2020年3月12日
*************************************************/
func (u *SysUser) GetMenusByUserId(userid ...int64) ([]*SysMenu, error) {
	o := orm.NewOrm()
	o.Using("default") //选定数据库
	if len(userid) > 0 {
		u.Id = userid[0]
	}
	var v []*SysMenu
	if _, err := o.QueryTable("SysMenu").Filter("Permissions__Permission__Roles__Role__Users__User__Id", u.Id).Filter("MenuType", 2).Filter("Pid__gt", 0).Distinct().OrderBy("Level", "Seq").All(&v); err == nil {
		return v, nil
	} else {
		return nil, err
	}
}

/*************************************************
功能:通过用户Id获取结构树(不带变量tag叶子节点)
输入:用户Id,是否需要作业层级之下的子节点
输出:结构树信息
说明:
编辑:wang_jp
时间:2020年3月12日
*************************************************/
func (u *SysUser) GetTreeNodesByUserId(needsub bool, userid ...int64) ([]*MsMiddleCorrelation, error) {
	o := orm.NewOrm()
	o.Using("default") //选定数据库
	if len(userid) > 0 {
		u.Id = userid[0]
	}
	var nodes []*MsMiddleCorrelation
	qt := o.QueryTable("MsMiddleCorrelation")
	if _, err := qt.Filter("Permissions__Permission__Roles__Role__Users__User__Id", u.Id).Filter("LevelCategory__Id__lte", 4).Distinct().OrderBy("TreeLevelCode", "Seq").All(&nodes); err == nil {
		if needsub == true {
			for _, node := range nodes { //遍历所有节点，寻找作业层级的子节点
				if node.LevelCategory.Id == 4 { //是作业层级
					md := new(MsMiddleCorrelation)
					if subnodes, err := md.GetSubNodesByStageLevelCode(node.TreeLevelCode); err == nil { //通过TreeLevelCode获取子节点
						for _, subnode := range subnodes { //将子节点追加到主节点传下
							nodes = append(nodes, subnode)
						}
					}
				}
			}
		}
		return nodes, nil
	} else {
		return nil, err
	}
}

/*************************************************
功能:通过用户Id获取实时监控结构树
输入:用户Id
输出:结构树信息
说明:
编辑:wang_jp
时间:2020年3月12日
*************************************************/
func (u *SysUser) GetMonitorNodesByUserId(userid ...int64) ([]*RealTimeMonitor, error) {
	o := orm.NewOrm()
	o.Using("default") //选定数据库
	if len(userid) > 0 {
		u.Id = userid[0]
	}
	var nodes []*RealTimeMonitor
	qt := o.QueryTable("RealTimeMonitor")
	if _, err := qt.Filter("Permissions__Permission__Roles__Role__Users__User__Id", u.Id).Distinct().OrderBy("Id", "Seq").All(&nodes); err == nil {
		return nodes, nil
	} else {
		return nil, err
	}
}

/*************************************************
功能:通过用户名从缓存表获取结构树(带有变量tag作为树的叶子节点)
输入:用户名
输出:结构树字符串,根节点id,是否成功
说明:
编辑:wang_jp
时间:2020年3月12日
*************************************************/
func (md *MsMiddleCorrelationString) GetTreeNodesFromeCache(idcode string) (string, int64, bool) {
	o := orm.NewOrm()
	o.Using("default") //选定数据库
	var nodes []*MsMiddleCorrelationString
	qt := o.QueryTable("MsMiddleCorrelationString")
	if n, err := qt.Filter("IdCode", idcode).All(&nodes); err == nil {
		if n > 0 {
			return nodes[0].MenuValue, nodes[0].RootPid, true
		} else {
			return "", 0, false
		}
	} else {
		return "", 0, false
	}
}

/*************************************************
功能:通过用户名更新缓存表
输入:用户名,树结构字符串,根节点id
输出:更新数据行数或者插入的数据的ID号,错误信息
说明:
编辑:wang_jp
时间:2020年3月20日
*************************************************/
func (m *MsMiddleCorrelationString) InsertOrUpdateTreeNodesInCache(idcode, menustring string, rootpid int64) (int64, error) {
	o := orm.NewOrm()
	o.Using("default") //选定数据库
	var nodes MsMiddleCorrelationString
	nodes.IdCode = idcode
	if o.Read(&nodes, "IdCode") == nil { //检查是否存在
		nodes.EditTime = time.Now().Format(EngineCfgMsg.Sys.TimeFormat)
		nodes.MenuValue = menustring
		nodes.RootPid = rootpid
		return o.Update(&nodes) //存在就更新
	} else {
		nodes.IdCode = idcode
		nodes.EditTime = time.Now().Format(EngineCfgMsg.Sys.TimeFormat)
		nodes.MenuValue = menustring
		nodes.RootPid = rootpid
		return o.Insert(&nodes)
	}
}

/*************************************************
功能:通过用户Id获取结构树(带有变量tag作为树的叶子节点)
输入:用户Id
输出:结构树信息
说明:
编辑:wang_jp
时间:2020年3月12日
*************************************************/
func (u *SysUser) GetTreeNodesWithTagsByUserId(userid ...int64) ([]*MsMiddleCorrelation, error) {
	o := orm.NewOrm()
	o.Using("default") //选定数据库
	if len(userid) > 0 {
		u.Id = userid[0]
	}
	var nodes []*MsMiddleCorrelation
	var maxnodeid int64 //最大节点id
	terminalCnt := 0    //终端节点数量
	qt := o.QueryTable("MsMiddleCorrelation")
	if _, err := qt.Filter("Permissions__Permission__Roles__Role__Users__User__Id", u.Id).Filter("LevelCategory__Id__lte", 4).Distinct().OrderBy("TreeLevelCode", "Seq").All(&nodes); err == nil {
		for _, node := range nodes { //遍历所有节点，寻找作业层级的子节点
			if node.LevelCategory.Id == 4 { //是作业层级
				if subnodes, err := node.GetSubNodesByStageLevelCode(node.TreeLevelCode); err == nil { //通过TreeLevelCode获取子节点
					for _, subnode := range subnodes { //将子节点追加到主节点传下
						nodes = append(nodes, subnode)
						if subnode.Id > maxnodeid {
							maxnodeid = subnode.Id
						}
						if subnode.LevelCategory.Id > 4 && subnode.LevelCategory.Id < 9 {
							terminalCnt += 1
						}
					}
				}
			}
			if node.Id > maxnodeid {
				maxnodeid = node.Id
			}
		}

		goCnt := 20
		var tags []*MsMiddleCorrelation
		loop := len(nodes) / goCnt

		for i := 0; i <= loop; i++ {
			wait := &sync.WaitGroup{}
			leafs := make(chan []*MsMiddleCorrelation, goCnt) //结果集chan
			var nds []*MsMiddleCorrelation
			if i == loop {
				nds = nodes[i*goCnt:]
			} else {
				nds = nodes[i*goCnt : (i+1)*goCnt]
			}
			cnt := 0
			for _, node := range nds { //遍历所有节点，寻找设备及其下层级各节点的直系变量节点
				var ndtype int64 = 0
				switch node.LevelCategory.Id {
				case 5: //设备
					ndtype = 4
				case 6: //仪表
					ndtype = 2
				case 7: //电机
					ndtype = 1
				case 8: //分析仪器
					ndtype = 3
				case 9: //取样点
					ndtype = 5
				case 22: //电能计量表
					ndtype = 6
				case 20, 21: //能源管理系统,供电线路
					ndtype = 7
				default: //其他
				}
				if ndtype > 0 {
					wait.Add(1)
					cnt += 1
					md := new(MsMiddleCorrelation)
					go md.GetDirectChildTagsOfNode(node.Id, ndtype, node.ItemIdInTable, maxnodeid, wait, leafs)
				}
			}
			for j := 0; j < cnt; j++ {
				tags = append(tags, <-leafs...)
			}
			wait.Wait()
			close(leafs) //关闭chan道
		}
		nodes = append(nodes, tags...)
		return nodes, nil
	} else {
		return nil, err
	}
}

/*************************************************
功能:通过节点类型和节点ID索引该节点下的直系Tag
输入:nodetype:节点类型,1=电机,2=仪表,3=分析仪,4=设备,5=取样点,
	6=电能计量表,7=供配电设备
	nodeid:节点在其所属类型表中的ID号
输出:Taglists
说明:节点类型随ResourceType种类的增减而增减
编辑:wang_jp
时间:2020年3月19日
*************************************************/
func (md *MsMiddleCorrelation) GetDirectChildTagsOfNode(pid, nodetype, nodeid, maxnodeid int64, wait *sync.WaitGroup, leafs chan []*MsMiddleCorrelation) {
	defer wait.Done()
	o := orm.NewOrm()
	o.Using("default") //选定数据库
	var nodes []*MsMiddleCorrelation
	var tags []*OreProcessDTaglist
	qt := o.QueryTable("OreProcessDTaglist")
	if _, err := qt.Filter("ResourceType", nodetype).Filter("ItemIDinTable", nodeid).OrderBy("Seq").All(&tags); err == nil {
		for i, tag := range tags {
			if fullname, err := tag.GetTagFullNameByID(); err == nil {
				tags[i].TagFullName = fullname
			}
		}
		lct := MineTableList{9999, "", 0, nil}
		for _, tag := range tags { //遍历找到的所有直接子节点
			leaf := new(MsMiddleCorrelation)
			leaf.Id = maxnodeid + tag.Id //叶子节点的id为虚拟id
			leaf.LevelCategory = &lct    //变量类型的category
			leaf.Pid = pid
			leaf.LevelName = tag.TagDescription
			leaf.ItemIdInTable = tag.Id //tag的id存储在ItemIdInTable中
			leaf.TreeLevelCode = tag.TagFullName
			leaf.ConstrutionCode = tag.DecimalNum
			leaf.ConstrutionTableCode = tag.Unit //tag.TagUnit
			switch strings.ToLower(tag.TagType) {
			case "bool":
				leaf.Seq = 1
			case "int", "int8", "int16", "int32", "int64":
				leaf.Seq = 2
			case "float", "float32", "float64", "double":
				leaf.Seq = 3
			default:
				leaf.Seq = 4
			}
			nodes = append(nodes, leaf)
		}
	}
	leafs <- nodes
}

/*************************************************
功能:通过作业层级的层级码获取该作业下的所有子节点
输入:作业层级码
输出:结构树子节点信息
说明:
编辑:wang_jp
时间:2020年3月17日
*************************************************/
func (md *MsMiddleCorrelation) GetSubNodesByStageLevelCode(levelcode string) ([]*MsMiddleCorrelation, error) {
	o := orm.NewOrm()
	o.Using("default") //选定数据库
	var nodes []*MsMiddleCorrelation
	qt := o.QueryTable("MsMiddleCorrelation")
	if _, err := qt.Filter("TreeLevelCode__istartswith", levelcode+"-").Filter("ConstrutionCode__gt", 0).Distinct().OrderBy("TreeLevelCode", "Seq").All(&nodes); err == nil {
		return nodes, nil
	} else {
		return nil, err
	}
}

/*************************************************
功能:更新引擎配置信息
输入:待更新的字段名
输出:错误信息
说明:
编辑:wang_jp
时间:2020年3月13日
*************************************************/
func (cfg *CalcKpiEngineConfig) UpdateEngineCfgInfo(fieldname string) error {
	o := orm.NewOrm()
	o.Using("default")
	EngineCfgMsg.CfgMsg.UpdateTime = TimeFormat(time.Now())
	_, err := o.Update(&EngineCfgMsg.CfgMsg, fieldname, "UpdateTime") //保存配置信息
	return err
}

/*************************************************
功能:更新用户日志
输入:待更新的日志信息
输出:错误信息
说明:
编辑:wang_jp
时间:2020年3月13日
*************************************************/
func (userlog *SysLog) InsertUserLog() error {
	o := orm.NewOrm()
	o.Using("default")
	_, err := o.Insert(userlog) //保存用户日志信息
	return err
}

/*************************************************
功能:通过节点的层级码获取该作业下的所有Tag
输入:节点层级码
输出:Taglists
说明:
编辑:wang_jp
时间:2020年3月17日
*************************************************/
func GetTaglistByNodeLevelCode(levelcode string) ([]*OreProcessDTaglist, error) {
	o := orm.NewOrm()
	o.Using("default") //选定数据库
	var tags []*OreProcessDTaglist
	qt := o.QueryTable("OreProcessDTaglist")
	if _, err := qt.Filter("TreeLevelCode__istartswith", levelcode).OrderBy("Seq").All(&tags); err == nil {
		for i, tag := range tags {
			if fullname, err := tag.GetTagFullNameByID(); err == nil {
				tags[i].TagFullName = fullname
			}
		}
		return tags, nil
	} else {
		return nil, err
	}
}

/*************************************************
功能:通过tag的dcs_id获取与之相关的实时数据库表和数据采集机的相关信息
输入:dcs_id
输出:RelevanceDcsToDbtable
说明:
编辑:wang_jp
时间:2020年3月17日
*************************************************/
func (r *RelevanceDcsToDbtable) GetTagsDataTableInfoByDcsId(dcs_id int64) ([]*RelevanceDcsToDbtable, error) {
	o := orm.NewOrm()
	o.Using("default") //选定数据库
	var rel []*RelevanceDcsToDbtable
	qt := o.QueryTable("RelevanceDcsToDbtable")
	if _, err := qt.Filter("Dcs__Id", dcs_id).RelatedSel().All(&rel); err == nil {
		return rel, nil
	} else {
		return nil, err
	}
}

/*************************************************
功能:从庚顿数据库读取数据并进行回归分析
输入:Y变量和X变量的全名以及时间参数
输出:回归分析结果
说明:
编辑:wang_jp
时间:2020年3月25日
*************************************************/
func (micgd *MicGolden) Regression(tagYname, tagXnames, begintime, endtime string, interval int64) (*regression.Regre, error) {
	//读取Y数据
	data, _, err := micgd.GoldenGetHistoryInterval(begintime, endtime, interval, tagYname)
	if err != nil {
		return nil, err
	}
	var dataY numgo.Array
	var dy []float64
	for _, dv := range data {
		for _, v := range dv {
			dy = append(dy, v.Value)
		}
	}
	dataY = dy //结果是Y数据
	ly := len(dataY)

	//格式化命令
	xnames := strings.Split(tagXnames, ",")
	hiss, _, err := micgd.GoldenGetHistoryInterval(begintime, endtime, interval, xnames...)
	if err != nil {
		return nil, err
	}

	//解析数据
	var dataXs numgo.Matrix
	dXs := make([][]float64, len(xnames))
	for tag, his := range hiss { //遍历结果集
		var data []float64
		for _, v := range his { //取出结果中的数据
			data = append(data, v.Value)
		}
		for i, xname := range xnames { //遍历x名称
			if xname == tag { //如果与结果名称相符
				dXs[i] = data //将结果存入相应位置
				break
			}
		}
	}

	for i, data := range dXs { //X数据矩阵
		if len(data) != ly { //数据长度不相等
			return nil, fmt.Errorf("因变量[%s]的数据长度[%d]与自变量[%s]的数据长度[%d]不一致,不能进行回归分析!", tagYname, ly, xnames[i], len(data))
		}
		dataXs = append(dataXs, data)
	}
	reg, err := regression.Regression(dataXs, dataY)
	return &reg, err
}

/*************************************************
功能:获取平台日志表中的数据总条数
输入:开始时间,结束时间,用户ID,系统类型
输出:int64, error
说明:
编辑:wang_jp
时间:2020年4月20日
*************************************************/
func (log *SysLog) GetSysLogRows(bgtime, endtime string, userid, systype, oprtype int64, desc ...string) (int64, error) {
	o := orm.NewOrm()
	o.Using("default") //选定数据库
	qt := o.QueryTable("SysLog").Filter("StartTime__gte", bgtime).Filter("StartTime__lte", endtime)
	if userid > 0 {
		qt = qt.Filter("User__Id", userid)
	}
	if systype > 0 {
		qt = qt.Filter("SysType", systype)
	}
	if oprtype > 0 {
		qt = qt.Filter("OprType", oprtype)
	}
	if len(desc) > 0 {
		qt = qt.Filter("Description__icontains", desc[0])
	}
	return qt.Count()
}

/*************************************************
功能:获取平台日志表中的数据
输入:开始时间,结束时间,用户ID,系统类型码,操作类型码,可选的描述
	bgtime:开始时间,格式:2006-01-02 15:04:05
	endtime:结束时间,格式:2006-01-02 15:04:05
	userid:用户ID,小于1时取所有用户ID
	systype:系统类型,1：PC端 2：移动端 3:计算服务,小于1时取前述所有
	oprtype:操作类型,0:其他,1:"添加",2:"删除",3:"更新",4:"查看",5:添加/更新",6:"登录"
			小于0时取上述所有。
	limit:每页显示的信息数,小于等于0时不限制
	pages:偏移页数,最小为0
	desc:描述,可选。
输出:总行数,最大页数, []*SysLog,error
说明:
编辑:wang_jp
时间:2020年4月20日
*************************************************/
func (log *SysLog) GetSysLog(bgtime, endtime string, userid, systype, oprtype, limit, pages int64, desc ...string) (int64, int64, []*SysLog, error) {
	o := orm.NewOrm()
	o.Using("default") //选定数据库
	var logs []*SysLog
	qt := o.QueryTable("SysLog").Filter("StartTime__gte", bgtime).Filter("StartTime__lte", endtime)
	if userid > 0 {
		qt = qt.Filter("User__Id", userid)
	}
	if systype > 0 {
		qt = qt.Filter("SysType", systype)
	}
	if oprtype >= 0 {
		qt = qt.Filter("OprType", oprtype)
	}
	if len(desc) > 0 { //有描述参数输入
		if len(desc[0]) > 0 { //且输入的参数长度大于0时才参与过滤
			qt = qt.Filter("Description__icontains", desc[0])
		}
	}
	totalRows, err := qt.Count() //总行数
	var maxpages int64 = 0
	if err == nil {
		if limit > 0 {
			maxpages = totalRows / limit
			if totalRows%limit > 0 {
				maxpages += 1
			}
			if pages > maxpages { //不能超过最大页数
				pages = maxpages
			}
			qt = qt.Limit(limit, pages*limit)
		}
		_, err = qt.RelatedSel().OrderBy("-Id").All(&logs)

		for i, log := range logs {
			user := new(SysUser)
			user.Id = log.User.Id
			user.Name = log.User.Name
			user.Username = log.User.Username
			log.User = user
			logs[i] = log
		}
	}
	return totalRows, maxpages, logs, err
}

/*************************************************
功能:获取平台日志表中的数据并分析
输入:开始时间,结束时间,用户ID,系统类型码,操作类型码,可选的描述
	bgtime:开始时间,格式:2006-01-02 15:04:05
	endtime:结束时间,格式:2006-01-02 15:04:05
	userid:用户ID,小于1时取所有用户ID
	systype:系统类型,1：PC端 2：移动端 3:计算服务,小于1时取前述所有
	oprtype:操作类型,0:其他,1:"添加",2:"删除",3:"更新",4:"查看",5:添加/更新",6:"登录"
			小于0时取上述所有。
	desc:描述,可选。
输出:总行数,最大页数, []*SysLog,error
说明:
编辑:wang_jp
时间:2020年4月20日
*************************************************/
func (slog *SysLog) GetSysLogAnalyse(bgtime, endtime string, userid, systype, oprtype int64, desc ...string) (interface{}, error) {
	o := orm.NewOrm()
	o.Using("default") //选定数据库
	type log struct {
		SysType int64
		ReqUrl  string
		Cnts    int64
	}
	var logs []log
	sqlstr := `SELECT
		log.sys_type,
		log.req_url,
		COUNT( log.req_url ) AS cnts
	FROM
		sys_log log
		LEFT JOIN sys_user user ON user.id=log.user_id
	WHERE
		start_time >= "%s" 
		AND start_time <= "%s"`
	sqlstr = fmt.Sprintf(sqlstr, bgtime, endtime)
	if userid > 0 {
		sqlstr += fmt.Sprintf(` AND user.id = %d`, userid)
	}
	if systype > 0 {
		sqlstr += fmt.Sprintf(` AND log.sys_type = %d`, systype)
	}
	if oprtype >= 0 {
		sqlstr += fmt.Sprintf(` AND log.opr_type = %d`, oprtype)
	}
	if len(desc) > 0 { //有描述参数输入
		if len(desc[0]) > 0 { //且输入的参数长度大于0时才参与过滤
			sqlstr += fmt.Sprintf(` AND log.description LIKE "%s"`, "%"+desc[0]+"%")
		}
	}
	sqlstr += `GROUP BY
			log.req_url,
			log.sys_type 
		ORDER BY
			log.sys_type`
	_, err := o.Raw(sqlstr).QueryRows(&logs)
	return logs, err
}
