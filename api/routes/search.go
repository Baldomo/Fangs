package routes

import "github.com/Baldomo/Claws/api/handlers"

var SearchRoutes = Routes{
	{
		"Movies",
		"GET",
		"/api/v1/movies",
		handlers.SearchHandler,
	},
}
