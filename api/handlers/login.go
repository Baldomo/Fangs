package handlers

import (
	"net/http"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) http.HandlerFunc {
	defer r.Body.Close()
}
