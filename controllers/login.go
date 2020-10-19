package controllers

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/bkzy-wangjp/MicEngine/models"
)

type LoginController struct {
	MyController
}

func (c *LoginController) Get() {
	//随机数产生验证码
	rand.Seed(time.Now().UnixNano())                //设定随机数种子
	captcha := fmt.Sprintf("%04d", rand.Intn(9999)) //产生四位随机数
	c.Data["Captcha"] = captcha
	for i, v := range captcha {
		namestr := fmt.Sprintf("Chkcod_%d", i+1)
		c.Data[namestr] = fmt.Sprintf("%c", v)
	}
	c.Data["WebTitle"] = models.EngineCfgMsg.CfgMsg.ProjectName
	c.Data["ProjectName"] = models.EngineCfgMsg.CfgMsg.ProjectName
	c.Data["IcoPic"] = models.EngineCfgMsg.CfgMsg.IcoPath
	var langType = c.Ctx.GetCookie("langType")
	var def_langType = models.EngineCfgMsg.Sys.Lang
	if len(langType) == 0 {
		langType = def_langType
	}
	c.Ctx.SetCookie("langType", langType, 60*60*60*24*3)

	c.Data["JsFileName"] = ""
	if langType != "zh-CN" {
		c.Data["JsFileName"] = langType + "/"
	}
	c.TplName = "login.tpl"
}

//用户登录逻辑
func (this *LoginController) Post() {
	//用户登录Form信息
	type userLoginInfo struct {
		Name    string //用户名
		Pwd     string `form:"Password"` //密码
		chkcod  string `form:-`          //输入的验证码,忽略之
		captcha string `form:-`          //验证码原值,忽略之
		tip     string `form:"tip"`      //登录IP和端口
	}

	var user userLoginInfo
	if err := this.ParseForm(&user); err != nil { //传入user指针
		this.Ctx.WriteString("出错了！")
		this.SaveUserActionMsg("登录失败-登录参数解析错误", 6, "Failed to login[用户登录失败]", string(this.Ctx.Input.RequestBody), err.Error()) //记录用动作作信息
	} else {
		u := new(models.SysUser)
		u.Username = user.Name

		iswan := false //是否广域网
		host := user.tip
		if len(user.tip) < 2 { //没有获取到浏览器地址
			host = this.Ctx.Request.Host
		}
		dic := new(models.SysDictionary)
		micenghost, err := dic.GetMicEngineHost(true) //计算服务访问地址
		if err == nil {
			iswan = strings.Contains(micenghost, host) //计算服务外网地址是否包含当前浏览器的访问地址
			if iswan == false {
				plathost, err := dic.GetPlatHost(true) //平台的外网访问地址
				if err == nil {
					iswan = strings.Contains(plathost, host) //平台外网地址是否包含当前浏览器的访问地址
				}
			}
		}

		if userMsg, err := u.UserLogIn(user.Pwd); err != nil { //读取数据库
			this.Ctx.WriteString(err.Error())
			this.SaveUserActionMsg("登录失败-数据库查询结果错误", 6, "Failed to login[用户登录失败]", string(this.Ctx.Input.RequestBody), err.Error()) //记录用动作作信息
		} else {
			if len(userMsg.Name) == 0 {
				userMsg.Name = userMsg.Username
			}
			this.SetSession("UId", userMsg.Id)
			this.SetSession("Name", userMsg.Name)
			this.SetSession("UserName", userMsg.Username)
			this.SetSession("UserMsg", userMsg)
			this.SetSession("SessionTag", models.Md5str(fmt.Sprintf("%d%s", userMsg.Id, userMsg.Username)))
			this.SetSession("BrowserHost", host)
			this.SetSession("FromWan", iswan) //访问来自外网

			this.SaveUserActionMsg(fmt.Sprintf("用户[ %s ]登录成功", userMsg.Name), 6) //记录用动作作信息

			//lasturl := this.GetSession("LastPage") //获取上一页面
			//if lasturl == nil || lasturl.(string) == "/login" { //如果上一页面为空或者为登录页面
			url := "/"
			if len(models.EngineCfgMsg.Sys.FirstPage) > 0 {
				url = "/" + models.EngineCfgMsg.Sys.FirstPage
			}
			this.Redirect(url, 302) //跳转到主页
			//} else {
			//	this.Redirect(lasturl.(string), 302) //返回上一页面
			//}
		}
	}
}

