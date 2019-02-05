package routes

import "github.com/Baldomo/Fangs/api/handlers"

var AuthRoutes = Routes{
	{
		"Login",
		"POST",
		"/api",
		handlers.LoginHandler(),
	},
}