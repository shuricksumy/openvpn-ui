package routers

import (
	"github.com/beego/beego/v2/server/web"
	"github.com/shuricksumy/openvpn-ui/lib"
)

func init() {
	web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:ClientsController"] =
		append(web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:ClientsController"],
			web.ControllerComments{
				Method:           "Get",
				Router:           lib.ParsePrefixURL() + `/clients`,
				AllowHTTPMethods: []string{"get"},
				Params:           nil})

	web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:ClientsController"] =
		append(web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:ClientsController"],
			web.ControllerComments{
				Method:           "NewClient",
				Router:           lib.ParsePrefixURL() + `/clients/newclient`,
				AllowHTTPMethods: []string{"post"},
				Params:           nil})

	web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:ClientsController"] =
		append(web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:ClientsController"],
			web.ControllerComments{
				Method:           "DelClient",
				Router:           lib.ParsePrefixURL() + `/clients/delclient/:key`,
				AllowHTTPMethods: []string{"get"},
				Params:           nil})

	web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:ClientsController"] =
		append(web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:ClientsController"],
			web.ControllerComments{
				Method:           "RenderModalRaw",
				Router:           lib.ParsePrefixURL() + `/clients/render_modal_raw`,
				AllowHTTPMethods: []string{"post"},
				Params:           nil})

	web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:ClientsController"] =
		append(web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:ClientsController"],
			web.ControllerComments{
				Method:           "RenderModal",
				Router:           lib.ParsePrefixURL() + `/clients/render_modal`,
				AllowHTTPMethods: []string{"post"},
				Params:           nil})

	web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:ClientsController"] =
		append(web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:ClientsController"],
			web.ControllerComments{
				Method:           "SaveClientDetailsData",
				Router:           lib.ParsePrefixURL() + `/clients/save_details_data`,
				AllowHTTPMethods: []string{"post"},
				Params:           nil})

	web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:ClientsController"] =
		append(web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:ClientsController"],
			web.ControllerComments{
				Method:           "UpdateFiles",
				Router:           lib.ParsePrefixURL() + `/clients/updatefiles`,
				AllowHTTPMethods: []string{"get"},
				Params:           nil})

	web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:ClientsController"] =
		append(web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:ClientsController"],
			web.ControllerComments{
				Method:           "SaveClientRawData",
				Router:           lib.ParsePrefixURL() + `/clients/save_client_data`,
				AllowHTTPMethods: []string{"post"},
				Params:           nil})

	web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:ClientsController"] =
		append(web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:ClientsController"],
			web.ControllerComments{
				Method:           "RenderAuthModal",
				Router:           lib.ParsePrefixURL() + `/clients/render_auth_modal/`,
				AllowHTTPMethods: []string{"post"},
				Params:           nil})

	web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:ClientsController"] =
		append(web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:ClientsController"],
			web.ControllerComments{
				Method:           "SaveClientAuthData",
				Router:           lib.ParsePrefixURL() + `/clients/save_auth_data/`,
				AllowHTTPMethods: []string{"post"},
				Params:           nil})

	web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:ClientsController"] =
		append(web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:ClientsController"],
			web.ControllerComments{
				Method:           "DeleteClientAuthData",
				Router:           lib.ParsePrefixURL() + `/clients/delete_auth_data/`,
				AllowHTTPMethods: []string{"post"},
				Params:           nil})

	web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:ClientsController"] =
		append(web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:ClientsController"],
			web.ControllerComments{
				Method:           "RenderModalClientRouting",
				Router:           lib.ParsePrefixURL() + `/clients/render_routing`,
				AllowHTTPMethods: []string{"post"},
				Params:           nil})

	web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:ClientsController"] =
		append(web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:ClientsController"],
			web.ControllerComments{
				Method:           "ResetCertificate",
				Router:           lib.ParsePrefixURL() + `/clients/reset_cert/:key`,
				AllowHTTPMethods: []string{"get"},
				Params:           nil})
}
