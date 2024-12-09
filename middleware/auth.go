package middleware

import (
	"fmt"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/context"
	"net/http"
	"sigmatech-xyz/pkg"
	"sigmatech-xyz/pkg/auth"
	"sigmatech-xyz/pkg/httpresponses"
)

func IsUserLoggedIn(ctx *context.Context) bool {
	return true
}

func AuthBasic(next beego.FilterFunc) beego.FilterFunc {
	return func(ctx *context.Context) {
		appG := httpresponses.Bee{
			Ctx: ctx,
		}
		//orm := orm2.NewOrmUsingDB("default")
		//var dataAuth models.Authorization
		//orm.QueryTable(new(models.Authorization)).Filter("isActive", 1).All(&dataAuth)
		//fmt.Println("AUTH:", dataAuth.Username, dataAuth.Password)
		user, password, hasAuth := ctx.Request.BasicAuth()
		if hasAuth && user == pkg.USERNAME && password == pkg.PASSWORD {
			fmt.Println("NEXT BASIC AUTH")
			//this.next.ServeHTTP(w, r)
			//this.next.ServeHTTP()
		} else {
			appG.Response(http.StatusBadRequest, "", "Basic Auth Failed", nil)
		}
		next(ctx)
	}
}

func AuthHeader(next beego.FilterFunc) beego.FilterFunc {
	return func(ctx *context.Context) {
		appG := httpresponses.Bee{
			Ctx: ctx,
		}
		key := ctx.Request.Header.Get("X-ORIGIN")
		if key != pkg.HEADERKEY {
			fmt.Println("header", key)
			appG.Response(http.StatusBadRequest, "", "invalid header", nil)
			return
		}
		next(ctx)

	}
}

func AuthorizeBaererJWT(next beego.FilterFunc) beego.FilterFunc {
	return func(ctx *context.Context) {
		appG := httpresponses.Bee{
			Ctx: ctx,
		}
		const BEARER_SCHEMA = "Bearer "
		authHeader := ctx.Request.Header.Get("Authorization")
		if authHeader == "" {
			appG.Response(http.StatusUnauthorized, "", "Invalide Authorization", nil)
			//ctx.Abort(http.StatusUnauthorized, "Invalid Auth")
			return
		}
		tokenString := authHeader[len(BEARER_SCHEMA):]
		if tokenString == "" {
			appG.Response(http.StatusUnauthorized, "", "Invalide Authorization", nil)
			//ctx.Abort(http.StatusUnauthorized, "Invalid Auth")
			return
		}

		metadata, errs := auth.ExtractedExt(ctx.Request, "")
		if metadata == nil {
			appG.Response(http.StatusUnauthorized, "", errs.Error(), metadata)
			//ctx.Abort(http.StatusUnauthorized, "")
			return
		}
		next(ctx)
	}
}
