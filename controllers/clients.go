package controllers

import (
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

	// // Specify the IDs of the routes to associate with the client
	routeIDs := []int{}

	if err := models.AddNewClient(new_client.ClientName, new_client.StaticIP, new_client.IsRouteDefault, new_client.IsRouter,
		new_client.Description, new_client.MD5Sum, routeIDs); err == nil {
		flash.Success("New client added successfully")
		flash.Store(&c.Controller)
	} else {
		flash.Error("Failed to add new client: ", err)
		flash.Store(&c.Controller)
	}

	c.ShowClients()
}
