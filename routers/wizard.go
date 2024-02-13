package routers

import (
	"github.com/beego/beego/v2/server/web"
	"github.com/shuricksumy/openvpn-ui/lib"
)

func init() {
	web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:WizardController"] =
		append(web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:WizardController"],
			web.ControllerComments{
				Method:           "Step1Get",
				Router:           lib.ParsePrefixURL() + `/wizard/step1`,
				AllowHTTPMethods: []string{"get"},
				Params:           nil})

	web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:WizardController"] =
		append(web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:WizardController"],
			web.ControllerComments{
				Method:           "Step1Post",
				Router:           lib.ParsePrefixURL() + `/wizard/step1`,
				AllowHTTPMethods: []string{"post"},
				Params:           nil})

	web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:WizardController"] =
		append(web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:WizardController"],
			web.ControllerComments{
				Method:           "Step2Get",
				Router:           lib.ParsePrefixURL() + `/wizard/step2`,
				AllowHTTPMethods: []string{"get"},
				Params:           nil})

	web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:WizardController"] =
		append(web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:WizardController"],
			web.ControllerComments{
				Method:           "Step2GetHmacAlg",
				Router:           lib.ParsePrefixURL() + `/wizard/step2/alg/:cipher/:selcted_hmac`,
				AllowHTTPMethods: []string{"get"},
				Params:           nil})

	web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:WizardController"] =
		append(web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:WizardController"],
			web.ControllerComments{
				Method:           "Step2GetCrtParam",
				Router:           lib.ParsePrefixURL() + `/wizard/step2/crtparam/:type/:selcted_option`,
				AllowHTTPMethods: []string{"get"},
				Params:           nil})

	web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:WizardController"] =
		append(web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:WizardController"],
			web.ControllerComments{
				Method:           "Step2GetCrtCipher",
				Router:           lib.ParsePrefixURL() + `/wizard/step2/crtcipher/:type/:selcted_option`,
				AllowHTTPMethods: []string{"get"},
				Params:           nil})

	web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:WizardController"] =
		append(web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:WizardController"],
			web.ControllerComments{
				Method:           "Step2GetDhParamr",
				Router:           lib.ParsePrefixURL() + `/wizard/step2/dhparam/:type/:selcted_option`,
				AllowHTTPMethods: []string{"get"},
				Params:           nil})

	web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:WizardController"] =
		append(web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:WizardController"],
			web.ControllerComments{
				Method:           "Step2Post",
				Router:           lib.ParsePrefixURL() + `/wizard/step2`,
				AllowHTTPMethods: []string{"post"},
				Params:           nil})

	web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:WizardController"] =
		append(web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:WizardController"],
			web.ControllerComments{
				Method:           "Step3Get",
				Router:           lib.ParsePrefixURL() + `/wizard/step3`,
				AllowHTTPMethods: []string{"get"},
				Params:           nil})

	web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:WizardController"] =
		append(web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:WizardController"],
			web.ControllerComments{
				Method:           "Setup",
				Router:           lib.ParsePrefixURL() + `/wizard/setup`,
				AllowHTTPMethods: []string{"get"},
				Params:           nil})
}
