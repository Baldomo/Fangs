package main

import (
	"github.com/Baldomo/Fangs/api/middleware"
	"github.com/Baldomo/Fangs/api/routes"
	"github.com/Baldomo/Fangs/logger"
	"github.com/Baldomo/Fangs/utils"
	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
	"net/http"

	_ "net/http/pprof"
)

func main() {
	router := fasthttprouter.New()

	for _, route := range routes.GET {
		router.GET(
			route.Pattern,
			middleware.Wrap(route.HandlerFunc),
		)
	}

	for _, route := range routes.POST {
		router.POST(
			route.Pattern,
			middleware.Wrap(route.HandlerFunc),
		)
	}

	if utils.IsDebug() {
		go func() {
			if err := http.ListenAndServe("0.0.0.0:6060", nil); err != nil {
				logger.DPanic("pprof failed: %v", err)
			}
		}()
	}

	logger.Debug("server starting")
	err := fasthttp.ListenAndServe("0.0.0.0:8080", router.Handler)
	if err != nil {
		logger.Fatal("server stopped!", "error", err)
	}
}
