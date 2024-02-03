package controllers

import (
	"bytes"
	"fmt"
	"path/filepath"
	"text/template"

	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/core/validation"
	"github.com/beego/beego/v2/server/web"
	"github.com/shuricksumy/openvpn-ui/lib"
	"github.com/shuricksumy/openvpn-ui/models"
	"github.com/shuricksumy/openvpn-ui/pkg/openvpn-server-config/client/config"
	"github.com/shuricksumy/openvpn-ui/state"
)

type NewCertParams struct {
	Name        string `form:"Name" valid:"Required;"`
	Description string `form:"description"`
	Passphrase  string `form:"passphrase"`
}

type CertificatesController struct {
	BaseController
	ConfigDir string
}

func (c *CertificatesController) NestPrepare() {
	if !c.IsLogin {
		c.Ctx.Redirect(302, c.LoginPath())
		return
	}
	settings := models.Settings{Profile: "default"}
	settings.Read("Profile")
	c.Data["Settings"] = &settings
	c.Data["breadcrumbs"] = &BreadCrumbs{
		Title: "Clients Certificates",
	}
}

// @router /certificates/:key [get]
func (c *CertificatesController) Download() {
	if !c.IsLogin {
		c.Ctx.Redirect(302, c.LoginPath())
		return
	}
	flash := web.NewFlash()
	c.TplName = "certificates.html"

	name := c.GetString(":key")
	filename := fmt.Sprintf("%s.ovpn", name)

	//keysPath := filepath.Join(state.GlobalCfg.OVConfigPath, "easy-rsa/pki/issued")

	cfgPath, err := c.saveClientConfig(name)
	if err != nil {
		logs.Error(err)
		flash.Error(err.Error())
		flash.Store(&c.Controller)
		c.showCerts()
		return
	}
	data, err := lib.RawReadFile(cfgPath)
	if err != nil {
		logs.Error(err)
		flash.Error(err.Error())
		flash.Store(&c.Controller)
		c.showCerts()
		return
	}

	c.Ctx.Output.Header("Content-Type", "application/octet-stream")
	c.Ctx.Output.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
	if _, err = c.Controller.Ctx.ResponseWriter.Write([]byte(data)); err != nil {
		logs.Error(err)
		flash.Error(err.Error())
		flash.Store(&c.Controller)
		c.showCerts()
	}
}

// @router /certificates [get]
func (c *CertificatesController) Get() {
	if !c.IsLogin {
		c.Ctx.Redirect(302, c.LoginPath())
		return
	}

	c.TplName = "certificates.html"
	c.showCerts()
}

func (c *CertificatesController) showCerts() {
	flash := web.NewFlash()
	path := filepath.Join(state.GlobalCfg.OVConfigPath, "easy-rsa/pki/index.txt")
	certs, err := lib.ReadCerts(path)
	if err != nil {
		logs.Error(err)
		flash.Error(err.Error())
		flash.Store(&c.Controller)
	}

	clientsNoCert, err := models.GetClientsDetailsWithoutCertificate()
	if err == nil {
		c.Data["Clients"] = &clientsNoCert
	} else {
		c.Data["Clients"] = map[string]string{"error": "Failed to get all ClientDetails"}
	}
	if err != nil {
		logs.Error(err)
		flash.Error(err.Error())
		flash.Store(&c.Controller)
	}

	md5Struct := lib.GetMD5StructureFromFS()

	c.TplName = "certificates.html"
	c.Data["certificates"] = &certs
	c.Data["MD5"] = &md5Struct

}

// @router /certificates [post]
func (c *CertificatesController) Post() {
	if !c.IsLogin {
		c.Ctx.Redirect(302, c.LoginPath())
		return
	}
	c.TplName = "certificates.html"
	flash := web.NewFlash()

	clientId := c.GetString("client_name")

	client, err_cl := models.GetClientDetailsById(clientId)
	if err_cl != nil {
		logs.Error(err_cl)
		flash.Error(err_cl.Error())
		flash.Store(&c.Controller)
	} else {
		if err := lib.CreateCertificate(client.ClientName, client.Passphrase); err != nil {
			logs.Error(err)
			flash.Error(err.Error())
			flash.Store(&c.Controller)
		} else {
			clName := lib.StringToNilString(client.ClientName)
			err_upd_cl := models.UpdateClientCertificateById(clientId, clName)
			if err_upd_cl != nil {
				logs.Error(err_cl)
				flash.Error(err_cl.Error())
				flash.Store(&c.Controller)
			} else {
				status := "Active"
				models.UpdateClientCertificateStatusById(clientId, &status)
				flash.Success("Success! Certificate for the name \"" + client.ClientName + "\" has been created")
				flash.Store(&c.Controller)
			}
		}
	}
	c.showCerts()
}

// @router /certificates/revoke/:key [get]
func (c *CertificatesController) Revoke() {
	if !c.IsLogin {
		c.Ctx.Redirect(302, c.LoginPath())
		return
	}
	c.TplName = "certificates.html"
	flash := web.NewFlash()
	name := c.GetString(":key")

	clientDetails, err_cl := models.GetClientDetailsByCertificate(name)
	if err_cl != nil {
		logs.Error(err_cl)
		flash.Error(err_cl.Error())
		flash.Store(&c.Controller)
	}

	if err := lib.RevokeCertificate(name); err != nil {
		logs.Error(err)
		flash.Error(err.Error())
		flash.Store(&c.Controller)
	} else {
		if clientDetails != nil {
			status := "Revoked"
			err_upd := models.UpdateClientCertificateStatusById(clientDetails.Id, &status)
			if err_upd != nil {
				logs.Error(err_upd)
				flash.Error(err_upd.Error())
				flash.Store(&c.Controller)
			}
		}
		flash.Warning("Success! Certificate for the name \"" + name + "\" has been revoked")
		flash.Store(&c.Controller)
	}
	c.showCerts()
}

