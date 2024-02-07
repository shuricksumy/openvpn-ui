package models

import (
	"encoding/json"
	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/core/validation"
	"github.com/google/uuid"
)

type ClientDetails struct {
	Id                string          `orm:"pk;type(uuid);default(uuid_generate_v4());unique"`
	ClientName        string          `orm:"unique" form:"client_name"`
	StaticIP          *string         `orm:"unique;null" form:"static_ip"`
	IsRouteDefault    bool            `valid:"Required" form:"use_def_routing"`
	IsRouter          bool            `valid:"Required" form:"client_is_router"`
	CertificateName   *string         `orm:"unique;null" form:"certificate_name"`
	CertificateStatus *string         `orm:"null"`
	Description       string          `form:"description"`
	Passphrase        string          `form:"cert_pass"`
	Routes            []*RouteDetails `orm:"rel(m2m)"`
	MD5Sum            string          `form:"md5sum"`
	OTPIsEnabled      bool            `valid:"Required" form:"otp_is_enabled"`
	OTPKey            *string         `orm:"unique;null" form:"otp_key"`
	OTPUserName       *string         `orm:"unique;null" form:"otp_user_name"`
	StaticPassIsUsed  bool            `valid:"Required" form:"static_pass_is_enabled"`
	StaticPass        *string         `orm:"null" form:"static_pass"`
}

