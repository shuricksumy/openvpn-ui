package controllers

import (
	"html/template"

	passlib "gopkg.in/hlandau/passlib.v1"

	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/core/validation"
	"github.com/beego/beego/v2/server/web"
	"github.com/shuricksumy/openvpn-ui/lib"
	"github.com/shuricksumy/openvpn-ui/models"
)

type ProfileController struct {
	BaseController
}

func (c *ProfileController) NestPrepare() {
	if !c.IsLogin {
		c.Ctx.Redirect(302, c.LoginPath())
		return
	}
	c.Data["breadcrumbs"] = &BreadCrumbs{
		Title: "Profile",
	}
}

func (c *ProfileController) Get() {
	if !c.IsLogin {
		c.Ctx.Redirect(302, c.LoginPath())
		return
	}

	c.Data["xsrfdata"] = template.HTML(c.XSRFFormHTML())
	c.Data["profile"] = c.Userinfo
	c.TplName = "profile.html"
}

func (c *ProfileController) Post() {
	if !c.IsLogin {
		c.Ctx.Redirect(302, c.LoginPath())
		return
	}

	c.TplName = "profile.html"
	c.Data["profile"] = c.Userinfo

	flash := web.NewFlash()

	user := models.User{}
	if err := c.ParseForm(&user); err != nil {
		logs.Error(err)
		flash.Error(err.Error())
		flash.Store(&c.Controller)
		return
	}
	user.Login = c.Userinfo.Login
	c.Data["profile"] = user

	if vMap := validateUser(user); vMap != nil {
		c.Data["validation"] = vMap
		return
	}

	hash, err := passlib.Hash(user.Password)
	if err != nil {
		flash.Error("Unable to hash password")
		flash.Store(&c.Controller)
		return
	}
	c.Userinfo.Email = user.Email
	c.Userinfo.Name = user.Name
	c.Userinfo.Password = hash
	o := orm.NewOrm()
	if _, err := o.Update(c.Userinfo); err != nil {
		flash.Error(err.Error())
	} else {
		flash.Success("Profile has been updated")
	}
	flash.Store(&c.Controller)
}

func validateUser(user models.User) map[string]map[string]string {
	valid := validation.Validation{}
	b, err := valid.Valid(&user)
	if err != nil {
		logs.Error(err)
		return nil
	}
	if !b {
		return lib.CreateValidationMap(valid)
	}
	return nil
}
