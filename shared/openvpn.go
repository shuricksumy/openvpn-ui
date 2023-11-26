package shared

import (
	"fmt"
	"os"
	"os/exec"
	"sync"
)

var (
	OpenVPNProcess   *exec.Cmd
	OpenVPNProcessID int
	OpenVPNStatus    string
	ProcessMutex     sync.Mutex
)

func IsOpenVPNRunning() bool {
	ProcessMutex.Lock()
	defer ProcessMutex.Unlock()

	return OpenVPNProcess != nil && OpenVPNProcess.ProcessState != nil && !OpenVPNProcess.ProcessState.Exited()
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
	cmd := exec.Command("/usr/sbin/openvpn", "--cd /etc/openvpn", "--config /etc/openvpn/server.conf")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		SetOpenVPNStatus(fmt.Sprintf("Failed to start OpenVPN: %s", err))
		return err
	}

	// Store the cmd object in the shared state for later use
	ProcessMutex.Lock()
	OpenVPNProcess = cmd
	ProcessMutex.Unlock()

	// Wait for the process to finish
	err = cmd.Wait()
	if err != nil {
		SetOpenVPNStatus(fmt.Sprintf("OpenVPN process exited with error: %s", err))
		return err
	}

	SetOpenVPNStatus("OpenVPN started successfully")
	return nil
}

func StopOpenVPN(pid int) error {
	cmd := exec.Command("kill", fmt.Sprintf("%d", pid))
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("Error stopping OpenVPN: %s", err)
	}

	ProcessMutex.Lock()
	OpenVPNProcess = nil
	ProcessMutex.Unlock()

	ProcessMutex.Lock()
	defer ProcessMutex.Unlock()

	// Check if the process state is available
	if OpenVPNProcess != nil && OpenVPNProcess.ProcessState != nil {
		// Check if the process has already completed
		if !OpenVPNProcess.ProcessState.Exited() {
			// Wait for the process to finish
			err := OpenVPNProcess.Wait()
			if err != nil {
				return fmt.Errorf("Error waiting for OpenVPN process to finish: %s", err)
			}
		}
	}

	SetOpenVPNStatus("OpenVPN stopped")
	return nil
}

func GetOpenVPNProcessID() int {
	ProcessMutex.Lock()
	defer ProcessMutex.Unlock()
	if OpenVPNProcess != nil && OpenVPNProcess.Process != nil {
		return OpenVPNProcess.Process.Pid
	}
	return 0
}
