package conf

import (
	"fmt"

	"github.com/beego/beego/v2/core/config"
)

type EnvServer struct {
	AppName          string
	HttpPort         string
	RunMode          string
	AutoRender       bool
	CopyRequestBody  bool
	EnableDocs       bool
	GRPCPort         string
	IvKeysDashboard  string
	KeyAesDashboadrd string
}

type EnvDB struct {
	Driver   string
	User     string
	Password string
	Host     string
	Name     string
}

type Env struct {
	Server EnvServer
	DB     EnvDB
	Ubuntu ConfigUbuntu
}

type ConfigUbuntu struct {
	IP   string
	Pass string
	User string
}

func ReadEnv() *Env {
	var result Env
	conf, err := config.NewConfig("ini", "conf/app.conf")
	if err != nil {
		fmt.Println("Error reading config file:", err)
		return nil
	}
	result.DB.Driver, _ = conf.String("db.driver")
	result.DB.User, _ = conf.String("db.user")
	result.DB.Password, _ = conf.String("db.password")
	result.DB.Host, _ = conf.String("db.url")
	result.DB.Name, _ = conf.String("db.name")

	result.Server.AppName, _ = config.String("service-name")
	result.Server.HttpPort, _ = config.String("service-port")
	result.DB.Name, _ = config.String("sqlconn")
	result.Server.RunMode, _ = config.String("service-env")
	result.Server.KeyAesDashboadrd, _ = config.String("service-dashboardkeyaes")
	result.Server.IvKeysDashboard, _ = config.String("service-dashboardivkeys")

	return &result
}
