package auth

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/sstp105/bangumi-cli/internal/bangumi"
	"github.com/sstp105/bangumi-cli/internal/config"
	"github.com/sstp105/bangumi-cli/internal/log"
	"github.com/sstp105/bangumi-cli/internal/ospath"
	"github.com/sstp105/bangumi-cli/internal/server"
	"github.com/sstp105/bangumi-cli/internal/utils"
)

var ch = make(chan os.Signal, 1)

func Handler() {
	if err := authenticate(); err != nil {
		log.Fatalf("error authenticating user:%s", err)
	}
}

func authenticate() error {
	var creds *bangumi.OAuthCredential
	var err error
	overwrite := true

	err = ospath.ReadJSONConfigFile(ospath.BangumiCredentialFile, &creds)

	switch {
	case err != nil && os.IsNotExist(err):
		log.Debugf("Bangumi credential config is not found, authenticating user.")
		creds, err = oauth()
	case err != nil:
		log.Errorf("error reading bangumi credential config:%s", err)
		return err
	case creds.IsExpired():
		log.Debugf("Access token already expired, re-authenticating user.")
		creds, err = oauth()
	case creds.ShouldRefresh():
		log.Debugf("Access token will expire soon, fetching new access token.")
		creds, err = refresh(creds.RefreshToken)
	default:
		log.Debugf("Bangumi credential config is valid.")
		overwrite = false
	}

	if err != nil {
		return err
	}

	if creds != nil && overwrite {
		if err := ospath.SaveJSONConfigFile(ospath.BangumiCredentialFile, creds); err != nil {
			log.Fatalf("error saving bangumi credentials:%s", err)
		}
	}

	return creds.Print()
}

func oauth() (*bangumi.OAuthCredential, error) {
	log.Info("Opening URL in the browser for bangumi.tv authentication...")

	clientID := config.BangumiClientID()
	url := fmt.Sprintf("https://bgm.tv/oauth/authorize?client_id=%s&response_type=code&redirect_uri=%s", clientID, config.LocalServerAddress())

	if err := utils.OpenBrowser(url); err != nil {
		return nil, err
	}

	signal.Notify(ch, os.Interrupt)

	var creds *bangumi.OAuthCredential
	go func() {
		server.Start(func(c *bangumi.OAuthCredential) {
			creds = c
			ch <- os.Interrupt
		})
	}()
	<-ch

	return creds, nil
}

func refresh(token string) (*bangumi.OAuthCredential, error) {
	client := bangumi.NewOAuthClient()
	creds, err := client.RefreshAccessToken(
		config.BangumiClientID(),
		config.BangumiClientSecret(),
		config.LocalServerAddress(),
		token,
	)
	if err != nil {
		return nil, err
	}

	return creds, nil
}
