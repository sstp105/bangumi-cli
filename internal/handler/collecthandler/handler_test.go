package collecthandler

import (
	"github.com/go-resty/resty/v2"
	"github.com/jarcoal/httpmock"
	"github.com/sstp105/bangumi-cli/internal/model"
	"os"
	"path/filepath"
	"testing"

	"github.com/sstp105/bangumi-cli/internal/bangumi"
	"github.com/sstp105/bangumi-cli/internal/path"
	"github.com/stretchr/testify/require"
)

type mockPathProvider struct {
	configPathFunc func() (string, error)
}

func (m mockPathProvider) ConfigPath() (string, error) {
	return m.configPathFunc()
}

func newMockClient() *bangumi.Client {
	client := resty.New()
	httpmock.ActivateNonDefault(client.GetClient())
	return bangumi.NewClient(bangumi.WithClient(client.GetClient()))
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

func removeFile(t *testing.T, path string) {
	t.Helper()
	require.NoError(t, os.Remove(path))
}

func TestNewHandler(t *testing.T) {
	originalRunningOS := path.RunningOS
	originalProviders := path.OSPathProviders
	defer func() {
		path.RunningOS = originalRunningOS
		path.OSPathProviders = originalProviders
	}()

	t.Run("empty username", func(t *testing.T) {
		h, err := NewHandler("", 3)
		require.Nil(t, h)
		require.EqualError(t, err, "username is empty")
	})

	t.Run("invalid collection type", func(t *testing.T) {
		h, err := NewHandler("mock-username", 0)
		require.Nil(t, h)
		require.EqualError(t, err, "invalid collection type 0")
	})

	t.Run("no subscription config", func(t *testing.T) {
		tmpDir := t.TempDir()
		setMockPathProvider(t, tmpDir)

		h, err := NewHandler("mock-username", 3)
		require.Nil(t, h)
		require.EqualError(t, err, "subscription config file is empty")
	})

	t.Run("invalid subscription file", func(t *testing.T) {
		tmpDir := t.TempDir()
		writeFile(t, filepath.Join(tmpDir, path.SubscriptionConfigFile), `{ invalid json ]`)
		setMockPathProvider(t, tmpDir)

		h, err := NewHandler("mock-username", 3)
		require.Nil(t, h)
		require.Error(t, err)
	})

	t.Run("missing credential config", func(t *testing.T) {
		tmpDir := t.TempDir()
		subscription := `[{"id":"233","name":"小林家的龙女仆","link":"https://mikan.example.com/bangumi/233"}]`
		writeFile(t, filepath.Join(tmpDir, path.SubscriptionConfigFile), subscription)
		setMockPathProvider(t, tmpDir)

		h, err := NewHandler("mock-username", 3)
		require.Nil(t, h)
		require.EqualError(t, err, "credential config file is empty")
	})

	t.Run("invalid credential config", func(t *testing.T) {
		tmpDir := t.TempDir()
		writeFile(t, filepath.Join(tmpDir, path.SubscriptionConfigFile), `[{"id":"233","name":"小林家的龙女仆","link":"https://mikan.example.com/bangumi/233"}]`)
		writeFile(t, filepath.Join(tmpDir, path.BangumiCredentialConfigFile), `{ invalid json ]`)
		setMockPathProvider(t, tmpDir)

		h, err := NewHandler("mock-username", 3)
		require.Nil(t, h)
		require.Error(t, err)
	})

	t.Run("successful init", func(t *testing.T) {
		tmpDir := t.TempDir()

		credential := `{
			"access_token": "mock-access-token",
			"refresh_token": "mock-refresh-token",
			"expires_in": 604800,
			"token_type": "Bearer",
			"expires_until": "2030-04-22T19:41:00.561143-07:00"
		}`
		writeFile(t, filepath.Join(tmpDir, path.BangumiCredentialConfigFile), credential)

		subscription := `[{"id":"233","name":"小林家的龙女仆","link":"https://mikan.example.com/bangumi/233"}]`
		writeFile(t, filepath.Join(tmpDir, path.SubscriptionConfigFile), subscription)

		setMockPathProvider(t, tmpDir)

		h, err := NewHandler("mock-username", 3)
		require.NoError(t, err)
		require.NotNil(t, h)
		require.Equal(t, "mock-username", h.username)
		require.Equal(t, bangumi.SubjectCollectionType(3), h.collectionType)
		require.NotNil(t, h.subscription)
		require.NotNil(t, h.client)
	})
}

func TestHandler_Run(t *testing.T) {
	tmpDir := t.TempDir()
	setup := func(t *testing.T) (*Handler, string) {
		// create necessary config files
		writeFile(t, filepath.Join(tmpDir, path.BangumiCredentialConfigFile), `{
			"access_token": "mock-access-token",
			"refresh_token": "mock-refresh-token",
			"expires_in": 604800,
			"token_type": "Bearer",
			"expires_until": "2030-04-22T19:41:00.561143-07:00"
		}`)

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
			"rss_link": "/RSS/Bangumi?bangumiId=3519\u0026subgroupid=382"
		}`)

		setMockPathProvider(t, tmpDir)

		h := &Handler{
			username:       "mock-username",
			collectionType: bangumi.SubjectCollectionType(3),
			subscription: []model.BangumiBase{
				{
					ID:   "3519",
					Name: "金牌得主",
					Link: "/Home/Bangumi/3519",
				},
			},
			client: newMockClient(),
		}

		return h, tmpDir
	}

	t.Run("returns error when bangumi config does not exist", func(t *testing.T) {
		h, _ := setup(t)
		defer httpmock.DeactivateAndReset()

		// remove from setup
		removeFile(t, filepath.Join(tmpDir, "3519.json"))

		err := h.Run()
		require.Error(t, err)
		require.Contains(t, err.Error(), "failed to read bangumi config file")
	})

	t.Run("returns error when bangumi id does not exist in config file", func(t *testing.T) {
		h, _ := setup(t)
		defer httpmock.DeactivateAndReset()

		// overwrite
		writeFile(t, filepath.Join(tmpDir, "3519.json"), `{
			"id": "3519",
			"name": "金牌得主",
			"link": "/Home/Bangumi/3519",
			"bangumi_id": "",
			"rss_link": "/RSS/Bangumi?bangumiId=3519\u0026subgroupid=382"
		}`)

		err := h.Run()
		require.Error(t, err)
		require.Contains(t, err.Error(), "bangumi id is empty")
	})

	t.Run("returns error when bangumi API fails to fetch collection status", func(t *testing.T) {
		h, _ := setup(t)
		defer httpmock.DeactivateAndReset()

		httpmock.RegisterResponder(
			"GET",
			"https://api.bgm.tv/v0/users/mock-username/collections/430699",
			httpmock.NewJsonResponderOrPanic(502, map[string]string{
				"title":       "Bad Gateway",
				"description": "The server is under maintenance",
			}),
		)

		err := h.Run()
		require.Error(t, err)
		require.Contains(t, err.Error(), "failed to fetch collection status")
	})

	t.Run("successfully create collection status", func(t *testing.T) {
		h, _ := setup(t)
		defer httpmock.DeactivateAndReset()

		httpmock.RegisterResponder(
			"GET",
			"https://api.bgm.tv/v0/users/mock-username/collections/430699",
			httpmock.NewJsonResponderOrPanic(404, map[string]string{
				"title":       "Not Found",
				"description": "The subject is not collected by user",
			}),
		)

		httpmock.RegisterResponder(
			"POST",
			"https://api.bgm.tv/v0/users/-/collections/430699",
			httpmock.NewStringResponder(200, ""),
		)

		err := h.Run()
		require.NoError(t, err)
	})

	t.Run("successfully updates collection status", func(t *testing.T) {
		h, _ := setup(t)
		defer httpmock.DeactivateAndReset()

		httpmock.RegisterResponder(
			"GET",
			"https://api.bgm.tv/v0/users/mock-username/collections/430699",
			httpmock.NewJsonResponderOrPanic(200, map[string]interface{}{
				"subject": map[string]interface{}{
					"id":      430699,
					"type":    2,
					"name":    "メダリスト",
					"name_cn": "金牌得主",
				},
				"type": 3,
			}),
		)

		httpmock.RegisterResponder(
			"PATCH",
			"https://api.bgm.tv/v0/users/-/collections/430699",
			httpmock.NewStringResponder(200, ""),
		)

		err := h.Run()
		require.NoError(t, err)
	})
}
