package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/salsa20"
)

type authResponse = struct {
	auth  bool
	token string
}

const authDelay = 10

func LoginHandler() http.HandlerFunc {
	isClientValid := false

	var key [32]byte
	copy(
		// Makes copy() think key is a slice
		key[:],

		// Make SECRENT_CLIENT_ID a slice of bytes and truncate it to 32
		[]byte(
			// Get first 32 characters of SECRET_CLIENT_ID
			os.Getenv("SECRET_CLIENT_ID")[:32],
		)[:32],
	)

	now := time.Now().Unix() / 1000

	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		buf := new(bytes.Buffer)
		buf.ReadFrom(r.Body)
		id := bytes.Split(buf.Bytes(), []byte("|"))
		iv, encrypted := id[0], id[1]

		// Decrypt message
		var messageBytes []byte
		salsa20.XORKeyStream(messageBytes, encrypted, iv, &key)
		messageString := string(messageBytes)

		for t := now; t >= now-authDelay && !isClientValid; t-- {
			isClientValid = messageString == fmt.Sprintf("%d|%s", t, os.Getenv("SECRET_CLIENT_ID"))
		}

		if !isClientValid {
			// If client is invalid, send 401 (unauthorized)
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			_ := json.NewEncoder(w).Encode(authResponse{false, ""})
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		// Create token
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"id":      "ApolloTV Official App",
			"message": "This better be from our app...",
			"ip":      r.RemoteAddr,

			// Expires in 1 hour
			"exp":     time.Now().Add(time.Hour),
		})

		// Try to sign token
		if tokenString, err := token.SignedString(os.Getenv("SECRET_CLIENT_ID")); err != nil {
			// If signing fails, return 500 (internal server error)
			w.WriteHeader(http.StatusInternalServerError)
			return
		} else {
			// Set header and send JSON response with token
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			if err = json.NewEncoder(w).Encode(authResponse{true, tokenString}); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusOK)
			return
		}
	}
}
