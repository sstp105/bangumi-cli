package path

import "os"

func (w WindowsPath) LogPath() ([]string, error) {
	dir := os.Getenv("APPDATA")
	if dir == "" {
		return nil, ErrWindowsAppDataEnvNotFound
	}
	return []string{dir, AppDir}, nil
}

func (l LinuxPath) LogPath() ([]string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	return []string{home, ".config", AppDir}, nil
}

func (m MacOSPath) LogPath() ([]string, error) {
	home, err := os.UserConfigDir()
	if err != nil {
		return nil, err
	}
	return []string{home, "Library", "Logs", AppDir}, nil
}
