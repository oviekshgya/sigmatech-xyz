package repositories

import (
	"fmt"
	database "sigmatech-xyz/db"
	"sigmatech-xyz/models"
	"sigmatech-xyz/models/users"
	"sigmatech-xyz/pkg"
	"sync"
	"time"

	beego "github.com/beego/beego/v2/server/web"
	"gorm.io/gorm"
)

type userRepositories struct {
	beego.Controller
	DbMain *gorm.DB
}

func StaticUserRepositoris() *userRepositories {
	return &userRepositories{
		DbMain: database.DBMain,
	}
}

type RequestOTPJobs struct {
	Result chan interface{}
	Error  chan error
}

const (
	NumWorkers = 5
	QueueSizes = 100
)

func (service userRepositories) RequestOTPCustomer(input models.JSONRequestOTP, jobs <-chan RequestOTPJobs, wg *sync.WaitGroup) {
	for i := 1; i <= NumWorkers; i++ {
		go func() {
			for job := range jobs {
				tx := service.DbMain.Begin()
				redisCon := pkg.InitializeRedis()
				kode := pkg.KodeVerify(6)
				key := fmt.Sprintf("%s-%s", input.Request.Email, "requestotp-registrasi")
				defer wg.Done()
				var idAkun int
				if tx.Table(pkg.AKUNCUSTOMER).Where("email = ? AND isActive = 1", input.Request.Email).Select("idAkun").Scan(&idAkun); idAkun != 0 {
					job.Error <- fmt.Errorf("akun hp %s already exists", input.Request.Email)
					return
				}
				var dataRedis map[string]int
				if get := redisCon.GetKey(key, &dataRedis); get != nil {
					fmt.Println("error get redis", get)
				}
				if dataRedis["email"] > 2 {
					job.Error <- fmt.Errorf("akun sudah melakukan request otp sebanyak 3x tunggu 30menit untuk mengirim otp ulang")
					return
				}

				if updated := tx.Table(pkg.OTPCUSTOMERS).Where("email = ? AND isUsed = 0 ", input.Request.Email).Updates(map[string]interface{}{
					"isUsed": 9,
				}); updated.Error != nil {
					job.Error <- fmt.Errorf("resend update otp error:%s", updated.Error.Error())
					return
				}

				if created := tx.Create(&users.OTPCustomer{
					Email:     input.Request.Email,
					IsUsed:    0,
					ExpiredAt: time.Now().Add(30 * time.Minute),
					Kode:      kode,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}); created.Error != nil {
					tx.Rollback()
					job.Error <- fmt.Errorf("insert error:%s", created.Error)
					return
				}

				if set := redisCon.SetKey(key, map[string]int{
					"email": dataRedis["email"] + 1,
				}, time.Duration(30*time.Minute)); set != nil {
					fmt.Println("error set", set.Error())
				}
				tx.Commit()
				job.Result <- map[string]interface{}{
					"kode": kode,
				}

			}

		}()
	}

}
