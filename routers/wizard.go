package routers

import (
	"github.com/beego/beego/v2/server/web"
)

func init() {
	web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:WizardController"] =
		append(web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:WizardController"],
			web.ControllerComments{
				Method:           "Step1Get",
				Router:           `/wizard/step1`,
				AllowHTTPMethods: []string{"get"},
				Params:           nil})

	web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:WizardController"] =
		append(web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:WizardController"],
			web.ControllerComments{
				Method:           "Step1Post",
				Router:           `/wizard/step1`,
				AllowHTTPMethods: []string{"post"},
				Params:           nil})

	web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:WizardController"] =
		append(web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:WizardController"],
			web.ControllerComments{
				Method:           "Step2Get",
				Router:           `/wizard/step2`,
				AllowHTTPMethods: []string{"get"},
				Params:           nil})

	web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:WizardController"] =
		append(web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:WizardController"],
			web.ControllerComments{
				Method:           "Step2GetHmacAlg",
				Router:           `/wizard/step2/alg/:cipher/:selcted_hmac`,
				AllowHTTPMethods: []string{"get"},
				Params:           nil})

	web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:WizardController"] =
		append(web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:WizardController"],
			web.ControllerComments{
				Method:           "Step2GetCrtParam",
				Router:           `/wizard/step2/crtparam/:type/:selcted_option`,
				AllowHTTPMethods: []string{"get"},
				Params:           nil})

	web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:WizardController"] =
		append(web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:WizardController"],
			web.ControllerComments{
				Method:           "Step2GetCrtCipher",
				Router:           `/wizard/step2/crtcipher/:type/:selcted_option`,
				AllowHTTPMethods: []string{"get"},
				Params:           nil})

	web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:WizardController"] =
		append(web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:WizardController"],
			web.ControllerComments{
				Method:           "Step2GetDhParamr",
				Router:           `/wizard/step2/dhparam/:type/:selcted_option`,
				AllowHTTPMethods: []string{"get"},
				Params:           nil})

	web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:WizardController"] =
		append(web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:WizardController"],
			web.ControllerComments{
				Method:           "Step2Post",
				Router:           `/wizard/step2`,
				AllowHTTPMethods: []string{"post"},
				Params:           nil})

	web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:WizardController"] =
		append(web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:WizardController"],
			web.ControllerComments{
				Method:           "Step3Get",
				Router:           `/wizard/step3`,
				AllowHTTPMethods: []string{"get"},
				Params:           nil})

	web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:WizardController"] =
		append(web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:WizardController"],
			web.ControllerComments{
				Method:           "Setup",
				Router:           `/wizard/setup`,
				AllowHTTPMethods: []string{"get"},
				Params:           nil})
}
