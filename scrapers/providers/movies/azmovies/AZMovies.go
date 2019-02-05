package azmovies

import "github.com/gocolly/colly"

var urls = []string{
	"https://azmovies.xyz",
}

var movieUri = "%s/watch.php?title=%s"

func scrape() {
	c := colly.NewCollector()
	c.OnHTML("#serverul li a", func(e *colly.HTMLElement) {

	})
}