//用户修改密码逻辑
func (this *LoginController) ApiUpdatePswd() {
	//修改密码用户信息
	type updatepswdInfo struct {
		UserName string `form:"UserName"` //旧密码
		OldPswd  string `form:"OldPswd"`  //旧密码
		NewPswd  string `form:"NewPswd"`  //新密码
		NewPswd2 string `form:"NewPswd2"` //新密码2
	}
	var userpswd updatepswdInfo
	username := this.GetSession("UserName")  //获取当前登录用户名
	name := this.GetSession("Name").(string) //获取当前登录用户名
	u := new(models.SysUser)
	u.Username = username.(string)
	u.Name = name
	if username != nil { //当前登录用户名不为空
		if err := this.ParseForm(&userpswd); err != nil { //传入user指针
			this.Data["json"] = "出错了！"
			this.SaveUserActionMsg("修改密码失败-解析参数错误", 3, "Failed to update password[修改密码失败]", string(this.Ctx.Input.RequestBody), err.Error()) //记录用动作作信息
		} else {
			if ok, err := u.UserUpdatePswd(userpswd.OldPswd, userpswd.NewPswd); err != nil { //读取数据库
				this.Data["json"] = err.Error()
				this.SaveUserActionMsg("修改密码失败-数据库查询错误", 3, "Failed to update password[修改密码失败]", string(this.Ctx.Input.RequestBody), err.Error()) //记录用动作作信息
			} else {
				if ok {
					this.Data["json"] = "1"                                                                                 //修改密码成功
					this.SaveUserActionMsg(fmt.Sprintf("用户[ %s ]修改密码成功", name), 3, "Password changed successfully[修改密码成功]") //记录用动作作信息
				} else {
					this.Data["json"] = "0"                                                                             //修改密码失败
					this.SaveUserActionMsg(fmt.Sprintf("用户[ %s ]修改密码失败", name), 3, "Failed to update password[修改密码失败]") //记录用动作作信息
				}
			}
		}
	} else { //没有当前登录用户
		this.Redirect("/login", 302) //跳转到登录页面
	}
	this.ServeJSON()
}

