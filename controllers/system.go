package controllers

import (
	"os"

	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/server/web"
	"github.com/d3vilh/openvpn-ui/lib"
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

	destFile, err := lib.Backup()
	if err != nil {
		logs.Error(err)
		flash.Success("Backup process has been broken")
		flash.Store(&c.Controller)
	}

	if _, err := os.Stat(destFile); err != nil {
		flash.Success("File not found to download")
		flash.Store(&c.Controller)
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
