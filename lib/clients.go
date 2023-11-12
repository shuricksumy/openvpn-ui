package lib

import (
	"bytes"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/beego/beego/v2/core/logs"
)

// Structure for using on WEB
type ClientDetails struct {
	ClientName          string   `form:"client_name" 		 	json:"ClientName"`
	StaticIP            string   `form:"static_ip" 			 	json:"StaticIP"`
	IsRouteDefault      bool     `form:"is_route_default" 	 	json:"IsRouteDefault"`
	IsRouter            bool     `form:"is_router" 			 	json:"IsRouter"`
	RouterSubnet        string   `form:"router_subnet"		 	json:"RouterSubnet"`
	RouterMask          string   `form:"router_mask" 		 	json:"RouterMask"`
	Description         string   `form:"description" 		 	json:"Description"`
	RouteList           []string `form:"route_list_selected"  	json:"RouteListSelected"`
	RouteListUnselected []string `form:"route_list_unselected" 	json:"RouteListUnSelected"`
	CSRFToken           string   `form:"csrftoken" 			 	json:"CSRFToken"`
	MD5Sum              string   `json:"MD5Sum"`
}

type NameSorterClientDetails []*ClientDetails

func (a NameSorterClientDetails) Len() int           { return len(a) }
func (a NameSorterClientDetails) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a NameSorterClientDetails) Less(i, j int) bool { return a[i].ClientName < a[j].ClientName }

// Structure to read easy-rsa index file
type Client struct {
	ClientName string
	StaticIP   string
}

type ClientMD5 struct {
	ClientName string
	MD5Hash    string
	IsValid    bool
}

// Read index easy-rsa index file
func ReadClientsFromIndexFile(path string) ([]*Client, error) {

	clients := make([]*Client, 0)

	text, err := os.ReadFile(path)
	if err != nil {
		return clients, err
	}
	lines := strings.Split(trim(string(text)), "\n")

	for i, line := range lines {
		// Skip first item - server cert
		if i == 0 {
			continue
		}

		fields := strings.Fields(trim(line))

		if fields[0] == "V" {
			c := &Client{
				ClientName: parseClientName(fields[4]),
				StaticIP:   parseClientIP(fields[4]),
			}
			clients = append(clients, c)
		} else if fields[0] == "R" {
			c := &Client{
				ClientName: parseClientName(fields[5]),
				StaticIP:   parseClientIP(fields[5]),
			}
			clients = append(clients, c)
		} else {
			c := &Client{
				ClientName: parseClientName(fields[4]),
				StaticIP:   parseClientIP(fields[4]),
			}
			clients = append(clients, c)
		}

	}
	return clients, nil
}

// Read JSON file
func ReadClientsFromJSONFile(path string) ([]*ClientDetails, error) {

	clientsDetails := make([]*ClientDetails, 0)

	// Open our jsonFile
	byteValue, err := os.ReadFile(path)
	if err != nil {
		return clientsDetails, err
	}

	json.Unmarshal(byteValue, &clientsDetails)

	return clientsDetails, nil
}

// Populate new file only clients from easy-rsa file
// if client described get its info
func CombineIndexJsonResults(clientsDetailsFromIndex []*Client, clientsDetailsFromJSON []*ClientDetails) []*ClientDetails {
	clientsDetailsResult := make([]*ClientDetails, 0)

	for _, indexClient := range clientsDetailsFromIndex {
		clientMatched := InitClientFromStructure(*indexClient, clientsDetailsFromJSON)
		clientsDetailsResult = append(clientsDetailsResult, &clientMatched)
	}

	return clientsDetailsResult
}

