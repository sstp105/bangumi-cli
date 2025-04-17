package torrent

import (
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

var (
	baseURL = "http://localhost:8080"

	mockUserName   = "mock-user"
	mockPassword   = "mock-password"
	mockAuthCookie = "mock-auth-cookie"
)

func TestAuthenticate_Error(t *testing.T) {
	withMockClient(t, func(c *QBittorrentClient) {
		httpmock.RegisterResponder("POST", QBittorrentAPILoginPath,
			httpmock.NewErrorResponder(errors.New("context canceled")),
		)

		cookie, err := c.authenticate()
		assert.Error(t, err)
		assert.Equal(t, "", cookie)
	})
}

func TestAuthenticate_NotAuthorized(t *testing.T) {
	withMockClient(t, func(c *QBittorrentClient) {
		httpmock.RegisterResponder("POST", QBittorrentAPILoginPath,
			httpmock.NewJsonResponderOrPanic(403, map[string]string{
				"title":       "Unauthorized",
				"description": "The user name and password are invalid",
			}))

		cookie, err := c.authenticate()
		assert.Error(t, err)
		assert.Equal(t, "", cookie)
	})
}

func TestAuthenticate_AuthCookieNotFound(t *testing.T) {
	withMockClient(t, func(c *QBittorrentClient) {
		httpmock.RegisterResponder("POST", QBittorrentAPILoginPath,
			func(req *http.Request) (*http.Response, error) {
				resp := httpmock.NewStringResponse(200, "")
				return resp, nil
			},
		)

		cookie, err := c.authenticate()
		assert.Error(t, err)
		assert.Equal(t, "", cookie)
	})
}

func TestAuthenticate_Success(t *testing.T) {
	withMockClient(t, func(c *QBittorrentClient) {
		httpmock.RegisterResponder("POST", QBittorrentAPILoginPath,
			func(req *http.Request) (*http.Response, error) {
				resp := httpmock.NewStringResponse(200, "")
				resp.Header.Set("Set-Cookie", fmt.Sprintf("SID=%s; Path=/; HttpOnly", mockAuthCookie))
				return resp, nil
			},
		)

		cookie, err := c.authenticate()
		assert.NoError(t, err)
		assert.Equal(t, mockAuthCookie, cookie)
	})
}

func newMockClient() *QBittorrentClient {
	client := resty.New()
	client.SetBaseURL(baseURL)

	cfg := QBittorrentClientConfig{
		Server:   baseURL,
		Username: mockUserName,
		Password: mockPassword,
	}

	httpmock.ActivateNonDefault(client.GetClient())

	return &QBittorrentClient{
		client: client,
		config: cfg,
	}
}

func withMockClient(t *testing.T, testFn func(c *QBittorrentClient)) {
	client := newMockClient()
	defer httpmock.DeactivateAndReset()
	testFn(client)
}
