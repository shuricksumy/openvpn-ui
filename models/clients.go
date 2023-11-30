package models

import (
	"encoding/json"

	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/core/validation"
)

type ClientDetails struct {
	Id              int             `orm:"auto;pk"`
	ClientName      string          `orm:"unique" form:"client_name"`
	StaticIP        *string         `orm:"unique;null" form:"static_ip"`
	IsRouteDefault  bool            `valid:"Required" form:"use_def_routing"`
	IsRouter        bool            `valid:"Required" form:"client_is_router"`
	CertificateName *string         `orm:"unique;null" form:"certificate_name"`
	Description     string          `form:"description"`
	Routes          []*RouteDetails `orm:"rel(m2m)"`
	MD5Sum          string          `form:"md5sum"`
}

// Validate function to perform custom validation
func (c *ClientDetails) Validate() error {
	valid := validation.Validation{}
	b, err := valid.Valid(c)
	if err != nil {
		return err
	}
	if !b {
		// If validation fails, return an error with validation messages
		for _, err := range valid.Errors {
			return err
		}
	}
	return nil
}

// GetRouteDetailsByID retrieves a RouteDetails by its ID
func GetRouteDetailsByID(routeDetailsID int) (*RouteDetails, error) {
	var routeDetails RouteDetails
	err := orm.NewOrm().QueryTable(new(RouteDetails)).Filter("Id", routeDetailsID).RelatedSel().One(&routeDetails)
	return &routeDetails, err
}

// UpdateRouteDetails updates a RouteDetails by its ID
func UpdateRouteDetails(routeDetailsID int, updatedDetails *RouteDetails) error {
	var routeDetails RouteDetails
	if err := orm.NewOrm().QueryTable(new(RouteDetails)).Filter("Id", routeDetailsID).One(&routeDetails); err == nil {
		// Update the RouteDetails attributes
		routeDetails.Name = updatedDetails.Name
		routeDetails.RouterName = updatedDetails.RouterName
		routeDetails.RouteIP = updatedDetails.RouteIP
		routeDetails.RouteMask = updatedDetails.RouteMask
		routeDetails.Description = updatedDetails.Description

		// Save the updated RouteDetails
		_, err := orm.NewOrm().Update(&routeDetails)
		return err
	}
	return nil
}

// AddNewClient creates a new client and adds it to the database
func AddNewClient(clientName string, staticIP *string, isRouteDefault, isRouter bool, description, md5Sum string, routeIDs []int) error {
	o := orm.NewOrm()

	client := &ClientDetails{
		ClientName:     clientName,
		StaticIP:       staticIP,
		IsRouteDefault: isRouteDefault,
		IsRouter:       isRouter,
		Description:    description,
		MD5Sum:         md5Sum,
	}

	// Add routes to the client
	for _, routeID := range routeIDs {
		route := &RouteDetails{Id: routeID}
		if err := o.Read(route); err == nil {
			route.Client = append(route.Client, client)
		}
	}

	// Save the client to the database
	_, err := o.Insert(client)
	return err
}

// AssignRoutesToClient assigns one or more routes to a specific client
func AssignRoutesToClient(clientID int, routeIDs []int) error {
	o := orm.NewOrm()

	client := &ClientDetails{Id: clientID}
	if err := o.Read(client); err != nil {
		return err
	}

	var routes []*RouteDetails
	for _, routeID := range routeIDs {
		route := &RouteDetails{Id: routeID}
		if err := o.Read(route); err == nil {
			routes = append(routes, route)
		}
	}

	// Add routes to the client
	o.QueryM2M(client, "RouteList").Add(routes)

	return nil
}

// DeleteClientDetailsByID deletes a client by its ID
func DeleteClientDetailsByID(clientID int) error {
	o := orm.NewOrm()

	client := &ClientDetails{Id: clientID}
	if err := o.Read(client); err == nil {
		// Remove the client from associated routes
		o.QueryM2M(client, "Routes").Clear()

		// Delete the client from the database
		_, err := o.Delete(client)
		return err
	}

	return nil
}

// GetRouterClients retrieves clients where IsRouter is true
func GetRouterClients() ([]*ClientDetails, error) {
	o := orm.NewOrm()

	var clients []*ClientDetails
	if _, err := o.QueryTable(new(ClientDetails)).Filter("IsRouter", true).All(&clients); err == nil {
		return clients, nil
	}

	return nil, nil
}

// GetConnectedRouteDetails retrieves all RouteDetails connected to a specific client
func GetConnectedRouteDetails(clientID int) ([]*RouteDetails, error) {
	o := orm.NewOrm()

	client := &ClientDetails{Id: clientID}
	if err := o.Read(client); err == nil {
		var routeDetails []*RouteDetails
		o.QueryTable(new(RouteDetails)).Filter("Client__ClientDetails__Id", clientID).Distinct().All(&routeDetails)
		return routeDetails, nil
	}

	return nil, nil
}

func GetDisconnectedRouteDetails(clientID int) ([]*RouteDetails, error) {
	o := orm.NewOrm()

	// Get the client
	client := &ClientDetails{Id: clientID}
	if err := o.Read(client); err != nil {
		return nil, err
	}

	// Load the associated RouteDetails
	o.LoadRelated(client, "Routes")

	// Get all routes
	allRoutes := make([]*RouteDetails, 0)
	_, err := o.QueryTable("route_details").All(&allRoutes)
	if err != nil {
		return nil, err
	}

	// Filter out routes where the client is a router
	disconnectedRoutes := make([]*RouteDetails, 0)
	for _, route := range allRoutes {

		if route.RouterName == client.ClientName {
			continue
		}

		routerFound := false
		for _, connectedRoute := range client.Routes {
			if route.Id == connectedRoute.Id {
				routerFound = true
				break
			}
		}

		if !routerFound {
			disconnectedRoutes = append(disconnectedRoutes, route)
		}
	}

	return disconnectedRoutes, nil
}

func ClientExistsByName(clientName string) bool {
	o := orm.NewOrm()

	return o.QueryTable(new(ClientDetails)).Filter("ClientName", clientName).Exist()
}

func GetAllClientDetails() ([]*ClientDetails, error) {
	o := orm.NewOrm()

	var clients []*ClientDetails
	if _, err := o.QueryTable(new(ClientDetails)).All(&clients); err == nil {
		return clients, nil
	}

	return nil, nil
}

// Custom function defined in the controller
func GetConnectedRoutes(inputId int) []*RouteDetails {
	// id, _ := strconv.Atoi(inputId)
	routes, _ := GetConnectedRouteDetails(inputId)
	return routes
}

// Custom function defined in the controller
func GetDisConnectedRoutes(inputId int) []*RouteDetails {
	// id, _ := strconv.Atoi(inputId)
	routes, _ := GetDisconnectedRouteDetails(inputId)
	return routes
}

// Dump any structure as json string
func Dump(obj interface{}) {
	result, _ := json.MarshalIndent(obj, "", "\t")
	logs.Debug(string(result))
}
