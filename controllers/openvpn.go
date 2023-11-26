package controllers

import (
	"fmt"

	"github.com/shuricksumy/openvpn-ui/shared"
)

type OpenVPNController struct {
	BaseController
}

func (c *OpenVPNController) Start() {
	// Start OpenVPN
	err := shared.StartOpenVPN()
	if !c.IsLogin {
		c.Ctx.Redirect(302, c.LoginPath())
		return
	}

	if err != nil {
		c.Ctx.WriteString(fmt.Sprintf("Error starting OpenVPN: %s", err))
		return
	}

	c.Ctx.WriteString("OpenVPN started")
}

func (c *OpenVPNController) Status() {
	if !c.IsLogin {
		c.Ctx.Redirect(302, c.LoginPath())
		return
	}

	status := shared.GetOpenVPNStatus()
	c.Ctx.WriteString(fmt.Sprintf("Status: %s", status))
}

func (c *OpenVPNController) Stop() {
	if !c.IsLogin {
		c.Ctx.Redirect(302, c.LoginPath())
		return
	}

	pid := shared.GetOpenVPNProcessID()
	if pid == 0 {
		c.Ctx.WriteString("OpenVPN is not running")
		return
	}

	err := shared.StopOpenVPN(pid)
	if err != nil {
		c.Ctx.WriteString(fmt.Sprintf("Error stopping OpenVPN: %s", err))
		return
	}

	c.Ctx.WriteString("OpenVPN stopped")
}

// Add any other controller functions related to OpenVPN control here
