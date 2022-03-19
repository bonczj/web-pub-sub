package handlers

import (
	"io/ioutil"
	"log"
	"net/http"

	"github.com/bonczj/web-pub-sub/internal/pubsub"
)

// Publish receives a new message in the body of the request. If there are
// any active subscribers, the message is published to each subscriber.
func Publish(rw http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("error reading body of message: %s", err)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	instance := pubsub.Instance()
	instance.Publish(body)

	rw.WriteHeader(http.StatusNoContent)
}
