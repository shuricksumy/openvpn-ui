package controllers

import (
	"html/template"
	"time"

	"github.com/beego/beego/v2/server/web"
	"github.com/shuricksumy/openvpn-ui/lib"
)

type LoginController struct {
	BaseController
}

func (c *LoginController) Login() {
	if c.IsLogin {
		c.Ctx.Redirect(302, c.URLFor("MainController.Get"))
		return
	}

	c.TplName = "login.html"
	c.Data["xsrfdata"] = template.HTML(c.XSRFFormHTML())
	if !c.Ctx.Input.IsPost() {
		return
	}

	flash := web.NewFlash()
	login := c.GetString("login")
	password := c.GetString("password")

	authType, err := web.AppConfig.String("AuthType")
	if err != nil {
		flash.Warning(err.Error())
		flash.Store(&c.Controller)
		return
	}
	user, err := lib.Authenticate(login, password, authType)

	if err != nil {
		flash.Warning(err.Error())
		flash.Store(&c.Controller)
		return
	}
	user.Lastlogintime = time.Now()
	err = user.Update("Lastlogintime")
	if err != nil {
		flash.Warning(err.Error())
		flash.Store(&c.Controller)
		return
	}
	flash.Success("Success logged in")
	flash.Store(&c.Controller)

	c.SetLogin(user)

	c.Redirect(c.URLFor("MainController.Get"), 303)
}

func (c *LoginController) Logout() {
	c.DelLogin()
	flash := web.NewFlash()
	flash.Success("Success logged out")
	flash.Store(&c.Controller)

	c.Ctx.Redirect(302, c.URLFor("LoginController.Login"))
}
