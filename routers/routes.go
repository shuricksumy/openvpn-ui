package routers

import (
	"github.com/beego/beego/v2/server/web"
	"github.com/shuricksumy/openvpn-ui/lib"
)

func init() {
	web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:RoutesController"] =
		append(web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:RoutesController"],
			web.ControllerComments{
				Method:           "Get",
				Router:           lib.ParsePrefixURL() + `/routes`,
				AllowHTTPMethods: []string{"get"},
				Params:           nil})

	web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:RoutesController"] =
		append(web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:RoutesController"],
			web.ControllerComments{
				Method:           "NewRoute",
				Router:           lib.ParsePrefixURL() + `/routes/newroute`,
				AllowHTTPMethods: []string{"post"},
				Params:           nil})

	web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:RoutesController"] =
		append(web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:RoutesController"],
			web.ControllerComments{
				Method:           "GetRouteDetails",
				Router:           lib.ParsePrefixURL() + `/routes/get`,
				AllowHTTPMethods: []string{"post"},
				Params:           nil})

	web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:RoutesController"] =
		append(web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:RoutesController"],
			web.ControllerComments{
				Method:           "Post",
				Router:           lib.ParsePrefixURL() + `/routes`,
				AllowHTTPMethods: []string{"post"},
				Params:           nil})

	web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:RoutesController"] =
		append(web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:RoutesController"],
			web.ControllerComments{
				Method:           "Delete",
				Router:           lib.ParsePrefixURL() + `/routes/delete/:key`,
				AllowHTTPMethods: []string{"get"},
				Params:           nil})

	web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:RoutesController"] =
		append(web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:RoutesController"],
			web.ControllerComments{
				Method:           "UpdateFiles",
				Router:           lib.ParsePrefixURL() + `/routes/updatefiles`,
				AllowHTTPMethods: []string{"get"},
				Params:           nil})
}
