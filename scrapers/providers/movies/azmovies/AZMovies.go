package azmovies

import "github.com/gocolly/colly"

var urls = []string{
	"https://azmovies.xyz",
}

func scrape() {
	c := colly.NewCollector()
	c.OnHTML("#serverul li a", func(e *colly.HTMLElement) {

	})
}
