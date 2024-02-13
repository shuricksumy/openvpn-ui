package lib

import (
	"bufio"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/netip"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/core/validation"
	"github.com/beego/beego/v2/server/web"
	"github.com/shuricksumy/openvpn-ui/state"
)

var PATH_INDEX string
var PATH_JSON string
var PATH_ROUTES_JSON string
var CCD_DIR_PATH string
var SERVER_CONFIG_PATH string

func InitGlobalVars() {
	PATH_INDEX = filepath.Join(state.GlobalCfg.OVConfigPath, "easy-rsa/pki/index.txt")
	PATH_JSON = filepath.Join(state.GlobalCfg.OVConfigPath, "clientDetails.json")
	PATH_ROUTES_JSON = filepath.Join(state.GlobalCfg.OVConfigPath, "routesDetails.json")
	CCD_DIR_PATH = filepath.Join(state.GlobalCfg.OVConfigPath, "ccd")
	SERVER_CONFIG_PATH = filepath.Join(state.GlobalCfg.OVConfigPath, "server.conf")
}

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

	dbName, _ := web.AppConfig.String("DbPath")
	backupFileNameSQL := filepath.Join(state.GlobalCfg.OVConfigPath, fmt.Sprintf("backup_%s.sql", timestamp))

	err := DumpSQLiteDatabaseToFile(backupFileNameSQL, dbName)
	if err != nil {
		// Handle error, e.g., command execution error
		logs.Error("Error dumping database")
		return dest, err
	} else {
		// Successful database dump
		logs.Warn("Database dumped successfully")
	}

	cmd := exec.Command("/bin/bash", "-c", fmt.Sprintf("/bin/tar -cjvf %s %s", dest, src))

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

func _isIPAddressValid(ip string) bool {
	addr, _ := netip.ParseAddr(ip)
	return addr.IsValid()
}

func _getNextIPAddress(ip string) string {
	addr, err := netip.ParseAddr(ip)
	if err != nil {
		logs.Error("IP is nov valid: ", ip)
		return ""
	}
	return addr.Next().String()
}

func GetExtIP() (string, error) {

	siteName := state.GlobalCfg.ServerName
	siteName = strings.ReplaceAll(siteName, " ", "")

	cmd := exec.Command("/usr/bin/curl", "-s", "https://api.ipify.org")
	cmd.Dir = state.GlobalCfg.OVConfigPath
	output, err := cmd.CombinedOutput()
	if err != nil {
		logs.Debug(string(output))
		logs.Error(err)
		return "", err
	}
	return string(output), err
}

func StringToNilString(input string) *string {

	if input == "" {
		return nil
	}
	return &input
}

func NilStringToString(input *string) string {
	if input == nil {
		return ""
	}
	pointerToString := &input
	return **pointerToString
}

// DumpSQLiteDatabaseToFile dumps the entire SQLite database to a file
func DumpSQLiteDatabaseToFile(backupFileName string, dbName string) error {

	// Check if the SQLite database file exists
	if _, err := os.Stat(dbName); os.IsNotExist(err) {
		logs.Error(fmt.Errorf("SQLite database file does not exist: %s", dbName))
		return fmt.Errorf("SQLite database file does not exist: %s", dbName)
	}

	// Prepare the sqlite3 command
	cmd := exec.Command("/bin/bash", "-c",
		fmt.Sprintf("/usr/bin/sqlite3 %s .dump", dbName))

	cmd.Dir = state.GlobalCfg.OVConfigPath
	outputFile, err := os.Create(backupFileName)
	if err != nil {
		logs.Error(fmt.Sprintf("/usr/bin/sqlite3 %s .dump", dbName))
		logs.Error(outputFile)
		logs.Error(err)
	}
	defer outputFile.Close()

	// Set the command's output to the file
	cmd.Stdout = outputFile

	// Run the sqlite3 command
	err = cmd.Run()
	if err != nil {
		logs.Error(err)
		return err
	}

	logs.Warning("Database dumped successfully to %s\n", backupFileName)
	return nil
}

func AppendStringToFile(filePath, content string) {
	// Open the file in append mode, create it if it doesn't exist
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	// Append the content to the file
	_, err = file.WriteString(content)
	if err != nil {
		fmt.Println("Error appending to file:", err)
		return
	}

	fmt.Println("String appended successfully!")
}

