package providers

import "github.com/Baldomo/Fangs/scrapers/providers/movies"

var c = make(chan interface{})

var Movies = []Provider{
	movies.AZMovies,
	movies.Afdah,
}

type Provider interface {
	Scrape(title string, c chan<- interface{})
}

func Provide(title string) {
	for _, movieProvider := range Movies {
		movieProvider.Scrape(title, c)
	}
}