//用户修改基本信息
func (this *LoginController) ApiUpdateEngineCfg() {
	username := this.GetSession("UserName") //获取当前登录用户名
	//this.SaveUserActionMsg("UpdateEngineCfg", string(this.Ctx.Input.RequestBody))
	if username != nil { //当前登录用户名不为空
		inputMsg := models.EngineCfgMsg.CfgMsg
		cfg := inputMsg
		if err := this.ParseForm(&inputMsg); err != nil { //传入user指针
			this.Data["json"] = "出错了！"
			this.SaveUserActionMsg("修改MicEngine配置信息-失败", 3, "解析信息不成功", string(this.Ctx.Input.RequestBody), err.Error()) //记录用动作作信息
		} else {
			cnt := 0
			if inputMsg.AuthCode != cfg.AuthCode { //授权码改变
				models.EngineCfgMsg.CfgMsg.AuthCode = inputMsg.AuthCode
				models.EngineCfgMsg.CfgMsg.UpdateEngineCfgInfo("AuthCode")
				this.SaveUserActionMsg("修改MicEngine配置信息", 3, "授权码由["+cfg.AuthCode+"]变为["+inputMsg.AuthCode+"]")
				cnt += 1
			}
			if inputMsg.Copyright != cfg.Copyright { //版权单位改变
				models.EngineCfgMsg.CfgMsg.Copyright = inputMsg.Copyright
				models.EngineCfgMsg.CfgMsg.UpdateEngineCfgInfo("Copyright")
				this.SaveUserActionMsg("修改MicEngine配置信息", 3, "版权单位由["+cfg.Copyright+"]变为["+inputMsg.Copyright+"]")
				cnt += 1
			}
			if inputMsg.Description != cfg.Description { //描述
				models.EngineCfgMsg.CfgMsg.Description = inputMsg.Description
				models.EngineCfgMsg.CfgMsg.UpdateEngineCfgInfo("Description")
				this.SaveUserActionMsg("修改MicEngine配置信息", 3, "描述由["+cfg.Description+"]变为["+inputMsg.Description+"]")
				cnt += 1
			}
			if inputMsg.DistributedId != cfg.DistributedId { //分布式ID
				models.EngineCfgMsg.CfgMsg.DistributedId = inputMsg.DistributedId
				models.EngineCfgMsg.CfgMsg.UpdateEngineCfgInfo("DistributedId")
				this.SaveUserActionMsg("修改MicEngine配置信息", 3, "分布式ID由["+fmt.Sprint(cfg.DistributedId)+"]变为["+fmt.Sprint(inputMsg.DistributedId)+"]")
				cnt += 1
			}
			if inputMsg.IcoPath != cfg.IcoPath { //ICO路径
				models.EngineCfgMsg.CfgMsg.IcoPath = inputMsg.IcoPath
				models.EngineCfgMsg.CfgMsg.UpdateEngineCfgInfo("IcoPath")
				this.SaveUserActionMsg("修改MicEngine配置信息", 3, "ICO路径由["+cfg.IcoPath+"]变为["+inputMsg.IcoPath+"]")
				cnt += 1
			}
			if inputMsg.LogoPath != cfg.LogoPath { //Logo路径
				models.EngineCfgMsg.CfgMsg.LogoPath = inputMsg.LogoPath
				models.EngineCfgMsg.CfgMsg.UpdateEngineCfgInfo("LogoPath")
				this.SaveUserActionMsg("修改MicEngine配置信息", 3, "Logo路径由["+cfg.LogoPath+"]变为["+inputMsg.LogoPath+"]")
			}
			if inputMsg.PlatPath != cfg.PlatPath { //平台路径
				models.EngineCfgMsg.CfgMsg.PlatPath = inputMsg.PlatPath
				models.EngineCfgMsg.CfgMsg.UpdateEngineCfgInfo("PlatPath")
				this.SaveUserActionMsg("修改MicEngine配置信息", 3, "平台路径由["+cfg.PlatPath+"]变为["+inputMsg.PlatPath+"]")
				cnt += 1
			}
			if inputMsg.ProjectName != cfg.ProjectName { //项目名称
				models.EngineCfgMsg.CfgMsg.ProjectName = inputMsg.ProjectName
				models.EngineCfgMsg.CfgMsg.UpdateEngineCfgInfo("ProjectName")
				this.SaveUserActionMsg("修改MicEngine配置信息", 3, "项目名称由["+cfg.ProjectName+"]变为["+inputMsg.ProjectName+"]")
				cnt += 1
			}

			if inputMsg.ResultdbClass != cfg.ResultdbClass { //结果数据库类型
				models.EngineCfgMsg.CfgMsg.ResultdbClass = inputMsg.ResultdbClass
				models.EngineCfgMsg.CfgMsg.UpdateEngineCfgInfo("ResultdbClass")
				this.SaveUserActionMsg("修改MicEngine配置信息", 3, "结果数据库类型由["+cfg.ResultdbClass+"]变为["+inputMsg.ResultdbClass+"]")
				cnt += 1
			}
			if inputMsg.ResultdbDbname != cfg.ResultdbDbname { //结果数据库名
				models.EngineCfgMsg.CfgMsg.ResultdbDbname = inputMsg.ResultdbDbname
				models.EngineCfgMsg.CfgMsg.UpdateEngineCfgInfo("ResultdbDbname")
				this.SaveUserActionMsg("修改MicEngine配置信息", 3, "结果数据库名由["+cfg.ResultdbDbname+"]变为["+inputMsg.ResultdbDbname+"]")
				cnt += 1
			}
			if inputMsg.ResultdbTbname != cfg.ResultdbTbname { //结果数据库表名
				models.EngineCfgMsg.CfgMsg.ResultdbTbname = inputMsg.ResultdbTbname
				models.EngineCfgMsg.CfgMsg.UpdateEngineCfgInfo("ResultdbTbname")
				this.SaveUserActionMsg("修改MicEngine配置信息", 3, "结果数据库名由["+cfg.ResultdbTbname+"]变为["+inputMsg.ResultdbTbname+"]")
				cnt += 1
			}
			if inputMsg.ResultdbServer != cfg.ResultdbServer { //结果数据库IP地址
				models.EngineCfgMsg.CfgMsg.ResultdbServer = inputMsg.ResultdbServer
				models.EngineCfgMsg.CfgMsg.UpdateEngineCfgInfo("ResultdbServer")
				this.SaveUserActionMsg("修改MicEngine配置信息", 3, "结果数据库IP地址由["+cfg.ResultdbServer+"]变为["+inputMsg.ResultdbServer+"]")
				cnt += 1
			}
			if inputMsg.ResultdbPort != cfg.ResultdbPort { //结果数据库端口
				models.EngineCfgMsg.CfgMsg.ResultdbPort = inputMsg.ResultdbPort
				models.EngineCfgMsg.CfgMsg.UpdateEngineCfgInfo("ResultdbPort")
				this.SaveUserActionMsg("修改MicEngine配置信息", 3, "结果数据库端口由["+fmt.Sprint(cfg.ResultdbPort)+"]变为["+fmt.Sprint(inputMsg.ResultdbPort)+"]")
				cnt += 1
			}
			if inputMsg.ResultdbUser != cfg.ResultdbUser { //结果数据库用户
				models.EngineCfgMsg.CfgMsg.ResultdbUser = inputMsg.ResultdbUser
				models.EngineCfgMsg.CfgMsg.UpdateEngineCfgInfo("ResultdbUser")
				this.SaveUserActionMsg("修改MicEngine配置信息", 3, "结果数据库用户名由["+cfg.ResultdbUser+"]变为["+inputMsg.ResultdbUser+"]")
				cnt += 1
			}
			if inputMsg.ResultdbPsw != cfg.ResultdbPsw { //结果数据库密码
				models.EngineCfgMsg.CfgMsg.ResultdbPsw = inputMsg.ResultdbPsw
				models.EngineCfgMsg.CfgMsg.UpdateEngineCfgInfo("ResultdbPsw")
				this.SaveUserActionMsg("修改MicEngine配置信息", 3, "结果数据库密码由["+cfg.ResultdbPsw+"]变为["+inputMsg.ResultdbPsw+"]")
				cnt += 1
			}

			if inputMsg.RtdbClass != cfg.RtdbClass { //实时数据库类型
				models.EngineCfgMsg.CfgMsg.RtdbClass = inputMsg.RtdbClass
				models.EngineCfgMsg.CfgMsg.UpdateEngineCfgInfo("RtdbClass")
				this.SaveUserActionMsg("修改MicEngine配置信息", 3, "实时数据库类型由["+cfg.RtdbClass+"]变为["+inputMsg.RtdbClass+"]")
				cnt += 1
			}
			if inputMsg.RtdbDbname != cfg.RtdbDbname { //实时数据库名
				models.EngineCfgMsg.CfgMsg.RtdbDbname = inputMsg.RtdbDbname
				models.EngineCfgMsg.CfgMsg.UpdateEngineCfgInfo("RtdbDbname")
				this.SaveUserActionMsg("修改MicEngine配置信息", 3, "实时数据库名由["+cfg.RtdbDbname+"]变为["+inputMsg.RtdbDbname+"]")
				cnt += 1
			}
			if inputMsg.RtdbTbname != cfg.RtdbTbname { //实时数据库表名
				models.EngineCfgMsg.CfgMsg.RtdbTbname = inputMsg.RtdbTbname
				models.EngineCfgMsg.CfgMsg.UpdateEngineCfgInfo("RtdbTbname")
				this.SaveUserActionMsg("修改MicEngine配置信息", 3, "实时数据库名由["+cfg.RtdbTbname+"]变为["+inputMsg.RtdbTbname+"]")
				cnt += 1
			}
			if inputMsg.RtdbServer != cfg.RtdbServer { //实时数据库IP地址
				models.EngineCfgMsg.CfgMsg.RtdbServer = inputMsg.RtdbServer
				models.EngineCfgMsg.CfgMsg.UpdateEngineCfgInfo("RtdbServer")
				this.SaveUserActionMsg("修改MicEngine配置信息", 3, "实时数据库IP地址由["+cfg.RtdbServer+"]变为["+inputMsg.RtdbServer+"]")
				cnt += 1
			}
			if inputMsg.RtdbPort != cfg.RtdbPort { //实时数据库端口
				models.EngineCfgMsg.CfgMsg.RtdbPort = inputMsg.RtdbPort
				models.EngineCfgMsg.CfgMsg.UpdateEngineCfgInfo("RtdbPort")
				this.SaveUserActionMsg("修改MicEngine配置信息", 3, "实时数据库端口由["+fmt.Sprint(cfg.RtdbPort)+"]变为["+fmt.Sprint(inputMsg.RtdbPort)+"]")
				cnt += 1
			}
			if inputMsg.RtdbUser != cfg.RtdbUser { //实时数据库用户
				models.EngineCfgMsg.CfgMsg.RtdbUser = inputMsg.RtdbUser
				models.EngineCfgMsg.CfgMsg.UpdateEngineCfgInfo("RtdbUser")
				this.SaveUserActionMsg("修改MicEngine配置信息", 3, "实时数据库用户名由["+cfg.RtdbUser+"]变为["+inputMsg.RtdbUser+"]")
				cnt += 1
			}
			if inputMsg.RtdbPsw != cfg.RtdbPsw { //实时数据库密码
				models.EngineCfgMsg.CfgMsg.RtdbPsw = inputMsg.RtdbPsw
				models.EngineCfgMsg.CfgMsg.UpdateEngineCfgInfo("RtdbPsw")
				this.SaveUserActionMsg("修改MicEngine配置信息", 3, "实时数据库密码由["+cfg.RtdbPsw+"]变为["+inputMsg.RtdbPsw+"]")
				cnt += 1
			}
			this.Data["json"] = fmt.Sprintf("更新了%d项配置", cnt)
		}
	} else { //没有当前登录用户
		this.Redirect("/login", 302) //跳转到登录页面
	}
	this.ServeJSON()
}

