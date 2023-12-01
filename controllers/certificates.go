package controllers

import (
	"bytes"
	"fmt"
	"path/filepath"
	"strconv"
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
	name := c.GetString(":key")
	filename := fmt.Sprintf("%s.ovpn", name)

	c.Ctx.Output.Header("Content-Type", "application/octet-stream")
	c.Ctx.Output.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))

	keysPath := filepath.Join(state.GlobalCfg.OVConfigPath, "easy-rsa/pki/issued")

	cfgPath, err := c.saveClientConfig(keysPath, name)
	if err != nil {
		logs.Error(err)
		return
	}
	data, err := lib.RawReadFile(cfgPath)
	if err != nil {
		logs.Error(err)
		return
	}
	if _, err = c.Controller.Ctx.ResponseWriter.Write([]byte(data)); err != nil {
		logs.Error(err)
	}
}

// @router /certificates [get]
func (c *CertificatesController) Get() {
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

	// //get clientsDetails from file
	// clientsDetails, err_read := lib.GetClientsDetailsFromFiles()
	// if err_read != nil {
	// 	logs.Error(err_read)
	// 	flash.Error("ERROR WHILE READING CLIENTS FROM FILE !")
	// 	flash.Store(&c.Controller)
	// }

	// // get md5 sums from file system
	// md5hashs := lib.GetMD5StructureFromFS(clientsDetails)
	// // lib.Dump(md5hashs)
	// c.Data["MD5"] = &md5hashs

	// lib.Dump(certs)
	c.TplName = "certificates.html"
	c.Data["certificates"] = &certs

}

// @router /certificates [post]
func (c *CertificatesController) Post() {
	c.TplName = "certificates.html"
	flash := web.NewFlash()

	clientIdStr := c.GetString("client_name")
	clientId, _ := strconv.Atoi(clientIdStr)

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
				flash.Success("Success! Certificate for the name \"" + client.ClientName + "\" has been created")
				flash.Store(&c.Controller)
			}
		}
	}
	c.showCerts()
}

// @router /certificates/revoke/:key [get]
func (c *CertificatesController) Revoke() {
	c.TplName = "certificates.html"
	flash := web.NewFlash()
	name := c.GetString(":key")
	if err := lib.RevokeCertificate(name); err != nil {
		logs.Error(err)
		flash.Error(err.Error())
		flash.Store(&c.Controller)
	} else {
		flash.Warning("Success! Certificate for the name \"" + name + "\" has been revoked")
		flash.Store(&c.Controller)
	}
	c.showCerts()
}

// @router /certificates/unrevoke/:key [get]
func (c *CertificatesController) UnRevoke() {
	c.TplName = "certificates.html"
	flash := web.NewFlash()
	name := c.GetString(":key")
	if err := lib.UnRevokeCertificate(name); err != nil {
		logs.Error(err)
		flash.Error(err.Error())
		flash.Store(&c.Controller)
	} else {
		flash.Warning("Success! Certificate for the name \"" + name + "\" has been UNrevoked")
		flash.Store(&c.Controller)
	}
	c.showCerts()
}

// @router /certificates/restart [get]
func (c *CertificatesController) Restart() {
	lib.Restart()
	c.Redirect(c.URLFor("CertificatesController.Get"), 302)
	// return
}

// @router /certificates/burn/:key/:serial [get]
func (c *CertificatesController) Burn() {
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

func (c *CertificatesController) saveClientConfig(keysPath string, name string) (string, error) {
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
