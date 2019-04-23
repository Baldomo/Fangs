package movies

import (
	"fmt"
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/extensions"
	"net/http"
	"strings"
)

var spaceReplacements = []string{
	`%20`,
	`+`,
	`%2B`,
}

type afdah struct {
	cookies  []*http.Cookie
	movieUri string
	urls     []string
}

var Afdah = afdah{
	cookies:  []*http.Cookie{},
	movieUri: "%s/?s=%s",
	urls: []string{
		"https://afdah.org",
		"https://genvideos.com",
		"https://genvideos.co",
		"https://watch32hd.co",
		"https://putlockerhd.co",
	},
}

// title, year, client ip
func (scraper afdah) Scrape(title string, c chan<- interface{}) {
	searchColl := colly.NewCollector(
		colly.Async(true),
	)
	searchColl.Limit(&colly.LimitRule{
		Parallelism: 2,
	})
	extensions.RandomUserAgent(searchColl)

	searchColl.OnError(func(resp *colly.Response, e error) {
		if len(spaceReplacements) == 0 {
			c <- e
			return
		}
		// Do things with resp
		formattedUrl := fmt.Sprintf(
			scraper.movieUri,
			resp.Request.URL.Host,
			formatTitle(title, spaceReplacements[0]),
		)
		searchColl.Visit(formattedUrl)
	})

	searchColl.OnHTML(".cell", func(e *colly.HTMLElement) {
		videoName := strings.TrimSpace(strings.ToLower(e.ChildText(".video_title")))
		videoYearAndQuality := strings.TrimSpace(e.ChildText(".video_title"))

	})

	for _, u := range scraper.urls {
		scraper.cookies = searchColl.Cookies(u)
	}
}

func formatTitle(title string, space string) string {
	spaceReplaced := strings.ReplaceAll(
		strings.ToLower(title),
		" ",
		space,
	)

	return strings.ReplaceAll(spaceReplaced, ":", "")
}

func (scraper afdah) String() string {
	return "Afdah"
}
