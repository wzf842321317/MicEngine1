package controllers

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/astaxie/beego"
	"github.com/bkzy-wangjp/MicEngine/models"
)

type MyController struct {
	beego.Controller
}

/***********************************************
功能:检查登录状态
输入:无
输出:登录用户的ID
说明:检查用户的登录状态,如果用户登录了,返回用户ID;如果没有登录,则跳转到登录页面
编辑:wang_jp
时间:2020年3月12日
************************************************/
func (c *MyController) CheckSession() int64 {
	v := c.GetSession("SessionTag") //获取Session状态值
	if v == nil {                   //如果获取到的值为空
		c.Redirect("/login", 302) //跳转到登录页面
		return 0
	}
	uid, _ := c.GetSession("UId").(int64) //如果获取到了登录状态,返回登录用户的ID值
	return uid
}

/***********************************************
功能:格式化导航栏
输入:导航栏名称
输出:无
说明:根据从Session中获取的用户ID，从数据库中查询授权给该用户的导航菜单；
    并将查询出来的原始数据转换为bootstrap导航栏的HTML字符串，写入给定的
    导航栏名称中(最多支持2级菜单)
编辑:wang_jp
时间:2020年3月15日
************************************************/
func (c *MyController) FormatNavs(navname string) {
	u := new(models.SysUser)
	u.Id = c.CheckSession()          //检查授权,返回授权的ID
	menus, _ := u.GetMenusByUserId() //从数据库读取授权给用户的菜单
	level_1 := 0                     //顶级菜单的数量
	drop_num := 0                    //下拉菜单的数量
	ids := make(map[int64]int)       //顶级菜单的序号列表,k为菜单在数据库中的id,v为其在navs中的序号
	for _, m := range menus {
		if m.Level == 1 { //顶层菜单
			ids[m.Id] = level_1
			level_1 += 1 //顶级菜单的数量
			if m.HasChild == 1 {
				drop_num += 1 //顶级下拉菜单的数量
			}
		}
	}

	dropIdList := make([]int, drop_num) //下拉菜单在navs中的序号列表
	navs := make([]string, level_1)     //顶级导航列表
	navId := 0
	dropId := 0
	thisurl := c.Ctx.Request.RequestURI
	isfindeMuneAsUrl := false //是否发现菜单中有当前URL
	for _, m := range menus {
		pagename := c.GetUrlPageName(m.Url)
		active := ""
		if thisurl == m.Url { //当前页面是否该菜单
			active = "active"
			isfindeMuneAsUrl = true
		}
		var langType = c.Ctx.GetCookie("langType")
		if langType == "en-US" {
			m.Name = m.NameEng
		}
		if m.Level == 1 { //顶层菜单
			if m.HasChild == 0 { //没有子菜单的
				layout_img := `
	<li class="nav-item">
      <a id="%s" class="nav-link %s" href="%s"><img src="%s" style="width:18px;">%s</a>
    </li>`
				layout_noimg := `
	<li class="nav-item">
      <a id="%s" class="nav-link %s" href="%s">%s</a>
    </li>`
				var str string

				if len(m.IconUrl) > 0 {
					str = fmt.Sprintf(layout_img, pagename, active, m.Url, m.IconUrl, m.Name)
				} else {
					str = fmt.Sprintf(layout_noimg, pagename, active, m.Url, m.Name)
				}
				navs[navId] = str

			} else { //有子菜单的
				for _, n := range menus { //检查当前活动页面是否本菜单下的子菜单
					if n.Url == thisurl {
						isfindeMuneAsUrl = true
						if n.Pid == m.Id {
							active = "active"
						}
						break
					}
				}
				layout_img := `
	<li class="nav-item dropdown">
		<a id="%s" class="nav-link dropdown-toggle %s" data-toggle="dropdown" href="%s"><img src="%s" style="width:18px;">%s</a>
	<div class="dropdown-menu">`
				layout_noimg := `
	<li class="nav-item dropdown">
		<a id="%s" class="nav-link dropdown-toggle %s" data-toggle="dropdown" href="%s">%s</a>
	<div class="dropdown-menu">`
				var str string
				if len(m.IconUrl) > 0 {
					str = fmt.Sprintf(layout_img, pagename, active, m.Url, m.IconUrl, m.Name)
				} else {
					str = fmt.Sprintf(layout_noimg, pagename, active, m.Url, m.Name)
				}
				navs[navId] = str
				dropIdList[dropId] = navId
				dropId += 1
			}
			navId += 1
		} else { //非顶层菜单
			layout_img := `
	<a id="%s" class="dropdown-item %s" href="%s"><img src="%s" style="width:18px;">%s</a>`
			layout_noimg := `
	<a id="%s" class="dropdown-item %s" href="%s">%s</a>`
			var str string
			if len(m.IconUrl) > 0 {
				str = fmt.Sprintf(layout_img, pagename, active, m.Url, m.IconUrl, m.Name)
			} else {
				str = fmt.Sprintf(layout_noimg, pagename, active, m.Url, m.Name)
			}
			//将子菜单追加到顶级菜单上
			for k, _ := range ids {
				if k == m.Pid {
					for nk, _ := range navs {
						if nk == ids[m.Pid] {
							navs[nk] += str
							break
						}
					}
					break
				}
			}
		}
	}
	//添加下拉菜单的结尾符号
	for _, nid := range dropIdList {
		navs[nid] += `</div></li>`
	}
	//将各个菜单字符串拼接成一个HTML字符串输出
	menu_str := ""
	for _, nav := range navs {
		menu_str += nav
	}
	c.Data[navname] = menu_str //写导航信息
	//判断是否访问了未经授权的页面
	if isfindeMuneAsUrl == false && thisurl != "/" && thisurl != "/usercenter" && thisurl != "/login" {
		c.SaveUserActionMsg("非法访问", 0, "非法访问,已拒绝")
		lasturl := c.GetSession("LastPage")                                                //获取上一页面
		if lasturl == nil || lasturl.(string) == "/login" || lasturl.(string) == thisurl { //如果上一页面为空或者为登录页面
			c.Redirect("/", 302) //跳转到主页
		} else {
			c.Redirect(lasturl.(string), 302) //返回上一页面
		}
	} else {
		c.SetSession("LastPage", thisurl) //保存本次URL为上次访问的页面
		c.SaveUserActionMsg("进入页面", 4)    //保存访问记录
	}
}

