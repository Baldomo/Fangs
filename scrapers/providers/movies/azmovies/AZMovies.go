package azmovies

import (
	"fmt"
	"github.com/Baldomo/Fangs/logger"
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/extensions"
	"net/http"
	"regexp"
)

var (
	urls = []string{
		"https://azmovies.xyz",
	}

	movieUri = "%s/watch.php?title=%s"

	cookies []*http.Cookie
)

func Scrape(movieTitle string) {
	collector := colly.NewCollector(
		colly.Async(true),
	)

	collector.Limit(&colly.LimitRule{
		Parallelism: 2,
	})

	extensions.RandomUserAgent(collector)

	collector.OnHTML("#serverul li a", func(e *colly.HTMLElement) {
		logger.Debug("Got link", "movie", e.Attr("href"))
		// TODO resolve
	})

	collector.OnHTML("script", scrapeCookies)

	for _, u := range urls {
		collector.SetCookies(u, cookies)
		collector.Visit(fmt.Sprintf(movieUri, u, movieTitle))
	}

	collector.Wait()
}

func scrapeCookies(e *colly.HTMLElement) {
	html := string(e.Text)
	re := regexp.MustCompile(`document\.cookie\s*=\s*"(.*)=(.*)";`)
	matches := re.FindAllStringSubmatch(html, -1)

	for _, match := range matches {
		cookies = append(cookies, &http.Cookie{
			Name:     match[0],
			Value:    match[1],
			HttpOnly: false,
		})
	}
}
