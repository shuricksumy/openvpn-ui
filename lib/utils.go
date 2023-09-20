package lib

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"os"
	"os/exec"
	"path"
	"strings"
	"time"

	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/core/validation"
	"github.com/beego/beego/v2/server/web"
	"github.com/shuricksumy/openvpn-ui/state"
)

// CreateValidationMap ranslates validation structure to map
// that can be easly presented in template
func CreateValidationMap(valid validation.Validation) map[string]map[string]string {
	v := make(map[string]map[string]string)
	/*
			{
				"email": {
					"Requrired" : "Can not be empty"
				},
				"password" :{

			  }
		  }
	*/
	for _, err := range valid.Errors {
		logs.Notice(err.Key, err.Message)
		k := strings.Split(err.Key, ".")
		var field, errorType string
		if len(k) > 1 {
			field = k[0]
			errorType = k[1]
		} else {
			field = err.Key
			errorType = " "
		}
		logs.Error(field)
		if _, ok := v[field]; !ok {
			v[field] = make(map[string]string)
		}
		v[field][errorType] = err.Message
	}
	return v

}

// Dump any structure as json string
func Dump(obj interface{}) {
	result, _ := json.MarshalIndent(obj, "", "\t")
	logs.Debug(string(result))
}

// Make bac file
func BackupFile(destPath string) error {
	orig, err := os.ReadFile(destPath)
	if err != nil {
		return err
	}
	now := time.Now().Format("20060102150405")
	newDest := destPath + "_bac_" + now
	return os.WriteFile(newDest, []byte(orig), 0644)
}

// SaveToFile RAW
func RawSaveToFile(destPath string, text string) error {
	// Replace Windows-style EOL (CRLF) with Linux-style EOL (LF).
	text = strings.ReplaceAll(text, "\r\n", "\n")
	return os.WriteFile(destPath, []byte(text), 0644)
}

// ReadGile RAW
func RawReadFile(destPath string) (string, error) {
	data, err := os.ReadFile(destPath)
	return string(data), err
}

func Restart() error {
	dockerName, err := web.AppConfig.String("OpenVpnServerDockerName")
	if err != nil {
		return err
	}
	cmd := exec.Command("/bin/bash", "-c",
		fmt.Sprintf(
			"cd /opt/scripts/ && export OPENVPN_SERVER_DOCKER_NAME="+
				dockerName+" && ./restart.sh"))
	cmd.Dir = state.GlobalCfg.OVConfigPath
	output, err := cmd.CombinedOutput()
	if err != nil {
		logs.Debug(string(output))
		logs.Error(err)
		return err
	}
	return nil
}

func Backup() (string, error) {

	siteName := state.GlobalCfg.ServerName
	siteName = strings.ReplaceAll(siteName, " ", "")

	timestamp := time.Now().Format("20060102150405")
	backupFileName := fmt.Sprintf("backup_%s_%s.tar.bz2", siteName, timestamp)

	src := state.GlobalCfg.OVConfigPath
	dest := path.Join("/tmp/", backupFileName)

	cmd := exec.Command("/bin/bash", "-c",
		fmt.Sprintf(
			"/bin/tar -cjvf "+dest+" "+src))
	cmd.Dir = state.GlobalCfg.OVConfigPath
	output, err := cmd.CombinedOutput()
	if err != nil {
		logs.Debug(string(output))
		logs.Error(err)
		return "", err
	}
	return dest, err
}

func GenRandomString(n int) string {

	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}

	return string(b)
}

func GetMD5SumFile(path string) (string, error) {

	file, err := os.Open(path)
	if err != nil {
		return "", nil
	}

	hash := md5.New()

	_, err = io.Copy(hash, file)
	if err != nil {
		return "", err
	}

	md5 := hash.Sum(nil)

	return fmt.Sprintf("%x", md5), nil
}

func trim(s string) string {
	return strings.Trim(strings.Trim(s, "\r\n"), "\n")
}
