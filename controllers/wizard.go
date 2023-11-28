package controllers

import (
	"bytes"
	"context"
	"encoding/gob"
	"fmt"
	"os/exec"
	"time"

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

	isSetuped := models.IsVPNConfigured()
	if isSetuped {
		c.TplName = "wizard/step_no.html"
		return
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
	//lib.Dump(ovpnWizardData)

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
	ovpnWizardData.TunNumber = c.GetString("tun_num")
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

	isSetuped := models.IsVPNConfigured()
	if isSetuped {
		c.TplName = "wizard/step_no.html"
		return
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

	isSetuped := models.IsVPNConfigured()
	if isSetuped {
		c.TplName = "wizard/step_no.html"
		return
	}

	var wizardByte = c.GetSession("ovpnWizardData").([]byte)
	ovpnWizardData := Decode(wizardByte)
	lib.Dump(ovpnWizardData)

	setupScript := GenerateEnvFile(ovpnWizardData)
	c.Data["EnvString"] = setupScript
	lib.RawSaveToFile("/tmp/setup.sh", setupScript)
}

func (c *WizardController) Setup() {
	// Create a new HTTP client with a custom timeout
	ctx, cancel := context.WithTimeout(context.Background(), 660*time.Second)
	defer cancel()

	// Run the Bash script within the context
	cmd := exec.CommandContext(ctx, "/bin/bash", "-c",
		fmt.Sprintf("cd /tmp/ && bash ./setup.sh"))
	output, err := cmd.CombinedOutput()
	if err != nil {
		// If there is an error, return 500 Internal Server Error with details
		c.CustomAbort(500, "Error running script: "+err.Error()+"\n"+string(output))
	} else {
		// If successful, return 200 OK with the script output
		c.Data["json"] = map[string]string{"output": string(output)}
		c.ServeJSON()
	}
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

func appendString(original string, addition string) string {
	return original + "\n" + addition
}

func GenerateEnvFile(ovpnWizardData models.OvpnServerBaseSetting) string {
	var envString = models.GetConstEnv()
	envString = appendString(envString, "export ENDPOINT=\""+ovpnWizardData.OvpnEndpoint+"\"")
	envString = appendString(envString, "export IP_RANGE=\""+ovpnWizardData.OvpnIPRange+"\"")
	envString = appendString(envString, "export PROTOCOL_CHOICE=\""+
		models.GetIndex(ovpnWizardData.OvpnProtocol)+"\" #"+ovpnWizardData.OvpnProtocol)
	envString = appendString(envString, "export PORT=\""+ovpnWizardData.OvpnPort+"\"")
	envString = appendString(envString, "export TUN_NUMBER=\""+ovpnWizardData.TunNumber+"\"")
	envString = appendString(envString, "export DNS1=\""+ovpnWizardData.OvpnDNS1+"\"")
	envString = appendString(envString, "export DNS2=\""+ovpnWizardData.OvpnDNS2+"\"")
	if ovpnWizardData.OvpnCompression == "Disabled" {
		envString = appendString(envString, "export COMPRESSION_ENABLED=\"n\"")
	} else {
		envString = appendString(envString, "export COMPRESSION_ENABLED=\"y\"")
		envString = appendString(envString, "export COMPRESSION_CHOICE=\""+
			models.GetIndex(ovpnWizardData.OvpnCompression)+"\" #"+ovpnWizardData.OvpnCompression)
	}
	envString = appendString(envString, "export CIPHER_CHOICE=\""+
		models.GetIndex(ovpnWizardData.CipherChoice)+"\" #"+ovpnWizardData.CipherChoice)
	envString = appendString(envString, "export CERT_TYPE=\""+
		models.GetIndex(ovpnWizardData.CertType)+"\" #"+ovpnWizardData.CertType)
	envString = appendString(envString, "export CERT_CURVE_CHOICE=\""+
		models.GetIndex(ovpnWizardData.CertCurve)+"\" #"+ovpnWizardData.CertCurve)
	envString = appendString(envString, "export RSA_KEY_SIZE_CHOICE=\""+
		models.GetIndex(ovpnWizardData.RSAKeySize)+"\" #"+ovpnWizardData.RSAKeySize)
	envString = appendString(envString, "export CC_CIPHER_CHOICE=\""+
		models.GetIndex(ovpnWizardData.CCCipherChoice)+"\" #"+ovpnWizardData.CCCipherChoice)
	envString = appendString(envString, "export DH_TYPE=\""+
		models.GetIndex(ovpnWizardData.DHType)+"\" #"+ovpnWizardData.DHType)
	envString = appendString(envString, "export DH_CURVE_CHOICE=\""+
		models.GetIndex(ovpnWizardData.DHCurve)+"\" #"+ovpnWizardData.DHCurve)
	envString = appendString(envString, "export DH_KEY_SIZE_CHOICE=\""+
		models.GetIndex(ovpnWizardData.DHKeySize)+"\" #"+ovpnWizardData.DHKeySize)
	envString = appendString(envString, "export HMAC_ALG_CHOICE=\""+
		models.GetIndex(ovpnWizardData.HMACAlgorithm)+"\" #"+ovpnWizardData.HMACAlgorithm)
	envString = appendString(envString, "export TLS_SIG=\""+
		models.GetIndex(ovpnWizardData.TLSsig)+"\" #"+ovpnWizardData.TLSsig)

	if ovpnWizardData.DisableDefRouteForClientsByDefault {
		envString = appendString(envString, "export DISABLE_DEF_ROUTE_FOR_CLIENTS=\"y\"")
	} else {
		envString = appendString(envString, "export DISABLE_DEF_ROUTE_FOR_CLIENTS=\"n\"")
	}

	if ovpnWizardData.ClientToClientConfigIsUsed {
		envString = appendString(envString, "export CLIENT_TO_CLIENT=\"y\"")
	} else {
		envString = appendString(envString, "export CLIENT_TO_CLIENT=\"n\"")
	}

	envString = appendString(envString, models.GetScriptEnv())

	return envString
}
