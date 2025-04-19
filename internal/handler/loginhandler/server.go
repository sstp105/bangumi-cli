package loginhandler

import (
	"fmt"
	"github.com/sstp105/bangumi-cli/internal/bangumi"
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

// oAuthCallback allows the caller to pass a variable and update it with the received bangumi credential.
type oAuthCallback func(*bangumi.OAuthCredential)

// oAuthHandler handles the callback request from bangumi.
// Reference: https://github.com/bangumi/api/blob/master/docs-raw/How-to-Auth.md
func oAuthHandler(w http.ResponseWriter, r *http.Request, callback oAuthCallback) {
	q := r.URL.Query()
	code := q.Get("code")
	if code == "" {
		http.Error(w, "Missing code parameter in callback", http.StatusBadRequest)
		return
	}

	client := bangumi.NewOAuthClient()
	credential, err := client.GetAccessToken(
		config.BangumiClientID(),
		config.BangumiClientSecret(),
		config.LocalServerAddress(),
		code,
	)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to authenticate bangumi:%s", err), http.StatusBadRequest)
		return
	}

	callback(credential)

	w.Write([]byte("You have successfully authenticated. You may close this window now."))
}
