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

	mockUserName     = "mock-user"
	mockPassword     = "mock-password"
	mockAuthCookie   = "mock-auth-cookie"
	mockTorrentURLs  = "https://example.com/torrent_1\nhttps://example.com/torrent_2"
	mockDownloadDest = "/mock-output"
)

func TestQBittorrent_Authenticate_Error(t *testing.T) {
	withMockClient(t, func(c *QBittorrentClient) {
		httpmock.RegisterResponder("POST", QBittorrentAPILoginPath,
			httpmock.NewErrorResponder(errors.New("context canceled")),
		)

		cookie, err := c.authenticate()

		assert.Error(t, err)
		assert.Equal(t, "", cookie)
	})
}

func TestQBittorrent_Authenticate_NotAuthorized(t *testing.T) {
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

func TestQBittorrent_Authenticate_AuthCookieNotFound(t *testing.T) {
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

func TestQBittorrent_Authenticate_Success(t *testing.T) {
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

func TestQBittorrent_Add_Error(t *testing.T) {
	withMockClient(t, func(c *QBittorrentClient) {
		httpmock.RegisterResponder("POST", QBittorrentAPIAddPath,
			httpmock.NewErrorResponder(errors.New("context canceled")),
		)

		err := c.Add(mockTorrentURLs, mockDownloadDest)

		assert.Error(t, err)
	})
}

func TestQBittorrent_Add_InternalServerResponse(t *testing.T) {
	withMockClient(t, func(c *QBittorrentClient) {
		httpmock.RegisterResponder("POST", QBittorrentAPIAddPath,
			httpmock.NewJsonResponderOrPanic(500, map[string]string{
				"title":       "Internal server error",
				"description": "The request failed due to an internal server error",
			}))

		err := c.Add(mockTorrentURLs, mockDownloadDest)

		assert.Error(t, err)
	})
}

func TestQBittorrent_Add_Success(t *testing.T) {
	withMockClient(t, func(c *QBittorrentClient) {
		httpmock.RegisterResponder("POST", QBittorrentAPIAddPath,
			func(req *http.Request) (*http.Response, error) {
				resp := httpmock.NewStringResponse(200, "")
				return resp, nil
			},
		)

		err := c.Add(mockTorrentURLs, mockDownloadDest)

		assert.NoError(t, err)
	})
}

func TestQBittorrent_Name(t *testing.T) {
	qbt := &QBittorrentClient{}
	name := qbt.Name()

	assert.Equal(t, QBittorrentClientName, name)
}

func TestQBittorrent_validate_InvalidServer(t *testing.T) {
	cfg := QBittorrentClientConfig{
		Server:   "",
		Username: mockUserName,
		Password: mockPassword,
	}

	err := cfg.validate()

	assert.Error(t, err)
	assert.Equal(t, ErrQBittorrentServerNotFound, err)
}

func TestQBittorrent_validate_InvalidUserName(t *testing.T) {
	cfg := QBittorrentClientConfig{
		Server:   baseURL,
		Username: "",
		Password: mockPassword,
	}

	err := cfg.validate()

	assert.Error(t, err)
	assert.Equal(t, ErrQBittorrentUserNameNotFound, err)
}

func TestQBittorrent_validate_InvalidPassword(t *testing.T) {
	cfg := QBittorrentClientConfig{
		Server:   baseURL,
		Username: mockUserName,
		Password: "",
	}

	err := cfg.validate()

	assert.Error(t, err)
	assert.Equal(t, ErrQBittorrentPasswordNotFound, err)
}

func TestQBittorrent_validate_ValidConfig(t *testing.T) {
	cfg := QBittorrentClientConfig{
		Server:   baseURL,
		Username: mockUserName,
		Password: mockPassword,
	}

	err := cfg.validate()

	assert.NoError(t, err)
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
