package middleware

import (
	"fmt"
	"github.com/beego/beego/v2/core/config"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/context"
	"net/http"
	"sigmatech-xyz/pkg/httpresponses"
	"strings"
)

func CORS(next beego.FilterFunc) beego.FilterFunc {
	return func(ctx *context.Context) {
		appB := httpresponses.Bee{
			Ctx: ctx,
		}
		ctx.ResponseWriter.Header().Set("Access-Control-Allow-Origin", "*")
		ctx.ResponseWriter.Header().Set("Access-Control-Allow-Credentials", "true")
		ctx.ResponseWriter.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, X-ORIGIN, X-TIMESTAMP, X-SIGNATURE, X-ACCESS")
		ctx.ResponseWriter.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")
		headerorigin := ctx.Request.Header.Get("Origin")

		if ctx.Request.Method == "OPTIONS" {
			ctx.ResponseWriter.WriteHeader(204)
			appB.Response(http.StatusNoContent, "", "CORS", nil)
			return
		}
		if Env == "dev" && !ContainsString(OriginDev, headerorigin) {
			appB.Response(http.StatusForbidden, "", "CORS ORIGIN FAILED", nil)
			return
		} else if Env == "stag" && !ContainsString(OriginStag, headerorigin) {
			appB.Response(http.StatusForbidden, "", "CORS ORIGIN FAILED", nil)
			return
		} else if Env == "prod" && !ContainsString(OriginProd, headerorigin) {
			appB.Response(http.StatusForbidden, "", "CORS ORIGIN FAILED", nil)
			return
		}

		next(ctx)
	}
}

func ContainsString(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

type OriginData struct {
	OriginNameDev   string `json:"origin_name_dev"`
	OriginNameStag  string `json:"origin_name_stag"`
	OriginNameLocal string `json:"origin_name_local"`
	OriginNameProd  string `json:"origin_name_prod"`
	ENV             string
}

var OriginDev []string
var OriginLocal []string
var OriginStag []string
var OriginProd []string
var Env string

func init() {
	data := ReadEnvOrigin()

	go func(originData OriginData) {
		originLocal := strings.Split(originData.OriginNameLocal, ",")
		if len(originLocal) > 0 {
			for i := 0; i < len(originLocal); i++ {
				OriginLocal = append(OriginLocal, originLocal[i])
			}
		}
	}(*data)

	go func(originData OriginData) {
		origindev := strings.Split(originData.OriginNameDev, ",")
		if len(origindev) > 0 {
			for i := 0; i < len(origindev); i++ {
				OriginDev = append(OriginDev, origindev[i])
			}
		}
	}(*data)

	go func(originData OriginData) {
		originstag := strings.Split(originData.OriginNameStag, ",")
		if len(originstag) > 0 {
			for i := 0; i < len(originstag); i++ {
				OriginStag = append(OriginStag, originstag[i])
			}
		}
	}(*data)

	go func(originData OriginData) {
		originprod := strings.Split(originData.OriginNameProd, ",")
		if len(originprod) > 0 {
			for i := 0; i < len(originprod); i++ {
				OriginProd = append(OriginProd, originprod[i])
			}
		}
	}(*data)

	go func(originData OriginData) {
		Env = originData.ENV
	}(*data)

}

func ReadEnvOrigin() *OriginData {
	var result OriginData
	conf, err := config.NewConfig("ini", "conf/app.conf")
	if err != nil {
		fmt.Println("Error reading config file:", err)
		return nil
	}
	result.OriginNameDev, _ = conf.String("serviceorigin.dev")
	result.OriginNameLocal, _ = conf.String("serviceorigin.local")
	result.OriginNameStag, _ = conf.String("serviceorigin.stag")
	result.OriginNameProd, _ = conf.String("serviceorigin.prod")
	result.ENV, _ = config.String("serviceorigin-env")

	return &result
}
