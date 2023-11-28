package routers

import (
	"github.com/beego/beego/v2/server/web"
)

func init() {

	web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:APISessionController"] =
		append(web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:APISessionController"],
			web.ControllerComments{
				Method:           "Get",
				Router:           `/`,
				AllowHTTPMethods: []string{"get"},
				Params:           nil})

	web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:APISessionController"] =
		append(web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:APISessionController"],
			web.ControllerComments{
				Method:           "Kill",
				Router:           `/`,
				AllowHTTPMethods: []string{"delete"},
				Params:           nil})

	web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:APISignalController"] =
		append(web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:APISignalController"],
			web.ControllerComments{
				Method:           "Send",
				Router:           `/`,
				AllowHTTPMethods: []string{"post"},
				Params:           nil})

	web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:APISysloadController"] =
		append(web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:APISysloadController"],
			web.ControllerComments{
				Method:           "Get",
				Router:           `/`,
				AllowHTTPMethods: []string{"get"},
				Params:           nil})

	web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:CertificatesController"] =
		append(web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:CertificatesController"],
			web.ControllerComments{
				Method:           "Download",
				Router:           `/certificates/ovpn/:key`,
				AllowHTTPMethods: []string{"get"},
				Params:           nil})

	web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:CertificatesController"] =
		append(web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:CertificatesController"],
			web.ControllerComments{
				Method:           "Get",
				Router:           `/certificates`,
				AllowHTTPMethods: []string{"get"},
				Params:           nil})

	web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:CertificatesController"] =
		append(web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:CertificatesController"],
			web.ControllerComments{
				Method:           "Post",
				Router:           `/certificates`,
				AllowHTTPMethods: []string{"post"},
				Params:           nil})

	web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:CertificatesController"] =
		append(web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:CertificatesController"],
			web.ControllerComments{
				Method:           "Revoke",
				Router:           `/certificates/revoke/:key`,
				AllowHTTPMethods: []string{"get"},
				Params:           nil})

	web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:CertificatesController"] =
		append(web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:CertificatesController"],
			web.ControllerComments{
				Method:           "UnRevoke",
				Router:           `/certificates/unrevoke/:key`,
				AllowHTTPMethods: []string{"get"},
				Params:           nil})

	web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:CertificatesController"] =
		append(web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:CertificatesController"],
			web.ControllerComments{
				Method:           "Restart",
				Router:           `/certificates/restart`,
				AllowHTTPMethods: []string{"get"},
				Params:           nil})

	web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:CertificatesController"] =
		append(web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:CertificatesController"],
			web.ControllerComments{
				Method:           "Burn",
				Router:           `/certificates/burn/:key/:serial`,
				AllowHTTPMethods: []string{"get"},
				Params:           nil})

	web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:CertificatesController"] =
		append(web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:CertificatesController"],
			web.ControllerComments{
				Method:           "Renew",
				Router:           `/certificates/renew/:key/:serial`,
				AllowHTTPMethods: []string{"get"},
				Params:           nil})

	web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:CertificatesController"] =
		append(web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:CertificatesController"],
			web.ControllerComments{
				Method:           "SaveClientRawData",
				Router:           `/certificates/save_client_data`,
				AllowHTTPMethods: []string{"post"},
				Params:           nil})

	web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:SystemController"] =
		append(web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:SystemController"],
			web.ControllerComments{
				Method:           "RestartLocalService",
				Router:           `/ov/system/restart`,
				AllowHTTPMethods: []string{"get"},
				Params:           nil})

	web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:SystemController"] =
		append(web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:SystemController"],
			web.ControllerComments{
				Method:           "Get",
				Router:           `/ov/system`,
				AllowHTTPMethods: []string{"get"},
				Params:           nil})

	web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:SystemController"] =
		append(web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:SystemController"],
			web.ControllerComments{
				Method:           "Backup",
				Router:           `/ov/system/backup`,
				AllowHTTPMethods: []string{"get"},
				Params:           nil})

	// web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:ClientsController"] =
	// 	append(web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:ClientsController"],
	// 		web.ControllerComments{
	// 			Method:           "Get",
	// 			Router:           `/clients`,
	// 			AllowHTTPMethods: []string{"get"},
	// 			Params:           nil})

	// web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:ClientsController"] =
	// 	append(web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:ClientsController"],
	// 		web.ControllerComments{
	// 			Method:           "RenderModalRaw",
	// 			Router:           `/clients/render_modal_raw`,
	// 			AllowHTTPMethods: []string{"post"},
	// 			Params:           nil})

	// web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:ClientsController"] =
	// 	append(web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:ClientsController"],
	// 		web.ControllerComments{
	// 			Method:           "RenderModal",
	// 			Router:           `/clients/render_modal`,
	// 			AllowHTTPMethods: []string{"post"},
	// 			Params:           nil})

	// web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:ClientsController"] =
	// 	append(web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:ClientsController"],
	// 		web.ControllerComments{
	// 			Method:           "SaveClientDetailsData",
	// 			Router:           `/clients/save_details_data`,
	// 			AllowHTTPMethods: []string{"post"},
	// 			Params:           nil})

	// web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:ClientsController"] =
	// 	append(web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:ClientsController"],
	// 		web.ControllerComments{
	// 			Method:           "UpdateFiles",
	// 			Router:           `/clients/updatefiles`,
	// 			AllowHTTPMethods: []string{"get"},
	// 			Params:           nil})

	// web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:ClientsController"] =
	// 	append(web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:ClientsController"],
	// 		web.ControllerComments{
	// 			Method:           "Restart",
	// 			Router:           `/clients/restart`,
	// 			AllowHTTPMethods: []string{"get"},
	// 			Params:           nil})

	// web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:ClientsController"] =
	// 	append(web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:ClientsController"],
	// 		web.ControllerComments{
	// 			Method:           "SaveClientRawData",
	// 			Router:           `/clients/save_client_data`,
	// 			AllowHTTPMethods: []string{"post"},
	// 			Params:           nil})

	// web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:RoutesController"] =
	// 	append(web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:RoutesController"],
	// 		web.ControllerComments{
	// 			Method:           "Get",
	// 			Router:           `/routes`,
	// 			AllowHTTPMethods: []string{"get"},
	// 			Params:           nil})

	// web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:RoutesController"] =
	// 	append(web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:RoutesController"],
	// 		web.ControllerComments{
	// 			Method:           "Post",
	// 			Router:           `/routes`,
	// 			AllowHTTPMethods: []string{"post"},
	// 			Params:           nil})

	// web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:RoutesController"] =
	// 	append(web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:RoutesController"],
	// 		web.ControllerComments{
	// 			Method:           "GetRoute",
	// 			Router:           `/routes/get/:key`,
	// 			AllowHTTPMethods: []string{"get"},
	// 			Params:           nil})

	// web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:RoutesController"] =
	// 	append(web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:RoutesController"],
	// 		web.ControllerComments{
	// 			Method:           "Delete",
	// 			Router:           `/routes/delete/:key`,
	// 			AllowHTTPMethods: []string{"get"},
	// 			Params:           nil})
}
