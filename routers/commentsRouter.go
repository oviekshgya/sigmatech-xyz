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

	beego.GlobalControllerRouter["sigmatech-xyz/controllers:UsersController"] = append(beego.GlobalControllerRouter["sigmatech-xyz/controllers:UsersController"],
		beego.ControllerComments{
			Method:           "ValidasiOTPCustomer",
			Router:           `/validasi-otp`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["sigmatech-xyz/controllers:UsersController"] = append(beego.GlobalControllerRouter["sigmatech-xyz/controllers:UsersController"],
		beego.ControllerComments{
			Method:           "Login",
			Router:           `/login`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["sigmatech-xyz/controllers:UsersController"] = append(beego.GlobalControllerRouter["sigmatech-xyz/controllers:UsersController"],
		beego.ControllerComments{
			Method:           "VerifikasiAkun",
			Router:           `/verifikasi`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["sigmatech-xyz/controllers:UsersController"] = append(beego.GlobalControllerRouter["sigmatech-xyz/controllers:UsersController"],
		beego.ControllerComments{
			Method:           "Profile",
			Router:           `/profile`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["sigmatech-xyz/controllers:UsersController"] = append(beego.GlobalControllerRouter["sigmatech-xyz/controllers:UsersController"],
		beego.ControllerComments{
			Method:           "MasterMerchant",
			Router:           `/master-merchant`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["sigmatech-xyz/controllers:UsersController"] = append(beego.GlobalControllerRouter["sigmatech-xyz/controllers:UsersController"],
		beego.ControllerComments{
			Method:           "SimulasiTransaksi",
			Router:           `/simulasi-transaksi`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

}
