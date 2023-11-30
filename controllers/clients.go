package controllers

import (
	"strconv"

	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/server/web"
	"github.com/shuricksumy/openvpn-ui/lib"
	"github.com/shuricksumy/openvpn-ui/models"
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

// @router /clients [get]
func (c *ClientsController) Get() {
	c.TplName = "clients.html"
	c.ShowClients()

}

func (c *ClientsController) ShowClients() {
	flash := web.NewFlash()
	clients, err := models.GetAllClientDetails()
	lib.Dump(clients)
	if err == nil {
		c.Data["Clients"] = &clients
	} else {
		c.Data["Clients"] = map[string]string{"error": "Failed to get all ClientDetails"}
	}
	// c.ServeJSON()
	if err != nil {
		flash.Error("Clients not found")
		flash.Store(&c.Controller)
	}
}

// @router /clients/newclient [post]
func (c *ClientsController) NewClient() {
	if !c.IsLogin {
		c.Ctx.Redirect(302, c.LoginPath())
		return
	}
	c.TplName = "clients.html"

	flash := web.NewFlash()
	new_client := models.ClientDetails{}
	err_parse := c.ParseForm(&new_client)
	if err_parse != nil {
		logs.Error(err_parse)
		flash.Error("ERROR PARSING !")
		flash.Store(&c.Controller)
	}

	new_client.StaticIP = lib.StringToNilString(c.GetString("static_ip"))

	// // Specify the IDs of the routes to associate with the client
	routeIDs := []int{}

	if err := models.AddNewClient(new_client.ClientName, new_client.StaticIP, new_client.IsRouteDefault, new_client.IsRouter,
		new_client.Description, new_client.MD5Sum, new_client.Passphrase, routeIDs); err == nil {
		flash.Success("New client added successfully")
		flash.Store(&c.Controller)
	} else {
		flash.Error("Failed to add new client: ", err)
		flash.Store(&c.Controller)
	}

	c.ShowClients()
}

// @router /clients/render_modal/ [post]
func (c *ClientsController) RenderModal() {
	flash := web.NewFlash()
	id, _ := c.GetInt("client-name")

	//get clientsDetails from file
	clientsDetails, err_read := models.GetClientDetailsById(id)
	if err_read != nil {
		logs.Error(err_read)
		flash.Error("ERROR WHILE READING CLIENTS FROM FILE !")
		flash.Store(&c.Controller)
	}

	c.Data["Client"] = &clientsDetails
	c.Data["ProvidedRoutes"], _ = models.GetAllRoutesProvided(id)

	c.TplName = "modalClientDetails.html"
	c.Render()
	c.ShowClients()
}

// @router /clients/save_details_data [post]
func (c *ClientsController) SaveClientDetailsData() {
	if !c.IsLogin {
		c.Ctx.Redirect(302, c.LoginPath())
		return
	}
	c.TplName = "clients.html"

	flash := web.NewFlash()

	clientID, _ := c.GetInt("client_id")
	staticIP := lib.StringToNilString(c.GetString("static_ip"))
	description := c.GetString("description")
	isRouteDefaultStr := c.GetString("is_route_default")
	isRouterStr := c.GetString("is_router")
	usedRoutes := c.GetStrings("route_list_selected")

	var isRouteDefault bool

	if isRouteDefaultStr == "true" {
		isRouteDefault = true
	} else {
		isRouteDefault = false
	}

	var isRouter bool
	if isRouterStr == "true" {
		isRouter = true
	} else {
		isRouter = false
	}

	// Specify the IDs of the routes to associate with the client
	var routeIDs []int
	for _, r := range usedRoutes {
		id, _ := strconv.Atoi(r)
		routeIDs = append(routeIDs, id)

	}

	//func UpdateClientDetails(clientID int, staticIP *string, routes []*RouteDetails, description string, isRouteDefault, isRouter bool) error {

	if err := models.UpdateClientDetails(clientID, staticIP, description, isRouteDefault, isRouter); err == nil {
		flash.Success("New client added successfully")
	} else {
		flash.Error("Failed to add new client: ", err)
		flash.Store(&c.Controller)
	}

	if err := models.UnassignAllRoutesFromClient(clientID); err == nil {
		flash.Success("New client added successfully")
		flash.Store(&c.Controller)
	} else {
		flash.Error("Failed to delete routes ", err)
		flash.Store(&c.Controller)
	}

	if err := models.AssignRoutesToClient(clientID, routeIDs); err == nil {
		flash.Success("New client added successfully")
		flash.Store(&c.Controller)
	} else {
		flash.Error("Failed to add routes ", err)
		flash.Store(&c.Controller)
	}

	c.ShowClients()
}
