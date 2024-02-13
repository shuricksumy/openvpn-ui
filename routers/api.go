package routers

import (
	"github.com/beego/beego/v2/server/web"
)

func init() {
	web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:APISessionController"] =
		append(web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:APISessionController"],
			web.ControllerComments{
				Method:           "Get",
				Router:           `/`,
				AllowHTTPMethods: []string{"get"},
				Params:           nil})

	web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:APISessionController"] =
		append(web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:APISessionController"],
			web.ControllerComments{
				Method:           "Kill",
				Router:           `/`,
				AllowHTTPMethods: []string{"delete"},
				Params:           nil})

	web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:APISignalController"] =
		append(web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:APISignalController"],
			web.ControllerComments{
				Method:           "Send",
				Router:           `/`,
				AllowHTTPMethods: []string{"post"},
				Params:           nil})

	web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:APISysloadController"] =
		append(web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:APISysloadController"],
			web.ControllerComments{
				Method:           "Get",
				Router:           `/`,
				AllowHTTPMethods: []string{"get"},
				Params:           nil})
}
