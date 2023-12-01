package lib

import (
	"errors"
	"fmt"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/beego/beego/v2/core/logs"
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
