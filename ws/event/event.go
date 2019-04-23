package event

import (
	"github.com/Baldomo/Fangs/utils"
	"github.com/satori/go.uuid"
	"net/http"
	"time"
)

type Event interface {
	EventType() string
}

type Scrape struct {
	Event string

	CookieRequired string
	Headers        *http.Header
	Resolver       string
	ScrapeId       string
	Source         string
	Target         string
}

func (s Scrape) EventType() string {
	return s.Event
}

type Status struct {
	Event string

	data []int64
}

type Result struct {
	Event            string
	Error            string `json:",omitempty"`
	Headers          *http.Header
	IsResultOfScrape bool
	Metadata

	file
}

func (r Result) EventType() string {
	return r.Event
}

type file struct {
	Data string
	Kind string
}

type Metadata struct {
	Quality    string
	Provider   string
	Source     string
	IsDownload bool
	Cookie     []*http.Cookie
}

type pairing struct {
	URL, VideoID, Target string
}

func NewScrape(cookieRequired string, header *http.Header, pairing *pairing, metadata Metadata) *Scrape {
	return &Scrape{
		Event:          "scrape",
		CookieRequired: cookieRequired,
		Headers:        header,
		Resolver:       metadata.Provider,
		ScrapeId:       uuid.NewV4().String(),
		Source:         metadata.Source,
		Target:         pairing.Target,
	}
}

func NewStatus() *Status {
	return &Status{
		Event: "status",
		data: []int64{
			time.Now().Unix(),
		},
	}
}

func NewResult(data string, error error, header *http.Header, isResultOfScrape bool, metadata Metadata) *Result {
	if error != nil {
		return &Result{
			Event:            "result",
			Error:            error.Error(),
			Headers:          header,
			IsResultOfScrape: isResultOfScrape,
			Metadata:         metadata,
			file: file{
				Data: data,
				Kind: utils.GetDataKind(data),
			},
		}
	}

	return &Result{
		Event:            "result",
		Headers:          header,
		IsResultOfScrape: isResultOfScrape,
		Metadata:         metadata,
		file: file{
			Data: data,
			Kind: utils.GetDataKind(data),
		},
	}
}

func NewMetadata(quality, provider, source string, isDownload bool, cookie []*http.Cookie) Metadata {
	return Metadata{quality, provider, source, isDownload, cookie}
}

func NewPairing(url, videoId, target string) *pairing {
	return &pairing{url, videoId, target}
}
