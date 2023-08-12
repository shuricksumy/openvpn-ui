package controllers

import (
	"html/template"
	"os"
	"path/filepath"

	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/server/web"
	mi "github.com/shuricksumy/openvpn-ui/pkg/openvpn-server-config/server/mi"
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
	c.TplName = "clientconfig.html"

	flash := web.NewFlash()

	destPathClientTempl := filepath.Join(state.GlobalCfg.OVConfigPath, "client-template.txt")
	err1 := lib.BackupFile(destPathClientTempl)
	if err1 != nil {
		logs.Error(err1)
		return
	}
	err2 := lib.RawSaveToFile(destPathClientTempl, c.GetString("ClientTemplate"))
	if err2 != nil {
		logs.Error(err2)
		return
	}else{
		flash.Success("Config has been updated")
		client := mi.NewClient(state.GlobalCfg.MINetwork, state.GlobalCfg.MIAddress)
		if err := client.Signal("SIGTERM"); err != nil {
		    flash.Warning("Config has been updated but OpenVPN server was NOT reloaded: " + err.Error())
		}
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

	clientTemplate, err := os.ReadFile(destPathClientTempl)

	if err != nil {
		logs.Error(err)
		return
	}
	c.Data["ClientTemplate"] = string(clientTemplate)

	flash.Store(&c.Controller)
}
