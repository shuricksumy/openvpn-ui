package controllers

import (
	"html/template"
	"os"
	"path/filepath"

	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/server/web"
	"github.com/shuricksumy/openvpn-ui/lib"
	"github.com/shuricksumy/openvpn-ui/state"
)

type ClientConfigController struct {
	BaseController
	ConfigDir string
}

func (c *ClientConfigController) NestPrepare() {
	if !c.IsLogin {
		c.Ctx.Redirect(302, c.LoginPath())
		return
	}

	c.Data["breadcrumbs"] = &BreadCrumbs{
		Title: "OpenVPN Client Configuration",
	}
}

func (c *ClientConfigController) Get() {
	if !c.IsLogin {
		c.Ctx.Redirect(302, c.LoginPath())
		return
	}

	c.TplName = "clientconfig.html"
	flash := web.NewFlash()
	destPathClientTempl := filepath.Join(state.GlobalCfg.OVConfigPath, "client-template.txt")
	clientTemplate, err := os.ReadFile(destPathClientTempl)
	if err != nil {
		logs.Error(err)
		flash.Error(err.Error())
		flash.Store(&c.Controller)
		return
	}
	c.Data["ClientTemplate"] = string(clientTemplate)

	c.Data["xsrfdata"] = template.HTML(c.XSRFFormHTML())
}

func (c *ClientConfigController) Post() {
	if !c.IsLogin {
		c.Ctx.Redirect(302, c.LoginPath())
		return
	}

	c.TplName = "clientconfig.html"

	flash := web.NewFlash()

	destPathClientTempl := filepath.Join(state.GlobalCfg.OVConfigPath, "client-template.txt")
	err1 := lib.BackupFile(destPathClientTempl)
	if err1 != nil {
		logs.Error(err1)
		flash.Error("Error with creating backup file: ", err1)
		flash.Store(&c.Controller)
		return
	}

	err2 := lib.RawSaveToFile(destPathClientTempl, c.GetString("ClientTemplate"))
	if err2 != nil {
		logs.Error(err2)
		flash.Error("Error with updating file: ", err2)
		flash.Store(&c.Controller)
		return
	}

	clientTemplate, err := os.ReadFile(destPathClientTempl)
	if err != nil {
		logs.Error(err)
		flash.Error("Error with reading new file: ", err)
		flash.Store(&c.Controller)
		return
	}

	c.Data["ClientTemplate"] = string(clientTemplate)

	flash.Success("Config has been updated")
	flash.Store(&c.Controller)
}
