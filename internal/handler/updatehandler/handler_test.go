package updatehandler

import (
	"github.com/sstp105/bangumi-cli/internal/mikan"
	"github.com/sstp105/bangumi-cli/internal/path"
	"github.com/stretchr/testify/require"
	"os"
	"path/filepath"
	"testing"
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
	cfg := mikan.ClientConfig{
		IdentityCookie: "mock-identity",
	}

	originalRunningOS := path.RunningOS
	originalProviders := path.OSPathProviders
	defer func() {
		path.RunningOS = originalRunningOS
		path.OSPathProviders = originalProviders
	}()

	t.Run("read subscription file failure", func(t *testing.T) {
		tmpDir := t.TempDir()

		subscription := `{invalid json}`
		writeFile(t, filepath.Join(tmpDir, path.SubscriptionConfigFile), subscription)

		setMockPathProvider(t, tmpDir)

		h, err := NewHandler(cfg)
		require.Error(t, err)
		require.Nil(t, h)
	})

	t.Run("subscription file not exist", func(t *testing.T) {
		tmpDir := t.TempDir()

		setMockPathProvider(t, tmpDir)

		h, err := NewHandler(cfg)
		require.Error(t, err)
		require.Contains(t, err.Error(), "no subscription config found")
		require.Nil(t, h)
	})

	t.Run("invalid client config", func(t *testing.T) {
		tmpDir := t.TempDir()

		subscription := `[{"id":"233","name":"小林家的龙女仆","link":"https://mikan.example.com/bangumi/233"}]`
		writeFile(t, filepath.Join(tmpDir, path.SubscriptionConfigFile), subscription)

		setMockPathProvider(t, tmpDir)

		h, err := NewHandler(mikan.ClientConfig{})
		require.Error(t, err)
		require.Contains(t, err.Error(), ".AspNetCore.Identity.Application is empty")
		require.Nil(t, h)
	})

	t.Run("success", func(t *testing.T) {
		tmpDir := t.TempDir()

		subscription := `[{"id":"233","name":"小林家的龙女仆","link":"https://mikan.example.com/bangumi/233"}]`
		writeFile(t, filepath.Join(tmpDir, path.SubscriptionConfigFile), subscription)

		setMockPathProvider(t, tmpDir)

		h, err := NewHandler(cfg)
		require.NoError(t, err)
		require.NotNil(t, h)
	})
}
