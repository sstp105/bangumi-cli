package login

import (
	"fmt"
	"github.com/sstp105/bangumi-cli/internal/console"
	"github.com/sstp105/bangumi-cli/internal/libs"
	"github.com/sstp105/bangumi-cli/internal/path"
	"os"
	"os/signal"

	"github.com/sstp105/bangumi-cli/internal/bangumi"
	"github.com/sstp105/bangumi-cli/internal/config"
	"github.com/sstp105/bangumi-cli/internal/log"
	"github.com/sstp105/bangumi-cli/internal/server"
)

var ch = make(chan os.Signal, 1)

// Handler handles bangumi login process.
func Handler() {
	if err := authenticate(); err != nil {
		console.Errorf("获取 bangumi 凭证失败:%s", err)
		log.Fatalf("error authenticating user:%s", err)
	}
}

// authenticate handles the process of checking for existing bangumi credentials,
// re-authenticating the user if necessary, and saving the credentials if they are refreshed.
//
// If new credentials are obtained, they are saved to a configuration file.
func authenticate() error {
	var credential *bangumi.OAuthCredential
	var err error
	overwrite := true

	err = path.ReadJSONConfigFile(path.BangumiCredentialFile, &credential)

	switch {
	case err != nil && os.IsNotExist(err):
		log.Debugf("bangumi credential config is not found, authenticating user.")
		credential, err = oauth()
	case err != nil:
		log.Errorf("error reading bangumi credential config:%s", err)
		return err
	case credential.IsExpired():
		log.Debugf("Access token already expired, re-authenticating user.")
		credential, err = oauth()
	case credential.ShouldRefresh():
		log.Debugf("Access token will expire soon, fetching new access token.")
		credential, err = refresh(credential.RefreshToken)
	default:
		log.Debugf("Bangumi credential config is valid.")
		overwrite = false
	}

	if err != nil {
		return err
	}

	if credential != nil && overwrite {
		if err := path.SaveJSONConfigFile(path.BangumiCredentialFile, credential); err != nil {
			log.Fatalf("error saving bangumi credentials:%s", err)
		}
	}

	return credential.Print()
}

// oauth starts the local server and listens for bangumi callback
// If the authentication succeeds, a bangumi credential will be returned.
func oauth() (*bangumi.OAuthCredential, error) {
	log.Info("Opening URL in the browser for bangumi.tv authentication...")

	clientID := config.BangumiClientID()
	url := fmt.Sprintf("https://bgm.tv/oauth/authorize?client_id=%s&response_type=code&redirect_uri=%s", clientID, config.LocalServerAddress())

	if err := libs.OpenBrowser(url); err != nil {
		return nil, err
	}

	signal.Notify(ch, os.Interrupt)

	var credential *bangumi.OAuthCredential
	go func() {
		server.Start(func(c *bangumi.OAuthCredential) {
			credential = c
			ch <- os.Interrupt
		})
	}()
	<-ch

	return credential, nil
}

// refresh uses token to request new bangumi access token.
func refresh(token string) (*bangumi.OAuthCredential, error) {
	client := bangumi.NewOAuthClient()
	credential, err := client.RefreshAccessToken(
		config.BangumiClientID(),
		config.BangumiClientSecret(),
		config.LocalServerAddress(),
		token,
	)
	if err != nil {
		return nil, err
	}

	return credential, nil
}
