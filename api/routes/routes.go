package routes

import (
	"github.com/Baldomo/Fangs/api/handlers"
	"github.com/valyala/fasthttp"
)

type Route struct {
	Name          string
	Pattern       string
	HandlerFunc   fasthttp.RequestHandler
}

var GET = []Route{
	{
		"Movies",
		"/api/v1/movies",
		handlers.SearchHandler,
	}, {
		"Status",
		"/api/v1/status",
		handlers.StatusHandler,
	},
}

var POST = []Route{
	{
		"Login",
		"/api/v1/login",
		handlers.LoginHandler,
	},
}
