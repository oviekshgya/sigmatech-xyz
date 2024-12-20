package controllers

import (
	"encoding/json"
	"fmt"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/go-playground/validator/v10"
	"net/http"
	"sigmatech-xyz/models"

	"sigmatech-xyz/pkg"
	"sigmatech-xyz/pkg/auth"
	"strconv"

	"sigmatech-xyz/pkg/httpresponses"
	"sigmatech-xyz/repositories"
	"sync"
)

var Validator = validator.New()

// UsersController UsersController operations for UsersController
type UsersController struct {
	beego.Controller
}

var validasiError []string

// RequestOTP
// @Description RequestOTP
// @Param	body	nil	true	"body nill"
// @Success 200 {int} interfaces{}
// @Failure 403 bodies are empty
// @router /request-otp [post]
func (controller UsersController) RequestOTP() {
	appB := httpresponses.Bee{
		Ctx: controller.Ctx,
	}
	var dataLogin models.JSONRequestOTP
	if err := json.Unmarshal([]byte((controller.Ctx.Input.GetData("RequestBody").(string))), &dataLogin); err != nil {
		fmt.Println("err", err)
		appB.Response(http.StatusUnprocessableEntity, "", err.Error(), nil)
		return
	}

	jobs := make(chan repositories.RequestOTPJobs, repositories.QueueSizes)
	results := make(chan interface{}, repositories.QueueSizes)
	errChan := make(chan error, 1)
	var wg sync.WaitGroup

	go repositories.StaticUserRepositoris().RequestOTPCustomer(dataLogin, jobs, &wg)

	wg.Add(1)
	jobs <- repositories.RequestOTPJobs{
		Result: results,
		Error:  errChan,
	}
	close(jobs)
	wg.Wait()

	select {
	case res := <-results:
		appB.Response(http.StatusCreated, "Success", "", res)
	case err := <-errChan:
		appB.Response(http.StatusInternalServerError, "", err.Error(), nil)
	default:
		appB.Response(http.StatusBadRequest, "Unknown Error", "", nil)
	}

}

// ValidasiOTPCustomer
// @Description ValidasiOTPCustomer
// @Param	body	nil	true	"body nill"
// @Success 200 {int} interfaces{}
// @Failure 403 bodies are empty
// @router /validasi-otp [post]
func (controller UsersController) ValidasiOTPCustomer() {
	appB := httpresponses.Bee{
		Ctx: controller.Ctx,
	}
	var input models.JSONValidasiOTP
	if err := json.Unmarshal([]byte((controller.Ctx.Input.GetData("RequestBody").(string))), &input); err != nil {
		fmt.Println("err", err)
		appB.Response(http.StatusUnprocessableEntity, "", err.Error(), nil)
		return
	}

	jobs := make(chan repositories.ValidasiOTPJobs, repositories.QueueSizes)
	results := make(chan interface{}, repositories.QueueSizes)
	errChan := make(chan error, 1)
	var wg sync.WaitGroup

	go repositories.StaticUserRepositoris().ValidasiOTPCustomers(input, jobs, &wg)

	wg.Add(1)
	jobs <- repositories.ValidasiOTPJobs{
		Result: results,
		Error:  errChan,
	}
	close(jobs)
	wg.Wait()

	select {
	case res := <-results:
		appB.Response(http.StatusCreated, "Success", "", res)
	case err := <-errChan:
		appB.Response(http.StatusInternalServerError, "", err.Error(), nil)
	default:
		appB.Response(http.StatusBadRequest, "Unknown Error", "", nil)
	}
}

// Login
// @Description Login
// @Param	body	nil	true	"body nill"
// @Success 200 {int} interfaces{}
// @Failure 403 bodies are empty
// @router /login [post]
func (controller UsersController) Login() {
	appB := httpresponses.Bee{
		Ctx: controller.Ctx,
	}
	var input models.JSONLogin
	if err := json.Unmarshal([]byte((controller.Ctx.Input.GetData("RequestBody").(string))), &input); err != nil {
		fmt.Println("err", err)
		appB.Response(http.StatusUnprocessableEntity, "", err.Error(), nil)
		return
	}

	result, err2 := repositories.StaticUserRepositoris().Login(input)
	if err2 != nil {
		appB.Response(http.StatusBadRequest, "", "email / password salah", nil)
		return
	}

	create, errCreate := pkg.CreateToken(result["idAkun"].(int))
	if err2 != nil {
		appB.Response(http.StatusBadRequest, "", errCreate.Error(), nil)
		return
	}
	appB.Response(http.StatusOK, "Success", "", map[string]interface{}{
		"accessToken":  create.AccessToken,
		"refrestToken": create.RefreshToken,
	})
	return
}

