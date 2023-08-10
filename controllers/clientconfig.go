package controllers

import (
	"html/template"
	"os"
	"path/filepath"

	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/server/web"
	mi "github.com/d3vilh/openvpn-server-config/server/mi"
	"github.com/d3vilh/openvpn-ui/lib"
	"github.com/d3vilh/openvpn-ui/models"
	"github.com/d3vilh/openvpn-ui/state"
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
	cfg := models.OVConfig{Profile: "default"}
	_ = cfg.Read("Profile")
	c.Data["Settings"] = &cfg

}

func (c *ClientConfigController) Post() {
	c.TplName = "clientconfig.html"
	flash := web.NewFlash()
	cfg := models.OVConfig{Profile: "default"}
	_ = cfg.Read("Profile")
	if err := c.ParseForm(&cfg); err != nil {
		logs.Warning(err)
		flash.Error(err.Error())
		flash.Store(&c.Controller)
		return
	}
	lib.Dump(cfg)
	c.Data["Settings"] = &cfg

	destPathClientTempl := filepath.Join(state.GlobalCfg.OVConfigPath, "client-template.txt")
	err3 := lib.BackupFile(destPathClientTempl)
	if err3 != nil {
		logs.Error(err3)
		return
	}
	err4 := lib.RawSaveToFile(destPathClientTempl, c.GetString("ClientTemplate"))
	if err4 != nil {
		logs.Error(err4)
		return
	}

	// DO NOT SAVE SERVER.CONF
	//destPath := filepath.Join(state.GlobalCfg.OVConfigPath, "server.conf")
	//err := config.SaveToFile(filepath.Join(c.ConfigDir, "openvpn-server-config.tpl"), cfg.Config, destPath)
	//if err != nil {
	//	logs.Warning(err)
	//	flash.Error(err.Error())
	//	flash.Store(&c.Controller)
	//	return
	//}

	o := orm.NewOrm()
	if _, err := o.Update(&cfg); err != nil {
		flash.Error(err.Error())
	} else {
		flash.Success("Config has been updated")
		client := mi.NewClient(state.GlobalCfg.MINetwork, state.GlobalCfg.MIAddress)
		if err := client.Signal("SIGTERM"); err != nil {
			flash.Warning("Config has been updated but OpenVPN server was NOT reloaded: " + err.Error())
		}
	}

	clientTemplate, err := os.ReadFile(destPathClientTempl)
	if err != nil {
		logs.Error(err)
		return
	}
	c.Data["ClientTemplate"] = string(clientTemplate)

	flash.Store(&c.Controller)
}
