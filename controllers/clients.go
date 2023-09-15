package controllers

import (
	"strings"

	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/server/web"
	"github.com/shuricksumy/openvpn-ui/lib"
	"github.com/shuricksumy/openvpn-ui/models"
	mi "github.com/shuricksumy/openvpn-ui/pkg/openvpn-server-config/server/mi"
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

// @router /clients [get]
func (c *ClientsController) Get() {
	c.TplName = "clients.html"
	c.showClients()
}

func (c *ClientsController) showClients() {
	flash := web.NewFlash()

	//get clientsDetails from file
	clientsDetails, err_read := lib.GetClientsDetailsFromFile()
	if err_read != nil {
		logs.Error(err_read)
		flash.Error("ERROR WHILE READING CLIENTS FROM FILE !")
		flash.Store(&c.Controller)
	}

	lib.Dump(clientsDetails)
	c.Data["clients"] = &clientsDetails
}

// @router /clients/render_modal/ [post]
func (c *ClientsController) RenderModal() {
	flash := web.NewFlash()
	clientName := c.GetString("client-name")

	//get clientsDetails from file
	clientsDetails, err_read := lib.GetClientsDetailsFromFile()
	if err_read != nil {
		logs.Error(err_read)
		flash.Error("ERROR WHILE READING CLIENTS FROM FILE !")
		flash.Store(&c.Controller)
	}

	foundClient, err_found := lib.GetClientFromStructure(clientsDetails, clientName)
	if err_found != nil {
		logs.Error(err_found)
		flash.Error("Cannot found " + clientName + " !")
		flash.Store(&c.Controller)
	}

	c.Data["ClientName"] = foundClient.ClientName
	c.Data["StaticIP"] = foundClient.StaticIP
	c.Data["IsRouteDefault"] = foundClient.IsRouteDefault
	c.Data["IsRouter"] = foundClient.IsRouter
	c.Data["RouterSubnet"] = foundClient.RouterSubnet
	c.Data["RouterMask"] = foundClient.RouterMask
	c.Data["Description"] = foundClient.Description
	c.Data["RouteListSelected"] = foundClient.RouteList
	c.Data["RouteListUnselected"] = foundClient.RouteListUnselected
	c.Data["CSRFToken"] = foundClient.CSRFToken

	c.TplName = "modalClientDetails.html"
	c.Render()
	c.showClients()
}

// @router /clients/save_details_data [post]
func (c *ClientsController) SaveClientDetailsData() {
	flash := web.NewFlash()
	wasError := false

	//get cleint detais from web form
	client := &lib.ClientDetails{}
	err_parse := c.ParseForm(client)
	if err_parse != nil {
		wasError = true
		logs.Error(err_parse)
		flash.Error("ERROR PARSING !")
		flash.Store(&c.Controller)
	}

	err_save := lib.AddClientToJsonFile(*client)
	if err_save != nil {
		wasError = true
		logs.Error(err_save)
		flash.Error("ERROR SAVINF TO JSON FILE !")
		flash.Store(&c.Controller)
	}

	if !wasError {
		// Redirect to the main page after successful file save.
		flash.Success("Settings are saved for " + string(client.ClientName) + " to file.")
		flash.Store(&c.Controller)
	}

	c.TplName = "clients.html"
	c.showClients()
}

// @router /clients/updatefiles [get]
func (c *ClientsController) UpdateFiles() {
	flash := web.NewFlash()
	wasError := false

	//update files
	err_save := lib.GenerateClientsFileToFS()
	if err_save != nil {
		logs.Error(err_save)
		flash.Error("ERROR SAVING CLIENTS TO FS !")
		flash.Store(&c.Controller)
		wasError = true
	}

	if !wasError {
		// Redirect to the main page after successful file save.
		flash.Success("Clients were updated. Please restart OPENVPN server!")
		flash.Store(&c.Controller)
		client := mi.NewClient(state.GlobalCfg.MINetwork, state.GlobalCfg.MIAddress)
		if err := client.Signal("SIGTERM"); err != nil {
			flash.Warning("Config has been updated but OpenVPN server was NOT reloaded: " + err.Error())
		}
	}

	c.TplName = "clients.html"
	c.showClients()
}

// @router /certificates/restart [get]
func (c *ClientsController) Restart() {
	lib.Restart()
	c.Redirect(c.URLFor("ClientsController.Get"), 302)
	// return
}

func trim(s string) string {
	return strings.Trim(strings.Trim(s, "\r\n"), "\n")
}
