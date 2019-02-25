package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/Baldomo/Fangs/logger"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/valyala/fasthttp"
	"golang.org/x/crypto/salsa20"
)

type authResponse = struct {
	auth  bool
	token string
}

const authDelay = 10

func LoginHandler(ctx *fasthttp.RequestCtx) {
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

	id := bytes.Split(ctx.Request.Body(), []byte("|"))
	iv, encrypted := id[0], id[1]

	// Decrypt message
	var messageBytes []byte
	salsa20.XORKeyStream(messageBytes, encrypted, iv, &key)
	messageString := string(messageBytes)

	for t := now; t >= now-authDelay && !isClientValid; t-- {
		isClientValid = messageString == fmt.Sprintf("%d|%s", t, os.Getenv("SECRET_CLIENT_ID"))
	}

	if !isClientValid {
		logger.Debug("Invalid client", "addr", ctx.RemoteAddr())
		// If client is invalid, send 401 (unauthorized)
		ctx.SetContentType("application/json; charset=UTF-8")
		if mess, err := json.Marshal(authResponse{false, ""}); err != nil {
			// If marshaling fails, return 500 (internal server error)
			ctx.SetStatusCode(fasthttp.StatusInternalServerError)
			return
		} else {
			ctx.SetStatusCode(http.StatusUnauthorized)
			ctx.SetBody(mess)
			return
		}
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":      "ApolloTV Official App",
		"message": "This better be from our app...",
		"ip":      ctx.RemoteAddr().String(),

		// Expires in 1 hour
		"exp": time.Now().Add(time.Hour),
	})

	// Try to sign token
	if tokenString, err := token.SignedString(os.Getenv("SECRET_CLIENT_ID")); err != nil {
		// If signing fails, return 500 (internal server error)
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return
	} else {
		// Set header and send JSON response with token
		ctx.SetContentType("application/json; charset=UTF-8")
		if mess, err := json.Marshal(authResponse{true, tokenString}); err != nil {
			// If marshaling fails, return 500 (internal server error)
			ctx.SetStatusCode(fasthttp.StatusInternalServerError)
			return
		} else {
			ctx.SetStatusCode(http.StatusOK)
			ctx.SetBody(mess)
			return
		}
	}
}
