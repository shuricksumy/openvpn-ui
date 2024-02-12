package routers

import "github.com/beego/beego/v2/server/web"

func init() {

	prefixURL, err := web.AppConfig.String("BaseURLPrefix")
	if err != nil {
		prefixURL = ""
	}

	web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:ClientsController"] =
		append(web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:ClientsController"],
			web.ControllerComments{
				Method:           "Get",
				Router:           prefixURL + `/clients`,
				AllowHTTPMethods: []string{"get"},
				Params:           nil})

	web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:ClientsController"] =
		append(web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:ClientsController"],
			web.ControllerComments{
				Method:           "NewClient",
				Router:           prefixURL + `/clients/newclient`,
				AllowHTTPMethods: []string{"post"},
				Params:           nil})

	web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:ClientsController"] =
		append(web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:ClientsController"],
			web.ControllerComments{
				Method:           "DelClient",
				Router:           prefixURL + `/clients/delclient/:key`,
				AllowHTTPMethods: []string{"get"},
				Params:           nil})

	web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:ClientsController"] =
		append(web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:ClientsController"],
			web.ControllerComments{
				Method:           "RenderModalRaw",
				Router:           prefixURL + `/clients/render_modal_raw`,
				AllowHTTPMethods: []string{"post"},
				Params:           nil})

	web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:ClientsController"] =
		append(web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:ClientsController"],
			web.ControllerComments{
				Method:           "RenderModal",
				Router:           prefixURL + `/clients/render_modal`,
				AllowHTTPMethods: []string{"post"},
				Params:           nil})

	web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:ClientsController"] =
		append(web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:ClientsController"],
			web.ControllerComments{
				Method:           "SaveClientDetailsData",
				Router:           prefixURL + `/clients/save_details_data`,
				AllowHTTPMethods: []string{"post"},
				Params:           nil})

	web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:ClientsController"] =
		append(web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:ClientsController"],
			web.ControllerComments{
				Method:           "UpdateFiles",
				Router:           prefixURL + `/clients/updatefiles`,
				AllowHTTPMethods: []string{"get"},
				Params:           nil})

	web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:ClientsController"] =
		append(web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:ClientsController"],
			web.ControllerComments{
				Method:           "SaveClientRawData",
				Router:           prefixURL + `/clients/save_client_data`,
				AllowHTTPMethods: []string{"post"},
				Params:           nil})

	web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:ClientsController"] =
		append(web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:ClientsController"],
			web.ControllerComments{
				Method:           "RenderAuthModal",
				Router:           prefixURL + `/clients/render_auth_modal/`,
				AllowHTTPMethods: []string{"post"},
				Params:           nil})

	web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:ClientsController"] =
		append(web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:ClientsController"],
			web.ControllerComments{
				Method:           "SaveClientAuthData",
				Router:           prefixURL + `/clients/save_auth_data/`,
				AllowHTTPMethods: []string{"post"},
				Params:           nil})

	web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:ClientsController"] =
		append(web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:ClientsController"],
			web.ControllerComments{
				Method:           "DeleteClientAuthData",
				Router:           prefixURL + `/clients/delete_auth_data/`,
				AllowHTTPMethods: []string{"post"},
				Params:           nil})

	web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:ClientsController"] =
		append(web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:ClientsController"],
			web.ControllerComments{
				Method:           "RenderModalClientRouting",
				Router:           prefixURL + `/clients/render_routing`,
				AllowHTTPMethods: []string{"post"},
				Params:           nil})

	web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:ClientsController"] =
		append(web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:ClientsController"],
			web.ControllerComments{
				Method:           "ResetCertificate",
				Router:           prefixURL + `/clients/reset_cert/:key`,
				AllowHTTPMethods: []string{"get"},
				Params:           nil})
}
