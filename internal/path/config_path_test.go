package path

import (
	"errors"
	"github.com/sstp105/bangumi-cli/internal/libs"
	"github.com/stretchr/testify/require"
	"os"
	"path/filepath"
	"testing"
)

func TestWindowsPath_ConfigPath(t *testing.T) {
	env := os.Getenv("AppData")

	t.Run("%AppData% is set", func(t *testing.T) {
		tmpDir := t.TempDir()
		_ = os.Setenv("AppData", tmpDir)
		defer os.Setenv("AppData", env)

		wp := WindowsPath{}
		got, err := wp.ConfigPath()

		require.NoError(t, err)
		require.Equal(t, filepath.Join(tmpDir, AppDir), got)
	})

	t.Run("%AppData% is not set", func(t *testing.T) {
		_ = os.Unsetenv("AppData")
		defer os.Setenv("AppData", env)

		wp := WindowsPath{}
		_, err := wp.ConfigPath()
		require.Error(t, err)
		require.EqualError(t, err, "%AppData% is not defined")
	})
}

func TestMacOSPath_ConfigPath(t *testing.T) {
	env := os.Getenv("HOME")

	t.Run("$HOME is set", func(t *testing.T) {
		tmpDir := t.TempDir()
		require.NoError(t, os.Setenv("Home", tmpDir))
		defer os.Setenv("HOME", env)

		m := MacOSPath{}
		got, err := m.ConfigPath()

		require.NoError(t, err)
		require.Equal(t, filepath.Join(tmpDir, ".config", AppDir), got)
	})

	t.Run("$HOME is not set", func(t *testing.T) {
		require.NoError(t, os.Unsetenv("Home"))
		defer os.Setenv("HoME", env)

		m := MacOSPath{}
		_, err := m.ConfigPath()

		require.Error(t, err)
		require.EqualError(t, err, "$HOME is not defined")
	})
}

type mockPathProvider struct {
	configPathFunc func() (string, error)
}

func (m mockPathProvider) ConfigPath() (string, error) {
	return m.configPathFunc()
}

func TestConfigPath(t *testing.T) {
	originalRunningOS := runningOS
	originalProviders := osPathProviders

	defer func() {
		runningOS = originalRunningOS
		osPathProviders = originalProviders
	}()

	t.Run("supported OS with successful path provider", func(t *testing.T) {
		runningOS = "macos"
		osPathProviders = map[string]Provider{
			"macos": mockPathProvider{
				configPathFunc: func() (string, error) {
					return "/Users/test/.config/mock-app-name", nil
				},
			},
		}

		path, err := configPath()

		require.NoError(t, err)
		require.Equal(t, "/Users/test/.config/mock-app-name", path)
	})

	t.Run("supported OS but provider returns error", func(t *testing.T) {
		runningOS = "windows"
		osPathProviders = map[string]Provider{
			"windows": mockPathProvider{
				configPathFunc: func() (string, error) {
					return "", errors.New("%AppData% is not defined")
				},
			},
		}

		_, err := configPath()

		require.Error(t, err)
		require.EqualError(t, err, "%AppData% is not defined")
	})

	t.Run("unsupported OS", func(t *testing.T) {
		runningOS = "plan9"
		osPathProviders = map[string]Provider{}

		_, err := configPath()

		require.Error(t, err)
		require.Equal(t, libs.ErrUnsupportedOS, err)
	})
}