//管理员重置用户密码
func (this *LoginController) ApiResetPswd() {
	//修改密码用户信息
	type updatepswdInfo struct {
		UserName string `form:"UserName"` //旧密码
		OldPswd  string `form:"OldPswd"`  //旧密码
		NewPswd  string `form:"NewPswd"`  //新密码
		NewPswd2 string `form:"NewPswd2"` //新密码2
	}
	var userpswd updatepswdInfo
	userid := this.GetSession("UId") //获取当前登录用户名
	permissioned := false            //具有权限标志
	if userid != nil {               //当前登录用户名不为空
		u := new(models.SysUser)
		u.Id = userid.(int64)
		menus, err := u.GetMenusByUserId()
		if err == nil { //获取到了管理员权限菜单列表
			for _, menu := range menus {
				if menu.Url == "/managerusers" { //必须具有用户管理权限
					permissioned = true                               //设置有权限标志
					if err := this.ParseForm(&userpswd); err != nil { //传入user指针
						this.Data["json"] = "出错了！"
						this.SaveUserActionMsg("重置密码失败-解析参数错误", 3, "Failed to reset password[修改密码失败]", string(this.Ctx.Input.RequestBody), err.Error()) //记录用动作作信息
					} else {
						u.Username = userpswd.UserName
						if ok, err := u.UserResetPswd(userpswd.NewPswd); err != nil { //读取数据库
							this.Data["json"] = err.Error()
							this.SaveUserActionMsg(fmt.Sprintf("为[ %s ]重置密码失败", userpswd.UserName), 3, "Failed to reset password[修改密码失败]", string(this.Ctx.Input.RequestBody), err.Error()) //记录用动作作信息
						} else {
							if ok {
								this.Data["json"] = "1"                                                                                           //修改密码成功
								this.SaveUserActionMsg(fmt.Sprintf("为[ %s ]重置密码成功", userpswd.UserName), 3, "Password reset successfully[修改密码成功]") //记录用动作作信息
							} else {
								this.Data["json"] = "0"                                                    //修改密码失败
								this.SaveUserActionMsg(fmt.Sprintf("为[ %s ]重置密码成功", userpswd.UserName), 3) //记录用动作作信息
							}
						}
					}
					break
				}
			}
		}
		if permissioned == false { //获取权限菜单失败
			this.Data["json"] = "您没有权限执行此操作！"
			this.SaveUserActionMsg("权限不足", 0, "You do not have permission to perform this operation![您没有权限执行此操作！]", string(this.Ctx.Input.RequestBody), err.Error()) //记录用动作作信息
		}
	} else { //没有当前登录用户
		this.Redirect("/login", 302) //跳转到登录页面
	}
	this.ServeJSON()
}

