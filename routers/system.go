package routers

import (
	"github.com/beego/beego/v2/server/web"
	"github.com/shuricksumy/openvpn-ui/lib"
)

func init() {
	web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:SystemController"] =
		append(web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:SystemController"],
			web.ControllerComments{
				Method:           "RestartLocalService",
				Router:           lib.ParsePrefixURL() + `/ov/system/restart`,
				AllowHTTPMethods: []string{"get"},
				Params:           nil})

	web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:SystemController"] =
		append(web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:SystemController"],
			web.ControllerComments{
				Method:           "Get",
				Router:           lib.ParsePrefixURL() + `/ov/system`,
				AllowHTTPMethods: []string{"get"},
				Params:           nil})

	web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:SystemController"] =
		append(web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:SystemController"],
			web.ControllerComments{
				Method:           "Backup",
				Router:           lib.ParsePrefixURL() + `/ov/system/backup`,
				AllowHTTPMethods: []string{"get"},
				Params:           nil})
}
