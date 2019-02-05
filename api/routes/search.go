package routes

import "github.com/Baldomo/Fangs/api/handlers"

var SearchRoutes = Routes{
	{
		"Movies",
		"GET",
		"/api/v1/movies",
		handlers.SearchHandler(),
	},
}
