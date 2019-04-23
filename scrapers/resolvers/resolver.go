package resolvers

import (
	"github.com/Baldomo/Fangs/utils"
	"github.com/Baldomo/Fangs/ws"
	"github.com/Baldomo/Fangs/ws/event"
	"github.com/valyala/fasthttp"
	"net/http"
	"os"
	"strings"
)

//type Resolver func(uri string, cookies []*http.Cookie, header *http.Header) error

type Resolver interface {
	Resolve(uri string, provider string, header *http.Header, cookie []*http.Cookie) error
	SupportsURI(uri string) bool
	IsLocked() bool
}

func ResolveHtml(html string, resolver string, header *http.Header, cookie []*http.Cookie) {}

func Resolve(uri string, provider string, quality string, ctx *fasthttp.RequestCtx, cookie []*http.Cookie) {
	ipLocked := os.ExpandEnv("CLAWS_ENV") == "server"

	if strings.Contains(uri, "openload.co") || strings.Contains(uri, "oload.coud") {
		videoId := strings.Split(uri, "/")[4]
		if strings.Contains(uri, "embed") {
			uri = "https://openload.co/embed/" + videoId
		}

		if !ipLocked {
			// Call Openload provider.Provide() to get data
			data := ""
			result := event.NewResult(
				data,
				nil,
				utils.FastHTTPHeaderToMap(ctx),
				false,
				event.NewMetadata(quality, provider, "Openload", false, cookie),
			)
			ws.SendEvent(result, ctx)
			return
		}

		scrape := event.NewScrape(
			"",
			utils.FastHTTPHeaderToMap(ctx),
			event.NewPairing("https://olpair.com", videoId, uri),
			event.NewMetadata(quality, provider, "Openload", false, cookie),
		)
		ws.SendEvent(scrape, ctx)
		return
	} else if strings.Contains(uri, "streamango.com") {

	}
}
