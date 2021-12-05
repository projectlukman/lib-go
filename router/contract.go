package router

import (
	"github.com/projectlukman/lib-go/router/response"
	"github.com/valyala/fasthttp"
)

type Handle func(*fasthttp.RequestCtx) *response.JSONResponse

type Router interface {
	GET(path string, reqHandler Handle)
	POST(path string, reqHandler Handle)
	PUT(path string, reqHandler Handle)
	PATCH(path string, reqHandler Handle)
	DELETE(path string, reqHandler Handle)

	Handle(path, method string, reqHandler Handle)
	Group(path string, fn func(r *MyRouter))
}
