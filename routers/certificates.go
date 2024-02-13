package routers

import (
	"github.com/beego/beego/v2/server/web"
	"github.com/shuricksumy/openvpn-ui/lib"
)

func init() {
	web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:CertificatesController"] =
		append(web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:CertificatesController"],
			web.ControllerComments{
				Method:           "Download",
				Router:           lib.ParsePrefixURL() + `/certificates/ovpn/:key`,
				AllowHTTPMethods: []string{"get"},
				Params:           nil})

	web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:CertificatesController"] =
		append(web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:CertificatesController"],
			web.ControllerComments{
				Method:           "Get",
				Router:           lib.ParsePrefixURL() + `/certificates`,
				AllowHTTPMethods: []string{"get"},
				Params:           nil})

	web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:CertificatesController"] =
		append(web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:CertificatesController"],
			web.ControllerComments{
				Method:           "Post",
				Router:           lib.ParsePrefixURL() + `/certificates`,
				AllowHTTPMethods: []string{"post"},
				Params:           nil})

	web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:CertificatesController"] =
		append(web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:CertificatesController"],
			web.ControllerComments{
				Method:           "Revoke",
				Router:           lib.ParsePrefixURL() + `/certificates/revoke/:key`,
				AllowHTTPMethods: []string{"get"},
				Params:           nil})

	web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:CertificatesController"] =
		append(web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:CertificatesController"],
			web.ControllerComments{
				Method:           "UnRevoke",
				Router:           lib.ParsePrefixURL() + `/certificates/unrevoke/:key`,
				AllowHTTPMethods: []string{"get"},
				Params:           nil})

	web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:CertificatesController"] =
		append(web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:CertificatesController"],
			web.ControllerComments{
				Method:           "Restart",
				Router:           lib.ParsePrefixURL() + `/certificates/restart`,
				AllowHTTPMethods: []string{"get"},
				Params:           nil})

	web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:CertificatesController"] =
		append(web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:CertificatesController"],
			web.ControllerComments{
				Method:           "Burn",
				Router:           lib.ParsePrefixURL() + `/certificates/burn/:key/:serial`,
				AllowHTTPMethods: []string{"get"},
				Params:           nil})

	web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:CertificatesController"] =
		append(web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:CertificatesController"],
			web.ControllerComments{
				Method:           "Renew",
				Router:           lib.ParsePrefixURL() + `/certificates/renew/:key/:serial`,
				AllowHTTPMethods: []string{"get"},
				Params:           nil})

	web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:CertificatesController"] =
		append(web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:CertificatesController"],
			web.ControllerComments{
				Method:           "SaveClientRawData",
				Router:           lib.ParsePrefixURL() + `/certificates/save_client_data`,
				AllowHTTPMethods: []string{"post"},
				Params:           nil})

	web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:CertificatesController"] =
		append(web.GlobalControllerRouter["github.com/shuricksumy/openvpn-ui/controllers:CertificatesController"],
			web.ControllerComments{
				Method:           "UpdateFiles",
				Router:           lib.ParsePrefixURL() + `/certificates/updatefiles`,
				AllowHTTPMethods: []string{"get"},
				Params:           nil})
}
