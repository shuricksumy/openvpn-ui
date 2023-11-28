package routers

import (
	"github.com/beego/beego/v2/server/web"
)

func init() {

	web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:SystemController"] =
		append(web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:SystemController"],
			web.ControllerComments{
				Method:           "RestartLocalService",
				Router:           `/ov/system/restart`,
				AllowHTTPMethods: []string{"get"},
				Params:           nil})

	web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:SystemController"] =
		append(web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:SystemController"],
			web.ControllerComments{
				Method:           "Get",
				Router:           `/ov/system`,
				AllowHTTPMethods: []string{"get"},
				Params:           nil})

	web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:SystemController"] =
		append(web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:SystemController"],
			web.ControllerComments{
				Method:           "Backup",
				Router:           `/ov/system/backup`,
				AllowHTTPMethods: []string{"get"},
				Params:           nil})
}