/***********************************************
功能:格式化树形导航栏
输入:页面中接收树形结构的变量名称,页面中接收根节点偏移量的变量名称,是否带变量节点,是否重新创建
输出:无
说明:根据从Session中获取的用户ID，从数据库中查询授权给该用户的树形结构；
    并将查询出来的原始数据转换为zTree的HTML字符串，写入给定的树结构名称中
编辑:wang_jp
时间:2020年3月15日
************************************************/
func (c *MyController) FormatTreeNode(treename, rootpidname string, iswithtag, recreat bool) {
	userid := c.CheckSession() //检查授权,返回授权的ID
	u := new(models.SysUser)
	u.Id = userid
	idcode, _ := c.GetSession("UserName").(string)
	var nodes []*models.MsMiddleCorrelation
	var nodestring string
	var ok bool
	var err error
	var rootpid int64 = 0 //根节点pid
	if iswithtag {
		idcode += "_with_tag"
		if recreat == false { //不重新创建的时候才读取已经缓存的
			md := new(models.MsMiddleCorrelationString)
			nodestring, rootpid, ok = md.GetTreeNodesFromeCache(idcode) //优先从缓存表读取
		}
		if ok == false { //缓存表中不存在
			nodes, err = u.GetTreeNodesWithTagsByUserId() //从数据库读取授权给用户的节点
		}
	} else {
		idcode += "_no_tag"
		if recreat == false { //不重新创建的时候才读取已经缓存的
			md := new(models.MsMiddleCorrelationString)
			nodestring, rootpid, ok = md.GetTreeNodesFromeCache(idcode) //优先从缓存表读取
		}
		if ok == false { //缓存表中不存在
			nodes, err = u.GetTreeNodesByUserId(true) //从数据库读取授权给用户的节点
		}
	}
	if ok == false && err == nil {
		nodelyout := `{id:%d,pId:%d,name:"%s",icon:"%s",open:%t,nodetype:%d,itemid:%d,dotnum:%d,seq:%d,treelevel:"%s",istag:%d,unit:"%s"}`
		tree_str := `<SCRIPT type="text/javascript"> 
	var zNodes =[`
		nodeslen := len(nodes) //节点总数

		rootid := 0             //根节点id
		secendNodeOpen := false //第二层的第一个节点是否打开标记

		for i, node := range nodes {
			istag := 0      //是否变量
			icon_path := "" //图标路径
			switch node.LevelCategory.Id {
			case 1: //矿
				icon_path = "static/img/home1.png"
			case 2: //厂
				//icon_path = "static/img/concentration.svg"
			case 3: //车间
				//icon_path = "static/img/workshop.png"
			case 4: //作业
				//icon_path = "static/img/stage.svg"
			case 5: //设备
				icon_path = "static/img/3.png"
			case 6: //仪表
				icon_path = "static/img/8.png"
			case 7: //电机
				//icon_path = "static/img/motor.svg"
			case 8: //分析仪器
				//icon_path = "static/img/8.png"
			case 9999: //变量
				icon_path = "static/img/2.png"
				istag = 1
			default: //其他

			}
			isopen := false //是否打开
			if i == 0 {     //第一层打开
				isopen = true
				rootid = int(node.Id) //根节点的ID
				rootpid = node.Pid    //根节点的PID
			}
			if int(node.Pid) == rootid && secendNodeOpen == false { //根节点的第一个子节点默认打开
				secendNodeOpen = true
				isopen = true
			}
			tree_str += fmt.Sprintf(nodelyout,
				node.Id,
				node.Pid,
				node.LevelName,
				icon_path,
				isopen,
				node.LevelCategory.Id,
				node.ItemIdInTable,
				node.ConstrutionCode,
				node.Seq,
				node.TreeLevelCode,
				istag,
				node.ConstrutionTableCode,
			)
			if i != nodeslen-1 {
				tree_str += `,`
			}
		}
		tree_str += `];
</SCRIPT>`
		if len(treename) > 0 { //正确设置了结构树名称的时候才写入
			c.Data[rootpidname] = rootpid //根节点修正
			c.Data[treename] = tree_str   //写导航信息
		}
		go func() {
			node := new(models.MsMiddleCorrelationString)
			node.InsertOrUpdateTreeNodesInCache(idcode, tree_str, rootpid)
		}()
	} else {
		c.Data[rootpidname] = rootpid //根节点修正
		c.Data[treename] = nodestring
	}

}

