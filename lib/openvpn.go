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
	OpenVPNStatus    string
	ProcessMutex     sync.Mutex
)

func IsOpenVPNRunning() bool {
	ProcessMutex.Lock()
	defer ProcessMutex.Unlock()

	return OpenVPNProcessID != 0
}

func SetOpenVPNStatus(status string) {
	ProcessMutex.Lock()
	OpenVPNStatus = status
	ProcessMutex.Unlock()
}

func GetOpenVPNStatus() string {
	ProcessMutex.Lock()
	defer ProcessMutex.Unlock()
	return OpenVPNStatus
}

func StartOpenVPN() error {
	cmd := exec.Command("/usr/sbin/openvpn", "--daemon", "openvpnserver", "--cd", "/etc/openvpn", "--config", "/etc/openvpn/server.conf")
	err := cmd.Run()
	if err != nil {
		SetOpenVPNStatus(fmt.Sprintf("Failed to start OpenVPN: %s", err))
		logs.Error(err)
		return err
	}

	// Get the PID of the newly started OpenVPN process
	pid, err := GetOpenVPNProcessIDFromPS()
	if err != nil {
		SetOpenVPNStatus(fmt.Sprintf("Error getting OpenVPN PID: %s", err))
		logs.Error(err)
		return err
	}

	ProcessMutex.Lock()
	OpenVPNProcessID = pid
	ProcessMutex.Unlock()

	SetOpenVPNStatus("OpenVPN is running now")
	return nil
}

func StopOpenVPN() error {
	if !IsOpenVPNRunning() {
		SetOpenVPNStatus("OpenVPN is stopped")
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

	SetOpenVPNStatus("OpenVPN stopped")
	return nil
}

func GetOpenVPNProcessIDFromPS() (int, error) {
	cmd := exec.Command("bash", "-c", "/usr/bin/ps -ef | /usr/bin/grep openvpnserver | /usr/bin/grep -v grep")
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
	InitGlobalVars()
	destPathClientConfig := filepath.Join(CCD_DIR_PATH, clientName)
	return RawReadFile(destPathClientConfig)
}

func ApplyClientsConfigToFS() error {
	InitGlobalVars()
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

		// Debug results
		// logs.Error("==========START==============")
		// logs.Error(buffer.String())
		// logs.Error("==========END================")

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

// // createFileMap creates a map with FileName as the key and FileData as the value.
// func createFileMap(files ...FileData) map[string]FileData {
// 	fileMap := make(map[string]FileData)
// 	for _, file := range files {
// 		fileMap[file.FileName] = file
// 	}
// 	return fileMap
// }

// // calculateMD5Sum calculates the MD5 sum of a file content (dummy implementation for illustration).
// func calculateMD5Sum(content string) string {
// 	hasher := md5.New()
// 	hasher.Write([]byte(content))
// 	return hex.EncodeToString(hasher.Sum(nil))
// }

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
		// Add result to the map
		result[c.ClientName] = isMD5Valid

	}

	return result

}

// func UpdateDBWithLatestMD5() error {
// 	clientsDetails, err_read := models.GetAllClientsWithCertificate()
// 	if err_read != nil {
// 		return err_read
// 	}

// 	md5hashs := GetMD5StructureFromFS(clientsDetails)

// 	for _, c := range clientsDetails {

// 		// models.UpdateMD5SumForClientDetails()

// 		for _, h := range md5hashs {
// 			if c.ClientName == h.ClientName {
// 				c.MD5Sum = h.MD5Hash
// 			}
// 		}
// 	}

// 	return nil
// }
