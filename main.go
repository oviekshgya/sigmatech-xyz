package main

import (
	"fmt"
	beego "github.com/beego/beego/v2/server/web"
	config "sigmatech-xyz/conf"
	database "sigmatech-xyz/db"
	"sigmatech-xyz/pkg"
	_ "sigmatech-xyz/routers"
)

func main() {
	connRedis := pkg.InitializeRedis()
	if connRedis == nil {
		fmt.Println("DISCONNECT REDIS. INSTALL OR START REDIS BEFORE RUN")
		return
	}
	fmt.Printf("REDIS CONNECTED")
	database.ConectionGORM()
	//go pkg.RunScriptShellBuild()
	beego.BConfig.RunMode = config.ReadEnv().Server.RunMode
	beego.Run("localhost:" + config.ReadEnv().Server.HttpPort)
}
