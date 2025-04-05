package path

import (
	"os"
	"path/filepath"
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