// VerifikasiAkun
// @Description VerifikasiAkun
// @Param	body	nil	true	"body nill"
// @Success 200 {int} interfaces{}
// @Failure 403 bodies are empty
// @router /verifikasi [post]
func (controller UsersController) VerifikasiAkun() {
	appB := httpresponses.Bee{
		Ctx: controller.Ctx,
	}

	meta, errMeta := auth.ExtractedExt(controller.Ctx.Request, "")
	if errMeta != nil {
		appB.Response(http.StatusUnauthorized, "", errMeta.Error(), nil)
		return
	}
	var input models.JSONVerifikasi
	if err := json.Unmarshal([]byte((controller.Ctx.Input.GetData("RequestBody").(string))), &input); err != nil {
		fmt.Println("err", err)
		appB.Response(http.StatusUnprocessableEntity, "", err.Error(), nil)
		return
	}

	jobs := make(chan repositories.VerifikasiJob, repositories.QueueSizes)
	results := make(chan interface{}, repositories.QueueSizes)
	errChan := make(chan error, 1)
	var wg sync.WaitGroup

	go repositories.StaticUserRepositoris().VerifikasiAkun(meta.Id, input, jobs, &wg)

	wg.Add(1)
	jobs <- repositories.VerifikasiJob{
		Result: results,
		Error:  errChan,
	}
	close(jobs)
	wg.Wait()

	select {
	case res := <-results:
		appB.Response(http.StatusCreated, "Success", "", res)
	case err := <-errChan:
		appB.Response(http.StatusInternalServerError, "", err.Error(), nil)
	default:
		appB.Response(http.StatusBadRequest, "Unknown Error", "", nil)
	}
}

// Profile
// @Description Profile
// @Param	body	nil	true	"body nill"
// @Success 200 {int} interfaces{}
// @Failure 403 bodies are empty
// @router /profile [post]
func (controller UsersController) Profile() {
	appB := httpresponses.Bee{
		Ctx: controller.Ctx,
	}

	meta, errMeta := auth.ExtractedExt(controller.Ctx.Request, "")
	if errMeta != nil {
		appB.Response(http.StatusUnauthorized, "", errMeta.Error(), nil)
		return
	}
	jobs := make(chan repositories.ProfileJob, repositories.QueueSizes)
	results := make(chan interface{}, repositories.QueueSizes)
	errChan := make(chan error, 1)
	var wg sync.WaitGroup

	go repositories.StaticUserRepositoris().Profile(meta.Id, jobs, &wg)

	wg.Add(1)
	jobs <- repositories.ProfileJob{
		Result: results,
		Error:  errChan,
	}
	close(jobs)
	wg.Wait()

	select {
	case res := <-results:
		appB.Response(http.StatusCreated, "Success", "", res)
	case err := <-errChan:
		appB.Response(http.StatusInternalServerError, "", err.Error(), nil)
	default:
		appB.Response(http.StatusBadRequest, "Unknown Error", "", nil)
	}
}

// MasterMerchant
// @Description MasterMerchant
// @Param	body	nil	true	"body nill"
// @Success 200 {int} interfaces{}
// @Failure 403 bodies are empty
// @router /master-merchant [get]
func (controller UsersController) MasterMerchant() {
	appB := httpresponses.Bee{
		Ctx: controller.Ctx,
	}

	_, errMeta := auth.ExtractedExt(controller.Ctx.Request, "")
	if errMeta != nil {
		appB.Response(http.StatusUnauthorized, "", errMeta.Error(), nil)
		return
	}
	page, _ := strconv.Atoi(controller.Ctx.Input.Query("page"))
	pageSize, _ := strconv.Atoi(controller.Ctx.Input.Query("pageSize"))
	result, err2 := repositories.StaticUserRepositoris().MasterMerchant(page, pageSize)
	if err2 != nil {
		appB.Response(http.StatusBadRequest, "", "email / password salah", nil)
		return
	}

	appB.Response(http.StatusOK, "Success", "", result)
	return
}

// SimulasiTransaksi
// @Description SimulasiTransaksi
// @Param	body	nil	true	"body nill"
// @Success 200 {int} interfaces{}
// @Failure 403 bodies are empty
// @router /simulasi-transaksi [post]
func (controller UsersController) SimulasiTransaksi() {
	appB := httpresponses.Bee{
		Ctx: controller.Ctx,
	}

	meta, errMeta := auth.ExtractedExt(controller.Ctx.Request, "")
	if errMeta != nil {
		appB.Response(http.StatusUnauthorized, "", errMeta.Error(), nil)
		return
	}

	var input models.JSONTransaksiSimulasi
	if err := json.Unmarshal([]byte((controller.Ctx.Input.GetData("RequestBody").(string))), &input); err != nil {
		fmt.Println("err", err)
		appB.Response(http.StatusUnprocessableEntity, "", err.Error(), nil)
		return
	}

	result, err2 := repositories.StaticUserRepositoris().SimulasiTransaksi(meta.Id, input)
	if err2 != nil {
		appB.Response(http.StatusBadRequest, "", err2.Error(), nil)
		return
	}

	appB.Response(http.StatusOK, "Success", "", result)
	return
}

