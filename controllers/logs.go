package controllers

import (
	"bufio"
	"os"
	"strings"

	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/server/web"
	"github.com/shuricksumy/openvpn-ui/lib"
	"github.com/shuricksumy/openvpn-ui/models"
)

type LogsController struct {
	BaseController
}

func (c *LogsController) NestPrepare() {
	if !c.IsLogin {
		c.Ctx.Redirect(302, c.LoginPath())
		return
	}
}

func ShowLogs(c *LogsController) {
	if !c.IsLogin {
		c.Ctx.Redirect(302, c.LoginPath())
		return
	}

	c.TplName = "logs.html"
	c.Data["breadcrumbs"] = &BreadCrumbs{
		Title: "Logs",
	}

	flash := web.NewFlash()

	settings := models.Settings{Profile: "default"}
	settings.Read("Profile")

	if err := settings.Read("OVConfigPath"); err != nil {
		logs.Error(err)
		return
	}

	fName := settings.OVConfigPath + "/openvpn.log"
	file, err := os.Open(fName)
	if err != nil {
		logs.Error(err)
		flash.Error(err.Error())
		flash.Store(&c.Controller)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	var logs []string
	for scanner.Scan() {
		line := scanner.Text()
		//	if strings.Index(line, " MANAGEMENT: ") == -1 {
		if !strings.Contains(line, " MANAGEMENT: ") {
			logs = append(logs, strings.Trim(line, "\t"))
		}
	}
	start := len(logs) - 500
	if start < 0 {
		start = 0
	}
	c.Data["logs"] = reverse(logs[start:])

}

func (c *LogsController) Get() {
	if !c.IsLogin {
		c.Ctx.Redirect(302, c.LoginPath())
		return
	}
	ShowLogs(c)
}

func reverse(lines []string) []string {
	for i := 0; i < len(lines)/2; i++ {
		j := len(lines) - i - 1
		lines[i], lines[j] = lines[j], lines[i]
	}
	return lines
}

func (c *LogsController) RestartLocalService() {
	if !c.IsLogin {
		c.Ctx.Redirect(302, c.LoginPath())
		return
	}

	flash := web.NewFlash()
	c.TplName = "logs.html"

	msg, err_r := lib.RestartOpenVpnUI()
	if err_r != nil {
		logs.Error(msg)
		flash.Error(msg)
		flash.Store(&c.Controller)
	} else {
		flash.Success(msg)
		flash.Store(&c.Controller)
	}

	ShowLogs(c)
}
