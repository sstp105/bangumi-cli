package server

import (
	"net/http"

	"github.com/sstp105/bangumi-cli/internal/bangumi"
	"github.com/sstp105/bangumi-cli/internal/config"
	"github.com/sstp105/bangumi-cli/internal/log"
)

type oAuthCallback func(*bangumi.OAuthCredential)

func oAuthHandler(w http.ResponseWriter, r *http.Request, callback oAuthCallback) {
	q := r.URL.Query()
	code := q.Get("code")
	if code == "" {
		http.Error(w, "Missing code parameter in callback", http.StatusBadRequest)
		return
	}

	client := bangumi.NewOAuthClient()
	creds, err := client.GetAccessToken(
		config.BangumiClientID(),
		config.BangumiClientSecret(),
		config.LocalServerAddress(),
		code,
	)
	if err != nil {
		log.Fatalf("error requesting bangumi access_token:%s", err)
	}

	callback(creds)

	w.Write([]byte("You have successfully authenticated. You can close this window now."))
}
