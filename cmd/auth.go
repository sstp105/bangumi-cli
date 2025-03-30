package cmd

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"

	"github.com/spf13/cobra"
	"github.com/sstp105/bangumi-cli/internal/bangumi"
	"github.com/sstp105/bangumi-cli/internal/config"
	"github.com/sstp105/bangumi-cli/internal/log"
	"github.com/sstp105/bangumi-cli/internal/ospath"
	"github.com/sstp105/bangumi-cli/internal/utils"
)

var authCmd = &cobra.Command{
	Use: "auth",
	Short: "Authenticate user and obtain banggumi API access token.",
	Run: func(cmd *cobra.Command, args []string) {
		authHandler()
	},
}

func init() {
	rootCmd.AddCommand(authCmd)
}

var creds *bangumi.OAuthCredential
var credsNew *bangumi.OAuthCredential
var ch = make(chan os.Signal, 1)

func authHandler() {
	if err := authenticate(); err != nil {
		log.Fatalf("error authenticating user:%s", err)
	}
}

func authenticate() error {
	var err error

	err = ospath.ReadJSONConfigFile(ospath.BangumiCredentialFile, &creds)

	switch {
	case err != nil && os.IsNotExist(err):
		err = oauth()
	case err != nil:
		return err
	case creds.IsExpired():
		err = oauth()
	case creds.ShouldRefresh():
		err = refreshCreds()
	}

	if err != nil {
		return err
	}

	if credsNew != nil {		
		if err := ospath.SaveJSONConfigFile(ospath.BangumiCredentialFile, credsNew); err != nil {
			log.Fatalf("error saving bangumi credentials:%s", err)
		}
		creds = credsNew
	}

	return printCreds(creds)
}

func printCreds(creds *bangumi.OAuthCredential) error {
	data, err := utils.MarshalJSONIndented(creds)
	if err != nil {
		return err
	}

	log.Infof("%s", string(data))
	
	return nil
}

func oauth() error {
	log.Info("Opening URL in the browser for bangumi.tv authentication...")

	clientID := config.BangumiClientID()
	url := fmt.Sprintf("https://bgm.tv/oauth/authorize?client_id=%s&response_type=code&redirect_uri=%s", clientID, config.LocalServerAddress())

	if err := utils.OpenBrowser(url); err != nil {
		return err
	}

	signal.Notify(ch, os.Interrupt)
	go func() {
		serve()
	}()
	<-ch

	return nil
}

func refreshCreds() error {
	var err error

	client := bangumi.NewOAuthClient()
	credsNew, err = client.RefreshAccessToken(
		config.BangumiClientID(),
		config.BangumiClientSecret(),
		config.LocalServerAddress(),
		creds.RefreshToken,
	)
	if err != nil {
		return err
	}

	return nil
}

func serve() {
	http.HandleFunc("/", listenCallback)

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

func listenCallback(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	code := q.Get("code")
	if code == "" {
		http.Error(w, "Missing code parameter in callback", http.StatusBadRequest)
		return
	}

	log.Debugf("Authorization code received from the callback: %s", code)

	var err error
	client := bangumi.NewOAuthClient()
	credsNew, err = client.GetAccessToken(
		config.BangumiClientID(),
		config.BangumiClientSecret(),
		config.LocalServerAddress(),
		code,
	)
	if err != nil {
		log.Fatalf("error requesting bangumi access_token:%s", err)
	}

	w.Write([]byte("You have successfully authenticated. You can close this window now."))

	ch <- os.Interrupt
}
