package repositories

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
	database "sigmatech-xyz/db"
	"sigmatech-xyz/models"
	"sigmatech-xyz/models/master"
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

type ValidasiOTPJobs struct {
	Result chan interface{}
	Error  chan error
}

func (service userRepositories) ValidasiOTPCustomers(input models.JSONValidasiOTP, jobs <-chan ValidasiOTPJobs, wg *sync.WaitGroup) {
	for i := 1; i <= NumWorkers; i++ {
		go func() {
			for job := range jobs {
				tx := service.DbMain.Begin()
				defer wg.Done()
				var data users.OTPCustomer
				tx.Table(pkg.OTPCUSTOMERS).Where("email = ? AND isUsed = 0 AND kode = ? ", input.Request.Email, input.Request.Otp).Select("idOtp, expiredAt").Take(&data)
				if data.IdOtp == 0 {
					tx.Rollback()
					job.Error <- fmt.Errorf("kode otp salah")
					return
				}

				if time.Now().After(data.ExpiredAt) {
					tx.Rollback()
					job.Error <- fmt.Errorf("kode otp expired")
					return
				}

				if updated := tx.Table(pkg.OTPCUSTOMERS).Where("email = ? AND isUsed = 0 AND kode = ?", input.Request.Email, input.Request.Otp).Updates(map[string]interface{}{
					"isUsed": 1,
				}); updated.Error != nil {
					tx.Rollback()
					job.Error <- fmt.Errorf("resend update otp error:%s", updated.Error.Error())
					return
				}

				hash, err := bcrypt.GenerateFromPassword([]byte(input.Request.Password), bcrypt.DefaultCost)
				if err != nil {
					tx.Rollback()
					job.Error <- fmt.Errorf("hash failed:%s", err.Error())
				}

				if created := tx.Create(&users.AkunCustomer{
					Email:     input.Request.Email,
					Hp:        input.Request.Hp,
					UpdatedAt: time.Now(),
					CreatedAt: time.Now(),
					IsActive:  1,
					Password:  string(hash),
				}); created.Error != nil {
					tx.Rollback()
					job.Error <- fmt.Errorf("created akun error:%s", created.Error.Error())
					return
				}

				tx.Commit()
				job.Result <- map[string]interface{}{
					"kode": true,
				}

			}

		}()
	}

}

func (service userRepositories) Login(input models.JSONLogin) (map[string]interface{}, error) {
	var data users.AkunCustomer
	service.DbMain.Where("email = ?", input.Request.Email).First(&data)
	if data.IdAkun != 0 {
		hash := []byte(data.Password)
		if err := bcrypt.CompareHashAndPassword(hash, []byte(input.Request.Password)); err != nil {
			return nil, fmt.Errorf("email / password salah")
		}
		return map[string]interface{}{
			"idAkun": data.IdAkun,
		}, nil
	}
	return nil, fmt.Errorf("email / password salah")
}

type VerifikasiJob struct {
	Result chan interface{}
	Error  chan error
}

func (service userRepositories) VerifikasiAkun(idAkun int, input models.JSONVerifikasi, job <-chan VerifikasiJob, wg *sync.WaitGroup) {
	for i := 1; i <= NumWorkers; i++ {
		go func() {
			for jobs := range job {
				tx := service.DbMain.Begin()
				defer wg.Done()
				if created := tx.Create(&users.UserCustomer{
					IdAkun:       idAkun,
					CreatedAt:    time.Now(),
					UpdatedAt:    time.Now(),
					Nik:          input.Request.Nik,
					FotoKtp:      input.Request.FotoKtp,
					FotoSelfie:   input.Request.FotoSelfie,
					IsAktivasi:   0,
					LegalName:    input.Request.LegalName,
					Salary:       input.Request.Salary,
					TanggalLahir: time.Now(),
					TempatLahir:  input.Request.TempatLahir,
				}); created.Error != nil {
					tx.Rollback()
					jobs.Error <- fmt.Errorf("insert error:%s", created.Error.Error())
					return
				}
				tx.Commit()
				jobs.Result <- map[string]interface{}{
					"insert": true,
				}
				return
			}
		}()
	}
}

