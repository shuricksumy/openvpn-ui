package routers

import (
	"github.com/beego/beego/v2/server/web"
)

func init() {

	prefixURL, err := web.AppConfig.String("BaseURLPrefix")
	if err != nil {
		prefixURL = ""
	}

	web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:SystemController"] =
		append(web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:SystemController"],
			web.ControllerComments{
				Method:           "RestartLocalService",
				Router:           prefixURL + `/ov/system/restart`,
				AllowHTTPMethods: []string{"get"},
				Params:           nil})

	web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:SystemController"] =
		append(web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:SystemController"],
			web.ControllerComments{
				Method:           "Get",
				Router:           prefixURL + `/ov/system`,
				AllowHTTPMethods: []string{"get"},
				Params:           nil})

	web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:SystemController"] =
		append(web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:SystemController"],
			web.ControllerComments{
				Method:           "Backup",
				Router:           prefixURL + `/ov/system/backup`,
				AllowHTTPMethods: []string{"get"},
				Params:           nil})
}