/***********************************************
功能:格式化监控页面树形导航栏
输入:页面中接收树形结构的变量名称,页面中接收根节点偏移量的变量名称,接收frameurl的变量名称
输出:无
说明:根据从Session中获取的用户ID，从数据库中查询授权给该用户的树形结构；
    并将查询出来的原始数据转换为zTree的HTML字符串，写入给定的树结构名称中
返回:无
编辑:wang_jp
时间:2020年3月15日
************************************************/
func (c *MyController) FormatMonitorTree(treename, rootpidname, frameurlname string) {
	userid := c.CheckSession()                 //检查授权,返回授权的ID
	iswan, _ := c.GetSession("FromWan").(bool) //用户来自内网还是外网
	dic := new(models.SysDictionary)
	gfhost, _ := dic.GetGrafanaHost(iswan) //获取Grafana的访问地址
	gfhost = strings.Replace(gfhost, "\n", "", -1)
	gfhost = strings.Replace(gfhost, "\r", "", -1)
	u := new(models.SysUser)
	u.Id = userid
	var rootpid int64 = 0                     //根节点pid
	nodes, err := u.GetMonitorNodesByUserId() //从数据库读取授权给用户的节点

	if err == nil {
		nodelyout := `{id:%d,pId:%d,name:"%s",icon:"%s",open:%t,isParent:%t,level:%d,href:"%s",remark:"%s",seq:%d}`
		tree_str := `<SCRIPT type="text/javascript"> 
	var zNodes =[`
		nodeslen := len(nodes)  //节点总数
		rootid := 0             //根节点id
		firsturl := ""          //第一个有效的URL
		secendNodeOpen := false //第二层的第一个节点是否打开标记
		for i, node := range nodes {
			isparent := false //是否变量
			icon_path := ""   //图标路径
			if node.Folder == 1 {
				isparent = true
			} else {
				icon_path = "static/img/8.png"
			}

			isopen := false //是否打开
			if i == 0 {     //第一层打开
				isopen = true
				rootid = int(node.Id) //根节点的ID
				rootpid = node.Pid    //根节点的PID
			}
			if int(node.Pid) == rootid && secendNodeOpen == false { //根节点的第一个子节点默认打开
				secendNodeOpen = true
				isopen = true
			}

			url := strings.Replace(node.Url, "\n", "", -1)
			url = strings.Replace(url, "\r", "", -1)
			if strings.Contains(url, "http") { //url包含http
				urls := strings.Split(url, "/d/") //用/d/分割
				if len(urls) > 1 {                //分割了至少2段
					url = gfhost
					for i := 1; i < len(urls); i++ {
						url += "/d/" + urls[i] //重新拼接
					}
				}
			} else { //不包含http
				if len(node.Url) > 0 { //配置了URL
					if url[0] != '/' { //第一个字符不是斜杠
						url = "/" + url
					}
					url = gfhost + url
				}
			}

			if len(firsturl) < 1 && len(url) > 0 { //获取第一个有效的url
				firsturl = url
			}
			tree_str += fmt.Sprintf(nodelyout,
				node.Id,
				node.Pid,
				node.Name,
				icon_path,
				isopen,
				isparent,
				node.Level,
				url,
				node.Remark,
				node.Seq,
			)
			if i != nodeslen-1 {
				tree_str += `,`
			}
		}
		tree_str += `];`
		tree_str += `IsWan=` + fmt.Sprint(iswan) + `;`
		tree_str += `BrowserHost="` + fmt.Sprint(c.GetSession("BrowserHost")) + `";`
		tree_str += `</SCRIPT>`
		c.Data[rootpidname] = rootpid   //根节点修正
		c.Data[treename] = tree_str     //写导航信息
		c.Data[frameurlname] = firsturl //写默认显示的Frame信息
	}
}

