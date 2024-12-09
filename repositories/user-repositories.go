package repositories

import (
	database "sigmatech-xyz/db"

	beego "github.com/beego/beego/v2/server/web"
	"gorm.io/gorm"
)

type userRepositories struct {
	beego.Controller
	DbMain *gorm.DB
}

func StaticUserRepositories() *userRepositories {
	return &userRepositories{
		DbMain: database.DBMain,
	}
}