func TestReadJSONConfigFile(t *testing.T) {
	originalRunningOS := runningOS
	originalProviders := osPathProviders

	defer func() {
		runningOS = originalRunningOS
		osPathProviders = originalProviders
	}()

	type Config struct {
		Test bool `json:"test"`
	}

	t.Run("successful read and unmarshal", func(t *testing.T) {
		tmpDir := t.TempDir()
		tmpFile := "test_config.json"
		fullPath := filepath.Join(tmpDir, tmpFile)

		expected := `{"test": true}`
		require.NoError(t, os.WriteFile(fullPath, []byte(expected), 0644))

		runningOS = "mockOS"
		osPathProviders = map[string]Provider{
			"mockOS": mockPathProvider{
				configPathFunc: func() (string, error) {
					return tmpDir, nil
				},
			},
		}

		var cfg Config
		err := ReadJSONConfigFile(tmpFile, &cfg)
		require.NoError(t, err)
		require.True(t, cfg.Test)
	})

	t.Run("configPath returns error", func(t *testing.T) {
		tmpDir := t.TempDir()

		runningOS = "mockOS"
		osPathProviders = map[string]Provider{
			"mockOS": mockPathProvider{
				configPathFunc: func() (string, error) {
					return tmpDir, os.ErrPermission
				},
			},
		}

		var cfg Config
		err := ReadJSONConfigFile("irrelevant.json", &cfg)
		require.Error(t, err)
		require.Equal(t, os.ErrPermission, err)
	})

	t.Run("file does not exist", func(t *testing.T) {
		tmpDir := t.TempDir()

		runningOS = "mockOS"
		osPathProviders = map[string]Provider{
			"mockOS": mockPathProvider{
				configPathFunc: func() (string, error) {
					return tmpDir, nil
				},
			},
		}

		var cfg Config
		err := ReadJSONConfigFile("nonexistent.json", &cfg)
		require.Error(t, err)
		require.True(t, os.IsNotExist(err))
	})

	t.Run("invalid JSON format", func(t *testing.T) {
		tmpDir := t.TempDir()
		tmpFile := "bad.json"
		fullPath := filepath.Join(tmpDir, tmpFile)

		require.NoError(t, os.WriteFile(fullPath, []byte(`{invalid_json}`), 0644))

		runningOS = "mockOS"
		osPathProviders = map[string]Provider{
			"mockOS": mockPathProvider{
				configPathFunc: func() (string, error) {
					return tmpDir, nil
				},
			},
		}

		var cfg Config
		err := ReadJSONConfigFile(tmpFile, &cfg)
		require.Error(t, err)
	})
}

func TestDeleteJSONConfigFile(t *testing.T) {
	originalRunningOS := runningOS
	originalProviders := osPathProviders

	defer func() {
		runningOS = originalRunningOS
		osPathProviders = originalProviders
	}()

	t.Run("successful delete", func(t *testing.T) {
		tmpDir := t.TempDir()
		tmpFile := "test_config.json"
		fullPath := filepath.Join(tmpDir, tmpFile)

		// create the tmp file to delete
		require.NoError(t, os.WriteFile(fullPath, []byte(`{"test": true}`), 0644))

		runningOS = "mockOS"
		osPathProviders = map[string]Provider{
			"mockOS": mockPathProvider{
				configPathFunc: func() (string, error) {
					return tmpDir, nil
				},
			},
		}

		err := DeleteJSONConfigFile(tmpFile)
		require.NoError(t, err)

		_, err = os.Stat(fullPath)
		require.True(t, os.IsNotExist(err))
	})

	t.Run("delete file errors", func(t *testing.T) {
		tmpDir := t.TempDir()

		runningOS = "mockOS"
		osPathProviders = map[string]Provider{
			"mockOS": mockPathProvider{
				configPathFunc: func() (string, error) {
					return tmpDir, os.ErrPermission
				},
			},
		}

		err := DeleteJSONConfigFile("file-require-permission.json")
		require.Error(t, err)
		require.Equal(t, os.ErrPermission, err)
	})

	t.Run("file does not exist", func(t *testing.T) {
		tmpDir := t.TempDir()

		runningOS = "mockOS"
		osPathProviders = map[string]Provider{
			"mockOS": mockPathProvider{
				configPathFunc: func() (string, error) {
					return tmpDir, nil
				},
			},
		}

		err := DeleteJSONConfigFile("non-existent-file.json")
		require.Error(t, err)
		require.True(t, os.IsNotExist(err))
	})
}

