package routers

import (
	"github.com/beego/beego/v2/server/web"
	"github.com/shuricksumy/openvpn-ui/lib"
)

func init() {
	web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:LogsController"] =
		append(web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:LogsController"],
			web.ControllerComments{
				Method:           "RestartLocalService",
				Router:           lib.ParsePrefixURL() + `/logs/restart`,
				AllowHTTPMethods: []string{"get"},
				Params:           nil})
}
