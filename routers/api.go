package routers

import (
	"github.com/beego/beego/v2/server/web"
)

func init() {
	prefixURL, err := web.AppConfig.String("BaseURLPrefix")
	if err != nil {
		prefixURL = ""
	}

	web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:APISessionController"] =
		append(web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:APISessionController"],
			web.ControllerComments{
				Method:           "Get",
				Router:           prefixURL + `/`,
				AllowHTTPMethods: []string{"get"},
				Params:           nil})

	web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:APISessionController"] =
		append(web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:APISessionController"],
			web.ControllerComments{
				Method:           "Kill",
				Router:           prefixURL + `/`,
				AllowHTTPMethods: []string{"delete"},
				Params:           nil})

	web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:APISignalController"] =
		append(web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:APISignalController"],
			web.ControllerComments{
				Method:           "Send",
				Router:           prefixURL + `/`,
				AllowHTTPMethods: []string{"post"},
				Params:           nil})

	web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:APISysloadController"] =
		append(web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:APISysloadController"],
			web.ControllerComments{
				Method:           "Get",
				Router:           prefixURL + `/`,
				AllowHTTPMethods: []string{"get"},
				Params:           nil})
}
