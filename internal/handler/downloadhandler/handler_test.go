package downloadhandler

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/jarcoal/httpmock"
	"github.com/sstp105/bangumi-cli/internal/model"
	"github.com/sstp105/bangumi-cli/internal/path"
	"github.com/sstp105/bangumi-cli/internal/torrent"
	"github.com/stretchr/testify/require"
)

type mockPathProvider struct {
	configPathFunc func() (string, error)
	downloadPathFunc func() (string, error)
}

func (m mockPathProvider) ConfigPath() (string, error) {
	return m.configPathFunc()
}


func (m mockPathProvider) DownloadPath() (string, error) {
	return m.downloadPathFunc()
}

func newMockClient() (*torrent.QBittorrentClient, error) {
	client := resty.New()
	httpmock.ActivateNonDefault(client.GetClient())
	httpmock.RegisterResponder("POST", "/api/v2/auth/login",
		func(req *http.Request) (*http.Response, error) {
			resp := httpmock.NewStringResponse(200, "")
			resp.Header.Set("Set-Cookie", fmt.Sprintf("SID=%s; Path=/; HttpOnly", "mock-auth-cookie"))
			return resp, nil
		},
	)
	cfg := torrent.QBittorrentClientConfig{
		Server:   "http://localhost:8888",
		Username: "admin",
		Password: "mock-password",
	}
	return torrent.NewQBittorrentClient(cfg, torrent.WithHTTPClient(client.GetClient()))
}

func setMockPathProvider(t *testing.T, configDir string) {
	t.Helper()
	path.RunningOS = "mockOS"
	path.OSPathProviders = map[string]path.Provider{
		"mockOS": mockPathProvider{
			configPathFunc: func() (string, error) {
				return configDir, nil
			},
		},
	}
}

func writeFile(t *testing.T, path, content string) {
	t.Helper()
	require.NoError(t, os.WriteFile(path, []byte(content), 0644))
}

func TestHandler_Run(t *testing.T) {
	tmpDir := t.TempDir()
	setup := func(t *testing.T) (*Handler, string) {
		writeFile(t, filepath.Join(tmpDir, path.SubscriptionConfigFile), `[{
			"id": "3519",
			"name": "金牌得主",
			"link": "/Home/Bangumi/3519"
		}]`)

		writeFile(t, filepath.Join(tmpDir, "3519.json"), `{
			"id": "3519",
			"name": "金牌得主",
			"link": "/Home/Bangumi/3519",
			"bangumi_id": "430699",
			"rss_link": "/RSS/Bangumi?bangumiId=3519\u0026subgroupid=382",
			"torrents": [
				"https://example.com/Download/1.torrent",
				"https://example.com/Download/2.torrent",
				"https://example.com/Download/3.torrent",
				"https://example.com/Download/4.torrent"
			]
		}`)

		setMockPathProvider(t, tmpDir)

		client, _ := newMockClient()
		h := &Handler{
			subscription: []model.BangumiBase{
				{
					ID:   "3519",
					Name: "金牌得主",
					Link: "/Home/Bangumi/3519",
				},
			},
			client: client,
		}

		return h, tmpDir
	}

	t.Run("returns error when bangumi config does not exist", func(t *testing.T) {
		h, _ := setup(t)
		defer httpmock.DeactivateAndReset()

		h.subscription = nil

		err := h.Run()
		require.NoError(t, err)
	})

	t.Run("failed to add torrents with qbittorrent api error", func(t *testing.T) {
		h, _ := setup(t)

		httpmock.RegisterResponder("POST", torrent.QBittorrentAPIAddPath,
			httpmock.NewJsonResponderOrPanic(500, map[string]string{
				"title":       "Internal server error",
				"description": "The request failed due to an internal server error",
			}))
		defer httpmock.DeactivateAndReset()

		err := h.Run()
		require.Error(t, err)
		require.Contains(t, err.Error(), "failed to add torrent")
	})

	t.Run("failed to read bangumi config file", func(t *testing.T) {
		h, _ := setup(t)

		defer httpmock.DeactivateAndReset()

		// overwrite
		writeFile(t, filepath.Join(tmpDir, "3519.json"), `{bad jason`)

		err := h.Run()
		require.Error(t, err)
		require.Contains(t, err.Error(), "failed to read config file")
	})

	t.Run("add torrents successfully", func(t *testing.T) {
		h, _ := setup(t)

		httpmock.RegisterResponder("POST", "/api/v2/torrents/add",
			func(req *http.Request) (*http.Response, error) {
				resp := httpmock.NewStringResponse(200, "")
				return resp, nil
			},
		)
		defer httpmock.DeactivateAndReset()

		err := h.Run()
		require.NoError(t, err)
	})
}
