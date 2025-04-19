package collecthandler

import (
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
