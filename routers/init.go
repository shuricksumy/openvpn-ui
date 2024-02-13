// Package routers defines application routes
// @APIVersion 1.0.0
// @Title OpenVPN API
// @Description REST API allows you to control and monitor your OpenVPN server
// @Contact adam.walach@gmail.com
// License Apache 2.0
// LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
	"github.com/beego/beego/v2/server/web"
	"github.com/shuricksumy/openvpn-ui/controllers"
)

func Init(configDir string, prefixURL string) {

	web.Router(prefixURL+"/", &controllers.MainController{})
	web.SetStaticPath(prefixURL+"/static", "static")
	web.SetStaticPath(prefixURL+"/swagger", "swagger")

	web.Router(prefixURL+"/login", &controllers.LoginController{}, "get,post:Login")
	web.Router(prefixURL+"/logout", &controllers.LoginController{}, "get:Logout")
	web.Router(prefixURL+"/profile", &controllers.ProfileController{})
	web.Router(prefixURL+"/settings", &controllers.SettingsController{})
	web.Router(prefixURL+"/ov/serverconfig", &controllers.ServerConfigController{ConfigDir: configDir})
	web.Router(prefixURL+"/ov/clientconfig", &controllers.ClientConfigController{ConfigDir: configDir})
	web.Router(prefixURL+"/ov/easyrsa", &controllers.EasyRSAConfigController{ConfigDir: configDir})
	web.Router(prefixURL+"/logs", &controllers.LogsController{})

	web.Include(&controllers.CertificatesController{ConfigDir: configDir})
	web.Include(&controllers.ClientsController{ConfigDir: configDir})
	web.Include(&controllers.RoutesController{ConfigDir: configDir})

	web.Include(&controllers.WizardController{ConfigDir: configDir})
	web.Include(&controllers.LogsController{})
	web.Include(&controllers.SystemController{})

	web.Router(prefixURL+"/openvpn/start", &controllers.OpenVPNController{}, "get:StartOpenVPN")
	web.Router(prefixURL+"/openvpn/stop", &controllers.OpenVPNController{}, "get:StopOpenVPN")
	web.Router(prefixURL+"/openvpn/restart", &controllers.OpenVPNController{}, "get:RestartOpenVPN")
	web.Router(prefixURL+"/openvpn/status", &controllers.OpenVPNController{}, "get:GetOpenVPNStatus")

	ns := web.NewNamespace(prefixURL+"/api/v1",
		web.NSNamespace("/session",
			web.NSInclude(
				&controllers.APISessionController{},
			),
		),
		web.NSNamespace("/sysload",
			web.NSInclude(
				&controllers.APISysloadController{},
			),
		),
		web.NSNamespace("/signal",
			web.NSInclude(
				&controllers.APISignalController{},
			),
		),
	)
	web.AddNamespace(ns)
}
