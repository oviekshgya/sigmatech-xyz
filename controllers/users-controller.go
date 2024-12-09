package controllers

import (
	beego "github.com/beego/beego/v2/server/web"
	"github.com/go-playground/validator/v10"
)

var Validator = validator.New()

// UsersController UsersController operations for UsersController
type UsersController struct {
	beego.Controller
}

var validasiError []string