// Transaksi
// @Description Transaksi
// @Param	body	nil	true	"body nill"
// @Success 200 {int} interfaces{}
// @Failure 403 bodies are empty
// @router /transaksi [post]
func (controller UsersController) Transaksi() {
	appB := httpresponses.Bee{
		Ctx: controller.Ctx,
	}

	meta, errMeta := auth.ExtractedExt(controller.Ctx.Request, "")
	if errMeta != nil {
		appB.Response(http.StatusUnauthorized, "", errMeta.Error(), nil)
		return
	}

	var input models.JSONTransaksiPinjaman
	if err := json.Unmarshal([]byte((controller.Ctx.Input.GetData("RequestBody").(string))), &input); err != nil {
		fmt.Println("err", err)
		appB.Response(http.StatusUnprocessableEntity, "", err.Error(), nil)
		return
	}

	result, err2 := repositories.StaticUserRepositoris().Transaksi(meta.Id, input)
	if err2 != nil {
		appB.Response(http.StatusBadRequest, "", err2.Error(), nil)
		return
	}

	appB.Response(http.StatusOK, "Success", "", result)
	return
}

// CheckPengajuan
// @Description CheckPengajuan
// @Param	body	nil	true	"body nill"
// @Success 200 {int} interfaces{}
// @Failure 403 bodies are empty
// @router /check-transaksi/:noKontrak [get]
func (controller UsersController) CheckPengajuan() {
	appB := httpresponses.Bee{
		Ctx: controller.Ctx,
	}

	meta, errMeta := auth.ExtractedExt(controller.Ctx.Request, "")
	if errMeta != nil {
		appB.Response(http.StatusUnauthorized, "", errMeta.Error(), nil)
		return
	}

	result, err2 := repositories.StaticUserRepositoris().CheckPengajuan(meta.Id, controller.Ctx.Input.Param(":noKontrak"))
	if err2 != nil {
		appB.Response(http.StatusBadRequest, "", err2.Error(), nil)
		return
	}

	appB.Response(http.StatusOK, "Success", "", result)
	return
}

// CheckPayment
// @Description CheckPayment
// @Param	body	nil	true	"body nill"
// @Success 200 {int} interfaces{}
// @Failure 403 bodies are empty
// @router /check-payment/:noKontrak/:angsuranke [get]
func (controller UsersController) CheckPayment() {
	appB := httpresponses.Bee{
		Ctx: controller.Ctx,
	}

	meta, errMeta := auth.ExtractedExt(controller.Ctx.Request, "")
	if errMeta != nil {
		appB.Response(http.StatusUnauthorized, "", errMeta.Error(), nil)
		return
	}
	angsuranKe, _ := strconv.Atoi(controller.Ctx.Input.Param(":angsuranke"))

	result, err2 := repositories.StaticUserRepositoris().CheckPayment(meta.Id, angsuranKe, controller.Ctx.Input.Param(":noKontrak"))
	if err2 != nil {
		appB.Response(http.StatusBadRequest, "", err2.Error(), nil)
		return
	}

	appB.Response(http.StatusOK, "Success", "", result)
	return
}

// PaymentAngsuran
// @Description PaymentAngsuran
// @Param	body	nil	true	"body nill"
// @Success 200 {int} interfaces{}
// @Failure 403 bodies are empty
// @router /payment [post]
func (controller UsersController) PaymentAngsuran() {
	appB := httpresponses.Bee{
		Ctx: controller.Ctx,
	}

	meta, errMeta := auth.ExtractedExt(controller.Ctx.Request, "")
	if errMeta != nil {
		appB.Response(http.StatusUnauthorized, "", errMeta.Error(), nil)
		return
	}

	var input models.JSONTransaksiPayment
	if err := json.Unmarshal([]byte((controller.Ctx.Input.GetData("RequestBody").(string))), &input); err != nil {
		fmt.Println("err", err)
		appB.Response(http.StatusUnprocessableEntity, "", err.Error(), nil)
		return
	}

	result, err2 := repositories.StaticUserRepositoris().PaymentAngsuran(meta.Id, input)
	if err2 != nil {
		appB.Response(http.StatusBadRequest, "", err2.Error(), nil)
		return
	}

	appB.Response(http.StatusOK, "Success", "", result)
	return
}
