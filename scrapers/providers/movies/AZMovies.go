package movies

import (
	"fmt"
	"github.com/Baldomo/Fangs/logger"
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/extensions"
	"net/http"
	"regexp"
)

type azmovies struct {
	cookies  []*http.Cookie
	movieUri string
	urls     []string
}

var AZMovies = azmovies{
	cookies:  []*http.Cookie{},
	movieUri: "%s/watch.php?title=%s",
	urls: []string{
		"https://azmovie.to",
	},
}

// req.client.remoteAddress.replace('::ffff:', '').replace('::1', '');

func (s azmovies) Scrape(title string, c chan<- interface{}) {
	coll := colly.NewCollector(
		colly.Async(true),
	)
	coll.Limit(&colly.LimitRule{
		Parallelism: 2,
	})
	extensions.RandomUserAgent(coll)

	coll.OnHTML("#serverul li a", func(e *colly.HTMLElement) {
		logger.Debug("Got link", "movie", e.Attr("href"))
	})

	coll.OnHTML("script", s.scrapeCookies)

	for _, u := range s.urls {
		coll.SetCookies(u, s.cookies)
		coll.Visit(fmt.Sprintf(s.movieUri, u, title))
	}

	coll.Wait()
}

func (s *azmovies) scrapeCookies(e *colly.HTMLElement) {
	html := string(e.Text)
	re := regexp.MustCompile(`document\.cookie\s*=\s*"(.*)=(.*)";`)
	matches := re.FindAllStringSubmatch(html, -1)

	for _, match := range matches {
		s.cookies = append(s.cookies, &http.Cookie{
			Name:     match[0],
			Value:    match[1],
			HttpOnly: false,
		})
	}
}

func (s azmovies) String() string {
	return "AZMovies"
}
