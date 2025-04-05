package bangumi

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/sstp105/bangumi-cli/internal/libs"
	"time"
)

// Constants
const (
	// oAuthBaseURL is the base URL for the Bangumi OAuth API.
	oAuthBaseURL string = "https://bgm.tv"
)

// OAuthErrorResponse represents the error response structure from the OAuth API.
type OAuthErrorResponse struct {
	ErrorCode        string `json:"error"`
	ErrorDescription string `json:"error_description"`
}

// Error implements the error interface for OAuthErrorResponse.
func (o *OAuthErrorResponse) Error() string {
	return fmt.Sprintf("bangumi oauth api error: %s - %s", o.ErrorCode, o.ErrorDescription)
}

// OAuthCredential represents the OAuth credentials.
type OAuthCredential struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresIn    int64     `json:"expires_in"`
	TokenType    string    `json:"token_type"`
	ExpiresUntil time.Time `json:"expires_until,omitempty"` // custom injected field
}

// IsValid checks if the OAuthCredential AccessToken is valid (not expired).
func (o *OAuthCredential) IsValid() bool {
	return time.Now().Before(o.ExpiresUntil)
}

// ShouldRefresh checks if the  OAuthCredential AccessToken should be refreshed.
func (o *OAuthCredential) ShouldRefresh() bool {
	return o.ExpiresUntil.Before(time.Now().Add(24 * time.Hour))
}

// IsExpired checks if the OAuthCredential AccessToken has expired.
func (o *OAuthCredential) IsExpired() bool {
	return o.ExpiresUntil.Before(time.Now())
}

// Print prints the OAuthCredential in JSON format.
func (o *OAuthCredential) Print() error {
	data, err := libs.MarshalJSONIndented(o)
	if err != nil {
		return err
	}
	fmt.Println(string(data))
	return nil
}

// setExpiresUntil sets the expiration time for the OAuth credential.
func (o *OAuthCredential) setExpiresUntil() {
	o.ExpiresUntil = time.Now().Add(time.Second * time.Duration(o.ExpiresIn))
}

// OAuthClient wraps a resty.Client for interacting with the Bangumi API.
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

// GetAccessToken gets an OAuthCredential using an authorization code.
func (o *OAuthClient) GetAccessToken(clientID, clientSecret, redirectURI, code string) (*OAuthCredential, error) {
	return o.requestOAuthToken("authorization_code", clientID, clientSecret, redirectURI, code)
}

// RefreshAccessToken refreshes the OAuthCredential using a refresh token.
func (o *OAuthClient) RefreshAccessToken(clientID, clientSecret, redirectURI, refreshToken string) (*OAuthCredential, error) {
	return o.requestOAuthToken("refresh_token", clientID, clientSecret, redirectURI, refreshToken)
}

// requestOAuthToken handles the logic for requesting OAuthCredential (both new and refreshed).
func (o *OAuthClient) requestOAuthToken(grantType, clientID, clientSecret, redirectURI, codeOrRefreshToken string) (*OAuthCredential, error) {
	credential := OAuthCredential{}
	errorResp := OAuthErrorResponse{}

	data := map[string]string{
		"grant_type":    grantType,
		"client_id":     clientID,
		"client_secret": clientSecret,
		"redirect_uri":  redirectURI,
	}
	
	if grantType == "authorization_code" {
		data["code"] = codeOrRefreshToken
	} else if grantType == "refresh_token" {
		data["refresh_token"] = codeOrRefreshToken
	}

	url := libs.FormatAPIPath(postAccessToken)

	resp, err := o.client.R().
		SetHeader("Content-Type", "application/x-www-form-urlencoded").
		SetFormData(data).
		SetResult(&credential).
		SetError(&errorResp).
		Post(url)

	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		return nil, &errorResp
	}

	credential.setExpiresUntil()

	return &credential, nil
}