func (this *LoginController) LogOut() { //用户退出
	name := this.GetSession("Name").(string)                     //获取当前登录用户名
	this.SaveUserActionMsg(fmt.Sprintf("用户[ %s ]退出登录", name), 6) //记录用动作作信息
	this.DelSession("SessionTag")                                //删除Session
	this.DestroySession()                                        //
	this.Redirect("/login", 302)
}

func (this *LoginController) ChangePswd() { //修改密码
	name := this.GetSession("Name").(string)                     //获取当前登录用户名
	this.SaveUserActionMsg(fmt.Sprintf("用户[ %s ]退出登录", name), 6) //记录用动作作信息
	this.DelSession("SessionTag")                                //删除Session
	this.DestroySession()                                        //
	this.Redirect("/login", 302)
}

func (this *LoginController) UserCenter() { //用户中心
	this.CheckSession() //检验权限
	userMsg := this.GetSession("UserMsg").(*models.SysUser)
	pagename := "usercenter"
	//this.Data["ModalSize"] = "modal-lg" //设置模态框大小
	this.Data["JsFileName"] = pagename
	this.InitPageTemplate(pagename) //载入模板数据

	this.Data["Username"] = userMsg.Username
	this.Data["Name"] = userMsg.Name
	this.Data["CardNo"] = userMsg.CardNo
	this.Data["Telephone"] = userMsg.Telephone
	this.Data["Email"] = userMsg.Email
	if userMsg.Sex == 1 {
		this.Data["Sex"] = "男"
	} else {
		this.Data["Sex"] = "女"
	}
	this.Data["Qq"] = userMsg.Qq
	this.Data["Weixin"] = userMsg.Weixin

	this.TplName = pagename + ".tpl"
}
