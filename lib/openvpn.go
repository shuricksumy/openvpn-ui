package lib

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/beego/beego/v2/core/logs"
	"github.com/shuricksumy/openvpn-ui/models"
)

var (
	OpenVPNProcessID int
	ProcessMutex     sync.Mutex
)

func IsOpenVPNRunning() bool {
	pid, err := GetOpenVPNProcessIDFromPS()
	if err != nil {
		fmt.Println(err)
		return false
	}

	return pid != 0
}

func GetOpenVPNStatus() string {
	pid, err := GetOpenVPNProcessIDFromPS()
	if err != nil {
		logs.Error(err)
		return fmt.Sprintf("Error getting OpenVPN status: %s", err)
	}

	ProcessMutex.Lock()
	defer ProcessMutex.Unlock()

	if pid != 0 {
		OpenVPNProcessID = pid
		return "OpenVPN is running"
	}

	OpenVPNProcessID = 0
	return "OpenVPN is stopped"
}

func StartOpenVPN() error {
	if IsOpenVPNRunning() {
		return nil
	}

	cmd := exec.Command("/usr/sbin/openvpn", "--daemon", "openvpnserver", "--cd", "/etc/openvpn", "--config", "/etc/openvpn/server.conf")
	err := cmd.Run()
	if err != nil {
		logs.Error(fmt.Sprintf("Failed to start OpenVPN: %s", err))
		logs.Error(err)
		return err
	}

	// Get the PID of the newly started OpenVPN process
	pid, err := GetOpenVPNProcessIDFromPS()
	if err != nil {
		logs.Error(fmt.Sprintf("Error getting OpenVPN PID: %s", err))
		logs.Error(err)
		return err
	}

	ProcessMutex.Lock()
	OpenVPNProcessID = pid
	ProcessMutex.Unlock()

	logs.Warn("OpenVPN is running now")
	return nil
}

func StopOpenVPN() error {
	if !IsOpenVPNRunning() {
		logs.Error("OpenVPN is stopped")
		return nil
	}

	cmd := exec.Command("kill", fmt.Sprintf("%d", OpenVPNProcessID))
	err := cmd.Run()
	if err != nil {
		logs.Error(err)
		return fmt.Errorf("Error stopping OpenVPN: %s", err)
	}

	ProcessMutex.Lock()
	OpenVPNProcessID = 0
	ProcessMutex.Unlock()

	return nil
}

func GetOpenVPNProcessIDFromPS() (int, error) {
	cmd := exec.Command("bash", "-c", "/bin/ps -ef | /bin/grep openvpnserver | /bin/grep -v grep")
	// logs.Error("CMD:", cmd)

	output, err := cmd.CombinedOutput()
	// logs.Error("OUTPUT:", output)
	if err != nil {
		logs.Error(err)
		logs.Error(output)
		return 0, fmt.Errorf("Error getting OpenVPN PID: %s\n%s", err, output)
	}

	// Extract the PID from the output
	pid, err := ExtractPIDFromPSOutput(string(output))
	if err != nil {
		logs.Error(err)
		logs.Error(output)
		return 0, fmt.Errorf("Error extracting OpenVPN PID: %s\n%s", err, output)
	}

	return pid, nil
}

func ExtractPIDFromPSOutput(output string) (int, error) {
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) > 1 {
			pid := fields[1]
			return strconv.Atoi(pid)
		}
	}
	logs.Error("PID not found in PS output")
	return 0, fmt.Errorf("PID not found in PS output")
}

func EnableFWRules() error {
	//bash /etc/openvpn/set_fw.sh
	cmd := exec.Command("/bin/bash", "-c", "/bin/bash /etc/openvpn/set_fw.sh")
	err := cmd.Run()
	if err != nil {
		logs.Error("SET FW RULES:", err)
		return fmt.Errorf("SET FW RULES: %s", err)
	}

	return nil
}

