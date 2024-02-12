package routers

import (
	"github.com/beego/beego/v2/server/web"
)

func init() {

	prefixURL, err := web.AppConfig.String("BaseURLPrefix")
	if err != nil {
		prefixURL = ""
	}

	web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:WizardController"] =
		append(web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:WizardController"],
			web.ControllerComments{
				Method:           "Step1Get",
				Router:           prefixURL + `/wizard/step1`,
				AllowHTTPMethods: []string{"get"},
				Params:           nil})

	web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:WizardController"] =
		append(web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:WizardController"],
			web.ControllerComments{
				Method:           "Step1Post",
				Router:           prefixURL + `/wizard/step1`,
				AllowHTTPMethods: []string{"post"},
				Params:           nil})

	web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:WizardController"] =
		append(web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:WizardController"],
			web.ControllerComments{
				Method:           "Step2Get",
				Router:           prefixURL + `/wizard/step2`,
				AllowHTTPMethods: []string{"get"},
				Params:           nil})

	web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:WizardController"] =
		append(web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:WizardController"],
			web.ControllerComments{
				Method:           "Step2GetHmacAlg",
				Router:           prefixURL + `/wizard/step2/alg/:cipher/:selcted_hmac`,
				AllowHTTPMethods: []string{"get"},
				Params:           nil})

	web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:WizardController"] =
		append(web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:WizardController"],
			web.ControllerComments{
				Method:           "Step2GetCrtParam",
				Router:           prefixURL + `/wizard/step2/crtparam/:type/:selcted_option`,
				AllowHTTPMethods: []string{"get"},
				Params:           nil})

	web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:WizardController"] =
		append(web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:WizardController"],
			web.ControllerComments{
				Method:           "Step2GetCrtCipher",
				Router:           prefixURL + `/wizard/step2/crtcipher/:type/:selcted_option`,
				AllowHTTPMethods: []string{"get"},
				Params:           nil})

	web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:WizardController"] =
		append(web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:WizardController"],
			web.ControllerComments{
				Method:           "Step2GetDhParamr",
				Router:           prefixURL + `/wizard/step2/dhparam/:type/:selcted_option`,
				AllowHTTPMethods: []string{"get"},
				Params:           nil})

	web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:WizardController"] =
		append(web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:WizardController"],
			web.ControllerComments{
				Method:           "Step2Post",
				Router:           prefixURL + `/wizard/step2`,
				AllowHTTPMethods: []string{"post"},
				Params:           nil})

	web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:WizardController"] =
		append(web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:WizardController"],
			web.ControllerComments{
				Method:           "Step3Get",
				Router:           prefixURL + `/wizard/step3`,
				AllowHTTPMethods: []string{"get"},
				Params:           nil})

	web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:WizardController"] =
		append(web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:WizardController"],
			web.ControllerComments{
				Method:           "Setup",
				Router:           prefixURL + `/wizard/setup`,
				AllowHTTPMethods: []string{"get"},
				Params:           nil})
}
