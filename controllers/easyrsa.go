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

type EasyRSAConfigController struct {
	BaseController
	ConfigDir string
}

func (c *EasyRSAConfigController) NestPrepare() {
	if !c.IsLogin {
		c.Ctx.Redirect(302, c.LoginPath())
		return
	}

	c.Data["breadcrumbs"] = &BreadCrumbs{
		Title: "OpenVPN EasyRSA Variables",
	}
}

func (c *EasyRSAConfigController) Get() {
	c.TplName = "easyrsa.html"
	flash := web.NewFlash()
	destPathEasyrsaVars := filepath.Join(state.GlobalCfg.OVConfigPath, "easy-rsa/vars")
	easyrsaVars, err := os.ReadFile(destPathEasyrsaVars)
	if err != nil {
		logs.Error(err)
		flash.Error(err.Error())
		flash.Store(&c.Controller)
		return
	}
	c.Data["EasyRSAenv"] = string(easyrsaVars)

	c.Data["xsrfdata"] = template.HTML(c.XSRFFormHTML())
}

func (c *EasyRSAConfigController) Post() {
	c.TplName = "easyrsa.html"

	flash := web.NewFlash()

	destPathEasyrsaVars := filepath.Join(state.GlobalCfg.OVConfigPath, "easy-rsa/vars")
	err1 := lib.BackupFile(destPathEasyrsaVars)
	if err1 != nil {
		logs.Error(err1)
		flash.Error("Error with creating backup file: ", err1)
		flash.Store(&c.Controller)
		return
	}

	err2 := lib.RawSaveToFile(destPathEasyrsaVars, c.GetString("EasyRSAVars"))
	if err2 != nil {
		logs.Error(err2)
		flash.Error("Error with updating file: ", err2)
		flash.Store(&c.Controller)
		return
	}

	easyrsaVars, err := os.ReadFile(destPathEasyrsaVars)
	if err != nil {
		logs.Error(err)
		flash.Error("Error with reading new file: ", err)
		flash.Store(&c.Controller)
		return
	}

	c.Data["EasyRSAenv"] = string(easyrsaVars)

	flash.Success("Config has been updated")
	flash.Store(&c.Controller)
}
