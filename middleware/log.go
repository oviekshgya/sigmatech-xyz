package middleware

import (
	"bytes"
	"fmt"
	"github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/context"
	"io/ioutil"
	"log"
	"net/http"
	"sigmatech-xyz/pkg/auth"
	"sigmatech-xyz/repositories"
)

func Log(next web.FilterFunc) web.FilterFunc {
	return func(c *context.Context) {
		//appG := httpresponses.Gin{C: c}
		metadata, errMetaData := auth.ExtractedExt(c.Request, "")
		if errMetaData != nil {
			log.Println("tidak login", errMetaData.Error())
			//return
		}
		if metadata != nil {
			_, err := repositories.StaticLogRepositories().LogCreated(metadata.Id, BeegoBodyLogMiddleware(c), BeegoBodyLogStatusCode(c), c)
			if err != nil {
				fmt.Println(err.Error)
				//appG.Response(err.Code, err.Message, nil)
				//c.Abort()
				//return
			}

		} else {
			_, err := repositories.StaticLogRepositories().LogCreated(0, BeegoBodyLogMiddleware(c), BeegoBodyLogStatusCode(c), c)
			if err != nil {
				fmt.Println(err.Error)
				//appG.Response(err.Code, err.Message, nil)
				//c.Abort()
				//return
			}
		}
		next(c)

	}
}

type bodyLogWriter struct {
	http.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)

	return w.ResponseWriter.Write(b)
}

func BeegoBodyLogMiddleware(ctx *context.Context) map[string]string {
	blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: ctx.ResponseWriter}
	//ctx.ResponseWriter = blw.ResponseWriter
	body, _ := ioutil.ReadAll(ctx.Request.Body)
	ctx.Request.Body = ioutil.NopCloser(bytes.NewBuffer([]byte(body)))
	fmt.Println("string(body) log middleware", string(body))
	ctx.Input.SetData("RequestBody", string(body))
	fmt.Println("ctx.ResponseWriter", ctx.ResponseWriter)
	fmt.Println("blw.body.String()", blw.body)
	return map[string]string{
		"bodyform":     string(body),
		"bodyresponse": blw.body.String(),
	}
}

func BeegoBodyLogStatusCode(ctx *context.Context) int {
	statusCode := ctx.ResponseWriter.Status
	return statusCode
}
