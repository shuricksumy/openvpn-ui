package controllers

import (
	"bytes"
	"encoding/base64"
	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/server/web"
	"github.com/shuricksumy/openvpn-ui/lib"
	"github.com/shuricksumy/openvpn-ui/models"
	"github.com/shuricksumy/openvpn-ui/state"
	"image/png"
	"path/filepath"
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
	if !c.IsLogin {
		c.Ctx.Redirect(302, c.LoginPath())
		return
	}

	c.TplName = "clients.html"
	c.ShowClients()

}

func (c *ClientsController) ShowClients() {
	if !c.IsLogin {
		c.Ctx.Redirect(302, c.LoginPath())
		return
	}

	flash := web.NewFlash()
	clients, err := models.GetAllClientDetails()
	// lib.Dump(clients)
	if err == nil {
		c.Data["Clients"] = &clients
	} else {
		c.Data["Clients"] = map[string]string{"error": "Failed to get all ClientDetails"}
	}

	//Get md5 validatios
	md5Struct := lib.GetMD5StructureFromFS()
	c.Data["MD5"] = &md5Struct

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
	new_client.StaticPass = lib.StringToNilString(c.GetString("new_static_pass"))

	// // Specify the IDs of the routes to associate with the client
	routeIDs := []string{}

	if err := models.AddNewClient(new_client.ClientName, new_client.StaticIP, new_client.IsRouteDefault, new_client.IsRouter,
		new_client.Description, new_client.MD5Sum, new_client.Passphrase, routeIDs, new_client.StaticPass); err == nil {
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
	if !c.IsLogin {
		c.Ctx.Redirect(302, c.LoginPath())
		return
	}

	flash := web.NewFlash()
	id := c.GetString("client-name")

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

	clientID := c.GetString("client_id")
	staticIP := lib.StringToNilString(c.GetString("static_ip"))
	description := c.GetString("description")
	isRouteDefaultStr := c.GetString("is_route_default")
	isRouterStr := c.GetString("is_router")
	usedRoutes := c.GetStrings("route_list_selected")
	passphrase := c.GetString("passphrase")

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
	var routeIDs []string
	for _, id := range usedRoutes {
		routeIDs = append(routeIDs, id)

	}

	if err := models.UpdateClientDetails(clientID, staticIP, description, isRouteDefault, isRouter); err == nil {
		models.UpdateMD5SumForClientDetailsByID(clientID, "edited")
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

	if err := models.UpdatePassphraseById(clientID, passphrase); err == nil {
		flash.Success("New client added successfully")
		flash.Store(&c.Controller)
	} else {
		flash.Error("Failed to update passphrase ", err)
		flash.Store(&c.Controller)
	}

	c.ShowClients()
}

// @router /clients/render_modal_raw/ [post]
func (c *ClientsController) RenderModalRaw() {
	if !c.IsLogin {
		c.Ctx.Redirect(302, c.LoginPath())
		return
	}
	flash := web.NewFlash()
	clientName := c.GetString("client-name")

	// Load data from the client-name.txt file.
	data, err := lib.RawReadClientFile(clientName)
	if err != nil {
		logs.Error(err)
		flash.Error("Cannot read " + clientName + " file !")
		flash.Store(&c.Controller)
	}

	client, err_client := models.GetClientDetailsByCertificate(clientName)
	if client != nil {
		providedRoutes, _ := models.GetAllRoutesProvided(client.Id)
		c.Data["RouterProvideRouts"] = &providedRoutes
	} else {
		logs.Error(err_client)
		flash.Error("Cannot find Client in DB")
		flash.Store(&c.Controller)
	}

	c.Data["Client"] = &client
	c.Data["ClientData"] = string(data)

	md5Struct := lib.GetMD5StructureFromFS()
	c.Data["MD5"] = &md5Struct

	c.TplName = "modalClientRaw.html"
	c.Render()
	c.ShowClients()
}

// @router /clients/save_client_data [post]
func (c *ClientsController) SaveClientRawData() {
	if !c.IsLogin {
		c.Ctx.Redirect(302, c.LoginPath())
		return
	}
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
	c.ShowClients()

}

// @router /clients/delclient [post]
func (c *ClientsController) DelClient() {
	if !c.IsLogin {
		c.Ctx.Redirect(302, c.LoginPath())
		return
	}
	c.TplName = "clients.html"

	flash := web.NewFlash()
	clientID := c.GetString(":key")
	// lib.Dump("---DELETE CLIENT")
	// lib.Dump(clientID)

	client, err := models.GetClientDetailsById(clientID)
	if err != nil {
		logs.Error(err)
		flash.Error("Client is not found")
		flash.Store(&c.Controller)
		c.ShowClients()
		return
	}
	providedRoutes, err := models.GetAllRoutesProvided(clientID)
	if err != nil {
		logs.Error(err)
		flash.Error("Error while getting client routes")
		flash.Store(&c.Controller)
		c.ShowClients()
		return
	}

	if client.CertificateName != nil {
		logs.Error("Client can not be deleted. First delete connected certificate")
		flash.Error("Client can not be deleted. First delete connected certificate")
		flash.Store(&c.Controller)
		c.ShowClients()
		return
	}

	if len(providedRoutes) > 0 {
		logs.Error("Client can not be deleted. First delete all provided routes")
		flash.Error("Client can not be deleted. First delete all provided routes")
		flash.Store(&c.Controller)
		c.ShowClients()
		return
	}

	err_del := models.DeleteClientDetailsByID(clientID)
	if err_del != nil {
		logs.Error("Error while deleteing client", err_del)
		flash.Error("Error while deleteing client", err_del)
		flash.Store(&c.Controller)
		c.ShowClients()
		return
	}

	flash.Success("Client was successfuly deleted:" + clientID)
	flash.Store(&c.Controller)

	c.ShowClients()
}

// @router /clients/updatefiles [get]
func (c *ClientsController) UpdateFiles() {
	if !c.IsLogin {
		c.Ctx.Redirect(302, c.LoginPath())
		return
	}
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

	// Update DB with new MD5
	err_upd_md5 := lib.UpdateDBWithLatestMD5()
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
		flash.Warning("Config has been updated but OpenVPN server was NOT reloaded")
	}

	c.TplName = "clients.html"
	c.ShowClients()
}

// @router /clients/render_auth_modal/ [post]
func (c *ClientsController) RenderAuthModal() {
	if !c.IsLogin {
		c.Ctx.Redirect(302, c.LoginPath())
		return
	}
	flash := web.NewFlash()
	id := c.GetString("client-name")

	clientOTP, key, isOTPNew, err_otp := lib.Get2FA(id)
	if err_otp != nil {
		logs.Error(err_otp)
	}

	var buf bytes.Buffer
	img, _ := clientOTP.Image(350, 350)
	png.Encode(&buf, img)
	img64 := base64.StdEncoding.EncodeToString(buf.Bytes())

	//get clientsDetails from file
	clientsDetails, err_read := models.GetClientDetailsById(id)
	if err_read != nil {
		logs.Error(err_read)
		flash.Error("ERROR WHILE READING CLIENTS FROM FILE !")
		flash.Store(&c.Controller)
	}

	c.Data["Client"] = &clientsDetails
	c.Data["IsOTPNew"] = isOTPNew
	c.Data["Key"] = key                   // KEY HEX TO STORE TO DB
	c.Data["Secret"] = clientOTP.Secret() // KEY SECRET FOR OTP GENERATOR
	c.Data["URL"] = clientOTP.URL()
	c.Data["IMG"] = img64
	c.Data["OtpUserName"] = clientOTP.AccountName()

	c.TplName = "modalClientAuth.html"
	c.Render()
	c.ShowClients()
}

// @router /clients/save_auth_data/ [post]
func (c *ClientsController) SaveClientAuthData() {
	if !c.IsLogin {
		c.Ctx.Redirect(302, c.LoginPath())
		return
	}
	c.TplName = "clients.html"
	flash := web.NewFlash()

	clientID := c.GetString("client_id")
	otpKey := c.GetString("otp_key")
	staticPass := c.GetString("static_pass")
	otpUserName := c.GetString("otp_user_name")
	OTPIsEnabled, _ := c.GetBool("otp_is_enabled")
	StaticPassIsUsed, _ := c.GetBool("static_pass_is_enabled")

	err_upd_otp := models.UpdateOTPDataByClientId(clientID, OTPIsEnabled, StaticPassIsUsed, otpKey, staticPass, otpUserName)
	if err_upd_otp != nil {
		logs.Error(err_upd_otp)
		flash.Error("Error while enabling Auth for " + otpUserName)
		flash.Store(&c.Controller)
		c.ShowClients()
		return
	}

	// Redirect to the main page after successful file save.
	flash.Warning("Auth enabled successfully to " + otpUserName + ". Don not forget apply configuration.")
	flash.Store(&c.Controller)

	c.TplName = "clients.html"
	c.ShowClients()

}

// @router /clients/delete_auth_data/ [post]
func (c *ClientsController) DeleteClientAuthData() {
	if !c.IsLogin {
		c.Ctx.Redirect(302, c.LoginPath())
		return
	}
	c.TplName = "clients.html"
	otpUserName := c.GetString("otp_user_name")

	flash := web.NewFlash()
	clientID := c.GetString("client_id")
	err_upd_otp := models.DisableOTPDataByClientId(clientID)
	if err_upd_otp != nil {
		logs.Error("Error while disabling Auth for", err_upd_otp)
		flash.Error("Error while disabling Auth for " + otpUserName)
		flash.Store(&c.Controller)
		c.ShowClients()
		return
	}

	// Redirect to the main page after successful file save.
	flash.Warning("Auth disabled successfully to " + otpUserName + ". Don not forget apply configuration.")
	flash.Store(&c.Controller)

	c.TplName = "clients.html"
	c.ShowClients()
}