func DisableFWRules() error {
	//bash /etc/openvpn/rm_fw.sh
	cmd := exec.Command("/bin/bash", "-c", "/bin/bash /etc/openvpn/rm_fw.sh")
	err := cmd.Run()
	if err != nil {
		logs.Error(err)
		return fmt.Errorf("DEL FW RULES: %s", err)
	}

	return nil
}

func RestartOpenVpnUI() (string, error) {

	errMsg := ""

	// Stop the process
	err := StopOpenVPN()
	if err != nil {
		logs.Error(err)
		errMsg = errMsg + "\n" + fmt.Sprintf("Error stopping OpenVPN: %s", err)
	}

	err_fw := DisableFWRules()
	if err_fw != nil {
		// c.Ctx.WriteString(fmt.Sprintf("Error deleting FireWall rules: %s", err_fw))
		logs.Error(err_fw)
		errMsg = errMsg + "\n" + fmt.Sprintf("Error deleting FireWall rules: %s", err_fw)
	}

	// Calling Sleep method
	time.Sleep(3 * time.Second)

	// Start the process again
	err = StartOpenVPN()
	if err != nil {
		logs.Error(err)
		errMsg = errMsg + "\n" + fmt.Sprintf("Error starting OpenVPN: %s", err)
	}

	err_fw = EnableFWRules()
	if err_fw != nil {
		// c.Ctx.WriteString(fmt.Sprintf("Error apply FireWall rules: %s", err_fw))
		logs.Error(err_fw)
		errMsg = errMsg + "\n" + fmt.Sprintf("Error apply FireWall rules: %s", err_fw)
	}

	if errMsg != "" {
		return errMsg, errors.New("Failed")
	}

	return "OpenVPN server has been restarted", nil
}

func RawReadClientFile(clientName string) (string, error) {
	//InitGlobalVars()
	destPathClientConfig := filepath.Join(CCD_DIR_PATH, clientName)
	return RawReadFile(destPathClientConfig)
}

func ApplyClientsConfigToFS() error {
	//InitGlobalVars()
	serverConf, err_read := RawReadFile(SERVER_CONFIG_PATH)
	if err_read != nil {
		logs.Error("Issue with reading config server file: ", SERVER_CONFIG_PATH)
		return err_read
	}

	isTopologyNet30, _ := regexp.MatchString(`\ntopology net30\s*\n`, serverConf)
	isTopologySubnet, _ := regexp.MatchString(`\ntopology subnet\s*\n`, serverConf)

	clientDetails, err_read_cl := models.GetAllClientsWithCertificate()
	if err_read_cl != nil {
		logs.Error("Issue with reading DB: ", err_read_cl)
		return err_read_cl
	}

	for _, client := range clientDetails {
		var buffer bytes.Buffer

		// 1. add init line
		buffer.WriteString("### Automatic generated \"" + client.ClientName + "\" settings file - " +
			time.Now().Format(time.RFC850) + "\n")

		// 2. add static IP
		if client.StaticIP != nil {
			staticIp := client.StaticIP
			staticIpStr := NilStringToString(staticIp)
			if _isIPAddressValid(staticIpStr) {
				var nextIP string

				if isTopologyNet30 {
					nextIP = _getNextIPAddress(staticIpStr) // --topology net30
				}

				if isTopologySubnet {
					nextIP = "255.255.255.0" // --topology subnet
				}

				if nextIP != "" {
					buffer.WriteString("ifconfig-push " + staticIpStr + " " + nextIP + "\n")
				} else {
					logs.Error("Cannot get nextIP or subnet topology is undefined in openvpn config file. Use as topology net30")
					nextIP = _getNextIPAddress(staticIpStr) // --topology net30
					buffer.WriteString("ifconfig-push " + staticIpStr + " " + nextIP + "\n")
				}
			}
		}
		// 3. if router add route to itself
		if client.IsRouter {
			routerRoutes, _ := models.GetAllRoutesProvided(client.Id)
			for _, r := range routerRoutes {
				if _isIPAddressValid(r.RouteIP) && _isIPAddressValid(r.RouteMask) {
					buffer.WriteString("iroute " + r.RouteIP + " " + r.RouteMask + "\n")
				}
			}
		}

		// 4. default router
		buffer.WriteString("\n")

		if client.IsRouteDefault {
			buffer.WriteString("# Set VPN as default route\n")
			buffer.WriteString("push \"redirect-gateway def1\"\n")
		} else {
			buffer.WriteString("# Set VPN as default route\n")
			buffer.WriteString("# push \"redirect-gateway def1\"\n")
		}
		buffer.WriteString("\n")

		// 5. if Routes is configured
		for _, r := range client.Routes {
			if _isIPAddressValid(r.RouteIP) && _isIPAddressValid(r.RouteMask) {
				buffer.WriteString("# Route to " + r.Name + " [" + r.Description + "] device internal subnet\n")
				buffer.WriteString("push \"route " + r.RouteIP + " " + r.RouteMask + "\"\n")
			}
		}

		//6. if Auth is used
		if client.OTPIsEnabled || client.StaticPassIsUsed {
			buffer.WriteString("\n#2FA OTP and/or static pass for auth.\n")
			buffer.WriteString("#2FA_KEY:" + NilStringToString(client.OTPKey) + "\n")
			buffer.WriteString("#2FA_USER:" + NilStringToString(client.OTPUserName) + "\n")
			buffer.WriteString("#STATIC_PASS:" + NilStringToString(client.StaticPass) + "\n")
		}

		//render file
		fileCCD := filepath.Join(CCD_DIR_PATH, client.ClientName)
		err_save := RawSaveToFile(fileCCD, buffer.String())
		if err_save != nil {
			logs.Error("Issue with saving client file: ", fileCCD)
			return err_save
		}
	}

	return nil
}

