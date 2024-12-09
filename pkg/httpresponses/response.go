package httpresponses

import (
	"github.com/beego/beego/v2/server/web/context"
)

type Bee struct {
	Ctx *context.Context
}
type DataResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func (g *Bee) Response(httpCode int, message, errorMessage string, data interface{}) {
	g.Ctx.Output.SetStatus(httpCode)
	if errorMessage == "" { //Response sukses
		g.Ctx.Output.JSON(DataResponse{
			Message: message,
			Data:    data,
		}, true, true)
		return
	} else { //Response error
		g.Ctx.Output.JSON(DataResponse{
			Message: errorMessage,
			Data:    data,
		}, true, true)
		return
	}

}