func TestSaveJSONConfigFile(t *testing.T) {
	originalRunningOS := runningOS
	originalProviders := osPathProviders

	defer func() {
		runningOS = originalRunningOS
		osPathProviders = originalProviders
	}()

	type Config struct {
		Enabled bool `json:"enabled"`
	}

	t.Run("successful save", func(t *testing.T) {
		tmpDir := t.TempDir()
		tmpFile := "config.json"
		fullPath := filepath.Join(tmpDir, tmpFile)

		runningOS = "mockOS"
		osPathProviders = map[string]Provider{
			"mockOS": mockPathProvider{
				configPathFunc: func() (string, error) {
					return tmpDir, nil
				},
			},
		}

		cfg := Config{Enabled: true}
		err := SaveJSONConfigFile(tmpFile, cfg)
		require.NoError(t, err)

		data, err := os.ReadFile(fullPath)
		require.NoError(t, err)
		require.JSONEq(t, `{"enabled": true}`, string(data))
	})

	t.Run("configPath returns error", func(t *testing.T) {
		runningOS = "mockOS"
		osPathProviders = map[string]Provider{
			"mockOS": mockPathProvider{
				configPathFunc: func() (string, error) {
					return "", os.ErrPermission
				},
			},
		}

		cfg := Config{Enabled: false}
		err := SaveJSONConfigFile("any.json", cfg)
		require.Error(t, err)
		require.Equal(t, os.ErrPermission, err)
	})

	t.Run("folder creation error", func(t *testing.T) {
		badPath := string([]byte{0}) // invalid path on most systems

		runningOS = "mockOS"
		osPathProviders = map[string]Provider{
			"mockOS": mockPathProvider{
				configPathFunc: func() (string, error) {
					return badPath, nil
				},
			},
		}

		cfg := Config{Enabled: false}
		err := SaveJSONConfigFile("file.json", cfg)
		require.Error(t, err)
	})

	t.Run("marshal error", func(t *testing.T) {
		runningOS = "mockOS"
		osPathProviders = map[string]Provider{
			"mockOS": mockPathProvider{
				configPathFunc: func() (string, error) {
					return t.TempDir(), nil
				},
			},
		}

		// Channels cannot be marshaled to JSON
		cfg := struct {
			Ch chan int `json:"ch"`
		}{Ch: make(chan int)}

		err := SaveJSONConfigFile("bad.json", cfg)
		require.Error(t, err)
	})
}

func TestReadSubscriptionConfigFile(t *testing.T) {
	originalRunningOS := runningOS
	originalProviders := osPathProviders

	defer func() {
		runningOS = originalRunningOS
		osPathProviders = originalProviders
	}()

	t.Run("successful read and parse", func(t *testing.T) {
		tmpDir := t.TempDir()
		tmpFile := filepath.Join(tmpDir, SubscriptionConfigFile)

		content := `[{
			"id": "233",
			"name": "小林家的龙女仆",
			"link": "https://mikan.example.com/bangumi/233"
		}]`

		require.NoError(t, os.WriteFile(tmpFile, []byte(content), 0644))

		runningOS = "mockOS"
		osPathProviders = map[string]Provider{
			"mockOS": mockPathProvider{
				configPathFunc: func() (string, error) {
					return tmpDir, nil
				},
			},
		}

		got, err := ReadSubscriptionConfigFile()
		require.NoError(t, err)
		require.Len(t, got, 1)
		require.Equal(t, "233", got[0].ID)
		require.Equal(t, "小林家的龙女仆", got[0].Name)
		require.Equal(t, "https://mikan.example.com/bangumi/233", got[0].Link)
	})

	t.Run("file not exist returns nil, nil", func(t *testing.T) {
		tmpDir := t.TempDir()

		runningOS = "mockOS"
		osPathProviders = map[string]Provider{
			"mockOS": mockPathProvider{
				configPathFunc: func() (string, error) {
					return tmpDir, nil
				},
			},
		}

		got, err := ReadSubscriptionConfigFile()
		require.NoError(t, err)
		require.Nil(t, got)
	})

	t.Run("returns error for unreadable or bad format file", func(t *testing.T) {
		tmpDir := t.TempDir()
		tmpFile := filepath.Join(tmpDir, SubscriptionConfigFile)

		require.NoError(t, os.WriteFile(tmpFile, []byte(`{ invalid json ]`), 0644))

		runningOS = "mockOS"
		osPathProviders = map[string]Provider{
			"mockOS": mockPathProvider{
				configPathFunc: func() (string, error) {
					return tmpDir, nil
				},
			},
		}

		got, err := ReadSubscriptionConfigFile()
		require.Error(t, err)
		require.Nil(t, got)
	})
}
