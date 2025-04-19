package collecthandler

import (
	"github.com/sstp105/bangumi-cli/internal/bangumi"
	"github.com/sstp105/bangumi-cli/internal/path"
	"github.com/stretchr/testify/require"
	"os"
	"path/filepath"
	"testing"
)

func TestNewHandler(t *testing.T) {
	originalRunningOS := path.RunningOS
	originalProviders := path.OSPathProviders

	defer func() {
		path.RunningOS = originalRunningOS
		path.OSPathProviders = originalProviders
	}()

	t.Run("username is empty", func(t *testing.T) {
		h, err := NewHandler("", 3)

		require.Nil(t, h)
		require.Error(t, err)
		require.Equal(t, "username is empty", err.Error())
	})

	t.Run("invalid bangumi subject collection type", func(t *testing.T) {
		h, err := NewHandler("mock-username", 0)

		require.Nil(t, h)
		require.Error(t, err)
		require.Contains(t, err.Error(), "invalid collection type")
	})

	t.Run("user has not subscribed any bangumi", func(t *testing.T) {
		tmpDir := t.TempDir()

		path.RunningOS = "mockOS"
		path.OSPathProviders = map[string]path.Provider{
			"mockOS": mockPathProvider{
				configPathFunc: func() (string, error) {
					return tmpDir, nil
				},
			},
		}

		h, err := NewHandler("mock-username", 3)

		require.Nil(t, h)
		require.Equal(t, err.Error(), "subscription config file is empty")
	})

	t.Run("invalid subscription file", func(t *testing.T) {
		tmpDir := t.TempDir()

		tmpSubscriptionFile := filepath.Join(tmpDir, path.SubscriptionConfigFile)
		require.NoError(t, os.WriteFile(tmpSubscriptionFile, []byte(`{ invalid json ]`), 0644))

		path.RunningOS = "mockOS"
		path.OSPathProviders = map[string]path.Provider{
			"mockOS": mockPathProvider{
				configPathFunc: func() (string, error) {
					return tmpDir, nil
				},
			},
		}

		h, err := NewHandler("mock-username", 3)

		require.Nil(t, h)
		require.Error(t, err)
	})

	t.Run("user has not authenticated to bangumi.tv", func(t *testing.T) {
		tmpDir := t.TempDir()

		tmpSubscriptionFile := filepath.Join(tmpDir, path.SubscriptionConfigFile)
		subscription := `[{
			"id": "233",
			"name": "小林家的龙女仆",
			"link": "https://mikan.example.com/bangumi/233"
		}]`

		require.NoError(t, os.WriteFile(tmpSubscriptionFile, []byte(subscription), 0644))

		path.RunningOS = "mockOS"
		path.OSPathProviders = map[string]path.Provider{
			"mockOS": mockPathProvider{
				configPathFunc: func() (string, error) {
					return tmpDir, nil
				},
			},
		}

		h, err := NewHandler("mock-username", 3)

		require.Nil(t, h)
		require.Equal(t, err.Error(), "credential config file is empty")
	})

	t.Run("invalid bangumi credential config file", func(t *testing.T) {
		tmpDir := t.TempDir()

		tmpSubscriptionFile := filepath.Join(tmpDir, path.SubscriptionConfigFile)
		subscription := `[{
			"id": "233",
			"name": "小林家的龙女仆",
			"link": "https://mikan.example.com/bangumi/233"
		}]`

		require.NoError(t, os.WriteFile(tmpSubscriptionFile, []byte(subscription), 0644))

		tmpBangumiCredentialFile := filepath.Join(tmpDir, path.BangumiCredentialConfigFile)
		require.NoError(t, os.WriteFile(tmpBangumiCredentialFile, []byte(`{ invalid json ]`), 0644))

		path.RunningOS = "mockOS"
		path.OSPathProviders = map[string]path.Provider{
			"mockOS": mockPathProvider{
				configPathFunc: func() (string, error) {
					return tmpDir, nil
				},
			},
		}

		h, err := NewHandler("mock-username", 3)

		require.Nil(t, h)
		require.Error(t, err)
	})

	t.Run("success init", func(t *testing.T) {
		tmpDir := t.TempDir()

		tmpBangumiCredentialFile := filepath.Join(tmpDir, path.BangumiCredentialConfigFile)
		credential := `{
			"access_token": "mock-access-token",
			"refresh_token": "mock-refresh-token",
			"expires_in": 604800,
			"token_type": "Bearer",
			"expires_until": "2030-04-22T19:41:00.561143-07:00"
		}`

		require.NoError(t, os.WriteFile(tmpBangumiCredentialFile, []byte(credential), 0644))

		tmpSubscriptionFile := filepath.Join(tmpDir, path.SubscriptionConfigFile)
		subscription := `[{
			"id": "233",
			"name": "小林家的龙女仆",
			"link": "https://mikan.example.com/bangumi/233"
		}]`

		require.NoError(t, os.WriteFile(tmpSubscriptionFile, []byte(subscription), 0644))

		path.RunningOS = "mockOS"
		path.OSPathProviders = map[string]path.Provider{
			"mockOS": mockPathProvider{
				configPathFunc: func() (string, error) {
					return tmpDir, nil
				},
			},
		}

		h, err := NewHandler("mock-username", 3)

		require.NoError(t, err)
		require.NotNil(t, h)
		require.Equal(t, h.username, "mock-username")
		require.Equal(t, h.collectionType, bangumi.SubjectCollectionType(3))
		require.NotNil(t, h.subscription)
		require.NotNil(t, h.client)
	})
}

type mockPathProvider struct {
	configPathFunc func() (string, error)
}

func (m mockPathProvider) ConfigPath() (string, error) {
	return m.configPathFunc()
}
