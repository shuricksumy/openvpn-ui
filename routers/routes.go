package routers

import "github.com/beego/beego/v2/server/web"

func init() {

	prefixURL, err := web.AppConfig.String("BaseURLPrefix")
	if err != nil {
		prefixURL = ""
	}

	web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:RoutesController"] =
		append(web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:RoutesController"],
			web.ControllerComments{
				Method:           "Get",
				Router:           prefixURL + `/routes`,
				AllowHTTPMethods: []string{"get"},
				Params:           nil})

	web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:RoutesController"] =
		append(web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:RoutesController"],
			web.ControllerComments{
				Method:           "NewRoute",
				Router:           prefixURL + `/routes/newroute`,
				AllowHTTPMethods: []string{"post"},
				Params:           nil})

	web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:RoutesController"] =
		append(web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:RoutesController"],
			web.ControllerComments{
				Method:           "GetRouteDetails",
				Router:           prefixURL + `/routes/get`,
				AllowHTTPMethods: []string{"post"},
				Params:           nil})

	web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:RoutesController"] =
		append(web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:RoutesController"],
			web.ControllerComments{
				Method:           "Post",
				Router:           prefixURL + `/routes`,
				AllowHTTPMethods: []string{"post"},
				Params:           nil})

	web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:RoutesController"] =
		append(web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:RoutesController"],
			web.ControllerComments{
				Method:           "Delete",
				Router:           prefixURL + `/routes/delete/:key`,
				AllowHTTPMethods: []string{"get"},
				Params:           nil})

	web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:RoutesController"] =
		append(web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:RoutesController"],
			web.ControllerComments{
				Method:           "UpdateFiles",
				Router:           prefixURL + `/routes/updatefiles`,
				AllowHTTPMethods: []string{"get"},
				Params:           nil})
}
