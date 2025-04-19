package path

import (
	"encoding/json"
	"errors"
	"github.com/sstp105/bangumi-cli/internal/bangumi"
	"github.com/sstp105/bangumi-cli/internal/libs"
	"github.com/sstp105/bangumi-cli/internal/model"
	"os"
	"path/filepath"
)

const (
	BangumiCredentialConfigFile = "bangumi_creds.json"
	SubscriptionConfigFile      = "subscriptions.json"
)

func (w WindowsPath) ConfigPath() (string, error) {
	dir := os.Getenv("AppData")
	if dir == "" {
		return "", errors.New("%AppData% is not defined")
	}
	return filepath.Join(dir, AppDir), nil
}

func (m MacOSPath) ConfigPath() (string, error) {
	dir := os.Getenv("HOME")
	if dir == "" {
		return "", errors.New("$HOME is not defined")
	}
	return filepath.Join(dir, ".config", AppDir), nil
}

func SaveJSONConfigFile(fn string, v any) error {
	data, err := libs.MarshalJSONIndented(v)
	if err != nil {
		return err
	}

	path, err := configPath()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(path, 0700); err != nil { // create the config folder if it does not exist
		return err
	}

	path = filepath.Join(path, fn)

	return os.WriteFile(path, data, 0600) // owner r&w
}

func ReadJSONConfigFile(fn string, v any) error {
	path, err := configPath()
	if err != nil {
		return err
	}

	path = filepath.Join(path, fn)
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
	path, err := configPath()
	if err != nil {
		return err
	}

	path = filepath.Join(path, fn)

	return os.Remove(path)
}

func ReadSubscriptionConfigFile() ([]model.BangumiBase, error) {
	var subscription []model.BangumiBase
	err := ReadJSONConfigFile(SubscriptionConfigFile, &subscription)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil
		}
		return nil, err
	}
	return subscription, nil
}

func ReadBangumiCredentialConfigFile() (*bangumi.OAuthCredential, error) {
	var credential bangumi.OAuthCredential
	err := ReadJSONConfigFile(BangumiCredentialConfigFile, &credential)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil
		}
		return nil, err
	}
	return &credential, nil
}

func configPath() (string, error) {
	provider, supported := OSPathProviders[RunningOS]
	if !supported {
		return "", libs.ErrUnsupportedOS
	}

	path, err := provider.ConfigPath()
	if err != nil {
		return "", err
	}

	return path, nil
}
