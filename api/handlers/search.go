package handlers

import (
	"github.com/valyala/fasthttp"
)

func SearchHandler(ctx *fasthttp.RequestCtx) {
	ctx.SetStatusCode(fasthttp.StatusOK)
}
