package lib

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/beego/beego/v2/core/logs"
)

// Structure for using on WEB
type ClientDetails struct {
	ClientName          string   `form:"client_name" 		 json:"ClientName"`
	StaticIP            string   `form:"static_ip" 			 json:"StaticIP"`
	IsRouteDefault      bool     `form:"is_route_default" 	 json:"IsRouteDefault"`
	IsRouter            bool     `form:"is_router" 			 json:"IsRouter"`
	RouterSubnet        string   `form:"router_subnet"		 json:"RouterSubnet"`
	RouterMask          string   `form:"router_mask" 		 json:"RouterMask"`
	Description         string   `form:"description" 		 json:"Description"`
	RouteList           []string `form:"route_list_selected"  json:"RouteListSelected"`
	RouteListUnselected []string `form:"route_list_unselected" json:"RouteListUnSelected"`
	CSRFToken           string   `form:"csrftoken" 			 json:"CSRFToken"`
}

// Structure to read easy-rsa index file
type Client struct {
	ClientName string
	StaticIP   string
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
func ReadJSONClientsDetailsFile(path string) ([]*ClientDetails, error) {

	clientsDetails := make([]*ClientDetails, 0)

	// Open our jsonFile
	byteValue, err := os.ReadFile(path)
	if err != nil {
		return clientsDetails, err
	}

	json.Unmarshal(byteValue, &clientsDetails)

	return clientsDetails, nil
}

// Combine two sourses: INDEX and JSON files and produce new structure for WEB
func GetClientsDetails(pathIndex string, pathJson string) ([]*ClientDetails, error) {
	// Crate empty list of clients
	clientsDetailsResult := make([]*ClientDetails, 0)

	// Get all describet clients from file
	clientsDetailsFromFile, errJson := ReadJSONClientsDetailsFile(pathJson)
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
		clientMatched := GetClientDetailsRaw(*indexClient, clientsDetailsFromFile)
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
				logs.Warn(fmt.Sprintf("Undefined entry: %s", line))
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
				logs.Warn(fmt.Sprintf("Undefined entry: %s", line))
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
func GetClientDetailsRaw(findClient Client, clients []*ClientDetails) ClientDetails {
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
		for _, r := range selected {
			// if not selected or not itself
			if route == r || route == clientName {
				routeExist = true
			}
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
func UpdateClientsDetails(clientsDetails []*ClientDetails, client ClientDetails) ([]*ClientDetails, error) {
	newClientsDetails := make([]*ClientDetails, 0)

	for _, c := range clientsDetails {
		if c.ClientName == client.ClientName {
			if c.CSRFToken == client.CSRFToken {
				client.CSRFToken = GenRandomString(32)
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

func SaveJsonFile(clientDetails []*ClientDetails, pathJson string) error {
	file, err := json.MarshalIndent(clientDetails, "", " ")
	if err != nil {
		return err
	}
	return os.WriteFile(pathJson, file, 0644)
}
