package handlers

import "github.com/valyala/fasthttp"

func StatusHandler(ctx *fasthttp.RequestCtx) {
	ctx.SetStatusCode(fasthttp.StatusOK)
}
