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
	ovpnWizardData := models.GetDefWizardSettings()
	lib.Dump(ovpnWizardData)

	//get extIP
	ipEndpoint, _ := lib.GetExtIP()
	c.Data["RouterRotes"] = ipEndpoint
	lib.Dump(ipEndpoint)

	//store to session
	wizardByte := Encode(ovpnWizardData)
	c.SetSession("ovpnWizardData", wizardByte.Bytes())
	//lib.Dump(wizardByte.Bytes())
}

func (c *WizardController) Step1Post() {
	//get data from session
	var wizardByte = c.GetSession("ovpnWizardData").([]byte)
	ovpnWizardData := Decode(wizardByte)
	//lib.Dump(ovpnWizardData)

	//MODIFY ovpnWizardData
	ovpnWizardData.OvpnEndpoint = "STEP2"
	//clientName := c.GetString("client-name")

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