func PatchFileAppendAfterLine(filePath, searchString, content string) error {
	// Open the file in read mode
	file, err := os.OpenFile(filePath, os.O_RDWR, 0644)
	if err != nil {
		return fmt.Errorf("Error opening file: %v", err)
	}
	defer file.Close()

	// Create a scanner to read the file line by line
	scanner := bufio.NewScanner(file)
	var lines []string
	lineNotFound := true

	// Read the file line by line
	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)

		// Check if the line contains the specified string
		if strings.Contains(line, searchString) {
			// Append the content after the line
			lines = append(lines, content)
			lineNotFound = false
		}
	}

	// Error if search searchString does not exist in file
	if lineNotFound == true {
		return fmt.Errorf("Searching line not found: %s", searchString)
	}

	// Check for errors during scanning
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("Error scanning file: %v", err)
	}

	// Rewind the file to the beginning
	file.Seek(0, 0)
	file.Truncate(0)

	// Write the modified lines back to the file
	writer := bufio.NewWriter(file)
	for _, line := range lines {
		fmt.Fprintln(writer, line)
	}

	// Flush and close the writer
	err = writer.Flush()
	if err != nil {
		return fmt.Errorf("Error writing to file: %v", err)
	}

	return nil
}

func PatchFileAppendBeforeAndAfterLine(filePath, searchString, contentBefore, contentAfter string) error {
	// Open the file in read mode
	file, err := os.OpenFile(filePath, os.O_RDWR, 0644)
	if err != nil {
		return fmt.Errorf("Error opening file: %v", err)
	}
	defer file.Close()

	// Create a scanner to read the file line by line
	scanner := bufio.NewScanner(file)
	var lines []string
	lineNotFound := true

	// Read the file line by line
	for scanner.Scan() {
		line := scanner.Text()

		// Check if the line contains the specified string
		if strings.Contains(line, searchString) {
			// Append contentBefore before the line
			lines = append(lines, contentBefore)
			lineNotFound = false
		}

		lines = append(lines, line)

		// Check if the line contains the specified string
		if strings.Contains(line, searchString) {
			// Append contentAfter after the line
			lines = append(lines, contentAfter)
			lineNotFound = false
		}
	}

	// Error if search searchString does not exist in file
	if lineNotFound == true {
		return fmt.Errorf("Searching line not found: %s", searchString)
	}

	// Check for errors during scanning
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("Error scanning file: %v", err)
	}

	// Rewind the file to the beginning
	file.Seek(0, 0)
	file.Truncate(0)

	// Write the modified lines back to the file
	writer := bufio.NewWriter(file)
	for _, line := range lines {
		fmt.Fprintln(writer, line)
	}

	// Flush and close the writer
	err = writer.Flush()
	if err != nil {
		return fmt.Errorf("Error writing to file: %v", err)
	}

	return nil
}

func PatchFileAppendBeforeLine(filePath, searchString, contentBefore string) error {
	// Open the file in read mode
	file, err := os.OpenFile(filePath, os.O_RDWR, 0644)
	if err != nil {
		return fmt.Errorf("Error opening file: %v", err)
	}
	defer file.Close()

	// Create a scanner to read the file line by line
	scanner := bufio.NewScanner(file)
	var lines []string
	lineNotFound := true

	// Read the file line by line
	for scanner.Scan() {
		line := scanner.Text()

		// Check if the line contains the specified string
		if strings.Contains(line, searchString) {
			// Append contentBefore before the line
			lines = append(lines, contentBefore)
			lineNotFound = false
		}

		lines = append(lines, line)
	}
	// Error if search searchString does not exist in file
	if lineNotFound == true {
		return fmt.Errorf("Searching line not found: %s", searchString)
	}

	// Check for errors during scanning
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("Error scanning file: %v", err)
	}

	// Rewind the file to the beginning
	file.Seek(0, 0)
	file.Truncate(0)

	// Write the modified lines back to the file
	writer := bufio.NewWriter(file)
	for _, line := range lines {
		fmt.Fprintln(writer, line)
	}

	// Flush and close the writer
	err = writer.Flush()
	if err != nil {
		return fmt.Errorf("Error writing to file: %v", err)
	}

	return nil
}

func ParsePrefixURL() string {
	prefixURL, err := web.AppConfig.String("BaseURLPrefix")
	if err != nil {
		prefixURL = ""
	}

	if !strings.HasPrefix(prefixURL, "/") {
		prefixURL = "/" + prefixURL
	}
	if strings.HasSuffix(prefixURL, "/") {
		prefixURL = strings.TrimSuffix(prefixURL, "/")
	}
	return prefixURL
}
