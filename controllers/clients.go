package controllers

import (
	"path/filepath"

	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/server/web"
	"github.com/shuricksumy/openvpn-ui/lib"
	"github.com/shuricksumy/openvpn-ui/models"
	"github.com/shuricksumy/openvpn-ui/state"
)

type ClientsController struct {
	BaseController
	ConfigDir string
}

func (c *ClientsController) NestPrepare() {
	if !c.IsLogin {
		c.Ctx.Redirect(302, c.LoginPath())
		return
	}
	settings := models.Settings{Profile: "default"}
	settings.Read("Profile")
	c.Data["Settings"] = &settings
	c.Data["breadcrumbs"] = &BreadCrumbs{
		Title: "Clients Details",
	}
}

// @router /routes [get]
func (c *ClientsController) Get() {
	c.TplName = "clients.html"
	c.showClients()
}

func (c *ClientsController) showClients() {
	flash := web.NewFlash()
	pathIndex := filepath.Join(state.GlobalCfg.OVConfigPath, "easy-rsa/pki/index.txt")
	pathJson := filepath.Join(state.GlobalCfg.OVConfigPath, "clientDetails.json")
	clientsDetails, err := lib.GetClientsDetails(pathIndex, pathJson)
	if err != nil {
		logs.Error(err)
		flash.Error(err.Error())
		flash.Store(&c.Controller)
	}
	lib.Dump(clientsDetails)
	c.Data["clients"] = &clientsDetails
}
