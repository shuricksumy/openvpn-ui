package routers

import (
	"github.com/beego/beego/v2/server/web"
)

func init() {

	prefixURL, err := web.AppConfig.String("BaseURLPrefix")
	if err != nil {
		prefixURL = ""
	}

	web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:CertificatesController"] =
		append(web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:CertificatesController"],
			web.ControllerComments{
				Method:           "Download",
				Router:           prefixURL + `/certificates/ovpn/:key`,
				AllowHTTPMethods: []string{"get"},
				Params:           nil})

	web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:CertificatesController"] =
		append(web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:CertificatesController"],
			web.ControllerComments{
				Method:           "Get",
				Router:           prefixURL + `/certificates`,
				AllowHTTPMethods: []string{"get"},
				Params:           nil})

	web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:CertificatesController"] =
		append(web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:CertificatesController"],
			web.ControllerComments{
				Method:           "Post",
				Router:           prefixURL + `/certificates`,
				AllowHTTPMethods: []string{"post"},
				Params:           nil})

	web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:CertificatesController"] =
		append(web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:CertificatesController"],
			web.ControllerComments{
				Method:           "Revoke",
				Router:           prefixURL + `/certificates/revoke/:key`,
				AllowHTTPMethods: []string{"get"},
				Params:           nil})

	web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:CertificatesController"] =
		append(web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:CertificatesController"],
			web.ControllerComments{
				Method:           "UnRevoke",
				Router:           prefixURL + `/certificates/unrevoke/:key`,
				AllowHTTPMethods: []string{"get"},
				Params:           nil})

	web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:CertificatesController"] =
		append(web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:CertificatesController"],
			web.ControllerComments{
				Method:           "Restart",
				Router:           prefixURL + `/certificates/restart`,
				AllowHTTPMethods: []string{"get"},
				Params:           nil})

	web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:CertificatesController"] =
		append(web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:CertificatesController"],
			web.ControllerComments{
				Method:           "Burn",
				Router:           prefixURL + `/certificates/burn/:key/:serial`,
				AllowHTTPMethods: []string{"get"},
				Params:           nil})

	web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:CertificatesController"] =
		append(web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:CertificatesController"],
			web.ControllerComments{
				Method:           "Renew",
				Router:           prefixURL + `/certificates/renew/:key/:serial`,
				AllowHTTPMethods: []string{"get"},
				Params:           nil})

	web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:CertificatesController"] =
		append(web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:CertificatesController"],
			web.ControllerComments{
				Method:           "SaveClientRawData",
				Router:           prefixURL + `/certificates/save_client_data`,
				AllowHTTPMethods: []string{"post"},
				Params:           nil})

	web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:CertificatesController"] =
		append(web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:CertificatesController"],
			web.ControllerComments{
				Method:           "UpdateFiles",
				Router:           prefixURL + `/certificates/updatefiles`,
				AllowHTTPMethods: []string{"get"},
				Params:           nil})
}
