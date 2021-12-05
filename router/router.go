package router

import (
	"context"
	"fmt"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/fasthttp/router"
	"github.com/projectlukman/lib-go/log"
	"github.com/projectlukman/lib-go/router/response"
	"github.com/valyala/fasthttp"
)

type MyRouter struct {
	Httprouter     *router.Router
	WrappedHandler http.Handler
	Options        *Options
}

type Options struct {
	Prefix string
}

func NewRouter(opt *Options) *MyRouter {
	return &MyRouter{
		Httprouter: router.New(),
		Options:    opt,
	}
}

func (mr *MyRouter) GET(path string, reqHandler Handle) {
	mr.Handle(path, http.MethodGet, reqHandler)
}

func (mr *MyRouter) POST(path string, reqHandler Handle) {
	mr.Handle(path, http.MethodPost, reqHandler)
}

func (mr *MyRouter) PUT(path string, reqHandler Handle) {
	mr.Handle(path, http.MethodPut, reqHandler)
}

func (mr *MyRouter) PATCH(path string, reqHandler Handle) {
	mr.Handle(path, http.MethodPatch, reqHandler)
}

func (mr *MyRouter) DELETE(path string, reqHandler Handle) {
	mr.Handle(path, http.MethodDelete, reqHandler)
}

func (mr *MyRouter) Handle(path, method string, reqHandler Handle) {
	fullPath := mr.Options.Prefix + path
	log.Println(fmt.Sprintf("%s: %s", method, fullPath))
	mr.Httprouter.Handle(method, fullPath, mr.handleNow(fullPath, reqHandler))
}

// func (mr *MyRouter) ServeFiles(path string, root http.FileSystem) {
// 	fullPath := mr.Options.Prefix + path
// 	mr.Httprouter.ServeFiles(fullPath, root)
// }

func (mr *MyRouter) Group(path string, fn func(r *MyRouter)) {
	sr := NewRouter(&Options{
		Prefix: mr.Options.Prefix + path,
	})
	fn(sr)
}

type panicObject struct {
	err        interface{}
	stackTrace string
}

func (mr *MyRouter) handleNow(fullPath string, handle Handle) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		t := time.Now()
		ctx.Response.Header.Set("routePath", fullPath)
		ctx.Response.Header.Set("Content-Type", "application/json")

		respChan := make(chan *response.JSONResponse)
		recovered := make(chan panicObject)

		go func() {
			defer func() {
				if err := recover(); err != nil {
					recovered <- panicObject{
						err:        err,
						stackTrace: string(debug.Stack()),
					}
				}
			}()
			respChan <- handle(ctx)
		}()

		select {
		case <-ctx.Done():
			if ctx.Err() == context.DeadlineExceeded {
				resp := response.NewJSONResponse().SetStatusCode(http.StatusRequestTimeout)
				ctx.Response.SetStatusCode(http.StatusRequestTimeout)
				ctx.Response.SetBody(resp.GetBody())
			}
		case cause := <-recovered:
			log.WithFields(log.Fields{
				"path":       string(ctx.Request.RequestURI()),
				"stackTrace": cause.stackTrace,
				"error":      fmt.Sprintf("%v", cause.err),
			}).Errorln("[Router] panic have occurred")
			resp := response.NewJSONResponse().SetStatusCode(http.StatusInternalServerError)
			ctx.Response.SetStatusCode(http.StatusInternalServerError)
			ctx.Response.SetBody(resp.GetBody())
		case resp := <-respChan:
			if resp != nil {
				resp.SetLatency(time.Since(t).Seconds() * 1000)
				ctx.Response.SetStatusCode(resp.StatusCode)
				ctx.Response.SetBody(resp.GetBody())
			} else {
				if fullPath != "/metrics" {
					resp := response.NewJSONResponse().SetStatusCode(http.StatusInternalServerError)
					ctx.Response.SetStatusCode(http.StatusInternalServerError)
					ctx.Response.SetBody(resp.GetBody())
				}
			}
		}
	}
}