/***********************************************
功能:格式化报表页面树形导航栏
输入:页面中接收树形结构的变量名称,页面中接收根节点偏移量的变量名称,接收frameurl的变量名称
输出:无
说明:根据从Session中获取的用户ID，从数据库中查询授权给该用户的树形结构；
    并将查询出来的原始数据转换为zTree的HTML字符串，写入给定的树结构名称中
返回:无
编辑:wang_jp
时间:2020年3月15日
************************************************/
func (c *MyController) FormatReportTree(treename, rootpidname string) {
	userid := c.CheckSession() //检查授权,返回授权的ID
	u := new(models.SysUser)
	u.Id = userid
	var rootpid int64 = 0                    //根节点pid
	nodes, err := u.GetReportNodesByUserId() //从数据库读取授权给用户的节点
	hosturl := "http://" + c.Ctx.Request.Host + "/"
	if err == nil {
		nodelyout := `{id:%d,pId:%d,name:"%s",icon:"%s",open:%t,isParent:%t,level:%d,pathurl:"%s",remark:"%s",seq:%d,levelcode:"%s",tplurl:"%s",tplname:"%s"}`
		tree_str := `<SCRIPT type="text/javascript"> 
	var zNodes =[`
		nodeslen := len(nodes)  //节点总数
		rootid := 0             //根节点id
		firsttplid := 0         //第一个有效的模板ID
		secendNodeOpen := false //第二层的第一个节点是否打开标记
		for i, node := range nodes {
			isparent := false //是否变量
			icon_path := ""   //图标路径
			if node.Folder == 1 {
				isparent = true
			} else {
				icon_path = "static/img/8.png"
			}
			if node.Pid == 0 {
				icon_path = "static/img/home1.png"
			}
			isopen := false //是否打开
			if i == 0 {     //第一层打开
				isopen = true
				rootid = int(node.Id) //根节点的ID
				rootpid = node.Pid    //根节点的PID
			}
			if int(node.Pid) == rootid && secendNodeOpen == false { //根节点的第一个子节点默认打开
				secendNodeOpen = true
				isopen = true
			}
			if firsttplid < 1 && node.Folder == 0 && node.Status == 1 { //获取第一个有效的报表模板
				firsttplid = int(node.Id)
			}
			tree_str += fmt.Sprintf(nodelyout,
				node.Id,
				node.Pid,
				node.Name,
				icon_path,
				isopen,
				isparent,
				node.Level,
				node.ResultUrl+"/"+fmt.Sprint(node.Id),
				node.Remark,
				node.Seq,
				node.LevelCode,
				node.TemplateUrl,
				node.TemplateFile,
			)
			if i != nodeslen-1 {
				tree_str += `,`
			}
		}
		tree_str += `];`
		tree_str += `var NODESMSG="";`
		tree_str += `var MICENGINEID=0;`
		tree_str += `var HOSTURL="";`
		if len(nodes) > 0 {
			firstnodejson, err := json.Marshal(nodes)
			if err == nil {
				tree_str += `NODESMSG = ` + string(firstnodejson) + `;`
				tree_str += `MICENGINEID = ` + fmt.Sprintf("%d", models.EngineCfgMsg.CfgMsg.Id) + `;`
				tree_str += `HOSTURL = "` + hosturl + `";`
			}
		}
		tree_str += `</SCRIPT>`
		c.Data[rootpidname] = rootpid //根节点修正
		c.Data[treename] = tree_str   //写导航信息
	}
}

