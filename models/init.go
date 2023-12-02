package models

import (
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
	orm.Debug = true
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
	hash, err := passlib.Hash(os.Getenv("OPENVPN_ADMIN_PASSWORD"))
	if err != nil {
		logs.Error("Unable to hash password", err)
	}
	user := User{
		Id:       1,
		Login:    os.Getenv("OPENVPN_ADMIN_USERNAME"),
		Name:     "Administrator",
		Email:    "root@localhost",
		Password: hash,
	}
	o := orm.NewOrm()
	if created, _, err := o.ReadOrCreate(&user, "Name"); err == nil {
		if created {
			logs.Info("Default admin account created")
		} else {
			logs.Debug(user)
		}
	}

}

func CreateDefaultSettings() (*Settings, error) {
	miAddress, err := web.AppConfig.String("OpenVpnManagementAddress")
	if err != nil {
		return nil, err
	}
	miNetwork, err := web.AppConfig.String("OpenVpnManagementNetwork")
	if err != nil {
		return nil, err
	}
	serverName, err := web.AppConfig.String("SiteName")
	if err != nil {
		return nil, err
	}
	ovConfigPath, err := web.AppConfig.String("OpenVpnPath")
	if err != nil {
		return nil, err
	}

	s := Settings{
		Profile:      "default",
		MIAddress:    miAddress,
		MINetwork:    miNetwork,
		ServerName:   serverName,
		OVConfigPath: ovConfigPath,
	}

	o := orm.NewOrm()
	if created, _, err := o.ReadOrCreate(&s, "Profile"); err == nil {
		if created {
			logs.Info("New settings profile created")
		} else {
			logs.Debug(s)
		}
		return &s, nil
	} else {
		return nil, err
	}
}

func GetBoolValueByKey(key string, m map[string]bool) bool {
	if m != nil {
		return m[key]
	}
	return false
}
