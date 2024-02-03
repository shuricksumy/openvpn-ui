package controllers

import (
	"html/template"

	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/server/web"
	"github.com/shuricksumy/openvpn-ui/models"
	"github.com/shuricksumy/openvpn-ui/state"
)

type SettingsController struct {
	BaseController
}

func (c *SettingsController) NestPrepare() {
	if !c.IsLogin {
		c.Ctx.Redirect(302, c.LoginPath())
		return
	}
	c.Data["breadcrumbs"] = &BreadCrumbs{
		Title: "Application Settings",
	}
}

func (c *SettingsController) Get() {
	if !c.IsLogin {
		c.Ctx.Redirect(302, c.LoginPath())
		return
	}
	c.TplName = "settings.html"
	c.Data["xsrfdata"] = template.HTML(c.XSRFFormHTML())
	settings := models.Settings{Profile: "default"}
	_ = settings.Read("Profile")
	c.Data["Settings"] = &settings
}

func (c *SettingsController) Post() {
	if !c.IsLogin {
		c.Ctx.Redirect(302, c.LoginPath())
		return
	}
	c.TplName = "settings.html"

	flash := web.NewFlash()
	settings := models.Settings{Profile: "default"}
	_ = settings.Read("Profile")
	if err := c.ParseForm(&settings); err != nil {
		logs.Warning(err)
		flash.Error(err.Error())
		flash.Store(&c.Controller)
		return
	}
	c.Data["Settings"] = &settings

	o := orm.NewOrm()
	if _, err := o.Update(&settings); err != nil {
		flash.Error(err.Error())
	} else {
		flash.Success("Settings has been updated")
		state.GlobalCfg = settings
	}
	flash.Store(&c.Controller)
}
