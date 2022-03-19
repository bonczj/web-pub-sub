package handlers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/bonczj/web-pub-sub/internal/pubsub"
	"github.com/gorilla/websocket"

	"github.com/google/uuid"
)

const (
	ContentTypeHeader    = "Content-Type"
	ContentTypeTextPlain = "text/plain"
)

// Subscribe registers a client as wanting to subscribe to new messages.
// Any new published message will be sent to the client.
//
// Subscribers can request converting an HTTP connection to a web socket
// connection. If the upgrade is requested and works, stream messages
// that way, otherwise, write new lines out over the HTTP connection.
func Subscribe(w http.ResponseWriter, r *http.Request) {
	id := uuid.New().String()
	ch := make(chan []byte, 10)
	instance := pubsub.Instance()

	err := instance.Subscribe(id, ch)
	defer instance.Unsubscribe(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if websocket.IsWebSocketUpgrade(r) {
		u := websocket.Upgrader{}
		c, err := u.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		defer c.Close()

		relayWebSocketMessages(c, ch)
	} else {
		relayHttpMessages(w, ch)
	}
}

func relayWebSocketMessages(c *websocket.Conn, ch chan []byte) {
	for {
		msg, ok := <-ch
		if !ok {
			// channel was closed, exit out
			return
		}

		if err := c.WriteMessage(websocket.TextMessage, msg); err != nil {
			log.Printf("error writing content over web socket: %s", err)
			return
		}
	}
}

func relayHttpMessages(w http.ResponseWriter, ch chan []byte) {
	w.Header().Set(ContentTypeHeader, ContentTypeTextPlain)
	w.WriteHeader(http.StatusAccepted)

	for {
		msg, ok := <-ch
		if !ok {
			// channel was closed, exit out
			return
		}

		if _, err := fmt.Fprintln(w, string(msg)); err != nil {
			log.Printf("error writing content out: %s", err)
			return
		}

		if f, ok := w.(http.Flusher); ok {
			f.Flush()
		}
	}
}
