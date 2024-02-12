package routers

import (
	"github.com/beego/beego/v2/server/web"
)

func init() {

	prefixURL, err := web.AppConfig.String("BaseURLPrefix")
	if err != nil {
		prefixURL = ""
	}

	web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:LogsController"] =
		append(web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:LogsController"],
			web.ControllerComments{
				Method:           "RestartLocalService",
				Router:           prefixURL + `/logs/restart`,
				AllowHTTPMethods: []string{"get"},
				Params:           nil})
}
