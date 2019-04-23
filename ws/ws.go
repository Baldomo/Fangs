package ws

import (
	"github.com/Baldomo/Fangs/logger"
	"github.com/Baldomo/Fangs/ws/event"
	"github.com/fasthttp/websocket"
	"github.com/valyala/fasthttp"
	"time"
)

var upgrader = websocket.FastHTTPUpgrader{
	// Keep a 5 second timeout for the upgrade to complete
	HandshakeTimeout: 5 * time.Second,
}

func SendEvent(event event.Event, ctx *fasthttp.RequestCtx) {
	err := upgrader.Upgrade(ctx, func(conn *websocket.Conn) {
		conn.WriteJSON(event)
	})
	if err != nil {
		logger.Error("Failed to send websocket event", "event_type", event.EventType())
	}
}
