package routers

import (
	"github.com/beego/beego/v2/server/web"
)

func init() {
	web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:WizardController"] =
		append(web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:WizardController"],
			web.ControllerComments{
				Method:           "Get",
				Router:           `/`,
				AllowHTTPMethods: []string{"get"},
				Params:           nil})
}
