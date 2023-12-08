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

type ServerConfigController struct {
	BaseController
	ConfigDir string
}

func (c *ServerConfigController) NestPrepare() {
	if !c.IsLogin {
		c.Ctx.Redirect(302, c.LoginPath())
		return
	}
	c.Data["breadcrumbs"] = &BreadCrumbs{
		Title: "Server Configuration",
	}
}

func (c *ServerConfigController) Get() {
	c.TplName = "serverconfig.html"
	flash := web.NewFlash()
	destPathServerConfig := filepath.Join(state.GlobalCfg.OVConfigPath, "server.conf")
	serverConfig, err := os.ReadFile(destPathServerConfig)
	if err != nil {
		logs.Error(err)
		flash.Error(err.Error())
		flash.Store(&c.Controller)
		return
	}
	c.Data["ServerConfig"] = string(serverConfig)
	c.Data["xsrfdata"] = template.HTML(c.XSRFFormHTML())
}

func (c *ServerConfigController) Post() {
	c.TplName = "serverconfig.html"
	flash := web.NewFlash()

	destPathServerConfig := filepath.Join(state.GlobalCfg.OVConfigPath, "server.conf")

	err1 := lib.BackupFile(destPathServerConfig)
	if err1 != nil {
		logs.Error(err1)
		flash.Error("Error creating backup file: ", err1)
		flash.Store(&c.Controller)
		return
	}

	err2 := lib.RawSaveToFile(destPathServerConfig, c.GetString("ServerConfig"))
	if err2 != nil {
		logs.Error(err2)
		flash.Error("Error saving config file: ", err2)
		flash.Store(&c.Controller)
		return
	}

	serverConfig, err := os.ReadFile(destPathServerConfig)
	if err != nil {
		logs.Error(err)
		flash.Error("Error with reading new config file: ", err)
		flash.Store(&c.Controller)
		return
	}

	// client := mi.NewClient(state.GlobalCfg.MINetwork, state.GlobalCfg.MIAddress)
	// if err := client.Signal("SIGUSR1"); err != nil {
	// 	flash.Warning("Config has been updated but OpenVPN server was NOT reloaded: " + err.Error())
	// } else {
	// 	flash.Success("Config were updated and OpenVPN reloaded it")
	// 	flash.Store(&c.Controller)
	// }

	flash.Warning("Config has been updated but OpenVPN server was NOT reloaded")

	c.Data["ServerConfig"] = string(serverConfig)
	flash.Store(&c.Controller)
}
