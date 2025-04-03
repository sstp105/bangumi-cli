package bangumi

import (
	"github.com/go-resty/resty/v2"
)

const (
	// oauthBaseURL is the base URL for the bangumi oauth API.
	oAuthBaseURL string = "https://bgm.tv"
)

// OAuthClient wraps a resty client for interacting with the bangumi API.
type OAuthClient struct {
	client *resty.Client
}

// NewOAuthClient creates and returns a new instance of the OAuthClient.
func NewOAuthClient() *OAuthClient {
	c := resty.New()
	c.SetBaseURL(oAuthBaseURL)

	return &OAuthClient{
		client: c,
	}
}

func (o *OAuthClient) GetAccessToken(clientID, clientSecret, redirectURI, code string) (*OAuthCredential, error) {
	creds := OAuthCredential{}
	errorResp := OAuthErrorResponse{}

	data := map[string]string{
		"grant_type":    "authorization_code",
		"client_id":     clientID,
		"client_secret": clientSecret,
		"code":          code,
		"redirect_uri":  redirectURI,
	}
	url := postAccessTokenURL()

	resp, err := o.client.R().
		SetHeader("Content-Type", "application/x-www-form-urlencoded").
		SetFormData(data).
		SetResult(&creds).
		SetError(&errorResp).
		Post(url)

	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		return nil, &errorResp
	}

	creds.setExpiresUntil()

	return &creds, nil
}

func (o *OAuthClient) RefreshAccessToken(clientID, clientSecret, redirectURI, refreshToken string) (*OAuthCredential, error) {
	creds := OAuthCredential{}
	errorResp := OAuthErrorResponse{}

	data := map[string]string{
		"grant_type":    "refresh_token",
		"client_id":     clientID,
		"client_secret": clientSecret,
		"refresh_token": refreshToken,
		"redirect_uri":  redirectURI,
	}
	url := postAccessTokenURL()

	resp, err := o.client.R().
		SetHeader("Content-Type", "application/x-www-form-urlencoded").
		SetFormData(data).
		SetResult(&creds).
		SetError(&errorResp).
		Post(url)

	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		return nil, &errorResp
	}

	creds.setExpiresUntil()

	return &creds, nil
}
