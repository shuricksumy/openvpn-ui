package routers

import "github.com/beego/beego/v2/server/web"

func init() {
	web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:ClientsController"] =
		append(web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:ClientsController"],
			web.ControllerComments{
				Method:           "Get",
				Router:           `/clients`,
				AllowHTTPMethods: []string{"get"},
				Params:           nil})

	web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:ClientsController"] =
		append(web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:ClientsController"],
			web.ControllerComments{
				Method:           "NewClient",
				Router:           `/clients/newclient`,
				AllowHTTPMethods: []string{"post"},
				Params:           nil})

	web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:ClientsController"] =
		append(web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:ClientsController"],
			web.ControllerComments{
				Method:           "DelClient",
				Router:           `/clients/delclient/:key`,
				AllowHTTPMethods: []string{"get"},
				Params:           nil})

	web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:ClientsController"] =
		append(web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:ClientsController"],
			web.ControllerComments{
				Method:           "RenderModalRaw",
				Router:           `/clients/render_modal_raw`,
				AllowHTTPMethods: []string{"post"},
				Params:           nil})

	web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:ClientsController"] =
		append(web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:ClientsController"],
			web.ControllerComments{
				Method:           "RenderModal",
				Router:           `/clients/render_modal`,
				AllowHTTPMethods: []string{"post"},
				Params:           nil})

	web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:ClientsController"] =
		append(web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:ClientsController"],
			web.ControllerComments{
				Method:           "SaveClientDetailsData",
				Router:           `/clients/save_details_data`,
				AllowHTTPMethods: []string{"post"},
				Params:           nil})

	web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:ClientsController"] =
		append(web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:ClientsController"],
			web.ControllerComments{
				Method:           "UpdateFiles",
				Router:           `/clients/updatefiles`,
				AllowHTTPMethods: []string{"get"},
				Params:           nil})

	// web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:ClientsController"] =
	// 	append(web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:ClientsController"],
	// 		web.ControllerComments{
	// 			Method:           "Restart",
	// 			Router:           `/clients/restart`,
	// 			AllowHTTPMethods: []string{"get"},
	// 			Params:           nil})

	web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:ClientsController"] =
		append(web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:ClientsController"],
			web.ControllerComments{
				Method:           "SaveClientRawData",
				Router:           `/clients/save_client_data`,
				AllowHTTPMethods: []string{"post"},
				Params:           nil})

	web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:ClientsController"] =
		append(web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:ClientsController"],
			web.ControllerComments{
				Method:           "RenderAuthModal",
				Router:           `/clients/render_auth_modal/`,
				AllowHTTPMethods: []string{"post"},
				Params:           nil})

	web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:ClientsController"] =
		append(web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:ClientsController"],
			web.ControllerComments{
				Method:           "SaveClientAuthData",
				Router:           `/clients/save_auth_data/`,
				AllowHTTPMethods: []string{"post"},
				Params:           nil})

	web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:ClientsController"] =
		append(web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:ClientsController"],
			web.ControllerComments{
				Method:           "DeleteClientAuthData",
				Router:           `/clients/delete_auth_data/`,
				AllowHTTPMethods: []string{"post"},
				Params:           nil})

	web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:ClientsController"] =
		append(web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:ClientsController"],
			web.ControllerComments{
				Method:           "RenderModalClientRouting",
				Router:           `/clients/render_routing`,
				AllowHTTPMethods: []string{"post"},
				Params:           nil})
}