// FileData represents the structure with FileName and MD5sum fields.
type FileData struct {
	FileName string
	MD5sum   string
}

func UpdateDBWithLatestMD5() error {
	clients, err_cl := models.GetAllClientsWithCertificate()
	if err_cl != nil {
		return err_cl
	}

	for _, c := range clients {

		pathFile := filepath.Join(CCD_DIR_PATH, c.ClientName)

		md5Client, err_get_md5 := GetMD5SumFile(pathFile)
		if err_get_md5 != nil {
			return err_get_md5
		}

		err_upd := models.UpdateMD5SumForClientDetails(c.ClientName, md5Client)
		if err_upd != nil {
			logs.Error("Issue with updating MD5 [", c.ClientName, "] :", err_upd)
		}

	}
	return nil

}

func GetMD5StructureFromFS() map[string]bool {
	//InitGlobalVars()
	result := make(map[string]bool)

	clients, err_cl := models.GetAllClientsWithCertificate()
	if err_cl != nil {
		logs.Error("Issue with reading from DB:", err_cl)
	}

	for _, c := range clients {
		pathFile := filepath.Join(CCD_DIR_PATH, c.ClientName)
		md5Client, err_get_md5 := GetMD5SumFile(pathFile)
		if err_get_md5 != nil {
			logs.Error("Issue with getting MD5 from FS:", err_get_md5)
		}

		// Compare MD5 sums
		isMD5Valid := c.MD5Sum == md5Client
		// Dump(c.MD5Sum)
		// Dump(md5Client)

		// Add result to the map
		result[c.ClientName] = isMD5Valid

	}

	return result

}

func Watchdog() {
	for {
		time.Sleep(60 * time.Second) // Check every 60 seconds

		if !IsOpenVPNRunning() {
			logs.Error("OpenVPN is not running. Restarting...")
			err := StartOpenVPN()
			if err != nil {
				logs.Error(err)
			} else {
				logs.Error("OpenVPN restarted successfully.")
			}
		}
	}
}