/***********************************************
功能:格式化质检树形导航栏
输入:页面中接收树形结构的变量名称,页面中接收根节点偏移量的变量名称
输出:无
说明:根据从Session中获取的用户ID，从数据库中查询授权给该用户的树形结构；
    并将查询出来的原始数据转换为zTree的HTML字符串，写入给定的树结构名称中
返回:无
编辑:wang_jp
时间:2020年4月9日
************************************************/
func (c *MyController) FormatSampleTree(treename, rootpidname string) {
	userid := c.CheckSession() //检查授权,返回授权的ID
	u := new(models.SysUser)
	u.Id = userid
	var rootpid int64 = 0                           //根节点pid
	nodes, err := u.GetSampleLabTreeNodesByUserId() //从数据库读取授权给用户的节点

	if err == nil {
		nodelyout := `{id:%d,pId:%d,name:"%s",icon:"%s",open:%t,nodetype:%d,itemid:%d,functype:%d,seq:%d,treelevel:"%s",isleaf:%t,funcname:"%s",sampsite:"%s",basetime:"%s",shifthour:%d,isregular:%d}`
		tree_str := `<SCRIPT type="text/javascript"> 
	var zNodes =[`
		nodeslen := len(nodes)  //节点总数
		rootid := 0             //根节点id
		secendNodeOpen := false //第二层的第一个节点是否打开标记

		for i, node := range nodes {
			icon_path := "" //图标路径
			switch node.NodeType {
			case 1: //矿
				icon_path = "static/img/home1.png"
			case 2: //厂
				//icon_path = "static/img/concentration.svg"
			case 3: //车间
				//icon_path = "static/img/workshop.png"
			case 4: //作业
				//icon_path = "static/img/stage.svg"
			case 5: //设备
				icon_path = "static/img/3.png"
			case 6: //仪表
				icon_path = "static/img/8.png"
			case 7: //电机
				//icon_path = "static/img/motor.svg"
			case 8: //分析仪器
				//icon_path = "static/img/8.png"
			case 9999: //变量
				icon_path = "static/img/2.png"
			default: //其他

			}
			isopen := false //是否打开
			if i == 0 {     //第一层打开
				isopen = true
				rootid = int(node.Id) //根节点的ID
				rootpid = node.Pid    //根节点的PID
			}
			if int(node.Pid) == rootid && secendNodeOpen == false { //根节点的第一个子节点默认打开
				secendNodeOpen = true
				isopen = true
			}
			tree_str += fmt.Sprintf(nodelyout,
				node.Id,
				node.Pid,
				node.Name,
				icon_path,
				isopen,
				node.NodeType,
				node.ItemId,
				node.Func,
				node.Seq,
				node.TreeLevel,
				node.IsLeaf,
				node.FuncName,
				node.SamplingSite,
				node.BaseTime,
				node.ShiftHour,
				node.IsRegular,
			)
			if i != nodeslen-1 {
				tree_str += `,`
			}
		}
		tree_str += `];
</SCRIPT>`
		c.Data[rootpidname] = rootpid //根节点修正
		c.Data[treename] = tree_str   //写导航信息
	}
}

/***********************************************
功能:格式化质检树形导航栏
输入:页面中接收树形结构的变量名称,页面中接收根节点偏移量的变量名称
输出:无
说明:根据从Session中获取的用户ID，从数据库中查询授权给该用户的树形结构；
    并将查询出来的原始数据转换为zTree的HTML字符串，写入给定的树结构名称中
返回:无
编辑:wang_jp
时间:2020年4月9日
************************************************/
func (c *MyController) FormatKpiTree(treename, rootpidname string) {
	// userid := c.CheckSession() //检查授权,返回授权的ID
	// u := new(models.SysUser)
	// u.Id = userid
	type treenode struct {
		Id        int
		Pid       int
		Name      string
		Desc      string
		Icon      string
		TagType   string
		BaseTime  string
		ShiftHour int
		IsOpen    bool
		IsForder  int
	}
	var nodes []treenode
	nd := treenode{Id: 0, Name: "KPI", IsForder: 1, IsOpen: true, Icon: "static/img/home1.png", Desc: ""}
	nodes = append(nodes, nd)

	dic := new(models.SysDictionary)
	kpi := new(models.CalcKpiConfigList) //获取Kpi Config List中的TagType列表
	dics, err := dic.GetKpiTagTypeList() //获取系统字典中的TagType列表
	if err == nil {
		typmp, _ := kpi.GetTagTypeLists()
		for _, dc := range dics {
			if _, ok := typmp[dc.DictionaryNameCode]; ok {
				var nd treenode
				nd.Id = int(dc.Id) + 10000000
				nd.IsForder = 1
				nd.Pid = 0
				nd.Name = dc.Name
				nd.Desc = ""
				nd.TagType = dc.DictionaryNameCode
				nodes = append(nodes, nd)
			}
		}
	}

	ids, _, typeids, shifthour, tagtypes, names, descs, basetime, err := kpi.GetKpiConfigBaseMsg()
	if err == nil {
		for i, id := range ids {
			var nd treenode
			nd.Id = id
			nd.Pid = typeids[i] + 10000000
			nd.Name = names[i]
			nd.Desc = descs[i]
			nd.IsForder = 0
			nd.IsOpen = false
			nd.Icon = "static/img/2.png"
			nd.TagType = tagtypes[i]
			nd.BaseTime = basetime[i]
			nd.ShiftHour = shifthour[i]
			nodes = append(nodes, nd)
		}
	}
	if err == nil {
		nodelyout := `{id:%d,pId:%d,name:"%s",icon:"%s",open:%t,isforder:%d,desc:"%s",tagtype:"%s",basetime:"%s",shifthour:%d}`
		tree_str := `<SCRIPT type="text/javascript"> 
	var zNodes =[`
		nodeslen := len(nodes) //节点总数
		rootid := 0            //根节点id

		for i, node := range nodes {
			tree_str += fmt.Sprintf(nodelyout,
				node.Id,
				node.Pid,
				node.Name,
				node.Icon,
				node.IsOpen,
				node.IsForder,
				node.Desc,
				node.TagType,
				node.BaseTime,
				node.ShiftHour,
			)
			if i != nodeslen-1 {
				tree_str += `,`
			}
		}
		tree_str += `];
</SCRIPT>`
		c.Data[rootpidname] = rootid //根节点修正
		c.Data[treename] = tree_str  //写导航信息
	}
}

