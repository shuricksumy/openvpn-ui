package controllers

import (
	"github.com/beego/beego/v2/server/web"
	"github.com/shuricksumy/openvpn-ui/lib"
	"github.com/shuricksumy/openvpn-ui/models"
	"github.com/shuricksumy/openvpn-ui/state"
)

type BaseController struct {
	web.Controller

	Userinfo *models.User
	IsLogin  bool
}

type NestPreparer interface {
	NestPrepare()
}

type NestFinisher interface {
	NestFinish()
}

func (c *BaseController) Prepare() {
	c.SetParams()

	c.Data["SiteName"] = state.GlobalCfg.ServerName

	c.IsLogin = c.GetSession("userinfo") != nil
	if c.IsLogin {
		c.Userinfo = c.GetLogin()
		lib.InitGlobalVars()
	}

	c.Data["IsLogin"] = c.IsLogin
	c.Data["Userinfo"] = c.Userinfo

	if app, ok := c.AppController.(NestPreparer); ok {
		app.NestPrepare()
	}
}

func (c *BaseController) Finish() {
	if app, ok := c.AppController.(NestFinisher); ok {
		app.NestFinish()
	}
}

func (c *BaseController) GetLogin() *models.User {
	u := &models.User{Id: c.GetSession("userinfo").(string)}
	u.Read()
	return u
}

func (c *BaseController) DelLogin() {
	c.DelSession("userinfo")
}

func (c *BaseController) SetLogin(user *models.User) {
	c.SetSession("userinfo", user.Id)
}

func (c *BaseController) LoginPath() string {
	return c.URLFor("LoginController.Login")
}

func (c *BaseController) SetParams() {
	c.Data["Params"] = make(map[string]string)
	input, err := c.Input()
	if err != nil {
		// handle the error
	}
	for k, v := range input {
		c.Data["Params"].(map[string]string)[k] = v[0]
	}
	settings := models.Settings{Profile: "default"}
	_ = settings.Read("Profile")
	c.Data["Settings"] = &settings
}

type BreadCrumbs struct {
	Title    string
	Subtitle string
}