type ProfileJob struct {
	Result chan interface{}
	Error  chan error
}

func (service userRepositories) Profile(idAkun int, job <-chan ProfileJob, wg *sync.WaitGroup) {
	for i := 1; i <= NumWorkers; i++ {
		go func() {
			for jobs := range job {
				defer wg.Done()
				var data users.UserCustomerData
				service.DbMain.Table(pkg.AKUNCUSTOMER).Preload("Data").Where("idAkun = ?", idAkun).First(&data)
				jobs.Result <- data
				return
			}
		}()
	}
}

func (service userRepositories) MasterMerchant(page, pageSize int) (interface{}, error) {
	var data []master.MasterMerchants
	var totalData int64
	var totalPage int

	switch {
	case pageSize > 100:
		pageSize = pageSize
	case pageSize <= 0:
		pageSize = 10
	}

	service.DbMain.Scopes(models.Pagination(pageSize, page)).Where("isActive = 1").Order("namaMerchant DESC").Find(&data)
	service.DbMain.Table(pkg.MERCHANT).Where("isActive = 1").Count(&totalData)

	if int(totalData) < pageSize {
		totalPage = 1
	} else {
		totalPage = int(totalData) / pageSize
		if (int(totalData) % pageSize) != 0 {
			totalPage = totalPage + 1
		}
	}

	if page == 0 {
		page = 1
	}

	return map[string]interface{}{
		"page":      page,
		"pageSize":  pageSize,
		"data":      data,
		"total":     totalData,
		"totalPage": totalPage,
	}, nil
}

func (service userRepositories) SimulasiTransaksi(idAkun int, input models.JSONTransaksiSimulasiPengajuan) (interface{}, error) {
	if input.Request.Tenor != 3 || input.Request.Tenor != 5 {
		return nil, fmt.Errorf("tenor hanya boleh 3 atau 5 tahun")
	}

	var dataMerchant master.MasterMerchants
	service.DbMain.Where("isActive = 1").Order("namaMerchant DESC").First(&dataMerchant)

	var dataRate master.MasterRates
	service.DbMain.Where("isActive = 1").First(&dataRate)

	var paymentSchedule []map[string]interface{}

	var a = 1
	tgljt := time.Now().AddDate(0, 0, 30)
	for i := 0; i < input.Request.Tenor; i++ {

		paymentSchedule = append(paymentSchedule, map[string]interface{}{
			"angsuranKe":    a + 1,
			"bunga":         dataRate.Rate / 12,
			"tanggal":       time.Now(),
			"tglJatuhTempo": tgljt,
			"jumlahCicilan": input.Request.OTR / float64(input.Request.Tenor),
			"totalCicilan":  (input.Request.OTR / float64(input.Request.Tenor)) + input.Request.OTR*((dataRate.Rate/12)/100)*float64(input.Request.Tenor),
			"status":        "ready",
		})
		tgljt = tgljt.AddDate(0, 0, 30)
	}

	return map[string]interface{}{
		"model":              input.Request.NamaAset,
		"merchan":            dataMerchant.NamaMerchant,
		"totalHutang":        (input.Request.OTR * (dataRate.Rate / 100) * float64(input.Request.Tenor)) + input.Request.OTR,
		"tenor":              input.Request.Tenor,
		"bunga":              input.Request.OTR * (dataRate.Rate / 100) * float64(input.Request.Tenor),
		"tglJatuhTempoBulan": time.Now().AddDate(0, 0, 30).Format("2006-01-02"),
		"tglJatuhTempo":      time.Now().AddDate(0, input.Request.Tenor, 0).Format("2006-01-02"),
		"paymanetSchedule":   paymentSchedule,
		"rate":               dataRate.Rate,
	}, nil
}
