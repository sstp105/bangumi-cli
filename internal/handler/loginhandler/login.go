package loginhandler

import (
	"fmt"
	"github.com/sstp105/bangumi-cli/internal/console"
	"github.com/sstp105/bangumi-cli/internal/libs"
	"github.com/sstp105/bangumi-cli/internal/path"
	"os"
	"os/signal"

	"github.com/sstp105/bangumi-cli/internal/bangumi"
	"github.com/sstp105/bangumi-cli/internal/config"
	"github.com/sstp105/bangumi-cli/internal/server"
)

var (
	ch = make(chan os.Signal, 1)

	oauthURL libs.APIPath = "https://bgm.tv/oauth/authorize?client_id=%s&response_type=code&redirect_uri=%s"
)

// Run triggers bangumi login process.
func Run() {
	if err := authenticate(); err != nil {
		console.Errorf("获取 bangumi 凭证失败:%s", err)
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

	err = path.ReadJSONConfigFile(path.BangumiCredentialConfigFile, &credential)

	switch {
	case err != nil && os.IsNotExist(err):
		credential, err = oauth()
	case err != nil:
		return err
	case credential.IsExpired():
		credential, err = oauth()
	case credential.ShouldRefresh():
		credential, err = refresh(credential.RefreshToken)
	default:
		overwrite = false
	}

	if err != nil {
		return err
	}

	if credential != nil && overwrite {
		if err := path.SaveJSONConfigFile(path.BangumiCredentialConfigFile, credential); err != nil {
			return fmt.Errorf("error saving bangumi credential file:%s", err)
		}
	}

	if err := credential.Print(); err != nil {
		return fmt.Errorf("error printing bangumi credential:%s", err)
	}

	return nil
}

// oauth starts a local server and listens for bangumi callback
// If the authentication succeeds, a bangumi credential will be returned.
func oauth() (*bangumi.OAuthCredential, error) {
	url := libs.FormatAPIPath(oauthURL, config.BangumiClientID(), config.LocalServerAddress())
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

// refresh uses refresh token to request new access token.
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