// Combine two sourses: INDEX and JSON files and produce new structure for WEB
func GetClientsDetails(pathIndex string, pathJson string) ([]*ClientDetails, error) {
	// Crate empty list of clients
	clientsDetailsResult := make([]*ClientDetails, 0)

	// Get all describet clients from file
	clientsDetailsFromFile, errJson := ReadClientsFromJSONFile(pathJson)
	if errJson != nil {
		return clientsDetailsResult, errJson
	}

	// Get all real existed clients from easy-rsa file
	clientsFromIndex, errIndex := ReadClientsFromIndexFile(pathIndex)
	if errIndex != nil {
		return clientsDetailsResult, errIndex
	}

	// Populate new file only clients from easy-rsa file
	// if client described get its info
	for _, indexClient := range clientsFromIndex {
		clientMatched := InitClientFromStructure(*indexClient, clientsDetailsFromFile)
		clientsDetailsResult = append(clientsDetailsResult, &clientMatched)
	}

	// get list of routers in result file
	var allRouters []string = GetRouterClients(clientsDetailsResult)

	// validate existed selected - remove not in client list
	// upend new non selected  - RouteListUnselected
	// exclude itself route
	clientsDetailsResult = UpdateRoutersToActual(clientsDetailsResult, allRouters)

	// return modified collection
	return clientsDetailsResult, nil

}

// Parse Index file - get client Name
func parseClientName(d string) string {
	lines := strings.Split(trim(d), "/")
	for _, line := range lines {
		if strings.Contains(line, "") {
			fields := strings.Split(trim(line), "=")
			switch fields[0] {
			case "name":
				return fields[1]
			case "CN":
				return fields[1]
			default:
				// logs.Warn(fmt.Sprintf("Undefined entry: %s", line))
			}
		}
	}
	return "" //todo
}

// Parse Index file - get client IP
func parseClientIP(d string) string {
	lines := strings.Split(trim(d), "/")
	for _, line := range lines {
		if strings.Contains(line, "") {
			fields := strings.Split(trim(line), "=")
			switch fields[0] {
			case "LocalIP":
				return fields[1]
			default:
				// logs.Warn(fmt.Sprintf("Undefined entry: %s", line))
			}
		}
	}
	return "" //todo
}

// Parse JSON file
func parseJSON(jline string) ClientDetails {
	jbyte := []byte(jline)
	var clientsDetails ClientDetails
	json.Unmarshal(jbyte, &clientsDetails)
	return clientsDetails
}

// Comapare INDEX vs JSON
// if exist - return client from JSON
// if not exist - generate new dummy
func InitClientFromStructure(findClient Client, clients []*ClientDetails) ClientDetails {
	for _, client := range clients {
		if client.ClientName == findClient.ClientName {
			return *client
		}
	}
	newClient := &ClientDetails{
		ClientName:     findClient.ClientName,
		StaticIP:       findClient.StaticIP,
		IsRouteDefault: false,
		IsRouter:       false,
		RouterSubnet:   "",
		RouterMask:     "",
		Description:    "New record from Index file",
		RouteList:      []string{},
		CSRFToken:      "",
	}

	return *newClient
}

func GetRouterClients(clients []*ClientDetails) []string {
	var routers []string
	for _, client := range clients {
		if client.IsRouter {
			routers = append(routers, client.ClientName)
		}
	}
	return routers
}

func GetSelectedClientRouters(clients []*ClientDetails, clientName string) []string {
	var selectedRoutes []string
	for _, client := range clients {
		if client.ClientName == clientName {
			return client.RouteList
		}
	}
	return selectedRoutes
}

func CombineSelectedRouters(selected []string, all []string) []string {
	var resultRoute []string
	for _, route := range selected {
		for _, r := range all {
			if route == r {
				resultRoute = append(resultRoute, route)
				continue
			}
		}

	}
	return resultRoute
}

// Populate unselected routers only by valid values
func CombineUnSelectedRouters(selected []string, all []string, clientName string) []string {
	var resultNewRoute []string
	for _, route := range all {
		var routeExist bool = false

		// skip in selected list
		for _, r := range selected {
			// if not selected or not itself
			if route == r {
				routeExist = true
			}
		}

		// skip itself
		if route == clientName {
			routeExist = true
		}

		if !routeExist {
			resultNewRoute = append(resultNewRoute, route)
		}
	}
	return resultNewRoute
}

// update Result structure with actual routers per client - selected/nonSelected
func UpdateRoutersToActual(clients []*ClientDetails, allRouters []string) []*ClientDetails {
	for _, client := range clients {
		var selectedRouters []string = GetSelectedClientRouters(clients, client.ClientName)
		var modSelectedRouters []string = CombineSelectedRouters(selectedRouters, allRouters)
		var nonSelectedRouters []string = CombineUnSelectedRouters(selectedRouters, allRouters, client.ClientName)

		client.RouteList = modSelectedRouters
		client.RouteListUnselected = nonSelectedRouters
	}
	return clients
}

