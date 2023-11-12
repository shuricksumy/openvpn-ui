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
	clientsDetails, err_cl_read := lib.GetClientsDetailsFromFiles()
	if err_cl_read != nil {
		logs.Error(err_cl_read)
		flash.Error("ERROR WHILE READING CLIENTS FROM FILE !")
		flash.Store(&c.Controller)
	}

	//get routeDetails from file
	routeDetails, err_read := lib.GetRoutesDetailsFromFiles()
	if err_read != nil {
		logs.Error(err_read)
		flash.Error("ERROR WHILE READING CLIENTS FROM FILE !")
		flash.Store(&c.Controller)
	}

	// lib.Dump(clientsDetails)
	//c.Data["clients"] = &clientsDetails

	// lib.GetRouterClients
	c.Data["routers"] = lib.GetRouterClients(clientsDetails)

	// lib.GetRoutesDetailsFromFiles
	c.Data["routes"] = &routeDetails
}

// @router /routes [post]
func (c *RoutesController) Post() {
	c.TplName = "routes.html"
	flash := web.NewFlash()
	wasError := false

	route := &lib.RouteDetails{}
	err_parse := c.ParseForm(route)
	if err_parse != nil {
		logs.Error(err_parse)
		flash.Error(err_parse.Error())
		flash.Store(&c.Controller)
	}

	err_save := lib.AddRouteToJsonFile(*route)
	if err_save != nil {
		wasError = true
		logs.Error(err_save)
		flash.Error("Error saving JSON file! " + string(err_save.Error()))
		flash.Store(&c.Controller)
	}

	if !wasError {
		// Redirect to the main page after successful file save.
		flash.Success("Settings are saved for " + string(route.RouterName) + " to file. " +
			"Do not forget to [Apply Configuration] for all clients at the end.")
		flash.Store(&c.Controller)
	}
	c.showRoutes()
}
