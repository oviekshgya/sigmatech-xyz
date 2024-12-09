package routers

import (
	beego "github.com/beego/beego/v2/server/web"
	"sigmatech-xyz/middleware"
	"time"
)

func init() {
	ipLimiter := middleware.NewIPLimiter(10, time.Minute) // Maksimum 10 permintaan per menit
	go ipLimiter.CleanUp()
	beego.InsertFilter("/*", beego.BeforeRouter, middleware.MiddlewareIPLimiter(ipLimiter))

	beego.InsertFilterChain("/*", middleware.CORS)
	beego.InsertFilterChain("/*", middleware.AuthHeader)
	beego.InsertFilterChain("/*", middleware.AuthBasic)
	beego.InsertFilterChain("/*", middleware.AuthorizeBaererJWT)

	ns := beego.NewNamespace("/v1") //beego.NSNamespace("/user",
	//	beego.NSInclude(&controllers.UsersController{}),
	//),

	beego.AddNamespace(ns)

	//user.POST("/registrasi", controllers.UsersController.NewRegistrasi)

}
