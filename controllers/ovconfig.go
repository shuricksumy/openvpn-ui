package controllers

import (
	"html/template"
	"os"
	"path/filepath"

	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/server/web"
	"github.com/shuricksumy/openvpn-ui/lib"
	"github.com/shuricksumy/openvpn-ui/models"
	mi "github.com/shuricksumy/openvpn-ui/pkg/openvpn-server-config/server/mi"
	"github.com/shuricksumy/openvpn-ui/state"
)

type OVConfigController struct {
	BaseController
	ConfigDir string
}

func (c *OVConfigController) NestPrepare() {
	if !c.IsLogin {
		c.Ctx.Redirect(302, c.LoginPath())
		return
	}
	c.Data["breadcrumbs"] = &BreadCrumbs{
		Title: "Server Configuration",
	}
}

func (c *OVConfigController) Get() {
	c.TplName = "ovconfig.html"

	destPathServerConfig := filepath.Join(state.GlobalCfg.OVConfigPath, "server.conf")
	serverConfig, err := os.ReadFile(destPathServerConfig)
	if err != nil {
		logs.Error(err)
		return
	}
	c.Data["ServerConfig"] = string(serverConfig)

	destPathClientTempl := filepath.Join(state.GlobalCfg.OVConfigPath, "client-template.txt")
	clientTemplate, err := os.ReadFile(destPathClientTempl)
	if err != nil {
		logs.Error(err)
		return
	}
	c.Data["ClientTemplate"] = string(clientTemplate)

	c.Data["xsrfdata"] = template.HTML(c.XSRFFormHTML())
	cfg := models.OVConfig{Profile: "default"}
	_ = cfg.Read("Profile")
	c.Data["Settings"] = &cfg

}

func (c *OVConfigController) Post() {
	c.TplName = "ovconfig.html"
	flash := web.NewFlash()
	cfg := models.OVConfig{Profile: "default"}
	_ = cfg.Read("Profile")
	if err := c.ParseForm(&cfg); err != nil {
		logs.Warning(err)
		flash.Error(err.Error())
		flash.Store(&c.Controller)
		return
	}
	//lib.Dump(cfg)
	c.Data["Settings"] = &cfg

	destPathServerConfig := filepath.Join(state.GlobalCfg.OVConfigPath, "server.conf")
	err1 := lib.BackupFile(destPathServerConfig)
	if err1 != nil {
		logs.Error(err1)
		return
	}
	err2 := lib.RawSaveToFile(destPathServerConfig, c.GetString("ServerConfig"))
	if err2 != nil {
		logs.Error(err2)
		return
	}

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

	o := orm.NewOrm()
	if _, err := o.Update(&cfg); err != nil {
		flash.Error(err.Error())
	} else {
		flash.Success("Config has been updated")
		client := mi.NewClient(state.GlobalCfg.MINetwork, state.GlobalCfg.MIAddress)
		if err := client.Signal("SIGUSR1"); err != nil {
			flash.Warning("Config has been updated but OpenVPN server was NOT reloaded: " + err.Error())
		}
	}

	serverConfig, err := os.ReadFile(destPathServerConfig)
	if err != nil {
		logs.Error(err)
		return
	}
	c.Data["ServerConfig"] = string(serverConfig)

	clientTemplate, err := os.ReadFile(destPathClientTempl)
	if err != nil {
		logs.Error(err)
		return
	}
	c.Data["ClientTemplate"] = string(clientTemplate)

	flash.Store(&c.Controller)
}
