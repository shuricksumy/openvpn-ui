package controllers

import (
	"fmt"
	"strings"

	"github.com/beego/beego/v2/core/logs"
	"github.com/shuricksumy/openvpn-ui/lib"
	"github.com/shuricksumy/openvpn-ui/models"
	mi "github.com/shuricksumy/openvpn-ui/pkg/openvpn-server-config/server/mi"
	"github.com/shuricksumy/openvpn-ui/state"
)

type MainController struct {
	BaseController
}

func (c *MainController) NestPrepare() {
	if !c.IsLogin {
		c.Ctx.Redirect(302, c.LoginPath())
		c.StopRun()
	}
	c.Data["breadcrumbs"] = &BreadCrumbs{
		Title: "",
	}
}

func (c *MainController) Get() {
	// Ensure trailing slash
	if !strings.HasSuffix(c.Ctx.Request.URL.Path, "/") {
		c.Ctx.Redirect(302, c.Ctx.Request.URL.Path+"/")
		c.StopRun()
		return
	}

	c.Data["sysinfo"] = lib.GetSystemInfo()
	//lib.Dump(lib.GetSystemInfo())
	client := mi.NewClient(state.GlobalCfg.MINetwork, state.GlobalCfg.MIAddress)
	status, err := client.GetStatus()
	if err != nil {
		logs.Error(err)
		logs.Warn(fmt.Sprintf("passed client line: %s", client))
		logs.Warn(fmt.Sprintf("error: %s", err))
	} else {
		c.Data["ovstatus"] = status
	}
	//lib.Dump(status)

	version, err := client.GetVersion()
	if err != nil {
		logs.Error(err)
	} else {
		c.Data["ovversion"] = version.OpenVPN
	}
	//lib.Dump(version)

	pid, err := client.GetPid()
	if err != nil {
		logs.Error(err)
	} else {
		c.Data["ovpid"] = pid
	}
	// lib.Dump(pid)

	loadStats, err := client.GetLoadStats()
	if err != nil {
		logs.Error(err)
	} else {
		c.Data["ovstats"] = loadStats
	}
	//lib.Dump(loadStats)

	c.Data["IsVPNSetup"] = models.IsVPNConfigured()

	c.TplName = "index.html"
}
