package routers

import "github.com/beego/beego/v2/server/web"

func init() {
	web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:RoutesController"] =
		append(web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:RoutesController"],
			web.ControllerComments{
				Method:           "Get",
				Router:           `/routes`,
				AllowHTTPMethods: []string{"get"},
				Params:           nil})

	web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:RoutesController"] =
		append(web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:RoutesController"],
			web.ControllerComments{
				Method:           "NewRoute",
				Router:           `/routes/newroute`,
				AllowHTTPMethods: []string{"post"},
				Params:           nil})

	web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:RoutesController"] =
		append(web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:RoutesController"],
			web.ControllerComments{
				Method:           "GetRouteDetails",
				Router:           `/routes/get/:key`,
				AllowHTTPMethods: []string{"get"},
				Params:           nil})

	web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:RoutesController"] =
		append(web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:RoutesController"],
			web.ControllerComments{
				Method:           "Post",
				Router:           `/routes`,
				AllowHTTPMethods: []string{"post"},
				Params:           nil})

	web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:RoutesController"] =
		append(web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:RoutesController"],
			web.ControllerComments{
				Method:           "Delete",
				Router:           `/routes/delete/:key`,
				AllowHTTPMethods: []string{"get"},
				Params:           nil})
}
