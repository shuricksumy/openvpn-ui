package controllers

import (
	"os"

	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/server/web"
	"github.com/shuricksumy/openvpn-ui/lib"
)

type SystemController struct {
	BaseController
	ConfigDir string
}

func (c *SystemController) NestPrepare() {
	if !c.IsLogin {
		c.Ctx.Redirect(302, c.LoginPath())
		return
	}
	c.Data["breadcrumbs"] = &BreadCrumbs{
		Title: "System Utils",
	}
}

// @router /ov/system [get]
func (c *SystemController) Get() {
	c.TplName = "system.html"
	//c.showCerts()
}

// @router /ov/system/backup [get]
func (c *SystemController) Backup() {
	flash := web.NewFlash()
	c.TplName = "system.html"

	destFile, err := lib.Backup()
	if err != nil {
		logs.Error(err)
		flash.Error("Backup process has been broken: ", err)
		flash.Store(&c.Controller)
		return
	}

	if _, err := os.Stat(destFile); err != nil {
		flash.Error("File not found to download:", err)
		flash.Store(&c.Controller)
		return
	}

	c.Ctx.Output.Header("Content-Disposition", "attachment; filename="+destFile)
	c.Ctx.Output.Download(destFile)

}

// @router /ov/system/restart [get]
func (c *SystemController) Restart() {
	lib.Restart()
	c.Redirect(c.URLFor("SystemController.Get"), 302)
	// return
}

func (c *SystemController) RestartLocalService() {
	if !c.IsLogin {
		c.Ctx.Redirect(302, c.LoginPath())
		return
	}

	flash := web.NewFlash()
	msg, err_r := lib.RestartOpenVpnUI()
	if err_r != nil {
		logs.Error(msg)
		flash.Error(msg)
		flash.Store(&c.Controller)
	} else {
		flash.Success(msg)
		flash.Store(&c.Controller)
	}

	c.TplName = "system.html"
}
