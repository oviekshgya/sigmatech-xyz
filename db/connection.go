package db

import (
	"fmt"
	"github.com/beego/beego/v2/core/config"
	_ "github.com/go-sql-driver/mysql"
	grmsql "gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"sigmatech-xyz/models/master"
	"sigmatech-xyz/models/transaksi"
	"sigmatech-xyz/models/users"
	"time"
)

type EnvRead struct {
	DbDrive    string
	DbUser     string
	DbPassword string
	DBHost     string
	DBName     string
}

func init() {

}

var DBMain *gorm.DB

func ConectionGORM() *gorm.DB {
	newLogger :=
		logger.New(

			log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
			logger.Config{
				LogLevel:                  logger.Info,
				IgnoreRecordNotFoundError: true, // Ignore ErrRecordNotFound error for logger
				Colorful:                  true, // Enable colorful output// Log level: logger.Silent, logger.Error, logger.Warn, logger.Info
			},
		)
	conf, err := config.NewConfig("ini", "conf/app.conf")
	if err != nil {
		fmt.Println("Error reading config file:", err)
		return nil
	}
	//dbDriver, _ := conf.String("db.driver")
	dbUser, _ := conf.String("db.user")
	dbPassword, _ := conf.String("db.password")
	dbURL, _ := conf.String("db.url")
	dbName, _ := conf.String("db.name")
	addr := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", dbUser, dbPassword, dbURL, "3306", dbName)
	DBMain, err = gorm.Open(grmsql.Open(addr), &gorm.Config{
		Logger: newLogger,
	})

	DBMain.AutoMigrate(&users.UserCustomer{}, &users.AkunCustomer{}, &transaksi.Transaksi{}, &transaksi.PaymentTransaksi{}, &master.MasterMerchant{}, &master.MasterRate{})

	// Konfigurasi pool koneksi
	sqlDB, _ := DBMain.DB()

	// Maksimal jumlah koneksi yang bisa dibuka
	sqlDB.SetMaxOpenConns(50)

	// Maksimal jumlah koneksi idle
	sqlDB.SetMaxIdleConns(10)

	// Maksimal durasi koneksi bisa digunakan
	sqlDB.SetConnMaxLifetime(time.Minute * 30)

	//defer sqlDB.Close()

	if err != nil {
		log.Println("err connect db main:", err)
	}
	fmt.Println("CONNECT MAIN DB GORM")

	//SetSessionWithDB(db, "maindb")
	return DBMain
}
