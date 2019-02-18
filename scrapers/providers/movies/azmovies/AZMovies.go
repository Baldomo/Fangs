package azmovies

import (
	"fmt"
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/extensions"
	"net/http/cookiejar"
	"net/url"
)

var urls = []string{
	"https://azmovies.xyz",
}

var movieUri = "%s/watch.php?title=%s"

func Scrape(movieTitle string) {
	c := colly.NewCollector(
		colly.AllowedDomains(urls...),
		colly.Async(true),
	)
	c.Limit(&colly.LimitRule{
		Parallelism: 2,
	})
	extensions.RandomUserAgent(c)
	c.OnHTML("#serverul li a", func(e *colly.HTMLElement) {
		// TODO: collect movie
	})
	c.OnResponse(func(r *colly.Response) {
		// r.Body
		//html := string(r.Body)
		//re := regexp.MustCompile(`document\.cookie\s*=\s*"(.*)=(.*)";`)
		//matches := re.FindAllStringSubmatch(html, -1)
	})
	for _, u := range urls {
		cookies := c.Cookies(u)
		jar, _ := cookiejar.New(nil)
		parsedUrl, _ := url.Parse(u)
		jar.SetCookies(parsedUrl, cookies)
		c.SetCookieJar(jar)
		// TODO: add cookies from parsed page
		//c.SetCookies(u, []*http.Cookie{
		//	{
		//
		//	},
		//})
		_ = c.Visit(fmt.Sprintf(movieUri, u, movieTitle))
	}
	c.Wait()
}