// Get Client from list
func GetClientFromStructure(clients []*ClientDetails, clientName string) (ClientDetails, error) {
	for _, client := range clients {
		if client.ClientName == clientName {
			return *client, nil
		}
	}
	var emptyClient ClientDetails
	return emptyClient, errors.New("Not Found client")
}

// compile new client from web
func UpdateClientsDetailsInStructure(clientsDetails []*ClientDetails, client ClientDetails) ([]*ClientDetails, error) {
	newClientsDetails := make([]*ClientDetails, 0)

	for _, c := range clientsDetails {
		if c.ClientName == client.ClientName {
			if c.CSRFToken == client.CSRFToken {
				client.CSRFToken = GenRandomString(32)
				client.MD5Sum = ""
				newClientsDetails = append(newClientsDetails, &client)
				continue
			} else {
				return newClientsDetails, errors.New("FILE WAS MODIFIED DURING YOUR OPPERATION - TRY AGAIN")
			}
		}
		c.CSRFToken = GenRandomString(32)
		newClientsDetails = append(newClientsDetails, c)
	}

	return newClientsDetails, nil
}

func GetClientsDetailsFromFiles() ([]*ClientDetails, error) {
	InitGlobalVars()
	//get clientsDetails from file
	clientsDetails, err := GetClientsDetails(PATH_INDEX, PATH_JSON)
	if err != nil {
		return clientsDetails, err
	}

	sort.Sort(NameSorterClientDetails(clientsDetails))

	return clientsDetails, nil
}

func AddClientToJsonFile(client ClientDetails) error {
	InitGlobalVars()
	wasError := false
	//get clientsDetails from file
	clientsDetails, err_read := GetClientsDetailsFromFiles()
	if err_read != nil {
		wasError = true
		logs.Error(err_read)
		logs.Error("ERROR WHILE READING CLIENTS FROM FILE !")
	}

	newClientDetails, err_upd := UpdateClientsDetailsInStructure(clientsDetails, client)
	if err_upd != nil {
		wasError = true
		logs.Error(err_upd)
		logs.Error("FILE WAS MODIFIED DURING YOU UPDATE - TRY AGAIN")
	}

	if !wasError {
		err_save := SaveJsonFile(newClientDetails, PATH_JSON)
		if err_save != nil {
			logs.Error(err_save)
			logs.Error("FILE WAS NOT SAVE")
			return errors.New("FILE WAS NOT SAVE")
		}
		return nil
	}

	return errors.New("Something goes wrong :(")
}

func SaveJsonFile(clientDetails []*ClientDetails, pathJson string) error {
	sort.Sort(NameSorterClientDetails(clientDetails))
	file, err := json.MarshalIndent(clientDetails, "", " ")
	if err != nil {
		return err
	}
	return os.WriteFile(pathJson, file, 0644)
}

// func ApplyClientsConfigToFS2() error {
// 	cmd := exec.Command("/bin/bash", "-c", " ./createClientFilesFromJSON.sh")
// 	cmd.Dir = state.GlobalCfg.OVConfigPath
// 	output, err := cmd.CombinedOutput()
// 	if err != nil {
// 		logs.Debug(string(output))
// 		logs.Error(err)
// 		return err
// 	}
// 	return nil
// }

func GetClientDetailsFieldValue(clientName string, fieldName string) string {
	allClients, _ := GetClientsDetailsFromFiles()
	client, _ := GetClientFromStructure(allClients, clientName)
	return reflect.ValueOf(&client).Elem().FieldByName(fieldName).String()
}

func GetMD5StructureFromFS(clients []*ClientDetails) []*ClientMD5 {
	md5Hashs := make([]*ClientMD5, 0)
	for _, c := range clients {
		isHashValid := false

		pathFile := filepath.Join(CCD_DIR_PATH, c.ClientName)
		md5Client, err_get_md5 := GetMD5SumFile(pathFile)

		if err_get_md5 != nil {
			isHashValid = false
		}

		if c.MD5Sum == md5Client {
			isHashValid = true
		}

		if c.MD5Sum == "" || md5Client == "" {
			isHashValid = false
		}

		newMD5 := &ClientMD5{
			ClientName: c.ClientName,
			MD5Hash:    md5Client,
			IsValid:    isHashValid,
		}
		md5Hashs = append(md5Hashs, newMD5)
	}

	return md5Hashs

}

