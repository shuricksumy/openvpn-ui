package routers

import (
	"github.com/beego/beego/v2/server/web"
	"github.com/shuricksumy/openvpn-ui/lib"
)

func init() {
	web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:APISessionController"] =
		append(web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:APISessionController"],
			web.ControllerComments{
				Method:           "Get",
				Router:           lib.ParsePrefixURL() + `/`,
				AllowHTTPMethods: []string{"get"},
				Params:           nil})

	web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:APISessionController"] =
		append(web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:APISessionController"],
			web.ControllerComments{
				Method:           "Kill",
				Router:           lib.ParsePrefixURL() + `/`,
				AllowHTTPMethods: []string{"delete"},
				Params:           nil})

	web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:APISignalController"] =
		append(web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:APISignalController"],
			web.ControllerComments{
				Method:           "Send",
				Router:           lib.ParsePrefixURL() + `/`,
				AllowHTTPMethods: []string{"post"},
				Params:           nil})

	web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:APISysloadController"] =
		append(web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:APISysloadController"],
			web.ControllerComments{
				Method:           "Get",
				Router:           lib.ParsePrefixURL() + `/`,
				AllowHTTPMethods: []string{"get"},
				Params:           nil})
}
