package routes

import "github.com/Baldomo/Claws/api/handlers"

var AuthRoutes = Routes{
	{
		"Login",
		"POST",
		"/api",
		handlers.LoginHandler,
	},
}