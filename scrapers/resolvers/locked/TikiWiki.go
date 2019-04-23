package locked

import (
	"github.com/gocolly/colly"
	"net/http"
	"strings"
)

type tikiwiki struct {
	cookies  []*http.Cookie
}

var TikiWiki = tikiwiki{
	cookies: []*http.Cookie{},
}

func (resolver tikiwiki) Resolve(uri string, provider string, header *http.Header, cookie []*http.Cookie) error {
	coll := colly.NewCollector(
		colly.Async(true),
	)
	coll.Limit(&colly.LimitRule{
		Parallelism: 2,
	})

	return nil
}

func (resolver tikiwiki) SupportsURI(uri string) bool {
	return strings.Contains(uri, "tiki.wiki/embed")
}

func (resolver tikiwiki) IsLocked() bool {
	return false
}
