package controllers

import (
	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/server/web"
	"github.com/shuricksumy/openvpn-ui/lib"
	"github.com/shuricksumy/openvpn-ui/models"
)

type RoutesController struct {
	BaseController
	ConfigDir string
}

func (c *RoutesController) NestPrepare() {
	if !c.IsLogin {
		c.Ctx.Redirect(302, c.LoginPath())
		return
	}

	settings := models.Settings{Profile: "default"}
	settings.Read("Profile")
	c.Data["Settings"] = &settings
	c.Data["breadcrumbs"] = &BreadCrumbs{
		Title: "Routes Details",
	}

}

// @router /routes [get]
func (c *RoutesController) Get() {
	c.TplName = "routes.html"
	c.showRoutes()
}

func (c *RoutesController) showRoutes() {
	flash := web.NewFlash()

	//get clientsDetails from file
	clientsDetails, err_read := lib.GetClientsDetailsFromFiles()
	if err_read != nil {
		logs.Error(err_read)
		flash.Error("ERROR WHILE READING CLIENTS FROM FILE !")
		flash.Store(&c.Controller)
	}

	// lib.Dump(clientsDetails)
	c.Data["clients"] = &clientsDetails

	// lib.GetRouterClients
	c.Data["routers"] = lib.GetRouterClients(clientsDetails)
}

// @router /routes [post]
func (c *RoutesController) Post() {
	c.TplName = "routes.html"
	flash := web.NewFlash()
	rParams := lib.RouteDetails{}
	if err := c.ParseForm(&rParams); err != nil {
		logs.Error(err)
		flash.Error(err.Error())
		flash.Store(&c.Controller)
	} else {
		//Validate
		//Store
		logs.Error(rParams.RouteID)
		logs.Error(rParams)
	}
	c.showRoutes()
}
