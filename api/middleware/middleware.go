package middleware

import (
	"encoding/binary"
	"github.com/Baldomo/Fangs/logger"
	"github.com/valyala/fasthttp"
)

type Handler func(fasthttp.RequestHandler) fasthttp.RequestHandler

var All = []Handler{
	AddHeaders,
	CheckSize,
	Compress,
	Logger,
}

func Wrap(h fasthttp.RequestHandler) fasthttp.RequestHandler {
	for _, m := range All {
		h = m(h)
	}

	return h
}

func AddHeaders(h fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		h(ctx)
		ctx.Response.Header.Set("Access-Control-Allow-Origin", "*")
		ctx.Response.Header.Set("Access-Control-Allow-Headers", "origin, X-Requested-With, Content-Type, Accept")
	}
}

func CheckSize(h fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		// Check if body is larger than 10mb (10'000'000 bytes)
		if binary.Size(ctx.Request.Body()) > 10000000 {
			// Return 413 (entitity too large)
			ctx.SetStatusCode(fasthttp.StatusRequestEntityTooLarge)
			return
		}

		h(ctx)
	}
}

func Compress(h fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		if ctx.Request.Header.HasAcceptEncoding("text/event-stream") || len(ctx.Request.Header.Peek("x-no-compression")) > 0 {
			h(ctx)
		} else {
			// Note: only compresses if 'Accept-Encoding' is 'gzip' or 'deflate'
			fasthttp.CompressHandler(h)
		}
	}
}

func Logger(h fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		logger.Debug("Got request",
			"method", string(ctx.Method()),
			"url", string(ctx.Path()),
			"raw_body", string(ctx.Request.Body()),
		)
		h(ctx)
	}
}
