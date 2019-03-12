package main

import (
	"github.com/Baldomo/Fangs/api/routes"
	"net/http"
	"strings"
	"testing"
)

var baseAddr = "http://localhost:8080"

func BenchmarkGet(b *testing.B) {
	go main()

	for _, endpoint := range routes.GET {
		for ind := 0; ind <= b.N; ind++ {
			if resp, err := http.Get(baseAddr + endpoint.Pattern); err != nil {
				b.Fatal(err)
			} else if resp.StatusCode != 200 {
				b.Errorf("got code %d on run %d", resp.StatusCode, ind)
			}
		}
	}
}

func BenchmarkPost(b *testing.B) {
	go main()

	for _, endpoint := range routes.POST {
		for ind := 0; ind <= b.N; ind++ {
			if resp, err := http.Post(
				baseAddr + endpoint.Pattern,
				"application/json",
				strings.NewReader(""),
			); err != nil {
				b.Fatal(err)
			} else if resp.StatusCode != 200 {
				b.Errorf("got code %d on run %d", resp.StatusCode, ind)
			}
		}
	}
}
