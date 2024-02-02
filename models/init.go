package models

import (
	"github.com/google/uuid"
	"os"

	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/server/web"
	"gopkg.in/hlandau/passlib.v1"
)

func InitDB() {
	err := orm.RegisterDriver("sqlite3", orm.DRSqlite)
	if err != nil {
		panic(err)
	}
	dbPath, err := web.AppConfig.String("dbPath")
	if err != nil {
		panic(err)
	}
	dbSource := "file:" + dbPath

	err = orm.RegisterDataBase("default", "sqlite3", dbSource)
	if err != nil {
		panic(err)
	}
	//orm.Debug = true
	orm.RegisterModel(
		new(User),
		new(Settings),
		new(OVConfig),
		new(RouteDetails),
		new(ClientDetails),
	)

	err = orm.RunSyncdb("default", false, true)
	if err != nil {
		logs.Error(err)
		return
	}
}

func CreateDefaultUsers() {
	// Check if the user already exists
	o := orm.NewOrm()
	user := User{Name: "Administrator"}
	if err := o.Read(&user, "Name"); err == nil {
		logs.Info("Default admin account already exists")
		return
	}

	// If the user doesn't exist, create it
	hash, err := passlib.Hash(os.Getenv("OPENVPN_ADMIN_PASSWORD"))
	if err != nil {
		logs.Error("Unable to hash password", err)
		return
	}

	user = User{
		Id:       uuid.New().String(),
		Login:    os.Getenv("OPENVPN_ADMIN_USERNAME"),
		Name:     "Administrator",
		Email:    "root@localhost",
		Password: hash,
	}

	if _, err := o.Insert(&user); err == nil {
		logs.Info("Default admin account created")
	} else {
		logs.Error("Error creating default admin account", err)
	}
}

func CreateDefaultSettings() (*Settings, error) {
	// Read configuration values
	miAddress := getConfigString("OpenVpnManagementAddress")
	miNetwork := getConfigString("OpenVpnManagementNetwork")
	serverName := getConfigString("SiteName")
	ovConfigPath := getConfigString("OpenVpnPath")

	// Create settings object
	s := Settings{
		Id:           uuid.New().String(),
		Profile:      "default",
		MIAddress:    miAddress,
		MINetwork:    miNetwork,
		ServerName:   serverName,
		OVConfigPath: ovConfigPath,
	}

	// Check if the settings already exist
	o := orm.NewOrm()
	if err := o.Read(&s, "Profile"); err == nil {
		logs.Info("Settings profile already exists")
		return &s, nil
	}

	// If the settings don't exist, create them
	if _, _, err := o.ReadOrCreate(&s, "Profile"); err == nil {
		logs.Info("New settings profile created")
		return &s, nil
	} else {
		return nil, err
	}
}

// getConfigString is a helper function to read a string configuration value
func getConfigString(key string) string {
	value, err := web.AppConfig.String(key)
	if err != nil {
		logs.Error("Error reading configuration:", err)
		// You might want to handle the error accordingly, e.g., log and return a default value
	}
	return value
}

func GetBoolValueByKey(key string, m map[string]bool) bool {
	if m != nil {
		return m[key]
	}
	return false
}