// @router /certificates/unrevoke/:key [get]
func (c *CertificatesController) UnRevoke() {
	if !c.IsLogin {
		c.Ctx.Redirect(302, c.LoginPath())
		return
	}
	c.TplName = "certificates.html"
	flash := web.NewFlash()
	name := c.GetString(":key")

	clientDetails, err_cl := models.GetClientDetailsByCertificate(name)
	if err_cl != nil {
		logs.Error(err_cl)
		flash.Error(err_cl.Error())
		flash.Store(&c.Controller)
	}

	if err := lib.UnRevokeCertificate(name); err != nil {
		logs.Error(err)
		flash.Error(err.Error())
		flash.Store(&c.Controller)
	} else {
		if clientDetails != nil {
			status := "Active"
			err_upd := models.UpdateClientCertificateStatusById(clientDetails.Id, &status)
			if err_upd != nil {
				logs.Error(err_upd)
				flash.Error(err_upd.Error())
				flash.Store(&c.Controller)
			}
		}
		flash.Warning("Success! Certificate for the name \"" + name + "\" has been UNrevoked")
		flash.Store(&c.Controller)
	}
	c.showCerts()
}

// @router /certificates/restart [get]
func (c *CertificatesController) Restart() {
	if !c.IsLogin {
		c.Ctx.Redirect(302, c.LoginPath())
		return
	}
	lib.Restart()
	c.Redirect(c.URLFor("CertificatesController.Get"), 302)
	// return
}

// @router /certificates/burn/:key/:serial [get]
func (c *CertificatesController) Burn() {
	if !c.IsLogin {
		c.Ctx.Redirect(302, c.LoginPath())
		return
	}
	c.TplName = "certificates.html"
	flash := web.NewFlash()
	CN := c.GetString(":key")
	serial := c.GetString(":serial")

	clientDetails, err_cl := models.GetClientDetailsByCertificate(CN)
	if err_cl != nil {
		logs.Error(err_cl)
		flash.Error(err_cl.Error())
		flash.Store(&c.Controller)
	}
	if err := lib.BurnCertificate(CN, serial); err != nil {
		logs.Error(err)
		flash.Error(err.Error())
		flash.Store(&c.Controller)
	} else {
		if clientDetails != nil {
			err_upd := models.ClearClientCertificateById(clientDetails.Id)
			if err_upd != nil {
				logs.Error(err_upd)
				flash.Error(err_upd.Error())
				flash.Store(&c.Controller)
			}
		}
		flash.Success("Success! Certificate for the name \"" + CN + "\" has been removed")
		flash.Store(&c.Controller)
	}
	c.showCerts()
}

// @router /certificates/revoke/:key [get]
func (c *CertificatesController) Renew() {
	if !c.IsLogin {
		c.Ctx.Redirect(302, c.LoginPath())
		return
	}

	c.TplName = "certificates.html"
	flash := web.NewFlash()
	name := c.GetString(":key")
	serial := c.GetString(":serial")
	if err := lib.RenewCertificate(name, serial); err != nil {
		logs.Error(err)
		flash.Error(err.Error())
		flash.Store(&c.Controller)
	} else {
		flash.Success("Success! Certificate for the name \"" + name + "\"  and \"" + serial + "\" has been renewed")
		flash.Store(&c.Controller)
	}
	c.showCerts()
}

// @router /certificates/save_client_data [post]
func (c *CertificatesController) SaveClientRawData() {
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

	c.TplName = "certificates.html"
	c.showCerts()

}

func validateCertParams(cert NewCertParams) map[string]map[string]string {
	valid := validation.Validation{}
	b, err := valid.Valid(&cert)
	if err != nil {
		logs.Error(err)
		return nil
	}
	if !b {
		return lib.CreateValidationMap(valid)
	}
	return nil
}

func (c *CertificatesController) saveClientConfig(name string) (string, error) {
	if !c.IsLogin {
		c.Ctx.Redirect(302, c.LoginPath())
		return "", nil
	}

	flash := web.NewFlash()
	destPath := filepath.Join(state.GlobalCfg.OVConfigPath, "clients", name+".ovpn")

	if err := lib.CreateOVPNFile(name); err != nil {
		logs.Error(err)
		flash.Error(err.Error())
		flash.Store(&c.Controller)
		return "", err
	} else {
		flash.Success("Success! Certificate for the name \"" + name + "\" has been created")
		flash.Store(&c.Controller)
	}

	client, _ := models.GetClientDetailsByCertificate(name)
	if client.OTPIsEnabled || client.StaticPassIsUsed {
		err_patch := lib.PatchFileAppendBeforeLine(destPath, "<ca>", "auth-user-pass\n")
		if err_patch != nil {
			logs.Error(err_patch)
			flash.Error(err_patch.Error())
			flash.Store(&c.Controller)
			return "", err_patch
		}
	}

	return destPath, nil
}

func GetText(tpl string, c config.Config) (string, error) {
	t := template.New("config")
	t, err := t.Parse(tpl)
	if err != nil {
		return "", err
	}
	buf := new(bytes.Buffer)
	err = t.Execute(buf, c)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

func SaveToFile(tplPath string, c config.Config, destPath string) error {
	tpl, err := lib.RawReadFile(tplPath)
	if err != nil {
		return err
	}

	str, err := GetText(string(tpl), c)
	if err != nil {
		return err
	}

	return lib.RawSaveToFile(destPath, str)
}

// @router /certificates/updatefiles [get]
func (c *CertificatesController) UpdateFiles() {
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

	c.TplName = "certificates.html"
	c.showCerts()
}
