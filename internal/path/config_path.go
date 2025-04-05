package path

import (
	"encoding/json"
	"github.com/sstp105/bangumi-cli/internal/libs"
	"os"
	"path/filepath"
	"runtime"
)

func (w WindowsPath) ConfigPath() ([]string, error) {
	appdata := os.Getenv("APPDATA")
	if appdata == "" {
		return nil, ErrWindowsAppDataEnvNotFound
	}
	return []string{appdata, AppDir}, nil
}

func (l LinuxPath) ConfigPath() ([]string, error) {
	dir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	return []string{dir, ".config", AppDir}, nil
}

func (m MacOSPath) ConfigPath() ([]string, error) {
	dir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	return []string{dir, ".config", AppDir}, nil
}

// SaveJSONConfigFile saves the value in json format under app's config directory.
// For Windows, the file will be saved under $env:APPDATA.
// For UNIX, the file will be saved under $HOME.
// An error will be returned for unsupported OS.
// Parameters:
//   - fn: The file name (e.g., "setting.json").
//   - v: The value to be saved to the file. It will be marshaled into JSON format.
//
// Returns:
//   - error: If any error occurs during the process, an error is returned.
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

func configPath(fn string) (string, error) {
	provider, supported := osPathProviders[runtime.GOOS]
	if !supported {
		return "", libs.ErrUnsupportedOS
	}

	path, err := provider.ConfigPath()
	if err != nil {
		return "", err
	}

	dir := filepath.Join(path...)
	if err := os.MkdirAll(dir, 0700); err != nil { // create the config folder if it does not exist
		return "", err
	}

	return filepath.Join(dir, fn), nil
}
