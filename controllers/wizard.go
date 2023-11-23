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

	c.Data["OvpnWizardData"] = ovpnWizardData
	c.Data["OvpnProtocolList"] = models.GetovpnProtocolList(ovpnWizardData.OvpnProtocol)
	selectedDNS := models.GetDNSProvidersNameByIP(ovpnWizardData.OvpnDNS1)
	c.Data["SelectedDNS"] = selectedDNS
	c.Data["DNSProvidersList"] = models.GetDNSProvidersList(selectedDNS.Name)

	//store to session
	wizardSaveByte := Encode(ovpnWizardData)
	c.SetSession("ovpnWizardData", wizardSaveByte.Bytes())
}

func (c *WizardController) Step1Post() {
	//get data from session
	var wizardByte = c.GetSession("ovpnWizardData").([]byte)
	ovpnWizardData := Decode(wizardByte)

	//MODIFY ovpnWizardData
	ovpnWizardData.OvpnEndpoint = c.GetString("ip_endpoint")
	ovpnWizardData.OvpnPort = c.GetString("port")
	ovpnWizardData.OvpnProtocol = c.GetString("ovpn_protocol")
	ovpnWizardData.OvpnIPRange = c.GetString("ovpn_ip_range")
	ovpnWizardData.OvpnDNS1 = c.GetString("dns_1")
	ovpnWizardData.OvpnDNS2 = c.GetString("dns_2")
	ovpnWizardData.DisableDefRouteForClientsByDefault, _ = c.GetBool("dis_def_client_routing")
	ovpnWizardData.ClientToClientConfigIsUsed, _ = c.GetBool("client_to_client")

	//save
	a := Encode(ovpnWizardData)
	c.SetSession("ovpnWizardData", a.Bytes())
	lib.Dump(ovpnWizardData)

	c.Redirect("/wizard/step2", 302)
}

func (c *WizardController) Step2Get() {
	c.TplName = "wizard/step_2.html"
	c.Data["breadcrumbs"] = &BreadCrumbs{
		Title: "Wizard: Step 2",
	}

	var wizardByte = c.GetSession("ovpnWizardData").([]byte)
	ovpnWizardData := Decode(wizardByte)

	c.Data["OvpnWizardData"] = ovpnWizardData
	c.Data["OvpnCompressionList"] = models.GetOvpnCompressionList(ovpnWizardData.OvpnCompression)
	c.Data["CipherChoiceList"] = models.GetCipherChoiceList(ovpnWizardData.CipherChoice)
	c.Data["HMACAlgorithmList"] = models.GetHMACAlgorithmList(ovpnWizardData.CipherChoice, ovpnWizardData.HMACAlgorithm)
	c.Data["CertTypeList"] = models.GetCertTypeList(ovpnWizardData.CertType)
	c.Data["CertParamList"] = models.GetCertParamList(ovpnWizardData.CertType, ovpnWizardData.CertCurve)
	c.Data["CCCipherChoiceList"] = models.GetCCCipherChoiceList(ovpnWizardData.CertType, ovpnWizardData.CCCipherChoice)
	c.Data["DHTypeList"] = models.GetDHTypeList(ovpnWizardData.DHType)
	c.Data["DHParamList"] = models.GetDHParamList(ovpnWizardData.DHType, ovpnWizardData.DHCurve)
	c.Data["TLSsigList"] = models.GetTLSsigList(ovpnWizardData.TLSsig)

	//store to session
	wizardSaveByte := Encode(ovpnWizardData)
	c.SetSession("ovpnWizardData", wizardSaveByte.Bytes())
}

func (c *WizardController) Step2Post() {
	//get data from session
	var wizardByte = c.GetSession("ovpnWizardData").([]byte)
	ovpnWizardData := Decode(wizardByte)

	//MODIFY ovpnWizardData
	ovpnWizardData.OvpnCompression = c.GetString("ovpn_compression")
	ovpnWizardData.CipherChoice = c.GetString("cipher_choice")
	ovpnWizardData.HMACAlgorithm = c.GetString("hmac_algorithm")
	ovpnWizardData.CertType = c.GetString("cert_type")
	ovpnWizardData.CertCurve = c.GetString("cert_params")
	ovpnWizardData.RSAKeySize = c.GetString("cert_params")
	ovpnWizardData.CCCipherChoice = c.GetString("cert_cipher")
	ovpnWizardData.DHType = c.GetString("dh_type")
	ovpnWizardData.DHCurve = c.GetString("dh_params")
	ovpnWizardData.DHKeySize = c.GetString("dh_params")
	ovpnWizardData.TLSsig = c.GetString("tls_sig")

	//save
	a := Encode(ovpnWizardData)
	c.SetSession("ovpnWizardData", a.Bytes())
	lib.Dump(ovpnWizardData)

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

// GET ENDPOINTS FOR JSON PARAMS

func (c *WizardController) Step2GetHmacAlg() {
	CipherChoice := c.GetString(":cipher")
	HMACAlgorithm := c.GetString(":selcted_hmac")
	response := models.GetHMACAlgorithmList(CipherChoice, HMACAlgorithm)
	c.Data["json"] = response
	c.ServeJSON()
}

func (c *WizardController) Step2GetCrtParam() {
	certType := c.GetString(":type")
	selectedParam := c.GetString(":selcted_option")
	response := models.GetCertParamList(certType, selectedParam)
	c.Data["json"] = response
	c.ServeJSON()
}

func (c *WizardController) Step2GetCrtCipher() {
	certType := c.GetString(":type")
	selectedParam := c.GetString(":selcted_option")
	response := models.GetCCCipherChoiceList(certType, selectedParam)
	c.Data["json"] = response
	c.ServeJSON()
}

func (c *WizardController) Step2GetDhParamr() {
	dhType := c.GetString(":type")
	selectedParam := c.GetString(":selcted_option")
	response := models.GetDHParamList(dhType, selectedParam)
	c.Data["json"] = response
	c.ServeJSON()
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
