package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/bonczj/web-pub-sub/internal/pubsub"

	"github.com/bonczj/web-pub-sub/internal/handlers"
	"github.com/gorilla/mux"
)

var wait time.Duration

func init() {
	wait = 15 * time.Second
}

// main initializes the web server used by this project.
func main() {
	r := mux.NewRouter()

	r.HandleFunc("/publish", handlers.Publish).Methods(http.MethodPost)
	r.HandleFunc("/subscribe", handlers.Subscribe).Methods(http.MethodGet)

	// mux middleware is different. Build one to just log requets
	m := accessLog(r)

	srv := &http.Server{
		Handler: m,
		Addr:    ":8000",
	}
	// Run our server in a goroutine so that it doesn't block.
	go func() {
		log.Printf("Starting server on %s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()
	// shutdown any subscribers that are waiting for data
	instance := pubsub.Instance()
	instance.Clear()

	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	err := srv.Shutdown(ctx)
	if err != nil {
		log.Printf("error while shutting down HTTP server: %s", err)
	}

	log.Println("shutting down")
	os.Exit(0)
}

// accessLog prints each request to a log to simulate an access log.
func accessLog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s", r.Method, r.URL.Path)
		next.ServeHTTP(rw, r)
	})
}