/***********************************************
功能:格式化物耗模块树形导航栏
输入:页面中接收树形结构的变量名称,页面中接收根节点偏移量的变量名称
输出:无
说明:根据从Session中获取的用户ID，从数据库中查询授权给该用户的树形结构；
    并将查询出来的原始数据转换为zTree的HTML字符串，写入给定的树结构名称中
返回:无
编辑:wang_jp
时间:2020年4月12日
************************************************/
func (c *MyController) FormatGoodsTree(treename, rootpidname string) {
	userid := c.CheckSession() //检查授权,返回授权的ID
	u := new(models.SysUser)
	u.Id = userid
	var rootpid int64 = 0                       //根节点pid
	nodes, err := u.GetGoodsTreeNodesByUserId() //从数据库读取授权给用户的节点

	if err == nil {
		nodelyout := `{id:%d,pId:%d,name:"%s",icon:"%s",open:%t,nodetype:%d,treelevel:"%s",isleaf:%t,basetime:"%s",shifthour:%d}`
		tree_str := `<SCRIPT type="text/javascript"> 
	var zNodes =[`
		nodeslen := len(nodes)  //节点总数
		rootid := 0             //根节点id
		secendNodeOpen := false //第二层的第一个节点是否打开标记

		for i, node := range nodes {
			icon_path := "" //图标路径
			switch node.NodeType {
			case 1: //矿
				icon_path = "static/img/home1.png"
			case 2: //厂
				//icon_path = "static/img/concentration.svg"
			case 3: //车间
				//icon_path = "static/img/workshop.png"
			case 4: //作业
				//icon_path = "static/img/stage.svg"
			case 5: //设备
				icon_path = "static/img/3.png"
			case 6: //仪表
				icon_path = "static/img/8.png"
			case 7: //电机
				//icon_path = "static/img/motor.svg"
			case 8: //分析仪器
				//icon_path = "static/img/8.png"
			case 9999: //变量
				icon_path = "static/img/2.png"
			default: //其他

			}
			isopen := false //是否打开
			if i == 0 {     //第一层打开
				isopen = true
				rootid = int(node.Id) //根节点的ID
				rootpid = node.Pid    //根节点的PID
			}
			if int(node.Pid) == rootid && secendNodeOpen == false { //根节点的第一个子节点默认打开
				secendNodeOpen = true
				isopen = true
			}
			tree_str += fmt.Sprintf(nodelyout,
				node.Id,
				node.Pid,
				node.Name,
				icon_path,
				isopen,
				node.NodeType,
				node.TreeLevel,
				node.IsLeaf,
				node.BaseTime,
				node.ShiftHour,
			)
			if i != nodeslen-1 {
				tree_str += `,`
			}
		}
		tree_str += `];
</SCRIPT>`
		c.Data[rootpidname] = rootpid //根节点修正
		c.Data[treename] = tree_str   //写导航信息
	}
}

