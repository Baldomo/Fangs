package handlers

import (
	"github.com/Baldomo/Fangs/api/sse"
	"net/http"
)

func SearchHandler() http.HandlerFunc {
	nc := sse.NewNotificationCenter()
	c := make(chan []byte)

	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")
		w.Header().Set("Access-Control-Allow-Origin", "*")

		unsubscribeFn, err := nc.Subscribe(c)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		select {
		case <-r.Context().Done():
			if err := nc.Notify([]byte("disconnected")); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			if err := unsubscribeFn(); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		default:

			if err := nc.Notify([]byte("done")); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
	}
}