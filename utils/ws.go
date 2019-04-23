package utils

import (
	"github.com/valyala/fasthttp"
	"net/http"
)

// Converts a fasthttp.RequestHeader to http.Header
func FastHTTPHeaderToMap(ctx *fasthttp.RequestCtx) (header *http.Header) {
	ctx.Request.Header.VisitAll(func(key, value []byte) {
		header.Set(string(key), string(value))
	})

	return
}

func MapToFastHTTPHeader(header *http.Header, ctx *fasthttp.RequestCtx) {
	for k, v := range *header {
		ctx.Request.Header.Set(k, v[0])
	}

	return
}