/***********************************************
功能:格式化巡检模块树形导航栏
输入:页面中接收树形结构的变量名称,页面中接收根节点偏移量的变量名称
输出:无
说明:根据从Session中获取的用户ID，从数据库中查询授权给该用户的树形结构；
    并将查询出来的原始数据转换为zTree的HTML字符串，写入给定的树结构名称中
返回:无
编辑:wang_jp
时间:2020年4月12日
************************************************/
func (c *MyController) FormatPatrolTree(treename, rootpidname string) {
	userid := c.CheckSession() //检查授权,返回授权的ID
	u := new(models.SysUser)
	u.Id = userid
	var rootpid int64 = 0                       //根节点pid
	nodes, err := u.GetCheckTreeNodesByUserId() //从数据库读取授权给用户的节点
	dep := new(models.MineDeptInfo)
	dpts, err := dep.GetMineDeptLists() //查询部门列表
	dptmp := make(map[int64]string)     //部门ID、部门名称map
	for _, dpt := range dpts {          //格式化map
		dptmp[dpt.Id] = dpt.DepartmentName
	}

	if err == nil {
		nodelyout := `{id:%d,pId:%d,name:"%s",siteid:%d,icon:"%s",open:%t,nodetype:%d,treelevel:"%s",isleaf:%t,basetime:"%s",shifthour:%d,lineid:%d,linename:"%s",deptid:%d,deptname:"%s"}`
		tree_str := `<SCRIPT type="text/javascript"> 
	var zNodes =[`
		nodeslen := len(nodes)  //节点总数
		rootid := 0             //根节点id
		secendNodeOpen := false //第二层的第一个节点是否打开标记

		for i, node := range nodes {
			icon_path := "" //图标路径
			switch node.NodeType {
			case 1: //矿
				icon_path = "static/img/home1.png"
			case 2: //厂
				//icon_path = "static/img/concentration.svg"
			case 3: //车间
				//icon_path = "static/img/workshop.png"
			case 4: //作业
				//icon_path = "static/img/stage.svg"
			case 5: //设备
				icon_path = "static/img/3.png"
			case 6: //仪表
				icon_path = "static/img/8.png"
			case 7: //电机
				//icon_path = "static/img/motor.svg"
			case 8: //分析仪器
				//icon_path = "static/img/8.png"
			case 9999: //变量
				icon_path = "static/img/2.png"
			default: //其他

			}
			isopen := false //是否打开
			if i == 0 {     //第一层打开
				isopen = true
				rootid = int(node.Id) //根节点的ID
				rootpid = node.Pid    //根节点的PID
			}
			if int(node.Pid) == rootid && secendNodeOpen == false { //根节点的第一个子节点默认打开
				secendNodeOpen = true
				isopen = true
			}

			tree_str += fmt.Sprintf(nodelyout,
				node.Id,
				node.Pid,
				node.Name, //名称
				node.SiteId,
				icon_path,
				isopen,
				node.NodeType,      //节点类型
				node.TreeLevel,     //层级码
				node.IsLeaf,        //叶子节点
				node.BaseTime,      //基准时间
				node.ShiftHour,     //每班时间
				node.LineId,        //线路ID
				node.LineName,      //线路名
				node.DeptId,        //所属部门ID
				dptmp[node.DeptId], //所属部门名称
			)
			if i != nodeslen-1 {
				tree_str += `,`
			}
		}
		tree_str += `];
</SCRIPT>`
		c.Data[rootpidname] = rootpid //根节点修正
		c.Data[treename] = tree_str   //写导航信息
	}
}

/***********************************************
功能:获取IP地址信息
输入:无
输出:IP地址
说明:
编辑:wang_jp
时间:2020年3月17日
************************************************/
func (c *MyController) GetIp() string {
	return c.Ctx.Request.RemoteAddr //nginx 中proxy_set_header 设置的值
}

