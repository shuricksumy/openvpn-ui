package lib

import (
	"encoding/json"
	"errors"
	"os"
	"sort"

	"github.com/beego/beego/v2/core/logs"
)

// Structure for using on WEB
type RouteDetails struct {
	RouteID       string `form:"route_id" 				json:"RouteID"`
	RouterName    string `form:"router_name"		 	json:"RouterName"`
	RouteIP       string `form:"route_ip"			 	json:"RouteIP"`
	RouteMask     string `form:"route_mask" 		 	json:"RouteMask"`
	Description   string `form:"description" 		 	json:"Description"`
	CSRFToken     string `form:"csrftoken" 			 	json:"CSRFToken"`
	RouterIsValid bool   `json:"RouterIsValid"`
	RouterIsUsed  bool   `json:"RouterIsUsed"`
	RouteIsUsedBy []string
}

type NameSorterRoutesDetails []*RouteDetails

func (a NameSorterRoutesDetails) Len() int           { return len(a) }
func (a NameSorterRoutesDetails) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a NameSorterRoutesDetails) Less(i, j int) bool { return a[i].RouterName < a[j].RouterName }

// Read Routes JSON file
func ReadRoutesFromJSONFile(path string) ([]*RouteDetails, error) {

	routesDetails := make([]*RouteDetails, 0)

	// Open our jsonFile
	byteValue, err := os.ReadFile(path)
	if err != nil {
		return routesDetails, err
	}

	json.Unmarshal(byteValue, &routesDetails)

	return routesDetails, nil
}

// Combine two sourses: INDEX and JSON files and produce new structure for WEB
func GetRoutesDetails(pathIndex string, pathJson string) ([]*RouteDetails, error) {
	// Crate empty list of routes
	routesDetailsResult := make([]*RouteDetails, 0)

	// Get all describet routes from file
	routesDetailsFromFile, errJson := ReadRoutesFromJSONFile(pathJson)
	if errJson != nil {
		return routesDetailsFromFile, errJson
	}

	// Get all real existed clients from easy-rsa file
	clientsFromIndex, errIndex := ReadClientsFromIndexFile(pathIndex)
	if errIndex != nil {
		return routesDetailsFromFile, errIndex
	}

	// Populate new file only clients from easy-rsa file
	// if client described get its info
	for _, iRoute := range routesDetailsFromFile {
		routeMatched, err_index := InitRouteFromStructure(*iRoute, clientsFromIndex)
		if err_index != nil {
			continue
		}
		routesDetailsResult = append(routesDetailsResult, &routeMatched)
	}

	// return modified collection
	return routesDetailsResult, nil

}

func GetRawRoutesDetailsFromFiles() ([]*RouteDetails, error) {
	InitGlobalVars()

	//get routesDetails from file
	routeDetails, err := GetRoutesDetails(PATH_INDEX, PATH_ROUTES_JSON)
	if err != nil {
		return routeDetails, err
	}

	sort.Sort(NameSorterRoutesDetails(routeDetails))

	return routeDetails, nil
}

func GetRoutesDetailsFromFiles() ([]*RouteDetails, error) {
	InitGlobalVars()

	//get routesDetails from file
	routeDetails, err := GetRoutesDetails(PATH_INDEX, PATH_ROUTES_JSON)
	if err != nil {
		return routeDetails, err
	}

	sort.Sort(NameSorterRoutesDetails(routeDetails))

	routeDetails, err_validation := ValidateRoutersInRouteList(routeDetails)
	if err_validation != nil {
		return nil, err_validation
	}
	return routeDetails, nil
}

func AddRouteToJsonFile(route RouteDetails) error {
	InitGlobalVars()
	wasError := false
	errMsg := ""

	//get routessDetails from file
	routeDetails, err_read := GetRoutesDetailsFromFiles()
	if err_read != nil {
		wasError = true
		errMsg = err_read.Error()
		logs.Error(err_read)
		logs.Error("ERROR WHILE READING ROUTES FROM FILE !")
	}

	newRouteDetails, err_upd := UpdateRoutessDetailsInStructure(routeDetails, route)
	if err_upd != nil {
		wasError = true
		errMsg = err_upd.Error()
		logs.Error(err_upd)
		logs.Error("FILE WAS MODIFIED DURING YOU UPDATE - TRY AGAIN")
	}

	if !wasError {
		err_save := SaveRouteJsonFile(newRouteDetails, PATH_ROUTES_JSON)
		if err_save != nil {
			logs.Error(err_save)
			logs.Error("FILE WAS NOT SAVE")
			return errors.New("FILE WAS NOT SAVE!" + err_save.Error())
		}
		return nil
	}

	return errors.New("Something goes wrong: " + errMsg)
}

