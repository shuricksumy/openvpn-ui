package controllers

import (
	"fmt"
	"os"
	"time"

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

func (c *SystemController) RestartLocalService() {
	flash := web.NewFlash()

	if !c.IsLogin {
		c.Ctx.Redirect(302, c.LoginPath())
		return
	}
	// Stop the process
	err := lib.StopOpenVPN()
	if err != nil {
		logs.Error(err)
		flash.Error(fmt.Sprintf("Error stopping OpenVPN: %s", err))
		flash.Store(&c.Controller)
	}

	err_fw := lib.DisableFWRules()
	if err_fw != nil {
		// c.Ctx.WriteString(fmt.Sprintf("Error deleting FireWall rules: %s", err_fw))
		logs.Error(err_fw)
		flash.Error(fmt.Sprintf("Error deleting FireWall rules: %s", err_fw))
		flash.Store(&c.Controller)
	}

	// Calling Sleep method
	time.Sleep(3 * time.Second)

	// Start the process again
	err = lib.StartOpenVPN()
	if err != nil {
		logs.Error(err)
		flash.Error(fmt.Sprintf("Error starting OpenVPN: %s", err))
		flash.Store(&c.Controller)
	}

	err_fw = lib.EnableFWRules()
	if err_fw != nil {
		// c.Ctx.WriteString(fmt.Sprintf("Error apply FireWall rules: %s", err_fw))
		logs.Error(err_fw)
		flash.Error(fmt.Sprintf("Error apply FireWall rules: %s", err_fw))
		flash.Store(&c.Controller)
	}

	// c.Redirect(c.URLFor("SystemController.Get"), 302)
	flash.Success("OpenVPN server has been restarted")
	flash.Store(&c.Controller)
	c.TplName = "system.html"
}