/***********************************************
功能:格式化智能计算页面树形结构
输入:页面中接收树形结构的变量名称,页面中接收根节点偏移量的变量名称
输出:无
说明:根据从Session中获取的用户ID，从数据库中查询授权给该用户的树形结构；
    并将查询出来的原始数据转换为zTree的HTML字符串，写入给定的树结构名称中
返回:无
编辑:wang_j
时间:2020年10月16日
************************************************/
func (c *MyController) FormatCalculateTree(){

}
/***********************************************
功能:记录用户访问信息
输入:无
输出:数据库日志
说明:
编辑:wang_jp
时间:2020年3月17日
************************************************/
func (c *MyController) SaveUserActionMsg(act string, oprtype int64, msg ...string) {
	log := new(models.SysLog)
	log.Description = act
	log.SysType = 3
	log.OprType = oprtype
	log.ReqMethod = c.Ctx.Request.Method
	log.ClassName = c.Ctx.Request.Host
	strs := strings.Split(c.GetUrlPageName(c.Ctx.Request.RequestURI), "?")
	str := ""
	if len(strs) > 0 {
		str = strs[0]
	}
	log.MethodName = str
	strs = strings.Split(c.Ctx.Request.RequestURI, "?")
	str = c.Ctx.Request.RequestURI
	if len(strs) > 0 {
		str = strs[0]
	}
	log.ReqUrl = str  //保存请求的url
	if len(msg) > 0 { //如果输入了信息
		for i, s := range msg {
			log.ReqParams += s
			if i < len(msg)-1 {
				log.ReqParams += ";"
			}
		}
	} else { //没有输入
		if len(strs) > 1 {
			str = strs[1]
		} else {
			str = string(c.Ctx.Input.RequestBody)
		}
		log.ReqParams = str //保存请求的url
	}
	uid := c.GetSession("UId")
	user := new(models.SysUser)
	if uid != nil {
		user.Id = c.GetSession("UId").(int64)
		log.User = user
		log.StartTime = models.TimeFormat(time.Now())
		log.RemoteIp = c.GetIp()
		if log.InsertUserLog() != nil {
		}
	} else {
		log.User = user
		log.StartTime = models.TimeFormat(time.Now())
		log.RemoteIp = c.GetIp()
		if log.InsertUserLog() != nil {

		}
	}
}

/***********************************************
功能:初始化页面模板
输入:无
输出:无
说明:
编辑:wang_jp
时间:2020年3月17日
************************************************/
func (c *MyController) InitPageTemplate(pagename ...string) {
	c.Data["WebTitle"] = models.EngineCfgMsg.CfgMsg.ProjectName    //页面标题
	c.Data["ProjectName"] = models.EngineCfgMsg.CfgMsg.ProjectName //项目名称
	c.Data["IcoPic"] = models.EngineCfgMsg.CfgMsg.IcoPath          //ICON
	c.Data["LogoPic"] = models.EngineCfgMsg.CfgMsg.LogoPath        //Logo
	c.FormatNavs("Navs")                                           //导航栏,含权限验证和进入页面日志记录
	c.Data["UserName"] = c.GetSession("Name")                      //导航栏上显示的用户信息
	c.Data["JsFileName"] = pagename
	langType := c.Ctx.GetCookie("langType")
	//pagename 数组 需判断pagename 是否为nil 后判断长度
	if pagename != nil && len(pagename) > 0 {
		if langType != "zh-CN" {
			c.Data["JsFileName"] = langType + "/" + pagename[0]
		} else {
			c.Data["JsFileName"] = pagename[0]
		}
	} else {
	}
}

/***********************************************
功能:获取URL中的最后一个字符串(即页面名称)
输入:需要分割的url
输出:页面名称
说明:
编辑:wang_jp
时间:2020年3月17日
************************************************/
func (c *MyController) GetUrlPageName(url ...string) string {
	var urls []string
	if len(url) > 0 { //有输入url参数的时候分割参数
		u := url[0]
		urls = strings.Split(u, "/")
	} else { //没有输入参数的时候分割控制器传来的参数
		urls = strings.Split(c.Ctx.Request.RequestURI, "/")
	}
	urllen := len(urls) //获取分割后的参数长度
	switch urllen {
	case 0:
		return ""
	case 1: //只有一个值的时候
		return urls[0]
	case 2: //有两个值的时候
		if len(urls[1]) > 0 { //检查最后一个值是否有数据
			return urls[1] //有就返回该数据
		} else { //没有数据
			return urls[0] //就返回第一个数据
		}
	default: //有多个值的时候
		if len(urls[urllen-1]) > 0 { //检查最后一个值是否有数据
			return urls[urllen-1] //有就返回该数据
		} else { //没有数据
			return urls[urllen-2] //就返回倒数第二个数据
		}
	}
}

/***********************************************
功能:获取平台路径
输入:无
输出:平台路径,以"/"结尾，比如:http://127.0.0.1:9999
说明:根据浏览器输入的IP地址判断是否广域网，并根据是否广域网获取外网地址或者内网地址
	如果sys_dictionary中没有配置，则以取calc_kpi_engine_config中的配置
编辑:wang_jp
时间:2020年8月14日
************************************************/
func (c *MyController) GetPlatHost() string {
	iswan, _ := c.GetSession("FromWan").(bool) //用户来自内网还是外网
	dic := new(models.SysDictionary)
	host, err := dic.GetPlatHost(iswan)
	if err != nil {
		host = models.EngineCfgMsg.CfgMsg.PlatPath
	}
	if len(host) > 1 {
		if host[len(host)-1] != '/' {
			host += "/"
		}
	}
	return host
}
