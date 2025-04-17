package loghandler

import (
	"net/http"

	"github.com/sstp105/bangumi-cli/internal/config"
	"github.com/sstp105/bangumi-cli/internal/log"
)

// Start starts a local HTTP server to handle requests (e.g. bangumi login)
func Start(callback oAuthCallback) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		oAuthHandler(w, r, callback)
	})

	port := config.Port()
	addr := ":" + port
	log.Debugf("Listening for callback at localhost%s", addr)

	server := &http.Server{
		Addr: addr,
	}
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("error starting the HTTP server locally: %s", err)
	}
}
