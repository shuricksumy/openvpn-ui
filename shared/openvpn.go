package shared

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"sync"
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
		return err
	}

	// Get the PID of the newly started OpenVPN process
	pid, err := GetOpenVPNProcessIDFromPS()
	if err != nil {
		SetOpenVPNStatus(fmt.Sprintf("Error getting OpenVPN PID: %s", err))
		return err
	}

	ProcessMutex.Lock()
	OpenVPNProcessID = pid
	ProcessMutex.Unlock()

	SetOpenVPNStatus("OpenVPN started successfully")
	return nil
}

func StopOpenVPN() error {
	if !IsOpenVPNRunning() {
		SetOpenVPNStatus("OpenVPN is not running")
		return nil
	}

	cmd := exec.Command("kill", fmt.Sprintf("%d", OpenVPNProcessID))
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("Error stopping OpenVPN: %s", err)
	}

	ProcessMutex.Lock()
	OpenVPNProcessID = 0
	ProcessMutex.Unlock()

	SetOpenVPNStatus("OpenVPN stopped")
	return nil
}

func GetOpenVPNProcessIDFromPS() (int, error) {
	cmd := exec.Command("ps", "-aux")
	grepCmd := exec.Command("grep", "openvpnserver")

	// Connect the output of ps to the input of grep
	grepCmd.Stdin, _ = cmd.StdoutPipe()

	// Capture the output of the grep command
	output, err := grepCmd.CombinedOutput()
	if err != nil {
		return 0, fmt.Errorf("Error getting OpenVPN PID: %s\n%s", err, output)
	}

	// Extract the PID from the output
	pid, err := ExtractPIDFromPSOutput(string(output))
	if err != nil {
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
	return 0, fmt.Errorf("PID not found in PS output")
}
