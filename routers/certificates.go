package routers

import (
	"github.com/beego/beego/v2/server/web"
)

func init() {

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
}
