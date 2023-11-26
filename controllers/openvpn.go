package controllers

import (
	"fmt"

	"github.com/shuricksumy/openvpn-ui/shared"
)

type OpenVPNController struct {
	BaseController
}

func (c *OpenVPNController) StartOpenVPN() {
	// Start OpenVPN
	err := shared.StartOpenVPN()
	if err != nil {
		c.Ctx.WriteString(fmt.Sprintf("Error starting OpenVPN: %s", err))
		return
	}

	c.Ctx.WriteString("OpenVPN started")
}

func (c *OpenVPNController) StopOpenVPN() {
	// Stop OpenVPN
	err := shared.StopOpenVPN()
	if err != nil {
		c.Ctx.WriteString(fmt.Sprintf("Error stopping OpenVPN: %s", err))
		return
	}

	c.Ctx.WriteString("OpenVPN stopped")
}

func (c *OpenVPNController) GetOpenVPNStatus() {
	// Get OpenVPN status
	status := shared.GetOpenVPNStatus()
	c.Ctx.WriteString(fmt.Sprintf("Status: %s", status))
}

// Add any other controller functions related to OpenVPN control here
