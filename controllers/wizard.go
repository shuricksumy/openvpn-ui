package controllers

import (
	"bytes"
	"encoding/gob"

	"github.com/shuricksumy/openvpn-ui/lib"
	"github.com/shuricksumy/openvpn-ui/models"
)

type WizardController struct {
	BaseController
	ConfigDir string
}

func (c *WizardController) NestPrepare() {
	if !c.IsLogin {
		c.Ctx.Redirect(302, c.LoginPath())
		return
	}
}

func (c *WizardController) Step1Get() {
	c.TplName = "wizard/step_1.html"
	c.Data["breadcrumbs"] = &BreadCrumbs{
		Title: "Wizard: Step 1",
	}

	//get def settings
	//TODO IF NO SESSION DATA INIT DEFAULT
	wizardByte, ok := c.GetSession("ovpnWizardData").([]byte)

	var ovpnWizardData models.OvpnServerBaseSetting
	if !ok || wizardByte == nil {
		ovpnWizardData = models.GetDefWizardSettings()
	} else {
		ovpnWizardData = Decode(wizardByte)
	}
	lib.Dump(ovpnWizardData)

	//get extIP
	ipEndpoint, _ := lib.GetExtIP()
	c.Data["IpEndpoint"] = ipEndpoint
	//lib.Dump(ipEndpoint)

	c.Data["OvpnWizardData"] = ovpnWizardData
	c.Data["OvpnProtocolList"] = models.GetovpnProtocolList(ovpnWizardData.OvpnProtocol)
	selectedDNS := models.GetDNSProvidersNameByIP(ovpnWizardData.OvpnDNS1)
	c.Data["SelectedDNS"] = selectedDNS
	c.Data["DNSProvidersList"] = models.GetDNSProvidersList(selectedDNS.Name)

	//store to session
	wizardSaveByte := Encode(ovpnWizardData)
	c.SetSession("ovpnWizardData", wizardSaveByte.Bytes())
	//lib.Dump(wizardByte.Bytes())
}

func (c *WizardController) Step1Post() {
	//get data from session
	var wizardByte = c.GetSession("ovpnWizardData").([]byte)
	ovpnWizardData := Decode(wizardByte)
	//lib.Dump(ovpnWizardData)

	//MODIFY ovpnWizardData
	ovpnWizardData.OvpnEndpoint = c.GetString("ip_endpoint")
	ovpnWizardData.OvpnPort = c.GetString("port")
	ovpnWizardData.OvpnProtocol = c.GetString("ovpn_protocol")
	ovpnWizardData.OvpnIPRange = c.GetString("ovpn_ip_range")
	ovpnWizardData.OvpnDNS1 = c.GetString("dns_1")
	ovpnWizardData.OvpnDNS2 = c.GetString("dns_2")

	//save
	a := Encode(ovpnWizardData)
	c.SetSession("ovpnWizardData", a.Bytes())

	c.Redirect("/wizard/step2", 302)
}

func (c *WizardController) Step2Get() {
	c.TplName = "wizard/step_2.html"
	c.Data["breadcrumbs"] = &BreadCrumbs{
		Title: "Wizard: Step 2",
	}

	var wizardByte = c.GetSession("ovpnWizardData").([]byte)
	ovpnWizardData := Decode(wizardByte)
	lib.Dump(ovpnWizardData)

}

func (c *WizardController) Step2Post() {
	//get data from session
	var wizardByte = c.GetSession("ovpnWizardData").([]byte)
	ovpnWizardData := Decode(wizardByte)
	//lib.Dump(ovpnWizardData)

	//MODIFY ovpnWizardData
	ovpnWizardData.OvpnEndpoint = "STEP3"

	//save
	a := Encode(ovpnWizardData)
	c.SetSession("ovpnWizardData", a.Bytes())

	c.Redirect("/wizard/step3", 302)
}

func (c *WizardController) Step3Get() {
	c.TplName = "wizard/step_3.html"
	c.Data["breadcrumbs"] = &BreadCrumbs{
		Title: "Wizard: Step 3",
	}

	var wizardByte = c.GetSession("ovpnWizardData").([]byte)
	ovpnWizardData := Decode(wizardByte)
	lib.Dump(ovpnWizardData)
}

func Encode(ovpnWizardData models.OvpnServerBaseSetting) bytes.Buffer {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	enc.Encode(ovpnWizardData)
	return buf
}

func Decode(wizardByte []byte) models.OvpnServerBaseSetting {
	buf := bytes.NewBuffer(wizardByte)
	var ovpnWizardData models.OvpnServerBaseSetting
	dec := gob.NewDecoder(buf)
	dec.Decode(&ovpnWizardData)
	return ovpnWizardData
}
