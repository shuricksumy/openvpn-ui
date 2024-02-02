package models

import (
	//Sqlite driver

	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/validation"
	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

type RouteDetails struct {
	Id          string           `orm:"pk;type(uuid);default(uuid_generate_v4());unique"`
	Name        string           `orm:"unique"`
	RouterName  string           `valid:"Required" form:"router_name"`
	RouterId    string           `valid:"Required"`
	RouteIP     string           `orm:"unique" form:"route_ip"`
	RouteMask   string           `form:"route_mask"`
	Description string           `form:"description"`
	Client      []*ClientDetails `orm:"reverse(many)"`
}

// Validate function to perform custom validation
func (r *RouteDetails) Validate() error {
	valid := validation.Validation{}
	b, err := valid.Valid(r)
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
func GetRouteDetailsByID(routeDetailsID string) (*RouteDetails, error) {
	var routeDetails RouteDetails
	err := orm.NewOrm().QueryTable(new(RouteDetails)).Filter("Id", routeDetailsID).RelatedSel().One(&routeDetails)
	return &routeDetails, err
}

// UpdateRouteDetails updates specified parameters for a RouteDetails instance by its ID
func UpdateRouteDetails(routeID string, routeIP, routeMask, description string) error {
	o := orm.NewOrm()

	// Get the existing route
	route := &RouteDetails{Id: routeID}
	err := o.Read(route)
	if err != nil {
		return err
	}

	// Update the specified parameters
	route.RouteIP = routeIP
	route.RouteMask = routeMask
	route.Description = description

	// Save the updated route
	_, err = o.Update(route)
	return err
}

// AddNewRouteDetails creates a new route and adds it to the database
func AddNewRouteDetails(name string, routerName string, routerId string, routeIP string, routeMask string, description string) error {
	o := orm.NewOrm()

	route := &RouteDetails{
		Id:          uuid.New().String(),
		Name:        name,
		RouterName:  routerName,
		RouterId:    routerId,
		RouteIP:     routeIP,
		RouteMask:   routeMask,
		Description: description,
	}

	_, err := o.Insert(route)
	return err
}

// UpdateRouteDetailsByID updates a route by its ID
func UpdateRouteDetailsByID(routeID string, updatedDetails *RouteDetails) error {
	o := orm.NewOrm()

	route := &RouteDetails{Id: routeID}
	if err := o.Read(route); err == nil {
		// Update the RouteDetails attributes
		route.Name = updatedDetails.Name
		route.RouterName = updatedDetails.RouterName
		route.RouteIP = updatedDetails.RouteIP
		route.RouteMask = updatedDetails.RouteMask
		route.Description = updatedDetails.Description

		// Save the updated RouteDetails
		_, err := o.Update(route)
		return err
	}
	return nil
}

// GetClientsForRouteID gets a list of client names that use a specific route by its ID
func GetClientsForRouteID(routeID string) ([]string, error) {
	o := orm.NewOrm()

	var clients []*ClientDetails
	if _, err := o.QueryTable(new(ClientDetails)).Filter("Routes__RouteDetails__Id", routeID).Distinct().All(&clients); err == nil {
		var clientNames []string
		for _, client := range clients {
			clientNames = append(clientNames, client.ClientName)
		}
		return clientNames, nil
	}

	return nil, nil
}

// DeleteRouteDetailsById deletes a RouteDetails instance by its ID
func DeleteRouteDetailsById(routeID string) error {
	o := orm.NewOrm()

	// Get the existing route
	route := &RouteDetails{Id: routeID}
	err := o.Read(route)
	if err != nil {
		return err
	}

	// Delete the route
	_, err = o.Delete(route)
	return err
}

// RouteExistsByIP checks if a route with the given IP exists
func RouteExistsByIP(routeIP string) bool {
	o := orm.NewOrm()

	return o.QueryTable(new(RouteDetails)).Filter("RouteIP", routeIP).Exist()
}

func GetAllRoutesDetails() ([]*RouteDetails, error) {
	o := orm.NewOrm()

	var routers []*RouteDetails
	if _, err := o.QueryTable(new(RouteDetails)).OrderBy("Name").All(&routers); err == nil {
		return routers, nil
	}

	return nil, nil
}

// Custom function defined in the controller
func RouteIsUsedBy(inputId string) []string {
	clients, _ := GetClientsForRouteID(inputId)
	return clients
}

// GetAllRoutesProvided retrieves all RouteDetails associated with a specific RouterId
func GetAllRoutesProvided(routerID string) ([]*RouteDetails, error) {
	o := orm.NewOrm()

	// Query RouteDetails with the given RouterId
	var routes []*RouteDetails
	_, err := o.QueryTable("route_details").Filter("RouterId", routerID).OrderBy("Name").All(&routes)

	if err != nil {
		// Handle error, e.g., database query error
		return nil, err
	}

	// Load related clients for each route
	for _, route := range routes {
		o.LoadRelated(route, "Client")
	}

	return routes, nil
}
