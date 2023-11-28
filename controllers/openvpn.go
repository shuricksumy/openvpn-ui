package controllers

import (
	"fmt"
	"time"

	"github.com/shuricksumy/openvpn-ui/lib"
)

type OpenVPNController struct {
	BaseController
}

func (c *OpenVPNController) StartOpenVPN() {
	if !c.IsLogin {
		c.Ctx.Redirect(302, c.LoginPath())
		return
	}
	// Start OpenVPN
	err := lib.StartOpenVPN()
	if err != nil {
		c.Ctx.WriteString(fmt.Sprintf("Error starting OpenVPN: %s", err))
		return
	}

	err_fw := lib.EnableFWRules()
	if err_fw != nil {
		c.Ctx.WriteString(fmt.Sprintf("Error apply FireWall rules: %s", err_fw))
		return
	}

	c.Ctx.WriteString("OpenVPN started")
}

func (c *OpenVPNController) StopOpenVPN() {
	if !c.IsLogin {
		c.Ctx.Redirect(302, c.LoginPath())
		return
	}
	// Stop OpenVPN
	err := lib.StopOpenVPN()
	if err != nil {
		c.Ctx.WriteString(fmt.Sprintf("Error stopping OpenVPN: %s", err))
		return
	}

	err_fw := lib.DisableFWRules()
	if err_fw != nil {
		c.Ctx.WriteString(fmt.Sprintf("Error deleting FireWall rules: %s", err_fw))
	}

	c.Ctx.WriteString("OpenVPN stopped")
}

func (c *OpenVPNController) RestartOpenVPN() {
	if !c.IsLogin {
		c.Ctx.Redirect(302, c.LoginPath())
		return
	}
	// Stop the process
	err := lib.StopOpenVPN()
	if err != nil {
		c.Ctx.WriteString(fmt.Sprintf("Error stopping OpenVPN: %s", err))
		return
	}

	err_fw := lib.DisableFWRules()
	if err_fw != nil {
		c.Ctx.WriteString(fmt.Sprintf("Error deleting FireWall rules: %s", err_fw))
	}

	// Calling Sleep method
	time.Sleep(3 * time.Second)

	// Start the process again
	err = lib.StartOpenVPN()
	if err != nil {
		c.Ctx.WriteString(fmt.Sprintf("Error starting OpenVPN: %s", err))
		return
	}

	err_fw = lib.EnableFWRules()
	if err_fw != nil {
		c.Ctx.WriteString(fmt.Sprintf("Error apply FireWall rules: %s", err_fw))
		return
	}

	c.Ctx.Output.Body([]byte("Process restarted."))
}

func (c *OpenVPNController) GetOpenVPNStatus() {
	if !c.IsLogin {
		c.Ctx.Redirect(302, c.LoginPath())
		return
	}
	// Get OpenVPN status
	status := lib.GetOpenVPNStatus()
	c.Ctx.WriteString(fmt.Sprintf("Status: %s", status))
}

// Add any other controller functions related to OpenVPN control here