func SaveRouteJsonFile(routeDetails []*RouteDetails, pathJson string) error {
	sort.Sort(NameSorterRoutesDetails(routeDetails))
	routeDetails, err_validation := ValidateRoutersInRouteList(routeDetails)
	if err_validation != nil {
		return err_validation
	}
	file, err := json.MarshalIndent(routeDetails, "", " ")
	if err != nil {
		return err
	}
	return os.WriteFile(pathJson, file, 0644)
}

// Comapare Route vs INDEX
// if exist - return true
// if not exist - return false
func InitRouteFromStructure(route RouteDetails, clients []*Client) (RouteDetails, error) {
	for _, client := range clients {
		if client.ClientName == route.RouterName {
			return route, nil
		}
	}
	return route, errors.New("Client does not exist")
}

// compile new client from web
func UpdateRoutessDetailsInStructure(routesDetails []*RouteDetails, route RouteDetails) ([]*RouteDetails, error) {
	newRoutesDetails := make([]*RouteDetails, 0)
	new := true

	for _, c := range routesDetails {
		if c.RouteID == route.RouteID {
			if c.CSRFToken == route.CSRFToken {
				route.CSRFToken = GenRandomString(32)
				newRoutesDetails = append(newRoutesDetails, &route)
				new = false
				continue
			} else {
				return newRoutesDetails, errors.New("FILE WAS MODIFIED DURING YOUR OPPERATION - TRY AGAIN")
			}
		}

		if c.RouteIP == route.RouteIP {
			return newRoutesDetails, errors.New("IP Route is duplicated")
		}

		c.CSRFToken = GenRandomString(32)
		newRoutesDetails = append(newRoutesDetails, c)
	}

	if new {
		route.CSRFToken = GenRandomString(32)
		newRoutesDetails = append(newRoutesDetails, &route)
	}

	return newRoutesDetails, nil
}

func GetRouteDetails(routeID string) *RouteDetails {
	routeDetails, _ := GetRoutesDetailsFromFiles()
	for _, r := range routeDetails {
		if r.RouteID == routeID {
			return r
		}
	}
	return nil
}

func DeleteRoute(routeID string) error {
	InitGlobalVars()

	routeDetails, err_read_file := GetRoutesDetailsFromFiles()
	if err_read_file != nil {
		return err_read_file
	}

	newRoutesDetails := make([]*RouteDetails, 0)
	for _, r := range routeDetails {
		if r.RouteID == routeID {
			continue
		}
		newRoutesDetails = append(newRoutesDetails, r)
	}

	err_save := SaveRouteJsonFile(newRoutesDetails, PATH_ROUTES_JSON)
	if err_save != nil {
		logs.Error(err_save)
		logs.Error("FILE WAS NOT SAVE")
		return errors.New("FILE WAS NOT SAVE!" + err_save.Error())
	}

	return nil
}

func ValidateRoutersInRouteList(routesDetails []*RouteDetails) ([]*RouteDetails, error) {
	newRoutesDetails := make([]*RouteDetails, 0)

	clientsDetails, err_read_file := GetClientsDetailsFromFiles()
	if err_read_file != nil {
		return nil, err_read_file
	}

	for _, r := range routesDetails {
		client, err_client := GetClientFromStructure(clientsDetails, r.RouterName)

		r.RouteIsUsedBy, _ = whoUseRoute(r.RouteID)

		if err_client != nil {
			r.RouterIsValid = false
			newRoutesDetails = append(newRoutesDetails, r)
			continue
		}

		if client.IsRouter == false {
			r.RouterIsValid = false
			newRoutesDetails = append(newRoutesDetails, r)
			continue
		}

		r.RouterIsValid = true
		newRoutesDetails = append(newRoutesDetails, r)
	}

	return newRoutesDetails, nil

}

func GetRouterRoutes(routerName string) ([]*RouteDetails, error) {
	//get routessDetails from file
	allRoutes, err_read := GetRoutesDetailsFromFiles()
	if err_read != nil {
		return nil, err_read
	}
	newRoutesDetails := make([]*RouteDetails, 0)
	for _, r := range allRoutes {
		if r.RouterName == routerName {
			newRoutesDetails = append(newRoutesDetails, r)
		}
	}
	return newRoutesDetails, nil
}

//TODO ROUTE IS USED
func whoUseRoute(RouteID string) ([]string, error) {

	var clientsName []string

	// Get all real existed clients from easy-rsa file
	clientsDetails, errIndex := GetClientsDetailsFromFiles()
	if errIndex != nil {
		return nil, errIndex
	}
	for _, c := range clientsDetails {
		for _, r := range c.RouteList {
			if r.RouteID == RouteID {
				clientsName = append(clientsName, c.ClientName)
				continue
			}
		}
	}

	return clientsName, nil
}
