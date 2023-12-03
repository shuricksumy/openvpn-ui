package controllers

import (
	"strconv"
	"strings"

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
	// lib.Dump(routes)
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
	// lib.Dump(routers)
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

	values := strings.Split(c.GetString("router_name"), ",")
	if len(values) >= 2 {
		new_route.RouterName = values[1]
		new_route.Name = values[1] + "_" + lib.GenRandomString(5)
		new_route.Id, _ = strconv.Atoi(values[0])
	} else {
		logs.Error(err_parse)
		flash.Error("ERROR PARSING ROUTERS PARAMS!")
		flash.Store(&c.Controller)
	}

	if err := models.AddNewRouteDetails(new_route.Name, new_route.RouterName, new_route.Id, new_route.RouteIP,
		new_route.RouteMask, new_route.Description); err == nil {

		flash.Success("New route added successfully")
		flash.Store(&c.Controller)
	} else {
		flash.Error("Failed to add new route: ", err)
		flash.Store(&c.Controller)
	}
	c.showRoutes()
}

// @router /routes/get/ID [get]
func (c *RoutesController) GetRouteDetails() {
	if !c.IsLogin {
		c.Ctx.Redirect(302, c.LoginPath())
		return
	}

	flash := web.NewFlash()
	routeIDstr := c.GetString(":key")
	routeID, _ := strconv.Atoi(routeIDstr)
	route, _ := models.GetRouteDetailsByID(routeID)

	if route == nil {
		c.TplName = "routes.html"
		flash.Error("Route: " + routeIDstr + " was found. ")
		flash.Store(&c.Controller)
		c.showRoutes()
	} else {
		c.Data["Route"] = &route
		c.TplName = "modalRouteEdit.html"
		c.Render()
	}
}

// @router /routes [post]
func (c *RoutesController) Post() {
	if !c.IsLogin {
		c.Ctx.Redirect(302, c.LoginPath())
		return
	}

	c.TplName = "routes.html"
	flash := web.NewFlash()

	routeID, _ := c.GetInt("route_id")
	routeIP := c.GetString("route_ip")
	routeMask := c.GetString("route_mask")
	description := c.GetString("description")

	err := models.UpdateRouteDetails(routeID, routeIP, routeMask, description)
	if err != nil {
		flash.Error("Error updating route details", err)
		flash.Store(&c.Controller)
		c.showRoutes()
	} else {
		flash.Success("Route was updates successfuly")
		flash.Store(&c.Controller)
		c.showRoutes()
	}
}

// @router /routes/delte/ID [get]
func (c *RoutesController) Delete() {
	if !c.IsLogin {
		c.Ctx.Redirect(302, c.LoginPath())
		return
	}
	flash := web.NewFlash()
	routeIDstr := c.GetString(":key")
	routeID, _ := strconv.Atoi(routeIDstr)

	routeIsUsed := models.RouteIsUsedBy(routeID)
	if len(routeIsUsed) != 0 {
		flash.Error("Route: " + routeIDstr + " was NOT deleted. It's Used !")
		flash.Store(&c.Controller)
	} else {
		err_del := models.DeleteRouteDetailsById(routeID)
		if err_del != nil {
			flash.Error("Route: " + routeIDstr + " was NOT deleted. " + string(err_del.Error()))
			flash.Store(&c.Controller)
		} else {
			flash.Success("Route: " + routeIDstr + " was successfully deleted. Please do not forget apply new config at the end of configuration!")
			flash.Store(&c.Controller)
		}
	}
	c.TplName = "routes.html"
	c.showRoutes()
}
