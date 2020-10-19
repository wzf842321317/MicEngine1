package controllers

import (
	"github.com/bkzy-wangjp/MicEngine/models"
)

type MainController struct {
	MyController //beego.Controller
}

func (c *MainController) Get() {

	c.Data["HeaderTitle"] = models.EngineCfgMsg.CfgMsg.ProjectName
	c.Data["Website"] = "mining-icloud.com"
	c.Data["Email"] = "server@mining-icloud.com"
	c.Data["Version"] = models.EngineCfgMsg.Version
	c.Data["Copyright"] = models.EngineCfgMsg.CfgMsg.Copyright
	c.Data["LogoImg"] = models.EngineCfgMsg.CfgMsg.LogoPath
	pagename := "index"
	c.Data["JsFileName"] = ""

	var langType = c.Ctx.GetCookie("langType")
	if langType != "zh-CN" {
		c.Data["JsFileName"] = langType + "/" + pagename

	} else {
		c.Data["JsFileName"] = pagename
	}
	c.TplName = pagename + ".tpl"
	c.InitPageTemplate(pagename) //载入模板数据

}
