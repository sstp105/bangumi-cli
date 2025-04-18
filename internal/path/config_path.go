package path

import (
	"encoding/json"
	"errors"
	"github.com/sstp105/bangumi-cli/internal/libs"
	"github.com/sstp105/bangumi-cli/internal/mikan"
	"os"
	"path/filepath"
	"runtime"
)

const (
	BangumiCredentialConfigFile = "bangumi_creds.json"
	SubscriptionConfigFile      = "subscriptions.json"
)

// ConfigPath returns the Windows path to the app's config directory in %AppData%\{APP_NAME}.
func (w WindowsPath) ConfigPath() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, AppDir), nil
}

// ConfigPath returns the Linux path to the app's config directory in $HOME/.config/<APP_NAME>.
func (l LinuxPath) ConfigPath() (string, error) {
	dir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, ".config", AppDir), nil
}

// ConfigPath returns the macOS path to the app's config directory in $HOME/.config/<APP_NAME>.
func (m MacOSPath) ConfigPath() (string, error) {
	dir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, ".config", AppDir), nil
}

func SaveJSONConfigFile(fn string, v any) error {
	path, err := configPath(fn)
	if err != nil {
		return err
	}

	data, err := libs.MarshalJSONIndented(v)
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0600) // owner r&w
}

func ReadJSONConfigFile(fn string, v any) error {
	path, err := configPath(fn)
	if err != nil {
		return err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, v)
	if err != nil {
		return err
	}

	return nil
}

func DeleteJSONConfigFile(fn string) error {
	path, err := configPath(fn)
	if err != nil {
		return err
	}

	return os.Remove(path)
}

func ReadSubscriptionConfigFile() ([]mikan.BangumiBase, error) {
	var subscription []mikan.BangumiBase
	err := ReadJSONConfigFile(SubscriptionConfigFile, &subscription)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil
		}
		return nil, err
	}
	return subscription, nil
}

func configPath(fn string) (string, error) {
	provider, supported := osPathProviders[runtime.GOOS]
	if !supported {
		return "", libs.ErrUnsupportedOS
	}

	path, err := provider.ConfigPath()
	if err != nil {
		return "", err
	}

	if err := os.MkdirAll(path, 0700); err != nil { // create the config folder if it does not exist
		return "", err
	}

	return filepath.Join(path, fn), nil
}
