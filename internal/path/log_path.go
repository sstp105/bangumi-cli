package path

import (
	"fmt"
	"github.com/sstp105/bangumi-cli/internal/libs"
	"os"
	"path/filepath"
	"time"
)

// LogPath returns the windows path to the app's log directory in %LocalAppData%\<APP_NAME>\logs.
func (w WindowsPath) LogPath() (string, error) {
	dir, err := os.UserCacheDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, AppDir, "logs"), nil
}

// LogPath returns the linux path to the app's log directory in $HOME/.local/share/<APP_NAME>/logs.
func (l LinuxPath) LogPath() (string, error) {
	dir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, ".local", "share", AppDir, "logs"), nil
}

// LogPath returns the macOS path to the app's log directory in $HOME/Library/Logs/<APP_NAME>.
func (m MacOSPath) LogPath() (string, error) {
	dir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "Library", "Logs", AppDir), nil
}

func ReadLogFile() (string, error) {
	fn := logFilename()
	path, err := logPath(fn)
	if err != nil {
		return "", err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func OpenLogFile() (*os.File, error) {
	fn := logFilename()
	dir, err := logPath(fn)
	if err != nil {
		return nil, err
	}

	f, err := os.OpenFile(dir, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	return f, nil
}

func logFilename() string {
	date := time.Now().Format("2006-01-02")
	fn := fmt.Sprintf("%s.log", date)
	return fn
}

func logPath(fn string) (string, error) {
	provider, supported := osPathProviders[runningOS]
	if !supported {
		return "", libs.ErrUnsupportedOS
	}

	path, err := provider.LogPath()
	if err != nil {
		return "", err
	}

	if err := os.MkdirAll(path, 0700); err != nil { // create the log folder if it does not exist
		return "", err
	}

	return filepath.Join(path, fn), nil
}