type ClientDetailsExtended struct {
	ClientDetails       ClientDetails
	CertificateHasIssue bool
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

// AddNewClient creates a new client and adds it to the database
func AddNewClient(clientName string, staticIP *string, isRouteDefault, isRouter bool, description, md5Sum string, routeIDs []string, staticPass *string) error {
	o := orm.NewOrm()

	client := &ClientDetails{
		Id:               uuid.New().String(),
		ClientName:       clientName,
		StaticIP:         staticIP,
		IsRouteDefault:   isRouteDefault,
		IsRouter:         isRouter,
		Description:      description,
		MD5Sum:           md5Sum,
		OTPUserName:      &clientName,
		StaticPassIsUsed: true,
		StaticPass:       staticPass,
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

// GetClientDetailsById retrieves a client by its ID from the database
func GetClientDetailsById(clientID string) (*ClientDetails, error) {
	o := orm.NewOrm()

	client := &ClientDetails{Id: clientID}
	err := o.Read(client)

	if err == orm.ErrNoRows {
		// Client with the given ID not found
		return nil, nil
	} else if err != nil {
		// Other error occurred
		return nil, err
	}

	// Load the associated RouteDetails
	o.LoadRelated(client, "Routes")

	return client, nil
}

// AssignRoutesToClient assigns one or more routes to a specific client
func AssignRoutesToClient(clientID string, routeIDs []string) error {
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
	o.QueryM2M(client, "Routes").Add(routes)

	return nil
}

// DeleteClientDetailsByID deletes a client by its ID
func DeleteClientDetailsByID(clientID string) error {
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
	if _, err := o.QueryTable(new(ClientDetails)).Filter("IsRouter", true).OrderBy("ClientName").All(&clients); err == nil {

		// Load the associated RouteDetails for each client
		for _, client := range clients {
			o.LoadRelated(client, "Routes")
		}

		return clients, nil
	}

	return nil, nil
}

// GetConnectedRouteDetails retrieves all RouteDetails connected to a specific client
func GetConnectedRouteDetails(clientID string) ([]*RouteDetails, error) {
	o := orm.NewOrm()

	client := &ClientDetails{Id: clientID}
	if err := o.Read(client); err == nil {
		var routeDetails []*RouteDetails
		o.QueryTable(new(RouteDetails)).Filter("Client__ClientDetails__Id", clientID).Distinct().OrderBy("Name").All(&routeDetails)
		return routeDetails, nil
	}

	return nil, nil
}

func GetDisconnectedRouteDetails(clientID string) ([]*RouteDetails, error) {
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
	_, err := o.QueryTable("route_details").OrderBy("Name").All(&allRoutes)
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
	if _, err := o.QueryTable(new(ClientDetails)).OrderBy("ClientName").All(&clients); err == nil {

		// Load the associated RouteDetails for each client
		for _, client := range clients {
			o.LoadRelated(client, "Routes")
		}

		return clients, nil
	}

	return nil, nil
}

// UpdateClientDetails updates specified parameters for a ClientDetails instance
func UpdateClientDetails(clientID string, staticIP *string, description string, isRouteDefault, isRouter bool) error {
	o := orm.NewOrm()

	// Get the existing client
	client := &ClientDetails{Id: clientID}
	err := o.Read(client)
	if err != nil {
		return err
	}

	// Update the specified parameters
	if staticIP != nil {
		client.StaticIP = staticIP
	}

	client.Description = description
	client.IsRouteDefault = isRouteDefault
	client.IsRouter = isRouter

	// Save the updated client
	_, err = o.Update(client)
	return err
}

// UnassignAllRoutesFromClient unassigns all routes from a specific client
func UnassignAllRoutesFromClient(clientID string) error {
	o := orm.NewOrm()

	client := &ClientDetails{Id: clientID}
	if err := o.Read(client); err != nil {
		return err
	}

	// Clear all routes from the client
	o.QueryM2M(client, "Routes").Clear()

	return nil
}

// GetClientDetailsByCertificate retrieves a client by its CertificateName from the database
func GetClientDetailsByCertificate(certificateName string) (*ClientDetails, error) {
	o := orm.NewOrm()

	var clients []*ClientDetails
	if _, err := o.QueryTable(new(ClientDetails)).Filter("CertificateName", certificateName).All(&clients); err == nil {
		if len(clients) > 0 {

			// Load the associated RouteDetails for each client
			for _, client := range clients {
				o.LoadRelated(client, "Routes")
			}

			return clients[0], nil
		}
	} else {
		return nil, err
	}

	return nil, nil
}

// GetClientCertificateByName retrieves the CertificateName for a client by its ClientName
func GetClientCertificateByName(clientName string) (*string, error) {
	o := orm.NewOrm()

	// Use QueryTable to construct a custom query
	query := o.QueryTable(new(ClientDetails)).Filter("ClientName", clientName)

	// Get the existing client
	client := &ClientDetails{}
	err := query.One(client)
	if err == orm.ErrNoRows {
		// Client with the given ClientName not found
		return nil, nil
	} else if err != nil {
		// Other error occurred
		return nil, err
	}

	// Return the CertificateName
	return client.CertificateName, nil
}

// GetClientsDetailsWithoutCertificate retrieves clients without a CertificateName from the database
func GetClientsDetailsWithoutCertificate() ([]*ClientDetails, error) {
	o := orm.NewOrm()

	// Query ClientDetails where CertificateName is nil
	var clients []*ClientDetails
	_, err := o.QueryTable("client_details").Filter("CertificateName__isnull", true).OrderBy("ClientName").All(&clients)

	if err != nil {
		// Handle error, e.g., database query error
		return nil, err
	}

	// Load the associated RouteDetails for each client
	for _, client := range clients {
		o.LoadRelated(client, "Routes")
	}

	return clients, nil
}

// UpdateClientCertificateById updates the CertificateName for a client by its ID
func UpdateClientCertificateById(clientID string, certificateName *string, passphrase string) error {
	o := orm.NewOrm()

	// Get the existing client
	client := &ClientDetails{Id: clientID}
	err := o.Read(client)
	if err != nil {
		return err
	}

	// Update the CertificateName
	client.CertificateName = certificateName
	client.Passphrase = passphrase
	client.MD5Sum = "NEW CERT CREATED"

	// Save the updated client
	_, err = o.Update(client)
	return err
}

func UpdateClientCertificateStatusById(clientID string, certificateStatus *string) error {
	o := orm.NewOrm()

	// Get the existing client
	client := &ClientDetails{Id: clientID}
	err := o.Read(client)
	if err != nil {
		return err
	}

	// Update the CertificateName
	client.CertificateStatus = certificateStatus

	// Save the updated client
	_, err = o.Update(client)
	return err
}

// ClearClientCertificateById clears the CertificateName for a client by its ID
func ClearClientCertificateById(clientID string) error {
	o := orm.NewOrm()

	// Get the existing client
	client := &ClientDetails{Id: clientID}
	err := o.Read(client)
	if err != nil {
		return err
	}

	// Clear the CertificateName
	client.CertificateName = nil
	client.CertificateStatus = nil
	client.Passphrase = ""

	// Save the updated client
	_, err = o.Update(client)
	return err
}

// UpdatePassphraseById updates the Passphrase for a client by its ID
func UpdatePassphraseById(clientID string, passphrase string) error {
	o := orm.NewOrm()

	// Get the existing client
	client := &ClientDetails{Id: clientID}
	err := o.Read(client)
	if err != nil {
		return err
	}

	// Update the Passphrase
	client.Passphrase = passphrase

	// Save the updated client
	_, err = o.Update(client)
	return err
}

// GetAllClientsWithCertificate retrieves all clients where CertificateName is not nil
func GetAllClientsWithCertificate() ([]*ClientDetails, error) {
	o := orm.NewOrm()

	// Query clients with non-nil CertificateName
	var clients []*ClientDetails
	_, err := o.QueryTable("client_details").
		Filter("CertificateName__isnull", false).OrderBy("ClientName").
		All(&clients)

	if err != nil {
		return nil, err
	}

	// Load the associated RouteDetails for each client
	for _, client := range clients {
		o.LoadRelated(client, "Routes")
	}

	return clients, nil
}

// UpdateMD5SumForClientDetails updates the MD5Sum for a client by its ClientName
func UpdateMD5SumForClientDetails(clientName, newMD5Sum string) error {
	o := orm.NewOrm()

	// Use QueryTable to construct a custom query
	query := o.QueryTable(new(ClientDetails)).Filter("ClientName", clientName)

	// Get the existing client
	client := &ClientDetails{}
	err := query.One(client)
	if err != nil {
		return err
	}

	// Update the MD5Sum with the provided value
	client.MD5Sum = newMD5Sum

	// Save the updated client
	_, err = o.Update(client)
	return err
}

// UpdateMD5SumForClientDetailsByID updates the MD5Sum for a client by its Id
func UpdateMD5SumForClientDetailsByID(clientID string, newMD5Sum string) error {
	o := orm.NewOrm()

	// Get the existing client by Id
	client := &ClientDetails{Id: clientID}
	err := o.Read(client)
	if err != nil {
		return err
	}

	// Update the MD5Sum with the provided value
	client.MD5Sum = newMD5Sum

	// Save the updated client
	_, err = o.Update(client)
	return err
}

// IsMD5SumValid checks if the provided MD5Sum matches the one stored in the database for a given client
func IsMD5SumValid(clientName, inputMD5Sum string) (bool, error) {
	o := orm.NewOrm()

	// Use QueryTable to construct a custom query
	query := o.QueryTable(new(ClientDetails)).Filter("ClientName", clientName)

	// Get the existing client
	client := &ClientDetails{}
	err := query.One(client)
	if err == orm.ErrNoRows {
		// Client with the given ClientName not found
		return false, nil
	} else if err != nil {
		// Other error occurred
		return false, err
	}

	// Compare the input MD5Sum with the stored MD5Sum
	isValid := client.MD5Sum == inputMD5Sum

	return isValid, nil
}

// Custom function defined in the controller
func GetConnectedRoutes(inputId string) []*RouteDetails {
	// id, _ := strconv.Atoi(inputId)
	routes, _ := GetConnectedRouteDetails(inputId)
	return routes
}

// Custom function defined in the controller
func GetDisConnectedRoutes(inputId string) []*RouteDetails {
	// id, _ := strconv.Atoi(inputId)
	routes, _ := GetDisconnectedRouteDetails(inputId)
	return routes
}

// Dump any structure as json string
func Dump(obj interface{}) {
	result, _ := json.MarshalIndent(obj, "", "\t")
	logs.Debug(string(result))
}

// GetOTPKeyByClientID retrieves the OTPKey by ClientName
func GetOTPKeyByClientID(id string) (*string, error) {
	o := orm.NewOrm()
	client := &ClientDetails{Id: id}
	err := o.Read(client, "Id")
	if err == nil {
		return client.OTPKey, nil
	}
	return nil, err
}

// GetOTPDetailsByClientName retrieves Key, StaticPass and OTPUserName by ClientName
func GetOTPDetailsByClientName(clientName string) (*string, *string, *string, error) {
	o := orm.NewOrm()
	client := &ClientDetails{ClientName: clientName}
	err := o.Read(client, "ClientName")
	if err == nil {
		return client.OTPKey, client.StaticPass, client.OTPUserName, nil
	}
	return nil, nil, nil, err
}

// UpdateOTPDataByClientID updates OTPKey, StaticPass, and OTPUserName by ClientID
func UpdateOTPDataByClientId(clientId string, OTPIsEnabled bool, StaticPassIsUsed bool, otpKey, staticPass, otpUserName string) error {
	o := orm.NewOrm()
	client := &ClientDetails{Id: clientId}
	err := o.Read(client)

	if err == nil {
		if !OTPIsEnabled {
			client.OTPKey = nil
		} else {
			client.OTPKey = &otpKey
		}
		client.StaticPass = &staticPass
		client.OTPUserName = &otpUserName
		client.OTPIsEnabled = OTPIsEnabled
		client.StaticPassIsUsed = StaticPassIsUsed
		client.MD5Sum = "Auth ADDED"
		_, err := o.Update(client, "OTPKey", "OTPIsEnabled", "StaticPassIsUsed", "StaticPass", "OTPUserName", "MD5Sum")
		return err
	}
	return err
}

// DisableOTPDataByClientID by ClientID
func DisableOTPDataByClientId(clientId string) error {
	o := orm.NewOrm()
	client := &ClientDetails{Id: clientId}
	err := o.Read(client)
	if err == nil {
		client.OTPKey = nil
		client.StaticPass = nil
		client.OTPUserName = nil
		client.OTPIsEnabled = false
		client.StaticPassIsUsed = false
		client.MD5Sum = "Auth DELETED"
		_, err := o.Update(client, "OTPIsEnabled", "OTPKey", "StaticPassIsUsed", "StaticPass", "OTPUserName", "MD5Sum")
		return err
	}
	return err
}

// GetIsAuthEnabledByClientName retrieves Is2FAEnabled by ClientName
func GetIsAuthEnabledByClientName(clientName string) (bool, error) {
	o := orm.NewOrm()
	client := &ClientDetails{ClientName: clientName}
	err := o.Read(client, "ClientName")
	if err == nil {
		if client.OTPIsEnabled {
			return true, nil
		}
	}
	return false, err
}

// Reset Client Certificate By Id
func ResetClientCertificateById(clientID string) error {
	o := orm.NewOrm()

	// Get the existing client by Id
	client := &ClientDetails{Id: clientID}
	err := o.Read(client)
	if err != nil {
		return err
	}

	// Reset certificate details
	client.CertificateName = nil
	client.CertificateStatus = nil
	client.Passphrase = ""
	client.MD5Sum = "CERTIFICATE IS RESET"

	// Save the updated client
	_, err = o.Update(client)
	return err
}
