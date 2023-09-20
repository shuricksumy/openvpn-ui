package controllers

import (
	"path/filepath"

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
	clientsDetails, err_read := lib.GetClientsDetailsFromFiles()
	if err_read != nil {
		logs.Error(err_read)
		flash.Error("ERROR WHILE READING CLIENTS FROM FILE !")
		flash.Store(&c.Controller)
	}

	// get md5 sums from file system
	md5hashs := lib.GetMD5StructureFromFS(clientsDetails)
	lib.Dump(md5hashs)
	c.Data["MD5"] = &md5hashs

	lib.Dump(clientsDetails)
	c.Data["clients"] = &clientsDetails
}

// @router /clients/render_modal/ [post]
func (c *ClientsController) RenderModal() {
	flash := web.NewFlash()
	clientName := c.GetString("client-name")

	//get clientsDetails from file
	clientsDetails, err_read := lib.GetClientsDetailsFromFiles()
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

// @router /clients/render_modal_raw/ [post]
func (c *ClientsController) RenderModalRaw() {
	flash := web.NewFlash()
	clientName := c.GetString("client-name")

	// Load data from the client-name.txt file.
	data, err := lib.RawReadClientFile(clientName)
	if err != nil {
		logs.Error(err)
		flash.Error("Cannot read " + clientName + " file !")
		flash.Store(&c.Controller)
	}

	// Pass the client name and data to the modal template for rendering.
	c.Data["ClientName"] = clientName
	c.Data["ClientData"] = string(data)

	clients, err_get := lib.GetClientsDetailsFromFiles()
	clientJSON, err_get := lib.GetClientFromStructure(clients, clientName)
	if err_get != nil {
		logs.Error(err_get)
		flash.Error("Cannot read " + clientName + " file !")
		flash.Store(&c.Controller)
	}
	c.Data["ClientJSON"] = clientJSON

	// get md5 sums from file system
	isMD5valid := lib.GetMD5StatusForClient(clients, clientName)
	c.Data["IsMD5Valid"] = isMD5valid

	c.Layout = ""
	c.TplName = "modalClientRaw.html"
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
		flash.Success("Settings are saved for " + string(client.ClientName) + " to file. " +
			"Do not forget to [Apply Configuration] for all clients at the end.")
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
	err_save := lib.ApplyClientsConfigToFS()
	if err_save != nil {
		logs.Error(err_save)
		flash.Error("ERROR SAVING CLIENTS TO FS !")
		flash.Store(&c.Controller)
		wasError = true
	}

	//UpdateJSON with new MD5
	err_upd_md5 := lib.UpdateJSONWithLatestMD5()
	if err_upd_md5 != nil {
		logs.Error(err_upd_md5)
		flash.Error("ERROR UPATING MD5 TO JSON ! ", err_upd_md5)
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

// @router /clients/restart [get]
func (c *ClientsController) Restart() {
	lib.Restart()
	c.Redirect(c.URLFor("ClientsController.Get"), 302)
	// return
}

// @router /clients/save_client_data [post]
func (c *ClientsController) SaveClientRawData() {
	flash := web.NewFlash()
	clientName := c.GetString("client_name")
	clientData := c.GetString("client_data")

	// Save the data to the client-name.txt file.
	destPathClientConfig := filepath.Join(state.GlobalCfg.OVConfigPath, "ccd", clientName)
	err := lib.RawSaveToFile(destPathClientConfig, clientData)
	if err != nil {
		logs.Error(err)
		flash.Error("Cannot save " + clientName + " file !")
		flash.Store(&c.Controller)
		return
	}

	// Redirect to the main page after successful file save.
	flash.Success("Settings are saved for " + clientName + " to file.")
	flash.Store(&c.Controller)

	c.TplName = "clients.html"
	c.showClients()

}
