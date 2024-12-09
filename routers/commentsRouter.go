package routers

import (
	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/context/param"
)

func init() {

	beego.GlobalControllerRouter["sigmatech-xyz/controllers:UsersController"] = append(beego.GlobalControllerRouter["sigmatech-xyz/controllers:UsersController"],
		beego.ControllerComments{
			Method:           "RequestOTP",
			Router:           `/request-otp`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

}
