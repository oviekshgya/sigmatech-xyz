package controllers

import (
	"encoding/json"
	"fmt"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/go-playground/validator/v10"
	"net/http"
	"sigmatech-xyz/models"
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
