package repositories

import (
	"bytes"
	"fmt"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/context"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"io/ioutil"
	"net/http/httputil"
	"sigmatech-xyz/conf"
	database "sigmatech-xyz/db"
	"sigmatech-xyz/models/logging"
	"time"
)

type logRepositories struct {
	beego.Controller
	DbMain *gorm.DB
}

func StaticLogRepositories() *logRepositories {
	return &logRepositories{
		DbMain: database.DBMain,
	}
}

func Checkendpoint(endpoint string, bodyform map[string]string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(bodyform["bodyform"]), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println("err log ", err.Error())
	}
	if endpoint == "/user/login" {
		bodyform["bodyform"] = string(hash)
	}
	return bodyform["bodyform"]
}

func (service logRepositories) LogCreated(id int, bodyrespon map[string]string, statuscode int, c *context.Context) (interface{}, error) {
	tx := service.DbMain.Begin()
	body, _ := ioutil.ReadAll(c.Request.Body)
	requestDumpHeader, _ := httputil.DumpRequest(c.Request, true)

	create := logging.MasterLog{
		Apps:      conf.ReadEnv().Server.AppName,
		IdUser:    id,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Method:    c.Request.Method,
		Tanggal:   time.Now(),
		Header:    string(requestDumpHeader),
		Body:      Checkendpoint(c.Request.RequestURI, bodyrespon),
		Endpoint:  c.Request.RequestURI,
		Ip:        "",
		MacAddr:   "",
		Response:  bodyrespon["bodyrespons"],
	}

	if created := tx.Table("log").Create(&create); created.Error != nil {
		tx.Rollback()
		return false, nil
	}
	tx.Commit()
	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer([]byte(body))) // Write body back
	return true, nil
}