func GetMD5StatusForClient(clients []*ClientDetails, clientName string) bool {

	for _, c := range clients {
		if c.ClientName == clientName {
			pathFile := filepath.Join(CCD_DIR_PATH, clientName)
			md5Client, err_get_md5 := GetMD5SumFile(pathFile)

			if err_get_md5 != nil {
				return false
			}

			if c.MD5Sum == "" || md5Client == "" {
				return false
			}

			if c.MD5Sum == md5Client {
				return true
			}

		}
	}

	return false
}

func UpdateJSONWithLatestMD5() error {
	clientsDetails, err_read := GetClientsDetailsFromFiles()
	if err_read != nil {
		return err_read
	}

	for _, c := range clientsDetails {
		md5hashs := GetMD5StructureFromFS(clientsDetails)

		for _, h := range md5hashs {
			if c.ClientName == h.ClientName {
				c.MD5Sum = h.MD5Hash
			}
		}
	}

	err_save_new_json := SaveJsonFile(clientsDetails, PATH_JSON)
	if err_save_new_json != nil {
		return err_save_new_json
	}

	return nil
}

func RawReadClientFile(clientName string) (string, error) {
	InitGlobalVars()
	destPathClientConfig := filepath.Join(CCD_DIR_PATH, clientName)
	return RawReadFile(destPathClientConfig)
}

func ApplyClientsConfigToFS() error {
	serverConf, err_read := RawReadFile(SERVER_CONFIG_PATH)
	if err_read != nil {
		logs.Error("Issue with reading config server file: ", SERVER_CONFIG_PATH)
		return err_read
	}

	isTopologyNet30, _ := regexp.MatchString(`\ntopology net30\s*\n`, serverConf)
	isTopologySubnet, _ := regexp.MatchString(`\ntopology subnet\s*\n`, serverConf)

	ClientDetails, err_read_json := ReadClientsFromJSONFile(PATH_JSON)
	if err_read_json != nil {
		logs.Error("Issue with reading JSON file: ", PATH_JSON)
		return err_read_json
	}

	for _, client := range ClientDetails {
		var buffer bytes.Buffer

		// 1. add init line
		buffer.WriteString("Automatic generated \"" + client.ClientName + "\" settings file - " +
			time.Now().Format(time.RFC850) + "\n")

		// 2. add static IP
		staticIp := client.StaticIP
		if _isIPAddressValid(staticIp) {
			var nextIP string

			if isTopologyNet30 {
				nextIP = _getNextIPAddress(staticIp) // --topology net30
			}

			if isTopologySubnet {
				nextIP = "255.255.255.0" // --topology subnet
			}

			if nextIP != "" {
				buffer.WriteString("ifconfig-push " + staticIp + " " + nextIP + "\n")
			} else {
				logs.Error("Cannot get nextIP or subnet topology is undefined in openvpn config file. Use as topology net30")
				nextIP = _getNextIPAddress(staticIp) // --topology net30
				buffer.WriteString("ifconfig-push " + staticIp + " " + nextIP + "\n")
			}
		}

		// 3. if router add route to itself
		if client.IsRouter {
			routerSubnet := client.RouterSubnet
			routerMask := client.RouterMask
			if _isIPAddressValid(routerSubnet) && _isIPAddressValid(routerMask) {
				buffer.WriteString("iroute " + routerSubnet + " " + routerMask + "\n")
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

		// 5. if Routes is xonfigured
		for _, route := range client.RouteList {
			routeDetails, _ := GetClientFromStructure(ClientDetails, route)
			routeSubnet := routeDetails.RouterSubnet
			routeMask := routeDetails.RouterMask

			if _isIPAddressValid(routeSubnet) && _isIPAddressValid(routeMask) {
				buffer.WriteString("# Route to " + route + " [" + routeDetails.Description + "] device internal subnet\n")
				buffer.WriteString("push \"route " + routeSubnet + " " + routeMask + "\"\n")
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