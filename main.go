package main

import (
	"flag"
	"fmt"
	"path/filepath"

	"github.com/beego/beego/v2/server/web"
	"github.com/shuricksumy/openvpn-ui/lib"
	"github.com/shuricksumy/openvpn-ui/models"
	"github.com/shuricksumy/openvpn-ui/routers"
	"github.com/shuricksumy/openvpn-ui/state"
)

func main() {
	// Start OpenVPN
	err := lib.StartOpenVPN()
	if err != nil {
		fmt.Println("Error starting OpenVPN:", err)
	}

	err_fw := lib.EnableFWRules()
	if err_fw != nil {
		fmt.Println("Error apply FireWall rules:", err)
	}

	configDir := flag.String("config", "conf", "Path to config dir")
	flag.Parse()

	configFile := filepath.Join(*configDir, "app.conf")
	fmt.Println("Config file:", configFile)

	if err := web.LoadAppConfig("ini", configFile); err != nil {
		panic(err)
	}

	models.InitDB()
	models.CreateDefaultUsers()
	defaultSettings, err := models.CreateDefaultSettings()
	if err != nil {
		panic(err)
	}

	state.GlobalCfg = *defaultSettings

	routers.Init(*configDir)

	lib.AddFuncMaps()

	//Add custom functions to web templates
	web.AddFuncMap("RouteIsUsedBy", models.RouteIsUsedBy)
	web.AddFuncMap("GetConnectedRoutes", models.GetConnectedRoutes)
	web.AddFuncMap("GetDisConnectedRoutes", models.GetDisConnectedRoutes)
	web.AddFuncMap("GetBoolValueByKey", models.GetBoolValueByKey)

	web.Run()
}
