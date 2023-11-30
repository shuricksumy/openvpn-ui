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

	routes, err := models.GetAllRoutesDetails()
	lib.Dump(routes)
	if err == nil {
		c.Data["Routes"] = &routes
	} else {
		c.Data["Routes"] = map[string]string{"error": "Failed to get all RoutesDetails"}
	}

	if err != nil {
		flash.Error("Routes not found")
		flash.Store(&c.Controller)
	}

	routers, err := models.GetRouterClients()
	lib.Dump(routers)
	if err == nil {
		c.Data["Routers"] = &routers
	} else {
		c.Data["Routers"] = map[string]string{"error": "Failed to get router clients"}
	}

	if err != nil {
		flash.Error("Routers not found")
		flash.Store(&c.Controller)
	}
}

// @router /routes/newroute [post]
func (c *RoutesController) NewRoute() {
	if !c.IsLogin {
		c.Ctx.Redirect(302, c.LoginPath())
		return
	}
	c.TplName = "routes.html"

	flash := web.NewFlash()
	new_route := models.RouteDetails{}
	err_parse := c.ParseForm(&new_route)
	if err_parse != nil {
		logs.Error(err_parse)
		flash.Error("ERROR PARSING !")
		flash.Store(&c.Controller)
	}
	new_route.Name = new_route.RouterName + "_" + lib.GenRandomString(5)

	if err := models.AddNewRouteDetails(new_route.Name, new_route.RouterName, new_route.RouteIP,
		new_route.RouteMask, new_route.RouteMask); err == nil {

		flash.Success("New route added successfully")
		flash.Store(&c.Controller)
	} else {
		flash.Error("Failed to add new route: ", err)
		flash.Store(&c.Controller)
	}
	c.showRoutes()
}
