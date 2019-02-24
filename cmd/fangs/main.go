package main

import (
	"github.com/Baldomo/Fangs/api/middleware"
	"github.com/Baldomo/Fangs/api/routes"
	"github.com/Baldomo/Fangs/logger"
	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
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

	err := fasthttp.ListenAndServe(":8080", router.Handler)
	logger.Fatal("server stopped!", "error", err)
}
