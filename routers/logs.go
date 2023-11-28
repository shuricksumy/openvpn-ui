package routers

import (
	"github.com/beego/beego/v2/server/web"
)

func init() {
	web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:LogsController"] =
		append(web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:LogsController"],
			web.ControllerComments{
				Method:           "RestartLocalService",
				Router:           `/logs/restart`,
				AllowHTTPMethods: []string{"get"},
				Params:           nil})
}
